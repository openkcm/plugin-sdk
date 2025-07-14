package catalog

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"sync"
	"time"

	"github.com/magodo/slog2hclog"
	"github.com/openkcm/plugin-sdk/api"
	"github.com/openkcm/plugin-sdk/internal/bootstrap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type BuiltIn struct {
	Name     string
	Tags     []string
	Plugin   api.PluginServer
	Services []api.ServiceServer
}

func MakeBuiltIn(name string, pluginServer api.PluginServer, serviceServers ...api.ServiceServer) BuiltIn {
	return BuiltIn{
		Name:     name,
		Plugin:   pluginServer,
		Services: serviceServers,
	}
}

type BuiltInConfig struct {
	// Logger is the logger to be wired to the external plugin.
	Logger *slog.Logger

	LogLevel string

	// HostServices are the host service servers provided to the plugin.
	HostServices []api.ServiceServer
}

//func LoadBuiltIn(ctx context.Context, builtIn BuiltIn, config BuiltInConfig) (_ *Plugin, err error) {
//	return loadBuiltIn(ctx, builtIn, config)
//}

func loadBuiltIn(ctx context.Context, builtIn BuiltIn, config BuiltInConfig) (_ *Plugin, err error) {
	dialer := &builtinDialer{
		pluginName:   builtIn.Name,
		log:          config.Logger,
		hostServices: config.HostServices,
	}

	var closers closerGroup
	defer func() {
		if err != nil {
			closers.Close()
		}
	}()
	closers = append(closers, dialer)

	builtinServer, serverCloser := newBuiltInServer(config.Logger)
	closers = append(closers, serverCloser)

	pluginServers := append([]api.ServiceServer{builtIn.Plugin}, builtIn.Services...)

	logLevelPlugin := new(slog.LevelVar)
	setLogLevel(logLevelPlugin, config.LogLevel)

	log := slog2hclog.New(config.Logger, logLevelPlugin)
	bootstrap.Register(builtinServer, pluginServers, log, dialer)

	builtinConn, err := startPipeServer(builtinServer, config.Logger)
	if err != nil {
		return nil, err
	}
	closers = append(closers, builtinConn)

	info := pluginInfo{
		name: builtIn.Name,
		typ:  builtIn.Plugin.Type(),
		tags: builtIn.Tags,
	}

	return newPlugin(ctx, builtinConn, info, config.Logger, closers, config.HostServices)
}

func newBuiltInServer(log *slog.Logger) (*grpc.Server, io.Closer) {
	drain := &drainHandlers{}
	return grpc.NewServer(
		grpc.ChainStreamInterceptor(drain.StreamServerInterceptor, streamPanicInterceptor(log)),
		grpc.ChainUnaryInterceptor(drain.UnaryServerInterceptor, unaryPanicInterceptor(log)),
	), closerFunc(drain.Wait)
}

type builtinDialer struct {
	pluginName   string
	log          *slog.Logger
	hostServices []api.ServiceServer
	conn         *pipeConn
}

func (d *builtinDialer) DialHost(context.Context) (grpc.ClientConnInterface, error) {
	if d.conn != nil {
		return d.conn, nil
	}
	server := newHostServer(d.log, d.pluginName)
	conn, err := startPipeServer(server, d.log)
	if err != nil {
		return nil, err
	}
	d.conn = conn
	return d.conn, nil
}

func (d *builtinDialer) Close() error {
	if d.conn != nil {
		return d.conn.Close()
	}
	return nil
}

type pipeConn struct {
	grpc.ClientConnInterface
	io.Closer
}

func startPipeServer(server *grpc.Server, log *slog.Logger) (_ *pipeConn, err error) {
	var closers closerGroup

	pipeNet := newPipeNet()
	closers = append(closers, pipeNet)

	var wg sync.WaitGroup
	closers = append(closers, closerFunc(wg.Wait), closerFunc(func() {
		if !gracefulStopWithTimeout(server, time.Minute) {
			log.Warn("Forced timed-out plugin server to stop")
		}
	}))

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := server.Serve(pipeNet); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			log.Error("Pipe server unexpectedly failed to serve", "error", err)
		}
	}()

	// Dial the server
	conn, err := grpc.NewClient(
		"passthrough:IGNORED",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithContextDialer(pipeNet.DialContext),
	)
	if err != nil {
		return nil, err
	}
	closers = append(closers, conn)

	return &pipeConn{
		ClientConnInterface: conn,
		Closer:              closers,
	}, nil
}

type drainHandlers struct {
	wg sync.WaitGroup
}

func (d *drainHandlers) Wait() {
	done := make(chan struct{})

	go func() {
		d.wg.Wait()
		close(done)
	}()

	t := time.NewTimer(time.Minute)
	defer t.Stop()

	select {
	case <-done:
	case <-t.C:
	}
}

func (d *drainHandlers) UnaryServerInterceptor(ctx context.Context, req any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	d.wg.Add(1)
	defer d.wg.Done()
	return handler(ctx, req)
}

func (d *drainHandlers) StreamServerInterceptor(srv any, ss grpc.ServerStream, _ *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	d.wg.Add(1)
	defer d.wg.Done()
	return handler(srv, ss)
}
