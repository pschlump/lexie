//
// M A T C H - Part of Lexie Lexical Generation System
//
// Copyright (C) Philip Schlump, 2014-2025.
// Version: 1.0.8
//

package dfa

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"sort"
	"strings"

	"github.com/pschlump/dbgo"
	"github.com/pschlump/filelib"
	"github.com/pschlump/lexie/com"
	"github.com/pschlump/lexie/in"
	"github.com/pschlump/lexie/nfa"
	"github.com/pschlump/lexie/pbread"
	"github.com/pschlump/lexie/re"
	"github.com/pschlump/lexie/tok"
)

type LexieChanelType struct {
	Token tok.Token
}

type LexieStackType struct {
	St int // Current State in
}

type Lexie struct {
	IsCompiled bool // Is it currently compiled into a DFA
	NFA        []*nfa.NFA_PoolType
	DFA        []*DFA_PoolType
	Machines   []int

	NFA_Machine []*nfa.NFA_PoolType
	DFA_Machine []*DFA_PoolType

	TokList *tok.TokenList // ATokBuffer TokenBuffer // Output Token Stuff

	Im *in.ImType

	InputReader  *pbread.PBReadType
	StartMachine string

	// Channel to return data on
	SendOnChanel bool
	Message      chan LexieChanelType
}

// -----------------------------------------------------------------------------------------------------------------------------

// See: /Users/corwin/Projects/pongo2/lexie/note/t1.lex.go.old

// ---------------------------------------------------------------------------------------------------------------------------------------------------------------
func NewLexie() *Lexie { // Create a new matcher pool
	return &Lexie{
		IsCompiled:   false,
		SendOnChanel: false,
		Message:      make(chan LexieChanelType),
	}
}

// ---------------------------------------------------------------------------------------------------------------------------------------------------------------
func (lex *Lexie) SetChanelOnOff(flag bool) {
	lex.SendOnChanel = flag
}

// ---------------------------------------------------------------------------------------------------------------------------------------------------------------

func ConvertStringActionToFlag(aa string) (rv int) {
	sa := strings.Split(aa, "|")
	rv = 0
	for _, tt := range sa {
		if x, ok := com.ReservedActionValues[tt]; ok {
			rv |= x
		}
	}
	return
}

// ---------------------------------------------------------------------------------------------------------------------------------------------------------------

func KeyIntMapStringSort(in map[int]string) []int {
	var rv []int
	for ii, _ := range in {
		rv = append(rv, ii)
	}
	return KeyIntSort(rv)
}

// ---------------------------------------------------------------------------------------------------------------------------------------------------------------

func KeyIntSort(in []int) (rv []int) {
	rv = in
	sort.Sort(sort.IntSlice(rv))
	return
}

// ---------------------------------------------------------------------------------------------------------------------------------------------------------------

// lex .DumpTokenBuffer()
func (lex *Lexie) DumpTokenBuffer(fo io.Writer) {
	lex.TokList.FDumpTokenBuffer(fo, false)
}
func (lex *Lexie) DumpTokenBuffer2(fo io.Writer) {
	lex.TokList.FDumpTokenBuffer(fo, true)
}

// ---------------------------------------------------------------------------------------------------------------------------------------------------------------

func (lex *Lexie) PrintActionFlags(dfa *DFA_PoolType, fmt string) {
	switch fmt {
	case "text":

		dbgo.DbPrintf("match", "Action Flags Are:\n")
		// com. ConvertActionFlagToString(kk int) (rv string) {
		dn := make(map[int]bool)
		for _, vv := range dfa.Pool {
			if vv.IsUsed {
				if vv.Info.Action != 0 {
					if _, ok := dn[vv.Info.Action]; !ok {
						dbgo.DbPrintf("match", "    %2x: %s\n", vv.Info.Action, com.ConvertActionFlagToString(vv.Info.Action))
						dn[vv.Info.Action] = true
					}
				}
			}
		}

	case "go":

		dbgo.Printf("/*\nAction Flags Are:\n")
		// com. ConvertActionFlagToString(kk int) (rv string) {
		dn := make(map[int]bool)
		for _, vv := range dfa.Pool {
			if vv.IsUsed {
				if vv.Info.Action != 0 {
					if _, ok := dn[vv.Info.Action]; !ok {
						dbgo.Printf("    %2x: %s\n", vv.Info.Action, com.ConvertActionFlagToString(vv.Info.Action))
						dn[vv.Info.Action] = true
					}
				}
			}
		}
		dbgo.Printf("*/\n\n")

	default:
		dbgo.Fprintf(os.Stderr, "Invalid format at:%(LF) >%s< \n", fmt)
		os.Exit(1)
	}
}

