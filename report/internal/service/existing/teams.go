package existing

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

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
func (self *Service) InsertTeams(client githubr.ReleaseClient, ghs githubr.ReleaseRepositoryDownloader, sq sqlr.Writer) (results []*sqlr.BoundStatement, err error) {
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
		Owner:      self.conf.Github.Organisation,
		Repository: self.conf.Github.Metadata.Repository,
		AssetName:  self.conf.Github.Metadata.Asset,
		UseRegex:   false,
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
func (self *Service) getTeamsFromMetadata(client githubr.ClientReleaseGetAndDownloader, ghs githubr.ReleaseRepositoryDownloader, options *teamDownloadOptions) (teams []*teamItem, err error) {
	var (
		asset        *github.ReleaseAsset
		fp           *os.File
		downloadedTo string
		accountFile  string = "accounts.json"
		downloadDir  string = filepath.Join(options.Dir, "download")
		extractDir   string = filepath.Join(options.Dir, "extract")
	)
	teams = []*teamItem{}
	// Download the metadata asset
	asset, downloadedTo, err = ghs.DownloadReleaseAssetByName(client,
		options.Owner,
		options.Repository,
		options.AssetName,
		options.UseRegex,
		downloadDir)

	if err != nil {
		self.log.Error("error downloading release by asset name", "err", err.Error())
		return
	}
	if asset == nil {
		err = fmt.Errorf("nil asset returned from DownloadReleaseAssetByName")
		self.log.Error("error with asset name", "err", err.Error())
		return
	}
	// remove the files on exit
	defer func() {
		os.RemoveAll(downloadDir)
		os.RemoveAll(extractDir)
	}()

	// deal with tar balls
	if strings.HasSuffix(*asset.Name, "tar.gz") {
		// extract the zip file
		fp, err = os.Open(downloadedTo)
		if err != nil {
			self.log.Error("error opening release downloaded file", "err", err.Error())
			return
		}
		err = utils.TarGzExtract(extractDir, fp)
		if err != nil {
			self.log.Error("error extracting downloaded release", "err", err.Error())
			return
		}
		// check the accounts json file exists
		accountFile = filepath.Join(extractDir, accountFile)
		if !utils.DirExists(extractDir) || !utils.FileExists(accountFile) {
			err = fmt.Errorf("directory or file not found")
			return
		}
		// read the json file into local struct
		err = utils.UnmarshalFile(accountFile, &teams)
	} else if strings.HasSuffix(*asset.Name, ".json") || strings.HasSuffix(*asset.Name, ".txt") {
		err = utils.UnmarshalFile(downloadedTo, &teams)
	} else {
		err = fmt.Errorf("unsupported file type [name: %s] [type: %s]", *asset.Name, *asset.ContentType)
	}
	return
}
