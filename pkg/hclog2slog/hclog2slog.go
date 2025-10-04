// Package hclog2slog provides an adapter that wraps a HashiCorp hclog.Logger
// and exposes it as a Go slog.Logger.
//
// This allows you to integrate libraries that use the modern slog API
// with systems or dependencies that still rely on HashiCorp's hclog.

// The adapter maps slog levels to hclog levels and preserves key/value fields.
package hclog2slog

import (
	"context"
	"log/slog"

	"github.com/hashicorp/go-hclog"
)

// levelMapToHclog converts slog.Level to hclog.Level.
var levelMapToHclog = map[slog.Level]hclog.Level{
	slog.LevelDebug - 4: hclog.Trace,
	slog.LevelDebug:     hclog.Debug,
	slog.LevelInfo:      hclog.Info,
	slog.LevelWarn:      hclog.Warn,
	slog.LevelError:     hclog.Error,
	slog.LevelError + 4: hclog.Off,
}

// hclogHandler implements slog.Handler by delegating to an hclog.Logger.
type hclogHandler struct {
	hcl  hclog.Logger
	lvl  slog.Level
	args []any
}

// Enabled checks whether the given level is enabled on the underlying hclog.Logger.
func (h *hclogHandler) Enabled(_ context.Context, level slog.Level) bool {
	hclLevel := levelMapToHclog[level]
	switch hclLevel {
	case hclog.Trace:
		return h.hcl.IsTrace()
	case hclog.Debug:
		return h.hcl.IsDebug()
	case hclog.Info:
		return h.hcl.IsInfo()
	case hclog.Warn:
		return h.hcl.IsWarn()
	case hclog.Error:
		return h.hcl.IsError()
	default:
		return true
	}
}

// Handle forwards the slog.Record to the underlying hclog.Logger.
func (h *hclogHandler) Handle(_ context.Context, r slog.Record) error {
	var kv []any

	r.Attrs(func(a slog.Attr) bool {
		kv = append(kv, a.Key, a.Value.Any())
		return true
	})

	// Include any pre-bound args
	if len(h.args) > 0 {
		kv = append(kv, h.args...)
	}

	hclLevel := levelMapToHclog[r.Level]
	h.hcl.Log(hclLevel, r.Message, kv...)

	return nil
}

// WithAttrs creates a new handler with additional attributes.
func (h *hclogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	args := make([]any, 0, len(attrs)*2)
	for _, a := range attrs {
		args = append(args, a.Key, a.Value.Any())
	}
	return &hclogHandler{
		hcl:  h.hcl,
		lvl:  h.lvl,
		args: append(h.args, args...),
	}
}

// WithGroup creates a new handler that represents a logical sub-logger.
func (h *hclogHandler) WithGroup(name string) slog.Handler {
	return &hclogHandler{
		hcl:  h.hcl.Named(name),
		lvl:  h.lvl,
		args: h.args,
	}
}

// New creates a new slog.Logger backed by an hclog.Logger.
//
// The returned *slog.Logger will send all log events to the given hclog.Logger
// with compatible levels and fields.
func New(hclLogger hclog.Logger) *slog.Logger {
	handler := &hclogHandler{
		hcl: hclLogger,
		lvl: slog.LevelInfo,
	}
	return slog.New(handler)
}

// AsSlog wraps an hclog.Logger as a *slog.Logger.
func AsSlog(h hclog.Logger) *slog.Logger {
	return New(h)
}
