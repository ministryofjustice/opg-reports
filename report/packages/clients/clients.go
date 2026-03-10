package clients

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"opg-reports/report/packages/logger"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/gofri/go-github-ratelimit/v2/github_ratelimit"
	"github.com/google/go-github/v84/github"
)

type SupportedClients interface {
	*github.Client |
		*s3.Client | *sts.Client | *costexplorer.Client | *cloudwatch.Client
}

var (
	ErrLoadingConfig   error = errors.New("error loading config.")
	ErrUnsupportedType error = errors.New("client type unsupported.")
)

// New generates a configured client of tye T using values required.
//
// `param` is presumed to be a token for github clients or a region
// for aws clients.
func New[T SupportedClients](ctx context.Context, param string) (T, error) {
	var (
		err    error
		awscfg aws.Config
		c      interface{}
		t      T
	)
	// deal with github client
	switch any(t).(type) {
	case *github.Client:
		c, err = ghClient(ctx, param)
		return c.(T), err
	}

	// deal with aws values
	awscfg, err = config.LoadDefaultConfig(ctx, config.WithRegion(param))
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
	return c.(T), err
}

// ghClient returns a token based ratelimited client for github usage
func ghClient(ctx context.Context, token string) (client *github.Client, err error) {
	var limited *http.Client
	var log *slog.Logger
	ctx, log = logger.Get(ctx)

	log.Debug("creating github client ...")
	if token == "" {
		log.Error("no token found for githug client")
		return
	}

	limited = github_ratelimit.NewClient(nil)
	client = github.NewClient(limited).WithAuthToken(token)

	log.Debug("github client completed..")
	return

}
