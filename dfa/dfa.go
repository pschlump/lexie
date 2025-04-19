//
// D F A - Part of Lexie Lexical Generation System
//
// Copyright (C) Philip Schlump, 2014-2025.
//
//
// DFA - Deterministic Finite Automata.
//

package dfa

import (
	"fmt"
	"io"
	"os"
	"unicode"
	"unicode/utf8"

	"github.com/pschlump/dbgo"
	"github.com/pschlump/lexie/com"
	"github.com/pschlump/lexie/in"
	"github.com/pschlump/lexie/nfa"
	"github.com/pschlump/lexie/re"
	"github.com/pschlump/lexie/smap"
	"github.com/pschlump/lexie/tok"
)

type DFA_Type struct {
	Next2      []nfa.TransitionType //
	Rv         int                  // 0 indicates not assigned, non-terminal
	Is0Ch      bool                 //	Tau
	Info       nfa.InfoType         //
	TRuleMatch int                  // may be a non-terminal that you want to know matched. -- A set of these is returned on matches or can be retrieved on fails too.
	NextFree   int                  //		For free list
	IsUsed     bool                 //		For Free list
	A_IAm      int                  //		Debug Usage
	LineNo     string               // LineNo where added
	StateName  string               // Used in NFA -> DFA
	StateSet   []int                //
	Visited    bool                 //
}

type DFA_PoolType struct {
	Pool      []DFA_Type     //
	Cur       int            //
	Top       int            //
	NextFree  int            //
	InitState int            //
	Sigma     string         //
	MTab      *dfaTable      //
	TokList   *tok.TokenList // ATokBuffer TokenBuffer // Output Token Stuff
	MachineId int
}

const InitDFASize = 3

// -----------------------------------------------------------------------------------------------------------------------------------------------
// Create a new DFA pool
func NewDFA_Pool() *DFA_PoolType {
	return &DFA_PoolType{
		Pool:     make([]DFA_Type, InitDFASize, InitDFASize),
		Cur:      0,
		Top:      InitDFASize,
		NextFree: -1,
	}
}

// -----------------------------------------------------------------------------------------------------------------------------------------------
func (dfa *DFA_PoolType) DumpTokenBuffer() {
	dfa.TokList.DumpTokenBuffer()
}

// -----------------------------------------------------------------------------------------------------------------------------------------------
// Allocate an DFA tree node
func (dfa *DFA_PoolType) GetDFA() int {
	//fmt.Printf("at %s\n", dbgo.LF())
	tmp := 0
	if dfa.Cur < dfa.Top && !dfa.Pool[dfa.Cur].IsUsed {
		//fmt.Printf("at %s\n", dbgo.LF())
		tmp = dfa.Cur
		dfa.Cur++
	} else if dfa.Cur >= dfa.Top || dfa.NextFree == -1 {
		//fmt.Printf("at %s, dfa.Cur=%d dfa.Top=%d dfa.NextFree=%d\n", dbgo.LF(), dfa.Cur, dfa.Top, dfa.NextFree)
		dfa.Top = 2 * dfa.Top
		newPool := make([]DFA_Type, dfa.Top, dfa.Top) // extend array
		copy(newPool, dfa.Pool)
		dfa.Pool = newPool
		tmp = dfa.Cur
		dfa.Cur++
	} else {
		//fmt.Printf("at %s\n", dbgo.LF())
		tmp = dfa.NextFree
		dfa.NextFree = dfa.Pool[tmp].NextFree
	}
	dfa.Pool[tmp].NextFree = -1
	dfa.Pool[tmp].Rv = 0
	dfa.Pool[tmp].Next2 = dfa.Pool[tmp].Next2[:0]
	dfa.Pool[tmp].IsUsed = true
	dfa.Pool[tmp].A_IAm = tmp
	dfa.Pool[tmp].LineNo = dbgo.LINE(2)
	return tmp
}

// -----------------------------------------------------------------------------------------------------------------------------------------------
// Free an DFA tree node
func (dfa *DFA_PoolType) FreeDFA(ii int) {
	dfa.Pool[ii].IsUsed = false
	dfa.Pool[ii].NextFree = dfa.NextFree
	dfa.Pool[ii].Rv = 0
	dfa.Pool[ii].Next2 = dfa.Pool[ii].Next2[:0]
	dfa.NextFree = ii
}

// -----------------------------------------------------------------------------------------------------------------------------------------------
// Return the start state number
func (dfa *DFA_PoolType) Pos0Start() int {
	return dfa.InitState
}

// -----------------------------------------------------------------------------------------------------------------------------------------------
func (dfa *DFA_PoolType) DiscardPool() {
	dfa.Pool = make([]DFA_Type, InitDFASize, InitDFASize)
	dfa.Cur = 0
	dfa.Top = InitDFASize
	dfa.NextFree = -1
}

