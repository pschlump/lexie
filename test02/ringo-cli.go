package main

//
// C L I / T E S T 2 - Part of Lexie Lexical Generation System
//
// Copyright (C) Philip Schlump, 2014-2025.
//

/*

 */

import (
	"fmt"
	"os"
	"strings"

	"github.com/jessevdk/go-flags"
	"github.com/pschlump/dbgo"
	"github.com/pschlump/filelib"
	"github.com/pschlump/lexie/dfa"
	"github.com/pschlump/lexie/pbread"
)

var opts struct {
	Config      string `short:"c" long:"config"      description:"Config Input File"       default:"./config.json"`     //
	Tokens      string `short:"t" long:"tokens"      description:"Token Output File"       default:""`                  //     *3* Output from running machine on Input File
	ReadMachine string `short:"r" long:"read"        description:"Machine Input File"      default:""`                  // <x>
	Input       string `short:"i" long:"input"       description:"Input File"              default:""`                  //
	Output      string `short:"o" long:"output"      description:"Output File"             default:""`                  //
	Machine     string `short:"m" long:"machine"     description:"Machine Output File"     default:""`                  // <x> Output in .mlex format
	Debug       string `short:"X" long:"debug"       description:"Debug Flags"             default:""`                  //
	Echo        string `short:"e" long:"echo"        description:"Output Machine "         default:""`                  //     *** <x> Output in .mlex format
	LexPat      string `short:"l" long:"lex"         description:"Lex Input File"          default:"../in/django3.lex"` //     *1* Input
	// TraceOut    string `short:"t" long:"trace"       description:"Trace Output"            default:""`                  //
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
		for _, s := range strings.Split(opts.Debug, ",") {
			dbgo.SetADbFlag(s, true)
		}
	}

	if opts.Echo != "" {
		dbgo.SetADbFlag("in-echo-machine", true) // Output machine
	}

	fmt.Fprintf(os.Stderr, "Test Matcher test from %s file, %s\n", opts.LexPat, dbgo.LF())

	// ------------------------------------------------------ setup Lexie --------------------------------------------------------------
	pt := NewParse2Type()
	pt.Lex = dfa.NewLexie()
	pt.Lex.SetChanelOnOff(true) // Set for getting back stuff via Chanel

	// ------------------------------------------------------ input machine  --------------------------------------------------------------
	if opts.LexPat != "" {
		pt.Lex.NewReadFile(opts.LexPat, "ringo") // pstk.Lex.NewReadFile("../in/django3.lex")
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
		fp, err = filelib.Fopen(opts.Tokens, "w")
		if err != nil {
			fmt.Fprintf(os.Stderr, "unable to open >%s< for output: %s\n", opts.Tokens, err)
			os.Exit(1)
		}
	} else {
		fp = os.Stdout
	}

	os.MkdirAll("./out", 0755)
	fCnst, err := filelib.Fopen("./out/const-stuff.go", "w")
	if err != nil {
		// xyzzy
	}
	fmt.Fprintf(fCnst, `
package Stuff

`)

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
			fmt.Printf("%s\n", dbgo.SVarI(xpt))
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
