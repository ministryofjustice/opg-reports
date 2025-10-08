package githubr

import (
	"fmt"

	"opg-reports/report/config"

	"github.com/gofri/go-github-ratelimit/v2/github_ratelimit"
	"github.com/google/go-github/v75/github"
)

// GetClient is an internal helper to handle creating the client
func GetClient(conf *config.Config) (client *github.Client, err error) {
	// handle empty configs
	if conf.Github == nil || conf.Github.Token == "" {
		return nil, fmt.Errorf("no github access token found in the config")
	}
	// get a rate limted version of the client
	limited := github_ratelimit.NewClient(nil)
	client = github.NewClient(limited).WithAuthToken(conf.Github.Token)

	return
}

func DefaultClient(conf *config.Config) (client *github.Client) {
	client, _ = GetClient(conf)
	return
}
