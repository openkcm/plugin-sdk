package system_information

import (
	"context"

	"buf.build/go/protovalidate"

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
	var grpcType string

	switch req.Type {
	case systeminformation.SystemType:
		grpcType = "system"
	case systeminformation.SubaccountType:
		grpcType = "subaccount"
	case systeminformation.AccountType:
		grpcType = "account"
	default:
		grpcType = "system"
	}

	in := &systeminformationv1.GetRequest{
		Id:   req.ID,
		Type: grpcType,
	}
	if err := protovalidate.Validate(in); err != nil {
		return nil, err
	}

	grpcResp, err := v1.Get(ctx, in)
	if err != nil {
		return nil, err
	}
	return &systeminformation.GetSystemInfoResponse{
		Metadata: grpcResp.GetMetadata(),
	}, nil
}
