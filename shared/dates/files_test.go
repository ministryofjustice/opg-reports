package dates_test

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/ministryofjustice/opg-reports/shared/dates"
)

func TestSharedDatesFileCreationTime(t *testing.T) {

	dir := t.TempDir()
	// dir := testhelpers.Dir()
	created := "2024-02-29T01:02:00"
	// expected := "2024-02-29T01:02:00"
	f := filepath.Join(dir, "test-creation")
	fmt.Println(f)
	cmd := exec.Command("touch", "-d", created, f)
	err := cmd.Run()
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
	if actual != created {
		t.Errorf("created date failed")
		fmt.Println(actual)
		fmt.Println(created)
	}

}
