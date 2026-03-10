package handler

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"opg-reports/report/packages/dbx"
	"opg-reports/report/packages/logger"
	"opg-reports/report/packages/respond"
	"opg-reports/report/packages/types/interfaces"
	"opg-reports/report/packages/types/models"
)

// API processes the request and provides a json result using
// interface based data.
//
// Before calling this function, the request & repsonse should be
// setup within the muc registration process to include the
// server values.
//
// Flow:
//   - Reset response
//   - Set type / label info
//   - Populate the request struct from the http request
//   - Attach the populated request to the response
//   - Create the filter data from the request
//   - Attach the filter to the response
//   - If the congifured statement is set ...
//     -- Run the select call
//     -- Add each record to the response
//   - TODO: data transformation - tabular conversion ??
//     -- HEADINGS as well
//   - Return response as json
//
// ADD MORE LOGGING
func API[T interfaces.Row](
	ctx context.Context,
	cfg interfaces.ApiConfiguration,
	request interfaces.ApiRequest,
	response interfaces.ApiResponse,
) (err error) {
	var (
		log     *slog.Logger
		sql     string
		filter  interfaces.Filterable
		records []T
	)
	ctx, log = logger.Get(ctx)
	log.Info("api: [" + cfg.Label() + "] running request ...")
	response.Reset()
	// set the label
	response.Typed(cfg.Label())
	// set the request details
	request.Populate(request.Request())
	response.Request(request)
	// setup the filter
	filter = request.Filter(request.Request())
	response.Filter(filter)
	// prep for running the sql statement
	if cfg.Statement() != nil {
		sql = cfg.Statement().SQL(filter)
		// run the select
		records, err = dbx.Select[T](ctx, sql, filter, cfg.DB())
		if err != nil {
			return
		}
		// set the database rows on the reponse
		for _, record := range records {
			response.Record(record)
		}
	}
	// attach record as a result
	for _, row := range response.Records() {
		response.Result(row.Result())
	}

	respond.AsJSON(ctx, request.Request(), response.Response(), response)

	return
}

// RegisterAPI attaches a handler setup (via API) to the endpoints passed.
//
// `res` &  `row` are unused and are there to resolve type resolution
// without having to have knowledge of the types within the main api
// command.
func RegisterAPI[T interfaces.Row](
	ctx context.Context,
	mux *http.ServeMux,
	cfg interfaces.ApiConfiguration,
	response interfaces.ApiResponse,
	row T,
	endpoints ...string,
) {
	var log *slog.Logger
	ctx, log = logger.Get(ctx)

	for _, ep := range endpoints {
		log.Info(fmt.Sprintf("api: [%s] registering endpoint [%s]", cfg.Label(), ep))

		// bind the handler to the endpoints and setup wrappers that we pass along
		mux.HandleFunc(ep, func(writer http.ResponseWriter, request *http.Request) {
			var r = &models.Request{}
			r.HttpRequest(request)
			response.Writer(writer)

			API[T](ctx, cfg, r, response)
		})
	}
}
