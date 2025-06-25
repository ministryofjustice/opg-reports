package team

// Team replaces the Unit structure as a grouping of services / accounts
type Team struct {
	ID        int    `json:"id" db:"id"`
	CreatedAt string `json:"created_at" db:"created_at" example:"2019-08-24T14:15:22Z"`
	Name      string `json:"name" db:"name" example:"Sirius"`
}
