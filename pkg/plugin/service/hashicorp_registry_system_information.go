package service

import (
	"fmt"

	"github.com/openkcm/plugin-sdk/api/service"
	"github.com/openkcm/plugin-sdk/api/service/systeminformation"
	systeminformationv1 "github.com/openkcm/plugin-sdk/proto/plugin/systeminformation/v1"
)

func (h *hashicorpPluginServiceRegistry) SystemInformationByName(name string) (systeminformation.SystemInformation, error) {
	plugin := h.catalog.LookupByTypeAndName(systeminformationv1.Type, name)
	if plugin == nil {
		return nil, fmt.Errorf("unable to find certificate issuer plugin %q", name)
	}
	return NewSystemInformationV1Plugin(plugin), nil
}

func (h *hashicorpPluginServiceRegistry) SystemInformationByNameAndVersion(version service.Version, name string) (systeminformation.SystemInformation, error) {
	switch version {
	case service.V1:
		plugin := h.catalog.LookupByTypeAndName(systeminformationv1.Type, name)
		if plugin == nil {
			return nil, fmt.Errorf("unable to find certificate issuer plugin %q", name)
		}
		return NewSystemInformationV1Plugin(plugin), nil
	}

	return nil, service.ErrVersionNotSupported
}

func (h *hashicorpPluginServiceRegistry) SystemInformationByVersion(version service.Version) ([]systeminformation.SystemInformation, error) {
	var issuers []systeminformation.SystemInformation
	switch version {
	case service.V1:
		plugins := h.catalog.LookupByType(systeminformationv1.Type)
		for _, plugin := range plugins {
			issuers = append(issuers, NewSystemInformationV1Plugin(plugin))
		}
		return issuers, nil
	}

	return issuers, service.ErrVersionNotSupported
}
