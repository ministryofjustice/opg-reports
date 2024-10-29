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
	"github.com/ministryofjustice/opg-reports/pkg/awscfg"
	"github.com/ministryofjustice/opg-reports/pkg/awsclient"
	"github.com/ministryofjustice/opg-reports/pkg/awssession"
	"github.com/ministryofjustice/opg-reports/pkg/consts"
	"github.com/ministryofjustice/opg-reports/pkg/convert"
	"github.com/ministryofjustice/opg-reports/sources/costs"
)

var (
	args   = &lib.Arguments{}
	awsCfg = awscfg.FromEnv()
)

func Run(args *lib.Arguments) (err error) {
	var (
		s         *session.Session
		client    *costexplorer.CostExplorer
		startDate time.Time
		endDate   time.Time
		raw       *costexplorer.GetCostAndUsageOutput
		data      []*costs.Cost
		content   []byte
		cfg       *awscfg.Config = awscfg.FromEnv()
	)

	if s, err = awssession.New(cfg); err != nil {
		slog.Error("[awscosts.main] aws session failed", slog.String("err", err.Error()))
		return
	}

	if client, err = awsclient.CostExplorer(s); err != nil {
		slog.Error("[awscosts.main] aws client failed", slog.String("err", err.Error()))
		return
	}

	if startDate, err = convert.ToTime(args.Month); err != nil {
		slog.Error("[awscosts.main] month conversion failed", slog.String("err", err.Error()))
		return
	}
	startDate = convert.DateResetMonth(startDate)
	// overwrite month with the parsed version
	args.Month = startDate.Format(consts.DateFormatYearMonth)
	endDate = startDate.AddDate(0, 1, 0)

	if raw, err = lib.CostData(client, startDate, endDate, costexplorer.GranularityDaily, consts.DateFormatYearMonthDay); err != nil {
		slog.Error("[awscosts.main] getting cost data failed", slog.String("err", err.Error()))
		return
	}
	if data, err = lib.Flat(raw, args); err != nil {
		slog.Error("[awscosts.main] flattening raw data to costs failed", slog.String("err", err.Error()))
		return
	}

	content, err = json.MarshalIndent(data, "", "  ")
	if err != nil {
		slog.Error("error marshaling", slog.String("err", err.Error()))
		os.Exit(1)
	}
	//
	lib.WriteToFile(content, args)

	return
}

func main() {
	var err error
	lib.SetupArgs(args)

	slog.Info("[awscosts.main] init...")
	slog.Debug("[awscosts.main]", slog.String("args", fmt.Sprintf("%+v", args)))
	slog.Debug("[awscosts.main]", slog.String("region", awsCfg.Region))

	if err = lib.ValidateArgs(args); err != nil {
		slog.Error("arg validation failed", slog.String("err", err.Error()))
		os.Exit(1)
	}

	Run(args)
	slog.Info("[awscosts.main] done.")

}
