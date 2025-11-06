package tables

import (
	"opg-reports/report/internal/utils"
	"testing"
)

type testRecord struct {
	Team   string `json:"team,omitempty"`
	Region string `json:"region,omitempty"`
	Date   string `json:"date,omitempty"`
	Cost   string `json:"cost,omitempty"`
}

// TestApiTablesListToTable does simple test to make sure table is generated
func TestApiTablesListToTable(t *testing.T) {
	var (
		err         error
		log         = utils.Logger("ERROR", "TEXT")
		expectedLen = 5
		dates       = []string{"2025-03", "2025-04", "2025-05"}
		identifiers = []string{"team", "region"}
		lcfg        = &ListToTableConfig{
			TextColumns:     identifiers,
			DataColumns:     dates,
			ValueField:      "cost",
			DataSourceField: "date",
			DefaultValue:    "0.000",
		}
		testRecords = []*testRecord{
			{
				Region: "eu-west-1",
				Team:   "T01",
				Date:   "2025-04",
				Cost:   "100",
			},
			{
				Region: "eu-west-1",
				Team:   "T01",
				Date:   "2025-05",
				Cost:   "100",
			},
			{
				Region: "eu-west-1",
				Team:   "T02",
				Date:   "2025-05",
				Cost:   "2000",
			},
			{
				Region: "eu-west-2",
				Team:   "T00",
				Date:   "2025-05",
				Cost:   "10",
			},
			{
				Region: "eu-west-2",
				Team:   "T01",
				Date:   "2025-04",
				Cost:   "100",
			},
			{
				Region: "eu-west-2",
				Team:   "T01",
				Date:   "2025-04",
				Cost:   "2000",
			},
			{
				Region: "eu-west-3",
				Team:   "T02",
				Date:   "2025-04",
				Cost:   "2000",
			},
		}
	)

	tbl, err := ListToTable(log, lcfg, testRecords)
	if err != nil {
		t.Errorf("unexpected error: %v", err.Error())
	}
	if len(tbl) != expectedLen {
		t.Errorf("unexpected len; expected [%d] actual [%v]", expectedLen, len(tbl))
	}

	// all 2025-03 records should be 0.000
	for _, row := range tbl {
		if v, ok := row["2025-03"]; !ok || v != lcfg.DefaultValue {
			t.Errorf("row 2025-03 entry is invalid: [%v]", row)
		}
	}
}

// TestApiTablesPopulate tests the populate is working with test data
// - only positive path
func TestApiTablesPopulate(t *testing.T) {
	var (
		expectedLen                     = 5
		dates                           = []string{"2025-03", "2025-04", "2025-05"}
		identifiers                     = []string{"team", "region"}
		records     []map[string]string = []map[string]string{}
		pcfg                            = &PopulateConfig{
			TextColumns:     identifiers,
			ValueField:      "cost",
			DataSourceField: "date",
			DefaultValue:    "0.000",
		}
		testRecords = []*testRecord{
			{
				Region: "eu-west-1",
				Team:   "T01",
				Date:   "2025-04",
				Cost:   "100",
			},
			{
				Region: "eu-west-1",
				Team:   "T01",
				Date:   "2025-05",
				Cost:   "100",
			},
			{
				Region: "eu-west-1",
				Team:   "T02",
				Date:   "2025-05",
				Cost:   "2000",
			},
			{
				Region: "eu-west-2",
				Team:   "T00",
				Date:   "2025-05",
				Cost:   "10",
			},
			{
				Region: "eu-west-2",
				Team:   "T01",
				Date:   "2025-04",
				Cost:   "100",
			},
			{
				Region: "eu-west-2",
				Team:   "T01",
				Date:   "2025-04",
				Cost:   "2000",
			},
			{
				Region: "eu-west-3",
				Team:   "T02",
				Date:   "2025-04",
				Cost:   "2000",
			},
		}
	)

	utils.Convert(testRecords, &records)
	keys, _ := PossibleCombinationsAsKeys(records, identifiers)
	pcfg.Skeleton = Skeleton(keys, dates, pcfg.DefaultValue)

	populated := Populate(pcfg, records)
	if len(populated) != expectedLen {
		t.Errorf("unexpected len; expected [%d] actual [%v]", expectedLen, len(populated))
	}

	// all 2025-03 records should be 0.000
	for _, row := range populated {
		if v, ok := row["2025-03"]; !ok || v != pcfg.DefaultValue {
			t.Errorf("row 2025-03 entry is invalid: [%v]", row)
		}
	}

}

