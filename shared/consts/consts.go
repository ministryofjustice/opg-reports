package consts

import (
	"time"
)

const (
	API_SCHEME  string        = "http"          // API default scheme is http
	API_ADDR    string        = ":8081"         // API default address is :8081
	API_TIMEOUT time.Duration = time.Second * 4 // Default timeout to use for API requests is 4 seconds
)

// BILLING_DATE is the day of the month billing data is accurate from
//
// Example:
// On 2024/05/14 the billing date should be 2024/03/31 as information for
// the month of 2024/04 is not yet ready.
const BILLING_DATE int = 15

// SQL_CONNECTION_PARAMS are our default connections used for performance etc
const SQL_CONNECTION_PARAMS string = "?_journal=WAL&_busy_timeout=5000&_vacuum=incremental&_synchronous=NORMAL&_cache_size=1000000000"
