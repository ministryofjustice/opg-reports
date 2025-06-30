package caller

import (
	"testing"

	"github.com/ministryofjustice/opg-reports/report/config"
	"github.com/ministryofjustice/opg-reports/report/internal/interfaces"
	"github.com/ministryofjustice/opg-reports/report/internal/utils"
)

var _ interfaces.Repository = &Repository{}
var _ interfaces.STSRepository = &Repository{}

func TestGetCallerID(t *testing.T) {
	var (
		err  error
		ctx  = t.Context()
		conf = config.NewConfig()
		log  = utils.Logger("ERROR", "TEXT")
	)

	if conf.Aws.GetToken() == "" {
		t.Skip("No AWS_SESSION_TOKEN, skipping test")
	}

	repo, err := New(ctx, log, conf)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
		t.FailNow()
	}

	account := repo.GetAccountID()
	if account == "" {
		t.Errorf("failed to find account")
	}

}
