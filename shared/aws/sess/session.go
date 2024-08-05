package sess

import (
	"log/slog"
	"opg-reports/shared/env"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
)

func AssumeRole(roleArn string, region string, sessionName string) (*session.Session, error) {
	baseSess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(region),
	}))
	stsSvc := sts.New(baseSess)
	assumedRole, err := stsSvc.AssumeRole(&sts.AssumeRoleInput{
		RoleArn:         aws.String(roleArn),
		RoleSessionName: aws.String(sessionName),
	})
	if err != nil {
		return nil, err
	}
	return NewSession(
		*assumedRole.Credentials.AccessKeyId,
		*assumedRole.Credentials.SecretAccessKey,
		*assumedRole.Credentials.SessionToken,
		region)

}

func NewSession(id string, secret string, token string, region string) (*session.Session, error) {
	return session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(id, secret, token),
		Region:      aws.String(region),
	})
}

func NewSessionFromEnv() (*session.Session, error) {
	return NewSessionFromEnvWithRegion(env.Get("AWS_DEFAULT_REGION", "eu-west-1"))
}

func NewSessionFromEnvWithRegion(region string) (*session.Session, error) {
	id := env.Get("AWS_ACCESS_KEY_ID", "")
	secret := env.Get("AWS_SECRET_ACCESS_KEY", "")
	token := env.Get("AWS_SESSION_TOKEN", "")
	slog.Info("AWS session", slog.String("region", region))
	return session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(id, secret, token),
		Region:      aws.String(region),
	})
}
