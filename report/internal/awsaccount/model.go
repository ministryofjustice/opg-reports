package awsaccount

type AwsAccount struct {
	ID          string `json:"id" db:"id"` // This is the AWS Account ID as a string
	CreatedAt   string `json:"created_at" db:"created_at" example:"2019-08-24T14:15:22Z"`
	Name        string `json:"name,omitempty" db:"name"`
	Label       string `json:"label,omitempty" db:"label"`
	Environment string `json:"environment,omitempty" db:"environment"`

	// join to teams
	TeamID int `json:"team_id,omitempty" db:"team_id"`
}

// AwsAccountImport captures an extra field from the metadata which
// is used in the stmtInsert to resolve join to team
type AwsAccountImport struct {
	AwsAccount
	BillingUnit string `json:"billing_unit,omitempty" db:"billing_unit"`
}
