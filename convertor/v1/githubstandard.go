package v1

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/ministryofjustice/opg-reports/internal/dateformats"
	"github.com/ministryofjustice/opg-reports/internal/dateutils"
	"github.com/ministryofjustice/opg-reports/models"
)

type GithubStandard struct {
	ID                             int    `json:"-"`
	Ts                             string `json:"ts"`
	CompliantBaseline              int    `json:"compliant_baseline"`
	CompliantExtended              int    `json:"compliant_extended"`
	CountOfClones                  int    `json:"count_of_clones"`
	CountOfForks                   int    `json:"count_of_forks"`
	CountOfPullRequests            int    `json:"count_of_pull_requests"`
	CountOfWebHooks                int    `json:"count_of_web_hooks"`
	CreatedAt                      string `json:"created_at"`
	DefaultBranch                  string `json:"default_branch"`
	FullName                       string `json:"full_name"`
	HasCodeOfConduct               int    `json:"has_code_of_conduct"`
	HasCodeownerApprovalRequired   int    `json:"has_codeowner_approval_required"`
	HasContributingGuide           int    `json:"has_contributing_guide"`
	HasDefaultBranchOfMain         int    `json:"has_default_branch_of_main"`
	HasDefaultBranchProtection     int    `json:"has_default_branch_protection"`
	HasDeleteBranchOnMerge         int    `json:"has_delete_branch_on_merge"`
	HasDescription                 int    `json:"has_description"`
	HasDiscussions                 int    `json:"has_discussions"`
	HasDownloads                   int    `json:"has_downloads"`
	HasIssues                      int    `json:"has_issues"`
	HasLicense                     int    `json:"has_license"`
	HasPages                       int    `json:"has_pages"`
	HasPullRequestApprovalRequired int    `json:"has_pull_request_approval_required"`
	HasReadme                      int    `json:"has_readme"`
	HasRulesEnforcedForAdmins      int    `json:"has_rules_enforced_for_admins"`
	HasVulnerabilityAlerts         int    `json:"has_vulnerability_alerts"`
	HasWiki                        int    `json:"has_wiki"`
	IsArchived                     int    `json:"is_archived"`
	IsPrivate                      int    `json:"is_private"`
	License                        string `json:"license"`
	LastCommitDate                 string `json:"last_commit_date"`
	Name                           string `json:"name"`
	Owner                          string `json:"owner"`
	Teams                          string `json:"teams"`
}

// MarshalJSON converts from current version to model version so
// when writing this struct to a json file it will take the form
// of a new model
func (self *GithubStandard) MarshalJSON() (bytes []byte, err error) {
	var (
		teams    []*models.GitHubTeam
		repo     *models.GitHubRepository
		ts       string
		standard *models.GitHubRepositoryStandard
		now      string = time.Now().UTC().Format(dateformats.Full)
	)
	if self.Ts == "" {
		self.Ts = now
	}
	ts = dateutils.Reformat(self.Ts, dateformats.Full)

	// setup new to include all of the old
	standard = &models.GitHubRepositoryStandard{
		Ts:                             self.Ts,
		CompliantBaseline:              uint8(self.CompliantBaseline),
		CompliantExtended:              uint8(self.CompliantExtended),
		CountOfClones:                  self.CountOfClones,
		CountOfForks:                   self.CountOfForks,
		CountOfPullRequests:            self.CountOfPullRequests,
		CountOfWebHooks:                self.CountOfWebHooks,
		DefaultBranch:                  self.DefaultBranch,
		HasCodeOfConduct:               uint8(self.HasCodeOfConduct),
		HasCodeownerApprovalRequired:   uint8(self.HasCodeownerApprovalRequired),
		HasContributingGuide:           uint8(self.HasContributingGuide),
		HasDefaultBranchOfMain:         uint8(self.HasDefaultBranchOfMain),
		HasDefaultBranchProtection:     uint8(self.HasDefaultBranchProtection),
		HasDeleteBranchOnMerge:         uint8(self.HasDeleteBranchOnMerge),
		HasDescription:                 uint8(self.HasDescription),
		HasDiscussions:                 uint8(self.HasDiscussions),
		HasDownloads:                   uint8(self.HasDownloads),
		HasIssues:                      uint8(self.HasIssues),
		HasLicense:                     uint8(self.HasLicense),
		HasPages:                       uint8(self.HasPages),
		HasPullRequestApprovalRequired: uint8(self.HasPullRequestApprovalRequired),
		HasReadme:                      uint8(self.HasReadme),
		HasRulesEnforcedForAdmins:      uint8(self.HasRulesEnforcedForAdmins),
		HasVulnerabilityAlerts:         uint8(self.HasVulnerabilityAlerts),
		HasWiki:                        uint8(self.HasWiki),
		IsArchived:                     uint8(self.IsArchived),
		IsPrivate:                      uint8(self.IsPrivate),
		License:                        self.License,
		LastCommitDate:                 self.LastCommitDate,
	}

	// create the team list
	for _, name := range strings.Split(self.Teams, "#") {
		if len(name) > 0 {
			team := &models.GitHubTeam{
				Ts:   ts,
				Slug: strings.ReplaceAll(strings.ToLower(name), " ", "-"),
			}
			team.Units = team.StandardUnits()
			teams = append(teams, team)
		}
	}

	repo = &models.GitHubRepository{
		Ts:             ts,
		Owner:          self.Owner,
		Name:           self.Name,
		FullName:       self.FullName,
		CreatedAt:      self.CreatedAt,
		DefaultBranch:  self.DefaultBranch,
		Archived:       uint8(self.IsArchived),
		Private:        uint8(self.IsPrivate),
		License:        self.License,
		LastCommitDate: self.LastCommitDate,
		GitHubTeams:    teams,
	}

	standard.GitHubRepositoryFullName = repo.FullName
	standard.GitHubRepository = (*models.GitHubRepositoryForeignKey)(repo)

	bytes, err = json.MarshalIndent(standard, "", "  ")

	return
}
