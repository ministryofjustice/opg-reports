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
}

// SetupArgs maps flag values to properies on the arg passed and runs
// flag.Parse to fetch values
func SetupArgs(args *Arguments) {

	flag.StringVar(&args.Organisation, "organisation", defOrg, "organisation slug.")
	flag.StringVar(&args.Team, "team", defTeam, "team slug")
	flag.StringVar(&args.StartDate, "start", defDay.Format(dateformats.YMD), "start date to fetch data for.")
	flag.StringVar(&args.EndDate, "end", defDay.Format(dateformats.YMD), "end date to fetch data for.")
	flag.StringVar(&args.OutputFile, "output", "./data/{start}_{end}_github_releases.json", "Filepath for the output")

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
		all = append(all, pg...)
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
func WorkflowRuns(ctx context.Context, client *github.Client, args *Arguments, repo *github.Repository) (all []*github.WorkflowRun, err error) {
	var (
		sdt, _         = dateutils.Time(args.StartDate)
		edt, _         = dateutils.Time(args.EndDate)
		startDay       = sdt.Format(dateformats.YMD)
		endDay         = edt.Format(dateformats.YMD)
		actionsService = client.Actions
		page           = 1
		opts           = &github.ListWorkflowRunsOptions{
			Branch:  *repo.DefaultBranch,
			Status:  "success",
			Created: fmt.Sprintf("%s..%s", startDay, endDay),
		}
	)
	opts.PerPage = 100
	all = []*github.WorkflowRun{}

	for page > 0 {
		var workflow *github.WorkflowRuns
		var resp *github.Response
		opts.Page = page

		workflow, resp, err = actionsService.ListRepositoryWorkflowRuns(ctx, args.Organisation, *repo.Name, opts)
		slog.Debug("[githubreleases] getting workflow runs",
			slog.String("day", opts.Created),
			slog.Int("page", opts.Page),
			slog.Int("total", *workflow.TotalCount),
			slog.String("repo", *repo.FullName))

		if err != nil {
			return
		}

		for _, run := range workflow.WorkflowRuns {
			var name = cleanWorkflowRunName(*run.Name)

			if strings.HasPrefix(name, pathToLive) {
				all = append(all, run)
			}
		}

		page = resp.NextPage
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