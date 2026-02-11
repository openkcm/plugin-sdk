package keystore_management

import (
	"github.com/openkcm/plugin-sdk/api/service/keystoremanagement"
)

type Repository struct {
	KeystoreManagements map[string]keystoremanagement.KeystoreManagement
}

func (repo *Repository) GetKeystoreManagements() map[string]keystoremanagement.KeystoreManagement {
	return repo.KeystoreManagements
}

func (repo *Repository) ListKeystoreManagement() []keystoremanagement.KeystoreManagement {
	list := make([]keystoremanagement.KeystoreManagement, 0, len(repo.KeystoreManagements))
	for _, management := range repo.KeystoreManagements {
		list = append(list, management)
	}
	return list
}

func (repo *Repository) AddKeystoreManagement(instance keystoremanagement.KeystoreManagement) {
	repo.KeystoreManagements[instance.ServiceInfo().Name()] = instance
}

func (repo *Repository) Clear() {
	repo.KeystoreManagements = nil
}