// ---------------------------------------------------------------------------------------------------------------------------------------------------------------

func NewContext(InitState int, dfa *DFA_PoolType) (rv *MatchContextType) {
	return &MatchContextType{
		St:  InitState,
		Dfa: dfa,
	}
}

// ---------------------------------------------------------------------------------------------------------------------------------------------------------------
/*
func (lex *Lexie) InitGetToken(rrr *pbread.PBReadType, sm string) (AToken tok.Token) {
	lex.InputReader = rrr
	lex.StartMachine = sm
	// xyzzy
	return
}
*/

// ---------------------------------------------------------------------------------------------------------------------------------------------------------------
/*
func (lex *Lexie) GetToken() (AToken tok.Token) {
	// xyzzy
	return
}
*/

// FinializeMachines will put the final touches on each macine.
func (lex *Lexie) FinializeMachines() {
	for ii := range lex.DFA_Machine {
		lex.DFA_Machine[ii].MachineId = ii // Assign unique ID to each table.
		dfa := lex.DFA_Machine[ii]
		dt := dfa.ConvertToTable() // Take state machine graph and covnert into tables.
		dfa.MTab = &dt
		lex.DFA_Machine[ii].MachineId = ii // Assign unique ID to each table.
	}
}

// ---------------------------------------------------------------------------------------------------------------------------------------------------------------
func convRuleToActionFlag(ww *in.ImRuleType) int {
	rv := 0
	if ww.Repl {
		rv |= com.A_Repl
	}
	if len(ww.CallName) > 0 {
		rv |= com.A_Push
	}
	if ww.Return {
		rv |= com.A_Pop
	}
	if ww.Reset {
		rv |= com.A_Reset
	}
	if ww.PatternType == 2 {
		rv |= com.A_EOF
	}
	if ww.Err {
		rv |= com.A_Reset
		rv |= com.A_Error
	}
	if ww.Warn {
		rv |= com.A_Warning
	}
	if ww.NotGreedy {
		rv |= com.A_NotGreedy
	}
	return rv
}

// ---------------------------------------------------------------------------------------------------------------------------------------------------------------

