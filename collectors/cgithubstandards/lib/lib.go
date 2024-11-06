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
	"github.com/ministryofjustice/opg-reports/pkg/consts"
	"github.com/ministryofjustice/opg-reports/pkg/convert"
	"github.com/ministryofjustice/opg-reports/sources/standards"
)

var (
	defOrg  = "ministryofjustice"
	defTeam = "opg"
)

// Arguments represents all the named arguments for this collector
type Arguments struct {
	Organisation string
	Team         string
	OutputFile   string
}

const readmePath string = "./README.md"
const codeOfConductPath string = "./CODE_OF_CONDUCT.md"
const contributingGuidePath string = "./CONTRIBUTING.md"

// SetupArgs maps flag values to properies on the arg passed and runs
// flag.Parse to fetch values
func SetupArgs(args *Arguments) {

	flag.StringVar(&args.Organisation, "organisation", defOrg, "organisation slug.")
	flag.StringVar(&args.Team, "team", defTeam, "team slug")
	flag.StringVar(&args.OutputFile, "output", "./data/github_standards.json", "Filepath for the output")

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

	return
}

// WriteToFile writes the content to the file
func WriteToFile(content []byte, args *Arguments) {
	var (
		filename string = args.OutputFile
		dir      string = filepath.Dir(args.OutputFile)
	)
	os.MkdirAll(dir, os.ModePerm)
	os.WriteFile(filename, content, os.ModePerm)

}

// AllRepos returns all accessible repos for the details passed
func AllRepos(ctx context.Context, client *github.Client, args *Arguments) (all []*github.Repository, err error) {
	var (
		org  string = args.Organisation
		team string = args.Team
		page int    = 1
	)

	all = []*github.Repository{}

	for page > 0 {
		slog.Info("getting repostiories", slog.Int("page", page))
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

// RepoToStandard generates a Standard item from api data
// - its a chunky one
func RepoToStandard(ctx context.Context, client *github.Client, repo *github.Repository) (g *standards.Standard) {

	g = &standards.Standard{
		Ts: time.Now().UTC().Format(consts.DateFormat),
	}

	g.DefaultBranch = repo.GetDefaultBranch()
	g.FullName = repo.GetFullName()
	g.Name = repo.GetName()
	g.Owner = repo.GetOwner().GetLogin()
	g.CreatedAt = repo.GetCreatedAt().Format(consts.DateFormat)
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
			g.Teams += "#" + strings.ToLower(*team.Name) + "#"
		}
	}
	// the GetDeleteBranchOnMerge seems to be empty and have to re-fetch the api to get result
	re, _, _ := client.Repositories.Get(ctx, g.Owner, g.Name)
	g.HasDeleteBranchOnMerge = convert.BoolToInt(re.GetDeleteBranchOnMerge())

	g.UpdateCompliance()
	return
}
