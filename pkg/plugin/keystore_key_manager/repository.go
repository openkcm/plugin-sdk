package keystore_key_manager

import (
	keystoreapi "github.com/openkcm/plugin-sdk/api/service/keystore"
)

type Repository struct {
	KeyManagers map[string]keystoreapi.KeyManager
}

func (repo *Repository) GetKeystoreKeyManagers() map[string]keystoreapi.KeyManager {
	return repo.KeyManagers
}

func (repo *Repository) ListKeystoreKeyManager() []keystoreapi.KeyManager {
	list := make([]keystoreapi.KeyManager, 0, len(repo.KeyManagers))
	for _, manager := range repo.KeyManagers {
		list = append(list, manager)
	}
	return list
}

func (repo *Repository) AddKeystoreKeyManager(instance keystoreapi.KeyManager) {
	repo.KeyManagers[instance.ServiceInfo().Name()] = instance
}

func (repo *Repository) Clear() {
	repo.KeyManagers = nil
}
