package catalog

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"sort"

	"google.golang.org/grpc"

	goplugin "github.com/hashicorp/go-plugin"

	"github.com/openkcm/plugin-sdk/api"
	"github.com/openkcm/plugin-sdk/internal/bootstrap"
	"github.com/openkcm/plugin-sdk/internal/slog2hclog"
)

type PluginConfigs []PluginConfig

type PluginConfig struct {
	// Name of the plugin
	Name string

	// Type is the plugin type
	Type string

	// Path is the path on disk to the plugin.
	Path string

	// Args are the command line arguments to supply to the plugin
	Args []string

	// Env is the environment variables to supply to the plugin
	Env map[string]string

	// Checksum is the hex-encoded SHA256 hash of the plugin binary.
	Checksum string

	Version uint32

	DataSource DataSource

	YamlConfiguration string

	LogLevel string

	Disabled bool

	Logger *slog.Logger

	HostServices []api.ServiceServer

	// Tags are the metadata associated with a plugin these can be used to filter plugins later e.g. ['FeatureA'] on client side.
	Tags []string
}

func (c *PluginConfig) IsExternal() bool {
	return c.Path != ""
}

func (c *PluginConfig) IsEnabled() bool {
	return !c.Disabled
}

type DataSource interface {
	Load() (string, error)
	IsDynamic() bool
}

type FixedData string

func (d FixedData) Load() (string, error) {
	return string(d), nil
}

func (d FixedData) IsDynamic() bool {
	return false
}

type FileData string

