package main

import (
	"opg-reports/report/internal/repository/githubr"
	"opg-reports/report/internal/repository/sqlr"
	"opg-reports/report/internal/service/api"

	"github.com/google/go-github/v77/github"
	"github.com/spf13/cobra"
)

const ghCodeOwnerLongDesc string = `
githubcodeowners will call the github api to first fetch all the active repositories and then find the code owners of each.

env variables used that can be adjusted:

	DATABASE_PATH
		The file path to the sqlite database that will be used

`

var (
	githubcodeownersCmd *cobra.Command = &cobra.Command{
		Use:   "githubcodeowners",
		Short: "githubcodeowners fetches data from the github api",
		Long:  ghCodeOwnerLongDesc,
		RunE:  githubCodeOwnerRunner,
	} // githubcodeownersCmd imports data from the github api
)

// used by the cobra command (awscostsCmd) to process the cli request to fetch data from
// the aws api and import to local database
func githubCodeOwnerRunner(cmd *cobra.Command, args []string) (err error) {
	var (
		repositories []*github.Repository
		org          = "ministryofjustice"
		parentTeam   = "opg"
		// clients
		ghClient   = githubr.DefaultClient(conf)
		ghStore    = githubr.Default(ctx, log, conf)
		sqClient   = sqlr.DefaultWithSelect[*api.GithubCodeOwner](ctx, log, conf)
		apiService = api.Default[*api.GithubCodeOwner](ctx, log, conf)
	)

	// get all the repos
	repositories, err = githubCodeOwnerGetRepos(ghClient.Teams, ghStore, org, parentTeam)
	if err != nil {
		return
	}

	// get code owners from all of the repos
	opts := &githubr.GetRepositoryOwnerOptions{
		FilterByParent: parentTeam, Exclude: []string{"ministryofjustice/opg-webops"},
	}
	owners, err := githubCodeOwnersFromRepos(ghClient.Repositories, ghStore, repositories, opts)
	if err != nil {
		return
	}
	// TODO - manually map github code owner to a team

	err = githubCodeOwnersInsert(sqClient, apiService, owners)
	if err != nil {
		log.Error("error inserting", "err", err.Error())
		return
	}

	return
}

func githubCodeOwnersInsert(
	client sqlr.RepositoryWriter, //*sqlr.RepositoryWithSelect[*api.GithubCodeOwner],
	service *api.Service[*api.GithubCodeOwner],
	owners []*api.GithubCodeOwner,
) (err error) {
	_, err = service.TruncateAndPutGithubCodeOwners(client, owners)
	return
}

func githubCodeOwnerGetRepos(
	client githubr.ClientTeamListRepositories,
	store githubr.RepositoryListByTeam,
	org string, parentTeam string) (repositories []*github.Repository, err error) {

	var opts = &githubr.GetRepositoriesForTeamOptions{ExcludeArchived: true}
	repositories, err = store.GetRepositoriesForTeam(client, org, parentTeam, opts)
	return
}

func githubCodeOwnersFromRepos(
	client githubr.ClientRepositoryOwnership,
	store githubr.RepositoryOwnerGetter,
	repositories []*github.Repository,
	coOptions *githubr.GetRepositoryOwnerOptions,
) (allOwners []*api.GithubCodeOwner, err error) {
	var defaultOwner string = "NONE"

	allOwners = []*api.GithubCodeOwner{}
	for _, repo := range repositories {
		var owners []string
		// get all the
		owners, err = store.GetRepositoryOwners(client, repo, coOptions)
		if err != nil {
			return
		}
		// add adefault owner
		if len(owners) == 0 {
			owners = append(owners, defaultOwner)
		}
		for _, ow := range owners {
			var co = &api.GithubCodeOwner{Repository: *repo.FullName, CodeOwner: ow}
			allOwners = append(allOwners, co)
		}
	}

	return

}
