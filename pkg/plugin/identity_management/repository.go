package identity_management

import (
	"github.com/openkcm/plugin-sdk/api/service/identitymanagement"
)

type Repository struct {
	IdentityManagement identitymanagement.IdentityManagement
}

func (repo *Repository) GetIdentityManagement() identitymanagement.IdentityManagement {
	return repo.IdentityManagement
}

func (repo *Repository) SetIdentityManagement(instance identitymanagement.IdentityManagement) {
	repo.IdentityManagement = instance
}

func (repo *Repository) Clear() {
	repo.IdentityManagement = nil
}
