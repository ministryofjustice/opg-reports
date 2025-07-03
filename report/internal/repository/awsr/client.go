package awsr

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

type client interface {
	*s3.Client | *sts.Client
}

// GetClient fetches a v2 version of the appropriate client from various AWS
// SDK libs.
//
// Supports: S3, STS
func GetClient[T client](ctx context.Context, region string) (T, error) {
	var err error
	var awscfg aws.Config
	var c interface{}
	var t T

	awscfg, err = config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return nil, err
	}

	switch any(t).(type) {
	case *sts.Client:
		c = sts.NewFromConfig(awscfg)
		return c.(T), nil
	case *s3.Client:
		// disable checksum warnings
		c = s3.NewFromConfig(awscfg, func(o *s3.Options) {
			o.DisableLogOutputChecksumValidationSkipped = true
		})
		return c.(T), nil
	default:
		err = fmt.Errorf("client type [%T] unsupported", t)
	}

	return nil, err

}

func DefaultClient[T client](ctx context.Context, region string) (c T) {
	c, _ = GetClient[T](ctx, region)
	return
}
