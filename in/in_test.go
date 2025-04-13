package in

import (
	"fmt"
	"testing"

	"github.com/pschlump/dbgo"
)

type In01TestType struct {
	Test       string
	Inp        string
	SkipTest   bool
	ResultTst2 string
	ResultTst3 []string
	ResultTst4 string
	ResultTst5 []string
	ResultTst6 []string
	ResultTst7 []string
}

var In01Test = []In01TestType{
	{Test: "0001", Inp: "abc : Rv(xxx1) ", SkipTest: false, ResultTst2: "[Rv(xxx1)]", ResultTst3: []string{"Rv", "xxx1"}},
	{Test: "0002", Inp: "abc : Rv(xxx1) Call(Machine)", SkipTest: false, ResultTst2: "[Rv(xxx1) Call(Machine)]", ResultTst3: []string{"Rv", "xxx1", "Call", "Machine"}},
	{Test: "0003", Inp: "abc					: Rv(Tok_OP_Var) Call(S_VAR)", SkipTest: false, ResultTst2: "[Rv(Tok_OP_Var) Call(S_VAR)]"},
	{Test: "0004", Inp: "abc					: Repl(`{{`)					", SkipTest: false, ResultTst2: "[Repl(`{{`)]", ResultTst3: []string{"Repl", `\{\{`}},
	{Test: "0005", Inp: "abc					: Rv(Tok_OP_Tag) Call(S_TAG)", SkipTest: false, ResultTst2: "[Rv(Tok_OP_Tag) Call(S_TAG)]"},
	{Test: "0006", Inp: "abc					: Repl(`{%`)					", SkipTest: false, ResultTst2: "[Repl(`{%`)]"},
	{Test: "0007", Inp: "abc						: Rv(Tok_HTML)", SkipTest: false, ResultTst2: "[Rv(Tok_HTML)]"},
	{Test: "0008", Inp: "abc	: Rv(Tok_ID) ReservedWord()		", SkipTest: false, ResultTst2: "[Rv(Tok_ID) ReservedWord()]"},
	{Test: "0009", Inp: "abc					: Rv(Tok_NUM)", SkipTest: false, ResultTst2: "[Rv(Tok_NUM)]"},
	{Test: "0010", Inp: "abc						: Rv(Tok_LE)", SkipTest: false, ResultTst2: "[Rv(Tok_LE)]"},
	{Test: "0011", Inp: "abc						: Rv(Tok_EQEQ)", SkipTest: false, ResultTst2: "[Rv(Tok_EQEQ)]"},
	{Test: "0012", Inp: "abc						: Rv(Tok_GE)", SkipTest: false, ResultTst2: "[Rv(Tok_GE)]"},
	{Test: "0013", Inp: "abc						: Rv(Tok_L_AND)", SkipTest: false, ResultTst2: "[Rv(Tok_L_AND)]"},
	{Test: "0014", Inp: "abc					: Rv(Tok_L_OR)", SkipTest: false, ResultTst2: "[Rv(Tok_L_OR)]"},
	{Test: "0015", Inp: "abc						: Rv(Tok_NE)", SkipTest: false, ResultTst2: "[Rv(Tok_NE)]"},
	{Test: "0016", Inp: "abc						: Rv(Tok_NE)", SkipTest: false, ResultTst2: "[Rv(Tok_NE)]"},
	{Test: "0017", Inp: "abc						: Call(S_Str0)					", SkipTest: false, ResultTst2: "[Call(S_Str0)]"},
	{Test: "0018", Inp: "abc						: Call(S_Str1)					", SkipTest: false, ResultTst2: "[Call(S_Str1)]"},
	{Test: "0019", Inp: "abc						: Call(S_Str2)					", SkipTest: false, ResultTst2: "[Call(S_Str2)]"},
	{Test: "0020", Inp: "abc						: Rv(Tok_CARROT)", SkipTest: false, ResultTst2: "[Rv(Tok_CARROT)]"},
	{Test: "0021", Inp: "abc						: Rv(Tok_OP_PAR)", SkipTest: false, ResultTst2: "[Rv(Tok_OP_PAR)]"},
	{Test: "0022", Inp: "abc						: Rv(Tok_CL_PAR)", SkipTest: false, ResultTst2: "[Rv(Tok_CL_PAR)]"},
	{Test: "0023", Inp: "abc						: Rv(Tok_PLUS)", SkipTest: false, ResultTst2: "[Rv(Tok_PLUS)]"},
	{Test: "0024", Inp: "abc						: Rv(Tok_MINUS)", SkipTest: false, ResultTst2: "[Rv(Tok_MINUS)]"},
	{Test: "0025", Inp: "abc						: Rv(Tok_STAR)", SkipTest: false, ResultTst2: "[Rv(Tok_STAR)]"},
	{Test: "0026", Inp: "abc						: Rv(Tok_LT)", SkipTest: false, ResultTst2: "[Rv(Tok_LT)]"},
	{Test: "0027", Inp: "abc						: Rv(Tok_GT)", SkipTest: false, ResultTst2: "[Rv(Tok_GT)]"},
	{Test: "0028", Inp: "abc						: Rv(Tok_PCT)", SkipTest: false, ResultTst2: "[Rv(Tok_PCT)]"},
	{Test: "0029", Inp: "abc						: Rv(Tok_EQ)", SkipTest: false, ResultTst2: "[Rv(Tok_EQ)]"},
	{Test: "0030", Inp: "abc			: Ignore()						", SkipTest: false, ResultTst2: "[Ignore()]", ResultTst3: []string{"Ignore", ""}},
	{Test: "0031", Inp: "abc						: Warn(Warn_Unrecog_Char)	Return()", SkipTest: false, ResultTst2: "[Warn(Warn_Unrecog_Char) Return()]"},
	{Test: "0032", Inp: "abc					: Return()", SkipTest: false, ResultTst2: "[Return()]"},
	{Test: "0033", Inp: "abc					: Warn(Warn_End_Var_Unexpected)", SkipTest: false, ResultTst2: "[Warn(Warn_End_Var_Unexpected)]"},
	{Test: "0034", Inp: "abc						: Return()", SkipTest: false, ResultTst2: "[Return()]"},
	{Test: "0035", Inp: "abc						: Call(S_Quote)", SkipTest: false, ResultTst2: "[Call(S_Quote)]"},
	{Test: "0036", Inp: "abc						: Rv(Tok_Str0)", SkipTest: false, ResultTst2: "[Rv(Tok_Str0)]"},
	{Test: "0037", Inp: "abc						: Accept()", SkipTest: false, ResultTst2: "[Accept()]"},
	{Test: "0038", Inp: "abc						: Rv(Tok_Str1)", SkipTest: false, ResultTst2: "[Rv(Tok_Str1)]"},
	{Test: "0039", Inp: "abc						: Call(S_Quote)", SkipTest: false, ResultTst2: "[Call(S_Quote)]"},
	{Test: "0040", Inp: "abc						: Accept()", SkipTest: false, ResultTst2: "[Accept()]"},
	{Test: "0041", Inp: "abc					: Repl(\"`\")", SkipTest: false, ResultTst2: "[Repl(\"`\")]"},
	{Test: "0042", Inp: "abc						: Rv(Tok_Str0)", SkipTest: false, ResultTst2: "[Rv(Tok_Str0)]"},
	{Test: "0043", Inp: "abc						: Accept()", SkipTest: false, ResultTst2: "[Accept()]"},
	{Test: "0044", Inp: "abc						: Rv(Tok_SLASH)", SkipTest: false, ResultTst2: "[Rv(Tok_SLASH)]"},
	{Test: "0045", Inp: "abc						: Rv(Tok_OR)", SkipTest: false, ResultTst2: "[Rv(Tok_OR)]"},
	{Test: "0046", Inp: "abc 	: Rv(Tok_ABC)", SkipTest: false, ResultTst2: "[Rv(Tok_ABC)]", ResultTst4: "abc"},
	{Test: "0047", Inp: "`abc`	: Rv(Tok_ABC)", SkipTest: false, ResultTst2: "[Rv(Tok_ABC)]", ResultTst4: "abc"},
	{Test: "0048", Inp: "`ab``c`	: Rv(Tok_ABC)", SkipTest: false, ResultTst2: "[Rv(Tok_ABC)]", ResultTst4: "ab`c"},
	{Test: "0049", Inp: "'abc'	: Rv(Tok_ABC)", SkipTest: false, ResultTst2: "[Rv(Tok_ABC)]", ResultTst4: "abc"},
	{Test: "0050", Inp: "'ab\\'c'	: Rv(Tok_ABC)", SkipTest: false, ResultTst2: "[Rv(Tok_ABC)]", ResultTst4: "ab'c"},
	{Test: "0051", Inp: "\"abc\"	: Rv(Tok_ABC)", SkipTest: false, ResultTst2: "[Rv(Tok_ABC)]", ResultTst4: "abc"},
	{Test: "0052", Inp: `"ab\"c"	: Rv(Tok_ABC)`, SkipTest: false, ResultTst2: "[Rv(Tok_ABC)]", ResultTst4: "ab\"c"},
	{Test: "0053", Inp: "abc : Rv(Tok_ABC)", SkipTest: false, ResultTst2: "[Rv(Tok_ABC)]", ResultTst4: "abc"},

	{Test: "0200", Inp: "abc=12", SkipTest: false, ResultTst5: []string{"abc", "12"}},
	{Test: "0201", Inp: "abc", SkipTest: false, ResultTst5: []string{"abc", ""}},
	{Test: "0202", Inp: "a", SkipTest: false, ResultTst5: []string{"a", ""}},
	{Test: "0203", Inp: "abc=", SkipTest: false, ResultTst5: []string{"abc", ""}},
	{Test: "0204", Inp: "abc=Tok_Name", SkipTest: false, ResultTst5: []string{"abc", "Tok_Name"}},

	{Test: "0300", Inp: "Rv(abc)", SkipTest: false, ResultTst6: []string{"Rv", "abc"}},
	{Test: "0301", Inp: "Rv()", SkipTest: false, ResultTst6: []string{"Rv", ""}},
	{Test: "0301", Inp: "Rv", SkipTest: false, ResultTst6: []string{"Rv", ""}},

	{Test: "0400", Inp: "abc, def, ghi", SkipTest: false, ResultTst7: []string{"abc", "def", "ghi"}},
	{Test: "0401", Inp: "abc, def, ", SkipTest: false, ResultTst7: []string{"abc", "def"}},
	{Test: "0402", Inp: "", SkipTest: false, ResultTst7: []string{}},
	{Test: "0403", Inp: "Tokens, Tok_null=0, Tok_ID=1", SkipTest: false, ResultTst7: []string{"Tokens", "Tok_null=0", "Tok_ID=1"}},
}

