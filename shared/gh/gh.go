package gh

import (
	"context"

	"github.com/google/go-github/v62/github"
)

func Org(ctx context.Context, client *github.Client, slug string) (org *github.Organization, err error) {
	org, _, err = client.Organizations.Get(ctx, slug)
	return
}

func Team(ctx context.Context, client *github.Client, org string, slug string) (team *github.Team, err error) {
	team, _, err = client.Teams.GetTeamBySlug(ctx, org, slug)
	return
}

func Branch(ctx context.Context, client *github.Client, owner string, repo string, branch string) (b *github.Branch, err error) {
	b, _, err = client.Repositories.GetBranch(ctx, owner, repo, branch, 1)
	return
}

func HasVulnerabilityAlerts(ctx context.Context, client *github.Client, owner string, repo string) (enabled bool, err error) {
	enabled, _, err = client.Repositories.GetVulnerabilityAlerts(ctx, owner, repo)
	return
}

func HasReadme(ctx context.Context, client *github.Client, owner string, repo string) (present bool, err error) {
	content, _, err := client.Repositories.GetReadme(ctx, owner, repo, nil)
	if content != nil {
		present = true
	}
	return
}

func HasFile(ctx context.Context, client *github.Client, owner string, repo string, path string) (present bool, err error) {
	content, _, _, err := client.Repositories.GetContents(ctx, owner, repo, path, nil)
	if content != nil {
		present = true
	}
	return
}
