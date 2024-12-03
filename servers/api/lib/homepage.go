package lib

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

// AddHomepage registers a simple home page as the default
// root of the api
func AddHomepage(api huma.API) {
	//
	huma.Register(api, huma.Operation{
		OperationID:   "get-homepage",
		Method:        http.MethodGet,
		Path:          "/",
		Summary:       "Home",
		Description:   "Operates as the root of the API and a simple endpoint to test connectivity with, but returns no reporting data.",
		DefaultStatus: http.StatusOK,
	}, func(ctx context.Context, input *struct{}) (homepage *HomepageResponse, err error) {
		homepage = &HomepageResponse{}
		homepage.Body.Message = "Successful connection"
		return
	})
}
