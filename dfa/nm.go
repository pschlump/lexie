//
// M A T C H - Part of Lexie Lexical Generation System
//
// Copyright (C) Philip Schlump, 2014-2025.
//

package dfa

import (
	"fmt"
	"os"

	"github.com/pschlump/dbgo"
	"github.com/pschlump/lexie/com"
	"github.com/pschlump/lexie/pbread"
	"github.com/pschlump/lexie/re"
	"github.com/pschlump/lexie/tok"
)

// MatchContextType maintain the current state of the machine.
type MatchContextType struct {
	St  int
	Dfa *DFA_PoolType
}

func (lex *Lexie) AssignMacineId(rrr *pbread.PBReadType, s_init string) {
	for ii := range lex.DFA_Machine {
		// dfa := lex.DFA_Machine[ii]
		// dfa.MachineId = ii
		lex.DFA_Machine[ii].MachineId = ii
	}
}

// MatcherLexieTable will use a push-back reader, `rrr`, a lexie machine, `lex`, and the name of a machine
// to match the input data and convert it into a stream of tokens.
func (lex *Lexie) MatcherLexieTable(rrr *pbread.PBReadType, s_init string) {

	var dfa *DFA_PoolType
	var ii, to, pos_no, col_no, line_no int
	var rn rune
	var SMatch, filename string

	lex.FinializeMachines()
	lex.AssignMacineId(rrr, s_init)

	init, err := lex.Im.LookupMachine(s_init)
	if err != nil {
		dbgo.Fprintf(os.Stderr, "Error: invalid machine name >%s< - not found, %s\n", s_init, err)
		return
	}
	if dbgo.IsDbOn("output-machine") {
		for ii := range lex.DFA_Machine {
			dbgo.Printf("\n\n%(blue)Machine[%d] =%(reset)\n", ii)
			dbgo.Printf("%(blue)~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~%(reset)\n", ii)
			dfa = lex.DFA_Machine[ii]
			dfa.PrintStateMachine(os.Stdout, "text")
			lex.OutputActionFlags(dfa)
		}
	}

	dfa = lex.DFA_Machine[init]

	// -----------------------------------------------------------------------------------------------------------
	// Machines are actually complete at this point - should be able to generate code.
	// 1.  xyzzy2145 - generate code.
	// -----------------------------------------------------------------------------------------------------------

	ctx := NewContext(dfa.MTab.InitState, nil)
	ctx_stack := make([]*MatchContextType, 0, 100)

	lex.TokList = tok.NewTokenList()
	ignoreToken, err := lex.Im.Lookup("Tokens", "Tok_Ignore")
	if err != nil {
		dbgo.Fprintf(os.Stderr, "Error: Missing required token Tok_Ingnore, %s\n", err)
		return
	}
	aTok_EOF, err := lex.Im.Lookup("Tokens", "Tok_EOF")
	if err != nil {
		dbgo.Fprintf(os.Stderr, "Error: Missing required token aTok_EOF, %s\n", err)
		return
	}
	lex.TokList.IgnoreToken = ignoreToken
	line_no = 1
	col_no = 1
	pos_no = 1
	filename = ""
	SMatch = ""
	AtEOF := false
	done := false
	TokStart := 0
	p_line_no := 1
	p_col_no := 1

	Next := func() (rn rune) {
		rn, done := rrr.NextRune()
		if done {
			AtEOF = true
			rn = re.R_EOF
			return
		}
		SMatch += string(rn)
		pos_no++
		return
	}

	Accept := func() {
		SMatch += string(rn)
	}
	_ = Accept

	Peek := func() (rn rune) {
		rn, done = rrr.PeekRune()
		if done {
			AtEOF = true
			rn = re.R_EOF
		}
		return
	}
	_ = Peek

	PeekPeek := func() (rn rune) {
		rn, done = rrr.PeekPeekRune()
		if done {
			rn = 0
		}
		return
	}
	_ = PeekPeek

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

		// This is the spot to convert from Tok_ID && ReservedWord -> a new ID
		// 		lex.Im.St.Lookup.LookupSymbol(SMatch)
		// 		func (st *SymbolTable) LookupSymbol(name string) (as *SymbolType, err error) {
		// 		vv=in.ImDefinedValueType {Seq:1 WhoAmI:ReservedWords NameValueStr:map[and:Tok_L_AND not:Tok_not as:Tok_as in:Tok_in bor:Tok_B_OR band:Tok_B_AND xor:Tok_XOR or:Tok_L_OR true:Tok_true false:Tok_false export:Tok_export] NameValue:map[and:4 true:32 as:34 bor:42 band:41 xor:64 or:5 false:33 not:31 export:35 in:28] Reverse:map[5:or 32:true 42:bor 31:not 41:band 35:export 33:false 28:in 64:xor 4:and 34:as] SeenAt:map[bor:{LineNo:[39] FileName:[unk-file]} band:{LineNo:[39] FileName:[unk-file]} and:{LineNo:[39] FileName:[unk-file]} true:{LineNo:[39] FileName:[unk-file]} export:{LineNo:[39] FileName:[unk-file]} in:{LineNo:[39] FileName:[unk-file]} as:{LineNo:[39] FileName:[unk-file]} or:{LineNo:[39] FileName:[unk-file]} false:{LineNo:[39] FileName:[unk-file]} not:{LineNo:[39 39] FileName:[unk-file unk-file]} xor:{LineNo:[39] FileName:[unk-file]}]}, File: /Users/corwin/Projects/pongo2/lexie/dfa/match.go LineNo:260
		if dfa.MTab.Machine[ctx.St].Info.ReservedWord {
			dbgo.DbPrintf("rw-lookup", "\n\nFound a  reserved word ------------------------------------------------------------------------ %(LF)\n")
			dbgo.DbPrintf("rw-lookup", "SMatch >%s<, %s\n", SMatch, dbgo.LF())

			vv, ok := lex.Im.Def.DefsAre["ReservedWords"].NameValue[SMatch]
			if ok {
				dbgo.DbPrintf("rw-lookup", "Substituting(2) -- Found it as %d !!!, %(LF)\n", vv)
				lRv = vv
			}

			dbgo.DbPrintf("rw-lookup", "\n\n")
		}

		dbgo.DbPrintf("match", "+==================================================================================\n")
		dbgo.DbPrintf("match", "+ SaveToken: Line:%d Col:%d\n", line_no, col_no)
		dbgo.DbPrintf("match", "+==================================================================================\n")
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

	// xyzzy3330 - StashToken
	// xyzzy3330 - HaveStashToken -> bool
	// xyzzy3330 - ReturnStashToken -- Will call FlushToken if HaveStashToken is true

	FlushToken := func(isHard bool) {
		if dfa.MTab.Machine[ctx.St].Rv == 9900 && (com.A_Repl&dfa.MTab.Machine[ctx.St].Info.Action) != 0 {
			dbgo.DbPrintf("match", " Doing Replace on Token, Len %d to ->%s<-\n ", dfa.MTab.Machine[ctx.St].Info.MatchLength, dfa.MTab.Machine[ctx.St].Info.ReplStr)
			lex.TokList.ReplaceToken(dfa.MTab.Machine[ctx.St].Info.MatchLength, dfa.MTab.Machine[ctx.St].Info.ReplStr)
		}
		dbgo.DbPrintf("match4", "FlushTokenStareted ------------------------------------------------------------------------------\n")
		beforeFlush := len(lex.TokList.TokenData)
		lex.TokList.FlushTokenBuffer(TokStart, isHard, AtEOF)
		afterFlush := len(lex.TokList.TokenData)

		if lex.SendOnChanel {
			for ll := beforeFlush; ll < afterFlush; ll++ {
				fmt.Printf("At: %s\n", dbgo.LF())

				tt := LexieChanelType{Token: lex.TokList.TokenData[ll]}
				lex.Message <- tt
			}
		}
		TokStart = 0
		SMatch = ""
		dbgo.DbPrintf("match4", "FlushTokenEnded ------------------------------------------------------------------------------\n")
	}

	DumpStack := func() {
		dbgo.DbPrintf("match4", "\tDump ctx_stack len=[%d]\n", len(ctx_stack))
		for ii := range ctx_stack {
			fmt.Printf("At: %s\n", dbgo.LF())
			dbgo.DbPrintf("match4", "\t\tDump ctx_stack[%d] = machine %d\n", ii, ctx_stack[ii].Dfa.MachineId)
		}
	}

	PushState := func(St int) {
		dbgo.DbPrintf("match", "\n-------------------------------------------------------------------------------\n")
		dbgo.DbPrintf("match", "Found a PUSH: to %d -- This should be subscript of new machine\n", dfa.MTab.Machine[St].Info.NextState)
		dbgo.DbPrintf("match", "-------------------------------------------------------------------------------\n")
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
			dbgo.DbPrintf("match", "\n-------------------------------------------------------------------------------\n")
			dbgo.DbPrintf("match", "Found a POP: State Poping like Corn, Machine now %d\n", dfa.MachineId)
			dbgo.DbPrintf("match", "-------------------------------------------------------------------------------\n")
			DumpStack()
		} else {
			dbgo.DbPrintf("match", "Error: Attempt to pop when stack is empty\n")
		}
		ctx.St = 0
	}

	ResetState := func() {
		dbgo.DbPrintf("match", "Found a RESET\n")
		ctx_stack = ctx_stack[:0]
		dfa = lex.DFA_Machine[init]
		ctx.St = 0
	}

	// xyzzy3330
	// If "greedy" in this state then need to look ahead, on a terminal state if there are future states
	// and if the next character will allow you to move to a future state then save (stash) the current
	// result and position and continue running the machine.  A future longer match will result in replacing
	// the current match.  If it is "greedy" then repeat the process.  If not "greedy" then return.
	// If you reach a fail-to-match ( current rune has no future state ), then return the saved state
	// and re-position to the point at which the saved state had saved it.

	for !done {
		dbgo.DbPrintf("match", "\n%(blue)**********************************************************************************************************%(reset)\n")
		dbgo.DbPrintf("match", "Top: (machine number %d) current state ctx.St:%d\n", dfa.MachineId, ctx.St)
		cur_st := ctx.St

		rn = Peek() // grab next token for a Peek()
		// rn = Next()                                          // grab next token
		ii = dfa.MTab.SMap.MapRune(rn)                       // Convet rune, 'rn', to a table state position
		to = dfa.MTab.Machine[cur_st].To[ii]                 // if -1, then no future state, if >= 0, then a potential future state
		hardMatch := dfa.MTab.Machine[cur_st].Info.HardMatch // if true, then this is a terminal state (match occurd)
		dbgo.DbPrintf("match", " (peek) machine=%d cur_st=%d ii=%d for rune rn=->%s<-  %(yellow)to=%d hard=%v%(reset)  Rv=%d Tau=%v Action=%x, ctx_stack=%d\n",
			dfa.MachineId, cur_st, ii, string(rn), to, hardMatch, dfa.MTab.Machine[cur_st].Rv, dfa.MTab.Machine[cur_st].Tau,
			dfa.MTab.Machine[cur_st].Info.Action, ctx_stack)

		if rn == re.R_EOF {
			dbgo.DbPrintf("match", " %(yellow)At: %(LF), matched re.R_EOF\n")
			FlushToken(true)
		} else if dfa.MTab.Machine[cur_st].Rv > 0 && dfa.MTab.Machine[cur_st].Tau && to == -1 {
			dbgo.DbPrintf("match", " %(yellow)At: %(LF) -- this one is imporant\n")
			SaveToken()
			FlushToken(hardMatch)
		} else if dfa.MTab.Machine[cur_st].Rv > 0 && hardMatch {
			if to == -1 {
				// no next state, so done at this point.
				dbgo.DbPrintf("match", " %(yellow)At: %(LF)\n")
				SaveToken()
				FlushToken(hardMatch)
			} else {
				// xyzzy3330 - I suspect that this is the problem spot, a to >= 0
				// xyzzy3330 - hard==true => terminal state, to >= 0 ==>> the problem we are seeing
				dbgo.DbPrintf("match", " %(cyan)At: %(LF), to=%d, indicates a potential next state\n", to)
				if hardMatch && to != -1 {
					dbgo.DbPrintf("match", " %(red)At: %(LF), hardMatch=%v to=%d, *** should push state ***\n", hardMatch, to)
				}
				// (peek) machine=0 cur_st=3 ii=5 for rune rn=-> <- to=0 hard=true Rv=8 Tau=true
				if (com.A_Push & dfa.MTab.Machine[cur_st].Info.Action) == 0 {
					dbgo.DbPrintf("match", " %(cyan)At: %s\n", dbgo.LF())
					rn = Next() // Remove makes beginning work
				}
				dbgo.DbPrintf("match", " %(yellow)At: %(LF)\n")
				SaveToken()
				FlushToken(hardMatch) // PJS added Wed Jul 15 15:43:28 MDT 2015
			}
			ctx.St = to
		} else if dfa.MTab.Machine[cur_st].Rv > 0 { // && !dfa.MTab.Machine[cur_st].Tau && to != -1 {
			dbgo.DbPrintf("match", " %(yellow)At: %(LF)\n")
			rn = Next()
			dbgo.DbPrintf("match", " %(yellow)At: %(LF)\n")
			SaveToken()
			ctx.St = to
		} else if ((com.A_Push | com.A_Pop | com.A_Reset) & dfa.MTab.Machine[cur_st].Info.Action) != 0 {
			dbgo.DbPrintf("match", "%(yellow)-- critical -- At: %(LF)\n")
			// (peek) machine=2 cur_st=21 ii=24 for rune rn=->{<- to=36 hard=false Rv=0 Tau=false
			// (peek) machine=2 cur_st=24 ii=4 for rune rn=-> <- to=-1 hard=true Rv=0 Tau=false
			if to >= 0 {
				rn = Next()
			}
			ctx.St = to
		} else {
			dbgo.DbPrintf("match", "%(yellow)-- other --  At: %(LF)\n")
			if to >= 0 {
				rn = Next()
			}
			ctx.St = to
		}

		if (com.A_Push & dfa.MTab.Machine[cur_st].Info.Action) != 0 { // 					xyzzy - total bullshit -- Allows to hapen on non-termaial in middle of stuff -- call non-terminal
			dbgo.DbPrintf("match", "PUSH: At: %(LF)\n")
			PushState(cur_st)
		} else if (com.A_Pop & dfa.MTab.Machine[cur_st].Info.Action) != 0 { // PJS modded Wed Jul 15 15:43:28 MDT 2015
			// } else if (com.A_Pop&dfa.MTab.Machine[cur_st].Info.Action) != 0 && to == -1 { // 	xyzzy - total bullshit -- Only on terminal state!
			dbgo.DbPrintf("match", "POP: cur_st=%d Rv=%d At: %(LF)\n", cur_st, dfa.MTab.Machine[cur_st].Rv)
			PopState()
		} else if (com.A_Reset & dfa.MTab.Machine[cur_st].Info.Action) != 0 {
			dbgo.DbPrintf("match", "RESET: At: %(LF)\n")
			ResetState()
		}

		dbgo.DbPrintf("match4", "Match, continuing to advance[ctx.St=%d] rn=%s", ctx.St, string(rn))
		dbgo.DbPrintf("match4", " ii=%d to=%d\n", ii, to)
		if to == -1 {
			dbgo.DbPrintf("match", "to=-1, At: %s\n", dbgo.LF())
			ctx.St = 0
		}
		if rn == re.R_EOF {
			dbgo.DbPrintf("match", " At: %s\n", dbgo.LF())
			// AtEOF = true
			done = true
		}

	}
	if AtEOF {
		dbgo.DbPrintf("match", " Reached EOF\n")
	} else {
		dbgo.DbPrintf("match", " Early Exit!!!, no match!!!\n")
	}

	// Mods July 14
	dbgo.DbPrintf("match", " send EOF \n")
	if lex.SendOnChanel {
		dbgo.DbPrintf("match", " At: %s\n", dbgo.LF())
		tt := LexieChanelType{Token: tok.Token{TokNo: aTok_EOF}}
		lex.Message <- tt
		//pjs July 3 - close(lex.Message)
	}
	dbgo.DbPrintf("match", " At: %s\n", dbgo.LF())
	dbgo.DbPrintf("match", " end of function \n")
	// lex.Message <- tt

}

const db1 = false

/*
./nm.go:60:17: assignment mismatch: 1 variable but lex.Im.Lookup returns 2 values
./nm.go:65:14: assignment mismatch: 1 variable but lex.Im.Lookup returns 2 values
*/
