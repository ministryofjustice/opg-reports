package lib_test

import (
	"testing"

	"github.com/ministryofjustice/opg-reports/internal/envar"
	"github.com/ministryofjustice/opg-reports/servers/api/lib"
)

// TestServersApiDownloadS3 pull a known file from the bucket
//   - only runs if the env is setup for it...
func TestServersApiDownloadS3(t *testing.T) {

	var (
		bucketName string = "report-data-development"
		bucketDB   string = "github_standards/github_standards.json"
		localPath  string = "./db/gs.json"
	)

	if envar.Get("AWS_SESSION_TOKEN", "") != "" {
		ok, err := lib.DownloadS3DB(bucketName, bucketDB, localPath)
		if err != nil {
			t.Fatalf(err.Error())
		}
		if !ok {
			t.Errorf("download returned not ok")
		}
	}
}
