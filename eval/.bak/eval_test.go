package eval

// xyzzy1857 - incomplete - len/length(), - auto tests, more tests
// xyzzy - test: -1 for array pos start, len()+1 for end.

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"strings"
	"testing"

	"../com"
	"../gen"
	"../tok"

	"../../../go-lib/sizlib"

	// "../../../go-lib/sizlib"
)

const (
	CmdEval            = 1
	CmdInsertInContext = 2
	CmdNop             = 3
	CmdDone            = 4
	CmdDumpContext     = 5
	CmdReadJson        = 6 // Read json file into context
	CmdCmpTo           = 7
	CmdValidateLnCol   = 8  // Check Line_no, Col_no, File_Name on Error Message
	CmdErrorsContain   = 9  // Errors have this error in it.
	CmdSetGlob         = 10 // Set global to NValue
	CmdCmpGlob         = 11 // Set global to NValue
	CmdArrayLen        = 12 //
	CmdArrayValue      = 13 // IValue is subscript, NValue is data
)

type ActType struct {
	OpCode   int         //
	Item     []tok.Token //
	Id       string      //
	DataType int         //
	NValue   int         //
	IValue   int         //
	SValue   string      //
	BValue   bool        //
	FValue   float64     //
	FxValue  interface{}
	Data     string //
	Line_No  int    //	For CmptValidateLnCol
	Col_No   int    //
	FileName string //
}

type Ev01TestType struct {
	Test     string    //
	SkipTest bool      //
	Actions  []ActType //
}

