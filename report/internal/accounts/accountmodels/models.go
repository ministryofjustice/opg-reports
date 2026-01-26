package accountmodels

// AwsAccount struct for db importing / seeding
type AwsAccountImport struct {
	ID          string `json:"id,omitempty" db:"id"`                   // This is the AWS Account ID as a string - they can have leading 0
	Name        string `json:"name,omitempty" db:"name"`               // account name as used internally
	Label       string `json:"label,omitempty" db:"label"`             // internal label
	Environment string `json:"environment,omitempty" db:"environment"` // environment type
	TeamName    string `json:"billing_unit,omitempty" db:"team_name"`  // team associated with the account
}
