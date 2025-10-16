package awsr

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

// SupportedClients is a type constraint on creatign clients
type SupportedClients interface {
	*s3.Client | *sts.Client | *costexplorer.Client | *cloudwatch.Client
}

// GetClient fetches a aws-sdk-v2 version of the appropriate client for T
//
// Supports: costexplorer, s3, sts, cloudwatch
func GetClient[T SupportedClients](ctx context.Context, region string) (T, error) {
	var err error
	var awscfg aws.Config
	var c interface{}
	var t T

	awscfg, err = config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return nil, err
	}

	switch any(t).(type) {
	case *costexplorer.Client:
		c = costexplorer.NewFromConfig(awscfg)
		return c.(T), nil
	case *sts.Client:
		c = sts.NewFromConfig(awscfg)
		return c.(T), nil
	case *s3.Client:
		// disable checksum warning outputs
		c = s3.NewFromConfig(awscfg, func(o *s3.Options) {
			o.DisableLogOutputChecksumValidationSkipped = true
		})
		return c.(T), nil
	case *cloudwatch.Client:
		c = cloudwatch.NewFromConfig(awscfg)
		return c.(T), nil
	default:
		err = fmt.Errorf("client type [%T] unsupported", t)
	}
	return nil, err

}

func DefaultClient[T SupportedClients](ctx context.Context, region string) (c T) {
	c, _ = GetClient[T](ctx, region)
	return
}