func (d FileData) Load() (string, error) {
	data, err := os.ReadFile(string(d))
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (d FileData) IsDynamic() bool {
	return true
}

type Build interface {
	SetValue(string)
}

type Plugin interface {
	io.Closer

	ClientConnection() grpc.ClientConnInterface
	Info() api.Info
	Logger() *slog.Logger
	GrpcServiceNames() []string
}

type pluginImpl struct {
	closerGroup

	conn             grpc.ClientConnInterface
	info             api.Info
	logger           *slog.Logger
	grpcServiceNames []string
}

func (p *pluginImpl) Close() error {
	return p.closerGroup.Close()
}
func (p *pluginImpl) ClientConnection() grpc.ClientConnInterface {
	return p.conn
}
func (p *pluginImpl) Info() api.Info {
	return p.info
}
func (p *pluginImpl) Logger() *slog.Logger {
	return p.logger
}
func (p *pluginImpl) GrpcServiceNames() []string {
	return p.grpcServiceNames
}

func loadPlugin(ctx context.Context, config PluginConfig) (*pluginImpl, error) {
	config.Logger.InfoContext(ctx, "Loading plugin", "name", config.Name, "path", config.Path)

	cmd := pluginCmd(config.Path, config.Args...)
	injectEnv(config, cmd)

	// Create the secure config based on the (optional) checksum
	seccfg, err := buildSecureConfig(config.Checksum)
	if err != nil {
		return nil, fmt.Errorf("invalid checksum: %w", err)
	}

	// Start the plugin client

	pluginClient := goplugin.NewClient(&goplugin.ClientConfig{
		SecureConfig: seccfg,
		Logger:       slog2hclog.NewWithLevel(config.Logger, config.LogLevel),
		HandshakeConfig: goplugin.HandshakeConfig{
			ProtocolVersion:  1,
			MagicCookieKey:   config.Type,
			MagicCookieValue: config.Type,
		},
		AutoMTLS:         true,
		Plugins:          map[string]goplugin.Plugin{config.Name: &HCRPCPlugin{config: config}},
		Cmd:              cmd,
		AllowedProtocols: []goplugin.Protocol{goplugin.ProtocolGRPC},
	})

	// Connect via RPC
	rpcClient, err := pluginClient.Client()
	if err != nil {
		pluginClient.Kill()
		return nil, err
	}

	// Request the plugin
	rawPlugin, err := rpcClient.Dispense(config.Name)
	if err != nil {
		pluginClient.Kill()
		return nil, err
	}
	plugin, ok := rawPlugin.(*HCPlugin)
	// Purely defensive. This should never happen since we control what
	// gets returned from hcClientPlugin.
	if !ok {
		return nil, fmt.Errorf("expected %T, got %T", plugin, rawPlugin)
	}

	// Plugin has been loaded and initialized. Ensure the plugin client is
	// killed when the plugin is closed.
	plugin.closers = append(plugin.closers, closerFunc(pluginClient.Kill))

	var version uint = 1
	if config.Version > 1 {
		version = uint(config.Version)
	}
	info := &pluginInfo{
		name:    config.Name,
		typ:     config.Type,
		tags:    config.Tags,
		version: version,
	}

	return newPlugin(ctx, plugin.conn, info, config.Logger, plugin.closers, config.HostServices)

}

// injectEnv injects the environment variables into the command
func injectEnv(config PluginConfig, cmd *exec.Cmd) {
	if len(config.Env) != 0 {
		for key, val := range config.Env {
			cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", key, val))
		}
	}
}

type pluginInfo struct {
	name      string
	typ       string
	buildInfo string
	tags      []string
	version   uint
}

func (info *pluginInfo) Name() string {
	return info.name
}

func (info *pluginInfo) Type() string {
	return info.typ
}

func (info *pluginInfo) Tags() []string { return info.tags }

func (info *pluginInfo) Build() string { return info.buildInfo }

func (info *pluginInfo) Version() uint {
	return info.version
}

func (info *pluginInfo) SetValue(value string) {
	info.buildInfo = value
}

type pluginCloser struct {
	plugin io.Closer
	log    *slog.Logger
}

func (c pluginCloser) Close() error {
	c.log.Info("Plugins unloading")
	if err := c.plugin.Close(); err != nil {
		c.log.Error("Failed to unload plugin", "error", err)
		return err
	}
	c.log.Info("Plugins unloaded")
	return nil
}

func buildSecureConfig(checksum string) (*goplugin.SecureConfig, error) {
	var seccfg *goplugin.SecureConfig
	if checksum == "" {
		return seccfg, nil
	}

	sum, err := hex.DecodeString(checksum)
	if err != nil {
		return nil, errors.New("checksum is not a valid hex string")
	}

	hash := sha256.New()
	if len(sum) != hash.Size() {
		return nil, fmt.Errorf("expected checksum of length %d; got %d", hash.Size()*2, len(sum)*2)
	}

	return &goplugin.SecureConfig{
		Checksum: sum,
		Hash:     sha256.New(),
	}, nil
}

func newPlugin(ctx context.Context, conn grpc.ClientConnInterface, info api.Info, logger *slog.Logger, closers closerGroup, hostServices []api.ServiceServer) (*pluginImpl, error) {
	grpcServiceNames, err := initPlugin(ctx, conn, hostServices)
	if err != nil {
		return nil, err
	}

	closers = append(closers, closerFunc(func() {
		ctx, cancel := context.WithTimeout(context.Background(), deinitTimeout)
		defer cancel()
		if err := bootstrap.Deinit(ctx, conn); err != nil {
			logger.ErrorContext(ctx, "Failed to deinitialize plugin", "error", err)
		} else {
			logger.Debug("Plugin deinitialized")
		}
	}))

	return &pluginImpl{
		closerGroup: closers,

		conn:             conn,
		info:             info,
		logger:           logger,
		grpcServiceNames: grpcServiceNames,
	}, nil
}

// Bind implements the Plugin interface method of the same name.
func (p *pluginImpl) Bind(facades ...api.Facade) (Configurer, error) {
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
	return configurer, nil
}

func initPlugin(ctx context.Context, conn grpc.ClientConnInterface, hostServices []api.ServiceServer) ([]string, error) {
	var hostServiceGRPCServiceNames []string
	for _, hostService := range hostServices {
		hostServiceGRPCServiceNames = append(hostServiceGRPCServiceNames, hostService.GRPCServiceName())
	}
	ctx, cancel := context.WithTimeout(ctx, initTimeout)
	defer cancel()
	return bootstrap.Init(ctx, conn, hostServiceGRPCServiceNames)
}

func (p *pluginImpl) makeConfigurer(grpcServiceNames map[string]struct{}) (Configurer, error) {
	repo := new(configurerRepo)
	bindable, err := makeBindableServiceRepo(repo)
	if err != nil {
		return nil, err
	}

	_, err = p.bindRepo(bindable, grpcServiceNames)
	if err != nil {
		return nil, err
	}

	return repo.configurer, nil
}

func (p *pluginImpl) bindRepo(repo bindableServiceRepo, grpcServiceNames map[string]struct{}) (any, error) {
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
			warnIfDeprecated(p.logger, version, versions[0])
			impl = p.bindFacade(repo, facade)

			if facade.Version() == p.info.Version() {
				break
			} else {
				impl = nil
			}
		}
	}

	if impl == nil {
		return nil, fmt.Errorf("requested `%d` version of plugin type `%s` wrapper binding implementation not found",
			p.info.Version(),
			p.info.Type(),
		)
	}

	return impl, nil
}

func (p *pluginImpl) bindFacade(repo bindable, facade api.Facade) any {
	impl := p.initFacade(facade)
	repo.bind(facade)
	return impl
}

func (p *pluginImpl) bindRepos(pluginRepo bindablePluginRepo, serviceRepos []bindableServiceRepo) (Configurer, error) {
	grpcServiceNames := grpcServiceNameSet(p.grpcServiceNames)

	impl, err := p.bindRepo(pluginRepo, grpcServiceNames)
	if err != nil {
		return nil, err
	}
	for _, serviceRepo := range serviceRepos {
		_, err := p.bindRepo(serviceRepo, grpcServiceNames)
		if err != nil {
			return nil, err
		}
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
			p.logger.With("plugin_service", grpcServiceName).Warn("Unsupported plugin service found")
		}
	}

	return configurer, nil
}

func warnIfDeprecated(log *slog.Logger, thisVersion, latestVersion api.Version) {
	if thisVersion.Deprecated() {
		log.Warn("Service is deprecated and will be removed in a future release")
	}
}

func (p *pluginImpl) initFacade(facade api.Facade) any {
	facade.InitInfo(p.info)
	facade.InitLog(p.logger)
	return facade.InitClient(p.conn)
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
