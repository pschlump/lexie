package nfa

/*
1. Change test output to put every output into a ./ref/XXXX.out file (each test has its own output fiel)
//    /opt/homebrew/bin/dot
*/

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"testing"

	"github.com/pschlump/dbgo"
	"github.com/pschlump/filelib"
	"github.com/pschlump/lexie/com"

	. "gopkg.in/check.v1"
)

// https://labix.org/gocheck

// Test the generation of NFA - Non-Deterministic Finite state Automata.
// ------------------------------------------------------------------------------------
//
// The tests are numberd.   This will match with test data and reference data that is
// in the ../ref direcotry.   Because these are graphs a .svg graph will be generated
// using 'dot' in GraphViz.  The output will be in the ../ref directory with a .svg
// extention for each test.  It is much easier to undersnd the graphs in this fasion.
//
// Steps:
// 1. Convert the regular exprssion to a NFA
// 2. Output the NFA to a file, ../ref/nfa_0000.tst
// 3. Read in the corect data from a file, ../ref/nfa_0000.ref
// 4. Compare the strings of the generate and the ference
// 5. Report errors.
//

// TODO - Figure out path to 'dot' so can run it.
// var dotPath = "/opt/homebrew/bin/dot"
var dotPath = "/usr/local/bin/dot"

type Lexie01DataType struct {
	Test         string
	Re           string
	Rv           int
	NExpectedErr int
	SkipTest     bool
	ELen         int
}