var Ev01Test = []Ev01TestType{
	{Test: "0001", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_NUM, Match: "1", DataType: CtxType_Int, CurValue: 1},
			tok.Token{TokNo: gen.Tok_PLUS, Match: "+"},
			tok.Token{TokNo: gen.Tok_NUM, Match: "2", DataType: CtxType_Int, CurValue: 2},
		}},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Int, NValue: 3},
		ActType{OpCode: CmdDone},
	}},
	{Test: "0002", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdInsertInContext, Id: "abc", DataType: CtxType_Int, NValue: 5},
		ActType{OpCode: CmdInsertInContext, Id: "def", DataType: CtxType_Int, NValue: 10},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "abc"},
			tok.Token{TokNo: gen.Tok_PLUS, Match: "+"},
			tok.Token{TokNo: gen.Tok_ID, Match: "def"},
		}},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Int, NValue: 15},
		ActType{OpCode: CmdDone},
	}},
	{Test: "0003", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdInsertInContext, Id: "abc", DataType: CtxType_Int, NValue: 5},
		ActType{OpCode: CmdInsertInContext, Id: "def", DataType: CtxType_Int, NValue: 10},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "abc"},
			tok.Token{TokNo: gen.Tok_PLUS, Match: "+"},
			tok.Token{TokNo: gen.Tok_ID, Match: "def"},
			tok.Token{TokNo: gen.Tok_PLUS, Match: "+"},
			tok.Token{TokNo: gen.Tok_NUM, Match: "3", DataType: CtxType_Int, CurValue: 3},
		}},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Int, NValue: 18},
		ActType{OpCode: CmdDone},
	}},
	{Test: "0004", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdInsertInContext, Id: "abc", DataType: CtxType_Int, NValue: 5},
		ActType{OpCode: CmdInsertInContext, Id: "def", DataType: CtxType_Int, NValue: 10},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "abc"},
			tok.Token{TokNo: gen.Tok_PLUS, Match: "+"},
			tok.Token{TokNo: gen.Tok_ID, Match: "def"},
			tok.Token{TokNo: gen.Tok_STAR, Match: "*"},
			tok.Token{TokNo: gen.Tok_NUM, Match: "3", DataType: CtxType_Int, CurValue: 3},
		}},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Int, NValue: 35},
		ActType{OpCode: CmdDone},
	}},

	// eval_test.go:1896: 0005: 4 error, expected 8 got 4, found it, File: /Users/corwin/Projects/pongo2/lexie/eval/eval_test.go LineNo:1896
	// Test Parens abc * ( 1 + 2 )
	{Test: "0005", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdInsertInContext, Id: "abc", DataType: CtxType_Int, NValue: 2},
		ActType{OpCode: CmdInsertInContext, Id: "def", DataType: CtxType_Int, NValue: 1},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "abc"},
			tok.Token{TokNo: gen.Tok_STAR, Match: "*"},
			tok.Token{TokNo: gen.Tok_OP, Match: "("},
			tok.Token{TokNo: gen.Tok_ID, Match: "def"},
			tok.Token{TokNo: gen.Tok_PLUS, Match: "+"},
			tok.Token{TokNo: gen.Tok_NUM, Match: "3", DataType: CtxType_Int, CurValue: 3},
			tok.Token{TokNo: gen.Tok_CL, Match: ")"},
		}},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Int, NValue: 8},
		ActType{OpCode: CmdDone},
	}},

	// Test Unary -
	{Test: "0006", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdInsertInContext, Id: "abc", DataType: CtxType_Int, NValue: 5},
		ActType{OpCode: CmdInsertInContext, Id: "def", DataType: CtxType_Int, NValue: 10},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "abc"},
			tok.Token{TokNo: gen.Tok_PLUS, Match: "+"},
			tok.Token{TokNo: gen.Tok_MINUS, Match: "-"},
			tok.Token{TokNo: gen.Tok_ID, Match: "def"},
		}},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Int, NValue: -5},
		ActType{OpCode: CmdDone},
	}},
	{Test: "0007", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdInsertInContext, Id: "abc", DataType: CtxType_Int, NValue: 5},
		ActType{OpCode: CmdInsertInContext, Id: "def", DataType: CtxType_Int, NValue: 10},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_MINUS, Match: "-"},
			tok.Token{TokNo: gen.Tok_ID, Match: "abc"},
			tok.Token{TokNo: gen.Tok_PLUS, Match: "+"},
			tok.Token{TokNo: gen.Tok_MINUS, Match: "-"},
			tok.Token{TokNo: gen.Tok_ID, Match: "def"},
		}},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Int, NValue: -15},
		ActType{OpCode: CmdDone},
	}},

	// Test assinging new value to ID
	{Test: "0008", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdInsertInContext, Id: "abc", DataType: CtxType_Int, NValue: 5},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "abc"},
			tok.Token{TokNo: gen.Tok_EQ, Match: "="},
			tok.Token{TokNo: gen.Tok_NUM, Match: "22", DataType: CtxType_Int, CurValue: 22},
		}},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Int, NValue: 22},
		ActType{OpCode: CmdDone},
	}},

	// tests automatic creation of ID on assignment, may want to change this to := for new IDs, and = for IDs that already exist
	{Test: "0009", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdInsertInContext, Id: "abc", DataType: CtxType_Int, NValue: 5},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "abc"},
			tok.Token{TokNo: gen.Tok_EQ, Match: "="},
			tok.Token{TokNo: gen.Tok_NUM, Match: "22", DataType: CtxType_Int, CurValue: 22},
		}},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Int, NValue: 22},
		ActType{OpCode: CmdDone},
	}},

	// eval.Mm[eval.Pos].CreateId = (eval.Pos+1 < len(eval.MM) && eval.Mm[eval.Pos+1] == gen.Tok_DCL_VAR)
	{Test: "0010", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "abc"},
			tok.Token{TokNo: gen.Tok_DCL_VAR, Match: ":="},
			tok.Token{TokNo: gen.Tok_NUM, Match: "22", DataType: CtxType_Int, CurValue: 22},
		}},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Int, NValue: 22},
		ActType{OpCode: CmdDone},
	}},

	// 1st array ref
	{Test: "0011", SkipTest: true, Actions: []ActType{
		ActType{OpCode: CmdReadJson, Id: "arr0", Data: "./arr1.json"},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "arr0"},
			tok.Token{TokNo: gen.Tok_OP_SQ, Match: "["},
			tok.Token{TokNo: gen.Tok_NUM, Match: "2", DataType: CtxType_Int, CurValue: 2},
			tok.Token{TokNo: gen.Tok_CL_SQ, Match: "]"},
		}},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Int, NValue: 88},
		ActType{OpCode: CmdDone},
	}},

	{Test: "0012", SkipTest: true, Actions: []ActType{
		ActType{OpCode: CmdReadJson, Id: "arr0", Data: "./arr2.json"},
		ActType{OpCode: CmdInsertInContext, Id: "bob", DataType: CtxType_Str, SValue: "xyz"},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "arr0"},
			tok.Token{TokNo: gen.Tok_OP_SQ, Match: "["},
			tok.Token{TokNo: gen.Tok_ID, Match: "bob"},
			tok.Token{TokNo: gen.Tok_CL_SQ, Match: "]"},
		}},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Str, SValue: "xyz1"},
		ActType{OpCode: CmdDone},
	}},

	{Test: "0013", SkipTest: true, Actions: []ActType{
		ActType{OpCode: CmdReadJson, Id: "arr0", Data: "./arr2.json"},
		ActType{OpCode: CmdInsertInContext, Id: "bob", DataType: CtxType_Str, SValue: "xyz"},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "arr0"},
			tok.Token{TokNo: gen.Tok_DOT, Match: "."},
			tok.Token{TokNo: gen.Tok_ID, Match: "abc"},
		}},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Str, SValue: "abc1"},
		ActType{OpCode: CmdDone},
	}},

	// Test Unary !
	{Test: "0014", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdInsertInContext, Id: "abc", DataType: CtxType_Bool, BValue: false},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_EXCLAM, Match: "!"},
			tok.Token{TokNo: gen.Tok_ID, Match: "abc"},
		}},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Bool, BValue: true},
		ActType{OpCode: CmdDone},
	}},

	{Test: "0015", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdInsertInContext, Id: "abc", DataType: CtxType_Bool, BValue: false},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_EXCLAM, Match: "!"},
			tok.Token{TokNo: gen.Tok_EXCLAM, Match: "!"},
			tok.Token{TokNo: gen.Tok_EXCLAM, Match: "!"},
			tok.Token{TokNo: gen.Tok_ID, Match: "abc"},
		}},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Bool, BValue: true},
		ActType{OpCode: CmdDone},
	}},

	{Test: "0016", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdInsertInContext, Id: "def", DataType: CtxType_Int, NValue: 10},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_EXCLAM, Match: "!"},
			tok.Token{TokNo: gen.Tok_EXCLAM, Match: "!"},
			tok.Token{TokNo: gen.Tok_EXCLAM, Match: "!"},
			tok.Token{TokNo: gen.Tok_EXCLAM, Match: "!"},
			tok.Token{TokNo: gen.Tok_EXCLAM, Match: "!"},
			tok.Token{TokNo: gen.Tok_ID, Match: "def"},
		}},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Bool, BValue: true},
		ActType{OpCode: CmdDone},
	}},

	{Test: "0017", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdInsertInContext, Id: "abc", DataType: CtxType_Int, NValue: 5},
		ActType{OpCode: CmdInsertInContext, Id: "def", DataType: CtxType_Int, NValue: 100},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "abc"},
			tok.Token{TokNo: gen.Tok_MINUS, Match: "-"},
			tok.Token{TokNo: gen.Tok_ID, Match: "def"},
			tok.Token{TokNo: gen.Tok_SLASH, Match: "/"},
			tok.Token{TokNo: gen.Tok_NUM, Match: "5", DataType: CtxType_Int, CurValue: 5}, // should be -15
			tok.Token{TokNo: gen.Tok_PCT, Match: "%"},
			tok.Token{TokNo: gen.Tok_NUM, Match: "10", DataType: CtxType_Int, CurValue: 10}, // should be -5
		}},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Int, NValue: 5},
		ActType{OpCode: CmdDone},
	}},

	{Test: "0018", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_Float, Match: "1.0", DataType: CtxType_Float, CurValue: 1.0},
			tok.Token{TokNo: gen.Tok_PLUS, Match: "+"},
			tok.Token{TokNo: gen.Tok_Float, Match: "2.0", DataType: CtxType_Float, CurValue: 2.0},
		}},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Float, FValue: 3},
		ActType{OpCode: CmdDone},
	}},

	{Test: "0019", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_NUM, Match: "1", DataType: CtxType_Int, CurValue: 1},
			tok.Token{TokNo: gen.Tok_PLUS, Match: "+"},
			tok.Token{TokNo: gen.Tok_Float, Match: "2.0", DataType: CtxType_Float, CurValue: 2.0},
		}},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Float, FValue: 3},
		ActType{OpCode: CmdDone},
	}},

	{Test: "0020", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_Float, Match: "1.0", DataType: CtxType_Float, CurValue: 1.0},
			tok.Token{TokNo: gen.Tok_PLUS, Match: "+"},
			tok.Token{TokNo: gen.Tok_NUM, Match: "2", DataType: CtxType_Int, CurValue: 2},
		}},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Float, FValue: 3},
		ActType{OpCode: CmdDone},
	}},

	{Test: "0021", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_NUM, Match: "1", DataType: CtxType_Int, CurValue: 1},
			tok.Token{TokNo: gen.Tok_S_L, Match: "<<"},
			tok.Token{TokNo: gen.Tok_NUM, Match: "2", DataType: CtxType_Int, CurValue: 2},
		}},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Int, NValue: 4},
		ActType{OpCode: CmdDone},
	}},

	{Test: "0022", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_NUM, Match: "16", DataType: CtxType_Int, CurValue: 16},
			tok.Token{TokNo: gen.Tok_S_R, Match: "<<"},
			tok.Token{TokNo: gen.Tok_NUM, Match: "2", DataType: CtxType_Int, CurValue: 2},
		}},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Int, NValue: 4},
		ActType{OpCode: CmdDone},
	}},

	{Test: "0023", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_NUM, Match: "1", DataType: CtxType_Int, CurValue: 1},
			tok.Token{TokNo: gen.Tok_NE, Match: "!="},
			tok.Token{TokNo: gen.Tok_NUM, Match: "2", DataType: CtxType_Int, CurValue: 2},
		}},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Bool, BValue: true},
		ActType{OpCode: CmdDone},
	}},

	{Test: "0024", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_NUM, Match: "1", DataType: CtxType_Int, CurValue: 1},
			tok.Token{TokNo: gen.Tok_NE, Match: "<>"},
			tok.Token{TokNo: gen.Tok_NUM, Match: "2", DataType: CtxType_Int, CurValue: 2},
		}},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Bool, BValue: true},
		ActType{OpCode: CmdDone},
	}},

	{Test: "0025", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_NUM, Match: "1", DataType: CtxType_Int, CurValue: 1},
			tok.Token{TokNo: gen.Tok_L_EQ, Match: "=="},
			tok.Token{TokNo: gen.Tok_NUM, Match: "2", DataType: CtxType_Int, CurValue: 2},
		}},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Bool, BValue: false},
		ActType{OpCode: CmdDone},
	}},

	{Test: "0026", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_NUM, Match: "1", DataType: CtxType_Int, CurValue: 1},
			tok.Token{TokNo: gen.Tok_L_EQ, Match: "=="},
			tok.Token{TokNo: gen.Tok_NUM, Match: "1", DataType: CtxType_Int, CurValue: 1},
		}},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Bool, BValue: true},
		ActType{OpCode: CmdDone},
	}},

	{Test: "0027", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_NUM, Match: "1", DataType: CtxType_Int, CurValue: 1},
			tok.Token{TokNo: gen.Tok_LE, Match: "<="},
			tok.Token{TokNo: gen.Tok_NUM, Match: "2", DataType: CtxType_Int, CurValue: 2},
			tok.Token{TokNo: gen.Tok_L_EQ, Match: "=="},
			tok.Token{TokNo: gen.Tok_ID, Match: "true"},
			tok.Token{TokNo: gen.Tok_NE, Match: "!="},
			tok.Token{TokNo: gen.Tok_ID, Match: "false"},
		}},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Bool, BValue: true},
		ActType{OpCode: CmdDone},
	}},

	{Test: "0028", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_NUM, Match: "1", DataType: CtxType_Int, CurValue: 1},
			tok.Token{TokNo: gen.Tok_GE, Match: ">="},
			tok.Token{TokNo: gen.Tok_NUM, Match: "2", DataType: CtxType_Int, CurValue: 2},
		}},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Bool, BValue: false},
		ActType{OpCode: CmdDone},
	}},

	{Test: "0029", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_NUM, Match: "1", DataType: CtxType_Int, CurValue: 1},
			tok.Token{TokNo: gen.Tok_GT, Match: ">"},
			tok.Token{TokNo: gen.Tok_NUM, Match: "1", DataType: CtxType_Int, CurValue: 1},
		}},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Bool, BValue: false},
		ActType{OpCode: CmdDone},
	}},

	{Test: "0030", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_NUM, Match: "2", DataType: CtxType_Int, CurValue: 2},
			tok.Token{TokNo: gen.Tok_GT, Match: ">"},
			tok.Token{TokNo: gen.Tok_NUM, Match: "1", DataType: CtxType_Int, CurValue: 1},
		}},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Bool, BValue: true},
		ActType{OpCode: CmdDone},
	}},

	// Test Unary -, +
	{Test: "0031", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdInsertInContext, Id: "abc", DataType: CtxType_Int, NValue: 5},
		ActType{OpCode: CmdInsertInContext, Id: "def", DataType: CtxType_Int, NValue: 10},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "abc"},
			tok.Token{TokNo: gen.Tok_PLUS, Match: "+"},
			tok.Token{TokNo: gen.Tok_MINUS, Match: "-"},
			tok.Token{TokNo: gen.Tok_MINUS, Match: "-"},
			tok.Token{TokNo: gen.Tok_PLUS, Match: "+"},
			tok.Token{TokNo: gen.Tok_MINUS, Match: "-"},
			tok.Token{TokNo: gen.Tok_ID, Match: "def"},
		}},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Int, NValue: -5},
		ActType{OpCode: CmdDone},
	}},

	// Test Unary -, + on a float
	// - use - (unary) - on float
	{Test: "0032", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdInsertInContext, Id: "abc", DataType: CtxType_Int, NValue: 5},
		ActType{OpCode: CmdInsertInContext, Id: "def", DataType: CtxType_Float, FValue: 10.0},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "abc"},
			tok.Token{TokNo: gen.Tok_PLUS, Match: "+"},
			tok.Token{TokNo: gen.Tok_MINUS, Match: "-"},
			tok.Token{TokNo: gen.Tok_MINUS, Match: "-"},
			tok.Token{TokNo: gen.Tok_PLUS, Match: "+"},
			tok.Token{TokNo: gen.Tok_MINUS, Match: "-"},
			tok.Token{TokNo: gen.Tok_ID, Match: "def"},
		}},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Float, FValue: -5.0},
		ActType{OpCode: CmdDone},
	}},

	// - use +=
	{Test: "0033", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdInsertInContext, Id: "abc", DataType: CtxType_Int, NValue: 5},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "abc"},
			tok.Token{TokNo: gen.Tok_PLUS_EQ, Match: "+="},
			tok.Token{TokNo: gen.Tok_NUM, Match: "22", DataType: CtxType_Int, CurValue: 22},
		}},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Int, NValue: 27},
		ActType{OpCode: CmdDone},
	}},

	// - use += (float)
	{Test: "0034", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdInsertInContext, Id: "abc", DataType: CtxType_Int, NValue: 5},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "abc"},
			tok.Token{TokNo: gen.Tok_PLUS_EQ, Match: "+="},
			tok.Token{TokNo: gen.Tok_Float, Match: "22.0", DataType: CtxType_Float, CurValue: 22.0},
		}},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Float, FValue: 27.0},
		ActType{OpCode: CmdDone},
	}},

	// - use *=
	{Test: "0035", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdInsertInContext, Id: "abc", DataType: CtxType_Int, NValue: 5},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "abc"},
			tok.Token{TokNo: gen.Tok_STAR_EQ, Match: "*="},
			tok.Token{TokNo: gen.Tok_NUM, Match: "2", DataType: CtxType_Int, CurValue: 2},
		}},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Int, NValue: 10},
		ActType{OpCode: CmdDone},
	}},

	// - use *= (float)
	{Test: "0036", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdInsertInContext, Id: "abc", DataType: CtxType_Int, NValue: 5},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "abc"},
			tok.Token{TokNo: gen.Tok_STAR_EQ, Match: "*="},
			tok.Token{TokNo: gen.Tok_Float, Match: "2.0", DataType: CtxType_Float, CurValue: 2.0},
		}},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Float, FValue: 10.0},
		ActType{OpCode: CmdDone},
	}},

	// - use /=
	{Test: "0037", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdInsertInContext, Id: "abc", DataType: CtxType_Int, NValue: 20},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "abc"},
			tok.Token{TokNo: gen.Tok_DIV_EQ, Match: "/="},
			tok.Token{TokNo: gen.Tok_NUM, Match: "2", DataType: CtxType_Int, CurValue: 2},
		}},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Int, NValue: 10},
		ActType{OpCode: CmdDone},
	}},

	// - use %=
	{Test: "0038", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdInsertInContext, Id: "abc", DataType: CtxType_Int, NValue: 21},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "abc"},
			tok.Token{TokNo: gen.Tok_MOD_EQ, Match: "%="},
			tok.Token{TokNo: gen.Tok_NUM, Match: "2", DataType: CtxType_Int, CurValue: 2},
		}},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Int, NValue: 1},
		ActType{OpCode: CmdDone},
	}},

	// - use ^=
	{Test: "0039", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdInsertInContext, Id: "abc", DataType: CtxType_Int, NValue: 8},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "abc"},
			tok.Token{TokNo: gen.Tok_CAROT_EQ, Match: "^="},
			tok.Token{TokNo: gen.Tok_NUM, Match: "2", DataType: CtxType_Int, CurValue: 2},
		}},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Int, NValue: 10},
		ActType{OpCode: CmdDone},
	}},

	// - use <<=
	{Test: "0040", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdInsertInContext, Id: "abc", DataType: CtxType_Int, NValue: 8},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "abc"},
			tok.Token{TokNo: gen.Tok_S_L_EQ, Match: "<<="},
			tok.Token{TokNo: gen.Tok_NUM, Match: "2", DataType: CtxType_Int, CurValue: 2},
		}},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Int, NValue: 32},
		ActType{OpCode: CmdDone},
	}},

	// - use >>=
	{Test: "0041", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdInsertInContext, Id: "abc", DataType: CtxType_Int, NValue: 8},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "abc"},
			tok.Token{TokNo: gen.Tok_S_R_EQ, Match: ">>="},
			tok.Token{TokNo: gen.Tok_NUM, Match: "2", DataType: CtxType_Int, CurValue: 2},
		}},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Int, NValue: 2},
		ActType{OpCode: CmdDone},
	}},

	// - 1st test of calling nullFunc
	{Test: "0042", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdInsertInContext, Id: "abc", DataType: CtxType_Int, NValue: 1},
		ActType{OpCode: CmdInsertInContext, Id: "def", DataType: CtxType_Int, NValue: 2},
		ActType{OpCode: CmdInsertInContext, Id: "ghi", DataType: CtxType_Int, NValue: 3},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "nullFuncSSN"},
			tok.Token{TokNo: gen.Tok_OP, Match: "("},
			tok.Token{TokNo: gen.Tok_Str0, Match: "p1", DataType: CtxType_Str, CurValue: "p1"},
			tok.Token{TokNo: gen.Tok_COMMA, Match: ","},
			tok.Token{TokNo: gen.Tok_Str0, Match: "p2", DataType: CtxType_Str, CurValue: "p2"},
			tok.Token{TokNo: gen.Tok_COMMA, Match: ","},
			tok.Token{TokNo: gen.Tok_NUM, Match: "3", DataType: CtxType_Int, CurValue: 3},
			tok.Token{TokNo: gen.Tok_CL, Match: ")"},
		}},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Int, NValue: 0},
		ActType{OpCode: CmdDone},
	}},

	// - 2nd test of calling nullFunc
	{Test: "0043", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdInsertInContext, Id: "abc", DataType: CtxType_Int, NValue: 1},
		ActType{OpCode: CmdInsertInContext, Id: "def", DataType: CtxType_Int, NValue: 2},
		ActType{OpCode: CmdInsertInContext, Id: "ghi", DataType: CtxType_Int, NValue: 3},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "nullFuncSSB"},
			tok.Token{TokNo: gen.Tok_OP, Match: "("},
			tok.Token{TokNo: gen.Tok_Str0, Match: "p1", DataType: CtxType_Str, CurValue: "p1"},
			tok.Token{TokNo: gen.Tok_COMMA, Match: ","},
			tok.Token{TokNo: gen.Tok_Str0, Match: "def", DataType: CtxType_ID, CurValue: "def"},
			tok.Token{TokNo: gen.Tok_COMMA, Match: ","},
			tok.Token{TokNo: gen.Tok_NUM, Match: "3", DataType: CtxType_Int, CurValue: 3},
			tok.Token{TokNo: gen.Tok_L_EQ, Match: "=="},
			tok.Token{TokNo: gen.Tok_NUM, Match: "1", DataType: CtxType_Int, CurValue: 1},
			tok.Token{TokNo: gen.Tok_PLUS, Match: "+"},
			tok.Token{TokNo: gen.Tok_NUM, Match: "2", DataType: CtxType_Int, CurValue: 2},
			tok.Token{TokNo: gen.Tok_COMMA, Match: ","},
			tok.Token{TokNo: gen.Tok_CL, Match: ")"},
		}},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Int, NValue: 0},
		ActType{OpCode: CmdDone},
	}},

	// - 3rd test of calling nullFunc
	{Test: "0044", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdInsertInContext, Id: "abc", DataType: CtxType_Int, NValue: 1},
		ActType{OpCode: CmdInsertInContext, Id: "def", DataType: CtxType_Int, NValue: 2},
		ActType{OpCode: CmdInsertInContext, Id: "ghi", DataType: CtxType_Int, NValue: 3},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "nullFunc"},
			tok.Token{TokNo: gen.Tok_OP, Match: "("},
			tok.Token{TokNo: gen.Tok_COMMA, Match: ","},
			tok.Token{TokNo: gen.Tok_CL, Match: ")"},
		}},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Int, NValue: 0},
		ActType{OpCode: CmdDone},
	}},

	// - 4th test of calling nullFunc
	{Test: "0045", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdInsertInContext, Id: "abc", DataType: CtxType_Int, NValue: 1},
		ActType{OpCode: CmdInsertInContext, Id: "def", DataType: CtxType_Int, NValue: 2},
		ActType{OpCode: CmdInsertInContext, Id: "ghi", DataType: CtxType_Int, NValue: 3},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "nullFunc"},
			tok.Token{TokNo: gen.Tok_OP, Match: "("},
			tok.Token{TokNo: gen.Tok_CL, Match: ")"},
		}},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Int, NValue: 0},
		ActType{OpCode: CmdDone},
	}},

	{Test: "0046", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_OP_SQ, Match: "["},
			tok.Token{TokNo: gen.Tok_NUM, Match: "1", DataType: CtxType_Int, CurValue: 1},
			tok.Token{TokNo: gen.Tok_COMMA, Match: ","},
			tok.Token{TokNo: gen.Tok_NUM, Match: "2", DataType: CtxType_Int, CurValue: 2},
			tok.Token{TokNo: gen.Tok_COMMA, Match: ","},
			tok.Token{TokNo: gen.Tok_NUM, Match: "3", DataType: CtxType_Int, CurValue: 3},
			tok.Token{TokNo: gen.Tok_CL_SQ, Match: "]"},
		}},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Bool, BValue: true},
		ActType{OpCode: CmdDone},
	}},

	{Test: "0047", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_OP_BRACE, Match: "{"},
			tok.Token{TokNo: gen.Tok_ID, Match: "aaa"},
			tok.Token{TokNo: gen.Tok_COLON, Match: ":"},
			tok.Token{TokNo: gen.Tok_NUM, Match: "1", DataType: CtxType_Int, CurValue: 1},
			tok.Token{TokNo: gen.Tok_COMMA, Match: ","},

			tok.Token{TokNo: gen.Tok_Str0, Match: "aab", DataType: CtxType_Str, CurValue: "aab"},
			tok.Token{TokNo: gen.Tok_COLON, Match: ":"},
			tok.Token{TokNo: gen.Tok_NUM, Match: "2", DataType: CtxType_Int, CurValue: 2},
			tok.Token{TokNo: gen.Tok_COMMA, Match: ","},

			tok.Token{TokNo: gen.Tok_ID, Match: "aac", DataType: CtxType_Str, CurValue: "aac"},
			tok.Token{TokNo: gen.Tok_COLON, Match: ":"},
			tok.Token{TokNo: gen.Tok_NUM, Match: "3", DataType: CtxType_Int, CurValue: 3},
			tok.Token{TokNo: gen.Tok_CL_BRACE, Match: "}"},
		}},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Bool, BValue: true},
		ActType{OpCode: CmdDone},
	}},

	{Test: "0048", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdInsertInContext, Id: "abc", DataType: CtxType_Float, FValue: 5},
		ActType{OpCode: CmdInsertInContext, Id: "def", DataType: CtxType_Float, FValue: 10},
		ActType{OpCode: CmdInsertInContext, Id: "iii", DataType: CtxType_Int, NValue: 2},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "iii"},
			tok.Token{TokNo: gen.Tok_PLUS, Match: "+"},
			tok.Token{TokNo: gen.Tok_ID, Match: "abc"}, // F: 7.0
			tok.Token{TokNo: gen.Tok_MINUS, Match: "-"},
			tok.Token{TokNo: gen.Tok_ID, Match: "def"}, // F: -3.0
			tok.Token{TokNo: gen.Tok_MINUS, Match: "-"},
			tok.Token{TokNo: gen.Tok_ID, Match: "iii"}, // F: -5.0
		}},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Float, FValue: -5},
		ActType{OpCode: CmdDone},
	}},

	{Test: "0049", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdInsertInContext, Id: "abc", DataType: CtxType_Float, FValue: 5},
		ActType{OpCode: CmdInsertInContext, Id: "def", DataType: CtxType_Float, FValue: 10},
		ActType{OpCode: CmdInsertInContext, Id: "iii", DataType: CtxType_Int, NValue: 2},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "iii"},
			tok.Token{TokNo: gen.Tok_MINUS, Match: "-"},
			tok.Token{TokNo: gen.Tok_ID, Match: "abc"},
			tok.Token{TokNo: gen.Tok_MINUS, Match: "-"},
			tok.Token{TokNo: gen.Tok_ID, Match: "def"},
			tok.Token{TokNo: gen.Tok_MINUS, Match: "-"},
			tok.Token{TokNo: gen.Tok_ID, Match: "iii"},
		}},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Float, FValue: -15},
		ActType{OpCode: CmdDone},
	}},

	{Test: "0050", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdInsertInContext, Id: "abc", DataType: CtxType_Float, FValue: 5},
		ActType{OpCode: CmdInsertInContext, Id: "def", DataType: CtxType_Float, FValue: 10},
		ActType{OpCode: CmdInsertInContext, Id: "iii", DataType: CtxType_Int, NValue: 2},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "iii"},
			tok.Token{TokNo: gen.Tok_STAR, Match: "*"},
			tok.Token{TokNo: gen.Tok_ID, Match: "abc"}, // 10
			tok.Token{TokNo: gen.Tok_SLASH, Match: "/"},
			tok.Token{TokNo: gen.Tok_ID, Match: "iii"}, // 5.0
			tok.Token{TokNo: gen.Tok_MINUS, Match: "-"},
			tok.Token{TokNo: gen.Tok_ID, Match: "iii"},
		}},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Float, FValue: 3},
		ActType{OpCode: CmdDone},
	}},

	// int / float, float / int etc.
	{Test: "0051", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdInsertInContext, Id: "abc", DataType: CtxType_Float, FValue: 5},
		ActType{OpCode: CmdInsertInContext, Id: "def", DataType: CtxType_Float, FValue: 10},
		ActType{OpCode: CmdInsertInContext, Id: "iii", DataType: CtxType_Int, NValue: 2},
		ActType{OpCode: CmdInsertInContext, Id: "one", DataType: CtxType_Float, FValue: 1},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "iii"},
			tok.Token{TokNo: gen.Tok_SLASH, Match: "/"},
			tok.Token{TokNo: gen.Tok_ID, Match: "abc"}, // 0.4
			tok.Token{TokNo: gen.Tok_STAR, Match: "*"},
			tok.Token{TokNo: gen.Tok_ID, Match: "iii"}, // 0.8
			tok.Token{TokNo: gen.Tok_STAR, Match: "*"},
			tok.Token{TokNo: gen.Tok_ID, Match: "one"}, // 0.8
			tok.Token{TokNo: gen.Tok_SLASH, Match: "/"},
			tok.Token{TokNo: gen.Tok_ID, Match: "one"}, // 0.8
		}},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Float, FValue: 0.8},
		ActType{OpCode: CmdDone},
	}},

	// Errors - check for errors 1st time
	{Test: "0052", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdInsertInContext, Id: "abc", DataType: CtxType_Float, FValue: 5},
		ActType{OpCode: CmdInsertInContext, Id: "def", DataType: CtxType_Float, FValue: 10},
		ActType{OpCode: CmdInsertInContext, Id: "iii", DataType: CtxType_Int, NValue: 2},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "iii"},
			tok.Token{TokNo: gen.Tok_PCT, Match: "%"},
			tok.Token{TokNo: gen.Tok_ID, Match: "abc"}, // error
		}},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Int, FValue: 0},
		ActType{OpCode: CmdErrorsContain, Id: "Eval00001", NValue: 1}, // One error containing string
		ActType{OpCode: CmdDone},
	}},

	// Errors - check for errors 1st time
	{Test: "0053", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdInsertInContext, Id: "abc", DataType: CtxType_Float, FValue: 5},
		ActType{OpCode: CmdInsertInContext, Id: "def", DataType: CtxType_Float, FValue: 10},
		ActType{OpCode: CmdInsertInContext, Id: "iii", DataType: CtxType_Int, NValue: 2},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "abc"}, // error
			tok.Token{TokNo: gen.Tok_PCT, Match: "%"},
			tok.Token{TokNo: gen.Tok_ID, Match: "iii"},
		}},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Int, FValue: 0},
		ActType{OpCode: CmdErrorsContain, Id: "Eval00001", NValue: 1}, // One error containing string
		ActType{OpCode: CmdDone},
	}},

	// Errors - check for errors 1st time
	{Test: "0054", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdInsertInContext, Id: "abc", DataType: CtxType_Float, FValue: 5},
		ActType{OpCode: CmdInsertInContext, Id: "def", DataType: CtxType_Float, FValue: 10},
		ActType{OpCode: CmdInsertInContext, Id: "iii", DataType: CtxType_Int, NValue: 2},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "abc"}, // error
			tok.Token{TokNo: gen.Tok_PCT, Match: "%"},
			tok.Token{TokNo: gen.Tok_ID, Match: "def"},
		}},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Int, FValue: 0},
		ActType{OpCode: CmdErrorsContain, Id: "Eval00001", NValue: 1}, // One error containing string
		ActType{OpCode: CmdDone},
	}},

	// Errors - check for errors 1st time
	{Test: "0055", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdInsertInContext, Id: "abc", DataType: CtxType_Float, FValue: 5},
		ActType{OpCode: CmdInsertInContext, Id: "def", DataType: CtxType_Float, FValue: 10},
		ActType{OpCode: CmdInsertInContext, Id: "iii", DataType: CtxType_Int, NValue: 2},
		ActType{OpCode: CmdInsertInContext, Id: "sss", DataType: CtxType_Str, SValue: "sss"},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "iii"}, // error
			tok.Token{TokNo: gen.Tok_PCT, Match: "%"},
			tok.Token{TokNo: gen.Tok_ID, Match: "sss"},
		}},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Int, FValue: 0},
		ActType{OpCode: CmdErrorsContain, Id: "Eval00002", NValue: 1}, // One error containing string
		ActType{OpCode: CmdDone},
	}},

	// + / - on strings
	{Test: "0056", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdInsertInContext, Id: "abc", DataType: CtxType_Float, FValue: 5},
		ActType{OpCode: CmdInsertInContext, Id: "def", DataType: CtxType_Float, FValue: 10},
		ActType{OpCode: CmdInsertInContext, Id: "iii", DataType: CtxType_Int, NValue: 2},
		ActType{OpCode: CmdInsertInContext, Id: "sss", DataType: CtxType_Str, SValue: "sss"},
		ActType{OpCode: CmdInsertInContext, Id: "ttt", DataType: CtxType_Str, SValue: "ttt"},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "sss"}, // error
			tok.Token{TokNo: gen.Tok_PLUS, Match: "+"},
			tok.Token{TokNo: gen.Tok_ID, Match: "ttt"},
		}},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Str, SValue: "sssttt"},
		// ActType{OpCode: CmdErrorsContain, Id: "Eval00002", NValue: 1}, // One error containing string
		ActType{OpCode: CmdDone},
	}},

	// - on strings
	{Test: "0057", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdInsertInContext, Id: "abc", DataType: CtxType_Float, FValue: 5},
		ActType{OpCode: CmdInsertInContext, Id: "def", DataType: CtxType_Float, FValue: 10},
		ActType{OpCode: CmdInsertInContext, Id: "iii", DataType: CtxType_Int, NValue: 2},
		ActType{OpCode: CmdInsertInContext, Id: "sss", DataType: CtxType_Str, SValue: "sss"},
		ActType{OpCode: CmdInsertInContext, Id: "ttt", DataType: CtxType_Str, SValue: "ttt"},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "sss"}, // error
			tok.Token{TokNo: gen.Tok_MINUS, Match: "-"},
			tok.Token{TokNo: gen.Tok_ID, Match: "ttt"},
		}},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Int, NValue: 0},
		ActType{OpCode: CmdErrorsContain, Id: "Eval00003", NValue: 1}, // One error containing string
		ActType{OpCode: CmdDone},
	}},

	// - on strings
	{Test: "0058", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdInsertInContext, Id: "abc", DataType: CtxType_Float, FValue: 5},
		ActType{OpCode: CmdInsertInContext, Id: "def", DataType: CtxType_Float, FValue: 10},
		ActType{OpCode: CmdInsertInContext, Id: "iii", DataType: CtxType_Int, NValue: 2},
		ActType{OpCode: CmdInsertInContext, Id: "sss", DataType: CtxType_Str, SValue: "sss"},
		ActType{OpCode: CmdInsertInContext, Id: "ttt", DataType: CtxType_Str, SValue: "ttt"},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "iii"}, // error
			tok.Token{TokNo: gen.Tok_MINUS, Match: "-"},
			tok.Token{TokNo: gen.Tok_ID, Match: "ttt"},
		}},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Int, NValue: 0},
		ActType{OpCode: CmdErrorsContain, Id: "Eval00004", NValue: 1}, // One error containing string
		ActType{OpCode: CmdDone},
	}},

	// case gen.Tok_S_L:
	// case gen.Tok_S_R:
	{Test: "0059", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdInsertInContext, Id: "abc", DataType: CtxType_Float, FValue: 5},
		ActType{OpCode: CmdInsertInContext, Id: "def", DataType: CtxType_Float, FValue: 10},
		ActType{OpCode: CmdInsertInContext, Id: "iii", DataType: CtxType_Int, NValue: 2},
		ActType{OpCode: CmdInsertInContext, Id: "jjj", DataType: CtxType_Int, NValue: 3},
		ActType{OpCode: CmdInsertInContext, Id: "sss", DataType: CtxType_Str, SValue: "sss"},
		ActType{OpCode: CmdInsertInContext, Id: "ttt", DataType: CtxType_Str, SValue: "ttt"},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "iii"}, // error
			tok.Token{TokNo: gen.Tok_S_L, Match: "<<"},
			tok.Token{TokNo: gen.Tok_ID, Match: "ttt"},
		}},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Int, NValue: 0},
		ActType{OpCode: CmdErrorsContain, Id: "Eval00005", NValue: 1}, // One error containing string
		ActType{OpCode: CmdDone},
	}},

	// if opTk != gen.Tok_LE && opTk != gen.Tok_LT && opTk != gen.Tok_GE && opTk != gen.Tok_GT {
	{Test: "0060", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdInsertInContext, Id: "abc", DataType: CtxType_Float, FValue: 5},
		ActType{OpCode: CmdInsertInContext, Id: "def", DataType: CtxType_Float, FValue: 10},
		ActType{OpCode: CmdInsertInContext, Id: "iii", DataType: CtxType_Int, NValue: 2},
		ActType{OpCode: CmdInsertInContext, Id: "jjj", DataType: CtxType_Int, NValue: 3},
		ActType{OpCode: CmdInsertInContext, Id: "sss", DataType: CtxType_Str, SValue: "sss"},
		ActType{OpCode: CmdInsertInContext, Id: "ttt", DataType: CtxType_Str, SValue: "ttt"},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "iii"}, // error
			tok.Token{TokNo: gen.Tok_LE, Match: "<="},
			tok.Token{TokNo: gen.Tok_ID, Match: "abc"},
		}},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Bool, BValue: true},
		// ActType{OpCode: CmdErrorsContain, Id: "Eval00005", NValue: 1}, // One error containing string
		ActType{OpCode: CmdDone},
	}},

	// if opTk != gen.Tok_LE && opTk != gen.Tok_LT && opTk != gen.Tok_GE && opTk != gen.Tok_GT {
	{Test: "0061", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdInsertInContext, Id: "abc", DataType: CtxType_Float, FValue: 5},
		ActType{OpCode: CmdInsertInContext, Id: "def", DataType: CtxType_Float, FValue: 10},
		ActType{OpCode: CmdInsertInContext, Id: "iii", DataType: CtxType_Int, NValue: 2},
		ActType{OpCode: CmdInsertInContext, Id: "jjj", DataType: CtxType_Int, NValue: 3},
		ActType{OpCode: CmdInsertInContext, Id: "sss", DataType: CtxType_Str, SValue: "sss"},
		ActType{OpCode: CmdInsertInContext, Id: "ttt", DataType: CtxType_Str, SValue: "ttt"},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "sss"}, // error
			tok.Token{TokNo: gen.Tok_LE, Match: "<="},
			tok.Token{TokNo: gen.Tok_ID, Match: "ttt"},
		}},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Bool, BValue: true},
		// ActType{OpCode: CmdErrorsContain, Id: "Eval00005", NValue: 1}, // One error containing string
		ActType{OpCode: CmdDone},
	}},

	// if opTk != gen.Tok_LE && opTk != gen.Tok_LT && opTk != gen.Tok_GE && opTk != gen.Tok_GT {
	{Test: "0062", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdInsertInContext, Id: "abc", DataType: CtxType_Float, FValue: 5},
		ActType{OpCode: CmdInsertInContext, Id: "def", DataType: CtxType_Float, FValue: 10},
		ActType{OpCode: CmdInsertInContext, Id: "iii", DataType: CtxType_Int, NValue: 2},
		ActType{OpCode: CmdInsertInContext, Id: "jjj", DataType: CtxType_Int, NValue: 3},
		ActType{OpCode: CmdInsertInContext, Id: "sss", DataType: CtxType_Str, SValue: "sss"},
		ActType{OpCode: CmdInsertInContext, Id: "ttt", DataType: CtxType_Str, SValue: "ttt"},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "sss"}, // error
			tok.Token{TokNo: gen.Tok_LT, Match: "<"},
			tok.Token{TokNo: gen.Tok_ID, Match: "ttt"},
		}},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Bool, BValue: true},
		// ActType{OpCode: CmdErrorsContain, Id: "Eval00005", NValue: 1}, // One error containing string
		ActType{OpCode: CmdDone},
	}},

	// if opTk != gen.Tok_LE && opTk != gen.Tok_LT && opTk != gen.Tok_GE && opTk != gen.Tok_GT {
	{Test: "0063", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdInsertInContext, Id: "abc", DataType: CtxType_Float, FValue: 5},
		ActType{OpCode: CmdInsertInContext, Id: "def", DataType: CtxType_Float, FValue: 10},
		ActType{OpCode: CmdInsertInContext, Id: "iii", DataType: CtxType_Int, NValue: 2},
		ActType{OpCode: CmdInsertInContext, Id: "jjj", DataType: CtxType_Int, NValue: 3},
		ActType{OpCode: CmdInsertInContext, Id: "sss", DataType: CtxType_Str, SValue: "sss"},
		ActType{OpCode: CmdInsertInContext, Id: "ttt", DataType: CtxType_Str, SValue: "ttt"},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "sss"}, // error
			tok.Token{TokNo: gen.Tok_GT, Match: ">"},
			tok.Token{TokNo: gen.Tok_ID, Match: "ttt"},
		}},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Bool, BValue: false},
		// ActType{OpCode: CmdErrorsContain, Id: "Eval00005", NValue: 1}, // One error containing string
		ActType{OpCode: CmdDone},
	}},

	// if opTk != gen.Tok_LE && opTk != gen.Tok_LT && opTk != gen.Tok_GE && opTk != gen.Tok_GT {
	{Test: "0064", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdInsertInContext, Id: "abc", DataType: CtxType_Float, FValue: 5},
		ActType{OpCode: CmdInsertInContext, Id: "def", DataType: CtxType_Float, FValue: 10},
		ActType{OpCode: CmdInsertInContext, Id: "iii", DataType: CtxType_Int, NValue: 2},
		ActType{OpCode: CmdInsertInContext, Id: "jjj", DataType: CtxType_Int, NValue: 3},
		ActType{OpCode: CmdInsertInContext, Id: "sss", DataType: CtxType_Str, SValue: "sss"},
		ActType{OpCode: CmdInsertInContext, Id: "ttt", DataType: CtxType_Str, SValue: "ttt"},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "sss"}, // error
			tok.Token{TokNo: gen.Tok_GE, Match: ">="},
			tok.Token{TokNo: gen.Tok_ID, Match: "ttt"},
		}},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Bool, BValue: false},
		// ActType{OpCode: CmdErrorsContain, Id: "Eval00005", NValue: 1}, // One error containing string
		ActType{OpCode: CmdDone},
	}},

	// if opTk != gen.Tok_LE && opTk != gen.Tok_LT && opTk != gen.Tok_GE && opTk != gen.Tok_GT {
	{Test: "0065", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdInsertInContext, Id: "abc", DataType: CtxType_Float, FValue: 5},
		ActType{OpCode: CmdInsertInContext, Id: "def", DataType: CtxType_Float, FValue: 10},
		ActType{OpCode: CmdInsertInContext, Id: "iii", DataType: CtxType_Int, NValue: 2},
		ActType{OpCode: CmdInsertInContext, Id: "jjj", DataType: CtxType_Int, NValue: 3},
		ActType{OpCode: CmdInsertInContext, Id: "sss", DataType: CtxType_Str, SValue: "sss"},
		ActType{OpCode: CmdInsertInContext, Id: "ttt", DataType: CtxType_Str, SValue: "ttt"},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "sss"}, // error
			tok.Token{TokNo: gen.Tok_GE, Match: ">="},
			tok.Token{TokNo: gen.Tok_ID, Match: "iii"},
		}},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Int, NValue: 0},
		ActType{OpCode: CmdErrorsContain, Id: "Eval00006", NValue: 1}, // One error containing string
		ActType{OpCode: CmdDone},
	}},

	// if opTk != gen.Tok_LE && opTk != gen.Tok_LT && opTk != gen.Tok_GE && opTk != gen.Tok_GT {
	{Test: "0066", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdInsertInContext, Id: "abc", DataType: CtxType_Float, FValue: 5},
		ActType{OpCode: CmdInsertInContext, Id: "def", DataType: CtxType_Float, FValue: 10},
		ActType{OpCode: CmdInsertInContext, Id: "iii", DataType: CtxType_Int, NValue: 2},
		ActType{OpCode: CmdInsertInContext, Id: "jjj", DataType: CtxType_Int, NValue: 3},
		ActType{OpCode: CmdInsertInContext, Id: "sss", DataType: CtxType_Str, SValue: "sss"},
		ActType{OpCode: CmdInsertInContext, Id: "ttt", DataType: CtxType_Str, SValue: "ttt"},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "abc"}, // error
			tok.Token{TokNo: gen.Tok_LT, Match: "<"},
			tok.Token{TokNo: gen.Tok_ID, Match: "def"},
		}},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Bool, BValue: true},
		// ActType{OpCode: CmdErrorsContain, Id: "Eval00005", NValue: 1}, // One error containing string
		ActType{OpCode: CmdDone},
	}},

	// if opTk != gen.Tok_LE && opTk != gen.Tok_LT && opTk != gen.Tok_GE && opTk != gen.Tok_GT {
	{Test: "0067", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdInsertInContext, Id: "abc", DataType: CtxType_Float, FValue: 5},
		ActType{OpCode: CmdInsertInContext, Id: "def", DataType: CtxType_Float, FValue: 10},
		ActType{OpCode: CmdInsertInContext, Id: "iii", DataType: CtxType_Int, NValue: 2},
		ActType{OpCode: CmdInsertInContext, Id: "jjj", DataType: CtxType_Int, NValue: 3},
		ActType{OpCode: CmdInsertInContext, Id: "sss", DataType: CtxType_Str, SValue: "sss"},
		ActType{OpCode: CmdInsertInContext, Id: "ttt", DataType: CtxType_Str, SValue: "ttt"},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "abc"}, // error
			tok.Token{TokNo: gen.Tok_GT, Match: ">"},
			tok.Token{TokNo: gen.Tok_ID, Match: "def"},
		}},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Bool, BValue: false},
		// ActType{OpCode: CmdErrorsContain, Id: "Eval00005", NValue: 1}, // One error containing string
		ActType{OpCode: CmdDone},
	}},

	// if opTk != gen.Tok_LE && opTk != gen.Tok_LT && opTk != gen.Tok_GE && opTk != gen.Tok_GT {
	{Test: "0068", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdInsertInContext, Id: "abc", DataType: CtxType_Float, FValue: 5},
		ActType{OpCode: CmdInsertInContext, Id: "def", DataType: CtxType_Float, FValue: 10},
		ActType{OpCode: CmdInsertInContext, Id: "iii", DataType: CtxType_Int, NValue: 2},
		ActType{OpCode: CmdInsertInContext, Id: "jjj", DataType: CtxType_Int, NValue: 3},
		ActType{OpCode: CmdInsertInContext, Id: "sss", DataType: CtxType_Str, SValue: "sss"},
		ActType{OpCode: CmdInsertInContext, Id: "ttt", DataType: CtxType_Str, SValue: "ttt"},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "def"},
			tok.Token{TokNo: gen.Tok_GE, Match: ">="},
			tok.Token{TokNo: gen.Tok_ID, Match: "iii"},
		}},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Bool, BValue: true},
		// ActType{OpCode: CmdErrorsContain, Id: "Eval00005", NValue: 1}, // One error containing string
		ActType{OpCode: CmdDone},
	}},

	// if opTk != gen.Tok_LE && opTk != gen.Tok_LT && opTk != gen.Tok_GE && opTk != gen.Tok_GT {
	{Test: "0069", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdInsertInContext, Id: "abc", DataType: CtxType_Float, FValue: 5},
		ActType{OpCode: CmdInsertInContext, Id: "def", DataType: CtxType_Float, FValue: 10},
		ActType{OpCode: CmdInsertInContext, Id: "iii", DataType: CtxType_Int, NValue: 2},
		ActType{OpCode: CmdInsertInContext, Id: "jjj", DataType: CtxType_Int, NValue: 3},
		ActType{OpCode: CmdInsertInContext, Id: "sss", DataType: CtxType_Str, SValue: "sss"},
		ActType{OpCode: CmdInsertInContext, Id: "ttt", DataType: CtxType_Str, SValue: "ttt"},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "iii"},
			tok.Token{TokNo: gen.Tok_LT, Match: "<"},
			tok.Token{TokNo: gen.Tok_ID, Match: "jjj"},
		}},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Bool, BValue: true},
		// ActType{OpCode: CmdErrorsContain, Id: "Eval00005", NValue: 1}, // One error containing string
		ActType{OpCode: CmdDone},
	}},

	// if opTk != gen.Tok_NE && opTk != gen.Tok_L_EQ {
	{Test: "0070", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdInsertInContext, Id: "abc", DataType: CtxType_Float, FValue: 5},
		ActType{OpCode: CmdInsertInContext, Id: "def", DataType: CtxType_Float, FValue: 10},
		ActType{OpCode: CmdInsertInContext, Id: "iii", DataType: CtxType_Int, NValue: 2},
		ActType{OpCode: CmdInsertInContext, Id: "jjj", DataType: CtxType_Int, NValue: 3},
		ActType{OpCode: CmdInsertInContext, Id: "sss", DataType: CtxType_Str, SValue: "sss"},
		ActType{OpCode: CmdInsertInContext, Id: "ttt", DataType: CtxType_Str, SValue: "ttt"},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "def"},
			tok.Token{TokNo: gen.Tok_NE, Match: "!="},
			tok.Token{TokNo: gen.Tok_ID, Match: "iii"},
		}},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Bool, BValue: true},
		// ActType{OpCode: CmdErrorsContain, Id: "Eval00005", NValue: 1}, // One error containing string
		ActType{OpCode: CmdDone},
	}},

	//
	{Test: "0071", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdInsertInContext, Id: "abc", DataType: CtxType_Float, FValue: 5},
		ActType{OpCode: CmdInsertInContext, Id: "def", DataType: CtxType_Float, FValue: 10},
		ActType{OpCode: CmdInsertInContext, Id: "iii", DataType: CtxType_Int, NValue: 2},
		ActType{OpCode: CmdInsertInContext, Id: "jjj", DataType: CtxType_Int, NValue: 3},
		ActType{OpCode: CmdInsertInContext, Id: "sss", DataType: CtxType_Str, SValue: "sss"},
		ActType{OpCode: CmdInsertInContext, Id: "ttt", DataType: CtxType_Str, SValue: "ttt"},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "jjj"},
			tok.Token{TokNo: gen.Tok_L_EQ, Match: "=="},
			tok.Token{TokNo: gen.Tok_ID, Match: "abc"},
		}},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Int, NValue: 0},
		// ActType{OpCode: CmdErrorsContain, Id: "Eval00007", NValue: 1}, // One error containing string
		ActType{OpCode: CmdDone},
	}},

	//
	{Test: "0072", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdInsertInContext, Id: "abc", DataType: CtxType_Float, FValue: 5},
		ActType{OpCode: CmdInsertInContext, Id: "def", DataType: CtxType_Float, FValue: 10},
		ActType{OpCode: CmdInsertInContext, Id: "iii", DataType: CtxType_Int, NValue: 2},
		ActType{OpCode: CmdInsertInContext, Id: "jjj", DataType: CtxType_Int, NValue: 3},
		ActType{OpCode: CmdInsertInContext, Id: "sss", DataType: CtxType_Str, SValue: "sss"},
		ActType{OpCode: CmdInsertInContext, Id: "ttt", DataType: CtxType_Str, SValue: "ttt"},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "sss"},
			tok.Token{TokNo: gen.Tok_L_EQ, Match: "=="},
			tok.Token{TokNo: gen.Tok_ID, Match: "ttt"},
		}},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Bool, BValue: false},
		// ActType{OpCode: CmdErrorsContain, Id: "Eval00005", NValue: 1}, // One error containing string
		ActType{OpCode: CmdDone},
	}},

	//
	{Test: "0073", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdInsertInContext, Id: "abc", DataType: CtxType_Float, FValue: 5},
		ActType{OpCode: CmdInsertInContext, Id: "def", DataType: CtxType_Float, FValue: 10},
		ActType{OpCode: CmdInsertInContext, Id: "iii", DataType: CtxType_Int, NValue: 2},
		ActType{OpCode: CmdInsertInContext, Id: "jjj", DataType: CtxType_Int, NValue: 3},
		ActType{OpCode: CmdInsertInContext, Id: "sss", DataType: CtxType_Str, SValue: "sss"},
		ActType{OpCode: CmdInsertInContext, Id: "ttt", DataType: CtxType_Str, SValue: "ttt"},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "sss"},
			tok.Token{TokNo: gen.Tok_NE, Match: "<>"},
			tok.Token{TokNo: gen.Tok_ID, Match: "ttt"},
		}},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Bool, BValue: true},
		// ActType{OpCode: CmdErrorsContain, Id: "Eval00005", NValue: 1}, // One error containing string
		ActType{OpCode: CmdDone},
	}},

	//
	{Test: "0074", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdInsertInContext, Id: "abc", DataType: CtxType_Float, FValue: 5},
		ActType{OpCode: CmdInsertInContext, Id: "def", DataType: CtxType_Float, FValue: 10},
		ActType{OpCode: CmdInsertInContext, Id: "iii", DataType: CtxType_Int, NValue: 2},
		ActType{OpCode: CmdInsertInContext, Id: "jjj", DataType: CtxType_Int, NValue: 3},
		ActType{OpCode: CmdInsertInContext, Id: "sss", DataType: CtxType_Str, SValue: "sss"},
		ActType{OpCode: CmdInsertInContext, Id: "ttt", DataType: CtxType_Str, SValue: "ttt"},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "jjj"},
			tok.Token{TokNo: gen.Tok_L_EQ, Match: "=="},
			tok.Token{TokNo: gen.Tok_ID, Match: "sss"},
		}},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Int, NValue: 0},
		ActType{OpCode: CmdErrorsContain, Id: "Eval00007", NValue: 1}, // One error containing string
		ActType{OpCode: CmdDone},
	}},

	//	if opTk != gen.Tok_B_AND {
	//	if opTk != gen.Tok_B_OR {
	{Test: "0075", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdInsertInContext, Id: "abc", DataType: CtxType_Float, FValue: 5},
		ActType{OpCode: CmdInsertInContext, Id: "def", DataType: CtxType_Float, FValue: 10},
		ActType{OpCode: CmdInsertInContext, Id: "iii", DataType: CtxType_Int, NValue: 2},
		ActType{OpCode: CmdInsertInContext, Id: "jjj", DataType: CtxType_Int, NValue: 3},
		ActType{OpCode: CmdInsertInContext, Id: "sss", DataType: CtxType_Str, SValue: "sss"},
		ActType{OpCode: CmdInsertInContext, Id: "ttt", DataType: CtxType_Str, SValue: "ttt"},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "iii"},
			tok.Token{TokNo: gen.Tok_B_AND, Match: "&"},
			tok.Token{TokNo: gen.Tok_ID, Match: "jjj"},
		}},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Int, NValue: 2},
		// ActType{OpCode: CmdErrorsContain, Id: "Eval00007", NValue: 1}, // One error containing string
		ActType{OpCode: CmdDone},
	}},

	//	if opTk != gen.Tok_B_AND {
	//	if opTk != gen.Tok_B_OR {
	{Test: "0076", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdInsertInContext, Id: "abc", DataType: CtxType_Float, FValue: 5},
		ActType{OpCode: CmdInsertInContext, Id: "def", DataType: CtxType_Float, FValue: 10},
		ActType{OpCode: CmdInsertInContext, Id: "iii", DataType: CtxType_Int, NValue: 2},
		ActType{OpCode: CmdInsertInContext, Id: "jjj", DataType: CtxType_Int, NValue: 3},
		ActType{OpCode: CmdInsertInContext, Id: "sss", DataType: CtxType_Str, SValue: "sss"},
		ActType{OpCode: CmdInsertInContext, Id: "ttt", DataType: CtxType_Str, SValue: "ttt"},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "iii"},
			tok.Token{TokNo: gen.Tok_B_OR, Match: "bor"},
			tok.Token{TokNo: gen.Tok_ID, Match: "jjj"},
		}},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Int, NValue: 3},
		// ActType{OpCode: CmdErrorsContain, Id: "Eval00007", NValue: 1}, // One error containing string
		ActType{OpCode: CmdDone},
	}},

	//	if opTk != gen.Tok_B_AND {
	//	if opTk != gen.Tok_B_OR {
	{Test: "0077", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdInsertInContext, Id: "abc", DataType: CtxType_Float, FValue: 5},
		ActType{OpCode: CmdInsertInContext, Id: "def", DataType: CtxType_Float, FValue: 10},
		ActType{OpCode: CmdInsertInContext, Id: "iii", DataType: CtxType_Int, NValue: 2},
		ActType{OpCode: CmdInsertInContext, Id: "jjj", DataType: CtxType_Int, NValue: 3},
		ActType{OpCode: CmdInsertInContext, Id: "sss", DataType: CtxType_Str, SValue: "sss"},
		ActType{OpCode: CmdInsertInContext, Id: "ttt", DataType: CtxType_Str, SValue: "ttt"},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "iii"},
			tok.Token{TokNo: gen.Tok_B_AND, Match: "&"},
			tok.Token{TokNo: gen.Tok_ID, Match: "abc"},
		}},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Int, NValue: 0},
		ActType{OpCode: CmdErrorsContain, Id: "Eval00008", NValue: 1}, // One error containing string
		ActType{OpCode: CmdDone},
	}},

	//	if opTk != gen.Tok_B_AND {
	//	if opTk != gen.Tok_B_OR {
	{Test: "0078", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdInsertInContext, Id: "abc", DataType: CtxType_Float, FValue: 5},
		ActType{OpCode: CmdInsertInContext, Id: "def", DataType: CtxType_Float, FValue: 10},
		ActType{OpCode: CmdInsertInContext, Id: "iii", DataType: CtxType_Int, NValue: 2},
		ActType{OpCode: CmdInsertInContext, Id: "jjj", DataType: CtxType_Int, NValue: 3},
		ActType{OpCode: CmdInsertInContext, Id: "sss", DataType: CtxType_Str, SValue: "sss"},
		ActType{OpCode: CmdInsertInContext, Id: "ttt", DataType: CtxType_Str, SValue: "ttt"},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "iii"},
			tok.Token{TokNo: gen.Tok_B_OR, Match: "bor"},
			tok.Token{TokNo: gen.Tok_ID, Match: "abc"},
		}},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Int, NValue: 0},
		ActType{OpCode: CmdErrorsContain, Id: "Eval00009", NValue: 1}, // One error containing string
		ActType{OpCode: CmdDone},
	}},

	//	^ xor
	{Test: "0079", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdInsertInContext, Id: "abc", DataType: CtxType_Float, FValue: 5},
		ActType{OpCode: CmdInsertInContext, Id: "def", DataType: CtxType_Float, FValue: 10},
		ActType{OpCode: CmdInsertInContext, Id: "iii", DataType: CtxType_Int, NValue: 2},
		ActType{OpCode: CmdInsertInContext, Id: "jjj", DataType: CtxType_Int, NValue: 3},
		ActType{OpCode: CmdInsertInContext, Id: "sss", DataType: CtxType_Str, SValue: "sss"},
		ActType{OpCode: CmdInsertInContext, Id: "ttt", DataType: CtxType_Str, SValue: "ttt"},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "iii"},
			tok.Token{TokNo: gen.Tok_CARRET, Match: "^"},
			tok.Token{TokNo: gen.Tok_ID, Match: "jjj"},
		}},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Int, NValue: 1},
		ActType{OpCode: CmdDone},
	}},

	//	^ xor
	{Test: "0080", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdInsertInContext, Id: "abc", DataType: CtxType_Float, FValue: 5},
		ActType{OpCode: CmdInsertInContext, Id: "def", DataType: CtxType_Float, FValue: 10},
		ActType{OpCode: CmdInsertInContext, Id: "iii", DataType: CtxType_Int, NValue: 2},
		ActType{OpCode: CmdInsertInContext, Id: "jjj", DataType: CtxType_Int, NValue: 3},
		ActType{OpCode: CmdInsertInContext, Id: "sss", DataType: CtxType_Str, SValue: "sss"},
		ActType{OpCode: CmdInsertInContext, Id: "ttt", DataType: CtxType_Str, SValue: "ttt"},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "iii"},
			tok.Token{TokNo: gen.Tok_CARRET, Match: "^"},
			tok.Token{TokNo: gen.Tok_ID, Match: "abc"},
		}},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Int, NValue: 0},
		ActType{OpCode: CmdErrorsContain, Id: "Eval00010", NValue: 1}, // One error containing string
		ActType{OpCode: CmdDone},
	}},

	// &&
	{Test: "0081", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdInsertInContext, Id: "abc", DataType: CtxType_Float, FValue: 5},
		ActType{OpCode: CmdInsertInContext, Id: "def", DataType: CtxType_Float, FValue: 10},
		ActType{OpCode: CmdInsertInContext, Id: "iii", DataType: CtxType_Int, NValue: 2},
		ActType{OpCode: CmdInsertInContext, Id: "jjj", DataType: CtxType_Int, NValue: 3},
		ActType{OpCode: CmdInsertInContext, Id: "kkk", DataType: CtxType_Int, NValue: 2},
		ActType{OpCode: CmdInsertInContext, Id: "sss", DataType: CtxType_Str, SValue: "sss"},
		ActType{OpCode: CmdInsertInContext, Id: "ttt", DataType: CtxType_Str, SValue: "ttt"},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "iii"},
			tok.Token{TokNo: gen.Tok_L_EQ, Match: "=="},
			tok.Token{TokNo: gen.Tok_ID, Match: "kkk"},
			tok.Token{TokNo: gen.Tok_L_AND, Match: "&&"},
			tok.Token{TokNo: gen.Tok_ID, Match: "sss"},
			tok.Token{TokNo: gen.Tok_LT, Match: "<"},
			tok.Token{TokNo: gen.Tok_ID, Match: "ttt"},
		}},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Bool, BValue: true},
		// ActType{OpCode: CmdErrorsContain, Id: "Eval00007", NValue: 1}, // One error containing string
		ActType{OpCode: CmdDone},
	}},

	// && -- error
	{Test: "0082", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdInsertInContext, Id: "abc", DataType: CtxType_Float, FValue: 5},
		ActType{OpCode: CmdInsertInContext, Id: "def", DataType: CtxType_Float, FValue: 10},
		ActType{OpCode: CmdInsertInContext, Id: "iii", DataType: CtxType_Int, NValue: 2},
		ActType{OpCode: CmdInsertInContext, Id: "jjj", DataType: CtxType_Int, NValue: 3},
		ActType{OpCode: CmdInsertInContext, Id: "kkk", DataType: CtxType_Int, NValue: 2},
		ActType{OpCode: CmdInsertInContext, Id: "sss", DataType: CtxType_Str, SValue: "sss"},
		ActType{OpCode: CmdInsertInContext, Id: "ttt", DataType: CtxType_Str, SValue: "ttt"},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "iii"},
			tok.Token{TokNo: gen.Tok_L_EQ, Match: "=="},
			tok.Token{TokNo: gen.Tok_ID, Match: "kkk"},
			tok.Token{TokNo: gen.Tok_L_AND, Match: "&&"},
			tok.Token{TokNo: gen.Tok_ID, Match: "sss"},
		}},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Int, NValue: 0},
		ActType{OpCode: CmdErrorsContain, Id: "Eval00011", NValue: 1}, // One error containing string
		ActType{OpCode: CmdDone},
	}},

	// ||
	{Test: "0083", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdInsertInContext, Id: "abc", DataType: CtxType_Float, FValue: 5},
		ActType{OpCode: CmdInsertInContext, Id: "def", DataType: CtxType_Float, FValue: 10},
		ActType{OpCode: CmdInsertInContext, Id: "iii", DataType: CtxType_Int, NValue: 2},
		ActType{OpCode: CmdInsertInContext, Id: "jjj", DataType: CtxType_Int, NValue: 3},
		ActType{OpCode: CmdInsertInContext, Id: "kkk", DataType: CtxType_Int, NValue: 2},
		ActType{OpCode: CmdInsertInContext, Id: "sss", DataType: CtxType_Str, SValue: "sss"},
		ActType{OpCode: CmdInsertInContext, Id: "ttt", DataType: CtxType_Str, SValue: "ttt"},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "iii"},
			tok.Token{TokNo: gen.Tok_L_EQ, Match: "=="},
			tok.Token{TokNo: gen.Tok_ID, Match: "kkk"},
			tok.Token{TokNo: gen.Tok_L_OR, Match: "&&"},
			tok.Token{TokNo: gen.Tok_ID, Match: "sss"},
			tok.Token{TokNo: gen.Tok_GT, Match: ">"},
			tok.Token{TokNo: gen.Tok_ID, Match: "ttt"},
		}},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Bool, BValue: true},
		// ActType{OpCode: CmdErrorsContain, Id: "Eval00007", NValue: 1}, // One error containing string
		ActType{OpCode: CmdDone},
	}},

	// || -- error
	{Test: "0084", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdInsertInContext, Id: "abc", DataType: CtxType_Float, FValue: 5},
		ActType{OpCode: CmdInsertInContext, Id: "def", DataType: CtxType_Float, FValue: 10},
		ActType{OpCode: CmdInsertInContext, Id: "iii", DataType: CtxType_Int, NValue: 2},
		ActType{OpCode: CmdInsertInContext, Id: "jjj", DataType: CtxType_Int, NValue: 3},
		ActType{OpCode: CmdInsertInContext, Id: "kkk", DataType: CtxType_Int, NValue: 2},
		ActType{OpCode: CmdInsertInContext, Id: "sss", DataType: CtxType_Str, SValue: "sss"},
		ActType{OpCode: CmdInsertInContext, Id: "ttt", DataType: CtxType_Str, SValue: "ttt"},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "iii"},
			tok.Token{TokNo: gen.Tok_L_EQ, Match: "=="},
			tok.Token{TokNo: gen.Tok_ID, Match: "kkk"},
			tok.Token{TokNo: gen.Tok_L_OR, Match: "||"},
			tok.Token{TokNo: gen.Tok_ID, Match: "sss"},
		}},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Int, NValue: 0},
		ActType{OpCode: CmdErrorsContain, Id: "Eval00012", NValue: 1}, // One error containing string
		ActType{OpCode: CmdDone},
	}},

	// - use -=
	{Test: "0085", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdInsertInContext, Id: "abc", DataType: CtxType_Int, NValue: 5},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "abc"},
			tok.Token{TokNo: gen.Tok_MINUS_EQ, Match: "-="},
			tok.Token{TokNo: gen.Tok_NUM, Match: "22", DataType: CtxType_Int, CurValue: 22},
		}},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Int, NValue: -17},
		ActType{OpCode: CmdDone},
	}},

	// - use += (float)
	{Test: "0086", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdInsertInContext, Id: "abc", DataType: CtxType_Int, NValue: 5},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "abc"},
			tok.Token{TokNo: gen.Tok_MINUS_EQ, Match: "-="},
			tok.Token{TokNo: gen.Tok_Float, Match: "22.0", DataType: CtxType_Float, CurValue: 22.0},
		}},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Float, FValue: -17.0},
		ActType{OpCode: CmdDone},
	}},

	// - use |=
	{Test: "0087", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdInsertInContext, Id: "abc", DataType: CtxType_Int, NValue: 5},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "abc"},
			tok.Token{TokNo: gen.Tok_B_OR_EQ, Match: "|="},
			tok.Token{TokNo: gen.Tok_NUM, Match: "16", DataType: CtxType_Int, CurValue: 16},
		}},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Int, NValue: 21},
		ActType{OpCode: CmdDone},
	}},

	// - use &=
	{Test: "0088", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdInsertInContext, Id: "abc", DataType: CtxType_Int, NValue: 5},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "abc"},
			tok.Token{TokNo: gen.Tok_B_AND_EQ, Match: "|="},
			tok.Token{TokNo: gen.Tok_NUM, Match: "17", DataType: CtxType_Int, CurValue: 17},
		}},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Int, NValue: 1},
		ActType{OpCode: CmdDone},
	}},

	// - use ~=
	{Test: "0089", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdInsertInContext, Id: "abc", DataType: CtxType_Int, NValue: 5},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "abc"},
			tok.Token{TokNo: gen.Tok_TILDE_EQ, Match: "~="},
			tok.Token{TokNo: gen.Tok_NUM, Match: "17", DataType: CtxType_Int, CurValue: 17},
		}},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Int, NValue: 0},
		ActType{OpCode: CmdErrorsContain, Id: "Eval00017", NValue: 1}, // One error containing string
		ActType{OpCode: CmdDone},
	}},

	// - use = (Float)
	{Test: "0090", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdInsertInContext, Id: "abc", DataType: CtxType_Float, FValue: 5.0},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "abc"},
			tok.Token{TokNo: gen.Tok_EQ, Match: "="},
			tok.Token{TokNo: gen.Tok_NUM, Match: "5", DataType: CtxType_Int, CurValue: 5},
		}},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Float, FValue: 5},
		ActType{OpCode: CmdDone},
	}},

	// - use := (float)
	{Test: "0091", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "abc"},
			tok.Token{TokNo: gen.Tok_DCL_VAR, Match: ":="},
			tok.Token{TokNo: gen.Tok_Float, Match: "22.0", DataType: CtxType_Float, CurValue: -22.0},
		}},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Float, FValue: -22.0},
		ActType{OpCode: CmdDone},
	}},

	// - use |=
	{Test: "0092", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdInsertInContext, Id: "abc", DataType: CtxType_Float, FValue: 5.0},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "abc"},
			tok.Token{TokNo: gen.Tok_B_OR_EQ, Match: "|="},
			tok.Token{TokNo: gen.Tok_NUM, Match: "16", DataType: CtxType_Int, CurValue: 16},
		}},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Int, NValue: 0},
		ActType{OpCode: CmdErrorsContain, Id: "Eval00022", NValue: 1}, // One error containing string
		ActType{OpCode: CmdDone},
	}},

	// - use &=
	{Test: "0093", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdInsertInContext, Id: "abc", DataType: CtxType_Float, FValue: 5.0},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "abc"},
			tok.Token{TokNo: gen.Tok_B_AND_EQ, Match: "|="},
			tok.Token{TokNo: gen.Tok_NUM, Match: "17", DataType: CtxType_Int, CurValue: 17},
		}},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Int, NValue: 0},
		ActType{OpCode: CmdErrorsContain, Id: "Eval00023", NValue: 1}, // One error containing string
		ActType{OpCode: CmdDone},
	}},

	// - use ~=
	{Test: "0094", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdInsertInContext, Id: "abc", DataType: CtxType_Float, FValue: 5.0},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "abc"},
			tok.Token{TokNo: gen.Tok_TILDE_EQ, Match: "~="},
			tok.Token{TokNo: gen.Tok_NUM, Match: "17", DataType: CtxType_Int, CurValue: 17},
		}},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Int, NValue: 0},
		ActType{OpCode: CmdErrorsContain, Id: "Eval00024", NValue: 1}, // One error containing string
		ActType{OpCode: CmdDone},
	}},

	// - use ^=
	{Test: "0095", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdInsertInContext, Id: "abc", DataType: CtxType_Float, FValue: 5.0},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "abc"},
			tok.Token{TokNo: gen.Tok_CAROT_EQ, Match: "^="},
			tok.Token{TokNo: gen.Tok_NUM, Match: "17", DataType: CtxType_Int, CurValue: 17},
		}},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Int, NValue: 0},
		ActType{OpCode: CmdErrorsContain, Id: "Eval00021", NValue: 1}, // One error containing string
		ActType{OpCode: CmdDone},
	}},

	// - use %=
	{Test: "0096", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdInsertInContext, Id: "abc", DataType: CtxType_Float, FValue: 5.0},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "abc"},
			tok.Token{TokNo: gen.Tok_MOD_EQ, Match: "%="},
			tok.Token{TokNo: gen.Tok_NUM, Match: "17", DataType: CtxType_Int, CurValue: 17},
		}},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Int, NValue: 0},
		ActType{OpCode: CmdErrorsContain, Id: "Eval00020", NValue: 1}, // One error containing string
		ActType{OpCode: CmdDone},
	}},

	// - use /= (Float)
	{Test: "0097", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdInsertInContext, Id: "abc", DataType: CtxType_Float, FValue: 5.0},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "abc"},
			tok.Token{TokNo: gen.Tok_DIV_EQ, Match: "/="},
			tok.Token{TokNo: gen.Tok_NUM, Match: "5", DataType: CtxType_Int, CurValue: 5},
		}},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Float, FValue: 1},
		ActType{OpCode: CmdDone},
	}},

	// - use /= (Float)
	{Test: "0098", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdInsertInContext, Id: "abc", DataType: CtxType_Float, FValue: 5.0},
		ActType{OpCode: CmdInsertInContext, Id: "iii", DataType: CtxType_Int, NValue: 5},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "iii"},
			tok.Token{TokNo: gen.Tok_DIV_EQ, Match: "/="},
			tok.Token{TokNo: gen.Tok_Float, Match: "5", DataType: CtxType_Float, CurValue: 5.0},
		}},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Float, FValue: 1},
		ActType{OpCode: CmdDone},
	}},

	// - use /= (str)
	{Test: "0099", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdInsertInContext, Id: "abc", DataType: CtxType_Float, FValue: 5.0},
		ActType{OpCode: CmdInsertInContext, Id: "iii", DataType: CtxType_Int, NValue: 5},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "iii"},
			tok.Token{TokNo: gen.Tok_DIV_EQ, Match: "/="},
			tok.Token{TokNo: gen.Tok_Str0, Match: "aab", DataType: CtxType_Str, CurValue: "aab"},
		}},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Int, NValue: 0},
		ActType{OpCode: CmdErrorsContain, Id: "Eval00014", NValue: 1}, // One error containing string
		ActType{OpCode: CmdDone},
	}},

	// - use /= (str)
	{Test: "0100", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdInsertInContext, Id: "abc", DataType: CtxType_Float, FValue: 5.0},
		ActType{OpCode: CmdInsertInContext, Id: "iii", DataType: CtxType_Int, NValue: 5},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "abc"},
			tok.Token{TokNo: gen.Tok_DIV_EQ, Match: "/="},
			tok.Token{TokNo: gen.Tok_Str0, Match: "aab", DataType: CtxType_Str, CurValue: "aab"},
		}},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Int, NValue: 0},
		ActType{OpCode: CmdErrorsContain, Id: "Eval00015", NValue: 1}, // One error containing string
		ActType{OpCode: CmdDone},
	}},

	// - use /= (Float,flaot)
	{Test: "0101", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdInsertInContext, Id: "abc", DataType: CtxType_Float, FValue: 5.0},
		ActType{OpCode: CmdInsertInContext, Id: "iii", DataType: CtxType_Int, NValue: 5},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "abc"},
			tok.Token{TokNo: gen.Tok_DIV_EQ, Match: "/="},
			tok.Token{TokNo: gen.Tok_Float, Match: "5", DataType: CtxType_Float, CurValue: 5.0},
		}},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Float, FValue: 1},
		ActType{OpCode: CmdDone},
	}},

	// - use <<= ( float - error )
	{Test: "0102", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdInsertInContext, Id: "abc", DataType: CtxType_Float, FValue: 8},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "abc"},
			tok.Token{TokNo: gen.Tok_S_L_EQ, Match: "<<="},
			tok.Token{TokNo: gen.Tok_NUM, Match: "2", DataType: CtxType_Int, CurValue: 2},
		}},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Int, NValue: 0},
		ActType{OpCode: CmdErrorsContain, Id: "Eval00025", NValue: 1}, // One error containing string
		ActType{OpCode: CmdDone},
	}},

	// - use >>= ( float - error )
	{Test: "0103", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdInsertInContext, Id: "abc", DataType: CtxType_Float, FValue: 8},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "abc"},
			tok.Token{TokNo: gen.Tok_S_R_EQ, Match: ">>="},
			tok.Token{TokNo: gen.Tok_NUM, Match: "2", DataType: CtxType_Int, CurValue: 2},
		}},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Int, NValue: 0},
		ActType{OpCode: CmdErrorsContain, Id: "Eval00026", NValue: 1}, // One error containing string
		ActType{OpCode: CmdDone},
	}},

	// - use <<= ( float - error )
	{Test: "0104", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdInsertInContext, Id: "abc", DataType: CtxType_Float, FValue: 8},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "abc"},
			tok.Token{TokNo: gen.Tok_S_L_EQ, Match: "<<="},
			tok.Token{TokNo: gen.Tok_Float, Match: "2", DataType: CtxType_Float, CurValue: 8.0},
		}},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Int, NValue: 0},
		ActType{OpCode: CmdErrorsContain, Id: "Eval00025", NValue: 1}, // One error containing string
		ActType{OpCode: CmdDone},
	}},

	// - use >>= ( float - error )
	{Test: "0105", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdInsertInContext, Id: "abc", DataType: CtxType_Float, FValue: 8},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "abc"},
			tok.Token{TokNo: gen.Tok_S_R_EQ, Match: ">>="},
			tok.Token{TokNo: gen.Tok_Float, Match: "2", DataType: CtxType_Float, CurValue: 2},
		}},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Int, NValue: 0},
		ActType{OpCode: CmdErrorsContain, Id: "Eval00026", NValue: 1}, // One error containing string
		ActType{OpCode: CmdDone},
	}},

	// - use := (str)
	{Test: "0106", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "abc"},
			tok.Token{TokNo: gen.Tok_DCL_VAR, Match: ":="},
			tok.Token{TokNo: gen.Tok_Str0, Match: "p1", DataType: CtxType_Str, CurValue: "p1"},
		}},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Str, SValue: "p1"},
		ActType{OpCode: CmdDone},
	}},

	// - use := (bool)
	{Test: "0107", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "abc"},
			tok.Token{TokNo: gen.Tok_DCL_VAR, Match: ":="},
			tok.Token{TokNo: gen.Tok_true, Match: "true", DataType: CtxType_Bool, CurValue: true},
		}},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Str, BValue: true},
		ActType{OpCode: CmdDone},
	}},

	// - use = (str)
	{Test: "0108", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdInsertInContext, Id: "abc", DataType: CtxType_Str, SValue: "xxxxxx"},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "abc"},
			tok.Token{TokNo: gen.Tok_EQ, Match: "="},
			tok.Token{TokNo: gen.Tok_Str0, Match: "p1", DataType: CtxType_Str, CurValue: "p1"},
		}},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Str, SValue: "p1"},
		ActType{OpCode: CmdDone},
	}},

	// - use = (bool)
	{Test: "0109", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdInsertInContext, Id: "abc", DataType: CtxType_Bool, BValue: false},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "abc"},
			tok.Token{TokNo: gen.Tok_EQ, Match: "="},
			tok.Token{TokNo: gen.Tok_true, Match: "true", DataType: CtxType_Bool, CurValue: true},
		}},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Str, BValue: true},
		ActType{OpCode: CmdDone},
	}},

	// - use = (bool)
	{Test: "0110", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdInsertInContext, Id: "abc", DataType: CtxType_Bool, BValue: false},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "abc"},
			tok.Token{TokNo: gen.Tok_MINUS_EQ, Match: "-="},
			tok.Token{TokNo: gen.Tok_true, Match: "true", DataType: CtxType_Bool, CurValue: true},
		}},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Int, NValue: 0},
		ActType{OpCode: CmdErrorsContain, Id: "Eval00029", NValue: 1}, // One error containing string
		ActType{OpCode: CmdDone},
	}},

	// - use = (string)
	{Test: "0111", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdInsertInContext, Id: "abc", DataType: CtxType_Str, SValue: "xxxxxx"},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "abc"},
			tok.Token{TokNo: gen.Tok_MINUS_EQ, Match: "-="},
			tok.Token{TokNo: gen.Tok_true, Match: "true", DataType: CtxType_Bool, CurValue: true},
		}},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Int, NValue: 0},
		ActType{OpCode: CmdErrorsContain, Id: "Eval00029", NValue: 1}, // One error containing string
		ActType{OpCode: CmdDone},
	}},

	// Test Unary -, + on a string
	{Test: "0112", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdInsertInContext, Id: "abc", DataType: CtxType_Int, NValue: 5},
		ActType{OpCode: CmdInsertInContext, Id: "def", DataType: CtxType_Str, SValue: "xxxxxx"},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "abc"},
			tok.Token{TokNo: gen.Tok_PLUS, Match: "+"},
			tok.Token{TokNo: gen.Tok_MINUS, Match: "-"},
			tok.Token{TokNo: gen.Tok_MINUS, Match: "-"},
			tok.Token{TokNo: gen.Tok_PLUS, Match: "+"},
			tok.Token{TokNo: gen.Tok_MINUS, Match: "-"},
			tok.Token{TokNo: gen.Tok_ID, Match: "def"},
		}},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Int, NValue: 0},
		ActType{OpCode: CmdErrorsContain, Id: "Eval00030", NValue: 1}, // One error containing string
		ActType{OpCode: CmdDone},
	}},

	// Test Unary - error, on right of 29*
	{Test: "0113", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdInsertInContext, Id: "abc", DataType: CtxType_Int, NValue: 5},
		ActType{OpCode: CmdInsertInContext, Id: "def", DataType: CtxType_Str, SValue: "xxxxxx"},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "abc"},
			tok.Token{TokNo: gen.Tok_STAR, Match: "*"},
			tok.Token{TokNo: gen.Tok_MINUS, Match: "-"},
			tok.Token{TokNo: gen.Tok_ID, Match: "def"},
		}},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Int, NValue: 0},
		ActType{OpCode: CmdErrorsContain, Id: "Eval00036", NValue: 1}, // One error containing string
		ActType{OpCode: CmdDone},
	}},

	// Test Unary ! error, on right of 29*
	{Test: "0114", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdInsertInContext, Id: "abc", DataType: CtxType_Int, NValue: 5},
		ActType{OpCode: CmdInsertInContext, Id: "def", DataType: CtxType_Str, SValue: "xxxxxx"},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "abc"},
			tok.Token{TokNo: gen.Tok_STAR, Match: "*"},
			tok.Token{TokNo: gen.Tok_EXCLAM, Match: "!"},
			tok.Token{TokNo: gen.Tok_ID, Match: "def"},
		}},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Int, NValue: 0},
		ActType{OpCode: CmdErrorsContain, Id: "Eval00031", NValue: 1}, // One error containing string
		ActType{OpCode: CmdDone},
	}},

	// xyzzy1857 - incomplete
	{Test: "0115", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdInsertInContext, Id: "abc", DataType: CtxType_ArrayOf},
		ActType{OpCode: CmdInsertInContext, Id: "def", DataType: CtxType_Int, NValue: 0},
		ActType{OpCode: CmdInsertInContext, Id: "ghi", DataType: CtxType_Int, NValue: 0},
		ActType{OpCode: CmdInsertInContext, Id: "jkl", DataType: CtxType_Int, NValue: 0},
		ActType{OpCode: CmdInsertInContext, Id: "m01", DataType: CtxType_Int, NValue: 0},
		ActType{OpCode: CmdInsertInContext, Id: "m02", DataType: CtxType_Int, NValue: 0},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "abc"},
			tok.Token{TokNo: gen.Tok_EQ, Match: "="},
			tok.Token{TokNo: gen.Tok_OP_SQ, Match: "["},
			tok.Token{TokNo: gen.Tok_NUM, Match: "1", DataType: CtxType_Int, CurValue: 1},
			tok.Token{TokNo: gen.Tok_COMMA, Match: ","},
			tok.Token{TokNo: gen.Tok_NUM, Match: "2", DataType: CtxType_Int, CurValue: 2},
			tok.Token{TokNo: gen.Tok_COMMA, Match: ","},
			tok.Token{TokNo: gen.Tok_NUM, Match: "3", DataType: CtxType_Int, CurValue: 3},
			tok.Token{TokNo: gen.Tok_CL_SQ, Match: "]"},
		}},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "def"},
			tok.Token{TokNo: gen.Tok_EQ, Match: "="},
			tok.Token{TokNo: gen.Tok_ID, Match: "abc"},
			tok.Token{TokNo: gen.Tok_OP_SQ, Match: "["},
			tok.Token{TokNo: gen.Tok_NUM, Match: "0", DataType: CtxType_Int, CurValue: 0},
			tok.Token{TokNo: gen.Tok_COLON, Match: ":"},
			tok.Token{TokNo: gen.Tok_NUM, Match: "1", DataType: CtxType_Int, CurValue: 1},
			tok.Token{TokNo: gen.Tok_CL_SQ, Match: "]"},
			//			//
		}},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "ghi"},
			tok.Token{TokNo: gen.Tok_EQ, Match: "="},
			tok.Token{TokNo: gen.Tok_ID, Match: "abc"},
			tok.Token{TokNo: gen.Tok_OP_SQ, Match: "["},
			tok.Token{TokNo: gen.Tok_NUM, Match: "0", DataType: CtxType_Int, CurValue: 0},
			tok.Token{TokNo: gen.Tok_CL_SQ, Match: "]"},
			//			//
		}},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "jkl"},
			tok.Token{TokNo: gen.Tok_EQ, Match: "="},
			tok.Token{TokNo: gen.Tok_ID, Match: "abc"},
			tok.Token{TokNo: gen.Tok_OP_SQ, Match: "["},
			tok.Token{TokNo: gen.Tok_COLON, Match: ":"},
			tok.Token{TokNo: gen.Tok_NUM, Match: "1", DataType: CtxType_Int, CurValue: 1},
			tok.Token{TokNo: gen.Tok_CL_SQ, Match: "]"},
			//			//
		}},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "m01"},
			tok.Token{TokNo: gen.Tok_EQ, Match: "="},
			tok.Token{TokNo: gen.Tok_ID, Match: "abc"},
			tok.Token{TokNo: gen.Tok_OP_SQ, Match: "["},
			tok.Token{TokNo: gen.Tok_COLON, Match: ":"},
			tok.Token{TokNo: gen.Tok_CL_SQ, Match: "]"},
			//			//
		}},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "m02"},
			tok.Token{TokNo: gen.Tok_EQ, Match: "="},
			tok.Token{TokNo: gen.Tok_ID, Match: "abc"},
			tok.Token{TokNo: gen.Tok_OP_SQ, Match: "["},
			tok.Token{TokNo: gen.Tok_NUM, Match: "1", DataType: CtxType_Int, CurValue: 1},
			tok.Token{TokNo: gen.Tok_COLON, Match: ":"},
			tok.Token{TokNo: gen.Tok_CL_SQ, Match: "]"},
			//
			// xyzzy1857 - incomplete
			//
			// m:n out of range
			// m:n wrong data type
			// m:n errors
			// automatic test on length of array that is the result - and values inside array
			//	 	CmdArrayLen         = 12 //
			// 		CmdArrayValue         = 13 // IValue is subscript, NValue is data
		}},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_ArrayOf, NValue: 3},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdDone},
	}},

	// - Call function - 2 params concat - check that fucntion gets called.
	{Test: "0116", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdSetGlob, NValue: 0},
		ActType{OpCode: CmdInsertInContext, Id: "cat2", DataType: CtxType_Func, FxValue: func(a, b string) string { g_global = 1922; return a + b }},
		ActType{OpCode: CmdInsertInContext, Id: "sss", DataType: CtxType_Str, SValue: "sss"},
		ActType{OpCode: CmdInsertInContext, Id: "ttt", DataType: CtxType_Str, SValue: "ttt"},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "cat2"},
			tok.Token{TokNo: gen.Tok_OP, Match: "("},
			tok.Token{TokNo: gen.Tok_ID, Match: "sss"},
			tok.Token{TokNo: gen.Tok_COMMA, Match: ")"},
			tok.Token{TokNo: gen.Tok_ID, Match: "ttt"},
			tok.Token{TokNo: gen.Tok_CL, Match: ")"},
		}},
		ActType{OpCode: CmdInsertInContext, Id: "cat2a", DataType: CtxType_Func, FxValue: func(a, b string) string { g_global = 1935; return a + b }},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Str, SValue: "sssttt"},
		ActType{OpCode: CmdCmpGlob, NValue: 1922},
		ActType{OpCode: CmdDone},
	}},

	// - Call function - At end to produce error
	{Test: "0117", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdSetGlob, NValue: 0},
		ActType{OpCode: CmdInsertInContext, Id: "cat2", DataType: CtxType_Func, FxValue: func(a, b string) string { g_global = 1945; return a + b }},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "cat2"},
		}},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Int, NValue: 0},
		ActType{OpCode: CmdErrorsContain, Id: "Eval00051", NValue: 1}, // One error containing string
		ActType{OpCode: CmdDone},
	}},

	// Test Unary - error, on right of 29*
	{Test: "0118", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdInsertInContext, Id: "abc", DataType: CtxType_Int, NValue: 5},
		ActType{OpCode: CmdInsertInContext, Id: "def", DataType: CtxType_Str, SValue: "xxxxxx"},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "abc"},
			tok.Token{TokNo: gen.Tok_MINUS, Match: "-"},
			tok.Token{TokNo: gen.Tok_ID, Match: "abc"},
			tok.Token{TokNo: gen.Tok_MINUS, Match: "-"},
			tok.Token{TokNo: gen.Tok_ID, Match: "def"},
		}},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Int, NValue: 0},
		ActType{OpCode: CmdErrorsContain, Id: "Eval00036", NValue: 1}, // One error containing string
		ActType{OpCode: CmdDone},
	}},

	//
	{Test: "0119", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdInsertInContext, Id: "abc", DataType: CtxType_Int, NValue: 1},
		ActType{OpCode: CmdInsertInContext, Id: "def", DataType: CtxType_Str, SValue: "xxxxxx"},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "abc"},
			tok.Token{TokNo: gen.Tok_S_L, Match: "<<"},
			tok.Token{TokNo: gen.Tok_ID, Match: "abc"},
			tok.Token{TokNo: gen.Tok_MINUS, Match: "-"},
			tok.Token{TokNo: gen.Tok_ID, Match: "def"},
		}},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Int, NValue: 0},
		ActType{OpCode: CmdErrorsContain, Id: "Eval00036", NValue: 1}, // One error containing string
		ActType{OpCode: CmdDone},
	}},

	{Test: "0120", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdInsertInContext, Id: "abc", DataType: CtxType_Int, NValue: 1},
		ActType{OpCode: CmdInsertInContext, Id: "def", DataType: CtxType_Str, SValue: "xxxxxx"},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "abc"},
			tok.Token{TokNo: gen.Tok_LE, Match: "<="},
			tok.Token{TokNo: gen.Tok_ID, Match: "abc"},
			tok.Token{TokNo: gen.Tok_MINUS, Match: "-"},
			tok.Token{TokNo: gen.Tok_ID, Match: "def"},
		}},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Int, NValue: 0},
		ActType{OpCode: CmdErrorsContain, Id: "Eval00036", NValue: 1}, // One error containing string
		ActType{OpCode: CmdDone},
	}},

	{Test: "0121", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdInsertInContext, Id: "abc", DataType: CtxType_Int, NValue: 1},
		ActType{OpCode: CmdInsertInContext, Id: "def", DataType: CtxType_Str, SValue: "xxxxxx"},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "abc"},
			tok.Token{TokNo: gen.Tok_L_EQ, Match: "=="},
			tok.Token{TokNo: gen.Tok_ID, Match: "abc"},
			tok.Token{TokNo: gen.Tok_MINUS, Match: "-"},
			tok.Token{TokNo: gen.Tok_ID, Match: "def"},
		}},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Int, NValue: 0},
		ActType{OpCode: CmdErrorsContain, Id: "Eval00036", NValue: 1}, // One error containing string
		ActType{OpCode: CmdDone},
	}},

	{Test: "0122", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdInsertInContext, Id: "abc", DataType: CtxType_Int, NValue: 1},
		ActType{OpCode: CmdInsertInContext, Id: "def", DataType: CtxType_Str, SValue: "xxxxxx"},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "abc"},
			tok.Token{TokNo: gen.Tok_B_AND, Match: "&"},
			tok.Token{TokNo: gen.Tok_ID, Match: "abc"},
			tok.Token{TokNo: gen.Tok_MINUS, Match: "-"},
			tok.Token{TokNo: gen.Tok_ID, Match: "def"},
		}},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Int, NValue: 0},
		ActType{OpCode: CmdErrorsContain, Id: "Eval00036", NValue: 1}, // One error containing string
		ActType{OpCode: CmdDone},
	}},

	{Test: "0123", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdInsertInContext, Id: "abc", DataType: CtxType_Int, NValue: 1},
		ActType{OpCode: CmdInsertInContext, Id: "def", DataType: CtxType_Str, SValue: "xxxxxx"},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "abc"},
			tok.Token{TokNo: gen.Tok_B_OR, Match: "bor"},
			tok.Token{TokNo: gen.Tok_ID, Match: "abc"},
			tok.Token{TokNo: gen.Tok_MINUS, Match: "-"},
			tok.Token{TokNo: gen.Tok_ID, Match: "def"},
		}},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Int, NValue: 0},
		ActType{OpCode: CmdErrorsContain, Id: "Eval0004", NValue: 1}, // One error containing string
		ActType{OpCode: CmdDone},
	}},

	{Test: "0124", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdInsertInContext, Id: "abc", DataType: CtxType_Int, NValue: 1},
		ActType{OpCode: CmdInsertInContext, Id: "def", DataType: CtxType_Str, SValue: "xxxxxx"},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "abc"},
			tok.Token{TokNo: gen.Tok_CARRET, Match: "^"},
			tok.Token{TokNo: gen.Tok_ID, Match: "abc"},
			tok.Token{TokNo: gen.Tok_MINUS, Match: "-"},
			tok.Token{TokNo: gen.Tok_ID, Match: "def"},
		}},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Int, NValue: 0},
		ActType{OpCode: CmdErrorsContain, Id: "Eval0004", NValue: 1}, // One error containing string
		ActType{OpCode: CmdDone},
	}},

	// || -- error
	{Test: "0125", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdInsertInContext, Id: "abc", DataType: CtxType_Float, FValue: 5},
		ActType{OpCode: CmdInsertInContext, Id: "def", DataType: CtxType_Float, FValue: 10},
		ActType{OpCode: CmdInsertInContext, Id: "iii", DataType: CtxType_Int, NValue: 2},
		ActType{OpCode: CmdInsertInContext, Id: "jjj", DataType: CtxType_Int, NValue: 3},
		ActType{OpCode: CmdInsertInContext, Id: "kkk", DataType: CtxType_Int, NValue: 2},
		ActType{OpCode: CmdInsertInContext, Id: "sss", DataType: CtxType_Str, SValue: "sss"},
		ActType{OpCode: CmdInsertInContext, Id: "ttt", DataType: CtxType_Str, SValue: "ttt"},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "iii"},
			tok.Token{TokNo: gen.Tok_NE, Match: "!="},
			tok.Token{TokNo: gen.Tok_ID, Match: "kkk"},
			tok.Token{TokNo: gen.Tok_L_OR, Match: "||"},
			tok.Token{TokNo: gen.Tok_ID, Match: "abc"},
			tok.Token{TokNo: gen.Tok_CARRET, Match: "^"},
			tok.Token{TokNo: gen.Tok_ID, Match: "def"},
		}},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Int, NValue: 0},
		ActType{OpCode: CmdErrorsContain, Id: "Eval00012", NValue: 1}, // One error containing string
		ActType{OpCode: CmdDone},
	}},

	// && -- error
	{Test: "0126", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdInsertInContext, Id: "abc", DataType: CtxType_Float, FValue: 5},
		ActType{OpCode: CmdInsertInContext, Id: "def", DataType: CtxType_Float, FValue: 10},
		ActType{OpCode: CmdInsertInContext, Id: "iii", DataType: CtxType_Int, NValue: 2},
		ActType{OpCode: CmdInsertInContext, Id: "jjj", DataType: CtxType_Int, NValue: 3},
		ActType{OpCode: CmdInsertInContext, Id: "kkk", DataType: CtxType_Int, NValue: 2},
		ActType{OpCode: CmdInsertInContext, Id: "sss", DataType: CtxType_Str, SValue: "sss"},
		ActType{OpCode: CmdInsertInContext, Id: "ttt", DataType: CtxType_Str, SValue: "ttt"},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "iii"},
			tok.Token{TokNo: gen.Tok_L_EQ, Match: "=="},
			tok.Token{TokNo: gen.Tok_ID, Match: "kkk"},
			tok.Token{TokNo: gen.Tok_L_AND, Match: "&&"},
			tok.Token{TokNo: gen.Tok_ID, Match: "abc"},
			tok.Token{TokNo: gen.Tok_CARRET, Match: "^"},
			tok.Token{TokNo: gen.Tok_ID, Match: "def"},
		}},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Int, NValue: 0},
		ActType{OpCode: CmdErrorsContain, Id: "Eval00012", NValue: 1}, // One error containing string
		ActType{OpCode: CmdDone},
	}},

	// Test assinging new value to ID
	{Test: "0127", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdInsertInContext, Id: "abc", DataType: CtxType_Int, NValue: 5},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "abc"},
			tok.Token{TokNo: gen.Tok_EQ, Match: "="},
			tok.Token{TokNo: gen.Tok_true, Match: "true", DataType: CtxType_Bool, CurValue: true},
			tok.Token{TokNo: gen.Tok_ID, Match: "abc"},
			tok.Token{TokNo: gen.Tok_EQ, Match: "="},
			tok.Token{TokNo: gen.Tok_true, Match: "true", DataType: CtxType_Bool, CurValue: true},
		}},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Bool, BValue: true},
		ActType{OpCode: CmdDone},
	}},

	// Test assinging new value to ID
	{Test: "0128", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdInsertInContext, Id: "abc", DataType: CtxType_Bool, BValue: false},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "abc"},
			tok.Token{TokNo: gen.Tok_EQ, Match: "="},
			tok.Token{TokNo: gen.Tok_true, Match: "true", DataType: CtxType_Bool, CurValue: true},
			tok.Token{TokNo: gen.Tok_ID, Match: "abc"},
			tok.Token{TokNo: gen.Tok_EQ, Match: "="},
			tok.Token{TokNo: gen.Tok_true, Match: "true", DataType: CtxType_Bool, CurValue: true},
		}},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Bool, BValue: true},
		ActType{OpCode: CmdDone},
	}},

	// Incorrect Results
	{Test: "0129", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_OP_BRACE, Match: "{"},
			tok.Token{TokNo: gen.Tok_ID, Match: "aaa"},
			tok.Token{TokNo: gen.Tok_COLON, Match: ":"},
			tok.Token{TokNo: gen.Tok_NUM, Match: "1", DataType: CtxType_Int, CurValue: 1},
			tok.Token{TokNo: gen.Tok_COMMA, Match: ","},

			tok.Token{TokNo: gen.Tok_Str0, Match: "aab", DataType: CtxType_Str, CurValue: "aab"},
			tok.Token{TokNo: gen.Tok_COLON, Match: ":"},
			tok.Token{TokNo: gen.Tok_OP_SQ, Match: "["},
			tok.Token{TokNo: gen.Tok_NUM, Match: "11", DataType: CtxType_Int, CurValue: 11},
			tok.Token{TokNo: gen.Tok_COMMA, Match: ","},
			tok.Token{TokNo: gen.Tok_NUM, Match: "12", DataType: CtxType_Int, CurValue: 12},
			tok.Token{TokNo: gen.Tok_COMMA, Match: ","},
			tok.Token{TokNo: gen.Tok_NUM, Match: "13", DataType: CtxType_Int, CurValue: 13},
			tok.Token{TokNo: gen.Tok_CL_SQ, Match: "]"},
			tok.Token{TokNo: gen.Tok_COMMA, Match: ","},

			tok.Token{TokNo: gen.Tok_ID, Match: "aac", DataType: CtxType_Str, CurValue: "aac"},
			tok.Token{TokNo: gen.Tok_COLON, Match: ":"},
			tok.Token{TokNo: gen.Tok_OP_BRACE, Match: "{"},
			tok.Token{TokNo: gen.Tok_ID, Match: "AAA"},
			tok.Token{TokNo: gen.Tok_COLON, Match: ":"},
			tok.Token{TokNo: gen.Tok_NUM, Match: "23", DataType: CtxType_Int, CurValue: 23},
			tok.Token{TokNo: gen.Tok_CL_BRACE, Match: "}"},
			tok.Token{TokNo: gen.Tok_CL_BRACE, Match: "}"},
		}},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Bool, BValue: true},
		ActType{OpCode: CmdDone},
	}},

	// xyzzy1857 - incomplete
	{Test: "0130", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdInsertInContext, Id: "abc", DataType: CtxType_ArrayOf},
		ActType{OpCode: CmdInsertInContext, Id: "def", DataType: CtxType_Int, NValue: 0},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "abc"},
			tok.Token{TokNo: gen.Tok_EQ, Match: "="},

			tok.Token{TokNo: gen.Tok_OP_BRACE, Match: "{"},

			tok.Token{TokNo: gen.Tok_Str0, Match: "aab", DataType: CtxType_Str, CurValue: "aab"},
			tok.Token{TokNo: gen.Tok_COLON, Match: ":"},
			tok.Token{TokNo: gen.Tok_NUM, Match: "1", DataType: CtxType_Int, CurValue: 1},
			tok.Token{TokNo: gen.Tok_COMMA, Match: ","},

			tok.Token{TokNo: gen.Tok_Str0, Match: "xxx", DataType: CtxType_Str, CurValue: "xxx"},
			tok.Token{TokNo: gen.Tok_COLON, Match: ":"},
			tok.Token{TokNo: gen.Tok_NUM, Match: "222", DataType: CtxType_Int, CurValue: 222},
			tok.Token{TokNo: gen.Tok_COMMA, Match: ","},

			tok.Token{TokNo: gen.Tok_Str0, Match: "aac", DataType: CtxType_Str, CurValue: "aac"},
			tok.Token{TokNo: gen.Tok_COLON, Match: ":"},
			tok.Token{TokNo: gen.Tok_NUM, Match: "3", DataType: CtxType_Int, CurValue: 3},

			tok.Token{TokNo: gen.Tok_CL_BRACE, Match: "}"},
		}},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "def"},
			tok.Token{TokNo: gen.Tok_EQ, Match: "="},
			tok.Token{TokNo: gen.Tok_ID, Match: "abc"},
			tok.Token{TokNo: gen.Tok_DOT, Match: "."},
			tok.Token{TokNo: gen.Tok_Str0, Match: "xxx"},
			//			//
		}},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Int, NValue: 222},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdDone},
	}},

	{Test: "0131", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdInsertInContext, Id: "abc", DataType: CtxType_ArrayOf},
		ActType{OpCode: CmdInsertInContext, Id: "def", DataType: CtxType_Int, NValue: 0},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "abc"},
			tok.Token{TokNo: gen.Tok_EQ, Match: "="},

			tok.Token{TokNo: gen.Tok_OP_BRACE, Match: "{"},

			tok.Token{TokNo: gen.Tok_Str0, Match: "aab", DataType: CtxType_Str, CurValue: "aab"},
			tok.Token{TokNo: gen.Tok_COLON, Match: ":"},
			tok.Token{TokNo: gen.Tok_NUM, Match: "1", DataType: CtxType_Int, CurValue: 1},
			tok.Token{TokNo: gen.Tok_COMMA, Match: ","},

			tok.Token{TokNo: gen.Tok_Str0, Match: "xxx", DataType: CtxType_Str, CurValue: "xxx"},
			tok.Token{TokNo: gen.Tok_COLON, Match: ":"},
			tok.Token{TokNo: gen.Tok_NUM, Match: "222", DataType: CtxType_Int, CurValue: 222},
			tok.Token{TokNo: gen.Tok_COMMA, Match: ","},

			tok.Token{TokNo: gen.Tok_Str0, Match: "aac", DataType: CtxType_Str, CurValue: "aac"},
			tok.Token{TokNo: gen.Tok_COLON, Match: ":"},
			tok.Token{TokNo: gen.Tok_NUM, Match: "3", DataType: CtxType_Int, CurValue: 3},

			tok.Token{TokNo: gen.Tok_CL_BRACE, Match: "}"},
		}},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "def"},
			tok.Token{TokNo: gen.Tok_EQ, Match: "="},
			tok.Token{TokNo: gen.Tok_ID, Match: "len"},
			tok.Token{TokNo: gen.Tok_OP, Match: "("},
			tok.Token{TokNo: gen.Tok_ID, Match: "abc"},
			tok.Token{TokNo: gen.Tok_CL, Match: ")"},
			//			//
		}},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Int, NValue: 3},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdDone},
	}},

	{Test: "0132", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdInsertInContext, Id: "abc", DataType: CtxType_ArrayOf},
		ActType{OpCode: CmdInsertInContext, Id: "def", DataType: CtxType_Int, NValue: 0},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "abc"},
			tok.Token{TokNo: gen.Tok_EQ, Match: "="},

			tok.Token{TokNo: gen.Tok_OP_BRACE, Match: "{"},

			tok.Token{TokNo: gen.Tok_Str0, Match: "aab", DataType: CtxType_Str, CurValue: "aab"},
			tok.Token{TokNo: gen.Tok_COLON, Match: ":"},
			tok.Token{TokNo: gen.Tok_NUM, Match: "1", DataType: CtxType_Int, CurValue: 1},
			tok.Token{TokNo: gen.Tok_COMMA, Match: ","},

			tok.Token{TokNo: gen.Tok_Str0, Match: "xxx", DataType: CtxType_Str, CurValue: "xxx"},
			tok.Token{TokNo: gen.Tok_COLON, Match: ":"},
			tok.Token{TokNo: gen.Tok_NUM, Match: "222", DataType: CtxType_Int, CurValue: 222},
			tok.Token{TokNo: gen.Tok_COMMA, Match: ","},

			tok.Token{TokNo: gen.Tok_Str0, Match: "aac", DataType: CtxType_Str, CurValue: "aac"},
			tok.Token{TokNo: gen.Tok_COLON, Match: ":"},
			tok.Token{TokNo: gen.Tok_NUM, Match: "3", DataType: CtxType_Int, CurValue: 3},

			tok.Token{TokNo: gen.Tok_CL_BRACE, Match: "}"},
		}},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "def"},
			tok.Token{TokNo: gen.Tok_EQ, Match: "="},
			tok.Token{TokNo: gen.Tok_ID, Match: "len"},
			tok.Token{TokNo: gen.Tok_OP, Match: "("},
			tok.Token{TokNo: gen.Tok_ID, Match: "abc"},
			tok.Token{TokNo: gen.Tok_CL, Match: ")"},
			//			//
		}},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Int, NValue: 3},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdDone},
	}},

	{Test: "0133", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdInsertInContext, Id: "abc", DataType: CtxType_ArrayOf},
		ActType{OpCode: CmdInsertInContext, Id: "arr", DataType: CtxType_ArrayOf},
		ActType{OpCode: CmdInsertInContext, Id: "def", DataType: CtxType_Int, NValue: 0},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "abc"},
			tok.Token{TokNo: gen.Tok_EQ, Match: "="},

			tok.Token{TokNo: gen.Tok_OP_BRACE, Match: "{"},

			tok.Token{TokNo: gen.Tok_Str0, Match: "aab", DataType: CtxType_Str, CurValue: "aab"},
			tok.Token{TokNo: gen.Tok_COLON, Match: ":"},
			tok.Token{TokNo: gen.Tok_NUM, Match: "1", DataType: CtxType_Int, CurValue: 1},
			tok.Token{TokNo: gen.Tok_COMMA, Match: ","},

			tok.Token{TokNo: gen.Tok_Str0, Match: "xxx", DataType: CtxType_Str, CurValue: "xxx"},
			tok.Token{TokNo: gen.Tok_COLON, Match: ":"},
			tok.Token{TokNo: gen.Tok_NUM, Match: "222", DataType: CtxType_Int, CurValue: 222},
			tok.Token{TokNo: gen.Tok_COMMA, Match: ","},

			tok.Token{TokNo: gen.Tok_Str0, Match: "aac", DataType: CtxType_Str, CurValue: "aac"},
			tok.Token{TokNo: gen.Tok_COLON, Match: ":"},
			tok.Token{TokNo: gen.Tok_NUM, Match: "3", DataType: CtxType_Int, CurValue: 3},

			tok.Token{TokNo: gen.Tok_CL_BRACE, Match: "}"},
		}},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "arr"},
			tok.Token{TokNo: gen.Tok_EQ, Match: "="},
			tok.Token{TokNo: gen.Tok_OP_SQ, Match: "["},
			tok.Token{TokNo: gen.Tok_NUM, Match: "1", DataType: CtxType_Int, CurValue: 1},
			tok.Token{TokNo: gen.Tok_COMMA, Match: ","},
			tok.Token{TokNo: gen.Tok_NUM, Match: "2", DataType: CtxType_Int, CurValue: 2},
			tok.Token{TokNo: gen.Tok_COMMA, Match: ","},
			tok.Token{TokNo: gen.Tok_NUM, Match: "3", DataType: CtxType_Int, CurValue: 3},
			tok.Token{TokNo: gen.Tok_CL_SQ, Match: "]"},
		}},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "def"},
			tok.Token{TokNo: gen.Tok_EQ, Match: "="},
			tok.Token{TokNo: gen.Tok_ID, Match: "len"},
			tok.Token{TokNo: gen.Tok_OP, Match: "("},
			tok.Token{TokNo: gen.Tok_NUM, Match: "3", DataType: CtxType_Int, CurValue: 3},
			tok.Token{TokNo: gen.Tok_CL, Match: ")"},
			//			//
		}},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Int, NValue: 0},
		ActType{OpCode: CmdErrorsContain, Id: "Eval01003", NValue: 1}, // One error containing string
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "def"},
			tok.Token{TokNo: gen.Tok_EQ, Match: "="},
			tok.Token{TokNo: gen.Tok_ID, Match: "len_e"},
			tok.Token{TokNo: gen.Tok_OP, Match: "("},
			tok.Token{TokNo: gen.Tok_NUM, Match: "3", DataType: CtxType_Int, CurValue: 3},
			tok.Token{TokNo: gen.Tok_CL, Match: ")"},
			//			//
		}},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Int, NValue: 0},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "def"},
			tok.Token{TokNo: gen.Tok_EQ, Match: "="},
			tok.Token{TokNo: gen.Tok_ID, Match: "len_e"},
			tok.Token{TokNo: gen.Tok_OP, Match: "("},
			tok.Token{TokNo: gen.Tok_ID, Match: "abc"},
			tok.Token{TokNo: gen.Tok_CL, Match: ")"},
			//			//
		}},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Int, NValue: 3},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "def"},
			tok.Token{TokNo: gen.Tok_EQ, Match: "="},
			tok.Token{TokNo: gen.Tok_ID, Match: "len"},
			tok.Token{TokNo: gen.Tok_OP, Match: "("},
			tok.Token{TokNo: gen.Tok_ID, Match: "abc"},
			tok.Token{TokNo: gen.Tok_CL, Match: ")"},
			//			//
		}},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Int, NValue: 3},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "def"},
			tok.Token{TokNo: gen.Tok_EQ, Match: "="},
			tok.Token{TokNo: gen.Tok_ID, Match: "len"},
			tok.Token{TokNo: gen.Tok_OP, Match: "("},
			tok.Token{TokNo: gen.Tok_ID, Match: "arr"},
			tok.Token{TokNo: gen.Tok_CL, Match: ")"},
			//			//
		}},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Int, NValue: 3},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "def"},
			tok.Token{TokNo: gen.Tok_EQ, Match: "="},
			tok.Token{TokNo: gen.Tok_ID, Match: "len_e"},
			tok.Token{TokNo: gen.Tok_OP, Match: "("},
			tok.Token{TokNo: gen.Tok_ID, Match: "arr"},
			tok.Token{TokNo: gen.Tok_CL, Match: ")"},
			//			//
		}},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Int, NValue: 3},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdDone},
	}},

	{Test: "0134", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdInsertInContext, Id: "abc", DataType: CtxType_ArrayOf},
		ActType{OpCode: CmdInsertInContext, Id: "def", DataType: CtxType_Int, NValue: 0},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "def"},
			tok.Token{TokNo: gen.Tok_EQ, Match: "="},
			tok.Token{TokNo: gen.Tok_ID, Match: "x_test_ret_float"},
			tok.Token{TokNo: gen.Tok_OP, Match: "("},
			tok.Token{TokNo: gen.Tok_NUM, Match: "3", DataType: CtxType_Int, CurValue: 3},
			tok.Token{TokNo: gen.Tok_CL, Match: ")"},
		}},
		//ActType{OpCode: CmdCmpTo, DataType: CtxType_Int, NValue: 0},
		//ActType{OpCode: CmdErrorsContain, Id: "Eval01003", NValue: 1}, // One error containing string
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "def"},
			tok.Token{TokNo: gen.Tok_EQ, Match: "="},
			tok.Token{TokNo: gen.Tok_ID, Match: "x_test_ret_bool"},
			tok.Token{TokNo: gen.Tok_OP, Match: "("},
			tok.Token{TokNo: gen.Tok_NUM, Match: "3", DataType: CtxType_Int, CurValue: 3},
			tok.Token{TokNo: gen.Tok_CL, Match: ")"},
		}},
		//ActType{OpCode: CmdCmpTo, DataType: CtxType_Int, NValue: 0},
		//ActType{OpCode: CmdErrorsContain, Id: "Eval01003", NValue: 1}, // One error containing string
		ActType{OpCode: CmdDone},
	}},

	{Test: "0135", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdInsertInContext, Id: "abc", DataType: CtxType_ArrayOf},
		ActType{OpCode: CmdInsertInContext, Id: "def", DataType: CtxType_Int, NValue: 0},
		ActType{OpCode: CmdInsertInContext, Id: "bbb", DataType: CtxType_Bool, BValue: true},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "def"},
			tok.Token{TokNo: gen.Tok_EQ, Match: "="},
			tok.Token{TokNo: gen.Tok_ID, Match: "float"},
			tok.Token{TokNo: gen.Tok_OP, Match: "("},
			tok.Token{TokNo: gen.Tok_NUM, Match: "3", DataType: CtxType_Int, CurValue: 3},
			tok.Token{TokNo: gen.Tok_CL, Match: ")"},
		}},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Float, FValue: 3},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "def"},
			tok.Token{TokNo: gen.Tok_EQ, Match: "="},
			tok.Token{TokNo: gen.Tok_ID, Match: "float"},
			tok.Token{TokNo: gen.Tok_OP, Match: "("},
			tok.Token{TokNo: gen.Tok_Float, Match: "3.0000", DataType: CtxType_Float, CurValue: 3.0},
			tok.Token{TokNo: gen.Tok_CL, Match: ")"},
		}},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Float, FValue: 3},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "def"},
			tok.Token{TokNo: gen.Tok_EQ, Match: "="},
			tok.Token{TokNo: gen.Tok_ID, Match: "float"},
			tok.Token{TokNo: gen.Tok_OP, Match: "("},
			tok.Token{TokNo: gen.Tok_ID, Match: "bbb"},
			tok.Token{TokNo: gen.Tok_CL, Match: ")"},
		}},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Float, FValue: 1},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "def"},
			tok.Token{TokNo: gen.Tok_EQ, Match: "="},
			tok.Token{TokNo: gen.Tok_ID, Match: "float"},
			tok.Token{TokNo: gen.Tok_OP, Match: "("},
			tok.Token{TokNo: gen.Tok_EXCLAM, Match: "!"},
			tok.Token{TokNo: gen.Tok_ID, Match: "bbb"},
			tok.Token{TokNo: gen.Tok_CL, Match: ")"},
		}},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Float, FValue: 0},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "def"},
			tok.Token{TokNo: gen.Tok_EQ, Match: "="},
			tok.Token{TokNo: gen.Tok_ID, Match: "float"},
			tok.Token{TokNo: gen.Tok_OP, Match: "("},
			tok.Token{TokNo: gen.Tok_Str0, Match: "", DataType: CtxType_Str, CurValue: ""},
			tok.Token{TokNo: gen.Tok_CL, Match: ")"},
		}},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Float, FValue: 1},
		ActType{OpCode: CmdDone},
	}},

	{Test: "0136", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdInsertInContext, Id: "abc", DataType: CtxType_ArrayOf},
		ActType{OpCode: CmdInsertInContext, Id: "hhh", DataType: CtxType_MapOf},
		ActType{OpCode: CmdInsertInContext, Id: "def", DataType: CtxType_Int, NValue: 0},
		ActType{OpCode: CmdInsertInContext, Id: "bbb", DataType: CtxType_Bool, BValue: true},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "def"},
			tok.Token{TokNo: gen.Tok_EQ, Match: "="},
			tok.Token{TokNo: gen.Tok_ID, Match: "bool"},
			tok.Token{TokNo: gen.Tok_OP, Match: "("},
			tok.Token{TokNo: gen.Tok_NUM, Match: "0", DataType: CtxType_Int, CurValue: 0},
			tok.Token{TokNo: gen.Tok_CL, Match: ")"},
		}},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Bool, BValue: false},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "def"},
			tok.Token{TokNo: gen.Tok_EQ, Match: "="},
			tok.Token{TokNo: gen.Tok_ID, Match: "bool"},
			tok.Token{TokNo: gen.Tok_OP, Match: "("},
			tok.Token{TokNo: gen.Tok_Float, Match: "3.0000", DataType: CtxType_Float, CurValue: 3.0},
			tok.Token{TokNo: gen.Tok_CL, Match: ")"},
		}},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Bool, BValue: true},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "def"},
			tok.Token{TokNo: gen.Tok_EQ, Match: "="},
			tok.Token{TokNo: gen.Tok_ID, Match: "bool"},
			tok.Token{TokNo: gen.Tok_OP, Match: "("},
			tok.Token{TokNo: gen.Tok_ID, Match: "bbb"},
			tok.Token{TokNo: gen.Tok_CL, Match: ")"},
		}},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Bool, BValue: true},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "def"},
			tok.Token{TokNo: gen.Tok_EQ, Match: "="},
			tok.Token{TokNo: gen.Tok_ID, Match: "bool"},
			tok.Token{TokNo: gen.Tok_OP, Match: "("},
			tok.Token{TokNo: gen.Tok_EXCLAM, Match: "!"},
			tok.Token{TokNo: gen.Tok_ID, Match: "bbb"},
			tok.Token{TokNo: gen.Tok_CL, Match: ")"},
		}},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Bool, BValue: false},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "def"},
			tok.Token{TokNo: gen.Tok_EQ, Match: "="},
			tok.Token{TokNo: gen.Tok_ID, Match: "bool"},
			tok.Token{TokNo: gen.Tok_OP, Match: "("},
			tok.Token{TokNo: gen.Tok_ID, Match: "abc"},
			tok.Token{TokNo: gen.Tok_CL, Match: ")"},
		}},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Bool, BValue: false},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "abc"},
			tok.Token{TokNo: gen.Tok_EQ, Match: "="},
			tok.Token{TokNo: gen.Tok_OP_SQ, Match: "["},
			tok.Token{TokNo: gen.Tok_NUM, Match: "1", DataType: CtxType_Int, CurValue: 1},
			tok.Token{TokNo: gen.Tok_COMMA, Match: ","},
			tok.Token{TokNo: gen.Tok_NUM, Match: "2", DataType: CtxType_Int, CurValue: 2},
			tok.Token{TokNo: gen.Tok_COMMA, Match: ","},
			tok.Token{TokNo: gen.Tok_NUM, Match: "3", DataType: CtxType_Int, CurValue: 3},
			tok.Token{TokNo: gen.Tok_CL_SQ, Match: "]"},
		}},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "def"},
			tok.Token{TokNo: gen.Tok_EQ, Match: "="},
			tok.Token{TokNo: gen.Tok_ID, Match: "bool"},
			tok.Token{TokNo: gen.Tok_OP, Match: "("},
			tok.Token{TokNo: gen.Tok_ID, Match: "abc"},
			tok.Token{TokNo: gen.Tok_CL, Match: ")"},
		}},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Bool, BValue: true},

		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "def"},
			tok.Token{TokNo: gen.Tok_EQ, Match: "="},
			tok.Token{TokNo: gen.Tok_ID, Match: "bool"},
			tok.Token{TokNo: gen.Tok_OP, Match: "("},
			tok.Token{TokNo: gen.Tok_ID, Match: "hhh"},
			tok.Token{TokNo: gen.Tok_CL, Match: ")"},
		}},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Bool, BValue: false},

		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "hhh"},
			tok.Token{TokNo: gen.Tok_EQ, Match: "="},

			tok.Token{TokNo: gen.Tok_OP_BRACE, Match: "{"},

			tok.Token{TokNo: gen.Tok_Str0, Match: "aab", DataType: CtxType_Str, CurValue: "aab"},
			tok.Token{TokNo: gen.Tok_COLON, Match: ":"},
			tok.Token{TokNo: gen.Tok_NUM, Match: "1", DataType: CtxType_Int, CurValue: 1},
			tok.Token{TokNo: gen.Tok_COMMA, Match: ","},

			tok.Token{TokNo: gen.Tok_Str0, Match: "xxx", DataType: CtxType_Str, CurValue: "xxx"},
			tok.Token{TokNo: gen.Tok_COLON, Match: ":"},
			tok.Token{TokNo: gen.Tok_NUM, Match: "222", DataType: CtxType_Int, CurValue: 222},
			tok.Token{TokNo: gen.Tok_COMMA, Match: ","},

			tok.Token{TokNo: gen.Tok_Str0, Match: "aac", DataType: CtxType_Str, CurValue: "aac"},
			tok.Token{TokNo: gen.Tok_COLON, Match: ":"},
			tok.Token{TokNo: gen.Tok_NUM, Match: "3", DataType: CtxType_Int, CurValue: 3},

			tok.Token{TokNo: gen.Tok_CL_BRACE, Match: "}"},
		}},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "def"},
			tok.Token{TokNo: gen.Tok_EQ, Match: "="},
			tok.Token{TokNo: gen.Tok_ID, Match: "bool"},
			tok.Token{TokNo: gen.Tok_OP, Match: "("},
			tok.Token{TokNo: gen.Tok_ID, Match: "hhh"},
			tok.Token{TokNo: gen.Tok_CL, Match: ")"},
		}},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Bool, BValue: true},

		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "def"},
			tok.Token{TokNo: gen.Tok_EQ, Match: "="},
			tok.Token{TokNo: gen.Tok_ID, Match: "bool"},
			tok.Token{TokNo: gen.Tok_OP, Match: "("},
			tok.Token{TokNo: gen.Tok_Str0, Match: "aac", DataType: CtxType_Str, CurValue: "aac"},
			tok.Token{TokNo: gen.Tok_CL, Match: ")"},
		}},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Bool, BValue: true},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "def"},
			tok.Token{TokNo: gen.Tok_EQ, Match: "="},
			tok.Token{TokNo: gen.Tok_ID, Match: "bool"},
			tok.Token{TokNo: gen.Tok_OP, Match: "("},
			tok.Token{TokNo: gen.Tok_Str0, Match: "", DataType: CtxType_Str, CurValue: ""},
			tok.Token{TokNo: gen.Tok_CL, Match: ")"},
		}},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Bool, BValue: false},
		ActType{OpCode: CmdDone},
	}},

	{Test: "0137", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdInsertInContext, Id: "abc", DataType: CtxType_ArrayOf},
		ActType{OpCode: CmdInsertInContext, Id: "def", DataType: CtxType_Int, NValue: 0},
		ActType{OpCode: CmdInsertInContext, Id: "bbb", DataType: CtxType_Bool, BValue: true},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "def"},
			tok.Token{TokNo: gen.Tok_EQ, Match: "="},
			tok.Token{TokNo: gen.Tok_ID, Match: "int"},
			tok.Token{TokNo: gen.Tok_OP, Match: "("},
			tok.Token{TokNo: gen.Tok_NUM, Match: "0", DataType: CtxType_Int, CurValue: 1},
			tok.Token{TokNo: gen.Tok_CL, Match: ")"},
		}},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Int, NValue: 1},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "def"},
			tok.Token{TokNo: gen.Tok_EQ, Match: "="},
			tok.Token{TokNo: gen.Tok_ID, Match: "int"},
			tok.Token{TokNo: gen.Tok_OP, Match: "("},
			tok.Token{TokNo: gen.Tok_Float, Match: "3.0000", DataType: CtxType_Float, CurValue: 3.0},
			tok.Token{TokNo: gen.Tok_CL, Match: ")"},
		}},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Int, NValue: 3},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "def"},
			tok.Token{TokNo: gen.Tok_EQ, Match: "="},
			tok.Token{TokNo: gen.Tok_ID, Match: "int"},
			tok.Token{TokNo: gen.Tok_OP, Match: "("},
			tok.Token{TokNo: gen.Tok_ID, Match: "bbb"},
			tok.Token{TokNo: gen.Tok_CL, Match: ")"},
		}},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Int, NValue: 1},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "def"},
			tok.Token{TokNo: gen.Tok_EQ, Match: "="},
			tok.Token{TokNo: gen.Tok_ID, Match: "int"},
			tok.Token{TokNo: gen.Tok_OP, Match: "("},
			tok.Token{TokNo: gen.Tok_EXCLAM, Match: "!"},
			tok.Token{TokNo: gen.Tok_ID, Match: "bbb"},
			tok.Token{TokNo: gen.Tok_CL, Match: ")"},
		}},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Int, NValue: 0},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "def"},
			tok.Token{TokNo: gen.Tok_EQ, Match: "="},
			tok.Token{TokNo: gen.Tok_ID, Match: "int"},
			tok.Token{TokNo: gen.Tok_OP, Match: "("},
			tok.Token{TokNo: gen.Tok_ID, Match: "abc"},
			tok.Token{TokNo: gen.Tok_CL, Match: ")"},
		}},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Int, NValue: 0},
		ActType{OpCode: CmdErrorsContain, Id: "Eval01006", NValue: 1}, // One error containing string
		ActType{OpCode: CmdDone},
	}},

	// Hash True if not empty - type cast
	{Test: "0138", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdInsertInContext, Id: "hhh", DataType: CtxType_MapOf},
		ActType{OpCode: CmdInsertInContext, Id: "def", DataType: CtxType_Bool, BValue: false},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "hhh"},
			tok.Token{TokNo: gen.Tok_EQ, Match: "="},

			tok.Token{TokNo: gen.Tok_OP_BRACE, Match: "{"},

			tok.Token{TokNo: gen.Tok_Str0, Match: "a", DataType: CtxType_Str, CurValue: "a"},
			tok.Token{TokNo: gen.Tok_COLON, Match: ":"},
			tok.Token{TokNo: gen.Tok_NUM, Match: "1", DataType: CtxType_Int, CurValue: 1},
			tok.Token{TokNo: gen.Tok_COMMA, Match: ","},

			tok.Token{TokNo: gen.Tok_CL_BRACE, Match: "}"},
		}},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "def"},
			tok.Token{TokNo: gen.Tok_EQ, Match: "="},
			tok.Token{TokNo: gen.Tok_ID, Match: "bool"},
			tok.Token{TokNo: gen.Tok_OP, Match: "("},
			tok.Token{TokNo: gen.Tok_ID, Match: "hhh"},
			tok.Token{TokNo: gen.Tok_CL, Match: ")"},
		}},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Bool, BValue: true},
		ActType{OpCode: CmdDone},
	}},

	// Conditional IF, ?:
	{Test: "0139", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdInsertInContext, Id: "xyz", DataType: CtxType_Bool, BValue: false},
		ActType{OpCode: CmdInsertInContext, Id: "def", DataType: CtxType_Bool, BValue: true},
		ActType{OpCode: CmdInsertInContext, Id: "rrr", DataType: CtxType_Int, NValue: 0},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "rrr"},
			tok.Token{TokNo: gen.Tok_EQ, Match: "="},
			tok.Token{TokNo: gen.Tok_OP, Match: "("},
			tok.Token{TokNo: gen.Tok_ID, Match: "xyz"},
			tok.Token{TokNo: gen.Tok_QUEST, Match: "?"},
			tok.Token{TokNo: gen.Tok_NUM, Match: "12", DataType: CtxType_Int, CurValue: 12},
			tok.Token{TokNo: gen.Tok_COLON, Match: ":"},
			tok.Token{TokNo: gen.Tok_NUM, Match: "44", DataType: CtxType_Int, CurValue: 44},
			tok.Token{TokNo: gen.Tok_CL, Match: ")"},
		}},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Int, NValue: 44},
		ActType{OpCode: CmdDone},
	}},

	// Conditional IF, ?:
	{Test: "0140", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdInsertInContext, Id: "xyz", DataType: CtxType_Bool, BValue: false},
		ActType{OpCode: CmdInsertInContext, Id: "def", DataType: CtxType_Bool, BValue: true},
		ActType{OpCode: CmdInsertInContext, Id: "rrr", DataType: CtxType_Int, NValue: 0},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "rrr"},
			tok.Token{TokNo: gen.Tok_EQ, Match: "="},
			tok.Token{TokNo: gen.Tok_OP, Match: "("},
			tok.Token{TokNo: gen.Tok_ID, Match: "def"},
			tok.Token{TokNo: gen.Tok_QUEST, Match: "?"},
			tok.Token{TokNo: gen.Tok_NUM, Match: "12", DataType: CtxType_Int, CurValue: 12},
			tok.Token{TokNo: gen.Tok_COLON, Match: ":"},
			tok.Token{TokNo: gen.Tok_NUM, Match: "44", DataType: CtxType_Int, CurValue: 44},
			tok.Token{TokNo: gen.Tok_CL, Match: ")"},
		}},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Int, NValue: 12},
		ActType{OpCode: CmdDone},
	}},

	// Conditional IF, ?:
	{Test: "0141", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdInsertInContext, Id: "xyz", DataType: CtxType_Bool, BValue: false},
		ActType{OpCode: CmdInsertInContext, Id: "def", DataType: CtxType_Bool, BValue: true},
		ActType{OpCode: CmdInsertInContext, Id: "rrr", DataType: CtxType_Int, NValue: 0},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "rrr"},
			tok.Token{TokNo: gen.Tok_EQ, Match: "="},
			tok.Token{TokNo: gen.Tok_OP, Match: "("},
			tok.Token{TokNo: gen.Tok_ID, Match: "def"},
			tok.Token{TokNo: gen.Tok_QUEST, Match: "?"},
			tok.Token{TokNo: gen.Tok_NUM, Match: "12", DataType: CtxType_Int, CurValue: 12},
			tok.Token{TokNo: gen.Tok_CL, Match: ")"},
		}},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Int, NValue: 0},
		ActType{OpCode: CmdErrorsContain, Id: "Eval00095", NValue: 1}, // One error containing string
		ActType{OpCode: CmdDone},
	}},

	// Conditional IF, ?:
	{Test: "0142", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdInsertInContext, Id: "xyz", DataType: CtxType_Bool, BValue: false},
		ActType{OpCode: CmdInsertInContext, Id: "def", DataType: CtxType_Bool, BValue: true},
		ActType{OpCode: CmdInsertInContext, Id: "rrr", DataType: CtxType_Int, NValue: 0},
		ActType{OpCode: CmdInsertInContext, Id: "iii", DataType: CtxType_Int, NValue: 2},
		ActType{OpCode: CmdInsertInContext, Id: "jjj", DataType: CtxType_Int, NValue: 3},
		ActType{OpCode: CmdInsertInContext, Id: "kkk", DataType: CtxType_Int, NValue: 2},
		ActType{OpCode: CmdInsertInContext, Id: "sss", DataType: CtxType_Str, SValue: "sss"},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "rrr"},
			tok.Token{TokNo: gen.Tok_EQ, Match: "="},
			tok.Token{TokNo: gen.Tok_OP, Match: "("},
			tok.Token{TokNo: gen.Tok_ID, Match: "xyz"},
			tok.Token{TokNo: gen.Tok_QUEST, Match: "?"},
			tok.Token{TokNo: gen.Tok_NUM, Match: "12", DataType: CtxType_Int, CurValue: 12},
			tok.Token{TokNo: gen.Tok_COLON, Match: ":"},
			tok.Token{TokNo: gen.Tok_ID, Match: "iii"},
			tok.Token{TokNo: gen.Tok_NE, Match: "!="},
			tok.Token{TokNo: gen.Tok_ID, Match: "sss"},
			tok.Token{TokNo: gen.Tok_CL, Match: ")"},
		}},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Int, NValue: 0},
		ActType{OpCode: CmdErrorsContain, Id: "Eval0007", NValue: 1}, // One error containing string
		ActType{OpCode: CmdDone},
	}},

	//eval_test.go:2954: 0142: 10 error, expected error - did not get one, found it, File: /Users/corwin/Projects/pongo2/lexie/eval/eval_test.go LineNo:2954
	//eval_test.go:2934: 0143: 9 error, expected 12 got 44, found it, File: /Users/corwin/Projects/pongo2/lexie/eval/eval_test.go LineNo:2934
	//eval_test.go:2954: 0143: 10 error, expected error - did not get one, found it, File: /Users/corwin/Projects/pongo2/lexie/eval/eval_test.go LineNo:2954

	// Conditional IF, ?:
	{Test: "0143", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdInsertInContext, Id: "xyz", DataType: CtxType_Bool, BValue: false},
		ActType{OpCode: CmdInsertInContext, Id: "def", DataType: CtxType_Bool, BValue: true},
		ActType{OpCode: CmdInsertInContext, Id: "rrr", DataType: CtxType_Int, NValue: 0},
		ActType{OpCode: CmdInsertInContext, Id: "iii", DataType: CtxType_Int, NValue: 2},
		ActType{OpCode: CmdInsertInContext, Id: "jjj", DataType: CtxType_Int, NValue: 3},
		ActType{OpCode: CmdInsertInContext, Id: "kkk", DataType: CtxType_Int, NValue: 2},
		ActType{OpCode: CmdInsertInContext, Id: "sss", DataType: CtxType_Str, SValue: "sss"},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "rrr"},
			tok.Token{TokNo: gen.Tok_EQ, Match: "="},
			tok.Token{TokNo: gen.Tok_OP, Match: "("},
			tok.Token{TokNo: gen.Tok_ID, Match: "def"},
			tok.Token{TokNo: gen.Tok_QUEST, Match: "?"},
			tok.Token{TokNo: gen.Tok_ID, Match: "iii"},
			tok.Token{TokNo: gen.Tok_NE, Match: "!="},
			tok.Token{TokNo: gen.Tok_ID, Match: "sss"},
			tok.Token{TokNo: gen.Tok_COLON, Match: ":"},
			tok.Token{TokNo: gen.Tok_NUM, Match: "44", DataType: CtxType_Int, CurValue: 44},
			tok.Token{TokNo: gen.Tok_CL, Match: ")"},
		}},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Int, NValue: 0},
		ActType{OpCode: CmdErrorsContain, Id: "Eval0007", NValue: 1}, // One error containing string
		ActType{OpCode: CmdDone},
	}},

	// Conditional IF, ?:
	{Test: "0144", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdInsertInContext, Id: "xyz", DataType: CtxType_Bool, BValue: false},
		ActType{OpCode: CmdInsertInContext, Id: "def", DataType: CtxType_Bool, BValue: true},
		ActType{OpCode: CmdInsertInContext, Id: "rrr", DataType: CtxType_Int, NValue: 0},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "rrr"},
			tok.Token{TokNo: gen.Tok_EQ, Match: "="},
			tok.Token{TokNo: gen.Tok_OP, Match: "("},
			tok.Token{TokNo: gen.Tok_ID, Match: "xyz"},
			tok.Token{TokNo: gen.Tok_QUEST, Match: "?"},
			tok.Token{TokNo: gen.Tok_NUM, Match: "12", DataType: CtxType_Int, CurValue: 12},
			tok.Token{TokNo: gen.Tok_CL, Match: ")"},
		}},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Int, NValue: 0},
		ActType{OpCode: CmdErrorsContain, Id: "Eval00095", NValue: 1}, // One error containing string
		ActType{OpCode: CmdDone},
	}},

	// Pipe Fitting 300
	{Test: "0145", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdInsertInContext, Id: "xyz", DataType: CtxType_Bool, BValue: false},
		ActType{OpCode: CmdInsertInContext, Id: "def", DataType: CtxType_Bool, BValue: true},
		ActType{OpCode: CmdInsertInContext, Id: "rrr", DataType: CtxType_Int, NValue: 0},
		ActType{OpCode: CmdInsertInContext, Id: "sss", DataType: CtxType_Str, SValue: "sss"},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "sss"},
			tok.Token{TokNo: gen.Tok_PIPE, Match: "|"},
			tok.Token{TokNo: gen.Tok_ID, Match: "StrToUpper"},
		}},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Str, SValue: "SSS"},
		ActType{OpCode: CmdDone},
	}},

	// Pipe Fitting 301
	{Test: "0146", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdInsertInContext, Id: "xyz", DataType: CtxType_Bool, BValue: false},
		ActType{OpCode: CmdInsertInContext, Id: "def", DataType: CtxType_Bool, BValue: true},
		ActType{OpCode: CmdInsertInContext, Id: "rrr", DataType: CtxType_Int, NValue: 0},
		ActType{OpCode: CmdInsertInContext, Id: "sss", DataType: CtxType_Str, SValue: "sss"},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "sss"},
			tok.Token{TokNo: gen.Tok_PIPE, Match: "|"},
			tok.Token{TokNo: gen.Tok_ID, Match: "x_test_str1"},
			tok.Token{TokNo: gen.Tok_OP, Match: "("},
			tok.Token{TokNo: gen.Tok_Str0, Match: "ted", DataType: CtxType_Str, CurValue: "ted"},
			tok.Token{TokNo: gen.Tok_CL, Match: ")"},
		}},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Str, SValue: "SSSted"},
		ActType{OpCode: CmdDone},
	}},

	// Pipe Fitting 302 - with a bad function (not a function)
	{Test: "0147", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdInsertInContext, Id: "xyz", DataType: CtxType_Bool, BValue: false},
		ActType{OpCode: CmdInsertInContext, Id: "def", DataType: CtxType_Bool, BValue: true},
		ActType{OpCode: CmdInsertInContext, Id: "rrr", DataType: CtxType_Int, NValue: 0},
		ActType{OpCode: CmdInsertInContext, Id: "sss", DataType: CtxType_Str, SValue: "sss"},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "sss"},
			tok.Token{TokNo: gen.Tok_PIPE, Match: "|"},
			tok.Token{TokNo: gen.Tok_ID, Match: "xyz"},
			tok.Token{TokNo: gen.Tok_OP, Match: "("},
			tok.Token{TokNo: gen.Tok_Str0, Match: "ted", DataType: CtxType_Str, CurValue: "ted"},
			tok.Token{TokNo: gen.Tok_CL, Match: ")"},
		}},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Int, NValue: 0},
		ActType{OpCode: CmdErrorsContain, Id: "Error00066", NValue: 1}, // One error containing string
		ActType{OpCode: CmdDone},
	}},

	// Pipe 303 - eval.SetErrorInfo(&TkL, "Error (Error00067): Function refrenced %s missing () to make call\n", match)
	// Pipe Fitting 302 - with a bad function (not a function)
	{Test: "0148", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdInsertInContext, Id: "xyz", DataType: CtxType_Bool, BValue: false},
		ActType{OpCode: CmdInsertInContext, Id: "def", DataType: CtxType_Bool, BValue: true},
		ActType{OpCode: CmdInsertInContext, Id: "rrr", DataType: CtxType_Int, NValue: 0},
		ActType{OpCode: CmdInsertInContext, Id: "sss", DataType: CtxType_Str, SValue: "sss"},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "sss"},
			tok.Token{TokNo: gen.Tok_PIPE, Match: "|"},
			tok.Token{TokNo: gen.Tok_ID, Match: "xxx"},
			tok.Token{TokNo: gen.Tok_OP, Match: "("},
			tok.Token{TokNo: gen.Tok_Str0, Match: "ted", DataType: CtxType_Str, CurValue: "ted"},
			tok.Token{TokNo: gen.Tok_CL, Match: ")"},
		}},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Int, NValue: 0},
		ActType{OpCode: CmdErrorsContain, Id: "Error00066", NValue: 1}, // One error containing string
		ActType{OpCode: CmdDone},
	}},

	//	Pipe 304 - eval.SetErrorInfo(&TkL, "Error (Error00068): Token found after pipe is not an ID, found %s\n", match)
	{Test: "0148", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdInsertInContext, Id: "xyz", DataType: CtxType_Bool, BValue: false},
		ActType{OpCode: CmdInsertInContext, Id: "def", DataType: CtxType_Bool, BValue: true},
		ActType{OpCode: CmdInsertInContext, Id: "rrr", DataType: CtxType_Int, NValue: 0},
		ActType{OpCode: CmdInsertInContext, Id: "sss", DataType: CtxType_Str, SValue: "sss"},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "sss"},
			tok.Token{TokNo: gen.Tok_PIPE, Match: "|"},
			tok.Token{TokNo: gen.Tok_NUM, Match: "12", DataType: CtxType_Int, CurValue: 12},
			tok.Token{TokNo: gen.Tok_OP, Match: "("},
			tok.Token{TokNo: gen.Tok_Str0, Match: "ted", DataType: CtxType_Str, CurValue: "ted"},
			tok.Token{TokNo: gen.Tok_CL, Match: ")"},
		}},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Int, NValue: 0},
		ActType{OpCode: CmdErrorsContain, Id: "Error00068", NValue: 1}, // One error containing string
		ActType{OpCode: CmdDone},
	}},

	// Pipe Fitting 305 //	eval.SetErrorInfo(&TkL, "Error (Eval00035): Extra tokens found at end of expressions, %s\n", eval.ListTokens())
	{Test: "0146", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdInsertInContext, Id: "xyz", DataType: CtxType_Bool, BValue: false},
		ActType{OpCode: CmdInsertInContext, Id: "def", DataType: CtxType_Bool, BValue: true},
		ActType{OpCode: CmdInsertInContext, Id: "rrr", DataType: CtxType_Int, NValue: 0},
		ActType{OpCode: CmdInsertInContext, Id: "sss", DataType: CtxType_Str, SValue: "sss"},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "sss"},
			tok.Token{TokNo: gen.Tok_PIPE, Match: "|"},
			tok.Token{TokNo: gen.Tok_ID, Match: "x_test_str1"},
			tok.Token{TokNo: gen.Tok_OP, Match: "("},
			tok.Token{TokNo: gen.Tok_Str0, Match: "ted", DataType: CtxType_Str, CurValue: "ted"},
			tok.Token{TokNo: gen.Tok_CL, Match: ")"},
			tok.Token{TokNo: gen.Tok_NUM, Match: "12", DataType: CtxType_Int, CurValue: 12},
			tok.Token{TokNo: gen.Tok_NUM, Match: "12", DataType: CtxType_Int, CurValue: 12},
		}},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Int, NValue: 0},
		ActType{OpCode: CmdErrorsContain, Id: "Eval00035", NValue: 1}, // One error containing string
		ActType{OpCode: CmdDone},
	}},

	// Pipe Fitting 306 - error at end of pipe
	{Test: "0145", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdInsertInContext, Id: "xyz", DataType: CtxType_Bool, BValue: false},
		ActType{OpCode: CmdInsertInContext, Id: "def", DataType: CtxType_Bool, BValue: true},
		ActType{OpCode: CmdInsertInContext, Id: "rrr", DataType: CtxType_Int, NValue: 0},
		ActType{OpCode: CmdInsertInContext, Id: "sss", DataType: CtxType_Str, SValue: "sss"},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "sss"},
			tok.Token{TokNo: gen.Tok_PIPE, Match: "|"},
			tok.Token{TokNo: gen.Tok_ID, Match: "StrToUpper"},
			tok.Token{TokNo: gen.Tok_PIPE, Match: "|"},
			tok.Token{TokNo: gen.Tok_ID, Match: "int"},
		}},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Int, NValue: 0},
		ActType{OpCode: CmdErrorsContain, Id: "Eval00035", NValue: 1}, // One error containing string
		ActType{OpCode: CmdDone},
	}},

	// Finish the assignment stuff - test cases
	{Test: "0146", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdInsertInContext, Id: "abc", DataType: CtxType_Int, NValue: 5},
		ActType{OpCode: CmdDumpContext},
		ActType{OpCode: CmdEval, Item: []tok.Token{
			tok.Token{TokNo: gen.Tok_ID, Match: "abc"},
			tok.Token{TokNo: gen.Tok_EQ, Match: "="},
			tok.Token{TokNo: gen.Tok_Str0, Match: "ted", DataType: CtxType_Str, CurValue: "ted"},
		}},
		ActType{OpCode: CmdDumpContext},
		// ActType{OpCode: CmdCmpTo, DataType: CtxType_Int, NValue: 22},
		ActType{OpCode: CmdCmpTo, DataType: CtxType_Str, SValue: "ted"},
		ActType{OpCode: CmdDone},
	}},

	// pass float to "bool(1.2)" see what happens
	// pass wrong number of parameters to type cast functions - check

	// - use := (string)
	// - use := (bool)
	// - use := [ a, b, c ]

	// xyzzy - test conversion of JSON data -> Internal Format
	// xyzzy - test JSON nested.

	// *4. type casts // xyzzy - if int/float - type cast then cast								xyzzy							ImplementTypeCast
	// xyzzy - use 0.0
	// xyzzy - use 0.0e0

	// xyzzy - use 0x	-- Not Implemented Yet
	// xyzzy - use ?:	-- Conditional if
	// xyzzy - use | 	-- Test/Write pipe

	// xyzzy - use float compare threshold

	// xyzzy - use - (binary) - on hash
	// xyzzy - use + (binary) - on hash
	// xyzzy - check for errors if = and not defined

	// ---------------------------------------------------------------------------------------------------------------------------------

	// xyzzy - use ~	 -- Not Implementd Yet (R.E. Compare)
	// xyzzy - use ~=	 -- Not Implementd Yet (R.E. Sbustitute)		a ~= "s/Re/Replace/"
	// xyzzy - mystry ops ~~, ?=, ~=, ===, =~=

	// ---------------------------------------------------------------------------------------------------------------------------------

}

