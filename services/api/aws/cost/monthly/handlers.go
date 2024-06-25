package monthly

import (
	"encoding/json"
	"net/http"
	"opg-reports/shared/aws/cost"
	"opg-reports/shared/data"
	"opg-reports/shared/dates"
	"opg-reports/shared/server"
	"strings"
	"time"
)

// func to return the date as a month (yyyy-mm)
var ym = func(i *cost.Cost) (string, string) {
	return "month", dates.Reformat(i.Date, dates.FormatYM)
}
var unit = func(i *cost.Cost) (string, string) {
	return "account_unit", strings.ToLower(i.AccountUnit)
}
var account_id = func(i *cost.Cost) (string, string) {
	return "account_id", i.AccountId
}
var account_env = func(i *cost.Cost) (string, string) {
	return "account_environment", strings.ToLower(i.AccountEnvironment)
}
var service = func(i *cost.Cost) (string, string) {
	return "service", strings.ToLower(i.Service)
}

// Group by month
var byMonth = func(item *cost.Cost) string {
	return data.ToIdxF(item, ym)
}
var byMonthUnit = func(item *cost.Cost) string {
	return data.ToIdxF(item, ym, unit)
}
var byMonthAccount = func(item *cost.Cost) string {
	return data.ToIdxF(item, ym, account_id, unit, account_env)
}
var byMonthAccountService = func(item *cost.Cost) string {
	return data.ToIdxF(item, ym, account_id, unit, account_env, service)
}

// Index: /aws/costs/v1/monthly
func (a *Api[V, F]) Index(w http.ResponseWriter, r *http.Request) {
	res := server.NewSimpleApiResponse()
	res.Start()
	res.SetStatus(http.StatusOK)
	res.End()
	content, _ := json.Marshal(res)
	a.Write(w, res.GetStatus(), content)
}

// Unit, Env & Service costs: /aws/costs/{version}/monthly/{start}/{end}/units/envs/services/{$}
// Previously "Detailed breakdown" sheet
func (a *Api[V, F]) UnitEnvironmentServices(w http.ResponseWriter, r *http.Request) {
	res := server.NewApiResponse[*cost.Cost, map[string][]*cost.Cost]()
	res.Start()
	store := a.store
	startDate, endDate := startAndEndDates(r, res)

	errs := res.GetErrors()
	if len(errs) == 0 {
		// Limit the items in the data store to those within the start & end date range
		//
		months := dates.Strings(dates.Months(startDate, endDate), dates.FormatYM)
		inMonth := func(item *cost.Cost) bool {
			return dates.InMonth(item.Date, months)
		}
		withinMonths := store.Filter(inMonth)
		items := map[string][]*cost.Cost{}

		for k, g := range withinMonths.Group(byMonthAccountService) {
			for _, v := range g.List() {
				items[k] = append(items[k], v)
			}
		}
		res.SetResult(items)
	}
	res.SetType()
	res.End()
	content, _ := json.Marshal(res)
	w.Header().Set("X-API-RES_TYPE", res.Type)
	a.Write(w, res.GetStatus(), content)

}

// Unit & Env costs: /aws/costs/{version}/monthly/{start}/{end}/units/envs/{$}
// Previously "Service And Environment" sheet
func (a *Api[V, F]) UnitEnvironments(w http.ResponseWriter, r *http.Request) {
	res := server.NewApiResponse[*cost.Cost, map[string][]*cost.Cost]()
	res.Start()
	store := a.store
	startDate, endDate := startAndEndDates(r, res)

	errs := res.GetErrors()
	if len(errs) == 0 {
		// Limit the items in the data store to those within the start & end date range
		//
		months := dates.Strings(dates.Months(startDate, endDate), dates.FormatYM)
		inMonth := func(item *cost.Cost) bool {
			return dates.InMonth(item.Date, months)
		}
		withinMonths := store.Filter(inMonth)
		items := map[string][]*cost.Cost{}

		for k, g := range withinMonths.Group(byMonthAccount) {
			for _, v := range g.List() {
				items[k] = append(items[k], v)
			}
		}
		res.SetResult(items)
	}
	res.SetType()
	res.End()
	content, _ := json.Marshal(res)
	w.Header().Set(server.ResponseTypeHeader, res.Type)
	a.Write(w, res.GetStatus(), content)

}

