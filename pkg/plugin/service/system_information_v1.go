package service

import (
	"context"

	"github.com/openkcm/plugin-sdk/api/service/systeminformation"
	"github.com/openkcm/plugin-sdk/pkg/catalog"
	systeminformationv1 "github.com/openkcm/plugin-sdk/proto/plugin/systeminformation/v1"
)

var _ systeminformation.SystemInformation = (*hashicorpSystemInformationV1Plugin)(nil)

type hashicorpSystemInformationV1Plugin struct {
	plugin     catalog.Plugin
	grpcClient systeminformationv1.SystemInformationServiceClient
}

func NewSystemInformationV1Plugin(plugin catalog.Plugin) systeminformation.SystemInformation {
	return &hashicorpSystemInformationV1Plugin{
		plugin:     plugin,
		grpcClient: systeminformationv1.NewSystemInformationServiceClient(plugin.ClientConnection()),
	}
}

func (h *hashicorpSystemInformationV1Plugin) Get(ctx context.Context, req *systeminformation.GetSystemInformationRequest) (*systeminformation.GetSystemInformationResponse, error) {
	in := &systeminformationv1.GetRequest{
		Id:   req.ID,
		Type: systeminformationv1.RequestType(req.Type),
	}
	grpcResp, err := h.grpcClient.Get(ctx, in)
	if err != nil {
		return nil, err
	}
	return &systeminformation.GetSystemInformationResponse{
		Metadata: grpcResp.GetMetadata(),
	}, nil
}
