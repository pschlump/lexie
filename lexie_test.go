package lexie

/*
1. Change test output to put every output into a ./ref/XXXX.out file (each test has its own output fiel)
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
	"github.com/pschlump/lexie/dfa"
	"github.com/pschlump/lexie/nfa"
)

type Lexie01DataType struct {
	Test         string
	Re           string
	Rv           int
	NExpectedErr int
	SkipTest     bool
	ELen         int
}

var Lexie01Data = []Lexie01DataType{
	{Test: "0000", Re: "(x|y)*abb", Rv: 1000, SkipTest: false, ELen: 3},                // Len(3)
	{Test: "0001", Re: "x*", Rv: 1001, SkipTest: false, ELen: 0},                       // Len(0)
	{Test: "0002", Re: "(xx)*", Rv: 1002, SkipTest: false, ELen: 0},                    // Len(0)
	{Test: "0003", Re: "(xx)+", Rv: 1003, SkipTest: false, ELen: 2},                    // Len(2)
	{Test: "0004", Re: "(xx)?", Rv: 1004, SkipTest: false, ELen: 0},                    // Len(0)
	{Test: "0005", Re: "(a|b)", Rv: 1005, SkipTest: false, ELen: 1},                    // Len(Min(len(1),Len(1)) = Len(1)
	{Test: "0006", Re: "(aa|bb)", Rv: 1006, SkipTest: false, ELen: 2},                  // Len(2)
	{Test: "0007", Re: "(a|b)*abb", Rv: 1007, SkipTest: false, ELen: 3},                // Len(3) Examle from Dragon Compiler Book and .pdf files
	{Test: "0008", Re: "(aa|bb|ccc)*abb", Rv: 1008, SkipTest: false, ELen: 3},          // Len(3)
	{Test: "0009", Re: "^abb$", Rv: 1009, SkipTest: false, ELen: 3},                    // Len(3)+Hard
	{Test: "0010", Re: "^abb", Rv: 1010, SkipTest: false, ELen: 3},                     // Len(3)+Hard
	{Test: "0011", Re: `a(bcd)*(ghi)+(jkl)*X`, Rv: 1011, SkipTest: false, ELen: 5},     // Len(1+3+1)
	{Test: "0012", Re: `a[.]d`, Rv: 1012, SkipTest: false, ELen: 3},                    // Len(3)
	{Test: "0013", Re: `a[^]d`, Rv: 1013, SkipTest: false, ELen: 0},                    // Len(?) TODO: -- Sigma should have an X_N_CCL char in it - missing
	{Test: "0014", Re: `a(def)*(klm(mno)+)?b`, Rv: 1014, SkipTest: false, ELen: 2},     // Len(2)
	{Test: "0015", Re: `a[a-zA-Z_][a-zA-Z_0-9]*d`, Rv: 1015, SkipTest: false, ELen: 3}, // Len(3)
	{Test: "0016", Re: `a.d`, Rv: 1016, SkipTest: false, ELen: 3},                      // Len(3)
	{Test: "0017", Re: "(aa|bb|ccc)abb", Rv: 1017, SkipTest: false, ELen: 5},           // Len(2+3=5)
	{Test: "0018", Re: "(||)", Rv: 1018, SkipTest: false, ELen: 0},                     // Len(0)
	{Test: "0019", Re: "||", Rv: 1019, SkipTest: false, ELen: 0},                       // Len(0)
	{Test: "0020", Re: "(||||||||||||||)", Rv: 1020, SkipTest: false, ELen: 0},         // Len(0)
	{Test: "0021", Re: "(||||||||a||||||)", Rv: 1021, SkipTest: false, ELen: 0},        // Len(0)
	{Test: "0022", Re: "(||||||||a|aa|||||)", Rv: 1022, SkipTest: false, ELen: 0},      // Len(0)
	{Test: "0023", Re: "(a|aa|aaa)", Rv: 1023, SkipTest: false, ELen: 1},               // Len(1)
	{Test: "0024", Re: "(ab|aab|aaab)", Rv: 1024, SkipTest: false, ELen: 2},            // Len(2)
	{Test: "0025", Re: "(ab|aab|aaab)c", Rv: 1025, SkipTest: false, ELen: 3},           // Len(3)
	{Test: "0026", Re: "(a*|aab|aaab)", Rv: 1026, SkipTest: false, ELen: 0},            // Len(0)
	{Test: "0027", Re: "(-=-|-=#|.*)", Rv: 1027, SkipTest: false, ELen: 0},             // Len(0)		<<<<< CRITICAL TEST >>>>>
	{Test: "0028", Re: "(a\u03bbb|a\u0428b|aaab)", Rv: 1028, SkipTest: false, ELen: 2}, // Len(2)
	// Test {m,n} stuff
	// Test [ccl] [[:alpha:]] ccl
}

// -----------------------------------------------------------------------------------------------------------------------------------------
// From: https://labix.org/gocheck
// -----------------------------------------------------------------------------------------------------------------------------------------
// Hook up gocheck into the "go test" runner.
// -----------------------------------------------------------------------------------------------------------------------------------------

// func (s *LexieTestSuite) TestLexie(c *C) {
func Test_DfaTestUsingDjango(t *testing.T) {

	fmt.Fprintf(os.Stderr, "Test Parsing of REs, %s\n", dbgo.LF())

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
		fmt.Printf("\n\n--- %d Test: %s -----------------------------------------------------------------------------\n\n", ii, vv.Test)

		Pool := nfa.NewNFA_Pool()
		Cur := Pool.GetNFA()
		Pool.InitState = Cur

		Pool.AddReInfo(vv.Re, "", 1, vv.Rv, nfa.InfoType{})
		Pool.Sigma = Pool.GenerateSigma()

		dbgo.DbPrintf("test7", "Pool=%s\n", dbgo.SVarI(Pool))

		Pool.DumpPool(false)
		Pool.DumpPoolJSON(os.Stdout, vv.Re, vv.Rv)

		fmt.Printf("Sigma: ->%s<-\n", Pool.Sigma)

		newFile := fmt.Sprintf("./ref/nfa_%s.tst", vv.Test)
		cmpFile := fmt.Sprintf("./ref/nfa_%s.ref", vv.Test)
		gvFile := fmt.Sprintf("./ref/nfa_%s.gv", vv.Test)
		svgFile := fmt.Sprintf("./ref/nfa_%s.svg", vv.Test)
		fp, err := filelib.Fopen(newFile, "w")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to open %s for output: %s\n", newFile, err)
			// t.Fatal("Invalid file name")
		}
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
				dbgo.Printf("%(red)Error%(reset): Test case %s failed to match\n", vv.Test)
				n_err++
			}
		} else {
			n_skip++
		}

		gv, err := filelib.Fopen(gvFile, "w")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to open %s for output: %s\n", gvFile, err)
			// t.Fatal("Invalid file name")
		}
		Pool.GenerateGVFile(gv, vv.Re, vv.Rv, vv.Re)
		gv.Close()

		out, err := exec.Command("/usr/local/bin/dot", "-Tsvg", "-o"+svgFile, gvFile).Output()
		if err != nil {
			fmt.Printf("Error from dot, %s, %s\n", err, dbgo.LF())
			fmt.Printf("Output: %s\n", out)
		}
	}
	if n_skip > 0 {
		dbgo.Fprintf(os.Stderr, "%(yellow)Skipped, # of files without automated checks = %d\n", n_skip)
	}
	if n_err > 0 {
		dbgo.Fprintf(os.Stderr, "%(red)Failed, # of errors = %d\n", n_err)
	} else {
		dbgo.Fprintf(os.Stderr, "%(green)PASS\n")
	}
}

// -----------------------------------------------------------------------------------------------------------------------------------------

type LambdaClosureTestSuite struct{}

var _ = Suite(&LambdaClosureTestSuite{})

func (s *LambdaClosureTestSuite) TestLexie(c *C) {

	return
	fmt.Fprintf(os.Stderr, "Test NFA generation from REs, %s\n", dbgo.LF())

	n_err := 0

	dbgo.SetADbFlag("db_NFA", true)
	dbgo.SetADbFlag("db_NFA_LnNo", true)
	dbgo.SetADbFlag("db_DumpPool", true)
	dbgo.SetADbFlag("parseExpression", true)

	// {Test: "0011", Re: `a(bcd)*(ghi)+(jkl)*X`, Rv: 1011},     //
	Pool := nfa.NewNFA_Pool()
	Cur := Pool.GetNFA()
	Pool.InitState = Cur

	vv := Lexie01Data[11]
	Pool.AddReInfo(vv.Re, "", 1, vv.Rv, nfa.InfoType{})
	Pool.Sigma = Pool.GenerateSigma()
	fmt.Printf("\n\nRe: %s\n", vv.Re)
	Pool.DumpPool(false)

	// -------------------------------------- test 1 -----------------------------------
	r1 := Pool.LambdaClosure([]int{1})
	fmt.Printf("\n\nr1(1)=%v\n", r1)

	if len(com.CompareSlices([]int{4, 1, 5, 6}, r1)) != 0 {
		dbgo.Printf("%(red)Error%(reset): Test case 1 failed to match\n")
		n_err++
	}
	// c.Check(r1, Equals, []int{4, 1, 5, 6})
	c.Check(len(com.CompareSlices([]int{4, 1, 5, 6}, r1)), Equals, 0)

	// -------------------------------------- test 2 -----------------------------------
	r2 := Pool.LambdaClosure([]int{5, 8, 11})
	fmt.Printf("\n\nr2(5,8,11)=%v\n", r2)

	if len(com.CompareSlices([]int{6}, r2)) != 0 {
		dbgo.Printf("%(red)Error%(reset): Test case 2 failed to match\n")
		n_err++
	}
	// c.Check(r2, Equals, []int{5, 8, 11})
	c.Check(len(com.CompareSlices([]int{6}, r2)), Equals, 0)

	// -------------------------------------- test 3 -----------------------------------
	r3 := Pool.LambdaClosure([]int{9})
	fmt.Printf("\n\nr3(9)=%v\n", r3)

	if len(com.CompareSlices([]int{6, 10, 13, 10, 14}, r3)) != 0 {
		dbgo.Printf("%(red)Error%(reset): Test case 3 failed to match\n")
		n_err++
	}
	// c.Check(r3, Equals, []int{5, 8, 11})
	c.Check(len(com.CompareSlices([]int{6, 10, 13, 10, 14}, r3)), Equals, 0)

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
	} else {
		dbgo.Fprintf(os.Stderr, "%(green)PASS\n")
	}
}

// -----------------------------------------------------------------------------------------------------------------------------------------

type NFA_to_DFA_TestSuite struct{}

var _ = Suite(&NFA_to_DFA_TestSuite{})

func (s *NFA_to_DFA_TestSuite) TestLexie(c *C) {

	fmt.Fprintf(os.Stderr, "Test NFA to DFA, %s\n", dbgo.LF())

	n_err := 0
	n_skip := 0

	dbgo.SetADbFlag("db_DFAGen", true)
	dbgo.SetADbFlag("db_DumpDFAPool", true)
	dbgo.SetADbFlag("db_DFA_LnNo", true)

	// -----------------------------------------------------------------------------------------------------------------
	// -----------------------------------------------------------------------------------------------------------------
	// -----------------------------------------------------------------------------------------------------------------
	// vv := Lexie01Data[11]
	for ii, vv := range Lexie01Data {
		if !vv.SkipTest {
			// {Test: "0011", Re: `a(bcd)*(ghi)+(jkl)*X`, Rv: 1011},     //
			Nfa := nfa.NewNFA_Pool()
			Cur := Nfa.GetNFA()
			Nfa.InitState = Cur

			// c.Log("\n\n--- %d NFA to DFA Test: %s -----------------------------------------------------------------------------\n\n", ii, vv.Test)
			fmt.Printf("\n\n--- %d NFA to DFA Test: %s -----------------------------------------------------------------------------\n\n", ii, vv.Test)
			newFile := fmt.Sprintf("./ref/dfa_%s.tst", vv.Test)
			cmpFile := fmt.Sprintf("./ref/dfa_%s.ref", vv.Test)
			gvFile := fmt.Sprintf("./ref/dfa_%s.gv", vv.Test)
			svgFile := fmt.Sprintf("./ref/dfa_%s.svg", vv.Test)
			tabFile := fmt.Sprintf("./ref/dfa_%s.tab", vv.Test)
			Nfa.AddReInfo(vv.Re, "", 1, vv.Rv, nfa.InfoType{})
			// Nfa.Sigma = Nfa.GenerateSigma()
			Nfa.Sigma = Nfa.GenerateSigma()
			fmt.Printf("\nSigma: ->%s<-\n", Nfa.Sigma)
			fmt.Printf("\n\nRe: %s\n", vv.Re)
			Nfa.DumpPool(false)

			Dfa := dfa.NewDFA_Pool()
			Dfa.ConvNDA_to_DFA(Nfa)

			fp, _ := filelib.Fopen(newFile, "w")
			Dfa.DumpPoolJSON(fp, vv.Re, vv.Rv)
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
			// -----------------------------------------------------------------------------------------------------------------
			// Post dump
			// -----------------------------------------------------------------------------------------------------------------
			fmt.Printf("----------------------------------------------------------------------------------------------------\n")
			fmt.Printf("Bottom .. machine is...\n")
			fmt.Printf("----------------------------------------------------------------------------------------------------\n")
			Dfa.DumpPool(false)

			// func (dfa *NFA_PoolType) GenerateGVFile(fo io.Writer, td string, tn int) {
			gv, _ := filelib.Fopen(gvFile, "w")
			Dfa.GenerateGVFile(gv, vv.Re, vv.Rv)
			gv.Close()

			tf, _ := filelib.Fopen(tabFile, "w")
			//Dfa.DumpPoolJSON(fp, vv.Re, vv.Rv)
			Dfa.OutputInFormat(tf, "text") // text, go-code, c-code, json, xml etc.
			tf.Close()

			out, err := exec.Command("/usr/local/bin/dot", "-Tsvg", "-o"+svgFile, gvFile).Output()
			if err != nil {
				fmt.Printf("Error from dot, %s, %s\n", err, dbgo.LF())
				fmt.Printf("Output: %s\n", out)
			}
		}
	}
	// -----------------------------------------------------------------------------------------------------------------
	// -----------------------------------------------------------------------------------------------------------------
	// -----------------------------------------------------------------------------------------------------------------

	if n_err > 0 {
		dbgo.Fprintf(os.Stderr, "%(red)Failed, # of errors = %d\n", n_err)
	} else {
		dbgo.Fprintf(os.Stderr, "%(green)PASS\n")
	}
}

// ------------------------------------------------------------------------------------------------------------------------------------------------------

/*
***** NOT WORKING *****

type Lexie02DataType struct {
	Test     string
	Inp      string
	Rv       int
	SkipTest bool
}

// 1. abc {{ == != }} {% id %} def
// 1. abc {{ == {% xy %} != }} {% id %} def

var Lexie02Data = []Lexie02DataType{
	{Test: "1000", Inp: "<BEF>{{ pp ! : {% qq ! : rr %} ss }}</aft>", Rv: 1000, SkipTest: false},
	{Test: "1001", Inp: "<bef>{{ pp != ! : {% qq != ! : rr %} ss }}</aft>", Rv: 1001, SkipTest: false},
}

type Reader_TestSuite struct{}

var _ = Suite(&Reader_TestSuite{})

func (s *Reader_TestSuite) TestLexie(c *C) {

	// return
	fmt.Fprintf(os.Stderr, "Test Matcher test from .json file, %s\n", dbgo.LF())

	dbOn["db_DumpDFAPool"] = true // DFA Dump Pool
	dbOn["db_DumpPool"] = true    // NFA Dump Pool
	dbOn["db_Matcher_02"] = true  // NFA Dump Pool

	lex := dfa.NewLexie()
	lex.ReadJSONSpec("./.json")

	if true {

		r := strings.NewReader("abcd{% simple {{ \u2022 }} stuff %} mOre")

		init := lex.DFA_Start["S_Init"]
		fmt.Printf("\n <> <> <> Machine Number: %d\n", init)
		dfa := lex.DFA_Machine[init]
		dfa.OutputInFormat(os.Stdout, "text")

		dfa.MatcherNewTab(r)

	} else {

		for ii, vv := range Lexie02Data {

			fmt.Printf("\n\nTest:%s ------------------------- Start --------------------------, %d, Input: -->>%s<<--\n", vv.Test, ii, vv.Inp)

			r := strings.NewReader(vv.Inp)
			lex.MatchInput2(r)

			fmt.Printf("Final Dump of TokenBuffer, Test:%s\n", vv.Test)
			lex.ATokList.DumpTokenBuffer()

			fmt.Printf("Test:%s ------------------------- End --------------------------\n\n", vv.Test)

		}

	}
}
*/

/* vim: set noai ts=4 sw=4: */
