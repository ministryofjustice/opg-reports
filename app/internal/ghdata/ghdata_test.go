package ghdata

import (
	"opg-reports/app/internal/ghdata/ghteamrepositories"
	"opg-reports/app/internal/ghdata/ghworkflowruns"

	"github.com/google/go-github/v87/github"
)

var (
	_ GitHubData[*github.TeamsService, *github.Repository]    = &ghteamrepositories.Source[*github.TeamsService, *github.Repository]{}
	_ GitHubData[*github.ActionsService, *github.WorkflowRun] = &ghworkflowruns.Source[*github.ActionsService, *github.WorkflowRun]{}
)
