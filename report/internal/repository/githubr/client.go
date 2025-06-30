package githubr

import (
	"fmt"

	"github.com/gofri/go-github-ratelimit/github_ratelimit"
	"github.com/google/go-github/v62/github"
)

// client is an internal helper to handle creating the client
func (self *Repository) client() (client *github.Client, err error) {
	// handle empty configs
	if self.conf.Github == nil || self.conf.Github.Token == "" {
		return nil, fmt.Errorf("no github access token found in the config")
	}
	// get a rate limted version of the client
	limited, err := github_ratelimit.NewRateLimitWaiterClient(nil)
	if err != nil {
		return
	}
	client = github.NewClient(limited).WithAuthToken(self.conf.Github.Token)

	return
}
