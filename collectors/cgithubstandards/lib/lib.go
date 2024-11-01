package lib

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/go-github/v62/github"
)

var (
	defOrg  = "ministryofjustice"
	defTeam = "opg"
)

// Arguments represents all the named arguments for this collector
type Arguments struct {
	Organisation string
	Team         string
	OutputFile   string
}

// SetupArgs maps flag values to properies on the arg passed and runs
// flag.Parse to fetch values
func SetupArgs(args *Arguments) {

	flag.StringVar(&args.Organisation, "organisation", defOrg, "organisation slug.")
	flag.StringVar(&args.Team, "team", defTeam, "team slug")
	flag.StringVar(&args.OutputFile, "output", "./data/github_standards.json", "Filepath for the output")

	flag.Parse()
}

// ValidateArgs checks rules and logic for the input arguments
// Make sure some have non empty values and apply default values to others
func ValidateArgs(args *Arguments) (err error) {
	failOnEmpty := map[string]string{
		"output": args.OutputFile,
	}
	for k, v := range failOnEmpty {
		if v == "" {
			err = errors.Join(err, fmt.Errorf("%s", k))
		}
	}
	if err != nil {
		err = fmt.Errorf("missing arguments: [%s]", strings.ReplaceAll(err.Error(), "\n", ", "))
	}

	if args.Organisation == "" {
		args.Organisation = defOrg
	}
	if args.Team == "" {
		args.Team = defTeam
	}

	return
}

// WriteToFile writes the content to the file
func WriteToFile(content []byte, args *Arguments) {
	var (
		filename string
		dir      string = filepath.Dir(args.OutputFile)
	)
	os.MkdirAll(dir, os.ModePerm)
	filename = args.OutputFile

	os.WriteFile(filename, content, os.ModePerm)

}

// AllRepos returns all accessible repos for the details passed
func AllRepos(ctx context.Context, client *github.Client, args *Arguments) (all []*github.Repository, err error) {
	var (
		org             string               = args.Organisation
		team            string               = args.Team
		includeArchived bool                 = false
		list            []*github.Repository = []*github.Repository{}
		page            int                  = 1
	)

	all = []*github.Repository{}

	for page > 0 {
		slog.Info("getting repostiories", slog.Int("page", page))
		pg, resp, e := client.Teams.ListTeamReposBySlug(ctx, org, team, &github.ListOptions{PerPage: 100, Page: page})
		if e != nil {
			err = e
			return
		}
		list = append(list, pg...)
		page = resp.NextPage
	}

	if !includeArchived {
		for _, r := range list {
			if !*r.Archived {
				all = append(all, r)
			}
		}
	} else {
		all = list
	}

	return
}
