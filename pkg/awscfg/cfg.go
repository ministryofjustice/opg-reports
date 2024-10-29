package awscfg

import "github.com/ministryofjustice/opg-reports/pkg/envar"

// Config contains env vars
type Config struct {
	Region          string
	AccessKeyID     string
	SecretAccessKey string
	SessionToken    string
}

func FromEnv() *Config {

	return &Config{
		Region:          envar.Get("AWS_DEFAULT_REGION", "eu-west-1"),
		AccessKeyID:     envar.Get("AWS_ACCESS_KEY_ID", ""),
		SecretAccessKey: envar.Get("AWS_SECRET_ACCESS_KEY", ""),
		SessionToken:    envar.Get("AWS_SESSION_TOKEN", ""),
	}
}
