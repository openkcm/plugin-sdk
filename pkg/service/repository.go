package service

import (
	"context"

	"github.com/openkcm/plugin-sdk/api"
	"github.com/openkcm/plugin-sdk/pkg/catalog"
)

const (
	certificateIssuerType = "CertificateIssuer"
	// TODO: value should go away after plugin proto refactoring
	certificateIssuerServiceType = "CertificateIssuerService"

	notificationType = "Notification"
	// TODO: value should go away after plugin proto refactoring
	notificationServiceType = "NotificationService"

	systemInformationType = "SystemInformation"
	// TODO: value should go away after plugin proto refactoring
	systemInformationServiceType = "SystemInformationService"

	identityManagementType = "IdentityManagement"
	// TODO: value should go away after plugin proto refactoring
	identityManagementServiceType = "IdentityManagementService"

	// TODO: value should become `KeystoreManagement` after plugin proto refactoring
	keystoreManagementType = "KeystoreProvider"
	// TODO: value should become `KeyManagement` after plugin proto refactoring
	keyManagementType = "KeystoreInstanceKeyOperation"
)

type Repository struct {
	identityManagementRepository
	certificateIssuerRepository
	notificationRepository
	systemInformationRepository
	keystoreManagementRepository
	keyManagementRepository

	RawCatalog *catalog.Catalog
}

func (repo *Repository) Plugins() map[string]api.PluginRepo {
	return map[string]api.PluginRepo{
		identityManagementType:        &repo.identityManagementRepository,
		identityManagementServiceType: &repo.identityManagementRepository,
		certificateIssuerType:         &repo.certificateIssuerRepository,
		certificateIssuerServiceType:  &repo.certificateIssuerRepository,
		notificationType:              &repo.notificationRepository,
		notificationServiceType:       &repo.notificationRepository,
		systemInformationType:         &repo.systemInformationRepository,
		systemInformationServiceType:  &repo.systemInformationRepository,
		keystoreManagementType:        &repo.keystoreManagementRepository,
		keyManagementType:             &repo.keyManagementRepository,
	}
}

func (repo *Repository) Services() []api.ServiceRepo {
	return nil
}

func (repo *Repository) Reconfigure(ctx context.Context) {
	repo.RawCatalog.Reconfigure(ctx)
}

func (repo *Repository) Close() error {
	if repo.RawCatalog == nil {
		return nil
	}

	return repo.RawCatalog.Close()
}

func CreateServiceRepository(ctx context.Context, config catalog.Config, builtIns ...catalog.BuiltInPlugin) (_ *Repository, err error) {
	repo := &Repository{}

	repo.RawCatalog, err = catalog.BuildCatalog(ctx, config, repo, builtIns...)
	if err != nil {
		return nil, err
	}

	return repo, nil
}
