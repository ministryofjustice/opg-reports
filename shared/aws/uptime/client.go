package uptime

import (
	"opg-reports/shared/aws/sess"

	"github.com/aws/aws-sdk-go/service/cloudwatch"
)

func ClientFromEnv(region string) (*cloudwatch.CloudWatch, error) {
	sess, err := sess.NewSessionFromEnvWithRegion(region)
	if err != nil {
		return nil, err
	}
	return cloudwatch.New(sess), nil
}
