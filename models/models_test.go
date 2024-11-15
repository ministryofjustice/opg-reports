package models_test

import (
	"context"
	"fmt"

	"github.com/ministryofjustice/opg-reports/internal/dbs"
	"github.com/ministryofjustice/opg-reports/internal/dbs/adaptors"
	"github.com/ministryofjustice/opg-reports/internal/dbs/crud"
)

// testDBbuilder spins up a table and inserts records and indexes into it
func testDBbuilder[T dbs.TableOfRecord](ctx context.Context, adaptor *adaptors.Sqlite, itemType T, insert []T) (results []T, err error) {

	_, err = crud.CreateTable(ctx, adaptor, itemType)
	if err != nil {
		return
	}

	_, err = crud.CreateIndexes(ctx, adaptor, itemType)
	if err != nil {
		return
	}

	results, err = crud.Insert(ctx, adaptor, itemType, insert...)
	if err != nil {
		return
	}

	if len(results) != len(insert) {
		err = fmt.Errorf("inserted record count mistmatch")
	}
	return
}
