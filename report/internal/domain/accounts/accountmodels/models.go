package accountmodels

// Account struct for db importing / seeding
type Account struct {
	ID          string `json:"id,omitempty" db:"id"`                   // This is the Account ID as a string - they can have leading 0
	Name        string `json:"name,omitempty" db:"name"`               // account name as used internally
	Label       string `json:"label,omitempty" db:"label"`             // internal label
	Environment string `json:"environment,omitempty" db:"environment"` // environment type
	TeamName    string `json:"billing_unit,omitempty" db:"team_name"`  // team associated with the account
}

// AccountRow is used in the api calls so the json name of TeamName can be corrected
type AccountRow struct {
	ID          string `json:"id" db:"id"`                   // This is the Account ID as a string - they can have leading 0
	Name        string `json:"name" db:"name"`               // account name as used internally
	Label       string `json:"label" db:"label"`             // internal label
	Environment string `json:"environment" db:"environment"` // environment type
	TeamName    string `json:"team_name" db:"team_name"`     // team associated with the account
}
