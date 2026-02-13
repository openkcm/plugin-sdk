package plugin

import (
	"fmt"
	"os"

	"github.com/openkcm/plugin-sdk/api"
	pluginoption "github.com/openkcm/plugin-sdk/api/plugin-option"
	"github.com/openkcm/plugin-sdk/internal/bootstrap"
)

func Serve(pluginServer api.PluginServer, serviceServers ...api.ServiceServer) {
	err := bootstrap.Serve(
		pluginoption.WithPluginServer(pluginServer),
		pluginoption.WithServiceServer(serviceServers...),
	)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "\nFailed to serve plugin: %v", err)
	}
}

func ServeOptions(options ...pluginoption.ServerOption) error {
	opts := make([]pluginoption.ServerOption, 0, len(options))
	opts = append(opts, options...)

	return bootstrap.Serve(opts...)
}
