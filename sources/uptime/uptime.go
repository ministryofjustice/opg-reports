// Package uptime provides wrappers for recording uptime
//
// Currently only from aws healthchecks
package uptime

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/ministryofjustice/opg-reports/pkg/datastore"
	"github.com/ministryofjustice/opg-reports/sources/uptime/uptimedb"
)

type Uptime struct {
	ID      int     `json:"id" db:"id" faker:"unique, boundary_start=1, boundary_end=2000000" doc:"Database primary key."`                         // ID is a generated primary key
	Ts      string  `json:"ts,omitempty" db:"ts" faker:"time_string" doc:"Time the record was created."`                                           // TS is timestamp when the record was created
	Unit    string  `json:"unit,omitempty" db:"unit" faker:"oneof: unitA, unitB, unitC" doc:"The name of the unit / team that owns this account."` // Unit is the team that owns this account, passed directly
	Date    string  `json:"date,omitempty" db:"date" faker:"date_string" doc:"Date this cost was generated."`                                      // The interval date for when this uptime was logged
	Average float64 `json:"average" db:"average" doc:"Percentage uptime average for this record period."`
}

// UID returns a unique uid
func (self *Uptime) UID() string {
	return fmt.Sprintf("%s-%d", "uptime", self.ID)
}

const RecordsToSeed int = (1440 * 7) // a week

var insert = uptimedb.InsertUptime
var creates = []datastore.CreateStatement{
	uptimedb.CreateUptimeTable,
	uptimedb.CreateUptimeTableDateIndex,
	uptimedb.CreateUptimeTableUnitDateIndex,
}

// Setup will ensure a database with records exists in the filepath requested.
// If there is no database at that location a new sqlite database will
// be created and populated with series of dummy data - helpful for local testing.
func Setup(ctx context.Context, dbFilepath string, seed bool) {
	datastore.Setup[Uptime](ctx, dbFilepath, insert, creates, seed, RecordsToSeed)
}

// CreateNewDB will create a new DB file and then
// try to run table and index creates
func CreateNewDB(ctx context.Context, dbFilepath string) (db *sqlx.DB, isNew bool, err error) {
	return datastore.CreateNewDB(ctx, dbFilepath, creates)
}
