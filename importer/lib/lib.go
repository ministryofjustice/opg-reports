package lib

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"strings"

	"github.com/ministryofjustice/opg-reports/internal/dbs"
	"github.com/ministryofjustice/opg-reports/internal/fileutils"
)

// Arguments represents all the named arguments for this collector
type Arguments struct {
	DatabasePath string
	SourceFile   string
	Type         string
}

type importerFunc func(ctx context.Context, adaptor dbs.Adaptor, path string) (any, error)

var TypeProcessors = map[string]importerFunc{
	"github-standards": processGithubStandards,
	"aws-uptime":       processAwsUptime,
}

var (
	defaultDBPath     string = "./api.db"
	defaultSourceFile string = "./data.json"
)

// SetupArgs maps flag values to properies on the arg passed and runs
// flag.Parse to fetch values
func SetupArgs(args *Arguments) {

	flag.StringVar(&args.DatabasePath, "database", defaultDBPath, "location of the database")
	flag.StringVar(&args.SourceFile, "file", defaultSourceFile, "File to import data from")
	flag.StringVar(&args.Type, "type", "github-standards", "Type of data in the source file")

	flag.Parse()
}

func ValidateArgs(args *Arguments) (err error) {
	failOnEmpty := map[string]string{
		"database": args.DatabasePath,
		"file":     args.SourceFile,
		"type":     args.Type,
	}
	for k, v := range failOnEmpty {
		if v == "" {
			err = errors.Join(err, fmt.Errorf("%s", k))
		}
	}
	if err != nil {
		err = fmt.Errorf("missing arguments: [%s]", strings.ReplaceAll(err.Error(), "\n", ", "))
	}

	if _, ok := TypeProcessors[args.Type]; !ok {
		err = errors.Join(fmt.Errorf("invalid type [%s]", args.Type), err)
	}

	if !fileutils.Exists(args.SourceFile) {
		err = errors.Join(fmt.Errorf("file not found [%s]", args.SourceFile), err)
	}

	return
}
