package keystore_management

import (
	keystoremanagementapi "github.com/openkcm/plugin-sdk/api/service/keystoremanagement"
)

type Repository struct {
	KeystoreManagements map[string]keystoremanagementapi.KeystoreManagement
}

func (repo *Repository) GetKeystoreManagements() map[string]keystoremanagementapi.KeystoreManagement {
	return repo.KeystoreManagements
}

func (repo *Repository) ListKeystoreManagement() []keystoremanagementapi.KeystoreManagement {
	list := make([]keystoremanagementapi.KeystoreManagement, 0, len(repo.KeystoreManagements))
	for _, management := range repo.KeystoreManagements {
		list = append(list, management)
	}
	return list
}

func (repo *Repository) AddKeystoreManagement(instance keystoremanagementapi.KeystoreManagement) {
	repo.KeystoreManagements[instance.ServiceInfo().Name()] = instance
}

func (repo *Repository) Clear() {
	repo.KeystoreManagements = nil
}
