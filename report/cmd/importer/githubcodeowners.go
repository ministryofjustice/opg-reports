package main

import (
	"fmt"
	"opg-reports/report/internal/repository/githubr"

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
		ghClient = githubr.DefaultClient(conf)
		ghStore  = githubr.Default(ctx, log, conf)
		// sqClient   = sqlr.DefaultWithSelect[*api.GithubCodeOwner](ctx, log, conf)
		// apiService = api.Default[*api.GithubCodeOwner](ctx, log, conf)

	)

	// get all the repos
	repositories, err = githubCodeOwnerGetRepos(ghClient.Teams, ghStore, org, parentTeam)
	if err != nil {
		return
	}
	// get code owners from all of the repos
	githubCodeOwnersFromRepos(ghClient.Repositories, ghStore, org, parentTeam, repositories)

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
	org string, parentTeam string,
	repositories []*github.Repository,
) (allOwners []string, err error) {

	var coOptions = &githubr.GetRepositoryOwnerOptions{FilterByParent: parentTeam}

	for _, repo := range repositories {
		var owners []string
		fmt.Println(*repo.FullName)

		// get all the
		owners, err = store.GetRepositoryOwners(client, repo, coOptions)
		if err != nil {
			return
		}

		for _, ow := range owners {
			fmt.Println(" --> " + ow)
		}
		fmt.Println("--")
	}

}
