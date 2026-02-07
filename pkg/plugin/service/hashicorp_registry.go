package service

import (
	apiservice "github.com/openkcm/plugin-sdk/api/service"
	"github.com/openkcm/plugin-sdk/pkg/catalog"
)

var _ apiservice.Registry = (*hashicorpPluginServiceRegistry)(nil)

type hashicorpPluginServiceRegistry struct {
	catalog *catalog.Catalog
}

func NewHashicorpPluginServiceRegistry(catalog *catalog.Catalog) apiservice.Registry {
	return &hashicorpPluginServiceRegistry{
		catalog: catalog,
	}
}
