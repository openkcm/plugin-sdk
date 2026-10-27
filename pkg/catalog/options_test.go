package catalog

import (
	"bytes"
	"context"
	"errors"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

func newTestLogger(buf *bytes.Buffer) *slog.Logger {
	return slog.New(slog.NewTextHandler(buf, nil))
}

func fakeInvoker(err error) grpc.UnaryInvoker {
	return func(_ context.Context, _ string, _, _ any, _ *grpc.ClientConn, _ ...grpc.CallOption) error {
		return err
	}
}

func TestWithErrorLogging(t *testing.T) {
	t.Parallel()

	t.Run("Should log on error", func(t *testing.T) {
		var buf bytes.Buffer
		cfg := &PluginConfig{Name: "my-plugin"}
		WithErrorLogging(newTestLogger(&buf))(cfg)

		err := cfg.unaryInterceptors[0](context.Background(), "/svc/Method", nil, nil, nil, fakeInvoker(errors.New("boom")))

		require.Error(t, err)
		assert.Contains(t, buf.String(), "my-plugin")
		assert.Contains(t, buf.String(), "boom")
		assert.Contains(t, buf.String(), "/svc/Method")
	})

	t.Run("Should not log on success", func(t *testing.T) {
		var buf bytes.Buffer
		cfg := &PluginConfig{Name: "my-plugin"}
		WithErrorLogging(newTestLogger(&buf))(cfg)

		require.NoError(t, cfg.unaryInterceptors[0](context.Background(), "/svc/Method", nil, nil, nil, fakeInvoker(nil)))
		assert.Empty(t, buf.String())
	})
}

func TestGetDialOptions_EmptyWithNoInterceptors(t *testing.T) {
	t.Parallel()

	assert.Empty(t, (&PluginConfig{}).getDialOptions())
}

func TestGetDialOptions_SingleChainForMultipleInterceptors(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	cfg := &PluginConfig{Name: "p"}
	WithErrorLogging(newTestLogger(&buf))(cfg)
	WithErrorLogging(newTestLogger(&buf))(cfg)

	assert.Len(t, cfg.unaryInterceptors, 2)
	assert.Len(t, cfg.getDialOptions(), 1)
}

func TestConfig_PluginOptionsAppliedToEachPlugin(t *testing.T) {
	t.Parallel()

	applied := make([]string, 0)
	cfg := Config{
		PluginOptions: []PluginOption{func(c *PluginConfig) { applied = append(applied, c.Name) }},
		PluginConfigs: []PluginConfig{{Name: "alpha"}, {Name: "beta"}},
	}

	for i := range cfg.PluginConfigs {
		for _, opt := range cfg.PluginOptions {
			opt(&cfg.PluginConfigs[i])
		}
	}

	require.Equal(t, []string{"alpha", "beta"}, applied)
}
