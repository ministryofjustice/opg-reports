package teams

import "github.com/ministryofjustice/opg-reports/report/internal/team"

// Team is a localised version of team.Team without fields removed for api
// sanitisation
type Team struct {
	team.Team
	CreatedAt string `json:"-"` // its not in the select, but blank the field incase
}

// GetTeamsAllResponse is response object used by the handler
type GetTeamsAllResponse struct {
	Body struct {
		Data []*Team
	}
}
