package codebasereleasesimport

import (
	"context"
	"fmt"
	"log/slog"
	"opg-reports/report/package/cntxt"
	"opg-reports/report/package/dbx"
	"opg-reports/report/package/repos"
	"opg-reports/report/package/times"
	"strings"
	"time"

	"github.com/google/go-github/v84/github"
)

// insert, but only those parts about releases
const InsertReleasesStatement string = `
INSERT INTO codebase_metrics (
	codebase,
	month,
	releases,
	releases_securityish
) VALUES (
	:codebase,
	:month,
	:releases,
	:releases_securityish
) ON CONFLICT (codebase,month) DO UPDATE SET
	releases=excluded.releases,
	releases_securityish=excluded.releases_securityish
RETURNING id
;
`

// teamClient wrapper around *github.TeamsService
type teamClient interface {
	ListTeamReposBySlug(ctx context.Context, org, slug string, opts *github.ListOptions) ([]*github.Repository, *github.Response, error)
}

// ActionClient wrapper for *github.ActionsService
type actionClient interface {
	ListRepositoryWorkflowRuns(ctx context.Context, owner, repo string, opts *github.ListWorkflowRunsOptions) (*github.WorkflowRuns, *github.Response, error)
	GetWorkflowRunUsageByID(ctx context.Context, owner, repo string, runID int64) (*github.WorkflowRunUsage, *github.Response, error)
}

// PR Client is a wrapper for *github.PullRequestsService
type prClient interface {
	List(ctx context.Context, owner string, repo string, opts *github.PullRequestListOptions) ([]*github.PullRequest, *github.Response, error)
}

type Args struct {
	DB           string    `json:"db"`             // database path
	Driver       string    `json:"driver"`         // database driver
	Params       string    `json:"params"`         // database connection params
	OrgSlug      string    `json:"org_slug"`       // github org name
	ParentSlug   string    `json:"parent_slug"`    // parent slug
	FilterByName string    `json:"filter_by_name"` // used to limit the repos to those that exactly match this name
	DateStart    time.Time `json:"date_start"`     // start date
	DateEnd      time.Time `json:"date_end"`       // end date
}

type Clients struct {
	Teams   teamClient   // *github.TeamsService
	Actions actionClient // *github.ActionsService
	PR      prClient     // *github.PullRequestService
}

// CodebaseMetric
//
// WorkflowRun data is used for average times, so a repo that does not use
// github actions will have an empty value
type CodebaseMetric struct {
	Codebase            string `json:"codebase"`             // full name of codebase
	Month               string `json:"month"`                // month as YYYY-MM string
	Releases            int    `json:"releases"`             // count of releases for this month
	ReleasesSecurityish int    `json:"releases_securityish"` // count of releases for this month that seem to be security related
}

// Import finds all github repositories and returns them for the moj/opg team
func Import(ctx context.Context, clients *Clients, in *Args) (err error) {
	var log *slog.Logger = cntxt.GetLogger(ctx).With("package", "codebasereleasesimport", "func", "Import")
	var repoList []*github.Repository
	var data = []*CodebaseMetric{}

	log.Info("starting ...")
	// fetch all the repos
	log.Debug("getting repository list ...")
	repoList, err = repos.GetList(ctx, clients.Teams, &repos.Args{
		OrgSlug:      in.OrgSlug,
		ParentSlug:   in.ParentSlug,
		FilterByName: in.FilterByName,
	})
	if err != nil {
		return
	}

	data, err = handler(ctx, clients, in, repoList)
	if err != nil {
		log.Error("error processing repos", "err", err.Error())
		return
	}

	// now write to db
	err = dbx.Insert(ctx, InsertReleasesStatement, data, &dbx.InsertArgs{
		DB:     in.DB,
		Driver: in.Driver,
		Params: in.Params,
	})
	if err != nil {
		log.Error("error write data during import", "err", err.Error())
		return
	}

	log.Info("complete.")
	return
}

// handler looks at both workflows & pull requests to get release data
func handler(ctx context.Context, clients *Clients, in *Args, repoList []*github.Repository) (data []*CodebaseMetric, err error) {
	var log *slog.Logger = cntxt.GetLogger(ctx).With("package", "codebasereleasesimport", "func", "Import")
	var byMonth = map[string][]*CodebaseMetric{}

	data = []*CodebaseMetric{}
	// now loop over each repo to then call helper methods
	for _, repo := range repoList {
		var lg = log.With("repo", *repo.Name)
		var found = []*CodebaseMetric{}
		log.Info(fmt.Sprintf("[%s]", *repo.Name))
		// dont get any release info on archived code bases
		if *repo.Archived {
			log.Warn("repository is archived, skipping fetching compliance details.")
			continue
		}
		found, err = workflowRunReleases(ctx, clients.Actions, in, repo)
		// return on error?
		if err != nil {
			return
		}
		// if found runs, then set the data and continue the loop
		if len(found) > 0 {
			lg.Info("found workflow runs ... ", "count", len(found))
		} else {
			lg.Info("no workflow runs found, looking for pull requests ... ")
			// get pull request data if theres no run data (runs will have triggered the break)
			found, err = mergedPullRequestReleases(ctx, clients.PR, in, repo)
			if err != nil {
				return
			}
			lg.Info("found pull requests ... ", "count", len(found))
		}
		// append the found data
		data = append(data, found...)
	}

	// flattern byMonth to a slice for insert
	for _, v := range byMonth {
		data = append(data, v...)
	}
	return
}