// "[Rv(Tok_PCT)]"},

// func ParseAction(ln string) [][]string {
func Test_ParseAction(t *testing.T) {
	tst1 := false
	dbgo.SetADbFlag("in-echo-machine", true) // Output machine
	for ii, vv := range In01Test {
		if !vv.SkipTest {
			// fmt.Printf("\nTest %s ------------------------------------------------------------------------------------ \n", vv.Test)
			if tst1 {
				x := ParseAction(vv.Inp)
				fmt.Printf("%3d=%s\n", ii, dbgo.SVarI(x))
			}
			if vv.ResultTst4 != "" {

				cls := clasifyLine(vv.Inp)
				// fmt.Printf("Test %s cls: %s for -->>%s<<--, %s\n", vv.Test, cls, vv.Inp, dbgo.LF())
				atFront, rest := PickOffPatternAtBeginning(cls, vv.Inp)
				_ = rest

				if vv.ResultTst2 != "" {
					_, _, opt := ParsePattern(cls, vv.Inp)
					r := fmt.Sprintf("%v", opt)
					if r != vv.ResultTst2 {
						t.Errorf("Test %s Failed, Expected ->%s<- Got ->%s<-, %s\n", vv.Test, vv.ResultTst2, r, dbgo.LF())
					}
				}
				if atFront != vv.ResultTst4 {
					t.Errorf("Test %s Failed, Expected ->%s<- Got ->%s<-, %s\n", vv.Test, vv.ResultTst4, atFront, dbgo.LF())
				}

			} else if vv.ResultTst2 != "" {
				cls := clasifyLine(vv.Inp)
				pat, flag, opt := ParsePattern(cls, vv.Inp)
				r := fmt.Sprintf("%v", opt)
				if r != vv.ResultTst2 {
					t.Errorf("Test %s Failed, Expected ->%s<- Got ->%s<-, %s\n", vv.Test, vv.ResultTst2, r, dbgo.LF())
				}
				if false {
					fmt.Printf("%3d: %v %v %v\n", ii, pat, flag, opt)
				}
			}

			if len(vv.ResultTst5) > 0 {
				name, value := ParseNameValue(vv.Inp)
				// fmt.Printf("name=%s value=%s\n", name, value)
				if name != vv.ResultTst5[0] {
					t.Errorf("Test %s Failed, Expected ->%s<- Got ->%s<-, %s\n", vv.Test, vv.ResultTst5[0], name, dbgo.LF())
				}
				if value != vv.ResultTst5[1] {
					t.Errorf("Test %s Failed, Expected ->%s<- Got ->%s<-, %s\n", vv.Test, vv.ResultTst5[1], value, dbgo.LF())
				}
			}

			if len(vv.ResultTst6) > 0 {
				name, value := ParseActionItem(vv.Inp)
				// fmt.Printf("name=%s value=%s\n", name, value)
				if name != vv.ResultTst6[0] {
					t.Errorf("Test %s Failed, Expected ->%s<- Got ->%s<-, %s\n", vv.Test, vv.ResultTst6[0], name, dbgo.LF())
				}
				if value != vv.ResultTst6[1] {
					t.Errorf("Test %s Failed, Expected ->%s<- Got ->%s<-, %s\n", vv.Test, vv.ResultTst6[1], value, dbgo.LF())
				}
			}

			if len(vv.ResultTst7) > 0 {
				pl := ParsePlist(vv.Inp)
				if len(pl) != len(vv.ResultTst7) {
					t.Errorf("Test %s Failed, Expected %d length Got %d, %s\n", vv.Test, len(vv.ResultTst7), len(pl), dbgo.LF())
				}
				tt := fmt.Sprintf("%s", pl)
				ss := fmt.Sprintf("%s", vv.ResultTst7)
				if ss != tt {
					t.Errorf("Test %s Failed, Expected ->%s<- Got ->%s<-, %s\n", vv.Test, ss, tt, dbgo.LF())
				}
			}
		}
	}

	//xp := ParsePlist("abc, def, ghi")
	//fmt.Printf("xp=%+v\n", xp)
}

