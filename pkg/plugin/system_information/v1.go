package system_information

import (
	"context"

	"github.com/openkcm/plugin-sdk/api"
	"github.com/openkcm/plugin-sdk/api/service/systeminformation"
	"github.com/openkcm/plugin-sdk/pkg/plugin"
	grpcsysteminformationv1 "github.com/openkcm/plugin-sdk/proto/plugin/systeminformation/v1"
)

type V1 struct {
	plugin.Facade
	grpcsysteminformationv1.SystemInformationServicePluginClient
}

func (v1 *V1) Version() uint {
	return 1
}

func (v1 *V1) ServiceInfo() api.Info {
	return v1.Info
}

func (v1 *V1) GetSystemInfo(ctx context.Context, req *systeminformation.GetSystemInfoRequest) (*systeminformation.GetSystemInfoResponse, error) {
	in := &grpcsysteminformationv1.GetRequest{
		Id:   req.ID,
		Type: toGRPCType(req.Type),
	}
	grpcResp, err := v1.Get(ctx, in)
	if err != nil {
		return nil, err
	}
	return &systeminformation.GetSystemInfoResponse{
		Metadata: grpcResp.GetMetadata(),
	}, nil
}

func toGRPCType(t systeminformation.Type) grpcsysteminformationv1.RequestType {
	switch t {
	case systeminformation.SystemType:
		return grpcsysteminformationv1.RequestType_REQUEST_TYPE_SYSTEM
	case systeminformation.SubaccountType:
		return grpcsysteminformationv1.RequestType_REQUEST_TYPE_SUBACCOUNT
	default:
		return grpcsysteminformationv1.RequestType_REQUEST_TYPE_UNSPECIFIED
	}
}
