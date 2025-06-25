package home

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/ministryofjustice/opg-reports/report/config"
)

type HomepageResponse struct {
	Body struct {
		Message string `json:"message" example:"Successful connection."`
	}
}

func RegisterHomepage(log *slog.Logger, conf *config.Config, api huma.API) {
	var operation = huma.Operation{
		OperationID:   "get-home",
		Method:        http.MethodGet,
		Path:          "/",
		Summary:       "Home",
		Description:   "Operates as the root of the API and a simple endpoint to test connectivity with, but returns no data.",
		DefaultStatus: http.StatusOK,
	}

	huma.Register(api, operation, func(ctx context.Context, input *struct{}) (homepage *HomepageResponse, err error) {
		homepage = &HomepageResponse{}
		homepage.Body.Message = "Successful connection"
		return
	})
}
