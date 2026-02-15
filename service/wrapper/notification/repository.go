package notification

import (
	"github.com/openkcm/plugin-sdk/service/api/notification"
)

type Repository struct {
	Instance notification.Notification
}

func (repo *Repository) Notification() (notification.Notification, bool) {
	return repo.Instance, repo.Instance != nil
}

func (repo *Repository) SetNotification(instance notification.Notification) {
	repo.Instance = instance
}

func (repo *Repository) Clear() {
	repo.Instance = nil
}
