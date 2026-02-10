package plugin

import (
	"log/slog"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/openkcm/plugin-sdk/api"
)

// PrefixMessage prefixes the given message with plugin information. The prefix
// is only applied if it is not already applied.
func PrefixMessage(pluginInfo api.Info, message string) string {
	message, _ = prefixMessage(pluginInfo, message)
	return message
}

// Facade is embedded by plugin interface facade implementations as a
// convenient way to embed Info but also provide a set of convenient
// functions for embellishing and generating errors that have the plugin
// name prefixed.
type Facade struct {
	api.Info
	Log *slog.Logger
}

// FixedFacade is a helper that creates a facade from fixed information, i.e.
// not the product of a loaded plugin.
func FixedFacade(pluginName, pluginType string, log *slog.Logger) Facade {
	return Facade{
		Info: pluginFacadeInfo{
			pluginName: pluginName,
			pluginType: pluginType,
		},
		Log: log,
	}
}

// InitInfo partially satisfies the catalog.Facade interface
func (f *Facade) InitInfo(pluginInfo api.Info) {
	f.Info = pluginInfo
}

// InitLog partially satisfies the catalog.Facade interface
func (f *Facade) InitLog(log *slog.Logger) {
	f.Log = log
}

// WrapErr wraps a given error such that it will be prefixed with the plugin
// name. This method should be used by facade implementations to wrap errors
// that come out of plugin implementations.
func (f *Facade) WrapErr(err error) error {
	if err == nil {
		return nil
	}

	// Embellish the gRPC status with the prefix, if necessary.
	if st, ok := status.FromError(err); ok {
		// Care must be taken to preserve any status details. Therefore, the
		// proto is embellished directly and a new status created from that
		// proto.
		pb := st.Proto()
		if message, ok := prefixMessage(f, pb.Message); ok {
			pb.Message = message
			return status.FromProto(pb).Err()
		}
		return err
	}

	// Embellish the normal error with the prefix, if necessary. This is a
	// defensive measure since plugins go over gRPC.
	if message, ok := prefixMessage(f, err.Error()); ok {
		return &facadeError{wrapped: err, message: message}
	}

	return err
}

// Error creates a gRPC status with the given code and message. The message
// will be prefixed with the plugin name.
func (f *Facade) Error(code codes.Code, message string) error {
	return status.Error(code, messagePrefix(f)+message)
}

// Errorf creates a gRPC status with the given code and
// formatted message. The message will be prefixed with the plugin name.
func (f *Facade) Errorf(code codes.Code, format string, args ...any) error {
	return status.Errorf(code, messagePrefix(f)+format, args...)
}

func prefixMessage(pluginInfo api.Info, message string) (string, bool) {
	prefix := messagePrefix(pluginInfo)

	if strings.HasPrefix(message, prefix) {
		return message, false
	}

	oldPrefix := pluginInfo.Name() + ": "
	return prefix + strings.TrimPrefix(message, oldPrefix), true
}

func messagePrefix(pluginInfo api.Info) string {
	return strings.ToLower(pluginInfo.Type()) + "(" + pluginInfo.Name() + "): "
}

type facadeError struct {
	wrapped error
	message string
}

func (e *facadeError) Error() string {
	return e.message
}

func (e *facadeError) Unwrap() error {
	return e.wrapped
}

type pluginFacadeInfo struct {
	pluginName string
	pluginType string
	pluginTags []string
	buildInfo  string
	version    uint
}

func (info pluginFacadeInfo) Name() string {
	return info.pluginName
}

func (info pluginFacadeInfo) Type() string {
	return info.pluginType
}

func (info pluginFacadeInfo) Tags() []string {
	return info.pluginTags
}

func (info pluginFacadeInfo) Build() string {
	return info.buildInfo
}

func (info pluginFacadeInfo) Version() uint {
	return info.version
}
