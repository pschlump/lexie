package pbread

import (
	"os"
	"strings"
	"testing"

	"../com"

	"../../../go-lib/tr"

	// "../../../go-lib/tr"
	// "../../../go-lib/sizlib"
)

type ActType struct {
	ActionType          int    // 1 = filename, 2 = string, 3 = rune, 4 = runeArray, 5 = get X from Next, 6 = peek, 7 = output, 8 = byte array, 9 = pbfile, 10 file seen, 11 get pos, 12 set pos, 20 pb x
	Fn                  string // 18 = do output, 19 = truckate output to make testing easier, 17 pushback macro test ( CH=.Rn to replace with Str=Data )
	Data                string //
	Rn                  rune   //
	RnS                 []rune //
	X                   int    //
	FileSeenFlag        bool   //
	IntermediateResults string //
	IntermediateRune    rune   //
	LineNo              int    //
	ColNo               int    //
}

type Pb01TestType struct {
	Test     string
	SkipTest bool
	Actions  []ActType
	Results  string
}

var Pb01Test = []Pb01TestType{

	// Simple test with 2 fiels
	{Test: "0001", SkipTest: false, Actions: []ActType{
		ActType{ActionType: 1, Fn: "test/t1.txt"},
		ActType{ActionType: 5, X: 4},
		ActType{ActionType: 1, Fn: "test/t2.txt"},
		ActType{ActionType: 7},
	}, Results: `a12
b34
ccc
ddd
EEE
FFF
GGG
`},

	// test with file and 2 push back runes
	{Test: "0002", SkipTest: false, Actions: []ActType{
		ActType{ActionType: 1, Fn: "test/t1.txt"},
		ActType{ActionType: 5, X: 4},
		ActType{ActionType: 3, Rn: 'x'},
		ActType{ActionType: 3, Rn: 'y'},
		ActType{ActionType: 7},
	}, Results: `a12
xyb34
ccc
ddd
`},

	// push back rune array
	{Test: "0003", SkipTest: false, Actions: []ActType{
		ActType{ActionType: 1, Fn: "test/t1.txt"},
		ActType{ActionType: 5, X: 4},
		ActType{ActionType: 4, RnS: []rune{'x', 'y', 'z', 'w', 'v'}},
		ActType{ActionType: 7},
	}, Results: `a12
xyzwvb34
ccc
ddd
`},

	// push back string
	{Test: "0004", SkipTest: false, Actions: []ActType{
		ActType{ActionType: 1, Fn: "test/t1.txt"},
		ActType{ActionType: 5, X: 4},
		ActType{ActionType: 2, Data: "xyzwv"},
		ActType{ActionType: 7},
	}, Results: `a12
xyzwvb34
ccc
ddd
`},

	// push back a number of items
	{Test: "0005", SkipTest: false, Actions: []ActType{
		ActType{ActionType: 1, Fn: "test/t1.txt"},
		ActType{ActionType: 6, Rn: 'a'},
		ActType{ActionType: 5, X: 4},
		ActType{ActionType: 2, Data: "xyzwv"},
		ActType{ActionType: 6, Rn: 'x'},
		ActType{ActionType: 5, X: 4},
		ActType{ActionType: 6, Rn: 'v'},
		ActType{ActionType: 7},
	}, Results: `a12
xyzwvb34
ccc
ddd
`},

	// Pull some off before push back
	{Test: "0006", SkipTest: false, Actions: []ActType{
		ActType{ActionType: 1, Fn: "test/t1.txt"},
		ActType{ActionType: 5, X: 4},
		ActType{ActionType: 8, Data: "XyZwV"},
		ActType{ActionType: 7},
	}, Results: `a12
XyZwVb34
ccc
ddd
`},

	// output empty buffer
	{Test: "0007", SkipTest: false, Actions: []ActType{
		ActType{ActionType: 7},
	}, Results: ``},

	// test file names
	{Test: "0008", SkipTest: false, Actions: []ActType{
		ActType{ActionType: 11, Fn: "", LineNo: 1, ColNo: 1},
		ActType{ActionType: 1, Fn: "test/t1.txt"},
		ActType{ActionType: 5, X: 4},
		ActType{ActionType: 11, Fn: "test/t1.txt", LineNo: 2, ColNo: 1},
		ActType{ActionType: 5, X: 1},
		ActType{ActionType: 11, Fn: "test/t1.txt", LineNo: 2, ColNo: 2},
		ActType{ActionType: 5, X: 3},
		ActType{ActionType: 11, Fn: "test/t1.txt", LineNo: 3, ColNo: 1},
		ActType{ActionType: 9, Fn: "test/t2.txt"},
		ActType{ActionType: 11, Fn: "test/t2.txt", LineNo: 1, ColNo: 1},
		ActType{ActionType: 10, FileSeenFlag: true, Fn: "test/t1.txt"},
		ActType{ActionType: 10, FileSeenFlag: true, Fn: "test/t2.txt"},
		ActType{ActionType: 10, FileSeenFlag: false, Fn: "test/t3.txt"},
		ActType{ActionType: 7},
		ActType{ActionType: 11, Fn: "test/t1.txt", LineNo: 5, ColNo: 1},
		ActType{ActionType: 12, Fn: "test/t8.txt", LineNo: 55, ColNo: 202},
		ActType{ActionType: 11, Fn: "test/t8.txt", LineNo: 55, ColNo: 202},
	}, Results: `a12
b34
EEE
FFF
GGG
ccc
ddd
`},

	// test with missing file - see if get error case - both open and PbFile
	{Test: "0009", SkipTest: false, Actions: []ActType{
		ActType{ActionType: 1, Fn: "test/x1.txt"},
		ActType{ActionType: 9, Fn: "test/x2.txt"},
	}, Results: ``},

	// test push back that overflows PbAFew's max
	{Test: "0010", SkipTest: false, Actions: []ActType{
		ActType{ActionType: 20, X: 622},           // Push back 622 chars ( more than buffer size )
		ActType{ActionType: 18},                   // Output
		ActType{ActionType: 1, Fn: "test/t1.txt"}, // Tak test/t1.txt on end
		ActType{ActionType: 5, X: 622},            // INput 622 chars
		ActType{ActionType: 20, X: 622},           // Push back 622 chars ( more than buffer size )
		ActType{ActionType: 5, X: 622},            // INput 622 chars
		ActType{ActionType: 19},                   // Discard current input so test is simpler
		ActType{ActionType: 7},                    // Get input - should be file
	}, Results: "a12\nb34\nccc\nddd\n"}, //

	{Test: "0011", SkipTest: false, Actions: []ActType{
		ActType{ActionType: 1, Fn: "test/t1.txt"},                          //	Input file
		ActType{ActionType: 5, X: 4},                                       //	Pick off 4 chars at beginning
		ActType{ActionType: 4, RnS: []rune{'x', 'y', 'z', 'w', 'v'}},       //  Push back "xyzwv" - 5 runes
		ActType{ActionType: 12, Fn: "test/t8.txt", LineNo: 55, ColNo: 202}, // Set file name to t8.txt,line55,col202
		ActType{ActionType: 11, Fn: "test/t8.txt", LineNo: 55, ColNo: 202}, // Check that we have those values
		ActType{ActionType: 5, X: 3},                                       //  Forward by 3
		ActType{ActionType: 11, Fn: "test/t8.txt", LineNo: 55, ColNo: 205}, //  Check col changed
		ActType{ActionType: 5, X: 4},                                       // Forward by 4 more runes
		ActType{ActionType: 18},                                            // Output before the final test -- Just one buffer left if OK
		ActType{ActionType: 11, Fn: "test/t1.txt", LineNo: 2, ColNo: 3},    //
		ActType{ActionType: 7},                                             //
	}, Results: `a12
xyzwvb34
ccc
ddd
`},

	// Simple macro replace
	{Test: "0012", SkipTest: false, Actions: []ActType{
		ActType{ActionType: 1, Fn: "test/t4.txt"}, // Input file
		ActType{ActionType: 17, Rn: 'M', Data: "aaaDvvv"},
		ActType{ActionType: 17, Rn: 'D', Data: "ddd"},
		ActType{ActionType: 17, Rn: 'E', Data: "e"},
		ActType{ActionType: 17, Rn: 'J', Data: "**error**"},
		ActType{ActionType: 16}, // abaaadddvvvcdddmem
	}, Results: `abaaavvvdddcdddmem
abdddcdddmÏm
`},

	// test combo files/ pb/ peek etc.
	// test with a list of files in order - pb.OpenFile()
	// test with a set of files in PbFile()
}

