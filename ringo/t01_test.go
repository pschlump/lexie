package ringo

// R I N G O - A template engine using the lexie system.

import (
	"fmt"
	"os"
	"testing"

	"github.com/pschlump/dbgo"
	"github.com/pschlump/filelib"
	"github.com/pschlump/lexie/com"
	"github.com/pschlump/lexie/pbread"

	"github.com/pschlump/lexie/dfa"
)

var data_01 = []string{
	//"./ref/test14.tpl",
	//"./ref/test15.tpl", // fiter...endfileter - test is not complete
	//"./ref/test16.tpl", // test comment...endcomment
	//"./ref/test17.tpl", // block item tests - no meaningufl tsts
	//"./ref/test18.tpl", // test just lorem
	//"./ref/test19.tpl", // causes infinite loop - error
	//"./ref/test20.tpl", // "for" test - core dump
	//"./ref/test28.tpl", // "for" test - mtest for expressions parsing corectly
	//"./ref/test30.tpl", // "for" test - empty array- with {%empty%} in it.
	// "./ref/test33.tpl", // {% csrf_token test %}
	//"./ref/test34.tpl", // {% dump_context %}
	// "./ref/test35.tpl", // {% extend "index.html" %}
	"./ref/test36.tpl", // {% extend "index2.html" %}
}

// ---------------------------------------------------------------------------------------------------------------------------------------
func Test_Test01_01(t *testing.T) {

	dbgo.SetADbFlag("trace-builtin", true)
	dbgo.SetADbFlag("match", true)

	dbgo.SetADbFlag("db_DumpDFAPool", true)
	dbgo.SetADbFlag("db_DumpPool", true)
	dbgo.SetADbFlag("db_Matcher_02", true)
	// dbgo.SetADbFlag("db_NFA_LnNo", true)
	dbgo.SetADbFlag("match", true)
	// dbgo.SetADbFlag("nfa3", true)
	dbgo.SetADbFlag("output-machine", true)
	dbgo.SetADbFlag("match", true)
	dbgo.SetADbFlag("match4", true)
	dbgo.SetADbFlag("match_x", true)
	// dbgo.SetADbFlag("nfa3", true)
	// dbgo.SetADbFlag("nfa4", true)
	// dbgo.SetADbFlag("db_DFAGen", true)
	// dbgo.SetADbFlag("pbbuf02", true)
	// dbgo.SetADbFlag("DumpParseNodes2", true)
	dbgo.SetADbFlag("db_FlushTokenBeforeBefore", true)
	dbgo.SetADbFlag("db_FlushTokenBeforeAfter", true)
	dbgo.SetADbFlag("db_tok01", true)
	dbgo.SetADbFlag("in-echo-machine", true) // Output machine

	Dbf = os.Stdout

	pt := NewParse2Type()
	pt.Lex = dfa.NewLexie()
	pt.Lex.SetChanelOnOff(true) // Set for getting back stuff via Chanel

	pt.Lex.NewReadFile("../in/django3.lex")

	pt.OpenLibraries("./tmpl")

	for _, fn := range data_01 {
		fn_o := com.RmExt(fn) + ".out"
		fn_r := com.RmExt(fn) + ".ref"
		pt.ReadFileAndRun(fn, fn_o)
		if !com.CompareFiles(fn_o, fn_r) {
			fmt.Printf("Files did not match %s %s\n", fn_o, fn_r)
			t.Errorf("%s error\n", fn_o)
		}
	}

}

// ---------------------------------------------------------------------------------------------------------------------------------------
func (pt *Parse2Type) ReadFileAndRun(fn, fn_o string) {

	go func() {
		r := pbread.NewPbRead()
		r.OpenFile(fn)
		pt.Lex.MatcherLexieTable(r, "S_Init")
	}()

	xpt := pt.GenParseTree(0)
	pt.TheTree = xpt
	pt.ExecuteFunctions(0)
	fmt.Printf("Tree Dump = %s\n", dbgo.SVarI(xpt))

	fp_o, err := filelib.Fopen(fn_o, "w")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
	} else {
		pt.OutputTree(fp_o, 0)
	}

	return
}
