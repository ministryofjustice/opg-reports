package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/google/go-github/v62/github"
	"github.com/google/uuid"
	"github.com/ministryofjustice/opg-reports/commands/shared/argument"
	"github.com/ministryofjustice/opg-reports/datastore/github_standards/ghs"
	"github.com/ministryofjustice/opg-reports/shared/env"
	"github.com/ministryofjustice/opg-reports/shared/github/cl"
	"github.com/ministryofjustice/opg-reports/shared/github/repos"
	"github.com/ministryofjustice/opg-reports/shared/logger"
)

const readmePath string = "./README.md"
const codeOfConductPath string = "./CODE_OF_CONDUCT.md"
const contributingGuidePath string = "./CONTRIBUTING.md"

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func mapFromApi(ctx context.Context, client *github.Client, r *github.Repository) (g *ghs.GithubStandard) {
	g = &ghs.GithubStandard{
		Uuid: uuid.NewString(),
		Ts:   time.Now().UTC().String(),
	}

	g.DefaultBranch = r.GetDefaultBranch()
	g.FullName = r.GetFullName()
	g.Name = r.GetName()
	g.Owner = r.GetOwner().GetLogin()
	r.GetCreatedAt()
	//
	if l := r.GetLicense(); l != nil {
		g.License = l.GetName()
	}
	//
	g.HasDefaultBranchProtection = 0
	if branch, _, err := client.Repositories.GetBranch(ctx, g.Owner, g.Name, g.DefaultBranch, 1); err == nil {
		g.LastCommitDate = branch.Commit.Commit.Author.Date.Time.String()
		if branch.GetProtected() {
			g.HasDefaultBranchProtection = 1
		}
	}
	// -- counters
	if clones, _, err := client.Repositories.ListTrafficClones(
		context.Background(), g.Owner, g.Name,
		&github.TrafficBreakdownOptions{Per: "day"}); err == nil {
		g.CountOfClones = clones.GetCount()
	}
	g.CountOfForks = r.GetForksCount()
	if prs, _, err := client.PullRequests.List(ctx, g.Owner, g.Name,
		&github.PullRequestListOptions{State: "open"}); err == nil {
		g.CountOfPullRequests = len(prs)
	}
	if hooks, _, err := client.Repositories.ListHooks(ctx, g.Owner, g.Name,
		&github.ListOptions{PerPage: 100}); err == nil {
		g.CountOfWebHooks = len(hooks)
	}
	// -- has
	g.HasCodeOfConduct = 0
	if _, _, _, err := client.Repositories.GetContents(ctx, g.Owner, g.Name,
		codeOfConductPath, nil); err == nil {
		g.HasCodeOfConduct = 1
	}
	g.HasContributingGuide = 0
	if _, _, _, err := client.Repositories.GetContents(ctx, g.Owner, g.Name,
		contributingGuidePath, nil); err == nil {
		g.HasContributingGuide = 1
	}
	g.HasReadme = 0
	if _, _, _, err := client.Repositories.GetContents(ctx, g.Owner, g.Name,
		readmePath, nil); err == nil {
		g.HasReadme = 1
	}
	g.HasVulnerabilityAlerts = 0
	if alerts, _, err := client.Repositories.GetVulnerabilityAlerts(ctx, g.Owner, g.Name); err == nil {
		g.HasVulnerabilityAlerts = boolToInt(alerts)
	}

	g.HasDefaultBranchOfMain = boolToInt((g.DefaultBranch == "main"))
	g.HasDeleteBranchOnMerge = boolToInt(r.GetDeleteBranchOnMerge())
	g.HasDescription = boolToInt((len(r.GetDescription()) > 0))
	g.HasDiscussions = boolToInt(r.GetHasDiscussions())
	g.HasDownloads = boolToInt(r.GetHasDownloads())
	g.HasIssues = boolToInt(r.GetHasIssues())
	g.HasLicense = boolToInt((len(g.License) > 0))
	g.HasPages = boolToInt(r.GetHasPages())
	g.HasWiki = boolToInt(r.GetHasWiki())

	if protection, _, err := client.Repositories.GetBranchProtection(ctx, g.Owner, g.Name,
		g.DefaultBranch); err == nil {
		g.HasRulesEnforcedForAdmins = boolToInt(protection.EnforceAdmins.Enabled)
		g.HasPullRequestApprovalRequired = boolToInt(protection.RequiredPullRequestReviews.RequiredApprovingReviewCount > 0)
		g.HasCodeownerApprovalRequired = boolToInt(protection.RequiredPullRequestReviews.RequireCodeOwnerReviews)
	}

	g.IsArchived = boolToInt(r.GetArchived())
	g.IsPrivate = boolToInt(r.GetPrivate())

	// -- teams
	g.Teams = ""
	if teams, _, err := client.Repositories.ListTeams(ctx, g.Owner, g.Name,
		&github.ListOptions{PerPage: 100}); err == nil {
		for _, team := range teams {
			g.Teams += *team.Name + "#"
		}
	}

	return
}

func main() {
	logger.LogSetup()
	token := env.Get("GITHUB_ACCESS_TOKEN", "")
	if token == "" {
		slog.Error("no github token found")
		return
	}

	group := flag.NewFlagSet("github_standards", flag.ExitOnError)
	org := argument.New(group, "organisation", "ministryofjustice", "Organisation slug")
	team := argument.New(group, "team", "opg", "Team slug")
	csv := argument.New(group, "output", "github_standards", "filename")
	d := argument.New(group, "dir", "github_standards", "sub dir")

	group.Parse(os.Args[1:])

	dir := fmt.Sprintf("./%s/", *d.Value)
	os.MkdirAll(dir, os.ModePerm)
	filename := dir + *csv.Value + ".csv"

	slog.Info("getting standards...",
		slog.String("org", *org.Value),
		slog.String("team", *team.Value),
		slog.String("csv", *csv.Value))

	ctx := context.Background()
	limiter, _ := cl.RateLimitedHttpClient()
	client := cl.Client(token, limiter)

	repositories, err := repos.All(ctx, client, *org.Value, *team.Value, true)
	if err != nil {
		slog.Error(err.Error())
		return
	}

	f, _ := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, os.ModePerm)
	defer f.Close()
	for i, repo := range repositories {
		slog.Info(fmt.Sprintf("[%d] %s", i+1, repo.GetFullName()))
		std := mapFromApi(ctx, client, repo)
		f.WriteString(std.ToCSV())
	}

}
