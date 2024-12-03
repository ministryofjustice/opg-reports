package transformers

import (
	"log/slog"

	"github.com/ministryofjustice/opg-reports/internal/transformers"
	"github.com/ministryofjustice/opg-reports/servers/inout"
)

// Transform takes the result from the api and converts
// the data into table rows that can be used for the front
// end.
func Transform(body interface{}) (result interface{}) {
	var err error
	var res map[string]map[string]interface{}
	result = body

	switch body.(type) {
	// -- AWS Costs
	case *inout.AwsCostsTaxesBody:
		var bdy = body.(*inout.AwsCostsTaxesBody)
		if res, err = transformers.ResultsToRows(bdy.Result, bdy.ColumnValues, bdy.DateRange); err == nil {
			bdy.TableRows = res
			result = bdy
		}
	case *inout.AwsCostsSumBody:
		var bdy = body.(*inout.AwsCostsSumBody)
		if res, err = transformers.ResultsToRows(bdy.Result, bdy.ColumnValues, bdy.DateRange); err == nil {
			bdy.TableRows = res
			result = bdy
		}
	case *inout.AwsCostsSumPerUnitBody:
		var bdy = body.(*inout.AwsCostsSumPerUnitBody)
		if res, err = transformers.ResultsToRows(bdy.Result, bdy.ColumnValues, bdy.DateRange); err == nil {
			bdy.TableRows = res
			result = bdy
		}
	case *inout.AwsCostsSumPerUnitEnvBody:
		var bdy = body.(*inout.AwsCostsSumPerUnitEnvBody)
		if res, err = transformers.ResultsToRows(bdy.Result, bdy.ColumnValues, bdy.DateRange); err == nil {
			bdy.TableRows = res
			result = bdy
		}
	case *inout.AwsCostsSumFullDetailsBody:
		var bdy = body.(*inout.AwsCostsSumFullDetailsBody)
		if res, err = transformers.ResultsToRows(bdy.Result, bdy.ColumnValues, bdy.DateRange); err == nil {
			bdy.TableRows = res
			result = bdy
		}
	// -- AWS Uptime
	case *inout.AwsUptimeAveragesBody:
		var bdy = body.(*inout.AwsUptimeAveragesBody)
		if res, err = transformers.ResultsToRows(bdy.Result, bdy.ColumnValues, bdy.DateRange); err == nil {
			bdy.TableRows = res
			result = bdy
		}
	case *inout.AwsUptimeAveragesPerUnitBody:
		var bdy = body.(*inout.AwsUptimeAveragesPerUnitBody)
		if res, err = transformers.ResultsToRows(bdy.Result, bdy.ColumnValues, bdy.DateRange); err == nil {
			bdy.TableRows = res
			result = bdy
		}
	// -- GitHub Releases
	case *inout.GitHubReleasesCountBody:
		var bdy = body.(*inout.GitHubReleasesCountBody)
		if res, err = transformers.ResultsToRows(bdy.Result, bdy.ColumnValues, bdy.DateRange); err == nil {
			bdy.TableRows = res
			result = bdy
		}
	case *inout.GitHubReleasesCountPerUnitBody:
		var bdy = body.(*inout.GitHubReleasesCountPerUnitBody)
		if res, err = transformers.ResultsToRows(bdy.Result, bdy.ColumnValues, bdy.DateRange); err == nil {
			bdy.TableRows = res
			result = bdy
		}

	}

	if err != nil {
		slog.Error("[transformers] aws costs transform error", slog.String("err", err.Error()))
	}

	return
}
