package fakes

import (
	"slices"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestSharedFakeString(t *testing.T) {

	str := String(5)

	if len(str) != 5 {
		t.Errorf("incorrect length")
	}
	// check all parts are within charset
	for i := 0; i < len(str); i++ {
		x := str[i:i]
		if !strings.Contains(charset, x) {
			t.Errorf("char not in charset: %s", x)
		}
	}
}

func TestSharedFakeInt(t *testing.T) {
	i := Int(-10, 10)
	if i > 10 || i < -10 {
		t.Errorf("int out of range: %d", i)
	}
}

func TestSharedFakeIntAsStr(t *testing.T) {
	is := IntAsStr(-10, 10)
	i, err := strconv.Atoi(is)
	if err != nil {
		t.Errorf("error converting back to int")
	}
	if i > 10 || i < -10 {
		t.Errorf("int out of range: %d", i)
	}
}

func TestSharedFakeFloat(t *testing.T) {
	f := Float(-2.0, 55.0)
	if f > 55.0 || f < -2.0 {
		t.Errorf("float out of range: %v", f)
	}
}

func TestSharedFakeFloatAsStr(t *testing.T) {
	fs := FloatAsStr(-10.0, 55.0)
	f, err := strconv.ParseFloat(fs, 10)
	if err != nil {
		t.Error("error converting back to float")
	}

	if f > 55.0 || f < -10.0 {
		t.Errorf("float out of range: %v", f)
	}
}

func TestSharedFakeDate(t *testing.T) {
	max := time.Now().UTC()
	min := time.Date(max.Year()-2, max.Month(), 1, 0, 0, 0, 0, time.UTC)

	d := Date(min, max)
	if d.Before(min) || d.After(max) {
		t.Errorf("date is out of range")
	}

}

func TestSharedFakeDateAsStr(t *testing.T) {
	max := time.Now().UTC()
	min := time.Date(max.Year()-2, max.Month(), 1, 0, 0, 0, 0, time.UTC)
	f := time.RFC3339
	ds := DateAsStr(min, max, f)
	d, er := time.Parse(f, ds)
	if er != nil {
		t.Errorf("date parse error")
	}
	if d.Before(min) || d.After(max) {
		t.Errorf("date is out of range")
	}

}

func TestSharedFakeChoices(t *testing.T) {
	single := []string{"one"}
	item := Choice(single)
	if item != single[0] {
		t.Errorf("failed to pick single choice")
	}

	many := []int{}
	for i := 0; i < 100; i++ {
		many = append(many, i)
	}
	for x := 0; x < 1000; x++ {
		picked := Choice(many)
		if !slices.Contains(many, picked) {
			t.Errorf("picked out of range")
		}
	}
}
