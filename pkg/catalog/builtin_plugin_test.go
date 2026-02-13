package catalog

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"testing"

	"google.golang.org/grpc"

	"github.com/openkcm/plugin-sdk/api"
)

//
// ─────────────────────────────────────────────
// Test doubles
// ─────────────────────────────────────────────
//

type fakePluginServer struct{}

var _ api.PluginServer = (*fakePluginServer)(nil)
var _ api.ServiceServer = (*fakePluginServer)(nil)

func (f *fakePluginServer) Type() string { return "fake" }
func (f *fakePluginServer) GRPCServiceName() string {
	return "fake"
}
func (f *fakePluginServer) RegisterServer(*grpc.Server) any {
	return nil
}

type fakeServiceServer struct{}

var _ api.ServiceServer = (*fakeServiceServer)(nil)

func (f *fakeServiceServer) GRPCServiceName() string {
	return "service"
}
func (f *fakeServiceServer) RegisterServer(*grpc.Server) any {
	return nil
}
func (f *fakeServiceServer) Register(*grpc.Server) {}

type panicPluginServer struct{}

var _ api.PluginServer = (*panicPluginServer)(nil)

func (p *panicPluginServer) Type() string {
	return "panic"
}
func (p *panicPluginServer) GRPCServiceName() string {
	return "panic"
}
func (p *panicPluginServer) RegisterServer(*grpc.Server) any {
	return errors.New("boom")
}

type discardWriter struct{}

func (discardWriter) Write(p []byte) (int, error) { return len(p), nil }

func TestBuiltInPluginStruct(t *testing.T) {
	t.Parallel()

	t.Run("all accessors and mutators", func(t *testing.T) {
		t.Parallel()

		plugin := &fakePluginServer{}
		service := &fakeServiceServer{}

		p := &builtInPluginStruct{
			name:      "test",
			tags:      []string{"a", "b"},
			plugin:    plugin,
			services:  []api.ServiceServer{service},
			buildInfo: "v1",
		}

		if got := p.Name(); got != "test" {
			t.Fatalf("Name(): want %q, got %q", "test", got)
		}
		if got := p.Type(); got != "fake" {
			t.Fatalf("Type(): want %q, got %q", "fake", got)
		}
		if got := p.Build(); got != "v1" {
			t.Fatalf("Build(): want %q, got %q", "v1", got)
		}
		if p.Plugin() != plugin {
			t.Fatalf("Plugin(): unexpected value")
		}
		if got := len(p.Services()); got != 1 {
			t.Fatalf("Services(): want 1, got %d", got)
		}
		if got := len(p.Tags()); got != 2 {
			t.Fatalf("Tags(): want 2, got %d", got)
		}

		p.SetValue("v2")
		if got := p.Build(); got != "v2" {
			t.Fatalf("SetValue(): want %q, got %q", "v2", got)
		}
	})
}

func TestMakeBuiltIn(t *testing.T) {
	t.Parallel()

	t.Run("creates builtin plugin wrapper", func(t *testing.T) {
		t.Parallel()

		plugin := &fakePluginServer{}
		s1 := &fakeServiceServer{}
		s2 := &fakeServiceServer{}

		b := MakeBuiltIn("plugin", plugin, s1, s2)

		if got := b.Name(); got != "plugin" {
			t.Fatalf("Name(): want %q, got %q", "plugin", got)
		}
		if b.Plugin() != plugin {
			t.Fatalf("Plugin(): unexpected value")
		}
		if got := len(b.Services()); got != 2 {
			t.Fatalf("Services(): want 2, got %d", got)
		}
	})
}

func TestBuiltinDialer(t *testing.T) {
	t.Parallel()

	t.Run("dial host reuses connection and closes", func(t *testing.T) {
		t.Parallel()

		log := slog.New(slog.NewTextHandler(discardWriter{}, nil))
		d := &builtinDialer{
			pluginName: "test",
			log:        log,
		}

		ctx := context.Background()

		conn1, err := d.DialHost(ctx)
		if err != nil {
			t.Fatalf("DialHost(): %v", err)
		}

		conn2, err := d.DialHost(ctx)
		if err != nil {
			t.Fatalf("DialHost() second call: %v", err)
		}

		if conn1 != conn2 {
			t.Fatalf("expected same connection on repeated DialHost")
		}

		if err := d.Close(); err != nil {
			t.Fatalf("Close(): %v", err)
		}
	})

	t.Run("close without connection is safe", func(t *testing.T) {
		t.Parallel()

		var d builtinDialer
		if err := d.Close(); err != nil {
			t.Fatalf("Close(): %v", err)
		}
	})
}

func TestNewBuiltInServer(t *testing.T) {
	t.Parallel()

	log := slog.New(slog.NewTextHandler(discardWriter{}, nil))

	server, closer := newBuiltInServer(log)
	if server == nil {
		t.Fatal("expected grpc.Server")
	}

	if err := closer.Close(); err != nil {
		t.Fatalf("Close(): %v", err)
	}
}

func TestStartPipeServer(t *testing.T) {
	t.Parallel()

	log := slog.New(slog.NewTextHandler(discardWriter{}, nil))
	server := grpc.NewServer()

	conn, err := startPipeServer(server, log)
	if err != nil {
		t.Fatalf("startPipeServer(): %v", err)
	}

	if conn.ClientConnInterface == nil {
		t.Fatal("expected grpc client connection")
	}

	if err := conn.Close(); err != nil {
		t.Fatalf("Close(): %v", err)
	}
}

func TestDrainHandlers(t *testing.T) {
	t.Parallel()

	t.Run("unary interceptor", func(t *testing.T) {
		t.Parallel()

		d := &drainHandlers{}

		_, err := d.UnaryServerInterceptor(
			context.Background(),
			nil,
			nil,
			func(ctx context.Context, req any) (any, error) {
				return "ok", nil
			},
		)
		if err != nil {
			t.Fatalf("UnaryServerInterceptor(): %v", err)
		}

		d.Wait()
	})

	t.Run("stream interceptor", func(t *testing.T) {
		t.Parallel()

		d := &drainHandlers{}

		err := d.StreamServerInterceptor(
			nil,
			nil,
			nil,
			func(any, grpc.ServerStream) error {
				return nil
			},
		)
		if err != nil {
			t.Fatalf("StreamServerInterceptor(): %v", err)
		}

		d.Wait()
	})
}

func TestCloserGroup(t *testing.T) {
	t.Parallel()

	t.Run("close empty group is safe", func(t *testing.T) {
		t.Parallel()

		var c closerGroup
		if err := c.Close(); err != nil && !errors.Is(err, io.EOF) {
			t.Fatalf("Close(): %v", err)
		}
	})
}
