/*
Import data into a database

	import [COMMAND]
*/
package main

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/ministryofjustice/opg-reports/report/config"
	"github.com/ministryofjustice/opg-reports/report/internal/repository/awsr"
	"github.com/ministryofjustice/opg-reports/report/internal/repository/githubr"
	"github.com/ministryofjustice/opg-reports/report/internal/repository/sqlr"
	"github.com/ministryofjustice/opg-reports/report/internal/service/api"
	"github.com/ministryofjustice/opg-reports/report/internal/service/existing"
	"github.com/ministryofjustice/opg-reports/report/internal/service/seed"
	"github.com/ministryofjustice/opg-reports/report/internal/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// set up in the init
var (
	conf      *config.Config
	viperConf *viper.Viper
	ctx       context.Context
	log       *slog.Logger
)

var (
	month string = ""
)

// root command
var rootCmd = &cobra.Command{
	Use:               "import",
	Short:             "Import",
	Long:              `import can populate database with fixture data ("fixtures"), fetch data from pre-existing json ("existing") or new data via specific external api's.`,
	CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
}

// existingCmd imports all the currently know and supported previous data
// from earlier versions of reporting that are mostly stored in s3 buckets
var existingCmd = &cobra.Command{
	Use:   "existing",
	Short: "existing imports all known existing data files.",
	Long:  `existing imports all known data files (generally json) from a mix of sources (github, s3 buckets) that covers current and prior reporting data to ensure completeness`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		var (
			githubClient                     = githubr.DefaultClient(conf)
			s3Client                         = awsr.DefaultClient[*s3.Client](ctx, "eu-west-1")
			githubStore  *githubr.Repository = githubr.Default(ctx, log, conf)
			s3Store      *awsr.Repository    = awsr.Default(ctx, log, conf)
			sqlStore     *sqlr.Repository    = sqlr.Default(ctx, log, conf)
			existService *existing.Service   = existing.Default(ctx, log, conf)
		)

		err = existingCmdRunner(githubClient.Repositories, githubStore, s3Client, s3Store, sqlStore, existService)

		return
	},
}

// seedCmd uses fixture / seed data to populate a fresh database which can then
// be used for local dev / testing
var seedCmd = &cobra.Command{
	Use:   "seed",
	Short: "seed inserts known test data.",
	Long:  `seed inserts known test data for use in development.`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		var (
			sqlStore    *sqlr.Repository = sqlr.Default(ctx, log, conf)
			seedService *seed.Service    = seed.Default(ctx, log, conf)
		)
		err = seedCmdRunner(sqlStore, seedService)

		return
	},
}

// dbDownloadCmd downloads the database from the s3 bucket to a temp file
// and then overwrites (using os.Rename) the configured database file.
var dbDownloadCmd = &cobra.Command{
	Use:   "dbdownload",
	Short: "dbdownload downloads the database from an s3 bucket to local file system",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		var (
			s3Client = awsr.DefaultClient[*s3.Client](ctx, "eu-west-1")
			awsStore = awsr.Default(ctx, log, conf)
		)
		err = dbDownloadCmdRunner(s3Client, awsStore)
		return
	},
}

var dbUploadCmd = &cobra.Command{
	Use:   "dbupload",
	Short: "dbupload uploads a local database to the s3 bucket",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		var (
			s3Client = awsr.DefaultClient[*s3.Client](ctx, "eu-west-1")
			awsStore = awsr.Default(ctx, log, conf)
		)
		err = dbUploadCmdRunner(s3Client, awsStore)
		return
	},
}

// awscostsCmd imports data from the cost explorer api directly
var awscostsCmd = &cobra.Command{
	Use:   "awscosts",
	Short: "awscosts fetches data from the cost explorer api",
	Long:  `awscosts will call the aws costexplorer api to retrieve data for period specific.`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		var (
			stsClient  = awsr.DefaultClient[*sts.Client](ctx, conf.Aws.GetRegion())
			ceClient   = awsr.DefaultClient[*costexplorer.Client](ctx, conf.Aws.GetRegion())
			awsStore   = awsr.Default(ctx, log, conf)
			sqClient   = sqlr.DefaultWithSelect[*api.AwsCost](ctx, log, conf)
			apiService = api.Default[*api.AwsCost](ctx, log, conf)
		)
		err = awscostsCmdRunner(stsClient, awsStore, ceClient, awsStore, sqClient, apiService)
		return
	},
}

