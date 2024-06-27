package compliance

import (
	"fmt"
	"log/slog"
	"opg-reports/shared/fake"

	"github.com/google/uuid"
)

// Fake returns a generated Cost item using fake data
// If you pass an existing cost item in, it will fill in blank fields only
func Fake(c *Compliance) (f *Compliance) {

	if c == nil {
		c = New(nil)
	}
	if c.UUID == "" {
		c.UUID = uuid.NewString()
	}

	if c.DefaultBranch == "" {
		c.DefaultBranch = fake.Choice[string]([]string{"main", "master"})
	}
	if c.Owner == "" {
		c.Owner = fake.String(12)
	}
	if c.Name == "" {
		c.Name = fake.String(10)
	}
	if c.FullName == "" {
		c.FullName = fmt.Sprintf("%s/%s", c.Owner, c.Name)
	}

	if c.License == "" {
		c.License = fake.Choice[string]([]string{"MIT", "GPL", ""})
	}

	if c.CountClones == "" {
		c.CountClones = fake.IntAsStr(1, 30)
	}
	if c.CountForks == "" {
		c.CountForks = fake.IntAsStr(0, 10)
	}
	if c.CountPullRequests == "" {
		c.CountPullRequests = fake.IntAsStr(0, 12)
	}
	if c.CountWebhooks == "" {
		c.CountWebhooks = fake.IntAsStr(0, 5)
	}
	// Booleans
	c.Archived = fake.Choice[string]([]string{"true", "false"})

	c.HasDescription = fake.Choice[string]([]string{"true", "false"})
	c.HasDiscussions = fake.Choice[string]([]string{"true", "false"})
	c.HasDownloads = fake.Choice[string]([]string{"true", "false"})
	c.HasIssues = fake.Choice[string]([]string{"true", "false"})
	c.HasPages = fake.Choice[string]([]string{"true", "false"})
	c.HasProjects = fake.Choice[string]([]string{"true", "false"})
	c.HasWiki = fake.Choice[string]([]string{"true", "false"})
	c.IsPrivate = fake.Choice[string]([]string{"true", "false"})

	c.HasCodeOfConduct = fake.Choice[string]([]string{"true", "false"})
	c.HasCodeOwnerApprovalRequired = fake.Choice[string]([]string{"true", "false"})
	c.HasContributingGuide = fake.Choice[string]([]string{"true", "false"})
	c.HasDefaultBranchProtection = fake.Choice[string]([]string{"true", "false"})
	c.HasReadme = fake.Choice[string]([]string{"true", "false"})
	c.HasRulesEnforcedForAdmins = fake.Choice[string]([]string{"true", "false"})
	c.HasPullRequestApprovalRequired = fake.Choice[string]([]string{"true", "false"})
	c.HasVulnerabilityAlerts = fake.Choice[string]([]string{"true", "false"})

	c.Baseline = bTs(c.CompliesWithBaseline())
	c.Extended = bTs(c.CompliesWithExtended())

	f = c
	slog.Debug("[aws/cost] fake", slog.String("UID", f.UID()))
	return
}
