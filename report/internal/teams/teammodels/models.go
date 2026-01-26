package teammodels

type Team struct {
	Name string `json:"name,omitempty" db:"name" example:"SREs"`
}
