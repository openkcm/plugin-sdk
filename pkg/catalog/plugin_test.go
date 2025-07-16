package catalog

import (
	"os/exec"
	"testing"
)

func TestClientConnection(t *testing.T) {
	// Arrange
	p := Plugin{}

	// Act
	got := p.ClientConnection()

	// Assert
	if got != nil {
		t.Errorf("Expected nil, got %v", got)
	}
}

func TestInfo(t *testing.T) {
	// Arrange
	p := Plugin{info: pluginInfo{name: "test"}}

	// Act
	got := p.Info()

	// Assert
	if got.Name() != "test" {
		t.Errorf("Expected name to be 'test', but got %s", got.Name())
	}
}

func TestGrpcServiceNames(t *testing.T) {
	// Arrange
	p := Plugin{grpcServiceNames: []string{"test"}}

	// Act
	got := p.GrpcServiceNames()

	// Assert
	if len(got) != 1 && got[0] != "test" {
		t.Errorf("Expected name to be 'test', but got %s", got[0])
	}
}

func TestGRPCServer(t *testing.T) {
	// Arrange
	p := HCRPCPlugin{}

	// Act
	err := p.GRPCServer(nil, nil)

	// Assert
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestName(t *testing.T) {
	// Arrange
	pi := pluginInfo{name: "test"}

	// Act
	got := pi.Name()

	// Assert
	if got != "test" {
		t.Errorf("Expected name to be 'test', but got %s", got)
	}
}

func TestType(t *testing.T) {
	// Arrange
	pi := pluginInfo{typ: "test"}

	// Act
	got := pi.Type()

	// Assert
	if got != "test" {
		t.Errorf("Expected name to be 'test', but got %s", got)
	}
}

func TestBuildSecureConfig(t *testing.T) {
	// create test cases
	tests := []struct {
		name      string
		checksum  string
		wantError bool
	}{
		{
			name: "zero values",
		}, {
			name:      "checksum not hex",
			checksum:  "1234567890abcdez",
			wantError: true,
		}, {
			name:      "checksum too short",
			checksum:  "1234567890abcdef1234567890abcdef",
			wantError: true,
		}, {
			name:     "valid checksum",
			checksum: "1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
		},
	}

	// run the tests
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			_, err := buildSecureConfig(tc.checksum)

			// Assert
			if tc.wantError && err != nil { // expected error and got it
				return
			} else if tc.wantError && err == nil { // expected error but did not get it
				t.Errorf("expected error value: %v, got: %s", tc.wantError, err)
			} else if !tc.wantError && err != nil { // got unexpected error
				t.Errorf("expected error value: %v, got: %s", tc.wantError, err)
			} else if !tc.wantError && err == nil {
				return
			}
		})
	}
}

func TestInjectEnv(t *testing.T) {
	// Arrange
	cmd := &exec.Cmd{}
	config := PluginConfig{
		Env: map[string]string{
			"KEY1": "VALUE1",
			"KEY2": "VALUE2",
		},
	}

	// Act
	injectEnv(config, cmd)

	// Assert
	expectedEnv := []string{"KEY1=VALUE1", "KEY2=VALUE2"}
	for _, expected := range expectedEnv {
		found := false
		for _, actual := range cmd.Env {
			if actual == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected environment variable %s not found in cmd.Env", expected)
		}
	}
}
