package service

import (
	"context"

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

type ServiceRepository struct {
	identityManagementRepository
	certificateIssuerRepository
	notificationRepository
	systemInformationRepository
	keystoreManagementRepository
	keyManagementRepository

	catalog *catalog.Catalog
}

func (repo *ServiceRepository) Plugins() map[string]catalog.PluginRepo {
	return map[string]catalog.PluginRepo{
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

func (repo *ServiceRepository) Services() []catalog.ServiceRepo {
	return nil
}

func (repo *ServiceRepository) Reconfigure(ctx context.Context) {
	repo.catalog.Reconfigure(ctx)
}

func (repo *ServiceRepository) Close() error {
	if repo.catalog == nil {
		return nil
	}

	return repo.catalog.Close()
}

func CreateCatalog(ctx context.Context, config catalog.Config, builtIns ...catalog.BuiltInPlugin) (_ *catalog.Catalog, err error) {
	repo := &ServiceRepository{}

	repo.catalog, err = catalog.BuildCatalog(ctx, config, repo, builtIns...)
	if err != nil {
		return nil, err
	}

	return repo.catalog, nil
}

func WrapAsServiceRepository(c *catalog.Catalog) *ServiceRepository {
	return &ServiceRepository{
		catalog: c,
	}
}
