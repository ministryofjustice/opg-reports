package client

import (
	"fmt"
	"opg-reports/report/internal/utils/logger"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

func TestUtilsAwsClient(t *testing.T) {
	var lg = logger.New("debug")
	var ctx = t.Context()

	c1, _ := New[*sts.Client](ctx, lg, "eu-west-1")
	if fmt.Sprintf("%T", c1) != "*sts.Client" {
		t.Errorf("incorrect client type")
	}

	c2, _ := New[*s3.Client](ctx, lg, "eu-west-1")
	if fmt.Sprintf("%T", c2) != "*s3.Client" {
		t.Errorf("incorrect client type")
	}

	c3, _ := New[*cloudwatch.Client](ctx, lg, "us-east-1")
	if fmt.Sprintf("%T", c3) != "*cloudwatch.Client" {
		t.Errorf("incorrect client type")
	}

}
