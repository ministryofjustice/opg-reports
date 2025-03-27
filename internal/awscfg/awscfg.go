// Package awscfg provides config methods for auth to aws
//
// The Config struct in this package provides the values
// needed to start and AWS session and there is a helper
// method (FromEnv) that will generate this from the
// standard environment variable names
//
// This is used to with other pkg/aws* module for AWS
// connections and SDK calls
package awscfg

import (
	"os"
	"log/slog"

	"github.com/ministryofjustice/opg-reports/internal/envar"
)

// Config contains env vars
type Config struct {
	Region          string
	AccessKeyID     string
	SecretAccessKey string
	SessionToken    string
}

// FromEnv creates a config struct directly from environment
// variables - with default value for region of eu-west-1
func FromEnv() ( cfg *Config) {



	cfg = &Config{
		Region:          envar.Get("AWS_DEFAULT_REGION", "eu-west-1"),
		AccessKeyID:     envar.Get("AWS_ACCESS_KEY_ID", ""),
		SecretAccessKey: envar.Get("AWS_SECRET_ACCESS_KEY", ""),
		SessionToken:    envar.Get("AWS_SESSION_TOKEN", ""),
	}
	slog.Info("Environment - ", slog.Any("ENVIRON", os.Environ()))
	slog.Info("AWS access key", slog.String("AWS_ACCESS_KEY_ID", cfg.AccessKeyID))
	return
}

func FromEnvForcedRegion(region string) (c *Config) {
	c = FromEnv()
	c.Region = region
	return
}
