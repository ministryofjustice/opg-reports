package uptimecli

import (
	"context"
	"log/slog"
	"opg-reports/report/internal/db/dbconnection"
	"opg-reports/report/internal/db/dbsetup"
	"opg-reports/report/internal/db/dbstmts"
	"opg-reports/report/internal/domain/uptime/uptimegetter"
	"opg-reports/report/internal/domain/uptime/uptimemodels"
	"opg-reports/report/internal/utils/awsclients"
	"opg-reports/report/internal/utils/awsid"
	"opg-reports/report/internal/utils/times"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/cobra"
)

const (
	name  string = "uptime"
	short string = `uptime fetches and imports uptime data from cloudwatch [needs AWS_SESSION].`
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
	Region string // fixed to us-east-1
}

// options to pass along to the getAndImport function
type options struct {
	AccountID string
	DateStart time.Time
	DateEnd   time.Time
}

// ctx / logs for the package
var (
	ctx       context.Context                     // default context
	log       *slog.Logger                        // default logger
	yesterday time.Time       = times.Yesterday() // today, used for time window setting
)

// default command options
var flags *cli = &cli{
	Region:   "us-east-1", // region is fixed to this value for getting heatlh checks
	DBPath:   "./database/api.db",
	DBDriver: "sqlite3",
	DateEnd:  times.AsYMDString(yesterday),
	DateStart: times.AsYMDString(
		times.Ago(times.ResetMonth(yesterday), 1, times.MONTH)),
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
}

// runCmd main runner
func runCmd(cmd *cobra.Command, args []string) (err error) {
	var db *sqlx.DB
	var client *cloudwatch.Client
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
	client, err = awsclients.New[*cloudwatch.Client](ctx, log, flags.Region)
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
func getAndImport(ctx context.Context, log *slog.Logger, client uptimegetter.AwsClient, db *sqlx.DB, params *options) (err error) {
	var (
		result []*dbstmts.Insert[*uptimemodels.Uptime, int]
		data   []*uptimemodels.Uptime = []*uptimemodels.Uptime{}
		lg     *slog.Logger           = log.With("func", "uptimecli.getAndImport", "account", params.AccountID)
	)

	lg.With("params", params).Info("starting uptime import ...")
	// run the data getter command
	data, err = uptimegetter.GetUptimeData(ctx, log, client, &uptimegetter.Options{
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
