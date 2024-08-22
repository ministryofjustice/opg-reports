package github_standards_test

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ministryofjustice/opg-reports/commands/seed/seeder"
	ghapi "github.com/ministryofjustice/opg-reports/servers/api/github_standards"
	"github.com/ministryofjustice/opg-reports/servers/front/config"
	"github.com/ministryofjustice/opg-reports/servers/front/config/navigation"
	"github.com/ministryofjustice/opg-reports/servers/front/config/src"
	"github.com/ministryofjustice/opg-reports/servers/front/getter"
	"github.com/ministryofjustice/opg-reports/servers/front/github_standards"
	"github.com/ministryofjustice/opg-reports/servers/front/template_helpers"
	"github.com/ministryofjustice/opg-reports/servers/shared/urls"
	"github.com/ministryofjustice/opg-reports/shared/convert"
	"github.com/ministryofjustice/opg-reports/shared/logger"
	"github.com/ministryofjustice/opg-reports/shared/testhelpers"
)

const realSchema string = "../../../datastore/github_standards/github_standards.sql"
const templateDir string = "../templates"

func TestServersFrontGithubStandards(t *testing.T) {
	logger.LogSetup()

	//--- spin up an api
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
	// set mock api
	ghapi.SetDBPath(dbF)
	ghapi.SetCtx(ctx)
	mockApi := testhelpers.MockServer(ghapi.ListHandler, "warn")
	defer mockApi.Close()

	// -- mock local server that calls the local api
	templates := template_helpers.GetTemplates(templateDir)
	cfg := &config.Config{Organisation: "TEST RESPONSE"}
	navItem := &navigation.NavigationItem{
		Name:     "test nav",
		Uri:      "/",
		Template: "github-standards",
		DataSources: map[string]src.ApiUrl{
			"list": src.ApiUrl(mockApi.URL),
		},
	}
	handler := func(w http.ResponseWriter, r *http.Request) {
		github_standards.ListHandler(w, r, templates, cfg, navItem)
	}

	mockFront := testhelpers.MockServer(handler, "warn")
	defer mockFront.Close()
	u := urls.Parse("", "", mockFront.URL)
	fmt.Println(mockFront.URL)
	r, _ := getter.GetUrl(u)

	str, _ := convert.Stringify(r)
	// now look in the string for expected data
	title := "<title>test nav - TEST RESPONSE Reports</title>"

	if !strings.Contains(str, title) {
		t.Errorf("expected to find known title, did not")
	}
}
