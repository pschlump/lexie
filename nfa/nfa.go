//
// N F A - Part of Lexie Lexical Generation System
//
// (C) Philip Schlump, 2014-2015.
// Version: 1.0.8
// BuildNo: 203
//
// Special Thanks to 2C-Why, LLC for supporting this project.
//

// xyzzy-NCCL - add in NCCL at this point

package nfa

import (
	"fmt"
	"io"
	"sort"
	"unicode/utf8"

	"github.com/pschlump/dbgo"
	"github.com/pschlump/lexie/com"
	"github.com/pschlump/lexie/re"
)

type TransitionType struct {
	On         string
	IsLambda   bool
	Is0ChMatch bool // false(default) => not set,       true(set) => match occures on any char that is not an accepted char (a*)-> next state on 'b'
	From       int
	To         int
	LineNo     string // LineNo where added
}

type InfoType struct {
	Action       int
	MatchLength  int
	NextState    int
	HardMatch    bool
	ReplStr      string
	ReservedWord bool //
}

type NFA_Type struct {
	Next2      []TransitionType //
	Rv         int              // 0 indicates not assigned, non-terminal
	Info       InfoType         //
	TRuleMatch int              // Order that stuff was creaed in
	NextFree   int              //		For free list
	IsUsed     bool             //		For Free list
	A_IAm      int              //		Debug Usage
	LineNo     string           // LineNo where added
	Visited    bool             // Visited Marker
	TokType    re.LR_TokType    // Type of token, marked for LR_STAR, LR_PLUS, LR_QUEST, LR_OP_BR etc.
}

type ReSetType struct {
	Re         string        //		String this RE is from
	Rv         int           //		Terminal value to return if it matches
	TRuleMatch int           // Order that stuff was created in
	IsDirty    bool          //		Indicates "Re" has been chagned without rebuilding ParsedRe
	HasDot     bool          //		Indicates that completed ParsedRe has a DOT in it
	HasNCCL    bool          //		Indicates that completed ParsedRe has a NCCL in it
	Deleted    bool          //		Inidiates that this item has been deleted
	ParsedRe   *re.LexReType //		Parsed Tree
	Info       InfoType      //
}

type NFA_PoolType struct {
	Pool      []NFA_Type
	Cur       int
	Top       int
	NextFree  int
	InitState int
	Sigma     string //
	ReSet     []*ReSetType
	IsDirty   bool
	//expand_dot_over_sigma bool
	//leave_dot_in_nfa      bool
	//gen_nccl              bool
}

const InitNFASize = 3

// Create a new NFA pool
func NewNFA_Pool() *NFA_PoolType {
	return &NFA_PoolType{
		Pool:     make([]NFA_Type, InitNFASize, InitNFASize),
		Cur:      0,
		Top:      InitNFASize,
		NextFree: -1,
		Sigma:    "",
		ReSet:    make([]*ReSetType, 0, 100),
		//expand_dot_over_sigma: false,
		//leave_dot_in_nfa:      true,
		//gen_nccl:              false,
	}
}

// Allocate an NFA tree node
func (nn *NFA_PoolType) GetNFA() int {
	//fmt.Printf("at %s\n", dbgo.LF())
	tmp := 0
	if nn.Cur < nn.Top && !nn.Pool[nn.Cur].IsUsed {
		//fmt.Printf("at %s\n", dbgo.LF())
		tmp = nn.Cur
		nn.Cur++
	} else if nn.Cur >= nn.Top || nn.NextFree == -1 {
		//fmt.Printf("at %s, nn.Cur=%d nn.Top=%d nn.NextFree=%d\n", dbgo.LF(), nn.Cur, nn.Top, nn.NextFree)
		nn.Top = 2 * nn.Top
		newPool := make([]NFA_Type, nn.Top, nn.Top) // extend array
		copy(newPool, nn.Pool)
		nn.Pool = newPool
		tmp = nn.Cur
		nn.Cur++
	} else {
		//fmt.Printf("at %s\n", dbgo.LF())
		tmp = nn.NextFree
		nn.NextFree = nn.Pool[tmp].NextFree
	}
	nn.Pool[tmp].NextFree = -1
	nn.Pool[tmp].Rv = 0
	nn.Pool[tmp].Next2 = nn.Pool[tmp].Next2[:0]
	nn.Pool[tmp].IsUsed = true
	nn.Pool[tmp].A_IAm = tmp
	nn.Pool[tmp].LineNo = dbgo.LINE(2)
	return tmp
}

// Free an NFA tree node
func (nn *NFA_PoolType) FreeNFA(ii int) {
	nn.Pool[ii].IsUsed = false
	nn.Pool[ii].NextFree = nn.NextFree
	nn.Pool[ii].Rv = 0
	nn.Pool[ii].Next2 = nn.Pool[ii].Next2[:0]
	nn.NextFree = ii
}

// Return the start state number
func (nn *NFA_PoolType) Pos0Start() int {
	return nn.InitState
}

