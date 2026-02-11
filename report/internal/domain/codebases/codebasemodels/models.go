package codebasemodels

import (
	"fmt"
	"opg-reports/report/internal/utils"
)

// Codebase is used for importing the flat db records
type Codebase struct {
	ID       int    `json:"id,omitempty" db:"id"`               // db key
	Name     string `json:"name,omitempty" db:"name"`           // short name of codebase (without owner)
	FullName string `json:"full_name,omitempty" db:"full_name"` // full name including the owner
	Url      string `json:"url,omitempty" db:"url"`             // url to access the codebase
}

// CodebaseAll is used in codebaseall api to return list of codebase data and the attached codeowners
type CodebaseAll struct {
	Name       string            `json:"name" db:"name"`                           // short name of codebase (without owner)
	FullName   string            `json:"full_name" db:"full_name"`                 // full name including the owner
	Url        string            `json:"url" db:"url"`                             // url to access the codebase
	Codeowners hasManyCodeowners `json:"codeowners,omitempty" db:"codeowner_list"` // list from the join
}

// joined teams is the codebase -> codeowners data
type joinedCodeowner struct {
	Name     string `json:"name" db:"name"`
	TeamName string `json:"team_name" db:"team_name"`
}

// hasManyCodeowners represents the join
// Interfaces:
//   - sql.Scanner
type hasManyCodeowners []*joinedCodeowner

// Scan handles the processing of the join data
func (self *hasManyCodeowners) Scan(src interface{}) (err error) {
	switch src.(type) {
	case []byte:
		err = utils.Unmarshal(src.([]byte), self)
	case string:
		err = utils.Unmarshal([]byte(src.(string)), self)
	default:
		err = fmt.Errorf("unsupported scan src type")
	}
	return
}
