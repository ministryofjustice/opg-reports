package awscosts

// Cost is used to store the data from cost explorer into the database
type Cost struct {
	ID int    `json:"id" db:"id"` // ID is a generated primary key
	Ts string `json:"ts" db:"ts"` // TS is timestamp when the record was created

	Organisation string `json:"organisation" db:"organisation"` // Organisation is part of the account details and string name
	AccountID    string `json:"account_id" db:"account_id"`     // AccountID is the aws account id this row is for
	AccountName  string `json:"account_name" db:"account_name"` // AccountName is a passed string used to represent the account purpose
	Unit         string `json:"unit" db:"unit"`                 // Unit is the team that owns this account, passed directly
	Label        string `json:"label" db:"label"`               // Label is passed string that sets a more exact name - so DB account production
	Environment  string `json:"environment" db:"environment"`   // Environment is passed along to show if this is production, development etc account

	Region  string `json:"region" db:"region"`   // From the cost data, this is the region the service cost aws generated in
	Service string `json:"service" db:"service"` // The AWS service name
	Date    string `json:"date" db:"date"`       // The data the cost was incurred - provided from the cost explorer result
	Cost    string `json:"cost" db:"cost"`       // The actual cost value as a string - without an currency, but is USD by default
}
