package systeminformation

import "context"

type SystemInformation interface {
	GetSystemInfo(ctx context.Context, req *GetSystemInfoRequest) (*GetSystemInfoResponse, error)
}

type RequestType int32

const (
	Unspecified RequestType = iota
	System
	Subaccount
)

type GetSystemInfoRequest struct {
	// V1 Fields
	ID   string
	Type RequestType
}

type GetSystemInfoResponse struct {
	// V1 Fields
	Metadata map[string]string
}
