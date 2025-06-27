package s3

import (
	"testing"

	"github.com/ministryofjustice/opg-reports/report/config"
	"github.com/ministryofjustice/opg-reports/report/internal/s3bucket"
	"github.com/ministryofjustice/opg-reports/report/internal/utils"
)

type testGH struct {
	Name  string `json:"name"`
	Owner string `json:"owner"`
}

// TestS3Service fetches known data from s3 bucket and
// converts it for testing
func TestS3Service(t *testing.T) {
	var (
		err    error
		dir    string = t.TempDir()
		bucket        = "report-data-development"
		prefix        = "github_standards/"
		ctx           = t.Context()
		cfg           = config.NewConfig()
		lg            = utils.Logger("WARN", "TEXT")
	)
	if cfg.Aws.Session.Token == "" {
		t.Skip("No AWS_SESSION_TOKEN, skipping test")
	}

	repo, err := s3bucket.New(ctx, lg, cfg)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}

	srv, err := NewService[*testGH](ctx, lg, cfg, repo)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
		t.FailNow()
	}
	defer srv.Close()

	// set tmp directory
	srv.SetDirectory(dir)

	res, err := srv.DownloadAndReturnData(bucket, prefix)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	utils.Debug(res)

	if len(res) <= 0 {
		t.Errorf("failed to convert data")
		t.FailNow()
	}

	if res[0].Name == "" {
		t.Errorf("possible data conversion issue")
	}

}
