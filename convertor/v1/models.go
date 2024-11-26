package v1

import (
	"strings"
	"time"

	"github.com/ministryofjustice/opg-reports/internal/dateformats"
	"github.com/ministryofjustice/opg-reports/internal/dateutils"
	"github.com/ministryofjustice/opg-reports/internal/structs"
	"github.com/ministryofjustice/opg-reports/models"
)

type AwsCost struct {
	ID           int    `json:"id"`
	Ts           string `json:"ts"`
	Organisation string `json:"organisation"`
	AccountID    string `json:"account_id"`
	AccountName  string `json:"account_name"`
	Unit         string `json:"unit"`
	Label        string `json:"label"`
	Environment  string `json:"environment"`
	Service      string `json:"service"`
	Region       string `json:"region"`
	Date         string `json:"date"`
	Cost         string `json:"cost"`
}

// V2 converts this version of AwsCosts into a models.AwsCost
func (self *AwsCost) V2() *models.AwsCost {
	var (
		unit    *models.Unit
		account *models.AwsAccount
		cost    *models.AwsCost
		ts      string
		now     string = time.Now().UTC().Format(dateformats.Full)
	)
	if self.Ts == "" {
		self.Ts = now
	}
	ts = dateutils.Reformat(self.Ts, dateformats.Full)

	unit = &models.Unit{
		Ts:   ts,
		Name: strings.ToLower(self.Unit),
	}
	account = &models.AwsAccount{
		Ts:          ts,
		Number:      self.AccountID,
		Name:        self.AccountName,
		Label:       self.Label,
		Environment: self.Environment,
		Unit:        (*models.UnitForeignKey)(unit),
	}
	cost = &models.AwsCost{
		Ts:         ts,
		Region:     self.Region,
		Service:    self.Service,
		Date:       self.Date,
		Cost:       self.Cost,
		AwsAccount: (*models.AwsAccountForeignKey)(account),
		Unit:       (*models.UnitForeignKey)(unit),
	}
	return cost
}

type AwsUptime struct {
	ID      int     `json:"id"`
	Ts      string  `json:"ts"`
	Unit    string  `json:"unit"`
	Average float64 `json:"average"`
	Date    string  `json:"date"`
}

// V2 converts
// Not a complete conversion, have to work out account data after
func (self *AwsUptime) V2() *models.AwsUptime {
	var (
		account *models.AwsAccount
		uptime  *models.AwsUptime
		unit    *models.Unit
		ts      string
		now     string = time.Now().UTC().Format(dateformats.Full)
	)
	if self.Ts == "" {
		self.Ts = now
	}
	ts = dateutils.Reformat(self.Ts, dateformats.Full)

	unit = &models.Unit{
		Ts:   ts,
		Name: strings.ToLower(self.Unit),
	}
	account = self.Account(unit)
	uptime = &models.AwsUptime{
		Ts:         ts,
		Date:       self.Date,
		Average:    self.Average,
		Unit:       (*models.UnitForeignKey)(unit),
		AwsAccount: (*models.AwsAccountForeignKey)(account),
	}
	return uptime
}

func (self *AwsUptime) Account(unit *models.Unit) (account *models.AwsAccount) {
	account = &models.AwsAccount{Environment: "production"}

	switch unit.Name {
	case "digideps":
		account.Number = "515688267891"
	case "make":
		account.Number = "980242665824"
	case "modernise":
		account.Number = "313879017102"
	case "serve":
		account.Number = "933639921819"
	case "sirius":
		account.Number = "649098267436"
	case "use":
		account.Number = "690083044361"
	}

	return
}

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

func (self *GithubStandard) V2() *models.GitHubRepositoryStandard {
	var (
		teams    []*models.GitHubTeam
		repo     *models.GitHubRepository
		standard = &models.GitHubRepositoryStandard{}
		ts       string
		now      string = time.Now().UTC().Format(dateformats.Full)
	)
	if self.Ts == "" {
		self.Ts = now
	}
	ts = dateutils.Reformat(self.Ts, dateformats.Full)

	// swap over
	structs.Convert(self, standard)
	// create the team list
	for _, name := range strings.Split(self.Teams, "#") {
		if len(name) > 0 {
			team := &models.GitHubTeam{
				Ts:   ts,
				Slug: strings.ReplaceAll(strings.ToLower(name), " ", "-"),
			}
			team.Units = self.TeamUnits(team)
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

	return standard
}

func (self *GithubStandard) TeamUnits(team *models.GitHubTeam) (units models.Units) {
	units = team.StandardUnits()
	return
}
