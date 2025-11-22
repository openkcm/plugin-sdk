package catalog

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os/exec"

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

func (c PluginConfig) IsEnabled() bool {
	return !c.Disabled
}

type Build interface {
	SetValue(string)
}

// PluginInfo provides the information for the loaded plugin.
type PluginInfo interface {
	// The name of the plugin
	Name() string

	// The type of the plugin
	Type() string

	// Tags associated with the plugin
	Tags() []string

	// Build of the plugin
	Build() string
}

type Plugin struct {
	closerGroup

	conn             grpc.ClientConnInterface
	info             PluginInfo
	logger           *slog.Logger
	grpcServiceNames []string
}

func (p *Plugin) ClientConnection() grpc.ClientConnInterface {
	return p.conn
}
func (p *Plugin) Info() PluginInfo {
	return p.info
}
func (p *Plugin) Logger() *slog.Logger {
	return p.logger
}
func (p *Plugin) GrpcServiceNames() []string {
	return p.grpcServiceNames
}

func loadPlugin(ctx context.Context, config PluginConfig) (*Plugin, error) {
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
	info := pluginInfo{
		name: config.Name,
		typ:  config.Type,
		tags: config.Tags,
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
}

func (info pluginInfo) Name() string {
	return info.name
}

func (info pluginInfo) Type() string {
	return info.typ
}

func (info pluginInfo) Tags() []string { return info.tags }

func (info pluginInfo) Build() string { return info.buildInfo }

//nolint:staticcheck
func (info pluginInfo) SetValue(value string) {
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

func newPlugin(ctx context.Context, conn grpc.ClientConnInterface, info PluginInfo, logger *slog.Logger, closers closerGroup, hostServices []api.ServiceServer) (*Plugin, error) {
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

	return &Plugin{
		closerGroup: closers,

		conn:             conn,
		info:             info,
		logger:           logger,
		grpcServiceNames: grpcServiceNames,
	}, nil
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
