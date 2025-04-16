package catalog

import "context"

type pluginNameKey struct{}

func WithPluginName(ctx context.Context, name string) context.Context {
	return context.WithValue(ctx, pluginNameKey{}, name)
}
