package crud_test

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/ministryofjustice/opg-reports/internal/dbs/adaptors"
	"github.com/ministryofjustice/opg-reports/internal/dbs/crud"
	"github.com/ministryofjustice/opg-reports/models"
)

func TestCRUDBootstrap(t *testing.T) {
	path := filepath.Join(t.TempDir(), "test.db")
	ctx := context.Background()
	ad, _ := adaptors.NewSqlite(path, false)
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

func TestCRUDBootstrapAll(t *testing.T) {
	all := models.All()
	path := filepath.Join(t.TempDir(), "test.db")
	ctx := context.Background()
	ad, _ := adaptors.NewSqlite(path, false)

	err := crud.Bootstrap(ctx, ad, all...)
	if err != nil {
		t.Errorf("unexpected boot strap error [%s]", err.Error())
	}

	tablesql := `SELECT name FROM sqlite_master WHERE type IN ('table','view') AND name NOT LIKE 'sqlite_%'`
	tables, err := crud.Select[string](ctx, ad, tablesql, nil)
	if err != nil {
		t.Errorf("unexpected boot strap error [%s]", err.Error())
	}

	if len(all) != len(tables) {
		t.Errorf("incorrect number of tables created.")
	}
}
