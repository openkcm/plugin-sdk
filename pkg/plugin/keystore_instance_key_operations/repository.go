package keystore_instance_key_operations

import (
	keystoreapi "github.com/openkcm/plugin-sdk/api/service/keystore"
)

type Repository struct {
	KeystoreInstanceKeyOperations keystoreapi.KeystoreInstanceKeyOperations
}

func (repo *Repository) GetKeystoreInstanceKeyOperations() keystoreapi.KeystoreInstanceKeyOperations {
	return repo.KeystoreInstanceKeyOperations
}

func (repo *Repository) SetKeystoreInstanceKeyOperations(instance keystoreapi.KeystoreInstanceKeyOperations) {
	repo.KeystoreInstanceKeyOperations = instance
}

func (repo *Repository) Clear() {
	repo.KeystoreInstanceKeyOperations = nil
}
