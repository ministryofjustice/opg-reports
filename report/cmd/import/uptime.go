package main

import (
	"context"
	"log/slog"
	"opg-reports/report/internal/db/dbconnection"
	"opg-reports/report/internal/db/dbmigrations"
	"opg-reports/report/internal/db/dbstatements"
	"opg-reports/report/internal/domain/uptime/uptime"
	"opg-reports/report/internal/domain/uptime/uptimeimports"
	"opg-reports/report/internal/domain/uptime/uptimemodels"
	"opg-reports/report/internal/utils/awsid"
	"opg-reports/report/internal/utils/debugger"
	"opg-reports/report/internal/utils/times"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/cobra"
)

const (
	uptimeShortDesc string = `uptime fetches and imports uptime data from cloudwatch.`
	uptimeLongDesc  string = `
uptime fetches and imports uptime data from cloudwatch for the month passed along.

Uptime is determined via route53 health check data.
`
)

var (
	uptimeDayFlag string = "" // represents --day="YYYY-MM-DD"
)

var (
	uptimeCmd *cobra.Command = &cobra.Command{
		Use:   "uptime",
		Short: uptimeShortDesc,
		Long:  uptimeLongDesc,
		RunE:  uptimeRunE,
	}
)

type UptimeOpts struct {
	AccountID string
	Day       string
}

// cobra compatabile func
func uptimeRunE(cmd *cobra.Command, args []string) (err error) {
	var db *sqlx.DB
	var client *cloudwatch.Client
	var region string = "us-east-1" // hard coded region for uptime
	// db connection
	db, err = dbconnection.Connection(ctx, log, cfg.DB.Driver, cfg.DB.ConnectionString())
	if err != nil {
		return
	}

	return uptimeImport(ctx, log, client, db, &UptimeOpts{
		AccountID: awsid.AccountID(ctx, log, region),
		Day:       uptimeDayFlag,
	})
}

// main import command
func uptimeImport(ctx context.Context, log *slog.Logger, client uptime.AwsClient, db *sqlx.DB, in *UptimeOpts) (err error) {
	var (
		result []*dbstatements.InsertStatement[*uptimemodels.Uptime, int]
		data   []*uptimemodels.Uptime = []*uptimemodels.Uptime{}
		opts   *uptime.Options        = &uptime.Options{
			AccountID: in.AccountID,
		}
	)

	log = log.With("package", "import", "func", "uptimeImport", "account", opts.AccountID)
	log.Info("starting uptime import command ...")
	// close the db
	defer db.Close()

	// work out dates
	opts.Start, err = times.FromString(in.Day)
	if err != nil {
		return
	}
	// reset
	opts.End = times.ResetDay(times.Add(opts.Start, 1, times.DAY))
	log.Debug("time period ...", "start", opts.Start, "end", opts.End)

	err = dbmigrations.Migrate(ctx, log, db)
	if err != nil {
		return
	}
	// fetch the data
	data, err = uptime.GetUptimeData(ctx, log, client, opts)
	if err != nil {
		return
	}
	debugger.Dump(data)
	// write the data
	result, err = uptimeimports.Import(ctx, log, db, data)
	if err != nil {
		return
	}
	log.With("count", len(result)).Info("completed.")

	return
}

// add params to the command
func init() {
	uptimeCmd.Flags().StringVar(&uptimeDayFlag, "day", times.AsYMDString(times.Yesterday()), "The day to get uptime data for. (YYYY-MM-DD)")
}
