// Package costsfront contains tools for the front end server to handle / transform
// api results.
package costsfront

import (
	"log/slog"

	"github.com/ministryofjustice/opg-reports/pkg/transformers"
	"github.com/ministryofjustice/opg-reports/sources/costs/costsio"
)

// TransformResult takes the result from the api and converts
// the data into table rows that can be used for the front
// end.
//
// `body` is one of 2 possible types:
//
//	TaxOverviewBody
//	StandardBody
//
// Any others will be ignored.
func TransformResult(body interface{}) (result interface{}) {
	var err error
	var res map[string]map[string]interface{}
	result = body

	switch body.(type) {
	case *costsio.CostsTaxOverviewBody:
		var taxBody = body.(*costsio.CostsTaxOverviewBody)
		if res, err = transformers.ResultsToRows(taxBody.Result, taxBody.ColumnValues, taxBody.DateRange); err == nil {
			taxBody.TableRows = res
			result = taxBody
		}
	case *costsio.CostsStandardBody:
		var standard = body.(*costsio.CostsStandardBody)
		if res, err = transformers.ResultsToRows(standard.Result, standard.ColumnValues, standard.DateRange); err == nil {
			standard.TableRows = res
			result = standard
		}
	}

	if err != nil {
		slog.Error("[costsfront.TransformResult] ", slog.String("err", err.Error()))
	}

	return
}