// ---------------------------------------------------------------------------------------------------------------------------------------------------------------
//  3. Input Parsing Issues/Errors/Warnings
//  1. Can not have A_Pop and A_Push at same time - Check for this.
//  2. A_Pop must be a hard match -
//  3. A_Pop must be a terminal! - Return() must have a Rv()
//  4. A_Push, A_Pop, A_Reset can not be ambiguous tokens, can not {% POP and {%= Return a value, won't work ( At least not yet )
func (dfa *DFA_PoolType) VerifyMachine() {
	for ii, vv := range dfa.Pool {
		if vv.IsUsed {
			if vv.Info.Action != 0 {
				if (vv.Info.Action&com.A_Pop) != 0 && (vv.Info.Action&com.A_Push) != 0 {
					com.StashError(fmt.Sprintf("Error: State[%d] has both a Push/Call and a Pop/Return optration at the same time.\n", ii))
				}
				if (vv.Info.Action&com.A_Pop) != 0 && (vv.Info.Action&com.A_Reset) != 0 {
					com.StashError(fmt.Sprintf("Error: State[%d] has both a Push/Call and a Pop/Return optration at the same time.\n", ii))
				}
				if (vv.Info.Action&com.A_Push) != 0 && (vv.Info.Action&com.A_Reset) != 0 {
					com.StashError(fmt.Sprintf("Error: State[%d] has both a Push/Call and a Pop/Return optration at the same time.\n", ii))
				}
				if ((vv.Info.Action&com.A_Pop) != 0 || (vv.Info.Action&com.A_Push) != 0 || (vv.Info.Action&com.A_Reset) != 0) && vv.Rv == 0 {
					com.StashError(fmt.Sprintf("Error: State[%d] a Push/Call or a Pop/Return/Reset optration Must be a terminal state with a Rv()\n", ii))
				}
				if ((vv.Info.Action&com.A_Pop) != 0 || (vv.Info.Action&com.A_Push) != 0 || (vv.Info.Action&com.A_Reset) != 0) && !vv.Info.HardMatch { // xyzzy8
					com.StashError(fmt.Sprintf("Info: State[%d] a Push/Call or a Pop/Return/Reset optration Must be a terminal state with fixed string matched, 'a*' or 'a?' is not a legitimate match.\n", ii))
				}

			}
		}
	}
}

// -----------------------------------------------------------------------------------------------------------------------------------------------
/*
type InfoType struct {
	Action      int
	MatchLength int
}

func (lex *Lexie) ConvertMachineNameToText(ee int) string {
	return lex.UseData.LexieData.MatchNames[ee]
}
*/

//func DumpInfo(info InfoType) string {
//	// xyzzy22 - convert info.NextState - > state name
//	sa := ConvertActionFlagToString(info.Action)
//	// nsname := lex.ConvertMachineNameToText(info.NextState)
//	return fmt.Sprintf("Action: %s(%02x), Ns:%d, MatchLength:%d\n", sa, info.Action, info.NextState, info.MatchLength)
//}

