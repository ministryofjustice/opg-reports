package aws_costs_test

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ministryofjustice/opg-reports/commands/seed/seeder"
	awapi "github.com/ministryofjustice/opg-reports/servers/api/aws_costs"
	"github.com/ministryofjustice/opg-reports/servers/front/aws_costs"
	"github.com/ministryofjustice/opg-reports/servers/front/config"
	"github.com/ministryofjustice/opg-reports/servers/front/config/navigation"
	"github.com/ministryofjustice/opg-reports/servers/front/config/src"
	"github.com/ministryofjustice/opg-reports/servers/front/getter"
	"github.com/ministryofjustice/opg-reports/servers/front/template_helpers"
	"github.com/ministryofjustice/opg-reports/servers/shared/urls"
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
	awapi.SetDBPath(dbF)
	awapi.SetCtx(ctx)
	mockApi := testhelpers.MockServer(awapi.StandardHandler, "warn")
	defer mockApi.Close()

	// -- mock local server that calls the local api
	templates := template_helpers.GetTemplates(templateDir)
	cfg := &config.Config{Organisation: "TEST RESPONSE"}
	navItem := &navigation.NavigationItem{
		Name:     "test nav",
		Uri:      "/",
		Template: "aws-costs-monthly",
		DataSources: map[string]src.ApiUrl{
			"list": src.ApiUrl(mockApi.URL),
		},
	}
	handler := func(w http.ResponseWriter, r *http.Request) {
		aws_costs.StandardHandler(w, r, templates, cfg, navItem)
	}

	mockFront := testhelpers.MockServer(handler, "warn")
	defer mockFront.Close()
	u := urls.Parse("", "", mockFront.URL)
	r, _ := getter.GetUrl(u)

	str, _ := convert.Stringify(r)
	// now look in the string for expected data
	title := "<title>test nav - TEST RESPONSE Reports</title>"

	if !strings.Contains(str, title) {
		t.Errorf("expected to find known title, did not")
	}
}

func TestServersFrontAwsCostsYtd(t *testing.T) {
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
	awapi.SetDBPath(dbF)
	awapi.SetCtx(ctx)
	mockApi := testhelpers.MockServer(awapi.YtdHandler, "warn")
	defer mockApi.Close()

	// -- mock local server that calls the local api
	templates := template_helpers.GetTemplates(templateDir)
	cfg := &config.Config{Organisation: "TEST RESPONSE"}
	navItem := &navigation.NavigationItem{
		Name:     "test ytd",
		Uri:      "/",
		Template: "aws-costs-index",
		DataSources: map[string]src.ApiUrl{
			"list": src.ApiUrl(mockApi.URL),
		},
	}
	handler := func(w http.ResponseWriter, r *http.Request) {
		aws_costs.YtdHandler(w, r, templates, cfg, navItem)
	}

	mockFront := testhelpers.MockServer(handler, "warn")
	defer mockFront.Close()
	u := urls.Parse("", "", mockFront.URL)
	r, _ := getter.GetUrl(u)

	str, _ := convert.Stringify(r)
	// now look in the string for expected data
	title := "<title>test ytd - TEST RESPONSE Reports</title>"

	if !strings.Contains(str, title) {
		t.Errorf("expected to find known title, did not")
	}
}

func TestServersFrontAwsCostsTax(t *testing.T) {
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
	awapi.SetDBPath(dbF)
	awapi.SetCtx(ctx)
	mockApi := testhelpers.MockServer(awapi.MonthlyTaxHandler, "warn")
	defer mockApi.Close()

	// -- mock local server that calls the local api
	templates := template_helpers.GetTemplates(templateDir)
	cfg := &config.Config{Organisation: "TEST RESPONSE"}
	navItem := &navigation.NavigationItem{
		Name:     "test tax",
		Uri:      "/",
		Template: "aws-costs-monthly-tax-totals",
		DataSources: map[string]src.ApiUrl{
			"list": src.ApiUrl(mockApi.URL),
		},
	}
	handler := func(w http.ResponseWriter, r *http.Request) {
		aws_costs.MonthlyTaxHandler(w, r, templates, cfg, navItem)
	}

	mockFront := testhelpers.MockServer(handler, "warn")
	defer mockFront.Close()
	u := urls.Parse("", "", mockFront.URL)
	r, _ := getter.GetUrl(u)

	str, _ := convert.Stringify(r)
	// now look in the string for expected data
	title := "<title>test tax - TEST RESPONSE Reports</title>"

	if !strings.Contains(str, title) {
		t.Errorf("expected to find known title, did not")
	}
}
