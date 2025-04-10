package main

/*
------------------ ------------------ ------------------ ------------------ ------------------ ------------------ ------------------ ------------------ ------------------ ------------------

Ringo - the web framework for developers with deadlines that care about web scale performance
	(i.e. get it done on time and scalable to the real world)

------------------ ------------------ ------------------ ------------------ ------------------ ------------------ ------------------ ------------------ ------------------ ------------------

http://twig.sensiolabs.org/ -- A Django Clone in PHP - Popular - Do a Performace Compare

http://twig.sensiolabs.org/documentation -- List of templates and what they do		-- Also see twig.pdf

------------------ ------------------ ------------------ ------------------ ------------------ ------------------ ------------------ ------------------ ------------------ ------------------

877: func (mt *MtType) EvalExpr(Context *eval.ContextType, n, m int) bool {

Need to build a set of expr - of token arrays and symbol tables and have a set of test for them.
	eval [ aa, bb, cc ] -- How are vectors of data represented
	eval vv = [ aa, bb, cc ]
Need to pull expression eval into own set of code.

1. Overall
	1. Stabalize lexie
		2. Get multi-part test cases in dfa_* - have a set of them (array of input files) with table-output as .tst/.ref
		1. Figure out what is going on with DOT, NCCLs, [^](tau)
	2. Get all of ringo to work
		1. if/else, for, extend, template
		2. pipes
		3. fucntion interface - Data structure interface fucntions - piped.
	3. Get the Ringo/Data/Func/->Theme->HTTP2.0 server working

1. Eval simple expressions (abc=false)
2. Eval More complex expression (abc=(1+2)==3)
3. Test with IF
4. Read Data
5. Eval expression leading to a range
6. Test For Loop (1st case)

0. All tests for stuff in place and automated
	+1. Need to be able to go back to each section RE, NFA, DFA, Match and run same test at each level
	+2. Need to be able to add/verify each test at each level.
	+3. Need test for templates at final level

+0. Expression Eval, {{ .. }}
	+1. Get types from tokens in place ( gen.Tok_ID, gen.Tok_NUM, ... )
		1. Conversion -> int/float etc.
		2. Line numbers -> MtType
		3. -tmp- list out tokens in EvalExpr with TokNo/Match/Line to check
	+2. Identiy ID, String(+fix), Int, Float
	3. Get set of tokens echoed in expression parser
	+4. Get line-no,col-no-file-name from tokens

------------------ ------------------ ------------------ ------------------ ------------------ ------------------ ------------------ ------------------ ------------------ ------------------
*/

import (
	"encoding/json"
	"fmt"
	"html"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/pschlump/dbgo"
	"github.com/pschlump/lexie/dfa"
	"github.com/pschlump/lexie/eval"
	"github.com/pschlump/lexie/gen"
	"github.com/pschlump/lexie/st"
	"github.com/pschlump/lexie/tok"
	"github.com/pschlump/uuid"
)

const (
	Fx_null = iota
	Fx_csrf_token
	Fx_cycle
	Fx_debug
	Fx_include
	Fx_load
	Fx_lorem
	Fx_now
	Fx_regroup
	Fx_templatetag
	Fx_url
	Fx_with
	Fx_withratio
	Fx_verbatim
	Fx_endverbatim
	Fx_spaceless
	Fx_endspaceless
	Fx_filter
	Fx_endfilter
	Fx_comment
	Fx_endcomment
	Fx_autoescape
	Fx_endautoescape
	Fx_block
	Fx_endblock
	Fx_dump_context
	Fx_set_context
	Fx_get_context
	Fx_If
	Fx_IfEqual
	Fx_IfNotEqual
	Fx_IfChanged
	Fx_IfNotChanged
	Fx_ElsIf
	Fx_Else
	Fx_EndIf
	Fx_For
	Fx_Empty
	Fx_EndFor
	Fx_ReadJson
)

type DS_ParamType struct {
	ParamName    string
	DefaultValue string
	Value        string
}

type DS_TemplateType struct {
	Name   string
	Body   string         // or  parse tree
	Params []DS_ParamType // NParams == len(Params)
}

type MtType struct {
	NodeType    int         // gen.Tok_XXX  ( may be token or parsed Sub-tree-reduced )
	EscapeIt    bool        //
	ID          string      //
	DataType    int         // eval.CtxType_* - data type
	XValue      interface{} //	// All values
	SVal        []string    // Parameters to this node {% cycle v1  v2  v3 as bob silent %} 0=>v1, 1=>v2, 2=>v3
	TVal        []int       // TokNo - The Token Number
	TokVal      []tok.Token // TokNo - The Token Number
	Error       bool        // True if error
	ErrorMsg    string      // The message
	DataVal     interface{} // Data that is used in eval of this node.
	FxId        int         // If Fx - this is the ID of the function to call
	HTML_Output string      //	Output at end.
	List        []*MtType   // N-Ary tree of children
	LineNo      int         // Where did I come from
	ColNo       int         //
	FileName    string      //
	LValue      bool        // Is this an L-Value (assignable, chanable)
}

type FxType struct {
	Fx     func(callNo int, pt *Parse2Type, Context *eval.ContextType, curTree *MtType) (err error)
	FxName string
}

type Parse2Type struct {
	Cc       tok.Token
	St       *st.SymbolTable
	Lex      *dfa.Lexie
	FxLookup map[int]*FxType
	TheTree  *MtType
	Ctx      *eval.ContextType
	x_walk   func(mt *MtType, pos, depth int)
	pos      int
	depth    int
}

func NewMtType(t int, v string) (mt *MtType) {
	mt = &MtType{
		NodeType:    t,
		EscapeIt:    true,
		DataType:    eval.CtxType_Str,
		XValue:      v,
		HTML_Output: "",
		LineNo:      0,
		ColNo:       0,
	}
	return
}

func (mt *MtType) DumpMtType(fo io.Writer, pos, depth int) {
	if pos > 0 {
		fmt.Fprintf(fo, "%s[%3d]", strings.Repeat("    ", depth), pos)
	} else {
		fmt.Fprintf(fo, "%s", strings.Repeat("    ", depth))
	}
	fmt.Fprintf(fo, "NodeType=%d value=%v\n", mt.NodeType, mt.XValue)
	for ii, vv := range mt.List {
		vv.DumpMtType(fo, ii, depth+1)
	}
}