func (dfa *DFA_PoolType) DumpPool(all bool) {
	if all {
		dbgo.DbPrintf("db_DumpDFAPool", "Cur: %d Top: %d NextFree %d\n", dfa.Cur, dfa.Top, dfa.NextFree)
	}
	dbgo.DbPrintf("db_DumpDFAPool", "\n---------------------------- DFA Output -----------------------------------------------\n")
	dbgo.DbPrintf("db_DumpDFAPool", "\nDFA InitState: %d, Sigma ->%s<-\n\n", dfa.InitState, dfa.Sigma)
	pLnNo := dbgo.IsDbOn("db_DFA_LnNo")
	IfLnNo := func(s string) string {
		if pLnNo {
			t := fmt.Sprintf("[%3s]", s)
			return t
		}
		return ""
	}
	dbgo.DbPrintf("db_DumpDFAPool", "%3s%s: ", "St", IfLnNo("/Ln"))
	dbgo.DbPrintf("db_DumpDFAPool", " %12s %12s \u2714              \tEdges", "StateName", "StateSet")
	dbgo.DbPrintf("db_DumpDFAPool", "\n\n")
	for ii, vv := range dfa.Pool {
		if all || vv.IsUsed {
			dbgo.DbPrintf("db_DumpDFAPool", "%3d%s: ", ii, IfLnNo(vv.LineNo))
			dbgo.DbPrintf("db_DumpDFAPool", " %12s %12s %s :", vv.StateName, dbgo.SVar(vv.StateSet), com.ChkOrBlank(vv.Visited))
			if vv.Rv > 0 {
				if vv.Is0Ch {
					dbgo.DbPrintf("db_DumpDFAPool", " \u03c4:Tau:%04d ", vv.Rv)
				} else {
					dbgo.DbPrintf("db_DumpDFAPool", " T:%04d ", vv.Rv)
				}
			} else {
				dbgo.DbPrintf("db_DumpDFAPool", "        ")
			}
			if dbgo.IsDbOn("db_DumpDFAPool") {
				fmt.Printf("\t E:")
				for _, ww := range vv.Next2 {
					if ww.Is0ChMatch {
						fmt.Printf("//Found  \u03c4 (%s) //", dbgo.LF()) // Show a Tau(t) for a lambda that matchiens on else conditions.
					}
					if ww.IsLambda {
						fmt.Printf("{  ERROR!! \u03bb  %2d -> %2d  %s}  ", ww.From, ww.To, ww.LineNo)
					} else {
						// fmt.Printf("{ \"%s\" %2d -> %2d  %s}  ", ww.On, ww.From, ww.To, IfLnNo(ww.LineNo))
						on, _ := utf8.DecodeRune([]byte(ww.On))
						son := fmt.Sprintf("%q", ww.On)
						switch on {
						case re.R_DOT: // = '\uF8FA' // Any char in Sigma
							son = "DOT/uF8FA"
						case re.R_BOL: // = '\uF8F3' // Beginning of line
							son = "BOL/uF8F3"
						case re.R_EOL: // = '\uF8F4' // End of line
							son = "EOL/uF8F4"
						case re.R_NUMERIC: // = '\uF8F5'
							son = "NUMERIC/uF8F5"
						case re.R_LOWER: // = '\uF8F6'
							son = "LOWER/uF8F6"
						case re.R_UPPER: // = '\uF8F7'
							son = "UPPER/uF8F7"
						case re.R_ALPHA: // = '\uF8F8'
							son = "ALPHA/uF8F8"
						case re.R_ALPHNUM: // = '\uF8F9'
							son = "ALPHANUM/uF8F9"
						case re.R_EOF: // = '\uF8FB'
							son = "EOF/uF8FB"
						case re.R_not_CH: // = '\uF8FC' // On input lookup if the char is NOT in Signa then it is returned as this.
							son = "else_CH/uF8Fc"
						case re.R_N_CCL: // = '\uF8FD' // If char is not matched in this state then take this path
							son = "N_CCL/uF8Fd"
						case re.R_LAMBDA_MATCH: // = '\uF8FE'
							son = "LambdaM/uF8FE"
						}
						fmt.Printf("{ %s  %2d -> %2d  %s}  ", son, ww.From, ww.To, IfLnNo(ww.LineNo))
					}
				}
				fmt.Printf("\n")
				if vv.Info.Action != 0 || vv.Info.MatchLength != 0 {
					// fmt.Printf("\t\t\tInfo: %s\n", dbgo.SVar(vv.Info))		// xyzzy - output Info
					// xyzzy - NextState info
					fmt.Printf("\t\t\tDFA.Info: %s", nfa.DumpInfo(vv.Info))
					// if ((vv.Info.Action&com.A_Pop) != 0 || (vv.Info.Action&com.A_Push) != 0 || (vv.Info.Action&com.A_Reset) != 0) && !vv.Info.HardMatch {		// xyzzy8
					fmt.Printf(" IsHard=%v (((false imples else case Rv!)))\n", vv.Info.HardMatch)
				}
				fmt.Printf("\n")
			}
		}
	}
}

// -----------------------------------------------------------------------------------------------------------------------------------------------
func (dfa *DFA_PoolType) DumpPoolOneState(ii int) string {
	pLnNo := dbgo.IsDbOn("db_DFA_LnNo")
	IfLnNo := func(s string) string {
		if pLnNo {
			t := fmt.Sprintf("[%3s]", s)
			return t
		}
		return ""
	}
	vv := dfa.Pool[ii]
	s := ""
	s += fmt.Sprintf("%3d%s: ", ii, IfLnNo(vv.LineNo))
	s += fmt.Sprintf(" %12s %12s %s :", vv.StateName, dbgo.SVar(vv.StateSet), com.ChkOrBlank(vv.Visited))
	if vv.Rv > 0 {
		s += fmt.Sprintf(" T:%04d ", vv.Rv)
	} else {
		s += fmt.Sprintf("        ")
	}
	s += fmt.Sprintf("\t E:")
	for _, ww := range vv.Next2 {
		if ww.IsLambda {
			s += fmt.Sprintf("{  ERROR!! \u03bb  %2d -> %2d  %s}  ", ww.From, ww.To, ww.LineNo)
		} else {
			s += fmt.Sprintf("{ \"%s\" %2d -> %2d  %s}  ", ww.On, ww.From, ww.To, IfLnNo(ww.LineNo))
		}
	}
	s += "\n"
	if vv.Info.Action != 0 || vv.Info.MatchLength != 0 {
		s += fmt.Sprintf("\t\t\tDFA.Info: %s", nfa.DumpInfo(vv.Info))
	}
	s += "\n"
	return s
}

