// Package govukassets provides utils for downloading and working with govuk front end
package govukassets

import (
	"os"
	"strings"
	"time"

	"github.com/ministryofjustice/opg-reports/pkg/consts"
	"github.com/ministryofjustice/opg-reports/pkg/fileutils"
)

const (
	frontEndLocation string = "https://github.com/alphagov/govuk-frontend/releases/download/{version}/release-{version}.zip"
)

type FrontEndConfig struct {
	Version          string // Version number to use
	Source           string // Source is the formatted string to fetch the release zip from
	DownloadLocation string // DownloadLocation is where the zip was downloaded to
}

// Url generates the url from the version and source strings
func (self *FrontEndConfig) Url() (url string) {

	url = strings.ReplaceAll(self.Source, "{version}", self.Version)
	return
}

// Download gets the zip file and downloads it into a temporary
// directoy and returns the full path to the zip
// This path can then be used with Extract
func (self *FrontEndConfig) Download(timeout time.Duration) (path string, err error) {
	var (
		url     string = self.Url()
		dir     string = os.TempDir()
		zipName string = "assets.zip"
	)
	path, err = fileutils.DownloadFromUrl(url, dir, zipName, timeout)
	self.DownloadLocation = path

	return
}

func (self *FrontEndConfig) Extract(zipFilepath string) (err error) {

	return
}

func FrontEnd() *FrontEndConfig {
	return &FrontEndConfig{
		Version: consts.GovUKFrontendVersion,
		Source:  frontEndLocation,
	}
}
