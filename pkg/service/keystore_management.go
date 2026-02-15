package service

import (
	"github.com/openkcm/plugin-sdk/api"
	"github.com/openkcm/plugin-sdk/pkg/catalog"
	"github.com/openkcm/plugin-sdk/pkg/service/internal/keystore_management"
)

type keystoreManagementRepository struct {
	keystore_management.Repository
}

func (repo *keystoreManagementRepository) Binder() any {
	return repo.AddKeystoreManagement
}

func (repo *keystoreManagementRepository) Constraints() catalog.Constraints {
	return catalog.ExactlyOne()
}

func (repo *keystoreManagementRepository) Versions() []api.Version {
	return []api.Version{keystoreManagementV1{}}
}

type keystoreManagementV1 struct{}

func (keystoreManagementV1) New() api.Facade  { return new(keystore_management.V1) }
func (keystoreManagementV1) Deprecated() bool { return false }
