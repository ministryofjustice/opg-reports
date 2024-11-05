package navigation

import (
	"testing"
)

func TestNavigationRoot(t *testing.T) {

	var (
		lv4  = New("1.1.1.1", "/1/1/1/1")
		lv3  = New("1.1.1", "/1/1/1", lv4)
		lv2  = New("1.1", "/1/1", lv3)
		lv1  = New("1", "/1", lv2)
		root *Navigation
	)

	root = Root(lv4)
	if root.Name != lv1.Name {
		t.Errorf("found incorrect root path")
	}

	root = Root(lv1)
	if root.Name != lv1.Name {
		t.Errorf("found incorrect root path from self")
	}

}
