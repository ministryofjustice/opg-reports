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

	DownloadAndReturn(owner string, repository string, assetName string, filename string) (data []T, err error)
}
