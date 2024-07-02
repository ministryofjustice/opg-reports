package main

import (
	"context"
	"fmt"
	"log/slog"
	"opg-reports/shared/data"
	"opg-reports/shared/env"
	"opg-reports/shared/files"
	"opg-reports/shared/github/cl"
	"opg-reports/shared/github/repos"
	"opg-reports/shared/github/std"
	"opg-reports/shared/logger"
	"opg-reports/shared/report"
)

var (
	orgArg  = report.NewArg("organisation", true, "Name of the organisation we'll get repositories for", "ministryofjustice")
	teamArg = report.NewArg("team", true, "Team within the <organisation> to fetch repositories for", "")
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

	repositories, err := repos.All(ctx, client, orgArg.Val(), teamArg.Val(), true)
	if err != nil {
		slog.Error("error getting repositories",
			slog.String("org", orgArg.Val()),
			slog.String("team", teamArg.Val()),
			slog.String("err", fmt.Sprintf("%v", err)),
		)
		return
	}
	slog.Info("repository count", slog.Int("count", len(repositories)))
	toStore := []*std.Repository{}
	for i, rep := range repositories {
		slog.Info(fmt.Sprintf("[%d] %s", i, rep.GetFullName()))

		cmp := std.NewWithR(nil, rep, client)
		toStore = append(toStore, cmp)
	}
	content, err := data.ToJsonList[*std.Repository](toStore)
	filename := r.Filename()
	err = files.WriteFile(dir, filename, content)

	fmt.Println(err)
}

func main() {
	logger.LogSetup()
	costReport := report.New(orgArg, teamArg)
	costReport.SetRunner(run)
	costReport.Run()

}
