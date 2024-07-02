package server

import (
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"net/url"
	"opg-reports/services/front/tmpl"
	"opg-reports/shared/dates"
	"opg-reports/shared/server/response"
	"strings"
	"time"
)

type dynamicHandlerFunc func(w http.ResponseWriter, r *http.Request)

const apiVersion string = "v1"
const billingDay int = 13

// Dynamic handles all end points that require data from the api.
func (s *FrontWebServer) Dynamic(w http.ResponseWriter, r *http.Request) {
	slog.Info("dynamic handler starting", slog.String("uri", r.RequestURI))

	data := s.Nav.Data(r)
	now := time.Now().UTC()
	end := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	if now.Day() < billingDay {
		end = end.AddDate(0, -2, 0)
	} else {
		end = end.AddDate(0, -1, 0)
	}
	start := end.AddDate(0, -11, 0)
	months := dates.Months(start, end)
	days := dates.Days(start, end)

	active := s.Nav.Active(r)
	// TODO - handle this issue with a redirect to root?
	if active == nil {
		slog.Error("active nil")
		return
	}

	// Setup data object for the templates
	data["Organisation"] = s.Config.Organisation
	data["PageTitle"] = active.Name + " - "
	data["Now"] = now
	data["StartDate"] = start
	data["EndDate"] = end
	data["Months"] = months
	data["Days"] = days
	data["Standards"] = s.Config.Standards
	data["Result"] = nil

	usePrefix := len(active.Api) > 1
	// Handle multiple api calls for one page
	for apiResultName, apiUrl := range active.Api {
		// Call API!
		path := urlParse(apiUrl, now, start, end)
		u := Url(s.ApiScheme, s.ApiAddr, path)
		apiResp, err := s.handleApiCall(u)
		// no error from the api, and no error from parsing the api resposne
		// into local data, then set to the data map ready for the template parsing
		if err == nil {
			apiData, err := s.parseResponse(apiResp)
			if err == nil {
				for key, val := range apiData {
					if usePrefix {
						key = fmt.Sprintf("%s_%s", apiResultName, key)
					}
					slog.Info("setting api res data", slog.String("key", key))
					data[key] = val

				}
			} else {
				slog.Error("dynamic handler error from parsing repsonse",
					slog.String("url", u.String()),
					slog.String("err", fmt.Sprintf("%v", err)),
				)
			}
		} else {
			slog.Error("dynamic handler error from api",
				slog.String("url", u.String()),
				slog.String("err", fmt.Sprintf("%v", err)),
			)
		}
	}

	t, err := template.New(active.TemplateName).Funcs(tmpl.Funcs()).ParseFiles(s.templateFiles...)
	if err != nil {
		slog.Error("dynamic error", slog.String("err", fmt.Sprintf("%v", err)))
		return
	}

	s.Write(w, 200, t, active.TemplateName, data)

}

func (s *FrontWebServer) parseResponse(apiResp *http.Response) (data map[string]interface{}, err error) {
	data = map[string]interface{}{}

	// _, b := response.Stringify(apiResp)
	resp := response.NewResponse[*response.Cell, *response.Row[*response.Cell]]()
	err = response.FromHttp(apiResp, resp)
	if err != nil {
		slog.Error("parse error")
		return
	}
	// get min / max times
	if min := resp.GetDataAgeMin(); min.Format(dates.FormatY) != dates.ErrYear {
		data["DataAgeMin"] = min
	}
	if max := resp.GetDataAgeMax(); max.Format(dates.FormatY) != dates.ErrYear {
		data["DataAgeMax"] = max
	}
	// If the result is nil (failed parsing), return
	result := resp.GetData()
	if result == nil {
		slog.Error("empty result")
		return
	}

	if heading := result.GetTableHead(); heading.GetHeadersCount() > 0 {
		data["Headings"] = heading.GetAll()
		data["HeadingsPre"] = heading.GetHeadersCount()
		data["HeadingsPost"] = heading.GetTotalCellCount() - heading.GetHeadersCount() - heading.GetSupplementaryCount()
	}

	if footer := result.GetTableFoot(); footer.GetHeadersCount() > 0 {
		data["Footer"] = footer.GetAll()
		data["FooterPre"] = footer.GetHeadersCount()
		data["FooterPost"] = footer.GetTotalCellCount() - footer.GetHeadersCount() - footer.GetSupplementaryCount()
	}

	// fetch the resulting rows, return if they are empty
	rows := result.GetTableBody()
	if rows != nil {
		data["Result"] = rows
	} else {
		slog.Warn("no row data")
	}

	return
}

func (s *FrontWebServer) handleApiCall(u *url.URL) (apiResp *http.Response, err error) {
	// call the api
	slog.Info("calling api", slog.String("url", u.String()))
	return GetUrl(u.String())
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
