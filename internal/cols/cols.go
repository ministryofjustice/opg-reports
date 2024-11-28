package cols

import (
	"log/slog"
	"slices"

	"github.com/ministryofjustice/opg-reports/internal/structs"
)

// Values finds all the unique values within rows passed for each of the columns, returning them
// as a map.
func Values[T any](rows []T, columns []string) (values map[string][]interface{}) {
	slog.Debug("[ColumnValues] get started ...")
	values = map[string][]interface{}{}

	for _, row := range rows {
		mapped, err := structs.ToMap(row)
		if err != nil {
			slog.Error("[ColumnValues] get to map failed", slog.String("err", err.Error()))
			return
		}

		for _, column := range columns {
			// if not set, set it
			if _, ok := values[column]; !ok {
				values[column] = []interface{}{}
			}
			// add the value into the slice
			if rowValue, ok := mapped[column]; ok {
				// if they arent in there already
				if !slices.Contains(values[column], rowValue) {
					values[column] = append(values[column], rowValue)
				}
			}

		}
	}
	return
}
