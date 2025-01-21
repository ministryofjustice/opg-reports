package lib

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/google/go-github/v62/github"
	"github.com/ministryofjustice/opg-reports/internal/dateformats"
	"github.com/ministryofjustice/opg-reports/internal/dateintervals"
	"github.com/ministryofjustice/opg-reports/internal/dateutils"
	"github.com/ministryofjustice/opg-reports/models"
)

// defaults
var (
	defOrg  = "ministryofjustice"
	defTeam = "opg"
	defDay  = dateutils.Reset(time.Now().UTC(), dateintervals.Day).AddDate(0, 0, -1)
)

var pathToLive = "path to live"

// Arguments represents all the named arguments for this collector
type Arguments struct {
	Organisation string // Organisation. Default "ministryofjustice"
	Team         string // Team. Default "opg"
	StartDate    string
	EndDate      string
	OutputFile   string // OutputFile is destination of data. Default "./data/{start}_{end}_github_releases.json"
	Repository   string // Optional repository filter
}

// SetupArgs maps flag values to properies on the arg passed and runs
// flag.Parse to fetch values
func SetupArgs(args *Arguments) {

	flag.StringVar(&args.Organisation, "organisation", defOrg, "organisation slug.")
	flag.StringVar(&args.Team, "team", defTeam, "team slug")
	flag.StringVar(&args.StartDate, "start", defDay.Format(dateformats.YMD), "start date to fetch data for.")
	flag.StringVar(&args.EndDate, "end", defDay.Format(dateformats.YMD), "end date to fetch data for.")
	flag.StringVar(&args.OutputFile, "output", "./data/{start}_{end}_github_releases.json", "Filepath for the output")
	flag.StringVar(&args.Repository, "repository", "", "Filter results to just this repository. Optional")
	flag.Parse()
}

// ValidateArgs checks rules and logic for the input arguments
// Make sure some have non empty values and apply default values to others
func ValidateArgs(args *Arguments) (err error) {
	failOnEmpty := map[string]string{
		"output": args.OutputFile,
	}
	for k, v := range failOnEmpty {
		if v == "" {
			err = errors.Join(err, fmt.Errorf("%s", k))
		}
	}
	if err != nil {
		err = fmt.Errorf("missing arguments: [%s]", strings.ReplaceAll(err.Error(), "\n", ", "))
	}

	if args.Organisation == "" {
		args.Organisation = defOrg
	}
	if args.Team == "" {
		args.Team = defTeam
	}

	if args.StartDate == "" || args.StartDate == "-" {
		args.StartDate = defDay.Format(dateformats.YMD)
	}
	if args.EndDate == "" || args.EndDate == "-" {
		args.EndDate = defDay.Format(dateformats.YMD)
	}

	return
}

// WriteToFile writes the content to the file
func WriteToFile(content []byte, args *Arguments) {
	var (
		filename string = args.OutputFile
		dir      string = filepath.Dir(args.OutputFile)
	)
	filename = strings.ReplaceAll(filename, "{start}", args.StartDate)
	filename = strings.ReplaceAll(filename, "{end}", args.EndDate)
	os.MkdirAll(dir, os.ModePerm)
	os.WriteFile(filename, content, os.ModePerm)

}

// TeamList generates a list of all teams attached to this repo
func TeamList(ctx context.Context, client *github.Client, organisation string, repoName string) (teams models.GitHubTeams, err error) {
	teams, err = models.NewGitHubTeams(ctx, client, organisation, repoName)
	return
}

// WorkflowRunsToReleases converts a slice of workflow runs into a slice of Release which can then be stored
func WorkflowRunsToReleases(ghRepo *models.GitHubRepository, teams models.GitHubTeams, runs []*github.WorkflowRun) (all []*models.GitHubRelease, err error) {

	all = []*models.GitHubRelease{}

	for _, run := range runs {
		var ts = time.Now().UTC().Format(dateformats.Full)
		var release = &models.GitHubRelease{
			Ts:               ts,
			Name:             *run.Name,
			Date:             run.CreatedAt.Format(dateformats.Full),
			SourceURL:        *run.HTMLURL,
			RelaseType:       models.GitHubWorkflowRelease,
			GitHubRepository: (*models.GitHubRepositoryForeignKey)(ghRepo),
			Count:            1,
		}

		all = append(all, release)
	}

	return
}

