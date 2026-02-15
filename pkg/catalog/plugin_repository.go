package catalog

import (
	"context"
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

type PluginRepository struct {
	identityManagementRepository
	certificateIssuerRepository
	notificationRepository
	systemInformationRepository
	keystoreManagementRepository
	keyManagementRepository

	catalog *Catalog
}

func (repo *PluginRepository) Plugins() map[string]PluginRepo {
	return map[string]PluginRepo{
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

func (repo *PluginRepository) Services() []ServiceRepo {
	return nil
}

func (repo *PluginRepository) Reconfigure(ctx context.Context) {
	repo.catalog.Reconfigure(ctx)
}

func (repo *PluginRepository) Close() error {
	if repo.catalog == nil {
		return nil
	}

	err := repo.catalog.Close()
	if err != nil {
		return err
	}

	return nil
}

func CreatePluginRepository(ctx context.Context, config Config, builtInPlugins ...BuiltInPlugin) (_ *PluginRepository, err error) {
	repo := &PluginRepository{}

	repo.catalog, err = buildCatalog(ctx, config, repo, builtInPlugins...)
	if err != nil {
		return nil, err
	}

	return repo, nil
}

func WrapAsPluginRepository(c *Catalog) *PluginRepository {
	return &PluginRepository{
		catalog: c,
	}
}
