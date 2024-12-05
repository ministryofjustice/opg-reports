// Package awssession provides wrapper funcs to create sdk sessions
//
// The standard New func uses the awscfg.Config struct to get
// all details for a standard session
package awssession

import (
	"log/slog"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/ministryofjustice/opg-reports/internal/awscfg"
)

// New creates a typical aws session from the config struct
func New(cfg *awscfg.Config) (sess *session.Session, err error) {
	slog.Debug("[awssession.New]", slog.String("region", cfg.Region))

	return session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(cfg.AccessKeyID, cfg.SecretAccessKey, cfg.SessionToken),
		Region:      aws.String(cfg.Region),
	})
}
