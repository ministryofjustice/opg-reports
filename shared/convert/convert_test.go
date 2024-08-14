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
