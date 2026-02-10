package keystore_management

import (
	"context"
	"fmt"

	"google.golang.org/protobuf/types/known/structpb"

	"github.com/openkcm/plugin-sdk/api"
	"github.com/openkcm/plugin-sdk/api/service/keystore"
	"github.com/openkcm/plugin-sdk/pkg/plugin"
	commonv1 "github.com/openkcm/plugin-sdk/proto/plugin/keystore/common/v1"
	managementv1 "github.com/openkcm/plugin-sdk/proto/plugin/keystore/management/v1"
)

type V1 struct {
	plugin.Facade
	managementv1.KeystoreProviderPluginClient
}

func (v1 *V1) Version() uint {
	return 1
}

func (v1 *V1) ServiceInfo() api.Info {
	return v1.Info
}

func (v1 *V1) CreateKeystore(ctx context.Context, req *keystore.CreateKeystoreRequest) (*keystore.CreateKeystoreResponse, error) {
	value, err := structpb.NewStruct(req.Values)
	if err != nil {
		return nil, fmt.Errorf("failed to parse values: %v", err)
	}

	in := &managementv1.CreateKeystoreRequest{
		Values: value,
	}
	grpcResp, err := v1.KeystoreProviderPluginClient.CreateKeystore(ctx, in)
	if err != nil {
		return nil, err
	}
	resp := &keystore.CreateKeystoreResponse{
		Config: keystore.InstanceConfig{
			Values: nil,
		},
	}
	if grpcResp.GetConfig() != nil || grpcResp.GetConfig().GetValues() != nil {
		resp.Config.Values = grpcResp.GetConfig().GetValues().AsMap()
	}
	return resp, nil
}

func (v1 *V1) DeleteKeystore(ctx context.Context, req *keystore.DeleteKeystoreRequest) (*keystore.DeleteKeystoreResponse, error) {
	value, err := structpb.NewStruct(req.Config.Values)
	if err != nil {
		return nil, fmt.Errorf("failed to parse values: %v", err)
	}
	in := &managementv1.DeleteKeystoreRequest{
		Config: &commonv1.KeystoreInstanceConfig{
			Values: value,
		},
	}
	_, err = v1.KeystoreProviderPluginClient.DeleteKeystore(ctx, in)
	if err != nil {
		return nil, err
	}
	return &keystore.DeleteKeystoreResponse{}, nil
}
