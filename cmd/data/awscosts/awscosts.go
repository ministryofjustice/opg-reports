package awscosts

import (
	"flag"
	"fmt"
	"log/slog"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/costexplorer"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/ministryofjustice/opg-reports/internal/services/awssdk/awscostexplorer"
	"github.com/ministryofjustice/opg-reports/internal/services/awssdk/awssts"
	"github.com/ministryofjustice/opg-reports/internal/utils/cliargs"
	"github.com/ministryofjustice/opg-reports/internal/utils/convert"
)

type AwsCostsArgs struct {
	StartDate   cliargs.DateTime
	EndDate     cliargs.DateTime
	Region      string
	Granularity string
}

var (
	now     = time.Now().UTC()
	start   = time.Date(now.Year(), now.Month()-1, 1, 0, 0, 0, 0, time.UTC)
	end     = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	flagSet = flag.NewFlagSet("awscosts", flag.ExitOnError)
)

func GetAWSCostData(logger *slog.Logger, s *session.Session, args *AwsCostsArgs) (response *awscostexplorer.Response, err error) {
	client := costexplorer.New(s)
	conn := awscostexplorer.NewConnection(client)
	srv := awscostexplorer.NewService(logger, conn)
	response, err = srv.GetData(&awscostexplorer.Parameters{
		StartDate:   *args.StartDate.Value,
		EndDate:     *args.EndDate.Value,
		Granularity: args.Granularity,
	})

	return
}

func GetAWSAccountData(logger *slog.Logger, s *session.Session) (response *awssts.Response, err error) {
	client := sts.New(s)

	conn := awssts.NewConnection(client)
	srv := awssts.NewService(logger, conn)
	response, err = srv.GetCallerID()

	return
}

// Flags is used to parse arguments for this command
func Flags(args *AwsCostsArgs) *flag.FlagSet {
	flagSet.Var(&args.StartDate, "start-date", "Start date for the command (YYYY-MM-DD)")
	flagSet.Var(&args.EndDate, "end-date", "End date for the command (YYYY-MM-DD)")
	flagSet.StringVar(&args.Granularity, "granularity", costexplorer.GranularityDaily, "DAILY|MONHTLY")

	return flagSet
}

func Run(logger *slog.Logger, argValues []string) (err error) {
	args := &AwsCostsArgs{
		StartDate:   cliargs.DateTime{Value: &start},
		EndDate:     cliargs.DateTime{Value: &end},
		Granularity: costexplorer.GranularityDaily,
		Region:      "eu-west-1",
	}
	// Check the arguments
	flagset := Flags(args)
	flagset.Parse(argValues)

	// AWS session using eu-west-1
	s, err := session.NewSession(&aws.Config{Region: &args.Region})
	if err != nil {
		return
	}
	// Check the account details so we know the account id
	account, err := GetAWSAccountData(logger, s)
	if err != nil {
		return
	}
	// get cost data from aws
	response, err := GetAWSCostData(logger, s, args)

	// convert from aws result to local
	if err == nil {
		costs := convert.FromAwsCostExplorerToAwsCosts(response.GetResult(), *account.GetResult().Account)
		fmt.Printf("%+v\n", costs)
	}

	return
}
