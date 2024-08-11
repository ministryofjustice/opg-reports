package config

import (
	"github.com/ministryofjustice/opg-reports/servers/front/config/navigation"
	"github.com/ministryofjustice/opg-reports/shared/convert"
	"github.com/ministryofjustice/opg-reports/shared/must"
)

type Config struct {
	Organisation string                       `json:"organisation"`
	Navigation   []*navigation.NavigationItem `json:"navigation"`
}

func NewConfig(content []byte) (*Config, error) {
	c := &Config{}
	return convert.Unmarshal[*Config](content, c)
}

func New(content []byte) *Config {
	return must.Must[*Config](NewConfig(content))

}