// -----------------------------------------------------------------------------------------------------------------------------------------------
func (dfa *DFA_PoolType) HaveStateAlready(inputSet []int) (loc int) {
	// fmt.Printf("HaveStateAlready: >>>>>>>>>>>>>>>>>>>>>>> 0 inputSet: %s\n", dbgo.SVar(inputSet))
	inputSet = com.USortIntSlice(inputSet) // Make set unique
	// fmt.Printf("HaveStateAlready: >>>>>>>>>>>>>>>>>>>>>>> 1 inputSet: %s\n", dbgo.SVar(inputSet))
	s := com.NameOf(inputSet)
	// fmt.Printf("HaveStateAlready: >>>>>>>>>>>>>>>>>>>>>>> Name searcing for is: %s\n", s)
	for ii, vv := range dfa.Pool {
		//		if ii >= dfa.Cur || ii >= dfa.Top {
		//			fmt.Printf("   Reached break\n")
		//			break
		//		}
		if vv.IsUsed {
			// fmt.Printf(" ****** Checking for match between dfa.Pool[%d].StateName ->%s<- and ->%s<-\n", ii, vv.StateName, s)
			if vv.StateName == s {
				// fmt.Printf("   Match found at %d\n", ii)
				return ii
			}
		}
	}
	// fmt.Printf("    no Match found, returing -1\n")
	return -1
}

// -----------------------------------------------------------------------------------------------------------------------------------------------
func (dfa *DFA_PoolType) NoneVisited() {
	for ii, vv := range dfa.Pool {
		if vv.IsUsed {
			dfa.Pool[ii].Visited = false
		}
	}
}

// -----------------------------------------------------------------------------------------------------------------------------------------------
func (dfa *DFA_PoolType) FindNotVisited() int {
	for ii, vv := range dfa.Pool {
		if vv.IsUsed && !vv.Visited {
			return ii
		}
	}
	return -1
}

// -----------------------------------------------------------------------------------------------------------------------------------------------
func (dfa *DFA_PoolType) GetDFAName(StateSet []int) int {
	A := dfa.GetDFA()
	StateSet = com.USortIntSlice(StateSet) // Make set unique
	dfa.Pool[A].StateSet = StateSet
	dfa.Pool[A].StateName = com.NameOf(StateSet)
	dfa.Pool[A].Visited = false
	return A
}

// -----------------------------------------------------------------------------------------------------------------------------------------------
func (dfa *DFA_PoolType) AddEdge(fr, to int, on string) {
	dfa.Pool[fr].Next2 = append(dfa.Pool[fr].Next2, nfa.TransitionType{IsLambda: false, On: on, To: to, From: fr, LineNo: dbgo.LINE(2)})
	// OLD:
	// Check if edge already exists - if so skip this
	// if !dfa.edgeExists(fr, to, on) {
	// 	dfa.Pool[fr].Next2 = append(dfa.Pool[fr].Next2, nfa.TransitionType{IsLambda: false, On: on, To: to, From: fr, LineNo: dbgo.LINE(2)})
	// }
}

// -----------------------------------------------------------------------------------------------------------------------------------------------
// func (dfa *DFA_PoolType) edgeExists(fr, to int, on string) bool {
// 	//	for _, vv := range dfa.Pool[fr].Next2 {
// 	//		if vv.To == to && vv.On == on {
// 	//			return true
// 	//		}
// 	//	}
// 	return false
// }

