package owner

import (
	"fmt"
	"log/slog"
	"os"
	"testing"

	"github.com/ministryofjustice/opg-reports/report/config"
	"github.com/ministryofjustice/opg-reports/report/internal/repository"
)

func TestOwnerServiceNew(t *testing.T) {
	var (
		err error
		dir = t.TempDir()
		ctx = t.Context()
		cfg = config.NewConfig()
	)
	cfg.Database.Path = fmt.Sprintf("%s/%s", dir, "test.db")

	lg := slog.New(slog.NewTextHandler(os.Stdout, nil))
	rep, _ := repository.New(ctx, lg, cfg)

	_, err = NewService(ctx, lg, rep)
	if err != nil {
		t.Errorf("unexpected error creating service: [%s]", err.Error())
	}

	_, err = NewService(ctx, nil, nil)
	if err == nil {
		t.Errorf("New service should have thrown error without a log or repository")
	}
	_, err = NewService(ctx, lg, nil)
	if err == nil {
		t.Errorf("New service should have thrown error without a repository")
	}

}
