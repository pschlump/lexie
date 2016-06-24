package test01

//
// C L I / T E S T 2 - Part of Lexie Lexical Generation System
//
// (C) Philip Schlump, 2014-2015.
// Version: 0.3.0
// BuildNo: 51
//
// Special Thanks to 2C-Why, LLC for supporting this project.
//

/*

	0. Issues (../cli/cli.go dup)
			0. Is showing up as >>index.html"<< - note the quote at end - why - where?
			>>7. Figure out why "index.html" is not showing up as a template in the symbol-table

		{% extend "index.html" %} produces error - never finds %} - stuck on end of string I think
			*Locks up in channel communication
			*1. Put this in the ../dfa/dfa_test.go - code - something is wrong with nm.go on the "string"
			2. Test with other strings '...' etc.
			3. Test with mal-formed strings '..."
			4. Look at parser FxExtend and find out if it will take a "list" of ID or Str0
			5. Verify that Str0 is returned not Str1, Str2
			6. Figure out if this is a down-stream effect of ./tmpl/index.html
			8. Add a 2nd template inside ./tmpl/library.tpl - and verify in symbol table
			9. test with {% extend "header.html" "body.html" "footer.html" %}
			10. Test with "cli.go" on the same

	0. Use It
		1. Where is my template - base file for style
		2. How do the components pull together
		3. Let's do it and see it.


------------------ ------------------ ------------------ ------------------ ------------------ ------------------ ------------------ ------------------ ------------------
------- eob ------ ------- eob ------ ------- eob ------ ------- eob ------ ------- eob ------ ------- eob ------ ------- eob ------ ------- eob ------ ------- eob ------
------------------ ------------------ ------------------ ------------------ ------------------ ------------------ ------------------ ------------------ ------------------

	0. VarExists ( ... )
	1. SetVar - set var in context -- SetVar ( context, value, Id1, Id2 ... )
		SerVar ( context, 1, "prev", "v" )
		SerVar ( context, 1, "prev.v" )
	2. CompareVar ( v1, v2 ) -- deep compare of variables for values
	3. CompareVarContext ( Name1, Name2 ) in context
		"prev.V" automatically set at bottom of for loop - used by ifchanged, ifnotchanged
		IfEqual, 2 expression, IfNotEqual 2 expresions


	1. 4 days to finish 14th - !!!! use it !!!!
	1. 3 days to finish 15th -
			+ifequal, +ifnotequal, +ifchanged, +ifnotchanged
			+cycle (32), +library (35)
	1. 2 days to finish 16th - include, load, verbatim, filter, +template, +extend, +block
	1. 1 days to finish 17th -

	// -- later -- 00 --

	// xyzzy-FxIfChanged (partially completed)			00
	// xyzzy-FxIfEqual (partially completed)			02
	// xyzzy-FxIfNotChanged (partially completed)		03
	// xyzzy-FxIfNotEqual (partially completed)			04
	// xyzzy-FxCycle (mostly completed)					01

	// xyzzy-FxFilter									06
	// xyzzy-FxInclude									07
	// xyzzy-FxLoad										08
	// xyzzy-FxTemplateTag								12

------------------ ------------------ ------------------ ------------------ ------------------ ------------------ ------------------ ------------------ ------------------
------------------ ------------------ ------------------ ------------------ ------------------ ------------------ ------------------ ------------------ ------------------

	// xyzzy-FxVerbatim									05	- Set new scanner state, on {%endverbatim%} set return state		-- similar to original "comment"..."endcomment"

	// -- later -- 01 --

	// xyzzy-FxNow										10
	// xyzzy-FxRegroup									11
	// xyzzy-FxUrl										13
	// xyzzy-FxWith										14
	// xyzzy-FxWithRatio								15

Issues:
	FOR
		1. for fails to set forloop.* variables
		2. for a, b in x - not setting / updating a, b
			{% mtest [ ] %}
			{% mtest reverse %}

Plan:
	6. Test with nested for loops and forloop.parentloop


	Ya
		+1. SearchPath ( path ) -> List of Files
		2. Allow library to have list of paths as params

	+0. Dissect and implement lexie website using this system

	7. pre-processing macros {%name%} -> pre-process? {$...$}???

	3. HTTP2.0 --- Get the Ringo/Data/Func/->Theme->HTTP2.0 server working

------------------ ------------------ ------------------ ------------------ ------------------ ------------------ ------------------ ------------------ ------------------

*/

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
	"strings"

	"github.com/pschlump/uuid"

	"github.com/pschlump/lexie/com"
	"github.com/pschlump/lexie/mt"

	"github.com/pschlump/lexie/dfa"
	"github.com/pschlump/lexie/eval"
	"github.com/pschlump/lexie/gen"
	"github.com/pschlump/lexie/st"
	"github.com/pschlump/lexie/tok"
)

var Dbf *os.File

type FxFuncType func(callNo int, pt *Parse2Type, Context *eval.ContextType, curTree **mt.MtType) (err error)

type FxType struct {
	Fx     FxFuncType //
	FxName string     //
}

type Parse2Type struct {
	Cc           tok.Token                            //
	St           *st.SymbolTable                      //
	Lex          *dfa.Lexie                           //
	FxLookup     map[int]*FxType                      //
	TheTree      *mt.MtType                           //
	Ctx          *eval.ContextType                    //
	x_walk       func(mt **mt.MtType, pos, depth int) //
	pos          int                                  //
	depth        int                                  //
	LibraryMode  bool
	TemplateName string
}

// -------------------------------------------------------------------------------------------------------------------------------------------------------------------------
func (pt *Parse2Type) DefineTemplate(name string, body *mt.MtType) {
	name = com.BasenameExt(name)
	ss := pt.St.DefineSymbol(name, "", []string{})
	ss.SymType = gen.Tok_Template
	ss.AnyData = body
}

// -------------------------------------------------------------------------------------------------------------------------------------------------------------------------
func NewParse2Type() (pt *Parse2Type) {
	pt = &Parse2Type{
		Cc:       tok.Token{},
		FxLookup: make(map[int]*FxType),
		Ctx:      eval.NewContextType(), // make(map[string]*ContextValueType),
	}
	pt.St = st.NewSymbolTable()

	Def := func(name string, st int, fxid int, Fx FxFuncType, FxName string) {
		ss := pt.St.DefineSymbol(name, "", []string{})
		ss.SymType = st
		ss.FxId = fxid
		pt.FxLookup[fxid] = &FxType{Fx: Fx, FxName: FxName}
	}

	Def("library", gen.Tok_Tree_Item, gen.Fx_Library, FxLibrary, "FxLibrary")
	Def("csrf_token", gen.Tok_Tree_Item, gen.Fx_csrf_token, FxCsrf_token, "FxCsrf_token")
	Def("dump_context", gen.Tok_Tree_Item, gen.Fx_dump_context, FxDump_context, "Fx_dump_context")
	Def("read_json", gen.Tok_Tree_Item, gen.Fx_ReadJson, FxReadJson, "Fx_ReadJson")
	Def("set_context", gen.Tok_Tree_Item, gen.Fx_set_context, FxSet_context, "Fx_set_context")
	Def("get_context", gen.Tok_Tree_Item, gen.Fx_get_context, FxGet_context, "Fx_get_context")
	Def("cycle", gen.Tok_Tree_Item, gen.Fx_cycle, FxCycle, "FxCycle")
	Def("debug", gen.Tok_Tree_Item, gen.Fx_debug, FxDebug, "FxDebug")
	//Def("include", gen.Tok_Tree_Item, gen.Fx_include, FxInclude, "DFxInclude")
	//Def("load", gen.Tok_Tree_Item, gen.Fx_load, FxLoad, "FxLoad")
	Def("lorem", gen.Tok_Tree_Item, gen.Fx_lorem, FxLorem, "FxLorem")
	//Def("now", gen.Tok_Tree_Item, gen.Fx_now, FxNow, "FxNow")
	//Def("regroup", gen.Tok_Tree_Item, gen.Fx_regroup, FxRegroup, "FxRegroup")
	//Def("templatetag", gen.Tok_Tree_Item, gen.Fx_templatetag, FxTemplateTag, "FxTemplateTag")
	//Def("url", gen.Tok_Tree_Item, gen.Fx_url, FxUrl, "FxUrl")
	//Def("with", gen.Tok_Tree_Item, gen.Fx_with, FxWith, "FxWith")
	//Def("withratio", gen.Tok_Tree_Item, gen.Fx_withratio, FxWithRatio, "FxWithRatio")
	Def("if", gen.Tok_Tree_If, gen.Fx_If, FxIf, "FxIf")
	Def("ifequal", gen.Tok_Tree_Ifequal, gen.Fx_IfEqual, FxIfEqual, "FxIfEqual")
	Def("ifnotequal", gen.Tok_Tree_Ifnotequal, gen.Fx_IfNotEqual, FxIfNotEqual, "FxIfNotEqual")
	Def("ifchanged", gen.Tok_Tree_Ifchanged, gen.Fx_IfChanged, FxIfChanged, "FxIfChanged")
	Def("ifnotchanged", gen.Tok_Tree_Ifnotchanged, gen.Fx_IfNotChanged, FxIfNotChanged, "FxIfNotChanged")

	Def("elsif", gen.Tok_Tree_ElsIf, gen.Fx_ElsIf, FxEmpty, "FxElsIf")
	Def("elseif", gen.Tok_Tree_ElsIf, gen.Fx_ElsIf, FxEmpty, "FxElsIf")
	Def("elif", gen.Tok_Tree_ElsIf, gen.Fx_ElsIf, FxEmpty, "FxElsIf")
	Def("eif", gen.Tok_Tree_ElsIf, gen.Fx_ElsIf, FxEmpty, "FxElsIf")
	Def("else", gen.Tok_Tree_Else, gen.Fx_Else, FxEmpty, "FxElse")
	Def("endif", gen.Tok_Tree_Endif, gen.Fx_EndIf, FxEmpty, "FxEndIf")
	Def("endifequal", gen.Tok_Tree_Endif, gen.Fx_EndIf, FxEmpty, "FxEndIf")
	Def("endifnotequal", gen.Tok_Tree_Endif, gen.Fx_EndIf, FxEmpty, "FxEndIf")
	Def("endifchanged", gen.Tok_Tree_Endif, gen.Fx_EndIf, FxEmpty, "FxEndIf")
	Def("endifnotchanged", gen.Tok_Tree_Endif, gen.Fx_EndIf, FxEmpty, "FxEndIf")

	Def("for", gen.Tok_Tree_For, gen.Fx_For, FxFor, "FxFor")
	Def("empty", gen.Tok_Tree_Empty, gen.Fx_Empty, FxEmpty, "FxEmpty")

	Def("endfor", gen.Tok_Tree_EndFor, gen.Fx_EndFor, FxEmpty, "FxEndFor")

	//Def("verbatim", gen.Tok_Tree_Begin, gen.Fx_verbatim, FxVerbatim, "FxVerbatim")
	//Def("endverbatim", gen.Tok_Tree_End, gen.Fx_endverbatim, FxEmpty, "FxEndVerbatim")
	Def("spaceless", gen.Tok_Tree_Begin, gen.Fx_spaceless, FxSpaceless, "FxSpaceless")
	Def("endspaceless", gen.Tok_Tree_End, gen.Fx_endspaceless, FxEmpty, "FxEndSpaceless")
	//Def("filter", gen.Tok_Tree_Begin, gen.Fx_filter, FxFilter, "FxFilter")
	//Def("endfilter", gen.Tok_Tree_End, gen.Fx_endfilter, FxEmpty, "FxEndFilter")

	// Def("comment", gen.Tok_Tree_Comment, gen.Fx_comment, FxComment, "FxComment")
	Def("comment", gen.Tok_Tree_Begin, gen.Fx_comment, FxComment, "FxComment")
	Def("endcomment", gen.Tok_Tree_End, gen.Fx_endcomment, FxEndComment, "FxEndComment")
	Def("autoescape", gen.Tok_Tree_Begin, gen.Fx_autoescape, FxAutoescape, "FxAutoescape")
	Def("endautoescape", gen.Tok_Tree_End, gen.Fx_endautoescape, FxEmpty, "FxEndAutoEscape")
	Def("block", gen.Tok_Tree_Begin, gen.Fx_block, FxBlock, "FxBlock")
	Def("endblock", gen.Tok_Tree_End, gen.Fx_endblock, FxEmpty, "FxEndBlock")

	// "extend" <name>
	// "extend" <name> as <newname>
	Def("extend", gen.Tok_Tree_Begin, gen.Fx_extend, FxExtend, "FxExtend")
	Def("endextend", gen.Tok_Tree_End, gen.Fx_endextend, FxEmpty, "FxEndExtend")

	// "template" <name> ... "endtemplate"
	// "template" <name> "extend" <name> ... "endtemplate"
	Def("template", gen.Tok_Tree_Begin, gen.Fx_template, FxTemplate, "FxTemplate")
	Def("endtemplate", gen.Tok_Tree_End, gen.Fx_endtemplate, FxEmpty, "FxEndTemplate")

	Def("render", gen.Tok_Tree_Item, gen.Fx_render, FxRender, "FxRender")
	Def("mtest", gen.Tok_Tree_Item, gen.Fx_Mtest, FxMtest, "FxMtest")

	return
}

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

	fmt.Fprintf(Dbf, "Token %d ->%s<- At: %s\n", tk.TokNo, tk.Match, com.LF())

	switch tk.TokNo {
	case gen.Tok_Str0:
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
			fmt.Fprintf(Dbf, "At: %s\n", com.LF())
			if t0, cv, f := pt.LookupReservedWord(tk.Match); f {
				tk.TokNo = t0
				tk.CurValue = cv
				fmt.Fprintf(Dbf, "TokNo changed to %d At: %s\n", t0, com.LF())
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

	case gen.Tok_EOF:
		// panic("this error")
		return
	}
	return
}

