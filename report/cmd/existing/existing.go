package main

import (
	"opg-reports/report/internal/repository/awsr"
	"opg-reports/report/internal/repository/githubr"
	"opg-reports/report/internal/repository/sqlr"
	"opg-reports/report/internal/service/existing"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/spf13/cobra"
)

// existingCmd imports all the currently know and supported previous data
// from earlier versions of reporting that are mostly stored in s3 buckets
var existingCmd = &cobra.Command{
	Use:   "existing",
	Short: "existing imports all known existing data files.",
	Long: `
existing imports all known data files (generally json) from a mix of sources (github, s3 buckets) that covers current and prior reporting data to ensure completeness.

env variables used that can be adjusted:

	EXISTING_COSTS_BUCKET
		The name of the bucket that current stores older aws cost data
	EXISTING_COSTS_PREFIX
		The bucket folder path (needs trailing /) with all cost data files within
	DATABASE_PATH
		The file path to the sqlite database that will be used
	GITHUB_ORGANISATION
		The name of the github organisation that owns the private repo
	METADATA_REPOSITORY
		The name of the repository to fetch release asset from for team / aws account lists
	METADATA_ASSET
		The name of the asset to download from the latest release on the repository
`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		var (
			githubClient                     = githubr.DefaultClient(conf)
			githubStore  *githubr.Repository = githubr.Default(ctx, log, conf)
			s3Client                         = awsr.DefaultClient[*s3.Client](ctx, "eu-west-1")
			s3Store      *awsr.Repository    = awsr.Default(ctx, log, conf)
			sqlStore     *sqlr.Repository    = sqlr.Default(ctx, log, conf)
			existService *existing.Service   = existing.Default(ctx, log, conf)
		)

		// TEAMS
		if _, err = existService.InsertTeams(githubClient.Repositories, githubStore, sqlStore); err != nil {
			return
		}
		// ACCOUNTS
		if _, err = existService.InsertAwsAccounts(githubClient.Repositories, githubStore, sqlStore); err != nil {
			return
		}
		// COSTS, only ig set
		if flagIncludeCosts {
			if _, err = existService.InsertAwsCosts(s3Client, s3Store, sqlStore); err != nil {
				return
			}
		}

		return
	},
}
