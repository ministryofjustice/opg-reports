package awscostexplorer

import (
	"log/slog"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/costexplorer"
)

// DateFormat is format used for calls to the aws api
const DateFormat string = "2006-01-02"

// Connection is used as a wrapper (to allow for interfaces)
// around the aws client
type Connection struct {
	client *costexplorer.CostExplorer
}

// GetClient returns the aws client
func (self *Connection) GetClient() *costexplorer.CostExplorer {
	return self.client
}

// NewConnection creates a new wrapper from the client
func NewConnection(client *costexplorer.CostExplorer) *Connection {
	return &Connection{client: client}
}

// Parameters are the variable parameters that are used to call
// the costexplorer sdk function
// Note: requested metrics and grouping is fixed
type Parameters struct {
	StartDate   time.Time
	EndDate     time.Time
	Granularity string
}

// ToSDK uses the fixed values for metrics and grouping with the struct values
// to generate an input structure for the SDK call
func (self *Parameters) ToSDK() *costexplorer.GetCostAndUsageInput {

	return &costexplorer.GetCostAndUsageInput{
		TimePeriod: &costexplorer.DateInterval{
			Start: aws.String(self.StartDate.Format(DateFormat)),
			End:   aws.String(self.EndDate.Format(DateFormat)),
		},
		Granularity: aws.String(self.Granularity),
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

// Response
type Response struct {
	response *costexplorer.GetCostAndUsageOutput
}

func (self *Response) GetResult() *costexplorer.GetCostAndUsageOutput {
	return self.response
}

// NewResponse returns a response
func NewResponse(response *costexplorer.GetCostAndUsageOutput) *Response {
	return &Response{
		response: response,
	}
}

// Service is used to call and fetch cost data
// from AWS costexplorer api
type Service struct {
	connection *Connection
	logger     *slog.Logger
}

// GetLogger returns the logger to use for any calls in this service
func (self *Service) GetLogger() *slog.Logger {
	return self.logger
}

// GetConnection returns the connection (in this case aws costexplorer) so it can
// then be used to via the SDK etc
func (self *Service) GetConnection() *Connection {
	return self.connection
}

// GetData uses the params passed to fetch data from the AWS costexplorer api
// and return the raw result
func (self *Service) GetData(params *Parameters) (result *Response, err error) {
	var (
		client = self.GetConnection().GetClient()
		input  = params.ToSDK()
		logger = self.GetLogger().With("operation", "GetData").With("params", params)
	)

	request, response := client.GetCostAndUsageRequest(input)

	err = request.Send()
	if err != nil {
		logger.Error("failed to fetch costexplorer data")
	} else {
		logger.Debug("fetched costexplorer data")
		result = NewResponse(response)
	}

	return
}

func NewService(logger *slog.Logger, connection *Connection) *Service {
	return &Service{
		logger:     logger.WithGroup("AWSCostExplorerService"),
		connection: connection,
	}
}
