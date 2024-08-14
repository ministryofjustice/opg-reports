package seeder

import (
	"database/sql"
)

type insertF func(fileContent []byte, db *sql.DB) error

var inserts map[string]insertF = map[string]insertF{
	"github_standards": func(fileContent []byte, db *sql.DB) (err error) {

		return
	},
}
