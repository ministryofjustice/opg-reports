package cost

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/costexplorer"
)

func CostAndUsageInput(startDate time.Time, endDate time.Time, granularity string, dateFormat string) *costexplorer.GetCostAndUsageInput {

	input := &costexplorer.GetCostAndUsageInput{
		TimePeriod: &costexplorer.DateInterval{
			Start: aws.String(startDate.Format(dateFormat)),
			End:   aws.String(endDate.Format(dateFormat)),
		},
		Granularity: aws.String(granularity),
		Metrics: []*string{
			aws.String("UNBLENDED_COST"),
		},
		GroupBy: []*costexplorer.GroupDefinition{
			{
				Type: aws.String("DIMENSION"),
				Key:  aws.String("SERVICE"),
			},
			{
				Type: aws.String("DIMENSION"),
				Key:  aws.String("REGION"),
			},
		},
	}
	return input

}

func CostAndUsage(roleArn string, startDate time.Time, endDate time.Time, region string, granularity string, dateFormat string) (*costexplorer.GetCostAndUsageOutput, error) {

	ceClient, err := Client(roleArn, region)
	if err != nil {
		slog.Error(fmt.Sprintf("error: CostAndUsage client: %v", err.Error()))
		return nil, err
	}

	sdkInput := CostAndUsageInput(startDate, endDate, granularity, dateFormat)
	slog.Debug("CostAndUsage",
		slog.String("start", *sdkInput.TimePeriod.Start),
		slog.String("end", *sdkInput.TimePeriod.End),
		slog.String("arn", roleArn),
	)
	request, response := ceClient.GetCostAndUsageRequest(sdkInput)
	err = request.Send()
	if err != nil {
		slog.Error(fmt.Sprintf("error: CostAndUsage request: %v", err.Error()))
		return nil, err
	}
	return response, nil

}
