// Package ghclient porivides consistent way of creating a client to use
// with api calls
package ghclient

import (
	"context"
	"os"

	"github.com/google/go-github/v87/github"
)

// Token helper to fetch value from env
func Token() string {
	var key string = "GH_TOKEN"
	return os.Getenv(key)
}

// New tries to return a new github client using the token value for auth
func New(ctx context.Context, token string) (client *github.Client, err error) {
	client, err = github.NewClient(github.WithAuthToken(token))
	return
}
