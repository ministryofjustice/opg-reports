// Package consts contains a series of const variables for use in the project.
//
// Capture things like DateTime formats and key values
package consts

import "time"

const ApiTimeout time.Duration = time.Second * 4

const (
	DateFormat             string = time.RFC3339
	DateFormatYear         string = "2006"
	DateFormatYearMonth    string = "2006-01"
	DateFormatYearMonthDay string = "2006-01-02"
	DateYear               string = "year"
	DateMonth              string = "month"
	DateDay                string = "day"
)

// CostsBillingDay is the day in the month where updated billing data has been
// fetched and we can use the month prior to the current.
//
// Example - 10th April would be looking at Feb costs, 15th April can see March costs
const CostsBillingDay int = 15