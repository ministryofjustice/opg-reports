package datatable

import (
	"opg-reports/report/internal/utils"
	"reflect"
	"testing"
)

var exampleApiResponse string = `
{
    "count": 152,
    "request": {
        "granularity": "monthly",
        "start_date": "2025-04-01",
        "end_date": "2025-05-31",
        "region": "true",
        "environment": "true"
    },
    "dates": [
        "2025-04",
        "2025-05"
    ],
    "groups": [
        "environment",
        "region"
    ],
    "data": [
        {
            "region": "eu-west-1",
            "date": "2025-04",
            "cost": "18698.651307699",
            "environment": "preproduction"
        },
        {
            "region": "global",
            "date": "2025-04",
            "cost": "112.8457268368",
            "environment": "preproduction"
        },
        {
            "region": "global",
            "date": "2025-05",
            "cost": "111.9029923845",
            "environment": "preproduction"
        },
        {
            "region": "global",
            "date": "2025-04",
            "cost": "274.8772830607",
            "environment": "development"
        },
        {
            "region": "eu-west-1",
            "date": "2025-05",
            "cost": "18472.3659787955",
            "environment": "preproduction"
        },
        {
            "region": "eu-west-2",
            "date": "2025-05",
            "cost": "3201.7159650691",
            "environment": "backup"
        },
        {
            "region": "eu-west-2",
            "date": "2025-04",
            "cost": "3169.8725900172",
            "environment": "backup"
        },
        {
            "region": "global",
            "date": "2025-05",
            "cost": "283.2623031041",
            "environment": "development"
        },
        {
            "region": "eu-west-1",
            "date": "2025-04",
            "cost": "16624.4592870909",
            "environment": "development"
        },
        {
            "region": "eu-west-1",
            "date": "2025-04",
            "cost": "122.2465100761",
            "environment": "backup"
        },
        {
            "region": "eu-west-1",
            "date": "2025-04",
            "cost": "23966.5839162545",
            "environment": "production"
        },
        {
            "region": "eu-west-1",
            "date": "2025-05",
            "cost": "24757.93159818",
            "environment": "production"
        },
        {
            "region": "eu-west-1",
            "date": "2025-05",
            "cost": "124.8369473436",
            "environment": "backup"
        },
        {
            "region": "eu-west-1",
            "date": "2025-05",
            "cost": "17437.4141169768",
            "environment": "development"
        },
        {
            "region": "us-east-1",
            "date": "2025-04",
            "cost": "13.5831226785",
            "environment": "preproduction"
        },
        {
            "region": "eu-west-2",
            "date": "2025-05",
            "cost": "5031.2779760823",
            "environment": "production"
        },
        {
            "region": "eu-west-2",
            "date": "2025-05",
            "cost": "2058.7912547951",
            "environment": "preproduction"
        },
        {
            "region": "eu-west-2",
            "date": "2025-04",
            "cost": "2111.4750480189",
            "environment": "preproduction"
        },
        {
            "region": "eu-west-2",
            "date": "2025-04",
            "cost": "4955.1675861698",
            "environment": "production"
        },
        {
            "region": "us-east-1",
            "date": "2025-05",
            "cost": "13.7084966662",
            "environment": "preproduction"
        },
        {
            "region": "eu-west-2",
            "date": "2025-04",
            "cost": "1315.3259356007",
            "environment": "development"
        },
        {
            "region": "eu-west-2",
            "date": "2025-05",
            "cost": "1397.6240849237",
            "environment": "development"
        },
        {
            "region": "us-east-1",
            "date": "2025-04",
            "cost": "28.2387595745",
            "environment": "development"
        },
        {
            "region": "us-east-1",
            "date": "2025-05",
            "cost": "27.9771124331",
            "environment": "development"
        },
        {
            "region": "us-east-1",
            "date": "2025-04",
            "cost": "0.283490873",
            "environment": "backup"
        },
        {
            "region": "eu-west-3",
            "date": "2025-05",
            "cost": "11.2224438",
            "environment": "production"
        },
        {
            "region": "eu-west-3",
            "date": "2025-04",
            "cost": "10.7722553",
            "environment": "production"
        },
        {
            "region": "eu-central-1",
            "date": "2025-05",
            "cost": "4.3365916",
            "environment": "preproduction"
        },
        {
            "region": "eu-west-3",
            "date": "2025-04",
            "cost": "0.4020508",
            "environment": "backup"
        },
        {
            "region": "eu-central-1",
            "date": "2025-04",
            "cost": "12.49873605",
            "environment": "production"
        },
        {
            "region": "eu-west-3",
            "date": "2025-04",
            "cost": "12.459335208",
            "environment": "development"
        },
        {
            "region": "eu-central-1",
            "date": "2025-04",
            "cost": "0.4094957",
            "environment": "backup"
        },
        {
            "region": "eu-central-1",
            "date": "2025-04",
            "cost": "12.735590808",
            "environment": "development"
        },
        {
            "region": "eu-central-1",
            "date": "2025-05",
            "cost": "12.9608380082",
            "environment": "production"
        },
        {
            "region": "eu-central-1",
            "date": "2025-04",
            "cost": "4.2100841",
            "environment": "preproduction"
        },
        {
            "region": "eu-west-3",
            "date": "2025-04",
            "cost": "4.1013085",
            "environment": "preproduction"
        },
        {
            "region": "eu-west-3",
            "date": "2025-05",
            "cost": "0.4071914",
            "environment": "backup"
        },
        {
            "region": "eu-west-3",
            "date": "2025-05",
            "cost": "12.583728984",
            "environment": "development"
        },
        {
            "region": "eu-central-1",
            "date": "2025-05",
            "cost": "0.4156414",
            "environment": "backup"
        },
        {
            "region": "eu-west-3",
            "date": "2025-05",
            "cost": "4.2261655000000005",
            "environment": "preproduction"
        },
        {
            "region": "eu-central-1",
            "date": "2025-05",
            "cost": "12.809816784",
            "environment": "development"
        },
        {
            "region": "us-east-1",
            "date": "2025-05",
            "cost": "0.2786651885",
            "environment": "backup"
        },
        {
            "region": "us-east-1",
            "date": "2025-04",
            "cost": "23.2672791857",
            "environment": "production"
        },
        {
            "region": "us-east-1",
            "date": "2025-05",
            "cost": "24.2866903039",
            "environment": "production"
        },
        {
            "region": "ap-northeast-2",
            "date": "2025-04",
            "cost": "2.1379418",
            "environment": "production"
        },
        {
            "region": "ap-northeast-2",
            "date": "2025-05",
            "cost": "1.6936004",
            "environment": "preproduction"
        },
        {
            "region": "ap-northeast-2",
            "date": "2025-04",
            "cost": "0.0452041",
            "environment": "backup"
        },
        {
            "region": "us-west-2",
            "date": "2025-04",
            "cost": "1.932466",
            "environment": "production"
        },
        {
            "region": "ap-northeast-2",
            "date": "2025-04",
            "cost": "3.9356709",
            "environment": "development"
        },
        {
            "region": "ap-northeast-2",
            "date": "2025-05",
            "cost": "2.2983818",
            "environment": "production"
        },
        {
            "region": "ap-northeast-2",
            "date": "2025-04",
            "cost": "1.6566362",
            "environment": "preproduction"
        },
        {
            "region": "us-west-2",
            "date": "2025-04",
            "cost": "0.0423435",
            "environment": "backup"
        },
        {
            "region": "ap-northeast-1",
            "date": "2025-04",
            "cost": "2.0962381",
            "environment": "production"
        },
        {
            "region": "us-west-2",
            "date": "2025-05",
            "cost": "1.501929",
            "environment": "preproduction"
        },
        {
            "region": "ap-south-1",
            "date": "2025-04",
            "cost": "1.9599762",
            "environment": "production"
        },
        {
            "region": "ca-central-1",
            "date": "2025-04",
            "cost": "2.0514714",
            "environment": "production"
        },
        {
            "region": "ap-northeast-1",
            "date": "2025-04",
            "cost": "0.0443681",
            "environment": "backup"
        },
        {
            "region": "eu-north-1",
            "date": "2025-05",
            "cost": "1.545215",
            "environment": "preproduction"
        },
        {
            "region": "sa-east-1",
            "date": "2025-04",
            "cost": "0.056763",
            "environment": "backup"
        },
        {
            "region": "eu-north-1",
            "date": "2025-04",
            "cost": "1.9133484",
            "environment": "production"
        },
        {
            "region": "us-west-2",
            "date": "2025-04",
            "cost": "3.439028",
            "environment": "development"
        },
        {
            "region": "ap-northeast-2",
            "date": "2025-05",
            "cost": "3.5600375",
            "environment": "development"
        },
        {
            "region": "sa-east-1",
            "date": "2025-04",
            "cost": "2.7945905",
            "environment": "production"
        },
        {
            "region": "us-east-2",
            "date": "2025-04",
            "cost": "0.039831",
            "environment": "backup"
        },
        {
            "region": "ap-southeast-1",
            "date": "2025-05",
            "cost": "1.6859561",
            "environment": "preproduction"
        },
        {
            "region": "ap-southeast-1",
            "date": "2025-04",
            "cost": "2.05099",
            "environment": "production"
        },
        {
            "region": "ap-southeast-2",
            "date": "2025-05",
            "cost": "1.6900976",
            "environment": "preproduction"
        },
        {
            "region": "ap-southeast-2",
            "date": "2025-04",
            "cost": "2.0544614",
            "environment": "production"
        },
        {
            "region": "sa-east-1",
            "date": "2025-05",
            "cost": "2.4958325",
            "environment": "preproduction"
        },
        {
            "region": "us-east-2",
            "date": "2025-05",
            "cost": "1.4822015",
            "environment": "preproduction"
        },
        {
            "region": "us-east-2",
            "date": "2025-04",
            "cost": "1.884009",
            "environment": "production"
        },
        {
            "region": "eu-north-1",
            "date": "2025-04",
            "cost": "0.0408989",
            "environment": "backup"
        },
        {
            "region": "us-west-2",
            "date": "2025-05",
            "cost": "2.0902125",
            "environment": "production"
        },
        {
            "region": "ap-northeast-1",
            "date": "2025-05",
            "cost": "1.73008088",
            "environment": "preproduction"
        },
        {
            "region": "ap-south-1",
            "date": "2025-05",
            "cost": "1.6116415",
            "environment": "preproduction"
        },
        {
            "region": "ca-central-1",
            "date": "2025-05",
            "cost": "1.6890786",
            "environment": "preproduction"
        },
        {
            "region": "ap-south-1",
            "date": "2025-04",
            "cost": "0.0411795",
            "environment": "backup"
        },
        {
            "region": "ca-central-1",
            "date": "2025-04",
            "cost": "0.0435587",
            "environment": "backup"
        },
        {
            "region": "ap-southeast-2",
            "date": "2025-04",
            "cost": "3.9168948",
            "environment": "development"
        },
        {
            "region": "ca-central-1",
            "date": "2025-04",
            "cost": "3.912648",
            "environment": "development"
        },
        {
            "region": "us-west-2",
            "date": "2025-04",
            "cost": "1.4591085",
            "environment": "preproduction"
        },
        {
            "region": "ap-northeast-1",
            "date": "2025-04",
            "cost": "4.01518396",
            "environment": "development"
        },
        {
            "region": "ap-south-1",
            "date": "2025-04",
            "cost": "3.7414669",
            "environment": "development"
        },
        {
            "region": "sa-east-1",
            "date": "2025-04",
            "cost": "5.8910775",
            "environment": "development"
        },
        {
            "region": "ap-southeast-1",
            "date": "2025-05",
            "cost": "2.2208021000000002",
            "environment": "production"
        },
        {
            "region": "ap-southeast-2",
            "date": "2025-04",
            "cost": "0.0435508",
            "environment": "backup"
        },
        {
            "region": "ap-southeast-2",
            "date": "2025-05",
            "cost": "2.2364092",
            "environment": "production"
        },
        {
            "region": "eu-north-1",
            "date": "2025-04",
            "cost": "3.5733411",
            "environment": "development"
        },
        {
            "region": "ca-central-1",
            "date": "2025-05",
            "cost": "2.2238581",
            "environment": "production"
        },
        {
            "region": "sa-east-1",
            "date": "2025-05",
            "cost": "3.024428",
            "environment": "production"
        },
        {
            "region": "us-east-2",
            "date": "2025-04",
            "cost": "3.4199135",
            "environment": "development"
        },
        {
            "region": "ap-southeast-1",
            "date": "2025-04",
            "cost": "3.9158584",
            "environment": "development"
        },
        {
            "region": "ap-southeast-1",
            "date": "2025-04",
            "cost": "0.0424913",
            "environment": "backup"
        },
        {
            "region": "ca-central-1",
            "date": "2025-04",
            "cost": "1.6419038",
            "environment": "preproduction"
        },
        {
            "region": "ap-northeast-3",
            "date": "2025-04",
            "cost": "0.0381995",
            "environment": "backup"
        },
        {
            "region": "ap-northeast-3",
            "date": "2025-04",
            "cost": "2.0462685",
            "environment": "production"
        },
        {
            "region": "us-east-2",
            "date": "2025-05",
            "cost": "2.04028",
            "environment": "production"
        },
        {
            "region": "eu-north-1",
            "date": "2025-04",
            "cost": "1.5076753999999999",
            "environment": "preproduction"
        },
        {
            "region": "sa-east-1",
            "date": "2025-04",
            "cost": "2.43026",
            "environment": "preproduction"
        },
        {
            "region": "ap-northeast-1",
            "date": "2025-04",
            "cost": "1.68239616",
            "environment": "preproduction"
        },
        {
            "region": "ap-northeast-1",
            "date": "2025-05",
            "cost": "2.27480652",
            "environment": "production"
        },
        {
            "region": "ap-south-1",
            "date": "2025-04",
            "cost": "1.567432",
            "environment": "preproduction"
        },
        {
            "region": "ap-northeast-3",
            "date": "2025-05",
            "cost": "1.774408",
            "environment": "preproduction"
        },
        {
            "region": "us-east-2",
            "date": "2025-04",
            "cost": "1.443094",
            "environment": "preproduction"
        },
        {
            "region": "eu-north-1",
            "date": "2025-05",
            "cost": "0.0350063",
            "environment": "backup"
        },
        {
            "region": "us-west-2",
            "date": "2025-05",
            "cost": "0.037053499999999996",
            "environment": "backup"
        },
        {
            "region": "us-west-2",
            "date": "2025-05",
            "cost": "3.118102",
            "environment": "development"
        },
        {
            "region": "ap-southeast-2",
            "date": "2025-04",
            "cost": "1.6431689",
            "environment": "preproduction"
        },
        {
            "region": "ap-northeast-3",
            "date": "2025-04",
            "cost": "4.2011235",
            "environment": "development"
        },
        {
            "region": "ca-central-1",
            "date": "2025-05",
            "cost": "3.5438642",
            "environment": "development"
        },
        {
            "region": "ap-northeast-3",
            "date": "2025-05",
            "cost": "2.214631",
            "environment": "production"
        },
        {
            "region": "ap-south-1",
            "date": "2025-05",
            "cost": "0.0353722",
            "environment": "backup"
        },
        {
            "region": "eu-north-1",
            "date": "2025-05",
            "cost": "3.2334554",
            "environment": "development"
        },
        {
            "region": "sa-east-1",
            "date": "2025-05",
            "cost": "0.0468055",
            "environment": "backup"
        },
        {
            "region": "us-east-2",
            "date": "2025-05",
            "cost": "0.034696",
            "environment": "backup"
        },
        {
            "region": "ap-southeast-2",
            "date": "2025-05",
            "cost": "0.0373348",
            "environment": "backup"
        },
        {
            "region": "ap-northeast-1",
            "date": "2025-05",
            "cost": "0.03814262",
            "environment": "backup"
        },
        {
            "region": "ap-northeast-2",
            "date": "2025-05",
            "cost": "0.0381745",
            "environment": "backup"
        },
        {
            "region": "ap-south-1",
            "date": "2025-05",
            "cost": "3.3909584",
            "environment": "development"
        },
        {
            "region": "ca-central-1",
            "date": "2025-05",
            "cost": "0.0375959",
            "environment": "backup"
        },
        {
            "region": "sa-east-1",
            "date": "2025-05",
            "cost": "5.323087",
            "environment": "development"
        },
        {
            "region": "ap-northeast-1",
            "date": "2025-05",
            "cost": "3.63467372",
            "environment": "development"
        },
        {
            "region": "ap-southeast-1",
            "date": "2025-05",
            "cost": "0.0367961",
            "environment": "backup"
        },
        {
            "region": "ap-southeast-1",
            "date": "2025-04",
            "cost": "1.6387347",
            "environment": "preproduction"
        },
        {
            "region": "ap-southeast-2",
            "date": "2025-05",
            "cost": "3.5453511",
            "environment": "development"
        },
        {
            "region": "us-east-2",
            "date": "2025-05",
            "cost": "3.0960235",
            "environment": "development"
        },
        {
            "region": "ap-southeast-1",
            "date": "2025-05",
            "cost": "3.5419899",
            "environment": "development"
        },
        {
            "region": "ap-south-1",
            "date": "2025-05",
            "cost": "2.1251171",
            "environment": "production"
        },
        {
            "region": "us-west-1",
            "date": "2025-05",
            "cost": "1.5616244",
            "environment": "preproduction"
        },
        {
            "region": "us-west-1",
            "date": "2025-04",
            "cost": "1.7816002",
            "environment": "production"
        },
        {
            "region": "us-west-1",
            "date": "2025-04",
            "cost": "0.0338689",
            "environment": "backup"
        },
        {
            "region": "ap-northeast-3",
            "date": "2025-05",
            "cost": "0.031333",
            "environment": "backup"
        },
        {
            "region": "us-west-1",
            "date": "2025-04",
            "cost": "3.694496",
            "environment": "development"
        },
        {
            "region": "ap-northeast-3",
            "date": "2025-05",
            "cost": "3.789459",
            "environment": "development"
        },
        {
            "region": "eu-north-1",
            "date": "2025-05",
            "cost": "2.0672513",
            "environment": "production"
        },
        {
            "region": "us-west-1",
            "date": "2025-05",
            "cost": "1.9504102",
            "environment": "production"
        },
        {
            "region": "ap-northeast-3",
            "date": "2025-04",
            "cost": "1.7243395",
            "environment": "preproduction"
        },
        {
            "region": "us-west-1",
            "date": "2025-04",
            "cost": "1.5200787",
            "environment": "preproduction"
        },
        {
            "region": "us-west-1",
            "date": "2025-05",
            "cost": "0.0283828",
            "environment": "backup"
        },
        {
            "region": "us-west-1",
            "date": "2025-05",
            "cost": "3.3410172",
            "environment": "development"
        },
        {
            "region": "global",
            "date": "2025-04",
            "cost": "0",
            "environment": "backup"
        },
        {
            "region": "global",
            "date": "2025-05",
            "cost": "0",
            "environment": "backup"
        },
        {
            "region": "global",
            "date": "2025-04",
            "cost": "140.0404518213",
            "environment": "production"
        },
        {
            "region": "global",
            "date": "2025-05",
            "cost": "152.1211034773",
            "environment": "production"
        },
        {
            "region": "NoRegion",
            "date": "2025-05",
            "cost": "-5291.496897623",
            "environment": "development"
        },
        {
            "region": "NoRegion",
            "date": "2025-04",
            "cost": "-7008.6244985329",
            "environment": "production"
        },
        {
            "region": "NoRegion",
            "date": "2025-05",
            "cost": "-7220.0696471628",
            "environment": "production"
        },
        {
            "region": "NoRegion",
            "date": "2025-04",
            "cost": "-5039.2425468039",
            "environment": "development"
        },
        {
            "region": "NoRegion",
            "date": "2025-04",
            "cost": "-757.557069693",
            "environment": "backup"
        },
        {
            "region": "NoRegion",
            "date": "2025-05",
            "cost": "-765.4609545169",
            "environment": "backup"
        },
        {
            "region": "NoRegion",
            "date": "2025-05",
            "cost": "-5034.7822818399",
            "environment": "preproduction"
        },
        {
            "region": "NoRegion",
            "date": "2025-04",
            "cost": "-5091.0731353969",
            "environment": "preproduction"
        }
    ]
}
`

