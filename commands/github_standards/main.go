package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/google/go-github/v62/github"
	"github.com/ministryofjustice/opg-reports/commands/shared/argument"
	"github.com/ministryofjustice/opg-reports/datastore/github_standards/ghs"
	"github.com/ministryofjustice/opg-reports/shared/convert"
	"github.com/ministryofjustice/opg-reports/shared/dates"
	"github.com/ministryofjustice/opg-reports/shared/env"
	"github.com/ministryofjustice/opg-reports/shared/github/cl"
	"github.com/ministryofjustice/opg-reports/shared/github/repos"
	"github.com/ministryofjustice/opg-reports/shared/logger"
)

const readmePath string = "./README.md"
const codeOfConductPath string = "./CODE_OF_CONDUCT.md"
const contributingGuidePath string = "./CONTRIBUTING.md"

// mapFromApi generates a GithubStandard item from api data
// - its a chunky one
func mapFromApi(ctx context.Context, client *github.Client, repo *github.Repository) (g *ghs.GithubStandard) {
	g = &ghs.GithubStandard{
		Ts: time.Now().UTC().String(),
	}

	g.DefaultBranch = repo.GetDefaultBranch()
	g.FullName = repo.GetFullName()
	g.Name = repo.GetName()
	g.Owner = repo.GetOwner().GetLogin()
	g.CreatedAt = repo.GetCreatedAt().Format(dates.Format)
	//
	if l := repo.GetLicense(); l != nil {
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
	g.CountOfForks = repo.GetForksCount()
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
		g.HasVulnerabilityAlerts = convert.BoolToInt(alerts)
	}

	g.HasDefaultBranchOfMain = convert.BoolToInt((g.DefaultBranch == "main"))
	g.HasDescription = convert.BoolToInt((len(repo.GetDescription()) > 0))
	g.HasDiscussions = convert.BoolToInt(repo.GetHasDiscussions())
	g.HasDownloads = convert.BoolToInt(repo.GetHasDownloads())
	g.HasIssues = convert.BoolToInt(repo.GetHasIssues())
	g.HasLicense = convert.BoolToInt((len(g.License) > 0))
	g.HasPages = convert.BoolToInt(repo.GetHasPages())
	g.HasWiki = convert.BoolToInt(repo.GetHasWiki())

	if protection, _, err := client.Repositories.GetBranchProtection(ctx, g.Owner, g.Name,
		g.DefaultBranch); err == nil {
		g.HasRulesEnforcedForAdmins = convert.BoolToInt(protection.EnforceAdmins.Enabled)
		g.HasPullRequestApprovalRequired = convert.BoolToInt(protection.RequiredPullRequestReviews.RequiredApprovingReviewCount > 0)
		g.HasCodeownerApprovalRequired = convert.BoolToInt(protection.RequiredPullRequestReviews.RequireCodeOwnerReviews)
	}

	g.IsArchived = convert.BoolToInt(repo.GetArchived())
	g.IsPrivate = convert.BoolToInt(repo.GetPrivate())

	// -- teams
	g.Teams = ""
	if teams, _, err := client.Repositories.ListTeams(ctx, g.Owner, g.Name,
		&github.ListOptions{PerPage: 100}); err == nil {
		for _, team := range teams {
			g.Teams += *team.Name + "#"
		}
	}
	// the GetDeleteBranchOnMerge seems to be empty and have to re-fetch the api to get result
	re, _, _ := client.Repositories.Get(ctx, g.Owner, g.Name)
	g.HasDeleteBranchOnMerge = convert.BoolToInt(re.GetDeleteBranchOnMerge())

	g.UpdateCompliance()
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
	out := argument.New(group, "output", "github_standards.json", "filename")
	d := argument.New(group, "dir", "data", "sub dir")

	group.Parse(os.Args[1:])

	dir := fmt.Sprintf("./%s/", *d.Value)
	os.MkdirAll(dir, os.ModePerm)
	filename := filepath.Join(dir, *out.Value)

	slog.Info("getting standards",
		slog.String("org", *org.Value),
		slog.String("team", *team.Value),
		slog.String("output file", *out.Value))

	ctx := context.Background()
	limiter, _ := cl.RateLimitedHttpClient()
	client := cl.Client(token, limiter)

	repositories, err := repos.All(ctx, client, *org.Value, *team.Value, true)
	if err != nil {
		slog.Error(err.Error())
		return
	}

	mapped := []ghs.GithubStandard{}
	total := len(repositories)
	for i, r := range repositories {
		if *r.FullName == "ministryofjustice/opg-reports" {

			g := *mapFromApi(ctx, client, r)
			fmt.Println(convert.PrettyString(g))
		}
		slog.Info(fmt.Sprintf("[%d/%d] %s", i+1, total, *r.FullName))
		// mapped = append(mapped, *mapFromApi(ctx, client, r))
	}

	content, err := convert.Marshals(mapped)
	if err != nil {
		slog.Error("error marshaling", slog.String("err", err.Error()))
		os.Exit(1)
	}

	os.WriteFile(filename, content, os.ModePerm)

}
