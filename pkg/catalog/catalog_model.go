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
}
