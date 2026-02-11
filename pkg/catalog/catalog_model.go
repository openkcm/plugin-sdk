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

// Repository is a set of plugin and service repositories.
type Repository interface {
	// Plugins returns a map of plugin repositories, keyed by the plugin type.
	Plugins() map[string]PluginRepo

	// Services returns service repositories.
	Services() []ServiceRepo
}

// PluginRepo is a repository of plugin facades for a given plugin type.
type PluginRepo interface {
	ServiceRepo

	// Constraints returns the constraints required by the plugin repository.
	// The Load function will ensure that these constraints are satisfied before
	// returning successfully.
	Constraints() Constraints

	// BuiltIns provides the list of built ins that are available for the
	// given plugin repository.
	BuiltIns() []BuiltInPlugin
}

// ServiceRepo is a repository for service facades for a given service.
type ServiceRepo interface {
	// Binder returns a function that is used by the catalog system to "bind"
	// the facade returned by selected version to the repository. It MUST
	// return void and take a single argument of type X, where X can be
	// assigned to by any of the facade implementation types returned by the
	// provided versions (see Versions).
	Binder() any

	// Versions returns the versions supported by the repository, ordered by
	// most to least preferred. The first version supported by the plugin will
	// be used. When a deprecated version is bound, warning messaging will
	// recommend the first version in the list as a replacement, unless it is
	// also deprecated.
	Versions() []api.Version

	// Clear is called when loading fails to clear the repository of any
	// previously bound facades.
	Clear()
}

type Config struct {
	// Logger is the logger. It is used for general purpose logging and also
	// provided to the plugins.
	Logger *slog.Logger

	// PluginConfigs is the list of plugin configurations.
	PluginConfigs []PluginConfig

	// HostServices are the servers for host services provided by SPIRE to
	// plugins.
	HostServices []api.ServiceServer

	RequiredPlugins []string
}
