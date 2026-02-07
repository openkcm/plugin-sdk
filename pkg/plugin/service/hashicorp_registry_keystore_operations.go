package service

import (
	"fmt"

	"github.com/openkcm/plugin-sdk/api/service"
	"github.com/openkcm/plugin-sdk/api/service/keystore"
	operationsv1 "github.com/openkcm/plugin-sdk/proto/plugin/keystore/operations/v1"
)

func (h *hashicorpPluginServiceRegistry) KeystoreOperationsByName(name string) (keystore.KeystoreOperations, error) {
	plugin := h.catalog.LookupByTypeAndName(operationsv1.Type, name)
	if plugin == nil {
		return nil, fmt.Errorf("unable to find certificate issuer plugin %q", name)
	}
	return NewKeystoreOperationsV1Plugin(plugin), nil
}

func (h *hashicorpPluginServiceRegistry) KeystoreOperationsByNameAndVersion(version service.Version, name string) (keystore.KeystoreOperations, error) {
	switch version {
	case service.V1:
		plugin := h.catalog.LookupByTypeAndName(operationsv1.Type, name)
		if plugin == nil {
			return nil, fmt.Errorf("unable to find certificate issuer plugin %q", name)
		}
		return NewKeystoreOperationsV1Plugin(plugin), nil
	}

	return nil, service.ErrVersionNotSupported
}

func (h *hashicorpPluginServiceRegistry) KeystoreOperationsByVersion(version service.Version) ([]keystore.KeystoreOperations, error) {
	var issuers []keystore.KeystoreOperations
	switch version {
	case service.V1:
		plugins := h.catalog.LookupByType(operationsv1.Type)
		for _, plugin := range plugins {
			issuers = append(issuers, NewKeystoreOperationsV1Plugin(plugin))
		}
		return issuers, nil
	}

	return issuers, service.ErrVersionNotSupported
}
