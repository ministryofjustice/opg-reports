package team

import (
	"fmt"
	"slices"
	"testing"

	"github.com/ministryofjustice/opg-reports/report/config"
	"github.com/ministryofjustice/opg-reports/report/internal/interfaces"
	"github.com/ministryofjustice/opg-reports/report/internal/repository/sqldb"
	"github.com/ministryofjustice/opg-reports/report/internal/utils"
)

var _ interfaces.Model = &TeamImport{}
var _ interfaces.MetadataService[*TeamImport] = &testMetaSrv[*TeamImport]{}

// testMetaSrv is a mock service to allow testing of the "Existing" function
// by returning a preset version of data to be inserted without need of
// fetching data from remote etc
type testMetaSrv[T interfaces.Model] struct{}

func (self *testMetaSrv[T]) Close() (err error) {
	return
}

// DownloadAndReturn returns a series of fake aws_account data that will generate series
// of team names from the billing_name value
// To do this we'll use the TeamImport struct
func (self *testMetaSrv[T]) DownloadAndReturn(owner string, repository string, assetName string, regex bool, filename string) (data []*TeamImport, err error) {

	data = []*TeamImport{
		{Name: "TeamZ"},
		{Name: "TeamA"},
		{Name: "TeamB"},
		{Name: "TeamZ"},
		{Name: "Team100"},
		{Name: "Team1001"},
		{Name: "TeamA"},
		{Name: "TeamA"},
		{Name: "Team_A"},
	}
	return
}

func TestTeamExisting(t *testing.T) {
	var (
		err  error
		dir  = t.TempDir()
		ctx  = t.Context()
		log  = utils.Logger("ERROR", "TEXT")
		conf = config.NewConfig()
		srv  = &testMetaSrv[*TeamImport]{}
	)
	conf.Database.Path = fmt.Sprintf("%s/%s", dir, "test-team-existing.db")
	// work out number of unique values
	data, _ := srv.DownloadAndReturn("", "", "", false, "")
	uni := []string{}
	for _, d := range data {
		uni = append(uni, d.Name)
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
	store, _ := sqldb.New[*Team](ctx, log, conf)
	tSrv, _ := NewService(ctx, log, conf, store)

	all, _ := tSrv.GetAllTeams()
	if len(uni) != len(all) {
		t.Errorf("mismatch number of records inserted")
	}

}
