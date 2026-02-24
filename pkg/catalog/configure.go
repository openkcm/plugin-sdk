package catalog

import (
	"context"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"io"
	"log/slog"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/openkcm/plugin-sdk/api"
	"github.com/openkcm/plugin-sdk/pkg/plugin"
	configv1 "github.com/openkcm/plugin-sdk/proto/service/common/config/v1"
)

type Configurer interface {
	Configure(ctx context.Context, configuration string) error
}

type ConfigurerFunc func(ctx context.Context, configuration string) error

func (fn ConfigurerFunc) Configure(ctx context.Context, configuration string) error {
	return fn(ctx, configuration)
}

func ConfigurePlugin(ctx context.Context, configurer Configurer, dataSource DataSource, lastHash string) (string, error) {
	data, err := dataSource.Load()
	if err != nil {
		return "", fmt.Errorf("failed to load plugin data: %w", err)
	}

	dataHash := hashData(data)
	if lastHash == "" || dataHash != lastHash {
		if err := configurer.Configure(ctx, data); err != nil {
			return "", err
		}
	}
	return dataHash, nil
}

func ReconfigureTask(log *slog.Logger, reconfigurer Reconfigurer) func(context.Context) error {
	return func(ctx context.Context) error {
		return ReconfigureOnSignal(ctx, log, reconfigurer)
	}
}

type Reconfigurer interface {
	Reconfigure(ctx context.Context)
}

type Reconfigurers []Reconfigurer

func (rs Reconfigurers) Reconfigure(ctx context.Context) {
	for _, r := range rs {
		r.Reconfigure(ctx)
	}
}

type Reconfigurable struct {
	Log        *slog.Logger
	Configurer Configurer
	DataSource DataSource
	LastHash   string
}

func (r *Reconfigurable) Reconfigure(ctx context.Context) {
	if r.DataSource == nil {
		return
	}

	if dataHash, err := ConfigurePlugin(ctx, r.Configurer, r.DataSource, r.LastHash); err != nil {
		r.Log.Error("Failed to reconfigure plugin", "error", err)
	} else if dataHash == r.LastHash {
		r.Log.With("hash", r.LastHash).Info("Plugin not reconfigured since the config is unchanged")
	} else {
		r.Log.With("old_hash", r.LastHash).With("new_hash", dataHash).Info("Plugin reconfigured")
		r.LastHash = dataHash
	}
}

func configurePlugin(ctx context.Context, pluginLog *slog.Logger, configurer Configurer, dataSource DataSource) (Reconfigurer, error) {
	switch {
	case configurer == nil && dataSource == nil:
		// The plugin doesn't support configuration and no data source was configured. Nothing to do.
		return nil, nil
	case configurer == nil && dataSource != nil:
		// The plugin does not support configuration but a data source was configured. Nothing to do.
		return nil, nil
	case configurer != nil && dataSource == nil:
		// The plugin supports configuration but no data source was configured. Default to an empty, fixed configuration.
		dataSource = FixedData("")
	case configurer != nil && dataSource != nil:
		// The plugin supports configuration and there was a data source.
	}

	dataHash, err := ConfigurePlugin(ctx, configurer, dataSource, "")
	if err != nil {
		return nil, err
	}

	if !dataSource.IsDynamic() {
		pluginLog.With("reconfigurable", false).Info("Configured plugin")

		return &Reconfigurable{
			Log:        pluginLog,
			Configurer: configurer,
		}, nil
	}

	pluginLog.With("reconfigurable", true).With("hash", dataHash).Info("Configured plugin")
	return &Reconfigurable{
		Log:        pluginLog,
		Configurer: configurer,
		DataSource: dataSource,
		LastHash:   dataHash,
	}, nil
}

type configurerRepo struct {
	configurer Configurer
}

func (repo *configurerRepo) Binder() any {
	return func(configurer Configurer) {
		repo.configurer = configurer
	}
}

func (repo *configurerRepo) Versions() []api.Version {
	return []api.Version{
		configurerV1Version{},
	}
}

func (repo *configurerRepo) Clear() {
	// This function is only for conforming to the Repo interface and isn't
	// expected to be called, but just in case, we'll do the right thing
	// and clear out the configurer that has been bound.
	repo.configurer = nil
}

type configurerV1Version struct{}

func (configurerV1Version) New() api.Facade  { return new(configurerV1) }
func (configurerV1Version) Deprecated() bool { return false }

type configurerV1 struct {
	plugin.Facade

	configv1.ConfigServiceClient

	metadata map[string]string
}

var _ Configurer = (*configurerV1)(nil)

func (v1 *configurerV1) InitInfo(api.Info) {
}

func (v1 *configurerV1) InitLog(*slog.Logger) {
}

func (v1 *configurerV1) Version() uint {
	return 1
}

func (v1 *configurerV1) GetMetadataByKey(key string) any {
	if metadata, ok := v1.metadata[key]; ok {
		return metadata
	}
	return nil
}

func (v1 *configurerV1) Configure(ctx context.Context, data string) error {
	resp, err := v1.ConfigServiceClient.Configure(ctx, &configv1.ConfigureRequest{
		YamlConfiguration: data,
	})
	switch status.Code(err) {
	case codes.OK:
		if v1.metadata == nil {
			v1.metadata = make(map[string]string)
		}
		v1.metadata[BuildInfoMetadata] = extractBuildInfo(resp)
	}

	return err
}

func hashData(data string) string {
	h := sha512.New()
	_, _ = io.Copy(h, strings.NewReader(data))
	return hex.EncodeToString(h.Sum(nil)[:16])
}

func extractBuildInfo(resp *configv1.ConfigureResponse) string {
	defer func() {
		_ = recover()
	}()

	if resp == nil {
		return ""
	}

	return strings.TrimSpace(resp.GetBuildInfo())
}
