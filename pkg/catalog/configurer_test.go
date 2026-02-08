package catalog

import "testing"

func TestMakeConfigurer(t *testing.T) {
	// Arrange
	plugin := &pluginStruct{}
	pluginConfig := PluginConfig{Name: "test"}

	// Act
	c := makeConfigurer(plugin, pluginConfig)

	// Assert
	if c.plugin != plugin {
		t.Errorf("plugin not set")
	}
	if c.pluginConfig.Name != pluginConfig.Name {
		t.Errorf("pluginConfig not set")
	}
}
