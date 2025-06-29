package interfaces

type Service interface {
	Close() (err error)
}

type DirService interface {
	SetDirectory(dir string)
	GetDirectory() string
}

type MetadataService[T Model] interface {
	Service

	DownloadAndReturn(owner string, repository string, assetName string, regex bool, filename string) (data []T, err error)
}

type S3Service[T Model] interface {
	Service

	Download(bucket string, prefix string) (downloaded []string, err error)
	DownloadAndReturnData(bucket string, prefix string) (data []T, err error)
}
