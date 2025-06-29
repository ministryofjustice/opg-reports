package awscost

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/ministryofjustice/opg-reports/report/config"
	"github.com/ministryofjustice/opg-reports/report/internal/awsaccount"
	"github.com/ministryofjustice/opg-reports/report/internal/interfaces"
	"github.com/ministryofjustice/opg-reports/report/internal/sqldb"
	"github.com/ministryofjustice/opg-reports/report/internal/team"
	"github.com/ministryofjustice/opg-reports/report/internal/utils"
)

var _ interfaces.Model = &AwsCostImport{}
var _ interfaces.S3Service[*AwsCostImport] = &testS3Srv[*AwsCostImport]{}

// testS3Srv is a mock service to allow testing of the "Existing" function
// by returning a preset version of data to be inserted without need of
// fetching data from remote etc
type testS3Srv[T interfaces.Model] struct {
	directory string
}

func (self *testS3Srv[T]) SetDirectory(dir string) {
	self.directory = dir
}

func (self *testS3Srv[T]) Close() (err error) {
	return
}

// Download creates a single dummy file from some test data
func (self *testS3Srv[T]) Download(bucket string, prefix string) (downloaded []string, err error) {
	var content string = `[
  {
    "id": 0,
    "ts": "2025-06-15T08:06:29Z",
    "organisation": "ORG",
    "account_id": "001A",
    "account_name": "Shared",
    "unit": "ORG",
    "label": "shared",
    "environment": "production",
    "service": "AWS Backup",
    "region": "NoRegion",
    "date": "2025-05-01",
    "cost": "-0.2059921617"
  },
  {
    "id": 0,
    "ts": "2025-06-15T08:06:29Z",
    "organisation": "ORG",
    "account_id": "001A",
    "account_name": "Shared",
    "unit": "ORG",
    "label": "shared",
    "environment": "production",
    "service": "AWS Backup",
    "region": "eu-west-1",
    "date": "2025-05-01",
    "cost": "0.8956180939"
  },
  {
    "id": 0,
    "ts": "2025-06-15T08:06:29Z",
    "organisation": "ORG",
    "account_id": "001A",
    "account_name": "Shared",
    "unit": "ORG",
    "label": "shared",
    "environment": "production",
    "service": "AWS CloudTrail",
    "region": "NoRegion",
    "date": "2025-05-01",
    "cost": "-0.274675545"
  }
]
`

	filename := filepath.Join(self.directory, "sample-file.json")
	os.WriteFile(filename, []byte(content), os.ModePerm)
	downloaded = []string{filename}
	return
}

// DownloadAndReturn does nothing for this test
func (self *testS3Srv[T]) DownloadAndReturnData(bucket string, prefix string) (data []T, err error) {
	return []T{}, nil
}

func TestAwsCostExisting(t *testing.T) {
	var (
		err       error
		dir       = t.TempDir()
		ctx       = t.Context()
		log       = utils.Logger("ERROR", "TEXT")
		conf      = config.NewConfig()
		srv       = &testS3Srv[*AwsCostImport]{}
		shouldGet = 3
	)
	conf.Database.Path = fmt.Sprintf("%s/%s", dir, "test-awscost-existing.db")
	// seed data
	team.Seed(ctx, log, conf, nil)
	awsaccount.Seed(ctx, log, conf, nil)
	// set the  temp dir
	srv.SetDirectory(dir)

	err = Existing(ctx, log, conf, srv)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}

	// check whats been entered
	store, _ := sqldb.New[*AwsCost](ctx, log, conf)
	tSrv, _ := NewService(ctx, log, conf, store)

	all, err := tSrv.GetAll()
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	if len(all) != shouldGet {
		t.Errorf("mismatch number of records inserted: [%d] [%d]", shouldGet, len(all))
	}

}
