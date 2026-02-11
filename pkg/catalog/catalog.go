package catalog

import (
	"context"
	"fmt"
	"io"
	"log/slog"

	"github.com/openkcm/plugin-sdk/api"
)

const (
	certificateIssuerType             = "CertificateIssuer"
	notificationType                  = "Notification"
	systemInformationType             = "SystemInformation"
	keystoreManagementType            = "KeystoreManagement"
	keystoreInstanceKeyOperationsType = "KeystoreInstanceKeyOperations"
)

type PluginRepository struct {
	certificateIssuerRepository
	notificationRepository
	systemInformationRepository
	keystoreManagementRepository
	keystoreInstanceKeyOperationsRepository

	log     *slog.Logger
	catalog *Catalog
}

func (repo *PluginRepository) Plugins() map[string]PluginRepo {
	return map[string]PluginRepo{
		certificateIssuerType:             &repo.certificateIssuerRepository,
		notificationType:                  &repo.notificationRepository,
		systemInformationType:             &repo.systemInformationRepository,
		keystoreManagementType:            &repo.systemInformationRepository,
		keystoreInstanceKeyOperationsType: &repo.keystoreInstanceKeyOperationsRepository,
	}
}

func (repo *PluginRepository) Services() []ServiceRepo {
	return nil
}

func (repo *PluginRepository) Reconfigure(ctx context.Context) {
	repo.catalog.Reconfigure(ctx)
}

func (repo *PluginRepository) Close() error {
	if repo.catalog == nil {
		return nil
	}

	err := repo.catalog.Close()
	if err != nil {
		return err
	}

	return nil
}

func (repo *PluginRepository) ListPluginInfo() []api.Info {
	var plugins []api.Info

	closerGr := repo.catalog.closers.(closerGroup)
	for _, cfger := range closerGr {
		pluginClose := cfger.(pluginCloser)
		plugin := pluginClose.plugin.(Plugin)
		plugins = append(plugins, plugin.Info())
	}
	return plugins
}

func CreateRegistry(ctx context.Context, config Config, builtIns ...BuiltInPlugin) (_ *PluginRepository, err error) {
	repo := &PluginRepository{
		log: config.Logger,
	}
	defer func() {
		if err != nil {
			_ = repo.Close()
		}
	}()

	if len(config.HostServices) == 0 {
		config.HostServices = make([]api.ServiceServer, 0)
	}

	repo.catalog, err = load(ctx, config, repo, builtIns...)
	if err != nil {
		return nil, err
	}

	return repo, nil
}

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

// Deprecated: [will be removed once switched to CreateRegistry function]
func (c *Catalog) LookupByType(pluginType string) []Plugin {
	var plugins []Plugin
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

// Deprecated: [will be removed once switched to CreateRegistry function]
func (c *Catalog) LookupByTypeAndName(pluginType, pluginName string) Plugin {
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

// Deprecated: [will be removed once switched to CreateRegistry function]
func (c *Catalog) ListPluginInfo() []api.Info {
	var plugins []api.Info

	closerGr := c.closers.(closerGroup)
	for _, cfger := range closerGr {
		pluginClose := cfger.(pluginCloser)
		plugin := pluginClose.plugin.(Plugin)
		plugins = append(plugins, plugin.Info())
	}
	return plugins
}

// Deprecated: [CreateRegistry function to be used instead]
func Load(ctx context.Context, config Config, builtIns ...BuiltInPlugin) (_ *Catalog, err error) {
	repo := &PluginRepository{
		log: config.Logger,
	}
	defer func() {
		if err != nil {
			_ = repo.Close()
		}
	}()

	if len(config.HostServices) == 0 {
		config.HostServices = make([]api.ServiceServer, 0)
	}

	repo.catalog, err = load(ctx, config, repo, builtIns...)
	if err != nil {
		return nil, err
	}

	return repo.catalog, nil
}

func load(ctx context.Context, config Config, repo Repository, builtIns ...BuiltInPlugin) (_ *Catalog, err error) {
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

		builtInPlugins := make([]BuiltInPlugin, 0, len(builtIns))
		builtInPlugins = append(builtInPlugins, builtIns...)
		builtInPlugins = append(builtInPlugins, pluginRepo.BuiltIns()...)

		plugin, err := loadPluginAs(ctx, config.Logger, pluginConfig, builtInPlugins...)
		if err != nil {
			return nil, err
		}

		closers = append(closers, pluginCloser{plugin: plugin, log: plugin.Logger()})

		cfrer, err := plugin.bindRepos(pluginRepo, serviceRepos)
		if err != nil {
			plugin.Logger().Error("Failed to bind plugin", "error", err)
			return nil, fmt.Errorf("failed to bind plugin %q: %w", pluginConfig.Name, err)
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

	requiringCheck := make(map[string]int)
	if len(config.RequiredPlugins) == 0 {
		requiringCheck = pluginCounts
	}
	for _, item := range config.RequiredPlugins {
		if v, ok := pluginCounts[item]; ok {
			requiringCheck[item] = v
		} else {
			requiringCheck[item] = 0
		}
	}

	// Make sure all plugin constraints are satisfied
	for pluginType, pluginRepo := range pluginRepos {
		if err := pluginRepo.Constraints().Check(requiringCheck[pluginType]); err != nil {
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