func NewParse2Type() (pt *Parse2Type) {
	pt = &Parse2Type{
		Cc:       tok.Token{},
		FxLookup: make(map[int]*FxType),
		Ctx:      eval.NewContextType(), // make(map[string]*ContextValueType),
	}
	pt.St = st.NewSymbolTable()

	ss := pt.St.DefineSymbol("csrf_token", "", []string{})
	ss.SymType = gen.Tok_Tree_Item
	ss.FxId = Fx_csrf_token
	pt.FxLookup[Fx_csrf_token] = &FxType{Fx: FxCsrf_token, FxName: "FxCsrf_token"}

	ss = pt.St.DefineSymbol("dump_context", "", []string{})
	ss.SymType = gen.Tok_Tree_Item
	ss.FxId = Fx_dump_context
	pt.FxLookup[Fx_dump_context] = &FxType{Fx: FxDump_context, FxName: "Fx_dump_context"}

	ss = pt.St.DefineSymbol("read_json", "", []string{})
	ss.SymType = gen.Tok_Tree_Item
	ss.FxId = Fx_ReadJson
	pt.FxLookup[Fx_ReadJson] = &FxType{Fx: FxReadJson, FxName: "Fx_ReadJson"}

	ss = pt.St.DefineSymbol("set_context", "", []string{})
	ss.SymType = gen.Tok_Tree_Item
	ss.FxId = Fx_set_context
	pt.FxLookup[Fx_set_context] = &FxType{Fx: FxSet_context, FxName: "Fx_set_context"}

	ss = pt.St.DefineSymbol("get_context", "", []string{})
	ss.SymType = gen.Tok_Tree_Item
	ss.FxId = Fx_get_context
	pt.FxLookup[Fx_get_context] = &FxType{Fx: FxGet_context, FxName: "Fx_get_context"}

	ss = pt.St.DefineSymbol("cycle", "", []string{})
	ss.SymType = gen.Tok_Tree_Item
	ss.FxId = Fx_cycle
	pt.FxLookup[Fx_cycle] = &FxType{Fx: FxCycle, FxName: "FxCycle"}

	ss = pt.St.DefineSymbol("debug", "", []string{})
	ss.SymType = gen.Tok_Tree_Item
	ss.FxId = Fx_debug

	ss = pt.St.DefineSymbol("include", "", []string{})
	ss.SymType = gen.Tok_Tree_Item
	ss.FxId = Fx_include

	ss = pt.St.DefineSymbol("load", "", []string{})
	ss.SymType = gen.Tok_Tree_Item
	ss.FxId = Fx_load

	ss = pt.St.DefineSymbol("lorem", "", []string{})
	ss.SymType = gen.Tok_Tree_Item
	ss.FxId = Fx_lorem

	ss = pt.St.DefineSymbol("now", "", []string{})
	ss.SymType = gen.Tok_Tree_Item
	ss.FxId = Fx_now

	ss = pt.St.DefineSymbol("regroup", "", []string{})
	ss.SymType = gen.Tok_Tree_Item
	ss.FxId = Fx_regroup

	ss = pt.St.DefineSymbol("templatetag", "", []string{})
	ss.SymType = gen.Tok_Tree_Item
	ss.FxId = Fx_templatetag

	ss = pt.St.DefineSymbol("url", "", []string{})
	ss.SymType = gen.Tok_Tree_Item
	ss.FxId = Fx_url

	ss = pt.St.DefineSymbol("with", "", []string{})
	ss.SymType = gen.Tok_Tree_Item
	ss.FxId = Fx_with

	ss = pt.St.DefineSymbol("withratio", "", []string{})
	ss.SymType = gen.Tok_Tree_Item
	ss.FxId = Fx_withratio

	ss = pt.St.DefineSymbol("if", "", []string{})
	ss.SymType = gen.Tok_Tree_If
	ss.FxId = Fx_If
	pt.FxLookup[Fx_If] = &FxType{Fx: FxIf, FxName: "FxIf"}
	ss = pt.St.DefineSymbol("ifequal", "", []string{})
	ss.SymType = gen.Tok_Tree_Ifequal
	ss.FxId = Fx_IfEqual
	pt.FxLookup[Fx_IfEqual] = &FxType{Fx: FxIfEqual, FxName: "FxIfEqual"}
	ss = pt.St.DefineSymbol("ifnotequal", "", []string{})
	ss.SymType = gen.Tok_Tree_Ifnotequal
	ss.FxId = Fx_IfNotEqual
	pt.FxLookup[Fx_IfNotEqual] = &FxType{Fx: FxIfNotEqual, FxName: "FxIfNotEqual"}
	ss = pt.St.DefineSymbol("ifchanged", "", []string{})
	ss.SymType = gen.Tok_Tree_Ifchanged
	ss.FxId = Fx_IfChanged
	pt.FxLookup[Fx_IfChanged] = &FxType{Fx: FxIfChanged, FxName: "FxIfChanged"}
	ss = pt.St.DefineSymbol("ifnotchanged", "", []string{})
	ss.SymType = gen.Tok_Tree_Ifnotchanged
	ss.FxId = Fx_IfNotChanged
	pt.FxLookup[Fx_IfNotChanged] = &FxType{Fx: FxIfNotChanged, FxName: "FxIfNotChanged"}

	ss = pt.St.DefineSymbol("elsif", "", []string{})
	ss.SymType = gen.Tok_Tree_ElsIf
	ss.FxId = Fx_ElsIf
	pt.FxLookup[Fx_ElsIf] = &FxType{Fx: FxElsIf, FxName: "FxElsIf"}
	ss = pt.St.DefineSymbol("elseif", "", []string{})
	ss.SymType = gen.Tok_Tree_ElsIf
	ss.FxId = Fx_ElsIf
	ss = pt.St.DefineSymbol("elif", "", []string{})
	ss.SymType = gen.Tok_Tree_ElsIf
	ss.FxId = Fx_ElsIf
	ss = pt.St.DefineSymbol("eif", "", []string{})
	ss.SymType = gen.Tok_Tree_ElsIf
	ss.FxId = Fx_ElsIf
	ss = pt.St.DefineSymbol("else", "", []string{})
	ss.SymType = gen.Tok_Tree_Else
	ss.FxId = Fx_Else
	pt.FxLookup[Fx_Else] = &FxType{Fx: FxElse, FxName: "FxElse"}

	ss = pt.St.DefineSymbol("endif", "", []string{})
	ss.SymType = gen.Tok_Tree_Endif
	ss.FxId = Fx_EndIf
	pt.FxLookup[Fx_EndIf] = &FxType{Fx: FxEndIf, FxName: "FxEndIf"}
	ss = pt.St.DefineSymbol("endifequal", "", []string{})
	ss.SymType = gen.Tok_Tree_Endif
	ss.FxId = Fx_EndIf
	ss = pt.St.DefineSymbol("endifnotequal", "", []string{})
	ss.SymType = gen.Tok_Tree_Endif
	ss.FxId = Fx_EndIf
	ss = pt.St.DefineSymbol("endifchanged", "", []string{})
	ss.SymType = gen.Tok_Tree_Endif
	ss.FxId = Fx_EndIf
	ss = pt.St.DefineSymbol("endifnotchanged", "", []string{})
	ss.SymType = gen.Tok_Tree_Endif
	ss.FxId = Fx_EndIf

	ss = pt.St.DefineSymbol("for", "", []string{})
	ss.SymType = gen.Tok_Tree_For
	ss.FxId = Fx_For
	pt.FxLookup[Fx_For] = &FxType{Fx: FxFor, FxName: "FxFor"}
	ss = pt.St.DefineSymbol("empty", "", []string{})
	ss.SymType = gen.Tok_Tree_Empty
	ss.FxId = Fx_Empty
	pt.FxLookup[Fx_Empty] = &FxType{Fx: FxEmpty, FxName: "FxEmpty"}
	ss = pt.St.DefineSymbol("endfor", "", []string{})
	ss.SymType = gen.Tok_Tree_EndFor
	ss.FxId = Fx_EndFor
	pt.FxLookup[Fx_EndFor] = &FxType{Fx: FxEndFor, FxName: "FxEndFor"}

	ss = pt.St.DefineSymbol("verbatim", "", []string{})
	ss.SymType = gen.Tok_Tree_Begin
	ss.FxId = Fx_verbatim

	ss = pt.St.DefineSymbol("endverbatim", "", []string{})
	ss.SymType = gen.Tok_Tree_End
	ss.FxId = Fx_endverbatim

	ss = pt.St.DefineSymbol("spaceless", "", []string{})
	ss.SymType = gen.Tok_Tree_Begin
	ss.FxId = Fx_spaceless

	ss = pt.St.DefineSymbol("endspaceless", "", []string{})
	ss.SymType = gen.Tok_Tree_End
	ss.FxId = Fx_endspaceless

	ss = pt.St.DefineSymbol("filter", "", []string{})
	ss.SymType = gen.Tok_Tree_Begin
	ss.FxId = Fx_filter

	ss = pt.St.DefineSymbol("endfilter", "", []string{})
	ss.SymType = gen.Tok_Tree_End
	ss.FxId = Fx_endfilter

	ss = pt.St.DefineSymbol("comment", "", []string{})
	ss.SymType = gen.Tok_Tree_Comment
	ss.FxId = Fx_comment

	ss = pt.St.DefineSymbol("endcomment", "", []string{})
	ss.SymType = gen.Tok_Tree_End
	ss.FxId = Fx_endcomment

	ss = pt.St.DefineSymbol("autoescape", "", []string{})
	ss.SymType = gen.Tok_Tree_Begin
	ss.FxId = Fx_autoescape
	pt.FxLookup[Fx_autoescape] = &FxType{Fx: FxAutoescape, FxName: "FxAutoescape"}

	ss = pt.St.DefineSymbol("endautoescape", "", []string{})
	ss.SymType = gen.Tok_Tree_End
	ss.FxId = Fx_endautoescape

	ss = pt.St.DefineSymbol("block", "", []string{})
	ss.SymType = gen.Tok_Tree_Begin
	ss.FxId = Fx_block

	ss = pt.St.DefineSymbol("endblock", "", []string{})
	ss.SymType = gen.Tok_Tree_End
	ss.FxId = Fx_endblock

	return
}

