package catalog

import (
	"github.com/openkcm/plugin-sdk/api"
	"github.com/openkcm/plugin-sdk/pkg/plugin/identity_management"
)

type identityManagementRepository struct {
	identity_management.Repository
}

func (repo *identityManagementRepository) Binder() any {
	return repo.SetIdentityManagement
}

func (repo *identityManagementRepository) Constraints() Constraints {
	return ExactlyOne()
}

func (repo *identityManagementRepository) Versions() []api.Version {
	return []api.Version{identityManagementV1{}}
}

func (repo *identityManagementRepository) BuiltIns() []BuiltInPlugin {
	return []BuiltInPlugin{}
}

type identityManagementV1 struct{}

func (identityManagementV1) New() api.Facade  { return new(identity_management.V1) }
func (identityManagementV1) Deprecated() bool { return false }
