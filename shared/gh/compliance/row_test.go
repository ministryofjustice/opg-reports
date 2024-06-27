package compliance

import (
	"opg-reports/shared/data"
	"testing"
)

func TestSharedGhComplianceToRow(t *testing.T) {
	item := Fake(nil)
	row := ToRow(item)

	nItem := FromRow(row)
	if nItem.UID() != item.UID() {
		t.Errorf("failed to mathc uid")
	}

	if nItem.Archived != item.Archived {
		t.Errorf("failed to match archived")
	}

	s1, _ := data.ToJson(item)
	s2, _ := data.ToJson(nItem)

	if string(s1) != string(s2) {
		t.Errorf("failed to match")
	}

}
