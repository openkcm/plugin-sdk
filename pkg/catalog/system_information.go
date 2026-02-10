package catalog

import (
	"github.com/openkcm/plugin-sdk/api"
	"github.com/openkcm/plugin-sdk/pkg/plugin/system_information"
)

type systemInformationRepository struct {
	system_information.Repository
}

func (repo *systemInformationRepository) Binder() any {
	return repo.SetSystemInformation
}

func (repo *systemInformationRepository) Constraints() Constraints {
	return ExactlyOne()
}

func (repo *systemInformationRepository) Versions() []api.Version {
	return []api.Version{systemInformationV1{}}
}

func (repo *systemInformationRepository) BuiltIns() []BuiltInPlugin {
	return []BuiltInPlugin{}
}

type systemInformationV1 struct{}

func (systemInformationV1) New() api.Facade  { return new(system_information.V1) }
func (systemInformationV1) Deprecated() bool { return false }
