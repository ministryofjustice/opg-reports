package server

import (
	"fmt"
	"net/http"
	"opg-reports/internal/testhelpers"
	"opg-reports/services/api/aws/cost/monthly"
	"opg-reports/services/front/cnf"
	"opg-reports/services/front/tmpl"
	"opg-reports/shared/aws/cost"
	"opg-reports/shared/data"
	"opg-reports/shared/dates"
	"opg-reports/shared/fake"
	"opg-reports/shared/files"
	"opg-reports/shared/logger"
	"opg-reports/shared/server/resp"
	"os"
	"strings"
	"testing"
)

func TestFrontServerDynamicHandlerTestAwsMonthlyTotalsTemplate(t *testing.T) {
	logger.LogSetup()
	// dates
	min, max, df := testhelpers.Dates()
	// -- get a faked api response
	// items
	count := 10
	store := data.NewStore[*cost.Cost]()
	services := []string{"ec2", "ecs", "tax", "rds", "r53"}
	for i := 0; i < count; i++ {
		c := cost.Fake(nil, min, max, df)
		c.Service = fake.Choice(services)
		store.Add(c)
	}
	// params
	allow := []string{"version", "start", "end"}
	params := map[string][]string{
		"start": {min.Format(dates.FormatYM)},
		"end":   {max.Format(dates.FormatYM)},
	}
	// functions
	var headerF = monthly.DisplayHeadFunctions(params)["monthlyTotals"]
	var rowF = monthly.DisplayRowFunctions(params)["monthlyTotals"]
	// setup route
	apiroute := fmt.Sprintf("/aws/costs/v1/monthly/%s/%s/",
		min.Format(dates.FormatYM),
		max.Format(dates.FormatYM))
	//
	w, r := testhelpers.WRGet(apiroute)
	r.SetPathValue("start", min.Format(dates.FormatYM))
	r.SetPathValue("end", max.Format(dates.FormatYM))
	ep := testhelpers.MockEndpoint[*cost.Cost](store, allow, headerF, rowF, w, r)
	ep.ProcessRequest(w, r)
	apiresponse, _ := resp.Stringify(w.Result())

	// -- test the front end response
	ms := testhelpers.MockServer(apiresponse, http.StatusOK)
	defer ms.Close()

	tDir := "../templates/"
	dfSys := os.DirFS(tDir).(files.IReadFS)
	f := files.NewFS(dfSys, tDir)
	templates := tmpl.Files(f, tDir)

	route := "/costs/aws/totals/"
	conf, _ := cnf.Load(testRealServerCnf())
	// create new
	s := New(conf, templates, "", "")
	// point the totals route to look at the test api
	s.Nav.Get(route).Api = map[string]string{"url": ms.URL}
	//
	mux := testhelpers.Mux()
	s.Register(mux)
	// now we fetch the local route, which should them call the mocked
	// server
	w, r = testhelpers.WRGet(route)
	mux.ServeHTTP(w, r)

	if w.Result().StatusCode != http.StatusOK {
		t.Errorf("unexpected error")
	}

	str, _ := resp.Stringify(w.Result())
	if !strings.Contains(str, "AWS Costs") {
		fmt.Println(str)
		t.Errorf("failed, costs header not found")
	}

}
