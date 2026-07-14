// Package ghdata
//
// direct calls
// - list of all repositories
// - list of workflow runs
// - list of merged pull requests
//
// filters
// - repos
//   - ignore archived repos
//   - name matched repo
//
// - workflows
//   - workflow run name
package ghdata

import (
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
//
// github.Repository
// Used for listing owned code bases.
//
// github.WorkflowRun
// Used to find path to live runs of a workflow for a set of repositories - part of code releases
type Result interface {
	*github.Repository | *github.WorkflowRun | *github.PullRequest
}

// GitHubData interface exposes the main methods used for fetching data
// and then filtering that information
type GitHubData[C Client, R Result] interface {
	// GetData returns slice of data from the api
	GetData() (results []R, skipped []any, err error)
	//
}
