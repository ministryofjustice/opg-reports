package existing

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ministryofjustice/opg-reports/report/internal/repository/githubr"
	"github.com/ministryofjustice/opg-reports/report/internal/repository/sqlr"
	"github.com/ministryofjustice/opg-reports/report/internal/utils"
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
type team struct {
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
func (self *Service) InsertTeams(ghs githubr.Releaser, sq sqlr.Writer) (results []*sqlr.BoundStatement, err error) {
	var dir string

	if ghs == nil || sq == nil {
		err = fmt.Errorf("params were nil")
		return
	}

	dir, err = os.MkdirTemp("./", "__download-gh-*")
	if err != nil {
		return
	}
	defer os.RemoveAll(dir)

	teams, err := self.getTeamsFromMetadata(ghs, &teamDownloadOptions{
		Owner:      self.conf.Github.Organisation,
		Repository: self.conf.Github.Metadata.Repository,
		AssetName:  self.conf.Github.Metadata.Asset,
		UseRegex:   false,
		Dir:        dir,
	})
	if err != nil {
		return
	}
	results, err = self.insertTeamsToDB(sq, teams)

	return
}

// insertTeams handles writing the records to the table
func (self *Service) insertTeamsToDB(sq sqlr.Writer, teams []*team) (statements []*sqlr.BoundStatement, err error) {
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
func (self *Service) getTeamsFromMetadata(ghs githubr.Releaser, options *teamDownloadOptions) (accounts []*team, err error) {
	var (
		fp           *os.File
		downloadedTo string
		accountFile  string = "accounts.json"
		downloadDir  string = filepath.Join(options.Dir, "download")
		extractDir   string = filepath.Join(options.Dir, "extract")
	)
	accounts = []*team{}
	// Download the metadata asset
	downloadedTo, err = ghs.DownloadReleaseAssetByName(options.Owner,
		options.Repository,
		options.AssetName,
		options.UseRegex,
		downloadDir)

	if err != nil {
		return
	}
	// remove the files on exit
	defer func() {
		os.RemoveAll(downloadDir)
		os.RemoveAll(extractDir)
	}()
	// extract the zip file
	fp, err = os.Open(downloadedTo)
	if err != nil {
		return
	}
	err = utils.TarGzExtract(extractDir, fp)
	if err != nil {
		return
	}
	// check the accounts json file exists
	accountFile = filepath.Join(extractDir, accountFile)
	if !utils.DirExists(extractDir) || !utils.FileExists(accountFile) {
		err = fmt.Errorf("directory or file not found")
		return
	}
	// read the json file into local struct
	err = utils.UnmarshalFile(accountFile, &accounts)
	return
}
