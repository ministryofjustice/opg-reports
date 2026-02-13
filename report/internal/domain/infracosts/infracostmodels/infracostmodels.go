package infracostmodels

// Cost represents a simple, joinless, db row in the cost table
//
// Used by imports and seeding commands
type Cost struct {
	ID        int    `json:"id,omitempty" db:"id"`
	Region    string `json:"region,omitempty" db:"region" example:"eu-west-1"`
	Service   string `json:"service,omitempty" db:"service"` // The AWS service name
	Date      string `json:"date,omitempty" db:"date" `      // The data the cost was incurred - provided from the cost explorer result
	Cost      string `json:"cost,omitempty" db:"cost" `      // The actual cost value as a string - without an currency, but is USD by default
	AccountID string `json:"account_id" db:"account_id"`     // the actual account id - string as it can have leading zeros. Use in joins as well
}

type CostData struct {
	Date        string  `json:"date,omitempty" db:"date" `                // The month for this cost
	Cost        float64 `json:"cost,omitempty" db:"cost" `                // The sum of all costs within this month
	TeamName    string  `json:"team,omitempty" db:"team"`                 // the team this costs is attached to
	AccountID   string  `json:"account_id,omitempty" db:"account_id"`     // the account this costs is attached to
	AccountName string  `json:"account_name,omitempty" db:"account_name"` // the account this costs is attached to
	Environment string  `json:"environment,omitempty" db:"environment"`   // the env used
	Service     string  `json:"service,omitempty" db:"service" `          // The aws service name

}
