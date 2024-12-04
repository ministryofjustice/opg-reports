package navigation

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNavigationActivateTree(t *testing.T) {
	var req = httptest.NewRequest(http.MethodGet, "/1/2/?test=true", nil)

	var kid = New("1.2.1", "/1/2/1")
	var shouldBeActive = New("1.2", "/1/2", kid)
	var node1 = New("1", "/1", New("1.1", "/1/1"), shouldBeActive)
	var node2 = New("2", "/2", New("2.1", "/2/1"))
	var nav = []*Navigation{
		node1,
		node2,
	}

	ActivateTree(nav, req)

	// the active item should be active
	if shouldBeActive.Display.IsActive != true {
		t.Errorf("error matching uri exactly to nav item")
	}
	if shouldBeActive.Display.InUri != true {
		t.Errorf("the active item should also be marked as being in the uri")
	}
	// the child of the active item should not be
	if kid.Display.IsActive || kid.Display.InUri {
		t.Errorf("child of active nav should not be active or shown as being in the active stack")
	}

	// the parent of the active item should be in the stakc, but not active
	if node1.Display.InUri != true {
		t.Errorf("parent of the active item should be marked as in the active stack")
	}
	if node1.Display.IsActive == true {
		t.Errorf("parent should not be considered active, only in the active stack")
	}
	// top level sibling should not be active at all
	if node2.Display.IsActive || node2.Display.InUri {
		t.Errorf("top level sibling should not be active in any way")
	}

}

// TestNavigationFlatPreDetermined checks that using a
// preset nav structure we get the correct number of
// elements in the flat version
func TestNavigationFlatPreDetermined(t *testing.T) {

	var expected = 8 // manually set so we dont have to recurse
	var flat = map[string]*Navigation{}
	var nav = []*Navigation{
		{
			Name: "root1",
			Uri:  "/1",
			children: []*Navigation{
				{
					Name: "1.1",
					Uri:  "/1/1",
				},
				{
					Name: "1.2",
					Uri:  "/1/2",
					children: []*Navigation{
						{
							Name: "1.2.1",
							Uri:  "/1/2/1",
						},
					},
				},
			},
		},
		{
			Name: "root2",
			Uri:  "/2",
			children: []*Navigation{
				{Name: "2.1", Uri: "/2/1"},
				{Name: "2.2", Uri: "/2/2"},
				{Name: "2.3", Uri: "/2/3"},
			},
		},
	}

	Flat(nav, flat)

	actual := len(flat)
	if expected != actual {
		t.Errorf("flat length in correct - expected [%d] actual [%v]", expected, actual)
	}
}
