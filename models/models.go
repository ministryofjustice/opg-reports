package models

import (
	"context"
	"strings"
	"time"

	"github.com/google/go-github/v62/github"
	"github.com/ministryofjustice/opg-reports/internal/bools"
	"github.com/ministryofjustice/opg-reports/internal/dateformats"
)

// Full returns the known models that require a database table to be created
func Full() []interface{} {

	return []interface{}{
		&Unit{},                       // Unit is the base grouping model
		&AwsAccount{},                 // AwsAccount details attached to other aws models
		&AwsCost{},                    // AwsCosts model
		&AwsUptime{},                  // AwsUptime tracking
		&Dataset{},                    // Single record table to say if data is real or not
		&GitHubRepositoryGitHubTeam{}, // Many to many join between repo and teams
		&GitHubTeamUnit{},             // Many to many between guthub team and base units
		&GitHubTeam{},                 // GitHub team models used on other github models
		&GitHubRepository{},           // Github repo model
		&GitHubRelease{},              // Release model
		&GitHubRepositoryStandard{},   // Standards model
	}

}

// NewRepositoryFromRemote converts a github.Repository over to local version and fetches some additional
// innformation like license name and list of teams
func NewRepository(ctx context.Context, client *github.Client, r *github.Repository) (repo *GitHubRepository) {
	var ts = time.Now().UTC().Format(dateformats.Full)
	repo = &GitHubRepository{
		Ts:            ts,
		Owner:         r.GetOwner().GetLogin(),
		Name:          r.GetName(),
		FullName:      r.GetFullName(),
		CreatedAt:     r.GetCreatedAt().Format(dateformats.Full),
		DefaultBranch: r.GetDefaultBranch(),
		Archived:      bools.Int(r.GetArchived()),
		Private:       bools.Int(r.GetPrivate()),
	}
	// get the license and attach the name
	if l := r.GetLicense(); l != nil {
		repo.License = l.GetName()
	}
	// get the default branch and grab the dates
	if branch, _, err := client.Repositories.GetBranch(ctx, repo.Owner, repo.Name, repo.DefaultBranch, 1); err == nil {
		repo.LastCommitDate = branch.Commit.Commit.Author.Date.Time.String()
	}

	repo.GitHubTeams, _ = NewGitHubTeams(ctx, client, repo.Owner, repo.Name)

	return
}

// NewRepositoryStandard uses the repository details passed to create a new Local GitHubRepositoryStandard with
// an attached GitHubRepository (which in turn will have GitHubTeams) and populates all the compliance related
// fields
func NewRepositoryStandard(ctx context.Context, client *github.Client, r *github.Repository) (g *GitHubRepositoryStandard) {
	const (
		readmePath            string = "./README.md"
		codeOfConductPath     string = "./CODE_OF_CONDUCT.md"
		contributingGuidePath string = "./CONTRIBUTING.md"
	)
	var repo = NewRepository(ctx, client, r)

	g = &GitHubRepositoryStandard{
		Ts:                       time.Now().UTC().Format(dateformats.Full),
		DefaultBranch:            repo.DefaultBranch,
		LastCommitDate:           repo.LastCommitDate,
		GitHubRepository:         (*GitHubRepositoryForeignKey)(repo),
		IsArchived:               repo.Archived,
		IsPrivate:                repo.Private,
		GitHubRepositoryFullName: repo.FullName,
	}

	//
	g.HasDefaultBranchProtection = 0
	if branch, _, err := client.Repositories.GetBranch(ctx, repo.Owner, repo.Name, g.DefaultBranch, 1); err == nil {
		if branch.GetProtected() {
			g.HasDefaultBranchProtection = 1
		}
	}
	// -- counters
	if clones, _, err := client.Repositories.ListTrafficClones(
		context.Background(), repo.Owner, repo.Name,
		&github.TrafficBreakdownOptions{Per: "day"}); err == nil {

		g.CountOfClones = clones.GetCount()
	}
	g.CountOfForks = r.GetForksCount()
	if prs, _, err := client.PullRequests.List(ctx, repo.Owner, repo.Name,
		&github.PullRequestListOptions{State: "open"}); err == nil {
		g.CountOfPullRequests = len(prs)
	}
	if hooks, _, err := client.Repositories.ListHooks(ctx, repo.Owner, repo.Name,
		&github.ListOptions{PerPage: 100}); err == nil {
		g.CountOfWebHooks = len(hooks)
	}
	// -- has
	g.HasCodeOfConduct = 0
	if _, _, _, err := client.Repositories.GetContents(ctx, repo.Owner, repo.Name, codeOfConductPath, nil); err == nil {
		g.HasCodeOfConduct = 1
	}
	g.HasContributingGuide = 0
	if _, _, _, err := client.Repositories.GetContents(ctx, repo.Owner, repo.Name, contributingGuidePath, nil); err == nil {
		g.HasContributingGuide = 1
	}
	g.HasReadme = 0
	if _, _, _, err := client.Repositories.GetContents(ctx, repo.Owner, repo.Name,
		readmePath, nil); err == nil {
		g.HasReadme = 1
	}

	g.HasVulnerabilityAlerts = 0
	if alerts, _, err := client.Repositories.GetVulnerabilityAlerts(ctx, repo.Owner, repo.Name); err == nil {
		g.HasVulnerabilityAlerts = bools.Int(alerts)
	}

	g.HasDefaultBranchOfMain = bools.Int((g.DefaultBranch == "main"))
	g.HasDescription = bools.Int((len(r.GetDescription()) > 0))
	g.HasDiscussions = bools.Int(r.GetHasDiscussions())
	g.HasDownloads = bools.Int(r.GetHasDownloads())
	g.HasIssues = bools.Int(r.GetHasIssues())
	g.HasLicense = bools.Int((len(repo.License) > 0))
	g.HasPages = bools.Int(r.GetHasPages())
	g.HasWiki = bools.Int(r.GetHasWiki())

	if protection, _, err := client.Repositories.GetBranchProtection(ctx, repo.Owner, repo.Name,
		g.DefaultBranch); err == nil {
		g.HasRulesEnforcedForAdmins = bools.Int(protection.EnforceAdmins.Enabled)
		g.HasPullRequestApprovalRequired = bools.Int(protection.RequiredPullRequestReviews.RequiredApprovingReviewCount > 0)
		g.HasCodeownerApprovalRequired = bools.Int(protection.RequiredPullRequestReviews.RequireCodeOwnerReviews)
	}

	// the GetDeleteBranchOnMerge seems to be empty and have to re-fetch the api to get result
	re, _, _ := client.Repositories.Get(ctx, repo.Owner, repo.Name)
	g.HasDeleteBranchOnMerge = bools.Int(re.GetDeleteBranchOnMerge())

	g.UpdateCompliance()
	return
}

// NewGitHubTeams returns a GitHubTeams ([]*GitHubTeam) of all the teams attached to the repo whose org and name is passed.
func NewGitHubTeams(ctx context.Context, client *github.Client, organisation string, repoName string) (teams GitHubTeams, err error) {
	teams = GitHubTeams{}
	opts := &github.ListOptions{PerPage: 100}

	if teamList, _, err := client.Repositories.ListTeams(ctx, organisation, repoName, opts); err == nil {
		for _, team := range teamList {
			var ts = time.Now().UTC().Format(dateformats.Full)
			teams = append(teams, &GitHubTeam{
				Ts:   ts,
				Slug: strings.ToLower(*team.Slug),
			})
		}
	}

	// Fetch and parse the code owner files - adding those values into teams

	return
}