// xyzzy
// if t0, cv, f := pt.LookupReservedWord(tk.Match); f {
// func (pt *Parse2Type) GetToken() (tk tok.Token) {
// -------------------------------------------------------------------------------------------------------------------------------------------------------------------------
// Use func (st *SymbolTable) DefineReservedWord(name string, fxid int) (ss *SymbolType) { to define
func (pt *Parse2Type) LookupReservedWord(name string) (tk int, cv int, f bool) {
	tk, cv, f = 0, 0, false
	ss, err := pt.St.LookupSymbol(name)
	if err != nil {
		return
	}
	tk = ss.FxId
	f = true
	return
}

// -------------------------------------------------------------------------------------------------------------------------------------------------------------------------

func (pt *Parse2Type) GetToken() (tk tok.Token) {
	var t1 dfa.LexieChanelType
	t1 = <-pt.Lex.Message
	tk = t1.Token
	switch tk.TokNo {
	case gen.Tok_Str0: // = 43
		tk.Match = tk.Match[0 : len(tk.Match)-1] // Remove trailing " ' ` from string - could use to indicate string type.
		tk.DataType = eval.CtxType_Str
		tk.CurValue = tk.Match

	case gen.Tok_ID:
		switch tk.Match {
		case "true":
			fallthrough
		case "TRUE":
			tk.TokNo = gen.Tok_Tree_Bool
			tk.DataType = eval.CtxType_Bool
			tk.CurValue = true
		case "false":
			fallthrough
		case "FALSE":
			tk.TokNo = gen.Tok_Tree_Bool
			tk.DataType = eval.CtxType_Bool
			tk.CurValue = false
		default:
			if t0, cv, f := pt.LookupReservedWord(tk.Match); f {
				tk.TokNo = t0
				tk.CurValue = cv
			}
		}

	case gen.Tok_NUM:

		n, err := strconv.ParseInt(tk.Match, 0, 0)
		if err == nil {
			tk.Error = false
			tk.DataType = eval.CtxType_Int
			tk.CurValue = n
		} else {
			tk.Error = true
			tk.ErrorMsg = fmt.Sprintf("%s", err)
		}

		//	n, err := strconv.Atoi(tk.Match)
		//	if err == nil {
		//		tk.Error = false
		//		tk.DataType = eval.CtxType_Int
		//		tk.CurValue = n
		//	} else {
		//		tk.Error = true
		//		tk.ErrorMsg = fmt.Sprintf("%s", err)
		//	}

	case gen.Tok_Float:
		f, err := strconv.ParseFloat(tk.Match, 64)
		if err == nil {
			tk.Error = false
			tk.DataType = eval.CtxType_Float
			tk.CurValue = f
		} else {
			tk.Error = true
			tk.ErrorMsg = fmt.Sprintf("%s", err)
		}

	}
	return
}

func (pt *Parse2Type) ScanToEndMarker(TokEnd int, mt *MtType) {
	// pt.Cc = pt.GetToken()
	jj := 0
	for pt.Cc.TokNo != TokEnd {
		pt.Cc = pt.GetToken()
		if pt.Cc.TokNo != TokEnd {
			fmt.Printf("  Scan Across [ %d ] = %s\n", jj, pt.Cc.Match)
			mt.SVal = append(mt.SVal, pt.Cc.Match)
			mt.TVal = append(mt.TVal, pt.Cc.TokNo)
			mt.TokVal = append(mt.TokVal, pt.Cc)
		}
		jj++
	}
}

func (pt *Parse2Type) ScanToNextMarker() {
	for pt.Cc.TokNo != gen.Tok_CL_BL {
		pt.Cc = pt.GetToken()
	}
}

