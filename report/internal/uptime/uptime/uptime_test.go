package uptime

import (
	"context"
	"log/slog"
	"opg-reports/report/internal/utils/awsclients"
	"opg-reports/report/internal/utils/logger"
	"os"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
)

func TestRedoUptimeWithoutMock(t *testing.T) {
	var (
		err    error
		client *cloudwatch.Client
		r      *cloudwatch.GetMetricStatisticsOutput
		ctx    context.Context = t.Context()
		log    *slog.Logger    = logger.New("error")
		now    time.Time       = time.Now().UTC()
		start  time.Time       = now.AddDate(0, -4, 0)
		end    time.Time       = now.AddDate(0, -3, 0)
	)

	if os.Getenv("AWS_SESSION_TOKEN") != "" {
		client, err = awsclients.New[*cloudwatch.Client](ctx, log, "us-east-1")
		if err != nil {
			t.Errorf("unexpected error:\n%s", err.Error())
			t.FailNow()
		}

		r, err = GetUptimeData(ctx, log, client, &GetUptimeDataOptions{Start: start, End: end})
		if err != nil {
			t.Errorf("unexpected error:\n%s", err.Error())
		}
		if len(r.Datapoints) <= 0 {
			t.Error("failed to find uptime data")
		}
	} else {
		t.Skip()
	}
}
