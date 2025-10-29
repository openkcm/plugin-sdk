package errors_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/status"

	keystoreErrs "github.com/openkcm/plugin-sdk/pkg/plugin/keystore/errors"
)

func TestNewGrpcErrorWithDetails(t *testing.T) {
	metadata := map[string]string{"foo": "bar"}
	err := keystoreErrs.NewGrpcErrorWithDetails(
		keystoreErrs.StatusInvalidKeyAccessData,
		keystoreErrs.ReasonMissingFields,
		metadata,
	)
	reason, meta := keystoreErrs.GetDetails(err)
	assert.Equal(t, keystoreErrs.ReasonMissingFields, reason)
	assert.Equal(t, metadata, meta)
}

func TestWithMetadata_NewMetadata(t *testing.T) {
	err := keystoreErrs.NewGrpcErrorWithDetails(
		keystoreErrs.StatusInvalidKeyAccessData,
		keystoreErrs.ReasonMissingFields,
		nil,
	)
	newMeta := map[string]string{
		keystoreErrs.MetadataMissingFields: "keyId,accessData",
	}
	mergedErr := keystoreErrs.WithMetadata(err, newMeta)
	_, meta := keystoreErrs.GetDetails(mergedErr)
	assert.Equal(t, map[string]string{
		"missingFields": "keyId,accessData",
	}, meta)
}

func TestWithMetadata_MergesMetadata(t *testing.T) {
	origMeta := map[string]string{"foo": "bar"}
	err := keystoreErrs.NewGrpcErrorWithDetails(
		keystoreErrs.StatusInvalidKeyAccessData,
		keystoreErrs.ReasonMissingFields,
		origMeta,
	)
	newMeta := map[string]string{
		keystoreErrs.MetadataMissingFields: "keyId,accessData",
	}
	mergedErr := keystoreErrs.WithMetadata(err, newMeta)
	_, meta := keystoreErrs.GetDetails(mergedErr)
	assert.Equal(t, map[string]string{
		"foo":           "bar",
		"missingFields": "keyId,accessData",
	}, meta)
}

func TestIsStatus(t *testing.T) {
	err := keystoreErrs.NewGrpcErrorWithDetails(
		keystoreErrs.StatusKeyNotFound,
		"NO_LONGER_EXISTS",
		nil,
	)
	assert.True(t, keystoreErrs.IsStatus(err, keystoreErrs.StatusKeyNotFound))
	assert.False(t, keystoreErrs.IsStatus(err, keystoreErrs.StatusInvalidKeyAccessData))
}

func TestGetDetails_NoDetails(t *testing.T) {
	err := status.New(1, "no details").Err()
	reason, meta := keystoreErrs.GetDetails(err)
	assert.Empty(t, reason)
	assert.Nil(t, meta)
}

func TestGetDetails_NormalError(t *testing.T) {
	err := errors.New("normal error")
	reason, meta := keystoreErrs.GetDetails(err)
	assert.Equal(t, "", reason)
	assert.Nil(t, meta)
}
