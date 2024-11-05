package uptimeapi

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/jmoiron/sqlx"
	"github.com/ministryofjustice/opg-reports/pkg/convert"
	"github.com/ministryofjustice/opg-reports/pkg/datastore"
	"github.com/ministryofjustice/opg-reports/sources/uptime"
	"github.com/ministryofjustice/opg-reports/sources/uptime/uptimedb"
	"github.com/ministryofjustice/opg-reports/sources/uptime/uptimeio"
)

const Segment string = "uptime"
const Tag string = "Uptime"

var overallDescription = `Returns lsit of uptime data grouped by interval for all services.`

// apiOverallHandler gets the overall uptime for all services within the start and end date
// range
func apiOverallHandler(ctx context.Context, input *uptimeio.UptimeInput) (response *uptimeio.UptimeOutput, err error) {
	var (
		result         []*uptime.Uptime
		db             *sqlx.DB
		dbFilepath     string                         = ctx.Value(Segment).(string)
		queryStatement datastore.NamedSelectStatement = uptimedb.UptimeByInterval
		bdy            *uptimeio.UptimeBody           = &uptimeio.UptimeBody{
			Request:     input,
			Type:        "overall",
			ColumnOrder: []string{"unit"},
			DateRange:   convert.DateRange(input.StartTime(), input.EndTime(), input.Interval),
		}
	)

	db, _, err = datastore.NewDB(ctx, datastore.Sqlite, dbFilepath)
	defer db.Close()
	if err != nil {
		return
	}

	if result, err = datastore.Select[[]*uptime.Uptime](ctx, db, queryStatement, input); err == nil {
		bdy.Result = result
		bdy.ColumnValues = datastore.ColumnValues(result, bdy.ColumnOrder)
	}

	response = &uptimeio.UptimeOutput{
		Body: bdy,
	}

	return
}

var unitDescription string = `Returns a list of all uptime data grouped by the interval and unit - with option of filtering by a unit.`

// apiUnitHandler gets uptime, groups that by the date interval and the unit name for all data between
// the start & end dates.
// Optional filter by unit name
func apiUnitHandler(ctx context.Context, input *uptimeio.UptimeInput) (response *uptimeio.UptimeOutput, err error) {

	var (
		result         []*uptime.Uptime
		db             *sqlx.DB
		dbFilepath     string                         = ctx.Value(Segment).(string)
		queryStatement datastore.NamedSelectStatement = uptimedb.UptimeByIntervalUnitAll
		bdy            *uptimeio.UptimeBody           = &uptimeio.UptimeBody{
			Request:     input,
			Type:        "unit",
			ColumnOrder: []string{"unit"},
			DateRange:   convert.DateRange(input.StartTime(), input.EndTime(), input.Interval),
		}
	)

	db, _, err = datastore.NewDB(ctx, datastore.Sqlite, dbFilepath)
	defer db.Close()
	if err != nil {
		return
	}

	// if there is a unit passed along, then use the query that will filter by that
	if input.Unit != "" {
		queryStatement = uptimedb.UptimeByIntervalUnitFiltered
	}

	if result, err = datastore.Select[[]*uptime.Uptime](ctx, db, queryStatement, input); err == nil {
		bdy.Result = result
		bdy.ColumnValues = datastore.ColumnValues(result, bdy.ColumnOrder)
	}

	response = &uptimeio.UptimeOutput{
		Body: bdy,
	}

	return
}

// Register attaches all the endpoints this module handles on the passed huma api
//
// Currently supports the following endpoints:
//   - /{version}/uptime/aws/overall/{start_date}/{end_date}/{interval}
//   - /{version}/uptime/aws/unit/{start_date}/{end_date}/{interval}?unit=<unit>
func Register(api huma.API) {

	huma.Register(api, huma.Operation{
		OperationID:   "get-uptime-aws",
		Method:        http.MethodGet,
		Path:          "/{version}/uptime/aws/overall/{start_date}/{end_date}/{interval}",
		Summary:       "Overal uptime per interval",
		Description:   overallDescription,
		DefaultStatus: http.StatusOK,
		Tags:          []string{Tag},
	}, apiOverallHandler)

	huma.Register(api, huma.Operation{
		OperationID:   "get-uptime-aws-by-unit",
		Method:        http.MethodGet,
		Path:          "/{version}/uptime/aws/unit/{start_date}/{end_date}/{interval}",
		Summary:       "Uptime for each unit",
		Description:   unitDescription,
		DefaultStatus: http.StatusOK,
		Tags:          []string{Tag},
	}, apiUnitHandler)
}
