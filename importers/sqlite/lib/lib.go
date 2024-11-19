package lib

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"path/filepath"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/ministryofjustice/opg-reports/pkg/convert"
	"github.com/ministryofjustice/opg-reports/pkg/datastore"
	"github.com/ministryofjustice/opg-reports/pkg/record"
	"github.com/ministryofjustice/opg-reports/sources/costs"
	"github.com/ministryofjustice/opg-reports/sources/costs/costsdb"
	"github.com/ministryofjustice/opg-reports/sources/releases"
	"github.com/ministryofjustice/opg-reports/sources/releases/releasesdb"
	"github.com/ministryofjustice/opg-reports/sources/standards"
	"github.com/ministryofjustice/opg-reports/sources/standards/standardsdb"
	"github.com/ministryofjustice/opg-reports/sources/uptime"
	"github.com/ministryofjustice/opg-reports/sources/uptime/uptimedb"
)

type Arguments struct {
	File     string
	Database string
	Type     string
}

// creatorF is a contstraint of the functions to call to create new DBs
type creatorF func(ctx context.Context, dbFilepath string) (db *sqlx.DB, isNew bool, err error)

// processorF is a type constraint for functions that can process a datafile into a set of db records
type processorF func(ctx context.Context, db *sqlx.DB, stmt datastore.InsertStatement, datafilepath string) (count int, err error)

// known is a struct to capture what we know for the type and how to use it
type known struct {
	CreateDB   creatorF
	InsertStmt datastore.InsertStatement
	Processor  processorF
}

var Known = map[string]known{
	"aws-costs": {
		CreateDB:   costs.CreateNewDB,
		InsertStmt: costsdb.InsertCosts,
		Processor:  processor[*costs.Cost],
	},
	"aws-uptime": {
		CreateDB:   uptime.CreateNewDB,
		InsertStmt: uptimedb.InsertUptime,
		Processor:  processor[*uptime.Uptime],
	},
	"github-standards": {
		CreateDB:   standards.CreateNewDB,
		InsertStmt: standardsdb.InsertStandard,
		Processor:  processor[*standards.Standard],
	},

	"github-releases": {
		CreateDB:   releases.CreateNewDB,
		InsertStmt: releasesdb.InsertRelease,
		Processor:  processor[*releases.Release],
	},
}

var defaultDB string = "./databases/api.db"

// SetupArgs setup flag args
func SetupArgs(args *Arguments) {
	flag.StringVar(&args.File, "file", "", "Path to single file with new data to add into the existing database.")
	flag.StringVar(&args.Type, "type", "", "Type of data to import. allowed: [aws-costs|aws-uptime|github-standards|github-releases]")
	flag.StringVar(&args.Database, "database", defaultDB, "Path to the database")
	flag.Parse()
}

// ValidateArgs make sure args are set as planned
func ValidateArgs(args *Arguments) (err error) {
	failOnEmpty := map[string]string{
		"file": args.File,
		"type": args.Type,
	}
	for k, v := range failOnEmpty {
		if v == "" {
			err = errors.Join(err, fmt.Errorf("%s", k))
		}
	}
	if err != nil {
		err = fmt.Errorf("missing arguments: [%s]", strings.ReplaceAll(err.Error(), "\n", ", "))
	}

	if args.Database == "" {
		args.Database = defaultDB
	}

	if _, ok := Known[args.Type]; !ok {
		err = fmt.Errorf("err: unsupported type [%s]", args.Type)
	}

	return
}

// GetDatabase creates and returns the db pointer
func GetDatabase(ctx context.Context, args *Arguments) (db *sqlx.DB, err error) {
	var ok bool
	var found known
	var path = args.Database
	if path == "" {
		err = fmt.Errorf("no database setup available for [%s]", args.Type)
		return
	}
	if found, ok = Known[args.Type]; !ok {
		err = fmt.Errorf("no database creator setup available for [%s]", args.Type)
		return
	}
	// Runs the create method, but disable seeding
	db, _, err = found.CreateDB(ctx, path)
	return
}

// processor allows generic handling for each known type
func processor[T record.Record](ctx context.Context, db *sqlx.DB, stmt datastore.InsertStatement, datafilepath string) (count int, err error) {
	var records []T
	var ids []int
	var base string = filepath.Base(datafilepath)
	var recordCount int

	records, err = convert.UnmarshalFile[[]T](datafilepath)
	if err != nil {
		return
	}

	recordCount = len(records)

	ids, err = datastore.InsertMany(ctx, db, stmt, records)
	if err != nil {
		return
	}

	count = len(ids)

	// if the number of original records does not match the count of
	// inserted records, add a custom error to flag that
	if count != recordCount {
		err = errors.Join(err, fmt.Errorf("inserted count doesn't match source data count - source [%d] inserted [%v]", recordCount, count))
	}

	slog.Info("[lib.processor] ",
		slog.String("type", fmt.Sprintf("%T", records)),
		slog.String("file", base),
		slog.Int("count", count),
		slog.Int("recordCount", recordCount))

	return

}

// ProcessDataFile operates in a go func to load all content in a json datafile
// and the insert this into the active database
func ProcessDataFile(ctx context.Context, db *sqlx.DB, args *Arguments, datafilepath string) (count int, err error) {
	var ok bool
	var found known
	var base string = filepath.Base(datafilepath)

	slog.Info("[lib.ProcessDataFile] starting ...", slog.String("file", base))

	if found, ok = Known[args.Type]; !ok {
		err = fmt.Errorf("cannot find known config for type [%s]", args.Type)
		return
	}

	if found.InsertStmt == "" {
		err = fmt.Errorf("cannot find insert statement for type [%s]", args.Type)
		return
	}

	if found.Processor == nil {
		err = fmt.Errorf("cannot find processor function for type [%s]", args.Type)
		return
	}

	count, err = found.Processor(ctx, db, found.InsertStmt, datafilepath)

	slog.Info("[lib.ProcessDataFile] done.", slog.String("file", base), slog.Int("count", count))
	return
}
