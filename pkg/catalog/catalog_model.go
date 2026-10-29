package catalog

import (
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

	// PluginOptions are applied to every PluginConfig when the catalog is
	// created. Use them to inject cross-cutting behaviour such as a middleware.
	PluginOptions []PluginOption
}

// PluginOption mutates a PluginConfig. Options supplied via Config.PluginOptions
// are applied to every plugin before it is loaded.
type PluginOption func(*PluginConfig)
