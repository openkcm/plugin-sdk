package errors_test

import (
	"testing"

	keystoreErrs "github.com/openkcm/plugin-sdk/pkg/plugin/keystore/errors"
)

func TestStatusProviderAuthenticationErrorWithDetails(t *testing.T) {
	err := keystoreErrs.NewGrpcErrorWithReason(
		keystoreErrs.StatusProviderAuthenticationError,
		"Invalid credentials",
	)

	extractedReason := keystoreErrs.GetReason(err)
	if extractedReason != "Invalid credentials" {
		t.Errorf("Expected reason 'Invalid credentials', got '%s'", extractedReason)
	}

	if !keystoreErrs.IsStatus(err, keystoreErrs.StatusProviderAuthenticationError) {
		t.Errorf("Error does not match StatusProviderAuthenticationError")
	}

}
