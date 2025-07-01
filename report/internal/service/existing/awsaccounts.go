package existing

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ministryofjustice/opg-reports/report/internal/repository/githubr"
	"github.com/ministryofjustice/opg-reports/report/internal/repository/sqlr"
	"github.com/ministryofjustice/opg-reports/report/internal/utils"
)

// stmtAwsAccountImport
const stmtAwsAccountImport string = `
INSERT INTO aws_accounts (
	id,
	name,
	label,
	environment,
	team_name
) SELECT
	:id,
	:name,
	:label,
	:environment,
	teams.name
FROM teams WHERE name=:team_name LIMIT 1
ON CONFLICT (id)
 	DO UPDATE SET
		name=excluded.name,
		label=excluded.label,
		environment=excluded.environment
RETURNING id;`

// awsAccount captures an extra field from the metadata which
// is used in the stmtInsert to create the initial join to team based
// on the billing_unit name
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
type awsAccount struct {
	ID          string `json:"id,omitempty" db:"id" example:"012345678910"` // This is the AWS Account ID as a string
	Name        string `json:"name,omitempty" db:"name" example:"Public API"`
	Label       string `json:"label,omitempty" db:"label" example:"aurora-cluster"`
	Environment string `json:"environment,omitempty" db:"environment" example:"development|preproduction|production"`
	TeamName    string `json:"billing_unit,omitempty" db:"team_name"`
}

// accountDownloadOptions used as just shorthand for passing lots of options around
type accountDownloadOptions struct {
	Owner      string
	Repository string
	AssetName  string
	Dir        string
	UseRegex   bool
}

// InsertTeams handles the inserting otf team data from opgmetadata reository
// into the local database service
func (self *Service) InsertAwsAccounts(client githubr.ReleaseClient, ghs githubr.ReleaseRepository, sq sqlr.Writer) (results []*sqlr.BoundStatement, err error) {
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

	teams, err := self.getAwsAccountsFromMetadata(client, ghs, &accountDownloadOptions{
		Owner:      self.conf.Github.Organisation,
		Repository: self.conf.Github.Metadata.Repository,
		AssetName:  self.conf.Github.Metadata.Asset,
		UseRegex:   false,
		Dir:        dir,
	})
	if err != nil {
		return
	}
	results, err = self.insertAwsAccountsToDB(sq, teams)

	return
}

// insertTeams handles writing the records to the table
func (self *Service) insertAwsAccountsToDB(sq sqlr.Writer, accounts []*awsAccount) (statements []*sqlr.BoundStatement, err error) {
	statements = []*sqlr.BoundStatement{}

	for _, acc := range accounts {
		statements = append(statements, &sqlr.BoundStatement{Data: acc, Statement: stmtAwsAccountImport})
	}
	err = sq.Insert(statements...)
	return
}

// getTeamsFromMetadata downloads the release asset from repository, extracts it locally and converts the files
// into []Team
//
// Removes directory and files on exit
func (self *Service) getAwsAccountsFromMetadata(client githubr.ReleaseClient, ghs githubr.ReleaseRepository, options *accountDownloadOptions) (accounts []*awsAccount, err error) {
	var (
		fp           *os.File
		downloadedTo string
		accountFile  string = "accounts.aws.json"
		downloadDir  string = filepath.Join(options.Dir, "download")
		extractDir   string = filepath.Join(options.Dir, "extract")
	)
	accounts = []*awsAccount{}
	// Download the metadata asset
	downloadedTo, err = ghs.DownloadReleaseAssetByName(client,
		options.Owner,
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
		err = fmt.Errorf("directory or file not found: [%s] or [%s]", extractDir, accountFile)
		return
	}
	// read the json file into local struct
	err = utils.UnmarshalFile(accountFile, &accounts)
	return
}
