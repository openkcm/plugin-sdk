package service

import (
	"fmt"

	"github.com/openkcm/plugin-sdk/api/service"
	"github.com/openkcm/plugin-sdk/api/service/identitymanagement"
	identity_managementv1 "github.com/openkcm/plugin-sdk/proto/plugin/identity_management/v1"
)

func (h *hashicorpPluginServiceRegistry) IdentityManagementByName(name string) (identitymanagement.IdentityManagement, error) {
	plugin := h.catalog.LookupByTypeAndName(identity_managementv1.Type, name)
	if plugin == nil {
		return nil, fmt.Errorf("unable to find certificate issuer plugin %q", name)
	}
	return NewIdentityManagementV1Plugin(plugin), nil
}

func (h *hashicorpPluginServiceRegistry) IdentityManagementByNameAndVersion(version service.Version, name string) (identitymanagement.IdentityManagement, error) {
	switch version {
	case service.V1:
		plugin := h.catalog.LookupByTypeAndName(identity_managementv1.Type, name)
		if plugin == nil {
			return nil, fmt.Errorf("unable to find certificate issuer plugin %q", name)
		}
		return NewIdentityManagementV1Plugin(plugin), nil
	}

	return nil, service.ErrVersionNotSupported
}

func (h *hashicorpPluginServiceRegistry) IdentityManagementByVersion(version service.Version) ([]identitymanagement.IdentityManagement, error) {
	var issuers []identitymanagement.IdentityManagement
	switch version {
	case service.V1:
		plugins := h.catalog.LookupByType(identity_managementv1.Type)
		for _, plugin := range plugins {
			issuers = append(issuers, NewIdentityManagementV1Plugin(plugin))
		}
		return issuers, nil
	}

	return issuers, service.ErrVersionNotSupported
}
