//
// C L I - Part of Lexie Lexical Generation System
//
// (C) Philip Schlump, 2014-2015.
// Version: 1.0.8
// BuildNo: 203
//
// Special Thanks to 2C-Why, LLC for supporting this project.
//

/*

Step 3
	1. Match up interface

Step 4
	1. Performance! Memory!
	2. Swap over


// 1. CSV file
// 2. JSONx
// 4. Go-Teplate scan
// 5. C scanner
// 6. Our own AWK like language ( Redis and PostgreSQL enabled )
*/

package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/jessevdk/go-flags"
	"github.com/pschlump/dbgo"
	"github.com/pschlump/filelib"
	"github.com/pschlump/lexie/com"
	"github.com/pschlump/lexie/dfa"
	"github.com/pschlump/lexie/in"
	"github.com/pschlump/lexie/pbread"
	"github.com/pschlump/lexie/tok"
)

var opts struct {
	LexPat      string `short:"l" long:"lex"         description:"Lex Input File"             default:""`                   //     *1* Input
	ReadMachine string `short:"r" long:"read"        description:"Machine Input File"         default:""`                   // <x>
	Echo        string `short:"e" long:"echo"        description:"Output Machine "            default:""`                   //     *** <x> Output in .mlex format
	Machine     string `short:"m" long:"machine"     description:"Machine Output File"        default:""`                   // <x> Output in .mlex format
	Output      string `short:"o" long:"output"      description:"Output File"                default:""`                   // <x>
	Format      string `short:"f" long:"fmt"         description:"Format of output File"      default:"go"`                 // <x>
	Input       string `short:"i" long:"input"       description:"Input File"                 default:""`                   //     *2* Match Input
	Tokens      string `short:"t" long:"tokens"      description:"Token Output File"          default:""`                   //     *3* Output from running machine on Input File
	DotPath     string `short:"D" long:"dotpath"     description:"Path to DOT"                default:"/usr/local/bin/dot"` // <x>
	DotOutput   string `short:"O" long:"dotoutput"   description:"Generate in DOT"            default:""`                   // <x> D, N, R - Dfa, Nfa, Re-Parse-Tree
	DotFN       string `short:"F" long:"dotfn"       description:"Generate in DOT file name"  default:""`                   // <x> Base file name to generate output in -F xx results in xx.D.svg, xx.N.sfg, xx.R.json
	Debug       string `short:"X" long:"debug"       description:"Debug Flags"                default:""`                   //     *
}

// out, err = exec.Command("/usr/local/bin/dot", "-Tsvg", "-o"+svgFile, gvFile).Output()

const (
	TokenError = iota
	EOF

	TokenHTML

	TokenKeyword
	TokenIdentifier
	TokenString
	TokenNumber
	TokenSymbol
)

const (
	Tok_ID     = 39
	Tok_HTML   = 40
	Tok_NUM    = 41
	Tok_Str0   = 42
	Tok_Str1   = 43
	Tok_Str2   = 44
	Tok_and    = 29
	Tok_as     = 34
	Tok_export = 35
	Tok_false  = 33
	Tok_in     = 28
	Tok_not    = 31
	Tok_or     = 28
	Tok_true   = 32
	Tok_X      = 38
)

type TokenMaperType map[int]int

var MapToken = TokenMaperType{
	Tok_ID:   TokenIdentifier,
	Tok_HTML: TokenHTML,
	Tok_NUM:  TokenNumber,
	Tok_Str0: TokenString,
	Tok_Str1: TokenString,
	Tok_Str2: TokenString,
	Tok_X:    TokenHTML,
}

// $def(Tokens, Tok_L_EQ=1, Tok_GE=2, Tok_LE=3, Tok_L_AND=4, Tok_L_OR=5 , Tok_OP_VAR=6, Tok_CL_VAR=7, Tok_OP_BL=8, Tok_CL_BL=9 , Tok_NE=10, Tok_NE_LG=11, Tok_OP=12, Tok_CL=13, Tok_PLUS=14, Tok_MINUS=15, Tok_STAR=16, Tok_LT=17, Tok_GT=18, Tok_SLASH=19, Tok_CARRET=20, Tok_COMMA=21, Tok_DOT=22, Tok_EXCLAM=23, Tok_OR=24, Tok_COLON=25, Tok_EQ=26, Tok_PCT=27, Tok_in=28, Tok_and=29, Tok_or=30, Tok_not=31, Tok_true=32, Tok_false=33, Tok_as=34, Tok_export=35, Tok_SS)
var MapRW = map[string]int{
	"and":    Tok_and,
	"as":     Tok_as,
	"export": Tok_export,
	"false":  Tok_false,
	"in":     Tok_in,
	"not":    Tok_not,
	"or":     Tok_or,
	"true":   Tok_true,
}

