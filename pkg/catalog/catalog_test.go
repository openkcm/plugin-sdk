package catalog

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/openkcm/plugin-sdk/api"
	"github.com/openkcm/plugin-sdk/pkg/plugin"
	testv1 "github.com/openkcm/plugin-sdk/proto/plugin/test/v1"
)

type testRepo struct {
	client *testServiceFacade
}

func (r *testRepo) Plugins() map[string]api.PluginRepo {
	return map[string]api.PluginRepo{testv1.Type: r}
}
func (r *testRepo) Services() []api.ServiceRepo  { return nil }
func (r *testRepo) Binder() any                  { return func(c *testServiceFacade) { r.client = c } }
func (r *testRepo) Versions() []api.Version      { return []api.Version{testServiceV1{}} }
func (r *testRepo) Constraints() api.Constraints { return api.ExactlyOne() }
func (r *testRepo) Clear()                       { r.client = nil }

type testServiceV1 struct{}

func (testServiceV1) New() api.Facade  { return new(testServiceFacade) }
func (testServiceV1) Deprecated() bool { return false }

type testServiceFacade struct {
	plugin.Facade
	testv1.TestServicePluginClient
}

func TestCatalog_WithErrorLogging(t *testing.T) {
	repo := &testRepo{}
	var buf bytes.Buffer

	cat, err := New(context.Background(), Config{
		Logger: newTestLogger(&buf),
		PluginConfigs: []PluginConfig{{
			Name: "testplugin",
			Type: testv1.Type,
			Path: "./testpluginbinary",
		}},
		PluginOptions: []PluginOption{WithErrorLogging(newTestLogger(&buf))},
	}, repo)
	require.NoError(t, err)
	defer cat.Close()

	_, err = repo.client.Test(context.Background(), &testv1.TestRequest{})
	require.NoError(t, err)

	assert.NotContains(t, buf.String(), "plugin method returned error")
}
