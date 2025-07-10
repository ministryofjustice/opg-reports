package endpoints

const API_VERSION string = "v1"

const HOME string = "/"

// All `teams` related endpoints
const (
	teams         string = "/" + API_VERSION + "/teams/"
	TEAMS_GET_ALL string = teams + "all"
)

// All `awsaccounts` related endpoints
const (
	awsaccounts         string = "/" + API_VERSION + "/awsaccounts/"
	AWSACCOUNTS_GET_ALL string = awsaccounts + "all"
)

// All `awscost` related endpoints
const (
	awscosts            string = "/" + API_VERSION + "/awscosts/"
	AWSCOSTS_GET_TOP_20 string = awscosts + "top20"
	AWSCOSTS_GROUPED    string = awscosts + "grouped/{granularity}/{start_date}/{end_date}"
)
