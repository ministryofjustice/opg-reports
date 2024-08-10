package ghs

import (
	"fmt"
)

func (g *GithubStandard) UID() string {
	return g.Uuid
}

func (g *GithubStandard) ToCSV() (line string) {

	line = fmt.Sprintf(`"%s","%s","%s","%s","%s","%s","%s","%s","%s",%d,%d,%d,%d,%d,%d,%d,%d,%d,%d,%d,%d,%d,%d,%d,%d,%d,%d,%d,%d,%d,%d,%d`,
		g.Uuid,
		g.Ts,
		g.DefaultBranch,
		g.FullName,
		g.Name,
		g.Owner,
		g.License,
		g.LastCommitDate,
		g.CreatedAt,
		g.CountOfClones,
		g.CountOfForks,
		g.CountOfPullRequests,
		g.CountOfWebHooks,
		g.HasCodeOfConduct,
		g.HasCodeownerApprovalRequired,
		g.HasContributingGuide,
		g.HasDefaultBranchOfMain,
		g.HasDefaultBranchProtection,
		g.HasDeleteBranchOnMerge,
		g.HasDescription,
		g.HasDiscussions,
		g.HasDownloads,
		g.HasIssues,
		g.HasLicense,
		g.HasPages,
		g.HasPullRequestApprovalRequired,
		g.HasReadme,
		g.HasRulesEnforcedForAdmins,
		g.HasVulnerabilityAlerts,
		g.HasWiki,
		g.IsArchived,
		g.IsPrivate,
	) + "\n"
	return
}
