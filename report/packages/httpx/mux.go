package httpx

import (
	"context"
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"opg-reports/report/internal/config"
	"opg-reports/report/packages/slogx"
)

// Mux is the main interface exposed and used
type Mux interface {
	http.Handler
	// Handle is the http.ServeMux handler exposed for static asset mapping
	Handle(pattern string, handler http.Handler)
	// HandleFunc is the http.ServeMux handler exposed for static asset mapping
	HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request))
	// Register is how api & front attach url pattern to handler (alternative HandleFunc)
	Register(pattern string, sources ...MuxResponseDataGetter)
	// RequestHandler
	RequestHandler(w http.ResponseWriter, r *http.Request, sources ...MuxResponseDataGetter)
}

// MuxConfig is used for exposing both connection and version details
// and a way to get all config details setup
type MuxConfig interface {
	Connection() *sql.DB
	Version() string
	Conf() *config.Config
}

// MuxResponseDataGetter is a function alias for how a package can process data from
// an api request
type MuxResponseDataGetter func(ctx context.Context, m Mux, r FitleredRequest, cfg MuxConfig, response *ResponseContent)

// writerF is alias for the `NewResponseWriter` function that we'll use to
// dynamically work out if this is html or not
type writerF func(w http.ResponseWriter, tmpl *template.Template) ResponseWriter

// Content is a general struct that is used to attach results for sending back
//
// The `Data` attribute should be added to by each MuxResponseDataGetter within
// their own scope, with json tags, to allow results to be returned as varying
// structs
type ResponseContent struct {
	Version string            `json:"version"`
	Request map[string]string `json:"request"`
	Data    map[string]any    `json:"data"`
}

// mux is used to wrap around the built-in ServeMux but attaches
// context and logger
type mux struct {
	*http.ServeMux
	ctx        context.Context    // context
	log        slogx.Logger       // logger
	newWriterF writerF            // function to create a new writer
	cfg        MuxConfig          // the version & db connection - used by api
	tmpl       *template.Template // the complied tempalte - used by the front end

}

// Register is an extended version of HandleFunc to allow passing along of data sources
//
// The sources are used within the handler function itself (`Handler`) to pull data
// from multiple sources and return the appropirate values
func (self *mux) Register(pattern string, sources ...MuxResponseDataGetter) {
	var handlerF = self.RequestHandler

	self.log.Info(self.ctx, fmt.Sprintf(`registering endpoint [%s] to a handler`, pattern))
	self.ServeMux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		handlerF(w, r, sources...)
	})
}

// Handler uses the data sources with the request data to generate a response.
//
// Allows for multiple data sources attaching to the same pattern / response
// so chaining is possible (formatting data to tables etc)
//
// DB connection is created once per call and closed within.
//
// A new response writer of the configured type (writerF) is created once per
// call to function
//
// Each source should contain any databaes calls / requirements needed
func (self *mux) RequestHandler(w http.ResponseWriter, r *http.Request, sources ...MuxResponseDataGetter) {
	var (
		filtered *Filter
		data     *ResponseContent
		content  []byte
		writer   = self.newWriterF(w, self.tmpl) // create a new writer using the attached func
		count    = len(sources)
	)

	self.log.Info(self.ctx, "handler triggered for request", "request", r.RequestURI)
	defer self.log.Info(self.ctx, "handler completed.", "request", r.RequestURI)
	// setup filter & request values from the original request
	filtered = RequestDataToFilter(ValuesFromRequest(r))

	// setup the standard request map
	data = &ResponseContent{
		Version: self.cfg.Version(),
		Request: filtered.RequestData().Map(),
		Data:    map[string]any{},
	}
	// loop over all the data getter sources and call each one in turn
	for i, srcF := range sources {
		self.log.Info(self.ctx, fmt.Sprintf("[%d/%d] getting data source content", i+1, count), "request", r.RequestURI)
		srcF(self.ctx, self, filtered, self.cfg, data)
	}
	// write the data
	content, _ = writer.BytesAndHeaders(data)
	writer.Write(content)
}

// NewMux create a new instance of mux
func NewMux(ctx context.Context, cfg MuxConfig, tmpl *template.Template) Mux {

	return &mux{
		ServeMux:   http.NewServeMux(),
		ctx:        ctx,
		log:        slogx.FromContext(ctx),
		cfg:        cfg,
		tmpl:       tmpl,
		newWriterF: NewResponseWriter,
	}
}
