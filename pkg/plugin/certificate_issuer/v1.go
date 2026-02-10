package certificate_issuer

import (
	"context"

	"github.com/openkcm/plugin-sdk/api"
	"github.com/openkcm/plugin-sdk/api/service/certificateissuer"
	"github.com/openkcm/plugin-sdk/pkg/plugin"
	certificate_issuerv1 "github.com/openkcm/plugin-sdk/proto/plugin/certificate_issuer/v1"
)

type V1 struct {
	plugin.Facade
	certificate_issuerv1.CertificateIssuerServicePluginClient
}

func (v1 *V1) Version() uint {
	return 1
}

func (v1 *V1) ServiceInfo() api.Info {
	return v1.Info
}

func (v1 *V1) IssueCertificate(ctx context.Context, req *certificateissuer.IssueCertificateRequest) (*certificateissuer.IssueCertificateResponse, error) {
	in := &certificate_issuerv1.GetCertificateRequest{
		CommonName: req.CommonName,
		Locality:   req.Localities,
		Validity:   CertificateValidityToGRPC(req.Validity),
		PrivateKey: CertificatePrivateKeyToGRPC(req.PrivateKey),
	}
	grpcResp, err := v1.GetCertificate(ctx, in)
	if err != nil {
		return nil, err
	}
	return &certificateissuer.IssueCertificateResponse{
		ChainPem: grpcResp.CertificateChain,
	}, nil
}

func CertificateValidityToGRPC(v *certificateissuer.CertificateValidity) *certificate_issuerv1.GetCertificateValidity {
	if v == nil {
		return nil
	}
	return &certificate_issuerv1.GetCertificateValidity{
		Value: v.Value,
		Type:  certificate_issuerv1.ValidityType(v.Type),
	}
}

func CertificatePrivateKeyToGRPC(pk *certificateissuer.CertificatePrivateKey) *certificate_issuerv1.PrivateKey {
	if pk == nil {
		return nil
	}
	return &certificate_issuerv1.PrivateKey{
		Data: pk.Data,
	}
}
