package frontmodels

type RegisterArgs struct {
	ApiHost      string `json:"api"`
	GovUKVersion string `json:"govuk_version"`
	SemVer       string `json:"semver"`
	RootDir      string `json:"root_dir"`
	TemplateDir  string `json:"template_dir"`
}

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
	// Releases
	Releases            int `json:"releases"`             // count releases
	ReleasesSecurityish int `json:"releases_securityish"` // count of releases if they likely sec
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

// Codebase - contains simple and stat fields used on the front end
type Codebase struct {
	Name       string `json:"name,omitempty"`        // short name of codebase (without owner)
	FullName   string `json:"full_name,omitempty" `  // full name including the owner
	Url        string `json:"url,omitempty" `        // url to access the codebase
	Visibility string `json:"visibility,omityempty"` // visibility status

	ComplianceLevel     string `json:"compliance_level,omitempty"`      // compliance level (moj based)
	ComplianceReportUrl string `json:"compliance_report_url,omitempty"` // compliance report url
	ComplianceBadge     string `json:"compliance_badge,omitempty"`      // compliance badge url
	ComplianceGrade     int    `json:"compliance_grade,omitempty"`

	TrivyUsage     int    `json:"trivy_usage"`      // boolean flag to show if the codebase is using trivy in workflows
	TrivySBOMUsage int    `json:"trivy_sbom_usage"` // boolean flag to show if trivy is being used to generate sboms
	TrivyLocations string `json:"trivy_locations"`  // tracks files where trivy has been utilised

}

// Codeowner
type Codeowner struct {
	FullName string   `json:"full_name,omitempty" ` // full name including the owner
	Name     string   `json:"name,omitempty"`       // short name of codebase (without owner)
	Url      string   `json:"url,omitempty" `       // url to access the codebase
	Owners   []*Owner `json:"owners"`               // list of codeowners
}

// joined teams is the codebase -> codeowners data
type Owner struct {
	Owner    string `json:"owner"`
	TeamName string `json:"team_name"`
}

// CodebaseData
type CodebaseData struct {
	Team       string
	Codebases  []*Codebase
	CodeOwners []*Codeowner
}

type ReleaseData struct {
	Team     string
	Releases []*Release
	Summary  *Release
}
type Release struct {
	Month               string `json:"month"`                // month as YYYY-MM string
	Releases            int    `json:"releases"`             // count of releases for this month
	ReleasesSecurityish int    `json:"releases_securityish"` // count of releases for this month that seem to be security related
}