// PullRequestsToReleases converts a set of pull requests into releases
func PullRequestsToReleases(ghRepo *models.GitHubRepository, teams models.GitHubTeams, prs []*github.PullRequest) (all []*models.GitHubRelease, err error) {

	all = []*models.GitHubRelease{}

	for _, pr := range prs {
		var ts = time.Now().UTC().Format(dateformats.Full)
		var release = &models.GitHubRelease{
			Ts:               ts,
			Name:             *pr.Title,
			Date:             pr.MergedAt.Format(dateformats.Full),
			SourceURL:        *pr.HTMLURL,
			RelaseType:       models.GitHubPRMergeRelease,
			GitHubRepository: (*models.GitHubRepositoryForeignKey)(ghRepo),
			Count:            1,
		}
		all = append(all, release)
	}

	return
}

// AllRepos returns all accessible repos for the details passed
// Uses page iterating for loop to handle api calls
func AllRepos(ctx context.Context, client *github.Client, args *Arguments) (all []*github.Repository, err error) {
	var (
		org  string = args.Organisation
		team string = args.Team
		page int    = 1
	)

	all = []*github.Repository{}

	for page > 0 {
		slog.Debug("[githubreleases] getting repostiories", slog.Int("page", page))
		pg, resp, e := client.Teams.ListTeamReposBySlug(ctx, org, team, &github.ListOptions{PerPage: 100, Page: page})
		if e != nil {
			err = e
			return
		}
		// If there is no filter, than attach all
		// otherwise, look specifically for it
		if args.Repository == "" {
			all = append(all, pg...)
		} else {
			slog.Debug("[githubreleases] filtering results", slog.String("repository", args.Repository))
			for _, repo := range pg {
				if *repo.FullName == args.Repository {
					all = append(all, repo)
				}
			}
		}
		page = resp.NextPage
	}

	return
}

func cleanWorkflowRunName(name string) (clean string) {
	clean = strings.ToLower(name)
	clean = strings.TrimPrefix(clean, "[workflow]")
	clean = strings.TrimPrefix(clean, "[job]")
	clean = strings.TrimSpace(clean)
	return
}

// WorkflowRuns returns all the workflow runs for this repos on the day (-day)
// requested.
// Looks for only successful runs and matchs the name against a standard prefix
// ('path to live')
// Cleans up the workflow name, removing some known starting elements such as
// `[Workflow]`, `[Job]` - and trims whitespace
// Uses page iterating for loop to handle api calls
//
// Note: API returns a max of 1k results in one call, so a long time range will
// likely cause multiple api calls to happen
//
// Note: makes async calls the github api to fetch api data
func WorkflowRuns(ctx context.Context, client *github.Client, args *Arguments, repo *github.Repository) (all []*github.WorkflowRun, err error) {
	var (
		opts           []*github.ListWorkflowRunsOptions
		total          int                    = 0
		actionsService *github.ActionsService = client.Actions

		mutex *sync.Mutex    = &sync.Mutex{}
		wg    sync.WaitGroup = sync.WaitGroup{}
	)
	all = []*github.WorkflowRun{}

	total, opts, err = decideWorkflowApiCall(ctx, client, args, repo)
	if err != nil {
		return
	}
	slog.Info("[githubreleases] decided on api calls required.",
		slog.Int("apiCallCount", len(opts)),
		slog.Int("total", total),
		slog.String("repo", *repo.FullName))

	for idx, opt := range opts {
		wg.Add(1)
		// use a go routine to make this call
		go func(o *github.ListWorkflowRunsOptions, i int) {
			slog.Debug("[githubreleases] getting workflow runs",
				slog.Int("i", i),
				slog.String("range", opt.Created),
				slog.String("repo", *repo.FullName))

			found, e := fetcher(ctx, repo, args, opt, actionsService)

			mutex.Lock()
			defer mutex.Unlock()
			if e != nil {
				err = errors.Join(e, err)
			} else {
				all = append(all, found...)
			}

			wg.Done()
		}(opt, idx)
	}
	wg.Wait()

	return
}