func CategorizeToken(tl []tok.Token) []tok.Token {
	for ii, vv := range tl {
		// Take strings and pre-process them - remove trailing quote mark.
		// fmt.Printf("CT: ii=%d TokNo=%d >%s<- ->%s<-\n", ii, vv.TokNo, vv.Val, vv.Match)
		if (vv.TokNo == Tok_Str0 || vv.TokNo == Tok_Str1 || vv.TokNo == Tok_Str2) && len(vv.Val) > 1 {
			vv.Val = vv.Val[0 : len(vv.Val)-1]
			vv.Match = vv.Match[0 : len(vv.Match)-1]
			// fmt.Printf("xx: ii=%d TokNo=%d >%s<- ->%s<-\n", ii, vv.TokNo, vv.Val, vv.Match)
		}
		// If it is an ID then see if this ID is a ReservedWord
		if vv.TokNo == Tok_ID {
			vv.Typ = TokenIdentifier
			if t, ok := MapRW[vv.Match]; ok {
				// vv.TokNo = tok.TokenNoType(t)
				vv.TokNo = int(t)
				vv.Typ = TokenKeyword
			}
			vv.TokNo = 0
		} else if x, ok := MapToken[int(vv.TokNo)]; ok { // If in the MapToken map then set the Typ
			vv.Typ = tok.TokenType(x)
			vv.TokNo = 0
		} else {
			vv.Typ = TokenSymbol
		}
		tl[ii] = vv
	}
	// Discard an initial empty token if one exists
	// fmt.Printf("yy: Typ=%d TokNo=%d >%s<- ->%s<-\n", tl[0].Typ, tl[0].TokNo, tl[0].Val, tl[0].Match)
	if len(tl) > 1 && tl[0].TokNo == 0 && tl[0].Typ == 2 && tl[0].Val == "" {
		// fmt.Printf("Match for yy\n")
		tl = tl[1:]
	}
	return tl
}

func main() {
	var fp *os.File

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
		dbgo.SetADbFlag(opts.Debug, true)
	}

	if opts.Echo != "" {
		com.DbOnFlags["in-echo-machine"] = true // Output machine
		dbgo.SetADbFlag("in-echo-machine", true)
	}

	fmt.Fprintf(os.Stderr, "Test Matcher test from %s file, %s\n", opts.LexPat, dbgo.LF())

	lex := dfa.NewLexie()
	if opts.LexPat != "" {
		lex.NewReadFile(opts.LexPat) // lex.NewReadFile("../in/django3.lex")
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

	if opts.Tokens != "" {
		fp, _ = filelib.Fopen(opts.Tokens, "w")
	} else {
		fp = os.Stdout
	}
	if opts.Input != "" {

		s := in.ReadFileIntoString(opts.Input)
		// 		r := strings.NewReader(s)
		r := pbread.NewPbRead() // func NewPbRead() (rv *PBReadType) { // PJS Sun Oct 31 13:47:42 MDT 2021
		r.PbString(s)

		lex.MatcherLexieTable(r, "S_Init")

		lex.TokList.TokenData = CategorizeToken(lex.TokList.TokenData)

		lex.DumpTokenBuffer(fp)

	} else {

		for ii, fn := range ifnList[1:] {
			s := in.ReadFileIntoString(fn)
			// r := strings.NewReader(s)
			r := pbread.NewPbRead() // func NewPbRead() (rv *PBReadType) { // PJS Sun Oct 31 13:47:42 MDT 2021
			r.PbString(s)
			lex.MatcherLexieTable(r, "S_Init")

			lex.TokList.TokenData = CategorizeToken(lex.TokList.TokenData)

			fmt.Fprintf(fp, "%d: %s -----Start---------------------------------------------------------------------------------\n", ii, fn)
			lex.DumpTokenBuffer2(fp)
			fmt.Fprintf(fp, "%d: %s -----End-----------------------------------------------------------------------------------\n\n\n", ii, fn)
		}

	}
	if opts.Tokens != "" {
		fp.Close()
	}

}
