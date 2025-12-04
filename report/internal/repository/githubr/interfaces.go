package githubr

import (
	"context"
	"io"
	"net/http"

	"github.com/google/go-github/v77/github"
)

//---- REPOSITORY OWNERSHIP

// RepositoryOwnership interface exposes all the functions needed to fetch
// list of owners for a repository based on both the CODEOWNER files and
// the attached team data
type RepositoryOwnership interface {
	RepositoryTeamList
	RepositoryCodeOwnerFiles
	RepositoryOwnerGetter
}

// RepositoryOwnerGetter interface exposes the func that will fetch and merge teams and the
// contents of codeowners as the overall oweners
type RepositoryOwnerGetter interface {
	GetRepositoryOwners(client ClientRepositoryOwnership, repo *github.Repository, options *GetRepositoryOwnerOptions) (owners []string, err error)
}

// RepositoryTeamList interface exposes functions to get list of teams
// that are atttached to the repository and then filter them in required
type RepositoryTeamList interface {
	GetTeamsForRepository(client ClientRepositoryTeamList, repo *github.Repository, options *GetRepositoryOwnerOptions) (teams []*github.Team, err error)
}

// RepositoryCodeOwnerFiles interface for repo that provides method
// to fetch codeowner file content as a set of strings
type RepositoryCodeOwnerFiles interface {
	GetCodeOwnersForRepository(client ClientRepositoryCodeOwnerDownload, repo *github.Repository, options *GetRepositoryOwnerOptions) (owners []string, err error)
}

// ClientRepositoryOwnership is used to fetch ownership via both team and
// codeowner files
// Note: wrapper of github.RepositoriesService
type ClientRepositoryOwnership interface {
	ClientRepositoryTeamList
	ClientRepositoryCodeOwnerDownload
}

// ClientRepositoryTeamList is used to fetch list of teams for the repository
// Note: wrapper of github.RepositoriesService.ListTeams
type ClientRepositoryTeamList interface {
	ListTeams(ctx context.Context, owner, repo string, opts *github.ListOptions) ([]*github.Team, *github.Response, error)
}

// ClientRepositoryCodeOwnerDownload is used to fetch content of the CODEOWNER files
// for the repository
// Note: wrapper of github.RepositoriesService.DownloadContents
type ClientRepositoryCodeOwnerDownload interface {
	DownloadContents(ctx context.Context, owner, repo, filepath string, opts *github.RepositoryContentGetOptions) (io.ReadCloser, *github.Response, error)
}

//---- REPOSITORY LISTS

// RepositoryListByTeam exposes an interface for how
// we fetch list of repositories based on the org and team
type RepositoryListByTeam interface {
	GetRepositoriesForTeam(client ClientTeamListRepositories, owner string, team string, options *GetRepositoriesForTeamOptions) (repositories []*github.Repository, err error)
}

// ClientTeamListRepositories
//
// note: func from *github.TeamsService
// interface for github.Team.ListTeamReposBySlug
type ClientTeamListRepositories interface {
	// ListTeamReposBySlug returns a list of all repositories the team within the org pass can see and match the options
	// passed.
	//
	// Currently filters by Archived only
	ListTeamReposBySlug(ctx context.Context, org string, team string, opts *github.ListOptions) ([]*github.Repository, *github.Response, error)
}

//---- REPOSITORY RELEASES

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
// client is interface for github.Client.Repositories.ListReleases
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
	ListReleases(ctx context.Context, owner string, repo string, opts *github.ListOptions) ([]*github.RepositoryRelease, *github.Response, error)
}

// ClientRepositoryReleaseDownloadReleaseAsset
//
// interface for github.Client.Repositories.ListReleases
type ClientRepositoryReleaseDownloadReleaseAsset interface {
	// DownloadReleaseAsset fetches the remote asset determined by <id> and returns an io pointer to its content
	// and any redirect string.
	//
	// Used to fetch a specifc asset from the a specific github release
	DownloadReleaseAsset(ctx context.Context, owner string, repo string, id int64, followRedirectsClient *http.Client) (rc io.ReadCloser, redirectURL string, err error)
}
