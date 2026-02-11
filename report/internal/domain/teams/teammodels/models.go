package teammodels

type Team struct {
	Name string `json:"name" db:"name" example:"SREs"`
}
