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
)

type dynamicHandlerFunc func(w http.ResponseWriter, r *http.Request)

const apiVersion string = "v1"
const billingDay int = 13

// Dynamic handles all end points that require data from the api.
func (s *FrontWebServer) Dynamic(w http.ResponseWriter, r *http.Request) {
	slog.Info("dynamic handler starting", slog.String("uri", r.RequestURI))
	data := s.Nav.Data(r)

	active := s.Nav.Active(r)
	// TODO - handle this issue with a redirect to root?
	if active == nil {
		slog.Error("active nil")
		return
	}

	// Setup data object for the templates
	data["Organisation"] = s.Config.Organisation
	data["PageTitle"] = active.Name + " - "
	data["Standards"] = s.Config.Standards
	data["Result"] = nil

	usePrefix := len(active.Api) > 1
	// Handle multiple api calls for one page
	urls, _ := active.ApiUrls()
	for apiResultName, apiUrl := range urls {
		// Call API!
		u := Url(s.ApiScheme, s.ApiAddr, apiUrl)
		slog.Info("calling api",
			slog.String("url", u.String()),
			slog.String("apiUrl", apiUrl))
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
					slog.Debug("setting api res data", slog.String("key", key))
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
	// get meta data
	if meta := resp.GetMetadata(); meta != nil {
		data["Metadata"] = meta
		start, sOk := meta["StartDate"]
		end, eOK := meta["EndDate"]
		if sOk && eOK {
			startDate, _ := dates.StringToDate(start.(string))
			endDate, _ := dates.StringToDate(end.(string))
			data["StartDate"] = startDate
			data["EndDate"] = endDate
			data["Months"] = dates.Months(startDate, endDate)

		}
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
