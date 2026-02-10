package keystore_management

import (
	keystoreapi "github.com/openkcm/plugin-sdk/api/service/keystore"
)

type Repository struct {
	KeystoreManagement keystoreapi.KeystoreManagement
}

func (repo *Repository) GetKeystoreManagement() keystoreapi.KeystoreManagement {
	return repo.KeystoreManagement
}

func (repo *Repository) SetKeystoreManagement(instance keystoreapi.KeystoreManagement) {
	repo.KeystoreManagement = instance
}

func (repo *Repository) Clear() {
	repo.KeystoreManagement = nil
}
