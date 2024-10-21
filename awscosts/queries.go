package awscosts

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/jmoiron/sqlx"
)

// SingularStatement is a string used as enum-esque
// type contraints for sql queries that return
// single value as a result and use series of optional
// arguments rather than named parameter
type SingularStatement string

// RowCount returns the total number of records within the database
const RowCount SingularStatement = `
SELECT
	count(*) as row_count
FROM aws_costs`

// TotalWithinDateRange returns the sum of cost field for all
// records with the date range passed (start_date, end_date)
const TotalInDateRange SingularStatement = `
SELECT
    coalesce(SUM(cost), 0) as total
FROM aws_costs
WHERE
    date >= ?
	AND date < ?
LIMIT 1
`
const TotalInDateRangeWithoutTax SingularStatement = `
SELECT
    coalesce(SUM(cost), 0) as total
FROM aws_costs
WHERE
    date >= ?
	AND date < ?
	AND service != 'Tax'
LIMIT 1
`

// ManyStatement is a string, used as a enum type
// for the various sql queries we want to run
// that allows named parameters etc and return
// multiple results
type ManyStatement string

// TotalsWithAndWithoutTax is used to calculate the total costs
// within the given date range (>= :start_date, < :end_date) and
// splits that based on the `service` being 'Tax' or not
// Used to show top line numbers without VAT etc
const TotalsWithAndWithoutTax ManyStatement = `
SELECT
    'Including Tax' as service,
    coalesce(SUM(cost), 0) as cost,
    strftime(:date_format, date) as date
FROM aws_costs as incTax
WHERE
    incTax.date >= :start_date
    AND incTax.date < :end_date
GROUP BY strftime(:date_format, incTax.date)
UNION ALL
SELECT
    'Excluding Tax' as service,
    coalesce(SUM(cost), 0) as cost,
    strftime(:date_format, date) as date
FROM aws_costs as excTax
WHERE
    excTax.date >= :start_date
    AND excTax.date < :end_date
	AND excTax.service != 'Tax'
GROUP BY strftime(:date_format, date)
ORDER by date ASC
`

// PerUnit groups the cost data by the time period and unit
// and limits the data to the date range specfied
// (>= :start_date, < :end_date) returning the SUM costs
// for each grouping
const PerUnit ManyStatement = `
SELECT
    unit,
    coalesce(SUM(cost), 0) as cost,
    strftime(:date_format, date) as date
FROM aws_costs
WHERE
    date >= :start_date
    AND date < :end_date
GROUP BY strftime(:date_format, date), unit
ORDER by strftime(:date_format, date), unit ASC
`

// PerUnitEnvironment groups cost date by the date period, unit
// and environment values in the row and restricts the dataset to the
// date range passed (>= :start_date, < :end_date) returning the
// SUM of each grouping as `cost`
// If the environment field is "null" then we default to "production"
// as several accounts (like managment / identity ) have only one
// version
const PerUnitEnvironment ManyStatement = `
SELECT
    unit,
	IIF(environment != "null", environment, "production") as environment,
    coalesce(SUM(cost), 0) as cost,
    strftime(:date_format, date) as date
FROM aws_costs
WHERE
    date >= :start_date
    AND date < :end_date
GROUP BY strftime(:date_format, date), unit, environment
ORDER by strftime(:date_format, date), unit, environment ASC
`

// Detailed is used to show the cost of each type of service per account and
// org for the time period passed along - allowing to track costs changes
// for s3 etc overtime at a granular level
// Limits the data to the date range expressed (>= :start_date, < :end_date)
const Detailed ManyStatement = `
SELECT
    unit,
	IIF(environment != "null", environment, "production") as environment,
	organisation,
	account_id,
	account_name,
	label,
	service,
    coalesce(SUM(cost), 0) as cost,
    strftime(:date_format, date) as date
FROM aws_costs
WHERE
    date >= :start_date
    AND date < :end_date
GROUP BY strftime(:date_format, date), unit, environment, organisation, account_id, service
ORDER by strftime(:date_format, date), unit, environment, account_id ASC
`

