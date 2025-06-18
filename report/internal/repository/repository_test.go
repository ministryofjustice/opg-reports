package repository

import (
	"fmt"
	"log/slog"
	"os"
	"testing"

	"github.com/ministryofjustice/opg-reports/report/config"
)

func TestRepositoryNew(t *testing.T) {
	var (
		err error
		dir = t.TempDir()
		ctx = t.Context()
		cfg = config.NewConfig()
	)
	cfg.Database.Path = fmt.Sprintf("%s/%s", dir, "test.db")
	lg := slog.New(slog.NewTextHandler(os.Stdout, nil))

	_, err = New(ctx, lg, cfg)

	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}

}
