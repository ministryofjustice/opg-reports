package githubr

import (
	"log/slog"
	"slices"

	"github.com/google/go-github/v75/github"
)

type GetTeamsForRepositoryOptions struct {
	FilterByParentSlugs []string // only return repositories whose parent.slug is present
}

// GetTeamsForRepository returns all the teams attached to the repository
//
// Note: client is interface for wrapper for *github.RepositoriesService
func (self *Repository) GetTeamsForRepository(
	client *github.RepositoriesService,
	repo *github.Repository,
	options *GetTeamsForRepositoryOptions,
) (teams []*github.Team, err error) {
	var (
		ctx                          = self.ctx
		org                          = *repo.Owner.Login
		repoName                     = *repo.Name
		page     int                 = 1
		opts     *github.ListOptions = &github.ListOptions{PerPage: 200}
		log      *slog.Logger        = self.log.With("organistion", org, "repo", repoName, "operation", "GetOwnersForRepository")
	)

	teams = []*github.Team{}
	// loop over paginations
	for page > 0 {
		var response *github.Response
		var list []*github.Team

		opts.Page = page
		log.With("page", page).Debug("getting team list ... ")

		list, response, err = client.ListTeams(ctx, org, repoName, opts)
		if err != nil {
			log.Error("error getting team list", "err", err.Error())
			return
		}
		log.With("page", page, "count", len(list)).Debug("found teams ... ")
		// add to team list if it meets criteria
		if len(list) > 0 {
			for _, item := range list {
				var include = repositoryTeamMeetsCriteria(item, options)
				log.With("include", include, "team", *item.Name).Info("team checked ... ")
				if include {
					teams = append(teams, item)
				}
			}
		}
		// pagination
		page = response.NextPage
	}

	return
}

func repositoryTeamMeetsCriteria(team *github.Team, criteria *GetTeamsForRepositoryOptions) (pass bool) {
	pass = true
	if criteria == nil {
		return
	}
	// check parent slugs
	if len(criteria.FilterByParentSlugs) > 0 {
		// if asked for one, flip state to false before checking
		pass = false
		// if there is a parent, and its within the list, it passes
		if team.Parent != nil && slices.Contains(criteria.FilterByParentSlugs, *team.Parent.Slug) {
			pass = true
		}
	}

	return
}
