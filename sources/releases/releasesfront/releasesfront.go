// Package releasesfront contains tools for the front end server to handle / transform
// api results.
package releasesfront

import (
	"log/slog"

	"github.com/ministryofjustice/opg-reports/pkg/transformers"
	"github.com/ministryofjustice/opg-reports/sources/releases/releasesio"
)

// TransformResult takes the result from the api and converts
// the data into table rows that can be used for the front
// end.
func TransformResult(body interface{}) (result interface{}) {
	var err error
	var res map[string]map[string]interface{}
	result = body

	switch body.(type) {
	case *releasesio.ReleasesBody:
		var bdy = body.(*releasesio.ReleasesBody)
		if res, err = transformers.ResultsToRows(bdy.Result, bdy.ColumnValues, bdy.DateRange); err == nil {
			bdy.TableRows = res
			result = bdy
		}

	}

	if err != nil {
		slog.Error("[releases.TransformResult] ", slog.String("err", err.Error()))
	}

	return
}
