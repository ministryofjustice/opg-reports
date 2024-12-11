package crud

import (
	"context"
	"fmt"

	"github.com/ministryofjustice/opg-reports/internal/dbs"
)

// errors
var (
	ErrAdaptorNotWritable         error = fmt.Errorf("adaptor is not configured for writing.")
	ErrNoSetupItems               error = fmt.Errorf("no setupItems passed.")
	ErrModelIsNotCorrectInterface error = fmt.Errorf("model does not implement dbs.CreateableTable")
)

// Bootstrap iterates over each setupItem (of type dbs.CreateableTable) and
// runs crud.CreateTable and crud.CreateIndexes.
//
// Returns an error when:
//   - the adaptor is not in readwrite mode
//   - there are no setupItems are passed
//   - a setupItem does not implement dbs.CreateableTable
//   - CreateTable or CreateIndexes returns and error
func Bootstrap[A dbs.Adaptor](ctx context.Context, adaptor A, setupItems ...interface{}) (err error) {

	if !adaptor.Mode().Write() {
		err = ErrAdaptorNotWritable
		return
	}

	if len(setupItems) <= 0 {
		err = ErrNoSetupItems
		return
	}
	// loop over each and create table and index
	// checking it implments CreateableTable
	for _, setup := range setupItems {
		ct, ok := setup.(dbs.CreateableTable)

		if !ok {
			err = ErrModelIsNotCorrectInterface
			return
		}
		// create table
		_, err = CreateTable(ctx, adaptor, ct)
		if err != nil {
			return
		}
		// create indexes
		_, err = CreateIndexes(ctx, adaptor, ct)
		if err != nil {
			return
		}
	}

	return
}
