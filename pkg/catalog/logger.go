package catalog

import (
	"log/slog"

	"github.com/hashicorp/go-hclog"
	"github.com/magodo/slog2hclog"
)

func newSlog2HClog(logger *slog.Logger, logLevel string) hclog.Logger {
	log := slog2hclog.New(logger, new(slog.LevelVar))
	log.SetLevel(hclog.LevelFromString(logLevel))
	return log
}
