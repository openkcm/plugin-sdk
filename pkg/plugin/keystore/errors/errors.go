package errors

import (
	"fmt"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	StatusProviderAuthenticationError = status.New(
		codes.InvalidArgument, "failed to authenticate with the keystore provider")
	StatusKeyNotFound = status.New(
		codes.NotFound, "key not found in the keystore provider")
)

// NewGrpcErrorWithReason creates a gRPC error with the given status and reason.
func NewGrpcErrorWithReason(st *status.Status, reason string) error {
	errInfo := &errdetails.ErrorInfo{
		Reason: reason,
	}
	st, err := st.WithDetails(errInfo)
	if err != nil {
		return fmt.Errorf("failed to add reason to status: %w", err)
	}
	return st.Err()
}

// IsStatus checks if the given error matches the provided gRPC status.
func IsStatus(err error, st *status.Status) bool {
	convertedErr := status.Convert(err)
	return convertedErr.Code() == st.Code() && convertedErr.Message() == st.Message()
}

// GetReason extracts the reason from the given gRPC error, if available.
func GetReason(err error) string {
	st, ok := status.FromError(err)
	if !ok {
		return ""
	}
	for _, detail := range st.Details() {
		if errInfo, ok := detail.(*errdetails.ErrorInfo); ok {
			return errInfo.Reason
		}
	}
	return ""
}
