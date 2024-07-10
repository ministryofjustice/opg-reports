package sess

import (
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

	return session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(
			*assumedRole.Credentials.AccessKeyId,
			*assumedRole.Credentials.SecretAccessKey,
			*assumedRole.Credentials.SessionToken),
		Region: aws.String(region),
	})
}
