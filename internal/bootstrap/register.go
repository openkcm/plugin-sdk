package bootstrap

import (
	"context"
	"io"

	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc"

	"github.com/openkcm/plugin-sdk/api"
	initv1 "github.com/openkcm/plugin-sdk/internal/proto/service/init/v1"
)

// HostDialer is a generic interface for dialing the host. This
// interface is only intended to be used internally.
type HostDialer interface {
	DialHost(ctx context.Context) (grpc.ClientConnInterface, error)
}

// register given servers with the gRPC server. The given dialer and logger will
// be used when the plugins are initialized.
func Register(s *grpc.Server, servers []api.ServiceServer, logger hclog.Logger, dialer HostDialer) {
	var names []string
	var impls []any
	for _, server := range servers {
		names = append(names, server.GRPCServiceName())
		impls = append(impls, server.RegisterServer(s))
	}

	initv1.RegisterBootstrapServer(s, &initService{
		logger: logger,
		names:  names,
		impls:  impls,
		dialer: dialer,
	})
}

type initService struct {
	initv1.UnimplementedBootstrapServer

	logger hclog.Logger
	names  []string
	impls  []any
	dialer HostDialer
}

func (s *initService) Init(ctx context.Context, req *initv1.InitRequest) (*initv1.InitResponse, error) {
	initted := map[any]struct{}{}
	for _, impl := range s.impls {
		// Wire up the logger and host service broker. Since the same
		// implementation might back more than one server, only initialize
		// once.
		if _, ok := initted[impl]; ok {
			continue
		}
		initted[impl] = struct{}{}

		if impl, ok := impl.(api.NeedsLogger); ok {
			impl.SetLogger(s.logger)
		}

		if impl, ok := impl.(api.NeedsHostServices); ok {
			conn, err := s.dialer.DialHost(ctx)
			if err != nil {
				return nil, err
			}
			broker := serviceBroker{conn: conn, hostServiceNames: req.HostServiceNames}
			if err := impl.BrokerHostServices(broker); err != nil {
				s.logger.Error("Plugin failed brokering host services", "error", err)
				return nil, err
			}
		}
	}

	return &initv1.InitResponse{
		PluginServiceNames: s.names,
	}, nil
}

func (s *initService) Deinit(ctx context.Context, req *initv1.DeinitRequest) (*initv1.DeinitResponse, error) {
	deinitted := map[any]struct{}{}
	for _, impl := range s.impls {
		// Deinitialize the implementation. Since the same
		// implementation might back more than one server, only deinitialize
		// once.
		if _, ok := deinitted[impl]; ok {
			continue
		}
		deinitted[impl] = struct{}{}

		if impl, ok := impl.(io.Closer); ok {
			if err := impl.Close(); err != nil {
				s.logger.Error("Plugin implementation failed to deinitialize", "error", err)
			}
		}
	}
	return &initv1.DeinitResponse{}, nil
}

type serviceBroker struct {
	conn             grpc.ClientConnInterface
	hostServiceNames []string
}

func (b serviceBroker) BrokerClient(client api.ServiceClient) bool {
	wants := client.GRPCServiceName()
	for _, has := range b.hostServiceNames {
		if wants == has {
			client.InitClient(b.conn)
			return true
		}
	}
	return false
}
