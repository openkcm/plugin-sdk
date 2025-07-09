package catalog

import (
	"context"
	"log/slog"
	"testing"
)

func TestLoad(t *testing.T) {
	// create test cases
	tests := []struct {
		name      string
		config    Config
		wantError bool
	}{
		{
			name: "zero values",
		}, {
			name: "load disabled",
			config: Config{
				Logger: slog.Default(),
				PluginConfigs: []PluginConfig{
					{
						Logger:   slog.Default(),
						Disabled: true,
					},
				},
				HostServices: nil,
			},
			wantError: false,
		}, {
			name: "invalid path",
			config: Config{
				Logger: slog.Default(),
				PluginConfigs: []PluginConfig{
					{
						Path:   "/does/not/exist",
						Logger: slog.Default(),
					},
				},
				HostServices: nil,
			},
			wantError: true,
		}, {
			name: "load testplugin",
			config: Config{
				Logger: slog.Default(),
				PluginConfigs: []PluginConfig{
					{
						// testpluginbinary is built in the TestMain function
						Path:   "./testpluginbinary",
						Type:   "TestService",
						Logger: slog.Default(),
					},
				},
				HostServices: nil,
			},
			wantError: false,
		},
		{
			name: "load testplugin feature",
			config: Config{
				Logger: slog.Default(),
				PluginConfigs: []PluginConfig{
					{
						// testpluginbinary is built in the TestMain function
						Path:        "./testpluginbinary",
						Type:        "TestService",
						Logger:      slog.Default(),
						HYOKEnabled: true,
					},
				},
				HostServices: nil,
			},
			wantError: false,
		},
	}

	// run the tests
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			catalog, err := Load(context.Background(), tc.config)
			defer func(c *Catalog) {
				if c != nil {
					c.Close()
				}
			}(catalog)

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

func TestLookupByType(t *testing.T) {
	// Arrange
	cfg := Config{
		Logger: slog.Default(),
		PluginConfigs: []PluginConfig{
			{
				// testpluginbinary is built in the TestMain function
				Path:   "./testpluginbinary",
				Type:   "TestService",
				Logger: slog.Default(),
			},
		},
		HostServices: nil,
	}
	catalog, err := Load(context.Background(), cfg)
	if err != nil {
		t.Fatalf("failed to load plugin: %v", err)
	}
	defer catalog.Close()

	// create test cases
	tests := []struct {
		name    string
		pName   string
		wantNil bool
	}{
		{
			name:    "unknown plugin",
			pName:   "foobar",
			wantNil: true,
		}, {
			name:    "known plugin",
			pName:   "TestService",
			wantNil: false,
		},
	}

	// run the tests
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			got := catalog.LookupByType(tc.pName)

			// Assert
			if tc.wantNil && got == nil { // expected nil and got it
				return
			} else if tc.wantNil && got != nil { // expected nil but did not get it
				t.Errorf("expected value: %v, got: %v", tc.wantNil, got)
			} else if !tc.wantNil && got == nil { // got unexpected nil
				t.Errorf("expected value: %v, got: %v", tc.wantNil, got)
			} else if !tc.wantNil && got != nil {
				return
			}
		})
	}
}
func TestLookupByTypeAndName(t *testing.T) {
	// Arrange
	cfg := Config{
		Logger: slog.Default(),
		PluginConfigs: []PluginConfig{
			{
				// testpluginbinary is built in the TestMain function
				Path:   "./testpluginbinary",
				Type:   "TestService",
				Name:   "TestPlugin",
				Logger: slog.Default(),
			},
		},
		HostServices: nil,
	}
	catalog, err := Load(context.Background(), cfg)
	if err != nil {
		t.Fatalf("failed to load plugin: %v", err)
	}
	defer catalog.Close()

	// create test cases
	tests := []struct {
		name    string
		pType   string
		pName   string
		wantNil bool
	}{
		{
			name:    "unknown plugin",
			pType:   "UnknownType",
			pName:   "UnknownPlugin",
			wantNil: true,
		}, {
			name:    "known plugin",
			pType:   "TestService",
			pName:   "TestPlugin",
			wantNil: false,
		},
	}

	// run the tests
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			got := catalog.LookupByTypeAndName(tc.pType, tc.pName)

			// Assert
			if tc.wantNil && got == nil { // expected nil and got it
				return
			} else if tc.wantNil && got != nil { // expected nil but did not get it
				t.Errorf("expected value: %v, got: %v", tc.wantNil, got)
			} else if !tc.wantNil && got == nil { // got unexpected nil
				t.Errorf("expected value: %v, got: %v", tc.wantNil, got)
			} else if !tc.wantNil && got != nil {
				return
			}
		})
	}
}

func TestPluginInfo_Features(t *testing.T) {
	tests := []struct {
		name         string
		pluginCfg    PluginConfig
		wantFeatures []string
	}{
		{
			name: "HYOK enabled",
			pluginCfg: PluginConfig{
				Path:        "./testpluginbinary",
				Type:        "TestService",
				Name:        "TestPluginHYOK",
				Logger:      slog.Default(),
				HYOKEnabled: true,
			},
			wantFeatures: []string{"HYOK"},
		},
		{
			name: "HYOK disabled",
			pluginCfg: PluginConfig{
				Path:        "./testpluginbinary",
				Type:        "TestService",
				Name:        "TestPluginNoHYOK",
				Logger:      slog.Default(),
				HYOKEnabled: false,
			},
			wantFeatures: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cfg := Config{
				Logger:        slog.Default(),
				PluginConfigs: []PluginConfig{tc.pluginCfg},
				HostServices:  nil,
			}
			catalog, err := Load(context.Background(), cfg)
			if err != nil {
				t.Fatalf("failed to load plugin: %v", err)
			}
			defer catalog.Close()
			plugin := catalog.LookupByTypeAndName(tc.pluginCfg.Type, tc.pluginCfg.Name)
			if plugin == nil {
				t.Fatalf("plugin not found")
			}
			features := plugin.Info().Features()
			if len(features) != len(tc.wantFeatures) {
				t.Errorf("expected features %v, got %v", tc.wantFeatures, features)
				return
			}
			for i, f := range tc.wantFeatures {
				if features[i] != f {
					t.Errorf("expected features %v, got %v", tc.wantFeatures, features)
					break
				}
			}
		})
	}
}
