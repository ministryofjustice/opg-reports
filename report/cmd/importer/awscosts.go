package main

import (
	"opg-reports/report/internal/repository/awsr"
	"opg-reports/report/internal/repository/sqlr"
	"opg-reports/report/internal/service/api"
	"opg-reports/report/internal/utils"

	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer/types"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/spf13/cobra"
)

// awscostsCmd imports data from the cost explorer api directly
var awscostsCmd = &cobra.Command{
	Use:   "awscosts",
	Short: "awscosts fetches data from the cost explorer api",
	Long: `
awscosts will call the aws costexplorer api to retrieve data for period specific.

env variables used that can be adjusted:

	DATABASE_PATH
		The file path to the sqlite database that will be used

`,
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

func awscostsCmdRunner(
	stsClient awsr.ClientSTSCaller,
	stsStore awsr.RepositorySTS,
	ceClient awsr.ClientCostExplorerGetter,
	ceStore awsr.RepositoryCostExplorerGetter,
	sqClient sqlr.RepositoryWriter,
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
		Granularity: types.GranularityMonthly,
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
