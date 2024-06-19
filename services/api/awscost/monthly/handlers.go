package monthly

import (
	"net/http"
	"opg-reports/shared/dates"
	"slices"
)

// Index: /aws/costs/v1/monthly
func (a *Api[V, F]) Index(w http.ResponseWriter, r *http.Request) {
	res := a.Response()
	res.Start()
	all := a.store.List()
	res.Set(all, http.StatusOK, nil)
	res.End()
	a.Write(w, res)
}

// Totals: /aws/costs/{version}/monthly/{start}/{end}/{$}
func (a *Api[V, F]) Totals(w http.ResponseWriter, r *http.Request) {
	errs := []error{}
	items := []V{}
	store := a.Store()
	res := a.Response()
	res.Start()

	// Get the dates
	startDate, err := dates.StringToDate(r.PathValue("start"))
	if err != nil {
		errs = append(errs, err)
	}
	endDate, err := dates.StringToDate(r.PathValue("end"))
	if err != nil {
		errs = append(errs, err)
	}

	if len(errs) == 0 {
		months := dates.Strings(dates.Months(startDate, endDate), dates.FormatYM)
		// Filter function to find items whose date is within the months range
		dateF := func(item V) bool {
			d, e := dates.StringToDate(item.Date)
			if e == nil && slices.Contains(months, d.Format(dates.FormatYM)) {
				return true
			}
			return false
		}
		store = store.Filter(dateF)
		items = store.List()
	}

	res.Set(items, http.StatusOK, errs)
	res.End()
	a.Write(w, res)
}
