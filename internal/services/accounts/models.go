package accounts

type AwsAccount struct {
	// Direct database fields
	CreatedAt string `json:"created_at,omitempty"` // Timestamp data entry is created
	// Fields
	ID          string `json:"id,omitempty"` // this is the AWS account ID and used as DB primary key
	Name        string `json:"name,omitempty"`
	Label       string `json:"label,omitempty"`
	Environment string `json:"environment,omitempty"`
	// Joins
	OwnerID int `json:"owner_id,omitempty"`
}
