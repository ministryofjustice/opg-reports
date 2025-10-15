package api

import (
	"fmt"
	"opg-reports/report/internal/utils"
	"testing"
)

// TODO - IMPROVE THIS TEST!!
// - TEST COST IMPORT
// - ADD GROUPED API FOR UPTIME (ONLY HAS TEAM FILTER)
func TestApiServiceSQLWithOptions(t *testing.T) {
	dummyValues := map[string]string{
		"team": "true",
	}

	fields := []*Field{
		&Field{
			Key:     "cost",
			Select:  "coalesce(SUM(cost), 0) as cost",
			OrderBy: "CAST(aws_costs.cost as REAL) DESC",
		},
		&Field{
			Key:     "date",
			Select:  "strftime(:date_format, date) as date",
			Where:   "date >= :start_date AND date <= :end_date",
			GroupBy: "strftime(:date_format, date)",
			OrderBy: "strftime(:date_format, date) ASC",
		},

		&Field{
			Key:     "team",
			Select:  "aws_accounts.team_name as team_name",
			Where:   "lower(aws_accounts.team_name)=lower(:team_name)",
			GroupBy: "aws_accounts.name",
			OrderBy: "aws_accounts.name ASC",
			Value:   utils.Ptr(dummyValues["team"]),
		},
	}

	s := generateSelect(fields...)
	w := generateWhere(fields...)
	g := generateGroupBy(fields...)
	o := generateOrderBy(fields...)

	fmt.Println("selects --")
	fmt.Println(s)
	fmt.Println("wheres --")
	fmt.Println(w)
	fmt.Println("groups --")
	fmt.Println(g)
	fmt.Println("orders --")
	fmt.Println(o)
	fmt.Println("==")

	sql := BuildSelectFromFields("aws_costs", "LEFT JOIN aws_accounts ON aws_accounts.id = aws_costs.aws_account_id", fields...)
	fmt.Println(sql)
	t.FailNow()
}
