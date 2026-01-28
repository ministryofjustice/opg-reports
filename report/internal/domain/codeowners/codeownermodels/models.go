package codeownermodels

type Codeowner struct {
	ID               int    `json:"id,omitempty" db:"id"`                                 // This is the AWS Account ID as a string
	Name             string `json:"name,omitempty" db:"name"`                             // Name of the codeowner
	CodebaseFullName string `json:"codebase_full_name,omitempty" db:"codebase_full_name"` // This is the full name of the repo, used in joins
	TeamName         string `json:"team_name,omitempty" db:"team_name"`                   // This is the associated team name, used in joins - can be empty
}
