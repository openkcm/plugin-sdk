package keystore_management

import (
	keystoreapi "github.com/openkcm/plugin-sdk/api/service/keystore"
)

type Repository struct {
	KeystoreManagements map[string]keystoreapi.KeystoreManagement
}

func (repo *Repository) GetKeystoreManagements() map[string]keystoreapi.KeystoreManagement {
	return repo.KeystoreManagements
}

func (repo *Repository) ListKeystoreManagement() []keystoreapi.KeystoreManagement {
	list := make([]keystoreapi.KeystoreManagement, 0, len(repo.KeystoreManagements))
	for _, management := range repo.KeystoreManagements {
		list = append(list, management)
	}
	return list
}

func (repo *Repository) AddKeystoreManagement(instance keystoreapi.KeystoreManagement) {
	repo.KeystoreManagements[instance.ServiceInfo().Name()] = instance
}

func (repo *Repository) Clear() {
	repo.KeystoreManagements = nil
}
