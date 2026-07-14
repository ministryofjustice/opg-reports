// Package ghdata
package ghdata

import (
	"opg-reports/app/internal/ghdata/ghmergedprs"
	"opg-reports/app/internal/ghdata/ghworkflowruns"

	"github.com/google/go-github/v87/github"
)

// Client is list of all github services that the data sources utilise:
//
// github.TeamsService
// Used for listing repositories within and organisation for a specific team
//
// github.ActionsService
// Used for finding workflow runs for a set of repositories - part of code releases metrics
//
// github.PullRequestsService
// Used for finding merged pull requests on a set of repositories - part of code release metrics
type Client interface {
	*github.TeamsService | *github.ActionsService | *github.PullRequestsService
}

// Result is the singular result type of all github data sources that
// are used.
type Result interface {
	*github.Repository | *ghworkflowruns.ResultData | *ghmergedprs.ResultData
}

// GitHubData interface exposes the main methods used for fetching data
// and then filtering that information
type GitHubData[C Client, R Result] interface {
	// GetData returns slice of data from the api
	GetData() (results []R, skipped []any, err error)
	//
}
