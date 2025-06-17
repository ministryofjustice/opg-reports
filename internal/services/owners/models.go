package owners

type Owner struct {
	// Direct database fields
	ID        int    `json:"id,omitempty"`         // Database id
	CreatedAt string `json:"created_at,omitempty"` // Timestamp data entry is created
	// Fields
	Name string `json:"name,omitempty"`
}