func (pt *Parse2Type) GenParseTree(depth int) (mt *MtType) {
	if mt == nil {
		mt = NewMtType(gen.Tok_Tree_List, "")
	}
	done := false
	fmt.Printf("At: %s\n", dbgo.LF())
	for !done {
		fmt.Printf("At: %s\n", dbgo.LF())
		pt.Cc = pt.GetToken()
		for pt.Cc.TokNo == gen.Tok_HTML {
			fmt.Printf("At: %s\n", dbgo.LF())
			mt.List = append(mt.List, NewMtType(gen.Tok_HTML, pt.Cc.Match))
			pt.Cc = pt.GetToken()
		}
		switch pt.Cc.TokNo {
		case gen.Tok_OP_BL: // Open Block, Tag, {%
			fmt.Printf("At: %s\n", dbgo.LF())
			pt.Cc = pt.GetToken()
			if pt.Cc.TokNo == gen.Tok_ID {
				fmt.Printf("******************** Lookup %s\n", pt.Cc.Match)
				sym, err := pt.St.LookupSymbol(pt.Cc.Match) // lookup and determine if it is an "Item" or a "Begin-Block" or a "End-Block"
				if err != nil {
					fmt.Printf("Error: Name Not Found ->%s<- in symbol table - invalid tag, Cc=%+v %s\n", pt.Cc.Match, pt.Cc, dbgo.LF())
					pt.ScanToNextMarker()
					// error - symbol not found - not defined
				} else {
					switch sym.SymType {
					case gen.Tok_Tree_Begin:
						fmt.Printf("------------------------------------------------------------------------\n")
						fmt.Printf("gen.Tok_Tree_Begin! Recursive Call Match ->%s<- At: %s\n", pt.Cc.Match, dbgo.LF())
						x := NewMtType(gen.Tok_Tree_Begin, pt.Cc.Match)
						x.FxId = sym.FxId
						pt.ScanToEndMarker(gen.Tok_CL_BL, x)              // scan across to %} - for entire begin-block/item
						x.List = append(x.List, pt.GenParseTree(depth+1)) // if block - then recursive call to "parse" // if end-block - then return
						mt.List = append(mt.List, x)
					case gen.Tok_Tree_End:
						fmt.Printf("At: %s\n", dbgo.LF())
						x := NewMtType(gen.Tok_Tree_End, pt.Cc.Match)
						x.FxId = sym.FxId
						mt.List = append(mt.List, x)
						pt.ScanToEndMarker(gen.Tok_CL_BL, x) // scan across to %} - for entire begin-block/item
						if depth != 0 {                      // if block - and name set - validate names match or warn - mis-matched names
							return
						}
					case gen.Tok_Tree_Item:
						fmt.Printf("At: %s\n", dbgo.LF())
						x := NewMtType(gen.Tok_Tree_Item, pt.Cc.Match)
						x.FxId = sym.FxId
						mt.List = append(mt.List, x)
						pt.ScanToEndMarker(gen.Tok_CL_BL, x) // scan across to %} - for entire begin-block/item

					case gen.Tok_Tree_If: //= 410
						fmt.Printf("------------------------------------------------------------------------\n")
						fmt.Printf("gen.Tok_Tree_If! Recursive Call Match ->%s<- At: %s\n", pt.Cc.Match, dbgo.LF())
						x := NewMtType(gen.Tok_Tree_If, pt.Cc.Match)
						x.FxId = sym.FxId
						pt.ScanToEndMarker(gen.Tok_CL_BL, x)              // scan across to %} - for entire begin-block/item
						x.List = append(x.List, pt.GenParseTree(depth+1)) // if block - then recursive call to "parse" // if end-block - then return
						mt.List = append(mt.List, x)
					case gen.Tok_Tree_Ifequal:
						fmt.Printf("------------------------------------------------------------------------\n")
						fmt.Printf("gen.Tok_Tree_IfEqual! Recursive Call Match ->%s<- At: %s\n", pt.Cc.Match, dbgo.LF())
						x := NewMtType(gen.Tok_Tree_Ifequal, pt.Cc.Match)
						x.FxId = sym.FxId
						pt.ScanToEndMarker(gen.Tok_CL_BL, x)              // scan across to %} - for entire begin-block/item
						x.List = append(x.List, pt.GenParseTree(depth+1)) // if block - then recursive call to "parse" // if end-block - then return
						mt.List = append(mt.List, x)
					case gen.Tok_Tree_Ifnotequal:
						fmt.Printf("------------------------------------------------------------------------\n")
						fmt.Printf("gen.Tok_Tree_IfEqual! Recursive Call Match ->%s<- At: %s\n", pt.Cc.Match, dbgo.LF())
						x := NewMtType(gen.Tok_Tree_Ifnotequal, pt.Cc.Match)
						x.FxId = sym.FxId
						pt.ScanToEndMarker(gen.Tok_CL_BL, x)              // scan across to %} - for entire begin-block/item
						x.List = append(x.List, pt.GenParseTree(depth+1)) // if block - then recursive call to "parse" // if end-block - then return
						mt.List = append(mt.List, x)
					case gen.Tok_Tree_ElsIf: //= 411
						fmt.Printf("At: %s\n", dbgo.LF())
						x := NewMtType(gen.Tok_Tree_ElsIf, pt.Cc.Match)
						x.FxId = sym.FxId
						mt.List = append(mt.List, x)
						pt.ScanToEndMarker(gen.Tok_CL_BL, x) // scan across to %} - for entire begin-block/item
					case gen.Tok_Tree_Else: //= 412
						fmt.Printf("At: %s\n", dbgo.LF())
						x := NewMtType(gen.Tok_Tree_Else, pt.Cc.Match)
						x.FxId = sym.FxId
						mt.List = append(mt.List, x)
						pt.ScanToEndMarker(gen.Tok_CL_BL, x) // scan across to %} - for entire begin-block/item
					case gen.Tok_Tree_Endif: //= 413
						fmt.Printf("At: %s\n", dbgo.LF())
						x := NewMtType(gen.Tok_Tree_Endif, pt.Cc.Match)
						x.FxId = sym.FxId
						mt.List = append(mt.List, x)
						pt.ScanToEndMarker(gen.Tok_CL_BL, x) // scan across to %} - for entire begin-block/item
						if depth != 0 {                      // if block - and name set - validate names match or warn - mis-matched names
							return
						}

					case gen.Tok_Tree_For:
						fmt.Printf("------------------------------------------------------------------------\n")
						fmt.Printf("gen.Tok_Tree_For! Recursive Call Match ->%s<- At: %s\n", pt.Cc.Match, dbgo.LF())
						x := NewMtType(gen.Tok_Tree_For, pt.Cc.Match)
						x.FxId = sym.FxId
						pt.ScanToEndMarker(gen.Tok_CL_BL, x)              // scan across to %} - for entire begin-block/item
						x.List = append(x.List, pt.GenParseTree(depth+1)) // if block - then recursive call to "parse" // if end-block - then return
						mt.List = append(mt.List, x)
					case gen.Tok_Tree_Empty:
						fmt.Printf("At: %s\n", dbgo.LF())
						x := NewMtType(gen.Tok_Tree_Empty, pt.Cc.Match)
						x.FxId = sym.FxId
						mt.List = append(mt.List, x)
						pt.ScanToEndMarker(gen.Tok_CL_BL, x) // scan across to %} - for entire begin-block/item
					case gen.Tok_Tree_EndFor:
						fmt.Printf("At: %s\n", dbgo.LF())
						x := NewMtType(gen.Tok_Tree_EndFor, pt.Cc.Match)
						x.FxId = sym.FxId
						mt.List = append(mt.List, x)
						pt.ScanToEndMarker(gen.Tok_CL_BL, x) // scan across to %} - for entire begin-block/item
						if depth != 0 {                      // if block - and name set - validate names match or warn - mis-matched names
							return
						}

					default:
						fmt.Printf("Error: Invalid SymbolTable.SymType=%d At: %s\n", sym.SymType, dbgo.LF())
					}
				}
			} else {
				fmt.Printf("Error: Tag must be followd by a name, %s/%d found instead, At: %s\n", pt.Cc.TokNo, pt.Cc.Match, dbgo.LF())
				// error
			}
		case gen.Tok_CL_BL: // Close Block, Tag, {%
			fmt.Printf("At: %s\n", dbgo.LF())
			//	// lookup and verify close block
			//	// scan across to %} - for entire begin-block/item
			if depth != 0 {
				return
			}
		case gen.Tok_OP_VAR:
			fmt.Printf("At: %s\n", dbgo.LF()) // Evaluate the VAR
		}
		fmt.Printf("At: %s\n", dbgo.LF())
		for pt.Cc.TokNo == gen.Tok_HTML {
			fmt.Printf("At: %s\n", dbgo.LF())
			mt.List = append(mt.List, NewMtType(gen.Tok_HTML, pt.Cc.Match))
			pt.Cc = pt.GetToken()
		}
		fmt.Printf("At: %s\n", dbgo.LF())
		if pt.Cc.TokNo == gen.Tok_EOF {
			fmt.Printf("At: %s\n", dbgo.LF())
			done = true
		}
	}
	return
}

