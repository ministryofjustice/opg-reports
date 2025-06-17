package awscostexplorer_test

import (
	"fmt"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/costexplorer"
	"github.com/ministryofjustice/opg-reports/internal/services/awssdk/awscostexplorer"
	"github.com/ministryofjustice/opg-reports/internal/utils/envar"
)

func TestCostExplorerWithAwsCreds(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	if envar.Get("AWS_SESSION_TOKEN", "") != "" {

		params := &awscostexplorer.Parameters{
			StartDate:   time.Date(2025, 5, 1, 0, 0, 0, 0, time.UTC),
			EndDate:     time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC),
			Granularity: costexplorer.GranularityMonthly,
		}

		s, _ := session.NewSession(&aws.Config{
			Region: aws.String("eu-west-1"),
		})
		client := costexplorer.New(s)
		conn := awscostexplorer.NewConnection(client)
		srv := awscostexplorer.NewService(logger, conn)

		a, b := srv.GetData(params)
		fmt.Printf("-->%+v", a.GetResult())
		fmt.Printf("-->%+v", b)

	} else {
		logger.Info("Skipping test as no AWS Creds found [TestServicesAwsSDKCostExplorerWithAwsCreds]")
	}
}
