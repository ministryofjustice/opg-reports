package aws_uptime

import (
	"net/http"

	_ "github.com/mattn/go-sqlite3"
	"github.com/ministryofjustice/opg-reports/servers/shared/srvr/api"
	"github.com/ministryofjustice/opg-reports/servers/shared/srvr/mw"
)

const url string = "/{version}/aws-uptime/{$}"

// Register sets the local context and database paths to the values passed and then
// attaches the local handles to the url patterns supported by aws_costs api
func Register(mux *http.ServeMux, apiServer *api.ApiServer) {
	mux.HandleFunc(url, mw.Middleware(api.Wrap(apiServer, Handler), mw.Logging, mw.SecurityHeaders))
}
