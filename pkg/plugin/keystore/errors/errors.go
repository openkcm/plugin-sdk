package errors

import (
	"fmt"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// List of common error reasons (sub-codes) that will be consumed by both client and server.
// Must be in UPPER_SNAKE_CASE according to google.golang.org/genproto/googleapis/rpc/errdetails

const (
	ReasonMissingFields   = "MISSING_FIELDS"    // Data is missing required fields
	ReasonCannotParseData = "CANNOT_PARSE_DATA" // Data could not be parsed
)

// List of common metadata keys for error details that will be consumed by both client and server.
// Must be in lowerCamelCase according to google.golang.org/genproto/googleapis/rpc/errdetails

const (
	MetadataMissingFields = "missingFields" // Comma-separated list of missing fields
)

// Predefined gRPC status errors for common keystore error scenarios.
var (
	// StatusProviderAuthenticationError indicates a failure to authenticate with the keystore provider.
	StatusProviderAuthenticationError = status.New(
		codes.InvalidArgument, "failed to authenticate with the keystore provider")

	// StatusInvalidKeyAccessData indicates that the provided key access data (management or crypto) is invalid.
	StatusInvalidKeyAccessData = status.New(
		codes.InvalidArgument, "invalid key access data")

	// StatusKeyNotFound indicates that the specified key was not found in the keystore provider.
	StatusKeyNotFound = status.New(
		codes.NotFound, "key not found in the keystore provider")
)

// NewGrpcErrorWithDetails creates a gRPC error with the given status and metadata mapping.
func NewGrpcErrorWithDetails(st *status.Status, reason string, metadata map[string]string) error {
	errInfo := &errdetails.ErrorInfo{
		Reason:   reason,
		Metadata: metadata,
	}
	st, err := st.WithDetails(errInfo)
	if err != nil {
		return fmt.Errorf("failed to add details to status: %w", err)
	}
	return st.Err()
}

// WithMetadata adds additional metadata to an existing gRPC error.
// If the error does not contain gRPC status information, it is returned unchanged.
// If the error already has metadata, the new metadata is merged in, with new values overwriting existing ones.
func WithMetadata(err error, metadata map[string]string) error {
	st, ok := status.FromError(err)
	if !ok {
		return err
	}

	var (
		reason      string
		newMetadata map[string]string
	)
	for _, detail := range st.Details() {
		if errInfo, ok := detail.(*errdetails.ErrorInfo); ok {
			reason = errInfo.Reason
			newMetadata = errInfo.Metadata
			break
		}
	}

	if newMetadata == nil {
		newMetadata = make(map[string]string)
	}
	for k, v := range metadata {
		newMetadata[k] = v
	}

	errInfo := &errdetails.ErrorInfo{
		Reason:   reason,
		Metadata: newMetadata,
	}

	st, err = status.New(st.Code(), st.Message()).WithDetails(errInfo)
	if err != nil {
		return fmt.Errorf("failed to add details to status: %w", err)
	}

	return st.Err()
}

// IsStatus checks if the given error matches the provided gRPC status.
func IsStatus(err error, st *status.Status) bool {
	convertedErr := status.Convert(err)
	return convertedErr.Code() == st.Code() && convertedErr.Message() == st.Message()
}

// GetDetails extracts the reason and metadata from the given gRPC error, if available.
func GetDetails(err error) (string, map[string]string) {
	st, ok := status.FromError(err)
	if !ok {
		return "", nil
	}
	for _, detail := range st.Details() {
		if errInfo, ok := detail.(*errdetails.ErrorInfo); ok {
			return errInfo.Reason, errInfo.Metadata
		}
	}
	return "", nil
}
