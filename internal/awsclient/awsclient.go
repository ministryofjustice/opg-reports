// Package awsclient contains all of the helper methods to generate
// the appropriate sdk client
package awsclient

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/costexplorer"
)

func CostExplorer(s *session.Session) (ce *costexplorer.CostExplorer, err error) {
	return costexplorer.New(s), nil
}

func CloudWatch(s *session.Session) (*cloudwatch.CloudWatch, error) {
	return cloudwatch.New(s), nil
}
