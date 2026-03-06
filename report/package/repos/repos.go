package repos

import (
	"context"
	"errors"
	"time"

	"github.com/google/go-github/v84/github"
)

// teamClient wrapper around *github.TeamsService
type teamClient interface {
	ListTeamReposBySlug(ctx context.Context, org, slug string, opts *github.ListOptions) ([]*github.Repository, *github.Response, error)
}

// actionClient wrapper for *github.ActionsService
type actionClient interface {
	// api docs - https://docs.github.com/en/rest/actions/workflow-runs#list-workflow-runs-for-a-repository
	ListRepositoryWorkflowRuns(ctx context.Context, owner, repo string, opts *github.ListWorkflowRunsOptions) (*github.WorkflowRuns, *github.Response, error)
	GetWorkflowRunUsageByID(ctx context.Context, owner, repo string, runID int64) (*github.WorkflowRunUsage, *github.Response, error)
}

// PR Client is a wrapper for *github.PullRequestsService
type prClient interface {
	List(ctx context.Context, owner string, repo string, opts *github.PullRequestListOptions) ([]*github.PullRequest, *github.Response, error)
}

type Args struct {
	OrgSlug      string    `json:"org_slug"`       // github org name
	ParentSlug   string    `json:"parent_slug"`    // parent slug
	FilterByName string    `json:"filter_by_name"` // used to limit the repos to those that exactly match this name
	DateStart    time.Time `json:"date_start"`     // start date
	DateEnd      time.Time `json:"date_end"`       // end date
}

var ErrFailedGettingRepositoryPage = errors.New("error getting page of repositories")
var ErrPROutOfDateRange = errors.New("pr outside of date range")
