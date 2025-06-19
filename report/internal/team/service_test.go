package team

import (
	"fmt"
	"log/slog"
	"os"
	"testing"

	"github.com/ministryofjustice/opg-reports/report/config"
	"github.com/ministryofjustice/opg-reports/report/internal/sqldb"
)

func TestTeamServiceNew(t *testing.T) {
	var (
		err error
		dir = t.TempDir()
		ctx = t.Context()
		cfg = config.NewConfig()
	)
	cfg.Database.Path = fmt.Sprintf("%s/%s", dir, "test.db")

	lg := slog.New(slog.NewTextHandler(os.Stdout, nil))
	rep, _ := sqldb.New[*Team](ctx, lg, cfg)

	_, err = NewService(ctx, lg, rep)
	if err != nil {
		t.Errorf("unexpected error creating service: [%s]", err.Error())
	}

	_, err = NewService[*Team](ctx, nil, nil)
	if err == nil {
		t.Errorf("New service should have thrown error without a log or repository")
	}
	_, err = NewService[*Team](ctx, lg, nil)
	if err == nil {
		t.Errorf("New service should have thrown error without a repository")
	}

}

func TestTeamServiceGetAll(t *testing.T) {
	var (
		err error
		dir = t.TempDir()
		ctx = t.Context()
		cfg = config.NewConfig()
		lg  = slog.New(slog.NewTextHandler(os.Stdout, nil))
	)
	cfg.Database.Path = fmt.Sprintf("%s/%s", dir, "test.db")

	rep, err := sqldb.New[*Team](ctx, lg, cfg)
	if err != nil {
		t.Errorf("unexpected error creating repository: [%s]", err.Error())
	}

	srv, err := NewService(ctx, lg, rep)
	if err != nil {
		t.Errorf("unexpected error creating service: [%s]", err.Error())
	}
	// insert standard items
	err = srv.Import()
	if err != nil {
		t.Errorf("unexpected error seeding service: [%s]", err.Error())
	}
	// fetch everything
	res, err := srv.GetAllTeams()
	if err != nil {
		t.Errorf("unexpected error getting data from service: [%s]", err.Error())
	}
	if len(res) <= 0 {
		t.Errorf("no results found in service")
	}
}
