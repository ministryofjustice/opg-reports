package awssession

import (
	"log/slog"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/ministryofjustice/opg-reports/pkg/awscfg"
)

func New(cfg *awscfg.Config) (sess *session.Session, err error) {
	slog.Debug("[awssession.New]", slog.String("region", cfg.Region))

	return session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(cfg.AccessKeyID, cfg.SecretAccessKey, cfg.SessionToken),
		Region:      aws.String(cfg.Region),
	})
}
