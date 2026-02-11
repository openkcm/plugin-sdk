package catalog

type BuiltInPluginRetriever interface {
	Retrieve() []BuiltInPlugin
}

type BuiltInPluginRegistry interface {
	BuiltInPluginRetriever

	Register(plugin BuiltInPlugin)
}

type buildInRegistry struct {
	plugins []BuiltInPlugin
}

func CreateBuiltInPluginRegistry() BuiltInPluginRegistry {
	return &buildInRegistry{
		plugins: make([]BuiltInPlugin, 0),
	}
}

func (r *buildInRegistry) Register(plugin BuiltInPlugin) {
	r.plugins = append(r.plugins, plugin)
}

func (r *buildInRegistry) Retrieve() []BuiltInPlugin {
	plugins := make([]BuiltInPlugin, 0, len(r.plugins))
	plugins = append(plugins, r.plugins...)
	return r.plugins
}