var Lexie01Data = []Lexie01DataType{
	{Test: "0000", Re: "(x|y)*abb", Rv: 1000, SkipTest: false, ELen: 3},                                                         // Len(3)
	{Test: "0001", Re: "x*", Rv: 1001, SkipTest: false, ELen: 0},                                                                // Len(0)
	{Test: "0002", Re: "(xx)*", Rv: 1002, SkipTest: false, ELen: 0},                                                             // Len(0)
	{Test: "0003", Re: "(xx)+", Rv: 1003, SkipTest: false, ELen: 2},                                                             // Len(2)
	{Test: "0004", Re: "(xx)?", Rv: 1004, SkipTest: false, ELen: 0},                                                             // Len(0)
	{Test: "0005", Re: "(a|b)", Rv: 1005, SkipTest: false, ELen: 1},                                                             // Len(Min(len(1),Len(1)) = Len(1)
	{Test: "0006", Re: "(aa|bb)", Rv: 1006, SkipTest: false, ELen: 2},                                                           // Len(2)
	{Test: "0007", Re: "(a|b)*abb", Rv: 1007, SkipTest: false, ELen: 3},                                                         // Len(3) Examle from Dragon Compiler Book and .pdf files
	{Test: "0008", Re: "(aa|bb|ccc)*abb", Rv: 1008, SkipTest: false, ELen: 3},                                                   // Len(3)
	{Test: "0009", Re: "^abb$", Rv: 1009, SkipTest: false, ELen: 3},                                                             // Len(3)+Hard
	{Test: "0010", Re: "^abb", Rv: 1010, SkipTest: false, ELen: 3},                                                              // Len(3)+Hard
	{Test: "0011", Re: `a(bcd)*(ghi)+(jkl)*X`, Rv: 1011, SkipTest: false, ELen: 5},                                              // Len(1+3+1)
	{Test: "0012", Re: `a[.]d`, Rv: 1012, SkipTest: false, ELen: 3},                                                             // Len(3)
	{Test: "0013", Re: `a[^]d`, Rv: 1013, SkipTest: false, ELen: 0},                                                             // Len(?) TODO: -- Sigma should have an X_N_CCL char in it - missing
	{Test: "0014", Re: `a(def)*(klm(mno)+)?b`, Rv: 1014, SkipTest: false, ELen: 2},                                              // Len(2)
	{Test: "0015", Re: `a[a-zA-Z_][a-zA-Z_0-9]*d`, Rv: 1015, SkipTest: false, ELen: 3},                                          // Len(3)
	{Test: "0016", Re: `a.d`, Rv: 1016, SkipTest: false, ELen: 3},                                                               // Len(3)
	{Test: "0017", Re: "(aa|bb|ccc)abb", Rv: 1017, SkipTest: false, ELen: 5},                                                    // Len(2+3=5)
	{Test: "0018", Re: "(||)", Rv: 1018, SkipTest: false, ELen: 0},                                                              // Len(0)
	{Test: "0019", Re: "||", Rv: 1019, SkipTest: false, ELen: 0},                                                                // Len(0)
	{Test: "0020", Re: "(||||||||||||||)", Rv: 1020, SkipTest: false, ELen: 0},                                                  // Len(0)
	{Test: "0021", Re: "(||||||||a||||||)", Rv: 1021, SkipTest: false, ELen: 0},                                                 // Len(0)
	{Test: "0022", Re: "(||||||||a|aa|||||)", Rv: 1022, SkipTest: false, ELen: 0},                                               // Len(0)
	{Test: "0023", Re: "(a|aa|aaa)", Rv: 1023, SkipTest: false, ELen: 1},                                                        // Len(1)
	{Test: "0024", Re: "(ab|aab|aaab)", Rv: 1024, SkipTest: false, ELen: 2},                                                     // Len(2)
	{Test: "0025", Re: "(ab|aab|aaab)c", Rv: 1025, SkipTest: false, ELen: 3},                                                    // Len(3)
	{Test: "0026", Re: "(a*|aab|aaab)", Rv: 1026, SkipTest: false, ELen: 0},                                                     // Len(0)
	{Test: "0027", Re: "(-=-|-=#|.*)", Rv: 1027, SkipTest: false, ELen: 0},                                                      // Len(0)		<<<<< CRITICAL TEST >>>>>
	{Test: "0028", Re: "(a\u03bbb|a\u0428b|aaab)", Rv: 1028, SkipTest: false, ELen: 2},                                          // Len(2)
	{Test: "0029", Re: "[0-1]*", Rv: 1003, SkipTest: false, ELen: 2},                                                            // Len(2)
	{Test: "0030", Re: "[0-1]+", Rv: 1003, SkipTest: false, ELen: 2},                                                            // Len(2)
	{Test: "0031", Re: "[0-1]?", Rv: 1003, SkipTest: false, ELen: 2},                                                            // Len(2)
	{Test: "0032", Re: "([0-9]*\\.[0-9]+([eE][0-9]+(\\.[0-9]*)?)?)|([a-zA-Z][a-zA-Z0-9]*)", Rv: 1003, SkipTest: false, ELen: 2}, // Len(2)
	{Test: "0033", Re: "aab{2,3}cc", Rv: 1003, SkipTest: false, ELen: 2},                                                        // Len(2)
	{Test: "0034", Re: "aab{,3}cc", Rv: 1003, SkipTest: false, ELen: 2},                                                         // Len(2)
	{Test: "0035", Re: "aab{2,}cc", Rv: 1003, SkipTest: false, ELen: 2},                                                         // Len(2)
	{Test: "0036", Re: "aab{0,3}cc", Rv: 1003, SkipTest: false, ELen: 2},                                                        // Len(2)
	{Test: "0037", Re: "aab{3,0}cc", Rv: 1003, SkipTest: false, ELen: 2},                                                        // Len(2)
	{Test: "0038", Re: "aab{3,3}cc", Rv: 1003, SkipTest: false, ELen: 2},                                                        // Len(2)
	{Test: "0039", Re: "[0-9]*\\.[0-9]+([eE][-+]?[0-9]+(\\.[0-9]*)?)?", Rv: 1003, SkipTest: false, ELen: 2},                     // Len(2)

	// Test {m,n} stuff -- Must add more complex cases - and hand check!

	// Test [ccl] [[:alpha:]] ccl
	// Test [ccl] [[:lower:]] ccl
	// Test [ccl] [[:upper:]] ccl
	// Test [ccl] [[:numeric:]] ccl
	// Test [^ccl] [^[:alpha:][:numeric:]] ccl
	/*
		wabi-sabi (侘寂?)
		"\\(a"
		"\\)a"
		"\\|a"
		"a\\*b"
	*/
}

// -----------------------------------------------------------------------------------------------------------------------------------------
// From: https://labix.org/gocheck
// -----------------------------------------------------------------------------------------------------------------------------------------
// Hook up gocheck into the "go test" runner.
// -----------------------------------------------------------------------------------------------------------------------------------------

func TestLexie(t *testing.T) { TestingT(t) }

type LexieTestSuite struct{}

var _ = Suite(&LexieTestSuite{})

