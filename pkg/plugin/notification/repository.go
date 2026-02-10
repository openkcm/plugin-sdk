package notification

import (
	notificationapi "github.com/openkcm/plugin-sdk/api/service/notification"
)

type Repository struct {
	Notification notificationapi.Notification
}

func (repo *Repository) GetNotification() notificationapi.Notification {
	return repo.Notification
}

func (repo *Repository) SetNotification(instance notificationapi.Notification) {
	repo.Notification = instance
}

func (repo *Repository) Clear() {
	repo.Notification = nil
}
