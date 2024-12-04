package lib

import (
	"bufio"
	"bytes"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"strings"
	"sync"

	"github.com/ministryofjustice/opg-reports/info"
	"github.com/ministryofjustice/opg-reports/internal/fetch"
	"github.com/ministryofjustice/opg-reports/internal/navigation"
	"github.com/ministryofjustice/opg-reports/internal/render"
	"github.com/ministryofjustice/opg-reports/internal/structs"
)

// Cfg contains the config data (address, mux server)
// for this server
type Cfg struct {
	Addr   string
	Mux    *http.ServeMux
	server *http.Server
}

// Server creates (or returns existing) http.Server
// using its own details
func (self *Cfg) Server() *http.Server {
	if self.server == nil {
		self.server = &http.Server{
			Addr:    self.Addr,
			Handler: self.Mux,
		}
	}
	return self.server
}

var errorTemplate = "error"

// Response contains info used for rendering the html
// page and its response details
type Response struct {
	Organisation string
	GovUKVersion string
	Templates    []string
	Funcs        template.FuncMap
	Errors       []error
	renderer     *render.Render
	headerSet    bool
	errCode      int
}

// Renderer gets the child template render helper
// This will deal with executing the correct template for the page
// using the data and functions
func (self *Response) Renderer() *render.Render {
	if self.renderer == nil {
		self.renderer = render.New(self.Templates, self.Funcs)
	}
	return self.renderer
}

// Write uses the template and data passed to execute the template stack
// (using renderer) and output the result using the writer
// This is the final step of processing the incoming request
// Sets http status and content type as text/html
// If there is an error executing the template then this will attempt to render
// the error template instead - adding a new error to the stack
// Uses an internal bool to avoid writing header status code more than once
func (self *Response) Write(templateName string, data map[string]any, writer http.ResponseWriter) {
	var rnd = self.Renderer()
	var buf = new(bytes.Buffer)
	var wr = bufio.NewWriter(buf)
	// inject errors to the page data
	if len(self.Errors) > 0 {
		data["Errors"] = self.ErrorList()
	}

	// use the renderer to execute the content and write the result to the http response
	if err := rnd.Write(templateName, data, wr); err != nil {
		slog.Error("[svr.Response] Write error", slog.String("err", err.Error()))
		if templateName != errorTemplate {
			slog.Error("[svr.Response] recovering from error, rendering error page")
			self.AddError(err)
			// recall self
			self.Write(errorTemplate, data, writer)
		}
	}
	// set the status code header
	if len(self.Errors) > 0 && !self.headerSet {
		writer.WriteHeader(self.errCode)
	} else if !self.headerSet {
		writer.WriteHeader(http.StatusOK)
	}
	// flush and write to real header
	writer.Header().Set("Content-Type", "text/html")
	wr.Flush()
	writer.Write(buf.Bytes())

	self.headerSet = true
	return
}

// WriteWithError ensures there is an error set and that the standard error
// code is returned (400)
// Will then render an error template to the page so the user sees something
// Calls Write after setting template and content
func (self *Response) WriteWithError(err error, writer http.ResponseWriter) {
	var data = map[string]interface{}{}
	self.AddError(err)
	self.Write(errorTemplate, data, writer)

}

// AddError method to add errors to the stack
func (self *Response) AddError(err error) {
	self.Errors = append(self.Errors, err)
}
func (self *Response) Reset() {
	self.Errors = []error{}
	self.headerSet = false
	self.errCode = http.StatusBadRequest
}

// ErrorList is helper to get string version of all errors
func (self *Response) ErrorList() (errs []string) {
	errs = []string{}
	for _, e := range self.Errors {
		errs = append(errs, e.Error())
	}
	return
}

// Nav holds the navigation data which will be used in the
// front end
type Nav struct {
	Tree []*navigation.Navigation
	flat map[string]*navigation.Navigation
}

// Flat returns a flat version of the navigation
func (self *Nav) Flat() map[string]*navigation.Navigation {
	if len(self.flat) <= 0 {
		self.flat = map[string]*navigation.Navigation{}
		navigation.Flat(self.Tree, self.flat)
	}
	return self.flat
}

// Api contains details on where the api is and
// how to connect to it
type Api struct {
	Version string
	Addr    string
}

// Svr is the main struct that renders the front end svr stack
type Svr struct {
	Cfg        *Cfg
	Response   *Response
	Navigation *Nav
	Api        *Api
}

