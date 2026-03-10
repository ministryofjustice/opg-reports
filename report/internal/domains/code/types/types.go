package types

import (
	"opg-reports/report/packages/convert"
	"opg-reports/report/packages/types/interfaces"
)

// Codebase is used in the importing of data from github api
//
// interfaces.Insertable
// interfaces.Selectable
type Codebase struct {
	FullName string `json:"full_name,omitempty" ` // full name including the owner - used as unique
	Name     string `json:"name,omitempty"`       // short name of codebase (without owner)
	Url      string `json:"url,omitempty" `       // url to access the codebase
	Archived int    `json:"archived"`             // int version of the archived flag on the repo
}

// Map returns a map of all fields on this struct
func (self *Codebase) Map() (m map[string]interface{}) {
	m = map[string]interface{}{}
	convert.Between(self, &m)
	return
}

// Sequence is used to return the columns in the order they are selected.
func (self *Codebase) Sequence() []any {
	return []any{
		&self.FullName,
		&self.Name,
		&self.Url,
		&self.Archived,
	}
}

// Select is the select statement used for pull data from the db.
//
// interfaces.Statement
type Select struct {
	Statement string
}

// SQL updates the sql statement to adjust for filtering
func (self *Select) SQL(filter interfaces.Filterable) (s string) {
	s = self.Statement
	return
}
