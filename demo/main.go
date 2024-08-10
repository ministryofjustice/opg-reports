package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/ministryofjustice/opg-reports/commands/shared/argument"
	"github.com/ministryofjustice/opg-reports/datastore/github_standards/ghs"
	"github.com/ministryofjustice/opg-reports/shared/fake"
	"github.com/ministryofjustice/opg-reports/shared/logger"
)

func exists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func main() {
	logger.LogSetup()
	group := flag.NewFlagSet("demo", flag.ExitOnError)
	which := argument.New(group, "which", "all", "")
	dir := argument.New(group, "out", ".", "")

	group.Parse(os.Args[1:])

	what := *which.Value

	// only generate data if it doesnt already exist
	if what == "github_standards" || what == "all" {
		counter := 1000000
		owner := fake.String(12)

		d := fmt.Sprintf("%s/github_standards", *dir.Value)
		file := d + "/github_standards.csv"

		if !exists(d) {
			os.MkdirAll(d, os.ModePerm)
			f, _ := os.OpenFile(file, os.O_APPEND|os.O_WRONLY|os.O_CREATE, os.ModePerm)
			defer f.Close()

			for x := 0; x < counter; x++ {
				slog.Debug(fmt.Sprintf("[%d/%d] generating %s", x+1, counter, "github_standards"),
					slog.String("which", what),
				)
				g := ghs.Fake()
				g.Owner = owner
				g.FullName = fmt.Sprintf("%s/%s", owner, g.Name)

				f.WriteString(g.ToCSV())
			}
			slog.Info(fmt.Sprintf("generated [%d] %s", counter, "github_standards"), slog.String("file", file))
		}
	}

}
