package lib

import (
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/costexplorer"
	"github.com/ministryofjustice/opg-reports/models"
	"github.com/ministryofjustice/opg-reports/pkg/consts"
	"github.com/ministryofjustice/opg-reports/pkg/convert"
)

var (
	defEnv   = "production"
	defMonth = convert.DateResetMonth(time.Now().UTC()).AddDate(0, -1, 0)
)

// Arguments represents all the named arguments for this collector
type Arguments struct {
	Month      string
	ID         string
	Name       string
	Label      string
	Env        string
	Unit       string
	OutputFile string
}

// SetupArgs maps flag values to properies on the arg passed and runs
// flag.Parse to fetch values
func SetupArgs(args *Arguments) {

	flag.StringVar(&args.Month, "month", defMonth.Format(consts.DateFormatYearMonthDay), "month to fetch data for.")
	flag.StringVar(&args.ID, "id", "", "AWS account id.")
	flag.StringVar(&args.Name, "name", "", "AWS account name.")
	flag.StringVar(&args.Label, "label", "", "Account label.")
	flag.StringVar(&args.Env, "environment", defEnv, "Account environment type.")

	flag.StringVar(&args.Unit, "unit", "", "Unit / team name.")
	flag.StringVar(&args.OutputFile, "output", "./data/{month}_{id}_aws_costs.json", "Filepath for the output")

	flag.Parse()
}

// ValidateArgs checks rules and logic for the input arguments
// Make sure some have non empty values and apply default values to others
func ValidateArgs(args *Arguments) (err error) {
	failOnEmpty := map[string]string{
		"month":  args.Month,
		"id":     args.ID,
		"name":   args.Name,
		"label":  args.Label,
		"unit":   args.Unit,
		"output": args.OutputFile,
	}
	for k, v := range failOnEmpty {
		if v == "" {
			err = errors.Join(err, fmt.Errorf("%s", k))
		}
	}
	if err != nil {
		err = fmt.Errorf("missing arguments: [%s]", strings.ReplaceAll(err.Error(), "\n", ", "))
	}

	if args.Month == "-" {
		args.Month = defMonth.Format(consts.DateFormat)
	}

	if args.Env == "" || args.Env == "-" {
		args.Env = defEnv
	}

	return
}

// WriteToFile writes the content to the file replacing values in
// the filename with values on arg
func WriteToFile(content []byte, args *Arguments) {
	var (
		filename string
		dir      string = filepath.Dir(args.OutputFile)
	)
	os.MkdirAll(dir, os.ModePerm)
	filename = args.OutputFile
	filename = strings.ReplaceAll(filename, "{month}", args.Month)
	filename = strings.ReplaceAll(filename, "{id}", args.ID)

	os.WriteFile(filename, content, os.ModePerm)

}

// SdkInput generates an input struct for usage
func SdkInput(startDate time.Time, endDate time.Time, granularity string, dateFormat string) *costexplorer.GetCostAndUsageInput {

	return &costexplorer.GetCostAndUsageInput{
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

}

// CostData makes the api call to fetch the data
func CostData(client *costexplorer.CostExplorer, startDate time.Time, endDate time.Time, granularity string, dateFormat string) (response *costexplorer.GetCostAndUsageOutput, err error) {

	sdkInput := SdkInput(startDate, endDate, granularity, dateFormat)
	slog.Debug("CostAndUsage",
		slog.String("start", *sdkInput.TimePeriod.Start),
		slog.String("end", *sdkInput.TimePeriod.End),
	)
	request, response := client.GetCostAndUsageRequest(sdkInput)
	err = request.Send()
	if err != nil {
		slog.Error(fmt.Sprintf("error: CostAndUsage request: %v", err.Error()))
		return nil, err
	}
	return response, nil

}

// Flat converts the multi dimensional raw data into a flat slice of costs.Cost entryies
func Flat(raw *costexplorer.GetCostAndUsageOutput, args *Arguments) (flat []*models.AwsCost, err error) {
	slog.Debug("Flattening cost data")
	now := time.Now().UTC().Format(consts.DateFormat)

	flat = []*models.AwsCost{}
	unit := &models.Unit{
		Ts:   now,
		Name: strings.ToLower(args.Unit),
	}
	account := &models.AwsAccount{
		Ts:          now,
		Number:      args.ID,
		Name:        args.Name,
		Label:       args.Label,
		Environment: args.Env,
	}

	for _, resultByTime := range raw.ResultsByTime {
		day := *resultByTime.TimePeriod.Start

		for _, costGroup := range resultByTime.Groups {
			service := *costGroup.Keys[0]
			region := *costGroup.Keys[1]

			for _, costMetric := range costGroup.Metrics {
				amount := *costMetric.Amount

				c := &models.AwsCost{
					Ts:         now,
					Region:     region,
					Service:    service,
					Date:       day,
					Cost:       amount,
					Unit:       (*models.UnitForeignKey)(unit),
					AwsAccount: (*models.AwsAccountForeignKey)(account),
				}

				flat = append(flat, c)
			}
		}
	}

	return
}
