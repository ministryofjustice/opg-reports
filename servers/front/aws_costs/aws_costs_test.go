package aws_costs_test

import (
	"context"
	"log"
	"log/slog"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ministryofjustice/opg-reports/commands/seed/seeder"
	awapi "github.com/ministryofjustice/opg-reports/servers/api/aws_costs"
	"github.com/ministryofjustice/opg-reports/servers/front/aws_costs"
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

const realSchema string = "../../../datastore/aws_costs/aws_costs.sql"
const templateDir string = "../templates"

func TestServersFrontAwsCostsStandard(t *testing.T) {
	logger.LogSetup()

	//--- spin up an api
	// seed
	ctx := context.TODO()
	N := 100
	dir := t.TempDir()
	dbF := filepath.Join(dir, "aws.db")
	schemaF := filepath.Join(dir, "aws.sql")
	dataF := filepath.Join(dir, "dummy.json")
	testhelpers.CopyFile(realSchema, schemaF)
	db, err := seeder.Seed(ctx, dbF, schemaF, dataF, "aws_costs", N)
	if err != nil {
		slog.Error(err.Error())
		log.Fatal(err.Error())
	}
	defer db.Close()
	// set mock api
	apiserver := api.New(ctx, dbF)
	apihandler := api.Wrap(apiserver, awapi.StandardHandler)
	mockApi := testhelpers.MockServer(apihandler, "warn")
	defer mockApi.Close()

	// -- mock local server that calls the local api
	templates := template.GetTemplates(templateDir)
	cfg := &config.Config{Organisation: "TEST RESPONSE"}
	navItem := &nav.Nav{
		Name:     "test nav",
		Uri:      "/",
		Template: "aws-costs-monthly",
		DataSources: map[string]string{
			"list": mockApi.URL,
		},
	}

	server := front.New(ctx, cfg, templates)
	handler := front.Wrap(server, navItem, aws_costs.Handler)

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