// TestApiTablesSkeleton check the skel has the correct number of rows
// and correct placeholder values
func TestApiTablesSkeleton(t *testing.T) {
	var (
		dates                             = []string{"2025-03", "2025-04", "2025-05"}
		expectedCount                     = (3 * 3) // 3 teams, 3 regions
		identifiers                       = []string{"team", "region"}
		records       []map[string]string = []map[string]string{}
		testRecords                       = []*testRecord{
			{
				Team:   "T00",
				Region: "eu-west-2",
				Date:   "2025-03",
				Cost:   "10",
			},
			{
				Team:   "T01",
				Region: "eu-west-1",
				Date:   "2025-04",
				Cost:   "100",
			},

			{
				Team:   "T01",
				Region: "eu-west-2",
				Date:   "2025-04",
				Cost:   "100",
			},
			{
				Team:   "T01",
				Region: "eu-west-1",
				Date:   "2025-05",
				Cost:   "100",
			},
			{
				Team:   "T02",
				Region: "eu-west-3",
				Date:   "2025-04",
				Cost:   "2000",
			},
			{
				Team:   "T02",
				Region: "eu-west-2",
				Date:   "2025-04",
				Cost:   "2000",
			},
			{
				Team:   "T02",
				Region: "eu-west-1",
				Date:   "2025-05",
				Cost:   "2000",
			},
		}
	)

	utils.Convert(testRecords, &records)
	keys, _ := PossibleCombinationsAsKeys(records, identifiers)

	skel := Skeleton(keys, dates, "--")
	if len(skel) != expectedCount {
		t.Errorf("unexpected row count in table, expected [%d] actual [%v]", expectedCount, len(skel))
	}
	// now check every row has every date key and its 0
	for _, row := range skel {
		for _, month := range dates {
			if val, ok := row[month]; !ok || val != "--" {
				t.Errorf("error in row with this month [%s] [%v]", month, row)
			}
		}
	}
}

// TestApiTablesPossibleCombinationsAsKeys simple test to check combinations are generated correctly
func TestApiTablesPossibleCombinationsAsKeys(t *testing.T) {
	var err error
	var testRecords = []*testRecord{
		{
			Team:   "T01",
			Region: "eu-west-1",
			Date:   "2025-04",
			Cost:   "100",
		},
		{
			Team:   "T01",
			Region: "eu-west-2",
			Date:   "2025-04",
			Cost:   "100",
		},
		{
			Team:   "T01",
			Region: "eu-west-1",
			Date:   "2025-05",
			Cost:   "100",
		},
		{
			Team:   "T02",
			Region: "eu-west-3",
			Date:   "2025-04",
			Cost:   "200",
		},
		{
			Team:   "T02",
			Region: "eu-west-2",
			Date:   "2025-04",
			Cost:   "200",
		},
		{
			Team:   "T02",
			Region: "eu-west-1",
			Date:   "2025-05",
			Cost:   "200",
		},
	}
	var keyLen = (3 * 2) // 3 regions used, 2 teams
	var identifiers = []string{"team", "region"}
	var records []map[string]string = []map[string]string{}

	err = utils.Convert(testRecords, &records)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
		t.FailNow()
	}

	keys, _ := PossibleCombinationsAsKeys(records, identifiers)
	if len(keys) != keyLen {
		t.Errorf("key combination mismatch, expected [%d] actual [%v]", keyLen, len(keys))
	}
}
