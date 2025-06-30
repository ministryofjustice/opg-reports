package awss3

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/ministryofjustice/opg-reports/report/config"
	"github.com/ministryofjustice/opg-reports/report/internal/interfaces"
	"github.com/ministryofjustice/opg-reports/report/internal/repository/awsr"
	"github.com/ministryofjustice/opg-reports/report/internal/utils"
)

// Service is used to download, covnert and return data files from within s3 buckets.
type Service[T interfaces.Model] struct {
	ctx       context.Context
	log       *slog.Logger
	conf      *config.Config
	store     awsr.S3er
	directory string
}

// SetDirectory changes the download / storage location
func (self *Service[T]) SetDirectory(dir string) {
	self.directory = dir
}

// GetDirectory
func (self *Service[T]) GetDirectory() string {
	if self.directory == "" {
		dir, _ := os.MkdirTemp("./", "__download_s3bucket_*")
		self.directory = dir
		defer os.RemoveAll(dir)
	}
	return self.directory
}

// Close cleans up
func (self *Service[T]) Close() (err error) {
	err = os.RemoveAll(self.GetDirectory())
	return
}

// Download uses the bucket & prefix name to fetch a list of all files stored
// and then downloads each file into a local directory, return a list of all
// local filepaths.
//
// If the number of downloaded files does not match the number listed in the
// bucket an error is returned.
func (self *Service[T]) Download(bucket string, prefix string) (downloaded []string, err error) {
	var log *slog.Logger = self.log.With("operation", "Download", "bucket", bucket, "prefix", prefix)

	downloaded, err = self.store.DownloadBucket(bucket, prefix, self.GetDirectory())
	if err != nil {
		return
	}

	log.With("downloadCount", len(downloaded)).Debug("downloaded")
	return
}

// DownloadAndReturnData downloads (via `.Download`) all files locally and then
// reads `.json` files, coverting the data into a slice T
//
// Warning: for large number of files, this can be very memory intensive
func (self *Service[T]) DownloadAndReturnData(bucket string, prefix string) (data []T, err error) {
	var (
		downloadedFiles []string     = []string{}
		log             *slog.Logger = self.log.With("operation", "DownloadAndReturnData", "bucket", bucket, "prefix", prefix)
	)

	data = []T{}

	downloadedFiles, err = self.Download(bucket, prefix)
	if err != nil {
		return
	}
	log.Debug("downloaded files, merging and converting to T ...")

	// each file contains a list of many T's
	for _, file := range downloadedFiles {
		var list = []T{}
		// only parse json filess
		if strings.HasSuffix(file, ".json") {
			log.With("file", file).Debug("reading file into struct")
			err = utils.StructFromJsonFile(file, &list)
			if err != nil {
				return
			}
			// append the file data into the existing data
			data = append(data, list...)
		}
	}
	log.With("count", len(data)).Debug("downloaded and converted")

	return
}

// NewService returns a configured s3 service object
func NewService[T interfaces.Model](ctx context.Context, log *slog.Logger, conf *config.Config, store awsr.S3er) (srv *Service[T], err error) {
	if log == nil {
		return nil, fmt.Errorf("no logger passed for s3 service")
	}
	if conf == nil {
		return nil, fmt.Errorf("no config passed for s3 service")
	}
	if conf.Aws == nil ||
		conf.Aws.Session == nil {
		return nil, fmt.Errorf("missing aws config details for s3 service")
	}
	if conf.Aws.Region == "" ||
		conf.Aws.Session.Token == "" {
		return nil, fmt.Errorf("missing aws config details for s3 service")
	}

	if store == nil {
		return nil, fmt.Errorf("no repository passed for s3 service")
	}

	srv = &Service[T]{
		ctx:   ctx,
		log:   log.With("service", "s3"),
		conf:  conf,
		store: store,
	}
	return
}

// Default generates the default gh repository and then the service
func Default[T interfaces.Model](ctx context.Context, log *slog.Logger, conf *config.Config) (srv *Service[T]) {

	store, err := awsr.New(ctx, log, conf)
	if err != nil {
		log.Error("error creating s3bucket repository for s3 service", "error", err.Error())
		return nil
	}
	srv, err = NewService[T](ctx, log, conf, store)
	if err != nil {
		log.Error("error creating s3 service", "error", err.Error())
		return nil
	}

	return
}
