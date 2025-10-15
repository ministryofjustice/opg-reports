package api

import (
	"fmt"
	"opg-reports/report/internal/utils"
	"strings"
	"testing"
)

type tStmtGen struct {
	Fields   []*Field
	Expected string
}

// TestApiServiceGenerateOrder tests various combinations fields for order bys
func TestApiServiceGenerateOrder(t *testing.T) {

	var tests = []*tStmtGen{
		// test a simple single field that should not generate an order by
		{
			Fields:   []*Field{{Select: "COUNT(*) as counter"}},
			Expected: "",
		},
		// simple order by that should always be there
		{
			Fields: []*Field{
				{
					Select:  "strftime(:date_format, date) as date",
					Where:   "(date >= :start_date AND date <= :end_date)",
					GroupBy: "strftime(:date_format, date)",
					OrderBy: "strftime(:date_format, date) ASC",
					Value:   utils.Ptr("true"),
				},
			},
			Expected: "strftime(:date_format, date) ASC",
		},
		// version thats always ignored as Value is empty string
		{
			Fields: []*Field{
				{
					Select:  "strftime(:date_format, date) as date",
					Where:   "(date >= :start_date AND date <= :end_date)",
					GroupBy: "strftime(:date_format, date)",
					OrderBy: "strftime(:date_format, date) ASC",
					Value:   utils.Ptr(""),
				},
			},
			Expected: "",
		},
		// version thats always ignored as Value is something other than true (so a filter)
		{
			Fields: []*Field{
				{
					Select:  "strftime(:date_format, date) as date",
					Where:   "(date >= :start_date AND date <= :end_date)",
					GroupBy: "strftime(:date_format, date)",
					OrderBy: "strftime(:date_format, date) ASC",
					Value:   utils.Ptr("2025-01"),
				},
			},
			Expected: "",
		},
	}

	for i, test := range tests {
		var actual = generateOrderBy(test.Fields...)
		actual = strings.ReplaceAll(actual, "\n", "")
		if test.Expected != actual {
			t.Errorf("[%d] orderby generation error, actual does not match expected:\nactual:[%s]\nexpected:[%s]\n", i, actual, test.Expected)
		}
	}
}

// TestApiServiceGenerateGroup tests various combinations fields for groups
func TestApiServiceGenerateGroup(t *testing.T) {

	var tests = []*tStmtGen{
		// test a simple single field that should not generate a group
		{
			Fields:   []*Field{{Select: "COUNT(*) as counter"}},
			Expected: "",
		},
		// should include a group by for the date
		{
			Fields: []*Field{
				{Select: "COUNT(*) as counter"},
				{
					Select:  "strftime(:date_format, date) as date",
					Where:   "(date >= :start_date AND date <= :end_date)",
					GroupBy: "strftime(:date_format, date)",
					OrderBy: "strftime(:date_format, date) ASC",
				},
			},
			Expected: "strftime(:date_format, date)",
		},
		// no grouping as value is set, but nil
		{
			Fields: []*Field{
				{Select: "COUNT(*) as counter"},
				{Select: "account", GroupBy: "account"},
				{
					Select:  "strftime(:date_format, date) as date",
					GroupBy: "strftime(:date_format, date)",
					Value:   utils.Ptr(""),
				}},
			Expected: "account",
		},
	}

	for i, test := range tests {
		var actual = generateGroupBy(test.Fields...)
		actual = strings.ReplaceAll(actual, "\n", "")
		if test.Expected != actual {
			t.Errorf("[%d] groupby generation error, actual does not match expected:\nactual:[%s]\nexpected:[%s]\n", i, actual, test.Expected)
		}
	}
}

// TestApiServiceGenerateWheres tests various combinations fields for wheres
func TestApiServiceGenerateWheres(t *testing.T) {

	var tests = []*tStmtGen{
		// test a simple single field that should not generate a where
		{
			Fields:   []*Field{{Select: "COUNT(*) as counter"}},
			Expected: "",
		},
		// simple where with two clauses
		{
			Fields: []*Field{
				{Select: "COUNT(*) as counter", Where: "COUNT(*) > 5"},
				{Select: "strftime(:date_format, date) as date", Where: "(date >= :start_date AND date <= :end_date)"},
			},
			Expected: "COUNT(*) > 5 AND (date >= :start_date AND date <= :end_date)",
		},
		// one field is not for wheres (so value == true)
		{
			Fields: []*Field{
				{Select: "COUNT(*) as counter", Where: "COUNT(*) > 5"},
				{Select: "strftime(:date_format, date) as date", Where: "(date >= :start_date AND date <= :end_date)", Value: utils.Ptr("true")},
			},
			Expected: "COUNT(*) > 5",
		},
	}

	for i, test := range tests {
		var actual = generateWhere(test.Fields...)
		actual = strings.ReplaceAll(actual, "\n", "")
		if test.Expected != actual {
			t.Errorf("[%d] where generation error, actual does not match expected:\nactual:[%s]\nexpected:[%s]\n", i, actual, test.Expected)
		}
	}
}

// TestApiServiceGenerateSelect tests various combinations of selects
func TestApiServiceGenerateSelect(t *testing.T) {

	var tests = []*tStmtGen{
		// test a simple single field
		{
			Fields: []*Field{
				{Select: "COUNT(*) as counter"},
			},
			Expected: "COUNT(*) as counter",
		},
		// test multiple fields, including duplicates, that should all be there
		{
			Fields: []*Field{
				{Select: "COUNT(*) as counter"},
				{Select: "COUNT(*) as counter"},
				{Select: "coalesce(SUM(cost), 0) as cost"},
				{Select: "strftime(:date_format, date) as date"},
			},
			Expected: "COUNT(*) as counter,COUNT(*) as counter,coalesce(SUM(cost), 0) as cost,strftime(:date_format, date) as date",
		},
		// test multiple fields, with one that should be ignored
		{
			Fields: []*Field{
				{Select: "COUNT(*) as counter"},
				{Select: "COUNT(*) as counter", Value: utils.Ptr("")},
				{Select: "coalesce(SUM(cost), 0) as cost"},
				{Select: "strftime(:date_format, date) as date"},
			},
			Expected: "COUNT(*) as counter,coalesce(SUM(cost), 0) as cost,strftime(:date_format, date) as date",
		},
	}

	for i, test := range tests {
		var actual = generateSelect(test.Fields...)
		actual = strings.ReplaceAll(actual, "\n", "")
		if test.Expected != actual {
			t.Errorf("[%d] select generation error, actual does not match expected:\n actual:[%s]\n expected:[%s]\n", i, actual, test.Expected)
		}
	}
}

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
