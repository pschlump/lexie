//
// M A T C H - Part of Lexie Lexical Generation System
//
// (C) Philip Schlump, 2014-2015.
// Version: 1.0.8
// BuildNo: 203
//
// Special Thanks to 2C-Why, LLC for supporting this project.
//

/*
// ------------------------------------------------------------------------------------------------------------------------------------------------
// Plan
// ------------------------------------------------------------------------------------------------------------------------------------------------

+1. Study Django templates
	+2. Understand what they do
	+2. Decide on replace or use code
		http://www.b-list.org/weblog/2007/sep/22/standalone-django-scripts/

+0. What I want is to be able to get a lit of tokens/class/values out as an output from the CLI

+0. Use CLI and turn off all extraneous output -- Validate matches on a bunch of machines

0. Output machine as a .mlex format - machine lex format -m <fn>
0. Output machine as a .go format --fmt go -o <fn>
0. add "$cfg" to input - for config params like size of smap
0. add a SetConfig() to code to set the config params

+1. Use in Ringo
+2. The Pongo2 v.s. Ringo test
+3. Coverage Testing
4. Benchmarks
4. Memory usage - eliminate dynamic allocations
+4. Internal comments (doc.go) etc.
+5. As a "go-routine"



// ------------------------------------------------------------------------------------------------------------------------------------------------
//
// TODO:
//		1. Failure to match causes a fall-out-the-bottom and quit behavior - change to a reset to 0 and continue to EOF
//				1. Warnings about bad tokens are not reported a {{ token should produce an error about invalid range.
// 					1. xyzzyStackEmpty - if stack not empty then - pop and continue running on different data??a
//		1. In nfa->dfa conversion may end up with non-distinct terminal values - use one that is first seen? -- Lex implies this.
//		1. Rune Fix ( ../re/re.go:800 ) + test cases
//					 xyzzyRune  TODO - this code is horribley non-utf8 compatable at this moment in time.
//
//	Err-Output:
//			1. Tests in ../in/in.go - are not automated - fix to check results of input to be correct
//				5. No test case for A_Reset
//				7. No Test case for this error: --- CCL --- 0-0 is an invalid CCL, 1-0 also etc. -- Report
//
//	Feature
//			1. Options on size of smap
//
// ------------------------------------------------------------------------------------------------------------------------------------------------
*/

package dfa

import (
	"fmt"
	"os"

	"github.com/pschlump/lexie/com"
	"github.com/pschlump/lexie/pbread"
	"github.com/pschlump/lexie/re"
	"github.com/pschlump/lexie/tok"
)

// ---------------------------------------------------------------------------------------------------------------------------------------------------------------
type MatchContextType struct {
	St  int
	Dfa *DFA_PoolType
}

