package httpx

import (
	"context"
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"opg-reports/report/packages/instance"
	"opg-reports/report/packages/slogx"
)

// MuxHook is called after all the MuxResponder slice have been; intended to provide
// a method of fetching same data form the same source for each page.
//
// Normally only used by the front / html server to fetch data for components that
// are on every page - such as the team based navigation.
type MuxHook func(ctx context.Context, cfg MuxConfigurer, r FitleredRequest, resp MuxResponseType)

// MuxConfigurer exposes configureation methods provided by `config.Config` struct
// that we want to be able to share within the both the html and json response
// types
type MuxConfigurer interface {
	// Connection returns a database connection that api can use
	Connection() *sql.DB
	// Version provides the version signature (semver only)
	Version() string
	// GovukVersion returns the govuk version string without the v prefix
	GovukVersion() string
	// Directories returns a map of directory paths that the
	// front end server would use
	Directories() map[string]string
	// Provides a template (possibly nil)
	Template(name string) (*template.Template, error)
	// ApiHostname
	ApiHostname() string
}

// MuxServer is a wrapper around `*http.ServeMux` so we can share some methods
type MuxServer interface {
	http.Handler
	//
	Handle(pattern string, handler http.Handler)
	// HandleFunc is the http.ServeMux handler exposed for static asset mapping
	HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request))
}

// MuxResponseType is the interface for the response data and methods that the
// mux would access to update / fetch data
type MuxResponseType interface {
	// TemplateName provides the template we want to use
	// For api / json response this is generally nil
	TemplateName() string
	// SetVersion pushes the version string into the response object
	SetVersion(v string)
	// SetGovukVersion sets the gov uk version
	SetGovukVersion(v string)
	// SetRequestData pushes the processed incoming request data
	// back on to the response
	SetRequestData(rd *RequestData)
	//
	SetTeams(teams []string)
}

// MuxResponder is typed version of a function that will deal with the incoming
// request, run queries etc and set values on the response object
type MuxResponder[T MuxResponseType] func(ctx context.Context, cfg MuxConfigurer, r FitleredRequest, response T)

// Register is a constrained  type function to register a handler to a url
// endpoint pattern
//
// The typing being at this level allows each http request to return a
// different type, so the mux can remain the same for each
//
//   - ctx is the active context, passed down from the starting command
//   - mux is the main mux that we'll call HandleFunc on
//   - cfg is configuration data used to fetch details that are set
//     within the environment / cli arguments
//   - hook is a post data source collection hook function, used to
//     provide a mechanism for consistent data; such as fetching team
//     names for the front end navigation thats used on every page
//   - dataSources are functions to call that will operate on the
//     repsonse data (T MuxResponseType) and update elements
//     within themselves
func Register[T MuxResponseType](
	ctx context.Context,
	mux MuxServer, cfg MuxConfigurer,
	pattern string, hook MuxHook, dataSources ...MuxResponder[T]) {

	var log slogx.Logger = slogx.FromContext(ctx)

	log.Info(ctx, fmt.Sprintf(`registering endpoint [%s] to a handler`, pattern))
	mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		RequestHandler(ctx, cfg, w, r, hook, dataSources...)
	})

}

// RequestHandler is called from Register to trigger the
// processing of an incoming request which will then call each
// for the dataSource functions to process and return the result.
//
// It converts the http request to a filter and sets the version
// and request data on the response, to help reduce repeating
// chunks of code.
//
//   - ctx is the active context, passed down from the starting command
//   - cfg is configuration data used to fetch details that are set
//     within the environment / cli arguments
//   - w & r are the standard http response writer and reader that
//     are being passed along from the `HandleFunc`
//   - hook is a post data source collection hook function, used to
//     provide a mechanism for consistent data; such as fetching team
//     names for the front end navigation thats used on every page
//   - dataSources are functions to call that will operate on the
//     repsonse data (T MuxResponseType) and update elements
//     within themselves - these can be chained to mutate results
//     with one function getting data from db and the next changing
//     it to be structured like a table etc.
func RequestHandler[T MuxResponseType](
	ctx context.Context,
	cfg MuxConfigurer,
	w http.ResponseWriter, r *http.Request,
	hook MuxHook,
	dataSources ...MuxResponder[T]) {

	var (
		err         error
		content     []byte
		writer      ResponseWriter
		tmpl        *template.Template
		data        T            = instance.Of[T]()
		count       int          = len(dataSources)
		log         slogx.Logger = slogx.FromContext(ctx)
		requestData *RequestData = ValuesFromRequest(r)
		filtered    *Filter      = RequestDataToFilter(requestData)
	)

	data.SetVersion(cfg.Version())
	data.SetGovukVersion(cfg.GovukVersion())
	data.SetRequestData(requestData)

	for i, sourceF := range dataSources {
		log.Info(ctx, fmt.Sprintf("[%d/%d] getting data source content", i+1, count), "request", r.RequestURI)
		sourceF(ctx, cfg, filtered, data)
	}

	// run hooks
	if hook != nil {
		hook(ctx, cfg, filtered, data)
	}
	// get the template
	tmpl, err = cfg.Template(data.TemplateName())
	if err != nil {
		log.Error(ctx, "error getting the template for the request", "err", err.Error())
	}
	// write the content out
	writer = NewResponseWriter(w, tmpl)
	content, _ = writer.BytesAndHeaders(data)
	writer.Write(content)
}

func NewMux() MuxServer {
	var s = http.NewServeMux()
	return s
}
