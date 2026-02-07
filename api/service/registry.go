package service

import (
	"errors"

	"github.com/openkcm/plugin-sdk/api/service/certificateissuer"
	"github.com/openkcm/plugin-sdk/api/service/identitymanagement"
	"github.com/openkcm/plugin-sdk/api/service/keystore"
	"github.com/openkcm/plugin-sdk/api/service/notification"
	"github.com/openkcm/plugin-sdk/api/service/systeminformation"
)

type Version int

const (
	V1 Version = iota + 1
)

var (
	ErrVersionNotSupported = errors.New("version not supported")
)

type Registry interface {
	CertificateIssuer
	IdentityManagement
	KeystoreManagement
	KeystoreOperations
	Notification
	SystemInformation
}

type CertificateIssuer interface {
	CertificateIssuerByName(name string) (certificateissuer.CertificateIssuer, error)
	CertificateIssuerByNameAndVersion(version Version, name string) (certificateissuer.CertificateIssuer, error)
	ListCertificateIssuerByVersion(version Version) ([]certificateissuer.CertificateIssuer, error)
}

type KeystoreManagement interface {
	KeystoreManagementByName(name string) (keystore.KeystoreManagement, error)
	KeystoreManagementByNameAndVersion(version Version, name string) (keystore.KeystoreManagement, error)
	KeystoreManagementByVersion(version Version) ([]keystore.KeystoreManagement, error)
}

type KeystoreOperations interface {
	KeystoreOperationsByName(name string) (keystore.KeystoreOperations, error)
	KeystoreOperationsByNameAndVersion(version Version, name string) (keystore.KeystoreOperations, error)
	KeystoreOperationsByVersion(version Version) ([]keystore.KeystoreOperations, error)
}

type IdentityManagement interface {
	IdentityManagementByName(name string) (identitymanagement.IdentityManagement, error)
	IdentityManagementByNameAndVersion(version Version, name string) (identitymanagement.IdentityManagement, error)
	IdentityManagementByVersion(version Version) ([]identitymanagement.IdentityManagement, error)
}

type Notification interface {
	NotificationByName(name string) (notification.Notification, error)
	NotificationByNameAndVersion(version Version, name string) (notification.Notification, error)
	NotificationByVersion(version Version) ([]notification.Notification, error)
}

type SystemInformation interface {
	SystemInformationByName(name string) (systeminformation.SystemInformation, error)
	SystemInformationByNameAndVersion(version Version, name string) (systeminformation.SystemInformation, error)
	SystemInformationByVersion(version Version) ([]systeminformation.SystemInformation, error)
}
