package monthly

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"opg-reports/shared/aws/cost"
	"opg-reports/shared/data"
	"opg-reports/shared/dates"
	"testing"
	"time"
)

func TestServicesApiAwsCostMonthlyHandlerIndex(t *testing.T) {
	fs := testFs()
	mux := testMux()
	store := data.NewStore[*cost.Cost]()
	api := New(store, fs)
	api.Register(mux)

	route := "/aws/costs/v1/monthly/"
	w, r := testWRGet(route)

	mux.ServeHTTP(w, r)

	res := decode[*ApiResponse](w.Result().Body)

	if len(res.Errors) != 0 {
		t.Errorf("found error when not expected")
	}

	if res.Result != nil || res.GetResults() != nil {
		t.Errorf("result should be empty")
	}

}

func TestServicesApiAwsCostMonthlyHandlerTotals(t *testing.T) {
	fs := testFs()
	mux := testMux()
	min, max, df := testDates()
	l := 5
	// out of bounds
	overm := time.Date(max.Year()+1, 1, 1, 0, 0, 0, 0, time.UTC)
	overmx := time.Date(max.Year()+2, 1, 1, 0, 0, 0, 0, time.UTC)

	store := data.NewStore[*cost.Cost]()
	for i := 0; i < l; i++ {
		store.Add(cost.Fake(nil, min, max, df))
	}
	for i := 0; i < 10; i++ {
		store.Add(cost.Fake(nil, overm, overmx, df))
	}

	api := New(store, fs)
	api.Register(mux)

	route := fmt.Sprintf("/aws/costs/v1/monthly/%s/%s/", min.Format(dates.FormatYM), max.Format(dates.FormatYM))
	w, r := testWRGet(route)
	mux.ServeHTTP(w, r)

	_, b := strResponse(w.Result())
	// fmt.Println(str)
	res := map[string]interface{}{}
	json.Unmarshal(b, &res)

	fmt.Printf("%+v\n", res)

	// is := res.Result.([]interface{})
	// items, _ := data.FromInterfaces[*cost.Cost](is)

	// if len(items) != 100 {
	// 	t.Errorf("incorrect number of items returned.")
	// }

}

func strResponse(r *http.Response) (string, []byte) {
	b, _ := io.ReadAll(r.Body)
	return string(b), b
}
