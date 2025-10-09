package api

import (
	"fmt"
	"opg-reports/report/internal/utils"
	"testing"
)

// coalesce(SUM(cost), 0) as cost,
// strftime(:date_format, date) as date
func TestApiServiceSQLWithOptions(t *testing.T) {
	dummyValues := map[string]string{
		"team": "true",
	}

	fields := []*Field{
		&Field{
			Key:       "cost",
			SelectAs:  "coalesce(SUM(cost), 0) as cost",
			OrderByAs: "CAST(aws_costs.cost as REAL) DESC",
		},
		&Field{
			Key:       "date",
			SelectAs:  "strftime(:date_format, date) as date",
			WhereAs:   "(date >= :start_date AND date <= :end_date)",
			GroupByAs: "strftime(:date_format, date)",
			OrderByAs: "strftime(:date_format, date) ASC",
		},
		&Field{
			Key:       "team",
			SelectAs:  "aws_accounts.team_name as team_name",
			WhereAs:   "lower(aws_accounts.team_name)=lower(:team_name)",
			GroupByAs: "aws_accounts.name",
			OrderByAs: "aws_accounts.name ASC",
			Value:     utils.Ptr(dummyValues["team"]),
		},
	}

	s := Selects(fields...)
	w := Wheres(fields...)
	g := Groups(fields...)
	o := Orders(fields...)

	fmt.Println("selects --")
	fmt.Println(s)
	fmt.Println("wheres --")
	fmt.Println(w)
	fmt.Println("groups --")
	fmt.Println(g)
	fmt.Println("orders --")
	fmt.Println(o)
	fmt.Println("==")

	sql := BuildGroupSelect("aws_costs", "LEFT JOIN aws_accounts ON aws_accounts.id = aws_costs.aws_account_id", fields...)
	fmt.Println(sql)
	t.FailNow()
}
