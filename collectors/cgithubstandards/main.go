/*
cgithubstandards fetches aws costs data for the month at a daily granularity.

Usage:

	cgithubstandards [flags]

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
	"fmt"
	"log/slog"
	"os"

	"github.com/google/go-github/v62/github"
	"github.com/ministryofjustice/opg-reports/collectors/cgithubstandards/lib"
	"github.com/ministryofjustice/opg-reports/pkg/githubcfg"
	"github.com/ministryofjustice/opg-reports/pkg/githubclient"
)

var (
	args = &lib.Arguments{}
)

func Run(args *lib.Arguments) (err error) {
	var (
		cfg          *githubcfg.Config = githubcfg.FromEnv()
		client       *github.Client    = githubclient.Client(cfg.Token)
		ctx          context.Context   = context.Background()
		repositories []*github.Repository
	)

	repositories, err = lib.AllRepos(ctx, client, args)
	if err != nil {
		return
	}
	fmt.Println(repositories)

	// var (
	// 	s         *session.Session
	// 	client    *costexplorer.CostExplorer
	// 	startDate time.Time
	// 	endDate   time.Time
	// 	raw       *costexplorer.GetCostAndUsageOutput
	// 	data      []*costs.Cost
	// 	content   []byte
	// 	cfg       *awscfg.Config = awscfg.FromEnv()
	// )

	// if s, err = awssession.New(cfg); err != nil {
	// 	slog.Error("[awscosts.main] aws session failed", slog.String("err", err.Error()))
	// 	return
	// }

	// if client, err = awsclient.CostExplorer(s); err != nil {
	// 	slog.Error("[awscosts.main] aws client failed", slog.String("err", err.Error()))
	// 	return
	// }

	// if startDate, err = convert.ToTime(args.Month); err != nil {
	// 	slog.Error("[awscosts.main] month conversion failed", slog.String("err", err.Error()))
	// 	return
	// }
	// startDate = convert.DateResetMonth(startDate)
	// // overwrite month with the parsed version
	// args.Month = startDate.Format(consts.DateFormatYearMonth)
	// endDate = startDate.AddDate(0, 1, 0)

	// if raw, err = lib.CostData(client, startDate, endDate, costexplorer.GranularityDaily, consts.DateFormatYearMonthDay); err != nil {
	// 	slog.Error("[awscosts.main] getting cost data failed", slog.String("err", err.Error()))
	// 	return
	// }
	// if data, err = lib.Flat(raw, args); err != nil {
	// 	slog.Error("[awscosts.main] flattening raw data to costs failed", slog.String("err", err.Error()))
	// 	return
	// }

	// content, err = json.MarshalIndent(data, "", "  ")
	// if err != nil {
	// 	slog.Error("error marshaling", slog.String("err", err.Error()))
	// 	os.Exit(1)
	// }
	// //
	// lib.WriteToFile(content, args)

	return
}

func main() {
	var err error
	lib.SetupArgs(args)

	slog.Info("[cgithubstandards.main] init...")
	slog.Debug("[cgithubstandards.main]", slog.String("args", fmt.Sprintf("%+v", args)))

	if err = lib.ValidateArgs(args); err != nil {
		slog.Error("arg validation failed", slog.String("err", err.Error()))
		os.Exit(1)
	}

	err = Run(args)
	if err != nil {
		panic(err)
	}

	slog.Info("[cgithubstandards.main] done.")

}
