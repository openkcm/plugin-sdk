package keystore_management

import (
	"context"
	"fmt"

	"google.golang.org/protobuf/types/known/structpb"

	"github.com/openkcm/plugin-sdk/api"
	"github.com/openkcm/plugin-sdk/api/service/common"
	"github.com/openkcm/plugin-sdk/api/service/keystoremanagement"
	"github.com/openkcm/plugin-sdk/pkg/plugin"
	grpccommonv1 "github.com/openkcm/plugin-sdk/proto/plugin/common/v1"
	grpckeystoremanagementv1 "github.com/openkcm/plugin-sdk/proto/plugin/keystore_management/v1"
)

type V1 struct {
	plugin.Facade
	grpckeystoremanagementv1.KeystoreManagementPluginClient
}

func (v1 *V1) Version() uint {
	return 1
}

func (v1 *V1) ServiceInfo() api.Info {
	return v1.Info
}

func (v1 *V1) CreateKeystore(ctx context.Context, req *keystoremanagement.CreateKeystoreRequest) (*keystoremanagement.CreateKeystoreResponse, error) {
	value, err := structpb.NewStruct(req.Values)
	if err != nil {
		return nil, fmt.Errorf("failed to parse values: %v", err)
	}

	in := &grpckeystoremanagementv1.CreateKeystoreRequest{
		Values: value,
	}
	grpcResp, err := v1.KeystoreManagementPluginClient.CreateKeystore(ctx, in)
	if err != nil {
		return nil, err
	}
	resp := &keystoremanagement.CreateKeystoreResponse{
		Config: common.InstanceConfig{
			Values: nil,
		},
	}
	if grpcResp.GetConfig() != nil || grpcResp.GetConfig().GetValues() != nil {
		resp.Config.Values = grpcResp.GetConfig().GetValues().AsMap()
	}
	return resp, nil
}

func (v1 *V1) DeleteKeystore(ctx context.Context, req *keystoremanagement.DeleteKeystoreRequest) (*keystoremanagement.DeleteKeystoreResponse, error) {
	value, err := structpb.NewStruct(req.Config.Values)
	if err != nil {
		return nil, fmt.Errorf("failed to parse values: %v", err)
	}
	in := &grpckeystoremanagementv1.DeleteKeystoreRequest{
		Config: &grpccommonv1.KeystoreInstanceConfig{
			Values: value,
		},
	}
	_, err = v1.KeystoreManagementPluginClient.DeleteKeystore(ctx, in)
	if err != nil {
		return nil, err
	}
	return &keystoremanagement.DeleteKeystoreResponse{}, nil
}
