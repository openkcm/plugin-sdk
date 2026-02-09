package systeminformation

import "context"

type SystemInformation interface {
	Get(ctx context.Context, req *GetSystemInformationRequest) (*GetSystemInformationResponse, error)
}

type RequestType int32

const (
	Unspecified RequestType = iota
	System
	Subaccount
)

type GetSystemInformationRequest struct {
	// V1 Fields
	ID   string
	Type RequestType
}

type GetSystemInformationResponse struct {
	// V1 Fields
	Metadata map[string]string
}
