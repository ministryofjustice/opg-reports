package args

type Args struct {
	DB     string `json:"db"`     // database path
	Driver string `json:"driver"` // database driver
	Params string `json:"params"` // database connection params

	OrgSlug    string `json:"org_slug"`    // github org name
	ParentSlug string `json:"parent_slug"` // parent slug

	IncludeStats      bool `json:"include_stats"`      // run the code base stats handler - stats are non-time boxed details
	IncludeCodeowners bool `json:"include_codeowners"` // option to fetch all codebases and then fetch codeowner data as well

	FilterByName string `json:"filter_by_name"` // used to limit the repos to those that exactly match this name

}
