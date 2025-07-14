package catalog

import (
	"context"
	"log/slog"

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

type Configurer interface {
	Configure(ctx context.Context, configuration string) error
}

type ConfigurerFunc func(ctx context.Context, configuration string) error

func (fn ConfigurerFunc) Configure(ctx context.Context, configuration string) error {
	return fn(ctx, configuration)
}

type configurerRepo struct {
	configurer Configurer
}

func (repo *configurerRepo) Binder() any {
	return func(configurer Configurer) {
		repo.configurer = configurer
	}
}

func (repo *configurerRepo) Versions() []Version {
	return []Version{
		configurerV1Version{},
	}
}

func (repo *configurerRepo) Clear() {
	// This function is only for conforming to the Repo interface and isn't
	// expected to be called, but just in case, we'll do the right thing
	// and clear out the configurer that has been bound.
	repo.configurer = nil
}

type configurerV1Version struct{}

func (configurerV1Version) New() Facade      { return new(configurerV1) }
func (configurerV1Version) Deprecated() bool { return false }

type configurerV1 struct {
	configv1.ConfigServiceClient
}

var _ Configurer = (*configurerV1)(nil)

func (v1 *configurerV1) InitInfo(PluginInfo) {
}

func (v1 *configurerV1) InitLog(*slog.Logger) {
}

func (v1 *configurerV1) Configure(ctx context.Context, yamlConfiguration string) error {
	_, err := v1.ConfigServiceClient.Configure(ctx, &configv1.ConfigureRequest{
		YamlConfiguration: yamlConfiguration,
	})
	return err
}

type configurerUnsupported struct{}

func (c configurerUnsupported) Configure(context.Context, string) error {
	return status.Error(codes.FailedPrecondition, "plugin does not support a configuration interface")
}
