package lib

import (
	"flag"
	"fmt"
	"path/filepath"

	v1 "github.com/ministryofjustice/opg-reports/convertor/v1"
	"github.com/ministryofjustice/opg-reports/internal/structs"
)

type Arguments struct {
	Type            string
	SourceDirectory string
	DestinationFile string
}

type V1er interface {
	MarshalJSON() (bytes []byte, err error)
}

// SetupArgs maps flag values to properies on the arg passed and runs
// flag.Parse to fetch values
func SetupArgs(args *Arguments) {

	flag.StringVar(&args.DestinationFile, "destination", "", "File to output v2 version of data")
	flag.StringVar(&args.SourceDirectory, "source", "", "Directory containing v1 versions of data")
	flag.StringVar(&args.Type, "type", "github-standards", "Type of data in the source file")

	flag.Parse()
}

func ReadAll[T V1er](directory string) (all []T, err error) {
	var (
		files   []string
		pattern = filepath.Join(directory, "*.json")
	)
	all = []T{}
	files, _ = filepath.Glob(pattern)
	for _, f := range files {
		many := []T{}
		err = structs.UnmarshalFile(f, &many)
		if err != nil {
			return
		}
		all = append(all, many...)
	}

	return
}

func Run(args *Arguments) (err error) {

	switch args.Type {
	case "aws-costs":
		var out []*v1.AwsCost
		out, err = ReadAll[*v1.AwsCost](args.SourceDirectory)
		if err != nil {
			return
		}
		err = structs.ToFile(out, args.DestinationFile)
	case "aws-uptime":
		var out []*v1.AwsUptime
		out, err = ReadAll[*v1.AwsUptime](args.SourceDirectory)
		if err != nil {
			return
		}
		err = structs.ToFile(out, args.DestinationFile)
	case "github-standards":
		var out []*v1.GithubStandard
		out, err = ReadAll[*v1.GithubStandard](args.SourceDirectory)
		if err != nil {
			return
		}
		err = structs.ToFile(out, args.DestinationFile)
	default:
		err = fmt.Errorf("unknown type [%s]", args.Type)
	}

	return
}
