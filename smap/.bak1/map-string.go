package smap

//
// SMap - Map from the input char set into an interal table representation
//
// An input rune 'a' gets maped into a 0..n value to be used in the table lookup DFA(t).
//
//
//

import (
	"fmt"
	"sort"
)

// ------------------------------------------------------------------------------------------------------------------------------------------------------
// 1. Build the maping info InputMap0, Ranges, etc.			1hr
// 2. Test this with current machines						1hr
// 3. Generate the array-based machine.						2hr
// 4. Hand verify - write an output function
// 5. Check this.										++	4hr
//
// Let's build the map stuff as it's own little project, ./smap
//	fx:		minv, maxv, m0, m1 := smap.BuildMapString ( Sigma, NoMapRn )
//	fx:		k := smap.MapRune ( rn )
// ------------------------------------------------------------------------------------------------------------------------------------------------------

type SMapType struct {
	MinV    int         // Subtract before procesisng - min(sigma)
	MaxV    int         // Range that iwill be maped in M0
	M0      []int       // Lower range
	M1      map[int]int // Lookup for outliers
	NoMap   int         // Rune used for no maping
	NoMapTo int         // Index for NoMap rune
	Len     int         //
}

func NewSMapType(sigma string, noMapCh rune) (rv *SMapType) {
	v := &SMapType{
		NoMap:   int(noMapCh),
		MinV:    0,
		MaxV:    0,
		NoMapTo: 0,
	}
	fmt.Printf("NewSMapType: %q %x\n", sigma, noMapCh)
	nCh, mCh, iCh, nnCh := 0, 0, 0, 0
	ts := make([]int, 0, len(sigma))
	for _, c := range sigma {
		ts = append(ts, int(c))
	}
	ts = KeyIntSort(ts)
	// fmt.Printf("ts=%v\n", ts) // Should be length of sigma in runes
	if len(ts) > 0 {
		v.MinV = ts[0]
		v.NoMapTo = len(ts)
		nCh, mCh, iCh, nnCh = 0, 0, 0, 0
		for ii, vv := range ts {
			if vv < 0xff {
				nCh++
				mCh = vv
				iCh = ii
			} else {
				break
			}
		}
		v.MaxV = mCh - v.MinV
		// fmt.Printf("n,m,iCh = %d, %d, %d len(v.M0)=%d\n", nCh, mCh, iCh, v.MaxV+1)
		v.M0 = make([]int, mCh-v.MinV+1, mCh-v.MinV+1)
		for kk := range v.M0 {
			v.M0[kk] = v.NoMapTo
		}
		// fmt.Printf("M0, initialized with subscript of no-match = %v\n", v.M0)
		for _, vv := range ts {
			if vv < 0xff {
				v.M0[vv-v.MinV] = nnCh
				nnCh++
			} else {
				break
			}
		}
		// fmt.Printf("M0, final = %+v\n", v.M0)
		v.M1 = make(map[int]int)
		// -------------------------------------------- ok to this poin ---------------------------------------
		v.M1[v.NoMap] = v.NoMapTo // Add the Special Char
		for jj := iCh + 1; jj < len(ts); jj++ {
			v.M1[ts[jj]] = nnCh
			// fmt.Printf("Adding %d to M1 as %d\n", ts[jj], nnCh)
			nnCh++
		}
		// fmt.Printf("M1 = %v\n", v.M1)
		// fmt.Printf("nnCh, final = %d\n", nnCh)
		v.Len = nnCh + 1
	}
	return v
}

// Map an input rune to one of the possible output subscripts
func (smap *SMapType) MapRune(rn rune) int {
	x := int(rn) - smap.MinV
	// fmt.Printf("x=%d, rn=%04x ( %s ),  %s\n", x, rn, string(rn), tr.LF())
	if x > smap.MaxV {
		// fmt.Printf("At %s\n", tr.LF())
		if y, ok := smap.M1[int(rn)]; ok {
			// fmt.Printf("At %s\n", tr.LF())
			return y
		}
		// fmt.Printf("At %s\n", tr.LF())
		return smap.NoMapTo
	} else if x >= 0 {
		// fmt.Printf("At %s\n", tr.LF())
		return smap.M0[x]
	} else {
		// fmt.Printf("At %s\n", tr.LF())
		return smap.NoMapTo
	}
}

// return the number of possible output subscripts.
func (smap *SMapType) Length() int {
	return smap.Len
}

// Display SmapType as a human readable content.
func (smap *SMapType) String() string {
	s := ""
	s += fmt.Sprintf("{{{ smap.MinV = %d, smap.MaxV = %d, smap.Len = %d\n", smap.MinV, smap.MaxV, smap.Len)
	s += fmt.Sprintf("smap.NoMap = %04x ( %s ), smap.NoMapTo = %d\n", smap.NoMap, string(smap.NoMap), smap.NoMapTo)
	s += fmt.Sprintf("smap.M0 = (%d items)\n", len(smap.M0))
	for ii := range smap.M0 {
		s += fmt.Sprintf("    %2d: [ %3d 0x%x %s ] = %d\n", ii, ii, ii+smap.MinV, string(rune(ii+smap.MinV)), smap.M0[ii])
	}
	s += fmt.Sprintf("smap.M1 = (%d items)\n", len(smap.M1))
	kv := KeyIntMapSort(smap.M1)
	for _, ii := range kv {
		vv := smap.M1[ii]
		s += fmt.Sprintf("    %2d: [ %3d 0x%x %s ] =  Value        %d\n", ii, ii, ii+smap.MinV, string(rune(ii+smap.MinV)), vv)
	}
	s += fmt.Sprintf("}}}\n")
	return s
}

func KeyIntMapSort(in map[int]int) []int {
	var rv []int
	for ii, _ := range in {
		rv = append(rv, ii)
	}
	return KeyIntSort(rv)
}

func KeyIntSort(in []int) (rv []int) {
	rv = in
	sort.Sort(sort.IntSlice(rv))
	return
}

/* vim: set noai ts=4 sw=4: */
