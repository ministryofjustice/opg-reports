package dates_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/ministryofjustice/opg-reports/shared/dates"
)

func TestSharedDatesFileCreationTime(t *testing.T) {

	dir := t.TempDir()
	f := filepath.Join(dir, "test-time")
	os.WriteFile(f, []byte(""), os.ModePerm)

	now := time.Now().UTC()
	m := time.Date(2024, 2, 29, 1, 1, 0, 0, time.UTC)

	err := os.Chtimes(f, now, m)

	if err != nil {
		t.Errorf("error forcing time change, this might be os related!")
		fmt.Println(err)
	}

	ts, err := dates.FileCreationTime(f)
	if err != nil {
		t.Errorf("error getting time, this might be os related!")
		fmt.Println(err)
	}
	actual := ts.Format(dates.FormatYMDHMS)
	expected := m.Format(dates.FormatYMDHMS)
	if actual != expected {
		t.Errorf("created date failed")
		fmt.Println(actual)
		fmt.Println(expected)
	}

}
