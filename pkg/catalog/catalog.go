package catalog

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"time"

	"github.com/openkcm/plugin-sdk/api"
	"github.com/openkcm/plugin-sdk/pkg/telemetry"
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
		pluginLog := config.Logger.With(
			telemetry.PluginName, pluginConfig.Name,
			telemetry.PluginType, pluginConfig.Type,
		)
		pluginConfig.Logger = pluginLog

		if pluginConfig.Disabled {
			config.Logger.Debug("Not loading plugin; disabled")
			continue
		}

		pluginConfig.HostServices = config.HostServices

		var plugin *Plugin
		if pluginConfig.IsExternal() {
			plugin, err = loadPlugin(ctx, pluginConfig)
			if err != nil {
				config.Logger.ErrorContext(ctx, "Failed to load plugin", telemetry.PluginName, pluginConfig.Name, "error", err)
				return nil, fmt.Errorf("failed to load plugin %q: %w", pluginConfig.Name, err)
			}
		} else {
			for _, builtin := range builtIns {
				if builtin.Name == pluginConfig.Name {
					plugin, err = loadBuiltIn(ctx, builtin, BuiltInConfig{
						Logger:       pluginConfig.Logger,
						LogLevel:     pluginConfig.LogLevel,
						HostServices: config.HostServices,
					})
					if err != nil {
						config.Logger.ErrorContext(ctx, "Failed to load builtin plugin", telemetry.PluginName, pluginConfig.Name, "error", err)
						return nil, fmt.Errorf("failed to load builtin plugin %q: %w", pluginConfig.Name, err)
					}
				}
			}
		}
		if plugin == nil {
			config.Logger.ErrorContext(ctx, "Failed to load external/builtin plugin", telemetry.PluginName, pluginConfig.Name)
			return nil, fmt.Errorf("failed to load external/builtin plugin %s", pluginConfig.Name)
		}

		closers = append(closers, pluginCloser{plugin: plugin, log: pluginLog})

		cfgurer := makeConfigurer(plugin, pluginConfig)
		err = cfgurer.Configure(ctx)
		if err != nil {
			config.Logger.ErrorContext(ctx, "Failed to configure plugin", telemetry.PluginName, pluginConfig.Name, "error", err)
			return nil, fmt.Errorf("failed to configure plugin %q of type %q; %v", pluginConfig.Name, pluginConfig.Type, err)
		}
		configurers = append(configurers, cfgurer)

		pluginLog.Info("Plugin loaded")
	}

	return &Catalog{
		closers:     closers,
		configurers: configurers,
	}, nil
}
