package awsclients

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

var (
	ErrLoadingConfig   error = errors.New("error loading config.")
	ErrUnsupportedType error = errors.New("client type unsupported.")
)

// SupportedClients is a type constraint on creatign clients
type SupportedClients interface {
	*s3.Client | *sts.Client | *costexplorer.Client | *cloudwatch.Client
}

// New fetches a aws-sdk-v2 version of the appropriate client for T
//
// Supports: *s3.Client | *sts.Client | *costexplorer.Client | *cloudwatch.Client
func New[T SupportedClients](ctx context.Context, log *slog.Logger, region string) (T, error) {
	var (
		err    error
		awscfg aws.Config
		c      interface{}
		t      T
		lg     *slog.Logger = log.With("func", "awsclients.New")
	)
	lg.Debug("starting ...")

	awscfg, err = config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		err = errors.Join(ErrLoadingConfig, err)
		return nil, err
	}

	switch any(t).(type) {
	case *costexplorer.Client:
		c = costexplorer.NewFromConfig(awscfg)
	case *sts.Client:
		c = sts.NewFromConfig(awscfg)
	case *s3.Client:
		// disable checksum warning outputs
		c = s3.NewFromConfig(awscfg, func(o *s3.Options) {
			o.DisableLogOutputChecksumValidationSkipped = true
		})
	case *cloudwatch.Client:
		c = cloudwatch.NewFromConfig(awscfg)
	default:
		err = errors.Join(ErrUnsupportedType, fmt.Errorf("client type [%T] is not supported.", t))
		return nil, err
	}
	lg.Debug("complete.")
	return c.(T), nil

}
