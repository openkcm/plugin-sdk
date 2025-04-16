package bootstrap

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	initv1 "github.com/openkcm/plugin-sdk/internal/proto/service/init/v1"
)

// Init initializes the plugin and advertises the given host service names to
// the plugin for brokering. The list of service names implemented by the
// plugin are returned. This function is only intended to be used internally.
func Init(ctx context.Context, conn grpc.ClientConnInterface, hostServiceNames []string) (pluginServiceNames []string, err error) {
	client := initv1.NewBootstrapClient(conn)
	resp, err := client.Init(ctx, &initv1.InitRequest{
		HostServiceNames: hostServiceNames,
	})
	switch status.Code(err) {
	case codes.Unimplemented:
		return []string{}, nil
	case codes.OK:
		return resp.PluginServiceNames, nil
	}
	return nil, err
}

// Deinit deinitializes the plugin. It should only be called right before the
// host unloads the plugin and will not be invoking any other plugin or service
// RPCs.
func Deinit(ctx context.Context, conn grpc.ClientConnInterface) error {
	client := initv1.NewBootstrapClient(conn)
	_, err := client.Deinit(ctx, &initv1.DeinitRequest{})
	switch status.Code(err) {
	case codes.OK, codes.Unimplemented:
		return nil
	default:
		return err
	}
}
