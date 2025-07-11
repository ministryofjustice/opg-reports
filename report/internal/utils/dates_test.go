package utils

import (
	"testing"
	"time"
)

type dateTest struct {
	Test     time.Time
	Expected time.Time
}

func TestLastDayOfMonth(t *testing.T) {
	var tests = []*dateTest{
		{
			Test:     time.Date(2024, 1, 15, 9, 0, 5, 0, time.UTC),
			Expected: time.Date(2024, 1, 31, 23, 59, 59, 0, time.UTC),
		},
		{
			Test:     time.Date(2024, 2, 15, 9, 0, 5, 0, time.UTC),
			Expected: time.Date(2024, 2, 29, 23, 59, 59, 0, time.UTC),
		},
	}

	for _, test := range tests {
		actual := LastDayOfMonth(test.Test)
		if test.Expected != actual {
			t.Errorf("last day mismatch, expected: [%s] actual [%s]", test.Expected.Format(DATE_FORMATS.Full), actual.Format(DATE_FORMATS.Full))
		}
	}
}

func TestBillingMonth(t *testing.T) {
	var tests = []*dateTest{
		{
			Test:     time.Date(2025, 7, 11, 9, 0, 0, 0, time.UTC),
			Expected: time.Date(2025, 5, 31, 23, 59, 59, 0, time.UTC),
		},
		{
			Test:     time.Date(2025, 7, 15, 9, 0, 0, 0, time.UTC),
			Expected: time.Date(2025, 6, 30, 23, 59, 59, 0, time.UTC),
		},
		{
			Test:     time.Date(2025, 7, 16, 9, 0, 0, 0, time.UTC),
			Expected: time.Date(2025, 6, 30, 23, 59, 59, 0, time.UTC),
		},
	}

	for _, test := range tests {
		actual := BillingMonth(test.Test, 15)
		if test.Expected != actual {
			t.Errorf("billing month mismatch, expected: [%s] actual [%s]", test.Expected.Format(DATE_FORMATS.Full), actual.Format(DATE_FORMATS.Full))
		}
	}
}
