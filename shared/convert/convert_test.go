package convert_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/ministryofjustice/opg-reports/shared/convert"
	"github.com/ministryofjustice/opg-reports/shared/testhelpers"
)

// test various forms of basic type swapping
func TestSharedConvertMarshalUnmarshal(t *testing.T) {

	// -- simple struct with public members, make sure match up
	now := time.Now().UTC()
	ts := &testhelpers.Ts{S: now, E: now}
	tsM, err := convert.Marshal(ts)
	if err != nil {
		t.Errorf("failed to marshal :%s", err.Error())
	}
	tsU, err := convert.Unmarshal[*testhelpers.Ts](tsM)
	if err != nil {
		t.Errorf("failed to unmarshal :%s", err.Error())
	}
	if tsU.E != ts.E || tsU.S != ts.S {
		t.Errorf("marshall, unmarshal failed")
		fmt.Printf("%+v\n", ts)
		fmt.Printf("%+v\n", tsU)
	}

	// -- handle array of simple structs
	aTs := []*testhelpers.Ts{
		{S: now, E: now},
		{S: now, E: now},
	}
	aTtsM, err := convert.Marshals(aTs)
	if err != nil {
		t.Errorf("failed to marshal :%s", err.Error())
	}
	aTsU, err := convert.Unmarshals[*testhelpers.Ts](aTtsM)
	if err != nil {
		t.Errorf("failed to unmarshal :%s", err.Error())
	}
	if len(aTsU) != len(aTs) {
		t.Errorf("failed to convert back")
		fmt.Printf("%+v\n", aTs)
		fmt.Printf("%+v\n", aTsU)
	}

}

// test converting struct back and forth from a map
func TestSharedConvertMapUnmap(t *testing.T) {
	now := time.Now().UTC()
	ts := &testhelpers.Ts{S: now, E: now}

	tsM, err := convert.Map(ts)
	if err != nil {
		t.Errorf("failed to map :%s", err.Error())
	}
	if _, ok := tsM["start"]; !ok {
		t.Errorf("failed to map")
	}
	tsU, err := convert.Unmap[*testhelpers.Ts](tsM)
	if err != nil {
		t.Errorf("failed to unmap :%s", err.Error())
	}

	if tsU.S != ts.S || tsU.E != ts.E {
		t.Errorf("failed to convert back")
		fmt.Printf("%+v\n", ts)
		fmt.Printf("%+v\n", tsU)
	}

	all := []*testhelpers.Ts{
		{S: now, E: now}, {S: now, E: now}, {S: now, E: now}, {S: now, E: now},
	}
	maps, err := convert.Maps(all)
	if err != nil {
		t.Errorf("failed to map :%s", err.Error())
	}
	if len(maps) != len(all) {
		t.Errorf("failed to convert slice of struct to slice of maps")
	}
	// now covnert back
	u, err := convert.Unmaps[*testhelpers.Ts](maps)
	if err != nil {
		t.Errorf("failed to unmap :%s", err.Error())
	}
	if len(u) != len(all) {
		t.Errorf("failed to unmap")
	}

}

func TestSharedConvertString(t *testing.T) {

	s := &testhelpers.Simple{Name: "test name"}
	expected := `{"name": "test name"}`
	str := convert.String(s)
	if expected != str {
		t.Errorf("got an invalid result")
		fmt.Println(str)
	}

}

func TestSharedConvertPrettyString(t *testing.T) {
	s := &testhelpers.Simple{Name: "test name"}
	str := convert.PrettyString(s)

	if len(str) <= 0 {
		t.Errorf("pretty string failed")
	}

}

func TestSharedConvertBools(t *testing.T) {

	for _, i := range []int{0, 2, 5, 10, 11} {
		if convert.IntToBool(i) {
			t.Errorf("should be false: [%d]", i)
		}
	}
	if !convert.IntToBool(1) {
		t.Errorf("should be true")
	}

	if convert.BoolToInt(true) != 1 {
		t.Error("true should be 1")
	}
	if convert.BoolToInt(false) == 1 {
		t.Error("false should not be 1")
	}

	if convert.BoolStringToInt("true") != 1 {
		t.Error("'true' should be 1")
	}
	if convert.BoolStringToInt("false") == 1 {
		t.Error("'false' should not be 1")
	}

}

func TestSharedConvertTitle(t *testing.T) {
	src := "unit_a"
	expected := "Unit A"
	actual := convert.Title(src)
	if actual != expected {
		t.Errorf("error converting to tile: [%s]", actual)
	}

}

func TestSharedConvertDict(t *testing.T) {

	d := convert.Dict(1, "test", 2)
	if len(d) > 0 {
		t.Errorf("dict should have failed with odd number of arguments")
	}

	d = convert.Dict("A", "test", "B", "foo")
	if len(d) != 2 {
		t.Errorf("dict did not create map correctly")
	}
	if v, ok := d["A"]; !ok || v != "test" {
		t.Errorf("dict did not return values correctly")
	}

}

func TestSharedConvertCurr(t *testing.T) {

	if convert.Curr(1, "$") != "$0.0" {
		t.Errorf("curr should return default when int is passed")
	}

	if a := convert.Curr("12.356", "£"); a != "£12.36" {
		t.Errorf("curr did not match expected value [%s]", a)
	}
	if a := convert.Curr(12.354, "£"); a != "£12.35" {
		t.Errorf("curr did not match expected value [%s]", a)
	}

}

func TestSharedConvertStripIntPrefix(t *testing.T) {

	if a := convert.StripIntPrefix("1.Test.3"); a != "Test3" {
		t.Errorf("stript int failed: [%s]", a)
	}
	if a := convert.StripIntPrefix("Test 1"); a != "Test 1" {
		t.Errorf("stript int failed with nothign to remove: [%s]", a)
	}

}

func TestSharedConvertPercent(t *testing.T) {

	if p := convert.Percent(50, 100); p != "50.00" {
		t.Errorf("convert failed: [%s]", p)
	}

}
