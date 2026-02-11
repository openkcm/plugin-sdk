package key_management

import (
	keymanagementapi "github.com/openkcm/plugin-sdk/api/service/keymanagement"
)

type Repository struct {
	KeyManagements map[string]keymanagementapi.KeyManagement
}

func (repo *Repository) GetKeyManagements() map[string]keymanagementapi.KeyManagement {
	return repo.KeyManagements
}

func (repo *Repository) ListKeyManagement() []keymanagementapi.KeyManagement {
	list := make([]keymanagementapi.KeyManagement, 0, len(repo.KeyManagements))
	for _, manager := range repo.KeyManagements {
		list = append(list, manager)
	}
	return list
}

func (repo *Repository) AddKeyManagement(instance keymanagementapi.KeyManagement) {
	repo.KeyManagements[instance.ServiceInfo().Name()] = instance
}

func (repo *Repository) Clear() {
	repo.KeyManagements = nil
}
