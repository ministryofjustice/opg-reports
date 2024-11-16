package crud_test

import (
	"context"
	"testing"

	"github.com/ministryofjustice/opg-reports/internal/dbs/adaptors"
	"github.com/ministryofjustice/opg-reports/internal/dbs/crud"
	"github.com/ministryofjustice/opg-reports/models"
)

func TestCRUDBootstrap(t *testing.T) {
	ctx := context.Background()
	ad, _ := adaptors.NewSqlite("./test.db", false)
	err := crud.Bootstrap(
		ctx,
		ad,
		&models.AwsAccount{},
		&models.AwsCost{},
	)
	if err != nil {
		t.Errorf("unexpected boot strap error [%s]", err.Error())
	}

}
