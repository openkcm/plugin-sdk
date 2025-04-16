package bootstrap

import "testing"

func TestNewLogger(t *testing.T) {
	// Act
	got := NewLogger()

	// Assert
	if got == nil {
		t.Errorf("Expected logger, got nil")
	}
}
