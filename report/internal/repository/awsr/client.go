package awsr

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

// GetClient fetches a aws-sdk-v2 version of the appropriate client for T
//
// Supports: S3, STS
func GetClient[T *s3.Client | *sts.Client](ctx context.Context, region string) (T, error) {
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
		// disable checksum warning outputs
		c = s3.NewFromConfig(awscfg, func(o *s3.Options) {
			o.DisableLogOutputChecksumValidationSkipped = true
		})
		return c.(T), nil
	default:
		err = fmt.Errorf("client type [%T] unsupported", t)
	}

	return nil, err

}

func DefaultClient[T *s3.Client | *sts.Client](ctx context.Context, region string) (c T) {
	c, _ = GetClient[T](ctx, region)
	return
}