// ConvertNFAToDFA will convert from a non-deterministic finite state automatata to a deterministic one.
func (dfa *DFA_PoolType) ConvertNFAToDFA(nn *nfa.NFA_PoolType) {
	// Note: http://www.w3schools.com/charsets/ref_utf_dingbats.asp
	//       http://www.utf8-chartable.de/unicode-utf8-table.pl?start=768
	// Tokens must contain - Start Line No, Col No, File Name, End*, Macro Translation(pushback)
	StartState := nn.InitState
	dfa.NoneVisited()

	nn.Sigma = nn.GenerateSigma()
	dfa.Sigma = nn.Sigma
	dbgo.DbPrintf("dfa2", "Sigma at top ->%s<-\n", dfa.Sigma)

	// Build initial state
	dfa_set := nn.LambdaClosure([]int{StartState}) // Find all the lambda closures from specified state
	dfa_set = append(dfa_set, StartState)          // Add in initial state
	dfa_set = com.USortIntSlice(dfa_set)           // Make set unique
	dbgo.DbPrintf("db_DFAGen", "\nStart: %s, \u03a3 =->%s<-, %s\n", dbgo.SVar(dfa_set), dfa.Sigma, dbgo.LF())
	A := dfa.GetDFAName(dfa_set)
	if r, is, info, Is0Ch := nn.IsTerminalState(dfa_set); is {
		dfa.Pool[A].Rv = r
		dfa.Pool[A].Is0Ch = Is0Ch
		dfa.Pool[A].Info = info
	} else {
		dfa.Pool[A].Info = info
	}
	if dbgo.IsDbOn("db_DFAGen") {
		dfa.DumpPool(false)
	}
	// Look at all the locaitons we can get to from this "state"
	for _, S := range dfa.Sigma {
		StateSet := nn.LambdaClosureSet(dfa_set, string(S))
		dbgo.DbPrintf("db_DFAGen", "FOR INITIAL state ->%s<- StateSet=%s, %s\n", string(S), dbgo.SVar(StateSet), dbgo.LF())
		if len(StateSet) > 0 {
			dbgo.DbPrintf("db_DFAGen", "Have a non-empty result, %s\n", dbgo.LF())
			dbgo.DbPrintf("db_DFAGen", "<><><> this is the point where we should check to see if 'S' is DOT or NCCL, %s\n", dbgo.LF())

			StateSetT := nn.LambdaClosure(StateSet) // need to lambda complete the state set
			StateSet = append(StateSet, StateSetT...)
			StateSet = com.USortIntSlice(StateSet) // Make set unique
			dbgo.DbPrintf("db_DFAGen", "    Output Is %s, %s\n", dbgo.SVar(StateSet), dbgo.LF())
			B := 0
			if t := dfa.HaveStateAlready(StateSet); t != -1 { // Have Already
				B = t
				dbgo.DbPrintf("db_DFAGen", "    Already have this state at location %d, %s\n", t, dbgo.LF())
			} else {
				B = dfa.GetDFAName(StateSet)
				dbgo.DbPrintf("db_DFAGen", "    *** New state %d, %s\n", B, dbgo.LF())
			}
			dfa.AddEdge(A, B, string(S))
			dbgo.DbPrintf("db_DFAGen", "    *** Before (top) %s\n", dbgo.LF())
			if r, is, info, Is0Ch := nn.IsTerminalState(StateSet); is {
				dfa.Pool[B].Rv = r
				dfa.Pool[B].Is0Ch = Is0Ch
				dfa.Pool[B].Info = info
				dbgo.DbPrintf("db_DFAGen", "    *** New state %d, %s\n", B, dbgo.LF())
			} else if _, is, info, Is0Ch := nn.IsNonTerminalPushPopState(StateSet); is {
				dfa.Pool[B].Is0Ch = Is0Ch
				dfa.Pool[B].Info = info
				dbgo.DbPrintf("db_DFAGen", "    *** New info for state %d, %s\n", B, dbgo.LF())
			} else {
				dfa.Pool[B].Info = info
				dbgo.DbPrintf("db_DFAGen", "    *** NO State Info for state %d, %s\n", B, dbgo.LF())
			}
			dbgo.DbPrintf("db_DFAGen", "    *** After (top) %s\n", dbgo.LF())
			if dbgo.IsDbOn("db_DFAGen") {
				fmt.Printf("for %s StateSet=%s, A=%d, B=%s %s\n", string(S), dbgo.SVar(StateSet), A, dbgo.SVar(B), dbgo.LF())
				dfa.DumpPool(false)
			}
		}
	}
	dfa.Pool[A].Visited = true

	dbgo.DbPrintf("db_DFAGen", "\nBefore Main Loop, %s\n", dbgo.LF())
	limit := 0
	for stateToDo := dfa.FindNotVisited(); stateToDo != -1; stateToDo = dfa.FindNotVisited() {
		dbgo.DbPrintf("db_DFAGen", "\nMain Loop: !!TOP!! State:%d\n", stateToDo)
		// -----------------------------------------------------------------------------------------------------------
		if !dfa.Pool[stateToDo].Visited {
			dfa_set := nn.LambdaClosure(dfa.Pool[stateToDo].StateSet)  // Find all the lambda closures from specified state
			dfa_set = append(dfa_set, dfa.Pool[stateToDo].StateSet...) // Add in initial state
			dfa_set = com.USortIntSlice(dfa_set)                       // Make set unique
			for _, S := range dfa.Sigma {
				StateSet := nn.LambdaClosureSet(dfa_set, string(S))
				dbgo.DbPrintf("db_DFAGen", "    for initial state %s StateSet=%s, %s\n", string(S), dbgo.SVar(StateSet), dbgo.LF())
				dbgo.DbPrintf("db_DFAGen", "<><><> this is the point where we should check to see if 'S' is DOT or NCCL, %s\n", dbgo.LF())
				if len(StateSet) > 0 {
					dbgo.DbPrintf("db_DFAGen", "    >>> Have a non-empty result, Input Is %s, %s\n", dbgo.SVar(StateSet), dbgo.LF())
					StateSetT := nn.LambdaClosure(StateSet) // need to lambda complete the state set
					StateSet = append(StateSet, StateSetT...)
					StateSet = com.USortIntSlice(StateSet) // Make set unique
					dbgo.DbPrintf("db_DFAGen", "    >>> Output Is %s, %s\n", dbgo.SVar(StateSet), dbgo.LF())
					B := 0
					if t := dfa.HaveStateAlready(StateSet); t != -1 { // Have Already
						B = t
						dbgo.DbPrintf("db_DFAGen", "    Already have this state at location %d, %s\n", t, dbgo.LF())
					} else {
						B = dfa.GetDFAName(StateSet)
					}
					dfa.AddEdge(stateToDo, B, string(S))
					dbgo.DbPrintf("db_DFAGen", "    *** Before %s\n", dbgo.LF())
					if r, is, info, Is0Ch := nn.IsTerminalState(StateSet); is {
						dfa.Pool[B].Rv = r
						dfa.Pool[B].Is0Ch = Is0Ch
						dfa.Pool[B].Info = info
						dbgo.DbPrintf("db_DFAGen", "    *** New state %d, %s\n", B, dbgo.LF())
					} else if _, is, info, Is0Ch := nn.IsNonTerminalPushPopState(StateSet); is {
						dfa.Pool[B].Is0Ch = Is0Ch
						dfa.Pool[B].Info = info
						dbgo.DbPrintf("db_DFAGen", "    *** New info for state %d, %s\n", B, dbgo.LF())
					} else {
						dbgo.DbPrintf("db_DFAGen", "    *** NO State Info for state %d, %s\n", B, dbgo.LF())
						dfa.Pool[B].Info = info
					}
					dbgo.DbPrintf("db_DFAGen", "    *** After %s\n", dbgo.LF())
					if dbgo.IsDbOn("db_DFAGen") {
						fmt.Printf("    Add New Edge on %s fr %d to %d, %s\n", string(S), stateToDo, B, dbgo.LF())
						fmt.Printf("    for %s StateSet=%s, A(stateToDo)=%d, %s\n", string(S), dbgo.SVar(StateSet), stateToDo, dbgo.LF())
						dfa.DumpPool(false)
					}
				}
			}
		}
		// -----------------------------------------------------------------------------------------------------------
		dfa.Pool[stateToDo].Visited = true
		limit++
		if limit > 50000 {
			break
		}
	}
	dfa.VerifyMachine()

}

