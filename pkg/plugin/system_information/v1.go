package system_information

import (
	"context"

	"github.com/openkcm/plugin-sdk/api"
	"github.com/openkcm/plugin-sdk/api/service/systeminformation"
	"github.com/openkcm/plugin-sdk/pkg/plugin"
	systeminformationv1 "github.com/openkcm/plugin-sdk/proto/plugin/systeminformation/v1"
)

type V1 struct {
	plugin.Facade
	systeminformationv1.SystemInformationServicePluginClient
}

func (v1 *V1) Version() uint {
	return 1
}

func (v1 *V1) ServiceInfo() api.Info {
	return v1.Info
}

func (v1 *V1) GetSystemInfo(ctx context.Context, req *systeminformation.GetSystemInfoRequest) (*systeminformation.GetSystemInfoResponse, error) {
	// Convert your API request type to the gRPC enum
	var grpcType systeminformationv1.RequestType

	switch req.Type {
	case systeminformation.UnspecifiedType:
		grpcType = systeminformationv1.RequestType_REQUEST_TYPE_SYSTEM
	case systeminformation.SubaccountType:
		grpcType = systeminformationv1.RequestType_REQUEST_TYPE_SUBACCOUNT
	default:
		grpcType = systeminformationv1.RequestType_REQUEST_TYPE_UNSPECIFIED
	}

	in := &systeminformationv1.GetRequest{
		Id:   req.ID,
		Type: grpcType,
	}
	grpcResp, err := v1.Get(ctx, in)
	if err != nil {
		return nil, err
	}
	return &systeminformation.GetSystemInfoResponse{
		Metadata: grpcResp.GetMetadata(),
	}, nil
}
