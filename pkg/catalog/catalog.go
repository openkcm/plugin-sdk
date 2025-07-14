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

func (c *Catalog) LookupByType(pluginType string) []*Plugin {
	var plugins []*Plugin
	for _, cfgr := range c.configurers {
		if cfgr.plugin.Info().Type() == pluginType {
			plugins = append(plugins, cfgr.plugin)
		}
	}
	return plugins
}

func (c *Catalog) LookupByTypeAndName(pluginType, pluginName string) *Plugin {
	for _, cfgr := range c.configurers {
		if cfgr.plugin.Info().Type() == pluginType && cfgr.plugin.Info().Name() == pluginName {
			return cfgr.plugin
		}
	}
	return nil
}

func Load(ctx context.Context, config Config, builtIns ...BuiltIn) (catalog *Catalog, err error) {
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

		closers = append(closers, pluginCloser{plugin: plugin, log: pluginConfig.Logger})

		cfgurer := makeConfigurer(plugin, pluginConfig)
		err = cfgurer.Configure(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to configure plugin %s of type %s; %v", pluginConfig.Name, pluginConfig.Type, err)
		}
		configurers = append(configurers, cfgurer)

		pluginConfig.Logger.Info("Plugin loaded")
	}

	return &Catalog{
		closers:     closers,
		configurers: configurers,
	}, nil
}

func loadPluginAs(ctx context.Context, logger *slog.Logger, pluginConfig PluginConfig, builtIns ...BuiltIn) (*Plugin, error) {
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

func loadPluginAsExternal(ctx context.Context, logger *slog.Logger, pluginConfig PluginConfig) (*Plugin, error) {
	pluginLog := logger.With(
		Name, pluginConfig.Name,
		Type, pluginConfig.Type,
	)
	pluginConfig.Logger = pluginLog

	return loadPlugin(ctx, pluginConfig)
}

func loadPluginAsBuiltIn(ctx context.Context, logger *slog.Logger, pluginConfig PluginConfig, builtIns ...BuiltIn) (*Plugin, error) {
	for _, builtin := range builtIns {
		if builtin.Name == pluginConfig.Name {
			pluginLog := logger.With(
				Name, pluginConfig.Name,
				Type, builtin.Plugin.Type(),
			)
			pluginConfig.Logger = pluginLog

			return loadBuiltIn(ctx, builtin, pluginConfig)
		}
	}
	return nil, fmt.Errorf("builtin plugin %q not found", pluginConfig.Name)
}
