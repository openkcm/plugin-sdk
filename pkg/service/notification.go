package service

import (
	"github.com/openkcm/plugin-sdk/api"
	"github.com/openkcm/plugin-sdk/pkg/catalog"
	"github.com/openkcm/plugin-sdk/pkg/service/internal/notification"
)

type notificationRepository struct {
	notification.Repository
}

func (repo *notificationRepository) Binder() any {
	return repo.SetNotification
}

func (repo *notificationRepository) Constraints() catalog.Constraints {
	return catalog.ExactlyOne()
}

func (repo *notificationRepository) Versions() []api.Version {
	return []api.Version{notificationV1{}}
}

type notificationV1 struct{}

func (notificationV1) New() api.Facade  { return new(notification.V1) }
func (notificationV1) Deprecated() bool { return false }
