package catalog

import (
	"github.com/openkcm/plugin-sdk/api"
	"github.com/openkcm/plugin-sdk/pkg/plugin/keystore_key_manager"
)

type keystoreKeyManagerRepository struct {
	keystore_key_manager.Repository
}

func (repo *keystoreKeyManagerRepository) Binder() any {
	return repo.AddKeystoreKeyManager
}

func (repo *keystoreKeyManagerRepository) Constraints() Constraints {
	return ExactlyOne()
}

func (repo *keystoreKeyManagerRepository) Versions() []api.Version {
	return []api.Version{keystoreKeyManagerV1{}}
}

func (repo *keystoreKeyManagerRepository) BuiltIns() []BuiltInPlugin {
	return []BuiltInPlugin{}
}

type keystoreKeyManagerV1 struct{}

func (keystoreKeyManagerV1) New() api.Facade {
	return new(keystore_key_manager.V1)
}
func (keystoreKeyManagerV1) Deprecated() bool { return false }
