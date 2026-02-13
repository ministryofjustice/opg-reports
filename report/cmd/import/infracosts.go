package main

import (
	"context"
	"log/slog"
	"opg-reports/report/internal/db/dbsetup"
	"opg-reports/report/internal/db/dbstmts"
	"opg-reports/report/internal/domain/infracosts/infracost"
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
	infracostsShortDesc string = `infracosts fetches and imports cost data from AWS cost explorer.`
	infracostsLongDesc  string = `
infracosts fetches and imports cost data from AWS cost explorer based on the month passed along.

By default, it will also fetch the previous months data as well to ensure billing stability.
`
)

var (
	includePreviousMonth bool   = true        // represents --include-previous-month
	costsMonth           string = ""          // represents --month="YYYY-MM-DD"
	costsRegion          string = "eu-west-1" // represents --region
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
	var region string = costsRegion

	// db connection
	db, err = dbconn(ctx, log)
	if err != nil {
		return
	}
	defer db.Close()

	client, err = awsclients.New[*costexplorer.Client](ctx, log, region)

	return importInfracosts(ctx, log, client, db, &InfraOpts{
		AccountID:            awsid.AccountID(ctx, log, region),
		EndDate:              costsMonth,
		IncludePreviousMonth: includePreviousMonth,
	})
}

// main import command
func importInfracosts(ctx context.Context, log *slog.Logger, client infracost.AwsClient, db *sqlx.DB, in *InfraOpts) (err error) {
	var (
		result []*dbstmts.Insert[*infracostmodels.Cost, int]
		data   []*infracostmodels.Cost = []*infracostmodels.Cost{}
		lg     *slog.Logger            = log.With("func", "import.importInfracosts", "account", in.AccountID)
		opts   *infracost.Options      = &infracost.Options{AccountID: in.AccountID}
	)

	lg.Info("starting infracost import command ...")
	// work out dates
	opts.End, err = times.FromString(in.EndDate)
	if err != nil {
		lg.Error("error parsing month value.", "err", err.Error())
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
	lg.Debug("time period ...", "start", opts.Start, "end", opts.End)

	// fetch the data
	data, err = infracost.GetCostData(ctx, log, client, opts)
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

// add params to the command
func init() {
	infracostsCmd.Flags().StringVar(&costsMonth, "month", times.AsYMDString(times.ResetMonth(time.Now().UTC())), "The month to get cost data for. (YYYY-MM-DD)")
	infracostsCmd.Flags().StringVar(&costsRegion, "region", costsRegion, "The AWS region to fetch data from.")
	infracostsCmd.Flags().BoolVar(&includePreviousMonth, "include-previous-month", true, "When enabled the command will also run the previous month to ensure costs are accurate.")
}
