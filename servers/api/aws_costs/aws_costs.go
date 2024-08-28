package aws_costs

import (
	"net/http"

	_ "github.com/mattn/go-sqlite3"
	"github.com/ministryofjustice/opg-reports/servers/shared/mw"
	"github.com/ministryofjustice/opg-reports/servers/shared/srvr/api"
)

// currently supported urls
const (
	ytdUrl      string = "/{version}/aws-costs/ytd/{$}"
	taxSplitUrl string = "/{version}/aws-costs/monthly-tax/{$}"
	standardUrl string = "/{version}/aws-costs/{$}"
)

// Register sets the local context and database paths to the values passed and then
// attaches the local handles to the url patterns supported by aws_costs api
func Register(mux *http.ServeMux, apiServer *api.ApiServer) {
	mux.HandleFunc(ytdUrl, mw.Middleware(api.Wrap(apiServer, YtdHandler), mw.Logging, mw.SecurityHeaders))
	mux.HandleFunc(taxSplitUrl, mw.Middleware(api.Wrap(apiServer, MonthlyTaxHandler), mw.Logging, mw.SecurityHeaders))
	mux.HandleFunc(standardUrl, mw.Middleware(api.Wrap(apiServer, StandardHandler), mw.Logging, mw.SecurityHeaders))
}
