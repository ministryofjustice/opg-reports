package server

import (
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"net/url"
	"opg-reports/services/front/tmpl"
	"opg-reports/shared/dates"
	"opg-reports/shared/server/resp"
)

type dynamicHandlerFunc func(w http.ResponseWriter, r *http.Request)

const apiVersion string = "v1"
const billingDay int = 13

// Dynamic handles all end points that require data from the api.
// The api data used here should match a resp.Response struct
//
// Uses data from the active navigation item to determine what api
// url's should be called, iterates over those urls and fetches
// the api data, calling `parseResponse` for each
//
// Generates a map (`data`) containing all details for rendering
// dynamic and static templates. Capitalises them to match struct naming
// conventions.
//
//   - Organisation: set to `organisation` value from config file
//   - PageTitle: currently active navigation items Name
//   - Standards: github repository standards from the config file
//
// To allow multiple api calls for one page, the following properties are
// generated per api url and all would have a `_<key>` suffix; where `<key>`
// relates to the key in the config for that url. If there is only one
// api url to call, then no suffixes are added.
//
//   - Result: data from the api, generally the `.Body`
//   - Metadata: setup within `parseResponse`, should contain `.Metadata`
//   - StartDate & EndDate: only set when a start & end date are found in the `.Metadata`
//   - Months: all months between the StartDate & EndDate
//   - DateAgeMin & DataAgeMax: timestamps showing the youngest and oldest data entry
//   - Headings & Footer: the header and footer rows of the dataset
//
// The above map is then passed into the template rendering and the result
// is then returned (via w.Write)
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
	// determine if we should use a a prefix on the result data
	usePrefix := len(active.Api) > 1
	// Handle multiple api calls for one page
	urls, _ := active.ApiUrls()
	for apiResultName, apiUrl := range urls {
		// Call API!
		u := Url(s.ApiScheme, s.ApiAddr, apiUrl)
		slog.Info("calling api",
			slog.String("url", u.String()),
			slog.String("apiUrl", apiUrl))
		// no error from the api, and no error from parsing the api resposne
		// into local data, then set to the data map ready for the template parsing
		if apiResp, err := s.handleApiCall(u); err == nil {
			// if there is no error, try to then parse the response
			if apiData, parseErr := s.parseResponse(apiResp); parseErr == nil {
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
					slog.String("err", fmt.Sprintf("%+v", parseErr)))
			}
		} else {
			slog.Error("dynamic handler error from api",
				slog.String("url", u.String()),
				slog.String("err", fmt.Sprintf("%+v", err)))
		}
	}

	t, err := template.New(active.TemplateName).Funcs(tmpl.Funcs()).ParseFiles(s.templateFiles...)
	if err != nil {
		slog.Error("dynamic error", slog.String("err", fmt.Sprintf("%v", err)))
		return
	}

	s.Write(w, 200, t, active.TemplateName, data)

}

// parseResponse takes a http.Response from the api call and returns useful data as a map
// Returns:
//   - Result: the `.Body` value
//   - Metadata: the `.Metadata` value
//   - StartDate & EndDate: only set when a start & end date are found in the `.Metadata`
//   - Months: all months between the StartDate & EndDate
//   - DateAgeMin & DataAgeMax: timestamps showing the youngest and oldest data entry
//   - Headings & Footer: the header and footer rows of the dataset
func (s *FrontWebServer) parseResponse(apiResp *http.Response) (data map[string]interface{}, err error) {
	data = map[string]interface{}{}

	res := resp.New()
	err = resp.FromHttp(apiResp, res)

	if err != nil {
		slog.Error("parse error")
		return
	}
	// get meta data
	meta := res.Metadata
	if len(meta) > 0 {
		data["Metadata"] = meta
		start, sOk := meta["startDate"]
		end, eOK := meta["endDate"]
		if sOk && eOK {
			startDate, _ := dates.StringToDate(start.(string))
			endDate, _ := dates.StringToDate(end.(string))
			data["StartDate"] = startDate
			data["EndDate"] = endDate
			data["Months"] = dates.Months(startDate, endDate)

		}
	}

	// get min / max times
	if min := res.GetDataAgeMin(); min.Format(dates.FormatY) != dates.ErrYear {
		data["DataAgeMin"] = min
	}
	if max := res.GetDataAgeMax(); max.Format(dates.FormatY) != dates.ErrYear {
		data["DataAgeMax"] = max
	}
	// If the result is nil (failed parsing), return
	result := res.Result
	if result == nil {
		slog.Error("empty result")
		return
	}

	heading := result.Head
	if heading != nil {
		data["Headings"] = heading.All()
	}

	footer := result.Foot
	if footer != nil {
		data["Footer"] = footer.All()
	}

	// fetch the resulting rows, return if they are empty
	rows := result.Body
	if rows != nil {
		data["Result"] = rows
	} else {
		slog.Warn("no row data")
	}

	return
}

func (s *FrontWebServer) handleApiCall(u *url.URL) (apiResp *http.Response, err error) {
	// call the api
	slog.Info("handle api call", slog.String("url", u.String()))
	return GetUrl(u.String())
}
