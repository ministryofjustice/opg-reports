package githubclient

import (
	"net/http"

	"github.com/gofri/go-github-ratelimit/github_ratelimit"
	"github.com/google/go-github/v62/github"
)

func RateLimitedHttpClient() (*http.Client, error) {
	return github_ratelimit.NewRateLimitWaiterClient(nil)
}

func Client(token string) *github.Client {
	httpClient, _ := RateLimitedHttpClient()
	return github.NewClient(httpClient).WithAuthToken(token)
}
