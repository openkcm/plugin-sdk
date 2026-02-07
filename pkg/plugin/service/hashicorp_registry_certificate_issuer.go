package service

import (
	"fmt"

	"github.com/openkcm/plugin-sdk/api/service"
	"github.com/openkcm/plugin-sdk/api/service/certificateissuer"
	certificate_issuerv1 "github.com/openkcm/plugin-sdk/proto/plugin/certificate_issuer/v1"
)

func (h *hashicorpPluginServiceRegistry) CertificateIssuerByName(name string) (certificateissuer.CertificateIssuer, error) {
	plugin := h.catalog.LookupByTypeAndName(certificate_issuerv1.Type, name)
	if plugin == nil {
		return nil, fmt.Errorf("unable to find certificate issuer plugin %q", name)
	}
	return NewCertificateIssuerV1Plugin(plugin), nil
}

func (h *hashicorpPluginServiceRegistry) CertificateIssuerByNameAndVersion(version service.Version, name string) (certificateissuer.CertificateIssuer, error) {
	switch version {
	case service.V1:
		plugin := h.catalog.LookupByTypeAndName(certificate_issuerv1.Type, name)
		if plugin == nil {
			return nil, fmt.Errorf("unable to find certificate issuer plugin %q", name)
		}
		return NewCertificateIssuerV1Plugin(plugin), nil
	}

	return nil, service.ErrVersionNotSupported
}

func (h *hashicorpPluginServiceRegistry) ListCertificateIssuerByVersion(version service.Version) ([]certificateissuer.CertificateIssuer, error) {
	var issuers []certificateissuer.CertificateIssuer
	switch version {
	case service.V1:
		plugins := h.catalog.LookupByType(certificate_issuerv1.Type)
		for _, plugin := range plugins {
			issuers = append(issuers, NewCertificateIssuerV1Plugin(plugin))
		}
		return issuers, nil
	}

	return issuers, service.ErrVersionNotSupported
}
