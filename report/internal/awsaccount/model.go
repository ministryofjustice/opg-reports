package awsaccount

type AwsAccount struct {
	ID          string `json:"id" db:"id"` // This is the AWS Account ID
	CreatedAt   string `json:"created_at" db:"created_at" example:"2019-08-24T14:15:22Z"`
	Name        string `json:"name,omitempty" db:"name"`
	Label       string `json:"label,omitempty" db:"label"`
	Environment string `json:"environment,omitempty" db:"environment"`
}
