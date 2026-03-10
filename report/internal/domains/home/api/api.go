// Package `api` handles the home page & ping responses.
//
// Minimal data is returned, just version & request.
package api

import (
	"opg-reports/report/internal/domains/home/types"
	"opg-reports/report/packages/args"
	"opg-reports/report/packages/handler"
	"opg-reports/report/packages/types/models"
)

const label string = `get-home`

func Config(opts *args.API) *handler.ApiConfig {
	return &handler.ApiConfig{
		Name:     label,
		Database: opts.DB,
		Selector: nil,
	}
}

func Response(opts *args.API) *models.ApiResponse {
	return &models.ApiResponse{
		Versions: opts.Versions,
	}
}

func T() *types.NilRow {
	return &types.NilRow{}
}
