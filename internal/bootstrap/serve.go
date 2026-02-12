package bootstrap

import (
	"context"
	"errors"
	"log/slog"

	"buf.build/go/protovalidate"
	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"

	goplugin "github.com/hashicorp/go-plugin"

	"github.com/openkcm/plugin-sdk/api"
	pluginerrors "github.com/openkcm/plugin-sdk/api/plugin-errors"
	pluginoption "github.com/openkcm/plugin-sdk/api/plugin-option"
	"github.com/openkcm/plugin-sdk/internal/consts"
)

// Serve serves the plugin with the given loggers and plugin/service servers and an optional test configuration.
func Serve(
	opts ...pluginoption.ServerOption,
) error {
	cfg := &pluginoption.ServerConfiguration{}
	for _, opt := range opts {
		opt(cfg)
	}

	if cfg.PluginServer == nil {
		return pluginerrors.ErrServerRequired
	}

	if cfg.Logger == nil {
		cfg.Logger = NewLogger()
	}

	if cfg.ValidateInput || cfg.ValidateOutput {
		if len(cfg.ServerOptions) == 0 {
			cfg.ServerOptions = make([]grpc.ServerOption, 0)
		}

		validator, err := protovalidate.New()
		if err != nil {
			slog.Error("failed to initialize validator", "error", err)
		} else {
			cfg.ServerOptions = append(cfg.ServerOptions, grpc.UnaryInterceptor(
				ValidationUnaryInterceptor(validator, cfg.ValidateInput, cfg.ValidateOutput),
			))
		}
	}

	goplugin.Serve(&goplugin.ServeConfig{
		HandshakeConfig: ServerHandshakeConfig(cfg.PluginServer),
		Plugins: map[string]goplugin.Plugin{
			cfg.PluginServer.Type(): newHCPlugin(cfg.Logger, cfg.PluginServer, cfg.ServiceServers),
		},
		Logger:     cfg.Logger,
		GRPCServer: customGRPCServer(cfg.ServerOptions),
		Test:       cfg.TestConfig,
	})

	return nil
}

type hcServer struct {
	goplugin.NetRPCUnsupportedPlugin
	logger  hclog.Logger
	servers []api.ServiceServer
}

func newHCPlugin(logger hclog.Logger, pluginServer api.PluginServer, serviceServers []api.ServiceServer) *hcServer {
	return &hcServer{
		logger:  logger,
		servers: append([]api.ServiceServer{pluginServer}, serviceServers...),
	}
}

func (p *hcServer) GRPCServer(broker *goplugin.GRPCBroker, server *grpc.Server) (err error) {
	Register(server, p.servers, p.logger, &hcDialer{broker: broker})
	return nil
}

func (p *hcServer) GRPCClient(context.Context, *goplugin.GRPCBroker, *grpc.ClientConn) (any, error) {
	return nil, errors.New("unimplemented")
}

type hcDialer struct {
	broker *goplugin.GRPCBroker
	conn   grpc.ClientConnInterface
}

func (d *hcDialer) DialHost(ctx context.Context) (grpc.ClientConnInterface, error) {
	if d.conn != nil {
		return d.conn, nil
	}

	conn, err := d.broker.Dial(consts.HostServiceProviderID)
	if err != nil {
		return nil, err
	}
	d.conn = conn
	return conn, nil
}

func customGRPCServer(
	base []grpc.ServerOption,
) func([]grpc.ServerOption) *grpc.Server {
	return func(opts []grpc.ServerOption) *grpc.Server {
		all := append(append([]grpc.ServerOption{}, opts...), base...)
		return grpc.NewServer(all...)
	}
}

func ValidationUnaryInterceptor(v protovalidate.Validator, request, response bool) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {

		if msg, ok := req.(proto.Message); request && ok {
			if err := v.Validate(msg); err != nil {
				return nil, status.Errorf(
					codes.InvalidArgument,
					"request validation failed: %v",
					err,
				)
			}
		}

		resp, err := handler(ctx, req)
		if err != nil {
			return nil, err
		}

		if msg, ok := resp.(proto.Message); response && ok {
			if err := v.Validate(msg); err != nil {
				return nil, status.Errorf(
					codes.Internal,
					"response validation failed: %v",
					err,
				)
			}
		}

		return resp, nil
	}
}