// ---------------------------------------------------------------------------------------------------------------------------------------------
func EscapeStr(v string, on bool) string {
	if on {
		return html.EscapeString(v)
	} else {
		return v
	}
}

// ---------------------------------------------------------------------------------------------------------------------------------------------
// Depth first across tree - run functions as necessary or report missing Fx
// ---------------------------------------------------------------------------------------------------------------------------------------------
func (pt *Parse2Type) ExecuteFunctions(depth int) {
	var walkTreeInit func(mt *MtType, pos, depth int)
	var walkTreePass1 func(mt *MtType, pos, depth int)
	var walkTreePass2 func(mt *MtType, pos, depth int)

	walkTreeInit = func(mt *MtType, pos, depth int) {
		mt.EscapeIt = true
		mt.HTML_Output = ""
		mt.Error = false
		mt.ErrorMsg = ""
		for ii, vv := range mt.List {
			walkTreeInit(vv, ii, depth+1)
		}
	}

	walkTreePass1 = func(mt *MtType, pos, depth int) {
		switch mt.NodeType {
		case gen.Tok_Tree_Item: // = 406 // An item like {% csrf_token %}
			if x, ok := pt.FxLookup[mt.FxId]; ok {
				fmt.Printf("Found item %s to execute, %s\n", x.FxName, dbgo.LF())
				x.Fx(0, pt, pt.Ctx, mt)
			}
		case gen.Tok_Tree_Begin: // = 407 // An begin token like {% block <name> %} <name>==Value, ID="block"
			if x, ok := pt.FxLookup[mt.FxId]; ok {
				fmt.Printf("Found item %s to begin-block, %s\n", x.FxName, dbgo.LF())
				x.Fx(1, pt, pt.Ctx, mt)
			}
		case gen.Tok_Tree_End: // = 408 // An begin token like {% endblock <name> %} <name>==Value, ID="block"
			if x, ok := pt.FxLookup[mt.FxId]; ok {
				fmt.Printf("Found item %s to end-block, %s\n", x.FxName, dbgo.LF())
				x.Fx(2, pt, pt.Ctx, mt)
			}
		}
		for ii, vv := range mt.List {
			walkTreePass1(vv, ii, depth+1)
		}
	}

	walkTreePass2 = func(mt *MtType, pos, depth int) {
		switch mt.NodeType {
		case gen.Tok_Tree_If: //           = 410
			fmt.Printf("Found IF, %s\n", dbgo.LF())
			if x, ok := pt.FxLookup[mt.FxId]; ok {
				// fmt.Printf("Run item %s if-block, %s\n", x.FxName, dbgo.LF())
				pt.pos = pos
				pt.depth = depth + 1
				x.Fx(11, pt, pt.Ctx, mt) // pass pt, walkTreePass2, pos, depth
			}

		case gen.Tok_Tree_For:
			fmt.Printf("Found FOR, %s\n", dbgo.LF())
			if x, ok := pt.FxLookup[mt.FxId]; ok {
				pt.pos = pos
				pt.depth = depth + 1
				x.Fx(11, pt, pt.Ctx, mt) // pass pt, walkTreePass2, pos, depth
			}

		case gen.Tok_Tree_Comment:
			fmt.Printf("Found Comment, %s\n", dbgo.LF())
		case gen.Tok_Tree_Item: // = 406 // An item like {% csrf_token %}
			for ii, vv := range mt.List {
				walkTreePass2(vv, ii, depth+1)
			}
			if x, ok := pt.FxLookup[mt.FxId]; ok {
				fmt.Printf("Run item %s to execute, %s\n", x.FxName, dbgo.LF())
				x.Fx(10, pt, pt.Ctx, mt)
			}
		case gen.Tok_Tree_Begin: // = 407 // An begin token like {% block <name> %} <name>==Value, ID="block"
			for ii, vv := range mt.List {
				walkTreePass2(vv, ii, depth+1)
			}
			if x, ok := pt.FxLookup[mt.FxId]; ok {
				fmt.Printf("Run item %s to begin-block, %s\n", x.FxName, dbgo.LF())
				x.Fx(11, pt, pt.Ctx, mt)
			}
		case gen.Tok_Tree_End: // = 408 // An begin token like {% endblock <name> %} <name>==Value, ID="block"
			//for ii, vv := range mt.List {
			//	walkTreePass2(vv, ii, depth+1)
			//}
			if x, ok := pt.FxLookup[mt.FxId]; ok {
				fmt.Printf("Run item %s to end-block, %s\n", x.FxName, dbgo.LF())
				x.Fx(12, pt, pt.Ctx, mt)
			}
		default:
			for ii, vv := range mt.List {
				walkTreePass2(vv, ii, depth+1)
			}
			fmt.Printf("Run default, %s\n", dbgo.LF())
			mt.HTML_Output = EscapeStr(fmt.Sprintf("%s", mt.XValue), mt.EscapeIt)
		}
	}

	pt.x_walk = walkTreePass2

	walkTreeInit(pt.TheTree, 0, 0)
	walkTreePass1(pt.TheTree, 0, 0)
	walkTreePass2(pt.TheTree, 0, 0)
}

func (pt *Parse2Type) OutputTree(fo io.Writer, depth int) {
	var walkTree func(mt *MtType, pos, depth int)
	walkTree = func(mt *MtType, pos, depth int) {
		fmt.Fprintf(fo, "%s", mt.HTML_Output)
		// fmt.Printf("%s", mt.Value)
		for ii, vv := range mt.List {
			walkTree(vv, ii, depth+1)
		}
	}
	walkTree(pt.TheTree, 0, 0)
}

func (pt *Parse2Type) OutputTree0(fo io.Writer, depth int) {
	var walkTree func(mt *MtType, pos, depth int)
	walkTree = func(mt *MtType, pos, depth int) {
		// fmt.Fprintf(fo, "%s", mt.HTML_Output)
		for ii, vv := range mt.List {
			walkTree(vv, ii, depth+1)
		}
	}
	walkTree(pt.TheTree, 0, 0)
}

