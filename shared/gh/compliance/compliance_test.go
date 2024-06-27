package compliance

import (
	"context"
	"opg-reports/shared/env"
	"opg-reports/shared/gh/cl"
	"opg-reports/shared/gh/repos"
	"testing"
)

func TestSharedGhComplianceRealData(t *testing.T) {
	owner := "ministryofjustice"
	testRepo := "opg-incident-response"
	token := env.Get("GITHUB_ACCESS_TOKEN", "")
	if token != "" {
		ctx := context.Background()
		limiter, _ := cl.RateLimitedHttpClient()
		client := cl.Client(token, limiter)
		r, _ := repos.Get(ctx, client, owner, testRepo)

		comp := NewWithR(nil, r, client)
		if !comp.Baseline {
			t.Errorf("failed baseline")
		}
	}
}
