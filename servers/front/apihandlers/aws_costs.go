package apihandlers

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/ministryofjustice/opg-reports/servers/api/aws_costs"
	"github.com/ministryofjustice/opg-reports/servers/shared/datarow"
	"github.com/ministryofjustice/opg-reports/servers/shared/srvr/httphandler"
	"github.com/ministryofjustice/opg-reports/shared/convert"
)

type AWSCosts struct{}

// IsYTD flags if this is the year to date endpoint, as
// this result should not get converted to rows
func (costs *AWSCosts) IsYTD(remote *httphandler.HttpHandler) bool {
	segment := "/ytd/"
	return strings.Contains(remote.Url.String(), segment)
}

func (costs *AWSCosts) Handles(url string) bool {
	segment := "/aws-costs/"
	return strings.Contains(url, segment)
}

func (costs *AWSCosts) Handle(data map[string]interface{}, key string, remote *httphandler.HttpHandler, w http.ResponseWriter, r *http.Request) {
	var (
		err      error
		response *http.Response
		api      *aws_costs.ApiResponse
	)
	response, err = remote.Result()
	slog.Info("called api",
		slog.Float64("duration", remote.Duration),
		slog.Int("status_code", remote.StatusCode),
		slog.String("uri", remote.Url.String()))

	if err != nil || remote.StatusCode != http.StatusOK {
		slog.Error("error with api response", slog.String("err", fmt.Sprintf("%+v", err)), slog.Int("status_code", remote.StatusCode))
		return
	}

	api, err = convert.UnmarshalR[*aws_costs.ApiResponse](response)
	if err != nil {
		slog.Error("error with unmarshal response", slog.String("err", err.Error()))
		return
	}
	// -- if its not the ytd endpoint, convert to rows
	if !costs.IsYTD(remote) {
		awsCostsToRows(api)
	}
	data[key] = api
}

func awsCostsToRows(re *aws_costs.ApiResponse) {
	mapped, _ := convert.Maps(re.Result)
	intervals := map[string][]string{"interval": re.DateRange}
	values := map[string]string{"interval": "total"}
	re.Rows = datarow.DataRows(mapped, re.Columns, intervals, values)
}
