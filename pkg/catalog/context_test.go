package catalog

import (
	"context"
	"testing"
)

func TestWithPluginName(t *testing.T) {
	// Act
	got := WithPluginName(context.Background(), "test")

	// Assert
	if got == nil {
		t.Errorf("expected a context.Context, got nil")
	}
	if got.Value(pluginNameKey{}) != "test" {
		t.Errorf("expected value to be 'test', got %v", got.Value(pluginNameKey{}))
	}
}
