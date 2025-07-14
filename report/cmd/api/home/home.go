package home

import (
	"context"
	"log/slog"
	"net/http"

	"opg-reports/report/config"
	"opg-reports/report/internal/endpoints"

	"github.com/danielgtaylor/huma/v2"
)

type HomepageResponse struct {
	Body struct {
		Message string `json:"message" example:"Successful connection."`
	}
}

func RegisterGetHomepage(log *slog.Logger, conf *config.Config, api huma.API) {
	var operation = huma.Operation{
		OperationID:   "get-home",
		Method:        http.MethodGet,
		Path:          endpoints.HOME,
		Summary:       "Home",
		Description:   "Operates as the root of the API and a simple endpoint to test connectivity with, but returns no data.",
		DefaultStatus: http.StatusOK,
		Tags:          []string{"Home"},
	}

	huma.Register(api, operation, func(ctx context.Context, input *struct{}) (homepage *HomepageResponse, err error) {
		log.Info("handling get-home")
		homepage = &HomepageResponse{}
		homepage.Body.Message = "Successful connection"
		return
	})
}
