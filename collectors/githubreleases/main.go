/*
githubreleases fetches release data (using mix of workflow runs and merges to main as a proxy).

Usage:

	githubreleases [flags]

The flags are:

	-organisation=<organisation>
		The name of the github organisation.
		Default: `ministryofjustice`
	-team=<unit>
		Team slug for whose repos to check.
		Default: `opg`
	-day=<yyyy-mm-dd>
		Day to fetch data for.
	-output=<path-pattern>
		Path (with magic values) to the output file
		Default: `./data/{day}_github_releases.json`

The command presumes an active, authorised session that can connect
to GitHub.
*/
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"

	"github.com/google/go-github/v62/github"
	"github.com/ministryofjustice/opg-reports/collectors/githubreleases/lib"
	"github.com/ministryofjustice/opg-reports/models"
	"github.com/ministryofjustice/opg-reports/pkg/consts"
	"github.com/ministryofjustice/opg-reports/pkg/convert"
	"github.com/ministryofjustice/opg-reports/pkg/githubcfg"
	"github.com/ministryofjustice/opg-reports/pkg/githubclient"
)

var (
	args = &lib.Arguments{}
)

func Run(args *lib.Arguments) (err error) {
	var (
		content      []byte
		total        int = 0
		repositories []*github.Repository
		cfg          *githubcfg.Config = githubcfg.FromEnv()
		client       *github.Client    = githubclient.Client(cfg.Token)
		ctx          context.Context   = context.Background()
		allReleases  []*models.GitHubRelease
	)
	if err = lib.ValidateArgs(args); err != nil {
		slog.Error("arg validation failed", slog.String("err", err.Error()))
		return
	}
	args.Day = convert.DateReformat(args.Day, consts.DateFormatYearMonthDay)
	// get all repos for the team
	repositories, err = lib.AllRepos(ctx, client, args)
	if err != nil {
		return
	}

	total = len(repositories)
	// Loop over all repos
	// - look for workflow runs on the day
	// - if non are found, then look for merges to main on the day
	for i, repo := range repositories {
		var (
			rels     []*models.GitHubRelease = []*models.GitHubRelease{}
			runs     []*github.WorkflowRun   = []*github.WorkflowRun{}
			merged   []*github.PullRequest   = []*github.PullRequest{}
			teams, _                         = lib.TeamList(ctx, client, args.Organisation, *repo.Name)
		)

		slog.Info(fmt.Sprintf("[%d/%d] %s", i+1, total, *repo.FullName))

		// Look for workflow runs first
		runs, err = lib.WorkflowRuns(ctx, client, args, repo)
		if err != nil {
			return
		}
		// if theres none, look for merges
		if len(runs) == 0 {
			merged, err = lib.MergedPullRequests(ctx, client, args, repo)
		}

		slog.Info("[githubreleases] found for day.", slog.String("day", args.Day), slog.Int("pull_requests", len(merged)), slog.Int("workflow_runs", len(runs)))
		// convert to releases
		if len(runs) > 0 {
			rels, err = lib.WorkflowRunsToReleases(repo, teams, runs)
		} else if len(merged) > 0 {
			rels, err = lib.PullRequestsToReleases(repo, teams, merged)
		}

		// attach the releases to the main set
		allReleases = append(allReleases, rels...)

	}

	// write to file
	content, err = json.MarshalIndent(allReleases, "", "  ")
	if err != nil {
		slog.Error("[githubreleases] error marshaling", slog.String("err", err.Error()))
		os.Exit(1)
	}
	lib.WriteToFile(content, args)

	return
}

func main() {
	var err error
	lib.SetupArgs(args)

	slog.Info("[githubreleases] starting ...")
	slog.Debug("[githubreleases]", slog.String("args", fmt.Sprintf("%+v", args)))

	err = Run(args)
	if err != nil {
		panic(err)
	}

	slog.Info("[githubreleases] done.")

}