func (s *LexieTestSuite) TestLexie(c *C) {

	// return
	fmt.Fprintf(os.Stderr, "Test Parsing of REs, Test genration of NFAs %s\n", dbgo.LF())

	dbgo.SetADbFlag("db_NFA", true)
	dbgo.SetADbFlag("db_NFA_LnNo", true)
	dbgo.SetADbFlag("db_DumpPool", true)
	dbgo.SetADbFlag("parseExpression", true)
	dbgo.SetADbFlag("CalcLength", true)

	// Add a test for any issue
	c.Check(42, Equals, 42)
	// c.Assert("nope", Matches, "hel.*there")
	fmt.Printf("**** In Test Issues\n")
	//x := test7GenDFA()
	//c.Check(x, Equals, 0)

	n_err := 0
	n_skip := 0

	for ii, vv := range Lexie01Data {
		if !vv.SkipTest {
			dbgo.Printf("\n\n%(yellow)--- %d Test: %s -----------------------------------------------------------------------------\n\n", ii, vv.Test)

			Pool := NewNFA_Pool()
			Cur := Pool.GetNFA()
			Pool.InitState = Cur

			Pool.AddReInfo(vv.Re, "", 1, vv.Rv, InfoType{})
			Pool.Sigma = Pool.GenerateSigma()

			if false {
				dbgo.DbPrintf("test7", "Pool=%s\n", dbgo.SVarI(Pool))
			}
			Pool.DumpPool(false)
			Pool.DumpPoolJSON(os.Stdout, vv.Re, vv.Rv)

			fmt.Printf("Sigma: ->%s<-\n", Pool.Sigma)

			newFile := fmt.Sprintf("../ref/nfa_%s.tst", vv.Test)
			cmpFile := fmt.Sprintf("../ref/nfa_%s.ref", vv.Test)
			gvFile := fmt.Sprintf("../ref/nfa_%s.gv", vv.Test)
			svgFile := fmt.Sprintf("../ref/nfa_%s.svg", vv.Test)

			fp, _ := filelib.Fopen(newFile, "w")
			Pool.DumpPoolJSON(fp, vv.Re, vv.Rv)
			fp.Close()
			newData, err := ioutil.ReadFile(newFile)
			if err != nil {
				panic("unable to read file, " + cmpFile)
			}

			if filelib.Exists(cmpFile) {
				ref, err := ioutil.ReadFile(cmpFile)
				if err != nil {
					panic("unable to read file, " + cmpFile)
				}
				if string(ref) != string(newData) {
					c.Check(string(newData), Equals, string(ref))
					dbgo.Printf("%(red)Error%(reset): Test case %s failed to match, cmpFile(.ref)=%s newFile(.tst)=%s to \n", vv.Test, cmpFile, newFile)
					n_err++
				}
			} else {
				n_skip++
			}

			gv, _ := filelib.Fopen(gvFile, "w")
			Pool.GenerateGVFile(gv, vv.Re, vv.Rv, "xyzzy401")
			gv.Close()

			// out, err := exec.Command("/usr/local/bin/dot", "-Tsvg", "-o"+svgFile, gvFile).Output()
			out, err := exec.Command(dotPath, "-Tsvg", "-o"+svgFile, gvFile).Output()
			if err != nil {
				dbgo.Printf("%(red)Error from dot, %s\n", err)
				dbgo.Printf("Output: %s\n", out)
				n_err++
			}
		}
	}
	if n_skip > 0 {
		dbgo.Fprintf(os.Stderr, "%(yellow)Skipped, # of files without automated checks = %d\n", n_skip)
		dbgo.DbPrintf("debug", "\n\n%(yellow)Skipped, # of files without automated checks = %d\n", n_skip)
	}
	if n_err > 0 {
		dbgo.Fprintf(os.Stderr, "%(red)Failed, # of errors = %d\n", n_err)
		dbgo.DbPrintf("debug", "\n\n%(red)Failed, # of errors = %d\n", n_err)
	} else {
		dbgo.Fprintf(os.Stderr, "%(green)PASS\n")
		dbgo.DbPrintf("debug", "\n\n%(green)PASS\n")
	}
}

// -----------------------------------------------------------------------------------------------------------------------------------------

type LambdaClosureTestSuite struct{}

var _ = Suite(&LambdaClosureTestSuite{})

