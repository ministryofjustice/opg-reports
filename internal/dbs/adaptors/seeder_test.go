package adaptors

import "testing"

func TestAdaptorsSeeder(t *testing.T) {
	var s = &Seed{}

	if s.Seedable() {
		t.Errorf("seedable should default to false")
	}

	s = &Seed{seedable: true}
	if !s.Seedable() {
		t.Errorf("seedable should return true")
	}

	s.Seeded()
	if s.Seedable() {
		t.Errorf("seedable should return false after being marked as seeded.")
	}

}
