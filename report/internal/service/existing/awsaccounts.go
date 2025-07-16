package existing

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"opg-reports/report/internal/repository/githubr"
	"opg-reports/report/internal/repository/sqlr"
	"opg-reports/report/internal/utils"

	"github.com/google/go-github/v62/github"
)

// stmtAwsAccountImport
const stmtAwsAccountImport string = `
INSERT INTO aws_accounts (
	id,
	name,
	label,
	environment,
	uptime_tracking,
	team_name
) VALUES (
	:id,
	:name,
	:label,
	:environment,
	:uptime_tracking,
	:team_name
)
ON CONFLICT (id)
 	DO UPDATE SET
		name=excluded.name,
		label=excluded.label,
		environment=excluded.environment,
		uptime_tracking=excluded.uptime_tracking
RETURNING id;`

const stmtAwsAccountUpdateEmptyEnvironments string = `
UPDATE aws_accounts
SET
	environment = "production"
WHERE
	environment = ""
`

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
	ID             string `json:"id,omitempty" db:"id" example:"012345678910"` // This is the AWS Account ID as a string
	Name           string `json:"name,omitempty" db:"name" example:"Public API"`
	Label          string `json:"label,omitempty" db:"label" example:"aurora-cluster"`
	Environment    string `json:"environment,omitempty" db:"environment" example:"development|preproduction|production"`
	UptimeTracking bool   `json:"uptime_tracking" db:"uptime_tracking" example:"1|0"`
	TeamName       string `json:"billing_unit,omitempty" db:"team_name"`
}

// accountDownloadOptions used as just shorthand for passing lots of options around
type accountDownloadOptions struct {
	Owner      string
	Repository string
	ReleaseTag string
	AssetName  string
	Dir        string
	UseRegex   bool
}

// InsertAwsAccounts handles the inserting of data from opgmetadata reository
// into the local database service.
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
func (self *Service) InsertAwsAccounts(
	client githubr.ClientRepositoryReleases,
	source githubr.RepositoryReleases,
	sq sqlr.RepositoryWriter) (results []*sqlr.BoundStatement, err error) {
	var dir string
	var sw = utils.Stopwatch()

	defer func() {
		self.log.With("seconds", sw.Stop().Seconds(), "inserted", len(results)).
			Info("[existing:AwsAccounts] existing func finished.")
	}()
	self.log.Info("[existing:AwsAccounts] starting existing records import ...")

	if source == nil {
		err = fmt.Errorf("source was nil")
		return
	}
	if sq == nil {
		err = fmt.Errorf("sq was nil")
		return
	}
	dir, err = os.MkdirTemp("./", "__download-gh-*")
	if err != nil {
		return
	}
	defer os.RemoveAll(dir)

	teams, err := self.getAwsAccountsFromMetadata(client, source, &accountDownloadOptions{
		Owner:      self.conf.Metadata.Owner,
		Repository: self.conf.Metadata.Repository,
		ReleaseTag: self.conf.Metadata.ReleaseTag,
		AssetName:  self.conf.Metadata.AssetName,
		UseRegex:   self.conf.Metadata.UseRegex,
		Dir:        dir,
	})
	if err != nil {
		self.log.Error("failed to get aws accounts", "err", err.Error())
		return
	}

	results, err = self.insertAwsAccountsToDB(sq, teams)
	if err != nil {
		self.log.Error("failed on insertAwsAccountsToDB", "err", err.Error())
		return
	}

	self.log.Info("[existing:AwsAccounts] existing records successful")
	return
}

// insertTeams handles writing the records to the table
func (self *Service) insertAwsAccountsToDB(sq sqlr.RepositoryWriter, accounts []*awsAccount) (statements []*sqlr.BoundStatement, err error) {
	statements = []*sqlr.BoundStatement{}

	for _, acc := range accounts {
		statements = append(statements, &sqlr.BoundStatement{Data: acc, Statement: stmtAwsAccountImport})
	}
	err = sq.Insert(statements...)
	if err != nil {
		return
	}
	// update empty environment values to default them to production
	_, err = sq.Exec(stmtAwsAccountUpdateEmptyEnvironments)
	return
}

// getAwsAccountsFromMetadata downloads the release asset from repository, extracts it locally and converts the files
// into []awsAccount
//
// Removes directory and files on exit
func (self *Service) getAwsAccountsFromMetadata(
	client githubr.ClientRepositoryReleases,
	source githubr.RepositoryReleases,
	options *accountDownloadOptions,
) (accounts []*awsAccount, err error) {
	var (
		fp           *os.File
		release      *github.RepositoryRelease
		asset        *github.ReleaseAsset
		downloadedTo string
		accountFile  string = "accounts.aws.json"
		downloadDir  string = filepath.Join(options.Dir, "download")
		extractDir   string = filepath.Join(options.Dir, "extract")
	)
	accounts = []*awsAccount{}

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
		self.log.Error("error downloading release asset", "err", err.Error())
		return
	}
	if asset == nil {
		err = fmt.Errorf("nil asset returned from DownloadLatestReleaseAssetByName")
		return
	}
	// remove the files on exit
	defer func() {
		os.RemoveAll(downloadDir)
		os.RemoveAll(extractDir)
	}()

	accounts, err = handleAsset[*awsAccount](self.log, asset, fp, extractDir, downloadedTo, accountFile)

	return
}

func handleAsset[T Model](
	log *slog.Logger,
	asset *github.ReleaseAsset,
	fp *os.File,
	extractDir string,
	downloadedTo string,
	dataFile string,
) (data []T, err error) {
	data = []T{}
	// deal with tar balls
	if strings.HasSuffix(*asset.Name, "tar.gz") {
		// extract the zip file
		fp, err = os.Open(downloadedTo)
		if err != nil {
			log.Error("error opening release downloaded file", "err", err.Error())
			return
		}
		err = utils.TarGzExtract(extractDir, fp)
		if err != nil {
			log.Error("error extracting downloaded release", "err", err.Error())
			return
		}
		// check the accounts json file exists
		dataFile = filepath.Join(extractDir, dataFile)
		if !utils.DirExists(extractDir) || !utils.FileExists(dataFile) {
			err = fmt.Errorf("directory or file not found")
			return
		}
		// read the json file into local struct
		err = utils.UnmarshalFile(dataFile, &data)
	} else if strings.HasSuffix(*asset.Name, ".json") || strings.HasSuffix(*asset.Name, ".txt") {
		err = utils.UnmarshalFile(downloadedTo, &data)
	} else {
		err = fmt.Errorf("unsupported file type [name: %s] [type: %s]", *asset.Name, *asset.ContentType)
	}

	return
}
