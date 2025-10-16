package awsr

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

func TestClientGeneric(t *testing.T) {

	c1, _ := GetClient[*sts.Client](context.TODO(), "eu-west-1")
	if fmt.Sprintf("%T", c1) != "*sts.Client" {
		t.Errorf("incorrect client type")
	}

	c2, _ := GetClient[*s3.Client](context.TODO(), "eu-west-1")
	if fmt.Sprintf("%T", c2) != "*s3.Client" {
		t.Errorf("incorrect client type")
	}

	c3, _ := GetClient[*cloudwatch.Client](context.TODO(), "us-east-1")
	if fmt.Sprintf("%T", c3) != "*cloudwatch.Client" {
		t.Errorf("incorrect client type")
	}

}
