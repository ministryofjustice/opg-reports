package githubr

import (
	"fmt"

	"github.com/gofri/go-github-ratelimit/github_ratelimit"
	"github.com/google/go-github/v62/github"
	"github.com/ministryofjustice/opg-reports/report/config"
)

// GetClient is an internal helper to handle creating the client
func GetClient(conf *config.Config) (client *github.Client, err error) {
	// handle empty configs
	if conf.Github == nil || conf.Github.Token == "" {
		return nil, fmt.Errorf("no github access token found in the config")
	}
	// get a rate limted version of the client
	limited, err := github_ratelimit.NewRateLimitWaiterClient(nil)
	if err != nil {
		return
	}
	client = github.NewClient(limited).WithAuthToken(conf.Github.Token)

	return
}

func DefaultClient(conf *config.Config) (client *github.Client) {
	client, _ = GetClient(conf)
	return
}
