package system_information

import (
	"github.com/openkcm/plugin-sdk/api/service/systeminformation"
)

type Repository struct {
	SystemInformation systeminformation.SystemInformation
}

func (repo *Repository) GetSystemInformation() systeminformation.SystemInformation {
	return repo.SystemInformation
}

func (repo *Repository) SetSystemInformation(sysinfo systeminformation.SystemInformation) {
	repo.SystemInformation = sysinfo
}

func (repo *Repository) Clear() {
	repo.SystemInformation = nil
}
