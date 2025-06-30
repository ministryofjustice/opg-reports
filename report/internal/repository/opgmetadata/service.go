package opgmetadata

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/google/go-github/v62/github"
	"github.com/ministryofjustice/opg-reports/report/config"
	"github.com/ministryofjustice/opg-reports/report/internal/interfaces"
	"github.com/ministryofjustice/opg-reports/report/internal/repository/gh"
	"github.com/ministryofjustice/opg-reports/report/internal/utils"
)

type Service[T interfaces.Model] struct {
	ctx       context.Context
	log       *slog.Logger
	conf      *config.Config
	store     *gh.Repository
	directory string
}

func (self *Service[T]) GetStore() *gh.Repository {
	return self.store
}

// SetDirectory changes the download / storage location
func (self *Service[T]) SetDirectory(dir string) {
	self.directory = dir
}

// GetDirectory
func (self *Service[T]) GetDirectory() string {
	if self.directory == "" {
		dir, _ := os.MkdirTemp("./", "__download_gh_*")
		self.directory = dir
		defer os.RemoveAll(dir)
	}
	return self.directory
}

// Close cleans up tmp values etc
func (self *Service[T]) Close() (err error) {
	err = os.RemoveAll(self.GetDirectory())
	return
}

// DownloadAndReturnAll downloads the asset from the github repository (<owner>/<repository>) to local
// storage, finds all `json` files in the extracted data and reads the content of those into a T[].
//
// Assumes each data file contains a list of T already.
func (self *Service[T]) DownloadAndReturnAll(owner string, repository string, assetName string, regex bool) (data []T, err error) {
	var (
		extractedDir string
		dataFiles    = []string{}
		log          = self.log
	)
	log = log.With("operation", "DownloadAndExtractAsset",
		"repository", repository,
		"assetName", assetName,
		"owner", owner)

	data = []T{}

	extractedDir, err = self.DownloadAndExtractAsset(owner, repository, assetName, regex)
	if err != nil {
		return
	}

	dataFiles = utils.FileList(extractedDir, ".json")
	// if there are now files, return early, as theres nothing else to do
	if len(dataFiles) <= 0 {
		return
	}
	// read each file, assume each contains many T and merge with main data
	for _, file := range dataFiles {
		var list = []T{}
		log.With("file", file).Debug("reading file into slice of T")
		err = utils.StructFromJsonFile(file, &list)
		if err != nil {
			return
		}
		data = append(data, list...)
	}
	log.With("count", len(data)).Debug("downloaded and converted")

	return
}

// DownloadAndReturnAll downloads the asset from the github repository (<owner>/<repository>) to local
// storage, finds the file matching `filename` and reads that into T[]
//
// filename should be relative path on where you would expect the file to be in the extracted data folder
func (self *Service[T]) DownloadAndReturn(owner string, repository string, assetName string, regex bool, filename string) (data []T, err error) {
	var (
		extractedDir string
		file         string
		log          = self.log
	)
	log = log.With("operation", "DownloadAndExtractAsset",
		"repository", repository,
		"assetName", assetName,
		"owner", owner)

	data = []T{}

	extractedDir, err = self.DownloadAndExtractAsset(owner, repository, assetName, regex)
	if err != nil {
		return
	}

	file = filepath.Join(extractedDir, filename)
	err = utils.StructFromJsonFile(file, &data)
	if err != nil {
		return
	}
	log.With("count", len(data)).Debug("downloaded and converted")
	return
}

// DownloadAndExtractAsset fetches the asset (assumed tar.gz) from the github repository (<owner>/<repository>)
// to a local directory and then extracts the content.
//
// The path to the extracted directory is returned
func (self *Service[T]) DownloadAndExtractAsset(owner string, repository string, assetName string, regex bool) (directoryPath string, err error) {
	var (
		asset          *github.ReleaseAsset
		downloadedFile *os.File
		downloadTo     string
		extractTo      string
		subDir         string
		log            *slog.Logger   = self.log
		ghs            *gh.Repository = self.store
		dir            string         = self.GetDirectory()
	)
	log = log.With("operation", "DownloadAndExtractAsset",
		"repository", repository,
		"assetName", assetName,
		"owner", owner)

	// download to a sub directory and use fixed paths underneath that
	// as assetName can be a regex pattern
	subDir, _ = os.MkdirTemp(dir, "*")
	dir = filepath.Join(dir, subDir)
	downloadTo = filepath.Join(dir, "downloaded")
	extractTo = filepath.Join(dir, "extract")

	// get the latest relase and the asset details that match the name
	log.Debug("Downloading the latest release asset ...")
	asset, err = ghs.GetLatestReleaseAsset(owner, repository, assetName, regex)
	if err != nil {
		log.Error("error getting latest release asset", "err", err.Error())
		return
	}
	// download this asset
	log.With("assetID", *asset.ID).Debug("downloading the latest release asset via repository")
	downloadedFile, err = ghs.DownloadReleaseAsset(owner, repository, *asset.ID, downloadTo)
	if err != nil {
		log.Error("error downloading the release asset", "err", err.Error())
		return
	}
	defer downloadedFile.Close()

	// now extract the tar.gz
	log.With("downloadedFile", downloadedFile, "extractTo", extractTo).Debug("extracting downloaded file ...")
	err = utils.TarGzExtract(extractTo, downloadedFile)
	if err != nil {
		log.Error("error extracting the asset", "err", err.Error())
		return
	}

	return extractTo, nil
}

// NewService returns a configured opgmetadata service object
func NewService[T interfaces.Model](ctx context.Context, log *slog.Logger, conf *config.Config, store *gh.Repository) (srv *Service[T], err error) {
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

	srv = &Service[T]{
		ctx:       ctx,
		log:       log.With("service", "opgmetadata"),
		conf:      conf,
		store:     store,
		directory: "",
	}
	return
}

// Default generates the default gh repository and then the service
func Default[T interfaces.Model](ctx context.Context, log *slog.Logger, conf *config.Config) (srv *Service[T]) {

	store, err := gh.New(ctx, log, conf)
	if err != nil {
		log.Error("error creating github repository for opgmetadata service", "error", err.Error())
		return nil
	}
	srv, err = NewService[T](ctx, log, conf, store)
	if err != nil {
		log.Error("error creating opgmetadata service", "error", err.Error())
		return nil
	}

	return
}