type testResponse struct {
	Count  int                 `json:"count,omityempty"`
	Dates  []string            `json:"dates,omitempty"`
	Groups []string            `json:"groups,omitempty"`
	Data   []map[string]string `json:"data"`
}

func (self *testResponse) DataHeaders() (dh []string) {
	return self.Dates
}
func (self *testResponse) DataRows() (data []map[string]string) {
	return self.Data
}
func (self *testResponse) PaddedDataRows() (all []map[string]string) {
	padding := utils.DummyRows(self.Dates, "date")
	all = append(self.Data, padding...)
	return
}
func (self *testResponse) Identifiers() (identifiers []string) {
	return self.Groups
}
func (self *testResponse) Cells() (cells []string) {
	return self.Dates
}
func (self *testResponse) TransformColumn() string {
	return "date"
}
func (self *testResponse) ValueColumn() string {
	return "cost"
}
func (self *testResponse) RowTotalKeyName() string {
	return "total"
}
func (self *testResponse) TrendKeyName() string {
	return "trend"
}
func (self *testResponse) SumColumns() (cols []string) {
	cols = self.Dates
	if c := self.RowTotalKeyName(); c != "" {
		cols = append(cols, c)
	}
	return
}

var _ ResponseBody = &testResponse{}

// TestDataTableNew the table generated is quite large, but we can check the totals
// row and make sure headers etc are set as expected
func TestDataTableNew(t *testing.T) {
	var err error
	var response = &testResponse{}

	// make a response from the sample data string
	err = utils.Unmarshal([]byte(exampleApiResponse), &response)
	if err != nil {
		t.Errorf("failed to unmarshal example data")
		t.FailNow()
	}

	dt, err := New(response)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
		t.FailNow()
	}
	// make sure row headers match expected values
	if len(dt.RowHeaders) != 2 {
		t.Errorf("incorrect row headers count")
	}
	if !reflect.DeepEqual(dt.RowHeaders, []string{"environment", "region"}) {
		t.Errorf("incorrect row headers")
	}
	// check the dates match
	if len(dt.DataHeaders) != 2 {
		t.Errorf("incorrect data headers count")
	}
	if !reflect.DeepEqual(dt.DataHeaders, []string{"2025-04", "2025-05"}) {
		t.Errorf("incorrect data headers")
	}

	// check the extra headers (total etc) match
	if len(dt.ExtraHeaders) != 2 {
		t.Errorf("incorrect extra headers count")
	}
	if !reflect.DeepEqual(dt.ExtraHeaders, []string{"trend", "total"}) {
		t.Errorf("incorrect extra headers")
	}

}

