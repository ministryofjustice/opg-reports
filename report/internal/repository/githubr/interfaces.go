package githubr

import (
	"context"
	"io"
	"net/http"
	"os"

	"github.com/google/go-github/v62/github"
)

type Model interface{}

// Releaser interface exposes correct methods for finding and
// downloading the release related assets such as tar.gz files
type Releaser interface {
	GetLatestReleaseAsset(client ReleaseClient, organisation string, repositoryName string, assetName string, regex bool) (asset *github.ReleaseAsset, err error)
	DownloadReleaseAsset(client ReleaseClient, organisation string, repositoryName string, assetID int64, destinationFilePath string) (destination *os.File, err error)
	DownloadReleaseAssetByName(client ReleaseClient, organisation string, repositoryName string, assetName string, regex bool, directory string) (downloadedTo string, err error)
}

type ReleaseClient interface {
	ListReleases(ctx context.Context, owner, repo string, opts *github.ListOptions) ([]*github.RepositoryRelease, *github.Response, error)
	GetLatestRelease(ctx context.Context, owner, repo string) (*github.RepositoryRelease, *github.Response, error)
	DownloadReleaseAsset(ctx context.Context, owner, repo string, id int64, followRedirectsClient *http.Client) (rc io.ReadCloser, redirectURL string, err error)
}
