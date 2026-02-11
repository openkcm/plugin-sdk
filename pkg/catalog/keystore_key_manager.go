package catalog

import (
	"github.com/openkcm/plugin-sdk/api"
	"github.com/openkcm/plugin-sdk/pkg/plugin/keystore_key_manager"
)

type keystoreInstanceKeyOperationsRepository struct {
	keystore_key_manager.Repository
}

func (repo *keystoreInstanceKeyOperationsRepository) Binder() any {
	return repo.AddKeystoreKeyManager
}

func (repo *keystoreInstanceKeyOperationsRepository) Constraints() Constraints {
	return ExactlyOne()
}

func (repo *keystoreInstanceKeyOperationsRepository) Versions() []api.Version {
	return []api.Version{keystoreKeyManagerV1{}}
}

func (repo *keystoreInstanceKeyOperationsRepository) BuiltIns() []BuiltInPlugin {
	return []BuiltInPlugin{}
}

type keystoreKeyManagerV1 struct{}

func (keystoreKeyManagerV1) New() api.Facade {
	return new(keystore_key_manager.V1)
}
func (keystoreKeyManagerV1) Deprecated() bool { return false }