// DetailedForUnit is an extension of Detailed and further restricts the data set
// to match the unit passed
const DetailedForUnit ManyStatement = `
SELECT
    unit,
	IIF(environment != "null", environment, "production") as environment,
	organisation,
	account_id,
	account_name,
	label,
	service,
    coalesce(SUM(cost), 0) as cost,
    strftime(:date_format, date) as date
FROM aws_costs
WHERE
    date >= :start_date
    AND date < :end_date
	AND unit = :unit
GROUP BY strftime(:date_format, date), unit, environment, organisation, account_id, service
ORDER by strftime(:date_format, date), unit, environment, account_id ASC
`

// GetOne returns a raw value from a query statments being used - this is typically a counter or the
// result of a sum operation ran against a series of rows
//
// Uses optional, ordered arguments instead of named parameter struct
func GetOne(ctx context.Context, db *sqlx.DB, query SingularStatement, args ...interface{}) (result interface{}, err error) {

	switch query {
	case RowCount:
		fallthrough
	case TotalInDateRange:
		fallthrough
	case TotalInDateRangeWithoutTax:
		err = db.GetContext(ctx, &result, string(query), args...)
	default:
		err = fmt.Errorf("unknown SingularStatement passed [%v]", query)
	}
	return
}

// Parameters are used to as named parameters on sqlx queries
// via the Query function and cover all possible
type Parameters struct {
	StartDate  string `json:"start_date,omitempty" db:"start_date"`   // StartDate is the lower bound of date based query
	EndDate    string `json:"end_date,omitempty" db:"end_date"`       // EndDate is the upper bound of date based query
	DateFormat string `json:"date_format,omitempty" db:"date_format"` // Date format to use for strftime with query
	Unit       string `json:"unit,omitempty" db:"unit"`               // Unit to filter the data by
}

// Map uses json marshal & unmarshal to return a map version of
// this struct
func (self *Parameters) Map() (m map[string]string) {
	m = map[string]string{}

	if bytes, err := json.Marshal(self); err == nil {
		json.Unmarshal(bytes, &m)
	}
	return
}

// GetMany runs the known statement against using the parameters as named values within them and returns the
// result as a slice of []*Cost
func GetMany(ctx context.Context, db *sqlx.DB, query ManyStatement, params *Parameters) (results []*Cost, err error) {
	var statement *sqlx.NamedStmt
	results = []*Cost{}
	// We have a switch as we do want to restrict what queries are allowed
	// so result and functions used in sqlx match
	switch query {
	case TotalsWithAndWithoutTax:
		fallthrough
	case PerUnit:
		fallthrough
	case PerUnitEnvironment:
		fallthrough
	case Detailed:
		fallthrough
	case DetailedForUnit:
		// Check the parameters passed is valid for the query
		if err = ValidateParameters(params, Needs(query)); err != nil {
			return
		}
		if statement, err = db.PrepareNamedContext(ctx, string(query)); err == nil {
			err = statement.SelectContext(ctx, &results, params)
		}
	default:
		err = fmt.Errorf("unknown ManyStatement passed [%v]", query)
	}
	return
}

// Needs is used in part of the validate check of the named parameters and returns
// the field names the ManyStatment passed in should have
// Uses a regex to find words starting with :
func Needs(query ManyStatement) (needs []string) {
	var namedParamPattern string = `(?m)(:[\w-]+)`
	var prefix string = ":"
	var re = regexp.MustCompile(namedParamPattern)
	for _, match := range re.FindAllString(string(query), -1) {
		needs = append(needs, strings.TrimPrefix(match, prefix))
	}
	return
}

// ValidateParameters checks if the parameters passed meets all the required
// needs for the query being run
func ValidateParameters(params *Parameters, needs []string) (err error) {
	mapped := params.Map()
	if len(mapped) == 0 {
		err = fmt.Errorf("parameters passed must contain at least one field")
		return
	}

	missing := []string{}
	// check each need if that exists as a key in the map
	for _, need := range needs {
		if _, ok := mapped[need]; !ok {
			missing = append(missing, need)
		}
	}
	// if any field is missing then set error
	if len(missing) > 0 {
		cols := strings.Join(missing, ",")
		err = fmt.Errorf("missing required fields for this query: [%s]", cols)
	}

	return
}
