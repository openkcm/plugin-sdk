package catalog

import (
	"github.com/openkcm/plugin-sdk/api"
	"github.com/openkcm/plugin-sdk/pkg/plugin/keystore_management"
	"github.com/openkcm/plugin-sdk/pkg/plugin/system_information"
)

type keystoreManagementRepository struct {
	keystore_management.Repository
}

func (repo *keystoreManagementRepository) Binder() any {
	return repo.SetKeystoreManagement
}

func (repo *keystoreManagementRepository) Constraints() Constraints {
	return ExactlyOne()
}

func (repo *keystoreManagementRepository) Versions() []api.Version {
	return []api.Version{keystoreManagementV1{}}
}

func (repo *keystoreManagementRepository) BuiltIns() []BuiltInPlugin {
	return []BuiltInPlugin{}
}

type keystoreManagementV1 struct{}

func (keystoreManagementV1) New() api.Facade  { return new(system_information.V1) }
func (keystoreManagementV1) Deprecated() bool { return false }
