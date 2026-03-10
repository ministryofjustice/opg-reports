package types

import (
	"opg-reports/report/packages/convert"
	"opg-reports/report/packages/types/interfaces"
	"opg-reports/report/packages/types/models"
)

// Team is used in the importing & api handler
//
// interfaces.Insertable
// interfaces.Selectable
type Team struct {
	Name string `json:"name" db:"name,omitempty"`
}

// Map returns a map of all fields on this struct
func (self *Team) Map() (m map[string]interface{}) {
	m = map[string]interface{}{}
	convert.Between(self, &m)
	return
}

// Sequence is used to return the columns in the order they are selected.
func (self *Team) Sequence() []any {
	return []any{
		&self.Name,
	}
}

// Select is the select statement used for pull data from the db.
//
// interfaces.Statement
type Select struct {
	Statement string
}

// SQL updates the sql statement to adjust for team filtering
func (self *Select) SQL(filter interfaces.Filterable) (s string) {
	var model = &models.Filter{}
	convert.Between(filter.Map(), &model)
	s = self.Statement
	return
}
