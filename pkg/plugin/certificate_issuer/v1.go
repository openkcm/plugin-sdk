package certificate_issuer

import (
	"context"

	"github.com/openkcm/plugin-sdk/api"
	"github.com/openkcm/plugin-sdk/api/service/certificateissuer"
	"github.com/openkcm/plugin-sdk/pkg/plugin"
	grpccertificateissuerv1 "github.com/openkcm/plugin-sdk/proto/plugin/certificate_issuer/v1"
)

type V1 struct {
	plugin.Facade
	grpccertificateissuerv1.CertificateIssuerPluginClient
}

func (v1 *V1) Version() uint {
	return 1
}

func (v1 *V1) ServiceInfo() api.Info {
	return v1.Info
}

func (v1 *V1) IssueCertificate(ctx context.Context, req *certificateissuer.IssueCertificateRequest) (*certificateissuer.IssueCertificateResponse, error) {
	in := &grpccertificateissuerv1.IssueCertificateRequest{
		CommonName: req.CommonName,
		Locality:   req.Localities,
		Validity:   CertificateValidityToGRPC(req.Validity),
		PrivateKey: CertificatePrivateKeyToGRPC(req.PrivateKey),
	}
	grpcResp, err := v1.CertificateIssuerPluginClient.IssueCertificate(ctx, in)
	if err != nil {
		return nil, err
	}
	return &certificateissuer.IssueCertificateResponse{
		ChainPem: grpcResp.CertificateChainPem,
	}, nil
}

func CertificateValidityToGRPC(v *certificateissuer.CertificateValidity) *grpccertificateissuerv1.CertificateValidity {
	if v == nil {
		return nil
	}
	return &grpccertificateissuerv1.CertificateValidity{
		Value: v.Value,
		Type:  grpccertificateissuerv1.ValidityType(v.Type),
	}
}

func CertificatePrivateKeyToGRPC(pk *certificateissuer.CertificatePrivateKey) *grpccertificateissuerv1.PrivateKey {
	if pk == nil {
		return nil
	}
	return &grpccertificateissuerv1.PrivateKey{
		Data: pk.Data,
	}
}
