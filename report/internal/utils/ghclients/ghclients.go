package ghclients

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/gofri/go-github-ratelimit/v2/github_ratelimit"
	"github.com/google/go-github/v81/github"
)

// New returns a token based client for github usage
func New(ctx context.Context, log *slog.Logger, token string) (client *github.Client, err error) {
	var limited *http.Client

	if token == "" {
		err = ErrNoToken
		log.Error("error creating github client", "err", err.Error())
		return
	}

	limited = github_ratelimit.NewClient(nil)
	client = github.NewClient(limited).WithAuthToken(token)

	return

}
