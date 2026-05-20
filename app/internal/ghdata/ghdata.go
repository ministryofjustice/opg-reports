package ghdata

import (
	"opg-reports/app/internal/ghdata/teamrepositories"

	"github.com/google/go-github/v87/github"
)

type Client interface {
	*github.TeamsService | *github.PullRequestsService
}

type Result interface {
	*github.Repository | *github.PullRequest
}

type GitHubData[C Client, R Result] interface {
	//
	GetData() (results []R, err error)
}

var _ GitHubData[*github.TeamsService, *github.Repository] = &teamrepositories.Source[*github.TeamsService, *github.Repository]{}
