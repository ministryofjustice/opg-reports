package main

import (
	"opg-reports/report/internal/repository/awsr"
	"opg-reports/report/internal/repository/sqlr"
	"opg-reports/report/internal/service/api"
	"opg-reports/report/internal/utils"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer/types"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/spf13/cobra"
)

const costsLongDesc string = `
awscosts will call the aws costexplorer api to retrieve data for specified period and the month before.

When stablise is enabled (true by default), the previous month is also run to ensure cost data is accurate. This is due to billing dates being middle of the next month to allow for discount plans and so on.

env variables used that can be adjusted:

	DATABASE_PATH
		The file path to the sqlite database that will be used

`

var (
	stabliseCosts  bool           = true // represents --stabilise
	costsMonthFlag string         = ""   // represents --month="YYYY-MM-DD"
	awscostsCmd    *cobra.Command = &cobra.Command{
		Use:   "awscosts",
		Short: "awscosts fetches data from the cost explorer api",
		Long:  costsLongDesc,
		RunE:  awsCostsRunner,
	} // awscostsCmd imports data from the cost explorer api directly
)

// awsCostsRunner used by the cobra command (awscostsCmd) to process the cli request to fetch data from
// the aws api and import to local database
//
// If `stablise` flag is enabled then it will also fetch the previous month cost data
func awsCostsRunner(cmd *cobra.Command, args []string) (err error) {

	var (
		requestedMonth   time.Time = utils.StringToTimeReset(costsMonthFlag, utils.TimeIntervalMonth) // start of the month
		previousMonth    time.Time = utils.TimeAdd(requestedMonth, -1, utils.TimeIntervalMonth)       // previous month
		getPreviousMonth bool      = stabliseCosts                                                    // should we fetch previous month as well
	)
	// get info for this month
	err = awsCostsForMonth(requestedMonth)
	// check if should be getting one for the last month
	if err == nil && getPreviousMonth == true {
		err = awsCostsForMonth(previousMonth)
	}
	return

}

// awsCostsForMonth finds the account account id from the awssion (expects aws-vault or similar)
// as well as the cost data for all servies in the month and inserts into the database (set in the config).
func awsCostsForMonth(month time.Time) (err error) {
	var (
		// clients / sessions
		stsClient          = awsr.DefaultClient[*sts.Client](ctx, conf.Aws.GetRegion())
		costexplorerClient = awsr.DefaultClient[*costexplorer.Client](ctx, conf.Aws.GetRegion())
		awsStore           = awsr.Default(ctx, log, conf)
		sqClient           = sqlr.DefaultWithSelect[*api.AwsCost](ctx, log, conf)
		apiService         = api.Default[*api.AwsCost](ctx, log, conf)
		// costs & account id
		accountID string              = ""                    // account id from the caller identity
		costs     []map[string]string = []map[string]string{} // api costs converted to map

	)
	// get the active account id
	accountID, err = awsAccountID(stsClient, awsStore)
	if err != nil {
		return
	}
	// get the requested months data
	costs, err = awsCostsGetData(costexplorerClient, awsStore, month)
	if err != nil {
		return
	}
	// insert
	err = awsCostsInsert(sqClient, apiService, accountID, costs)
	return
}

// awsCostsGetData gets the raw cost data
func awsCostsGetData(
	client awsr.ClientCostExplorerGetter,
	store awsr.RepositoryCostExplorerGetter,
	start time.Time,
) (costs []map[string]string, err error) {
	var (
		end          = start.AddDate(0, 1, 0)
		startStr     = start.Format(utils.DATE_FORMATS.YMD)
		endStr       = end.Format(utils.DATE_FORMATS.YMD)
		groupService = "SERVICE"
		groupRegion  = "REGION"
		options      = &costexplorer.GetCostAndUsageInput{
			Granularity: types.GranularityMonthly,
			TimePeriod: &types.DateInterval{
				Start: &startStr,
				End:   &endStr,
			},
			Metrics: []string{"UnblendedCost"},
			GroupBy: []types.GroupDefinition{
				{Type: types.GroupDefinitionTypeDimension, Key: &groupService},
				{Type: types.GroupDefinitionTypeDimension, Key: &groupRegion},
			},
		}
	)
	log.With("start", start, "end", end).Info("Getting costs between dates ... ")
	// get the raw data from the api
	costs, err = store.GetCostData(client, options)
	return
}

// awsCostsInsert adds new data into the existing database for aws costs
func awsCostsInsert(
	client sqlr.RepositoryWriter,
	service *api.Service[*api.AwsCost],
	accountID string,
	apiCosts []map[string]string,
) (err error) {

	var dbCosts = []*api.AwsCost{}

	// convert to AwsCosts struct
	err = utils.Convert(apiCosts, &dbCosts)
	if err != nil {
		log.Error("error converting", "err", err.Error())
		return
	}

	// add account id into each row
	for _, c := range dbCosts {
		c.AwsAccountID = accountID
	}

	// insert
	_, err = service.PutAwsCosts(client, dbCosts)
	if err != nil {
		log.Error("error inserting", "err", err.Error())
		return
	}

	return
}

func init() {
	awscostsCmd.Flags().StringVar(&costsMonthFlag, "month", utils.StartOfMonth().Format(utils.DATE_FORMATS.YMD), "The month to get cost data for. (YYYY-MM-DD)")
	awscostsCmd.Flags().BoolVar(&stabliseCosts, "stablise", true, "When enabled the command will also run the previous month to ensure costs are accurate.")
}
