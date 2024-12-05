package transformers

import (
	"fmt"
	"log/slog"
	"slices"

	"github.com/ministryofjustice/opg-reports/internal/structs"
)

type dateDeepTable interface {
	DateDeepDateColumn() string
}

// recordToDeepRow takes the row and appends the data into the existingData based on the key
// which allows appending (and addition of value) into the dataset
func recordToDeepRow[T dateDeepTable](item T, columns []string, existingData map[string]map[string]interface{}) (key string, err error) {
	var (
		ok          bool
		asMap       map[string]interface{} = map[string]interface{}{}
		existingRow map[string]interface{} = map[string]interface{}{}
	)

	if asMap, err = structs.ToMap(item); err != nil {
		slog.Error("[transformers] recordToDeepRow convert failed")
		return
	}
	// this is generated id this cost item would use in the possible list
	key = RowKV(columns, asMap)
	// look for the existing data within the dataset
	// - if cant find it, error
	existingRow, ok = existingData[key]
	if !ok {
		err = fmt.Errorf("failed to find existing data with key [%s]", key)
		return
	}

	// now we try and set values on the existingRow from this cost by using the columns
	for _, field := range columns {
		existingRow[field] = asMap[field]
	}

	return
}

// ResultsToDeepRows uses the columnValues (from the api) and appends the date field (which will be missing) into that
// data set so we can show a depth rather than width (so date column rather than a column per date value)
func ResultsToDeepRows[T dateDeepTable](apiData []T, columnValues map[string][]interface{}, dateRange []string) (dataAsMap map[string]map[string]interface{}, err error) {

	var (
		// columns is sorted column names only - this is to ensure 'key' order is a match
		columns    []string = SortedColumnNames(columnValues)
		dateColumn string
		// found tracks which 'key' has real data and inserted in to the data map
		// so anything that is not in this list can be removed - as it
		// will not have and values
		found []string = []string{}
	)

	if len(apiData) > 0 {
		// append the date column
		dateColumn = apiData[0].DateDeepDateColumn()
		columns = append(columns, dateColumn)
		slices.Sort(columns)
		// append the dates in column values
		columnValues[dateColumn] = []interface{}{}
		for _, date := range dateRange {
			columnValues[dateColumn] = append(columnValues[dateColumn], date)
		}
	}

	dataAsMap = TableSkeleton(columnValues)

	for _, item := range apiData {
		rowKey, e := recordToDeepRow(item, columns, dataAsMap)
		if e != nil {
			slog.Error("[transformers] recordToDeepRow failed", slog.String("err", e.Error()))
			return
		}
		// insert to the list of done rows
		if !slices.Contains(found, rowKey) {
			found = append(found, rowKey)
		}
	}

	// remove any row that has not been marked as 'done' - these are empty combinations
	for key := range dataAsMap {
		if !slices.Contains(found, key) {
			delete(dataAsMap, key)
		}
	}

	return
}
