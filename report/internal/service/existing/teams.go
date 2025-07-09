package existing

import (
	"fmt"
	"os"
	"path/filepath"

	"opg-reports/report/internal/repository/githubr"
	"opg-reports/report/internal/repository/sqlr"
	"opg-reports/report/internal/utils"

	"github.com/google/go-github/v62/github"
)

// stmtTeamImport
const stmtTeamImport string = `
INSERT INTO teams (name)
VALUES (:name)
ON CONFLICT (name) DO UPDATE SET name=excluded.name RETURNING name;
`

// Team captures the team name under its prior version of
// `billing_unit` ready for inserting to db
//
// Example account from the opg-metadata source file:
//
//	{
//		"id": "500000067891",
//		"name": "My production",
//		"billing_unit": "Team A",
//		"label": "prod",
//		"environment": "production",
//		"type": "aws",
//		"uptime_tracking": true
//	}
//
// We only want `billing_unit` field and are ignoring the others
type teamItem struct {
	Name string `json:"billing_unit,omitempty" db:"name"`
}

// teamDownloadOptions used as just shorthand for passing lots of options around
type teamDownloadOptions struct {
	Owner      string
	Repository string
	AssetName  string
	ReleaseTag string
	Dir        string
	UseRegex   bool
}

// InsertTeams handles the inserting otf team data from opgmetadata reository
// into the local database service
//
// Example account from the opg-metadata source file:
//
//	{
//		"id": "500000067891",
//		"name": "My production",
//		"billing_unit": "Team A",
//		"label": "prod",
//		"environment": "production",
//		"type": "aws",
//		"uptime_tracking": true
//	}
//
// We only want `billing_unit` field and are ignoring the others
func (self *Service) InsertTeams(
	client githubr.ClientRepositoryReleases,
	ghs githubr.RepositoryReleases,
	sq sqlr.Writer,
) (results []*sqlr.BoundStatement, err error) {

	var dir string
	var sw = utils.Stopwatch()

	defer func() {
		self.log.With("seconds", sw.Stop().Seconds(), "inserted", len(results)).
			Info("[existing:Teams] existing func finished.")
	}()
	self.log.Info("[existing:Teams] starting existing records import ...")

	if ghs == nil {
		err = fmt.Errorf("ghs was nil")
		return
	}
	if sq == nil {
		err = fmt.Errorf("sq was nil")
		return
	}

	dir, err = os.MkdirTemp("./", "__download-gh-*")
	if err != nil {
		self.log.Error("mkdir error")
		return
	}
	defer os.RemoveAll(dir)

	teams, err := self.getTeamsFromMetadata(client, ghs, &teamDownloadOptions{
		Owner:      self.conf.Metadata.Owner,
		Repository: self.conf.Metadata.Repository,
		AssetName:  self.conf.Metadata.AssetName,
		UseRegex:   self.conf.Metadata.UseRegex,
		Dir:        dir,
	})

	if err != nil {
		self.log.Error("error getting team metadata")
		return
	}

	results, err = self.insertTeamsToDB(sq, teams)
	if err != nil {
		self.log.Error("error inserting teams")
		return
	}

	self.log.Info("[existing:Teams] existing records successful")
	return
}

// insertTeams handles writing the records to the table
func (self *Service) insertTeamsToDB(sq sqlr.Writer, teams []*teamItem) (statements []*sqlr.BoundStatement, err error) {
	statements = []*sqlr.BoundStatement{}

	for _, team := range teams {
		statements = append(statements, &sqlr.BoundStatement{Data: team, Statement: stmtTeamImport})
	}
	err = sq.Insert(statements...)
	return
}

// getTeamsFromMetadata downloads the release asset from repository, extracts it locally and converts the files
// into []Team
//
// Removes directory and files on exit
func (self *Service) getTeamsFromMetadata(
	client githubr.ClientRepositoryReleases,
	source githubr.RepositoryReleases,
	options *teamDownloadOptions,
) (teams []*teamItem, err error) {
	var (
		release      *github.RepositoryRelease
		asset        *github.ReleaseAsset
		fp           *os.File
		downloadedTo string
		accountFile  string = "accounts.json"
		downloadDir  string = filepath.Join(options.Dir, "download")
		extractDir   string = filepath.Join(options.Dir, "extract")
	)
	teams = []*teamItem{}
	ropts := &githubr.GetRepositoryReleaseOptions{
		ExcludePrereleases: true,
		ExcludeDraft:       true,
		ExcludeNoAssets:    true,
		ReleaseTag:         options.ReleaseTag,
		UseRegex:           options.UseRegex,
	}
	// find the release
	release, err = source.GetRepositoryRelease(
		client,
		options.Owner,
		options.Repository,
		ropts)
	if err != nil {
		self.log.Error("error finding repository release", "err", err.Error())
		return
	}
	if release == nil {
		err = fmt.Errorf("failed to find repository release")
		self.log.Error("failed finding repository release", "err", err.Error())
		return
	}
	// find the asset on the release
	asset, downloadedTo, err = source.DownloadRepositoryReleaseAsset(
		client,
		options.Owner,
		options.Repository,
		release,
		downloadDir,
		&githubr.DownloadRepositoryReleaseAssetOptions{
			AssetName: options.AssetName,
			UseRegex:  options.UseRegex,
		})
	if err != nil {
		self.log.Error("error downloading release by asset name", "err", err.Error())
		return
	}
	if asset == nil {
		err = fmt.Errorf("nil asset returned from DownloadLatestReleaseAssetByName")
		self.log.Error("error with asset name", "err", err.Error())
		return
	}
	// remove the files on exit
	defer func() {
		os.RemoveAll(downloadDir)
		os.RemoveAll(extractDir)
	}()

	teams, err = handleAsset[*teamItem](self.log, asset, fp, extractDir, downloadedTo, accountFile)

	return
}
