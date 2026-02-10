package catalog

import (
	"context"
	"errors"
	"log/slog"
	"os/exec"
	"testing"

	"google.golang.org/grpc"

	"github.com/openkcm/plugin-sdk/api"
)

//
// ─────────────────────────────────────────────
// Test doubles
// ─────────────────────────────────────────────
//

type fakeCloser struct {
	closed bool
	err    error
}

func (f *fakeCloser) Close() error {
	f.closed = true
	return f.err
}

type fakePluginServiceServer struct {
	name string
}

func (f *fakePluginServiceServer) GRPCServiceName() string { return f.name }
func (f *fakePluginServiceServer) RegisterServer(*grpc.Server) any {
	return nil
}

var _ api.ServiceServer = (*fakePluginServiceServer)(nil)

type discardPluginWriter struct{}

func (discardPluginWriter) Write(p []byte) (int, error) { return len(p), nil }

//
// ─────────────────────────────────────────────
// PluginConfig
// ─────────────────────────────────────────────
//

func TestPluginConfigFlags(t *testing.T) {
	t.Parallel()

	t.Run("IsExternal", func(t *testing.T) {
		t.Parallel()

		if (&PluginConfig{Path: ""}).IsExternal() {
			t.Fatal("expected internal plugin")
		}
		if !(&PluginConfig{Path: "/bin/plugin"}).IsExternal() {
			t.Fatal("expected external plugin")
		}
	})

	t.Run("IsEnabled", func(t *testing.T) {
		t.Parallel()

		if (&PluginConfig{Disabled: true}).IsEnabled() {
			t.Fatal("expected disabled plugin")
		}
		if !(&PluginConfig{Disabled: false}).IsEnabled() {
			t.Fatal("expected enabled plugin")
		}
	})
}

//
// ─────────────────────────────────────────────
// injectEnv
// ─────────────────────────────────────────────
//

func TestInjectEnv(t *testing.T) {
	t.Parallel()

	cmd := exec.Command("test")
	cmd.Env = []string{
		"PATH=/usr/bin:/bin",
	}

	cfg := PluginConfig{
		Env: map[string]string{
			"A": "1",
			"B": "2",
		},
	}

	injectEnv(cfg, cmd)

	if len(cmd.Env) != 3 {
		t.Fatalf("expected 3 env vars, got %d", len(cmd.Env))
	}
}

//
// ─────────────────────────────────────────────
// pluginInfo
// ─────────────────────────────────────────────
//

func TestPluginInfo(t *testing.T) {
	t.Parallel()

	info := &pluginInfo{
		name: "n",
		typ:  "t",
		tags: []string{"a", "b"},
	}

	if info.Name() != "n" {
		t.Fatal("Name mismatch")
	}
	if info.Type() != "t" {
		t.Fatal("Type mismatch")
	}
	if len(info.Tags()) != 2 {
		t.Fatal("Tags mismatch")
	}

	info.SetValue("v1")
	if info.Build() != "v1" {
		t.Fatal("Build mismatch")
	}
}

//
// ─────────────────────────────────────────────
// pluginImpl
// ─────────────────────────────────────────────
//

func TestPluginStruct(t *testing.T) {
	t.Parallel()

	closer := &fakeCloser{}
	var cg closerGroup
	cg = append(cg, closer)

	p := &pluginImpl{
		closerGroup:      cg,
		conn:             nil,
		info:             &pluginInfo{name: "p"},
		logger:           slog.New(slog.NewTextHandler(discardPluginWriter{}, nil)),
		grpcServiceNames: []string{"svc1", "svc2"},
	}

	if p.Info().Name() != "p" {
		t.Fatal("Info mismatch")
	}
	if len(p.GrpcServiceNames()) != 2 {
		t.Fatal("service names mismatch")
	}

	if err := p.Close(); err != nil {
		t.Fatalf("Close failed: %v", err)
	}
	if !closer.closed {
		t.Fatal("expected closer to be called")
	}
}

