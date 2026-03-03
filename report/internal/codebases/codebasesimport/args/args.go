package args

import "time"

type Args struct {
	DB     string `json:"db"`     // database path
	Driver string `json:"driver"` // database driver
	Params string `json:"params"` // database connection params

	DateStart time.Time `json:"date_start"` // start date
	DateEnd   time.Time `json:"date_end"`   // end date

	OrgSlug    string `json:"org_slug"`    // github org name
	ParentSlug string `json:"parent_slug"` // parent slug

	IncludeCodeowners bool `json:"include_codeowners"` // option to fetch all codebases and then fetch codeowner data as well
	IncludeStats      bool `json:"include_stats"`      // run the code base stats handler - stats are non-time boxed details
	IncludeMetrics    bool `json:"include_metrics"`    // metrics are time based stats

	FilterByName string `json:"filter_by_name"` // used to limit the repos to those that exactly match this name

}
