package sqlx

import "testing"

type testJSONWith struct {
	ID        int `json:"id"`
	privateID int `json:"-"`
}

type testJSONWithout struct {
	ID        int
	privateID int `json:"-"`
}

type testJSONNestedWith struct {
	ID   string       `json:"-"`
	With testJSONWith `json:"with"`
}
type testJSONNestedWithout struct {
	ID      string          `json:"-"`
	With    testJSONWith    `json:"with"`
	Without testJSONWithout `json:"without"`
}

func TestSQLxAllFieldsHaveJSONTags(t *testing.T) {
	var valid bool
	// this struct has all tags - test it as ptr and value
	valid = AllFieldsHaveJSONTags(NewReflection(&testJSONWith{}))
	if !valid {
		t.Errorf("unexpected error - struct should have all json tags")
	}
	valid = AllFieldsHaveJSONTags(NewReflection(testJSONWith{}))
	if !valid {
		t.Errorf("unexpected error - struct should have all json tags")
	}
	// now try nested structs that works
	valid = AllFieldsHaveJSONTags(NewReflection(&testJSONNestedWith{}))
	if !valid {
		t.Errorf("unexpected error - struct should have all json tags")
	}
	valid = AllFieldsHaveJSONTags(NewReflection(testJSONNestedWith{}))
	if !valid {
		t.Errorf("unexpected error - struct should have all json tags")
	}

	// this does not have all json tags
	valid = AllFieldsHaveJSONTags(NewReflection(&testJSONWithout{}))
	if valid {
		t.Errorf("unexpected error - struct does not have all json tags")
	}
	valid = AllFieldsHaveJSONTags(NewReflection(testJSONWithout{}))
	if valid {
		t.Errorf("unexpected error - struct does not all json tags")
	}
	// nested struct that is missing some should not work
	valid = AllFieldsHaveJSONTags(NewReflection(&testJSONNestedWithout{}))
	if valid {
		t.Errorf("unexpected error - struct does not have all json tags")
	}
	valid = AllFieldsHaveJSONTags(NewReflection(testJSONNestedWithout{}))
	if valid {
		t.Errorf("unexpected error - struct does not all json tags")
	}

	// a non struct should also fail
	valid = AllFieldsHaveJSONTags(NewReflection("not a struct"))
	if valid {
		t.Errorf("unexpected error - string should not pass this check")
	}

}
