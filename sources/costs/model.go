package costs

import (
	"fmt"
	"strconv"

	"github.com/ministryofjustice/opg-reports/pkg/record"
)

// Cost is both the database model and the result struct for the api
type Cost struct {
	ID           int    `json:"id,omitempty" db:"id" faker:"-" doc:"Database primary key."`                                                                              // ID is a generated primary key
	Ts           string `json:"ts,omitempty" db:"ts"  faker:"time_string" doc:"Time the record was created."`                                                            // TS is timestamp when the record was created
	Organisation string `json:"organisation,omitempty" db:"organisation" faker:"oneof: foo, bar, foobar" doc:"Name of the organisation."`                                // Organisation is part of the account details and string name
	AccountID    string `json:"account_id,omitempty" db:"account_id" faker:"oneof: 101, 102, 201, 202, 301, 302" doc:"Account ID this cost comes from."`                 // AccountID is the aws account id this row is for
	AccountName  string `json:"account_name,omitempty" db:"account_name" faker:"word" doc:"A simple name for the account this cost came from."`                          // AccountName is a passed string used to represent the account purpose
	Unit         string `json:"unit,omitempty" db:"unit" faker:"oneof: unitA, unitB, unitC" doc:"The name of the unit / team that owns this account."`                   // Unit is the team that owns this account, passed directly
	Label        string `json:"label,omitempty" db:"label" faker:"word" doc:"A supplimental lavel to provide extra detail on the account type."`                         // Label is passed string that sets a more exact name - so DB account production
	Environment  string `json:"environment,omitempty" db:"environment" faker:"oneof: production, pre-production, development" doc:"Environment type."`                   // Environment is passed along to show if this is production, development etc account
	Region       string `json:"region,omitempty" db:"region" faker:"oneof: NoRegion, eu-west-1, eu-west-2, us-east-2" doc:"Region this cost was generated within."`      // From the cost data, this is the region the service cost aws generated in
	Service      string `json:"service,omitempty" db:"service" faker:"oneof: Tax, ecs, ec2, s3, sqs, waf, ses, rds" doc:"Name of the service that generated this cost."` // The AWS service name
	Date         string `json:"date,omitempty" db:"date" faker:"date_string" doc:"Date this cost was generated."`                                                        // The data the cost was incurred - provided from the cost explorer result
	Cost         string `json:"cost,omitempty" db:"cost" faker:"float_string" doc:"Cost value."`                                                                         // The actual cost value as a string - without an currency, but is USD by default
}

// New
// Record interface
func (self *Cost) New() record.Record {
	return &Cost{}
}

// UID
// Record interface
func (self *Cost) UID() string {
	return fmt.Sprintf("%s-%d", "costs", self.ID)
}

// SetID
// Record interface
func (self *Cost) SetID(id int) {
	self.ID = id
}

// TDate
// Transformable interface
func (self *Cost) TDate() string {
	return self.Date
}

// TValue
// Transformable interface
func (self *Cost) TValue() string {
	return self.Cost
}

// Value handles converting the string value of Cost into a float64
func (self *Cost) Value() (cost float64) {
	if floated, err := strconv.ParseFloat(self.Cost, 10); err == nil {
		cost = floated
	}
	return
}
