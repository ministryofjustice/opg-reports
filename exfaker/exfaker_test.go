package exfaker_test

import (
	"testing"

	"github.com/go-faker/faker/v4"
	"github.com/ministryofjustice/opg-reports/exfaker"
)

type fTest struct {
	ID   int     `json:"id,omitempty" db:"id" faker:"unique, boundary_start=1, boundary_end=20"`
	Ts   string  `json:"ts,omitempty" db:"ts"  faker:"time_string"`
	Cost float64 `json:"cost,omitempty" db:"cost" faker:"float"`
}

// TestExFakerExtended checks that ID is within bounds and
// the custom float faker stays within its bounds as well
func TestExFakerExtended(t *testing.T) {
	exfaker.AddProviders()

	var f = &fTest{}
	faker.FakeData(f)

	if f.ID < 1 || f.ID > 20 {
		t.Errorf("ID out of allowed range: [%v]", f.ID)
	}
	if f.Cost < exfaker.FloatMin || f.Cost > exfaker.FloatMax {
		t.Errorf("cost generated out of bounds")
	}

}

// TestExFakerMany checks that many generates different
// structs each time by generation almost as many records
// as ID field has range for and checks that there are
// no repeating IDs
func TestExFakerMany(t *testing.T) {
	exfaker.AddProviders()
	var n = 15

	many := exfaker.Many[fTest](n)

	ids := map[int]int{}
	for _, m := range many {
		if _, ok := ids[m.ID]; !ok {
			ids[m.ID] = 1
		} else {
			ids[m.ID] += 1
		}
	}

	for id, count := range ids {
		if count > 1 {
			t.Errorf("ID [%d] used more than one [%d], should be unique", id, count)
		}
	}
}
