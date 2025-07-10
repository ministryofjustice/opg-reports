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

func TestFrontGetFromAPI(t *testing.T) {
	var (
		err    error
		addr   = `https://www.gov.uk`
		ep     = "bank-holidays.json?month=10"
		ctx    = context.TODO()
		log    = utils.Logger("ERROR", "TEXT")
		conf   = config.NewConfig()
		client = restr.Default(ctx, log, conf)
		res    = &bankHols{}
	)

	conf.Servers.Api.Addr = addr
	conf.Servers.Front.Timeout = (2 * time.Second)

	srv := Default[*bankHols](ctx, log, conf)
	res, err = srv.GetFromAPI(client, ep, func(res *bankHols) (e error) {
		if res == nil {
			return
		}
		res.EnglandAndWales.Division = "TEST"
		return
	})

	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	if res == nil {
		t.Error("unexpected empty result")
	}
	if res.EnglandAndWales.Division != "TEST" {
		t.Error("post processer failed to change values in resulting data")
	}

}
