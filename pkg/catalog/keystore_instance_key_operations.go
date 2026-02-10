package catalog

import (
	"github.com/openkcm/plugin-sdk/api"
	"github.com/openkcm/plugin-sdk/pkg/plugin/keystore_instance_key_operations"
)

type keystoreInstanceKeyOperationsRepository struct {
	keystore_instance_key_operations.Repository
}

func (repo *keystoreInstanceKeyOperationsRepository) Binder() any {
	return repo.SetKeystoreInstanceKeyOperations
}

func (repo *keystoreInstanceKeyOperationsRepository) Constraints() Constraints {
	return ExactlyOne()
}

func (repo *keystoreInstanceKeyOperationsRepository) Versions() []api.Version {
	return []api.Version{keystoreInstanceKeyOperationsV1{}}
}

func (repo *keystoreInstanceKeyOperationsRepository) BuiltIns() []BuiltInPlugin {
	return []BuiltInPlugin{}
}

type keystoreInstanceKeyOperationsV1 struct{}

func (keystoreInstanceKeyOperationsV1) New() api.Facade {
	return new(keystore_instance_key_operations.V1)
}
func (keystoreInstanceKeyOperationsV1) Deprecated() bool { return false }
