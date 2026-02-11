package notification

import (
	"github.com/openkcm/plugin-sdk/api/service/notification"
)

type Repository struct {
	Notification notification.Notification
}

func (repo *Repository) GetNotification() notification.Notification {
	return repo.Notification
}

func (repo *Repository) SetNotification(instance notification.Notification) {
	repo.Notification = instance
}

func (repo *Repository) Clear() {
	repo.Notification = nil
}
