package systeminformation

import "context"

type SystemInformation interface {
	Get(ctx context.Context, req *GetRequest) (*GetResponse, error)
}

type RequestType int32

const (
	Unspecified RequestType = iota
	System
	Subaccount
)

type GetRequest struct {
	// V1 Fields
	ID   string
	Type RequestType
}

type GetResponse struct {
	// V1 Fields
	Metadata map[string]string
}
