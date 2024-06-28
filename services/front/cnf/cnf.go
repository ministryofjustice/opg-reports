package cnf

import "encoding/json"

type RepoStandards struct {
	Baseline    []string `json:"baseline"`
	Extended    []string `json:"extended"`
	Information []string `json:"information"`
}
type StandardsCnf struct {
	Repository RepoStandards `json:"repository"`
}

type SiteSection struct {
	Name     string         `json:"name"`
	Href     string         `json:"href"`
	Header   bool           `json:"header"`
	Sections []*SiteSection `json:"sections"`

	Exclude bool `json:"exclude"`

	Api             map[string]string `json:"api"`
	ResponseHandler string            `json:"handler"`
	TemplateName    string            `json:"template"`

	Registered bool `json:"-"`
}
type Config struct {
	Organisation string         `json:"organisation"`
	Sections     []*SiteSection `json:"sections"`
	// Standards is used ot allow customisable fields for baseline repo standards
	Standards *StandardsCnf `json:"standards"`
}

func Load(content []byte) (*Config, error) {
	cfg := &Config{}
	err := json.Unmarshal(content, &cfg)
	return cfg, err
}

func FlatSections(sects []*SiteSection, flat map[string]*SiteSection) {

	for _, sect := range sects {
		flat[sect.Href] = sect
		if len(sect.Sections) > 0 {
			FlatSections(sect.Sections, flat)
		}
	}
	return
}
