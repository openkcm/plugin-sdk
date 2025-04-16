// This package implements a test plugin for go test.
package main

import (
	"context"

	"github.com/openkcm/plugin-sdk/pkg/plugin"
	testv1 "github.com/openkcm/plugin-sdk/proto/plugin/test/v1"
	configv1 "github.com/openkcm/plugin-sdk/proto/service/common/config/v1"
)

type TestPlugin struct {
	testv1.UnsafeTestServiceServer
	configv1.UnsafeConfigServer
}

func (p *TestPlugin) Test(ctx context.Context, req *testv1.TestRequest) (*testv1.TestResponse, error) {
	return &testv1.TestResponse{Response: "test"}, nil
}

func (p *TestPlugin) Configure(ctx context.Context, req *configv1.ConfigureRequest) (*configv1.ConfigureResponse, error) {
	return &configv1.ConfigureResponse{}, nil
}

// main() serves the plugin. Serve() will not return. If there is a
// failure, the process will exit with a non-zero exit code.
func main() {
	p := &TestPlugin{}
	plugin.Serve(
		testv1.TestServicePluginServer(p),
		configv1.ConfigServiceServer(p),
	)
}
