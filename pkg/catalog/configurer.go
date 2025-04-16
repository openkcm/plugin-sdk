package catalog

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	configv1 "github.com/openkcm/plugin-sdk/proto/service/common/config/v1"
)

type Configurers []*configurer

type configurer struct {
	plugin       *Plugin
	pluginConfig PluginConfig
}

func makeConfigurer(plugin *Plugin, pluginConfig PluginConfig) *configurer {
	return &configurer{
		plugin:       plugin,
		pluginConfig: pluginConfig,
	}
}

func (c *configurer) Configure(ctx context.Context) error {
	client := configv1.NewConfigClient(c.plugin.ClientConnection())
	_, err := client.Configure(ctx, &configv1.ConfigureRequest{
		YamlConfiguration: c.pluginConfig.YamlConfiguration,
	})
	switch status.Code(err) {
	case codes.Unimplemented:
		return nil
	case codes.OK:
		return nil
	}
	return err
}
