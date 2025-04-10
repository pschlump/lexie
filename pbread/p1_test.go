package pbread

//
// P B B U F F E R - Push back buffer - test.
//
// (C) Philip Schlump, 2013-2015.
// Version: 1.0.0
// BuildNo: 28
//
// Special Thanks to 2C-Why, LLC for supporting this project.
//

import (
	"os"
	"strings"
	"testing"

	"github.com/pschlump/dbgo"
	"github.com/pschlump/filelib"
)

const (
	CmdOpenFile        = 1
	CmdPbString        = 2
	CmdPbRune          = 3
	CmdPbRuneArray     = 4
	CmdNextNChar       = 5
	CmdPeek            = 6
	CmdOutputToEof     = 7
	CmdPbByteArray     = 8
	CmdPbFile          = 9
	CmdFileSeen        = 10
	CmdGetPos          = 11
	CmdSetPos          = 12
	CmdResetST         = 15
	CmdMacroProc       = 16
	CmdDefineMacro     = 17
	CmdDumpBuffer      = 18
	CmdResetOutput     = 19
	CmdPushBackXCopies = 20 // Special test to pub back more than buffer of 'x'
)

type ActType struct {
	OpCode       int    // 1/CmdOpenFile = Fn:filename, 2/CmdPbString = Data:string, 3/CmdPbRune = Rn:rune, 4/CmdPbRuneArray = RnS:runeArray, 5/CmdNextNChar = get X from Next, 6/CmdPeek = X:n-to-peek,
	Fn           string // 7/CmdOutputToEof = do output, 8/CmdPbByteArray = Data:byte array with type cast - input as string, 9/CmdPbFile = Fn:pbfile, 10/CmdFileSeen FileSeenFlag:t/f to match,
	Data         string // 11/CmdGetPos get filenae+lineno+colno Fn:Filename LineNo:number ColNo:number, 12/CmdSetPos set pos (see CmdGetPos), 15/CmdResetSt reset symbol table to empty,
	Rn           rune   // 16/CmdMacroProc process macros to eof, 18/CmdDumpBuffer = do output, if Fn is set send to file else stdout,  17/CmdDefiineMacro pushback macro test ( CH=.Rn to replace with Str=Data )
	RnS          []rune // 19/CmdResetOutput discard output to make compare easier, 20/CmdPushBackXCopies of Data
	X            int    //
	FileSeenFlag bool   //
	LineNo       int    //
	ColNo        int    //
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
		ActType{OpCode: CmdOpenFile, Fn: "test/t1.txt"},
		ActType{OpCode: CmdNextNChar, X: 4},
		ActType{OpCode: CmdOpenFile, Fn: "test/t2.txt"},
		ActType{OpCode: CmdOutputToEof},
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
		ActType{OpCode: CmdOpenFile, Fn: "test/t1.txt"},
		ActType{OpCode: CmdNextNChar, X: 4},
		ActType{OpCode: CmdPbRune, Rn: 'y'},
		ActType{OpCode: CmdPbRune, Rn: 'x'},
		ActType{OpCode: CmdOutputToEof},
	}, Results: `a12
xyb34
ccc
ddd
`},

	// push back rune array
	{Test: "0003", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdOpenFile, Fn: "test/t1.txt"},
		ActType{OpCode: CmdNextNChar, X: 4},
		ActType{OpCode: CmdPbRuneArray, RnS: []rune{'x', 'y', 'z', 'w', 'v'}},
		ActType{OpCode: CmdOutputToEof},
	}, Results: `a12
xyzwvb34
ccc
ddd
`},

	// push back string
	{Test: "0004", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdOpenFile, Fn: "test/t1.txt"},
		ActType{OpCode: CmdNextNChar, X: 4},
		ActType{OpCode: CmdPbString, Data: "xyzwv"},
		ActType{OpCode: CmdOutputToEof},
	}, Results: `a12
xyzwvb34
ccc
ddd
`},

	// push back a number of items
	{Test: "0005", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdOpenFile, Fn: "test/t1.txt"},
		ActType{OpCode: CmdPeek, Rn: 'a'},
		ActType{OpCode: CmdNextNChar, X: 4},
		ActType{OpCode: CmdPbString, Data: "xyzwv"},
		ActType{OpCode: CmdPeek, Rn: 'x'},
		ActType{OpCode: CmdNextNChar, X: 4},
		ActType{OpCode: CmdPeek, Rn: 'v'},
		ActType{OpCode: CmdOutputToEof},
	}, Results: `a12
xyzwvb34
ccc
ddd
`},

	// Pull some off before push back
	{Test: "0006", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdOpenFile, Fn: "test/t1.txt"},
		ActType{OpCode: CmdNextNChar, X: 4},
		ActType{OpCode: 8, Data: "XyZwV"},
		ActType{OpCode: CmdOutputToEof},
	}, Results: `a12
XyZwVb34
ccc
ddd
`},

	// output empty buffer
	{Test: "0007", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdOutputToEof},
	}, Results: ``},

	// test file names
	{Test: "0008", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdGetPos, Fn: "", LineNo: 1, ColNo: 1},
		ActType{OpCode: CmdOpenFile, Fn: "test/t1.txt"},
		ActType{OpCode: CmdNextNChar, X: 4},
		ActType{OpCode: CmdGetPos, Fn: "test/t1.txt", LineNo: 2, ColNo: 1},
		ActType{OpCode: CmdNextNChar, X: 1},
		ActType{OpCode: CmdGetPos, Fn: "test/t1.txt", LineNo: 2, ColNo: 2},
		ActType{OpCode: CmdNextNChar, X: 3},
		ActType{OpCode: CmdGetPos, Fn: "test/t1.txt", LineNo: 3, ColNo: 1},
		ActType{OpCode: CmdPbFile, Fn: "test/t2.txt"},
		ActType{OpCode: CmdGetPos, Fn: "test/t2.txt", LineNo: 1, ColNo: 1},
		ActType{OpCode: CmdFileSeen, FileSeenFlag: true, Fn: "test/t1.txt"},
		ActType{OpCode: CmdFileSeen, FileSeenFlag: true, Fn: "test/t2.txt"},
		ActType{OpCode: CmdFileSeen, FileSeenFlag: false, Fn: "test/t3.txt"},
		ActType{OpCode: CmdOutputToEof},
		ActType{OpCode: CmdGetPos, Fn: "test/t1.txt", LineNo: 5, ColNo: 1},
		ActType{OpCode: CmdSetPos, Fn: "test/t8.txt", LineNo: 55, ColNo: 202},
		ActType{OpCode: CmdGetPos, Fn: "test/t8.txt", LineNo: 55, ColNo: 202},
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
		ActType{OpCode: CmdOpenFile, Fn: "test/x1.txt"},
		ActType{OpCode: CmdPbFile, Fn: "test/x2.txt"},
	}, Results: ``},

	// test push back that overflows PbAFew's max
	{Test: "0010", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdPushBackXCopies, X: 622, Data: "x"}, // Push back 622 chars ( more than buffer size )
		ActType{OpCode: CmdDumpBuffer, Fn: "test/10-1.out"},    // Output
		ActType{OpCode: CmdOpenFile, Fn: "test/t1.txt"},        // Tak test/t1.txt on end
		ActType{OpCode: CmdNextNChar, X: 622},                  // INput 622 chars
		ActType{OpCode: CmdPushBackXCopies, X: 622, Data: "x"}, // Push back 622 chars ( more than buffer size )
		ActType{OpCode: CmdNextNChar, X: 622},                  // INput 622 chars
		ActType{OpCode: CmdResetOutput},                        // Discard current input so test is simpler
		ActType{OpCode: CmdOutputToEof},                        // Get input - should be file
	}, Results: "a12\nb34\nccc\nddd\n"}, //

	{Test: "0011", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdOpenFile, Fn: "test/t1.txt"},                       //	Input file
		ActType{OpCode: CmdNextNChar, X: 4},                                   //	Pick off 4 chars at beginning
		ActType{OpCode: CmdPbRuneArray, RnS: []rune{'x', 'y', 'z', 'w', 'v'}}, //  Push back "xyzwv" - 5 runes
		ActType{OpCode: CmdSetPos, Fn: "test/t8.txt", LineNo: 55, ColNo: 202}, // Set file name to t8.txt,line55,col202
		ActType{OpCode: CmdGetPos, Fn: "test/t8.txt", LineNo: 55, ColNo: 202}, // Check that we have those values
		ActType{OpCode: CmdNextNChar, X: 3},                                   //  Forward by 3
		ActType{OpCode: CmdGetPos, Fn: "test/t8.txt", LineNo: 55, ColNo: 205}, //  Check col changed
		ActType{OpCode: CmdNextNChar, X: 4},                                   // Forward by 4 more runes
		ActType{OpCode: CmdDumpBuffer, Fn: "test/11-9.out"},                   // Output before the final test -- Just one buffer left if OK
		ActType{OpCode: CmdGetPos, Fn: "test/t1.txt", LineNo: 2, ColNo: 3},    //
		ActType{OpCode: CmdOutputToEof},                                       //
	}, Results: `a12
xyzwvb34
ccc
ddd
`},

	// Simple macro replace
	{Test: "0012", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdOpenFile, Fn: "test/t4.txt"}, // Input file
		ActType{OpCode: CmdResetST},
		ActType{OpCode: CmdDefineMacro, Rn: 'M', Data: "aaaDvvv"},
		ActType{OpCode: CmdDefineMacro, Rn: 'D', Data: "ddd"},
		ActType{OpCode: CmdDefineMacro, Rn: 'E', Data: "e"},
		ActType{OpCode: CmdDefineMacro, Rn: 'J', Data: "**error**"},
		ActType{OpCode: CmdMacroProc},
	}, Results: `abaaadddvvvcdddmem
abdddcdddmÏm
`},

	// test combo files/ pb/ peek etc.  // (fixed??) test it - xyzzyErr1 - what happens if we PbFile when NPbAFew is > 0 - don't we need to pre-buffer the push back?
	// test with a list of files in order - pb.OpenFile()
	// test with a set of files in PbFile()
	{Test: "0013", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdOpenFile, Fn: "test/t1.txt"},
		ActType{OpCode: CmdNextNChar, X: 4},
		ActType{OpCode: CmdPbRune, Rn: 'Y'},
		ActType{OpCode: CmdPbRune, Rn: 'X'},
		ActType{OpCode: CmdPbFile, Fn: "test/t2.txt"},
		ActType{OpCode: CmdPbRune, Rn: 'Z'},
		ActType{OpCode: CmdPbRune, Rn: '8'},
		ActType{OpCode: CmdResetST},
		ActType{OpCode: CmdDefineMacro, Rn: 'Z', Data: "aaaDvvv"},
		ActType{OpCode: CmdDefineMacro, Rn: 'u', Data: "((it))"},
		ActType{OpCode: CmdDefineMacro, Rn: 'D', Data: "777u777"},
		ActType{OpCode: CmdMacroProc},
	}, Results: `a12
8aaa777((it))777vvvEEE
FFF
GGG
XYb34
ccc
ddd
`},
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

			// Implement a quick - fetch execute macine to test - the PbBuffer - commands/opcodes are the Cmd* constants above.
			ss := ""
			pb := NewPbRead()
			ResetST()
			for pc, ww := range vv.Actions {

				switch ww.OpCode {
				case CmdOpenFile: // Open a file , at the tail end of list of input
					pb.OpenFile(ww.Fn)
					dbgo.DbPrintf("testCode", "Open file %s At: %s\n", ww.Fn, dbgo.LF())
					if dbgo.DbOn("testDump") {
						pb.Dump01(os.Stdout)
					}
				case CmdPbString: // Push back a string
					pb.PbString(ww.Data)
					if dbgo.DbOn("testDump") {
						pb.Dump01(os.Stdout)
					}
				case CmdPbRune: // Push back a rune
					pb.PbRune(ww.Rn)
					if dbgo.DbOn("testDump") {
						pb.Dump01(os.Stdout)
					}
				case CmdPbRuneArray: // Push back a rune array
					pb.PbRuneArray(ww.RnS)
					if dbgo.DbOn("testDump") {
						pb.Dump01(os.Stdout)
					}
				case CmdNextNChar:
					for ll := 0; ll < ww.X; ll++ {
						rn, done := pb.NextRune()
						if !done {
							ss = ss + string(rn)
						}
						dbgo.DbPrintf("testCode", "Case 5: At: ->%s<- ll=%d ss >>>%s<<< %s\n", string(rn), ll, ss, dbgo.LF())
					}
					if dbgo.DbOn("testDump") {
						pb.Dump01(os.Stdout)
					}
				case CmdPeek:
					rn, done := pb.PeekRune()
					if done || rn != ww.Rn {
						t.Errorf("%04s: Peek at [pc=%d] in test [%s] did not work, got %s expected %s, done=%v\n", pc, ii, string(rn), string(ww.Rn), done)
					}
				case CmdOutputToEof:
					dbgo.DbPrintf("testCode", "All Done: ss >>>%s<<< before At: %s\n", ss, dbgo.LF())
					if dbgo.DbOn("testDump") {
						pb.Dump01(os.Stdout)
					}
					for rn, done := pb.NextRune(); !done; rn, done = pb.NextRune() {
						ss = ss + string(rn)
					}
					dbgo.DbPrintf("testCode", "All Done: ss >>>%s<<< after At: %s\n", ss, dbgo.LF())
				case CmdPbByteArray: // Push back a byte array
					pb.PbByteArray([]byte(ww.Data))
					if dbgo.DbOn("testDump") {
						pb.Dump01(os.Stdout)
					}
				case CmdPbFile: // Open file and push contents back onto input at head of list. (Macro file, Include, Require)
					pb.PbFile(ww.Fn)
					dbgo.DbPrintf("testCode", "Pb file %s At: %s\n", ww.Fn, dbgo.LF())
					if dbgo.DbOn("testDump") {
						pb.Dump01(os.Stdout)
					}
				case CmdFileSeen:
					fs := pb.FileSeen(ww.Fn)
					if fs != ww.FileSeenFlag {
						t.Errorf("%04s: Peek at [pc=%d] in test [%s] did not work, got %v expected %s for file seen flagv\n", pc, ii, fs, ww.FileSeenFlag)
					}
				case CmdGetPos: // Check get file name
					ln, cn, fn := pb.GetPos()
					dbgo.DbPrintf("testCode", "fn=%s ln=%d cn=%d\n", fn, ln, cn)
					if ln != ww.LineNo {
						t.Errorf("%04s: %d: did not match line no Expected ->%d<-, Got ->%d<-\n", vv.Test, pc, ww.LineNo, ln)
					}
					if cn != ww.ColNo {
						t.Errorf("%04s: %d: did not match col no Expected ->%d<-, Got ->%d<-\n", vv.Test, pc, ww.ColNo, cn)
					}
					if fn != ww.Fn {
						t.Errorf("%04s: %d: did not match file name Expected ->%s<-, Got ->%s<-\n", vv.Test, pc, ww.Fn, fn)
					}
				case CmdSetPos: // Check get file name
					pb.SetPos(ww.LineNo, ww.ColNo, ww.Fn)
				case CmdResetST: // Reset symbol table
					ResetST()
				case CmdMacroProc: // Apply 1 char macros to input and process
					for rn, done := pb.NextRune(); !done; rn, done = pb.NextRune() {
						if m_body, m_found := HaveMacro(rn); m_found {
							pb.PbString(m_body)
						} else {
							ss = ss + string(rn)
						}
					}
				case CmdDefineMacro: // Define
					Define(ww.Rn, ww.Data)
				case CmdDumpBuffer: // Dump the buffer  - debuging
					if ww.Fn == "" {
						pb.Dump01(os.Stdout)
					} else {
						fp, err := filelib.Fopen(ww.Fn, "w")
						if err == nil {
							pb.Dump01(fp)
							fp.Close()
						} else {
							pb.Dump01(os.Stdout)
							t.Errorf("%04s: Unable to open file for output ->%s<-, error: %s\n", vv.Test, ww.Fn, err)
						}
					}
				case CmdResetOutput: // Reset output
					ss = ""
				case CmdPushBackXCopies: // Special test to pub back more than buffer of 'x'
					x := strings.Repeat(ww.Data, ww.X)
					pb.PbString(x)
				}
			}
			if ss != vv.Results {
				t.Errorf("%04s: did not match Expected ->%s<-, Got ->%s<-\n", vv.Test, vv.Results, ss)
			} else {
				dbgo.DbPrintf("testCode", "%04s: Passed ------------------------------------------------------------------------------------------------\n\n", vv.Test)
			}
		}
	}

}