// -------------------------------------------------------------------------------------------------------------------------------------------------------------------------
func (pt *Parse2Type) ScanToEndMarker(TokEnd int, mt *mt.MtType) {
	// pt.Cc = pt.GetToken()
	jj := 0
	for pt.Cc.TokNo != TokEnd && pt.Cc.TokNo != gen.Tok_EOF {
		pt.Cc = pt.GetToken()
		if pt.Cc.TokNo != TokEnd {
			fmt.Fprintf(Dbf, "  Scan Across [ %d ] = %s\n", jj, pt.Cc.Match)
			mt.SVal = append(mt.SVal, pt.Cc.Match)
			mt.TVal = append(mt.TVal, pt.Cc.TokNo)
			mt.TokVal = append(mt.TokVal, pt.Cc)
		}
		jj++
	}
}

// -------------------------------------------------------------------------------------------------------------------------------------------------------------------------
func (pt *Parse2Type) ScanToNextMarker() {
	for pt.Cc.TokNo != gen.Tok_CL_BL {
		fmt.Fprintf(Dbf, "At: %s\n", com.LF())
		pt.Cc = pt.GetToken()
	}
}

// -------------------------------------------------------------------------------------------------------------------------------------------------------------------------
func (pt *Parse2Type) GenParseTree(depth int) (mtv *mt.MtType) {
	if mtv == nil {
		mtv = mt.NewMtType(gen.Tok_Tree_List, "")
	}
	done := false
	fmt.Fprintf(Dbf, "At: %s\n", com.LF())
	for !done {
		fmt.Fprintf(Dbf, "At: %s\n", com.LF())
		pt.Cc = pt.GetToken()
		fmt.Fprintf(Dbf, "At: %s\n", com.LF())
		if pt.Cc.TokNo == gen.Tok_EOF {
			done = true
		}
		for pt.Cc.TokNo == gen.Tok_HTML {
			fmt.Fprintf(Dbf, "At: %s\n", com.LF())
			mtv.List = append(mtv.List, mt.NewMtType(gen.Tok_HTML, pt.Cc.Match))
			pt.Cc = pt.GetToken()
			if pt.Cc.TokNo == gen.Tok_EOF {
				done = true
			}
		}
		fmt.Fprintf(Dbf, "At: %s\n", com.LF())
		switch pt.Cc.TokNo {
		case gen.Tok_OP_BL: // Open Block, Tag, {%
			fmt.Fprintf(Dbf, "At: %s\n", com.LF())
			pt.Cc = pt.GetToken()
			if pt.Cc.TokNo == gen.Tok_ID || pt.Cc.TokNo >= 500 {
				fmt.Fprintf(Dbf, "******************** Lookup %s\n", pt.Cc.Match)
				sym, err := pt.St.LookupSymbol(pt.Cc.Match) // lookup and determine if it is an "Item" or a "Begin-Block" or a "End-Block"
				if err != nil {
					fmt.Fprintf(Dbf, "Error: Name Not Found ->%s<- in symbol table - invalid tag, Cc=%+v %s\n", pt.Cc.Match, pt.Cc, com.LF())
					pt.ScanToNextMarker()
					// error - symbol not found - not defined
				} else {
					switch sym.SymType {
					case gen.Tok_Tree_Begin:
						fmt.Fprintf(Dbf, "------------------------------------------------------------------------\n")
						fmt.Fprintf(Dbf, "gen.Tok_Tree_Begin! Recursive Call Match ->%s<- At: %s\n", pt.Cc.Match, com.LF())
						x := mt.NewMtType(gen.Tok_Tree_Begin, pt.Cc.Match)
						x.FxId = sym.FxId
						pt.ScanToEndMarker(gen.Tok_CL_BL, x)              // scan across to %} - for entire begin-block/item
						x.List = append(x.List, pt.GenParseTree(depth+1)) // if block - then recursive call to "parse" // if end-block - then return
						mtv.List = append(mtv.List, x)
					case gen.Tok_Tree_End:
						fmt.Fprintf(Dbf, "At: %s\n", com.LF())
						x := mt.NewMtType(gen.Tok_Tree_End, pt.Cc.Match)
						x.FxId = sym.FxId
						mtv.List = append(mtv.List, x)
						pt.ScanToEndMarker(gen.Tok_CL_BL, x) // scan across to %} - for entire begin-block/item
						if depth != 0 {                      // if block - and name set - validate names match or warn - mis-matched names
							return
						}
					case gen.Tok_Tree_Item:
						fmt.Fprintf(Dbf, "ITEM %s At: %s\n", pt.Cc.Match, com.LF())
						x := mt.NewMtType(gen.Tok_Tree_Item, pt.Cc.Match)
						x.FxId = sym.FxId
						fmt.Printf("FxId = %d\n", sym.FxId)
						mtv.List = append(mtv.List, x)
						pt.ScanToEndMarker(gen.Tok_CL_BL, x) // scan across to %} - for entire begin-block/item

					case gen.Tok_Tree_If: //= 410
						fmt.Fprintf(Dbf, "------------------------------------------------------------------------\n")
						fmt.Fprintf(Dbf, "gen.Tok_Tree_If! Recursive Call Match ->%s<- At: %s\n", pt.Cc.Match, com.LF())
						x := mt.NewMtType(gen.Tok_Tree_If, pt.Cc.Match)
						x.FxId = sym.FxId
						pt.ScanToEndMarker(gen.Tok_CL_BL, x)              // scan across to %} - for entire begin-block/item
						x.List = append(x.List, pt.GenParseTree(depth+1)) // if block - then recursive call to "parse" // if end-block - then return
						mtv.List = append(mtv.List, x)
					case gen.Tok_Tree_Ifequal:
						fmt.Fprintf(Dbf, "------------------------------------------------------------------------\n")
						fmt.Fprintf(Dbf, "gen.Tok_Tree_IfEqual! Recursive Call Match ->%s<- At: %s\n", pt.Cc.Match, com.LF())
						x := mt.NewMtType(gen.Tok_Tree_Ifequal, pt.Cc.Match)
						x.FxId = sym.FxId
						pt.ScanToEndMarker(gen.Tok_CL_BL, x)              // scan across to %} - for entire begin-block/item
						x.List = append(x.List, pt.GenParseTree(depth+1)) // if block - then recursive call to "parse" // if end-block - then return
						mtv.List = append(mtv.List, x)
					case gen.Tok_Tree_Ifnotequal:
						fmt.Fprintf(Dbf, "------------------------------------------------------------------------\n")
						fmt.Fprintf(Dbf, "gen.Tok_Tree_IfEqual! Recursive Call Match ->%s<- At: %s\n", pt.Cc.Match, com.LF())
						x := mt.NewMtType(gen.Tok_Tree_Ifnotequal, pt.Cc.Match)
						x.FxId = sym.FxId
						pt.ScanToEndMarker(gen.Tok_CL_BL, x)              // scan across to %} - for entire begin-block/item
						x.List = append(x.List, pt.GenParseTree(depth+1)) // if block - then recursive call to "parse" // if end-block - then return
						mtv.List = append(mtv.List, x)
					case gen.Tok_Tree_ElsIf: //= 411
						fmt.Fprintf(Dbf, "At: %s\n", com.LF())
						x := mt.NewMtType(gen.Tok_Tree_ElsIf, pt.Cc.Match)
						x.FxId = sym.FxId
						mtv.List = append(mtv.List, x)
						pt.ScanToEndMarker(gen.Tok_CL_BL, x) // scan across to %} - for entire begin-block/item
					case gen.Tok_Tree_Else: //= 412
						fmt.Fprintf(Dbf, "At: %s\n", com.LF())
						x := mt.NewMtType(gen.Tok_Tree_Else, pt.Cc.Match)
						x.FxId = sym.FxId
						mtv.List = append(mtv.List, x)
						pt.ScanToEndMarker(gen.Tok_CL_BL, x) // scan across to %} - for entire begin-block/item
					case gen.Tok_Tree_Endif: //= 413
						fmt.Fprintf(Dbf, "At: %s\n", com.LF())
						x := mt.NewMtType(gen.Tok_Tree_Endif, pt.Cc.Match)
						x.FxId = sym.FxId
						mtv.List = append(mtv.List, x)
						pt.ScanToEndMarker(gen.Tok_CL_BL, x) // scan across to %} - for entire begin-block/item
						if depth != 0 {                      // if block - and name set - validate names match or warn - mis-matched names
							return
						}

					case gen.Tok_Tree_For:
						fmt.Fprintf(Dbf, "------------------------------------------------------------------------\n")
						fmt.Fprintf(Dbf, "gen.Tok_Tree_For! Recursive Call Match ->%s<- At: %s\n", pt.Cc.Match, com.LF())
						x := mt.NewMtType(gen.Tok_Tree_For, pt.Cc.Match)
						x.FxId = sym.FxId
						pt.ScanToEndMarker(gen.Tok_CL_BL, x)              // scan across to %} - for entire begin-block/item
						x.List = append(x.List, pt.GenParseTree(depth+1)) // if block - then recursive call to "parse" // if end-block - then return
						mtv.List = append(mtv.List, x)
					case gen.Tok_Tree_Empty:
						fmt.Fprintf(Dbf, "At: %s\n", com.LF())
						x := mt.NewMtType(gen.Tok_Tree_Empty, pt.Cc.Match)
						x.FxId = sym.FxId
						mtv.List = append(mtv.List, x)
						pt.ScanToEndMarker(gen.Tok_CL_BL, x) // scan across to %} - for entire begin-block/item
					case gen.Tok_Tree_EndFor:
						fmt.Fprintf(Dbf, "At: %s\n", com.LF())
						x := mt.NewMtType(gen.Tok_Tree_EndFor, pt.Cc.Match)
						x.FxId = sym.FxId
						mtv.List = append(mtv.List, x)
						pt.ScanToEndMarker(gen.Tok_CL_BL, x) // scan across to %} - for entire begin-block/item
						if depth != 0 {                      // if block - and name set - validate names match or warn - mis-matched names
							return
						}
					case gen.Tok_EOF:
						done = true

					default:
						fmt.Fprintf(Dbf, "Error: Invalid SymbolTable.SymType=%d At: %s\n", sym.SymType, com.LF())
					}
				}
			} else {
				fmt.Fprintf(Dbf, "Error: Tag must be followd by a name, %d/%s found instead, At: %s\n", pt.Cc.TokNo, pt.Cc.Match, com.LF())
				// error
			}
		case gen.Tok_CL_BL: // Close Block, Tag, {%
			fmt.Fprintf(Dbf, "At: %s\n", com.LF())
			//	// lookup and verify close block
			//	// scan across to %} - for entire begin-block/item
			if depth != 0 {
				return
			}
		case gen.Tok_OP_VAR:
			fmt.Fprintf(Dbf, "At: %s\n", com.LF()) // Evaluate the VAR
		}
		fmt.Fprintf(Dbf, "At: %s\n", com.LF())
		for pt.Cc.TokNo == gen.Tok_HTML {
			fmt.Fprintf(Dbf, "At: %s\n", com.LF())
			mtv.List = append(mtv.List, mt.NewMtType(gen.Tok_HTML, pt.Cc.Match))
			pt.Cc = pt.GetToken()
		}
		fmt.Fprintf(Dbf, "At: %s\n", com.LF())
		if pt.Cc.TokNo == gen.Tok_EOF {
			fmt.Fprintf(Dbf, "At: %s\n", com.LF())
			done = true
		}
	}
	return
}

