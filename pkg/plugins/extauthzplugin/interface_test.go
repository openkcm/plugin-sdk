package extauthzplugin

import (
	"testing"

	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
)

func TestGRPCServer(t *testing.T) {
	broker := &plugin.GRPCBroker{}
	server := grpc.NewServer()
	plugin := AuthZgRPCPlugin{}
	if err := plugin.GRPCServer(broker, server); err != nil {
		t.Errorf("Error: %v", err)
	}
}

func TestGRPCClient(t *testing.T) {
	plugin := AuthZgRPCPlugin{}
	_, err := plugin.GRPCClient(nil, nil, nil)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
}
