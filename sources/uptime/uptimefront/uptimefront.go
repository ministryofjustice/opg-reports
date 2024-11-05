// Package uptimefront handles transforming api data into tabluar
package uptimefront

import (
	"log/slog"

	"github.com/ministryofjustice/opg-reports/pkg/transformers"
	"github.com/ministryofjustice/opg-reports/sources/uptime/uptimeio"
)

// TransformResult takes the result from the api and converts
// the data into table rows that can be used for the front
// end.
//
// `body` can only be `UptimeBody` - any others will be ignored.
func TransformResult(body interface{}) (result interface{}) {
	var err error
	var res map[string]map[string]interface{}
	result = body

	switch body.(type) {
	case *uptimeio.UptimeBody:
		var uptime = body.(*uptimeio.UptimeBody)
		if res, err = transformers.ResultsToRows(uptime.Result, uptime.ColumnValues, uptime.DateRange); err == nil {
			uptime.TableRows = res
			result = uptime
		}
	}

	if err != nil {
		slog.Error("[uptimefront.TransformResult] ", slog.String("err", err.Error()))
	}

	return
}
