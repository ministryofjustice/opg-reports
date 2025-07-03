package team

import (
	"fmt"
	"testing"

	"github.com/ministryofjustice/opg-reports/report/config"
	"github.com/ministryofjustice/opg-reports/report/internal/repository/sqlr"
	"github.com/ministryofjustice/opg-reports/report/internal/service/seed"
	"github.com/ministryofjustice/opg-reports/report/internal/utils"
)

func TestTeamServiceNew(t *testing.T) {
	var (
		err error
		srv *Service[*Team]
		dir = t.TempDir()
		ctx = t.Context()
		cfg = config.NewConfig()
	)
	cfg.Database.Path = fmt.Sprintf("%s/%s", dir, "test-team-connection.db")

	lg := utils.Logger("ERROR", "TEXT")
	rep, _ := sqlr.NewWithSelect[*Team](ctx, lg, cfg)

	srv, err = NewService(ctx, lg, cfg, rep)
	if err != nil {
		t.Errorf("unexpected error creating service: [%s]", err.Error())
	}
	defer srv.Close()

	srv, err = NewService[*Team](ctx, nil, nil, nil)
	if err == nil {
		t.Errorf("New service should have thrown error without a log or repository")
	}
	defer srv.Close()

	srv, err = NewService[*Team](ctx, lg, nil, nil)
	if err == nil {
		t.Errorf("New service should have thrown error without a repository")
	}
	defer srv.Close()

}

func TestTeamServiceGetAll(t *testing.T) {
	var (
		err  error
		dir  = t.TempDir()
		ctx  = t.Context()
		conf = config.NewConfig()
		log  = utils.Logger("ERROR", "TEXT")
	)
	conf.Database.Path = fmt.Sprintf("%s/%s", dir, "test-team-getall.db")

	sqc := sqlr.Default(ctx, log, conf)
	seeder := seed.Default(ctx, log, conf)
	insert, _ := seeder.Teams(sqc)

	// generate the service
	srv := Default[*Team](ctx, log, conf)
	if srv == nil {
		t.Errorf("unexpected error creating service: [%s]", err.Error())
		t.FailNow()
	}
	defer srv.Close()

	// fetch everything
	res, err := srv.GetAllTeams()
	if err != nil {
		t.Errorf("unexpected error getting data from service: [%s]", err.Error())
	}
	if len(res) != len(insert) {
		t.Errorf("no results found in service")
	}
}
