package catalog

import (
	"context"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	configv1 "github.com/openkcm/plugin-sdk/proto/service/common/config/v1"
)

const (
	defaultEmptyBuildInfo = "{}"
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
	resp, err := client.Configure(ctx, &configv1.ConfigureRequest{
		YamlConfiguration: c.pluginConfig.YamlConfiguration,
	})
	switch status.Code(err) {
	case codes.Unimplemented:
		return nil
	case codes.OK:
		sbi, ok := c.plugin.Info().(Build)
		if ok {
			sbi.SetValue(extractBuildInfo(resp))
		}
		return nil
	}
	return err
}

func extractBuildInfo(resp *configv1.ConfigureResponse) string {
	defer func() {
		_ = recover()
	}()

	if resp == nil {
		return defaultEmptyBuildInfo
	}

	value := strings.TrimSpace(resp.GetBuildInfo())
	if value == "" {
		return defaultEmptyBuildInfo
	}
	return value
}
