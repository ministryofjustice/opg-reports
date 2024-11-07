package datastore_test

import (
	"context"
	"database/sql"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/ministryofjustice/opg-reports/pkg/convert"
	"github.com/ministryofjustice/opg-reports/pkg/datastore"
	"github.com/ministryofjustice/opg-reports/pkg/exfaker"
	"github.com/ministryofjustice/opg-reports/pkg/record"
)

var transactionOptions *sql.TxOptions = &sql.TxOptions{ReadOnly: false, Isolation: sql.LevelDefault}

type address struct {
	ID     int    `json:"id,omitempty" db:"id" faker:"unique, boundary_start=1, boundary_end=2000000" doc:"Database primary key."` // ID is a generated primary key
	Street string `json:"street,omitempty" db:"street" faker:"word"`
}

func (self *address) New() record.Record {
	return &address{}
}
func (self *address) UID() string {
	return fmt.Sprintf("%s-%d", "releases", self.ID)
}
func (self *address) SetID(id int) {
	self.ID = id
}

type address_person struct {
	ID        int `json:"id,omitempty" db:"id" faker:"unique, boundary_start=1, boundary_end=2000000" doc:"Database primary key."`
	PersonID  int `json:"person_id,omitempty" db:"person_id"`
	AddressID int `json:"address_id,omitempty" db:"address_id"`
}

func (self *address_person) New() record.Record {
	return &address_person{}
}
func (self *address_person) UID() string {
	return fmt.Sprintf("%s-%d", "releases", self.ID)
}
func (self *address_person) SetID(id int) {
	self.ID = id
}

type person struct {
	ID        int        `json:"id,omitempty" db:"id" faker:"unique, boundary_start=1, boundary_end=2000000" doc:"Database primary key."` // ID is a generated primary key
	Name      string     `json:"name,omitempty" db:"name" faker:"first_name"`
	Addresses []*address `json:"addresses,omitempty" db:"-" faker:"slice_len=2" doc:"pulled from a many to many join table"`
}

func (self *person) New() record.Record {
	return &person{}
}
func (self *person) UID() string {
	return fmt.Sprintf("%s-%d", "releases", self.ID)
}

func (self *person) SetID(id int) {
	self.ID = id
}

func (self *person) Addrs(ctx context.Context, db *sqlx.DB) (addrs []*address, err error) {
	var stmt datastore.NamedSelectStatement = `
SELECT
	addresses.id as id,
	addresses.street as street
FROM addresses
LEFT JOIN addresses_people ON addresses.id = addresses_people.address_id
WHERE
	addresses_people.person_id = :id
`

	addrs, err = datastore.Select[*address](ctx, db, stmt, self)

	return
}

func (self *person) ProcessJoins(ctx context.Context, db *sqlx.DB, tx *sqlx.Tx) (err error) {
	var transaction *sqlx.Tx = tx
	// create own txn if we havent got one
	if tx == nil {
		transaction = db.MustBeginTx(ctx, transactionOptions)
	}
	// for each address attached, deal with finding, inserting and
	// joining rows between the tables
	for _, addr := range self.Addresses {
		var addressID int = 0
		var joinID int = 0
		var er error

		// try to find the address
		addressID, er = datastore.Get[int](ctx, db, stmtAddressSelect, addr.Street)
		// if theres an error, and its not missing rows.. return the error
		if er != nil && er != sql.ErrNoRows {
			return er
		}
		// if address is not found, we create one
		if addressID == 0 {
			addressID, er = datastore.InsertOne(ctx, db, stmtInsertAddress, addr, transaction)

		}
		// now we find the join
		joinID, er = datastore.Get[int](ctx, db, stmtJoinSelect, self.ID, addressID)
		// if theres an error, and its not missing rows.. return the error
		if er != nil && er != sql.ErrNoRows {
			return er
		}
		if joinID == 0 {
			join := &address_person{PersonID: self.ID, AddressID: addressID}
			joinID, er = datastore.InsertOne(ctx, db, stmtInsertJoin, join, transaction)
		}

		// last error handler
		if er != nil && er != sql.ErrNoRows {
			return er
		}

	}
	// if we used our own tx, then commit
	if tx == nil {
		err = transaction.Commit()
	}

	return
}