// Unit costs: /aws/costs/{version}/monthly/{start}/{end}/units/{$}
// Previously "Service" sheet
func (a *Api[V, F]) Units(w http.ResponseWriter, r *http.Request) {
	res := server.NewApiResponse[*cost.Cost, map[string][]*cost.Cost]()
	res.Start()
	store := a.store
	startDate, endDate := startAndEndDates(r, res)

	errs := res.GetErrors()
	if len(errs) == 0 {
		// Limit the items in the data store to those within the start & end date range
		//
		months := dates.Strings(dates.Months(startDate, endDate), dates.FormatYM)
		inMonth := func(item *cost.Cost) bool {
			return dates.InMonth(item.Date, months)
		}
		withinMonths := store.Filter(inMonth)
		items := map[string][]*cost.Cost{}

		for k, g := range withinMonths.Group(byMonthUnit) {
			for _, v := range g.List() {
				items[k] = append(items[k], v)
			}
		}
		res.SetResult(items)
	}
	res.SetType()
	res.End()
	content, _ := json.Marshal(res)

	w.Header().Set(server.ResponseTypeHeader, res.Type)
	a.Write(w, res.GetStatus(), content)

}

// Totals: /aws/costs/{version}/monthly/{start}/{end}/{$}
// Returns cost data split into with & without tax segments, then grouped by the month
// Previously "Totals" sheet
// Note: if {start} or {end} are "-" it uses current month
func (a *Api[V, F]) Totals(w http.ResponseWriter, r *http.Request) {
	res := server.NewApiResponse[*cost.Cost, map[string]map[string][]*cost.Cost]()
	res.Start()
	store := a.store
	startDate, endDate := startAndEndDates(r, res)

	errs := res.GetErrors()
	if len(errs) == 0 {
		// Limit the items in the data store to those within the start & end date range
		//
		months := dates.Strings(dates.Months(startDate, endDate), dates.FormatYM)
		inMonth := func(item *cost.Cost) bool {
			return dates.InMonth(item.Date, months)
		}
		// Filter out tax from the cost data
		//
		notTax := func(item *cost.Cost) bool {
			return strings.ToLower(item.Service) != "tax"
		}

		withTax := store.Filter(inMonth)
		withoutTax := withTax.Filter(notTax)

		items := map[string]map[string][]*cost.Cost{
			"without_tax": {},
			"with_tax":    {},
		}

		// for with & without tax we now group them by their
		// yyyy-mm
		for with, itemSet := range items {
			s := withoutTax
			if with == "with_tax" {
				s = withTax
			}

			for k, g := range s.Group(byMonth) {
				for _, v := range g.List() {
					itemSet[k] = append(itemSet[k], v)
				}
			}
		}
		// set the result
		res.SetResult(items)

	}
	res.SetType()
	res.End()
	content, _ := json.Marshal(res)
	w.Header().Set(server.ResponseTypeHeader, res.Type)
	a.Write(w, res.GetStatus(), content)
}

func startAndEndDates(r *http.Request, res server.IApiResponseBase) (startDate time.Time, endDate time.Time) {
	var err error
	now := time.Now().UTC().Format(dates.FormatYM)
	// Get the dates
	startDate, err = dates.StringToDateDefault(r.PathValue("start"), "-", now)
	if err != nil {
		res.AddErrorWithStatus(err, http.StatusConflict)
	}
	endDate, err = dates.StringToDateDefault(r.PathValue("end"), "-", now)
	if err != nil {
		res.AddErrorWithStatus(err, http.StatusConflict)
	}
	return

}
