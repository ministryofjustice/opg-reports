package compliance

import (
	"fmt"
	"testing"
)

func TestSharedGhComplianceToRow(t *testing.T) {
	item := Fake(nil)

	fmt.Printf("%+v\n", item)

	row := ToRow(item)
	for _, c := range row.GetCells() {
		fmt.Printf("[%s]=%s\n", c.Name, c.Value)
	}

}
