package key_management

import (
	"log/slog"

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
	if repo.KeyManagements == nil {
		repo.KeyManagements = make(map[string]keymanagementapi.KeyManagement)
	}

	info := instance.ServiceInfo()
	if info == nil {
		slog.Error("FATAL:Service info of KeyManagement is required!")
		return
	}

	repo.KeyManagements[info.Name()] = instance
}

func (repo *Repository) Clear() {
	repo.KeyManagements = nil
}
