package system_information

import (
	systeminformationapi "github.com/openkcm/plugin-sdk/api/service/systeminformation"
)

type Repository struct {
	SystemInformation systeminformationapi.SystemInformation
}

func (repo *Repository) GetSystemInformation() systeminformationapi.SystemInformation {
	return repo.SystemInformation
}

func (repo *Repository) SetSystemInformation(sysinfo systeminformationapi.SystemInformation) {
	repo.SystemInformation = sysinfo
}

func (repo *Repository) Clear() {
	repo.SystemInformation = nil
}
