//
// S M A P - Part of Lexie Lexical Generation System
//
// (C) Philip Schlump, 2014-2015.
// Version: 1.0.8
// BuildNo: 203
//
// Special Thanks to 2C-Why, LLC for supporting this project.
//

package smap

//
// SMap - Map from the input char set into an interal table representation
//
// An input rune 'a' gets maped into a 0..n value to be used in the table lookup DFA(t).
//

import (
	"fmt"
	"sort"
	"unicode"

	"github.com/pschlump/dbgo"
	"github.com/pschlump/lexie/re"
)

// ------------------------------------------------------------------------------------------------------------------------------------------------------

type SMapType struct {
	MinV    int          // Subtract before procesisng - min(sigma)
	MaxV    int          // Range that iwill be maped in M0
	M0      []int        // Lower range
	M1      map[rune]int // Lookup for outliers
	NoMap   rune         // Rune used for no maping
	NoMapTo int          // Index for NoMap rune
	Len     int          //
	SigmaRN []rune
}

func NewSMapType(sigma string, noMapCh rune) (rv *SMapType) {
	v := &SMapType{
		NoMap:   noMapCh,
		MinV:    0,
		MaxV:    0,
		NoMapTo: 0,
		SigmaRN: make([]rune, 0, len(sigma)),
	}
	//for _, rn := range sigma {
	//	v.SigmaRN = append(v.SigmaRN, rn)
	//}
	dbgo.DbPrintf("smap", "NewSMapType: %q %x\n", sigma, noMapCh)
	nCh, iCh, nnCh := 0, 0, 0
	var mCh rune
	ts := make([]rune, 0, len(sigma))
	for _, c := range sigma {
		ts = append(ts, c)
	}
	ts = KeyRuneSort(ts)
	v.SigmaRN = ts
	v.SigmaRN = append(v.SigmaRN, re.R_else_CH)
	// fmt.Printf("ts=%v\n", ts) // Should be length of sigma in runes
	if len(ts) > 0 {
		v.MinV = int(ts[0])
		//if v.MinV > 0x80 {
		//	v.MinV = 0x80
		//}
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
		v.MaxV = int(mCh) - v.MinV
		// fmt.Printf("n,m,iCh = %d, %d, %d len(v.M0)=%d\n", nCh, mCh, iCh, v.MaxV+1)
		if (int(mCh) - v.MinV + 1) > 0 {
			v.M0 = make([]int, int(mCh)-v.MinV+1, int(mCh)-v.MinV+1)
		} else {
			v.M0 = make([]int, 1, 1)
		}
		for kk := range v.M0 {
			v.M0[kk] = v.NoMapTo
		}
		// fmt.Printf("M0, initialized with subscript of no-match = %v\n", v.M0)
		for _, vv := range ts {
			if vv < 0xff {
				v.M0[int(vv)-v.MinV] = nnCh
				nnCh++
			} else {
				break
			}
		}
		// fmt.Printf("M0, final = %+v\n", v.M0)
		v.M1 = make(map[rune]int)
		// -------------------------------------------- ok to this poin ---------------------------------------
		if v.MinV > 0x80 {
			iCh--
		}
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

func (smap *SMapType) ReverseMapRune(x int) rune {
	if x > smap.MaxV {
		for k, v := range smap.M1 {
			if v == x {
				return k
			}
		}
	}
	return rune(x + smap.MinV)
}

// Map an input rune to one of the possible output subscripts
func (smap *SMapType) MapRune(rn rune) int {
	x := int(rn) - smap.MinV
	// fmt.Printf("x=%d, rn=%04x ( %s ),  %s\n", x, rn, string(rn), dbgo.LF())
	if x > smap.MaxV {
		// fmt.Printf("At %s\n", dbgo.LF())
		if y, ok := smap.M1[rn]; ok {
			// fmt.Printf("At %s\n", dbgo.LF())
			return y
		}
		// fmt.Printf("At %s\n", dbgo.LF())
		return smap.NoMapTo
	} else if x >= 0 {
		dbgo.DbPrintf("smap", "********************************** this one **************************, %s\n", string(rn))
		v := smap.M0[x]
		if v == smap.NoMapTo {
			dbgo.DbPrintf("smap", "********************************** No MAP - 28 case \n")
			if unicode.IsDigit(rn) {
				if y, ok := smap.M1[re.R_NUMERIC]; ok {
					dbgo.DbPrintf("smap", "********************************** case - numeric \n")
					return y
				}
			} else if unicode.IsUpper(rn) {
				if y, ok := smap.M1[re.R_UPPER]; ok {
					dbgo.DbPrintf("smap", "********************************** case - upper\n")
					return y
				}
			} else if unicode.IsLower(rn) {
				if y, ok := smap.M1[re.R_LOWER]; ok {
					dbgo.DbPrintf("smap", "********************************** case - lower\n")
					return y
				}
			}
		}
		// fmt.Printf("At %s\n", dbgo.LF())
		return v
	} else {
		// fmt.Printf("At %s\n", dbgo.LF())
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
		s += fmt.Sprintf("    %2d: [ %3d 0x%x %7q ] = %d\n", ii, ii, ii+smap.MinV, string(rune(ii+smap.MinV)), smap.M0[ii])
	}
	s += fmt.Sprintf("smap.M1 = (%d items)\n", len(smap.M1))
	kv := KeyRuneMapSort(smap.M1)
	for _, ii := range kv {
		vv := smap.M1[ii]
		s += fmt.Sprintf("    %s: %d", string(ii), vv)
		switch ii {
		case re.R_DOT:
			s += " X_DOT     \\uF8FA " // Any char in Sigma
		case re.R_BOL:
			s += " X_BOL     \\uF8F3 " // Beginning of line
		case re.R_EOL:
			s += " X_EOL     \\uF8F4 " // End of line
		case re.R_NUMERIC:
			s += " X_NUMERIC \\uF8F5 "
		case re.R_LOWER:
			s += " X_LOWER   \\uF8F6 "
		case re.R_UPPER:
			s += " X_UPPER   \\uF8F7 "
		case re.R_ALPHA:
			s += " X_ALPHA   \\uF8F8 "
		case re.R_ALPHNUM:
			s += " X_ALPHNUM \\uF8F9 "
		case re.R_EOF:
			s += " X_EOF     \\uF8FB "
		//case re.R_not_CH:
		//	s += " X_not_CH  \\uF8FC " // On input lookup if the char is NOT in Signa then it is returned as this.
		case re.R_else_CH:
			s += " X_else_CH \\uF8FC " // If char is not matched in this state then take this path

		}
		s += "\n"
	}
	s += fmt.Sprintf("}}}\n")
	return s
}

// Display SmapType as a human readable content.
//func (smap *SMapType) StringOrig() string {
//	s := ""
//	s += fmt.Sprintf("{{{ smap.MinV = %d, smap.MaxV = %d, smap.Len = %d\n", smap.MinV, smap.MaxV, smap.Len)
//	s += fmt.Sprintf("smap.NoMap = %04x ( %s ), smap.NoMapTo = %d\n", smap.NoMap, string(smap.NoMap), smap.NoMapTo)
//	s += fmt.Sprintf("smap.M0 = (%d items)\n", len(smap.M0))
//	for ii := range smap.M0 {
//		s += fmt.Sprintf("    %2d: [ %3d 0x%x %s ] = %d\n", ii, ii, ii+smap.MinV, string(rune(ii+smap.MinV)), smap.M0[ii])
//	}
//	s += fmt.Sprintf("smap.M1 = (%d items)\n", len(smap.M1))
//	kv := KeyIntMapSort(smap.M1)
//	for _, ii := range kv {
//		vv := smap.M1[ii]
//		s += fmt.Sprintf("    %s: %d\n", string(ii), vv )
//	}
//	s += fmt.Sprintf("}}}\n")
//	return s
//}

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

func KeyRuneMapSort(in map[rune]int) (rv []rune) {
	var rt []int
	for ii, _ := range in {
		rt = append(rt, int(ii))
	}
	KeyIntSort(rt)
	for _, vv := range rt {
		rv = append(rv, rune(vv))
	}
	return
}

func KeyRuneSort(in []rune) (rv []rune) {
	rr := make([]int, len(in), len(in))
	for i := 0; i < len(in); i++ {
		rr[i] = int(in[i])
	}
	sort.Sort(sort.IntSlice(rr))
	for i := 0; i < len(rr); i++ {
		rv = append(rv, rune(rr[i]))
	}
	return
}

/* vim: set noai ts=4 sw=4: */
