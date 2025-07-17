package awsr

import (
	"testing"
	"time"

	"opg-reports/report/config"
	"opg-reports/report/internal/utils"

	"github.com/aws/aws-sdk-go-v2/service/costexplorer/types"
)

func TestCEGetCosts(t *testing.T) {
	var (
		err  error
		ctx  = t.Context()
		conf = config.NewConfig()
		log  = utils.Logger("ERROR", "TEXT")
	)
	client := &mockClientCostExplorerGetter{}
	sv := Default(ctx, log, conf)
	data, err := sv.GetCostData(client, &GetCostDataOptions{
		StartDate:   utils.TimeReset(time.Now().UTC().AddDate(0, -4, 0), "month").Format(utils.DATE_FORMATS.YMD),
		EndDate:     utils.TimeReset(time.Now().UTC().AddDate(0, -3, 0), "month").Format(utils.DATE_FORMATS.YMD),
		Granularity: types.GranularityMonthly,
	})
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
		t.FailNow()
	}
	if len(data) <= 0 {
		t.Errorf("should return dummy cost values")
	}

}
