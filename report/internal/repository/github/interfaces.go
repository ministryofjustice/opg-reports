package github

import (
	"os"

	"github.com/google/go-github/v62/github"
)

// Releaser interface exposes correct methods for finding and
// downloading the release related assets such as tar.gz files
type Releaser interface {
	GetLatestReleaseAsset(organisation string, repositoryName string, assetName string, regex bool) (asset *github.ReleaseAsset, err error)
	DownloadReleaseAsset(organisation string, repositoryName string, assetID int64, destinationFilePath string) (destination *os.File, err error)
}