//
// ─────────────────────────────────────────────
// pluginCloser
// ─────────────────────────────────────────────
//

func TestPluginCloser(t *testing.T) {
	t.Parallel()

	t.Run("successful close", func(t *testing.T) {
		t.Parallel()

		c := &fakeCloser{}
		pc := pluginCloser{
			plugin: c,
			log:    slog.New(slog.NewTextHandler(discardPluginWriter{}, nil)),
		}

		if err := pc.Close(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !c.closed {
			t.Fatal("plugin was not closed")
		}
	})

	t.Run("close error is propagated", func(t *testing.T) {
		t.Parallel()

		c := &fakeCloser{err: errors.New("boom")}
		pc := pluginCloser{
			plugin: c,
			log:    slog.New(slog.NewTextHandler(discardPluginWriter{}, nil)),
		}

		if err := pc.Close(); err == nil {
			t.Fatal("expected error")
		}
	})
}

//
// ─────────────────────────────────────────────
// buildSecureConfig
// ─────────────────────────────────────────────
//

func TestBuildSecureConfig(t *testing.T) {
	t.Parallel()

	t.Run("empty checksum", func(t *testing.T) {
		t.Parallel()

		cfg, err := buildSecureConfig("")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cfg != nil {
			t.Fatal("expected nil secure config")
		}
	})

	t.Run("invalid hex", func(t *testing.T) {
		t.Parallel()

		_, err := buildSecureConfig("zzz")
		if err == nil {
			t.Fatal("expected error")
		}
	})

	t.Run("wrong length", func(t *testing.T) {
		t.Parallel()

		_, err := buildSecureConfig("deadbeef")
		if err == nil {
			t.Fatal("expected error")
		}
	})

	t.Run("valid checksum", func(t *testing.T) {
		t.Parallel()

		_ = make([]byte, 32)
		hex := make([]byte, 64)
		for i := range hex {
			hex[i] = 'a'
		}

		cfg, err := buildSecureConfig(string(hex))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cfg == nil {
			t.Fatal("expected secure config")
		}
	})
}

//
// ─────────────────────────────────────────────
// initPlugin (failure path)
// ─────────────────────────────────────────────
//

func TestInitPluginFailure(t *testing.T) {
	t.Parallel()

	var err error
	defer func() {
		if r := recover(); r == nil {
			err = errors.New("nil grpc connection")
		}
	}()

	_, err = initPlugin(
		context.Background(),
		nil, // invalid conn triggers failure
		[]api.ServiceServer{
			&fakePluginServiceServer{name: "svc"},
		},
	)

	if err == nil {
		t.Fatal("expected error")
	}
}

//
// ─────────────────────────────────────────────
// newPlugin (failure path)
// ─────────────────────────────────────────────
//

func TestNewPluginInitFailure(t *testing.T) {
	t.Parallel()

	var err error
	defer func() {
		if r := recover(); r == nil {
			err = errors.New("nil grpc connection")
		}
	}()

	_, err = newPlugin(
		context.Background(),
		nil,
		&pluginInfo{name: "p"},
		slog.New(slog.NewTextHandler(discardPluginWriter{}, nil)),
		nil,
		nil,
	)

	if err == nil {
		t.Fatal("expected error")
	}
}

//
// ─────────────────────────────────────────────
// loadPlugin (early failure)
// ─────────────────────────────────────────────
//

func TestLoadPluginInvalidChecksum(t *testing.T) {
	t.Parallel()

	cfg := PluginConfig{
		Name:     "p",
		Type:     "t",
		Path:     "/does/not/exist",
		Checksum: "invalid",
		Logger:   slog.New(slog.NewTextHandler(discardPluginWriter{}, nil)),
	}

	_, err := loadPlugin(context.Background(), cfg)
	if err == nil {
		t.Fatal("expected error")
	}
}
