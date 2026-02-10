package system_information

import (
	"context"

	"github.com/openkcm/plugin-sdk/api/service/systeminformation"
	"github.com/openkcm/plugin-sdk/pkg/plugin"
	systeminformationv1 "github.com/openkcm/plugin-sdk/proto/plugin/systeminformation/v1"
)

type V1 struct {
	plugin.Facade
	systeminformationv1.SystemInformationServicePluginClient
}

func (v1 *V1) GetSystemInfo(ctx context.Context, req *systeminformation.GetSystemInfoRequest) (*systeminformation.GetSystemInfoResponse, error) {
	in := &systeminformationv1.GetRequest{
		Id:   req.ID,
		Type: systeminformationv1.RequestType(req.Type),
	}
	grpcResp, err := v1.Get(ctx, in)
	if err != nil {
		return nil, err
	}
	return &systeminformation.GetSystemInfoResponse{
		Metadata: grpcResp.GetMetadata(),
	}, nil
}
