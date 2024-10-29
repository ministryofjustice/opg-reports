package awsclient

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/costexplorer"
)

func CostExplorer(s *session.Session) (ce *costexplorer.CostExplorer, err error) {
	return costexplorer.New(s), nil
}
