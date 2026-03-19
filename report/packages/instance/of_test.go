package instance

import (
	"fmt"
	"testing"
)

type testS struct {
	Name string `json:"name"`
}

func TestPackagesInstanceOf(t *testing.T) {

	// test simple pointer
	dest := Of[*testS]()
	if fmt.Sprintf("%T", dest) != "*instance.testS" {
		t.Errorf("type mismatch")
		t.FailNow()
	}
	dest.Name = "test"

	if dest.Name != "test" {
		t.Errorf("setting value failed")
	}

	// test a base type
	i := Of[int]()
	// fmt.Printf("%T => %v\n", i, i)
	if fmt.Sprintf("%T", i) != "int" {
		t.Errorf("type mismatch")
		t.FailNow()
	}

}
