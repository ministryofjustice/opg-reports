package getter

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/ministryofjustice/opg-reports/servers/front/config"
	"github.com/ministryofjustice/opg-reports/servers/front/config/navigation"
	"github.com/ministryofjustice/opg-reports/servers/shared/resp"
	"github.com/ministryofjustice/opg-reports/servers/shared/urls"
	"github.com/ministryofjustice/opg-reports/shared/consts"
	"github.com/ministryofjustice/opg-reports/shared/convert"
	"github.com/ministryofjustice/opg-reports/shared/dates"
	"github.com/ministryofjustice/opg-reports/shared/env"
	"github.com/ministryofjustice/opg-reports/shared/must"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func GetUrl(url *url.URL) (resp *http.Response, err error) {
	req, err := http.NewRequest(http.MethodGet, url.String(), nil)
	if err != nil {
		return
	}
	apiClient := http.Client{Timeout: consts.API_TIMEOUT}
	resp, err = apiClient.Do(req)
	return
}

// Api uses the passed in nav item to work out the api urls to call, fetches the data
// and returns map of that
// - api urls have substitutions like {month} resolved before calling
// - multiple apis urls to call generate data with prefix based on their config name
//   - "list": "/url"
//   - "home": "/home"
//     would generate List<Key> and Home<key> data
func Api(conf *config.Config, nav *navigation.NavigationItem, r *http.Request) (data map[string]interface{}, dataErr error) {
	var apiScheme = env.Get("API_SCHEME", consts.API_SCHEME)
	var apiAddr = env.Get("API_ADDR", consts.API_ADDR)
	// for title case
	c := cases.Title(language.English)

	data = map[string]interface{}{}
	dataSources := nav.DataSources

	withPrefix := len(dataSources) > 1

	for key, source := range dataSources {
		endpoint := source.Parsed()
		url := urls.Parse(apiScheme, apiAddr, endpoint)

		s := time.Now().UTC()
		slog.Debug("getting data from api",
			slog.String("key", key),
			slog.String("endpoint", endpoint),
			slog.String("url", url.String()))

		httpResponse, err := GetUrl(url)
		e := time.Now().UTC()
		duration := e.Sub(s)

		// if failed, skip rest of loop
		if err != nil || httpResponse.StatusCode != http.StatusOK {
			dataErr = err
			slog.Error("api call failed",
				slog.String("err", fmt.Sprintf("%+v", err)),
				slog.String("key", key),
				slog.Int("status", httpResponse.StatusCode),
				slog.String("endpoint", endpoint),
				slog.String("url", url.String()))
			continue
		}

		_, bytes := convert.Stringify(httpResponse)
		// response := resp.New()
		response, err := convert.Unmarshal[*resp.Response](bytes)
		if err != nil {
			dataErr = err
			slog.Error("failed to unmarshal api result")
			continue
		}

		slog.Info("api call details",
			slog.Float64("duration (s)", duration.Seconds()),
			slog.String("key", key),
			slog.String("endpoint", endpoint),
			slog.String("url", url.String()))

		prefix := c.String(key)
		if !withPrefix {
			prefix = ""
		}
		// setup key naming
		for k, v := range must.Must(convert.Map(response)) {
			f := c.String(strings.ReplaceAll(k, "_", " "))
			field := fmt.Sprintf("%s%s", prefix, f)
			field = strings.ReplaceAll(field, " ", "")
			slog.Debug("api result mapping: " + field)
			data[f] = v
		}

		// setup common values
		if response.DataAge.Max != "" {
			data["DataAgeMax"] = response.DataAge.Max
		}
		if response.DataAge.Min != "" {
			data["DataAgeMin"] = response.DataAge.Min
		}

		// some common metadata items
		if _, ok := data["Metadata"]; ok {
			metadata := data["Metadata"].(map[string]interface{})
			// start & end date
			if sd, ok := metadata["start_date"]; ok {
				data["StartDate"] = dates.Time(sd.(string))
			}
			if ed, ok := metadata["end_date"]; ok {
				data["EndDate"] = dates.Time(ed.(string))
			}
			// date range
			if dr, ok := metadata["date_range"]; ok {
				dataRange := []string{}
				for _, dr := range dr.([]interface{}) {
					dataRange = append(dataRange, dr.(string))
				}
				data["DateRange"] = dataRange
			}
			// columns - version one, just based on the raw versions
			if metaCols, ok := metadata["columns"]; ok {
				cols := metaCols.(map[string]interface{})
				data["Columns"] = cols
				detailed := map[string][]interface{}{}
				for col, val := range cols {
					detailed[col] = val.([]interface{})
				}
				data["ColumnsDetailed"] = detailed
			}
			// column ordering
			if order, ok := metadata["column_ordering"]; ok {
				colNames := []string{}
				for _, col := range order.([]interface{}) {
					colNames = append(colNames, col.(string))
				}
				data["ColumnsOrdered"] = colNames
			}

		}

	}

	return
}
