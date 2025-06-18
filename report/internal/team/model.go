package owner

type Owner struct {
	ID        int    `json:"id" db:"id"`
	CreatedAt string `json:"created_at" db:"created_at"`
	Name      string `json:"name" db:"name"`
}
