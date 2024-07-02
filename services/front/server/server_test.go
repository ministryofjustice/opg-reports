package server

import (
	"fmt"
	th "opg-reports/internal/testhelpers"
	"opg-reports/services/api/aws/cost/monthly"
	"opg-reports/shared/aws/cost"
	"opg-reports/shared/data"
	"opg-reports/shared/dates"
	"opg-reports/shared/fake"
	"opg-reports/shared/server/response"
	"time"
)

func mockAwsCostMonthlyUnitsEnvsServices() string {
	min, max, df := th.Dates()
	route := fmt.Sprintf("/aws/costs/v1/monthly/%s/%s/units/envs/services/", min.Format(dates.FormatYM), max.Format(dates.FormatYM))
	return mockAwsCostApiResponse(route, min, max, df)
}
func mockAwsCostMonthlyUnitsEnvs() string {
	min, max, df := th.Dates()
	route := fmt.Sprintf("/aws/costs/v1/monthly/%s/%s/units/envs/", min.Format(dates.FormatYM), max.Format(dates.FormatYM))
	return mockAwsCostApiResponse(route, min, max, df)
}
func mockAwsCostMonthlyUnits() string {
	min, max, df := th.Dates()
	route := fmt.Sprintf("/aws/costs/v1/monthly/%s/%s/units/", min.Format(dates.FormatYM), max.Format(dates.FormatYM))
	return mockAwsCostApiResponse(route, min, max, df)
}
func mockAwsCostMonthlyTotals() string {
	min, max, df := th.Dates()
	route := fmt.Sprintf("/aws/costs/v1/monthly/%s/%s/", min.Format(dates.FormatYM), max.Format(dates.FormatYM))
	return mockAwsCostApiResponse(route, min, max, df)
}

func mockAwsCostApiResponse(route string, min time.Time, max time.Time, df string) string {

	store := data.NewStore[*cost.Cost]()
	count := fake.Int(1, 5)
	for i := 0; i < count; i++ {
		fk := cost.Fake(nil, min, max, df)
		store.Add(fk)
	}
	fs := th.Fs()
	mux := th.Mux()

	resp := response.NewResponse[response.ICell, response.IRow[response.ICell]]()
	api := monthly.New(store, fs, resp)
	api.Register(mux)
	w, r := th.WRGet(route)
	mux.ServeHTTP(w, r)
	str, _ := response.Stringify(w.Result())
	return str
}
