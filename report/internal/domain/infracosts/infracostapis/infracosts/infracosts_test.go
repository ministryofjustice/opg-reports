package infracosts

import (
	"fmt"
	"opg-reports/report/internal/utils/debugger"
	"opg-reports/report/internal/utils/marshal"
	"opg-reports/report/internal/utils/query"
	"testing"
)

func TestInfraCostsRequestToJson(t *testing.T) {

	req := &InfracostRequest{
		DateRange:   "2025-11..2026-02",
		Team:        "sirius",
		Account:     "true",
		Environment: "true",
		Service:     "true",
	}

	debugger.Dump(req.Months())

	reqAsJson := map[string]interface{}{}
	marshal.Convert(req, &reqAsJson)

	stmt := &query.Select{
		From:     "infracosts as base",
		Joins:    "LEFT JOIN accounts ON accounts.id = base.account_id",
		Segments: selectQuery.Segments,
	}
	q := stmt.FromRequest(reqAsJson)

	fmt.Println(q)

	t.FailNow()

}
