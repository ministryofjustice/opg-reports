package githubr

import (
	"context"
	"io"
	"net/http"

	"github.com/google/go-github/v74/github"
)

type RepositoryReleases interface {
	RepositoryReleasesGetMany
	RepositoryReleasesGetOne
	RepositoryReleasesDownloadReleaseAsset
}

// RepositoryReleasesGetMany exposes an interface for how this repository
// would return mutliple matching releases based on the options provided
//
// interface for github.Client.Repositories.ListReleases
type RepositoryReleasesGetMany interface {
	GetRepositoryReleases(client ClientRepositoryReleaseListReleases, owner string, repository string, options *GetRepositoryReleaseOptions) (releases []*github.RepositoryRelease, err error)
}

type RepositoryReleasesGetOneDownloader interface {
	RepositoryReleasesGetOne
	RepositoryReleasesDownloadReleaseAsset
}

// RepositoryReleasesGetOne exposes an interface for how this repository
// would return a single matching release based on the options provided.
//
// interface for github.Client.Repositories.ListReleases
type RepositoryReleasesGetOne interface {
	GetRepositoryRelease(client ClientRepositoryReleaseListReleases, owner string, repository string, options *GetRepositoryReleaseOptions) (release *github.RepositoryRelease, err error)
}

// RepositoryReleasesDownloadReleaseAsset exposes an interface for how
// this repositories used to download an asset from a known release
// to the local filesystem
//
// interface for github.Client.Repositories.ListReleases
type RepositoryReleasesDownloadReleaseAsset interface {
	DownloadRepositoryReleaseAsset(client ClientRepositoryReleaseDownloadReleaseAsset, owner string, repository string, release *github.RepositoryRelease, destination string, options *DownloadRepositoryReleaseAssetOptions) (asset *github.ReleaseAsset, path string, err error)
}

type ClientRepositoryReleases interface {
	ClientRepositoryReleaseListReleases
	ClientRepositoryReleaseDownloadReleaseAsset
}

// ClientRepositoryReleaseListReleases
//
// interface for github.Client.Repositories.ListReleases
type ClientRepositoryReleaseListReleases interface {
	ListReleases(ctx context.Context, owner, repo string, opts *github.ListOptions) ([]*github.RepositoryRelease, *github.Response, error)
}

// ClientRepositoryReleaseDownloadReleaseAsset
//
// interface for github.Client.Repositories.ListReleases
type ClientRepositoryReleaseDownloadReleaseAsset interface {
	// DownloadReleaseAsset fetches the remote asset determined by <id> and returns an io pointer to its content
	// and any redirect string.
	//
	// Used to fetch a specifc asset from the a specific github release
	DownloadReleaseAsset(ctx context.Context, owner, repo string, id int64, followRedirectsClient *http.Client) (rc io.ReadCloser, redirectURL string, err error)
}
