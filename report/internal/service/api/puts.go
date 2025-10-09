package api

import (
	"fmt"
	"log/slog"
	"opg-reports/report/internal/repository/sqlr"
)

// Put inserts new records into the table.
//
// Note: Dont expose to the api endpoints
func (self *Service[T]) Put(store sqlr.RepositoryWriter, insertStmt string, data []T) (results []*sqlr.BoundStatement, err error) {
	var (
		dataType string                 = fmt.Sprintf("%T", data)
		inserts  []*sqlr.BoundStatement = []*sqlr.BoundStatement{}
		log      *slog.Logger           = self.log.With("operation", "Put").With("type", dataType)
	)
	results = []*sqlr.BoundStatement{}

	log.Debug("generating insert statements")
	// for each cost item generate the insert
	for _, row := range data {
		inserts = append(inserts, &sqlr.BoundStatement{Data: row, Statement: insertStmt})
	}
	log.With("count", len(inserts)).Debug("inserting records ...")

	// run inserts
	if err = store.Insert(inserts...); err != nil {
		log.Error("error inserting", "err", err.Error())
		return
	}
	// only merge in the items with return values
	for _, in := range inserts {
		if in.Returned != nil {
			results = append(results, in)
		}
	}
	if len(results) != len(data) {
		err = fmt.Errorf("not all records were inserted; expected [%d] actual [%d]", len(data), len(results))
		return
	}

	log.With("inserted", len(results)).Info("inserting records successful")
	return
}
