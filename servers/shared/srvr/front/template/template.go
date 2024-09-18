package template

import (
	"log/slog"
	"net/http"
	"strings"
	"text/template"

	"github.com/ministryofjustice/opg-reports/servers/api/aws_costs"
	"github.com/ministryofjustice/opg-reports/shared/convert"
	"github.com/ministryofjustice/opg-reports/shared/dates"
)

// Declare all of the template exposed functions

var (
	currency      = func(s interface{}, symbol string) string { return convert.Curr(s, symbol) }             // convert interface into string formatted like a currency
	add           = func(a float64, b float64) float64 { return a + b }                                      // add 2 floats together
	addInt        = func(a int, b int) int { return a + b }                                                  // add two integers together
	uptimePercent = func(s interface{}) string { return convert.FloatString(s, "%.4f") + " %" }              // return uptime as a % with 4 decimals
	percent       = func(got int, total int) string { return convert.Percent(got, total) }                   // work out percentage
	title         = func(s string) string { return convert.Title(s) }                                        // create a string suitable for use in page title
	titles        = func(strs ...string) (str string) { return convert.Titles(strs...) }                     // create a single tile string from a series of strings
	lower         = func(s string) string { return strings.ToLower(s) }                                      // lowercase the string passed
	stripI        = func(s string) string { return convert.StripIntPrefix(s) }                               // remove integer prefix (1.Name => Name)
	day           = func(t string) string { return dates.Time(t).Format(dates.FormatYMD) }                   // Format the string into a time in YYYY-MM-DD
	dayBefore     = func(t string) string { return dates.Time(t).AddDate(0, 0, -1).Format(dates.FormatYMD) } // Format the string into a time in YYYY-MM-DD and go back one day
	costIdx       = func(set []*aws_costs.CommonResult, i int) any { return set[i] }                         // Return the entry in the slice at index i
	dict          = func(values ...any) map[string]any { return convert.Dict(values...) }                    // generate a map[string] from series of passed params
	col           = func(i string, mapped map[string]string) string { return mapped[i] }
)

// functionMap contains all the function to expose within the templates
var functionMap map[string]interface{} = map[string]interface{}{
	"currency":      currency,
	"add":           add,
	"addInt":        addInt,
	"percent":       percent,
	"uptimePercent": uptimePercent,
	"title":         title,
	"titles":        titles,
	"lower":         lower,
	"stripI":        stripI,
	"day":           day,
	"dayBefore":     dayBefore,
	"costIdx":       costIdx,
	"dict":          dict,
	"col":           col,
}

// ---------

// Template handles the output of the server front template
type Template struct {
	TemplateName string
	Template     *template.Template
	AllTemplates []string
	Writer       http.ResponseWriter
}

// response sets the response status and header and executes the template
// as well. Cpatures error and logs to slog
func (f *Template) response(data any, status int) {
	f.Writer.WriteHeader(status)
	f.Writer.Header().Set("Content-Type", "text/html")
	if err := f.Template.ExecuteTemplate(f.Writer, f.TemplateName, data); err != nil {
		slog.Error("[front] template execute failed",
			slog.String("templateName", f.TemplateName),
			slog.String("error", err.Error()))
	}
}

// Run generates a new template (as .TemplateName) adds in the functions from functionMap
// and includes AllTemplates to be parsed.
//
// The template will then be executed and the result written in to the http.ResponseWriter
// with a status header and content type values set
func (f *Template) Run(data any) {
	status := http.StatusOK
	tmpl, err := template.New(f.TemplateName).Funcs(functionMap).ParseFiles(f.AllTemplates...)

	if err != nil {
		slog.Error("[front] template run error", slog.String("err", err.Error()))
		status = http.StatusBadGateway
	}
	f.Template = tmpl
	f.response(data, status)
}

// NewTemplate returns a prepared template for handling a front end call
func New(templateName string, templates []string, writer http.ResponseWriter) *Template {
	return &Template{
		TemplateName: templateName,
		AllTemplates: templates,
		Writer:       writer,
	}
}