func (s *LambdaClosureTestSuite) TestLexie(c *C) {

	// return
	fmt.Fprintf(os.Stderr, "Test NFA generation from REs, %s\n", dbgo.LF())

	n_err := 0

	dbgo.SetADbFlag("db_NFA", true)
	dbgo.SetADbFlag("db_NFA_LnNo", true)
	dbgo.SetADbFlag("db_DumpPool", true)
	dbgo.SetADbFlag("parseExpression", true)

	// {Test: "0011", Re: `a(bcd)*(ghi)+(jkl)*X`, Rv: 1011},     //
	Pool := NewNFA_Pool()
	Cur := Pool.GetNFA()
	Pool.InitState = Cur

	vv := Lexie01Data[11]
	Pool.AddReInfo(vv.Re, "", 1, vv.Rv, InfoType{})
	Pool.Sigma = Pool.GenerateSigma()
	fmt.Printf("\n\nRe: %s\n", vv.Re)
	Pool.DumpPool(false)

	// -------------------------------------- test 1 -----------------------------------
	r1 := Pool.LambdaClosure([]int{4, 1, 5})
	fmt.Printf("\n\nr1(4,1,5)=%v\n", r1)

	if len(com.CompareSlices([]int{4, 1, 5}, r1)) != 0 {
		dbgo.Printf("%(red)Error%(reset): Test case 1 failed to match\n")
		n_err++
	}
	c.Check(len(com.CompareSlices([]int{4, 1, 5}, r1)), Equals, 0)

	// -------------------------------------- test 2 -----------------------------------
	r2 := Pool.LambdaClosure([]int{12, 9, 13})
	fmt.Printf("\n\nr2(5,9,12,9,13)=%v\n", r2)

	if len(com.CompareSlices([]int{12, 9, 13}, r2)) != 0 {
		dbgo.Printf("%(red)Error%(reset): Test case 2 failed to match\n")
		n_err++
	}
	c.Check(len(com.CompareSlices([]int{12, 9, 13}, r2)), Equals, 0)

	// -------------------------------------- test 3 -----------------------------------
	r3 := Pool.LambdaClosure([]int{12, 9, 13})
	fmt.Printf("\n\nr3(12,9,13)=%v\n", r3)

	if len(com.CompareSlices([]int{12, 9, 13}, r3)) != 0 {
		dbgo.Printf("%(red)Error%(reset): Test case 3 failed to match\n")
		n_err++
	}
	c.Check(len(com.CompareSlices([]int{12, 9, 13}, r3)), Equals, 0)

	// ------------------------- understand test runner ---------------------------------------
	if false {
		c.Check(1, Equals, 0)
		// c.Assert(2, Equals, 0) // Failure of an assert ends test (exit)
		sss := c.GetTestLog()
		fp, err := filelib.Fopen(",g", "w")
		c.Check(err, Equals, nil)
		fmt.Fprintf(fp, "c.GetTestLog: ->%s<-\n", sss)
	}

	// ------------------------- eval results now ---------------------------------------

	if n_err > 0 {
		dbgo.Fprintf(os.Stderr, "%(red)Failed, # of errors = %d\n", n_err)
		dbgo.DbPrintf("debug", "\n\n%(red)Failed, # of errors = %d\n", n_err)
	} else {
		dbgo.Fprintf(os.Stderr, "%(green)PASS\n")
		dbgo.DbPrintf("debug", "\n\n%(green)PASS\n")
	}
}

// -----------------------------------------------------------------------------------------------------------------------------------------

type NFATest_03Type struct {
	Re  string
	Sp  string
	Rv  int
	Act int
	Ns  string
}

type NFATest_02DataType struct {
	Test     string
	Data     []NFATest_03Type
	SkipTest bool
	fullRv   string
}

var NFATest_02Data = []NFATest_02DataType{
	{Test: "2000", SkipTest: false, Data: []NFATest_03Type{
		NFATest_03Type{Re: "{{", Rv: 100, Act: com.A_Push, Ns: "2001"},
		NFATest_03Type{Re: "{%", Rv: 101, Act: com.A_Push, Ns: "2002"},
		NFATest_03Type{Re: ".*", Rv: 102},
		NFATest_03Type{Sp: "EOF", Rv: 104},
	},
	},
	{Test: "2001", SkipTest: false, Data: []NFATest_03Type{
		NFATest_03Type{Re: "}}", Rv: 110, Act: com.A_Pop},
		NFATest_03Type{Re: ".*", Rv: 112}, // Implicitly .*/}}
		NFATest_03Type{Sp: "EOF", Rv: 114},
	},
	},
	{Test: "2002", SkipTest: false, Data: []NFATest_03Type{
		NFATest_03Type{Re: "%}", Rv: 120, Act: com.A_Pop},
		NFATest_03Type{Re: ".*", Rv: 122}, // Implicitly .*/%}
		NFATest_03Type{Sp: "EOF", Rv: 124},
	},
	},
}

type NFA_Multi_Part_TestSuite struct{}

var _ = Suite(&NFA_Multi_Part_TestSuite{})

