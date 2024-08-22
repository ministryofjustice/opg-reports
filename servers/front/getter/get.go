package getter

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strings"

	"github.com/ministryofjustice/opg-reports/servers/front/config/navigation"
	"github.com/ministryofjustice/opg-reports/servers/shared/resp"
	"github.com/ministryofjustice/opg-reports/servers/shared/urls"
	"github.com/ministryofjustice/opg-reports/shared/consts"
	"github.com/ministryofjustice/opg-reports/shared/convert"
	"github.com/ministryofjustice/opg-reports/shared/dates"
	"github.com/ministryofjustice/opg-reports/shared/env"
	"github.com/ministryofjustice/opg-reports/shared/must"
	"github.com/ministryofjustice/opg-reports/shared/timer"
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

// ApiResponses calls all the api urls for the nav item
func ApiResponses(nav *navigation.NavigationItem, r *http.Request) (responses map[string]*resp.Response, requestErr error) {
	apiScheme := env.Get("API_SCHEME", consts.API_SCHEME)
	apiAddr := env.Get("API_ADDR", consts.API_ADDR)
	dataSources := nav.DataSources
	responses = map[string]*resp.Response{}

	for key, source := range dataSources {
		endpoint := source.Parsed()
		url := urls.Parse(apiScheme, apiAddr, endpoint)
		tick := timer.New()

		slog.Debug("getting data from api",
			slog.String("key", key),
			slog.String("endpoint", endpoint),
			slog.String("url", url.String()))

		httpResponse, err := GetUrl(url)
		tick.Stop()
		if err != nil || httpResponse.StatusCode != http.StatusOK {
			requestErr = err
			slog.Error("api call failed",
				slog.String("err", fmt.Sprintf("%+v", err)),
				slog.String("key", key),
				slog.Int("status", httpResponse.StatusCode),
				slog.String("endpoint", endpoint),
				slog.String("url", url.String()))
			continue
		}

		_, bytes := convert.Stringify(httpResponse)
		response, err := convert.Unmarshal[*resp.Response](bytes)
		if err != nil {
			requestErr = err
			slog.Error("failed to unmarshal api result")
			continue
		}

		slog.Info("api call details",
			slog.String("duration (s)", tick.Seconds()),
			slog.String("key", key),
			slog.String("endpoint", endpoint),
			slog.String("url", url.String()))

		responses[key] = response
	}
	return

}

// keyRemap maps the top level keys from the data to be TitleCase for the templates
func keyRemap(data map[string]interface{}, response *resp.Response, prefix string) {
	c := cases.Title(language.English)
	mapped := must.Must(convert.Map(response))

	for k, v := range mapped {
		f := c.String(strings.ReplaceAll(k, "_", " "))
		field := fmt.Sprintf("%s%s", prefix, f)
		field = strings.ReplaceAll(field, " ", "")
		slog.Debug("api result mapping: " + field)
		data[f] = v
	}
}

func parseMetadata(data map[string]interface{}) {
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

// ParseApiResponse handles parsing the repsonse object into data values
func ParseApiResponse(responses map[string]*resp.Response) (data map[string]interface{}) {
	// for title case
	c := cases.Title(language.English)
	withPrefix := len(responses) > 1
	data = map[string]interface{}{}

	for key, response := range responses {
		prefix := c.String(key)
		if !withPrefix {
			prefix = ""
		}

		keyRemap(data, response, prefix)
		// setup common values
		if response.DataAge.Max != "" {
			data["DataAgeMax"] = response.DataAge.Max
		}
		if response.DataAge.Min != "" {
			data["DataAgeMin"] = response.DataAge.Min
		}
		// parse metadata
		if _, ok := data["Metadata"]; ok {
			parseMetadata(data)
		}
	}
	return
}
