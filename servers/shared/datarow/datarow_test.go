package datarow_test

import (
	"testing"
	"time"

	"github.com/ministryofjustice/opg-reports/servers/shared/datarow"
	"github.com/ministryofjustice/opg-reports/shared/dates"
)

func TestSharedDatarowRowsSkeleton(t *testing.T) {

	s := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	e := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)
	months := dates.Strings(dates.Range(s, e, dates.MONTH), dates.FormatYM)
	intervals := map[string][]string{"interval": months}

	columns := map[string][]string{
		"unit":    {"foo"},
		"env":     {"prod"},
		"account": {"1"},
	}
	skel := datarow.Skeleton(columns, intervals)
	if len(skel) != 1 {
		t.Errorf("permuations incorrect")
	}

	columns = map[string][]string{
		"unit":    {"foo", "bar"},
		"env":     {"prod", "dev"},
		"account": {"1", "2", "3"},
	}
	// math to determine combinations
	l := len(columns["unit"]) * len(columns["env"]) * len(columns["account"])
	skel = datarow.Skeleton(columns, intervals)
	if len(skel) != l {
		t.Errorf("permuations incorrect")
	}

}

func TestSharedDatarowRowsDataToRows(t *testing.T) {
	s := time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC)
	e := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)
	months := dates.Strings(dates.Range(s, e, dates.MONTH), dates.FormatYM)
	intervals := map[string][]string{"interval": months}

	columns := map[string][]string{
		"unit":    {"foo", "bar"},
		"env":     {"prod", "dev"},
		"account": {"1", "2", "3"},
	}

	data := []map[string]interface{}{
		map[string]interface{}{
			"unit":     "foo",
			"env":      "prod",
			"account":  "1",
			"interval": "2024-05",
			"cost":     10.5,
		},
		map[string]interface{}{
			"unit":     "foo",
			"env":      "dev",
			"account":  "3",
			"interval": "2024-06",
			"cost":     65,
		},
		map[string]interface{}{
			"unit":     "bar",
			"env":      "dev",
			"account":  "2",
			"interval": "2024-06",
			"cost":     50,
		},
	}

	values := map[string]string{
		"interval": "cost",
	}

	rows := datarow.DataToRows(data, columns, intervals, values)

	l := len(columns["unit"]) * len(columns["env"]) * len(columns["account"])
	if len(rows) != l {
		t.Errorf("permuations incorrect")
	}

}

