// Package awssession provides wrapper funcs to create sdk sessions
//
// The standard New func uses the awscfg.Config struct to get
// all details for a standard session
package awssession

import (
	"log/slog"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/ministryofjustice/opg-reports/internal/envar"
)

// New creates a typical aws session from the config struct
func New(region string) (sess *session.Session, err error) {
	if region == "" {
		region = envar.Get("AWS_DEFAULT_REGION", "eu-west-1")
	}
	slog.Debug("[awssession.New]", slog.String("region", region))

	return session.NewSession(&aws.Config{
		Region:      aws.String(region),
	})
}
