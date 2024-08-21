package consts

import (
	"time"
)

const BILLING_DATE int = 15

const API_SCHEME string = "http"
const API_ADDR string = ":8081"
const API_TIMEOUT time.Duration = time.Second * 3

// -- sql connections
const SQL_CONNECTION_PARAMS string = "?_journal=WAL&_busy_timeout=5000&_vacuum=incremental&_synchronous=NORMAL&_cache_size=1000000000"
