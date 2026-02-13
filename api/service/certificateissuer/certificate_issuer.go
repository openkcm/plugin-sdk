package certificateissuer

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/openkcm/plugin-sdk/api"
)

type CertificateIssuer interface {
	ServiceInfo() api.Info

	IssueCertificate(ctx context.Context, req *IssueCertificateRequest) (*IssueCertificateResponse, error)
}
type CertificateFormat int

const (
	CertificateFormatUnspecified CertificateFormat = iota
	CertificateFormatPEM
	CertificateFormatDER
	CertificateFormatPKCS7
)

type KeyFormat int

const (
	KeyFormatUnspecified KeyFormat = iota
	KeyFormatPKCS1
	KeyFormatPKCS8
	KeyFormatSEC1
)

func (k KeyFormat) String() string {
	switch k {
	case KeyFormatPKCS1:
		return "PKCS1"
	case KeyFormatPKCS8:
		return "PKCS8"
	case KeyFormatSEC1:
		return "SEC1"
	default:
		return "UNSPECIFIED"
	}
}

type ValidityUnit int

const (
	ValidityUnitUnspecified ValidityUnit = iota
	ValidityUnitDays
	ValidityUnitMonths
	ValidityUnitYears
)

// Domain Models

type RelativeValidity struct {
	Value int32
	Unit  ValidityUnit
}

// CertificateLifetime represents the oneof field.
// In pure Go, using pointers allows us to check which field is active (not nil).
type CertificateLifetime struct {
	Duration *time.Duration
	NotAfter *time.Time
	Relative *RelativeValidity
}

type Subject struct {
	CommonName         string
	SerialNumber       *string
	Country            []string
	Organization       []string
	OrganizationalUnit []string
	Locality           []string
	Province           []string
	StreetAddress      []string
	PostalCode         []string
}

type PrivateKey struct {
	Data   []byte
	Format KeyFormat
}

type IssueCertificateRequest struct {
	Lifetime        CertificateLifetime
	Subject         Subject
	PrivateKey      PrivateKey
	PreferredFormat CertificateFormat
}

type IssueCertificateResponse struct {
	CertificateData []byte
	Format          CertificateFormat
	CAChain         [][]byte
}

type SupportedKeyFormatsError struct {
	RejectedFormat   string
	SupportedFormats []KeyFormat
	Reason           string
}

func (e *SupportedKeyFormatsError) Error() string {
	var formats []string
	for _, f := range e.SupportedFormats {
		formats = append(formats, f.String())
	}

	return fmt.Sprintf("unsupported private key format '%s': %s (supported formats: %s)",
		e.RejectedFormat,
		e.Reason,
		strings.Join(formats, ", "),
	)
}