func (nn *NFA_PoolType) DiscardPool() {
	nn.Pool = make([]NFA_Type, InitNFASize, InitNFASize)
	nn.Cur = 0
	nn.Top = InitNFASize
	nn.NextFree = -1
}

// Add a new string to the NFA graph
//func (nn *NFA_PoolType) AddString(buf string, tRv int) {
//	Cur := nn.Pos0Start() // Start new NFA at Init Pos
//	for ii := range buf {
//		ss := buf[ii : ii+1]
//		if y, hasIt := nn.Pool[Cur].Next[ss]; hasIt {
//			Cur = y
//		} else {
//			x := nn.GetNFA()
//			nn.Pool[Cur].Next[ss] = x
//			Cur = x
//		}
//	}
//	nn.Pool[Cur].Rv = tRv
//}

//
// Take a regular expression and add it to the NFA.  The tRuleMatchId is the unique ID for this R.E.
// Usually the tRuleMatchId is the subscript of the array where the R.W. came from.
//
// Params:
//	re				The regular expression string.
//	reFlags			Flags for parsing/understanding the r.e.
//						'c' - case insensitive
//						'n' - '.' will not match \n, defualt is cross new line matches
//	tRuleMatchId	A number that will be returned with the match of this r.e. -- this is returned as an array of rule matches.
//	tRv				if non-zero, then this is a terminal rule and it returns this token number.
//

func DumpInfo(info InfoType) string {
	// xyzzy22 - convert info.NextState - > state name
	sa := com.ConvertActionFlagToString(info.Action)
	// nsname := lex.ConvertMachineNameToText(info.NextState)
	return fmt.Sprintf("Action: %s(%02x), Ns:%d, MatchLength:%d ReplStr:->%s<-", sa, info.Action, info.NextState, info.MatchLength, info.ReplStr)
}

func (nn *NFA_PoolType) AddLambda(fr, to int) {
	if fr == to {
		dbgo.DbPrintf("nfa4", "!!! Rejecting lambda(\u03bb) loop %d to %d, called from %s\n", fr, to, dbgo.LF(2))
		// fmt.Printf("!!! Rejecting lambda(L) loop %d to %d, called from %s\n", fr, to, dbgo.LF(2))
		return
	}
	if !nn.EdgeExists(fr, to, "", true) {
		nn.Pool[fr].Next2 = append(nn.Pool[fr].Next2, TransitionType{IsLambda: true, To: to, From: fr, LineNo: dbgo.LINE(2)})
	}
}

func (nn *NFA_PoolType) AddLambdaSpecial(fr, to int) {
	if fr == to {
		dbgo.DbPrintf("nfa4", "!!! Rejecting Tau(\u03c4) loop %d to %d, called from %s -- This could be a serious error if at the end of a match, dropping a terminal state\n", fr, to, dbgo.LF(2))
		// fmt.Printf("!!! Rejecting lambda(L) loop %d to %d, called from %s\n", fr, to, dbgo.LF(2))
		return
	}
	if !nn.EdgeExists(fr, to, "", true) {
		// Is0ChMatch bool // false(default) => not set,       true(set) => match occures on any char that is not an accepted char (a*)-> next state on 'b'
		nn.Pool[fr].Next2 = append(nn.Pool[fr].Next2, TransitionType{IsLambda: true, Is0ChMatch: true, To: to, From: fr, LineNo: dbgo.LINE(2)})
	}
}

func (nn *NFA_PoolType) Mark(cur int, t re.LR_TokType) {
	nn.Pool[cur].TokType = t
}

func (nn *NFA_PoolType) AddEdge(fr, to int, on string) {
	// Check if edge already exists - if so skip this
	if !nn.EdgeExists(fr, to, on, false) {
		nn.Pool[fr].Next2 = append(nn.Pool[fr].Next2, TransitionType{IsLambda: false, On: on, To: to, From: fr, LineNo: dbgo.LINE(2)})
	}
}

func (nn *NFA_PoolType) EdgeExists(fr, to int, on string, lambda bool) bool {
	//	for _, vv := range nn.Pool[fr].Next2 {
	//		if vv.To == to && vv.IsLambda && lambda {
	//			return true
	//		} else if vv.To == to && vv.On == on && !vv.IsLambda {
	//			return true
	//		}
	//	}
	return false
}

func runeInString(ww rune, vv string) bool {
	for _, ss := range vv {
		if ss == ww {
			return true
		}
	}
	return false
}

