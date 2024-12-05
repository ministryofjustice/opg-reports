/*
githubstandards fetches repo standards data.

Usage:

	githubstandards [flags]

The flags are:

	-organisation=<organisation>
		The name of the github organisation.
		Default: `ministryofjustice`
	-team=<unit>
		Team slug for whose repos to check.
		Default: `opg`
	-output=<path-pattern>
		Path (with magic values) to the output file
		Default: `./data/{month}_{id}_aws_costs.json`

The command presumes an active, autherised session that can connect
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
	"github.com/ministryofjustice/opg-reports/collectors/githubstandards/lib"
	"github.com/ministryofjustice/opg-reports/internal/githubcfg"
	"github.com/ministryofjustice/opg-reports/internal/githubclient"
	"github.com/ministryofjustice/opg-reports/models"
)

var (
	args = &lib.Arguments{}
)

func Run(args *lib.Arguments) (err error) {
	var (
		content      []byte
		cfg          *githubcfg.Config                  = githubcfg.FromEnv()
		client       *github.Client                     = githubclient.Client(cfg.Token)
		ctx          context.Context                    = context.Background()
		stndrds      []*models.GitHubRepositoryStandard = []*models.GitHubRepositoryStandard{}
		repositories []*github.Repository
	)
	// get all repos for the team
	repositories, err = lib.AllRepos(ctx, client, args)
	if err != nil {
		return
	}
	// convert each to a standards entry
	total := len(repositories)
	for i, repo := range repositories {
		slog.Info(fmt.Sprintf("[%d/%d] %s", i+1, total, *repo.FullName))

		var std = lib.RepoToStandard(ctx, client, repo)
		stndrds = append(stndrds, std)
	}
	// write to file
	content, err = json.MarshalIndent(stndrds, "", "  ")
	if err != nil {
		slog.Error("error marshaling", slog.String("err", err.Error()))
		os.Exit(1)
	}
	lib.WriteToFile(content, args)

	return
}

func main() {
	var err error
	lib.SetupArgs(args)

	slog.Info("[githubstandards] starting ...")
	slog.Debug("[githubstandards]", slog.String("args", fmt.Sprintf("%+v", args)))

	if err = lib.ValidateArgs(args); err != nil {
		slog.Error("arg validation failed", slog.String("err", err.Error()))
		os.Exit(1)
	}

	err = Run(args)
	if err != nil {
		panic(err)
	}

	slog.Info("[githubstandards] done.")

}