// -------------------------------------------------------------------------------------------------------------------------------------------------------------------------
// Depth first across tree - run functions as necessary or report missing Fx
func (pt *Parse2Type) ExecuteFunctions(depth int) {
	var walkTreeInit func(mtv **mt.MtType, pos, depth int)
	var walkTreePass1 func(mtv **mt.MtType, pos, depth int)
	var walkTreePass2 func(mtv **mt.MtType, pos, depth int)
	var walkTreePass3 func(mtv **mt.MtType, pos, depth int)

	walkTreeInit = func(mtv **mt.MtType, pos, depth int) {
		fmt.Printf("bob: %d %d %+v\n", pos, depth, mtv)
		(*mtv).EscapeIt = false
		(*mtv).HTML_Output = ""
		(*mtv).Error = false
		(*mtv).ErrorMsg = ""
		for ii, _ := range (*mtv).List {
			walkTreeInit(&((*mtv).List[ii]), ii, depth+1)
		}
	}

	walkTreePass1 = func(mtv **mt.MtType, pos, depth int) {
		switch (*mtv).NodeType {
		case gen.Tok_Tree_Item: // = 406 // An item like {% csrf_token %}
			if x, ok := pt.FxLookup[(*mtv).FxId]; ok {
				fmt.Fprintf(Dbf, "Found item %s to execute, %s\n", x.FxName, com.LF())
				x.Fx(0, pt, pt.Ctx, mtv)
			}
		case gen.Tok_Tree_Begin: // = 407 // An begin token like {% block <name> %} <name>==Value, ID="block"
			if x, ok := pt.FxLookup[(*mtv).FxId]; ok {
				fmt.Fprintf(Dbf, "Found item %s to begin-block, %s\n", x.FxName, com.LF())
				x.Fx(1, pt, pt.Ctx, mtv)
			}
		case gen.Tok_Tree_End: // = 408 // An begin token like {% endblock <name> %} <name>==Value, ID="block"
			if x, ok := pt.FxLookup[(*mtv).FxId]; ok {
				fmt.Fprintf(Dbf, "Found item %s to end-block, %s\n", x.FxName, com.LF())
				x.Fx(2, pt, pt.Ctx, mtv)
			}
		}
		for ii, _ := range (*mtv).List {
			walkTreePass1(&((*mtv).List[ii]), ii, depth+1)
		}
	}

	walkTreePass2 = func(mtv **mt.MtType, pos, depth int) {
		switch (*mtv).NodeType {
		case gen.Tok_Tree_If: //           = 410
			fmt.Fprintf(Dbf, "Found IF, %s\n", com.LF())
			if x, ok := pt.FxLookup[(*mtv).FxId]; ok {
				// fmt.Fprintf(Dbf,"Run item %s if-block, %s\n", x.FxName, com.LF())
				pt.pos = pos
				pt.depth = depth + 1
				x.Fx(11, pt, pt.Ctx, mtv) // pass pt, walkTreePass2, pos, depth
			}

		case gen.Tok_Tree_For:
			fmt.Fprintf(Dbf, "Found FOR, %s\n", com.LF())
			if x, ok := pt.FxLookup[(*mtv).FxId]; ok {
				pt.pos = pos
				pt.depth = depth + 1
				x.Fx(11, pt, pt.Ctx, mtv) // pass pt, walkTreePass2, pos, depth
			}

		//case gen.Tok_Tree_Comment:
		//	fmt.Fprintf(Dbf, "Found Comment, %s\n", com.LF())
		case gen.Tok_Tree_Item: // = 406 // An item like {% csrf_token %}
			fmt.Fprintf(Dbf, "At: %s\n", com.LF())
			for ii, _ := range (*mtv).List {
				walkTreePass2(&((*mtv).List[ii]), ii, depth+1)
			}
			fmt.Fprintf(Dbf, "At: %s\n", com.LF())
			if x, ok := pt.FxLookup[(*mtv).FxId]; ok {
				fmt.Fprintf(Dbf, "Run item %s to execute, %s\n", x.FxName, com.LF())
				x.Fx(10, pt, pt.Ctx, mtv)
			}
			fmt.Fprintf(Dbf, "At: %s\n", com.LF())
		case gen.Tok_Tree_Begin: // = 407 // An begin token like {% block <name> %} <name>==Value, ID="block"
			fmt.Fprintf(Dbf, "At: %s\n", com.LF())
			for ii, _ := range (*mtv).List {
				walkTreePass2(&((*mtv).List[ii]), ii, depth+1)
			}
			fmt.Fprintf(Dbf, "At: %s\n", com.LF())
			if x, ok := pt.FxLookup[(*mtv).FxId]; ok {
				fmt.Fprintf(Dbf, "********************************************** Run item %s to begin-block, %s\n", x.FxName, com.LF())
				x.Fx(11, pt, pt.Ctx, mtv)
			}
			fmt.Fprintf(Dbf, "At: %s\n", com.LF())
		case gen.Tok_Tree_End: // = 408 // An begin token like {% endblock <name> %} <name>==Value, ID="block"
			//for ii, vv := range mt.List {
			//	walkTreePass2(&((*mtv).List[ii]), ii, depth+1)
			//}
			fmt.Fprintf(Dbf, "At: %s\n", com.LF())
			if x, ok := pt.FxLookup[(*mtv).FxId]; ok {
				fmt.Fprintf(Dbf, "Run item %s to end-block, %s\n", x.FxName, com.LF())
				x.Fx(12, pt, pt.Ctx, mtv)
			}
			fmt.Fprintf(Dbf, "At: %s\n", com.LF())
		default:
			for ii, _ := range (*mtv).List {
				walkTreePass2(&((*mtv).List[ii]), ii, depth+1)
			}
			fmt.Fprintf(Dbf, "Run default, %s\n", com.LF())
			(*mtv).HTML_Output = com.EscapeStr(fmt.Sprintf("%s", (*mtv).XValue), (*mtv).EscapeIt)
		}
	}

	walkTreePass3 = func(mtv **mt.MtType, pos, depth int) {
		switch (*mtv).NodeType {
		case gen.Tok_Tree_If: //           = 410
			fmt.Fprintf(Dbf, "Found IF, %s\n", com.LF())
			if x, ok := pt.FxLookup[(*mtv).FxId]; ok {
				// fmt.Fprintf(Dbf,"Run item %s if-block, %s\n", x.FxName, com.LF())
				pt.pos = pos
				pt.depth = depth + 1
				x.Fx(21, pt, pt.Ctx, mtv) // pass pt, walkTreePass2, pos, depth
			}

		case gen.Tok_Tree_For:
			fmt.Fprintf(Dbf, "Found FOR, %s\n", com.LF())
			if x, ok := pt.FxLookup[(*mtv).FxId]; ok {
				pt.pos = pos
				pt.depth = depth + 1
				x.Fx(21, pt, pt.Ctx, mtv) // pass pt, walkTreePass2, pos, depth
			}

		case gen.Tok_Tree_Comment:
			fmt.Fprintf(Dbf, "Found Comment, %s\n", com.LF())
		case gen.Tok_Tree_Item: // = 406 // An item like {% csrf_token %}
			fmt.Fprintf(Dbf, "At: %s\n", com.LF())
			for ii, _ := range (*mtv).List {
				walkTreePass3(&((*mtv).List[ii]), ii, depth+1)
			}
			fmt.Fprintf(Dbf, "At: %s\n", com.LF())
			if x, ok := pt.FxLookup[(*mtv).FxId]; ok {
				fmt.Fprintf(Dbf, "Run item %s to execute, %s\n", x.FxName, com.LF())
				x.Fx(20, pt, pt.Ctx, mtv)
			}
			fmt.Fprintf(Dbf, "At: %s\n", com.LF())
		case gen.Tok_Tree_Begin: // = 407 // An begin token like {% block <name> %} <name>==Value, ID="block"
			fmt.Fprintf(Dbf, "At: %s\n", com.LF())
			for ii, _ := range (*mtv).List {
				walkTreePass3(&((*mtv).List[ii]), ii, depth+1)
			}
			fmt.Fprintf(Dbf, "At: %s\n", com.LF())
			if x, ok := pt.FxLookup[(*mtv).FxId]; ok {
				fmt.Fprintf(Dbf, "Run item %s to begin-block, %s\n", x.FxName, com.LF())
				x.Fx(21, pt, pt.Ctx, mtv)
			}
			fmt.Fprintf(Dbf, "At: %s\n", com.LF())
		case gen.Tok_Tree_End: // = 408 // An begin token like {% endblock <name> %} <name>==Value, ID="block"
			//for ii, vv := range mt.List {
			//	walkTreePass2(mtv.List[ii], ii, depth+1)
			//}
			fmt.Fprintf(Dbf, "At: %s\n", com.LF())
			if x, ok := pt.FxLookup[(*mtv).FxId]; ok {
				fmt.Fprintf(Dbf, "Run item %s to end-block, %s\n", x.FxName, com.LF())
				x.Fx(22, pt, pt.Ctx, mtv)
			}
			fmt.Fprintf(Dbf, "At: %s\n", com.LF())
		default:
			for ii, _ := range (*mtv).List {
				walkTreePass3(&((*mtv).List[ii]), ii, depth+1)
			}
			fmt.Fprintf(Dbf, "Run default, %s\n", com.LF())
			(*mtv).HTML_Output = com.EscapeStr(fmt.Sprintf("%s", (*mtv).XValue), (*mtv).EscapeIt)
		}
	}

	pt.x_walk = walkTreePass2

	fmt.Fprintf(Dbf, "Tree Is: %s\n", com.SVarI(pt.TheTree))

	if pt.TheTree == nil {
		return
	}

	fmt.Fprintf(Dbf, "\n\n-----------------------------------------------------------------------------------------------------------------------------\n")
	fmt.Fprintf(Dbf, "Before Pass 0 At: %s\n", com.LF())
	walkTreeInit(&pt.TheTree, 0, 0)
	fmt.Fprintf(Dbf, "\n\n-----------------------------------------------------------------------------------------------------------------------------\n")
	fmt.Fprintf(Dbf, "Before Pass 1 At: %s\n", com.LF())
	walkTreePass1(&pt.TheTree, 0, 0)
	fmt.Fprintf(Dbf, "\n\n-----------------------------------------------------------------------------------------------------------------------------\n")
	fmt.Fprintf(Dbf, "Before Pass 2 At: %s\n", com.LF())
	walkTreePass2(&pt.TheTree, 0, 0)
	fmt.Fprintf(Dbf, "\n\n-----------------------------------------------------------------------------------------------------------------------------\n")
	fmt.Fprintf(Dbf, "Before Pass 3 At: %s\n", com.LF())
	walkTreePass3(&pt.TheTree, 0, 0)
	fmt.Fprintf(Dbf, "All Done At: %s\n", com.LF())
}