func (lex *Lexie) MatcherLexieTable(rrr *pbread.PBReadType, s_init string) {

	var dfa *DFA_PoolType
	var ii, to, pos_no, col_no, line_no int
	var rn rune
	var SMatch, filename string
	init := 0

	lex.FinializeMachines()

	init = lex.Im.LookupMachine(s_init)
	for ii := range lex.DFA_Machine {
		fmt.Printf("At: %s\n", com.LF())
		com.DbPrintf("match", "Machine[%d] =\n", ii)
		dfa = lex.DFA_Machine[ii]
		dfa.MachineId = ii
		dfa.OutputInFormat(os.Stdout, "text")

		lex.OutputActionFlags(dfa)
	}
	dfa = lex.DFA_Machine[init]

	ctx := NewContext(dfa.MTab.InitState, nil)
	ctx_stack := make([]*MatchContextType, 0, 100)

	lex.TokList = tok.NewTokenList()
	ignoreToken := lex.Im.Lookup("Tokens", "Tok_Ignore")
	aTok_EOF := lex.Im.Lookup("Tokens", "Tok_EOF")
	lex.TokList.IgnoreToken = ignoreToken
	line_no = 1
	col_no = 1
	pos_no = 1
	filename = ""
	SMatch = ""
	AtEOF := false
	done := false
	TokStart := 0

	Next := func() (rn rune) {
		rn, done := rrr.NextRune()
		if done {
			AtEOF = true
			rn = re.R_EOF
			return
		} else {
			SMatch += string(rn)
			pos_no++
		}
		return
	}

	Peek := func() (rn rune) {
		rn, done = rrr.PeekRune()
		if done {
			AtEOF = true
			rn = re.R_EOF
		}
		return
	}

	var p_line_no int = 1
	var p_col_no int = 1
	SaveToken := func() {
		start := pos_no - dfa.MTab.Machine[ctx.St].Info.MatchLength
		if dfa.MTab.Machine[ctx.St].Info.MatchLength == 0 {
			start = 0
		}
		line_no, col_no, filename = rrr.GetPos()
		if col_no > dfa.MTab.Machine[ctx.St].Info.MatchLength {
			col_no -= dfa.MTab.Machine[ctx.St].Info.MatchLength
		}

		lRv := dfa.MTab.Machine[ctx.St].Rv

		// xyzzy - this is the spot to convert from Tok_ID && ReservedWord -> a new ID
		// lex.Im.St.Lookup.LookupSymbol(SMatch)
		// func (st *SymbolTable) LookupSymbol(name string) (as *SymbolType, err error) {
		// vv=in.ImDefinedValueType {Seq:1 WhoAmI:ReservedWords NameValueStr:map[and:Tok_L_AND not:Tok_not as:Tok_as in:Tok_in bor:Tok_B_OR band:Tok_B_AND xor:Tok_XOR or:Tok_L_OR true:Tok_true false:Tok_false export:Tok_export] NameValue:map[and:4 true:32 as:34 bor:42 band:41 xor:64 or:5 false:33 not:31 export:35 in:28] Reverse:map[5:or 32:true 42:bor 31:not 41:band 35:export 33:false 28:in 64:xor 4:and 34:as] SeenAt:map[bor:{LineNo:[39] FileName:[unk-file]} band:{LineNo:[39] FileName:[unk-file]} and:{LineNo:[39] FileName:[unk-file]} true:{LineNo:[39] FileName:[unk-file]} export:{LineNo:[39] FileName:[unk-file]} in:{LineNo:[39] FileName:[unk-file]} as:{LineNo:[39] FileName:[unk-file]} or:{LineNo:[39] FileName:[unk-file]} false:{LineNo:[39] FileName:[unk-file]} not:{LineNo:[39 39] FileName:[unk-file unk-file]} xor:{LineNo:[39] FileName:[unk-file]}]}, File: /Users/corwin/Projects/pongo2/lexie/dfa/match.go LineNo:260
		if dfa.MTab.Machine[ctx.St].Info.ReservedWord {
			com.DbPrintf("match4", "\n\nFound a  reserved word ------------------------------------------------------------------------ CCddEE, %s\n", com.LF())
			com.DbPrintf("match4", "SMatch >%s<, %s\n", SMatch, com.LF())

			vv, ok := lex.Im.Def.DefsAre["ReservedWords"].NameValue[SMatch]
			if ok {
				com.DbPrintf("match4", "Substituting(2) -- xyzzy hard coded for test --- Found it as %d !!!, %s\n", vv, com.LF())
				lRv = vv
			}

			com.DbPrintf("match4", "\n\n")
		}

		com.DbPrintf("match", "+==================================================================================\n")
		com.DbPrintf("match", "+ SaveToken: Line:%d Col:%d\n", line_no, col_no)
		com.DbPrintf("match", "+==================================================================================\n")
		lex.TokList.AddTokenToBuffer(tok.Token{
			Match:    SMatch,
			Val:      SMatch,
			LineNo:   p_line_no,
			ColNo:    p_col_no,
			FileName: filename,
			TokNo:    lRv,
		}, start, pos_no, 1)
		p_line_no = line_no
		p_col_no = col_no
	}

	FlushToken := func(isHard bool) {
		if dfa.MTab.Machine[ctx.St].Rv == 9900 && (com.A_Repl&dfa.MTab.Machine[ctx.St].Info.Action) != 0 {
			com.DbPrintf("match", " Doing Replace on Token, Len %d to ->%s<-\n ", dfa.MTab.Machine[ctx.St].Info.MatchLength, dfa.MTab.Machine[ctx.St].Info.ReplStr)
			lex.TokList.ReplaceToken(dfa.MTab.Machine[ctx.St].Info.MatchLength, dfa.MTab.Machine[ctx.St].Info.ReplStr)
		}
		com.DbPrintf("match4", "FlushTokenStareted ------------------------------------------------------------------------------\n")
		beforeFlush := len(lex.TokList.TokenData)
		lex.TokList.FlushTokenBuffer(TokStart, isHard, AtEOF)
		afterFlush := len(lex.TokList.TokenData)

		if lex.SendOnChanel {
			for ll := beforeFlush; ll < afterFlush; ll++ {
				fmt.Printf("At: %s\n", com.LF())

				tt := LexieChanelType{Token: lex.TokList.TokenData[ll]}
				lex.Message <- tt
			}
		}
		TokStart = 0
		SMatch = ""
		com.DbPrintf("match4", "FlushTokenEnded ------------------------------------------------------------------------------\n")
	}

	DumpStack := func() {
		com.DbPrintf("match4", "\tDump ctx_stack len=[%d]\n", len(ctx_stack))
		for ii := range ctx_stack {
			fmt.Printf("At: %s\n", com.LF())
			com.DbPrintf("match4", "\t\tDump ctx_stack[%d] = machine %d\n", ii, ctx_stack[ii].Dfa.MachineId)
		}
	}

	PushState := func(St int) {
		com.DbPrintf("match", "\n-------------------------------------------------------------------------------\n")
		com.DbPrintf("match", "Found a PUSH: to %d -- This should be subscript of new machine\n", dfa.MTab.Machine[St].Info.NextState)
		com.DbPrintf("match", "-------------------------------------------------------------------------------\n")
		t_ctx := NewContext(0, dfa)

		ns := dfa.MTab.Machine[St].Info.NextState
		ctx_stack = append(ctx_stack, t_ctx)
		DumpStack()

		ctx.St = 0
		dfa = lex.DFA_Machine[ns]
	}

	PopState := func() {
		if len(ctx_stack) >= 1 {
			t_ctx := ctx_stack[len(ctx_stack)-1]
			ctx_stack = ctx_stack[0 : len(ctx_stack)-1]
			dfa = t_ctx.Dfa
			com.DbPrintf("match", "\n-------------------------------------------------------------------------------\n")
			com.DbPrintf("match", "Found a POP: State Poping like Corn, Machine now %d\n", dfa.MachineId)
			com.DbPrintf("match", "-------------------------------------------------------------------------------\n")
			DumpStack()
		} else {
			com.DbPrintf("match", "Error: Attempt to pop when stack is empty\n")
		}
		ctx.St = 0
	}

	ResetState := func() {
		com.DbPrintf("match", "Found a RESET\n")
		ctx_stack = ctx_stack[:0]
		dfa = lex.DFA_Machine[init]
		ctx.St = 0
	}

	for !done {
		com.DbPrintf("match", "\n**********************************************************************************************************\n")
		com.DbPrintf("match", "Top: (machine number %d) ctx.St:%d\n", dfa.MachineId, ctx.St)
		cur_st := ctx.St

		rn = Peek()
		ii = dfa.MTab.SMap.MapRune(rn)
		to = dfa.MTab.Machine[cur_st].To[ii]
		hh := dfa.MTab.Machine[cur_st].Info.HardMatch
		com.DbPrintf("match", " (peek) machine=%d cur_st=%d ii=%d for rune rn=->%s<- to=%d hard=%v Rv=%d Tau=%v Action=%x, ctx_stack=%d\n",
			dfa.MachineId, cur_st, ii, string(rn), to, hh, dfa.MTab.Machine[cur_st].Rv, dfa.MTab.Machine[cur_st].Tau,
			dfa.MTab.Machine[cur_st].Info.Action, ctx_stack)
		if rn == re.R_EOF {
			com.DbPrintf("match", " At: %s\n", com.LF())
			FlushToken(true)
		} else if dfa.MTab.Machine[cur_st].Rv > 0 && dfa.MTab.Machine[cur_st].Tau && to == -1 {
			com.DbPrintf("match", " At: %s\n", com.LF())
			SaveToken()
			FlushToken(hh)
		} else if dfa.MTab.Machine[cur_st].Rv > 0 && hh {
			if to == -1 {
				com.DbPrintf("match", " At: %s\n", com.LF())
				SaveToken()
				FlushToken(hh)
			} else {
				com.DbPrintf("match", " At: %s\n", com.LF())
				// (peek) machine=0 cur_st=3 ii=5 for rune rn=-> <- to=0 hard=true Rv=8 Tau=true
				if (com.A_Push & dfa.MTab.Machine[cur_st].Info.Action) == 0 {
					com.DbPrintf("match", " At: %s\n", com.LF())
					rn = Next() // Remove makes beginning work
				}
				SaveToken()
				FlushToken(hh) // PJS added Wed Jul 15 15:43:28 MDT 2015
			}
			ctx.St = to
		} else if dfa.MTab.Machine[cur_st].Rv > 0 { // && !dfa.MTab.Machine[cur_st].Tau && to != -1 {
			com.DbPrintf("match", " At: %s\n", com.LF())
			rn = Next()
			SaveToken()
			ctx.St = to
		} else if ((com.A_Push | com.A_Pop | com.A_Reset) & dfa.MTab.Machine[cur_st].Info.Action) != 0 {
			com.DbPrintf("match", "-- critical -- At: %s\n", com.LF())
			// (peek) machine=2 cur_st=21 ii=24 for rune rn=->{<- to=36 hard=false Rv=0 Tau=false
			// (peek) machine=2 cur_st=24 ii=4 for rune rn=-> <- to=-1 hard=true Rv=0 Tau=false
			if to >= 0 {
				rn = Next()
			}
			ctx.St = to
		} else {
			com.DbPrintf("match", "-- other --  At: %s\n", com.LF())
			if to >= 0 {
				rn = Next()
			}
			ctx.St = to
		}

		if (com.A_Push & dfa.MTab.Machine[cur_st].Info.Action) != 0 { // 					xyzzy - total bullshit -- Allows to hapen on non-termaial in middle of stuff -- call non-terminal
			com.DbPrintf("match", "PUSH: At: %s\n", com.LF())
			PushState(cur_st)
		} else if (com.A_Pop & dfa.MTab.Machine[cur_st].Info.Action) != 0 { // PJS modded Wed Jul 15 15:43:28 MDT 2015
			// } else if (com.A_Pop&dfa.MTab.Machine[cur_st].Info.Action) != 0 && to == -1 { // 	xyzzy - total bullshit -- Only on terminal state!
			com.DbPrintf("match", "POP: cur_st=%d Rv=%d At: %s\n", cur_st, dfa.MTab.Machine[cur_st].Rv, com.LF())
			PopState()
		} else if (com.A_Reset & dfa.MTab.Machine[cur_st].Info.Action) != 0 {
			com.DbPrintf("match", "RESET: At: %s\n", com.LF())
			ResetState()
		}

		com.DbPrintf("match4", "Match, continuing to advance[ctx.St=%d] rn=%s", ctx.St, string(rn))
		com.DbPrintf("match4", " ii=%d to=%d\n", ii, to)
		if to == -1 {
			com.DbPrintf("match", "to=-1, At: %s\n", com.LF())
			ctx.St = 0
		}
		if rn == re.R_EOF {
			com.DbPrintf("match", " At: %s\n", com.LF())
			// AtEOF = true
			done = true
		}

	}
	if AtEOF {
		com.DbPrintf("match", " Reached EOF\n")
	} else {
		com.DbPrintf("match", " Early Exit!!!, no match!!!\n")
	}

	// Mods July 14
	com.DbPrintf("match", " send EOF \n")
	if lex.SendOnChanel {
		com.DbPrintf("match", " At: %s\n", com.LF())
		tt := LexieChanelType{Token: tok.Token{TokNo: aTok_EOF}}
		lex.Message <- tt
		//pjs July 3 - close(lex.Message)
	}
	com.DbPrintf("match", " At: %s\n", com.LF())
	com.DbPrintf("match", " end of function \n")
	// lex.Message <- tt

}
