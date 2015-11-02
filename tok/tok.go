package tok

import (
	"fmt"
	"io"

	"github.com/pschlump/lexie/com"
)

// -----------------------------------------------------------------------------------------------------------------------------
type TokenType int

// type TokenNoType int

type Token struct {
	FileName     string      // What file are we reading from
	Typ          TokenType   //
	Match        string      // What matched
	Val          string      // Value of what matched (can be a replaced string)
	TokNo        int         // // TokNo    TokenNoType //
	LineNo       int         // Where did it start
	ColNo        int         //
	State        int         // What is the "S_State" from the scanner
	IsRepl       bool        //
	ReplStr      string      //
	DataType     int         //	../eval/lst.go CtxType_*
	CurValue     interface{} //
	Error        bool        // True if error
	ErrorMsg     string      // The message
	LValue       bool        //
	CreateId     bool        //
	CoceLocation string      // Locaiton in code to note (usually for errors)
}

// -----------------------------------------------------------------------------------------------------------------------------

type TokenListItem struct {
	Start  int
	End    int
	NHard  int
	AToken Token
}

type TokenList struct {
	TL          []TokenListItem
	EndIdx      map[int]int
	TokenData   []Token
	IgnoreToken int // IgnoreToken TokenNoType
}

// -----------------------------------------------------------------------------------------------------------------------------

//type TokenBuffer struct {
//	TokenData []Token
//}

func NewTokenList() *TokenList {
	return &TokenList{
		TL:          make([]TokenListItem, 0, 100),
		EndIdx:      make(map[int]int),
		TokenData:   make([]Token, 0, 100),
		IgnoreToken: -3,
	}
}

func (tl *TokenList) AddTokenToBuffer(A Token, Start int, End int, NHard int) {
	tl.TL = append(tl.TL, TokenListItem{AToken: A, Start: Start, End: End, NHard: NHard})
	tl.EndIdx[End] = len(tl.TL) - 1
}

//func (tl *TokenList) SetTokenLocation(line_no int, col_no int) {
//	if len(tl.TL) > 0 {
//		p := len(tl.TL) - 1
//		tl.TL.AToken[p].Line = line_no
//		tl.TL.AToken[p].Col = col_no
//	}
//}

// dfa.TokList.ReplaceToken ( dfa.MTab.Machine[ctx.St].Info.MatchLength, dfa.MTab.Machine[ctx.St].Info.ReplStr )
func (tl *TokenList) ReplaceToken(l int, s string) {
	ii := len(tl.TL) - 1
	lv := len(tl.TL[ii].AToken.Val)
	tl.TL[ii].AToken.IsRepl = true
	tl.TL[ii].AToken.ReplStr = s
	if lv-l >= 1 {
		tl.TL[ii].AToken.Val = tl.TL[ii].AToken.Val[0:lv-l] + s
	} else {
		com.DbPrintf("db_tok01", "Error: ReplaceToken has invalid data, %s\n", com.LF())
	}
	com.DbPrintf("db_tok01", "ReplaceToken: Match: ->%s<- Val: ->%s<-\n", tl.TL[ii].AToken.Match, tl.TL[ii].AToken.Val)
}

