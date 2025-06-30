package aws

import (
	"testing"

	"github.com/ministryofjustice/opg-reports/report/config"
	"github.com/ministryofjustice/opg-reports/report/internal/utils"
)

func TestSTSCallerIdentity(t *testing.T) {
	var (
		err        error
		repository *Repository
		ctx        = t.Context()
		conf       = config.NewConfig()
		log        = utils.Logger("ERROR", "TEXT")
	)

	if conf.Aws.GetToken() == "" {
		t.Skip("No AWS_SESSION_TOKEN, skipping test")
	}

	repository, err = New(ctx, log, conf)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
		t.FailNow()
	}
	caller, err := repository.GetCallerIdentity()
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	if caller.Account == nil || *(caller.Account) == "" {
		t.Errorf("no caller found")
	}

}
