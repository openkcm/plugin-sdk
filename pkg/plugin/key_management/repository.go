package key_management

import (
	keymanagementapi "github.com/openkcm/plugin-sdk/api/service/keymanagement"
)

type Repository struct {
	KeyManagers map[string]keymanagementapi.KeyManagement
}

func (repo *Repository) GetKeystoreKeyManagers() map[string]keymanagementapi.KeyManagement {
	return repo.KeyManagers
}

func (repo *Repository) ListKeystoreKeyManager() []keymanagementapi.KeyManagement {
	list := make([]keymanagementapi.KeyManagement, 0, len(repo.KeyManagers))
	for _, manager := range repo.KeyManagers {
		list = append(list, manager)
	}
	return list
}

func (repo *Repository) AddKeystoreKeyManager(instance keymanagementapi.KeyManagement) {
	repo.KeyManagers[instance.ServiceInfo().Name()] = instance
}

func (repo *Repository) Clear() {
	repo.KeyManagers = nil
}
