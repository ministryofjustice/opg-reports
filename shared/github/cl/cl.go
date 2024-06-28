package cl

import (
	"net/http"

	"github.com/gofri/go-github-ratelimit/github_ratelimit"
	"github.com/google/go-github/v62/github"
)

func RateLimitedHttpClient() (*http.Client, error) {
	return github_ratelimit.NewRateLimitWaiterClient(nil)
}

func Client(token string, httpClient *http.Client) *github.Client {
	return github.NewClient(httpClient).WithAuthToken(token)
}
