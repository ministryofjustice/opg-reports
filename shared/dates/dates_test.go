package dates

import (
	"testing"
	"time"
)

func TestSharedDatesGetFormat(t *testing.T) {

	test := time.Now().UTC().Format(Format)

	if GetFormat(test[0:4]) != "2006" || GetFormat(FormatY) != FormatY {
		t.Errorf("year format mismatch")
	}
	if GetFormat(test[0:7]) != "2006-01" || GetFormat(FormatYM) != FormatYM {
		t.Errorf("year month format mismatch")
	}
	if GetFormat(test[0:10]) != "2006-01-02" || GetFormat(FormatYMD) != FormatYMD {
		t.Errorf("year month day format mismatch")
	}
	ltest := Format + ".000"
	if GetFormat(ltest) != Format {
		t.Errorf("max length swap failed: [%v]->[%v] [%v]", ltest, GetFormat(ltest), Format)
	}

}

func TestSharedDatesStringToDate(t *testing.T) {
	valid := []string{
		"2024", "2024-01", "2024-02-29", "2024-03-01T00:00",
	}
	for _, val := range valid {
		_, err := StringToDate(val)
		if err != nil {
			t.Errorf("error converting string to date: [%s] -> [%v]", val, err.Error())
		}
	}

	invalid := []string{
		"2023-02-29", "2024-03-01T00:60",
	}
	for _, val := range invalid {
		_, err := StringToDate(val)
		if err == nil {
			t.Errorf("did not get expected error when converting [%v]", val)
		}
	}
}

func TestSharedDatesStrings(t *testing.T) {
	dates := []time.Time{
		time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2023, 9, 1, 0, 0, 0, 0, time.UTC),
	}

	months := Strings(dates, FormatYM)

	if months[0] != "2024-01" {
		t.Errorf("first date incorrect")
	}
	if months[1] != "2023-09" {
		t.Errorf("second date incorrect")
	}
}

func TestSharedDatestringToDateDefault(t *testing.T) {

	ds := time.Now().UTC().Format(FormatYMD)

	d, _ := StringToDateDefault("-", "-", ds)
	if d.Format(FormatYMD) != ds {
		t.Errorf("default date failed")
	}

}
