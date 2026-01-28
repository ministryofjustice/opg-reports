package codebasemodels

type Codebase struct {
	ID       int    `json:"id,omitempty" db:"id"`               // db key
	Name     string `json:"name,omitempty" db:"name"`           // short name of codebase (without owner)
	FullName string `json:"full_name,omitempty" db:"full_name"` // full name including the owner
	Url      string `json:"url,omitempty" db:"url"`             // url to access the codebase
}
