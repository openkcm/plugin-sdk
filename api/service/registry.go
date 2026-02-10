package service

import (
	"io"

	"github.com/openkcm/plugin-sdk/api/service/certificateissuer"
	"github.com/openkcm/plugin-sdk/api/service/keystore"
	"github.com/openkcm/plugin-sdk/api/service/notification"
	"github.com/openkcm/plugin-sdk/api/service/systeminformation"
)

type Catalog interface {
	io.Closer

	GetCertificateIssuer() certificateissuer.CertificateIssuer
	GetNotification() notification.Notification
	GetSystemInformation() systeminformation.SystemInformation
	GetKeystoreManagement() keystore.KeystoreManagement
	GetKeystoreInstanceKeyOperations() keystore.KeystoreInstanceKeyOperations
}
