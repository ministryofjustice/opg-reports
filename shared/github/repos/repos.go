package repos

import (
	"context"
	"log/slog"

	"github.com/google/go-github/v62/github"
)

// Get a specific repository
func Get(ctx context.Context, client *github.Client, owner string, slug string) (r *github.Repository, err error) {
	r, _, err = client.Repositories.Get(ctx, owner, slug)
	return
}

// Return all the repositories for the org & team passed
// This handles the pagination
func All(ctx context.Context, client *github.Client, org string, team string, includeArchived bool) (all []*github.Repository, err error) {
	all = []*github.Repository{}
	list := []*github.Repository{}
	page := 1

	for page > 0 {
		slog.Info("getting repostiories", slog.Int("page", page))
		pg, resp, e := client.Teams.ListTeamReposBySlug(ctx, org, team, &github.ListOptions{PerPage: 100, Page: page})
		if e != nil {
			err = e
			return
		}
		list = append(list, pg...)
		page = resp.NextPage
	}

	if !includeArchived {
		for _, r := range list {
			if !*r.Archived {
				all = append(all, r)
			}
		}
	} else {
		all = list
	}

	return
}
