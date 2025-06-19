package opgmetadata

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"slices"

	"github.com/google/go-github/v62/github"
	"github.com/ministryofjustice/opg-reports/report/config"
	"github.com/ministryofjustice/opg-reports/report/internal/gh"
	"github.com/ministryofjustice/opg-reports/report/internal/utils"
)

const (
	dataRepo  string = "opg-metadata"
	assetName string = "metadata.tar.gz"
)

type Service struct {
	ctx        context.Context
	log        *slog.Logger
	conf       *config.Config
	store      *gh.Repository
	directory  string
	downloaded bool
}

// SetDirectory changes the download / storage location
func (self *Service) SetDirectory(dir string) {
	self.directory = dir
}

// Download fetches the data from the repostitory asset and extracts the zip file to
// the local filesystem
func (self *Service) Download() (err error) {
	var (
		asset          *github.ReleaseAsset
		downloadedFile *os.File
		dir            string         = self.directory
		org            string         = self.conf.Github.Organisation
		log            *slog.Logger   = self.log.With("operation", "Download", "dataRepo", dataRepo, "assetName", assetName, "org", org)
		downloadTo     string         = filepath.Join(dir, assetName)
		extractTo      string         = filepath.Join(dir, dataRepo)
		gh             *gh.Repository = self.store
	)
	// if already downloaded, skip calling again
	if self.downloaded {
		return
	}
	// get the latest relase and the asset details that match the name
	log.Info("Downloading the latest release asset ...")
	asset, err = gh.GetLatestReleaseAsset(org, dataRepo, assetName, false)
	if err != nil {
		return
	}
	// download this asset
	log.With("assetID", *asset.ID).Debug("downloading the latest release asset via repository")
	downloadedFile, err = gh.DownloadReleaseAsset(org, dataRepo, *asset.ID, downloadTo)
	if err != nil {
		return
	}
	defer downloadedFile.Close()

	// now extract the tar.gz
	log.With("downloadedFile", downloadedFile, "extractTo", extractTo).Debug("extracting downloaded file...")
	err = utils.TarGzExtract(downloadedFile, extractTo)
	// set download flag
	self.downloaded = (err == nil)
	return
}

// GetAllAccounts returns all accounts from the meta data set which can then be used for import
// into the accounts table
func (self *Service) GetAllAccounts() (accounts []map[string]interface{}, err error) {
	var (
		dir         string = self.directory
		accountFile string = filepath.Join(dir, dataRepo, "accounts.json")
	)
	accounts = []map[string]interface{}{}
	// download the repo artifact
	err = self.Download()
	if err != nil {
		return
	}

	// check the account file exists
	if !utils.FileExists(accountFile) {
		err = fmt.Errorf("account data file not found [%s]", accountFile)
		return
	}
	// unmarshal the account data
	err = utils.UnmarshalFile(accountFile, &accounts)

	return
}

// GetAllTeams uses the list of all accounts to return the billing_unit names which
// are now used as team names for grouping of data
func (self *Service) GetAllTeams() (teams []map[string]interface{}, err error) {

	teams = []map[string]interface{}{}
	all := []string{}
	// get all the accounts
	accounts, err := self.GetAllAccounts()
	// get all the billing units from the accounts and make that a team
	for _, acc := range accounts {
		if val, ok := acc["billing_unit"]; ok {
			all = append(all, val.(string))
		}
	}
	// remove duplicates and create the output
	slices.Sort(all)
	all = slices.Compact(all)
	for _, nm := range all {
		teams = append(teams, map[string]interface{}{"name": nm})
	}

	return
}

// NewService returns a configured opgmetadata service object
func NewService(ctx context.Context, log *slog.Logger, conf *config.Config, store *gh.Repository) (srv *Service, err error) {
	if log == nil {
		return nil, fmt.Errorf("no logger passed for opgmetadata service")
	}
	if conf == nil {
		return nil, fmt.Errorf("no config passed for opgmetadata service")
	}
	if conf.Github == nil || conf.Github.Organisation == "" || conf.Github.Token == "" {
		return nil, fmt.Errorf("no github config details passed for opgmetadata service")
	}

	if store == nil {
		return nil, fmt.Errorf("no repository passed for opgmetadata service")
	}

	srv = &Service{
		ctx:       ctx,
		log:       log.With("service", "opgmetadata"),
		conf:      conf,
		store:     store,
		directory: "./__downloads/github/",
	}
	return
}
