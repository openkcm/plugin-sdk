package catalog

type BuiltInPluginRegistry interface {
	Register(plugin BuiltIn)
	Get() []BuiltIn
}

type buildInRegistry struct {
	plugins []BuiltIn
}

func DefaultBuiltInPluginRegistry() BuiltInPluginRegistry {
	return &buildInRegistry{
		plugins: make([]BuiltIn, 0),
	}
}

func (r *buildInRegistry) Register(plugin BuiltIn) {
	r.plugins = append(r.plugins, plugin)
}

func (r *buildInRegistry) Get() []BuiltIn {
	plugins := make([]BuiltIn, 0, len(r.plugins))
	plugins = append(plugins, r.plugins...)
	return plugins
}
