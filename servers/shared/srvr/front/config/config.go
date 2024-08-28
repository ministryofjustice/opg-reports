package config

import (
	"github.com/ministryofjustice/opg-reports/servers/shared/srvr/front/config/nav"
	"github.com/ministryofjustice/opg-reports/shared/convert"
)

type Config struct {
	Organisation string
	Navigation   []*nav.Nav
}

func New(content []byte) (cfg *Config) {
	var err error
	cfg, err = convert.Unmarshal[*Config](content)
	if err != nil {
		cfg = &Config{}
	}
	return
}
