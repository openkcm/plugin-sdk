package key_management

import (
	"github.com/openkcm/plugin-sdk/api/service/keymanagement"
)

type Repository struct {
	KeyManagers map[string]keymanagement.KeyManagement
}

func (repo *Repository) GetKeystoreKeyManagers() map[string]keymanagement.KeyManagement {
	return repo.KeyManagers
}

func (repo *Repository) ListKeystoreKeyManager() []keymanagement.KeyManagement {
	list := make([]keymanagement.KeyManagement, 0, len(repo.KeyManagers))
	for _, manager := range repo.KeyManagers {
		list = append(list, manager)
	}
	return list
}

func (repo *Repository) AddKeystoreKeyManager(instance keymanagement.KeyManagement) {
	repo.KeyManagers[instance.ServiceInfo().Name()] = instance
}

func (repo *Repository) Clear() {
	repo.KeyManagers = nil
}
