/*
awscosts fetches aws costs data for the month at a daily granularity.

Usage:

	awscosts [flags]

The flags are:

	-month=<yyyy-mm-dd>
		The month (formated as YYYY-MM-DD) to fetch data for.
		If set to "-", uses the current month.
		Defaults to the current month.
	-id=<account-id>
		The AWS account id as a string.
	-name=<name>
		The free entry string used for this account id.
		Example: TeamA production
	-label=<label>
		A string to describe what this account is in more detail.
		Example: TeamA production databases
	-environment=<development|pre-production|production>
		One of the following: development, pre-production, production.
		Default: production
	-unit=<unit>
		Team name for who owns the account.
	-organisation=<organisation>
		Name of the organsiation that looks after the account
	-output=<path-pattern>
		Path (with magic values) to the output file
		Default: `./data/{month}_{id}_aws_costs.json`

The command presumes an active, autherised session that can connect
to AWS cost explorer for the account specified. These are dynamically
fetched from environment variables.
*/
package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/costexplorer"
	"github.com/ministryofjustice/opg-reports/collectors/awscosts/lib"
	"github.com/ministryofjustice/opg-reports/internal/awsclient"
	"github.com/ministryofjustice/opg-reports/internal/awssession"
	"github.com/ministryofjustice/opg-reports/internal/dateformats"
	"github.com/ministryofjustice/opg-reports/internal/dateintervals"
	"github.com/ministryofjustice/opg-reports/internal/dateutils"
	"github.com/ministryofjustice/opg-reports/internal/envar"
	"github.com/ministryofjustice/opg-reports/models"
)

var (
	args   = &lib.Arguments{}
	region = envar.Get("AWS_DEFAULT_REGION", "eu-west-1")
)

func Run(args *lib.Arguments) (err error) {
	var (
		s         *session.Session
		client    *costexplorer.CostExplorer
		startDate time.Time
		endDate   time.Time
		raw       *costexplorer.GetCostAndUsageOutput
		data      []*models.AwsCost
		content   []byte
	)

	if s, err = awssession.New(); err != nil {
		slog.Error("[awscosts] aws session failed", slog.String("err", err.Error()))
		return
	}

	if client, err = awsclient.CostExplorer(s); err != nil {
		slog.Error("[awscosts] aws client failed", slog.String("err", err.Error()))
		return
	}

	if startDate, err = dateutils.Time(args.Month); err != nil {
		slog.Error("[awscosts] month conversion failed", slog.String("err", err.Error()))
		return
	}
	startDate = dateutils.Reset(startDate, dateintervals.Month)
	// overwrite month with the parsed version
	args.Month = startDate.Format(dateformats.YMD)
	endDate = startDate.AddDate(0, 1, 0)

	if raw, err = lib.CostData(client, startDate, endDate, costexplorer.GranularityDaily, dateformats.YMD); err != nil {
		slog.Error("[awscosts] getting cost data failed", slog.String("err", err.Error()))
		return
	}
	if data, err = lib.Flat(raw, args); err != nil {
		slog.Error("[awscosts] flattening raw data to costs failed", slog.String("err", err.Error()))
		return
	}

	content, err = json.MarshalIndent(data, "", "  ")
	if err != nil {
		slog.Error("[awscosts] error marshaling", slog.String("err", err.Error()))
		os.Exit(1)
	}
	//
	lib.WriteToFile(content, args)

	return
}

func main() {
	var err error
	lib.SetupArgs(args)

	slog.Info("[awscosts] starting...")
	slog.Debug("[awscosts]", slog.String("args", fmt.Sprintf("%+v", args)))
	slog.Debug("[awscosts]", slog.String("region", region))

	if err = lib.ValidateArgs(args); err != nil {
		slog.Error("arg validation failed", slog.String("err", err.Error()))
		os.Exit(1)
	}

	Run(args)
	slog.Info("[awscosts] done.")

}
