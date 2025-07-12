package front

import (
	"context"
	"opg-reports/report/config"
	"opg-reports/report/internal/repository/restr"
	"opg-reports/report/internal/utils"
	"testing"
	"time"
)

type holiday struct {
	Title string `json:"title"`
}
type holidays struct {
	Division string     `json:"division"`
	Events   []*holiday `json:"events"`
}
type bankHols struct {
	EnglandAndWales *holidays `json:"england-and-wales"`
	Scotland        *holidays `json:"scotland"`
	NorthernIreland *holidays `json:"northern-ireland"`
}

func bhGet(resp *bankHols) (result *holidays, err error) {
	if resp == nil {
		return
	}
	result = resp.EnglandAndWales
	result.Division = "TEST"
	return
}

func TestFrontGetFromAPI(t *testing.T) {
	var (
		err    error
		addr   = `https://www.gov.uk`
		ep     = "bank-holidays.json?month=10"
		ctx    = context.TODO()
		log    = utils.Logger("ERROR", "TEXT")
		conf   = config.NewConfig()
		client = restr.Default(ctx, log, conf)
		res    = &holidays{}
	)

	conf.Servers.Api.Addr = addr
	conf.Servers.Front.Timeout = (2 * time.Second)

	res, err = getFromAPI[*bankHols, *holidays](ctx, log, conf, client, ep, bhGet)

	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	if res == nil {
		t.Error("unexpected empty result")
	}
	if res.Division != "TEST" {
		t.Error("post processer failed to change values in resulting data")
	}

}