// -------------------------------------------------------------------------------------------------------------------------------------------------------------------------
// This is the function that outputs the HTML
// -------------------------------------------------------------------------------------------------------------------------------------------------------------------------
func (pt *Parse2Type) OutputTree(fo io.Writer, depth int) {
	var walkTree func(mtv *mt.MtType, pos, depth int)
	walkTree = func(mtv *mt.MtType, pos, depth int) {
		fmt.Fprintf(fo, "%s", mtv.HTML_Output)
		// fmt.Fprintf(Dbf,"%s", mt.Value)
		for ii, _ := range mtv.List {
			walkTree(mtv.List[ii], ii, depth+1)
		}
	}
	fmt.Printf("------------------------------------- before output ------------------------------------- \n")
	walkTree(pt.TheTree, 0, 0)
	fmt.Printf("------------------------------------- after output ------------------------------------- \n")
}

// -------------------------------------------------------------------------------------------------------------------------------------------------------------------------
func (pt *Parse2Type) OutputTreeToString(ss string, depth int) {
	var walkTree func(mtv *mt.MtType, pos, depth int)
	walkTree = func(mtv *mt.MtType, pos, depth int) {
		ss += mtv.HTML_Output
		// fmt.Fprintf(Dbf,"%s", mt.Value)
		for ii, _ := range mtv.List {
			walkTree(mtv.List[ii], ii, depth+1)
		}
	}
	ss = ""
	walkTree(pt.TheTree, 0, 0)
}

// -------------------------------------------------------------------------------------------------------------------------------------------------------------------------
func (pt *Parse2Type) CollectErrorNodes(depth int) (rv []*mt.MtType) {
	var walkTree func(mtv *mt.MtType, pos, depth int)
	walkTree = func(mtv *mt.MtType, pos, depth int) {
		if mtv.Error {
			rv = append(rv, mtv)
		}
		for ii, _ := range mtv.List {
			walkTree(mtv.List[ii], ii, depth+1)
		}
	}
	if rv != nil {
		rv = rv[:0]
	}
	walkTree(pt.TheTree, 0, 0)
	return
}

// -------------------------------------------------------------------------------------------------------------------------------------------------------------------------
func (pt *Parse2Type) CollectTree(mtv *mt.MtType, depth int) (rv string) {
	var walkTree func(mtv *mt.MtType, pos, depth int)
	walkTree = func(mtv *mt.MtType, pos, depth int) {
		rv += mtv.HTML_Output
		for ii, _ := range mtv.List {
			walkTree(mtv.List[ii], ii, depth+1)
		}
	}
	walkTree(pt.TheTree, 0, 0)
	return
}

// -------------------------------------------------------------------------------------------------------------------------------------------------------------------------
// -------------------------------------------------------------------------------------------------------------------------------------------------------------------------
// -------------------------------------------------------------------------------------------------------------------------------------------------------------------------
//                             Fx Funcitons
// -------------------------------------------------------------------------------------------------------------------------------------------------------------------------
// -------------------------------------------------------------------------------------------------------------------------------------------------------------------------
// -------------------------------------------------------------------------------------------------------------------------------------------------------------------------
type FxCycleDataType struct {
	As     string
	Silent bool
	Opts   []string
	CurPos int
}

// xyzzy-FxCycle (mostly completed)			01
func FxCycle(callNo int, pt *Parse2Type, Context *eval.ContextType, curTree **mt.MtType) (err error) {
	fmt.Fprintf(Dbf, "Cycle Called, %d\n", callNo)
	if callNo == 0 {
		if !(*curTree).MoreThan(2) {
		} else {
			x := &FxCycleDataType{CurPos: 0}
			ne := 0                                                          // Number at end we have processed.
			no := len((*curTree).SVal)                                       // Number of Options
			if (*curTree).MoreThan(1) && (*curTree).SVal[no-1] == "silent" { // extract silent
				x.Silent = true
				ne = 1
			}
			if (*curTree).MoreThan(2+ne) && (*curTree).SVal[no-1-ne] == "as" { // extract as <id>
				x.As = (*curTree).SVal[no-ne-1]
				ne += 2
			}
			if (*curTree).MoreThan(2 + ne) {
				x.Opts = (*curTree).EvalVars((*curTree).SVal[0 : no-ne]) // eval options -> values
				x.CurPos = 0                                             // establish data and position
			}
			(*curTree).DataVal = x
		}
	}
	if callNo == 10 {
		x := (*curTree).DataVal.(*FxCycleDataType)
		if x.Silent {
			(*curTree).HTML_Output = ""
		} else {
			(*curTree).HTML_Output = x.Opts[x.CurPos] // Return Value
		}
		if x.As != "" {
			Context.SetInContext(x.As, eval.CtxType_Str, x.Opts[x.CurPos]) // xyzzy - type may not be correct
		}
		x.CurPos = (x.CurPos + 1) % len(x.Opts) // Increment Postion Mod Length
		(*curTree).DataVal = x
	}
	return
}

// -------------------------------------------------------------------------------------------------------------------------------------------------------------------------
func FxDump_context(callNo int, pt *Parse2Type, Context *eval.ContextType, curTree **mt.MtType) (err error) {
	fmt.Fprintf(Dbf, "DumpContext Called, %d\n", callNo)
	fmt.Printf("--- Context -------------------------------------------------------------\n")
	Context.DumpContext()
	fmt.Printf("--- Symbol Table -------------------------------------------------------------\n")
	// pt.St = st.NewSymbolTable()
	pt.St.DumpSymbolTable(os.Stdout)
	fmt.Printf("--- End -------------------------------------------------------------\n")
	return
}

// -------------------------------------------------------------------------------------------------------------------------------------------------------------------------
func FxSet_context(callNo int, pt *Parse2Type, Context *eval.ContextType, curTree **mt.MtType) (err error) {
	fmt.Fprintf(Dbf, "SetContext Called, %d\n", callNo)
	if callNo == 10 {
		fmt.Fprintf(Dbf, "EvalExpr(context,0,0)=%v\n", (*curTree).EvalExpr(Context, 0, 0))
		if !(*curTree).MoreThan(2) {
		} else {
			id := (*curTree).SVal[0]
			val := (*curTree).SVal[1]
			fmt.Printf("val=>>%s<<\n", val)
			// xyzzy - needs to call eval for functions not just true/false
			// xyzzy - needs to call eval for functions not just true/false
			// xyzzy - needs to call eval for functions not just true/false
			if val == "true" || val == "TRUE" || val == "True" {
				Context.SetInContext(id, eval.CtxType_Bool, true) // xyzzy - type may not be correct
			} else {
				Context.SetInContext(id, eval.CtxType_Bool, false) // xyzzy - type may not be correct
			}
		}
	}
	return
}

// -------------------------------------------------------------------------------------------------------------------------------------------------------------------------
func FxGet_context(callNo int, pt *Parse2Type, Context *eval.ContextType, curTree **mt.MtType) (err error) {
	fmt.Fprintf(Dbf, "GetContext Called, %d\n", callNo)
	if callNo == 10 {
		if !(*curTree).NOptions(1) {
		} else {
			id := (*curTree).SVal[0] // xyzzy - should be an EvalExpr of ...
			val, typ, fnd := Context.GetFromContext(id)
			if fnd {
				fmt.Fprintf(Dbf, "Found! id=%s typ=%d = ->%s<-\n", id, typ, val)
				// xyzzy - needs to call eval for functions not just true/false
				// xyzzy - needs to call eval for functions not just true/false
				// xyzzy - needs to call eval for functions not just true/false
				switch val.(type) {
				case bool:
					if val.(bool) {
						(*curTree).HTML_Output = "true"
					} else {
						(*curTree).HTML_Output = "false"
					}
				case int:
					(*curTree).HTML_Output = fmt.Sprintf("%d", val)
				case int64:
					(*curTree).HTML_Output = fmt.Sprintf("%d", val)
				case int32:
					(*curTree).HTML_Output = fmt.Sprintf("%d", val)
				case float64:
					(*curTree).HTML_Output = fmt.Sprintf("%f", val)
				case string:
					(*curTree).HTML_Output = val.(string)
				}
			} else {
				fmt.Fprintf(Dbf, "Not Found %s\n", id)
				(*curTree).HTML_Output = ""
			}
		}
	}
	return
}

// -------------------------------------------------------------------------------------------------------------------------------------------------------------------------
// consider adding "as" id silent to this.
// what should this actually do? - set a token in a header?
func FxCsrf_token(callNo int, pt *Parse2Type, Context *eval.ContextType, curTree **mt.MtType) (err error) {
	fmt.Fprintf(Dbf, "Csrf_token Called, %d\n", callNo)
	if callNo == 10 {
		if len((*curTree).SVal) > 0 && (*curTree).SVal[0] == "__test__" {
			s_id := "a954d701-4f31-46c0-75a1-59827ffbe530"
			(*curTree).HTML_Output = s_id
		} else {
			id, _ := uuid.NewV4()
			s_id := id.String()
			(*curTree).HTML_Output = s_id
		}
	}
	return
}

// -------------------------------------------------------------------------------------------------------------------------------------------------------------------------
func FxAutoescape(callNo int, pt *Parse2Type, Context *eval.ContextType, curTree **mt.MtType) (err error) {
	fmt.Fprintf(Dbf, "Autoescape Called, %d\n", callNo)
	var walkTree func(mtv *mt.MtType, pos, depth int)
	setTo := false
	walkTree = func(mtv *mt.MtType, pos, depth int) {
		mtv.EscapeIt = setTo
		for ii, _ := range mtv.List {
			walkTree(mtv.List[ii], ii, depth+1)
		}
	}
	if callNo == 1 {
		if !(*curTree).NOptions(1) {
		} else if !(*curTree).OptInList(0, "on", "off") {
		} else {
			setTo = ((*curTree).SVal[0] == "on")
			walkTree((*curTree), 0, 0)
		}
	}
	return
}