func (tl *TokenList) FlushTokenBuffer(TokStart int, isHard bool, atEof bool) {
	ok := true

	if !isHard {
		n := len(tl.TL) - 1
		if tl.TL[n].AToken.TokNo != tl.IgnoreToken {
			tl.TokenData = append(tl.TokenData, tl.TL[n].AToken)
		}
		tl.TL = tl.TL[:0]
		tl.EndIdx = make(map[int]int)
		return
	}

	if com.DbOn("db_FlushTokenBeforeBefore") {
		fmt.Printf("Before Flush, TokStart=%d, eof = %v\n", TokStart, atEof)
		tl.DumpTokenBuffer()
	}

	keepTok := make([]int, 0, len(tl.TL)) // Set of tokens to keepTokerve

	//com.DbPrintf("db_tok01","tl.EndIdx = %+v\n", tl.EndIdx)
	// Walk Backwards Creating List
	limit := 15
	for ii := len(tl.TL) - 1; ii >= 0 && limit > 0; limit-- {
		top_ii := ii
		com.DbPrintf("db_tok01", "At top of loop, ii = %d, tl.TL[ii]=%+v\n", ii, tl.TL[ii])
		keepTok = append(keepTok, ii)
		if ii <= 0 {
			break
		}
		end := tl.TL[ii].Start
		ii, ok = tl.EndIdx[end]
		if atEof {
			ii--
		}
		com.DbPrintf("db_tok01", "CCC ii=%d ok=%v end=%d atEof=%v\n", ii, ok, end, atEof)
		if !ok {
			if com.DbOn("OutputErrors") {
				fmt.Printf("Note: Failed to flush - invalid token set, ii=%v, end=%v, %s\n", ii, end, com.LF())
			}
			break
		}
		if end == 0 {
			fmt.Printf("Reached and end of 0\n")
			break
		}

		l := tl.TL[top_ii].End - tl.TL[top_ii].Start
		d := len(tl.TL[top_ii].AToken.Match)
		//com.DbPrintf("db_tok01","DDD top_ii=%d l=%d d=%d before ->%s<-\n", top_ii, l, d, tl.TL[top_ii].AToken.Match)
		if d-l >= 0 && d-l < d && top_ii >= 0 && top_ii < len(tl.TL) {
			tl.TL[top_ii].AToken.Match = tl.TL[top_ii].AToken.Match[d-l:]
		}
		//com.DbPrintf("db_tok01","EEE top_ii=%d l=%d d=%d after  ->%s<-\n", top_ii, l, d, tl.TL[top_ii].AToken.Match)
		if tl.TL[top_ii].AToken.IsRepl {
			tl.TL[top_ii].AToken.Val = tl.TL[top_ii].AToken.ReplStr
		} else {
			tl.TL[top_ii].AToken.Val = tl.TL[top_ii].AToken.Val[d-l:]
		}
	}

	last := len(tl.TokenData) - 1
	com.DbPrintf("db_tok01", "keepTok = %s, last=%d\n", com.SVar(keepTok), last)
	for ii := len(keepTok) - 1; ii >= 0; ii-- {
		vv := keepTok[ii]
		// xyzzy - if not an Ignore token then
		if tl.TL[vv].AToken.TokNo != tl.IgnoreToken {
			tl.TokenData = append(tl.TokenData, tl.TL[vv].AToken)
		}
	}

	// send EOF token -- Handeled in matcher better
	//if atEof {
	//	tl.TokenData = append(tl.TokenData, Token{TokNo: 37}) // tl.TokenData = append(tl.TokenData, Token{TokNo: Tok_EOF})
	//}

	tl.TL = tl.TL[:0]
	tl.EndIdx = make(map[int]int)

	if com.DbOn("db_FlushTokenBeforeAfter") {
		fmt.Printf("After Flush\n")
		fmt.Printf("--------------------------------------------------------\n")
		tl.DumpTokenBuffer()
		fmt.Printf("--------------------------------------------------------\n")
	}
}

/*
TokenBuffer: {
	"TokenData": [
		{
			"Filename": "",
			"Typ": 1,
			"Match": "abcd",
			"Val": "abcd",
			"TokNo": 12,
			"Line": 0,
			"Col": 0,
			"State": 0,
			"EndFilename": ""
		},
*/
func (tl *TokenList) DumpTokenBuffer() {
	// com.DbPrintf("db_tok01","TokenList: %s \n", com.SVarI(tl.TL))
	fmt.Printf("TokenList:\n")
	fmt.Printf("\t%3s %-5s %-5s %-4s %-5s %8s %-20s \n", "row", "Start", "End", "Hard", "TokNo", "sL/C", "Match")
	for ii, jj := range tl.TL {
		fmt.Printf("\t%3d %5d %5d %4d %5d %3d/%4d -->>%s<<--\n", ii, jj.Start, jj.End, jj.NHard, jj.AToken.TokNo, jj.AToken.LineNo, jj.AToken.ColNo, jj.AToken.Match)
	}
	fmt.Printf("\n")
	for ii, jj := range tl.EndIdx {
		fmt.Printf("\tEndIdx = %d -> %d\n", ii, jj)
	}
	fmt.Printf("TokenBuffer:\n")
	fmt.Printf("\t%3s %-5s %8s %-20s %s\n", "row", "TokNo", "sL/C", "Match", "Val")
	for ii, jj := range tl.TokenData {
		fmt.Printf("\t%3d %5d %3d/%4d -->>%s<<-- -->%s<-\n", ii, jj.TokNo, jj.LineNo, jj.ColNo, jj.Match, jj.Val)
		// , jj.Match)
	}
}

func (tl *TokenList) FDumpTokenBuffer(fo io.Writer, sc bool) {
	fmt.Fprintf(fo, "TokenBuffer:\n")
	if !sc {
		fmt.Fprintf(fo, "\t%3s %-5s %-20s\n", "row", "TokNo", "Match")
		for ii, jj := range tl.TokenData {
			fmt.Fprintf(fo, "\t%3d %5d -->>%s<<--\n", ii, jj.TokNo, jj.Match)
		}
	} else {
		fmt.Fprintf(fo, "\t%3s %-6s %-5s %-20s\n", "row", "TokTyp", "TokNo", "Match")
		for ii, jj := range tl.TokenData {
			fmt.Fprintf(fo, "\t%3d %6d %5d -->>%s<<--\n", ii, jj.Typ, jj.TokNo, jj.Match)
		}
	}
}