var g_global = 0

var db_test01 = false

func Test_St01(t *testing.T) {

	// SymbolTable := NewSymbolTable()

	for ii, vv := range Ev01Test {
		_ = ii

		if !vv.SkipTest {

			fmt.Printf("\n\n*************************************************************************************************************\n")
			fmt.Printf("* Test %d/%s Start ----------------------------------------------------- \n", ii+1, vv.Test)
			fmt.Printf("*************************************************************************************************************\n")
			Context := NewContextType()

			var tr tok.Token
			var evalData *EvalType

			// Implement a quick - fetch execute macine to test - the SymbolTable
			for pc, ww := range vv.Actions {

				switch ww.OpCode {
				case CmdNop:
					// Do Nothing - comment out command
				case CmdDone:
					fmt.Printf("All Done\n")
				case CmdEval:
					// as, err := SymbolTable.LookupSymbol(ww.Item)
					evalData = &EvalType{
						Pos: 0,
						Ctx: Context,
						Mm:  ww.Item, // mt.TokVal[n:m], // []tok.Token
					}
					evalData.InitFunctions()
					tr = evalData.PresTop()
					fmt.Printf("Results %s ----------------------------------------------------- \n", sizlib.SVarI(tr))
					// fmt.Printf("!!! %v %T\n", tr.CurValue, tr.CurValue)
				case CmdCmpTo:
					switch tr.CurValue.(type) {
					case int:
						if tr.CurValue.(int) != ww.NValue {
							t.Errorf("%04s: %d error, expected %d got %d, found it, %s\n", vv.Test, pc, ww.NValue, tr.CurValue, com.LF())
						}
					case float64:
						if tr.CurValue.(float64) != ww.FValue {
							t.Errorf("%04s: %d error, expected %v got %v, found it, %s\n", vv.Test, pc, ww.FValue, tr.CurValue, com.LF())
						}
					case bool:
						if tr.CurValue.(bool) != ww.BValue {
							t.Errorf("%04s: %d error, expected %v got %v, found it, %s\n", vv.Test, pc, ww.BValue, tr.CurValue, com.LF())
						}
					case string:
						if tr.CurValue.(string) != ww.SValue {
							t.Errorf("%04s: %d error, expected %v got %v, found it, %s\n", vv.Test, pc, ww.SValue, tr.CurValue, com.LF())
						}
					}

				// ActType{OpCode: CmdErrorsContain, Id:"Eval00001", NValue:1},		// One error containing string
				case CmdErrorsContain:
					fmt.Printf("***************** Check for errors %s count %d\n", ww.Id, ww.NValue)
					if !tr.Error {
						t.Errorf("%04s: %d error, expected error - did not get one, found it, %s\n", vv.Test, pc, com.LF())
					} // else if strings.Contains(tr.ErrorMsg, ww.Id) {
					// 	t.Errorf("%04s: %d error, expected error with %s in it - did not find that, got %s, found it, %s\n", vv.Test, pc, ww.Id, tr.ErrorMsg, com.LF())
					// }

				case CmdInsertInContext:
					switch ww.DataType {
					case CtxType_Int:
						Context.SetInContext(ww.Id, ww.DataType, ww.NValue)
					case CtxType_Float:
						Context.SetInContext(ww.Id, ww.DataType, ww.FValue)
					case CtxType_Bool:
						Context.SetInContext(ww.Id, ww.DataType, ww.BValue)
					case CtxType_Str:
						Context.SetInContext(ww.Id, ww.DataType, ww.SValue)
					case CtxType_Func:
						if evalData != nil {
							evalData.DclFunction(ww.Id, ww.FxValue)
						} else {
							Context.SetInContext(ww.Id, ww.DataType, ww.FxValue)
						}
					case CtxType_ArrayOf:
						Context.SetInContext(ww.Id, ww.DataType, []tok.Token{})
					case CtxType_MapOf:
						Context.SetInContext(ww.Id, ww.DataType, make(map[string]tok.Token))
					}
				case CmdDumpContext:
					fmt.Printf("Dump Context - Dump is:\n")
					Context.DumpContext()

				case CmdSetGlob: //   = 10 // Set global to NValue
					g_global = ww.NValue
				case CmdCmpGlob: //  = 11 // Set global to NValue
					if g_global != ww.NValue {
						t.Errorf("%04s: %d error, expected %d in g_global, got %d, found it, %s\n", vv.Test, pc, ww.NValue, g_global, com.LF())
					}

				case CmdReadJson:
					id := ww.Id
					path := ww.Data
					file, err := ioutil.ReadFile(path)
					if err != nil {
						fmt.Printf("Error(10014): %v, %s, Config File:%s\n", err, com.LF(), path)
						return
					}
					file = []byte(strings.Replace(string(file), "\t", " ", -1)) // file = []byte(ReplaceString(string(file), "^[ \t][ \t]*//.*$", ""))

					// Check beginning of file if "{" then MapOf, if "[" Array, else look at single value
					if strings.HasPrefix(string(file), "{") {

						jsonData := make(map[string]interface{})

						err = json.Unmarshal(file, &jsonData)
						if err != nil {
							fmt.Printf("Error(10012): %v, %s, Config File:%s\n", err, com.LF(), path)
							return
						}

						Context.SetInContext(id, CtxType_MapOf, jsonData)

					} else if strings.HasPrefix(string(file), "[") {

						jsonData := make([]interface{}, 0, 100)

						err = json.Unmarshal(file, &jsonData)
						if err != nil {
							fmt.Printf("Error(10012): %v, %s, Config File:%s\n", err, com.LF(), path)
							return
						}

						Context.SetInContext(id, CtxType_ArrayOf, jsonData)

					} else {

						var jsonData interface{}

						err = json.Unmarshal(file, &jsonData)
						if err != nil {
							fmt.Printf("Error(10012): %v, %s, Config File:%s\n", err, com.LF(), path)
							return
						}

						switch jsonData.(type) {
						case float64:
							if math.Floor(jsonData.(float64)) == jsonData.(float64) {
								Context.SetInContext(id, CtxType_Int, int(jsonData.(float64)))
							} else {
								Context.SetInContext(id, CtxType_Float, jsonData)
							}
						case string:
							Context.SetInContext(id, CtxType_Str, jsonData)
						case bool:
							Context.SetInContext(id, CtxType_Bool, jsonData)
						}

					}
				}
				// xyzzy - add to conctext - remove from context - setup context commands
			}
		}
	}

}
