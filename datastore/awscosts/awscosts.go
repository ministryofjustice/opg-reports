// Package awscosts provides struct and database methods for handling cost explorer data
// that is then used by the api
package awscosts

import (
	"strconv"
)

// Cost is used to store the data from cost explorer into the database
type Cost struct {
	ID int    `json:"id,omitempty" db:"id"` // ID is a generated primary key
	Ts string `json:"ts,omitempty" db:"ts"` // TS is timestamp when the record was created

	Organisation string `json:"organisation,omitempty" db:"organisation"` // Organisation is part of the account details and string name
	AccountID    string `json:"account_id,omitempty" db:"account_id"`     // AccountID is the aws account id this row is for
	AccountName  string `json:"account_name,omitempty" db:"account_name"` // AccountName is a passed string used to represent the account purpose
	Unit         string `json:"unit,omitempty" db:"unit"`                 // Unit is the team that owns this account, passed directly
	Label        string `json:"label,omitempty" db:"label"`               // Label is passed string that sets a more exact name - so DB account production
	Environment  string `json:"environment,omitempty" db:"environment"`   // Environment is passed along to show if this is production, development etc account

	Region  string `json:"region,omitempty" db:"region"`   // From the cost data, this is the region the service cost aws generated in
	Service string `json:"service,omitempty" db:"service"` // The AWS service name
	Date    string `json:"date,omitempty" db:"date"`       // The data the cost was incurred - provided from the cost explorer result
	Cost    string `json:"cost,omitempty" db:"cost"`       // The actual cost value as a string - without an currency, but is USD by default
}

// Value handles converting the string value of Cost into a float64
func (self *Cost) Value() (cost float64) {
	if floated, err := strconv.ParseFloat(self.Cost, 10); err == nil {
		cost = floated
	}
	return
}
