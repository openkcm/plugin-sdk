package service

import (
	"fmt"

	"github.com/openkcm/plugin-sdk/api/service"
	"github.com/openkcm/plugin-sdk/api/service/keystore"
	managementv1 "github.com/openkcm/plugin-sdk/proto/plugin/keystore/management/v1"
)

func (h *hashicorpPluginServiceRegistry) KeystoreManagementByName(name string) (keystore.KeystoreManagement, error) {
	plugin := h.catalog.LookupByTypeAndName(managementv1.Type, name)
	if plugin == nil {
		return nil, fmt.Errorf("unable to find certificate issuer plugin %q", name)
	}
	return NewKeystoreManagementV1Plugin(plugin), nil
}

func (h *hashicorpPluginServiceRegistry) KeystoreManagementByNameAndVersion(version service.Version, name string) (keystore.KeystoreManagement, error) {
	switch version {
	case service.V1:
		plugin := h.catalog.LookupByTypeAndName(managementv1.Type, name)
		if plugin == nil {
			return nil, fmt.Errorf("unable to find certificate issuer plugin %q", name)
		}
		return NewKeystoreManagementV1Plugin(plugin), nil
	}

	return nil, service.ErrVersionNotSupported
}

func (h *hashicorpPluginServiceRegistry) KeystoreManagementByVersion(version service.Version) ([]keystore.KeystoreManagement, error) {
	var issuers []keystore.KeystoreManagement
	switch version {
	case service.V1:
		plugins := h.catalog.LookupByType(managementv1.Type)
		for _, plugin := range plugins {
			issuers = append(issuers, NewKeystoreManagementV1Plugin(plugin))
		}
		return issuers, nil
	}

	return issuers, service.ErrVersionNotSupported
}