func existingCmdRunner(
	githubClient githubr.ReleaseClient,
	githubStore githubr.ReleaseRepository,
	s3Client awsr.ClientS3ListAndGetter,
	s3Store awsr.RepositoryS3BucketDownloader,
	sqlStore sqlr.Writer,
	existService *existing.Service,
) (err error) {

	// TEAMS
	if _, err = existService.InsertTeams(githubClient, githubStore, sqlStore); err != nil {
		return
	}
	// ACCOUNTS
	if _, err = existService.InsertAwsAccounts(githubClient, githubStore, sqlStore); err != nil {
		return
	}
	// COSTS
	if _, err = existService.InsertAwsCosts(s3Client, s3Store, sqlStore); err != nil {
		return
	}

	return
}

func seedCmdRunner(
	sqlStore sqlr.Writer,
	seedService *seed.Service,
) (err error) {

	// TEAMS
	if _, err = seedService.Teams(sqlStore); err != nil {
		return
	}
	// ACCOUNTS
	if _, err = seedService.AwsAccounts(sqlStore); err != nil {
		return
	}
	// COSTS
	if _, err = seedService.AwsCosts(sqlStore); err != nil {
		return
	}
	return
}

func awscostsCmdRunner(
	stsClient awsr.ClientSTSCaller,
	stsStore awsr.RepositorySTS,
	ceClient awsr.ClientCostExplorerGetter,
	ceStore awsr.RepositoryCostExplorerGetter,
	sqClient sqlr.Writer,
	apiService *api.Service[*api.AwsCost],

) (err error) {
	var (
		costs     = []*api.AwsCost{}
		caller, _ = stsStore.GetCallerIdentity(stsClient)
		start     = utils.StringToTimeReset(month, utils.TimeIntervalMonth)
		end       = start.AddDate(0, 1, 0)
	)
	opts := &awsr.GetCostDataOptions{
		StartDate:   start.Format(utils.DATE_FORMATS.YMD),
		EndDate:     end.Format(utils.DATE_FORMATS.YMD),
		Granularity: types.GranularityDaily,
	}
	// get the raw data from the api
	data, err := ceStore.GetCostData(ceClient, opts)
	if err != nil {
		return
	}
	// convert to AwsCosts struct
	err = utils.Convert(data, &costs)
	if err != nil {
		log.Error("error converting", "err", err.Error())
		return
	}
	// inject the account id into the cost records
	if caller != nil {
		for _, c := range costs {
			c.AwsAccountID = *caller.Account
		}
	}
	// insert
	_, err = apiService.PutAwsCosts(sqClient, costs)
	if err != nil {
		log.Error("error inserting", "err", err.Error())
		return
	}
	return
}

func dbDownloadCmdRunner(
	client awsr.ClientS3Getter,
	store awsr.RepositoryS3BucketItemDownloader,
) (err error) {
	var (
		dir, _ = os.MkdirTemp("./", "__download-s3-*")
		local  string
	)
	defer os.RemoveAll(dir)
	local, err = store.DownloadItemFromBucket(client, conf.Aws.Buckets.DB.Name, conf.Aws.Buckets.DB.Path(), dir)
	if err != nil {
		return
	}
	err = os.Rename(local, conf.Database.Path)
	return
}

func dbUploadCmdRunner(
	client awsr.ClientS3Putter,
	store awsr.RepositoryS3BucketItemUploader,
) (err error) {
	var (
		dir, _       = os.MkdirTemp("./", "__upload-s3-*")
		copyFrom     = conf.Database.Path
		copyTo       = filepath.Join(dir, filepath.Base(conf.Database.Path))
		targetBucket = conf.Aws.Buckets.DB.Name
		targetKey    = conf.Aws.Buckets.DB.Path()
		src          *os.File
	)
	// open the existing db file & copy to the new location
	src, err = os.Open(copyFrom)
	if err != nil {
		return
	}
	// copy...
	err = utils.FileCopy(src, copyTo)
	if err != nil {
		return
	}
	targetKey = "database/api2.db"
	// now upload the copy
	_, err = store.UploadItemToBucket(client, targetBucket, targetKey, copyTo)

	return
}

// init
func init() {
	conf, viperConf = config.New()
	ctx = context.Background()
	log = utils.Logger(conf.Log.Level, conf.Log.Type)

	// extra options that aren't handled via config env values
	// awscosts - sync-db
	awscostsCmd.Flags().StringVar(&month, "month", utils.Month(-2), "The month to get cost data for. (YYYY-MM-DD)")
}

func main() {
	rootCmd.AddCommand(
		existingCmd, seedCmd,
		dbDownloadCmd, dbUploadCmd,
		awscostsCmd)
	rootCmd.Execute()

}
