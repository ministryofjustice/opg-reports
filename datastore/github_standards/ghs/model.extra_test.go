package ghs_test

import (
	"os"
	"testing"

	"github.com/ministryofjustice/opg-reports/datastore/github_standards/ghs"
)

func TestDataGithubStandardsCSV(t *testing.T) {
	create := 1000010
	file := "test.csv"
	f, _ := os.OpenFile(file, os.O_APPEND|os.O_WRONLY|os.O_CREATE, os.ModePerm)
	defer f.Close()

	for i := 0; i < create; i++ {
		g := ghs.Fake()
		content := g.ToCSV()
		f.WriteString(content)
	}
}
