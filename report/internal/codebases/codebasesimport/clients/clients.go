package clients

import (
	"context"
	"io"

	"github.com/google/go-github/v84/github"
)

// TeamClient wrapper around *github.TeamsService
type TeamClient interface {
	ListTeamReposBySlug(ctx context.Context, org, slug string, opts *github.ListOptions) ([]*github.Repository, *github.Response, error)
}

// RepoClient wrapper around *github.RepositoriesService
type RepoClient interface {
	// fetch attached teams (*github.RepositoriesService)
	ListTeams(ctx context.Context, owner, repo string, opts *github.ListOptions) ([]*github.Team, *github.Response, error)
	// fetch file content (*github.RepositoriesService)
	DownloadContents(ctx context.Context, owner, repo, filepath string, opts *github.RepositoryContentGetOptions) (io.ReadCloser, *github.Response, error)
	// get contenst method to fetch directory content
	GetContents(ctx context.Context, owner, repo, path string, opts *github.RepositoryContentGetOptions) (fileContent *github.RepositoryContent, directoryContent []*github.RepositoryContent, resp *github.Response, err error)
}
