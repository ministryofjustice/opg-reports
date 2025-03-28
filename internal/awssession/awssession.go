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

// New creates a typical aws session with the default region
func New() (sess *session.Session, err error) {
	var (
		region string
	)
	region = envar.Get("AWS_DEFAULT_REGION", "eu-west-1")

	slog.Debug("[awssession.New]", slog.String("region", region))

	return session.NewSession(&aws.Config{
		Region:      aws.String(region),
	})
}


// New creates a typical aws session with a set region
func NewWithRegion(region string) (sess *session.Session, err error) {
	slog.Debug("[awssession.New]", slog.String("region", region))

	return session.NewSession(&aws.Config{
		Region:      aws.String(region),
	})
}
