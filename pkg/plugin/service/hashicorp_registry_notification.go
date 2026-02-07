package service

import (
	"fmt"

	"github.com/openkcm/plugin-sdk/api/service"
	"github.com/openkcm/plugin-sdk/api/service/notification"
	notificationv1 "github.com/openkcm/plugin-sdk/proto/plugin/notification/v1"
)

func (h *hashicorpPluginServiceRegistry) NotificationByName(name string) (notification.Notification, error) {
	plugin := h.catalog.LookupByTypeAndName(notificationv1.Type, name)
	if plugin == nil {
		return nil, fmt.Errorf("unable to find certificate issuer plugin %q", name)
	}
	return NewNotificationV1Plugin(plugin), nil
}

func (h *hashicorpPluginServiceRegistry) NotificationByNameAndVersion(version service.Version, name string) (notification.Notification, error) {
	switch version {
	case service.V1:
		plugin := h.catalog.LookupByTypeAndName(notificationv1.Type, name)
		if plugin == nil {
			return nil, fmt.Errorf("unable to find certificate issuer plugin %q", name)
		}
		return NewNotificationV1Plugin(plugin), nil
	}

	return nil, service.ErrVersionNotSupported
}

func (h *hashicorpPluginServiceRegistry) NotificationByVersion(version service.Version) ([]notification.Notification, error) {
	var issuers []notification.Notification
	switch version {
	case service.V1:
		plugins := h.catalog.LookupByType(notificationv1.Type)
		for _, plugin := range plugins {
			issuers = append(issuers, NewNotificationV1Plugin(plugin))
		}
		return issuers, nil
	}

	return issuers, service.ErrVersionNotSupported
}