// fetcher called in go routine to get data
func fetcher(ctx context.Context, repo *github.Repository, args *Arguments, opt *github.ListWorkflowRunsOptions, actionsService *github.ActionsService) (found []*github.WorkflowRun, err error) {
	var page int = 1
	found = []*github.WorkflowRun{}

	for page > 0 {
		var (
			runs *github.WorkflowRuns
			resp *github.Response
		)
		opt.Page = page
		runs, resp, err = actionsService.ListRepositoryWorkflowRuns(ctx, args.Organisation, *repo.Name, opt)
		slog.Debug("[githubreleases] workflow run results",
			slog.Int("currentPage", opt.Page),
			slog.Int("nextPage", resp.NextPage),
			slog.Int("total", *runs.TotalCount),
			slog.String("repo", *repo.FullName))
		// error
		if err != nil {
			fmt.Println(err)
			return
		}
		// loop over all workflow run names
		for _, run := range runs.WorkflowRuns {
			var name = cleanWorkflowRunName(*run.Name)
			slog.Debug("[githubreleases] workflow name",
				slog.String("created", run.CreatedAt.String()),
				slog.String("workflow", name),
				slog.String("repo", *repo.FullName))

			if strings.HasPrefix(name, pathToLive) {
				found = append(found, run)
			}
		}

		page = resp.NextPage

	}
	return
}

// decideWorkflowApiCall determines how many times call this api end point
//
// ListRepositoryWorkflowRuns when called returns at most 1k records at a time, so
// we make the first call to the api to figure out the number of results for the
// entire time range and then decide if we need to make multiple calls.
//
// When asking for a long time range this could well be split down to monthly or weekly
// date ranges to call the api in.
//
// NOTE: If the repo has more than 1k workflow runs per day, there is currently nothing to handle that!
func decideWorkflowApiCall(ctx context.Context, client *github.Client, args *Arguments, repo *github.Repository) (resultCount int, allOpts []*github.ListWorkflowRunsOptions, err error) {
	var (
		apiResultLimit = 1000
		perPage        = 100
		page           = 1

		start, _ = dateutils.Time(args.StartDate)
		end, _   = dateutils.Time(args.EndDate)

		actionsService = client.Actions
		workflowRuns   *github.WorkflowRuns
		resp           *github.Response

		months = dateutils.CountInRange(start, end, dateintervals.Month)
		days   = dateutils.CountInRange(start, end, dateintervals.Day)
		weeks  = (days / 7)

		opts = &github.ListWorkflowRunsOptions{
			Branch:              *repo.DefaultBranch,
			ExcludePullRequests: true,
			Status:              "success",
		}
	)

	allOpts = []*github.ListWorkflowRunsOptions{}

	// see how many total pages there are between these dates
	opts.Page = page
	opts.PerPage = perPage
	opts.Created = fmt.Sprintf("%s..%s", start.Format(dateformats.YMD), end.Format(dateformats.YMD))

	workflowRuns, resp, err = actionsService.ListRepositoryWorkflowRuns(ctx, args.Organisation, *repo.Name, opts)
	if err != nil {
		return
	}

	resultCount = *workflowRuns.TotalCount
	monthPages := (resultCount / months)
	weekPages := (resultCount / weeks)
	dayPages := (resultCount / days)

	slog.Debug("[githubreleases] repository workflow counts",
		slog.String("dateRange", opts.Created),
		slog.Int("months", months),
		slog.Int("weeks", weeks),
		slog.Int("days", days),
		slog.Int("apiResultLimit", apiResultLimit),
		slog.Int("lastPage", resp.LastPage),
		slog.Int("resultCount", resultCount),
		slog.Int("monthPages", monthPages),
		slog.Int("weekPages", weekPages),
		slog.Int("dayPages", dayPages),
	)

	created := []string{}

	if resultCount <= apiResultLimit {
		slog.Debug("[githubreleases] one api call for ALL")
		created = append(created, opts.Created)
	} else if monthPages <= perPage {
		slog.Debug("[githubreleases] api calls per MONTH")
		created = createdStrings(end, dateutils.Times(start, end, dateintervals.Month))
	} else if weekPages <= perPage {
		slog.Debug("[githubreleases] api calls per WEEK")
		created = createdStrings(end, dateutils.TimesI(start, end, dateintervals.Day, 7))
	} else {
		slog.Debug("[githubreleases] api calls per DAY")
		created = createdStrings(end, dateutils.Times(start, end, dateintervals.Day))
	}
	// generate the list of opts to call to catch as much as we can
	for _, str := range created {
		opt := &github.ListWorkflowRunsOptions{
			Branch:              *repo.DefaultBranch,
			ExcludePullRequests: true,
			Status:              "success",
			Created:             str,
		}
		opt.Page = page
		opt.PerPage = perPage
		allOpts = append(allOpts, opt)
	}

	return
}

