package frontmodels

// HeadlineData
type HeadlineData struct {
	Team      string
	DateStart string `json:"date_start"`
	DateEnd   string `json:"date_end"`
	// costs
	TotalCost           float64 `json:"total_cost"`             // total cost result
	AverageCostPerMonth float64 `json:"average_cost_per_month"` // average cost per month
	// uptime
	OverallUptime float64 `json:"overall_uptime"` // uptime
	// codebases
	CodebaseCount  int     `json:"codebase_count"`  // total number of codebases
	CodebasePassed float64 `json:"codebase_passed"` // % that have a passing status
}

// TableHeaders
type TableHeaders struct {
	Labels []string `json:"labels"`
	Data   []string `json:"data"`
	Extra  []string `json:"extra"`
	End    []string `json:"end"`
}

// TableData is used to handle the cost table data construct
type TableData struct {
	Headers    *TableHeaders            `json:"headers"` // headers contains details for table headers / rendering
	Data       []map[string]interface{} `json:"data"`    // the actual data results
	Summary    map[string]interface{}   `json:"summary"` // used to contain table totals etc
	BillingDay int                      `json:"billing_day"`
}

// DateRanges
type DateRanges struct {
	Months    []string
	DateStart string
	DateEnd   string
}

// DateComparision
type DateComparision struct {
	Months  []string
	Changes []string
	Change  string
	DateA   string
	DateB   string
}

// Codebase
type Codebase struct {
	Name                string `json:"name,omitempty"`                  // short name of codebase (without owner)
	FullName            string `json:"full_name,omitempty" `            // full name including the owner
	Url                 string `json:"url,omitempty" `                  // url to access the codebase
	Visibility          string `json:"visibility,omityempty"`           // visibility status
	ComplianceLevel     string `json:"compliance_level,omitempty"`      // compliance level (moj based)
	ComplianceReportUrl string `json:"compliance_report_url,omitempty"` // compliance report url
	ComplianceBadge     string `json:"compliance_badge,omitempty"`      // compliance badge url
}

// CodebaseData
type CodebaseData struct {
	Codebases []*Codebase
}
