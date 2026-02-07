package service

import (
	"context"

	"github.com/openkcm/plugin-sdk/api/service/certificateissuer"
	"github.com/openkcm/plugin-sdk/pkg/catalog"
	certificate_issuerv1 "github.com/openkcm/plugin-sdk/proto/plugin/certificate_issuer/v1"
)

var _ certificateissuer.CertificateIssuer = (*hashicorpCertificateIssuerV1Plugin)(nil)

type hashicorpCertificateIssuerV1Plugin struct {
	plugin     *catalog.Plugin
	grpcClient certificate_issuerv1.CertificateIssuerServiceClient
}

func NewCertificateIssuerV1Plugin(plugin *catalog.Plugin) certificateissuer.CertificateIssuer {
	return &hashicorpCertificateIssuerV1Plugin{
		plugin:     plugin,
		grpcClient: certificate_issuerv1.NewCertificateIssuerServiceClient(plugin.ClientConnection()),
	}
}

func (h *hashicorpCertificateIssuerV1Plugin) GetCertificate(ctx context.Context, req *certificateissuer.GetCertificateRequest) (*certificateissuer.GetCertificateResponse, error) {
	in := &certificate_issuerv1.GetCertificateRequest{
		CommonName: req.CommonName,
		Locality:   req.Localities,
		Validity:   CertificateValidityToGRPC(req.Validity),
		PrivateKey: CertificatePrivateKeyToGRPC(req.PrivateKey),
	}
	grpcResp, err := h.grpcClient.GetCertificate(ctx, in)
	if err != nil {
		return nil, err
	}
	return &certificateissuer.GetCertificateResponse{
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
