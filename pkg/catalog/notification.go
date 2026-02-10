package catalog

import (
	"github.com/openkcm/plugin-sdk/api"
	"github.com/openkcm/plugin-sdk/pkg/plugin/certificate_issuer"
	"github.com/openkcm/plugin-sdk/pkg/plugin/notification"
)

type notificationRepository struct {
	notification.Repository
}

func (repo *notificationRepository) Binder() any {
	return repo.SetNotification
}

func (repo *notificationRepository) Constraints() Constraints {
	return ExactlyOne()
}

func (repo *notificationRepository) Versions() []api.Version {
	return []api.Version{notificationV1{}}
}

func (repo *notificationRepository) BuiltIns() []BuiltInPlugin {
	return []BuiltInPlugin{}
}

type notificationV1 struct{}

func (notificationV1) New() api.Facade  { return new(certificate_issuer.V1) }
func (notificationV1) Deprecated() bool { return false }
