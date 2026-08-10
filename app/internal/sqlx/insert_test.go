package sqlx

import (
	"errors"
	"path/filepath"
	"testing"
)

type tInsertComms struct {
	Email string `json:"email"`
	Phone string `json:"phone"`
}
type tInsertPerson struct {
	ID    int           `json:":id"`
	Comms *tInsertComms `json:"comms"`
}

var (
	tCreatePersonTable      string = `CREATE TABLE IF NOT EXISTS test_person_table(id INTEGER PRIMARY KEY, email TEXT NOT NULL, phone TEXT NOT NULL) STRICT;`
	tInsertPersonSQLWorking string = `INSERT INTO test_person_table(email, phone) VALUES(:email, :phone) RETURNING id;`
	tInsertPersonSQLFail    string = `INSERT INTO test_person_table(email, phone) VALUES(:email, :telephone) RETURNING id;`
)

func TestSQLxInsertSimple(t *testing.T) {
	var (
		sq  *Sqlite
		err error
		ctx        = t.Context()
		dir string = t.TempDir()
	)
	// test connecting to the db that works
	sq = NewSQLite(filepath.Join(dir, "test-insert.db"), false)
	_, err = sq.Open()
	if err != nil {
		t.Errorf("unexpected error calling open [%s]", err.Error())
		t.FailNow()
	}
	defer sq.Close()
	_, err = Exec(ctx, sq, tCreatePersonTable)
	if err != nil {
		t.Errorf("unexpected error calling exec [%s]", err.Error())
		t.FailNow()
	}
	// test a working insert
	p := &tInsertPerson{Comms: &tInsertComms{
		Email: "test@example.com",
		Phone: "00441215123456",
	}}
	_, err = Insert(ctx, sq, tInsertPersonSQLWorking, []*tInsertPerson{p})
	if err != nil {
		t.Errorf("unexpected error calling insert [%s]", err.Error())
		t.FailNow()
	}

	// test an insert with incorrect fields
	_, err = Insert(ctx, sq, tInsertPersonSQLFail, []*tInsertPerson{p})
	if err == nil {
		t.Errorf("expected an error to be returned about binding and struct fields.")
		t.FailNow()
	}
	if !errors.Is(err, ErrBindingNoKey) {
		t.Errorf("expected abinding error to be returned.")
	}

	//t.FailNow()
}
