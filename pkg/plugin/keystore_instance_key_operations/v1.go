package keystore_instance_key_operations

import (
	"context"
	"fmt"

	"google.golang.org/protobuf/types/known/structpb"

	"github.com/openkcm/plugin-sdk/api"
	"github.com/openkcm/plugin-sdk/api/service/keystore"
	"github.com/openkcm/plugin-sdk/pkg/plugin"
	commonv1 "github.com/openkcm/plugin-sdk/proto/plugin/keystore/common/v1"
	operationsv1 "github.com/openkcm/plugin-sdk/proto/plugin/keystore/operations/v1"
)

type V1 struct {
	plugin.Facade
	operationsv1.KeystoreInstanceKeyOperationPluginClient
}

func (v1 *V1) Version() uint {
	return 1
}

func (v1 *V1) ServiceInfo() api.Info {
	return v1.Info
}

func (v1 *V1) GetKey(ctx context.Context, req *keystore.GetKeyRequest) (*keystore.GetKeyResponse, error) {
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
	grpcResp, err := v1.KeystoreInstanceKeyOperationPluginClient.GetKey(ctx, in)
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

func (v1 *V1) CreateKey(ctx context.Context, req *keystore.CreateKeyRequest) (*keystore.CreateKeyResponse, error) {
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
	grpcResp, err := v1.KeystoreInstanceKeyOperationPluginClient.CreateKey(ctx, in)
	if err != nil {
		return nil, err
	}

	return &keystore.CreateKeyResponse{
		KeyID:  grpcResp.GetKeyId(),
		Status: grpcResp.GetStatus(),
	}, nil
}

func (v1 *V1) DeleteKey(ctx context.Context, req *keystore.DeleteKeyRequest) (*keystore.DeleteKeyResponse, error) {
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
	_, err = v1.KeystoreInstanceKeyOperationPluginClient.DeleteKey(ctx, in)
	if err != nil {
		return nil, err
	}

	return &keystore.DeleteKeyResponse{}, nil
}

func (v1 *V1) EnableKey(ctx context.Context, req *keystore.EnableKeyRequest) (*keystore.EnableKeyResponse, error) {
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
	_, err = v1.KeystoreInstanceKeyOperationPluginClient.EnableKey(ctx, in)
	if err != nil {
		return nil, err
	}

	return &keystore.EnableKeyResponse{}, nil
}

func (v1 *V1) GetImportParameters(ctx context.Context, req *keystore.GetImportParametersRequest) (*keystore.GetImportParametersResponse, error) {
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
	grpcResp, err := v1.KeystoreInstanceKeyOperationPluginClient.GetImportParameters(ctx, in)
	if err != nil {
		return nil, err
	}

	return &keystore.GetImportParametersResponse{
		KeyID:            grpcResp.GetKeyId(),
		ImportParameters: grpcResp.GetImportParameters().AsMap(),
	}, nil
}

func (v1 *V1) ImportKeyMaterial(ctx context.Context, req *keystore.ImportKeyMaterialRequest) (*keystore.ImportKeyMaterialResponse, error) {
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
	_, err = v1.KeystoreInstanceKeyOperationPluginClient.ImportKeyMaterial(ctx, in)
	if err != nil {
		return nil, err
	}

	return &keystore.ImportKeyMaterialResponse{}, nil
}

func (v1 *V1) ValidateKey(ctx context.Context, req *keystore.ValidateKeyRequest) (*keystore.ValidateKeyResponse, error) {
	in := &operationsv1.ValidateKeyRequest{
		KeyType:     operationsv1.KeyType(req.KeyType),
		Algorithm:   operationsv1.KeyAlgorithm(req.KeyAlgorithm),
		Region:      req.Region,
		NativeKeyId: req.NativeKeyID,
	}
	grpcResp, err := v1.KeystoreInstanceKeyOperationPluginClient.ValidateKey(ctx, in)
	if err != nil {
		return nil, err
	}

	return &keystore.ValidateKeyResponse{
		IsValid: grpcResp.GetIsValid(),
		Message: grpcResp.GetMessage(),
	}, nil
}

func (v1 *V1) ValidateKeyAccessData(ctx context.Context, req *keystore.ValidateKeyAccessDataRequest) (*keystore.ValidateKeyAccessDataResponse, error) {
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
	grpcResp, err := v1.KeystoreInstanceKeyOperationPluginClient.ValidateKeyAccessData(ctx, in)
	if err != nil {
		return nil, err
	}

	return &keystore.ValidateKeyAccessDataResponse{
		IsValid: grpcResp.GetIsValid(),
		Message: grpcResp.GetMessage(),
	}, nil
}

func (v1 *V1) TransformCryptoAccessData(ctx context.Context, req *keystore.TransformCryptoAccessDataRequest) (*keystore.TransformCryptoAccessDataResponse, error) {
	in := &operationsv1.TransformCryptoAccessDataRequest{
		NativeKeyId: req.NativeKeyID,
		AccessData:  req.AccessData,
	}
	grpcResp, err := v1.KeystoreInstanceKeyOperationPluginClient.TransformCryptoAccessData(ctx, in)
	if err != nil {
		return nil, err
	}

	return &keystore.TransformCryptoAccessDataResponse{
		TransformedAccessData: grpcResp.GetTransformedAccessData(),
	}, nil
}

func (v1 *V1) ExtractKeyRegion(ctx context.Context, req *keystore.ExtractKeyRegionRequest) (*keystore.ExtractKeyRegionResponse, error) {
	management, err := structpb.NewStruct(req.ManagementAccessData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse values: %v", err)
	}

	in := &operationsv1.ExtractKeyRegionRequest{
		NativeKeyId:          req.NativeKeyID,
		ManagementAccessData: management,
	}
	grpcResp, err := v1.KeystoreInstanceKeyOperationPluginClient.ExtractKeyRegion(ctx, in)
	if err != nil {
		return nil, err
	}

	return &keystore.ExtractKeyRegionResponse{
		Region: grpcResp.GetRegion(),
	}, nil
}
