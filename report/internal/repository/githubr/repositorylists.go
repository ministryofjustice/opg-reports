package githubr

import (
	"log/slog"

	"github.com/google/go-github/v77/github"
)

// GetRepositoriesForTeamOptions used to decide what repositories
// should be returned
type GetRepositoriesForTeamOptions struct {
	ExcludeArchived bool // if a repo is marked as archived, it will be excluded from the data
}

// GetRepositoriesForTeam returns all repositories for a team within the organisation passed along
// and can filter those results afterwards via the `options`
//
// Note: client is interface for a *github.TeamsService
func (self *Repository) GetRepositoriesForTeam(
	client ClientTeamListRepositories, // *github.TeamsService,
	organisation string,
	team string,
	options *GetRepositoriesForTeamOptions,
) (repositories []*github.Repository, err error) {

	var (
		page int                 = 1
		opts *github.ListOptions = &github.ListOptions{PerPage: 200}
		log  *slog.Logger        = self.log.With("organistion", organisation, "team", team, "operation", "GetRepositoriesForTeam")
	)
	repositories = []*github.Repository{}
	// loop over all results and handle the pagination on the data set
	for page > 0 {
		var response *github.Response
		var list []*github.Repository

		opts.Page = page
		log.With("page", page).Debug("getting repository list ... ")

		list, response, err = client.ListTeamReposBySlug(self.ctx, organisation, team, opts)
		if err != nil {
			log.Error("error getting repository list", "err", err.Error())
			return
		}

		log.With("page", page, "count", len(list)).Debug("found repositories ... ")
		// process returned items
		if len(list) > 0 {
			for _, item := range list {
				var include = repositoryMeetsCriteria(item, options)
				log.With("include", include, "repo", *item.FullName).Info("repo checked ... ")
				if include {
					repositories = append(repositories, item)
				}
			}
		}
		// page iteration - will be 0 when none left
		page = response.NextPage
	}

	return
}

func repositoryMeetsCriteria(repo *github.Repository, criteria *GetRepositoriesForTeamOptions) (pass bool) {
	pass = true
	if criteria == nil {
		return
	}
	// if we're exlcudeding archived & this is archived, dont include
	if criteria.ExcludeArchived && repo.Archived != nil && *repo.Archived == true {
		pass = false
	}
	return
}
