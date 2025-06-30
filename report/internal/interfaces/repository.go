package interfaces

type Repository interface{}

type S3Repository interface {
	ListBucket(bucket string, prefix string) (fileList []string, err error)
	Download(bucket string, files []string, localDir string) (downloadedFiles []string, err error)
}
