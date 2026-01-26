package cemodels

// AwsCost represents a simple, joinless, db row in the cost table
//
// Used by imports and seeding commands
type AwsCost struct {
	ID        int    `json:"id,omitempty" db:"id"`
	Region    string `json:"region,omitempty" db:"region" example:"eu-west-1"`
	Service   string `json:"service,omitempty" db:"service"` // The AWS service name
	Date      string `json:"date,omitempty" db:"date" `      // The data the cost was incurred - provided from the cost explorer result
	Cost      string `json:"cost,omitempty" db:"cost" `      // The actual cost value as a string - without an currency, but is USD by default
	AccountID string `json:"account_id" db:"aws_account_id"` // the actual aws account id - string as it can have leading zeros. Use in joins as well
}