// Handler is used to process each non-static / non-redirect request that comes
// to the front server.
//
// It uses the request url to find the active item in the navigation tree.
// It then uses that navigation item to fetch page data from the api. If the
// data sources are empty, no data is fetched, but pageData will contain
// default info for being able to render the outer template of the website.
//
// Each call the api (handled inside FetchDataForPage) is run within a go func
// to allow faster concurrent calls.
//
// The returned data is merged into the pageData map under the namespace
// tracked in the navigation item.
//
// That data has been casted into the .Body attribute, so will be of the
// expected struct within the template (generally a `StandardBody`).
func (self *Svr) Handler(writer http.ResponseWriter, request *http.Request) {
	slog.Info("[svr.Handler] uri: " + request.URL.String())
	var (
		activePage *navigation.Navigation
		flat       = self.Navigation.Flat()
		pageData   = map[string]interface{}{
			//
			"Signature": info.BuildInfo(),
			// org is used in the header
			"Organisation": self.Response.Organisation,
			"Path":         request.URL.Path,
			// used for path to the css etc
			"GovUKVersion": self.Response.GovUKVersion,
			// empty placeholders
			"NavigationActive":  nil,
			"NavigationRoot":    nil,
			"NavigationSidebar": []*navigation.Navigation{},
			"NavigationTopbar":  []*navigation.Navigation{},
		}
	)
	self.Response.Reset()
	// activate items in the stack
	activePage = navigation.ActivateFlat(flat, request)
	if activePage == nil {
		e := fmt.Errorf("requested url [%s] does not match any navigation item.", request.URL.String())
		// Return with error
		self.Response.WriteWithError(e, writer)
		return
	}
	slog.Info("[svr.Handler] activePage: " + activePage.Name)
	pageData["NavigationActive"] = activePage
	// top nav bar
	top := navigation.Level(self.Navigation.Tree)
	pageData["NavigationTopbar"] = top

	// get the navigation sidebar
	root := navigation.Root(activePage)
	if root != nil && len(root.Children()) > 0 {
		pageData["NavigationRoot"] = root
		pageData["NavigationSidebar"] = root.Children()
	}

	// get the data for each endpoint - use go routines
	if len(activePage.Data) > 0 {
		FetchDataForPage(self.Api, activePage, pageData, request)
	}
	self.Response.Write(activePage.Display.PageTemplate, pageData, writer)
	return
}

// Register iterates over the configured navigation structure
// and attaches each to the `Handler` method and appends the
// `/{$}` pattern for the go router to allow trailing slashes.
func (self *Svr) Register() {
	var suffix = "{$}"
	for _, nav := range self.Navigation.Flat() {
		var uri = nav.Uri
		uri = strings.TrimSuffix(uri, "/") + "/" + suffix
		slog.Info("[svr] registering", slog.String("uri", uri))

		self.Cfg.Mux.HandleFunc(uri, self.Handler)
		nav.Display.Registered = true
	}
}

// Run is called to start up the server
// This triggers the ListenAnServe method as well
// as setting up static pages and the homepage
// redirect
func (self *Svr) Run() {
	// setup the static redirects
	Statics(self.Cfg.Mux)
	// setup homepage redirect
	HomepageRedirect(self.Cfg.Mux, self.Navigation.Flat(), self.Navigation.Tree[0])
	// Register all the urls
	self.Register()

	slog.Info("Starting front server...")
	slog.Info(fmt.Sprintf("FRONT: [http://%s/]", self.Cfg.Addr))

	self.Cfg.Server().ListenAndServe()

}

// FetchDataForPage iterates over the data items for this navigation item
// and tries to fetch them from the api (by calling Fetcher).
// The result is then merged into the pageData map for processing in the
// front end templates.
//
// If the navigation data source has a transform function (`.Transformer`)
// that will be called to process the api result, otherwise the api
// data is directly inserted to pageData. This allows complex end points
// to have raw data converted before the front end - a semi middleware
// approach. Generally used to convert raw data into tabular rows for
// easier front end display.
//
// It uses go func call with a mutex and waitgroup to handle calling
// the api multiple times concurrently - this will mean pages with
// multiple blocks should perform better.
func FetchDataForPage(api *Api, activePage *navigation.Navigation, pageData map[string]interface{}, request *http.Request) {
	var (
		mutex     *sync.Mutex    = &sync.Mutex{}
		waitgroup sync.WaitGroup = sync.WaitGroup{}
	)
	for _, navData := range activePage.Data {
		waitgroup.Add(1)
		// use a go func to get the api data and process
		go func(a *Api, nd *navigation.Data) {
			mutex.Lock()
			defer mutex.Unlock()
			Fetcher(a, nd, request, nd.Body)
			// deal with the data
			// - if we have a transformer, set it
			// - otherwise pass it raw
			if nd.Body != nil && nd.Transformer != nil {
				slog.Info("Setting result to namespace via transformer", slog.String("namespace", nd.Namespace))
				pageData[nd.Namespace] = nd.Transformer(nd.Body)
			} else if nd.Body != nil {
				slog.Info("Setting result to namespace.", slog.String("namespace", nd.Namespace))
				pageData[nd.Namespace] = nd.Body
			}
			waitgroup.Done()
		}(api, navData)
	}
	// wait for the
	waitgroup.Wait()

}

// Fetcher uses the api & navigation data to go and fetch the data from the api
// and set the result item.
// Uses `fetch.Fetch` to call the api, and then uses Unmarshalling to convert
// the string data from the response into a struct.
func Fetcher[T any](api *Api, nav *navigation.Data, request *http.Request, result T) {
	var (
		host   = fmt.Sprintf("http://%s", api.Addr)
		source = nav.Source
	)
	if nav.Body == nil {
		slog.Error("no return type set on this nav item", slog.String("source", source.String()))
		return
	}
	content, code, err := fetch.Fetch(host, source, request)
	// error checks
	if err != nil || code != http.StatusOK {
		slog.Error("error calling api",
			slog.Int("statusCode", code),
			slog.String("err", fmt.Sprintf("%v", err)),
			slog.String("source", source.Parse(request)),
			slog.String("host", host))
		return
	}

	if err := structs.Unmarshal(content, result); err != nil {
		slog.Error("error converting into body", slog.String("err", err.Error()))
		return
	}
	return
}

// NewSvr creates a svr setup that can then be executed
// with `.Run()`
func NewSvr(cfg *Cfg, resp *Response, nv *Nav, api *Api) (s *Svr) {
	s = &Svr{
		Cfg:        cfg,
		Response:   resp,
		Navigation: nv,
		Api:        api,
	}

	return
}
