package lib

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/go-github/v62/github"
	"github.com/ministryofjustice/opg-reports/internal/dateformats"
	"github.com/ministryofjustice/opg-reports/internal/dateutils"
	"github.com/ministryofjustice/opg-reports/internal/githubcfg"
	"github.com/ministryofjustice/opg-reports/internal/githubclient"
	"github.com/ministryofjustice/opg-reports/models"
)

func Run(args *Arguments) (err error) {
	var (
		content      []byte
		total        int = 0
		repositories []*github.Repository
		cfg          *githubcfg.Config = githubcfg.FromEnv()
		client       *github.Client    = githubclient.Client(cfg.Token)
		ctx          context.Context   = context.Background()
		allReleases  []*models.GitHubRelease
	)
	if err = ValidateArgs(args); err != nil {
		slog.Error("arg validation failed", slog.String("err", err.Error()))
		return
	}
	slog.Info("[githubreleases] running with args.",
		slog.String("Organisation", args.Organisation),
		slog.String("Team", args.Team),
		slog.String("StartDate", args.StartDate),
		slog.String("EndDate", args.EndDate),
		slog.String("Repository", args.Repository),
		slog.String("OutputFile", args.OutputFile),
	)
	args.StartDate = dateutils.Reformat(args.StartDate, dateformats.YMD)
	args.EndDate = dateutils.Reformat(args.EndDate, dateformats.YMD)

	// get all repos for the team
	repositories, err = AllRepos(ctx, client, args)
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
			teams, _                         = TeamList(ctx, client, args.Organisation, *repo.Name)
		)

		slog.Info(fmt.Sprintf("[%d/%d] %s", i+1, total, *repo.FullName))
		// rough api rate limiting help - TODO: make better
		time.Sleep(1 * time.Second)
		// Look for workflow runs first
		runs, err = WorkflowRuns(ctx, client, args, repo)
		// if theres none, look for merges
		if len(runs) == 0 {
			merged, err = MergedPullRequests(ctx, client, args, repo)
		}
		if err != nil {
			slog.Error("error with workflows or merges..", slog.String("err", err.Error()))
			return
		}

		slog.Debug("[githubreleases] found for dates.",
			slog.String("repository", *repo.FullName),
			slog.String("StartDate", args.StartDate),
			slog.String("EndDate", args.EndDate),
			slog.Int("pull_requests", len(merged)),
			slog.Int("workflow_runs", len(runs)))

		repoModel := models.NewRepository(ctx, client, repo)
		// convert to releases
		if len(runs) > 0 {
			rels, err = WorkflowRunsToReleases(repoModel, teams, runs)
		} else if len(merged) > 0 {
			rels, err = PullRequestsToReleases(repoModel, teams, merged)
		}

		if err != nil {
			slog.Error("error converting to releases..", slog.String("err", err.Error()))
			return
		}

		// attach the releases to the main set
		allReleases = append(allReleases, rels...)
		// Save to file during loop
		// Save(allReleases, args)
	}

	// write to file
	content, err = json.MarshalIndent(allReleases, "", "  ")
	if err != nil {
		slog.Error("[githubreleases] error marshaling", slog.String("err", err.Error()))
		return
	}
	WriteToFile(content, args)

	return
}

// func Save(allReleases []*models.GitHubRelease, args *Arguments) {
// 	// write to file
// 	content, err := json.MarshalIndent(allReleases, "", "  ")
// 	if err != nil {
// 		slog.Error("[githubreleases] error marshaling", slog.String("err", err.Error()))
// 		os.Exit(1)
// 	}
// 	WriteToFile(content, args)
// }
