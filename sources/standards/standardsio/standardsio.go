package standardsio

import (
	"github.com/danielgtaylor/huma/v2"
	"github.com/ministryofjustice/opg-reports/pkg/convert"
	"github.com/ministryofjustice/opg-reports/sources/standards"
)

type StandardsInput struct {
	Version  string `json:"version" path:"version" required:"true" doc:"Version prefix for the api" default:"v1" enum:"v1"`
	Archived bool   `json:"archived" path:"archived" default:"false" doc:"Returns results with is_archived matching this value"`
	Unit     string `json:"unit" db:"unit" query:"unit" doc:"optional unit name to filter data by"`

	IsArchived uint8  `json:"is_archived" db:"is_archived"`
	Teams      string `json:"teams" db:"teams"`
}

func (self *StandardsInput) Resolve(ctx huma.Context) []error {
	self.IsArchived = convert.BoolToInt(self.Archived)

	if self.Unit != "" {
		self.Teams = `%#` + self.Unit + `#%`
	}
	return nil
}

var _ huma.Resolver = (*StandardsInput)(nil)

type Counters struct {
	Total                  int `json:"total" doc:"Overall total number of records in the database."`
	TotalBaselineCompliant int `json:"total_baseline_compliant" doc:"Overall number of records that are baseline compliant."`
	TotalExtendedCompliant int `json:"total_extended_compliant" doc:"Overall number of records that are extended compliant."`
	TotalArchived          int `json:"total_archived" doc:"Overall number of archived records."`

	Count             int `json:"count" doc:"Number of results returned in this query."`
	BaselineCompliant int `json:"baseline_compliant" doc:"Number of results in this query that are baseline compliant."`
	ExtendedCompliant int `json:"extended_compliant" doc:"Number of results in this query that are extended compliant."`
}

type StandardsBody struct {
	Type     string                `json:"type" doc:"States what type of data this is for front end handling"`
	Result   []*standards.Standard `json:"result" doc:"List of all matching repository data."`
	Request  *StandardsInput       `json:"request" doc:"The original request."`
	Counters *Counters             `json:"counters" doc:"Count information."`
}

type StandardsOutput struct {
	Body *StandardsBody
}
