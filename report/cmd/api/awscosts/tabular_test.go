package awscosts

import (
	"opg-reports/report/internal/service/api"
	"opg-reports/report/internal/utils"
	"testing"
)

type tTabularTest struct {
	Columns     []string
	Dates       []string
	Costs       []*api.AwsCostGrouped
	ExpectedLen int
}

func TestApiAwsCostsTabularA(t *testing.T) {
	var (
		err    error
		sorted []map[string]string
		log    = utils.Logger("ERROR", "TEXT")
	)
	var testA *tTabularTest = &tTabularTest{
		Columns: []string{
			"environment",
			"region",
		},
		Dates: []string{
			"2025-04",
			"2025-05",
			"2025-03",
		},
		ExpectedLen: 3,
		Costs: []*api.AwsCostGrouped{
			{
				Region:                "eu-west-1",
				AwsAccountEnvironment: "development",
				Date:                  "2025-04",
				Cost:                  "10",
			},
			{
				Region:                "eu-west-1",
				AwsAccountEnvironment: "preproduction",
				Date:                  "2025-04",
				Cost:                  "100",
			},
			{
				Region:                "eu-west-1",
				AwsAccountEnvironment: "production",
				Date:                  "2025-04",
				Cost:                  "1000",
			},
			{
				Region:                "eu-west-1",
				AwsAccountEnvironment: "development",
				Date:                  "2025-05",
				Cost:                  "20",
			},
			{
				Region:                "eu-west-1",
				AwsAccountEnvironment: "preproduction",
				Date:                  "2025-05",
				Cost:                  "200",
			},
			{
				Region:                "eu-west-1",
				AwsAccountEnvironment: "production",
				Date:                  "2025-05",
				Cost:                  "2000",
			},
		},
	}

	sorted, _, err = TabulateGroupedCosts(log, testA.Columns, testA.Dates, testA.Costs)
	if err != nil {
		t.Errorf("unexpected error: %v", err.Error())
	}
	// check length matches exptectations
	if len(sorted) != testA.ExpectedLen {
		t.Errorf("unexpected length; expected %d, actual :%v", testA.ExpectedLen, len(sorted))
	}
	// all 2025-03 fields should be 0
	for i, row := range sorted {
		if v, ok := row["2025-03"]; !ok || v != "0.00" {
			t.Errorf("[%d] 2025-03 value incorrect: [%v]", i, row)
		}
	}
}