func (nn *NFA_PoolType) ConvParsedReToNFA(depth int, lr *re.LexReType, CurIn int, Children []re.ReTreeNodeType) (int, int) {

	if depth > 3000 {
		panic("infinite recursion")
	}

	Cur := CurIn
	if CurIn == -1 {
		dbgo.DbPrintf("nfa2", " Initialize with CurIn of -1\n")
		Cur = nn.GetNFA()
		CurIn = Cur
	}
	dbgo.DbPrintf("db_NFA", "CurIn: %d TOP, depth=%d\n", Cur, depth)

	for _, vv := range Children {

		dbgo.DbPrintf("db_NFA", "CurIn: %d Loop TOP, vv=%s, %s\n", Cur, dbgo.SVarI(vv), dbgo.LF())

		switch vv.LR_Tok {

		// ------------------------------------------------------------------------------------------------------
		// Text Items
		// ------------------------------------------------------------------------------------------------------
		case re.LR_CL_BR: // }
			fallthrough
		case re.LR_COMMA: // ,
			fallthrough
		case re.LR_CL_PAR: //  - a text node only
			fallthrough
		case re.LR_E_CCL: // ] - a text node only
			fallthrough
		case re.LR_MINUS: // -
			fallthrough
		case re.LR_Text: //			-- Add a node to list, move right
			dbgo.DbPrintf("nfa2", "Text Add Item ->%s<-\n", vv.Item)
			for _, ss := range vv.Item {
				x := nn.GetNFA()
				nn.AddEdge(Cur, x, string(ss))
				Cur = x
			}
		case re.LR_CARROT: // ^		-- BOL Like Text - special
			x := nn.GetNFA()
			nn.AddEdge(Cur, x, re.X_BOL)
			Cur = x
		case re.LR_DOLLAR: // $		-- EOL Like Text - special
			x := nn.GetNFA()
			nn.AddEdge(Cur, x, re.X_EOL)
			Cur = x

		// ------------------------------------------------------------------------------------------------------
		// CCLs
		// ------------------------------------------------------------------------------------------------------
		case re.LR_CCL: // [...]		-- CCL Node (Above)
			x := nn.GetNFA()
			for _, ss := range vv.Item {
				nn.AddEdge(Cur, x, string(ss))
			}
			Cur = x
		case re.LR_DOT: // .			-- Like speical CCL
			x := nn.GetNFA()
			nn.AddEdge(Cur, x, re.X_DOT)
			Cur = x
		case re.LR_N_CCL: // [^...]	-- N_CCL Node, Hm....
			x := nn.GetNFA()
			nn.AddEdge(Cur, x, re.X_N_CCL)
			Cur = x

		// ------------------------------------------------------------------------------------------------------
		// Loops ( *, +, ? )
		// ------------------------------------------------------------------------------------------------------
		case re.LR_STAR: // *			-- Error if 1st char, else take prev item from list, star and replace it.
			childStart, childEnd := nn.ConvParsedReToNFA(depth+1, lr, Cur, vv.Children[0:1])
			dbgo.DbPrintf("db_NFA", "Return from STAR call, %d, %d, %s\n", childStart, childEnd, dbgo.LF())
			nn.Mark(Cur, re.LR_STAR)
			nn.AddLambda(Cur, childEnd)
			nn.AddLambda(childEnd, Cur)
			tail := nn.GetNFA()
			nn.AddLambdaSpecial(childEnd, tail) // Is0ChMatch bool // false(default) => not set,       true(set) => match occures on any char that is not an accepted char (a*)-> next state on 'b'
			Cur = tail
		case re.LR_PLUS: // +			-- Like *
			childStart, childEnd := nn.ConvParsedReToNFA(depth+1, lr, Cur, vv.Children[0:1])
			dbgo.DbPrintf("db_NFA", "Return from PLUS call, %d, %d, %s\n", childStart, childEnd, dbgo.LF())
			nn.Mark(Cur, re.LR_PLUS)
			// nn.AddLambda(Cur, childStart)
			nn.AddLambda(childEnd, childStart)
			tail := nn.GetNFA()
			nn.AddLambdaSpecial(childEnd, tail)
			Cur = tail
		case re.LR_QUEST: // ?			-- Like *
			childStart, childEnd := nn.ConvParsedReToNFA(depth+1, lr, Cur, vv.Children[0:1])
			dbgo.DbPrintf("db_NFA", "Return from QUESTION call, %d, %d, %s\n", childStart, childEnd, dbgo.LF())
			nn.Mark(Cur, re.LR_QUEST)
			nn.AddLambda(childStart, childEnd)
			// xyzzy Is0ChMatch -- will need to add an additional state to make this work!
			tail := nn.GetNFA()
			nn.AddLambdaSpecial(childEnd, tail)
			Cur = tail
		case re.LR_OP_BR: // {			-- xyzzy - need to add this
			if vv.Mm >= 0 && vv.Mm < vv.Nn {
				nn.Mark(Cur, re.LR_OP_BR)
				pCur, childStart, childEnd := Cur, Cur, Cur
				for kk := 0; kk < vv.Mm; kk++ {
					childStart, childEnd = nn.ConvParsedReToNFA(depth+1, lr, pCur, vv.Children[0:1])
					pCur = childEnd
				}
				if vv.Nn == re.InfiniteIteration {
					nn.AddLambda(childEnd, childStart)
				} else {
					for kk := vv.Mm; kk < vv.Nn; kk++ {
						childStart, childEnd = nn.ConvParsedReToNFA(depth+1, lr, pCur, vv.Children[0:1])
						nn.AddLambda(childStart, childEnd)
						pCur = childEnd
					}
				}
				Cur = pCur
			}
			// xyzzy Is0ChMatch -- will need to add an additional state to make this work!
			tail := nn.GetNFA()
			nn.AddLambdaSpecial(Cur, tail)
			Cur = tail

		// ------------------------------------------------------------------------------------------------------
		// OR, |
		// ------------------------------------------------------------------------------------------------------

		case re.LR_OR: // |
			if len(vv.Children) > 0 {
				dest := nn.GetNFA()
				beg := Cur
				nn.Mark(Cur, re.LR_OR)
				dbgo.DbPrintf("db_NFA", "*************************** OR dest = %d, Tree is (vv)=%s, %s\n", dest, dbgo.SVarI(vv), dbgo.LF())
				dbgo.DbPrintf("db_NFA", "*************************** This would be the point to add (N) terms to graph\n")
				for jj := range vv.Children {
					// dbgo.DbPrintf("db_NFA", "OR Before recursive call\n")
					childStart, childEnd := nn.ConvParsedReToNFA(depth+1, lr, beg, vv.Children[jj:jj+1])
					dbgo.DbPrintf("db_NFA", "OR After recursive call, %d -> %d, beg=%d dest=%d\n", childStart, childEnd, beg, dest)
					//nn.AddLambda(beg, childStart)
					nn.AddLambda(childEnd, dest)
				}
				dbgo.DbPrintf("db_NFA", "OR Final Dest  %d \n", dest)
				Cur = dest // Note - inside loop, if no iterations then Cur stays unchanged.
			}

		// ------------------------------------------------------------------------------------------------------
		// Grouping ()
		// ------------------------------------------------------------------------------------------------------
		case re.LR_null: //
			dbgo.DbPrintf("db_NFA", "LR_null\n")
			if false {
				dest := nn.GetNFA()
				// childStart, childEnd := nn.ConvParsedReToNFA(depth+1, lr, -1, vv.Children)
				childStart, childEnd := nn.ConvParsedReToNFA(depth+1, lr, Cur, vv.Children)
				nn.AddLambda(Cur, childStart)
				nn.AddLambda(childEnd, dest)
				Cur = dest
			} else {
				dbgo.DbPrintf("db_NFA", "Return (%d) LR_null  %d, %d, %s\n", depth, CurIn, Cur, dbgo.LF())
				return nn.ConvParsedReToNFA(depth+1, lr, Cur, vv.Children)
			}

		case re.LR_OP_PAR: // (		-- Start of Sub_Re
			dbgo.DbPrintf("db_NFA", "LR_OP_PAR\n")
			// dest := nn.GetNFA()
			_ /*childStart*/, childEnd := nn.ConvParsedReToNFA(depth+1, lr, Cur, vv.Children)
			// nn.AddLambda(childEnd, dest)
			// Cur = dest
			Cur = childEnd

		// ------------------------------------------------------------------------------------------------------
		// Error Cases
		// ------------------------------------------------------------------------------------------------------
		case re.LR_EOF: //  Should case a return.
			fallthrough
		default: // Hm...
			dbgo.DbPrintf("db_NFA", "Return (%d) via default or EOF,  %d, %d, %s\n", depth, CurIn, Cur, dbgo.LF())
			return CurIn, Cur
		}
	}

	dbgo.DbPrintf("db_NFA", "Return (%d) via bottom of code, loop ended,  %d, %d, %s\n", depth, CurIn, Cur, dbgo.LF())
	return CurIn, Cur

}

