package comp

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

	if c.CountClones == 0 {
		c.CountClones = fake.Int(1, 30)
	}
	if c.CountForks == 0 {
		c.CountForks = fake.Int(0, 10)
	}
	if c.CountPullRequests == 0 {
		c.CountPullRequests = fake.Int(0, 12)
	}
	if c.CountWebhooks == 0 {
		c.CountWebhooks = fake.Int(0, 5)
	}
	// Booleans
	c.Archived = fake.Choice[bool]([]bool{true, false})

	c.HasDescription = fake.Choice[bool]([]bool{true, false})
	c.HasDiscussions = fake.Choice[bool]([]bool{true, false})
	c.HasDownloads = fake.Choice[bool]([]bool{true, false})
	c.HasIssues = fake.Choice[bool]([]bool{true, false})
	c.HasPages = fake.Choice[bool]([]bool{true, false})
	c.HasProjects = fake.Choice[bool]([]bool{true, false})
	c.HasWiki = fake.Choice[bool]([]bool{true, false})
	c.IsPrivate = fake.Choice[bool]([]bool{true, false})

	c.HasCodeOfConduct = fake.Choice[bool]([]bool{true, false})
	c.HasCodeOwnerApprovalRequired = fake.Choice[bool]([]bool{true, false})
	c.HasContributingGuide = fake.Choice[bool]([]bool{true, false})
	c.HasDefaultBranchProtection = fake.Choice[bool]([]bool{true, false})
	c.HasReadme = fake.Choice[bool]([]bool{true, false})
	c.HasRulesEnforcedForAdmins = fake.Choice[bool]([]bool{true, false})
	c.HasPullRequestApprovalRequired = fake.Choice[bool]([]bool{true, false})
	c.HasVulnerabilityAlerts = fake.Choice[bool]([]bool{true, false})

	c.Baseline = c.CompliesWithBaseline()
	c.Extended = c.CompliesWithExtended()

	f = c
	slog.Debug("[aws/cost] fake", slog.String("UID", f.UID()))
	return
}
