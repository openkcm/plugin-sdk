package system_information

import (
	"context"

	"github.com/openkcm/plugin-sdk/api"
	"github.com/openkcm/plugin-sdk/api/service/systeminformation"
	"github.com/openkcm/plugin-sdk/pkg/plugin"
	grpcsysteminformationv1 "github.com/openkcm/plugin-sdk/proto/plugin/system_information/v1"
)

type V1 struct {
	plugin.Facade
	grpcsysteminformationv1.SystemInformationPluginClient
}

func (v1 *V1) Version() uint {
	return 1
}

func (v1 *V1) ServiceInfo() api.Info {
	return v1.Info
}

func (v1 *V1) GetSystemInfo(ctx context.Context, req *systeminformation.GetSystemInfoRequest) (*systeminformation.GetSystemInfoResponse, error) {
	in := &grpcsysteminformationv1.GetInfoRequest{
		Id: req.ID,
	}
	setTypeValue(in, req.Type)

	grpcResp, err := v1.GetInfo(ctx, in)
	if err != nil {
		return nil, err
	}
	return &systeminformation.GetSystemInfoResponse{
		Metadata: grpcResp.GetMetadata(),
	}, nil
}

func setTypeValue(req *grpcsysteminformationv1.GetInfoRequest, requestType systeminformation.Type) {
	switch requestType {
	case systeminformation.UnspecifiedType:
		req.TypeValue = &grpcsysteminformationv1.GetInfoRequest_Unspecified{}
	case systeminformation.SystemType:
		req.TypeValue = &grpcsysteminformationv1.GetInfoRequest_System{}
	case systeminformation.SubaccountType:
		req.TypeValue = &grpcsysteminformationv1.GetInfoRequest_Subaccount{}
	case systeminformation.AccountType:
		req.TypeValue = &grpcsysteminformationv1.GetInfoRequest_Account{}
	}

	req.TypeValue = &grpcsysteminformationv1.GetInfoRequest_Unspecified{}
}
