package adaptors

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/ministryofjustice/opg-reports/internal/dbs"
)

// SqlxDB provides methods to get and close sqlx database connections
//
// Implements dbs.DBer
type SqlxDB struct {
	db *sqlx.DB
}

// GetDB returns a db pointer for use in queries etc
func (self *SqlxDB) Get(ctx context.Context, connector dbs.Connector) (db *sqlx.DB, err error) {
	var (
		driver = connector.DriverName()
		source = connector.String()
	)

	if self.db == nil {
		self.db, err = sqlx.ConnectContext(ctx, driver, source)
	}
	db = self.db
	return
}

// CloseDB will close the current database connection if its present
func (self *SqlxDB) Close() (err error) {
	if self.db != nil {
		err = self.db.Close()
	}
	self.db = nil
	return
}
