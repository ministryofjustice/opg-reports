package main

import (
	"context"
	"fmt"
	"log/slog"
	"opg-reports/shared/data"
	"opg-reports/shared/env"
	"opg-reports/shared/files"
	"opg-reports/shared/gh/cl"
	"opg-reports/shared/gh/compliance"
	"opg-reports/shared/gh/repos"
	"opg-reports/shared/logger"
	"opg-reports/shared/report"
	"strconv"
)

var (
	includeArchivedArg = report.NewArg("include-archived", false, "When true, includes archived repositories as well", "false")
	orgArg             = report.NewArg("organisation", true, "Name of the organisation we'll get respotiroies for", "ministryofjustice")
	teamArg            = report.NewArg("team", true, "Team within the <organisation> to fetch repositories for", "")
)

const dir string = "data"

func run(r report.IReport) {
	ctx := context.Background()
	token := env.Get("GITHUB_ACCESS_TOKEN", "")
	if token == "" {
		slog.Error("not github token found")
		return
	}

	limiter, _ := cl.RateLimitedHttpClient()
	client := cl.Client(token, limiter)

	includeArchived, _ := strconv.ParseBool(includeArchivedArg.Val())

	repositories, err := repos.All(ctx, client, orgArg.Val(), teamArg.Val(), includeArchived)
	if err != nil {
		slog.Error("error getting repositories",
			slog.String("org", orgArg.Val()),
			slog.String("team", teamArg.Val()),
			slog.Bool("includeArchived", includeArchived),
			slog.String("err", fmt.Sprintf("%v", err)),
		)
		return
	}
	slog.Info("repository count", slog.Int("count", len(repositories)))

	toStore := []*compliance.Compliance{}
	for i, rep := range repositories {
		slog.Info(fmt.Sprintf("[%d] %s", i, rep.GetFullName()))
		cmp := compliance.NewFromR(ctx, nil, rep, client)
		toStore = append(toStore, cmp)
	}
	content, err := data.ToJsonList[*compliance.Compliance](toStore)
	filename := r.Filename()
	err = files.WriteFile(dir, filename, content)

	fmt.Println(err)
}

func main() {
	logger.LogSetup()
	costReport := report.New(includeArchivedArg, orgArg, teamArg)
	costReport.SetRunner(run)
	costReport.Run()

}