func (nn *NFA_PoolType) AddReInfo(Re string, reFlags string, tRuleMatchId int, tRv int, info InfoType) int {
	nn.IsDirty = true
	p := len(nn.ReSet)
	nn.ReSet = append(nn.ReSet, &ReSetType{Re: Re, Rv: tRv, IsDirty: false, HasDot: false, HasNCCL: false, Deleted: false, TRuleMatch: tRuleMatchId})

	lr := re.NewLexReType()
	lr.ParseRe(Re)
	lr.Sigma = lr.GenerateSigma() // This only generats Sigma for this sub-set of this RE - not for the entire set of REs
	reLen, isHard := lr.CalcLength()
	nn.ReSet[p].ParsedRe = lr

	dbgo.DbPrintf("DumpParseNodes2", "DumpParseNodes: Pass 1 in AddReInfo %s\n\n", dbgo.SVarI(lr))
	dbgo.DbPrintf("DumpParseNodes2", "DumpParseNodes: As Passed %s\n\n", dbgo.SVarI(lr.Tree.Children))

	// ------------------------------------------------------------------------------------------------------------
	// xyzzy - add HasDot, HasNCCL to this.
	// ------------------------------------------------------------------------------------------------------------

	CurIn := nn.Pos0Start() // Start new NFA at Init Pos
	CurStart, Cur := nn.ConvParsedReToNFA(0, lr, CurIn, lr.Tree.Children)

	dbgo.DbPrintf("db_NFA", "CurStar: %d CurEnd: %d, reLen=%d\n", CurStart, Cur, reLen)

	nn.Pool[Cur].Rv = tRv
	nn.Pool[Cur].TRuleMatch = tRuleMatchId
	nn.Pool[Cur].Info.MatchLength = reLen
	nn.Pool[Cur].Info.HardMatch = isHard
	nn.Pool[Cur].Info.Action = info.Action
	nn.Pool[Cur].Info.NextState = info.NextState
	nn.Pool[Cur].Info.ReplStr = info.ReplStr
	nn.ReSet[p].Info = nn.Pool[Cur].Info
	return Cur
}

