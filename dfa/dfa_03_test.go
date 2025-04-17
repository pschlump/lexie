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
		Result: []Lr2Type{
			Lr2Type{StrTokNo: "Tok_CL", Match: "%}"}, // 0
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
	lex.NewReadFile("../in/test03_dfa.lex", "pct")

	in.DumpTokenMap()

	for ii, vv := range Lexie03Data {

		if vv.SkipTest {
			continue
		}

		dbgo.Printf("\n\n%(yellow)Test:%s ------------------------- Start --------------------------, %d, Input: -->>%s<<--\n", vv.Test, ii, vv.Inp)

		// r := strings.NewReader(vv.Inp)
		r := pbread.NewPbRead()                                                                                  // Create a push-back buffer
		dbgo.DbPrintf("trace-dfa-01 (../in/test03_dfa.lex scanner model)", "At: %(LF), Input: ->%s<-\n", vv.Inp) //
		r.PbString(vv.Inp)                                                                                       // set the input to the string
		r.SetPos(1, 1, fmt.Sprintf("sf-%d.txt", ii))                                                             // simulate  file = sf-%d.txt, set line to 1

		dbgo.DbPrintf("trace-dfa-01", "At: %(LF)\n")
		lex.MatcherLexieTable(r, "S_Init")
		dbgo.DbPrintf("trace-dfa-01", "At: %(LF)\n")

		if len(vv.Result) > 0 {
			dbgo.DbPrintf("trace-dfa-01", "At: %(LF)\n")
			if len(lex.TokList.TokenData) != len(vv.Result) {
				// fmt.Printf("Lengths did not match, %s", dbgo.SVarI(lex.TokList.TokenData))
				// c.Check(len(lex.TokList.TokenData), Equals, len(vv.Result)) // xyzzy
				dbgo.DbPrintf("trace-dfa-01", "At: %(LF)\n")
				t.Errorf("Length did not match, expected %d tokens, got %d\n", len(lex.TokList.TokenData), len(vv.Result))
			} else {
				for i := 0; i < len(vv.Result); i++ {
					dbgo.DbPrintf("trace-dfa-01", "At: %(LF)\n")
					if vv.Result[i].StrTokNo != "" {
						// func in.LookupTokenName(Tok int) (rv string) { -- use to repace token numbers '38' with Token Name and lookup for test
						// c.Check(vv.Result[i].StrTokNo, Equals, in.LookupTokenName(int(lex.TokList.TokenData[i].TokNo))) // xyzzy
						dbgo.DbPrintf("trace-dfa-01", "At: %(LF)\n")
						if vv.Result[i].StrTokNo != in.LookupTokenName(int(lex.TokList.TokenData[i].TokNo)) {
							t.Errorf("Invalid token found.  Expected %d/%s got %d/%s\n", lex.TokList.TokenData[i].TokNo, in.LookupTokenName(int(lex.TokList.TokenData[i].TokNo)),
								int(lex.TokList.TokenData[i].TokNo), in.LookupTokenName(int(lex.TokList.TokenData[i].TokNo)))
						}
					} else {
						dbgo.DbPrintf("trace-dfa-01", "At: %(LF)\n")
						// c.Check(vv.Result[i].TokNo, Equals, int(lex.TokList.TokenData[i].TokNo)) // xyzzy
						if vv.Result[i].TokNo != int(lex.TokList.TokenData[i].TokNo) {
							t.Errorf("Invalid token found.  Expected %d/%s got %d/%s\n", lex.TokList.TokenData[i].TokNo, in.LookupTokenName(int(lex.TokList.TokenData[i].TokNo)),
								int(lex.TokList.TokenData[i].TokNo), in.LookupTokenName(int(lex.TokList.TokenData[i].TokNo)))
						}
					}
					/*
						// c.Check(vv.Result[i].Match, Equals, lex.TokList.TokenData[i].Match) // xyzzy
						dbgo.DbPrintf("trace-dfa-01", "At: %(LF)\n")
						if vv.Result[i].LineNo > 0 {
							dbgo.DbPrintf("trace-dfa-01", "At: %(LF)\n")
							// c.Check(vv.Result[i].LineNo, Equals, lex.TokList.TokenData[i].LineNo) // xyzzy
						}
						if vv.Result[i].ColNo > 0 {
							dbgo.DbPrintf("trace-dfa-01", "At: %(LF)\n")
							// c.Check(vv.Result[i].ColNo, Equals, lex.TokList.TokenData[i].ColNo) // xyzzy
						}
						if vv.Result[i].FileName != "" {
							dbgo.DbPrintf("trace-dfa-01", "At: %(LF)\n")
							// c.Check(vv.Result[i].FileName, Equals, lex.TokList.TokenData[i].FileName) // xyzzy
						}
					*/
				}
			}
		}

		dbgo.DbPrintf("trace-dfa-01", "At: %(LF)\n")
		fmt.Printf("Test:%s ------------------------- End --------------------------\n\n", vv.Test)

	}

}