// mergedPullRequestReleases finds release data based on pull requests merged into the default branch
// as a proxy measure for the path to live
func mergedPullRequestReleases(ctx context.Context, client prClient, in *Args, repo *github.Repository) (metrics []*CodebaseMetric, err error) {
	var log *slog.Logger = cntxt.GetLogger(ctx).With("package", "codebasereleasesimport", "func", "mergedPullRequestReleases")
	var prs []*github.PullRequest
	var byMonth = map[string]*CodebaseMetric{}

	metrics = []*CodebaseMetric{}
	prs, err = repos.GetMergedPRs(ctx, client, repo, &repos.Args{
		OrgSlug:    in.OrgSlug,
		ParentSlug: in.ParentSlug,
		DateStart:  in.DateStart,
		DateEnd:    in.DateEnd,
	})
	if err != nil {
		log.Error("error getting merged prs", "err", err.Error())
		return
	}

	for _, pr := range prs {
		var when = times.AsYMString(pr.MergedAt.Time)
		if _, ok := byMonth[when]; !ok {
			byMonth[when] = emptyMetric(repo, when)
		}
		byMonth[when].Releases += 1
		byMonth[when].ReleasesSecurityish += isSecurityishPR(pr)
	}

	for _, v := range byMonth {
		metrics = append(metrics, v)
	}

	return

}

// workflowRunReleases fetches workflows that run against main with name of path to live only and counts that as a release
func workflowRunReleases(ctx context.Context, client actionClient, in *Args, repo *github.Repository) (metrics []*CodebaseMetric, err error) {
	var log *slog.Logger = cntxt.GetLogger(ctx).With("package", "codebasereleasesimport", "func", "workflowRunReleases", "repo", *repo.Name)
	var runs []*github.WorkflowRun
	var byMonth = map[string]*CodebaseMetric{}

	metrics = []*CodebaseMetric{}
	// get just the release / path to live workflow runs
	runs, err = repos.GetWorkflowRuns(ctx, client, repo, &repos.Args{
		OrgSlug:      in.OrgSlug,
		ParentSlug:   in.ParentSlug,
		DateStart:    in.DateStart,
		DateEnd:      in.DateEnd,
		FilterByName: "path to live",
	}, true)

	if err != nil {
		log.Error("error getting release workflow runs", "err", err.Error())
		return
	} else if len(runs) == 0 {
		return
	}
	// group workflow data by month
	for _, run := range runs {
		var when = times.AsYMString(run.CreatedAt.Time)
		if _, ok := byMonth[when]; !ok {
			byMonth[when] = emptyMetric(repo, when)
		}
		log.Debug("adding stats for workflow run ...", "when", when)
		// get stats
		byMonth[when].Releases += 1
		byMonth[when].ReleasesSecurityish += isSecurityishRun(run)
		// updated[when].Dur += runDuration(ctx, client, repo, run)
	}
	// flattern the months
	for _, v := range byMonth {
		metrics = append(metrics, v)
	}

	// // work out averages
	// for _, v := range updated {
	// 	var avg = v.Dur / int64(v.Releases)
	// 	v.ReleasesAverageTime = fmt.Sprintf("%d", avg)
	// }

	// dump out workflow data
	// dump.Now(updated)

	return
}

// // runDuration returns the total run time of the job in milliseconds (as a time.Duration)
// func runDuration(ctx context.Context, client actionClient, repo *github.Repository, run *github.WorkflowRun) (duration int64) {
// 	var err error
// 	var usage *github.WorkflowRunUsage
// 	usage, _, err = client.GetWorkflowRunUsageByID(ctx, *repo.Owner.Login, *repo.Name, *run.ID)

// 	if err != nil {
// 		return
// 	}
// 	duration = *usage.RunDurationMS
// 	// dur = time.Millisecond * time.Duration(ms)
// 	return
// }

// isSecurityishRun returns a 0 or 1 to say if its likely to be security related
// workflow.
//
// returned int is added to a counter
func isSecurityishRun(run *github.WorkflowRun) (securityish int) {
	var msg string
	var author string
	securityish = 0
	// if we dont have any of the required fields to check, return and therefore not be security
	if run.HeadCommit == nil || run.HeadCommit.Message == nil ||
		run.HeadCommit.Author == nil || run.HeadCommit.Author.Name == nil {
		return
	}
	// check the commit for a security related content
	msg = strings.ToLower(*run.HeadCommit.Message)
	if strings.Contains(msg, "security") || strings.Contains(msg, "vuln") {
		securityish = 1
	}
	// if the head commit is from a bot, then presume security related
	author = strings.ToLower(*run.HeadCommit.Author.Name)
	if strings.Contains(author, "dependabot") || strings.Contains(author, "renovate") {
		securityish = 1
	}

	return
}

func isSecurityishPR(pr *github.PullRequest) (securityish int) {
	var title string = ""
	var body string = ""
	var author string = ""
	var authorType string = ""

	securityish = 0

	if pr.Title != nil {
		title = strings.ToLower(*pr.Title)
	}
	if pr.Body != nil {
		body = strings.ToLower(*pr.Body)
	}
	if pr.User != nil && pr.User.Login != nil {
		author = strings.ToLower(*pr.User.Login)
	}
	if pr.User != nil && pr.User.Type != nil {
		authorType = strings.ToLower(*pr.User.Type)
	}

	if strings.Contains(title, "vuln") || strings.Contains(title, "security") {
		securityish = 1
	}
	if strings.Contains(body, "vuln") || strings.Contains(body, "security") {
		securityish = 1
	}
	if strings.Contains(author, "dependabot") || strings.Contains(author, "renovate") {
		securityish = 1
	}
	if strings.Contains(authorType, "bot") {
		securityish = 1
	}

	return
}

func emptyMetric(repo *github.Repository, month string) *CodebaseMetric {
	return &CodebaseMetric{
		Codebase:            *repo.FullName,
		Month:               month,
		Releases:            0,
		ReleasesSecurityish: 0,
	}
}