// -----------------------------------------------------------------------------------------------------------------------------------------------
//

func (dfa *DFA_PoolType) DumpPoolJSON(fo io.Writer, td string, tn int) {
	fmt.Fprintf(fo, `{"Input":%q, "Rv":%d, "Start": %d, "States":[%s`, td, tn, dfa.InitState, "\n")
	for ii, vv := range dfa.Pool {
		if vv.IsUsed {
			fmt.Fprintf(fo, ` { "Sn":%d, `, ii)
			if vv.Rv > 0 {
				fmt.Fprintf(fo, ` "Term":%d, `, vv.Rv)
			}
			fmt.Fprintf(fo, ` "Edge":[ `)
			com := ""
			for _, ww := range vv.Next2 {
				fmt.Fprintf(fo, "%s{ \"On\":\"%s\", \"Fr\":%d, \"To\":%d }", com, ww.On, ww.From, ww.To)
				com = ", "
			}
			fmt.Fprintf(fo, "]}\n")
		}
	}
	fmt.Fprintf(fo, "]}\n")
}

// -----------------------------------------------------------------------------------------------------------------------------------------------
//

func (dfa *DFA_PoolType) GenerateGVFile(fo io.Writer, td string, tn int) {
	// fmt.Fprintf(fo, `{"Input":%q, "Rv":%d, "Start": %d, "Sigma":%q, "States":[%s`, td, tn, dfa.InitState, dfa.Sigma, "\n")
	siz := "5"
	if len(dfa.Pool) > 100 {
		siz = "50"
	} else if len(dfa.Pool) > 50 {
		siz = "40"
	} else if len(dfa.Pool) > 30 {
		siz = "55"
	} else if len(dfa.Pool) > 20 {
		siz = "10"
	} else if len(dfa.Pool) > 10 {
		siz = "8"
	}
	fmt.Fprintf(fo,
		`digraph finite_state_machine {
	rankdir=LR;
	size="18,%s"
`, siz)
	// size= for bigger graph - should be configurable with tests -

	var term []int
	for ii, vv := range dfa.Pool {
		if vv.IsUsed {
			if vv.Rv > 0 {
				term = append(term, ii)
			}
		}
	}
	s := ""
	cc := ""
	for _, vv := range term {
		s += cc + fmt.Sprintf("s%d", vv)
		cc = " "
	}
	fmt.Fprintf(fo,
		`	node [shape = doublecircle]; %s;
	node [shape = circle];
`, s)

	for _, vv := range dfa.Pool {
		if vv.IsUsed {
			for _, ww := range vv.Next2 {
				if ww.On[0] <= ' ' {
					fmt.Fprintf(fo, "	s%d -> s%d [ label = %q ];\n", ww.From, ww.To, re.EscapeStrForGV(ww.On))
				} else {
					fmt.Fprintf(fo, "	s%d -> s%d [ label = \"%s\" ];\n", ww.From, ww.To, re.EscapeStr(ww.On))
				}
			}
		}
	}
	fmt.Fprintf(fo, "}\n")
}

