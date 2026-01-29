package ghclients

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/gofri/go-github-ratelimit/v2/github_ratelimit"
	"github.com/google/go-github/v81/github"
)

var ErrNoToken = errors.New("missing required token.")

// New returns a token based client for github usage
func New(ctx context.Context, log *slog.Logger, token string) (client *github.Client, err error) {
	var limited *http.Client
	var lg *slog.Logger = log.With("func", "utils.ghclients.New")

	lg.Debug("starting ...")
	if token == "" {
		err = ErrNoToken
		lg.Error("error creating github client", "err", err.Error())
		return
	}

	limited = github_ratelimit.NewClient(nil)
	client = github.NewClient(limited).WithAuthToken(token)
	lg.Debug("complete.")
	return

}
