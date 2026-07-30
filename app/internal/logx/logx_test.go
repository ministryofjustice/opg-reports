package logx_test

import (
	"log/slog"
	"opg-reports/app/internal/logx"
	"testing"
)

func TestLogxSet(t *testing.T) {
	var ctx = t.Context()

	logx.Set(slog.LevelError)
	lg := slog.Default()

	if lg.Enabled(ctx, slog.LevelWarn) {
		t.Errorf("only error and upwards should be set.")
	}
	if !lg.Enabled(ctx, slog.LevelError) {
		t.Errorf("error should be enabled.")
	}

}
