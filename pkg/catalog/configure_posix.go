//go:build !windows

package catalog

import (
	"context"
	"log/slog"
	"os"
	"os/signal"

	"golang.org/x/sys/unix"
)

func ReconfigureOnSignal(ctx context.Context, log *slog.Logger, reconfigurer Reconfigurer) error {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, unix.SIGUSR1)
	defer signal.Stop(ch)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ch:
			log.Info("Reconfigure signal received")
			reconfigurer.Reconfigure(ctx)
		}
	}
}
