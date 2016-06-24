package eval

import (
	"testing"

	"github.com/pschlump/lexie/tok"
	// "../../../go-lib/sizlib"
)

func Test_St02(t *testing.T) {
	// func MapIsEmpty(t map[string]tok.Token) bool {
	var tt map[string]tok.Token

	b := MapIsEmpty(tt)
	if !b {
		t.Errorf("misc_test: failed to return true on an un-initialized map, MapIsEmpty()\n")
	}

	tt = make(map[string]tok.Token)
	b = MapIsEmpty(tt)
	if !b {
		t.Errorf("misc_test: failed to return true on an initialized empty map, MapIsEmpty()\n")
	}

	tt["x"] = tok.Token{}
	b = MapIsEmpty(tt)
	if b {
		t.Errorf("misc_test: failed to return false on aninitialized non-empty map, MapIsEmpty()\n")
	}

	// t = BoundArrayIndex ( t, 0, len(eval.Mm) )
	// func BoundArrayIndex ( i, min, max int ) int {

	ii := -1
	ii = BoundArrayIndex(ii, 0, 5)
	if ii != 0 {
		t.Errorf("misc_test: Bound failed to fix, expected 0 got %d\n", ii)
	}

	ii = 0
	ii = BoundArrayIndex(ii, 0, 5)
	if ii != 0 {
		t.Errorf("misc_test: Bound failed to fix, expected 0 got %d\n", ii)
	}

	ii = 4
	ii = BoundArrayIndex(ii, 0, 5)
	if ii != 4 {
		t.Errorf("misc_test: Bound failed to fix, expected 4 got %d\n", ii)
	}

	ii = 5
	ii = BoundArrayIndex(ii, 0, 5)
	if ii != 4 {
		t.Errorf("misc_test: Bound failed to fix, expected 4 got %d\n", ii)
	}
}
