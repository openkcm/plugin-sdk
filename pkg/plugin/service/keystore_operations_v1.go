package service

import (
	"context"
	"fmt"

	"google.golang.org/protobuf/types/known/structpb"

	"github.com/openkcm/plugin-sdk/api/service/keystore"
	"github.com/openkcm/plugin-sdk/pkg/catalog"
	commonv1 "github.com/openkcm/plugin-sdk/proto/plugin/keystore/common/v1"
	operationsv1 "github.com/openkcm/plugin-sdk/proto/plugin/keystore/operations/v1"
)

var _ keystore.KeystoreOperations = (*hashicorpKeystoreOperationsV1Plugin)(nil)

type hashicorpKeystoreOperationsV1Plugin struct {
	plugin     *catalog.Plugin
	grpcClient operationsv1.KeystoreInstanceKeyOperationClient
}

func NewKeystoreOperationsV1Plugin(plugin *catalog.Plugin) keystore.KeystoreOperations {
	return &hashicorpKeystoreOperationsV1Plugin{
		plugin:     plugin,
		grpcClient: operationsv1.NewKeystoreInstanceKeyOperationClient(plugin.ClientConnection()),
	}
}

func (h *hashicorpKeystoreOperationsV1Plugin) GetKey(ctx context.Context, req *keystore.GetKeyRequest) (*keystore.GetKeyResponse, error) {
	value, err := structpb.NewStruct(req.Parameters.Config.Values)
	if err != nil {
		return nil, fmt.Errorf("failed to parse values: %v", err)
	}

	in := &operationsv1.GetKeyRequest{
		Parameters: &operationsv1.RequestParameters{
			Config: &commonv1.KeystoreInstanceConfig{
				Values: value,
			},
			KeyId: req.Parameters.KeyID,
		},
	}
	grpcResp, err := h.grpcClient.GetKey(ctx, in)
	if err != nil {
		return nil, err
	}

	return &keystore.GetKeyResponse{
		KeyID:        grpcResp.GetKeyId(),
		KeyAlgorithm: keystore.KeyAlgorithm(grpcResp.GetAlgorithm()),
		Status:       grpcResp.GetStatus(),
		Usage:        grpcResp.GetUsage(),
	}, nil
}

func (h *hashicorpKeystoreOperationsV1Plugin) CreateKey(ctx context.Context, req *keystore.CreateKeyRequest) (*keystore.CreateKeyResponse, error) {
	value, err := structpb.NewStruct(req.Config.Values)
	if err != nil {
		return nil, fmt.Errorf("failed to parse values: %v", err)
	}

	in := &operationsv1.CreateKeyRequest{
		Config: &commonv1.KeystoreInstanceConfig{
			Values: value,
		},
		Algorithm: operationsv1.KeyAlgorithm(req.KeyAlgorithm),
		Id:        req.ID,
		Region:    req.Region,
		KeyType:   operationsv1.KeyType(req.KeyType),
	}
	grpcResp, err := h.grpcClient.CreateKey(ctx, in)
	if err != nil {
		return nil, err
	}

	return &keystore.CreateKeyResponse{
		KeyID:  grpcResp.GetKeyId(),
		Status: grpcResp.GetStatus(),
	}, nil
}

func (h *hashicorpKeystoreOperationsV1Plugin) DeleteKey(ctx context.Context, req *keystore.DeleteKeyRequest) (*keystore.DeleteKeyResponse, error) {
	value, err := structpb.NewStruct(req.Parameters.Config.Values)
	if err != nil {
		return nil, fmt.Errorf("failed to parse values: %v", err)
	}

	in := &operationsv1.DeleteKeyRequest{
		Parameters: &operationsv1.RequestParameters{
			Config: &commonv1.KeystoreInstanceConfig{
				Values: value,
			},
			KeyId: req.Parameters.KeyID,
		},
		Window: req.Window,
	}
	_, err = h.grpcClient.DeleteKey(ctx, in)
	if err != nil {
		return nil, err
	}

	return &keystore.DeleteKeyResponse{}, nil
}

func (h *hashicorpKeystoreOperationsV1Plugin) EnableKey(ctx context.Context, req *keystore.EnableKeyRequest) (*keystore.EnableKeyResponse, error) {
	value, err := structpb.NewStruct(req.Parameters.Config.Values)
	if err != nil {
		return nil, fmt.Errorf("failed to parse values: %v", err)
	}

	in := &operationsv1.EnableKeyRequest{
		Parameters: &operationsv1.RequestParameters{
			Config: &commonv1.KeystoreInstanceConfig{
				Values: value,
			},
			KeyId: req.Parameters.KeyID,
		},
	}
	_, err = h.grpcClient.EnableKey(ctx, in)
	if err != nil {
		return nil, err
	}

	return &keystore.EnableKeyResponse{}, nil
}

