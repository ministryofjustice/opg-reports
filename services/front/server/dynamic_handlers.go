package server

import (
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"opg-reports/services/front/tmpl"
	"opg-reports/shared/dates"
	"strings"
	"time"
)

type dynamicHandlerFunc func(w http.ResponseWriter, r *http.Request)

const apiVersion string = "v1"

func (s *FrontWebServer) Dynamic(w http.ResponseWriter, r *http.Request) {
	slog.Info("dynamic handler starting", slog.String("uri", r.RequestURI))

	data := s.Nav.Data(r)
	now := time.Now().UTC()
	end := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	start := end.AddDate(0, -9, 0)

	active := s.Nav.Active(r)
	if active == nil {
		return
	}

	// Call API!
	path := urlParse(active.Api, now, start, end)
	u := Url(s.ApiScheme, s.ApiAddr, path)
	slog.Info("calling api", slog.String("url", u.String()))

	resp, err := GetFromApi(u.String())

	if err != nil {
		return
	}

	headings := resp.GetResult().GetHeadings()
	// figure out how many headings in a row
	m := dates.Months(start, end)
	l := len(headings.GetCells()) - len(m)
	if l <= 0 {
		l = 1
	}

	data["Headings"] = headings.GetCells()
	data["HeadingCounter"] = l
	data["Result"] = resp.GetResult().GetRows()

	// Setup data object for the templates
	data["Organisation"] = s.Config.Organisation
	data["PageTitle"] = active.Name + " - "
	data["Now"] = now
	data["StartDate"] = start
	data["EndDate"] = end
	data["Months"] = dates.Months(start, end)
	data["Days"] = dates.Days(start, end)

	t, err := template.New(active.TemplateName).Funcs(tmpl.Funcs()).ParseFiles(s.templateFiles...)
	if err != nil {
		slog.Error("dynamic error", slog.String("err", fmt.Sprintf("%v", err)))
		return
	}

	s.Write(w, 200, t, active.TemplateName, data)

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
