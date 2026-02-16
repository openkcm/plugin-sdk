package system_information

import (
	"context"
	"fmt"

	"buf.build/go/protovalidate"

	"github.com/openkcm/plugin-sdk/api"
	"github.com/openkcm/plugin-sdk/pkg/plugin"
	grpcsysteminformationv1 "github.com/openkcm/plugin-sdk/proto/plugin/system_information/v1"
	"github.com/openkcm/plugin-sdk/service/api/systeminformation"
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
		Id:   req.ID,
		Type: req.Type,
	}
	if err := protovalidate.Validate(in); err != nil {
		return nil, fmt.Errorf("failed validation: %v", err)
	}

	grpcResp, err := v1.GetInfo(ctx, in)
	if err != nil {
		return nil, err
	}
	return &systeminformation.GetSystemInfoResponse{
		Metadata: grpcResp.GetMetadata(),
	}, nil
}
