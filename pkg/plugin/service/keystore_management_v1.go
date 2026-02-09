package service

import (
	"context"
	"fmt"

	"google.golang.org/protobuf/types/known/structpb"

	"github.com/openkcm/plugin-sdk/api/service/keystore"
	"github.com/openkcm/plugin-sdk/pkg/catalog"
	commonv1 "github.com/openkcm/plugin-sdk/proto/plugin/keystore/common/v1"
	managementv1 "github.com/openkcm/plugin-sdk/proto/plugin/keystore/management/v1"
)

var _ keystore.KeystoreManagement = (*hashicorpKeystoreManagementV1Plugin)(nil)

type hashicorpKeystoreManagementV1Plugin struct {
	plugin     catalog.Plugin
	grpcClient managementv1.KeystoreProviderClient
}

func NewKeystoreManagementV1Plugin(plugin catalog.Plugin) keystore.KeystoreManagement {
	return &hashicorpKeystoreManagementV1Plugin{
		plugin:     plugin,
		grpcClient: managementv1.NewKeystoreProviderClient(plugin.ClientConnection()),
	}
}

func (h *hashicorpKeystoreManagementV1Plugin) CreateKeystore(ctx context.Context, req *keystore.CreateKeystoreRequest) (*keystore.CreateKeystoreResponse, error) {
	value, err := structpb.NewStruct(req.Values)
	if err != nil {
		return nil, fmt.Errorf("failed to parse values: %v", err)
	}

	in := &managementv1.CreateKeystoreRequest{
		Values: value,
	}
	grpcResp, err := h.grpcClient.CreateKeystore(ctx, in)
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

func (h *hashicorpKeystoreManagementV1Plugin) DeleteKeystore(ctx context.Context, req *keystore.DeleteKeystoreRequest) (*keystore.DeleteKeystoreResponse, error) {
	value, err := structpb.NewStruct(req.Config.Values)
	if err != nil {
		return nil, fmt.Errorf("failed to parse values: %v", err)
	}
	in := &managementv1.DeleteKeystoreRequest{
		Config: &commonv1.KeystoreInstanceConfig{
			Values: value,
		},
	}
	_, err = h.grpcClient.DeleteKeystore(ctx, in)
	if err != nil {
		return nil, err
	}
	return &keystore.DeleteKeystoreResponse{}, nil
}
