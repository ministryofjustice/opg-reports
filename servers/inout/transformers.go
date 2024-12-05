package inout

import (
	"fmt"
	"log/slog"

	"github.com/ministryofjustice/opg-reports/internal/transformers"
)

// TransformToDateWideTable takes the result from the api and converts
// the data into table rows that can be used for the front
// end with dates as column headers, so merging items into less rows
func TransformToDateWideTable(body interface{}) (result interface{}) {
	var err error
	var res map[string]map[string]interface{}
	result = body

	switch body.(type) {
	// -- AWS Costs
	case *AwsCostsTaxesBody:
		var bdy = body.(*AwsCostsTaxesBody)
		if res, err = transformers.ResultsToDateRows(bdy.Result, bdy.ColumnValues, bdy.DateRange); err == nil {
			bdy.TableRows = res
			result = bdy
		}
	case *AwsCostsSumBody:
		var bdy = body.(*AwsCostsSumBody)
		if res, err = transformers.ResultsToDateRows(bdy.Result, bdy.ColumnValues, bdy.DateRange); err == nil {
			bdy.TableRows = res
			result = bdy
		}
	case *AwsCostsSumPerUnitBody:
		var bdy = body.(*AwsCostsSumPerUnitBody)
		if res, err = transformers.ResultsToDateRows(bdy.Result, bdy.ColumnValues, bdy.DateRange); err == nil {
			bdy.TableRows = res
			result = bdy
		}
	case *AwsCostsSumPerUnitEnvBody:
		var bdy = body.(*AwsCostsSumPerUnitEnvBody)
		if res, err = transformers.ResultsToDateRows(bdy.Result, bdy.ColumnValues, bdy.DateRange); err == nil {
			bdy.TableRows = res
			result = bdy
		}
	case *AwsCostsSumFullDetailsBody:
		var bdy = body.(*AwsCostsSumFullDetailsBody)
		if res, err = transformers.ResultsToDateRows(bdy.Result, bdy.ColumnValues, bdy.DateRange); err == nil {
			bdy.TableRows = res
			result = bdy
		}
	// -- AWS Uptime
	case *AwsUptimeAveragesBody:
		var bdy = body.(*AwsUptimeAveragesBody)
		if res, err = transformers.ResultsToDateRows(bdy.Result, bdy.ColumnValues, bdy.DateRange); err == nil {
			bdy.TableRows = res
			result = bdy
		}
	case *AwsUptimeAveragesPerUnitBody:
		var bdy = body.(*AwsUptimeAveragesPerUnitBody)
		if res, err = transformers.ResultsToDateRows(bdy.Result, bdy.ColumnValues, bdy.DateRange); err == nil {
			bdy.TableRows = res
			result = bdy
		}
	// -- GitHub Releases
	case *GitHubReleasesCountBody:
		var bdy = body.(*GitHubReleasesCountBody)
		if res, err = transformers.ResultsToDateRows(bdy.Result, bdy.ColumnValues, bdy.DateRange); err == nil {
			bdy.TableRows = res
			result = bdy
		}
	case *GitHubReleasesCountPerUnitBody:
		var bdy = body.(*GitHubReleasesCountPerUnitBody)
		if res, err = transformers.ResultsToDateRows(bdy.Result, bdy.ColumnValues, bdy.DateRange); err == nil {
			bdy.TableRows = res
			result = bdy
		}

	}

	if err != nil {
		slog.Error("[transformers] api transform error", slog.String("err", err.Error()), slog.String("type", fmt.Sprintf("%T", body)))
	}

	return
}

// func TransformToDateDeepTable(body interface{}) (result interface{}) {
// 	var err error
// 	var res map[string]map[string]interface{}
// 	result = body

// }
