package infracostcli

import (
	"context"
	"log/slog"
	"opg-reports/report/internal/db/dbconnection"
	"opg-reports/report/internal/db/dbsetup"
	"opg-reports/report/internal/db/dbstmts"
	"opg-reports/report/internal/domain/infracosts/infracostgetter"
	"opg-reports/report/internal/domain/infracosts/infracostmodels"
	"opg-reports/report/internal/utils/awsclients"
	"opg-reports/report/internal/utils/awsid"
	"opg-reports/report/internal/utils/times"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/cobra"
)

const (
	name  string = "infracosts"
	short string = `infracosts fetches and imports cost data from AWS cost explorer [needs AWS_SESSION].`
)

// the command flags used on the import cli tool
type cli struct {
	// data base
	DBPath   string // represents --db
	DBDriver string // represents --driver
	// date ranges
	DateEnd   string // represents --end; is the end date for date we're getting data for, generally this morning
	DateStart string // represents --start; the first part of the date period
	// AWS related
	Region string // represents --region
}

// options to pass along to the getAndImport function
type options struct {
	AccountID string
	DateStart time.Time
	DateEnd   time.Time
}

// ctx / logs for the package
var (
	ctx   context.Context                 // default context
	log   *slog.Logger                    // default logger
	today time.Time       = times.Today() // today, used for time window setting
)

// default command options
var flags *cli = &cli{
	DBPath:   "./database/api.db",
	DBDriver: "sqlite3",
	Region:   "eu-west-1",
	DateEnd:  times.AsYMDString(today),
	DateStart: times.AsYMDString(
		times.Ago(times.ResetMonth(today), 2, times.MONTH)),
}

// main command
var cmd = &cobra.Command{
	Use:   name,
	Short: short,
	RunE:  runCmd,
}

func CMD(c context.Context, l *slog.Logger) *cobra.Command {
	ctx = c
	log = l
	return cmd
}

func init() {
	cmd.Flags().StringVar(&flags.DateEnd, "end", flags.DateEnd, "End date")
	cmd.Flags().StringVar(&flags.DateStart, "start", flags.DateStart, "Start date")
	cmd.Flags().StringVar(&flags.DBPath, "db", flags.DBPath, "Database path")
	cmd.Flags().StringVar(&flags.DBDriver, "driver", flags.DBDriver, "Database driver")
	cmd.Flags().StringVar(&flags.Region, "region", flags.Region, "AWS region")
}

// runCmd main runner
func runCmd(cmd *cobra.Command, args []string) (err error) {
	var db *sqlx.DB
	var client *costexplorer.Client
	// db connection
	db, err = dbconnection.Connection(ctx, log, flags.DBDriver, flags.DBPath)
	if err != nil {
		return
	}
	defer db.Close()
	// db migration before import
	err = dbsetup.Migrate(ctx, log, db)
	if err != nil {
		return
	}
	// aws client
	client, err = awsclients.New[*costexplorer.Client](ctx, log, flags.Region)
	if err != nil {
		return
	}
	return getAndImport(ctx, log, client, db, &options{
		AccountID: awsid.AccountID(ctx, log, flags.Region),
		DateStart: times.MustFromString(flags.DateStart),
		DateEnd:   times.MustFromString(flags.DateEnd),
	})

}

// getAndImport uses package getter to fetch and then insert data into the passed database
func getAndImport(ctx context.Context, log *slog.Logger, client infracostgetter.AwsClient, db *sqlx.DB, params *options) (err error) {
	var (
		result []*dbstmts.Insert[*infracostmodels.Cost, int]
		data   []*infracostmodels.Cost = []*infracostmodels.Cost{}
		lg     *slog.Logger            = log.With("func", "infracostcli.getAndImport", "account", params.AccountID)
	)

	lg.With("params", params).Info("starting infracost import ...")
	// run the data getter command
	data, err = infracostgetter.GetCostData(ctx, log, client, &infracostgetter.Options{
		AccountID: params.AccountID,
		Start:     params.DateStart,
		End:       params.DateEnd,
	})
	if err != nil {
		return
	}
	// write the data
	result, err = dbsetup.Import[int](ctx, log, db, data, nil)
	if err != nil {
		return
	}
	lg.With("count", len(result)).Info("complete.")

	return
}
