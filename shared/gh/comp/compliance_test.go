package comp

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

		com := NewWithR(nil, r, client)
		if comply, _, _ := com.Compliant(defaultBaselineCompliance); comply != true {
			t.Errorf("error with compliance")
		}
	}
}
func TestSharedGhComplianceComply(t *testing.T) {
	c := FakeCompliant(nil, defaultBaselineCompliance)
	if comply, _, _ := c.Compliant(defaultBaselineCompliance); !comply {
		t.Errorf("compliance failed")
	}

	c = FakeCompliant(nil, defaultExtendedCompliance)
	if comply, _, _ := c.Compliant(defaultExtendedCompliance); !comply {
		t.Errorf("compliance failed")
	}

}