type MacroDefTestType struct {
	Rn   rune   // Defined Macro
	Body string // Replace with
}

func Test_PbBufer01(t *testing.T) {

	SymbolTable := make([]*MacroDefTestType, 0, 100)
	Define := func(name rune, body string) {
		for ii := 0; ii < len(SymbolTable); ii++ {
			// fmt.Printf("Search at %d, for %s\n", ii, string(name))
			if SymbolTable[ii] != nil && SymbolTable[ii].Rn == name {
				SymbolTable[ii].Body = body
				return
			}
		}
		// fmt.Printf("Append\n")
		SymbolTable = append(SymbolTable, &MacroDefTestType{Rn: name, Body: body})
	}
	ResetST := func() {
		// SymbolTable = make([]*MacroDefTestType, 0, 100)
		SymbolTable = SymbolTable[:1]
	}
	HaveMacro := func(name rune) (body string, found bool) {
		body = ""
		found = false
		for ii := 0; ii < len(SymbolTable); ii++ {
			if SymbolTable[ii] != nil && SymbolTable[ii].Rn == name {
				body, found = SymbolTable[ii].Body, true
				return
			}
		}
		return
	}

	for ii, vv := range Pb01Test {

		if !vv.SkipTest {

			// Implement a quick - fetch execute macine to test - the PbBuffer
			ss := ""
			pb := NewPbRead()
			ResetST()
			for kk, ww := range vv.Actions {

				switch ww.ActionType {
				case 1:
					pb.OpenFile(ww.Fn)
					com.DbPrintf("testCode", "Open file %s At: %s\n", ww.Fn, tr.LF())
					if com.DbOn("testDump") {
						pb.Dump01(os.Stdout)
					}
				case 2: // Push back a rune array
					pb.PbString(ww.Data)
					if com.DbOn("testDump") {
						pb.Dump01(os.Stdout)
					}
				case 3: // Push back a rune
					pb.PbRune(ww.Rn)
					if com.DbOn("testDump") {
						pb.Dump01(os.Stdout)
					}
				case 4: // Push back a rune array
					pb.PbRuneArray(ww.RnS)
					if com.DbOn("testDump") {
						pb.Dump01(os.Stdout)
					}
				case 5:
					for ll := 0; ll < ww.X; ll++ {
						rn, done := pb.NextRune()
						if !done {
							ss = ss + string(rn)
						}
						com.DbPrintf("testCode", "Case 5: At: ->%s<- ll=%d ss >>>%s<<< %s\n", string(rn), ll, ss, tr.LF())
					}
					if com.DbOn("testDump") {
						pb.Dump01(os.Stdout)
					}
				case 6:
					rn, done := pb.PeekRune()
					if done || rn != ww.Rn {
						t.Errorf("%04s: Peek at [kk=%d] in test [%s] did not work, got %s expected %s, done=%v\n", kk, ii, string(rn), string(ww.Rn), done)
					}
				case 7:
					com.DbPrintf("testCode", "All Done: ss >>>%s<<< before At: %s\n", ss, tr.LF())
					if com.DbOn("testDump") {
						pb.Dump01(os.Stdout)
					}
					for rn, done := pb.NextRune(); !done; rn, done = pb.NextRune() {
						ss = ss + string(rn)
					}
					com.DbPrintf("testCode", "All Done: ss >>>%s<<< after At: %s\n", ss, tr.LF())
				case 8: // Push back a rune array
					pb.PbByteArray([]byte(ww.Data))
					if com.DbOn("testDump") {
						pb.Dump01(os.Stdout)
					}
				case 9:
					pb.PbFile(ww.Fn)
					com.DbPrintf("testCode", "Pb file %s At: %s\n", ww.Fn, tr.LF())
					if com.DbOn("testDump") {
						pb.Dump01(os.Stdout)
					}
				case 10:
					fs := pb.FileSeen(ww.Fn)
					if fs != ww.FileSeenFlag {
						t.Errorf("%04s: Peek at [kk=%d] in test [%s] did not work, got %v expected %s for file seen flagv\n", kk, ii, fs, ww.FileSeenFlag)
					}
				case 11: // Check get file name
					ln, cn, fn := pb.GetPos()
					com.DbPrintf("testCode", "fn=%s ln=%d cn=%d\n", fn, ln, cn) // xyzzy - definitly not working correctly
					if ln != ww.LineNo {
						t.Errorf("%04s: %d: did not match line no Expected ->%d<-, Got ->%d<-\n", vv.Test, kk, ww.LineNo, ln)
					}
					if cn != ww.ColNo {
						t.Errorf("%04s: %d: did not match col no Expected ->%d<-, Got ->%d<-\n", vv.Test, kk, ww.ColNo, cn)
					}
					if fn != ww.Fn {
						t.Errorf("%04s: %d: did not match file name Expected ->%s<-, Got ->%s<-\n", vv.Test, kk, ww.Fn, fn)
					}
				case 12: // Check get file name
					pb.SetPos(ww.LineNo, ww.ColNo, ww.Fn)
				case 15: // Reset
					ResetST()
				case 16: // Reset
					// ProcessInput()
					for rn, done := pb.NextRune(); !done; rn, done = pb.NextRune() {
						if m_body, m_found := HaveMacro(rn); m_found {
							pb.PbString(m_body)
						} else {
							ss = ss + string(rn)
						}
					}
				case 17: // Define
					Define(ww.Rn, ww.Data)
				case 18: // Dump the buffer  - debuging
					pb.Dump01(os.Stdout)
				case 19: // Reset output
					ss = ""
				case 20: // Special test to pub back more than buffer of 'x'
					x := strings.Repeat("x", ww.X)
					pb.PbString(x)
				}
			}
			if ss != vv.Results {
				t.Errorf("%04s: did not match Expected ->%s<-, Got ->%s<-\n", vv.Test, vv.Results, ss)
			} else {
				com.DbPrintf("testCode", "%04s: Passed ------------------------------------------------------------------------------------------------\n\n", vv.Test)
			}
		}
	}

}
