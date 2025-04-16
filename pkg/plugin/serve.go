package plugin

import (
	"github.com/openkcm/plugin-sdk/api"
	"github.com/openkcm/plugin-sdk/internal/bootstrap"
)

func Serve(pluginServer api.PluginServer, serviceServers ...api.ServiceServer) {
	logger := bootstrap.NewLogger()
	bootstrap.Serve(logger, logger, pluginServer, serviceServers, nil)
}
