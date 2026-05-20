package ghconfig

// Config is used to capture information needed for the api call
type Config struct {
	OrganisationSlug string // the slug of the organisation
	TeamSlug         string // slug of the team whose repositories we're returning
}

// New returns simple Config value with org & team defaults
func New(org string, team string) *Config {
	return &Config{
		OrganisationSlug: org,
		TeamSlug:         team,
	}
}
