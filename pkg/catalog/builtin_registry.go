package catalog

type builtInPluginRetriever interface {
	retrieve() []BuiltInPlugin
}

type BuiltInPluginRegistry interface {
	builtInPluginRetriever

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

func (r *buildInRegistry) retrieve() []BuiltInPlugin {
	return r.plugins
}
