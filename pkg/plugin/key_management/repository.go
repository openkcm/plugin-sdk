package key_management

import (
	"github.com/openkcm/plugin-sdk/api/service/keymanagement"
)

type Repository struct {
	KeyManagements map[string]keymanagement.KeyManagement
}

func (repo *Repository) GetKeyManagements() map[string]keymanagement.KeyManagement {
	return repo.KeyManagements
}

func (repo *Repository) ListKeyManagement() []keymanagement.KeyManagement {
	list := make([]keymanagement.KeyManagement, 0, len(repo.KeyManagements))
	for _, manager := range repo.KeyManagements {
		list = append(list, manager)
	}
	return list
}

func (repo *Repository) AddKeyManagement(instance keymanagement.KeyManagement) {
	repo.KeyManagements[instance.ServiceInfo().Name()] = instance
}

func (repo *Repository) Clear() {
	repo.KeyManagements = nil
}
