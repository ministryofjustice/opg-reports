package ghdata

import (
	"opg-reports/app/internal/ghdata/ghmergedprs"
	"opg-reports/app/internal/ghdata/ghteamrepositories"
	"opg-reports/app/internal/ghdata/ghworkflowruns"

	"github.com/google/go-github/v87/github"
)

var (
	_ GitHubData[*github.TeamsService, *github.Repository]             = &ghteamrepositories.Source[*github.TeamsService, *github.Repository]{}       // this would get all repositories
	_ GitHubData[*github.ActionsService, *ghworkflowruns.ResultData]   = &ghworkflowruns.Source[*github.ActionsService, *ghworkflowruns.ResultData]{} // this would get all workflow runs for a list of repos
	_ GitHubData[*github.PullRequestsService, *ghmergedprs.ResultData] = &ghmergedprs.Source[*github.PullRequestsService, *ghmergedprs.ResultData]{}  // this would get all pull requests for a list of repos
)
