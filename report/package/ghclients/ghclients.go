package ghclients

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"opg-reports/report/package/cntxt"

	"github.com/gofri/go-github-ratelimit/v2/github_ratelimit"
	"github.com/google/go-github/v81/github"
)

var ErrNoToken = errors.New("missing required token.")

// New returns a token based client for github usage
func New(ctx context.Context, token string) (client *github.Client, err error) {
	var limited *http.Client
	var log *slog.Logger = cntxt.GetLogger(ctx).With("package", "ghclients", "func", "New")

	log.Debug("starting ...")
	if token == "" {
		err = ErrNoToken
		log.Error("error creating github client", "err", err.Error())
		return
	}

	limited = github_ratelimit.NewClient(nil)
	client = github.NewClient(limited).WithAuthToken(token)
	log.Debug("complete.")
	return

}
