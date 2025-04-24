package smap

import (
	"fmt"
	"testing"
)

// https://labix.org/gocheck
// . "gopkg.in/check.v1"

type Ri struct {
	R rune
	I int
}

type SMap01DataType struct {
	Test     string
	Sigma    string
	NoMatch  rune
	SkipTest bool
	RIData   []Ri
	Len      int
}

var SMap01Data = []SMap01DataType{
	{Test: "0000", Sigma: "ab", NoMatch: rune(1000), SkipTest: false, Len: 3, RIData: []Ri{Ri{'a', 0}, Ri{'b', 1}, Ri{'z', 2}, Ri{'\u2022', 2}}},
	{Test: "0001", Sigma: "a\u2022z", NoMatch: rune(1000), SkipTest: false, Len: 4, RIData: []Ri{Ri{'a', 0}, Ri{'b', 3}, Ri{'z', 1}, Ri{'\u2022', 2}}},
	{Test: "0002", Sigma: "abcdef\u2022\u2716xyz", NoMatch: rune(0xf812), SkipTest: false, Len: 12, RIData: []Ri{Ri{'a', 0}, Ri{'\u2716', 10}, Ri{'\uf0f8', 11}, Ri{'{', 11}}},
	{Test: "0003", Sigma: "defacb\u2716\u2022xyz", NoMatch: rune(0xf812), SkipTest: false, Len: 12, RIData: []Ri{Ri{'a', 0}, Ri{'b', 1}, Ri{'A', 11}, Ri{'~', 11}, Ri{'`', 11}}},
}

/*

->#}£»ï<-
->£»ïï£º<-
->%\{£»ïï£º<-
->[]£»ï<-
->"\£»ïï£º<-
->'\£»ïï£º<-
->`£»ïï£º<-
->!"#%&')*+,-./<=>\^`{|}£»ïï£µï£¸ï£¹<-
->!"#%&')*+,-./:<=>\^`{|}£»ïï£µï£¸ï£¹<-

*/

func TestLexie(t *testing.T) {

	t_debug := false

	for ii, vv := range SMap01Data {
		if !vv.SkipTest {
			if t_debug {
				fmt.Printf("Test: %d ---------------------------------------------------------------\n", ii)
			}
			t1 := NewSMapType(vv.Sigma, vv.NoMatch)
			if t_debug {
				fmt.Printf("Results: %v\n", t1)
			}
			for jj, ww := range vv.RIData {
				// t1. MapRune(rn rune) int {
				rx := t1.MapRune(ww.R)
				if rx != ww.I {
					if t_debug {
						fmt.Printf("  Faild %d,%d : rn=%04x %s expected = %d result = %d\n", ii, jj, ww.R, string(ww.R), ww.I, rx)
					}
					// c.Check(ww.I, Equals, rx)
					if ww.I != rx {
						t.Errorf("Failed %d,%d:  Expected %v got %v\n", ii, jj, ww.I, rx)
					}
				}
			}
			// func (smap *SMapType) Length() int {
			if vv.Len > 0 {
				// c.Check(vv.Len, Equals, t1.Length())
				if vv.Len != t1.Length() {
					t.Errorf("Failed:  Length did not match, expeced, %d got %d\n", vv.Len, t1.Length())
				}
			}
		}
	}
}

/* vim: set noai ts=4 sw=4: */
