package mt

// M T - N Ary Tree

import (
	"fmt"
	"io"
	"strings"

	"github.com/pschlump/dbgo"
	"github.com/pschlump/lexie/eval"
	"github.com/pschlump/lexie/gen"
	"github.com/pschlump/lexie/tok"
)

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
	HTML_Output string      //	// Output at end.
	List        []*MtType   // N-Ary tree of children
	LineNo      int         // Where did I come from
	ColNo       int         //
	FileName    string      //
	LValue      bool        // Is this an L-Value (assignable, chanable)
}

// -------------------------------------------------------------------------------------------------------------------------------------------------------------------------
func NewMtType(t int, v string) (mt *MtType) {
	mt = &MtType{
		NodeType:    t,
		EscapeIt:    false,
		DataType:    eval.CtxType_Str,
		XValue:      v,
		HTML_Output: "",
		LineNo:      0,
		ColNo:       0,
	}
	return
}

// -------------------------------------------------------------------------------------------------------------------------------------------------------------------------
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

// -------------------------------------------------------------------------------------------------------------------------------------------------------------------------
func (mt *MtType) MoreThan(n int) bool {
	if len(mt.SVal) <= n {
		mt.Error = true
		mt.ErrorMsg = fmt.Sprintf("Error: incorrect number of options supplied - should have more than %d options, %s", n, dbgo.LF(2))
		return false
	}
	return true
}

// -------------------------------------------------------------------------------------------------------------------------------------------------------------------------
//
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
//

// n from beginning to m from end.
func (mt *MtType) EvalExpr(Context *eval.ContextType, n, m int) bool {
	m = len(mt.SVal) - m // Convert to Pos
	sv := mt.SVal[n:m]   // Slice of params to eval.
	fmt.Printf("mt.EvalExpr - TOP: n=%d m=%d, Range [n:m]\n", n, m)
	fmt.Printf("mt.EvalExpr - TOP: ealuate sv=%+v (Subset) ----------------------------------------------------- \n", sv)
	fmt.Printf("mt.EvalExpr - TOP: ealuate mt.SVal=%+v (Orig)----------------------------------------------------- \n", mt.SVal)
	// xyzzy -- temporary -- incomplete!!!!!!!!!!!!!!!!!!!!!!
	evalData := &eval.EvalType{
		Pos:           0,
		Ctx:           Context,
		Mm:            mt.TokVal[n:m], // []tok.Token
		PrintErrorMsg: true,
	}

	// hot patch - xyzzy
	//if evalData.Mm[0].Match == "[" {
	// evalData.Mm[0].TokNo = 38
	//}
	// hot patch - xyzzy

	evalData.InitFunctions()
	fmt.Printf("INPUT m=%d n=%d, %s ----------------------------------------------------- \n", m, n, dbgo.SVarI(evalData))
	tr := evalData.PresTop()
	fmt.Printf("BOTTOM: %s ----------------------------------------------------- \n", dbgo.SVarI(tr))
	//s := sv[0]
	//v, t, _ := Context.GetFromContext(s)
	//fmt.Printf("At: %s - in EvalExpr, v=%v t=%v for >%s<-\n", dbgo.LF(), v, t, s)
	mt.DataType = tr.DataType
	mt.XValue = tr.CurValue
	fmt.Printf("At bottom of EvalExpr - Type = %d == %T, value = %v\n", tr.DataType, tr.CurValue, tr.CurValue)
	return true
}

// -------------------------------------------------------------------------------------------------------------------------------------------------------------------------
// Evaluate a set of contents or variables in 'in' - retuning an array of values in rv
func (mt *MtType) EvalVars(in []string) (rv []string) {
	rv = in
	// xyzzy
	return
}

// -------------------------------------------------------------------------------------------------------------------------------------------------------------------------
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
func FindTags(mtv *MtType, tag ...int) (rv []int) {
	for ii, vv := range mtv.List {
		for _, ww := range tag {
			if vv.NodeType == ww {
				rv = append(rv, ii)
			}
		}
	}
	return
}

// -------------------------------------------------------------------------------------------------------------------------------------------------------------------------
// Duplicate a tree
func DuplicateTree(tree *MtType) (rv *MtType) {

	var walkTree func(mtv *MtType, pos, depth int, rv **MtType)
	walkTree = func(mtv *MtType, pos, depth int, rv **MtType) {
		*rv = &MtType{
			NodeType:    mtv.NodeType,
			EscapeIt:    mtv.EscapeIt,
			ID:          mtv.ID,
			DataType:    mtv.DataType,
			XValue:      mtv.XValue,
			SVal:        mtv.SVal,
			TVal:        mtv.TVal,
			TokVal:      mtv.TokVal,
			Error:       mtv.Error,
			ErrorMsg:    mtv.ErrorMsg,
			DataVal:     mtv.DataVal,
			FxId:        mtv.FxId,
			HTML_Output: mtv.HTML_Output,
			LineNo:      mtv.LineNo,
			ColNo:       mtv.ColNo,
			FileName:    mtv.FileName,
			LValue:      mtv.LValue,
		}
		if len(mtv.List) > 0 {
			(*rv).List = make([]*MtType, len(mtv.List), len(mtv.List))
			for ii, vv := range mtv.List {
				walkTree(vv, ii, depth+1, &((*rv).List[ii]))
			}
		}
	}

	walkTree(tree, 0, 0, &rv)

	return
}

// -------------------------------------------------------------------------------------------------------------------------------------------------------------------------
// Find the nodes in tree (blocks) with stated name and replace them with new
func ReplaceBlocksWithNew(search_in_tree **MtType, new_block *MtType) {

	block_name := new_block.SVal[0]

	var walkTree func(mt **MtType, pos, depth int)
	walkTree = func(mt **MtType, pos, depth int) {
		for ii := range (*mt).List {
			//if vv.FxId == gen.Fx_block && block_name == vv.SVal[0] {
			//	fmt.Printf("FxExtend Replace: [%d] found block with name >%s<-, %s\n", ii, vv.SVal[0], dbgo.LF())
			//	*mt = new_block
			//}
			walkTree(&((*mt).List[ii]), ii, depth+1)
		}
		if len((*mt).SVal) > 0 {
			fmt.Printf("FxExtend Before FxId = %d, looking for %d (*mt).SVal[0] = >%s< looking for %s, %s\n", (*mt).FxId, gen.Fx_block, (*mt).SVal[0], block_name, dbgo.LF())
			if (*mt).FxId == gen.Fx_block && block_name == (*mt).SVal[0] {
				fmt.Printf("FxExtend Replace: found block with name >%s<-, %s\n", (*mt).SVal[0], dbgo.LF())
				*mt = new_block
			}
		}
	}

	walkTree(search_in_tree, 0, 0)

	return
}

/* vim: set noai ts=4 sw=4: */
