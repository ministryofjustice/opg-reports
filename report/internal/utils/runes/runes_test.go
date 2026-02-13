package runes

import "testing"

type tester struct {
	Start    rune
	Expected rune
}

func TestRunesNext(t *testing.T) {

	var tests = []*tester{
		{Start: 'Z', Expected: 'A'},
		{Start: 'A', Expected: 'B'},
	}

	for _, test := range tests {
		actual := Next(test.Start)
		if test.Expected != actual {
			t.Errorf("next rune error; expected [%v] actual [%v]", test.Expected, actual)
		}
	}
}
