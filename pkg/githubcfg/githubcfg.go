package githubcfg

import "github.com/ministryofjustice/opg-reports/pkg/envar"

// Config contains env vars
type Config struct {
	Token string
}

// FromEnv creates a config struct directly from environment
// variables - with default value for region of eu-west-1
func FromEnv() *Config {

	return &Config{
		Token: envar.Get("GITHUB_ACCESS_TOKEN", ""),
	}
}
