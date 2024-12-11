package lib

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/ministryofjustice/opg-reports/internal/dbs"
	"github.com/ministryofjustice/opg-reports/internal/dbs/adaptors"
	"github.com/ministryofjustice/opg-reports/internal/dbs/crud"
	"github.com/ministryofjustice/opg-reports/internal/fakerextensions/fakerextras"
	"github.com/ministryofjustice/opg-reports/internal/fakerextensions/fakermany"
	"github.com/ministryofjustice/opg-reports/internal/structs"
	"github.com/ministryofjustice/opg-reports/models"
)

func Test_processAwsUptime(t *testing.T) {
	var (
		adaptor dbs.Adaptor
		err     error
		res     any
		ok      bool
		ctx     = context.Background()
		// dir        string = "./"
		// sourceFile string = "../../convertor/converted/aws_uptime.json"
		dir        string = t.TempDir()
		sourceFile string = filepath.Join(dir, "data.json")
		dbFile     string = filepath.Join(dir, "test.db")
		units      []*models.Unit
		accounts   []*models.AwsAccount
		uptime     []*models.AwsUptime
		result     []*models.AwsUptime
	)

	fakerextras.AddProviders()

	adaptor, err = adaptors.NewSqlite(dbFile, false)
	if err != nil {
		t.Fatalf(err.Error())
	}
	defer adaptor.DB().Close()

	// make some fake units
	units = fakermany.Fake[*models.Unit](3)
	// some fake accounts
	accounts = fakermany.Fake[*models.AwsAccount](3)
	// join them up
	for i, ac := range accounts {
		ac.Unit = (*models.UnitForeignKey)(units[i])
	}
	// some fake uptimes
	uptime = fakermany.Fake[*models.AwsUptime](3)
	// join the accounts and units
	for i, up := range uptime {
		acc := accounts[i]
		up.AwsAccount = (*models.AwsAccountForeignKey)(acc)
		up.Unit = acc.Unit
	}

	structs.ToFile(uptime, sourceFile)
	// bootstrap the database - this will now recreate the standards table
	err = crud.Bootstrap(ctx, adaptor, models.Full()...)
	if err != nil {
		return
	}

	res, err = processAwsUptime(ctx, adaptor, sourceFile)
	if err != nil {
		t.Fatalf(err.Error())
	}

	if result, ok = res.([]*models.AwsUptime); !ok {
		t.Errorf("failed to change result to type")
	}

	if len(uptime) != len(result) {
		t.Errorf("number of returned results dont match originals")
	}

}