var stmtCreateDBs = []datastore.CreateStatement{
	`CREATE TABLE IF NOT EXISTS people (id INTEGER PRIMARY KEY, name TEXT NOT NULL) STRICT;`,
	`CREATE TABLE IF NOT EXISTS addresses (id INTEGER PRIMARY KEY, street TEXT NOT NULL) STRICT;`,
	`CREATE TABLE IF NOT EXISTS addresses_people (id INTEGER PRIMARY KEY, person_id INTEGER NOT NULL, address_id INTEGER NOT NULL) STRICT;`,
}

var (
	stmtInsertPerson  datastore.InsertStatement = `INSERT INTO people (name) VALUES (:name) RETURNING id;`
	stmtInsertAddress datastore.InsertStatement = `INSERT INTO addresses (street) VALUES (:street) RETURNING id;`
	stmtInsertJoin    datastore.InsertStatement = `INSERT INTO addresses_people (person_id, address_id) VALUES (:person_id, :address_id) RETURNING id;`
)
var (
	stmtCountPeople datastore.SelectStatement = `SELECT count(*) as row_count FROM people LIMIT 1`
	stmtCountJoins  datastore.SelectStatement = `SELECT count(*) as row_count FROM addresses_people LIMIT 1`
)

var (
	stmtAddressSelect datastore.SelectStatement = `SELECT id FROM addresses WHERE street = ? LIMIT 1`
	stmtJoinSelect    datastore.SelectStatement = `SELECT id FROM addresses_people WHERE person_id = ? AND address_id = ? LIMIT 1`
)

// make sure that person meets both needs
var _ record.JoinedRecord = &person{}
var _ record.Record = &person{}

func TestDatastoreJoins(t *testing.T) {
	var err error
	var db *sqlx.DB
	var dir string = t.TempDir()
	var dbFile string = filepath.Join(dir, "test.db")
	var ctx context.Context = context.Background()
	var n int = 2
	var people []*person = exfaker.Many[*person](n)

	db, _, err = datastore.NewDB(ctx, datastore.Sqlite, dbFile)
	if err != nil {
		t.Errorf("error [%s]", err.Error())
	}
	defer db.Close()
	datastore.Create(ctx, db, stmtCreateDBs)

	// -- insert generated people
	// insert many should trigger the JoinMany as well
	ids, err := datastore.InsertMany(ctx, db, stmtInsertPerson, people)
	if err != nil {
		t.Errorf("error [%s]", err.Error())
	}
	// make sure id cound matches
	if len(ids) != n {
		t.Errorf("failed to insert rows: [%v]", len(ids))
	}
	// make sure a fresh db call shows the same number of people
	count, err := datastore.Get[int](ctx, db, stmtCountPeople)
	if err != nil {
		t.Errorf("error [%s]", err.Error())
	}
	if count != n {
		t.Errorf("failed to create rows: expected [%d] actual [%v] ", n, count)
	}
	// now check the join numbers are correct
	streets := map[string]int{}
	for _, p := range people {
		for _, a := range p.Addresses {
			streets[a.Street] = 1
		}
	}
	expectedJoins := len(streets)
	joinCount, err := datastore.Get[int](ctx, db, stmtCountJoins)
	if err != nil {
		t.Errorf("error [%s]", err.Error())
	}
	if joinCount != expectedJoins {
		t.Errorf("join count mismatch - expected [%d] actual [%v]", expectedJoins, joinCount)
	}

	// now we want to check join accuracy
	for _, p := range people {
		ogAddrs := p.Addresses
		fetchedAddrs, err := p.Addrs(ctx, db)
		if err != nil {
			t.Errorf("error [%s]", err.Error())
		}
		// length match
		if len(ogAddrs) != len(fetchedAddrs) {
			t.Errorf("number of addresses did not match")
		}
		// content match
		for _, addr := range ogAddrs {
			found := false
			for _, faddr := range fetchedAddrs {
				if faddr.Street == addr.Street {
					found = true
				}
			}
			if !found {
				t.Errorf("addresses did not match")
				convert.PrettyPrint(ogAddrs)
				convert.PrettyPrint(fetchedAddrs)
			}
		}

	}

}
