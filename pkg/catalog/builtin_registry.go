package catalog

type BuiltInPluginRegistry interface {
	Register(plugin BuiltInPlugin)
	Get() []BuiltInPlugin
}

type buildInRegistry struct {
	plugins []BuiltInPlugin
}

func DefaultBuiltInPluginRegistry() BuiltInPluginRegistry {
	return &buildInRegistry{
		plugins: make([]BuiltInPlugin, 0),
	}
}

func (r *buildInRegistry) Register(plugin BuiltInPlugin) {
	r.plugins = append(r.plugins, plugin)
}

func (r *buildInRegistry) Get() []BuiltInPlugin {
	plugins := make([]BuiltInPlugin, 0, len(r.plugins))
	plugins = append(plugins, r.plugins...)
	return plugins
}
