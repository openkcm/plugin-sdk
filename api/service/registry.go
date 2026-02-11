package serviceapi

import (
	"io"

	"github.com/openkcm/plugin-sdk/api/service/certificateissuer"
	"github.com/openkcm/plugin-sdk/api/service/identitymanagement"
	"github.com/openkcm/plugin-sdk/api/service/keymanagement"
	"github.com/openkcm/plugin-sdk/api/service/keystoremanagement"
	"github.com/openkcm/plugin-sdk/api/service/notification"
	"github.com/openkcm/plugin-sdk/api/service/systeminformation"
)

type Registry interface {
	io.Closer

	GetCertificateIssuer() certificateissuer.CertificateIssuer
	GetNotification() notification.Notification
	GetSystemInformation() systeminformation.SystemInformation
	GetIdentityManagement() identitymanagement.IdentityManagement

	GetKeystoreManagements() map[string]keystoremanagement.KeystoreManagement
	ListKeystoreManagement() []keystoremanagement.KeystoreManagement

	GetKeystoreKeyManagers() map[string]keymanagement.KeyManagement
	ListKeystoreKeyManager() []keymanagement.KeyManagement
}
