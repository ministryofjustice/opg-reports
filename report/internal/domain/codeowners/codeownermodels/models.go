package codeownermodels

// Codeowner struct is used for importing / seeding flat database records
type Codeowner struct {
	ID               int    `json:"id,omitempty" db:"id"`
	Name             string `json:"name,omitempty" db:"name"`                             // Name of the codeowner
	CodebaseFullName string `json:"codebase_full_name,omitempty" db:"codebase_full_name"` // This is the full name of the repo, used in joins
	TeamName         string `json:"team_name,omitempty" db:"team_name"`                   // This is the associated team name
}

// CodeownerData
type CodeownerData struct {
	Codeowner string `json:"codeowner" db:"codeowner"` // Name of the codeowner
	Codebase  string `json:"codebase" db:"codebase"`   // This is the full name of the repo
	TeamName  string `json:"team" db:"team"`           // This is the associated team name
	Url       string `json:"url" db:"url"`             // This is the url of the associate codebase via the left join
}
