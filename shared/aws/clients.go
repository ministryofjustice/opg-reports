package aws

import (
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/costexplorer"
)

const defaultRegion string = "eu-west-1"

// CEClientFromEnv returns a cost explorer client
func CEClientFromEnv() (*costexplorer.CostExplorer, error) {
	sess, err := NewSessionFromEnv()
	if err != nil {
		return nil, err
	}
	return costexplorer.New(sess), nil
}

// CWClientFromEnv returns a cloudwatch client
func CWClientFromEnv(region string) (*cloudwatch.CloudWatch, error) {
	sess, err := NewSessionFromEnvWithRegion(region)
	if err != nil {
		return nil, err
	}
	return cloudwatch.New(sess), nil
}
