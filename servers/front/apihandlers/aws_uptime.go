package apihandlers

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/ministryofjustice/opg-reports/servers/api/aws_uptime"
	"github.com/ministryofjustice/opg-reports/servers/shared/datarow"
	"github.com/ministryofjustice/opg-reports/servers/shared/srvr/httphandler"
	"github.com/ministryofjustice/opg-reports/shared/convert"
)

type AWSUptime struct{}

func (up *AWSUptime) Handles(url string) bool {
	segment := "/aws-uptime/"
	return strings.Contains(url, segment)
}
func (up *AWSUptime) Handle(data map[string]interface{}, key string, remote *httphandler.HttpHandler, w http.ResponseWriter, r *http.Request) {
	var (
		err      error
		response *http.Response
		api      *aws_uptime.ApiResponse
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

	api, err = convert.UnmarshalR[*aws_uptime.ApiResponse](response)
	if err != nil {
		slog.Error("error with unmarshal response", slog.String("err", err.Error()))
		return
	}
	awsUptimeToRows(api)
	data[key] = api
}

func awsUptimeToRows(re *aws_uptime.ApiResponse) {
	mapped, _ := convert.Maps(re.Result)
	intervals := map[string][]string{"interval": re.DateRange}
	values := map[string]string{"interval": "average"}
	re.Rows = datarow.DataRows(mapped, re.Columns, intervals, values)
}
