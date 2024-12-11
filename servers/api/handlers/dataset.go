package handlers

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/ministryofjustice/opg-reports/internal/dbs"
	"github.com/ministryofjustice/opg-reports/internal/dbs/adaptors"
	"github.com/ministryofjustice/opg-reports/internal/dbs/crud"
	"github.com/ministryofjustice/opg-reports/models"
	"github.com/ministryofjustice/opg-reports/servers/inout"
)

var (
	DatasetsSegment string   = "dataset"
	DatasetTags     []string = []string{"Datasets"}
)

const DatasetsListDescription string = `Returns the dataset.`
const DatasetListOperationID string = "get-dataset-list"
const datasetListSQL string = `
SELECT
	dataset.*
FROM dataset
LIMIT 1;
`

// ApiDatasetsListHandler
//
// Endpoints:
//
//	/version/datasets/list
func ApiDatasetsListHandler(ctx context.Context, input *inout.VersionInput) (response *inout.DatasetsListResponse, err error) {
	var (
		adaptor dbs.Adaptor
		results []*models.Dataset       = []*models.Dataset{}
		dbPath  string                  = ctx.Value(dbPathKey).(string)
		body    *inout.DatasetsListBody = inout.NewDatasetsListBody()
	)
	body.Request = input
	body.Operation = DatasetListOperationID
	// setup response
	response = &inout.DatasetsListResponse{}
	// hook up adaptor
	adaptor, err = adaptors.NewSqlite(dbPath, false)
	if err != nil {
		slog.Error("[api] dataset list adaptor error", slog.String("err", err.Error()))
	}
	defer adaptor.DB().Close()
	// get the data and attach results / errors to the response
	results, err = crud.Select[*models.Dataset](ctx, adaptor, datasetListSQL, nil)
	if err != nil {
		slog.Error("[api] dataset list select error", slog.String("err", err.Error()))
		body.Errors = append(body.Errors, fmt.Errorf("dataset list selection failed."))
	} else {
		body.Result = results
	}
	response.Body = body
	return
}

// Register attaches the handler to the main api
func RegisterDatasets(api huma.API) {
	var uri string = "/{version}/" + DatasetsSegment + "/list"

	slog.Info("[api] handler register ", slog.String("uri", uri))
	huma.Register(api, huma.Operation{
		OperationID:   DatasetListOperationID,
		Method:        http.MethodGet,
		Path:          uri,
		Summary:       "List dataset yupe",
		Description:   DatasetsListDescription,
		DefaultStatus: http.StatusOK,
		Tags:          DatasetTags,
	}, ApiDatasetsListHandler)

}