// createdStrings uses the dates passed to created date ranges for use in api calls
func createdStrings(end time.Time, dates []time.Time) (created []string) {
	created = []string{}

	l := len(dates)
	for i, date := range dates {
		e := date
		if i+1 < l {
			e = dates[i+1]
		} else {
			e = end
		}
		created = append(
			created,
			fmt.Sprintf("%s..%s", date.Format(dateformats.YMD), e.Format(dateformats.YMD)),
		)
	}
	return
}

// MergedPullRequests finds all closed pull requests that were merged on the day (-day) we've asked for.
// Requests pull requests to be returned in `updated` date descending order, and then assumes that when we
// find a pr whose updated time is before the day we asked for, we can then skip all the rest of the results
// as anything more should also be older. There might be some oddities here with pr's closed and re-opened.
// Only when the merged time is on the same day as the (`-day`) we asked for will it be added to the
// returned data.
//
// Uses page iterating for loop to handle api calls
func MergedPullRequests(ctx context.Context, client *github.Client, args *Arguments, repo *github.Repository) (all []*github.PullRequest, err error) {
	var (
		start, _  = dateutils.Time(args.StartDate)
		end, _    = dateutils.Time(args.EndDate)
		prService = client.PullRequests
		page      = 1
		opts      = &github.PullRequestListOptions{
			State:     "closed",
			Sort:      "updated",
			Direction: "desc",
			Base:      *repo.DefaultBranch,
		}
	)
	opts.PerPage = 100
	all = []*github.PullRequest{}

	// loops over the page numbers (github api calls are paginated)
	for page > 0 {
		var prs = []*github.PullRequest{}
		var resp *github.Response

		opts.Page = page

		prs, resp, err = prService.List(ctx, args.Organisation, *repo.Name, opts)
		slog.Debug("[githubreleases] getting pull requests",
			slog.String("state", opts.State),
			slog.Int("page", opts.Page),
			slog.Int("count", len(prs)),
			slog.String("repo", *repo.FullName))

		if err != nil {
			return
		}
		// loop over the current block of pull requests
		for _, pr := range prs {
			var (
				mergedAt    time.Time
				when        time.Time = pr.UpdatedAt.Time
				old         bool      = when.Before(start)
				merged      bool      = false
				withinRange bool
			)
			if pr.MergeCommitSHA != nil {
				merged = len(*pr.MergeCommitSHA) > 0
			}
			// if its older than we want we can skip all the rest of the records
			// as results are in date descending
			if old {
				return
			}

			if merged && pr.MergedAt != nil {
				mergedAt = dateutils.Reset(pr.MergedAt.Time, dateintervals.Day)
				afterStart := mergedAt.After(start) || mergedAt.Equal(start)
				beforeEnd := mergedAt.Before(end) || mergedAt.Equal(end)

				withinRange = afterStart && beforeEnd

				slog.Debug("[githubreleases] pull request",
					slog.String("repo", *repo.Name),
					slog.Bool("withinDateRange", withinRange),
					slog.String("name", *pr.Title),
					slog.String("url", pr.GetURL()),
					slog.String("updatedAt", when.String()),
					slog.Bool("merged", merged),
					slog.Any("mergedAt", mergedAt),
					slog.Bool("old", old),
				)

				if withinRange {
					all = append(all, pr)
				}
			}

		}

		page = resp.NextPage
	}

	return

}
