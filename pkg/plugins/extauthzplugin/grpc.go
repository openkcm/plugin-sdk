package extauthzplugin

import (
	"context"

	authzpluginv1 "github.com/openkcm/plugin-sdk/proto/kms/plugin/extauthz/v1"
)

// GRPCClient is an implementation of AuthZ that talks over RPC.
type GRPCClient struct {
	client authzpluginv1.ExternalAuthZPluginClient
}

func (c *GRPCClient) Check(req CheckRequest) (CheckResponse, error) {
	resp, err := c.client.Check(context.Background(), &authzpluginv1.CheckRequest{
		Subject: req.Subject,
		Object:  req.Object,
		Action:  req.Action,
	})
	if err != nil {
		return CheckResponse{}, err
	}

	return CheckResponse{
		Allowed: resp.Allowed,
		Message: resp.Message,
	}, nil
}

// GRPCServer is the gRPC server that GRPCClients talk to.
type GRPCServer struct {
	authzpluginv1.UnimplementedExternalAuthZPluginServer
	Impl AuthZ
}

func (s *GRPCServer) Check(ctx context.Context, req *authzpluginv1.CheckRequest) (*authzpluginv1.CheckResponse, error) {
	resp, err := s.Impl.Check(CheckRequest{
		Subject: req.Subject,
		Object:  req.Object,
		Action:  req.Action,
	})
	if err != nil {
		return nil, err
	}

	return &authzpluginv1.CheckResponse{
		Allowed: resp.Allowed,
		Message: resp.Message,
	}, nil
}
