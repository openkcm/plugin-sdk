package catalog

import (
	"context"
	"log/slog"
	"testing"
)

func TestNewHostServer(t *testing.T) {
	// Act
	got := newHostServer(nil, "test")

	// Assert
	if got == nil {
		t.Errorf("expected a *grpc.Server, got nil")
	}
}

func TestConvertPanic(t *testing.T) {
	// Act
	err := convertPanic(slog.Default(), "test")

	// Assert
	if err == nil {
		t.Errorf("expected an error, got nil")
	}
}

func TestContext(t *testing.T) {
	// Arrange
	sw := streamWrapper{ctx: context.Background()}

	// Act
	got := sw.Context()

	// Assert
	if got == nil {
		t.Errorf("expected a context.Context, got nil")
	}
}
