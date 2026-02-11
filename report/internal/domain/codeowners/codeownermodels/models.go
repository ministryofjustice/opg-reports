package codeownermodels

// Codeowner struct is used for importing / seeding flat database records
type Codeowner struct {
	ID               int    `json:"id,omitempty" db:"id"`
	Name             string `json:"name,omitempty" db:"name"`                             // Name of the codeowner
	CodebaseFullName string `json:"codebase_full_name,omitempty" db:"codebase_full_name"` // This is the full name of the repo, used in joins
	TeamName         string `json:"team_name,omitempty" db:"team_name"`                   // This is the associated team name
}

// CodeownerAll is used in the codeowner all endpoint to return all codeowners and the url of the codebase
type CodeownerAll struct {
	Name             string `json:"name" db:"name"`                             // Name of the codeowner
	CodebaseFullName string `json:"codebase_full_name" db:"codebase_full_name"` // This is the full name of the repo
	TeamName         string `json:"team_name" db:"team_name"`                   // This is the associated team name
	Url              string `json:"url" db:"url"`                               // This is the url of the associate codebase via the left join
}

// CodeownerForTeam is used in the codeowner all endpoint to return all codeowners and the url of the codebase
type CodeownerForTeam struct {
	Name             string `json:"name" db:"name"`                             // Name of the codeowner
	CodebaseFullName string `json:"codebase_full_name" db:"codebase_full_name"` // This is the full name of the repo
	TeamName         string `json:"team_name" db:"team_name"`                   // This is the associated team name
	Url              string `json:"url" db:"url"`                               // This is the url of the associate codebase via the left join
}