func (h *hashicorpKeystoreOperationsV1Plugin) GetImportParameters(ctx context.Context, req *keystore.GetImportParametersRequest) (*keystore.GetImportParametersResponse, error) {
	value, err := structpb.NewStruct(req.Parameters.Config.Values)
	if err != nil {
		return nil, fmt.Errorf("failed to parse values: %v", err)
	}

	in := &operationsv1.GetImportParametersRequest{
		Parameters: &operationsv1.RequestParameters{
			Config: &commonv1.KeystoreInstanceConfig{
				Values: value,
			},
			KeyId: req.Parameters.KeyID,
		},
		Algorithm: operationsv1.KeyAlgorithm(req.KeyAlgorithm),
	}
	grpcResp, err := h.grpcClient.GetImportParameters(ctx, in)
	if err != nil {
		return nil, err
	}

	return &keystore.GetImportParametersResponse{
		KeyID:            grpcResp.GetKeyId(),
		ImportParameters: grpcResp.GetImportParameters().AsMap(),
	}, nil
}

func (h *hashicorpKeystoreOperationsV1Plugin) ImportKeyMaterial(ctx context.Context, req *keystore.ImportKeyMaterialRequest) (*keystore.ImportKeyMaterialResponse, error) {
	value, err := structpb.NewStruct(req.Parameters.Config.Values)
	if err != nil {
		return nil, fmt.Errorf("failed to parse values: %v", err)
	}

	importParams, err := structpb.NewStruct(req.ImportParameters)
	if err != nil {
		return nil, fmt.Errorf("failed to parse values: %v", err)
	}

	in := &operationsv1.ImportKeyMaterialRequest{
		Parameters: &operationsv1.RequestParameters{
			Config: &commonv1.KeystoreInstanceConfig{
				Values: value,
			},
			KeyId: req.Parameters.KeyID,
		},
		ImportParameters: importParams,
	}
	_, err = h.grpcClient.ImportKeyMaterial(ctx, in)
	if err != nil {
		return nil, err
	}

	return &keystore.ImportKeyMaterialResponse{}, nil
}

func (h *hashicorpKeystoreOperationsV1Plugin) ValidateKey(ctx context.Context, req *keystore.ValidateKeyRequest) (*keystore.ValidateKeyResponse, error) {
	in := &operationsv1.ValidateKeyRequest{
		KeyType:     operationsv1.KeyType(req.KeyType),
		Algorithm:   operationsv1.KeyAlgorithm(req.KeyAlgorithm),
		Region:      req.Region,
		NativeKeyId: req.NativeKeyID,
	}
	grpcResp, err := h.grpcClient.ValidateKey(ctx, in)
	if err != nil {
		return nil, err
	}

	return &keystore.ValidateKeyResponse{
		IsValid: grpcResp.GetIsValid(),
		Message: grpcResp.GetMessage(),
	}, nil
}

func (h *hashicorpKeystoreOperationsV1Plugin) ValidateKeyAccessData(ctx context.Context, req *keystore.ValidateKeyAccessDataRequest) (*keystore.ValidateKeyAccessDataResponse, error) {
	management, err := structpb.NewStruct(req.Management)
	if err != nil {
		return nil, fmt.Errorf("failed to parse values: %v", err)
	}
	crypto, err := structpb.NewStruct(req.Crypto)
	if err != nil {
		return nil, fmt.Errorf("failed to parse values: %v", err)
	}

	in := &operationsv1.ValidateKeyAccessDataRequest{
		Management: management,
		Crypto:     crypto,
	}
	grpcResp, err := h.grpcClient.ValidateKeyAccessData(ctx, in)
	if err != nil {
		return nil, err
	}

	return &keystore.ValidateKeyAccessDataResponse{
		IsValid: grpcResp.GetIsValid(),
		Message: grpcResp.GetMessage(),
	}, nil
}

func (h *hashicorpKeystoreOperationsV1Plugin) TransformCryptoAccessData(ctx context.Context, req *keystore.TransformCryptoAccessDataRequest) (*keystore.TransformCryptoAccessDataResponse, error) {
	in := &operationsv1.TransformCryptoAccessDataRequest{
		NativeKeyId: req.NativeKeyID,
		AccessData:  req.AccessData,
	}
	grpcResp, err := h.grpcClient.TransformCryptoAccessData(ctx, in)
	if err != nil {
		return nil, err
	}

	return &keystore.TransformCryptoAccessDataResponse{
		TransformedAccessData: grpcResp.GetTransformedAccessData(),
	}, nil
}

func (h *hashicorpKeystoreOperationsV1Plugin) ExtractKeyRegion(ctx context.Context, req *keystore.ExtractKeyRegionRequest) (*keystore.ExtractKeyRegionResponse, error) {
	management, err := structpb.NewStruct(req.ManagementAccessData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse values: %v", err)
	}

	in := &operationsv1.ExtractKeyRegionRequest{
		NativeKeyId:          req.NativeKeyID,
		ManagementAccessData: management,
	}
	grpcResp, err := h.grpcClient.ExtractKeyRegion(ctx, in)
	if err != nil {
		return nil, err
	}

	return &keystore.ExtractKeyRegionResponse{
		Region: grpcResp.GetRegion(),
	}, nil
}