func (pt *Parse2Type) CollectErrorNodes(depth int) (rv []*MtType) {
	var walkTree func(mt *MtType, pos, depth int)
	walkTree = func(mt *MtType, pos, depth int) {
		if mt.Error {
			rv = append(rv, mt)
		}
		for ii, vv := range mt.List {
			walkTree(vv, ii, depth+1)
		}
	}
	if rv != nil {
		rv = rv[:0]
	}
	walkTree(pt.TheTree, 0, 0)
	return
}

func (pt *Parse2Type) CollectTree(mt *MtType, depth int) (rv string) {
	var walkTree func(mt *MtType, pos, depth int)
	walkTree = func(mt *MtType, pos, depth int) {
		rv += mt.HTML_Output
		for ii, vv := range mt.List {
			walkTree(vv, ii, depth+1)
		}
	}
	walkTree(pt.TheTree, 0, 0)
	return
}

// -------------------------------------------------------------------------------------------------------------------------------------------------------------------------
func (mt *MtType) NOptions(n int) bool {
	if len(mt.SVal) != n {
		mt.Error = true
		mt.ErrorMsg = fmt.Sprintf("Error: incorrect number of  options supplied - should have %d options, found %d, %s", n, len(mt.SVal), dbgo.LF(2))
		// fmt.Printf("Error: incorrect number of  options supplied - should have %d options, found %d (%v), %s", n, len(mt.SVal), mt.SVal, dbgo.LF(2))
		return false
	}
	return true
}

func (mt *MtType) MoreThan(n int) bool {
	if len(mt.SVal) <= n {
		mt.Error = true
		mt.ErrorMsg = fmt.Sprintf("Error: incorrect number of options supplied - should have more than %d options, %s", n, dbgo.LF(2))
		return false
	}
	return true
}

// Return True if expression is syntactically correct and convert it into type and value
// information in the mt node of the tree.  Convert the final expression to a "string"
// in HTML_Output if this is an expression evaluation node.  It the type is a boolean
// then set mt.BValue, if float then .FValue, if string then Value etc.   If the
// epression evaluates to a list then set .NValue to # of items in the set, and put
// the values into []AValue - for a simple set this is the set of values.
//
// Example
//  {{ = 1 + 2 }}
//  {{ id = 2 + 3 }}
//  {{ id = {{ select name, id from emp }} | sort:"1" }}

// n from beginning to m from end.
func (mt *MtType) EvalExpr(Context *eval.ContextType, n, m int) bool {
	m = len(mt.SVal) - m // Convert to Pos
	sv := mt.SVal[n:m]   // Slice of params to eval.
	fmt.Printf("mt.EvalExpr - TOP: ealuate sv=%+v ----------------------------------------------------- \n", sv)
	fmt.Printf("mt.EvalExpr - TOP: ealuate mt.SVal=%+v ----------------------------------------------------- \n", mt.SVal)
	// xyzzy -- temporary -- incomplete!!!!!!!!!!!!!!!!!!!!!!
	evalData := &eval.EvalType{
		Pos: 0,
		Ctx: Context,
		Mm:  mt.TokVal[n:m], // []tok.Token
	}
	fmt.Printf("INPUT m=%d n=%d, %s ----------------------------------------------------- \n", m, n, dbgo.SVarI(evalData))
	tr := evalData.Pres2()
	fmt.Printf("BOTTOM: %s ----------------------------------------------------- \n", dbgo.SVarI(tr))
	s := sv[0]
	v, t, _ := Context.GetFromContext(s)
	fmt.Printf("At: %s - in EvalExpr, v=%v t=%v for >%s<-\n", dbgo.LF(), v, t, s)
	// xyzzy nil, 9 -- 9 is error, not found
	if t == eval.CtxType_Bool {
		fmt.Printf("Setting bool to true\n")
		mt.DataType = t
		mt.XValue = v
	}
	return true
}

// Evaluate a set of contents or variables in 'in' - retuning an array of values in rv
func (mt *MtType) EvalVars(in []string) (rv []string) {
	rv = in
	// xyzzy
	return
}

func (mt *MtType) OptInList(pos int, s ...string) bool {
	for _, vv := range s {
		if mt.SVal[pos] == vv {
			return true
		}
	}
	mt.Error = true
	mt.ErrorMsg = fmt.Sprintf("Error: incorrect options supplied - should be one of %s", s)
	return false
}

// -------------------------------------------------------------------------------------------------------------------------------------------------------------------------
func FindTags(mt *MtType, tag ...int) (rv []int) {
	for ii, vv := range mt.List {
		for _, ww := range tag {
			if vv.NodeType == ww {
				rv = append(rv, ii)
			}
		}
	}
	return
}

// -------------------------------------------------------------------------------------------------------------------------------------------------------------------------
// -------------------------------------------------------------------------------------------------------------------------------------------------------------------------
//
//	Fx Funcitons
//
// -------------------------------------------------------------------------------------------------------------------------------------------------------------------------
// -------------------------------------------------------------------------------------------------------------------------------------------------------------------------
type FxCycleDataType struct {
	As     string
	Silent bool
	Opts   []string
	CurPos int
}

func FxCycle(callNo int, pt *Parse2Type, Context *eval.ContextType, curTree *MtType) (err error) {
	fmt.Printf("Cycle Called, %d\n", callNo)
	if callNo == 0 {
		if !curTree.MoreThan(2) {
		} else {
			x := &FxCycleDataType{CurPos: 0}
			ne := 0                                                    // Number at end we have processed.
			no := len(curTree.SVal)                                    // Number of Options
			if curTree.MoreThan(1) && curTree.SVal[no-1] == "silent" { // extract silent
				x.Silent = true
				ne = 1
			}
			if curTree.MoreThan(2+ne) && curTree.SVal[no-1-ne] == "as" { // extract as <id>
				x.As = curTree.SVal[no-ne-1]
				ne += 2
			}
			if curTree.MoreThan(2 + ne) {
				x.Opts = curTree.EvalVars(curTree.SVal[0 : no-ne]) // eval options -> values
				x.CurPos = 0                                       // establish data and position
			}
			curTree.DataVal = x
		}
	}
	if callNo == 10 {
		x := curTree.DataVal.(*FxCycleDataType)
		if x.Silent {
			curTree.HTML_Output = ""
		} else {
			curTree.HTML_Output = x.Opts[x.CurPos] // Return Value
		}
		if x.As != "" {
			Context.SetInContext(x.As, eval.CtxType_Str, x.Opts[x.CurPos]) // xyzzy - type may not be correct
		}
		x.CurPos = (x.CurPos + 1) % len(x.Opts) // Increment Postion Mod Length
		curTree.DataVal = x
	}
	return
}

// -------------------------------------------------------------------------------------------------------------------------------------------------------------------------
func FxDump_context(callNo int, pt *Parse2Type, Context *eval.ContextType, curTree *MtType) (err error) {
	fmt.Printf("DumpContext Called, %d\n", callNo)
	Context.DumpContext()
	return
}