// -------------------------------------------------------------------------------------------------------------------------------------------------------------------------
func FxIf(callNo int, pt *Parse2Type, Context *eval.ContextType, curTree **mt.MtType) (err error) {
	com.DbFprintf("trace-builtin", Dbf, "Fx_If Called, %d\n", callNo)
	fmt.Fprintf(Dbf, "---------------------------------------------------------------------------- if tree -------------------------------------------------------------------------\n")
	(*curTree).DumpMtType(os.Stdout, 0, 0)

	tmpMt := func(ss []*mt.MtType) (rv *mt.MtType) {
		rv = &mt.MtType{
			NodeType: gen.Tok_Tree_List,
			List:     make([]*mt.MtType, 0, len(ss)),
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
		fmt.Fprintf(Dbf, "n options = %d, opts = %v AT: %s\n", len((*curTree).SVal), (*curTree).SVal, com.LF())
		if !(*curTree).MoreThan(0) {
		} else {
			ifp := mt.FindTags((*curTree).List[0], gen.Tok_Tree_ElsIf, gen.Tok_Tree_Else, gen.Tok_Tree_Endif) // find parts of if/else
			fmt.Fprintf(Dbf, "ifp=%+v, 1st expr = %v\n", ifp, (*curTree).EvalExpr(Context, 0, 0))
			// xyzzy - should check order of ElsIf...Else...EndIf
			if (*curTree).EvalExpr(Context, 0, 0) {
				if (*curTree).DataType == eval.CtxType_Bool && (*curTree).XValue.(bool) {
					x := tmpMt((*curTree).List[0].List[0:ifp[0]])
					pt.x_walk(&x, pt.pos, pt.depth) // xyzzy
					return
				}
			}
			fmt.Fprintf(Dbf, "At AT: %s\n", com.LF())
			for i := 0; i < len(ifp)-1; i++ {
				ct := (*curTree).List[0].List[ifp[i]]
				fmt.Fprintf(Dbf, "At AT: %s\n", com.LF())
				if ct.NodeType == gen.Tok_Tree_ElsIf {
					fmt.Fprintf(Dbf, "At AT, it is (((ElsIf))): %s\n", com.LF()) //
					if ct.EvalExpr(Context, 0, 0) {                              // expression is correct
						fmt.Fprintf(Dbf, "At AT: %s, ct=%+v\n", com.LF(), ct) //
						// -- xyzzy - use native .(type) and a switch
						if ct.DataType == eval.CtxType_Bool && ct.XValue.(bool) { // If true value for expression
							x := tmpMt((*curTree).List[0].List[ifp[i]+1 : ifp[i+1]])
							fmt.Fprintf(Dbf, "At -- Need to collect results -- AT: %s -------- elsif sub-tree Range[%d,%d] is %s\n", com.LF(), ifp[i]+1, ifp[i+1], com.SVarI(x))
							pt.x_walk(&x, pt.pos, pt.depth) // xyzzy
							return
						}
					}
				} else if ct.NodeType == gen.Tok_Tree_Else {
					fmt.Fprintf(Dbf, "At AT, it is (((Else))): %s\n", com.LF())
					x := tmpMt((*curTree).List[0].List[ifp[i]+1 : ifp[i+1]])
					pt.x_walk(&x, pt.pos, pt.depth) // xyzzy
					return
				}
			}
		}
	}
	return
}

// -------------------------------------------------------------------------------------------------------------------------------------------------------------------------
type AMatch struct {
	Type   int    // Type of this match
	AToken int    // Token Number
	AStr   string // String
}

type AMatchSet struct {
	Tm []AMatch
}

func MatchAtBeg(curTree *mt.MtType, pat []AMatchSet) (found bool, match int) {
	match = 0
	p := len(curTree.SVal) - 1
	fmt.Fprintf(Dbf, "At: %s\n", com.LF())
	for jj, ww := range pat {
		fmt.Fprintf(Dbf, "At: %s\n", com.LF())
		for ii, vv := range ww.Tm {
			if ii <= p {
				fmt.Fprintf(Dbf, "At: %s, p=%d ii=%d, vv.AToken=%d, vv.AStr=>>>%s<<<\n", com.LF(), p, ii, vv.AToken, vv.AStr)
			}
			if ii <= p &&
				((vv.AToken == gen.Tok_ID_or_Str && curTree.SVal[ii] == vv.AStr) ||
					(vv.AToken == gen.Tok_ID_or_Str && curTree.TVal[ii] == gen.Tok_ID) ||
					(vv.AToken == gen.Tok_Match_Str && curTree.SVal[ii] == vv.AStr) ||
					(vv.AToken == gen.Tok_Expr) ||
					(curTree.TVal[ii] == vv.AToken)) {
				fmt.Fprintf(Dbf, "At: %s - matched, loop on\n", com.LF())
			} else {
				fmt.Fprintf(Dbf, "At: %s - next\n", com.LF())
				goto next
			}
		}
		fmt.Fprintf(Dbf, "At: %s\n", com.LF())
		found = true
		match = jj
		return
	next:
	}
	fmt.Fprintf(Dbf, "At: %s\n", com.LF())
	return
}

func MatchAtEnd(curTree *mt.MtType, pat []AMatchSet) (found bool, match int) {
	match = 0
	p := len(curTree.SVal) - 1
	fmt.Fprintf(Dbf, "At: %s\n", com.LF())
	for jj, ww := range pat {
		fmt.Fprintf(Dbf, "At: %s\n", com.LF())
		lm := p - (len(ww.Tm) - 1)
		for ii, vv := range ww.Tm {
			fmt.Fprintf(Dbf, "At: %s, p=%d ii=%d, (lm+ii)=%d len(ww.Tm)=%d vv.AToken=%d, vv.AStr=>>>%s<<<\n", com.LF(), p, ii, lm+ii, len(ww.Tm), vv.AToken, vv.AStr)
			if lm+ii <= p && lm+ii >= 0 {
				fmt.Fprintf(Dbf, "... curTree.SVal[%d]->>>%s<<<\n", lm-ii, curTree.SVal[lm-ii])
			}
			fmt.Fprintf(Dbf, "At: %s\n", com.LF())
			if lm+ii <= p && lm+ii >= 0 && // (vv.AToken == gen.Tok_ID_or_Str && curTree.SVal[lm+ii] == vv.AStr) || (vv.AToken == gen.Tok_Expr) || (curTree.TVal[lm+ii] == vv.AToken) {
				((vv.AToken == gen.Tok_ID_or_Str && curTree.SVal[lm+ii] == vv.AStr) ||
					(vv.AToken == gen.Tok_ID_or_Str && curTree.TVal[lm+ii] == gen.Tok_ID) ||
					(vv.AToken == gen.Tok_Match_Str && curTree.SVal[lm+ii] == vv.AStr) ||
					(vv.AToken == gen.Tok_Expr) ||
					(curTree.TVal[ii] == vv.AToken)) {
				fmt.Fprintf(Dbf, "At: %s - matched, loop on\n", com.LF())
			} else {
				fmt.Fprintf(Dbf, "At: %s - next\n", com.LF())
				goto next
			}
		}
		fmt.Fprintf(Dbf, "At: %s\n", com.LF())
		found = true
		match = jj
		return
	next:
	}
	fmt.Fprintf(Dbf, "At: %s\n", com.LF())
	return
}

// -------------------------------------------------------------------------------------------------------------------------------------------------------------------------
func FxMtest(callNo int, pt *Parse2Type, Context *eval.ContextType, curTree **mt.MtType) (err error) {
	com.DbFprintf("trace-builtin", Dbf, "Fx_For Called, %d\n", callNo)

	if callNo == 10 {
		f, m := MatchAtBeg((*curTree), []AMatchSet{
			AMatchSet{Tm: []AMatch{
				AMatch{0, gen.Tok_ID, ""},
				AMatch{0, gen.Tok_COMMA, ""},
				AMatch{0, gen.Tok_ID, ""},
				AMatch{0, gen.Tok_ID_or_Str, "in"},
				AMatch{0, gen.Tok_Expr, ""},
			}},
			AMatchSet{Tm: []AMatch{
				AMatch{0, gen.Tok_ID, ""},
				AMatch{0, gen.Tok_ID_or_Str, "in"},
				AMatch{0, gen.Tok_Expr, ""},
			}},
			AMatchSet{Tm: []AMatch{
				AMatch{0, gen.Tok_Expr, ""},
			}},
		})
		fmt.Printf("\n* ====================================================================\n")
		fmt.Printf("* MatchAtBeg: Found=%v Loc=%d\n", f, m)
		fmt.Printf("* ====================================================================\n\n")
		f, m = MatchAtEnd((*curTree), []AMatchSet{
			AMatchSet{Tm: []AMatch{
				AMatch{0, gen.Tok_Expr, ""},
				AMatch{0, gen.Tok_Match_Str, "keysort"},
				AMatch{0, gen.Tok_Match_Str, "reverse"},
			}},
			AMatchSet{Tm: []AMatch{
				AMatch{0, gen.Tok_Expr, ""},
				AMatch{0, gen.Tok_Match_Str, "reverse"},
			}},
			AMatchSet{Tm: []AMatch{
				AMatch{0, gen.Tok_Expr, ""},
				AMatch{0, gen.Tok_Match_Str, "keysort"},
			}},
			AMatchSet{Tm: []AMatch{
				AMatch{0, gen.Tok_Expr, ""},
			}},
		})
		fmt.Printf("\n* ====================================================================\n")
		fmt.Printf("* MatchAtEnd: Found=%v Loc=%d\n", f, m)
		fmt.Printf("* ====================================================================\n\n")
	}

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
//		<table>
//			<tr> <th> Variable </th> 			<th> Value </th> 							<th> Description </th> </tr>
//
//			<tr><td>forloop.counter</td>		<td>{{ forloop.counter }}</td>				<td>The current iteration of the loop (1-indexed)</td> </tr>
//			<tr><td>forloop.counter0</td>		<td>{{ forloop.counter0 }}</td>				<td>The current iteration of the loop (0-indexed)</td> </tr>
//			<tr><td>forloop.revcounter</td>		<td>{{ forloop.revcounter }}</td>			<td>The number of iterations from the end of the loop (1-indexed)</td> </tr>
//			<tr><td>forloop.revcounter0</td>	<td>{{ forloop.revcounter0 }}</td>			<td>The number of iterations from the end of the loop (0-indexed)</td> </tr>
//			<tr><td>forloop.first</td>			<td>{{ forloop.first }}</td>				<td>True if this is the first time through the loop</td> </tr>
//			<tr><td>forloop.last</td>			<td>{{ forloop.last }}</td>					<td>True if this is the last time through the loop</td> </tr>
//			<tr><td>forloop.isempty</td>		<td>{{ forloop.isempty }}</td>				<td>True if this is the empty case in the for with no trips through the loop.</td> </tr>
//			<tr><td>forloop.parentloop</td>		<td>{{ forloop.parentloop }}</td>			<td>For nested loops, this is the loop surrounding the current one</td> </tr>
//
//		</table>
//
// -------------------------------------------------------------------------------------------------------------------------------------------------------------------------
//
// for Id , Vl in <expr>
// for Id , Vl in <expr> revserse
// for Id , Vl in <expr> keysort reverse
// for Id , Vl in <expr> keysort
//
// for _ , Vl in <expr>
// for _ , Vl in <expr> revserse
// for _ , Vl in <expr> keysort reverse
// for _ , Vl in <expr> keysort
//
// for Id in <expr>
// for Id in <expr> revserse
// for Id in <expr> keysort reverse
// for Id in <expr> keysort
//
// for _ in <expr>
// for _ in <expr> revserse
// for _ in <expr> keysort reverse
// for _ in <expr> keysort
//
// -------------------------------------------------------------------------------------------------------------------------------------------------------------------------

// -------------------------------------------------------------------------------------------------------------------------------------------------------------------------
func FxFor(callNo int, pt *Parse2Type, Context *eval.ContextType, curTree **mt.MtType) (err error) {
	com.DbFprintf("trace-builtin", Dbf, "Fx_For Called, %d\n", callNo)
	// fmt.Fprintf(Dbf, "---------------------------------------------------------------------------- for tree -------------------------------------------------------------------------\n")
	// if false {
	// 	fmt.Fprintf(Dbf, "%s\n\n", com.SVarI((*curTree)))
	// } else {
	// 	(*curTree).DumpMtType(os.Stdout, 0, 0)
	// }

	tmpMt := func(ss []*mt.MtType) (rv *mt.MtType) {
		rv = &mt.MtType{
			NodeType: gen.Tok_Tree_List,
			List:     make([]*mt.MtType, 0, len(ss)),
			LineNo:   ss[0].LineNo,
			ColNo:    ss[0].ColNo,
			FileName: ss[0].FileName,
		}
		for _, vv := range ss {
			rv.List = append(rv.List, vv)
		}
		return
	}

	var walkTreeEmptyOutput func(mtv *mt.MtType, pos, depth int)
	walkTreeEmptyOutput = func(mtv *mt.MtType, pos, depth int) {
		mtv.HTML_Output = ""
		for ii, _ := range mtv.List {
			walkTreeEmptyOutput(mtv.List[ii], ii, depth+1)
		}
	}

	var walkTreePass2 func(mtv **mt.MtType, pos, depth int)
	walkTreePass2 = func(mtv **mt.MtType, pos, depth int) {
		switch (*mtv).NodeType {
		case gen.Tok_Tree_If: //           = 410
			fmt.Fprintf(Dbf, "Found IF, %s\n", com.LF())
			if x, ok := pt.FxLookup[(*mtv).FxId]; ok {
				// fmt.Fprintf(Dbf,"Run item %s if-block, %s\n", x.FxName, com.LF())
				pt.pos = pos
				pt.depth = depth + 1
				x.Fx(11, pt, pt.Ctx, mtv) // pass pt, walkTreePass2, pos, depth
			}

		case gen.Tok_Tree_For:
			fmt.Fprintf(Dbf, "Found FOR, %s\n", com.LF())
			if x, ok := pt.FxLookup[(*mtv).FxId]; ok {
				pt.pos = pos
				pt.depth = depth + 1
				x.Fx(11, pt, pt.Ctx, mtv) // pass pt, walkTreePass2, pos, depth
			}

		//case gen.Tok_Tree_Comment:
		//	fmt.Fprintf(Dbf, "Found Comment, %s\n", com.LF())
		case gen.Tok_Tree_Item: // = 406 // An item like {% csrf_token %}
			fmt.Fprintf(Dbf, "At: %s\n", com.LF())
			for ii, _ := range (*mtv).List {
				walkTreePass2(&((*mtv).List[ii]), ii, depth+1)
			}
			fmt.Fprintf(Dbf, "At: %s\n", com.LF())
			if x, ok := pt.FxLookup[(*mtv).FxId]; ok {
				fmt.Fprintf(Dbf, "Run item %s to execute, %s\n", x.FxName, com.LF())
				x.Fx(10, pt, pt.Ctx, mtv)
			}
			fmt.Fprintf(Dbf, "At: %s\n", com.LF())
		case gen.Tok_Tree_Begin: // = 407 // An begin token like {% block <name> %} <name>==Value, ID="block"
			fmt.Fprintf(Dbf, "At: %s\n", com.LF())
			for ii, _ := range (*mtv).List {
				walkTreePass2(&((*mtv).List[ii]), ii, depth+1)
			}
			fmt.Fprintf(Dbf, "At: %s\n", com.LF())
			if x, ok := pt.FxLookup[(*mtv).FxId]; ok {
				fmt.Fprintf(Dbf, "Run item %s to begin-block, %s\n", x.FxName, com.LF())
				x.Fx(11, pt, pt.Ctx, mtv)
			}
			fmt.Fprintf(Dbf, "At: %s\n", com.LF())
		case gen.Tok_Tree_End: // = 408 // An begin token like {% endblock <name> %} <name>==Value, ID="block"
			//for ii, vv := range mt.List {
			//	walkTreePass2(&((*mtv).List[ii]), ii, depth+1)
			//}
			fmt.Fprintf(Dbf, "At: %s\n", com.LF())
			if x, ok := pt.FxLookup[(*mtv).FxId]; ok {
				fmt.Fprintf(Dbf, "Run item %s to end-block, %s\n", x.FxName, com.LF())
				x.Fx(12, pt, pt.Ctx, mtv)
			}
			fmt.Fprintf(Dbf, "At: %s\n", com.LF())
		default:
			for ii, _ := range (*mtv).List {
				walkTreePass2(&((*mtv).List[ii]), ii, depth+1)
			}
			fmt.Fprintf(Dbf, "Run default, %s\n", com.LF())
			(*mtv).HTML_Output = com.EscapeStr(fmt.Sprintf("%s", (*mtv).XValue), (*mtv).EscapeIt)
		}
	}
	pt.x_walk = walkTreePass2

	VMatchAtBeg := []AMatchSet{
		AMatchSet{Tm: []AMatch{
			AMatch{0, gen.Tok_ID, ""},
			AMatch{0, gen.Tok_COMMA, ""},
			AMatch{0, gen.Tok_ID, ""},
			AMatch{0, gen.Tok_ID_or_Str, "in"},
			AMatch{0, gen.Tok_Expr, ""},
		}},
		AMatchSet{Tm: []AMatch{
			AMatch{0, gen.Tok_ID, ""},
			AMatch{0, gen.Tok_ID_or_Str, "in"},
			AMatch{0, gen.Tok_Expr, ""},
		}},
		AMatchSet{Tm: []AMatch{
			AMatch{0, gen.Tok_Expr, ""},
		}},
	}

	VMatchAtEnd := []AMatchSet{
		AMatchSet{Tm: []AMatch{
			AMatch{0, gen.Tok_Expr, ""},
			AMatch{0, gen.Tok_Match_Str, "keysort"},
			AMatch{0, gen.Tok_Match_Str, "reverse"},
		}},
		AMatchSet{Tm: []AMatch{
			AMatch{0, gen.Tok_Expr, ""},
			AMatch{0, gen.Tok_Match_Str, "reverse"},
		}},
		AMatchSet{Tm: []AMatch{
			AMatch{0, gen.Tok_Expr, ""},
			AMatch{0, gen.Tok_Match_Str, "keysort"},
		}},
		AMatchSet{Tm: []AMatch{
			AMatch{0, gen.Tok_Expr, ""},
		}},
	}

	if callNo == 11 {
		if !(*curTree).MoreThan(1) {
		} else {
			fb, mb := MatchAtBeg((*curTree), VMatchAtBeg)
			// fmt.Printf("\n* ====================================================================\n")
			// fmt.Printf("* MatchAtBeg: Found=%v Loc=%d\n", fb, mb)
			// fmt.Printf("* ====================================================================\n\n")
			fe, me := MatchAtEnd((*curTree), VMatchAtEnd)
			// fmt.Printf("\n* ====================================================================\n")
			// fmt.Printf("* MatchAtEnd: Found=%v Loc=%d\n", fe, me)
			// fmt.Printf("* ====================================================================\n\n")
			ifp := mt.FindTags((*curTree).List[0], gen.Tok_Tree_Empty, gen.Tok_Tree_EndFor) // find parts of for loop
			// fmt.Printf("ifp = %+v\n", ifp)
			len_front := len(VMatchAtBeg[mb].Tm) - 1
			if !fb {
				len_front = 0
			}
			len_end := len(VMatchAtEnd[me].Tm) - 1
			if !fe {
				len_end = 0
			}
			// fmt.Fprintf(Dbf, "ifp=%+v, mb=%v,%d,%d me=%v,%d,%d\n", ifp, fb, mb, len(VMatchAtBeg[mb].Tm), fe, me, len(VMatchAtEnd[me].Tm))

			if (*curTree).EvalExpr(Context, len_front, len_end) {
				// fmt.Printf("FxFor 1 After Eval %+v\n", (*curTree))
				xx := tmpMt((*curTree).List[0].List[0:ifp[0]])
				(*curTree).HTML_Output = ""
				// fmt.Printf("FxFor: 0: Type=%T\n", (*curTree).XValue)
				switch (*curTree).XValue.(type) {
				case int:
					(*curTree).XValue = []interface{}{(*curTree).XValue.(int)}
				case int32:
					(*curTree).XValue = []interface{}{int((*curTree).XValue.(int32))}
				case int64:
					(*curTree).XValue = []interface{}{int((*curTree).XValue.(int32))}
				case float32:
					(*curTree).XValue = []interface{}{float64((*curTree).XValue.(float32))}
				case float64:
					(*curTree).XValue = []interface{}{(*curTree).XValue.(float64)}
				case []tok.Token:
					y := make([]interface{}, 0, len((*curTree).XValue.([]tok.Token)))
					for _, vv := range (*curTree).XValue.([]tok.Token) {
						y = append(y, vv.CurValue)
					}
					(*curTree).XValue = y
					// fmt.Printf("Len of loop = %d\n", len(y))
				}
				loop_occured := false
				tt := ""
				for ii, vv := range (*curTree).XValue.([]interface{}) {
					loop_occured = true
					Context.SetInContext("$index", eval.CtxType_Int, ii) // xyzzy - conversion to string not correct -- needs to push $index - on endfor pop
					Context.SetInContext("$value", eval.CtxType_Str, vv) // xyzzy - conversion to string not correct
					//Context.SetInContext("key", fmt.Sprintf("%d", ii))    // xyzzy - conversion to string not correct	 -- key should be ID, Value too.
					//Context.SetInContext("value", fmt.Sprintf("%v", vv))  // xyzzy - conversion to string not correct
					pt.x_walk(&xx, pt.pos, pt.depth) // xyzzy
					ss := pt.CollectTree(xx, 0)      // Need to collect HTML_Output and append it to (*curTree).HTML_Output
					// fmt.Printf("Inside loop, iteration %d ->%s<-\n", ii, ss)
					tt += ss
				}
				mx := len(ifp)
				if loop_occured {
					(*curTree).HTML_Output = tt
					// fmt.Printf("HTML_Output after loop = ---->>>>%s<<<<----\n", (*curTree).HTML_Output)
				} else if mx > 1 {
					// fmt.Printf("mx=%d At:%s\n", mx, com.LF())
					for kk := 0; kk < mx; kk++ {
						if kk+1 < mx && (*curTree).List[0].List[ifp[kk]].NodeType == gen.Tok_Tree_Empty && (*curTree).List[0].List[ifp[kk+1]].NodeType == gen.Tok_Tree_EndFor && ifp[kk]+1 <= ifp[kk+1]-1 {
							// fmt.Printf("Inside NO loop, i=%d\n", kk)
							zz := tmpMt((*curTree).List[0].List[ifp[kk]+1 : ifp[kk+1]])
							pt.x_walk(&zz, pt.pos, pt.depth)               // xyzzy
							(*curTree).HTML_Output = pt.CollectTree(zz, 0) // Need to collect HTML_Output and append it to (*curTree).HTML_Output
							kk = mx                                        // exit loop
						}
					}
				}
				walkTreeEmptyOutput((*curTree).List[0], 0, 0) // set children's HTML_Output to ""
			}
		}
	}
	return
}

// -------------------------------------------------------------------------------------------------------------------------------------------------------------------------
func FxEmpty(callNo int, pt *Parse2Type, Context *eval.ContextType, curTree **mt.MtType) (err error) {
	com.DbFprintf("trace-builtin", Dbf, "Fx_Empty Called, %d - error\n", callNo)
	return
}

// -------------------------------------------------------------------------------------------------------------------------------------------------------------------------
func FxReadJson(callNo int, pt *Parse2Type, Context *eval.ContextType, curTree **mt.MtType) (err error) {
	com.DbFprintf("trace-builtin", Dbf, "Fx_ReadJson Called, %d\n", callNo)
	// {% read_json ID "file_name.json" %} (config allows url:// not just file"

	if callNo == 0 {

		if !(*curTree).NOptions(2) {
			// xyzzy } else if !(*curTree).IsId(0) {		// -- implement to check that [0] is an ID
		} else {
			id := (*curTree).SVal[0]
			path := (*curTree).SVal[1]
			// path = path[0 : len(path)-1]
			err = nil
			// var jsonData map[string]SQLOne
			var file []byte
			file, err = ioutil.ReadFile(path)
			if err != nil {
				fmt.Fprintf(Dbf, "Error(10014): %v, %s, Config File:%s\n", err, com.LF(), path)
				return
			}
			file = []byte(strings.Replace(string(file), "\t", " ", -1)) // file = []byte(ReplaceString(string(file), "^[ \t][ \t]*//.*$", ""))

			// Check beginning of file if "{" then MapOf, if "[" Array, else look at single value
			if strings.HasPrefix(string(file), "{") {

				jsonData := make(map[string]interface{})

				err = json.Unmarshal(file, &jsonData)
				if err != nil {
					fmt.Fprintf(Dbf, "Error(10012): %v, %s, Config File:%s\n", err, com.LF(), path)
					return
				}

				Context.SetInContext(id, eval.CtxType_MapOf, jsonData)

			} else {

				jsonData := make([]interface{}, 0, 100)

				err = json.Unmarshal(file, &jsonData)
				if err != nil {
					fmt.Fprintf(Dbf, "Error(10012): %v, %s, Config File:%s\n", err, com.LF(), path)
					return
				}

				Context.SetInContext(id, eval.CtxType_ArrayOf, jsonData)

			}

		}
	}

	return
}

// -- Compares based on looping ------------------------------------------------------------------------------------------------------------

type FxDataChangeType struct {
	InitFlag bool
	TheData  interface{}
}

// xyzzy-FxIfChanged (partially completed)			00
func FxIfChanged(callNo int, pt *Parse2Type, Context *eval.ContextType, curTree **mt.MtType) (err error) {
	com.DbFprintf("trace-builtin", Dbf, "Fx_IfChanged Called, %d\n", callNo)

	tmpMt := func(ss []*mt.MtType) (rv *mt.MtType) {
		rv = &mt.MtType{
			NodeType: gen.Tok_Tree_List,
			List:     make([]*mt.MtType, 0, len(ss)),
			LineNo:   ss[0].LineNo,
			ColNo:    ss[0].ColNo,
			FileName: ss[0].FileName,
		}
		for _, vv := range ss {
			rv.List = append(rv.List, vv)
		}
		return
	}

	// No NO NO!
	// Eval expression after if - if it is different than previous value stored in store - or if store not defined -then true-
	// else false.

	if callNo == 1 {
		x := &FxDataChangeType{InitFlag: false}
		(*curTree).DataVal = x
	}
	if callNo == 11 {
		dv := (*curTree).DataVal.(FxDataChangeType)
		isInit := false
		if !dv.InitFlag {
			dv.InitFlag = true
			isInit = true
			dv.TheData = (*curTree).XValue
			(*curTree).DataVal = dv
		}

		if !(*curTree).MoreThan(0) {
		} else {
			ifp := mt.FindTags((*curTree).List[0], gen.Tok_Tree_ElsIf, gen.Tok_Tree_Else, gen.Tok_Tree_Endif) // find parts of if/else
			fmt.Fprintf(Dbf, "ifp=%+v, 1st expr = %v\n", ifp, (*curTree).EvalExpr(Context, 0, 0))
			// xyzzy - should check order of ElsIf...Else...EndIf
			if !isInit && (*curTree).EvalExpr(Context, 0, 0) {
				if ValuesSame(dv.TheData, (*curTree).XValue) {
					x := tmpMt((*curTree).List[0].List[0:ifp[0]])
					pt.x_walk(&x, pt.pos, pt.depth) // xyzzy
					return
				}
			}
			dv.TheData = (*curTree).XValue
			(*curTree).DataVal = dv
			fmt.Fprintf(Dbf, "At AT: %s\n", com.LF())
			for i := 0; i < len(ifp)-1; i++ {
				ct := (*curTree).List[0].List[ifp[i]]
				fmt.Fprintf(Dbf, "At AT: %s\n", com.LF())
				if ct.NodeType == gen.Tok_Tree_ElsIf {
					fmt.Fprintf(Dbf, "At AT, it is (((ElsIf))): %s\n", com.LF()) //
					if ct.EvalExpr(Context, 0, 0) {                              // expression is correct
						fmt.Fprintf(Dbf, "At AT: %s, ct=%+v\n", com.LF(), ct) //
						// -- xyzzy - use native .(type) and a switch
						if ct.DataType == eval.CtxType_Bool && ct.XValue.(bool) { // If true value for expression
							x := tmpMt((*curTree).List[0].List[ifp[i]+1 : ifp[i+1]])
							fmt.Fprintf(Dbf, "At -- Need to collect results -- AT: %s -------- elsif sub-tree Range[%d,%d] is %s\n", com.LF(), ifp[i]+1, ifp[i+1], com.SVarI(x))
							pt.x_walk(&x, pt.pos, pt.depth) // xyzzy
							return
						}
					}
				} else if ct.NodeType == gen.Tok_Tree_Else {
					fmt.Fprintf(Dbf, "At AT, it is (((Else))): %s\n", com.LF())
					x := tmpMt((*curTree).List[0].List[ifp[i]+1 : ifp[i+1]])
					pt.x_walk(&x, pt.pos, pt.depth) // xyzzy
					return
				}
			}
		}
	}

	return
}

func ValuesSame(x interface{}, y interface{}) bool {
	switch x.(type) {
	case int:
		a := x.(int)
		switch y.(type) {
		case int:
			b := y.(int)
			return a == b
		case int32:
			b := int(y.(int32))
			return a == b
		case int64:
			b := int(y.(int64))
			return a == b
		default:
			fmt.Printf("Error: Incompatible types %T and %T\n", x, y)
			return false
		}
	case int32:
		a := int(x.(int32))
		switch y.(type) {
		case int:
			b := y.(int)
			return a == b
		case int32:
			b := int(y.(int32))
			return a == b
		case int64:
			b := int(y.(int64))
			return a == b
		default:
			fmt.Printf("Error: Incompatible types %T and %T\n", x, y)
			return false
		}
	case int64:
		a := int(x.(int64))
		switch y.(type) {
		case int:
			b := y.(int)
			return a == b
		case int32:
			b := int(y.(int32))
			return a == b
		case int64:
			b := int(y.(int64))
			return a == b
		default:
			fmt.Printf("Error: Incompatible types %T and %T\n", x, y)
			return false
		}
	case float64:
		a := x.(float64)
		switch y.(type) {
		case int:
			b := float64(y.(int))
			return a == b
		case int32:
			b := float64(y.(int32))
			return a == b
		case int64:
			b := float64(y.(int64))
			return a == b
		case float64:
			b := y.(float64)
			return a == b
		default:
			fmt.Printf("Error: Incompatible types %T and %T\n", x, y)
			return false
		}
	case string:
		a := x.(string)
		switch y.(type) {
		case string:
			b := y.(string)
			return a == b
		default:
			fmt.Printf("Error: Incompatible types %T and %T\n", x, y)
			return false
		}
	default:
		fmt.Printf("Error: Incompatible types %T \n", x)
		return false
	}
	return false
}

// -------------------------------------------------------------------------------------------------------------------------------------------------------------------------
func FxIfEqual(callNo int, pt *Parse2Type, Context *eval.ContextType, curTree **mt.MtType) (err error) {
	com.DbFprintf("trace-builtin", Dbf, "Fx_IfEqual Called, %d\n", callNo)

	// xyzzy-FxIfEqual (partially completed)

	return
}

// -------------------------------------------------------------------------------------------------------------------------------------------------------------------------
func FxIfNotChanged(callNo int, pt *Parse2Type, Context *eval.ContextType, curTree **mt.MtType) (err error) {
	com.DbFprintf("trace-builtin", Dbf, "Fx_IfNotChanged Called, %d\n", callNo)

	// xyzzy-FxIfNotChanged (partially completed)

	return
}

// -------------------------------------------------------------------------------------------------------------------------------------------------------------------------
func FxIfNotEqual(callNo int, pt *Parse2Type, Context *eval.ContextType, curTree **mt.MtType) (err error) {
	com.DbFprintf("trace-builtin", Dbf, "Fx_IfNotEqual Called, %d\n", callNo)

	// xyzzy-FxIfNotEqual (partially completed)

	return
}

// -------------------------------------------------------------------------------------------------------------------------------------------------------------------------
func FxDebug(callNo int, pt *Parse2Type, Context *eval.ContextType, curTree **mt.MtType) (err error) {
	com.DbFprintf("trace-builtin", Dbf, "FxDebug Called, %d\n", callNo)

	if callNo == 0 || callNo == 1 {
		l := len((*curTree).SVal)
		b := true
		if l > 0 {
			if l > 1 {
				b = com.ParseBool((*curTree).SVal[1]) // b = (*curTree).SVal[1] == "on" || (*curTree).SVal[1] == "yes"
			}
			name := (*curTree).SVal[0]
			com.DbOnFlags[name] = b
		}
	}

	return
}

// -------------------------------------------------------------------------------------------------------------------------------------------------------------------------
func FxInclude(callNo int, pt *Parse2Type, Context *eval.ContextType, curTree **mt.MtType) (err error) {
	com.DbFprintf("trace-builtin", Dbf, "FxInclude Called, %d\n", callNo)

	// xyzzy-FxInclude

	return
}

// -------------------------------------------------------------------------------------------------------------------------------------------------------------------------
func FxLoad(callNo int, pt *Parse2Type, Context *eval.ContextType, curTree **mt.MtType) (err error) {
	com.DbFprintf("trace-builtin", Dbf, "FxLoad Called, %d\n", callNo)

	// xyzzy-FxLoad

	return
}

// -------------------------------------------------------------------------------------------------------------------------------------------------------------------------
func FxLorem(callNo int, pt *Parse2Type, Context *eval.ContextType, curTree **mt.MtType) (err error) {
	com.DbFprintf("trace-builtin", Dbf, "FxLorem Called, %d\n", callNo)

	mode_count := 1      // 1st Pram (expression)
	method := "b"        // w, p, b - one of
	mode_random := false // "random" - if marked

	// If last == "random" -- then pull off and mark
	// If last == "w"|"p"|"b" - then pull off and mark
	end := len((*curTree).SVal)
	fmt.Printf("Lorem end=%d at top\n", end)
	nu := 0
	if end > 1 {
		last := (*curTree).SVal[end-1]
		end--
		if last == "random" {
			mode_random = true // "random" - if marked
			nu++
		}
		if last == "w" || last == "b" || last == "p" {
			method = last
			nu++
		} else if end > 1 {
			last := (*curTree).SVal[end-1]
			end--
			if last == "w" || last == "b" || last == "p" {
				method = last
				nu++
			}
		}
	}

	// Eval expression for count - must result in number
	if end > 0 {
		(*curTree).EvalExpr(Context, 0, nu)
		fmt.Printf("Lorem - evaluating expression from 0 to nu=%d\n", nu)
		switch (*curTree).XValue.(type) {
		case int:
			mode_count = (*curTree).XValue.(int)
		case int64:
			mode_count = int((*curTree).XValue.(int64))
		case int32:
			mode_count = int((*curTree).XValue.(int32))
		default:
			fmt.Printf("Error: Invalid expression\n") // xyzzy -- Catch/report errors
		}
	}

	fmt.Printf("Lorem Control Values, count=%d method=%s random=%v\n", mode_count, method, mode_random)

	s := ""
	com := ""

	switch method {
	case "b":
		if mode_random {
			for i := 0; i < mode_count; i++ {
				s = s + com + tagLoremParagraphs[rand.Intn(len(tagLoremParagraphs))]
				com = "\n"
			}
		} else {
			for i := 0; i < mode_count; i++ {
				s = s + com + tagLoremParagraphs[i%len(tagLoremParagraphs)]
				com = "\n"
			}
		}
	case "w":
		if mode_random {
			for i := 0; i < mode_count; i++ {
				s = s + com + tagLoremWords[rand.Intn(len(tagLoremWords))]
				com = " "
			}
		} else {
			for i := 0; i < mode_count; i++ {
				s = s + com + tagLoremWords[i%len(tagLoremWords)]
				com = " "
			}
		}
	case "p":
		if mode_random {
			for i := 0; i < mode_count; i++ {
				s = s + com + "<p>" + tagLoremParagraphs[rand.Intn(len(tagLoremParagraphs))] + "</p>"
				com = "\n"
			}
		} else {
			for i := 0; i < mode_count; i++ {
				s = s + com + "<p>" + tagLoremParagraphs[i%len(tagLoremParagraphs)] + "</p>"
				com = "\n"
			}
		}
	}
	fmt.Printf("Lorem: %s\n", s)
	(*curTree).XValue = s
	(*curTree).HTML_Output = s

	return
}

// -------------------------------------------------------------------------------------------------------------------------------------------------------------------------
func FxNow(callNo int, pt *Parse2Type, Context *eval.ContextType, curTree **mt.MtType) (err error) {
	com.DbFprintf("trace-builtin", Dbf, "FxNow Called, %d\n", callNo)

	// xyzzy-FxNow

	return
}

// -------------------------------------------------------------------------------------------------------------------------------------------------------------------------
func FxRegroup(callNo int, pt *Parse2Type, Context *eval.ContextType, curTree **mt.MtType) (err error) {
	com.DbFprintf("trace-builtin", Dbf, "FxRegroup Called, %d\n", callNo)

	// xyzzy-FxRegroup

	return
}

// -------------------------------------------------------------------------------------------------------------------------------------------------------------------------
func FxTemplateTag(callNo int, pt *Parse2Type, Context *eval.ContextType, curTree **mt.MtType) (err error) {
	com.DbFprintf("trace-builtin", Dbf, "FxTemplateTag Called, %d\n", callNo)

	// xyzzy-FxTemplateTag

	return
}

// -------------------------------------------------------------------------------------------------------------------------------------------------------------------------
func FxUrl(callNo int, pt *Parse2Type, Context *eval.ContextType, curTree **mt.MtType) (err error) {
	com.DbFprintf("trace-builtin", Dbf, "FxUrl Called, %d\n", callNo)

	// xyzzy-FxUrl

	return
}

// -------------------------------------------------------------------------------------------------------------------------------------------------------------------------
func FxWith(callNo int, pt *Parse2Type, Context *eval.ContextType, curTree **mt.MtType) (err error) {
	com.DbFprintf("trace-builtin", Dbf, "FxWith Called, %d\n", callNo)

	// xyzzy-FxWith

	return
}

// -------------------------------------------------------------------------------------------------------------------------------------------------------------------------
func FxWithRatio(callNo int, pt *Parse2Type, Context *eval.ContextType, curTree **mt.MtType) (err error) {
	com.DbFprintf("trace-builtin", Dbf, "FxWithRatio Called, %d\n", callNo)

	// xyzzy-FxWithRatio

	return
}

// -------------------------------------------------------------------------------------------------------------------------------------------------------------------------
func FxVerbatim(callNo int, pt *Parse2Type, Context *eval.ContextType, curTree **mt.MtType) (err error) {
	com.DbFprintf("trace-builtin", Dbf, "FxVerbatim Called, %d\n", callNo)

	// xyzzy-FxVerbatim

	return
}

// -------------------------------------------------------------------------------------------------------------------------------------------------------------------------
func FxSpaceless(callNo int, pt *Parse2Type, Context *eval.ContextType, curTree **mt.MtType) (err error) {
	com.DbFprintf("trace-builtin", Dbf, "FxSpaceless Called, %d\n", callNo)
	fmt.Printf("This One: !!!\n")

	var walkTree func(mtv **mt.MtType, pos, depth int)
	walkTree = func(mtv **mt.MtType, pos, depth int) {
		lines := strings.Split(string((*mtv).XValue.(string)), "\n")
		for ii, ss := range lines {
			tt := strings.TrimSpace(ss)
			fmt.Printf("This One: before >%s< after >%s<\n", ss, tt)
			lines[ii] = tt
		}
		(*mtv).XValue = strings.Join(lines, "\n")
		for ii, _ := range (*mtv).List {
			walkTree(&((*mtv).List[ii]), ii, depth+1)
		}
	}
	if callNo == 11 {
		walkTree(curTree, 0, 0)
	}

	return
}

// -------------------------------------------------------------------------------------------------------------------------------------------------------------------------
func FxFilter(callNo int, pt *Parse2Type, Context *eval.ContextType, curTree **mt.MtType) (err error) {
	com.DbFprintf("trace-builtin", Dbf, "FxFilter Called, %d\n", callNo)

	// xyzzy-FxFilter

	return
}

// -------------------------------------------------------------------------------------------------------------------------------------------------------------------------
func FxComment(callNo int, pt *Parse2Type, Context *eval.ContextType, curTree **mt.MtType) (err error) {
	com.DbFprintf("trace-builtin", Dbf, "FxComment Called, %d\n", callNo)

	if callNo == 11 {
		(*curTree).List = (*curTree).List[:0] // remove tree?
	}

	return
}

// -------------------------------------------------------------------------------------------------------------------------------------------------------------------------
func FxEndComment(callNo int, pt *Parse2Type, Context *eval.ContextType, curTree **mt.MtType) (err error) {
	com.DbFprintf("trace-builtin", Dbf, "FxEndComment Called, %d\n", callNo)

	// done

	return
}

// -------------------------------------------------------------------------------------------------------------------------------------------------------------------------
func FxBlock(callNo int, pt *Parse2Type, Context *eval.ContextType, curTree **mt.MtType) (err error) {
	com.DbFprintf("trace-builtin", Dbf, "FxBlock Called, %d\n", callNo)

	// done - Acts as a named marker in the template - no implementaiton

	return
}

// -------------------------------------------------------------------------------------------------------------------------------------------------------------------------
func FxExtend(callNo int, pt *Parse2Type, Context *eval.ContextType, curTree **mt.MtType) (err error) {
	com.DbFprintf("trace-builtin", Dbf, "FxExtend Called, %d\n", callNo)

	e_block := make([]*mt.MtType, 0, 20)

	var walkTree func(mtv *mt.MtType, pos, depth int)
	walkTree = func(mtv *mt.MtType, pos, depth int) {
		for ii, vv := range mtv.List {
			if vv.FxId == gen.Fx_block {
				com.DbFprintf("trace-builtin", Dbf, "FxExtend [%d] found block with name >%s<-, %s\n", ii, vv.SVal[0], com.LF())
				e_block = append(e_block, vv)
			}
			walkTree(mtv.List[ii], ii, depth+1)
		}
	}

	if callNo == 11 {
		// if vv, ok := pt.PMatch ( 11, gen.Tok_ID ) ; ok { 		// {% extend <name> %}
		// } else if vv, ok := pt.PMatch ( 11, gen.Tok_Str0 ) ; ok { 		// {% extend "string" %}
		// } else if vv, ok := pt.PMatch ( 11, gen.Tok_ID, gen.Tok_as, gen.Tok_ID_or_Str  ) ; ok { 		// {% extend <name> as <name> %}
		// } else if vv, ok := pt.PMatch ( 11, gen.Tok_ID, gen.Tok_as, gen.Tok_Expr  ) ; ok { 		// {% extend <name> as <name> %}

		// xyzzy - needs to be a loop to pull out set of template names
		name := (*curTree).SVal[0]
		thisname := ""
		output_template := false
		if len((*curTree).SVal) > 2 && (*curTree).SVal[1] == "as" {
			fmt.Fprintf(Dbf, "At: %s\n", com.LF())
			thisname = (*curTree).SVal[2]
		} else {
			fmt.Fprintf(Dbf, "At: %s\n", com.LF())
			thisname = (*curTree).FileName
			if len(thisname) >= 2 && thisname[0:2] == "./" {
				fmt.Fprintf(Dbf, "At: %s\n", com.LF())
				thisname = thisname[2:]
			}
			fmt.Fprintf(Dbf, "At: %s\n", com.LF())
			output_template = true
		}
		// xyzzy - hot patch to test extra quote removal ------------------------------------  Wed Jul 15 20:00:12 MDT 2015
		if name[len(name)-1:] == `"` {
			fmt.Fprintf(Dbf, "At: %s\n", com.LF())
			name = name[:len(name)-1]
		}
		fmt.Fprintf(Dbf, "E2 Extending Template >%s<- with new name of ->%s<-, about to do lookup of it, should output %v\n", name, thisname, output_template)
		// xyzzy - added to debug
		// Context.DumpContext()
		// xyzzy - added to debug
		ss, err := pt.St.LookupSymbol(name)
		if err != nil {
			fmt.Fprintf(Dbf, "Error: Template %s is not defined\n", name)
		} else {
			fmt.Fprintf(Dbf, "Good: Template %s found, --- lookup tree is ---\n\n%s\n\n---end---\n", name, com.SVarI(ss.AnyData.(*mt.MtType)))
			walkTree((*curTree), 0, 0)
			newtree := mt.DuplicateTree(ss.AnyData.(*mt.MtType)) // Make copy of item -- - change with extention blocks - --
			fmt.Fprintf(Dbf, "---- newtree ---\n%s\n\n--end--\n", com.SVarI(newtree))
			for _, ww := range e_block {
				fmt.Fprintf(Dbf, "At: %s\n", com.LF())
				mt.ReplaceBlocksWithNew(&newtree, ww)
			}
			fmt.Fprintf(Dbf, "---- block repalced  ---\n%s\n\n--end--\n", com.SVarI(newtree))
			if !output_template {
				fmt.Fprintf(Dbf, "At: %s\n", com.LF())
				pt.DefineTemplate(thisname, newtree)  // pt.DefineTemplate(name string, (*curTree) *mt.MtType)
				(*curTree).List = (*curTree).List[:0] // remove tree?
			} else {
				fmt.Fprintf(Dbf, "At: %s\n", com.LF())
				(*curTree) = newtree
				// xyzzy - do I need to iterate over tree at this point? or at call point
				// xyzzy - Should the definition of the Block/End or Item have a "walk-child" flag or a "child-handles-subtree" flag
				// xyzzy - "child-handles-subtree" is for if/for or if items
				// xyzzy - the rest all have the caller just process the subtree
			}

		}
	}
	return
}

// -------------------------------------------------------------------------------------------------------------------------------------------------------------------------
// {% template <name> %} ... {% endtemplate %}
func FxTemplate(callNo int, pt *Parse2Type, Context *eval.ContextType, curTree **mt.MtType) (err error) {
	com.DbFprintf("trace-builtin", Dbf, "FxTemplate Called, %d\n", callNo)

	if callNo == 11 {
		newtree := mt.DuplicateTree((*curTree))        //
		pt.DefineTemplate((*curTree).SVal[0], newtree) // pt.DefineTemplate(name string, (*curTree) *mt.MtType)
		// pt.DefineTemplate((*curTree).SVal[0], (*curTree)) // pt.DefineTemplate(name string, (*curTree) *mt.MtType)
		(*curTree).List = (*curTree).List[:0] // remove tree?
	}

	return
}

// -------------------------------------------------------------------------------------------------------------------------------------------------------------------------
func FxRender(callNo int, pt *Parse2Type, Context *eval.ContextType, curTree **mt.MtType) (err error) {
	com.DbFprintf("trace-builtin", Dbf, "FxRender Called, %d\n", callNo)

	if callNo == 20 {

		name := (*curTree).SVal[0]

		com.DbFprintf("trace-builtin", Dbf, "FxRender E2 - pass is 20, name=%s\n", name)

		// Lookup the item
		ss, err := pt.St.LookupSymbol(name)
		if err != nil {
			fmt.Fprintf(Dbf, "E2 Attempt to render non-existent template >>%s<<\n", name)
		} else {
			// Copy and replace {%render name%} with body of render
			// (*curTree).List = append((*curTree).List, ss.AnyData.(*mt.MtType))
			(*curTree).List = append((*curTree).List, ss.AnyData.(*mt.MtType))
		}

	}

	return
}

// -------------------------------------------------------------------------------------------------------------------------------------------------------------------------
func FxLibrary(callNo int, pt *Parse2Type, Context *eval.ContextType, curTree **mt.MtType) (err error) {

	if callNo == 0 {
		for _, path := range (*curTree).SVal { // open each file // parse each template and store
			list := com.AllFilesInPath(path)
			for _, fn := range list {
				_ = pt.ReadFileAsTemplate(fn)
			}
		}
	}

	return
}
