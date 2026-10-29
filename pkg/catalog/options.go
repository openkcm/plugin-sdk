package catalog

import (
	"context"
	"log/slog"

	"google.golang.org/grpc"
)

// WithErrorLogging returns a PluginOption that logs every gRPC error returned by the plugin
func WithErrorLogging(logger *slog.Logger) PluginOption {
	return func(cfg *PluginConfig) {
		cfg.unaryInterceptors = append(cfg.unaryInterceptors, func(
			ctx context.Context,
			method string,
			req, reply any,
			cc *grpc.ClientConn,
			invoker grpc.UnaryInvoker,
			opts ...grpc.CallOption,
		) error {
			err := invoker(ctx, method, req, reply, cc, opts...)
			if err != nil {
				logger.ErrorContext(ctx, "plugin method returned error",
					slog.String("pluginName", cfg.Name),
					slog.String("method", method),
					slog.String("error", err.Error()),
				)
			}
			return err
		})
	}
}
