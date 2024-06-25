package server

import (
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"opg-reports/services/front/cnf"
	"opg-reports/services/front/tmpl"
	"opg-reports/shared/aws/cost"
	"opg-reports/shared/dates"
	"opg-reports/shared/server"
	"strings"
	"time"
)

type dynamicHandlerFunc func(w http.ResponseWriter, r *http.Request)

const apiScheme string = "http"
const apiVersion string = "v1"

func (s *FrontWebServer) Dynamic(w http.ResponseWriter, r *http.Request) {
	data := s.Nav.Data(r)

	now := time.Now().UTC()
	end := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	start := end.AddDate(-1, -1, 0)

	slog.Info("dynamic handler starting", slog.String("uri", r.RequestURI))

	active := s.Nav.Active(r)
	if active == nil {
		return
	}

	// Call API!
	path := urlParse(active.Api, now, start, end)
	resultT, result, err := getFromApi(s.ApiScheme, s.ApiAddr, path)
	if err != nil {
		return
	}

	_, content := server.ResponseAsStrings(result)

	// Setup data object for the templates
	data["Organisation"] = s.Config.Organisation
	data["PageTitle"] = active.Name + " - "
	data["Now"] = now
	data["StartDate"] = start
	data["EndDate"] = end
	data["Months"] = dates.Months(start, end)
	data["Days"] = dates.Days(start, end)

	// map[string]R | map[string][]R | map[string]map[string][]R | []R
	switch active.ResponseHandler {
	case "Map":
	case "MapSlice":
	case "MapMapSlice":
		s.renderMapMapSlice(w, r, resultT, content, active, data)
	case "Slice":

	}

}

func (s *FrontWebServer) renderMapMapSlice(
	w http.ResponseWriter,
	r *http.Request,
	resultT server.ApiResponseConstraintString,
	content []byte,
	active *cnf.SiteSection,
	data map[string]interface{}) {

	slog.Debug("renderMapMapSlice starting")

	parsed := map[string]map[string]interface{}{}

	switch resultT.R() {
	case "cost.Cost":
		resp, _ := cost.ResponseAsStringsMapMapSlice(content)
		fmt.Printf("%+v\n", resp.Result)
		for key, values := range resp.Result {
			parsed[key] = map[string]interface{}{}
			for k, v := range values {
				parsed[key][k] = cost.Total(v)
			}
		}
	}

	data["Result"] = parsed
	t, err := template.New(active.TemplateName).Funcs(tmpl.Funcs()).ParseFiles(s.templateFiles...)
	if err != nil {
		slog.Error("renderMapMapSlice", slog.String("err", fmt.Sprintf("%v", err)))
		return
	}

	s.Write(w, 200, t, active.TemplateName, data)
}

func getFromApi(scheme string, addr string, path string) (resultType server.ApiResponseConstraintString, result *http.Response, err error) {

	u := Url(scheme, addr, path)
	slog.Info("calling api", slog.String("url", u.String()))

	resultType, result, err = GetFromApi(u.String())

	slog.Debug("api result",
		slog.String("resultType", string(resultType)),
		slog.String("err", fmt.Sprintf("%v", err)),
		slog.String("url", u.String()))
	return
}

// Parse out segements of the url that we typically replace with real values
func urlParse(url string, now time.Time, start time.Time, end time.Time) string {

	replacers := map[string]string{
		"apiVersion": apiVersion,
		"nowYM":      now.Format(dates.FormatYM),
		"nowYMD":     now.Format(dates.FormatYMD),
		"startYM":    start.Format(dates.FormatYM),
		"startYMD":   start.Format(dates.FormatYMD),
		"endYM":      end.Format(dates.FormatYM),
		"endYMD":     end.Format(dates.FormatYMD),
	}
	for key, value := range replacers {
		url = strings.ReplaceAll(url, fmt.Sprintf("{%s}", key), value)
	}

	return url
}
