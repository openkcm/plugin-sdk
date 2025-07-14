package catalog

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"sort"

	"github.com/openkcm/plugin-sdk/api"
	"github.com/openkcm/plugin-sdk/internal/bootstrap"
	"github.com/openkcm/plugin-sdk/pkg/telemetry"
	"google.golang.org/grpc"
)

// Plugin is a loaded plugin.
type PluginBuildIn interface {
	// Closer is used to unload the plugin. Any facades initialized by the
	// call to bind are considered invalidated after the plugin is closed.
	io.Closer

	// Bind binds the given facades to the plugin. It also returns a Configurer
	// that can be used to configure the plugin. If the plugin does not support
	// a given facade, an error will be returned. This function is designed
	// only for use by unit-tests for built-in plugin implementations or fake
	// facade implementations that rely on built-ins.
	Bind(facades ...Facade) (Configurer, error)
}

type pluginImpl struct {
	closerGroup

	conn             grpc.ClientConnInterface
	info             PluginInfo
	log              *slog.Logger
	grpcServiceNames []string
}

func newPluginBuiltIn(ctx context.Context, conn grpc.ClientConnInterface, info PluginInfo, log *slog.Logger, closers closerGroup, hostServices []api.ServiceServer) (*pluginImpl, error) {
	grpcServiceNames, err := initPlugin(ctx, conn, hostServices)
	if err != nil {
		return nil, err
	}

	closers = append(closers, closerFunc(func() {
		ctx, cancel := context.WithTimeout(context.Background(), deinitTimeout)
		defer cancel()
		if err := bootstrap.Deinit(ctx, conn); err != nil {
			log.Error("Failed to deinitialize plugin", "error", err)
		} else {
			log.Debug("Plugin deinitialized")
		}
	}))

	return &pluginImpl{
		conn:             conn,
		info:             info,
		log:              log,
		closerGroup:      closers,
		grpcServiceNames: grpcServiceNames,
	}, nil
}

// Bind implements the Plugin interface method of the same name.
func (p *pluginImpl) Bind(facades ...Facade) (Configurer, error) {
	grpcServiceNames := grpcServiceNameSet(p.grpcServiceNames)

	for _, facade := range facades {
		if _, ok := grpcServiceNames[facade.GRPCServiceName()]; !ok {
			return nil, fmt.Errorf("plugin does not support facade service %q", facade.GRPCServiceName())
		}
		p.initFacade(facade)
	}

	configurer, err := p.makeConfigurer(grpcServiceNames)
	if err != nil {
		return nil, err
	}
	if configurer == nil {
		configurer = configurerUnsupported{}
	}
	return configurer, nil
}

func (p *pluginImpl) bindFacade(repo bindable, facade Facade) any {
	impl := p.initFacade(facade)
	repo.bind(facade)
	return impl
}

func (p *pluginImpl) initFacade(facade Facade) any {
	facade.InitInfo(p.info)
	facade.InitLog(p.log)
	return facade.InitClient(p.conn)
}

func (p *pluginImpl) bindRepos(pluginRepo bindablePluginRepo, serviceRepos []bindableServiceRepo) (Configurer, error) {
	grpcServiceNames := grpcServiceNameSet(p.grpcServiceNames)

	impl := p.bindRepo(pluginRepo, grpcServiceNames)
	for _, serviceRepo := range serviceRepos {
		p.bindRepo(serviceRepo, grpcServiceNames)
	}

	configurer, err := p.makeConfigurer(grpcServiceNames)
	if err != nil {
		return nil, err
	}

	switch {
	case impl == nil:
		return nil, fmt.Errorf("no supported plugin interface found in: %q", p.grpcServiceNames)
	case len(grpcServiceNames) > 0:
		for _, grpcServiceName := range sortStringSet(grpcServiceNames) {
			p.log.With(telemetry.PluginService, grpcServiceName).Warn("Unsupported plugin service found")
		}
	}

	return configurer, nil
}

func (p *pluginImpl) makeConfigurer(grpcServiceNames map[string]struct{}) (Configurer, error) {
	repo := new(configurerRepo)
	bindable, err := makeBindableServiceRepo(repo)
	if err != nil {
		return nil, err
	}
	p.bindRepo(bindable, grpcServiceNames)
	return repo.configurer, nil
}

func (p *pluginImpl) bindRepo(repo bindableServiceRepo, grpcServiceNames map[string]struct{}) any {
	versions := repo.Versions()

	var impl any
	for _, version := range versions {
		facade := version.New()
		if _, ok := grpcServiceNames[facade.GRPCServiceName()]; ok {
			delete(grpcServiceNames, facade.GRPCServiceName())
			// Use the first matching version (in case the plugin implements
			// more than one). The rest will be removed from the list of
			// service names above so we can properly warn of unhandled
			// services without false negatives.
			if impl != nil {
				continue
			}
			warnIfDeprecated(p.log, version, versions[0])
			impl = p.bindFacade(repo, facade)
		}
	}
	return impl
}

func warnIfDeprecated(log *slog.Logger, thisVersion, latestVersion Version) {
	if thisVersion.Deprecated() {
		log = log.With(telemetry.PluginName, thisVersion.New().GRPCServiceName())
		if !latestVersion.Deprecated() {
			log = log.With(telemetry.PluginName, latestVersion.New().GRPCServiceName())
		}
		log.Warn("Service is deprecated and will be removed in a future release")
	}
}

func grpcServiceNameSet(grpcServiceNames []string) map[string]struct{} {
	set := make(map[string]struct{})
	for _, grpcServiceName := range grpcServiceNames {
		set[grpcServiceName] = struct{}{}
	}
	return set
}

func sortStringSet(set map[string]struct{}) []string {
	ss := make([]string, 0, len(set))
	for s := range set {
		ss = append(ss, s)
	}
	sort.Strings(ss)
	return ss
}