// larger dataset
func TestSharedDatarowRowsDataToRowsRealistic(t *testing.T) {
	values := map[string]string{
		"interval": "total",
	}
	months := []string{
		"2023-11",
		"2023-12",
		"2024-01",
		"2024-02",
		"2024-03",
		"2024-04",
		"2024-05",
		"2024-06",
		"2024-07",
	}
	intervals := map[string][]string{"interval": months}

	columns := map[string][]string{
		"unit": []string{
			"Foo1",
			"OldFoo",
			"FooBar",
			"BarFoo",
			"FooOwner",
			"WhichFoo",
			"StarFoo",
			"WhatBar",
		},
	}

	data := []map[string]interface{}{
		map[string]interface{}{
			"interval": "2023-11",
			"total":    3681.1324005098,
			"unit":     "Foo1",
		},
		map[string]interface{}{
			"interval": "2023-11",
			"total":    175.01637666780002,
			"unit":     "OldFoo",
		},
		map[string]interface{}{
			"interval": "2023-11",
			"total":    2575.3563783735,
			"unit":     "FooBar",
		},
		map[string]interface{}{
			"interval": "2023-11",
			"total":    6609.6168719351,
			"unit":     "BarFoo",
		},
		map[string]interface{}{
			"interval": "2023-11",
			"total":    5224.5261678834,
			"unit":     "FooOwner",
		},
		map[string]interface{}{
			"interval": "2023-11",
			"total":    1399.9822676583,
			"unit":     "WhichFoo",
		},
		map[string]interface{}{
			"interval": "2023-11",
			"total":    32976.8213704883,
			"unit":     "StarFoo",
		},
		map[string]interface{}{
			"interval": "2023-11",
			"total":    2579.7163149466,
			"unit":     "WhatBar",
		},
		map[string]interface{}{
			"interval": "2023-12",
			"total":    3021.8102651718,
			"unit":     "Foo1",
		},
		map[string]interface{}{
			"interval": "2023-12",
			"total":    179.5049332877,
			"unit":     "OldFoo",
		},
		map[string]interface{}{
			"interval": "2023-12",
			"total":    1880.2939871138,
			"unit":     "FooBar",
		},
		map[string]interface{}{
			"interval": "2023-12",
			"total":    2478.8424690372,
			"unit":     "BarFoo",
		},
		map[string]interface{}{
			"interval": "2023-12",
			"total":    4433.7531660984005,
			"unit":     "FooOwner",
		},
		map[string]interface{}{
			"interval": "2023-12",
			"total":    929.1960804954,
			"unit":     "WhichFoo",
		},
		map[string]interface{}{
			"interval": "2023-12",
			"total":    33801.3687152191,
			"unit":     "StarFoo",
		},
		map[string]interface{}{
			"interval": "2023-12",
			"total":    2232.5727563481,
			"unit":     "WhatBar",
		},
		map[string]interface{}{
			"interval": "2024-01",
			"total":    3060.9926280196,
			"unit":     "Foo1",
		},
		map[string]interface{}{
			"interval": "2024-01",
			"total":    180.2262436464,
			"unit":     "OldFoo",
		},
		map[string]interface{}{
			"interval": "2024-01",
			"total":    2105.2164487917,
			"unit":     "FooBar",
		},
		map[string]interface{}{
			"interval": "2024-01",
			"total":    2581.0161401119,
			"unit":     "BarFoo",
		},
		map[string]interface{}{
			"interval": "2024-01",
			"total":    5121.0144704104005,
			"unit":     "FooOwner",
		},
		map[string]interface{}{
			"interval": "2024-01",
			"total":    946.2213903969,
			"unit":     "WhichFoo",
		},
		map[string]interface{}{
			"interval": "2024-01",
			"total":    30672.4060697091,
			"unit":     "StarFoo",
		},
		map[string]interface{}{
			"interval": "2024-01",
			"total":    2795.9872205306,
			"unit":     "WhatBar",
		},
		map[string]interface{}{
			"interval": "2024-02",
			"total":    3046.3105692933,
			"unit":     "Foo1",
		},
		map[string]interface{}{
			"interval": "2024-02",
			"total":    177.8805391635,
			"unit":     "OldFoo",
		},
		map[string]interface{}{
			"interval": "2024-02",
			"total":    1929.2532781536,
			"unit":     "FooBar",
		},
		map[string]interface{}{
			"interval": "2024-02",
			"total":    2823.1421482973,
			"unit":     "BarFoo",
		},
		map[string]interface{}{
			"interval": "2024-02",
			"total":    4809.6650860542,
			"unit":     "FooOwner",
		},
		map[string]interface{}{
			"interval": "2024-02",
			"total":    884.9780573077,
			"unit":     "WhichFoo",
		},
		map[string]interface{}{
			"interval": "2024-02",
			"total":    29696.516459829898,
			"unit":     "StarFoo",
		},
		map[string]interface{}{
			"interval": "2024-02",
			"total":    3082.0846058193,
			"unit":     "WhatBar",
		},
		map[string]interface{}{
			"interval": "2024-03",
			"total":    3094.5595165393,
			"unit":     "Foo1",
		},
		map[string]interface{}{
			"interval": "2024-03",
			"total":    183.8363489543,
			"unit":     "OldFoo",
		},
		map[string]interface{}{
			"interval": "2024-03",
			"total":    2181.4068602739,
			"unit":     "FooBar",
		},
		map[string]interface{}{
			"interval": "2024-03",
			"total":    4975.7330031414,
			"unit":     "BarFoo",
		},
		map[string]interface{}{
			"interval": "2024-03",
			"total":    4334.4893436939,
			"unit":     "FooOwner",
		},
		map[string]interface{}{
			"interval": "2024-03",
			"total":    953.4483683153,
			"unit":     "WhichFoo",
		},
		map[string]interface{}{
			"interval": "2024-03",
			"total":    32873.0957829386,
			"unit":     "StarFoo",
		},
		map[string]interface{}{
			"interval": "2024-03",
			"total":    2903.9871663564,
			"unit":     "WhatBar",
		},
		map[string]interface{}{
			"interval": "2024-04",
			"total":    2967.6449618315,
			"unit":     "Foo1",
		},
		map[string]interface{}{
			"interval": "2024-04",
			"total":    184.5640378139,
			"unit":     "OldFoo",
		},
		map[string]interface{}{
			"interval": "2024-04",
			"total":    2142.5678722909,
			"unit":     "FooBar",
		},
		map[string]interface{}{
			"interval": "2024-04",
			"total":    5423.3646691351005,
			"unit":     "BarFoo",
		},
		map[string]interface{}{
			"interval": "2024-04",
			"total":    4362.074764721,
			"unit":     "FooOwner",
		},
		map[string]interface{}{
			"interval": "2024-04",
			"total":    961.8384561468,
			"unit":     "WhichFoo",
		},
		map[string]interface{}{
			"interval": "2024-04",
			"total":    33939.5472399987,
			"unit":     "StarFoo",
		},
		map[string]interface{}{
			"interval": "2024-04",
			"total":    2901.1051098394,
			"unit":     "WhatBar",
		},
		map[string]interface{}{
			"interval": "2024-05",
			"total":    3020.5348264837,
			"unit":     "Foo1",
		},
		map[string]interface{}{
			"interval": "2024-05",
			"total":    191.4955916516,
			"unit":     "OldFoo",
		},
		map[string]interface{}{
			"interval": "2024-05",
			"total":    2170.2902759644,
			"unit":     "FooBar",
		},
		map[string]interface{}{
			"interval": "2024-05",
			"total":    4627.3278626821,
			"unit":     "BarFoo",
		},
		map[string]interface{}{
			"interval": "2024-05",
			"total":    4439.7370499576,
			"unit":     "FooOwner",
		},
		map[string]interface{}{
			"interval": "2024-05",
			"total":    785.1327349612,
			"unit":     "WhichFoo",
		},
		map[string]interface{}{
			"interval": "2024-05",
			"total":    34737.1321425157,
			"unit":     "StarFoo",
		},
		map[string]interface{}{
			"interval": "2024-05",
			"total":    3159.978892263,
			"unit":     "WhatBar",
		},
		map[string]interface{}{
			"interval": "2024-06",
			"total":    2697.4038323496,
			"unit":     "Foo1",
		},
		map[string]interface{}{
			"interval": "2024-06",
			"total":    193.3475672784,
			"unit":     "OldFoo",
		},
		map[string]interface{}{
			"interval": "2024-06",
			"total":    1948.8610525628,
			"unit":     "FooBar",
		},
		map[string]interface{}{
			"interval": "2024-06",
			"total":    3987.3137863393,
			"unit":     "BarFoo",
		},
		map[string]interface{}{
			"interval": "2024-06",
			"total":    3940.4556692908,
			"unit":     "FooOwner",
		},
		map[string]interface{}{
			"interval": "2024-06",
			"total":    765.8581850346,
			"unit":     "WhichFoo",
		},
		map[string]interface{}{
			"interval": "2024-06",
			"total":    34369.0074814626,
			"unit":     "StarFoo",
		},
		map[string]interface{}{
			"interval": "2024-06",
			"total":    3473.9052610188,
			"unit":     "WhatBar",
		},
		map[string]interface{}{
			"interval": "2024-07",
			"total":    2698.1273535541,
			"unit":     "Foo1",
		},
		map[string]interface{}{
			"interval": "2024-07",
			"total":    197.0553485397,
			"unit":     "OldFoo",
		},
		map[string]interface{}{
			"interval": "2024-07",
			"total":    2128.8244090243,
			"unit":     "FooBar",
		},
		map[string]interface{}{
			"interval": "2024-07",
			"total":    5141.5668000052,
			"unit":     "BarFoo",
		},
		map[string]interface{}{
			"interval": "2024-07",
			"total":    5032.8104955813,
			"unit":     "FooOwner",
		},
		map[string]interface{}{
			"interval": "2024-07",
			"total":    773.7805316959,
			"unit":     "WhichFoo",
		},
		map[string]interface{}{
			"interval": "2024-07",
			"total":    37371.2112964623,
			"unit":     "StarFoo",
		},
		map[string]interface{}{
			"interval": "2024-07",
			"total":    3701.7667495253,
			"unit":     "WhatBar",
		},
	}

	rows := datarow.DataToRows(data, columns, intervals, values)

	// for _, row := range rows {
	// 	js := convert.PrettyString(row)
	// 	fmt.Println(js)
	// }
	l := len(columns["unit"])
	if len(rows) != l {
		t.Errorf("permuations incorrect")
	}
}