func (s *NFA_Multi_Part_TestSuite) TestLexie(c *C) {

	fmt.Fprintf(os.Stderr, "Test NFA Multi-Part RE - NFA test %s\n", dbgo.LF())
	n_err := 0
	n_skip := 0

	// ------------------------- ------------------------- --------------------------------------- ---------------------------------------
	// Test as sections
	// ------------------------- ------------------------- --------------------------------------- ---------------------------------------

	for ii, vv := range NFATest_02Data {
		s := ""
		for _, ww := range vv.Data {
			s += ww.Re + "\n"
		}
		NFATest_02Data[ii].fullRv = s
	}

	for ii, vv := range NFATest_02Data {
		dbgo.Printf("\n\n%(yellow)--- %2d Test: %4s RE: ->%s<- -------------------------------------------------------------------------------\n", ii, vv.Test, vv.fullRv)
		Pool := NewNFA_Pool()
		Cur := Pool.GetNFA()
		Pool.InitState = Cur
		for jj, ww := range vv.Data {
			fmt.Printf("\n\n--- %2d Test: %4s Part %2d ----------------------------------------------------------------------\n\n", ii, vv.Test, jj)
			// Add in components
			Pool.AddReInfo(ww.Re, "", 1, ww.Rv, InfoType{})
			// 		Dum out parsed REs
			// 		Dum out parsed NFAs along the way
		}
		Pool.FinializeNFA()  // Fnialize
		Pool.DumpPool(false) // Dump out NFA - check it.

		// ------------------------- --------------------------------------- ---------------------------------------
		// Test these also
		// 		func (nn *NFA_PoolType) DeleteRe(oldRe string) {
		// 		func (nn *NFA_PoolType) ChangeRe(oldRe string, newRe string) {
		// ------------------------- --------------------------------------- ---------------------------------------

		newFile := fmt.Sprintf("../ref/n2_%s.tst", vv.Test)
		cmpFile := fmt.Sprintf("../ref/n2_%s.ref", vv.Test)
		gvFile := fmt.Sprintf("../ref/n2_%s.gv", vv.Test)
		svgFile := fmt.Sprintf("../ref/n2_%s.svg", vv.Test)
		fp, _ := filelib.Fopen(newFile, "w")
		Pool.DumpPoolJSON(fp, vv.Test, 0)
		fp.Close()
		newData, err := ioutil.ReadFile(newFile)
		if err != nil {
			panic("unable to read file, " + cmpFile)
		}

		if filelib.Exists(cmpFile) {
			ref, err := ioutil.ReadFile(cmpFile)
			if err != nil {
				panic("unable to read file, " + cmpFile)
			}
			if string(ref) != string(newData) {
				c.Check(string(newData), Equals, string(ref))
				dbgo.Printf("%(red)Error%(reset): Test case %s failed to match\n", vv.Test)
				n_err++
			}
		} else {
			n_skip++
		}

		gv, err := filelib.Fopen(gvFile, "w")
		if err != nil {
			dbgo.Fprintf(os.Stderr, "%(red)Error - unable to open %s for output, error: %s\n", gvFile, err)
			n_err++
			continue
		}
		defer gv.Close()
		Pool.GenerateGVFile(gv, vv.Test, 0, vv.fullRv)

		// _, err = exec.Command("/usr/local/bin/dot", "-Tsvg", "-o"+svgFile, gvFile).Output()
		out, err := exec.Command(dotPath, "-Tsvg", "-o"+svgFile, gvFile).Output()
		if err != nil {
			dbgo.Printf("%(red)Error from dot, %s\n", err)
			dbgo.Printf("Output: %s\n", out)
			n_err++
		}
	}

	// ------------------------- ------------------------- --------------------------------------- ---------------------------------------
	// Test as a single machine
	// ------------------------- ------------------------- --------------------------------------- ---------------------------------------

	if n_skip > 0 {
		dbgo.Fprintf(os.Stderr, "%(yellow)Skipped, # of files without automated checks = %d\n", n_skip)
		dbgo.DbPrintf("debug", "\n\n%(yellow)Skipped, # of files without automated checks = %d\n", n_skip)
	}
	if n_err > 0 {
		dbgo.Fprintf(os.Stderr, "%(red)Failed, # of errors = %d\n", n_err)
		dbgo.DbPrintf("debug", "\n\n%(red)Failed, # of errors = %d\n", n_err)
	} else {
		dbgo.Fprintf(os.Stderr, "%(green)PASS\n")
		dbgo.DbPrintf("debug", "\n\n%(green)PASS\n")
	}

	_ = n_skip
}
