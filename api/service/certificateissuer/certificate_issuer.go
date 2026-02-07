package certificateissuer

import (
	"context"
)

type CertificateIssuer interface {
	GetCertificate(ctx context.Context, req *GetCertificateRequest) (*GetCertificateResponse, error)
}

type ValidityType int32

const (
	Unspecified ValidityType = iota
	Days
	Months
	Years
)

type GetCertificateRequest struct {
	// V1 Fields
	CommonName string
	Localities []string
	Validity   *CertificateValidity
	PrivateKey *CertificatePrivateKey
}

type GetCertificateResponse struct {
	// V1 Fields
	ChainPem string
}

type CertificateValidity struct {
	// V1 Fields
	Value int64
	Type  ValidityType
}

type CertificatePrivateKey struct {
	// V1 Fields
	Data []byte
}