func (nn *NFA_PoolType) SetReservedWord(Cur int) {
	nn.Pool[Cur].Info.ReservedWord = true
}

func (nn *NFA_PoolType) FinializeNFA() {
	nn.Sigma = nn.GenerateSigma()

	for _, vv := range nn.Pool {
		if vv.IsUsed {
			for _, ww := range vv.Next2 {
				if ww.IsLambda {
				} else if ww.On == re.X_DOT {
					// xyzzy-NCCL - add in NCCL at this point
					for _, rn := range nn.Sigma {
						if rn != re.R_DOT {
							nn.AddEdge(ww.From, ww.To, string(rn))
						}
					}
				}
			}
		}
	}

	nn.IsDirty = false // Mark as completed and ready to use. (( Should be checed in DFA generation - don't generate until finialized ))
}

func (nn *NFA_PoolType) ChangeRe(oldRe string, newRe string) {
	for ii, vv := range nn.ReSet {
		if vv.Re == oldRe {
			vv.Re = newRe
			vv.IsDirty = true
			vv.HasDot = false
			vv.Deleted = false
			vv.HasNCCL = false
			vv.ParsedRe = nil
			nn.ReSet[ii] = vv
			nn.IsDirty = true
			nn.ReSet[ii] = vv
			break
		}
	}
}

func (nn *NFA_PoolType) DeleteRe(oldRe string) {
	for ii, vv := range nn.ReSet {
		if vv.Re == oldRe {
			vv.Deleted = true
			nn.ReSet[ii] = vv
			break
		}
	}
}

func (nn *NFA_PoolType) UnDeleteRe(oldRe string) {
	for ii, vv := range nn.ReSet {
		if vv.Re == oldRe {
			vv.Deleted = false
			vv.IsDirty = true
			nn.ReSet[ii] = vv
			break
		}
	}
}

