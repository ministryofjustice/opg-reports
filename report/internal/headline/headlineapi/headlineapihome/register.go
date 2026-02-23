package headlineapihome

import (
	"context"
	"fmt"
	"net/http"
	"opg-reports/report/package/cntxt"
)

const ENDPOINT string = "/v1/headlines/{date_start}/{date_end}"

// Config contains required values for DB and others to generate a response
type Config struct {
	DB      string `json:"db"`
	Driver  string `json:"driver"`
	Params  string `json:"params"`
	Version string `json:"version"`
	SHA     string `json:"sha"`
}

// Register wraps the handle func with a local version that also gets additional config
// details
func Register(ctx context.Context, mux *http.ServeMux, config *Config) {
	var log = cntxt.GetLogger(ctx)

	log.Info("registering handler ... ", "endpoint", ENDPOINT)
	mux.HandleFunc(fmt.Sprintf("%s/{$}", ENDPOINT), func(writer http.ResponseWriter, request *http.Request) {
		Responder(ctx, config, request, writer)
	})

}
