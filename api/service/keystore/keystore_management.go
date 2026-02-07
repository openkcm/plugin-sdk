package keystore

import (
	"context"
)

type KeystoreManagement interface {
	CreateKeystore(ctx context.Context, req *CreateKeystoreRequest) (*CreateKeystoreResponse, error)
	DeleteKeystore(ctx context.Context, req *DeleteKeystoreRequest) (*DeleteKeystoreResponse, error)
}

type CreateKeystoreRequest struct {
	// V1 Fields
	Values map[string]any
}

type CreateKeystoreResponse struct {
	// V1 Fields
	Config KeystoreInstanceConfig
}

type DeleteKeystoreRequest struct {
	// V1 Fields
	Config KeystoreInstanceConfig
}

type DeleteKeystoreResponse struct{}
