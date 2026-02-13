package certificate_issuer

import (
	"github.com/openkcm/plugin-sdk/api/service/certificateissuer"
)

type Repository struct {
	CertificateIssuer certificateissuer.CertificateIssuer
}

func (repo *Repository) GetCertificateIssuer() certificateissuer.CertificateIssuer {
	return repo.CertificateIssuer
}

func (repo *Repository) SetCertificateIssuer(instance certificateissuer.CertificateIssuer) {
	repo.CertificateIssuer = instance
}

func (repo *Repository) Clear() {
	repo.CertificateIssuer = nil
}
