package extauthzplugin

import (
	"context"

	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"

	authzpluginv1 "github.com/openkcm/plugin-sdk/proto/kms/plugin/extauthz/v1"
)

var Handshake = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "AUTHZ_PLUGIN",
	MagicCookieValue: "authz",
}

type AuthZ interface {
	Check(CheckRequest) (CheckResponse, error)
}

type CheckRequest struct {
	Subject string
	Object  string
	Action  string
}

type CheckResponse struct {
	Allowed bool
	Message string
}

// This is the implementation of plugin.GRPCPlugin so we can serve/consume this.
type AuthZgRPCPlugin struct {
	// AuthZgRPCPlugin must implement the Plugin interface
	plugin.Plugin
	Impl AuthZ
}

func (p *AuthZgRPCPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	authzpluginv1.RegisterExternalAuthZPluginServer(s, &GRPCServer{Impl: p.Impl})
	return nil
}

func (p *AuthZgRPCPlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &GRPCClient{client: authzpluginv1.NewExternalAuthZPluginClient(c)}, nil
}
