package main

//
// C L I / T E S T 2 - Part of Lexie Lexical Generation System
//
// (C) Philip Schlump, 2014-2015.
// Version: 0.0.8
// BuildNo: 3
//
// Special Thanks to 2C-Why, LLC for supporting this project.
//

/*

 */

import (
	"fmt"
	"os"
	"strings"

	"github.com/pschlump/lexie/com"
	"github.com/pschlump/lexie/dfa"
	"github.com/pschlump/lexie/flags"
	"github.com/pschlump/lexie/pbread"
)

var opts struct {
	Config   string `short:"c" long:"config"      description:"Config Input File"       default:"./config.json"` //
	Input    string `short:"i" long:"input"       description:"Input File"              default:""`              //
	Output   string `short:"o" long:"output"      description:"Output File"             default:""`              //
	TraceOut string `short:"t" long:"trace"       description:"Trace Output"            default:""`              //
	Debug    string `short:"X" long:"debug"       description:"Debug Flags"             default:""`              //
}

func main() {
	var fp *os.File

	// ------------------------------------------------------ cli processing --------------------------------------------------------------
	ifnList, err := flags.ParseArgs(&opts, os.Args)

	if err != nil {
		fmt.Printf("Invalid Command Line: %s\n", err)
		os.Exit(1)
	}

	if opts.Debug != "" {
		s := strings.Split(opts.Debug, ",")
		com.DbOnFlags[opts.Debug] = true
		for _, v := range s {
			com.DbOnFlags[v] = true
		}
	}

	if opts.Echo != "" {
		com.DbOnFlags["in-echo-machine"] = true // Output machine
	}

	fmt.Fprintf(os.Stderr, "Test Matcher test from %s file, %s\n", opts.LexPat, com.LF())

	// ------------------------------------------------------ setup Lexie --------------------------------------------------------------
	pt := NewParse2Type()
	pt.Lex = dfa.NewLexie()
	pt.Lex.SetChanelOnOff(true) // Set for getting back stuff via Chanel

	// ------------------------------------------------------ input machine  --------------------------------------------------------------
	if opts.LexPat != "" {
		pt.Lex.NewReadFile(opts.LexPat) // pstk.Lex.NewReadFile("../in/django3.lex")
	} else if opts.ReadMachine != "" {
		fmt.Printf("Should input machine at this point\n")
		// xyzzy
	} else {
		fmt.Printf("Fatal: Must have -l <fn> or -r <fn>, neither supplied.\n")
		os.Exit(1)
	}

	if opts.Machine != "" {
		fmt.Printf("Should output machine at this point\n")
		// xyzzy
	}

	// -------------------------------------------------- start scanning process  ----------------------------------------------------------
	if opts.Tokens != "" {
		fp, _ = com.Fopen(opts.Tokens, "w")
	} else {
		fp = os.Stdout
	}

	if opts.Input != "" {

		go func() {
			r := pbread.NewPbRead()
			r.OpenFile(opts.Input)
			pt.Lex.MatcherLexieTable(r, "S_Init")
		}()

	} else {

		go func() {
			r := pbread.NewPbRead()
			for _, fn := range ifnList[1:] {
				r.OpenFile(fn)
			}
			pt.Lex.MatcherLexieTable(r, "S_Init")
		}()

	}

	// ------------------------------------------------------ process tokens --------------------------------------------------------------
	if false {
		// just print tokens out to check the scanning prcess and CLI options
		for msg := range pt.Lex.Message {
			fmt.Fprintf(fp, "%+v\n", msg)
		}
	} else {
		// Generate a parse tree and print out.
		xpt := pt.GenParseTree(0)
		pt.TheTree = xpt
		xpt.DumpMtType(fp, 0, 0)
		fmt.Printf("----------------------------------- start execute ----------------------------------------------------\n")
		pt.ExecuteFunctions(0)
		fmt.Printf("----------------------------------- debug output ----------------------------------------------------\n")
		if true {
			fmt.Printf("%s\n", com.SVarI(xpt))
		}
		fmt.Printf("----------------------------------- output ----------------------------------------------------\n")
		for i := 0; i < 1000000; i++ {
			pt.OutputTree0(fp, 0)
		}
		fmt.Printf("----------------------------------- errors ----------------------------------------------------\n")
		pp := pt.CollectErrorNodes(0)
		for ii, vv := range pp {
			fmt.Printf("Error [%3d]: msg=%s\n", ii, vv.ErrorMsg)
		}
		fmt.Printf("----------------------------------- final template results  ----------------------------------------------------\n")
		pt.OutputTree(fp, 0)
	}

	if opts.Tokens != "" {
		fp.Close()
	}

}

// ---------------------------------------------------------------------------------------------------------------------------------------
//
// expr ::= I
//		| B expr E
//		| html* expr
//		;
//
// ---------------------------------------------------------------------------------------------------------------------------------------
//
// Template is defined to be
//		Name
//		List of Named Paramters with default values
//		Set of tokens making up teplate { HTML*, Param(name), Template }
//
// Template defined by [ something ] extends [ template ] required [ list ]
//
// ---------------------------------------------------------------------------------------------------------------------------------------
