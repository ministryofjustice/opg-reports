package lib

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"path/filepath"
	"strings"
	"sync"

	"github.com/jmoiron/sqlx"
	"github.com/ministryofjustice/opg-reports/pkg/convert"
	"github.com/ministryofjustice/opg-reports/pkg/datastore"
	"github.com/ministryofjustice/opg-reports/sources/costs"
	"github.com/ministryofjustice/opg-reports/sources/costs/costsdb"
)

type Arguments struct {
	Directory string
	Type      string
}

type dbInfo struct {
	File   string
	Create func(ctx context.Context, dbFilepath string) (db *sqlx.DB, isNew bool, err error)
}

var dbs = map[string]*dbInfo{
	"costs": {File: "./databases/costs.db", Create: costs.CreateNewDB},
}

// SetupArgs setup flag args
func SetupArgs(args *Arguments) {
	flag.StringVar(&args.Directory, "directory", "", "Directory to fetch data files from.")
	flag.StringVar(&args.Type, "type", "costs", "Type of import to do.")
	flag.Parse()
}

// ValidateArgs make sure args are set as planned
func ValidateArgs(args *Arguments) (err error) {
	failOnEmpty := map[string]string{
		"directory": args.Directory,
		"type":      args.Type,
	}
	for k, v := range failOnEmpty {
		if v == "" {
			err = errors.Join(err, fmt.Errorf("%s", k))
		}
	}
	if err != nil {
		err = fmt.Errorf("missing arguments: [%s]", strings.ReplaceAll(err.Error(), "\n", ", "))
	}
	return
}

// GetDatabase creates and returns the db pointer
func GetDatabase(ctx context.Context, args *Arguments) (db *sqlx.DB, err error) {

	info, ok := dbs[args.Type]
	if !ok {
		err = fmt.Errorf("no database setup available for [%s]", args.Type)
		return
	}
	// Runs the create method, but disable seeding
	db, _, err = info.Create(ctx, info.File)
	if err != nil {
		return
	}

	return
}

// ProcessDataFile operates in a go func to load all content in a json datafile
// and the insert this into the active database
func ProcessDataFile(ctx context.Context, wg *sync.WaitGroup, db *sqlx.DB, args *Arguments, datafilepath string) (err error) {
	var base string = filepath.Base(datafilepath)

	slog.Info("[lib.ProcessDataFile] starting ...", slog.String("file", base))
	switch args.Type {
	case "costs":
		records, e := convert.UnmarshalFile[[]*costs.Cost](datafilepath)
		slog.Info("[lib.ProcessDataFile] ", slog.String("file", base), slog.Int("count", len(records)))
		if e == nil {
			_, err = datastore.InsertMany(ctx, db, costsdb.InsertCosts, records)
		}
	default:
		err = fmt.Errorf("err: unsupported type [%s]", args.Type)
	}
	slog.Info("[lib.ProcessDataFile] done ...", slog.String("file", base))
	wg.Done()
	return
}
