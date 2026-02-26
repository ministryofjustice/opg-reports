package global

type ImportArgs struct {
	DB         string `json:"db"`         // DB related (--db)
	Driver     string `json:"driver"`     // DB related (--driver)
	Params     string `json:"params"`     // DB related (--params)
	Region     string `json:"region"`     // AWS related (--region)
	DateStart  string `json:"date_start"` // Date ranges (--start)
	DateEnd    string `json:"date_end"`   // Date ranges (--end)
	SrcFile    string `json:"src-file"`   // File based import (--src-file)
	OrgSlug    string `json:"org"`        // github org (--org)
	ParentSlug string `json:"parent"`     // github parent team (--parent)
}
