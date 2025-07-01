package githubr

import (
	"context"
	"io"
	"net/http"
	"os"

	"github.com/google/go-github/v62/github"
)

// ReleaseRepository interface exposes the methods used to
// get information about releases from the github api for
// a particular repository.
//
// Uses a series fo client interfaces (`ClientRelease*`) to
// require the appropriate methods for for each method which
// can then be mocked in testing easier to check capability
// without calling the api
type ReleaseRepository interface {
	GetAllReleases(client ClientReleaseLister, organisation string, repositoryName string) (releases []*github.RepositoryRelease, err error)
	GetReleases(client ClientReleaseLister, organisation string, repositoryName string, options *GetReleaseOptions) (releases []*github.RepositoryRelease, err error)
	GetLatestReleaseAsset(client ClientReleaseGetter, organisation string, repositoryName string, assetName string, regex bool) (asset *github.ReleaseAsset, err error)
	DownloadReleaseAsset(client ClientReleaseDownloader, organisation string, repositoryName string, assetID int64, destinationFilePath string) (destination *os.File, err error)
	DownloadReleaseAssetByName(client ReleaseGetAndDownloader, organisation string, repositoryName string, assetName string, regex bool, directory string) (downloadedTo string, err error)
}

// ClientReleaseLister represents the methods required for the
// github client to list all releases and is used by:
//
//   - GetAllReleases
//   - GetReleases
type ClientReleaseLister interface {
	ListReleases(ctx context.Context, owner, repo string, opts *github.ListOptions) ([]*github.RepositoryRelease, *github.Response, error)
}

// ClientReleaseGetter represents the methods used for github
// client to get the latest release. Used by
//
//   - GetLatestReleaseAsset
type ClientReleaseGetter interface {
	GetLatestRelease(ctx context.Context, owner, repo string) (*github.RepositoryRelease, *github.Response, error)
}

// ClientReleaseDownloader represents the methods used for github
// client to download a release. Used by
//
//   - DownloadReleaseAsset
type ClientReleaseDownloader interface {
	// DownloadReleaseAsset fetches the remote asset determined by <id> and returns a file pointer to its location
	// and any redirect string.
	//
	// Used to fetch a specifc asset from the a specific github release
	DownloadReleaseAsset(ctx context.Context, owner, repo string, id int64, followRedirectsClient *http.Client) (rc io.ReadCloser, redirectURL string, err error)
}

// ReleaseGetAndDownloader represents the methods used for github
// client to get and download a release
//
//   - ReleaseGetAndDownloader
type ReleaseGetAndDownloader interface {
	ClientReleaseGetter
	ClientReleaseDownloader
}

// ReleaseClient is the merged interface that represents all
// the methods used for release related github clients
type ReleaseClient interface {
	ClientReleaseLister
	ClientReleaseGetter
	ClientReleaseDownloader
}
