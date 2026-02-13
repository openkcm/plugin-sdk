package pluginoption

import (
	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc"

	goplugin "github.com/hashicorp/go-plugin"

	"github.com/openkcm/plugin-sdk/api"
)

type ServerConfiguration struct {
	PluginServer   api.PluginServer
	ServerOptions  []grpc.ServerOption
	ServiceServers []api.ServiceServer

	Logger hclog.Logger

	TestConfig *goplugin.ServeTestConfig

	ValidateInput  bool
	ValidateOutput bool
}

type ServerOption func(*ServerConfiguration)

func EnableInputValidation() ServerOption {
	return func(gs *ServerConfiguration) {
		gs.ValidateInput = true
	}
}

func EnableOutputValidation() ServerOption {
	return func(gs *ServerConfiguration) {
		gs.ValidateOutput = true
	}
}

func SetServerOption(opts ...grpc.ServerOption) ServerOption {
	return func(gs *ServerConfiguration) {
		gs.ServerOptions = opts
	}
}

func WithServiceServer(serviceServers ...api.ServiceServer) ServerOption {
	return func(gs *ServerConfiguration) {
		gs.ServiceServers = serviceServers
	}
}

func WithLogger(logger hclog.Logger) ServerOption {
	return func(gs *ServerConfiguration) {
		gs.Logger = logger
	}
}

func WithPluginServer(server api.PluginServer) ServerOption {
	return func(gs *ServerConfiguration) {
		gs.PluginServer = server
	}
}

func WithTestConfig(config *goplugin.ServeTestConfig) ServerOption {
	return func(gs *ServerConfiguration) {
		gs.TestConfig = config
	}
}
