package aws_costs

import (
	"context"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
	"github.com/ministryofjustice/opg-reports/servers/shared/mw"
)

// currently supported urls
const (
	ytdUrl      string = "/{version}/aws-costs/ytd/{$}"
	taxSplitUrl string = "/{version}/aws-costs/monthly-tax/{$}"
	standardUrl string = "/{version}/aws-costs/{$}"
)

// db and context
var (
	apiCtx    context.Context
	apiDbPath string
)

// Register sets the local context and database paths to the values passed and then
// attaches the local handles to the url patterns supported by aws_costs api
func Register(ctx context.Context, mux *http.ServeMux, dbPath string) (err error) {
	SetCtx(ctx)
	SetDBPath(dbPath)
	// -- registers
	mux.HandleFunc(taxSplitUrl, mw.Middleware(MonthlyTaxHandler, mw.Logging, mw.SecurityHeaders))
	mux.HandleFunc(ytdUrl, mw.Middleware(YtdHandler, mw.Logging, mw.SecurityHeaders))
	mux.HandleFunc(standardUrl, mw.Middleware(StandardHandler, mw.Logging, mw.SecurityHeaders))
	return nil
}

func SetDBPath(path string) {
	apiDbPath = path
}
func SetCtx(ctx context.Context) {
	apiCtx = ctx
}
