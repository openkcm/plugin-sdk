package certificateissuer

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"

	"github.com/cloudflare/cfssl/crypto/pkcs7"
)

// ParseCertificateChain extracts the leaf certificate and its accompanying CA chain.
// Cognitive Complexity: drastically reduced by delegating parsing to parseBytes.
func (resp *IssueCertificateResponse) ParseCertificateChain() ([]*x509.Certificate, error) {
	if len(resp.CertificateData) == 0 {
		return nil, errors.New("certificate data is empty")
	}

	var rawCerts []*x509.Certificate

	// Parse main payload
	certs, err := parseBytes(resp.CertificateData, resp.Format)
	if err != nil {
		return nil, fmt.Errorf("failed to parse certificate data: %w", err)
	}
	rawCerts = append(rawCerts, certs...)

	// Parse CA chain elements
	for i, caBytes := range resp.CAChain {
		certs, err := parseBytes(caBytes, resp.Format)
		if err != nil {
			return nil, fmt.Errorf("failed to parse CA chain at index %d: %w", i, err)
		}
		rawCerts = append(rawCerts, certs...)
	}

	return orderCertificateChain(rawCerts), nil
}

// parseBytes is a dedicated helper that handles the format switching.
// This isolates the branching logic from the iteration logic.
func parseBytes(data []byte, format CertificateFormat) ([]*x509.Certificate, error) {
	switch format {
	case CertificateFormatDER:
		cert, err := x509.ParseCertificate(data)
		if err != nil {
			return nil, err
		}
		return []*x509.Certificate{cert}, nil

	case CertificateFormatPEM, CertificateFormatUnspecified:
		block, _ := pem.Decode(data)
		if block == nil || block.Type != "CERTIFICATE" {
			return nil, errors.New("failed to decode valid PEM CERTIFICATE block")
		}
		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return nil, err
		}
		return []*x509.Certificate{cert}, nil

	case CertificateFormatPKCS7:
		p7, err := pkcs7.ParsePKCS7(data)
		if err != nil {
			return nil, err
		}
		if len(p7.Content.SignedData.Certificates) == 0 {
			return nil, errors.New("no certificates found in the PKCS7 container")
		}
		return p7.Content.SignedData.Certificates, nil

	default:
		return nil, errors.New("unsupported certificate format")
	}
}

// orderCertificateChain logically orders an unsorted slice of certificates.
// Cognitive Complexity: drastically reduced by using maps instead of nested loops.
func orderCertificateChain(certs []*x509.Certificate) []*x509.Certificate {
	if len(certs) <= 1 {
		return certs
	}

	subjects := make(map[string]*x509.Certificate)
	issuers := make(map[string]bool)

	// Build O(1) lookup maps
	for _, cert := range certs {
		subjects[string(cert.RawSubject)] = cert
		issuers[string(cert.RawIssuer)] = true
	}

	leaf := findLeaf(certs, issuers)
	if leaf == nil {
		return certs // Fallback if no clean leaf is found
	}

	return buildChainFromLeaf(leaf, subjects, len(certs))
}

// findLeaf identifies the certificate that hasn't issued any other certificate in the pool.
func findLeaf(certs []*x509.Certificate, issuers map[string]bool) *x509.Certificate {
	for _, cert := range certs {
		if !issuers[string(cert.RawSubject)] {
			return cert
		}
	}
	return nil
}

// buildChainFromLeaf walks up the cryptographic chain using the subjects map.
func buildChainFromLeaf(leaf *x509.Certificate, subjects map[string]*x509.Certificate, total int) []*x509.Certificate {
	ordered := make([]*x509.Certificate, 0, total)
	ordered = append(ordered, leaf)

	// Track added certs to prevent infinite loops (e.g., self-signed roots)
	added := make(map[string]bool)
	added[string(leaf.RawSubject)] = true

	current := leaf
	for len(ordered) < total {
		parent, exists := subjects[string(current.RawIssuer)]

		// Break if the chain breaks or we hit a circular loop/self-signed root
		if !exists || added[string(parent.RawSubject)] {
			break
		}

		ordered = append(ordered, parent)
		added[string(parent.RawSubject)] = true
		current = parent
	}

	// Append any remaining disconnected certificates
	for _, cert := range subjects {
		if !added[string(cert.RawSubject)] {
			ordered = append(ordered, cert)
		}
	}

	return ordered
}
