package apihandlers

import (
	"log/slog"
	"net/http"

	"github.com/ministryofjustice/opg-reports/servers/shared/srvr/httphandler"
)

var HANDLERS = map[string]ApiHandler{
	"aws_uptime":       &AWSUptime{},
	"aws_costs":        &AWSCosts{},
	"github_standards": &GithubStandards{},
}

type ApiHandler interface {
	Handles(uri string) bool
	Handle(data map[string]interface{}, key string, remote *httphandler.HttpHandler, w http.ResponseWriter, r *http.Request)
}

// Get returns the struct to call to process the api call requried
func Get(uri string) (apiHandler ApiHandler) {
	for key, handler := range HANDLERS {
		if handler.Handles(uri) {
			slog.Info("found handler", slog.String("handler", key), slog.String("uri", uri))
			apiHandler = handler
		}
	}
	return
}
