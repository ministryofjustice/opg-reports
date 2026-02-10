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
	ID       int          `json:"id,omitempty" db:"id"`               // db key
	Name     string       `json:"name,omitempty" db:"name"`           // short name of codebase (without owner)
	FullName string       `json:"full_name,omitempty" db:"full_name"` // full name including the owner
	Url      string       `json:"url,omitempty" db:"url"`             // url to access the codebase
	Teams    hasManyTeams `json:"teams" db:"team_list"`
}

// joined teams is the codebase -> codeowners.team_name data
type joinedTeam struct {
	Name string `json:"name,omitempty" db:"name"`
}

// hasManyTeams represents the join
// Interfaces:
//   - sql.Scanner
type hasManyTeams []*joinedTeam

// Scan handles the processing of the join data
func (self *hasManyTeams) Scan(src interface{}) (err error) {
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
