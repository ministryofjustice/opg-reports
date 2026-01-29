package main

import (
	"context"
	"log/slog"
	"opg-reports/report/internal/db/dbconnection"
	"opg-reports/report/internal/db/dbmigrations"
	"opg-reports/report/internal/db/dbstatements"
	"opg-reports/report/internal/domain/infracosts/infracost"
	"opg-reports/report/internal/domain/infracosts/infracostimports"
	"opg-reports/report/internal/domain/infracosts/infracostmodels"
	"opg-reports/report/internal/utils/awsid"
	"opg-reports/report/internal/utils/times"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/cobra"
)

const (
	infracostsShortDesc string = `infracosts fetches and imports cost data from AWS cost explorer.`
	infracostsLongDesc  string = `
infracosts fetches and imports cost data from AWS cost explorer based on the month passed along.

By default, it will also fetch the previous months data as well to ensure billing stability.
`
)

var (
	includePreviousMonth bool   = true // represents --include-previous-month
	costsMonthFlag       string = ""   // represents --month="YYYY-MM-DD"
)

var (
	infracostsCmd *cobra.Command = &cobra.Command{
		Use:   "infracosts",
		Short: infracostsShortDesc,
		Long:  infracostsLongDesc,
		RunE:  infracostsRunE,
	}
)

type InfraOpts struct {
	AccountID            string
	EndDate              string
	IncludePreviousMonth bool
}

// wrapper to use with cobra
func infracostsRunE(cmd *cobra.Command, args []string) (err error) {
	var db *sqlx.DB
	var client *costexplorer.Client
	var region string = cfg.AWS.Region
	// db connection
	db, err = dbconnection.Connection(ctx, log, cfg.DB.Driver, cfg.DB.ConnectionString())
	if err != nil {
		return
	}

	return infracostsImport(ctx, log, client, db, &InfraOpts{
		AccountID:            awsid.AccountID(ctx, log, region),
		EndDate:              costsMonthFlag,
		IncludePreviousMonth: includePreviousMonth,
	})
}

// main import command
func infracostsImport(ctx context.Context, log *slog.Logger, client infracost.AwsClient, db *sqlx.DB, in *InfraOpts) (err error) {
	var (
		result []*dbstatements.InsertStatement[*infracostmodels.Cost, int]
		data   []*infracostmodels.Cost = []*infracostmodels.Cost{}
		opts   *infracost.Options      = &infracost.Options{
			AccountID: in.AccountID,
		}
	)

	log = log.With("package", "import", "func", "infracostsImport", "account", opts.AccountID)
	log.Info("starting infracost import command ...")
	// close the db
	defer db.Close()

	// work out dates
	opts.End, err = times.FromString(in.EndDate)
	if err != nil {
		return
	}
	// reset to start of the month
	opts.End = times.ResetMonth(opts.End)
	// fix the start date
	if in.IncludePreviousMonth {
		opts.Start = times.Add(opts.End, -2, times.MONTH)
	} else {
		opts.Start = times.Add(opts.End, -1, times.MONTH)
	}
	log.Debug("time period ...", "start", opts.Start, "end", opts.End)

	err = dbmigrations.Migrate(ctx, log, db)
	if err != nil {
		return
	}
	// fetch the data
	data, err = infracost.GetCostData(ctx, log, client, opts)
	if err != nil {
		return
	}
	// write the data
	result, err = infracostimports.Import(ctx, log, db, data)
	if err != nil {
		return
	}
	log.With("count", len(result)).Info("completed.")

	return
}

// add params to the command
func init() {
	infracostsCmd.Flags().StringVar(&costsMonthFlag, "month", times.AsYMDString(times.ResetMonth(time.Now().UTC())), "The month to get cost data for. (YYYY-MM-DD)")
	infracostsCmd.Flags().BoolVar(&includePreviousMonth, "include-previous-month", true, "When enabled the command will also run the previous month to ensure costs are accurate.")
}
