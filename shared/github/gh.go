// Package gh provides top level helpers for github
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
