package existing

import (
	"path/filepath"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/ministryofjustice/opg-reports/report/config"
	"github.com/ministryofjustice/opg-reports/report/internal/repository/awsr"
	"github.com/ministryofjustice/opg-reports/report/internal/repository/sqlr"
	"github.com/ministryofjustice/opg-reports/report/internal/service/awsaccount"
	"github.com/ministryofjustice/opg-reports/report/internal/service/team"
	"github.com/ministryofjustice/opg-reports/report/internal/utils"
)

// TODO - swap to mocked service!
func TestAwsCostsInsert(t *testing.T) {
	var (
		err  error
		dir  string = t.TempDir()
		ctx         = t.Context()
		conf        = config.NewConfig()
		log         = utils.Logger("INFO", "TEXT")
	)
	if conf.Aws.GetToken() == "" {
		t.Skip("No AWS_SESSION_TOKEN, skipping test")
	}
	// set config values
	conf.Database.Path = filepath.Join(dir, "./existing-awscosts.db")
	// TODO - swap to new method when in place
	// seed data
	team.Seed(ctx, log, conf, nil)
	awsaccount.Seed(ctx, log, conf, nil)

	awc, _ := awsr.New(ctx, log, conf)
	sq, _ := sqlr.New(ctx, log, conf)
	client, _ := awsr.GetClient[*s3.Client](ctx, "eu-west-1")
	// existing srv
	srv, err := New(ctx, log, conf)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}

	stmts, err := srv.InsertAwsCosts(client, awc, sq)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	if len(stmts) <= 0 {
		t.Errorf("inserts failed")
	}

}