// ------------------------------------------------------------------------------------------------------------------------------------------------------
// 1. Build the maping info InputMap0, Ranges, etc.			1hr
// 2. Test this with current machines						1hr
// 3. Generate the array-based machine.						2hr
// 4. Hand verify - write an output function
// 5. Check this.										++	4hr
//
// Let's build the map stuff as it's own little project, ./smap
//	fx:		minv, maxv, m0, m1 := smap.BuildMapString ( Sigma, NoMapRn )
//	fx:		k := smap.MapRune ( rn )
// ------------------------------------------------------------------------------------------------------------------------------------------------------

type MachineStatesType struct {
	Rv   int          //
	Tau  bool         //
	Info nfa.InfoType //
	To   []int        //
}

type dfaTable struct {
	InitState int                  //
	N_States  int                  //
	Width     int                  //
	SMap      *smap.SMapType       //
	Machine   []*MachineStatesType //
}

// -----------------------------------------------------------------------------------------------------------------------------------------------
func (dfa *DFA_PoolType) NumberOfStates() (N_States int) {
	N_States = 0
	for _, vv := range dfa.Pool {
		if vv.IsUsed {
			N_States++
		}
	}
	return
}

// -----------------------------------------------------------------------------------------------------------------------------------------------
func (dfa *DFA_PoolType) ConvertToTable() (rv dfaTable) {
	rv.InitState = dfa.InitState
	rv.SMap = smap.NewSMapType(dfa.Sigma, re.R_not_CH)
	rv.N_States = dfa.NumberOfStates()
	rv.Machine = make([]*MachineStatesType, rv.N_States, rv.N_States)
	XLen := rv.SMap.Length()
	rv.Width = XLen

	jj := 0
	for ii, vv := range dfa.Pool {
		if vv.IsUsed {
			rv.Machine[jj] = &MachineStatesType{Rv: vv.Rv, Info: vv.Info, Tau: vv.Is0Ch}
			rv.Machine[jj].To = make([]int, XLen, XLen)
			for kk := 0; kk < XLen; kk++ {
				rv.Machine[jj].To[kk] = -1
			}
			dot := false
			dotTo := 0
			alpha := false
			alphaTo := 0
			num := false
			numTo := 0
			for _, ww := range vv.Next2 {
				//if s == re.X_DOT {
				// return "\u2022"	// Middle Bullet
				if ww.On == re.X_DOT {
					// fmt.Printf("Found a dot, at jj=%d\n", jj)
					dot = true
					dotTo = ww.To
				} else if ww.On == re.X_ALPHA {
					// fmt.Printf("Found a alpha, at jj=%d\n", jj)
					alpha = true
					alphaTo = ww.To
					rr, _ := utf8.DecodeRune([]byte(ww.On))
					xx := rv.SMap.MapRune(rr)
					rv.Machine[jj].To[xx] = ww.To
				} else if ww.On == re.X_NUMERIC {
					// dbgo.DbPrintf("dfa2", "Found a numeric, at jj=%d\n", jj)
					num = true
					numTo = ww.To
					rr, _ := utf8.DecodeRune([]byte(ww.On))
					xx := rv.SMap.MapRune(rr)
					rv.Machine[jj].To[xx] = ww.To
				} else {
					rr, _ := utf8.DecodeRune([]byte(ww.On))
					xx := rv.SMap.MapRune(rr)
					rv.Machine[jj].To[xx] = ww.To
				}
			}
			if num {
				for ii := 0; ii < rv.Width; ii++ {
					// rn := rune(ii + rv.SMap.MinV)
					rn := rv.SMap.SigmaRN[ii]
					// func (smap *SMapType) ReverseMapRune(x int) rune {
					// dbgo.DbPrintf("dfa2", "numeric: ->%s<-\n", string(rn))
					if rv.Machine[jj].To[ii] == -1 && unicode.IsDigit(rn) {
						rv.Machine[jj].To[ii] = numTo
					}
				}
			}
			if alpha {
				for ii := 0; ii < rv.Width; ii++ {
					// rn := rune(ii + rv.SMap.MinV)
					rn := rv.SMap.SigmaRN[ii]
					if rv.Machine[jj].To[ii] == -1 && (unicode.IsUpper(rn) || unicode.IsLower(rn)) {
						rv.Machine[jj].To[ii] = alphaTo
					}
				}
			}
			if dot {
				for ii := 0; ii < rv.Width; ii++ {
					if rv.Machine[jj].To[ii] == -1 {
						rv.Machine[jj].To[ii] = dotTo
					}
				}
			}
			if ((vv.Info.Action&com.A_Pop) != 0 || (vv.Info.Action&com.A_Push) != 0 || (vv.Info.Action&com.A_Reset) != 0) && !vv.Info.HardMatch { // xyzzy8
				dbgo.DbPrintf("dfa7", "Info2-in-TabGen: State[%d] a Push/Call or a Pop/Return/Reset optration Must be a terminal state with fixed string matched, 'a*' or 'a?' is not a legitimate match.\n", ii)
				//for ii := 0; ii < rv.Width; ii++ {
				//	if rv.Machine[jj].To[ii] == -1 {
				//		rv.Machine[jj].To[ii] = dotTo
				//	}
				//}
			}
			jj++
		}
	}
	return
}

