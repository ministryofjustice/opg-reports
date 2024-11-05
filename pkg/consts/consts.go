// Package consts contains a series of const variables for use in the project.
//
// Capture things like DateTime formats and key values
package consts

import "time"

// FetchTimeout is the duration used when calling the api endpoints
const FetchTimeout time.Duration = time.Second * 4

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

// GovUKFrontendVersion is the current release being used for css etc
const GovUKFrontendVersion string = "5.7.1"

const (
	ServerDefaultFrontAddr string = "localhost:8080"
	ServerDefaultApiAddr   string = "localhost:8081"
)

const (
	DefaultFloatString string = "0.0000"
)
