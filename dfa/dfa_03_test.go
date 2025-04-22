package dfa

import (
	"fmt"
	"os"
	"testing"

	"github.com/pschlump/dbgo"
	"github.com/pschlump/lexie/in"
	"github.com/pschlump/lexie/pbread"
)

var Lexie03Data = []Lexie02DataType{
	{Test: "4300", Inp: "%}", Rv: 4300, SkipTest: false,
		ExpectedTokens: []Lr2Type{
			// Lr2Type{StrTokNo: "Tok_CL", Match: "%}"}, // 0 ***correct*** xyzzy847
			Lr2Type{StrTokNo: "Tok_PCT", Match: "%}"}, // 0
		},
	},
}

func Test_03_DfaTest03(t *testing.T) {

	dbgo.Fprintf(os.Stderr, "\n\n%(cyan)Test Matcher test from ../in/test03_dfa.lex file, %(LF)\n========================================================================\n\n")

	dbgo.SetADbFlag("db_DumpDFAPool", true)
	dbgo.SetADbFlag("db_DumpPool", true)
	dbgo.SetADbFlag("db_Matcher_02", true)
	// dbgo.SetADbFlag("db_NFA_LnNo", true)
	dbgo.SetADbFlag("match", true)
	dbgo.SetADbFlag("nfa3", true)
	dbgo.SetADbFlag("output-machine", true)
	dbgo.SetADbFlag("match", true)
	dbgo.SetADbFlag("match_x", true)
	dbgo.SetADbFlag("nfa3", true)
	dbgo.SetADbFlag("nfa4", true)
	// dbgo.SetADbFlag("db_DFAGen", true)
	// dbgo.SetADbFlag("pbbuf02", true)
	// dbgo.SetADbFlag("DumpParseNodes2", true)
	dbgo.SetADbFlag("db_FlushTokenBeforeBefore", true)
	dbgo.SetADbFlag("db_FlushTokenBeforeAfter", true)
	dbgo.SetADbFlag("db_tok01", true)
	dbgo.SetADbFlag("in-echo-machine", true) // Output machine

	lex := NewLexie()
	Machine := "../in/test03_dfa.lex"
	lex.NewReadFile(Machine, "pct")

	in.DumpTokenMap()

	lex.GenerateTokenMap("./out/token_map.go")

	for ii, vv := range Lexie03Data {
		if vv.SkipTest {
			continue
		}

		dbgo.Printf("\n\n%(yellow)Test:%s ------------------------- Start [%s] --------------------------, %d, Input: -->>%s<<--\n", vv.Test, Machine, ii, vv.Inp)

		// ---------------------------------------------------------------------------------
		// Read in input
		// ---------------------------------------------------------------------------------
		// r := strings.NewReader(vv.Inp)
		// r := pbread.NewStringReader(vv.Inp)	// todo.
		r := pbread.NewPbRead()                                                                                  // Create a push-back buffer
		dbgo.DbPrintf("trace-dfa-03 (../in/test03_dfa.lex scanner model)", "At: %(LF), Input: ->%s<-\n", vv.Inp) //
		r.PbString(vv.Inp)                                                                                       // set the input to the string
		r.SetPos(1, 1, fmt.Sprintf("sf-%d.txt", ii))                                                             // simulate  file = sf-%d.txt, set line to 1, this takes input from a string instead of from a file.

		// ---------------------------------------------------------------------------------
		// Generate machine
		// ---------------------------------------------------------------------------------
		dbgo.DbPrintf("trace-dfa-03", "At: %(LF) --- generate machine ---\n") //
		lex.FinializeMachines()

		// ---------------------------------------------------------------------------------
		// Run interepreted matcher
		// ---------------------------------------------------------------------------------
		dbgo.DbPrintf("trace-dfa-03", "At: %(LF) --- run matcher -- \n") //
		lex.MatcherLexieTable(r, "S_Init")                               // Run the matcing machine

		// Results are in lex.TokList.TokenData ************************************************************************************
		if len(vv.ExpectedTokens) > 0 {
			dbgo.DbPrintf("trace-dfa-03", "At: %(LF)\n")
			if len(lex.TokList.TokenData) != len(vv.ExpectedTokens) {
				// fmt.Printf("Lengths did not match, %s", dbgo.SVarI(lex.TokList.TokenData))
				// c.Check(len(lex.TokList.TokenData), Equals, len(vv.ExpectedTokens)) // xyzzy
				dbgo.DbPrintf("trace-dfa-03", "At: %(LF)\n")
				t.Errorf("Length did not match, expected %d tokens, got %d\n", len(lex.TokList.TokenData), len(vv.ExpectedTokens))
			} else {
				for i := 0; i < len(vv.ExpectedTokens); i++ {
					dbgo.DbPrintf("trace-dfa-03", "At: %(LF)\n")
					if vv.ExpectedTokens[i].StrTokNo != "" {
						// func in.LookupTokenName(Tok int) (rv string) { -- use to repace token numbers '38' with Token Name and lookup for test
						// c.Check(vv.ExpectedTokens[i].StrTokNo, Equals, in.LookupTokenName(int(lex.TokList.TokenData[i].TokNo))) // xyzzy
						dbgo.DbPrintf("trace-dfa-03", "At: %(LF)\n")
						if vv.ExpectedTokens[i].StrTokNo != in.LookupTokenName(int(lex.TokList.TokenData[i].TokNo)) {
							t.Errorf("Invalid token found.  Expected %d/%s got %d/%s\n",
								vv.ExpectedTokens[i].TokNo, in.LookupTokenName(int(vv.ExpectedTokens[i].TokNo)),
								int(lex.TokList.TokenData[i].TokNo), in.LookupTokenName(int(lex.TokList.TokenData[i].TokNo)),
							)
						}
					} else {
						dbgo.DbPrintf("trace-dfa-03", "At: %(LF)\n")
						// c.Check(vv.ExpectedTokens[i].TokNo, Equals, int(lex.TokList.TokenData[i].TokNo)) // xyzzy
						if vv.ExpectedTokens[i].TokNo != int(lex.TokList.TokenData[i].TokNo) {
							t.Errorf("Invalid token found.  Expected %d/%s got %d/%s\n",
								int(vv.ExpectedTokens[i].TokNo), in.LookupTokenName(int(vv.ExpectedTokens[i].TokNo)),
								lex.TokList.TokenData[i].TokNo, in.LookupTokenName(int(lex.TokList.TokenData[i].TokNo)),
							)
						}
					}
					/*
						// c.Check(vv.ExpectedTokens[i].Match, Equals, lex.TokList.TokenData[i].Match) // xyzzy
						dbgo.DbPrintf("trace-dfa-03", "At: %(LF)\n")
						if vv.ExpectedTokens[i].LineNo > 0 {
							dbgo.DbPrintf("trace-dfa-03", "At: %(LF)\n")
							// c.Check(vv.ExpectedTokens[i].LineNo, Equals, lex.TokList.TokenData[i].LineNo) // xyzzy
						}
						if vv.ExpectedTokens[i].ColNo > 0 {
							dbgo.DbPrintf("trace-dfa-03", "At: %(LF)\n")
							// c.Check(vv.ExpectedTokens[i].ColNo, Equals, lex.TokList.TokenData[i].ColNo) // xyzzy
						}
						if vv.ExpectedTokens[i].FileName != "" {
							dbgo.DbPrintf("trace-dfa-03", "At: %(LF)\n")
							// c.Check(vv.ExpectedTokens[i].FileName, Equals, lex.TokList.TokenData[i].FileName) // xyzzy
						}
					*/
				}
			}
		}

		dbgo.DbPrintf("trace-dfa-03", "At: %(LF)\n")
		fmt.Printf("Test:%s ------------------------- End --------------------------\n\n", vv.Test)

	}

}