// -----------------------------------------------------------------------------------------------------------------------------------------------
// func (dfa *DFA_PoolType) FinializeDFA() {
// 	dt := dfa.ConvertToTable()
// 	dfa.MTab = &dt
// }

// -----------------------------------------------------------------------------------------------------------------------------------------------
func (dfa *DFA_PoolType) PrintStateMachine(fo io.Writer, format string) {

	dt := dfa.MTab

	switch format {
	case "text":

		fmt.Fprintf(fo, `
Sigma = %q
InitState = %d
N_States = %d
Width = %d
`, dfa.Sigma, dt.InitState, dt.N_States, dt.Width)

		SigmaArray := make([]rune, dt.Width, dt.Width)
		pp := 0
		for jj := 0; jj < dt.Width; jj++ {
			if pp < len(dfa.Sigma) {
				rn, sz := utf8.DecodeRune([]byte(dfa.Sigma[pp:]))
				SigmaArray[jj] = rn
				pp += sz
			} else {
				// SigmaArray[jj] = rune(0xFBAD)
				SigmaArray[jj] = dt.SMap.NoMap
			}
		}

		fmt.Fprintf(fo, "SMap = %+v\n", dt.SMap)
		fmt.Fprintf(fo, "%-6s : %-20s  %-5s %-5s %-4s %-4s    | ", "State", "Rv/Name", "Actn", "Hard", "Next", "Leng")
		for jj := 0; jj < dt.Width; jj++ {
			fmt.Fprintf(fo, "   %2d", jj)
		}
		fmt.Fprintf(fo, "\n")
		fmt.Fprintf(fo, "%-6s : %-20s  %-5s %-5s %-4s %-4s    | ", "======", "====/===============", "-----", "-----", "----", "----")
		for jj := 0; jj < dt.Width; jj++ {
			fmt.Fprintf(fo, "   %2s", "--")
		}
		fmt.Fprintf(fo, "\n")

		fmt.Fprintf(fo, "%-6s : %-4s/%-15s  %-5s %-5s %-4s %-4s    | ", " ", " ", " ", " ", " ", " ", " ")
		for jj := 0; jj < dt.Width; jj++ {
			if SigmaArray[jj] < ' ' {
				s := fmt.Sprintf("%q", string(SigmaArray[jj]))
				s = s[1:]
				fmt.Fprintf(fo, "   %2s", s[0:len(s)-1])
			} else {
				fmt.Fprintf(fo, "   %2s", string(SigmaArray[jj]))
			}
		}
		fmt.Fprintf(fo, "\n")

		fmt.Fprintf(fo, "%-6s : %-20s  %-5s %-5s %-4s %-4s    | ", "======", "====/===============", "-----", "-----", "----", "----")
		for jj := 0; jj < dt.Width; jj++ {
			fmt.Fprintf(fo, "   %2s", "--")
		}
		fmt.Fprintf(fo, "\n")

		fx := func(rv int) string {
			return in.LookupTokenName(rv)
		}

		for ii, vv := range dt.Machine {
			tau := " "
			if vv.Tau {
				tau = "\u03c4"
			}
			if vv.Info.Action == 0 {
				fmt.Fprintf(fo, "m[%3d] : %4d/%-15s%s %5s %5v %4d %4d    | ", ii, vv.Rv, fx(vv.Rv), tau, "", vv.Info.HardMatch, vv.Info.NextState, vv.Info.MatchLength)
			} else {
				fmt.Fprintf(fo, "m[%3d] : %4d/%-15s%s %5x %5v %4d %4d    | ", ii, vv.Rv, fx(vv.Rv), tau, vv.Info.Action, vv.Info.HardMatch, vv.Info.NextState, vv.Info.MatchLength)
			}
			for _, ww := range vv.To {
				if ww == -1 {
					// fmt.Fprintf(fo, "   %2s", "\u26d4") // No Entry
					fmt.Fprintf(fo, "   %2s", "\u2629")
				} else {
					fmt.Fprintf(fo, "   %2d", ww)
				}
			}
			fmt.Fprintf(fo, "\n")
		}
		fmt.Fprintf(fo, "\n")

		// should have 'json' or 'go' for output format?
	default:
		fmt.Fprintf(os.Stderr, "Invalid Output Format for dfa.OutputInFomrat, %s,  should be 'text'.\n", format)
		os.Exit(1)
	}
}

/* vim: set noai ts=4 sw=4: */
