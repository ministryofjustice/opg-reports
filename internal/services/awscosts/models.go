package awscosts

// Cost represents a single line item from our cost records
// with a join to the account that it came from
type Cost struct {
	// Direct database fields
	ID        int    `json:"id,omitempty"`         // Database id
	CreatedAt string `json:"created_at,omitempty"` // Timestamp data entry is created
	// Fields specific to costs
	Region  string `json:"region,omitempty"`  // AWS region name
	Service string `json:"service,omitempty"` // AWS service name
	Date    string `json:"date,omitempty"`    // YYYY-MM-DD formatted string for when the cost was incurred
	Cost    string `json:"cost,omitempty"`    // Raw cost value without any currency - should be USD
	// Joins
	AccountID string `json:"account_id,omitempty"` // AWS Account ID used on the account table
}
