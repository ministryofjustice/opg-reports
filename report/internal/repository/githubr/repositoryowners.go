package githubr

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"slices"
	"strings"

	"github.com/google/go-github/v77/github"
)

// GetRepositoryOwnerOptions allow to filter team data
type GetRepositoryOwnerOptions struct {
	FilterByParent string   // only return repositories whose parent.slug is present - this is ignored for CODEOWNER files
	Exclude        []string // list of teams to ignore - in particular used to skip default teams that are added to everything
}

// GetRepositoryOwners returns a combined list of teams from the repository and the
// CODEOWNER file entries of into a slice of strings with org prefix
//
// It calls `GetTeamsForRepository` & `GetCodeOwnersForRepository`, merging those
// results and removing duplicates
//
// Does not check if the team or CODEOWNERS are valid, just parses the data
func (self *Repository) GetRepositoryOwners(
	client ClientRepositoryOwnership, // client *github.RepositoriesService,
	repo *github.Repository,
	options *GetRepositoryOwnerOptions) (owners []string, err error) {

	var (
		org      string       = *repo.Owner.Login
		repoName string       = *repo.Name
		log      *slog.Logger = self.log.With("organistion", org, "repo", repoName, "operation", "GetRepositoryOwners")
	)
	owners = []string{}
	// first, find the teams and merge those in
	log.Debug("getting teams ... ")
	teams, _ := self.GetTeamsForRepository(client, repo, options)
	for _, team := range teams {
		owners = append(owners, teamSlug(team))
	}
	// now add in code owners
	log.Debug("getting CODEOWNERS ... ")
	codeowners, _ := self.GetCodeOwnersForRepository(client, repo, options)
	owners = append(owners, codeowners...)
	// remove duplicates
	slices.Sort(owners)
	owners = slices.Compact(owners)

	return
}

// GetTeamsForRepository returns all the teams attached to the repository
//
// Note: client is interface for wrapper for *github.RepositoriesService
func (self *Repository) GetTeamsForRepository(
	client ClientRepositoryTeamList, // client *github.RepositoriesService,
	repo *github.Repository,
	options *GetRepositoryOwnerOptions,
) (teams []*github.Team, err error) {
	var (
		ctx      context.Context     = self.ctx
		org      string              = *repo.Owner.Login
		repoName string              = *repo.Name
		page     int                 = 1
		opts     *github.ListOptions = &github.ListOptions{PerPage: 200}
		log      *slog.Logger        = self.log.With("organistion", org, "repo", repoName, "operation", "GetTeamsForRepository")
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
				log.With("include", include, "team", *item.Name).Debug("team checked ... ")
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

// GetCodeOwnersForRepository returns all the CODEOWNERS found from within repo
// which can then be used along with teams to determine ownersip
//
// Note: client is interface wrapper for client *github.RepositoriesService
func (self *Repository) GetCodeOwnersForRepository(
	client ClientRepositoryCodeOwnerDownload, // client *github.RepositoriesService,
	repo *github.Repository,
	options *GetRepositoryOwnerOptions,
) (owners []string, err error) {
	var (
		ctx            context.Context = self.ctx
		org            string          = *repo.Owner.Login
		repoName       string          = *repo.Name
		log            *slog.Logger    = self.log.With("organistion", org, "repo", repoName, "operation", "GetCodeOwnersForRepository")
		codeOwnerFiles []string        = []string{"./CODEOWNERS", "./.github/CODEOWNERS"}
	)
	owners = []string{}

	for _, codeOwnerFile := range codeOwnerFiles {
		var lines []string
		log.With("codeOwnerFile", codeOwnerFile).Debug("trying to fetch CODEOWNERS")

		lines, err = getFileContent(client.DownloadContents(ctx, org, repoName, codeOwnerFile, nil))
		if err == nil && len(lines) > 0 {
			codeowners := codeOwnersFromLines(lines)
			// owners = append(owners, ...)
			for _, co := range codeowners {
				var include = repositoryCodeOwnerMeetsCriteria(co, options)
				log.With("include", include, "codeowner", co).Debug("codeowner checked ... ")
				if include {
					owners = append(owners, co)
				}
			}

		}
	}
	if err != nil {
		return
	}
	// remove duplicates
	slices.Sort(owners)
	owners = slices.Compact(owners)
	return
}

// repositoryTeamMeetsCriteria checks if the team settings meeting the asked for values.
// Normally used to do filtering that isnt supported at the api for the end point
func repositoryTeamMeetsCriteria(team *github.Team, criteria *GetRepositoryOwnerOptions) (pass bool) {
	pass = true
	if criteria == nil {
		return
	}
	// check parent slugs
	if len(criteria.FilterByParent) > 0 {
		// if asked for one, flip state to false before checking
		pass = false
		// if there is a parent, and its within the list, it passes
		if team.Parent != nil && criteria.FilterByParent == *team.Parent.Slug {
			pass = true
		}
	}
	// if the team is within the ignore list it doesnt pass
	if len(criteria.Exclude) > 0 {
		for _, skip := range criteria.Exclude {
			var slug = teamSlug(team)
			if skip == slug {
				pass = false
			}
		}
	}

	return
}

func teamSlug(team *github.Team) string {
	var teamSlug = *team.Slug
	if o := team.GetOrganization(); o != nil && o.Login != nil {
		teamSlug = fmt.Sprintf("%s/%s", *o.Login, *team.Slug)
	}
	return teamSlug
}

func repositoryCodeOwnerMeetsCriteria(owner string, criteria *GetRepositoryOwnerOptions) (pass bool) {
	pass = true
	if criteria == nil {
		return
	}
	// if the team is within the ignore list it doesnt pass
	if len(criteria.Exclude) > 0 {
		for _, skip := range criteria.Exclude {
			if skip == owner {
				pass = false
			}
		}
	}

	return
}

// getFileContent reads the content from the remote buffer for a github file
func getFileContent(content io.ReadCloser, response *github.Response, e error) (lines []string, err error) {

	if e != nil {
		err = e
		return
	}
	if response.StatusCode != http.StatusOK {
		err = fmt.Errorf("error from http request [%v] [%v]", response.StatusCode, response.Status)
		return
	}

	b, err := io.ReadAll(content)
	if err != nil {
		return
	}
	err = content.Close()
	if err != nil {
		return
	}
	// trim the last new line from the file
	lines = strings.Split(strings.TrimRight(string(b), "\n"), "\n")
	return
}

// codeOwnersFromLines find all the code owners from lines in the codeowners file
// Note: strips lead @ from the team slug name
func codeOwnersFromLines(lines []string) (owners []string) {
	owners = []string{}

	for _, line := range lines {
		exploded := strings.Split(line, " ")
		if len(exploded) > 1 {
			for _, segment := range exploded[1:] {
				if len(segment) > 0 && segment != " " {
					owners = append(owners, segment)
				}
			}
		}
	}
	// remove duplicates
	slices.Sort(owners)
	owners = slices.Compact(owners)
	// remove the @prefix
	for i, o := range owners {
		if o[0] == '@' {
			owners[i] = o[1:]
		}
	}
	return
}
