package awsr

import (
	"context"
	"testing"

	"opg-reports/report/config"
	"opg-reports/report/internal/utils"

	"github.com/aws/aws-sdk-go-v2/service/sts"
)

var accountId = "001A"
var accountArn = "arn:aws:iam::001A:user/person/test.name"
var accountUser = "AIDZFKC"

// mockSTSCaller used to avoid making api calls
type mockSTSCaller struct{}

// GetCallerIdentity always returns dummy creds
func (self *mockSTSCaller) GetCallerIdentity(ctx context.Context, params *sts.GetCallerIdentityInput, optFns ...func(*sts.Options)) (*sts.GetCallerIdentityOutput, error) {
	return &sts.GetCallerIdentityOutput{
		Account: &accountId,
		Arn:     &accountArn,
		UserId:  &accountUser,
	}, nil
}

func TestSTSCallerIdentity(t *testing.T) {
	var (
		err        error
		repository *Repository
		ctx        = t.Context()
		conf       = config.NewConfig()
		log        = utils.Logger("ERROR", "TEXT")
	)

	repository, err = New(ctx, log, conf)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
		t.FailNow()
	}
	caller, err := repository.GetCallerIdentity(&mockSTSCaller{})
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	if caller.Account == nil || *(caller.Account) == "" {
		t.Errorf("no caller found")
	}

}
