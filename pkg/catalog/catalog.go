package catalog

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"strings"

	"github.com/openkcm/plugin-sdk/api"
)

type Catalog struct {
	closers     io.Closer
	configurers Reconfigurers
}

func (c *Catalog) Close() error {
	return c.closers.Close()
}

func (c *Catalog) Reconfigure(ctx context.Context) {
	c.configurers.Reconfigure(ctx)
}

func (c *Catalog) LookupByType(pluginType string) []Plugin {
	var plugins []Plugin

	defer func() {
		if err := recover(); err != nil {
			slog.Error("Panic: Plugin lookup failed with plugin type",
				"pluginType", pluginType, "error", err)
		}
	}()

	closerGr := c.closers.(closerGroup)
	for _, cfger := range closerGr {
		pluginClose := cfger.(pluginCloser)
		plugin := pluginClose.plugin.(Plugin)
		if plugin.Info().Type() == pluginType {
			plugins = append(plugins, plugin)
		}
	}
	return plugins
}

func (c *Catalog) LookupByTypeAndName(pluginType, pluginName string) Plugin {
	defer func() {
		if err := recover(); err != nil {
			slog.Error("Panic: Plugin lookup failed with plugin type",
				"pluginType", pluginType, "pluginName", pluginName, "error", err)
		}
	}()

	closerGr := c.closers.(closerGroup)
	for _, cfger := range closerGr {
		pluginClose := cfger.(pluginCloser)
		plugin := pluginClose.plugin.(Plugin)
		if plugin.Info().Type() == pluginType && plugin.Info().Name() == pluginName {
			return plugin
		}
	}
	return nil
}

func (c *Catalog) ListPluginInfo() []api.Info {
	var plugins []api.Info

	defer func() {
		if err := recover(); err != nil {
			slog.Error("Panic: List plugin info", "error", err)
		}
	}()

	closerGr := c.closers.(closerGroup)
	for _, cfger := range closerGr {
		pluginClose := cfger.(pluginCloser)
		plugin := pluginClose.plugin.(Plugin)
		plugins = append(plugins, plugin.Info())
	}
	return plugins
}

func New(ctx context.Context, config Config, repo api.Repository, builtIns ...BuiltInPlugin) (_ *Catalog, err error) {
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

	pluginRepos, err := makeBindablePluginRepos(repo.Plugins())
	if err != nil {
		return nil, err
	}
	serviceRepos, err := makeBindableServiceRepos(repo.Services())
	if err != nil {
		return nil, err
	}

	// in case if configuration logger is not set get the default one
	if config.Logger == nil {
		config.Logger = slog.Default()
	}

	pluginCounts := make(map[string]int)
	var reconfigurers Reconfigurers

	for _, pluginConfig := range config.PluginConfigs {
		if pluginConfig.Disabled {
			config.Logger.Debug("Not loading plugin; disabled")
			continue
		}

		pluginConfig.HostServices = config.HostServices

		pluginRepo, ok := pluginRepos[pluginConfig.Type]
		if !ok {
			slog.Error("Unsupported plugin type")
			return nil, fmt.Errorf("unsupported plugin type %q", pluginConfig.Type)
		}

		plugin, err := loadPluginAs(ctx, config.Logger, pluginConfig, builtIns...)
		if err != nil {
			return nil, err
		}

		closers = append(closers, pluginCloser{plugin: plugin, log: plugin.Logger()})

		cfrer, err := plugin.bindRepos(pluginRepo, serviceRepos)
		if err != nil {
			plugin.Logger().Error("Failed to bind plugin", "error", err)
			return nil, fmt.Errorf("failed to bind plugin %q: %w", pluginConfig.Name, err)
		}

		externalYamlConfiguration := strings.TrimSpace(pluginConfig.YamlConfiguration)
		if pluginConfig.DataSource == nil && len(externalYamlConfiguration) > 0 {
			pluginConfig.DataSource = FixedData(pluginConfig.YamlConfiguration)
		}

		reconfigurer, err := configurePlugin(ctx, plugin.Logger(), cfrer, pluginConfig.DataSource)
		if err != nil {
			plugin.Logger().Error("Failed to configure plugin", "error", err)
			return nil, fmt.Errorf("failed to configure plugin %q: %w", pluginConfig.Name, err)
		}
		if reconfigurer != nil {
			reconfigurers = append(reconfigurers, reconfigurer)
		}

		plugin.Logger().Info("Loaded plugin")
		pluginCounts[pluginConfig.Type]++
	}

	impl := &Catalog{
		closers:     closers,
		configurers: reconfigurers,
	}

	// Make sure all plugin constraints are satisfied
	for pluginType, pluginRepo := range pluginRepos {
		if _, ok := pluginCounts[pluginType]; !ok {
			continue
		}

		if err := pluginRepo.Constraints().Check(pluginCounts[pluginType]); err != nil {
			return nil, fmt.Errorf("plugin type %q constraint not satisfied: %w", pluginType, err)
		}
	}

	return impl, nil
}

func loadPluginAs(ctx context.Context, logger *slog.Logger, pluginConfig PluginConfig, builtIns ...BuiltInPlugin) (*pluginImpl, error) {
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

func loadPluginAsExternal(ctx context.Context, logger *slog.Logger, pluginConfig PluginConfig) (*pluginImpl, error) {
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

func loadPluginAsBuiltIn(ctx context.Context, logger *slog.Logger, pluginConfig PluginConfig, builtIns ...BuiltInPlugin) (*pluginImpl, error) {
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
	return nil, fmt.Errorf("builtin plugin %q not found", pluginConfig.Name)
}