func (lex *Lexie) NewReadFile(path, outTag string) (err error) {
	lex.Im = in.ImReadFile(path)

	lex.NFA_Machine = make([]*nfa.NFA_PoolType, 0, 100)
	lex.DFA_Machine = make([]*DFA_PoolType, 0, 100)

	// vv=in.ImDefinedValueType {Seq:1 WhoAmI:ReservedWords NameValueStr:map[and:Tok_L_AND not:Tok_not as:Tok_as in:Tok_in bor:Tok_B_OR band:Tok_B_AND xor:Tok_XOR or:Tok_L_OR true:Tok_true false:Tok_false export:Tok_export] NameValue:map[and:4 true:32 as:34 bor:42 band:41 xor:64 or:5 false:33 not:31 export:35 in:28] Reverse:map[5:or 32:true 42:bor 31:not 41:band 35:export 33:false 28:in 64:xor 4:and 34:as] SeenAt:map[bor:{LineNo:[39] FileName:[unk-file]} band:{LineNo:[39] FileName:[unk-file]} and:{LineNo:[39] FileName:[unk-file]} true:{LineNo:[39] FileName:[unk-file]} export:{LineNo:[39] FileName:[unk-file]} in:{LineNo:[39] FileName:[unk-file]} as:{LineNo:[39] FileName:[unk-file]} or:{LineNo:[39] FileName:[unk-file]} false:{LineNo:[39] FileName:[unk-file]} not:{LineNo:[39 39] FileName:[unk-file unk-file]} xor:{LineNo:[39] FileName:[unk-file]}]}, File: /Users/corwin/Projects/pongo2/lexie/dfa/match.go LineNo:260
	for _, vv := range lex.Im.Def.DefsAre {
		// ["ReservedWords"] {
		// func (st *SymbolTable) DefineReservedWord(name string, fxid int) (ss *SymbolType) {
		dbgo.DbPrintf("dfa5", "vv=%T %+v, %s\n", vv, vv, dbgo.LF())
	}

	for ii, vv := range lex.Im.Machine {
		nm := vv.Name

		Nfa := nfa.NewNFA_Pool()
		Cur := Nfa.GetNFA()
		Nfa.InitState = Cur
		for jj, ww := range vv.Rules {
			rVx := ww.Rv
			if ww.ReservedWord {
				dbgo.DbPrintf("dfa5", "This rule rv=%d is a reserved word rule, AAbbCC\n", rVx)
			}
			ww_A := convRuleToActionFlag(ww)
			if ww.Repl {
				rVx = 9900 // 9900 is replace
				dbgo.DbPrintf("match", "###################################### ww.Replace: ii=%d jj=%d ->%s<-, %s\n", ii, jj, ww.ReplString, dbgo.LF())
			}
			cur := -1
			if ww.PatternType == 2 {
				dbgo.DbPrintf("db_Matcher_02", "ADDING AT %2d RE: %-30s (Rv:%2d, final=%4d), %s\n", jj, "<M_EOF>", ww.Rv, rVx, dbgo.LF())
				cur = Nfa.AddReInfo(re.X_EOF, "", jj+1, rVx, nfa.InfoType{Action: ww_A, NextState: ww.Call})
			} else {
				dbgo.DbPrintf("db_Matcher_02", "M= ->%s<- Adding at %2d RE: %-30s (Rv:%2d, final=%4d), %s\n", ww.Pattern, jj, ww.Pattern, ww.Rv, rVx, dbgo.LF())
				cur = Nfa.AddReInfo(ww.Pattern, "", jj+1, rVx, nfa.InfoType{Action: ww_A, NextState: ww.Call, ReplStr: ww.ReplString})
			}
			if ww.ReservedWord {
				Nfa.SetReservedWord(cur)
			}
		}

		dbgo.DbPrintf("match", "BuildDFA_2: Nfa.Sigma Before Finialize->%s<-\n", Nfa.Sigma)
		if dbgo.IsDbOn("db_Matcher_02") {
			dbgo.DbPrintf("match", "NFA for (Before Finialize) ->%s<-\n", nm)
			Nfa.DumpPool(false)
		}

		Nfa.FinializeNFA()

		dbgo.DbPrintf("match", "BuildDFA_2: Nfa.Sigma ->%s<-\n", Nfa.Sigma)
		if dbgo.IsDbOn("db_Matcher_02") {
			dbgo.DbPrintf("match", "Final NFA for ->%s<-\n", nm)
			Nfa.DumpPool(false)
		}
		lex.NFA_Machine = append(lex.NFA_Machine, Nfa)

		Dfa := NewDFA_Pool()
		Dfa.ConvertNFAToDFA(Nfa)
		if dbgo.IsDbOn("db_Matcher_02") {
			dbgo.DbPrintf("match", "Final DFA for ->%s<-\n", nm)
			Dfa.DumpPool(false)
		}
		lex.DFA_Machine = append(lex.DFA_Machine, Dfa)

		if dbgo.IsDbOn("db_Matcher_02") {

			// lex.Im.InputLines // has the entire input.

			last := len(lex.DFA_Machine) - 1

			newFile := fmt.Sprintf("../ref/%s_%s_%d.tst", outTag, "machine", last)
			gvFile := fmt.Sprintf("../ref/%s_%s_%d.gv", outTag, "machine", last)
			svgFile := fmt.Sprintf("../ref/%s_%s_%d.svg", outTag, "machine", last)

			fp, err := filelib.Fopen(newFile, "w")
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error opening otput file %s, error: %s\n", newFile, err)
				return err
			}
			fmt.Fprintf(fp, "Machine:\n")
			for _, ss := range lex.Im.InputLines {
				fmt.Fprintf(fp, "\t%s\n", ss)
			}
			lex.DFA_Machine[last].DumpPoolJSON(fp, fmt.Sprintf("Lex-Machine-%d", last), 1)
			fp.Close()

			gv, err := filelib.Fopen(gvFile, "w")
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error opening otput file %s, error: %s\n", gvFile, err)
				return err
			}
			lex.DFA_Machine[last].GenerateGVFile(gv, fmt.Sprintf("Lex-Machine-%d", last), 1)
			gv.Close()

			// out, err := exec.Command("/usr/local/bin/dot", "-Tsvg", "-o"+svgFile, gvFile).Output()
			out, err := exec.Command(com.DotPath, "-Tsvg", "-o"+svgFile, gvFile).Output()
			if err != nil {
				dbgo.DbPrintf("match", "%(red)Error%(reset) from dot, %s, %(LF)\n", err)
				dbgo.DbPrintf("match", "Output: %s\n", out)
			}
		}
	}
	return nil
}

/* vim: set noai ts=4 sw=4: */