func (nn *NFA_PoolType) DumpPool(all bool) {
	if all {
		dbgo.DbPrintf("db_DumpPool", "Cur: %d Top: %d NextFree %d\n", nn.Cur, nn.Top, nn.NextFree)
	}
	dbgo.DbPrintf("db_DumpPool", "\n-------------------------------------- Modified for New Rule Order -----------------------------------------\n")
	dbgo.DbPrintf("db_DumpPool", "\nNFA InitState: %d\n\n", nn.InitState)
	pLnNo := dbgo.IsDbOn("db_NFA_LnNo")
	IfLnNo := func(s string) string {
		if pLnNo {
			t := fmt.Sprintf("[%3s]", s)
			return t
		}
		return ""
	}
	for ii, vv := range nn.Pool {
		if all || vv.IsUsed {
			dbgo.DbPrintf("db_DumpPool", "%3d%s: ", ii, IfLnNo(vv.LineNo))
			if vv.Rv > 0 {
				dbgo.DbPrintf("db_DumpPool", " T:%04d ", vv.Rv)
			} else {
				dbgo.DbPrintf("db_DumpPool", "        ")
			}
			// dbgo.DbPrintf("db_DumpPool", ` Edges: %s`, dbgo.SVar(vv.Next2))
			if dbgo.IsDbOn("db_DumpPool") {
				fmt.Printf("\t E:")
				for _, ww := range vv.Next2 {
					if ww.IsLambda {
						if ww.Is0ChMatch {
							fmt.Printf("{  \u03c4  %2d -> %2d  %s}  ", ww.From, ww.To, IfLnNo(ww.LineNo)) // Show a Tau(t) for a lambda that matchiens on else conditions.
						} else {
							fmt.Printf("{  \u03bb  %2d -> %2d  %s}  ", ww.From, ww.To, IfLnNo(ww.LineNo))
						}
						// fmt.Printf("{ (L) %2d -> %2d  %s}  ", ww.From, ww.To, IfLnNo(ww.LineNo))
					} else {
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
				if vv.Info.Action != 0 || vv.Info.MatchLength != 0 {
					fmt.Printf("\t\t\tNFA.Info: %s, PositionNumber:%d\n", DumpInfo(vv.Info), vv.TRuleMatch)
				}
				fmt.Printf("\n\n")
			}
		}
	}
}

func (nn *NFA_PoolType) DumpPoolJSON(fo io.Writer, td string, tn int) {
	fmt.Fprintf(fo, `{"Input":%q, "Rv":%d, "Start": %d, "Sigma":%q, "States":[%s`, td, tn, nn.InitState, nn.Sigma, "\n")
	for ii, vv := range nn.Pool {
		if vv.IsUsed {
			fmt.Fprintf(fo, ` { "Sn":%d, "Info":%s, `, ii, dbgo.SVar(vv.Info))
			if vv.Rv > 0 {
				fmt.Fprintf(fo, ` "Term":%d, `, vv.Rv)
			}
			fmt.Fprintf(fo, ` "Edge":[ `)
			cc := ""
			for _, ww := range vv.Next2 {
				if ww.IsLambda {
					if ww.Is0ChMatch {
						fmt.Fprintf(fo, "%s{ \"On\":\"\u03c4\", \"Fr\":%d, \"To\":%d }", cc, ww.From, ww.To)
					} else {
						fmt.Fprintf(fo, "%s{ \"On\":\"\u03bb\", \"Fr\":%d, \"To\":%d }", cc, ww.From, ww.To)
					}
					// fmt.Fprintf(fo, "%s{ \"On\":\"(L)\", \"Fr\":%d, \"To\":%d }", cc, ww.From, ww.To)
				} else {
					fmt.Fprintf(fo, "%s{ \"On\":\"%s\", \"Fr\":%d, \"To\":%d }", cc, ww.On, ww.From, ww.To)
				}
				cc = ", "
			}
			fmt.Fprintf(fo, "]}\n")
		}
	}
	fmt.Fprintf(fo, "]}\n")
}

// Set of possible input tokens
// Walk the NFA and collect all unique tokens that are not lambda and have a transition
func (nn *NFA_PoolType) GenerateSigma() (s string) {
	uniq := make(map[string]bool)
	s = ""
	for _, vv := range nn.Pool {
		if vv.IsUsed {
			for _, ww := range vv.Next2 {
				if !ww.IsLambda {
					rr, _ := utf8.DecodeRune([]byte(ww.On))
					uniq[string(rr)] = true
				}
			}

		}
	}

	dbgo.DbPrintf("nfa2", "NFA.GenerateSigma: uniq=%v, %s\n", uniq, dbgo.LF())

	// To store the keys in slice in sorted order
	var keys []string
	for k := range uniq {
		keys = append(keys, k)
	}
	dbgo.DbPrintf("nfa2", "NFA.GenerateSigma: keys (unsorted)=%v\n", keys)
	sort.Strings(keys)
	dbgo.DbPrintf("nfa2", "NFA.GenerateSigma: keys (sorted)=%v\n", keys)

	for _, k := range keys {
		s += k
	}

	return
}

func (nn *NFA_PoolType) UniqSigma() (rv string) {
	rv = nn.Sigma

	uniq := make(map[string]bool)
	for _, rn := range rv {
		uniq[string(rn)] = true
	}

	//fmt.Printf("NFA:UniqSigma: uniq=%v\n", uniq)

	// To store the keys in slice in sorted order
	var keys []string
	for k := range uniq {
		keys = append(keys, k)
	}
	//fmt.Printf("NFA:UniqSigma: keys (unsorted)=%v\n", keys)
	sort.Strings(keys)
	//fmt.Printf("NFA:UniqSigma: keys (sorted)=%v\n", keys)

	rv = ""
	for _, k := range keys {
		rv += k
	}

	// lr.Sigma = rv
	return
}

/*
digraph finite_state_machine {
	rankdir=LR;
	size="8,5"
	node [shape = doublecircle]; LR_0 LR_3 LR_4 LR_8;
	node [shape = circle];
	LR_0 -> LR_2 [ label = "SS(B)" ];
	LR_0 -> LR_1 [ label = "SS(S)" ];
	LR_1 -> LR_3 [ label = "S($end)" ];
	LR_2 -> LR_6 [ label = "SS(b)" ];
	LR_2 -> LR_5 [ label = "SS(a)" ];
	LR_2 -> LR_4 [ label = "S(A)" ];
	LR_5 -> LR_7 [ label = "S(b)" ];
	LR_5 -> LR_5 [ label = "S(a)" ];
	LR_6 -> LR_6 [ label = "S(b)" ];
	LR_6 -> LR_5 [ label = "S(a)" ];
	LR_7 -> LR_8 [ label = "S(b)" ];
	LR_7 -> LR_5 [ label = "S(a)" ];
	LR_8 -> LR_6 [ label = "S(b)" ];
	LR_8 -> LR_5 [ label = "S(a)" ];
}
*/

func (nn *NFA_PoolType) GenerateGVFile(fo io.Writer, td string, tn int) {
	// fmt.Fprintf(fo, `{"Input":%q, "Rv":%d, "Start": %d, "Sigma":%q, "States":[%s`, td, tn, nn.InitState, nn.GenerateSigma(), "\n")
	fmt.Fprintf(fo,
		`digraph finite_state_machine {
	rankdir=LR;
	size="8,5"
`)
	// size= for bigger graph - should be configurable with tests -

	var term []int
	for ii, vv := range nn.Pool {
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

	for _, vv := range nn.Pool {
		if vv.IsUsed {
			for _, ww := range vv.Next2 {
				if ww.IsLambda {
					if ww.Is0ChMatch {
						fmt.Fprintf(fo, "	s%d -> s%d [ label = \"%s\" ];\n", ww.From, ww.To, "\u03c4")
					} else {
						fmt.Fprintf(fo, "	s%d -> s%d [ label = \"%s\" ];\n", ww.From, ww.To, "\u03bb")
					}
				} else {
					fmt.Fprintf(fo, "	s%d -> s%d [ label = \"%s\" ];\n", ww.From, ww.To, re.EscapeStr(ww.On))
				}
			}
		}
	}
	fmt.Fprintf(fo, "}\n")
}

func (nn *NFA_PoolType) NoneVisited() {
	for ii, vv := range nn.Pool {
		if vv.IsUsed {
			nn.Pool[ii].Visited = false
		}
	}
}

// Given an initial set of startState, calculate the set of states that can be
// reached via a lambda (empty string).
func (nn *NFA_PoolType) LambdaClosure(startState []int) (setLambda []int) {
	nn.NoneVisited()
	return nn.lambdaClosureR(startState)
}

func (nn *NFA_PoolType) lambdaClosureR(startState []int) (setLambda []int) {
	// setLambda = setLambda[:0]
	for _, st := range startState {
		// fmt.Printf("StartState[%d]=%s, %s\n", st, dbgo.SVar(nn.Pool[st]), dbgo.LF())
		if !nn.Pool[st].Visited {
			// fmt.Printf("    ! visited, doing it, %s\n", dbgo.LF())
			nn.Pool[st].Visited = true
			vv := nn.Pool[st]
			if vv.IsUsed {
				// fmt.Printf("    ! IsUsed - that's good, %s, Edges Are:%s\n", dbgo.LF(), dbgo.SVar(vv.Next2))
				for _, ee := range vv.Next2 {
					if ee.IsLambda {
						setLambda = append(setLambda, ee.To)
						setLambda = append(setLambda, nn.lambdaClosureR([]int{ee.To})...)
					}
				}
			}
		}
	}
	return
}

// t1 := nn.LambdaClosureSet ( dfa_set, string(S) )
func (nn *NFA_PoolType) LambdaClosureSet(startState []int, S string) (setLambda []int) {
	for _, st := range startState {
		if nn.Pool[st].IsUsed {
			for _, vv := range nn.Pool[st].Next2 {
				if !vv.IsLambda && vv.On == S {
					setLambda = append(setLambda, vv.To)
				}
			}
		}
	}
	return
}

type NNPairType struct {
	StateSetIdx   int
	TRuleMatchVal int
	MatchLength   int
}

// ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
func IsInArray(n int, arr []int) bool {
	for _, v := range arr {
		if v == n {
			return true
		}
	}
	return false
}

// ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
func (nn *NFA_PoolType) HasTauEdge(StateSet []int) (Is0Ch bool) {
	Is0Ch = false
	for ii, vv := range nn.Pool {
		if nn.Pool[ii].IsUsed {
			for _, ee := range vv.Next2 {
				if ee.IsLambda && ee.Is0ChMatch && IsInArray(ee.To, StateSet) {
					// if "to" -> StateSet && StateSet [to] . Rv > 0
					if nn.Pool[ee.To].Rv > 0 {
						Is0Ch = true
						return
					}
				}
			}
		}
	}
	return
}

// ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
// Find the longest match length in the StateSet, with Rv > 0
// Find the item that is highes up(min value) on TRuleMatch
// Return that one

func (nn *NFA_PoolType) IsTerminalState(StateSet []int) (rv int, is bool, info InfoType, Is0Ch bool) {
	Is0Ch = false
	Is0Ch = nn.HasTauEdge(StateSet)
	dbgo.DbPrintf("nfa4", "IsTau-Term: StateSet[%v] = %v\n", StateSet, Is0Ch)
	// fmt.Printf("IsTermailal: Input %s\n", dbgo.SVar(StateSet))
	x := 0  // xyzzy - I don't think that this is really correct.  But it seems to work for now
	ns := 0 // xyzzy - I don't think that this is really correct.  But it seems to work for now
	set := make([]NNPairType, 0, len(StateSet))
	max_len := 0
	dbgo.DbPrintf("nfa3", "\nSet: (IsTerminalState - top) StateSet = %v\n", StateSet)
	for iii, v := range StateSet {
		if nn.Pool[v].Rv > 0 && nn.Pool[v].IsUsed {
			set = append(set, NNPairType{StateSetIdx: iii, TRuleMatchVal: nn.Pool[v].TRuleMatch, MatchLength: nn.Pool[v].Info.MatchLength})
			if max_len < nn.Pool[v].Info.MatchLength {
				max_len = nn.Pool[v].Info.MatchLength
			}
		}
	}
	dbgo.DbPrintf("nfa3", "\nSet: Set = %v, max_len = %d\n", set, max_len)

	// Find the Rv/Info with the lowest TRuleMatch Number (TRuleMatchVal)
	min_StateSetIdx := 0
	min_TRuleMatchVal := 999999999
	for i, v := range set {
		if v.TRuleMatchVal < min_TRuleMatchVal {
			min_TRuleMatchVal = v.TRuleMatchVal
			min_StateSetIdx = i
		}
	}
	dbgo.DbPrintf("nfa3", "min TRuleMatchVal = %d, at subscript %d, StateSet: %+v, %s\n", min_TRuleMatchVal, min_StateSetIdx, StateSet, dbgo.LF()) // Correct results at this point

	// Search each of the StateSetIdx's for the 1st Rv > 0
	for _, v := range StateSet {
		if nn.Pool[v].TRuleMatch == min_TRuleMatchVal {
			dbgo.DbPrintf("nfa3", "min_TRuleMatchVal == v == %d, match found for state %d, %s\n", min_TRuleMatchVal, v, dbgo.LF())
			x |= nn.Pool[v].Info.Action
			if nn.Pool[v].Info.NextState != 0 {
				ns = nn.Pool[v].Info.NextState
			}
			dbgo.DbPrintf("nfa3", "IsTermailal: Found at %d, value = %d\n", v, nn.Pool[v].Rv)
			return nn.Pool[v].Rv, true, nn.Pool[v].Info, Is0Ch
		}
	}
	dbgo.DbPrintf("nfa3", "At %s\n", dbgo.LF())

	return -1, false, InfoType{x, 0, ns, false, "", false}, Is0Ch
}

// ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
func (nn *NFA_PoolType) IsNonTerminalPushPopState(StateSet []int) (rv int, is bool, info InfoType, Is0Ch bool) {
	Is0Ch = false
	Is0Ch = nn.HasTauEdge(StateSet)
	dbgo.DbPrintf("nfa4", "IsTau-NonTerm: StateSet[%v] = %v\n", StateSet, Is0Ch)
	// fmt.Printf("IsTermailal: Input %s\n", dbgo.SVar(StateSet))
	x := 0  // xyzzy - I don't think that this is really correct.  But it seems to work for now
	ns := 0 // xyzzy - I don't think that this is really correct.  But it seems to work for now
	set := make([]NNPairType, 0, len(StateSet))
	max_len := 0
	dbgo.DbPrintf("nfa4", "\nSet: (IsNonTerminalPushPopState - top) StateSet = %v\n", StateSet)
	for iii, v := range StateSet {
		if nn.Pool[v].Rv == 0 && nn.Pool[v].Info.Action != 0 {
			set = append(set, NNPairType{StateSetIdx: iii, TRuleMatchVal: nn.Pool[v].TRuleMatch, MatchLength: nn.Pool[v].Info.MatchLength})
			if max_len < nn.Pool[v].Info.MatchLength {
				max_len = nn.Pool[v].Info.MatchLength
			}
		}
	}
	dbgo.DbPrintf("nfa4", "\nSet: Set = %v, max_len = %d\n", set, max_len)

	// Find the Rv/Info with the lowest TRuleMatch Number (TRuleMatchVal)
	min_StateSetIdx := 0
	min_TRuleMatchVal := 999999999
	for i, v := range set {
		if v.TRuleMatchVal < min_TRuleMatchVal {
			min_TRuleMatchVal = v.TRuleMatchVal
			min_StateSetIdx = i
		}
	}
	dbgo.DbPrintf("nfa4", "min TRuleMatchVal = %d, at subscript %d, StateSet: %+v, %s\n", min_TRuleMatchVal, min_StateSetIdx, StateSet, dbgo.LF()) // Correct results at this point

	// Search each of the StateSetIdx's for the 1st Rv > 0
	for _, v := range StateSet {
		if nn.Pool[v].TRuleMatch == min_TRuleMatchVal {
			dbgo.DbPrintf("nfa4", "min_TRuleMatchVal == v == %d, match found for state %d, %s\n", min_TRuleMatchVal, v, dbgo.LF())
			x |= nn.Pool[v].Info.Action
			if nn.Pool[v].Info.NextState != 0 {
				ns = nn.Pool[v].Info.NextState
			}
			dbgo.DbPrintf("nfa4", "IsNonTermailal: Found at %d, value = %d\n", v, nn.Pool[v].Rv)
			return nn.Pool[v].Rv, true, nn.Pool[v].Info, Is0Ch
		}
	}
	dbgo.DbPrintf("nfa4", "At %s\n", dbgo.LF())

	return -1, false, InfoType{x, 0, ns, false, "", false}, Is0Ch
}

/* vim: set noai ts=4 sw=4: */
