package cols_test

import (
	"fmt"
	"testing"

	"github.com/ministryofjustice/opg-reports/internal/cols"
)

type dummy struct {
	ID           int    `json:"id,omitempty" db:"id"`
	Label        string `json:"label,omitempty" db:"label" faker:"word"`
	Organisation string `json:"organisation,omitempty"`
	Unit         string `json:"unit,omitempty"`
}

// checks that all permutations of column
// values are found correctly for the data
func TestColsValues(t *testing.T) {

	simple := []*dummy{
		{Organisation: "A", Unit: "A", Label: "One"},
		{Organisation: "A", Unit: "B", Label: "One"},
		{Organisation: "A", Unit: "C", Label: "One"},
		{Organisation: "A", Unit: "1", Label: "Three"},
		{Organisation: "A", Unit: "Z", Label: "One"},
		{Organisation: "A", Unit: "Z", Label: "Four"},
	}

	list := []string{"organisation", "unit", "label"}

	actual := cols.Values(simple, list)
	expected := map[string]int{"organisation": 1, "unit": 5, "label": 3}
	for col, count := range expected {
		l := len(actual[col])
		if l != count {
			t.Errorf("[%s] values dont match - expected [%d] actual [%v]", col, count, l)
			fmt.Println(actual[col])
		}
	}

}
