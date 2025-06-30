package awsaccount

import (
	"fmt"
	"slices"
	"testing"

	"github.com/ministryofjustice/opg-reports/report/config"
	"github.com/ministryofjustice/opg-reports/report/internal/interfaces"
	"github.com/ministryofjustice/opg-reports/report/internal/repository/sqldb"
	"github.com/ministryofjustice/opg-reports/report/internal/service/team"
	"github.com/ministryofjustice/opg-reports/report/internal/utils"
)

var _ interfaces.Model = &AwsAccountImport{}
var _ interfaces.MetadataService[*AwsAccountImport] = &testMetaSrv[*AwsAccountImport]{}

// testMetaSrv is a mock service to allow testing of the "Existing" function
// by returning a preset version of data to be inserted without need of
// fetching data from remote etc
type testMetaSrv[T interfaces.Model] struct{}

func (self *testMetaSrv[T]) Close() (err error) {
	return
}

// DownloadAndReturn returns a series of fake aws_account data
func (self *testMetaSrv[T]) DownloadAndReturn(owner string, repository string, assetName string, regex bool, filename string) (data []*AwsAccountImport, err error) {

	data = []*AwsAccountImport{
		{AwsAccount: AwsAccount{ID: "001A", Name: "TestAcc01", Label: "label-1", Environment: "production"}, TeamName: "TeamA"},
		{AwsAccount: AwsAccount{ID: "002A", Name: "TestAcc02", Label: "label-2", Environment: "production"}, TeamName: "TeamB"},
		{AwsAccount: AwsAccount{ID: "003A", Name: "TestAcc03", Label: "label-3", Environment: "production"}, TeamName: "TeamC"},
		{AwsAccount: AwsAccount{ID: "004A", Name: "TestAcc04", Label: "label-4", Environment: "production"}, TeamName: "TeamA"},
		{AwsAccount: AwsAccount{ID: "005A", Name: "TestAcc05", Label: "label-5", Environment: "production"}, TeamName: "TeamD"},
		{AwsAccount: AwsAccount{ID: "006A", Name: "TestAcc05", Label: "label-6", Environment: "production"}, TeamName: "TeamA"},
	}
	return
}

func TestAwsAccountExisting(t *testing.T) {
	var (
		err  error
		dir  = t.TempDir()
		ctx  = t.Context()
		log  = utils.Logger("ERROR", "TEXT")
		conf = config.NewConfig()
		srv  = &testMetaSrv[*AwsAccountImport]{}
	)
	conf.Database.Path = fmt.Sprintf("%s/%s", dir, "test-awsaccount-existing.db")
	// seed team data
	team.Seed(ctx, log, conf, nil)

	// work out number of unique values
	data, _ := srv.DownloadAndReturn("", "", "", false, "")
	uni := []string{}
	for _, d := range data {
		uni = append(uni, d.ID)
	}
	slices.Sort(uni)
	uni = slices.Compact(uni)

	// run the import
	err = Existing(ctx, log, conf, srv)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	// now fetch the rows to make sure count matches
	// grab team store setup to check after the insert
	store, _ := sqldb.New[*AwsAccountImport](ctx, log, conf)
	tSrv, _ := NewService(ctx, log, conf, store)

	all, _ := tSrv.GetAllAccounts()
	if len(uni) != len(all) {
		t.Errorf("mismatch number of records inserted")
	}

}
