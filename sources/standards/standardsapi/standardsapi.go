package standardsapi

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/jmoiron/sqlx"
	"github.com/ministryofjustice/opg-reports/pkg/datastore"
	"github.com/ministryofjustice/opg-reports/sources/standards"
	"github.com/ministryofjustice/opg-reports/sources/standards/standardsdb"
	"github.com/ministryofjustice/opg-reports/sources/standards/standardsio"
)

const Segment string = "standards"
const Tag string = "Standards"

var description string = `Returns a list of repository informations relating to the standards eahc has met and their current status.`

// apiHandler default handler for all github standards info
func apiHandler(ctx context.Context, input *standardsio.StandardsInput) (response *standardsio.StandardsOutput, err error) {

	var (
		result         []*standards.Standard
		db             *sqlx.DB
		dbFilepath     string = ctx.Value(Segment).(string)
		queryStatement        = standardsdb.FilterByIsArchived
		counters              = &standardsio.Counters{BaselineCompliant: 0, ExtendedCompliant: 0}
		bdy                   = &standardsio.StandardsBody{
			Request: input,
			Type:    "default",
		}
	)

	db, _, err = datastore.NewDB(ctx, datastore.Sqlite, dbFilepath)
	defer db.Close()
	if err != nil {
		return
	}

	if input.Unit != "" {
		queryStatement = standardsdb.FilterByIsArchivedAndTeam
	}

	// -- Do the overall count querys
	// Get total number of row
	counters.Total, err = datastore.Get[int](ctx, db, standardsdb.RowCount)
	if err != nil {
		return
	}
	// get total archived
	counters.TotalArchived, err = datastore.Get[int](ctx, db, standardsdb.ArchivedCount)
	if err != nil {
		return
	}
	// get total baseline count
	counters.TotalBaselineCompliant, err = datastore.Get[int](ctx, db, standardsdb.CompliantBaselineCount)
	if err != nil {
		return
	}
	// get total extended compliance count
	counters.TotalExtendedCompliant, err = datastore.Get[int](ctx, db, standardsdb.CompliantExtendedCount)
	if err != nil {
		return
	}
	// now run the main query
	if result, err = datastore.SelectMany[*standards.Standard](ctx, db, queryStatement, input); err == nil {
		bdy.Result = result
		// Add the last counter details
		counters.Count = len(result)
		for _, row := range result {
			if row.IsCompliantBaseline() {
				counters.BaselineCompliant += 1
			}
			if row.IsCompliantExtended() {
				counters.ExtendedCompliant += 1
			}
		}

	}

	bdy.Counters = counters
	response = &standardsio.StandardsOutput{
		Body: bdy,
	}
	return
}

// Register attaches all the endpoints this module handles on the passed huma api
//
// Currently supports the following endpoints:
//   - /{version}/standards/github/{archived}?unit=<team>
func Register(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID:   "get-standards-github",
		Method:        http.MethodGet,
		Path:          "/{version}/standards/github/{archived}",
		Summary:       "Github Standards",
		Description:   description,
		DefaultStatus: http.StatusOK,
		Tags:          []string{Tag},
	}, apiHandler)
}
