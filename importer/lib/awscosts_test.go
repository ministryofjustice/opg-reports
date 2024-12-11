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

func Test_processAwsCosts(t *testing.T) {
	var (
		adaptor dbs.Adaptor
		err     error
		res     any
		ok      bool
		ctx     context.Context = context.Background()
		// dir        string          = "./"
		// sourceFile string          = "../../convertor/converted/aws_costs.json"
		dir        string = t.TempDir()
		sourceFile string = filepath.Join(dir, "data.json")
		dbFile     string = filepath.Join(dir, "test.db")
		units      []*models.Unit
		accounts   []*models.AwsAccount
		costs      []*models.AwsCost
		result     []*models.AwsCost
	)
	// structs.UnmarshalFile(sourceFile, &costs)

	fakerextras.AddProviders()

	adaptor, err = adaptors.NewSqlite(dbFile, false)
	if err != nil {
		t.Fatalf(err.Error())
	}
	defer adaptor.DB().Close()
	err = crud.Bootstrap(ctx, adaptor, models.Full()...)
	if err != nil {
		t.Fatalf(err.Error())
	}

	// make some fake units
	units = fakermany.Fake[*models.Unit](3)
	// some fake accounts
	accounts = fakermany.Fake[*models.AwsAccount](3)
	// join them up
	for i, ac := range accounts {
		ac.Unit = (*models.UnitForeignKey)(units[i])
	}
	// some fake uptimes
	costs = fakermany.Fake[*models.AwsCost](3)
	// join the accounts and units
	for i, c := range costs {
		acc := accounts[i]
		c.AwsAccount = (*models.AwsAccountForeignKey)(acc)
		c.Unit = acc.Unit
	}

	structs.ToFile(costs, sourceFile)

	res, err = processAwsCosts(ctx, adaptor, sourceFile)
	if err != nil {
		t.Fatalf(err.Error())
	}

	if result, ok = res.([]*models.AwsCost); !ok {
		t.Errorf("failed to change result to type")
	}

	if len(costs) != len(result) {
		t.Errorf("number of returned results dont match originals - expected [%d] actual [%v]", len(costs), len(result))
	}

}