// -------------------------------------------------------------------------------------------------------------------------------------------------------------------------
func FxSet_context(callNo int, pt *Parse2Type, Context *eval.ContextType, curTree *MtType) (err error) {
	fmt.Printf("SetContext Called, %d\n", callNo)
	if callNo == 10 {
		fmt.Printf("EvalExpr(context,0,0)=%v\n", curTree.EvalExpr(Context, 0, 0))
		if !curTree.MoreThan(2) {
		} else {
			id := curTree.SVal[0]
			val := curTree.SVal[1]
			if val == "true" || val == "TRUE" {
				Context.SetInContext(id, eval.CtxType_Bool, true) // xyzzy - type may not be correct
			} else {
				Context.SetInContext(id, eval.CtxType_Bool, false) // xyzzy - type may not be correct
			}
		}
	}
	return
}

// -------------------------------------------------------------------------------------------------------------------------------------------------------------------------
func FxGet_context(callNo int, pt *Parse2Type, Context *eval.ContextType, curTree *MtType) (err error) {
	fmt.Printf("GetContext Called, %d\n", callNo)
	if callNo == 10 {
		if !curTree.NOptions(1) {
		} else {
			id := curTree.SVal[0] // xyzzy - should be an EvalExpr of ...
			val, typ, fnd := Context.GetFromContext(id)
			if fnd {
				fmt.Printf("Found! id=%s typ=%d = ->%s<-\n", id, typ, val)
				if typ == eval.CtxType_Bool { // xyzzy - other types ...
					if val.(bool) {
						curTree.HTML_Output = "true"
					} else {
						curTree.HTML_Output = "false"
					}
				}
			} else {
				fmt.Printf("Not Found %s\n", id)
				curTree.HTML_Output = ""
			}
		}
	}
	return
}

// -------------------------------------------------------------------------------------------------------------------------------------------------------------------------
// consider adding "as" id silent to this.
// what should this actually do? - set a token in a header?
func FxCsrf_token(callNo int, pt *Parse2Type, Context *eval.ContextType, curTree *MtType) (err error) {
	fmt.Printf("Csrf_token Called, %d\n", callNo)
	if callNo == 10 {
		id, _ := uuid.NewV4()
		s_id := id.String()
		curTree.HTML_Output = s_id
	}
	return
}

// -------------------------------------------------------------------------------------------------------------------------------------------------------------------------
func FxAutoescape(callNo int, pt *Parse2Type, Context *eval.ContextType, curTree *MtType) (err error) {
	fmt.Printf("Autoescape Called, %d\n", callNo)
	var walkTree func(mt *MtType, pos, depth int)
	setTo := false
	walkTree = func(mt *MtType, pos, depth int) {
		mt.EscapeIt = setTo
		for ii, vv := range mt.List {
			walkTree(vv, ii, depth+1)
		}
	}
	if callNo == 1 {
		if !curTree.NOptions(1) {
		} else if !curTree.OptInList(0, "on", "off") {
		} else {
			setTo = (curTree.SVal[0] == "on")
			walkTree(curTree, 0, 0)
		}
	}
	return
}

// -------------------------------------------------------------------------------------------------------------------------------------------------------------------------
func FxIf(callNo int, pt *Parse2Type, Context *eval.ContextType, curTree *MtType) (err error) {
	fmt.Printf("Fx_If Called, %d\n", callNo)
	fmt.Printf("---------------------------------------------------------------------------- if tree -------------------------------------------------------------------------\n")
	if false {
		fmt.Printf("%s\n\n", dbgo.SVarI(curTree))
	} else {
		curTree.DumpMtType(os.Stdout, 0, 0)
	}

	tmpMt := func(ss []*MtType) (rv *MtType) {
		rv = &MtType{
			NodeType: gen.Tok_Tree_List,
			List:     make([]*MtType, 0, len(ss)),
			LineNo:   ss[0].LineNo,
			ColNo:    ss[0].ColNo,
			FileName: ss[0].FileName,
		}
		for _, vv := range ss {
			rv.List = append(rv.List, vv)
		}
		return
	}

	if callNo == 11 {
		fmt.Printf("n options = %d, opts = %v AT: %s\n", len(curTree.SVal), curTree.SVal, dbgo.LF())
		if !curTree.MoreThan(0) {
		} else {
			ifp := FindTags(curTree.List[0], gen.Tok_Tree_ElsIf, gen.Tok_Tree_Else, gen.Tok_Tree_Endif) // find parts of if/else
			fmt.Printf("ifp=%+v, 1st expr = %v\n", ifp, curTree.EvalExpr(Context, 0, 0))
			// xyzzy - should check order of ElsIf...Else...EndIf
			if curTree.EvalExpr(Context, 0, 0) {
				if curTree.DataType == eval.CtxType_Bool && curTree.XValue.(bool) {
					x := tmpMt(curTree.List[0].List[0:ifp[0]])
					pt.x_walk(x, pt.pos, pt.depth)
					return
				}
			}
			fmt.Printf("At AT: %s\n", dbgo.LF())
			for i := 0; i < len(ifp)-1; i++ {
				ct := curTree.List[0].List[ifp[i]]
				fmt.Printf("At AT: %s\n", dbgo.LF())
				if ct.NodeType == gen.Tok_Tree_ElsIf {
					fmt.Printf("At AT, it is (((ElsIf))): %s\n", dbgo.LF()) //
					if ct.EvalExpr(Context, 0, 0) {                         // expression is correct
						fmt.Printf("At AT: %s, ct=%+v\n", dbgo.LF(), ct)          //
						if ct.DataType == eval.CtxType_Bool && ct.XValue.(bool) { // If true value for expression
							x := tmpMt(curTree.List[0].List[ifp[i]+1 : ifp[i+1]])
							fmt.Printf("At -- Need to collect results -- AT: %s -------- elsif sub-tree Range[%d,%d] is %s\n", dbgo.LF(), ifp[i]+1, ifp[i+1], dbgo.SVarI(x))
							pt.x_walk(x, pt.pos, pt.depth)
							return
						}
					}
				} else if ct.NodeType == gen.Tok_Tree_Else {
					fmt.Printf("At AT, it is (((Else))): %s\n", dbgo.LF())
					x := tmpMt(curTree.List[0].List[ifp[i]+1 : ifp[i+1]])
					pt.x_walk(x, pt.pos, pt.depth)
					return
				}
			}
		}
	}
	return
}

// -------------------------------------------------------------------------------------------------------------------------------------------------------------------------
func FxElsIf(callNo int, pt *Parse2Type, Context *eval.ContextType, curTree *MtType) (err error) {
	fmt.Printf("Fx_ElsIf Called, %d - error\n", callNo)

	return
}

func FxElse(callNo int, pt *Parse2Type, Context *eval.ContextType, curTree *MtType) (err error) {
	fmt.Printf("Fx_Else Called, %d - error\n", callNo)
	return
}

func FxEndFor(callNo int, pt *Parse2Type, Context *eval.ContextType, curTree *MtType) (err error) {
	fmt.Printf("Fx_EndFor Called, %d - error\n", callNo)
	return
}

func FxEndIf(callNo int, pt *Parse2Type, Context *eval.ContextType, curTree *MtType) (err error) {
	fmt.Printf("Fx_EndIf Called, %d - error\n", callNo)
	return
}

