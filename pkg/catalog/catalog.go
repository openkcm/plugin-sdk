package catalog

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"time"

	"github.com/openkcm/plugin-sdk/api"
)

const (
	deinitTimeout = 1 * time.Minute
	initTimeout   = 10 * time.Minute
)

type Config struct {
	// Logger is the logger. It is used for general purpose logging and also
	// provided to the plugins.
	Logger *slog.Logger

	// PluginConfigs is the list of plugin configurations.
	PluginConfigs []PluginConfig

	// HostServices are the servers for host services provided by SPIRE to
	// plugins.
	HostServices []api.ServiceServer
}

type Catalog struct {
	closers     io.Closer
	configurers Configurers
}

func (c *Catalog) Close() error {
	return c.closers.Close()
}

func (c *Catalog) LookupByType(pluginType string) []Plugin {
	var plugins []Plugin
	for _, cfgr := range c.configurers {
		if cfgr.plugin.Info().Type() == pluginType {
			plugins = append(plugins, cfgr.plugin)
		}
	}
	return plugins
}

func (c *Catalog) LookupByTypeAndName(pluginType, pluginName string) Plugin {
	for _, cfgr := range c.configurers {
		if cfgr.plugin.Info().Type() == pluginType && cfgr.plugin.Info().Name() == pluginName {
			return cfgr.plugin
		}
	}
	return nil
}

func (c *Catalog) ListPluginInfo() []PluginInfo {
	var plugins []PluginInfo
	for _, cfgr := range c.configurers {
		plugins = append(plugins, cfgr.plugin.Info())
	}
	return plugins
}

func Load(ctx context.Context, config Config, builtIns ...BuiltInPlugin) (catalog *Catalog, err error) {
	closers := make(closerGroup, 0)
	defer func() {
		// If loading fails, clear out the catalog and close down all plugins
		// that have been loaded thus far.
		if err != nil {
			if err2 := closers.Close(); err2 != nil {
				config.Logger.ErrorContext(ctx, "Failed to close plugins", "error", err2)
			}
		}
	}()

	// in case if configuration logger is not set get the default one
	if config.Logger == nil {
		config.Logger = slog.Default()
	}

	configurers := make(Configurers, 0)
	for _, pluginConfig := range config.PluginConfigs {
		if pluginConfig.Disabled {
			config.Logger.Debug("Not loading plugin; disabled")
			continue
		}

		pluginConfig.HostServices = config.HostServices

		plugin, err := loadPluginAs(ctx, config.Logger, pluginConfig, builtIns...)
		if err != nil {
			return nil, err
		}

		closers = append(closers, pluginCloser{plugin: plugin, log: plugin.Logger()})

		cfgurer := makeConfigurer(plugin, pluginConfig)
		err = cfgurer.Configure(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to configure plugin %s of type %s; %v", pluginConfig.Name, pluginConfig.Type, err)
		}
		configurers = append(configurers, cfgurer)

		plugin.Logger().Info("Loaded plugin")
	}

	return &Catalog{
		closers:     closers,
		configurers: configurers,
	}, nil
}

func loadPluginAs(ctx context.Context, logger *slog.Logger, pluginConfig PluginConfig, builtIns ...BuiltInPlugin) (Plugin, error) {
	if pluginConfig.IsExternal() {
		plugin, err := loadPluginAsExternal(ctx, logger, pluginConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to load external plugin %s: %w", pluginConfig.Name, err)
		}
		return plugin, nil
	}

	plugin, err := loadPluginAsBuiltIn(ctx, logger, pluginConfig, builtIns...)
	if err != nil {
		return nil, fmt.Errorf("failed to load builtin plugin %s: %w", pluginConfig.Name, err)
	}
	return plugin, nil
}

func loadPluginAsExternal(ctx context.Context, logger *slog.Logger, pluginConfig PluginConfig) (Plugin, error) {
	if pluginConfig.Name == "" {
		return nil, fmt.Errorf("failed to load external plugin, missing name")
	}

	if pluginConfig.Type == "" {
		return nil, fmt.Errorf("failed to load external plugin %s, missing type", pluginConfig.Name)
	}

	pluginConfig.Logger = logger.With(
		Name, pluginConfig.Name,
		Type, pluginConfig.Type,
	)

	return loadPlugin(ctx, pluginConfig)
}

func loadPluginAsBuiltIn(ctx context.Context, logger *slog.Logger, pluginConfig PluginConfig, builtIns ...BuiltInPlugin) (Plugin, error) {
	if pluginConfig.Name == "" {
		return nil, fmt.Errorf("failed to load builtin plugin, missing name")
	}

	for _, builtin := range builtIns {
		if builtin.Name() == pluginConfig.Name {
			pluginConfig.Logger = logger.With(
				Name, pluginConfig.Name,
				Type, builtin.Type(),
			)

			return loadBuiltInPlugin(ctx, builtin, pluginConfig)
		}
	}
	return loadPluginAsExternal(ctx, logger, pluginConfig)
}
