package service

import (
	"io"

	"github.com/openkcm/plugin-sdk/api/service/certificateissuer"
	"github.com/openkcm/plugin-sdk/api/service/identitymanagement"
	"github.com/openkcm/plugin-sdk/api/service/keystore"
	"github.com/openkcm/plugin-sdk/api/service/notification"
	"github.com/openkcm/plugin-sdk/api/service/systeminformation"
)

type Registry interface {
	io.Closer

	GetCertificateIssuer() certificateissuer.CertificateIssuer
	GetNotification() notification.Notification
	GetSystemInformation() systeminformation.SystemInformation
	GetIdentityManagement() identitymanagement.IdentityManagement

	GetKeystoreManagements() map[string]keystore.KeystoreManagement
	ListKeystoreManagement() []keystore.KeystoreManagement

	GetKeystoreKeyManagers() map[string]keystore.KeyManager
	ListKeystoreKeyManager() []keystore.KeyManager
}