// -------------------------------------------------------------------------------------------------------------------------------------------------------------------------
//
// Django
// 	{% for key, value in data.Items %}
// 	{% endfor %}
//
// Extended
// 	{% for key, value, pos in data.Items %}
// 	{% for key, _ in data.Items %}
// 	{% for _, _ in data.Items %}
// 	{% for _, value in data.Items %}
// 	{% for i = 0; i < x; i++ %}
// 	{% for i = 0; i < x; i += 4 %}
// 	{% for i = 0; i < x; %}
// 	{% for expr/call %}
// 	{% for true %}
//
// -------------------------------------------------------------------------------------------------------------------------------------------------------------------------

func FxFor(callNo int, pt *Parse2Type, Context *eval.ContextType, curTree *MtType) (err error) {
	fmt.Printf("Fx_For Called, %d\n", callNo)
	fmt.Printf("---------------------------------------------------------------------------- for tree -------------------------------------------------------------------------\n")
	if false {
		fmt.Printf("%s\n\n", dbgo.SVarI(curTree))
	} else {
		curTree.DumpMtType(os.Stdout, 0, 0)
	}

	tmpMt := func(ss []*MtType) (rv *MtType) {
		rv = &MtType{
			NodeType: gen.Tok_Tree_List,
			List:     make([]*MtType, 0, len(ss)),
			LineNo:   ss[0].LineNo,
			ColNo:    ss[0].ColNo,
			FileName: ss[0].FileName,
		}
		for _, vv := range ss {
			rv.List = append(rv.List, vv)
		}
		return
	}

	var walkTreeEmptyOutput func(mt *MtType, pos, depth int)
	walkTreeEmptyOutput = func(mt *MtType, pos, depth int) {
		mt.HTML_Output = ""
		for ii, vv := range mt.List {
			walkTreeEmptyOutput(vv, ii, depth+1)
		}
	}

	if callNo == 11 {
		if !curTree.MoreThan(1) {
		} else {
			ifp := FindTags(curTree.List[0], gen.Tok_Tree_Empty, gen.Tok_Tree_EndFor) // find parts of for loop
			fmt.Printf("ifp=%+v\n", ifp)
			if curTree.EvalExpr(Context, 0, 0) {
				x := tmpMt(curTree.List[0].List[0:ifp[0]])
				curTree.HTML_Output = ""
				// xyzzy - check type
				for ii, vv := range curTree.XValue.([]interface{}) {
					Context.SetInContext("$index", eval.CtxType_Int, ii) // xyzzy - conversion to string not correct -- needs to push $index - on endfor pop
					Context.SetInContext("$value", eval.CtxType_Str, vv) // xyzzy - conversion to string not correct
					//Context.SetInContext("key", fmt.Sprintf("%d", ii))    // xyzzy - conversion to string not correct	 -- key should be ID, Value too.
					//Context.SetInContext("value", fmt.Sprintf("%v", vv))  // xyzzy - conversion to string not correct
					pt.x_walk(x, pt.pos, pt.depth)
					curTree.HTML_Output += pt.CollectTree(x, 0) // Need to collect HTML_Output and append it to curTree.HTML_Output
				}
				mx := len(ifp)
				// xyzzy - check type
				if len(curTree.XValue.([]interface{})) == 0 && mx > 1 && curTree.List[0].List[ifp[mx-1]].NodeType == gen.Tok_Tree_Empty {
					i := mx - 1
					x := tmpMt(curTree.List[0].List[ifp[i]+1 : ifp[i+1]])
					pt.x_walk(x, pt.pos, pt.depth)
					curTree.HTML_Output += pt.CollectTree(x, 0) // Need to collect HTML_Output and append it to curTree.HTML_Output
				}
				walkTreeEmptyOutput(curTree.List[0], 0, 0) // set children's HTML_Output to ""
			}
		}
	}
	return
}

func FxEmpty(callNo int, pt *Parse2Type, Context *eval.ContextType, curTree *MtType) (err error) {
	fmt.Printf("Fx_Empty Called, %d - error\n", callNo)
	return
}

// ----------------------------------------------------------------------------------------------------------------------------------------
func FxReadJson(callNo int, pt *Parse2Type, Context *eval.ContextType, curTree *MtType) (err error) {
	fmt.Printf("Fx_ReadJson Called, %d\n", callNo)
	// {% read_json ID "file_name.json" %} (config allows url:// not just file"

	if callNo == 0 {

		if !curTree.NOptions(2) {
			// xyzzy } else if !curTree.IsId(0) {		// -- implement to check that [0] is an ID
		} else {
			id := curTree.SVal[0]
			path := curTree.SVal[1]
			// path = path[0 : len(path)-1]
			err = nil
			// var jsonData map[string]SQLOne
			var file []byte
			file, err = ioutil.ReadFile(path)
			if err != nil {
				fmt.Printf("Error(10014): %v, %s, Config File:%s\n", err, dbgo.LF(), path)
				return
			}
			file = []byte(strings.Replace(string(file), "\t", " ", -1)) // file = []byte(ReplaceString(string(file), "^[ \t][ \t]*//.*$", ""))

			// Check beginning of file if "{" then MapOf, if "[" Array, else look at single value
			if strings.HasPrefix(string(file), "{") {

				jsonData := make(map[string]interface{})

				err = json.Unmarshal(file, &jsonData)
				if err != nil {
					fmt.Printf("Error(10012): %v, %s, Config File:%s\n", err, dbgo.LF(), path)
					return
				}

				Context.SetInContext(id, eval.CtxType_MapOf, jsonData)

			} else {

				jsonData := make([]interface{}, 0, 100)

				err = json.Unmarshal(file, &jsonData)
				if err != nil {
					fmt.Printf("Error(10012): %v, %s, Config File:%s\n", err, dbgo.LF(), path)
					return
				}

				Context.SetInContext(id, eval.CtxType_ArrayOf, jsonData)

			}

		}
	}

	return
}

// -- Compares based on looping ------------------------------------------------------------------------------------------------------------

func FxIfChanged(callNo int, pt *Parse2Type, Context *eval.ContextType, curTree *MtType) (err error) {
	fmt.Printf("Fx_IfChanged Called, %d\n", callNo)
	// xyzzy

	return
}

func FxIfEqual(callNo int, pt *Parse2Type, Context *eval.ContextType, curTree *MtType) (err error) {
	fmt.Printf("Fx_IfEqual Called, %d\n", callNo)
	// xyzzy

	return
}

func FxIfNotChanged(callNo int, pt *Parse2Type, Context *eval.ContextType, curTree *MtType) (err error) {
	fmt.Printf("Fx_IfNotChanged Called, %d\n", callNo)
	// xyzzy

	return
}

func FxIfNotEqual(callNo int, pt *Parse2Type, Context *eval.ContextType, curTree *MtType) (err error) {
	fmt.Printf("Fx_IfNotEqual Called, %d\n", callNo)
	// xyzzy

	return
}

// --------------------------- Tempalte ------------------------------------------------------------------------------------------------------

// xyzzy - TODO

// --------------------------- Template File - name is basename(no-extention) on file --------------------------------------------------------

// xyzzy - TODO
