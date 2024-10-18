package awscosts

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

// Parameters are used to as named parameters on sqlx queries
// via the Query function and cover all possible
type Parameters struct {
	StartDate  string `json:"start_date,omitempty" db:"start_date"`    // StartDate is the lower bound of date based query
	EndDate    string `json:"end_date,omitempty" db:"end_date"`        // EndDate is the upper bound of date based query
	DateFormat string `json:"date_format,,omitempty" db:"date_format"` // Date format to use for strftime with query
}

// Statment is a string, used as a enum style type
// for the various sql queries we want to run
// that allows named parameters etc
type Statement string

// TotalsWithAndWithoutTax is used to calculate the total costs
// within the given date range (>= :start_date, < :end_date) and
// splits that based on the `service` being 'Tax' or not
// Used to show top line numbers without VAT etc
const TotalsWithAndWithoutTax Statement = `
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
    excTax.service != 'Tax'
    AND excTax.date >= :start_date
    AND excTax.date < :end_date
GROUP BY strftime(:date_format, date)
ORDER by date ASC;
`

// Query runs the known statement against using the parameters as named values within them and returns the
// result as a slice of []*Cost
func Query(ctx context.Context, db *sqlx.DB, query Statement, params *Parameters) (results []*Cost, err error) {
	var statement *sqlx.NamedStmt
	results = []*Cost{}

	switch query {
	case TotalsWithAndWithoutTax:
		if statement, err = db.PrepareNamedContext(ctx, string(query)); err == nil {
			err = statement.SelectContext(ctx, &results, params)
		}
	default:
		err = fmt.Errorf("unknown statement passed [%v]", query)
	}

	return
}
