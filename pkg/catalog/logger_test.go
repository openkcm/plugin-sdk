package catalog

import (
	"log/slog"
	"testing"

	"github.com/hashicorp/go-hclog"
)

func TestSetLogLevel(t *testing.T) {
	// create test cases
	tests := []struct {
		name  string
		input string
		want  hclog.Level
	}{
		{
			name: "zero values",
			want: hclog.Info,
		}, {
			name:  "invalid input",
			input: "invalid",
			want:  hclog.Info,
		}, {
			name:  "debug lower case",
			input: "debug",
			want:  hclog.Debug,
		}, {
			name:  "debug upper case",
			input: "DEBUG",
			want:  hclog.Debug,
		}, {
			name:  "debug mixed case",
			input: "DEbug",
			want:  hclog.Debug,
		}, {
			name:  "info lower case",
			input: "info",
			want:  hclog.Info,
		}, {
			name:  "info upper case",
			input: "INFO",
			want:  hclog.Info,
		}, {
			name:  "info mixed case",
			input: "INfo",
			want:  hclog.Info,
		}, {
			name:  "warn lower case",
			input: "warn",
			want:  hclog.Warn,
		}, {
			name:  "warn upper case",
			input: "WARN",
			want:  hclog.Warn,
		}, {
			name:  "warn mixed case",
			input: "WArn",
			want:  hclog.Warn,
		}, {
			name:  "error lower case",
			input: "error",
			want:  hclog.Error,
		}, {
			name:  "error upper case",
			input: "ERROR",
			want:  hclog.Error,
		}, {
			name:  "error mixed case",
			input: "ERRor",
			want:  hclog.Error,
		},
	}

	// run the tests
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			log := newHClogFromSlog(slog.Default(), tc.input)

			// Assert
			if got := log.GetLevel(); got != tc.want {
				t.Errorf("expected  value: %v, got: %s", tc.want, got)
			}
		})
	}
}
