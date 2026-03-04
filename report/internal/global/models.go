package global

type ImportArgs struct {
	DB     string `json:"db"`     // DB related (--db)
	Driver string `json:"driver"` // DB related (--driver)
	Params string `json:"params"` // DB related (--params)
	// used over many import
	DateStart string `json:"date_start"` // Date ranges (--start)
	DateEnd   string `json:"date_end"`   // Date ranges (--end)
	// github related
	OrgSlug    string `json:"org"`    // github org (--org)
	ParentSlug string `json:"parent"` // github parent team (--parent)
	// file based imports for teams / accounts
	SrcFile string `json:"src-file"` // File based import (--src-file)
	// aws region
	Region string `json:"region"` // AWS related (--region)
	// aws costs
	DateStartCosts string `json:"date_start_costs"` // Date ranges (--start-costs)
	// general filter
	Filter string `json:"filter"` // generic filter option (--filter)
	// code owner command only
	Owners  bool `json:"owners"`  // --owners
	Stats   bool `json:"stats"`   // --stats
	Metrics bool `json:"metrics"` // --metrics
}
