package std

import (
	"context"
	"opg-reports/shared/env"
	"opg-reports/shared/github/cl"
	"opg-reports/shared/github/repos"
	"opg-reports/shared/logger"
	"testing"
)

func TestSharedGithubStandardsRealData(t *testing.T) {
	logger.LogSetup()
	owner := "ministryofjustice"
	testRepo := "opg-incident-response"
	token := env.Get("GITHUB_ACCESS_TOKEN", "")
	if token != "" {
		ctx := context.Background()
		limiter, _ := cl.RateLimitedHttpClient()
		client := cl.Client(token, limiter)
		r, _ := repos.Get(ctx, client, owner, testRepo)

		com := NewWithR(nil, r, client)
		if comply, _, _ := com.Compliant(DefaultBaselineCompliance); comply != true {
			t.Errorf("error with compliance")
		}
	}
}
func TestSharedGithubStandardsComply(t *testing.T) {
	logger.LogSetup()
	c := FakeCompliant(nil, DefaultBaselineCompliance)
	if comply, _, _ := c.Compliant(DefaultBaselineCompliance); !comply {
		t.Errorf("compliance failed")
	}

	c = FakeCompliant(nil, DefaultExtendedCompliance)
	if comply, _, _ := c.Compliant(DefaultExtendedCompliance); !comply {
		t.Errorf("compliance failed")
	}

}
