package key_management

import (
	"context"
	"fmt"

	"buf.build/go/protovalidate"
	"google.golang.org/protobuf/types/known/structpb"

	"github.com/openkcm/plugin-sdk/api"
	"github.com/openkcm/plugin-sdk/pkg/plugin"
	grpccommonv1 "github.com/openkcm/plugin-sdk/proto/plugin/common/v1"
	grpckeymanagementv1 "github.com/openkcm/plugin-sdk/proto/plugin/key_management/v1"
	"github.com/openkcm/plugin-sdk/service/api/keymanagement"
)

const (
	msgFailedValidation = "failed to validate request object"
	msgFailedParseValue = "failed to parse values"
)

type V1 struct {
	plugin.Facade
	grpckeymanagementv1.KeyManagementPluginClient
}

func (v1 *V1) Version() uint {
	return 1
}

func (v1 *V1) ServiceInfo() api.Info {
	return v1.Info
}

func (v1 *V1) GetKey(ctx context.Context, req *keymanagement.GetKeyRequest) (*keymanagement.GetKeyResponse, error) {
	value, err := structpb.NewStruct(req.Parameters.Config.Values)
	if err != nil {
		return nil, failureError(msgFailedParseValue, err)
	}

	in := &grpckeymanagementv1.GetKeyRequest{
		Parameters: &grpckeymanagementv1.RequestParameters{
			KeystoreConfig: &grpccommonv1.KeystoreConfig{
				ConfigurationParameters: value,
			},
			KeyId: req.Parameters.KeyID,
		},
	}
	if err := protovalidate.Validate(in); err != nil {
		return nil, failureError(msgFailedValidation, err)
	}
	grpcResp, err := v1.KeyManagementPluginClient.GetKey(ctx, in)
	if err != nil {
		return nil, err
	}

	return &keymanagement.GetKeyResponse{
		KeyID:        grpcResp.GetKeyId(),
		KeyAlgorithm: keymanagement.KeyAlgorithm(grpcResp.GetAlgorithm()),
		Status:       grpcResp.GetStatus(),
		Usage:        grpcResp.GetUsage(),
	}, nil
}

func (v1 *V1) CreateKey(ctx context.Context, req *keymanagement.CreateKeyRequest) (*keymanagement.CreateKeyResponse, error) {
	value, err := structpb.NewStruct(req.Config.Values)
	if err != nil {
		return nil, failureError(msgFailedParseValue, err)
	}

	in := &grpckeymanagementv1.CreateKeyRequest{
		KeystoreConfig: &grpccommonv1.KeystoreConfig{
			ConfigurationParameters: value,
		},
		Algorithm: grpckeymanagementv1.KeyAlgorithm(req.KeyAlgorithm),
		Id:        req.ID,
		Region:    req.Region,
		KeyType:   grpckeymanagementv1.KeyType(req.KeyType),
	}
	if err := protovalidate.Validate(in); err != nil {
		return nil, failureError(msgFailedValidation, err)
	}
	grpcResp, err := v1.KeyManagementPluginClient.CreateKey(ctx, in)
	if err != nil {
		return nil, err
	}

	return &keymanagement.CreateKeyResponse{
		KeyID:  grpcResp.GetKeyId(),
		Status: grpcResp.GetStatus(),
	}, nil
}

func (v1 *V1) DeleteKey(ctx context.Context, req *keymanagement.DeleteKeyRequest) (*keymanagement.DeleteKeyResponse, error) {
	value, err := structpb.NewStruct(req.Parameters.Config.Values)
	if err != nil {
		return nil, failureError(msgFailedParseValue, err)
	}

	in := &grpckeymanagementv1.DeleteKeyRequest{
		Parameters: &grpckeymanagementv1.RequestParameters{
			KeystoreConfig: &grpccommonv1.KeystoreConfig{
				ConfigurationParameters: value,
			},
			KeyId: req.Parameters.KeyID,
		},
		Window: req.Window,
	}
	if err := protovalidate.Validate(in); err != nil {
		return nil, failureError(msgFailedValidation, err)
	}
	_, err = v1.KeyManagementPluginClient.DeleteKey(ctx, in)
	if err != nil {
		return nil, err
	}

	return &keymanagement.DeleteKeyResponse{}, nil
}

func (v1 *V1) EnableKey(ctx context.Context, req *keymanagement.EnableKeyRequest) (*keymanagement.EnableKeyResponse, error) {
	value, err := structpb.NewStruct(req.Parameters.Config.Values)
	if err != nil {
		return nil, failureError(msgFailedParseValue, err)
	}

	in := &grpckeymanagementv1.EnableKeyRequest{
		Parameters: &grpckeymanagementv1.RequestParameters{
			KeystoreConfig: &grpccommonv1.KeystoreConfig{
				ConfigurationParameters: value,
			},
			KeyId: req.Parameters.KeyID,
		},
	}
	if err := protovalidate.Validate(in); err != nil {
		return nil, failureError(msgFailedValidation, err)
	}
	_, err = v1.KeyManagementPluginClient.EnableKey(ctx, in)
	if err != nil {
		return nil, err
	}

	return &keymanagement.EnableKeyResponse{}, nil
}