func Test_ParseFile(t *testing.T) {
	if false {
		Im := NewIm()
		fd := ReadFileIntoLines("./django3.lex")
		if len(fd) > 0 {
			Im.ParseFile(fd)
		}
		fmt.Printf("%+v\n", Im)
		Im.OutputImType()
	} else {
		Im := ImReadFile("./django3.lex")
		_ = Im
	}
}

// func ReadFileIntoString(fn string) string {
func Test_ReadFileIntoString(t *testing.T) {
	ss := ReadFileIntoString("./django3.lex")
	_ = ss
}

func Test_ReadFileIntoStringError(t *testing.T) {
	ss := ReadFileIntoString("./,,,,,")
	if ss != "" {
		t.Errorf("Expected error, got >%s<\n", ss)
	}
}

func Test_ReadFileIntoLinesError(t *testing.T) {
	ss := ReadFileIntoLines("./,,,,,")
	if len(ss) != 0 {
		t.Errorf("Expected error, got >%s<\n", ss)
	}
}

func Test_MiscTests(t *testing.T) {
	Add_Lookup_Token(1, "bobbob")
	x := Lookup_Tok_Name(1)
	if x != "bobbob" {
		t.Errorf("Token Look bad\n")
	}
	y := Lookup_Tok_Name(99999)
	if y == "bobbob" {
		t.Errorf("Token Look bad\n")
	}

	Im := NewIm()
	Im.LookupMachine("AName")
}

func Test_ClsString(t *testing.T) {
	s := fmt.Sprintf("%s", ClsPattern)
	_ = s
}

/* vim: set noai ts=4 sw=4: */
