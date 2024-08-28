package github_standards_test

import (
	"context"
	"log"
	"log/slog"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ministryofjustice/opg-reports/commands/seed/seeder"
	ghapi "github.com/ministryofjustice/opg-reports/servers/api/github_standards"
	"github.com/ministryofjustice/opg-reports/servers/front/github_standards"
	"github.com/ministryofjustice/opg-reports/servers/shared/srvr/api"
	"github.com/ministryofjustice/opg-reports/servers/shared/srvr/front"
	"github.com/ministryofjustice/opg-reports/servers/shared/srvr/front/config"
	"github.com/ministryofjustice/opg-reports/servers/shared/srvr/front/config/nav"
	"github.com/ministryofjustice/opg-reports/servers/shared/srvr/front/template"
	"github.com/ministryofjustice/opg-reports/servers/shared/srvr/httphandler"
	"github.com/ministryofjustice/opg-reports/shared/convert"
	"github.com/ministryofjustice/opg-reports/shared/logger"
	"github.com/ministryofjustice/opg-reports/shared/testhelpers"
)

const realSchema string = "../../../datastore/github_standards/github_standards.sql"
const templateDir string = "../templates"

func TestServersFrontGithubStandards(t *testing.T) {
	logger.LogSetup()
	// seed
	ctx := context.TODO()
	N := 10
	dir := t.TempDir()
	dbF := filepath.Join(dir, "ghs.db")
	schemaF := filepath.Join(dir, "ghs.sql")
	dataF := filepath.Join(dir, "dummy.json")
	testhelpers.CopyFile(realSchema, schemaF)
	db, err := seeder.Seed(ctx, dbF, schemaF, dataF, "github_standards", N)
	if err != nil {
		slog.Error(err.Error())
		log.Fatal(err.Error())
	}
	defer db.Close()

	apiServer := api.New(ctx, dbF)
	apihandler := api.Wrap(apiServer, ghapi.ListHandler)
	// -- setup a mock api thats bound to the correct handler func
	mockApi := testhelpers.MockServer(apihandler, "warn")
	defer mockApi.Close()

	// -- mock local server that calls the local api
	templates := template.GetTemplates(templateDir)
	// cfg := config.Config
	cfg := &config.Config{Organisation: "TEST RESPONSE"}
	navItem := &nav.Nav{
		Name:     "test nav",
		Uri:      "/",
		Template: "github-standards",
		DataSources: map[string]string{
			"list": mockApi.URL,
		},
	}

	server := front.New(ctx, cfg, templates)
	handler := front.Wrap(server, navItem, github_standards.ListHandler)

	mockFront := testhelpers.MockServer(handler, "warn")
	defer mockFront.Close()
	resp, err := httphandler.Get("", "", mockFront.URL)
	if err != nil {
		t.Errorf("error getting url: %s", err.Error())
	}

	str, _ := convert.Stringify(resp.Response)
	// now look in the string for expected data
	title := "<title>test nav - TEST RESPONSE - Reports</title>"

	if !strings.Contains(str, title) {
		t.Errorf("expected to find known title, did not")
	}
}
