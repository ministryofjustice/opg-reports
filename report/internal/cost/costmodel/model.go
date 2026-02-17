package costmodel

// Import represents a simple, joinless, db row in the cost table; used by imports and seeding commands
type Import struct {
	Region    string `json:"region,omitempty"`  // AWS Region
	Service   string `json:"service,omitempty"` // The AWS service name
	Month     string `json:"month,omitempty"`   // The data the cost was incurred - provided from the cost explorer result
	Cost      string `json:"cost,omitempty"`    // The actual cost value as a string - without an currency, but is USD by default
	AccountID string `json:"account_id"`        // the actual account id - string as it can have leading zeros. Use in joins as well
}