func (v1 *V1) GetImportParameters(ctx context.Context, req *keymanagement.GetImportParametersRequest) (*keymanagement.GetImportParametersResponse, error) {
	value, err := structpb.NewStruct(req.Parameters.Config.Values)
	if err != nil {
		return nil, failureError(msgFailedParseValue, err)
	}

	in := &grpckeymanagementv1.GetImportParametersRequest{
		Parameters: &grpckeymanagementv1.RequestParameters{
			KeystoreConfig: &grpccommonv1.KeystoreConfig{
				ConfigurationParameters: value,
			},
			KeyId: req.Parameters.KeyID,
		},
		Algorithm: grpckeymanagementv1.KeyAlgorithm(req.KeyAlgorithm),
	}
	if err := protovalidate.Validate(in); err != nil {
		return nil, failureError(msgFailedValidation, err)
	}

	grpcResp, err := v1.KeyManagementPluginClient.GetImportParameters(ctx, in)
	if err != nil {
		return nil, err
	}

	return &keymanagement.GetImportParametersResponse{
		KeyID:            grpcResp.GetKeyId(),
		ImportParameters: grpcResp.GetImportParameters().AsMap(),
	}, nil
}

func (v1 *V1) ImportKeyMaterial(ctx context.Context, req *keymanagement.ImportKeyMaterialRequest) (*keymanagement.ImportKeyMaterialResponse, error) {
	value, err := structpb.NewStruct(req.Parameters.Config.Values)
	if err != nil {
		return nil, failureError(msgFailedParseValue, err)
	}

	importParams, err := structpb.NewStruct(req.ImportParameters)
	if err != nil {
		return nil, failureError(msgFailedParseValue, err)
	}

	in := &grpckeymanagementv1.ImportKeyMaterialRequest{
		Parameters: &grpckeymanagementv1.RequestParameters{
			KeystoreConfig: &grpccommonv1.KeystoreConfig{
				ConfigurationParameters: value,
			},
			KeyId: req.Parameters.KeyID,
		},
		ImportParameters: importParams,
	}
	if err := protovalidate.Validate(in); err != nil {
		return nil, failureError(msgFailedValidation, err)
	}

	_, err = v1.KeyManagementPluginClient.ImportKeyMaterial(ctx, in)
	if err != nil {
		return nil, err
	}

	return &keymanagement.ImportKeyMaterialResponse{}, nil
}

func (v1 *V1) ValidateKey(ctx context.Context, req *keymanagement.ValidateKeyRequest) (*keymanagement.ValidateKeyResponse, error) {
	in := &grpckeymanagementv1.ValidateKeyRequest{
		KeyType:     grpckeymanagementv1.KeyType(req.KeyType),
		Algorithm:   grpckeymanagementv1.KeyAlgorithm(req.KeyAlgorithm),
		Region:      req.Region,
		NativeKeyId: req.NativeKeyID,
	}
	if err := protovalidate.Validate(in); err != nil {
		return nil, failureError(msgFailedValidation, err)
	}

	grpcResp, err := v1.KeyManagementPluginClient.ValidateKey(ctx, in)
	if err != nil {
		return nil, err
	}

	return &keymanagement.ValidateKeyResponse{
		IsValid: grpcResp.GetIsValid(),
		Message: grpcResp.GetMessage(),
	}, nil
}

func (v1 *V1) ValidateKeyAccessData(ctx context.Context, req *keymanagement.ValidateKeyAccessDataRequest) (*keymanagement.ValidateKeyAccessDataResponse, error) {
	management, err := structpb.NewStruct(req.Management)
	if err != nil {
		return nil, failureError(msgFailedParseValue, err)
	}
	crypto, err := structpb.NewStruct(req.Crypto)
	if err != nil {
		return nil, failureError(msgFailedParseValue, err)
	}

	in := &grpckeymanagementv1.ValidateKeyAccessDataRequest{
		Management: management,
		Crypto:     crypto,
	}
	if err := protovalidate.Validate(in); err != nil {
		return nil, failureError(msgFailedValidation, err)
	}

	grpcResp, err := v1.KeyManagementPluginClient.ValidateKeyAccessData(ctx, in)
	if err != nil {
		return nil, err
	}

	return &keymanagement.ValidateKeyAccessDataResponse{
		IsValid: grpcResp.GetIsValid(),
		Message: grpcResp.GetMessage(),
	}, nil
}

func (v1 *V1) TransformCryptoAccessData(ctx context.Context, req *keymanagement.TransformCryptoAccessDataRequest) (*keymanagement.TransformCryptoAccessDataResponse, error) {
	in := &grpckeymanagementv1.TransformCryptoAccessDataRequest{
		NativeKeyId: req.NativeKeyID,
		AccessData:  req.AccessData,
	}
	if err := protovalidate.Validate(in); err != nil {
		return nil, failureError(msgFailedValidation, err)
	}

	grpcResp, err := v1.KeyManagementPluginClient.TransformCryptoAccessData(ctx, in)
	if err != nil {
		return nil, err
	}

	return &keymanagement.TransformCryptoAccessDataResponse{
		TransformedAccessData: grpcResp.GetTransformedAccessData(),
	}, nil
}

func (v1 *V1) ExtractKeyRegion(ctx context.Context, req *keymanagement.ExtractKeyRegionRequest) (*keymanagement.ExtractKeyRegionResponse, error) {
	management, err := structpb.NewStruct(req.ManagementAccessData)
	if err != nil {
		return nil, failureError(msgFailedParseValue, err)
	}

	in := &grpckeymanagementv1.ExtractKeyRegionRequest{
		NativeKeyId:          req.NativeKeyID,
		ManagementAccessData: management,
	}
	if err := protovalidate.Validate(in); err != nil {
		return nil, failureError(msgFailedValidation, err)
	}

	grpcResp, err := v1.KeyManagementPluginClient.ExtractKeyRegion(ctx, in)
	if err != nil {
		return nil, err
	}

	return &keymanagement.ExtractKeyRegionResponse{
		Region: grpcResp.GetRegion(),
	}, nil
}

func failureError(msg string, err error) error {
	return fmt.Errorf(msg+": %v", err)
}