// TestTableGeneration runs though what New would call but passes
// some fixed values and runs it step by step
func TestTableGeneration(t *testing.T) {
	var err error
	var response = &testResponse{}
	var identifiers = []string{"environment", "region"}
	var dates = []string{"2025-03", "2025-04", "2025-05"}

	// make a data source
	err = utils.Unmarshal([]byte(exampleApiResponse), &response)
	if err != nil {
		t.Errorf("failed to unmarshal example data")
		t.FailNow()
	}

	possibles, u := PossibleCombinationsAsKeys(response.Data, identifiers)

	// the number of possible combinations should be based on the number
	// of unique values per column
	expectedCount := 1
	for _, pv := range u {
		expectedCount = expectedCount * len(pv)
	}

	if expectedCount != len(possibles) {
		t.Errorf("permutations are incorrect")
	}

	// now generated a table
	table := SkeletonTable(possibles, dates)
	ptable := PopulateTable(response.Data, table, identifiers, "date", "cost")

	if len(table) == 0 || len(ptable) == 0 {
		t.Error("table creation failed")
	}
	AddRowTotals(ptable, identifiers, "total")
	cols := append(dates, "total")
	totals := ColumnTotals(ptable, cols)

	if len(totals) == 0 {
		t.Error("failed to parse totals")
	}

}

func TestTransformersPermutations(t *testing.T) {

	var checkN = []string{"1", "2", "3"}
	var checkA = []string{"A", "B"}
	var check = [][]string{checkN, checkA}
	var expected = len(checkN) * len(checkA)

	p := Permutations(check...)
	actual := len(p)
	if actual != expected {
		t.Errorf("permutation length mismtach")
	}
}
