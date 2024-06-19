package cost

import (
	"opg-reports/shared/data"
	"testing"
)

func TestSharedAwsCostUUID(t *testing.T) {
	cost := New(nil)
	if len(cost.UUID) <= 0 {
		t.Errorf("failed to create uuid")
	}

	if cost.UID() != cost.UUID {
		t.Errorf("idx doesnt match uuid")
	}

	u := "000-01-01"
	cost = New(&u)
	if cost.UUID != u {
		t.Errorf("UUID not set correctly")
	}
}

func TestSharedAwsCostValid(t *testing.T) {
	cost := New(nil)

	if cost.Valid() {
		t.Errorf("empty, should not be valid")
	}
	cost.AccountOrganisation = "org"
	if cost.Valid() {
		t.Errorf("mostly empty, should not be valid")
	}

	cost.AccountId = "01"
	cost.AccountEnvironment = "dev"
	cost.AccountLabel = "test"
	cost.AccountName = "random"
	cost.AccountUnit = "unit"
	cost.Service = "test-service"
	cost.Date = "2024-01"
	cost.Region = "eu"

	if cost.Valid() {
		t.Errorf("mostly full, should not be valid")
	}

	cost.Cost = "1"
	if !cost.Valid() {
		t.Errorf("complete, should not be valid")
	}

}

func testIsI[V data.IEntry](i V) bool {
	return i.UID() != ""
}
func TestSharedAwsCostInterface(t *testing.T) {
	c := &Cost{}
	if testIsI[*Cost](c) {
		t.Errorf("should not be nil")
	}
}
