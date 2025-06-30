package interfaces

import (
	"os"

	"github.com/google/go-github/v62/github"
)

type Repository interface{}

type GithubRepository interface {
	GetLatestReleaseAsset(organisation string, repositoryName string, assetName string, regex bool) (asset *github.ReleaseAsset, err error)
	DownloadReleaseAsset(organisation string, repositoryName string, assetID int64, destinationFilePath string) (destination *os.File, err error)
}

type S3Repository interface {
	ListBucket(bucket string, prefix string) (fileList []string, err error)
	Download(bucket string, files []string, localDir string) (downloadedFiles []string, err error)
}

type STSRepository interface {
	GetAccountID() (accountID string)
}
