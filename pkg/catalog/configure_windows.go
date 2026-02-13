package catalog

import (
	"context"
	"log/slog"
)

func ReconfigureOnSignal(ctx context.Context, _ *slog.Logger, _ Reconfigurer) error {
	// TODO: maybe drive this using an event?
	<-ctx.Done()
	return ctx.Err()
}
