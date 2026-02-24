package catalog

const (
	BuildInfoMetadata = "BuildInfo"
)

type MetadataRetriever interface {
	GetMetadataByKey(key string) any
}
