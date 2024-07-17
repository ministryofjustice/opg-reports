package cost

import (
	"opg-reports/shared/aws/sess"

	"github.com/aws/aws-sdk-go/service/costexplorer"
)

const defaultRegion string = "eu-west-1"

func ClientWithAssumedRole(roleArn string, region string) (*costexplorer.CostExplorer, error) {
	sessionName := "cost-retrival"
	if region == "" {
		region = defaultRegion
	}

	sess, err := sess.AssumeRole(roleArn, region, sessionName)
	if err != nil {
		return nil, err
	}
	return costexplorer.New(sess), nil
}

func ClientFromEnv() (*costexplorer.CostExplorer, error) {
	sess, err := sess.NewSessionFromEnv()
	if err != nil {
		return nil, err
	}
	return costexplorer.New(sess), nil
}
