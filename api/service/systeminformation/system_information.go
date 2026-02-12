package systeminformation

import (
	"context"

	"github.com/openkcm/plugin-sdk/api"
)

type SystemInformation interface {
	ServiceInfo() api.Info

	GetSystemInfo(ctx context.Context, req *GetSystemInfoRequest) (*GetSystemInfoResponse, error)
}

type Type int32

const (
	SystemType Type = iota
	SubaccountType
	AccountType
)

type GetSystemInfoRequest struct {
	// V1 Fields
	ID   string
	Type Type
}

type GetSystemInfoResponse struct {
	// V1 Fields
	Metadata map[string]string
}
