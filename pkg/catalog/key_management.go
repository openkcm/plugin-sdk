package catalog

import (
	"github.com/openkcm/plugin-sdk/api"
	"github.com/openkcm/plugin-sdk/pkg/plugin/key_management"
)

type keyManagementRepository struct {
	key_management.Repository
}

func (repo *keyManagementRepository) Binder() any {
	return repo.AddKeystoreKeyManager
}

func (repo *keyManagementRepository) Constraints() Constraints {
	return ExactlyOne()
}

func (repo *keyManagementRepository) Versions() []api.Version {
	return []api.Version{keyManagementV1{}}
}

func (repo *keyManagementRepository) BuiltIns() []BuiltInPlugin {
	return []BuiltInPlugin{}
}

type keyManagementV1 struct{}

func (keyManagementV1) New() api.Facade {
	return new(key_management.V1)
}
func (keyManagementV1) Deprecated() bool { return false }
