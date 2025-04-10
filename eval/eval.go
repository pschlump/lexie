package eval

/*

	1. Check and validate error messages
	2. add to eval results tests - figure out what each of the test SHOULD result int - what should be in the symbol table
		and add test cases / eval results - change Error found stuff to work and validate with correct errors


-- by eob --

	Improve Code:
		1. Get rid of "DataType" and use xxx.(type)

	Feature:
		1. Fancy Pipes!

	Feature:
		1. Fetch data directly inside the templates		{% x = ( readUrl("....") | XML_To_Json ) %}

	Test:
		1. Build the lexie, rigno websites with it.

	Feature:																																1hr
		3. // xyzzyPMHash - operators +, - on Arrays, Hashes -- Set operations on arrays/hashes/maps
			a - b
			a == b
			a + b


------------------ ------------------ ------------------ ------------------ ------------------ ------------------ ------------------ ------------------ ------------------ ------------------
------------------ ------------------ ------------------ ------------------ ------------------ ------------------ ------------------ ------------------ ------------------ ------------------

	Implement As Macros:
				__line__, __file__
				__s_line__		"123"
				__ss_line__		'123'
				__sd_line__		"123"
				__sh_line__		`123`
				__sbs_line__	\"123\"


	Simple:
		Tokens not used ++, --, ~, ~=, ===, =~=
		; - to split statements?

	Simple:
		1. Add option to generate ../gen/gen_tok.go (with name at top) from lexie generator

	Harder:
		1. Fix/Test NCCLs

	Optimize:
		1. Do a performance tests between ".Call" arbitrary and known type calls

		   Any function that is F(s,...) -> s

		   once you have the Type of the method or function in hand, you can then call on of those:

				// In returns the type of a function type's i'th input parameter.
				// It panics if the type's Kind is not Func.
				// It panics if i is not in the range [0, NumIn()).
				In(i int) Type

				// NumIn returns a function type's input parameter count.
				// It panics if the type's Kind is not Func.
				NumIn() int

				// NumOut returns a function type's output parameter count.
				// It panics if the type's Kind is not Func.
				NumOut() int

				// Out returns the type of a function type's i'th output parameter.
				// It panics if the type's Kind is not Func.
				// It panics if i is not in the range [0, NumOut()).
				Out(i int) Type

		2. Look at the stuff in string.* package and do calls to all of them with .Call






------------------ ------------------ ------------------ ------------------ ------------------ ------------------ ------------------ ------------------ ------------------ ------------------
------------------ ------------------ ------------------ ------------------ ------------------ ------------------ ------------------ ------------------ ------------------ ------------------

Ringo - the web framework for developers with deadlines that care about performance
	(i.e. get it done on time and scalable to the real world)

------------------ ------------------ ------------------ ------------------ ------------------ ------------------ ------------------ ------------------ ------------------ ------------------
------------------ ------------------ ------------------ ------------------ ------------------ ------------------ ------------------ ------------------ ------------------ ------------------

Steve Jobs, Elon Musk ... obsess over making “insanely great products”
======================================================================

You miss 100% of the shots you don't take
=========================================

...your business at its core must wage a constant battle to find the
truth. And for a business, "the truth" means value. As an entrepreneur,
you should be spending every day getting to the core of what "value"
means to your customers.

Your job is to build a venture that will search for real value and
"the truth" of what your client really needs.
=======================================================================



------------------ ------------------ ------------------ ------------------ ------------------ ------------------ ------------------ ------------------ ------------------ ------------------
------------------ ------------------ ------------------ ------------------ ------------------ ------------------ ------------------ ------------------ ------------------ ------------------
1. Do actual templates
	2. Implment "extend"
	3. Implment "import_library"
	4. Implment "template"..."endtemplate"

	CLI -c config.json -i input.tmpl.html  -o result.output.html -t command-trace.out

		command-trace.out - Included Dependent Files -- in JSON format
		../test02


------------------ ------------------ ------------------ ------------------ ------------------ ------------------ ------------------ ------------------ ------------------ ------------------
------------------ ------------------ ------------------ ------------------ ------------------ ------------------ ------------------ ------------------ ------------------ ------------------
Goals: 1. 100% coverage of eval
	2. Works for all operators defined on all data types or gives appropriate errors
	3. Returns errors for all non-syntacticly correct castes - check for extra tokens at end
		ad "as" Var - right hand assignment
	4. Works with different flags set for behavior - all cases
	5. Works with time data type also
	5. Works with strings data type also

------------------ ------------------ ------------------ ------------------ ------------------ ------------------ ------------------ ------------------ ------------------ ------------------
------------------ ------------------ ------------------ ------------------ ------------------ ------------------ ------------------ ------------------ ------------------ ------------------


*/

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/pschlump/dbgo"
	"github.com/pschlump/lexie/com"
	"github.com/pschlump/lexie/gen"
	"github.com/pschlump/lexie/tok"
	// "../../../go-lib/sizlib"
)

type EvalType struct {
	Pos     int          //
	Ctx     *ContextType //
	Mm      []tok.Token  //
	Lax_DCL bool         // If true, then auto declare on all assignments
	// Lax_ArrayHash bool         // If true, then [a] == a, { "x":a } == a
	// FloatEQThrUse bool    // If true then compare floats using threshold
	// FloatEQThr    float64 //
	PrintErrorMsg bool
	TestCase      string
}

// -------------------------------------------------------------------------------------------------------------------------------------------------------
//func (eval *EvalType) SetConfig(laxDcl bool, laxType bool, laxArrayHash bool, useFloatThr bool, floatThr float64) {
//	eval.Lax_DCL = laxDcl
//	eval.Lax_ArrayHash = laxArrayHash
//	eval.FloatEQThrUse = useFloatThr
//	eval.FloatEQThr = floatThr
//}

// -------------------------------------------------------------------------------------------------------------------------------------------------------
func (eval *EvalType) SetErrorInfo(rv *tok.Token, format string, args ...interface{}) {
	t := eval.Pos
	if t >= len(eval.Mm) {
		t = len(eval.Mm) - 1
	} // else if t < 0 {
	// 	t = 0
	// }
	if len(eval.Mm) > 0 {
		rv.LineNo = eval.Mm[t].LineNo
		rv.ColNo = eval.Mm[t].ColNo
		rv.FileName = eval.Mm[t].FileName
	}
	rv.CurValue = 0 // Add in a magic 0 or something for doing ops off of end of tokens
	rv.DataType = CtxType_Int
	rv.CoceLocation = dbgo.LF(2)
	rv.LValue = false
	rv.Error = true
	rv.ErrorMsg = fmt.Sprintf(format, args...)
	if eval.PrintErrorMsg {
		fmt.Printf("%s", rv.ErrorMsg)
	}
}

// -------------------------------------------------------------------------------------------------------------------------------------------------------
func (eval *EvalType) ParsePlist() (rv []tok.Token) {
	rv = make([]tok.Token, 0, 5)
	opTk := eval.Mm[eval.Pos].TokNo
	if opTk == gen.Tok_OP {
		eval.Pos++
		// fmt.Printf("IN ParsePlist -00-: %s\n", dbgo.LF())
	}
	for eval.Pos < len(eval.Mm) {
		// fmt.Printf("IN ParsePlist -loop top-: %s\n", dbgo.LF())
		opTk := eval.Mm[eval.Pos].TokNo
		if opTk != gen.Tok_CL && opTk != gen.Tok_COMMA { // not ) not ,
			// fmt.Printf("IN ParsePlist -return-: %s\n", dbgo.LF())
			Tk := eval.PresG()
			rv = append(rv, Tk)
			if Tk.Error {
				return
			}
			if eval.Pos < len(eval.Mm) {
				opTk = eval.Mm[eval.Pos].TokNo
				if opTk == gen.Tok_COMMA {
					// fmt.Printf("IN ParsePlist -BB-: %s\n", dbgo.LF())
					eval.Pos++
					opTk = eval.Mm[eval.Pos].TokNo
				} else if opTk == gen.Tok_CL {
					eval.Pos++
					// fmt.Printf("IN ParsePlist -CC-: %s\n", dbgo.LF())
					return
				}
			} // else {
			// 	return
			// }
		} else {
			Tk := tok.Token{Error: true, ErrorMsg: "Error (Eval00001): Invalid Parameter List\n"}
			rv = append(rv, Tk)
			return
		}
		// fmt.Printf("IN ParsePlist -A-: %s\n", dbgo.LF())
	}
	return
}

func BoundArrayIndex(i, min, max int) int {
	if i < min {
		i = min
	}
	if i >= max {
		i = max - 1
	}
	//if t >= len(eval.Mm) {
	//	t = len(eval.Mm) - 1
	//} else if t < 0 {
	//	t = 0
	//}
	return i
}

// -------------------------------------------------------------------------------------------------------------------------------------------------------
func (eval *EvalType) CallFunction(match string, plist []tok.Token) (rv tok.Token) {

	// fmt.Printf("=============================================================================================================\n")
	// fmt.Printf("IN CallFunction: %s\n", dbgo.LF())
	// fmt.Printf("IN CallFunction: Call >>>%s<<< with %s\n", match, com.SVarI(plist))
	// fmt.Printf("=============================================================================================================\n")

	t := eval.Pos - 1
	t = BoundArrayIndex(t, 0, len(eval.Mm))
	if len(eval.Mm) > 0 {
		rv.LineNo = eval.Mm[t].LineNo
		rv.ColNo = eval.Mm[t].ColNo
		rv.FileName = eval.Mm[t].FileName
	}
	rv.CoceLocation = dbgo.LF(2)
	rv.LValue = false
	rv.Error = false
	rv.ErrorMsg = ""
	rv.DataType = CtxType_Int
	rv.CurValue = 0

	// func (ctx *ContextType) Call(name string, params ...interface{}) (result []reflect.Value, err error) {

	p_plist := make([]interface{}, len(plist))
	for ii, vv := range plist {
		// xyzzy - coerce params to correct type?
		p_plist[ii] = vv.CurValue
	}

	ww, err := eval.Ctx.Call(match, p_plist...)
	if err != nil {
		eval.SetErrorInfo(&rv, "%s", err)
	} else {
		err_pos := -1
		var err error
		// Pass 1 check for error return type - if found the if not nil then process as error else pass 2 - skiping error
		for ii, val := range ww {
			i := val.Interface()
			switch val.Type().String() {
			case "error":
				err_pos = ii
				err = nil
				if i != nil {
					err = errors.New(fmt.Sprintf("%v", val.Interface()))
				}
			}
		}
		if err_pos >= 0 && err != nil {
			// fmt.Printf("Found Error Not Nil - At: %s\n", dbgo.LF())
			rv.DataType = CtxType_Int
			rv.CurValue = 0
			rv.CoceLocation = "Call to funciton " + match
			rv.LValue = false
			rv.Error = true
			rv.ErrorMsg = fmt.Sprintf("%s", err)
			if eval.PrintErrorMsg {
				fmt.Printf("%s", rv.ErrorMsg)
			}
		} else {
			for ii, val := range ww {
				if ii != err_pos {
					switch val.Kind() {
					case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
						rv.DataType = CtxType_Int
						rv.CurValue = int(val.Int())
					case reflect.String:
						rv.DataType = CtxType_Str
						rv.CurValue = val.String()
					case reflect.Bool:
						rv.DataType = CtxType_Bool
						rv.CurValue = val.Bool()
					case reflect.Float32, reflect.Float64:
						rv.DataType = CtxType_Float
						rv.CurValue = val.Float()
					default:
						// fmt.Printf("Error (Eval00002): Can't handle type: %s as a return type from a function.", val.Type().String())
						rv.DataType = CtxType_Int
						rv.CurValue = 0
						rv.CoceLocation = dbgo.LF(1)
						rv.LValue = false
						rv.Error = true
						rv.ErrorMsg = fmt.Sprintf("Error (Eval00003): Can't handle type: %s as a return type from a function.", val.Type().String())
						if eval.PrintErrorMsg {
							fmt.Printf("%s", rv.ErrorMsg)
						}
					}
				}
			}
		}
	}

	return
}

/*
From: https://code.google.com/p/go/source/browse/src/pkg/fmt/scan.go?name=release-branch.go1.1#994

				switch v := ptr.Elem(); v.Kind() {
                case reflect.Bool:
                        v.SetBool(s.scanBool(verb))
                case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
                        v.SetInt(s.scanInt(verb, v.Type().Bits()))
                case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
                        v.SetUint(s.scanUint(verb, v.Type().Bits()))
                case reflect.String:
                        v.SetString(s.convertString(verb))
                case reflect.Slice:
                        // For now, can only handle (renamed) []byte.
                        typ := v.Type()
                        if typ.Elem().Kind() != reflect.Uint8 {
                                s.errorString("Scan: can't handle type: " + val.Type().String())
                        }
                        str := s.convertString(verb)
                        v.Set(reflect.MakeSlice(typ, len(str), len(str)))
                        for i := 0; i < len(str); i++ {
                                v.Index(i).SetUint(uint64(str[i]))
                        }
                case reflect.Float32, reflect.Float64:
                        s.skipSpace(false)
                        s.notEOF()
                        v.SetFloat(s.convertFloat(s.floatToken(), v.Type().Bits()))
                case reflect.Complex64, reflect.Complex128:
                        v.SetComplex(s.scanComplex(verb, v.Type().Bits()))
                default:
                        s.errorString("Scan: can't handle type: " + val.Type().String())
*/

// -------------------------------------------------------------------------------------------------------------------------------------------------------
// Create the function with Fx ptr in the symbol table.
func (eval *EvalType) DclFunction(id string, fx interface{}) {
	eval.Ctx.SetInContext(id, CtxType_Func, fx)
	return
}

// -------------------------------------------------------------------------------------------------------------------------------------------------------
func MapIsEmpty(t map[string]tok.Token) bool {
	if t == nil {
		return true
	}
	for _, _ = range t {
		return false
	}
	return true
}

// -------------------------------------------------------------------------------------------------------------------------------------------------------
//
//	rv.CurValue = mapLength(x.CurValue.(map[string]tok.Token))
func MapLength(t map[string]tok.Token) (l int) {
	l = 0
	// fmt.Printf("At: %s\n", dbgo.LF())
	// fmt.Printf("t=%s\n", com.SVarI(t))
	for _, _ = range t {
		// fmt.Printf("At: %s\n", dbgo.LF())
		l++
	}
	// fmt.Printf("Returing %d At: %s\n", l, dbgo.LF())
	return
}

// -------------------------------------------------------------------------------------------------------------------------------------------------------
// ToDo:
//
//	IsEmpty - array/map - if len == 0, rueturn true (faster than len)
//
// -------------------------------------------------------------------------------------------------------------------------------------------------------
func x_len_e(x interface{}) (rv int) {
	// fmt.Printf("len type=%T\n", x)
	switch x.(type) {
	case []tok.Token:
		// fmt.Printf("At: %s\n", dbgo.LF())
		rv = len(x.([]tok.Token))
	case map[string]tok.Token:
		// fmt.Printf("At: %s\n", dbgo.LF())
		rv = MapLength(x.(map[string]tok.Token))
	default:
		// fmt.Printf("At: %s\n", dbgo.LF())
		rv = 0
	}
	return
}

func x_len(x interface{}) (rv int, e error) {
	// fmt.Printf("len type=%T\n", x)
	switch x.(type) {
	case []tok.Token:
		// fmt.Printf("At: %s\n", dbgo.LF())
		rv = len(x.([]tok.Token))
	case map[string]tok.Token:
		// fmt.Printf("At: %s\n", dbgo.LF())
		rv = MapLength(x.(map[string]tok.Token))
	default:
		// fmt.Printf("At: %s\n", dbgo.LF())
		e = errors.New(fmt.Sprintf("Error (Eval00004): Invalid type. len() works on arrays and maps.  Type supplied %T\n", x))
		rv = 0
	}
	return
}

func x_test_ret_float(i int) float32 {
	return 1.2
}
func x_test_ret_bool(i int) bool {
	return false
}

func x_float_type_cast(x interface{}) (rv float64, e error) {
	e = nil
	switch x.(type) {
	case int:
		rv = float64(x.(int))
	case float64:
		rv = x.(float64)
	case bool:
		if x.(bool) {
			rv = 1
		} else {
			rv = 0
		}
	default:
		e = errors.New(fmt.Sprintf("Error (Eval00005): Invalid type conversion, attempt to convert from %T to float\n", x))
		rv = 0
	}
	return
}

func x_bool_type_cast(x interface{}) (rv bool, e error) {
	// fmt.Printf("BOOL: x=%v, %T %s\n", x, x, dbgo.LF())
	rv = false
	// fmt.Printf("At: %s\n", dbgo.LF())
	switch x.(type) {
	case int:
		// fmt.Printf("At: %s\n", dbgo.LF())
		if x.(int) == 0 {
			rv = false
		}
	case bool:
		// fmt.Printf("At: %s\n", dbgo.LF())
		rv = x.(bool)
	case string:
		// fmt.Printf("At: %s\n", dbgo.LF())
		if len(x.(string)) != 0 {
			rv = true
		}
	case []tok.Token:
		// fmt.Printf("At: %s\n", dbgo.LF())
		if len(x.([]tok.Token)) != 0 {
			rv = true
		}
	case map[string]tok.Token:
		// fmt.Printf("At: %s\n", dbgo.LF())
		if !MapIsEmpty(x.(map[string]tok.Token)) {
			// fmt.Printf("At: %s\n", dbgo.LF())
			rv = true
		}
	default:
		e = errors.New(fmt.Sprintf("Error (Eval00006): Invalid type conversion, attempt to convert from %T to bool\n", x))
	}
	// fmt.Printf("At: %s\n", dbgo.LF())
	return
}

func x_int_type_cast(x interface{}) (rv int, e error) {
	rv = 0
	switch x.(type) {
	case int:
		rv = x.(int)
	case bool:
		if x.(bool) {
			rv = 1
		}
	case float64:
		rv = int(x.(float64))
	default:
		e = errors.New(fmt.Sprintf("Error (Eval00007): Invalid type conversion, attempt to convert from %T to int\n", x))
	}
	return
}

func x_StrToUpper(s string) string {
	// fmt.Printf("At: %s\n", dbgo.LF())
	return strings.ToUpper(s)
}

func x_test_str1(s string, t string) string {
	// fmt.Printf("At: %s\n", dbgo.LF())
	return strings.ToUpper(s) + t
}

type BobType struct {
	Ooops  int
	Ooopsy int
}

func x_test_str2(s string, t string) (rv BobType) {
	return
}

// -------------------------------------------------------------------------------------------------------------------------------------------------------
func (eval *EvalType) InitFunctions() {
	eval.DclFunction("float", x_float_type_cast)           //
	eval.DclFunction("int", x_int_type_cast)               //
	eval.DclFunction("bool", x_bool_type_cast)             //
	eval.DclFunction("len", x_len)                         // length of item, must be array or hash
	eval.DclFunction("len_e", x_len_e)                     // Length ignoring errors, 0 returned instead
	eval.DclFunction("StrToUpper", x_StrToUpper)           //
	eval.DclFunction("yStrToUpper", strings.ToUpper)       //
	eval.DclFunction("x_test_ret_float", x_test_ret_float) // Test code
	eval.DclFunction("x_test_ret_bool", x_test_ret_bool)   // Test code
	eval.DclFunction("x_test_str1", x_test_str1)           //
	eval.DclFunction("x_test_str2", x_test_str2)           // Returns an un-handlable type
}

// -------------------------------------------------------------------------------------------------------------------------------------------------------
func (eval *EvalType) Pres0() (rv tok.Token) {
	// fmt.Printf("Pos:%d At: %s\n", eval.Pos, dbgo.LF())
	if eval.Pos < len(eval.Mm) {
		switch eval.Mm[eval.Pos].TokNo {
		case gen.Tok_ID:
			// fmt.Printf("At: %s\n", dbgo.LF())
			// fmt.Printf("Doing a lookup of >>>%s<<<\n", eval.Mm[eval.Pos].Match)
			match := eval.Mm[eval.Pos].Match
			if val0, t, f := eval.Ctx.GetFromContext(eval.Mm[eval.Pos].Match); f {
				// fmt.Printf("   Found, Type=%d/%s, %s\n", t, eval.Ctx.NameOfType(t), dbgo.LF())
				switch t {
				case CtxType_Func:
					// fmt.Printf("At: %s\n", dbgo.LF())
					eval.Pos++
					if eval.Pos < len(eval.Mm) && eval.Mm[eval.Pos].TokNo == gen.Tok_OP {
						Plist := eval.ParsePlist()
						rv = eval.CallFunction(match, Plist)
						if rv.Error {
							return
						}
						// fmt.Printf("Found a Func\n")
					} else {
						eval.SetErrorInfo(&rv, "Error (Eval00008): Function refrenced %s missing () to make call\n", match)
						return
					}
				case CtxType_Int, CtxType_Str, CtxType_Bool, CtxType_Float, CtxType_ArrayOf, CtxType_MapOf, CtxType_SMapOf, CtxType_KMapOf:
					// fmt.Printf("At: %s\n", dbgo.LF())
					rv.Match = eval.Mm[eval.Pos].Match
					rv.CurValue = val0
					rv.DataType = t
					rv.LValue = true

					//default:
					//	eval.SetErrorInfo(&rv, "Error (Eval00009): ID Invalid type %d/%s for %s\n", t, eval.Ctx.NameOfType(t), eval.Mm[eval.Pos].Match)
					//	return
				}
			} else if eval.Mm[eval.Pos].CreateId {
				eval.Ctx.SetInContext(eval.Mm[eval.Pos].Match, CtxType_Int, 0)
				rv.Match = eval.Mm[eval.Pos].Match
				rv.CurValue = 0
				rv.DataType = CtxType_Int
				rv.LValue = true
			} else {
				eval.SetErrorInfo(&rv, "Error (Eval00010): ID Missing %s\n", eval.Mm[eval.Pos].Match)
				return
			}
		case gen.Tok_OP:
			// fmt.Printf("At: %s\n", dbgo.LF())
			eval.Pos++ // Step past the '('
			rv = eval.PresG()
			if rv.Error {
				return
			}
			if eval.Pos < len(eval.Mm) && eval.Mm[eval.Pos].TokNo == gen.Tok_CL {
				eval.Pos++
				// fmt.Printf("Return from recrusive () call -- Sitting on Close Paren\n")
			} else {
				eval.SetErrorInfo(&rv, "Error (Eval00011): missing ')' Close Paren\n")
				return
			}
		case gen.Tok_OP_SQ:
			// fmt.Printf("JSON data found (array) : At: %s\n", dbgo.LF())
			rv = eval.ParseJSON()
			// fmt.Printf("TestCase=%s, eval.Pos=%d\n", eval.TestCase, eval.Pos)
			if rv.Error {
				return
			}
		case gen.Tok_OP_BRACE:
			// fmt.Printf("JSON data found (hash) : At: %s\n", dbgo.LF())
			rv = eval.ParseJSON()
			if rv.Error {
				return
			}
		case gen.Tok_Str0:
			rv.CurValue = eval.Mm[eval.Pos].CurValue
			rv.DataType = eval.Mm[eval.Pos].DataType
			rv.LValue = false
		case gen.Tok_NUM:
			rv.CurValue = eval.Mm[eval.Pos].CurValue
			// fmt.Printf("Num = %d, %s\n", rv.CurValue.(int), dbgo.LF())
			rv.DataType = eval.Mm[eval.Pos].DataType
			rv.LValue = false
		case gen.Tok_Float:
			// fmt.Printf("At: %s\n", dbgo.LF())
			rv.CurValue = eval.Mm[eval.Pos].CurValue
			rv.DataType = eval.Mm[eval.Pos].DataType
			rv.LValue = false
		case gen.Tok_true, gen.Tok_false, gen.Tok_Tree_Bool:
			// fmt.Printf("At: %s\n", dbgo.LF())
			rv.DataType = CtxType_Bool
			rv.CurValue = eval.Mm[eval.Pos].CurValue
			rv.LValue = false
		}
		eval.Pos++
	}
	return
}

// -------------------------------------------------------------------------------------------------------------------------------------------------------
func (eval *EvalType) ParseJSON() (rv tok.Token) {
	// Have { or [ leading token, remember to pick off matching one at end.
	// if [ then array, if { then hash
	opTk := eval.Mm[eval.Pos].TokNo

	Adv := func() {
		if eval.Pos < len(eval.Mm) {
			eval.Pos++
		}
		opTk = gen.Tok_null
		if eval.Pos < len(eval.Mm) {
			opTk = eval.Mm[eval.Pos].TokNo
		}
	}
	Set := func() {
		opTk = gen.Tok_null
		if eval.Pos < len(eval.Mm) {
			opTk = eval.Mm[eval.Pos].TokNo
		}
	}

	switch opTk {
	case gen.Tok_OP_SQ:
		rv.DataType = CtxType_ArrayOf
		// opCLTK_s := "[]"
		r0 := make([]tok.Token, 0, 5)
		rv.CurValue = r0

		Adv()
		// fmt.Printf("IN ParseJSON-[] -00-: %s\n", dbgo.LF())

		for eval.Pos < len(eval.Mm) {
			// fmt.Printf("IN ParseJSON-%s -loop top-: %s\n", opCLTK_s, dbgo.LF())
			opTk = eval.Mm[eval.Pos].TokNo
			if opTk == gen.Tok_COMMA { //                                                   	Walk over ,
				// fmt.Printf("IN ParseJSON-%s -BB-: %s\n", opCLTK_s, dbgo.LF())
				Adv()
			}
			if opTk == gen.Tok_CL_SQ { //                                                   	Walk over ], then return
				Adv()
				// fmt.Printf("IN ParseJSON-%s -CC-: %s\n", opCLTK_s, dbgo.LF())
				return
			}
			if opTk == gen.Tok_OP_SQ || opTk == gen.Tok_OP_BRACE { //  Found [ or {
				// fmt.Printf("IN ParseJSON-%s -A-: %s\n", opCLTK_s, dbgo.LF())
				Tk := eval.ParseJSON()
				Set()
				r0 = append(r0, Tk)
				rv.CurValue = r0
				if Tk.Error {
					rv.Error = true
					rv.ErrorMsg = Tk.ErrorMsg
				}
			} else {
				// fmt.Printf("IN ParseJSON-%s -A-: %s\n", opCLTK_s, dbgo.LF())
				Tk := eval.PresG()
				Set()
				r0 = append(r0, Tk)
				rv.CurValue = r0
				if Tk.Error {
					rv.Error = true
					rv.ErrorMsg = Tk.ErrorMsg
				}
			}
			if eval.Pos < len(eval.Mm) {
				opTk = eval.Mm[eval.Pos].TokNo
			} else {
				opTk = gen.Tok_CL_SQ
				// error/recovered - added CL to close plist!
				eval.SetErrorInfo(&rv, "Error (Eval00012):  invalid JSON data.\n")
				return
			}
			// fmt.Printf("IN ParseJSON-%s -A-: %s\n", opCLTK_s, dbgo.LF())
		}

	case gen.Tok_OP_BRACE:
		rv.DataType = CtxType_MapOf
		// opCLTK_s := "{}"
		r1 := make(map[string]tok.Token)
		rv.CurValue = r1

		name := ""
		Adv()
		// fmt.Printf("IN ParseJSON-{} -00-: %s\n", dbgo.LF())

		for eval.Pos < len(eval.Mm) {
			// fmt.Printf("IN ParseJSON-%s -loop top-: %s\n", opCLTK_s, dbgo.LF())
			haveOne := true

			// fmt.Printf("IN ParseJSON-%s -value-: %s\n", opCLTK_s, dbgo.LF())
			opTk = eval.Mm[eval.Pos].TokNo

			if opTk == gen.Tok_ID { // xyzzy - ID : -- handle this also -> Name
				// fmt.Printf("IN ParseJSON-%s -ID-: %s\n", opCLTK_s, dbgo.LF())
				name = eval.Mm[eval.Pos].Match
				Adv()
			} else if opTk == gen.Tok_Str0 { // xyzzy - "name" :  -- Handle this
				// fmt.Printf("IN ParseJSON-%s -STR-: %s\n", opCLTK_s, dbgo.LF())
				name = eval.Mm[eval.Pos].Match
				Adv()
			} else if opTk == gen.Tok_CL_BRACE { //                                                   	Walk over ], then return
				Adv()
				// fmt.Printf("IN ParseJSON-%s -CC-: %s\n", opCLTK_s, dbgo.LF())
				return
			} else {
				name = ""
				// fmt.Printf("************* setting to empty *************\n")
				haveOne = false
				eval.SetErrorInfo(&rv, "Error (Eval00013):  invalid JSON data.\n")
				return
			}

			// fmt.Printf("************* name=>%s<- *************\n", name)
			if opTk == gen.Tok_COLON {
				// fmt.Printf("IN ParseJSON-%s -colon-: %s\n", opCLTK_s, dbgo.LF())
				Adv()
			}

			if opTk == gen.Tok_OP_SQ || opTk == gen.Tok_OP_BRACE { // [ or { found
				// fmt.Printf("IN ParseJSON-%s -Sub-JSON-: %s\n", opCLTK_s, dbgo.LF())
				Tk := eval.ParseJSON()
				Set()
				if haveOne {
					r1[name] = Tk
					rv.CurValue = r1
				}
				if Tk.Error {
					rv.Error = true
					rv.ErrorMsg = Tk.ErrorMsg
				}
			} else {
				// fmt.Printf("IN ParseJSON-%s -expression-: %s\n", opCLTK_s, dbgo.LF())
				Tk := eval.PresG()
				Set()
				if haveOne {
					r1[name] = Tk
					rv.CurValue = r1
				}
				// fmt.Printf("TestCase=%s, eval.Pos=%d, %s\n", eval.TestCase, eval.Pos, dbgo.LF())
				if Tk.Error {
					rv.Error = true
					rv.ErrorMsg = Tk.ErrorMsg
				}
			}

			Set()
			if opTk == gen.Tok_COMMA { //                                                   	Walk over ,
				// fmt.Printf("IN ParseJSON-%s -BB-: %s\n", opCLTK_s, dbgo.LF())
				Adv()
			} else if opTk == gen.Tok_CL_BRACE { //                                                   	Walk over }, then return
				Adv()
				// fmt.Printf("IN ParseJSON-%s -CC-: %s\n", opCLTK_s, dbgo.LF())
				return
			} else {
				// fmt.Printf("IN ParseJSON-%s -EvalThis 555 555 5555 -: %s, Pos=%d\n", opCLTK_s, dbgo.LF(), eval.Pos)
				if eval.Pos < len(eval.Mm) {
					// fmt.Printf("IN At:%s\n", dbgo.LF())
					Tk := eval.PresG()
					Set()
					if haveOne {
						// fmt.Printf("IN At:%s\n", dbgo.LF())
						r1[name] = Tk
						rv.CurValue = r1
					}
					// fmt.Printf("TestCase=%s, eval.Pos=%d, %s\n", eval.TestCase, eval.Pos, dbgo.LF())
					if Tk.Error {
						rv.Error = true
						rv.ErrorMsg = Tk.ErrorMsg
					}
				} else {
					return
				}
			}
			// fmt.Printf("IN At:%s\n", dbgo.LF())
		}

	}

	return
}

// -------------------------------------------------------------------------------------------------------------------------------------------------------
func (eval *EvalType) ParseArrayRef() (rv tok.Token) {
	return eval.PresG()
}

// -------------------------------------------------------------------------------------------------------------------------------------------------------
//
// Array Slice Ops
//
//	Expr [ : Expr ]
//	Expr [ Expr : ]
//	Expr [ : ]
//	Expr [ ExprID ]							-- 			-- Map Ref
//	Expr . Name								-- 			-- Map Ref xyzzy
//
// ArrayId(expr) [ ... ]
// MapId(expr) . Name List
//
// xyzzy1857 - not working
//
// Array is just a postfix '[' or '.' operator (Could have been -> too)
// -->> search for occurrence of "name" in hash
func (eval *EvalType) Pres1() (TkL tok.Token) {
	var Tk0, Tk1 tok.Token
	opTk := eval.Mm[eval.Pos].TokNo
	TkL = eval.Pres0()
	if TkL.Error {
		return
	}

	Adv := func() {
		if eval.Pos < len(eval.Mm) {
			eval.Pos++
		}
		if eval.Pos < len(eval.Mm) {
			opTk = eval.Mm[eval.Pos].TokNo
		} else {
			opTk = gen.Tok_null
		}
	}

	// XXX if !(TkL.CurValue.(type) == []tok.Token || TkL.CurValue.(type) == []interface{}) {
	// XXX if !eval.IsArray(TkL) {
	if TkL.DataType != CtxType_ArrayOf && TkL.DataType != CtxType_MapOf {
		return
	}

	// match := ""
	if eval.Pos < len(eval.Mm) {
		opTk = eval.Mm[eval.Pos].TokNo
	}

	// fmt.Printf("Before Loop TkL = %s, %s\n", com.SVarI(TkL), dbgo.LF())
	for eval.Pos < len(eval.Mm) {

		// fmt.Printf("At loop top eval.Pos=%d, %s\n", eval.Pos, dbgo.LF())
		if TkL.DataType != CtxType_ArrayOf && TkL.DataType != CtxType_MapOf {
			return
		}

		// Array [ a : b ] [ ...
		// Array [ : b ] . id ...
		// Array [ : ]
		// Array [ x ]
		// Array . ID/Name
		//       ^----------------------------------- Pos

		if TkL.DataType == CtxType_ArrayOf {

			// fetch data from expression for Array // output to verify
			// get array or hash from symboltable or expression and put into token

			// fmt.Printf("At 1st Adv(), after eval.Pos=%d, TokNo=%d %s\n", eval.Pos, opTk, dbgo.LF())
			// Array [ a : b ]
			// Array [ : b ]
			// Array [ : ]
			// Array [ x ]
			//       ^----------------------------------- Pos

			if opTk != gen.Tok_OP_SQ { // if it is an Array but not a '[' to ref it, then return
				return
			}

			Adv()
			// Array [ a : b ]
			// Array [ : b ]
			// Array [ : ]
			// Array [ x ]
			//         ^----------------------------------- Pos

			isRef := false
			// if eval.Mm[eval.Pos].TokNo == gen.Tok_COLON {
			if eval.Mm[eval.Pos].TokNo == gen.Tok_COLON {
				Adv()
				// Array [ : b ]
				// Array [ : ]
				//           ^----------------------------------- Pos
				// fmt.Printf("At Step Over Colon Adv(), after eval.Pos=%d, TokNo=%d %s\n", eval.Pos, opTk, dbgo.LF())
				Tk0.TokNo = gen.Tok_NUM
				Tk0.CurValue = 0
			} else {
				Tk0 = eval.ParseArrayRef()
				// Array [ a : b ]
				// Array [ x ]
				//           ^----------------------------------- Pos
				if Tk0.Error {
					TkL = Tk0
					return
				}
				if Tk0.DataType != CtxType_Int {
					TkL = Tk0
					eval.SetErrorInfo(&TkL, "Error (Eval00014):  Attempted to index array with non-integer data.\n")
					//TkL.CurValue = 0 // Add in a magic 0 or something for doing ops off of end of tokens
					//TkL.DataType = CtxType_Int
					return
				}
				// fmt.Printf("eval.Pos=%d token=%d\n", eval.Pos, eval.Mm[eval.Pos].TokNo)
				if eval.Mm[eval.Pos].TokNo == gen.Tok_COLON {
					Adv()
					// Array [ a : b ]
					//             ^----------------------------------- Pos
				} else if eval.Mm[eval.Pos].TokNo == gen.Tok_CL_SQ {
					// Array [ x ]
					//           ^----------------------------------- Pos
					isRef = true
					Adv()
					// Array [ x ]
					//             ^----------------------------------- Pos
					Tk1.TokNo = gen.Tok_NUM
					switch TkL.CurValue.(type) {
					case []tok.Token:
						Tk1.CurValue = len(TkL.CurValue.([]tok.Token))
						//case []interface{}:
						//	Tk1.CurValue = len(TkL.CurValue.([]interface{}))
					}
					goto evalIt
				} else {
					eval.SetErrorInfo(&TkL, "Error (Eval00015): invalid array reference\n")
					return
				}
			}
			// Array [ a : ]
			// Array [   : ]
			// Array [ a : b ]
			//             ^----------------------------------- Pos
			if eval.Mm[eval.Pos].TokNo == gen.Tok_CL_SQ {
				Adv()
				// Array [ a : ]
				// Array [   : ]
				//                 ^----------------------------------- Pos
				Tk1.TokNo = gen.Tok_NUM
				switch TkL.CurValue.(type) {
				case []tok.Token:
					Tk1.CurValue = len(TkL.CurValue.([]tok.Token))
					//case []interface{}:
					//	Tk1.CurValue = len(TkL.CurValue.([]interface{}))
				}
				goto evalIt
			}

			Tk1 = eval.ParseArrayRef()
			// Array [ a : b ]
			//               ^----------------------------------- Pos
			if Tk1.Error {
				TkL = Tk1
				return
			}
			if Tk1.DataType != CtxType_Int {
				TkL = Tk1
				eval.SetErrorInfo(&TkL, "Error (Eval00016):  Attempted to index array with non-integer data.\n")
				//TkL.CurValue = 0 // Add in a magic 0 or something for doing ops off of end of tokens
				//TkL.DataType = CtxType_Int
				return
			}
			if eval.Mm[eval.Pos].TokNo == gen.Tok_CL_SQ {
				Adv()
				// Array [ a : b ]
				//                 ^----------------------------------- Pos
				// goto evalIt
			} else {
				eval.SetErrorInfo(&TkL, "Error (Eval00017): invalid array reference\n")
				return
			}
		evalIt:

			n := Tk0.CurValue.(int)
			if isRef {
				switch TkL.CurValue.(type) {
				case []tok.Token:
					x := TkL.CurValue.([]tok.Token)
					TkL.DataType = x[n].DataType
					TkL.CurValue = x[n].CurValue
					TkL.LValue = true
					//case []interface{}:
					//	x := TkL.CurValue.([]interface{})
					//	TkL.CurValue = x[n]
				}
			} else {
				m := Tk1.CurValue.(int)
				// fmt.Printf("N=%d M=%d CurValue=%s\n", n, m, com.SVar(TkL.CurValue))
				switch TkL.CurValue.(type) {
				case []tok.Token:
					x := TkL.CurValue.([]tok.Token)
					TkL.DataType = CtxType_ArrayOf
					TkL.CurValue = x[n:m]
					//case []interface{}:
					//	x := TkL.CurValue.([]interface{})
					//	TkL.DataType = CtxType_ArrayOf
					//	TkL.CurValue = x[n:m]
				}
			}

			// Array . ID/Name
			//       ^----------------------------------- Pos
		} else if TkL.DataType == CtxType_MapOf {

			name := ""

			if opTk != gen.Tok_DOT { // if it is an Map but not a '.' to ref it, then return
				return
			}

			Adv()
			// Map . ID/Name
			//       ^----------------------------------- Pos

			if opTk == gen.Tok_ID { // xyzzy - ID : -- handle this also -> Name
				// fmt.Printf("IN MapRef -ID-: %s\n", dbgo.LF())
				name = eval.Mm[eval.Pos].Match
				Adv()
				// Map . ID/Name
				//               ^----------------------------------- Pos
			} else if opTk == gen.Tok_Str0 { // xyzzy - "name" :  -- Handle this
				// fmt.Printf("IN MapRef -STR-: %s\n", dbgo.LF())
				name = eval.Mm[eval.Pos].Match
				Adv()
				// Map . ID/Name
				//               ^----------------------------------- Pos
			} else {
				eval.SetErrorInfo(&TkL, "Error (Eval00018):  Attempted to reference into map with non-string type\n")
				return
			}

			// fmt.Printf("name=%s CurValue=%s\n", name, com.SVar(TkL.CurValue))
			switch TkL.CurValue.(type) {
			case map[string]tok.Token:
				x := TkL.CurValue.(map[string]tok.Token)
				TkL.DataType = x[name].DataType
				TkL.CurValue = x[name].CurValue
				//case map[string]interface{}:
				//	x := TkL.CurValue.(map[string]interface{})
				//	// TkL.DataType = x[name].DataType
				//	TkL.CurValue = x[name]
			}

		}
	}
	return
}

// -------------------------------------------------------------------------------------------------------------------------------------------------------
// Unary !
func (eval *EvalType) Pres2() (TkL tok.Token) {
	// fmt.Printf("Pos:%d At: %s\n", eval.Pos, dbgo.LF())
	neg := false
	found := false
	if eval.Pos < len(eval.Mm) {
		for (eval.Pos < len(eval.Mm)) && (eval.Mm[eval.Pos].TokNo == gen.Tok_EXCLAM) {
			found = true
			neg = !neg
			eval.Pos++
		}
		TkL = eval.Pres1()
		if TkL.Error {
			return
		}
		if found {
			// if type is string && loose -> "" is false
			// if type is array && loose -> [] is false
			// if type is map/hash && loose -> {} is false
			// if type is float && loose && thr -> |x|<Thr is false
			// if type is float && loose && !thr -> x == 0 is false
			switch TkL.CurValue.(type) {
			case int:
				ii := TkL.CurValue.(int)
				TkL.CurValue = (ii == 0)
				// fmt.Printf("IN NOT(int): neg=%v\n", neg)
				if neg {
					TkL.CurValue = !TkL.CurValue.(bool)
				}
				// fmt.Printf("IN NOT(int), result: %v\n", TkL.CurValue.(bool))
			case bool:
				// fmt.Printf("IN NOT: neg=%v\n", neg)
				if neg {
					TkL.CurValue = !TkL.CurValue.(bool)
					// fmt.Printf("IN NOT, result: %v\n", TkL.CurValue.(bool))
				}
			default:
				TkL.CurValue = false
				TkL.DataType = CtxType_Bool
				eval.SetErrorInfo(&TkL, "Error (Eval00019):  Attempted to use '!' operator on invalid type\n")
				return
			}
		}
	}
	return
}

// -------------------------------------------------------------------------------------------------------------------------------------------------------
// Unary +, -
func (eval *EvalType) Pres3() (TkL tok.Token) {
	// fmt.Printf("Unary-+: Pos:%d At: %s\n", eval.Pos, dbgo.LF())
	found := false
	neg := false
	if eval.Pos < len(eval.Mm) {
		for (eval.Pos < len(eval.Mm)) && (eval.Mm[eval.Pos].TokNo == gen.Tok_MINUS || eval.Mm[eval.Pos].TokNo == gen.Tok_PLUS) {
			// fmt.Printf(" Top Of Loop TokNo=%d, Pos:%d\n", eval.Mm[eval.Pos].TokNo, eval.Pos)
			found = true
			if eval.Mm[eval.Pos].TokNo == gen.Tok_MINUS {
				neg = !neg
			}
			eval.Pos++
		}
		TkL = eval.Pres2()
		if found {
			// fmt.Printf("Unary-+: Found is true - :%d At: %s\n", eval.Pos, dbgo.LF())
			switch TkL.DataType {
			case CtxType_Int:
				if neg {
					TkL.CurValue = -TkL.CurValue.(int)
				}
			case CtxType_Float:
				if neg {
					TkL.CurValue = -TkL.CurValue.(float64)
				}
			default:
				// fmt.Printf("Unary-+: error type - :%d At: %s\n", eval.Pos, dbgo.LF())
				//TkL.CurValue = 0
				//TkL.DataType = CtxType_Int
				eval.SetErrorInfo(&TkL, "Error (Eval00020):  Attempted to use '+', '-' operator on invalid type\n")
				return
			}
		}
	}
	return
}

// -------------------------------------------------------------------------------------------------------------------------------------------------------
// Parse multiplication and division, *, /, %
func (eval *EvalType) Pres4() (TkL tok.Token) {
	// fmt.Printf("Pos:%d At: %s\n", eval.Pos, dbgo.LF())
	TkL = eval.Pres3()
	if TkL.Error {
		return
	}
	for eval.Pos < len(eval.Mm) {
		opPos := eval.Pos
		opTk := eval.Mm[opPos].TokNo
		if opTk != gen.Tok_STAR && opTk != gen.Tok_SLASH && opTk != gen.Tok_PCT {
			return TkL
		}
		eval.Pos++
		TkR := eval.Pres3()
		if TkR.Error {
			TkL = TkR
			return
		}
		if TkL.DataType == CtxType_Int && TkR.DataType == CtxType_Int {
			switch opTk {
			case gen.Tok_STAR:
				// fmt.Printf("IN MUL: %d * %d = %d\n", TkL.CurValue.(int), TkR.CurValue.(int), TkL.CurValue.(int)*TkR.CurValue.(int))
				TkL.CurValue = TkL.CurValue.(int) * TkR.CurValue.(int)
			case gen.Tok_SLASH:
				// fmt.Printf("IN DIV: %d / %d = %d\n", TkL.CurValue.(int), TkR.CurValue.(int), TkL.CurValue.(int)/TkR.CurValue.(int))
				TkL.CurValue = TkL.CurValue.(int) / TkR.CurValue.(int) // xyzzy zero devide
			case gen.Tok_PCT:
				// fmt.Printf("IN MOD: %d %% %d = %d\n", TkL.CurValue.(int), TkR.CurValue.(int), TkL.CurValue.(int)%TkR.CurValue.(int))
				TkL.CurValue = TkL.CurValue.(int) % TkR.CurValue.(int)
			}
		} else if (TkL.DataType == CtxType_Float && TkR.DataType == CtxType_Int) || (TkL.DataType == CtxType_Int && TkR.DataType == CtxType_Float) || (TkL.DataType == CtxType_Float && TkR.DataType == CtxType_Float) {

			if TkL.DataType == CtxType_Int {
				TkL.DataType = CtxType_Float
				TkL.CurValue = float64(TkL.CurValue.(int))
			}

			if TkR.DataType == CtxType_Int {
				TkR.DataType = CtxType_Float
				TkR.CurValue = float64(TkR.CurValue.(int))
			}

			switch opTk {
			case gen.Tok_STAR:
				TkL.CurValue = TkL.CurValue.(float64) * TkR.CurValue.(float64)
				TkL.DataType = CtxType_Float
			case gen.Tok_SLASH:
				TkL.CurValue = TkL.CurValue.(float64) / TkR.CurValue.(float64)
				TkL.DataType = CtxType_Float
			case gen.Tok_PCT:
				//TkL.CurValue = 0
				//TkL.DataType = CtxType_Int
				eval.SetErrorInfo(&TkL, "Error (Eval00021):  '%%' (modulous devision) not defined on float data.\n")
			}
		} else {
			//TkL.CurValue = 0
			//TkL.DataType = CtxType_Int
			eval.SetErrorInfo(&TkL, "Error (Eval00022):  Attempted to use '/', '%%', '*' operator on invalid or mixed data types\n")
		}
	}
	return
}

// -------------------------------------------------------------------------------------------------------------------------------------------------------
// Parse addition and subtraction, -, + (binary)
func (eval *EvalType) Pres5() (TkL tok.Token) {
	// fmt.Printf("At: %s\n", dbgo.LF())
	TkL = eval.Pres4()
	if TkL.Error {
		return
	}
	for eval.Pos < len(eval.Mm) {
		opPos := eval.Pos
		opTk := eval.Mm[opPos].TokNo
		if opTk != gen.Tok_PLUS && opTk != gen.Tok_MINUS {
			return TkL
		}
		eval.Pos++
		TkR := eval.Pres4()
		if TkR.Error {
			TkL = TkR
			return
		}
		TkL = eval.TypeConvertToNumber(TkL)
		TkR = eval.TypeConvertToNumber(TkR)
		if TkL.DataType == CtxType_Int && TkR.DataType == CtxType_Int {
			switch opTk {
			case gen.Tok_PLUS:
				// fmt.Printf("IN ADD: %d + %d = %d\n", TkL.CurValue.(int), TkR.CurValue.(int), TkL.CurValue.(int)+TkR.CurValue.(int))
				TkL.CurValue = TkL.CurValue.(int) + TkR.CurValue.(int)
			case gen.Tok_MINUS:
				// fmt.Printf("IN SUB: %d - %d = %d\n", TkL.CurValue.(int), TkR.CurValue.(int), TkL.CurValue.(int)-TkR.CurValue.(int))
				TkL.CurValue = TkL.CurValue.(int) - TkR.CurValue.(int)
			}
		} else if TkL.DataType == CtxType_Str && TkR.DataType == CtxType_Str {
			switch opTk {
			case gen.Tok_PLUS:
				TkL.CurValue = TkL.CurValue.(string) + TkR.CurValue.(string)
			case gen.Tok_MINUS:
				//TkL.CurValue = 0
				//TkL.DataType = CtxType_Int
				eval.SetErrorInfo(&TkL, "Error (Eval00023):  Attempted to use '-', operator on two strings - not defined.\n")
				return
			}
		} else if (TkL.DataType == CtxType_Float && TkR.DataType == CtxType_Int) || (TkL.DataType == CtxType_Int && TkR.DataType == CtxType_Float) || (TkL.DataType == CtxType_Float && TkR.DataType == CtxType_Float) {

			if TkL.DataType == CtxType_Int {
				TkL.DataType = CtxType_Float
				TkL.CurValue = float64(TkL.CurValue.(int))
			}

			if TkR.DataType == CtxType_Int {
				TkR.DataType = CtxType_Float
				TkR.CurValue = float64(TkR.CurValue.(int))
			}

			switch opTk {
			case gen.Tok_PLUS:
				TkL.CurValue = TkL.CurValue.(float64) + TkR.CurValue.(float64)
				TkL.DataType = CtxType_Float
			case gen.Tok_MINUS:
				TkL.CurValue = TkL.CurValue.(float64) - TkR.CurValue.(float64)
				TkL.DataType = CtxType_Float
			}

			// xyzzyPMHash - operators +, - on Arrays, Hashes

		} else {
			//TkL.CurValue = 0
			//TkL.DataType = CtxType_Int
			eval.SetErrorInfo(&TkL, "Error (Eval00024):  Attempted to use '+', '-' operator on invalid or mixed data types, %d, %d\n", TkL.DataType, TkR.DataType)
		}
	}
	return
}

// -------------------------------------------------------------------------------------------------------------------------------------------------------
// Parse addition and subtraction, <<, >> (binary)
func (eval *EvalType) Pres5a() (TkL tok.Token) {
	// fmt.Printf("At: %s\n", dbgo.LF())
	TkL = eval.Pres5()
	if TkL.Error {
		return
	}
	for eval.Pos < len(eval.Mm) {
		opPos := eval.Pos
		opTk := eval.Mm[opPos].TokNo
		if opTk != gen.Tok_S_L && opTk != gen.Tok_S_R {
			return TkL
		}
		eval.Pos++
		TkR := eval.Pres5()
		if TkR.Error {
			TkL = TkR
			return
		}
		TkL = eval.TypeConvertToNumber(TkL)
		TkR = eval.TypeConvertToNumber(TkR)
		if TkL.DataType == CtxType_Int && TkR.DataType == CtxType_Int {
			switch opTk {
			case gen.Tok_S_L:
				// fmt.Printf("IN S-L: %d << %d = %d\n", TkL.CurValue.(int), TkR.CurValue.(int), TkL.CurValue.(int)<<uint(TkR.CurValue.(int)))
				TkL.CurValue = TkL.CurValue.(int) << uint(TkR.CurValue.(int))
			case gen.Tok_S_R:
				// fmt.Printf("IN S-R: %d >> %d = %d\n", TkL.CurValue.(int), TkR.CurValue.(int), TkL.CurValue.(int)>>uint(TkR.CurValue.(int)))
				TkL.CurValue = TkL.CurValue.(int) >> uint(TkR.CurValue.(int))
			}
		} else {
			//TkL.CurValue = 0
			//TkL.DataType = CtxType_Int
			eval.SetErrorInfo(&TkL, "Error (Eval00025):  Attempted to use '<<', '>>' operator on invalid or mixed data types\n")
		}
	}
	return
}

// -------------------------------------------------------------------------------------------------------------------------------------------------------
func (eval *EvalType) TypeConvertToNumber(TkL tok.Token) (rv tok.Token) {
	rv = TkL
	fmt.Printf("Just before got a %T for type\n", TkL.CurValue)
	switch TkL.CurValue.(type) {
	//case byte:
	//	TkL.DataType = CtxType_Int
	//	TkL.CurValue = int(TkL.CurValue.(byte))
	case int64:
		rv.DataType = CtxType_Int
		rv.CurValue = int(TkL.CurValue.(int64))
	case int32:
		rv.DataType = CtxType_Int
		rv.CurValue = int(TkL.CurValue.(int32))
	case float32:
		rv.DataType = CtxType_Float
		rv.CurValue = float64(TkL.CurValue.(float32))
	}
	return
}

// -------------------------------------------------------------------------------------------------------------------------------------------------------
// Parse Compare OPS, <, <=, >, >=
func (eval *EvalType) Pres6() (TkL tok.Token) {
	// fmt.Printf("At: %s\n", dbgo.LF())
	TkL = eval.Pres5a()
	if TkL.Error {
		return
	}
	for eval.Pos < len(eval.Mm) {
		opPos := eval.Pos
		opTk := eval.Mm[opPos].TokNo
		if opTk != gen.Tok_LE && opTk != gen.Tok_LT && opTk != gen.Tok_GE && opTk != gen.Tok_GT {
			return TkL
		}
		eval.Pos++
		TkR := eval.Pres5a()
		if TkR.Error {
			TkL = TkR
			return
		}

		TkL = eval.TypeConvertToNumber(TkL)
		TkR = eval.TypeConvertToNumber(TkR)

		if TkL.DataType == CtxType_Int && TkR.DataType == CtxType_Int {
			switch opTk {
			case gen.Tok_LT:
				TkL.CurValue = TkL.CurValue.(int) < TkR.CurValue.(int)
			case gen.Tok_LE:
				TkL.CurValue = TkL.CurValue.(int) <= TkR.CurValue.(int)
			case gen.Tok_GT:
				TkL.CurValue = TkL.CurValue.(int) > TkR.CurValue.(int)
			case gen.Tok_GE:
				TkL.CurValue = TkL.CurValue.(int) >= TkR.CurValue.(int)
			}
			TkL.DataType = CtxType_Bool
		} else if (TkL.DataType == CtxType_Float && TkR.DataType == CtxType_Int) || (TkL.DataType == CtxType_Int && TkR.DataType == CtxType_Float) || (TkL.DataType == CtxType_Float && TkR.DataType == CtxType_Float) {

			if TkL.DataType == CtxType_Int {
				TkL.DataType = CtxType_Float
				TkL.CurValue = float64(TkL.CurValue.(int))
			}

			if TkR.DataType == CtxType_Int {
				TkR.DataType = CtxType_Float
				TkR.CurValue = float64(TkR.CurValue.(int))
			}

			switch opTk {
			case gen.Tok_LT:
				TkL.CurValue = TkL.CurValue.(float64) < TkR.CurValue.(float64)
			case gen.Tok_LE:
				TkL.CurValue = TkL.CurValue.(float64) <= TkR.CurValue.(float64)
			case gen.Tok_GT:
				TkL.CurValue = TkL.CurValue.(float64) > TkR.CurValue.(float64)
			case gen.Tok_GE:
				TkL.CurValue = TkL.CurValue.(float64) >= TkR.CurValue.(float64)
			}

			TkL.DataType = CtxType_Bool
		} else if TkL.DataType == CtxType_Str && TkR.DataType == CtxType_Str {
			switch opTk {
			case gen.Tok_LT:
				TkL.CurValue = TkL.CurValue.(string) < TkR.CurValue.(string)
			case gen.Tok_LE:
				TkL.CurValue = TkL.CurValue.(string) <= TkR.CurValue.(string)
			case gen.Tok_GT:
				TkL.CurValue = TkL.CurValue.(string) > TkR.CurValue.(string)
			case gen.Tok_GE:
				TkL.CurValue = TkL.CurValue.(string) >= TkR.CurValue.(string)
			}
			TkL.DataType = CtxType_Bool
		} else {
			//TkL.CurValue = 0
			//TkL.DataType = CtxType_Int
			eval.SetErrorInfo(&TkL, "Error (Eval00026):  Attempted to use '!=', '<>', '==' operator on invalid or mixed data types\n")
		}
	}
	return
}

// -------------------------------------------------------------------------------------------------------------------------------------------------------
// !=, ==, <>
func (eval *EvalType) Pres7() (TkL tok.Token) {
	// fmt.Printf("At: %s\n", dbgo.LF())
	TkL = eval.Pres6()
	if TkL.Error {
		return
	}
	for eval.Pos < len(eval.Mm) {
		opPos := eval.Pos
		opTk := eval.Mm[opPos].TokNo
		if opTk != gen.Tok_NE && opTk != gen.Tok_L_EQ {
			return TkL
		}
		eval.Pos++
		TkR := eval.Pres6()
		if TkR.Error {
			TkL = TkR
			return
		}
		TkL = eval.TypeConvertToNumber(TkL)
		TkR = eval.TypeConvertToNumber(TkR)
		if TkL.DataType == CtxType_Int && TkR.DataType == CtxType_Int {
			switch opTk {
			case gen.Tok_NE:
				TkL.CurValue = TkL.CurValue.(int) != TkR.CurValue.(int)
			case gen.Tok_L_EQ:
				TkL.CurValue = TkL.CurValue.(int) == TkR.CurValue.(int)
			}
			TkL.DataType = CtxType_Bool
		} else if (TkL.DataType == CtxType_Float && TkR.DataType == CtxType_Int) || (TkL.DataType == CtxType_Int && TkR.DataType == CtxType_Float) || (TkL.DataType == CtxType_Float && TkR.DataType == CtxType_Float) {
			// fmt.Printf("At: %s\n", dbgo.LF())

			if TkL.DataType == CtxType_Int {
				// fmt.Printf("At: %s\n", dbgo.LF())
				TkL.DataType = CtxType_Float
				TkL.CurValue = float64(TkL.CurValue.(int))
			}

			if TkR.DataType == CtxType_Int {
				// fmt.Printf("At: %s\n", dbgo.LF())
				TkR.DataType = CtxType_Float
				TkR.CurValue = float64(TkR.CurValue.(int))
			}

			switch opTk {
			case gen.Tok_NE:
				// fmt.Printf("At: %s\n", dbgo.LF())
				TkL.CurValue = TkL.CurValue.(float64) != TkR.CurValue.(float64)
			case gen.Tok_L_EQ:
				// fmt.Printf("At: %s\n", dbgo.LF())
				TkL.CurValue = TkL.CurValue.(float64) == TkR.CurValue.(float64)
			}
			// fmt.Printf("At: %s\n", dbgo.LF())
			TkL.DataType = CtxType_Bool
		} else if TkL.DataType == CtxType_Bool && TkR.DataType == CtxType_Bool {
			switch opTk {
			case gen.Tok_NE:
				TkL.CurValue = TkL.CurValue.(bool) != TkR.CurValue.(bool)
			case gen.Tok_L_EQ:
				TkL.CurValue = TkL.CurValue.(bool) == TkR.CurValue.(bool)
			}
			TkL.DataType = CtxType_Bool
		} else if TkL.DataType == CtxType_Str && TkR.DataType == CtxType_Str {
			switch opTk {
			case gen.Tok_NE:
				TkL.CurValue = TkL.CurValue.(string) != TkR.CurValue.(string)
			case gen.Tok_L_EQ:
				TkL.CurValue = TkL.CurValue.(string) == TkR.CurValue.(string)
			}
			TkL.DataType = CtxType_Bool
		} else {
			//TkL.CurValue = 0
			//TkL.DataType = CtxType_Int
			eval.SetErrorInfo(&TkL, "Error (Eval00027):  Attempted to use '/', '%%', '*' operator on invalid or mixed data types\n")
		}
	}
	return
}

// -------------------------------------------------------------------------------------------------------------------------------------------------------
// &
func (eval *EvalType) Pres8() (TkL tok.Token) {
	// fmt.Printf("At: %s\n", dbgo.LF())
	TkL = eval.Pres7()
	if TkL.Error {
		return
	}
	for eval.Pos < len(eval.Mm) {
		opPos := eval.Pos
		opTk := eval.Mm[opPos].TokNo
		if opTk != gen.Tok_B_AND {
			return TkL
		}
		eval.Pos++
		TkR := eval.Pres7()
		if TkR.Error {
			TkL = TkR
			return
		}
		if TkL.DataType == CtxType_Int && TkR.DataType == CtxType_Int {
			switch opTk {
			case gen.Tok_B_AND:
				TkL.CurValue = TkL.CurValue.(int) & TkR.CurValue.(int)
			}
		} else {
			//TkL.CurValue = 0
			//TkL.DataType = CtxType_Int
			eval.SetErrorInfo(&TkL, "Error (Eval00028):  Attempted to use '&' (bit-and),  operator on invalid or mixed data types\n")
		}
	}
	return
}

// -------------------------------------------------------------------------------------------------------------------------------------------------------
// |		-- Change to "bor" token/ID
func (eval *EvalType) Pres9() (TkL tok.Token) {
	// fmt.Printf("At: %s\n", dbgo.LF())
	TkL = eval.Pres8()
	for eval.Pos < len(eval.Mm) {
		opPos := eval.Pos
		opTk := eval.Mm[opPos].TokNo
		if opTk != gen.Tok_B_OR {
			return TkL
		}
		eval.Pos++
		TkR := eval.Pres8()
		if TkR.Error {
			TkL = TkR
			return
		}
		if TkL.DataType == CtxType_Int && TkR.DataType == CtxType_Int {
			switch opTk {
			case gen.Tok_B_OR:
				TkL.CurValue = TkL.CurValue.(int) | TkR.CurValue.(int)
			}
		} else {
			//TkL.CurValue = 0
			//TkL.DataType = CtxType_Int
			eval.SetErrorInfo(&TkL, "Error (Eval00029):  Attempted to use 'bor',  operator on invalid or mixed data types\n")
		}
	}
	return
}

// -------------------------------------------------------------------------------------------------------------------------------------------------------
// ^
func (eval *EvalType) PresA() (TkL tok.Token) {
	// fmt.Printf("At: %s\n", dbgo.LF())
	TkL = eval.Pres9()
	if TkL.Error {
		return
	}
	for eval.Pos < len(eval.Mm) {
		opPos := eval.Pos
		opTk := eval.Mm[opPos].TokNo
		if opTk != gen.Tok_CARRET {
			return TkL
		}
		eval.Pos++
		TkR := eval.Pres9()
		if TkR.Error {
			TkL = TkR
			return
		}
		if TkL.DataType == CtxType_Int && TkR.DataType == CtxType_Int {
			switch opTk {
			case gen.Tok_CARRET:
				TkL.CurValue = TkL.CurValue.(int) ^ TkR.CurValue.(int)
			}
		} else {
			//TkL.CurValue = 0
			//TkL.DataType = CtxType_Int
			eval.SetErrorInfo(&TkL, "Error (Eval00030):  Attempted to use '^' (bit-xor),  operator on invalid or mixed data types\n")
		}
	}
	return
}

// -------------------------------------------------------------------------------------------------------------------------------------------------------
// &&
func (eval *EvalType) PresB() (TkL tok.Token) {
	// fmt.Printf("At: %s\n", dbgo.LF())
	TkL = eval.PresA()
	for eval.Pos < len(eval.Mm) {
		opPos := eval.Pos
		opTk := eval.Mm[opPos].TokNo
		if opTk != gen.Tok_L_AND {
			return TkL
		}
		eval.Pos++
		TkR := eval.PresA()
		if TkR.Error {
			TkL = TkR
			return
		}
		if TkL.DataType == CtxType_Bool && TkR.DataType == CtxType_Bool {
			switch opTk {
			case gen.Tok_L_AND:
				TkL.CurValue = TkL.CurValue.(bool) && TkR.CurValue.(bool)
			}
		} else {
			//TkL.CurValue = 0
			//TkL.DataType = CtxType_Int
			eval.SetErrorInfo(&TkL, "Error (Eval00031):  Attempted to use '&&' (logical and),  operator on invalid or mixed data types\n")
		}
	}
	return
}

// -------------------------------------------------------------------------------------------------------------------------------------------------------
// ||
func (eval *EvalType) PresC() (TkL tok.Token) {
	// fmt.Printf("At: %s\n", dbgo.LF())
	TkL = eval.PresB()
	if TkL.Error {
		return
	}
	for eval.Pos < len(eval.Mm) {
		opPos := eval.Pos
		opTk := eval.Mm[opPos].TokNo
		if opTk != gen.Tok_L_OR {
			return TkL
		}
		eval.Pos++
		TkR := eval.PresB()
		if TkR.Error {
			TkL = TkR
			return
		}
		if TkL.DataType == CtxType_Bool && TkR.DataType == CtxType_Bool {
			switch opTk {
			case gen.Tok_L_OR:
				TkL.CurValue = TkL.CurValue.(bool) || TkR.CurValue.(bool)
			}
		} else {
			//TkL.CurValue = 0
			//TkL.DataType = CtxType_Int
			eval.SetErrorInfo(&TkL, "Error (Eval00032):  Attempted to use '&&' (logical and),  operator on invalid or mixed data types\n")
		}
	}
	return
}

// -------------------------------------------------------------------------------------------------------------------------------------------------------
// Not Used Uet
func (eval *EvalType) PresD() (TkL tok.Token) {
	TkL = eval.PresC() // TkL = eval.ParseJsonExpression()
	return
}

// -------------------------------------------------------------------------------------------------------------------------------------------------------
// = (assignment), += etc. :=??
// Right To Left - how, Recursion?? array and append??
// xyzzy R->L not implemented - just one assignemtn at the moment
func (eval *EvalType) PresE() (TkL tok.Token) {
	// fmt.Printf("At: %s\n", dbgo.LF())

	if eval.Mm[eval.Pos].TokNo == gen.Tok_ID {
		eval.Mm[eval.Pos].CreateId = (eval.Pos+1 < len(eval.Mm) && eval.Mm[eval.Pos+1].TokNo == gen.Tok_DCL_VAR)
		// fmt.Printf("CreateID is set to %v, for %s\n", eval.Mm[eval.Pos].CreateId, eval.Mm[eval.Pos].Match)
	}
	TkL = eval.PresD()
	if TkL.Error {
		return
	}
	if TkL.LValue {
		// fmt.Printf("IS Lvalue TkL=%+v At: %s\n", TkL, dbgo.LF())
		if eval.Pos < len(eval.Mm) {
			opPos := eval.Pos
			opTk := eval.Mm[opPos].TokNo
			if opTk != gen.Tok_EQ && opTk != gen.Tok_DCL_VAR && opTk != gen.Tok_PLUS_EQ && opTk != gen.Tok_MINUS_EQ && opTk != gen.Tok_STAR_EQ && opTk != gen.Tok_DIV_EQ && opTk != gen.Tok_MOD_EQ && opTk != gen.Tok_CAROT_EQ && opTk != gen.Tok_B_OR_EQ && opTk != gen.Tok_B_AND_EQ && opTk != gen.Tok_TILDE_EQ && opTk != gen.Tok_S_L_EQ && opTk != gen.Tok_S_R_EQ {
				return TkL
			}
			eval.Pos++
			TkR := eval.PresD()
			if TkR.Error {
				TkL = TkR
				return
			}
			isInt := true
			isAssign := false
			// xyzzy - hand check that all types have been placed in each and every switch
			switch TkL.CurValue.(type) {
			case int:
				switch TkR.CurValue.(type) {
				case int:
				case float64:
					isInt = false
					TkL.DataType = CtxType_Float
					TkL.CurValue = float64(TkL.CurValue.(int))
				case []interface{}, []tok.Token, map[string]interface{}, map[string]tok.Token, tok.Token, string, bool:
					if opTk == gen.Tok_EQ || opTk == gen.Tok_DCL_VAR {
						isAssign = true
					} else {
						isInt = false
						//TkL.CurValue = 0
						//TkL.DataType = CtxType_Int
						eval.SetErrorInfo(&TkL, "Error (Eval00033):  Invalid type conversion. Type=%s\n", fmt.Sprintf("%T = %T", TkL.CurValue, TkR.CurValue))
						return
					}
				}
			case float64:
				isInt = false
				switch TkR.CurValue.(type) {
				case int:
					TkR.DataType = CtxType_Float
					TkR.CurValue = float64(TkR.CurValue.(int))
				case float64:
				case []interface{}, []tok.Token, map[string]interface{}, map[string]tok.Token, tok.Token, string, bool:
					if opTk == gen.Tok_EQ || opTk == gen.Tok_DCL_VAR {
						isAssign = true
					} else {
						isInt = false
						//TkL.CurValue = 0
						//TkL.DataType = CtxType_Int
						eval.SetErrorInfo(&TkL, "Error (Eval00034):  Invalid type conversion. Type=%s\n", fmt.Sprintf("%T = %T", TkL.CurValue, TkR.CurValue))
						return
					}
				}
			case string:
				switch TkR.CurValue.(type) {
				case []interface{}, []tok.Token, map[string]interface{}, map[string]tok.Token, tok.Token, string, bool, float64, int:
					if opTk == gen.Tok_EQ || opTk == gen.Tok_DCL_VAR {
						isAssign = true
					} else {
						isInt = false
						//TkL.CurValue = 0
						//TkL.DataType = CtxType_Int
						eval.SetErrorInfo(&TkL, "Error (Eval00035):  Invalid type conversion. Type=%s\n", fmt.Sprintf("%T = %T", TkL.CurValue, TkR.CurValue))
						return
					}
				}
			case bool:
				switch TkR.CurValue.(type) {
				case []interface{}, []tok.Token, map[string]interface{}, map[string]tok.Token, tok.Token, string, bool, float64, int:
					if opTk == gen.Tok_EQ || opTk == gen.Tok_DCL_VAR {
						isAssign = true
					} else {
						isInt = false
						//TkL.CurValue = 0
						//TkL.DataType = CtxType_Int
						eval.SetErrorInfo(&TkL, "Error (Eval00036):  Invalid type conversion. Type=%s\n", fmt.Sprintf("%T = %T", TkL.CurValue, TkR.CurValue))
						return
					}
				}
			case []interface{}, []tok.Token, map[string]interface{}, map[string]tok.Token, tok.Token:
				switch TkR.CurValue.(type) {
				case []interface{}, []tok.Token, map[string]interface{}, map[string]tok.Token, tok.Token, string, bool, float64, int:
					if opTk == gen.Tok_EQ || opTk == gen.Tok_DCL_VAR {
						isAssign = true
					} else {
						isInt = false
						//TkL.CurValue = 0
						//TkL.DataType = CtxType_Int
						eval.SetErrorInfo(&TkL, "Error (Eval00037):  Invalid type conversion. Type=%s\n", fmt.Sprintf("%T = %T", TkL.CurValue, TkR.CurValue))
						return
					}
				}
			}
			if isAssign {
				switch opTk {
				case gen.Tok_EQ:
					// fmt.Printf("IS = Setting [[[%s]]] to %v At: %s\n", TkL.Match, TkR.CurValue, dbgo.LF())
					TkL.CurValue = TkR.CurValue
				case gen.Tok_DCL_VAR:
					// fmt.Printf("IS = Setting [[[%s]]] to %v At: %s\n", TkL.Match, TkR.CurValue, dbgo.LF())
					TkL.CurValue = TkR.CurValue
				}
			} else if isInt {
				switch opTk {
				case gen.Tok_EQ:
					// fmt.Printf("IS = Setting [[[%s]]] to %v At: %s\n", TkL.Match, TkR.CurValue, dbgo.LF())
					TkL.CurValue = TkR.CurValue
				case gen.Tok_DCL_VAR:
					// fmt.Printf("IS = Setting [[[%s]]] to %v At: %s\n", TkL.Match, TkR.CurValue, dbgo.LF())
					TkL.CurValue = TkR.CurValue
				case gen.Tok_PLUS_EQ:
					// fmt.Printf("IS = Setting [[[%s]]] to %v At: %s\n", TkL.Match, TkR.CurValue, dbgo.LF())
					TkL.CurValue = TkL.CurValue.(int) + TkR.CurValue.(int)
				case gen.Tok_MINUS_EQ:
					// fmt.Printf("IS = Setting [[[%s]]] to %v At: %s\n", TkL.Match, TkR.CurValue, dbgo.LF())
					TkL.CurValue = TkL.CurValue.(int) - TkR.CurValue.(int)
				case gen.Tok_STAR_EQ:
					// fmt.Printf("IS = Setting [[[%s]]] to %v At: %s\n", TkL.Match, TkR.CurValue, dbgo.LF())
					TkL.CurValue = TkL.CurValue.(int) * TkR.CurValue.(int)
				case gen.Tok_DIV_EQ:
					// fmt.Printf("IS = Setting [[[%s]]] to %v At: %s\n", TkL.Match, TkR.CurValue, dbgo.LF())
					TkL.CurValue = TkL.CurValue.(int) / TkR.CurValue.(int)
				case gen.Tok_MOD_EQ:
					// fmt.Printf("IS = Setting [[[%s]]] to %v At: %s\n", TkL.Match, TkR.CurValue, dbgo.LF())
					TkL.CurValue = TkL.CurValue.(int) % TkR.CurValue.(int)
				case gen.Tok_CAROT_EQ:
					// fmt.Printf("IS = Setting [[[%s]]] to %v At: %s\n", TkL.Match, TkR.CurValue, dbgo.LF())
					TkL.CurValue = TkL.CurValue.(int) ^ TkR.CurValue.(int)
				case gen.Tok_B_OR_EQ:
					// fmt.Printf("IS = Setting [[[%s]]] to %v At: %s\n", TkL.Match, TkR.CurValue, dbgo.LF())
					TkL.CurValue = TkL.CurValue.(int) | TkR.CurValue.(int)
				case gen.Tok_B_AND_EQ:
					// fmt.Printf("IS = Setting [[[%s]]] to %v At: %s\n", TkL.Match, TkR.CurValue, dbgo.LF())
					TkL.CurValue = TkL.CurValue.(int) & TkR.CurValue.(int)
				case gen.Tok_TILDE_EQ:
					//TkL.CurValue = 0
					//TkL.DataType = CtxType_Int
					eval.SetErrorInfo(&TkL, "Error (Eval00038):  ~= operator not implemented yet.\n")
					return
				case gen.Tok_S_L_EQ:
					// fmt.Printf("IS = Setting [[[%s]]] to %v At: %s\n", TkL.Match, TkR.CurValue, dbgo.LF())
					TkL.CurValue = TkL.CurValue.(int) << uint(TkR.CurValue.(int))
				case gen.Tok_S_R_EQ:
					// fmt.Printf("IS = Setting [[[%s]]] to %v At: %s\n", TkL.Match, TkR.CurValue, dbgo.LF())
					TkL.CurValue = TkL.CurValue.(int) >> uint(TkR.CurValue.(int))
				}
			} else {
				TkR.DataType = CtxType_Float
				switch opTk {
				case gen.Tok_EQ:
					// fmt.Printf("IS = Setting [[[%s]]] to %v At: %s\n", TkL.Match, TkR.CurValue, dbgo.LF())
					TkL.CurValue = TkR.CurValue
				case gen.Tok_DCL_VAR:
					// fmt.Printf("IS = Setting [[[%s]]] to %v At: %s\n", TkL.Match, TkR.CurValue, dbgo.LF())
					TkL.CurValue = TkR.CurValue
				case gen.Tok_PLUS_EQ:
					// fmt.Printf("IS = Setting [[[%s]]] to %v At: %s\n", TkL.Match, TkR.CurValue, dbgo.LF())
					TkL.CurValue = TkL.CurValue.(float64) + TkR.CurValue.(float64)
				case gen.Tok_MINUS_EQ:
					// fmt.Printf("IS = Setting [[[%s]]] to %v At: %s\n", TkL.Match, TkR.CurValue, dbgo.LF())
					TkL.CurValue = TkL.CurValue.(float64) - TkR.CurValue.(float64)
				case gen.Tok_STAR_EQ:
					// fmt.Printf("IS = Setting [[[%s]]] to %v At: %s\n", TkL.Match, TkR.CurValue, dbgo.LF())
					TkL.CurValue = TkL.CurValue.(float64) * TkR.CurValue.(float64)
				case gen.Tok_DIV_EQ:
					// fmt.Printf("IS = Setting [[[%s]]] to %v At: %s\n", TkL.Match, TkR.CurValue, dbgo.LF())
					TkL.CurValue = TkL.CurValue.(float64) / TkR.CurValue.(float64)
				case gen.Tok_MOD_EQ:
					// fmt.Printf("IS = Setting [[[%s]]] to %v At: %s\n", TkL.Match, TkR.CurValue, dbgo.LF())
					//TkL.CurValue = 0
					//TkL.DataType = CtxType_Int
					eval.SetErrorInfo(&TkL, "Error (Eval00039):  %= operator not defined for floating point data.\n")
					return
				case gen.Tok_CAROT_EQ:
					// fmt.Printf("IS = Setting [[[%s]]] to %v At: %s\n", TkL.Match, TkR.CurValue, dbgo.LF())
					//TkL.CurValue = 0
					//TkL.DataType = CtxType_Int
					eval.SetErrorInfo(&TkL, "Error (Eval00040):  ^= operator not defined for floating point data.\n")
					return
				case gen.Tok_B_OR_EQ:
					// fmt.Printf("IS = Setting [[[%s]]] to %v At: %s\n", TkL.Match, TkR.CurValue, dbgo.LF())
					//TkL.CurValue = 0
					//TkL.DataType = CtxType_Int
					eval.SetErrorInfo(&TkL, "Error (Eval00041):  |= operator not defined for floating point data.\n")
					return
				case gen.Tok_B_AND_EQ:
					// fmt.Printf("IS = Setting [[[%s]]] to %v At: %s\n", TkL.Match, TkR.CurValue, dbgo.LF())
					//TkL.CurValue = 0
					//TkL.DataType = CtxType_Int
					eval.SetErrorInfo(&TkL, "Error (Eval00042):  &= operator not defined for floating point data.\n")
					return
				case gen.Tok_TILDE_EQ:
					// fmt.Printf("IS = Setting [[[%s]]] to %v At: %s\n", TkL.Match, TkR.CurValue, dbgo.LF())
					//TkL.DataType = TkR.DataType
					//TkL.CurValue = TkL.CurValue.(float64) ~ TkR.CurValue.(float64)
					//eval.Ctx.SetInContext(TkL.Match, TkL.DataType, TkL.CurValue)
					// xyzzy - error
					//TkL.CurValue = 0
					//TkL.DataType = CtxType_Int
					eval.SetErrorInfo(&TkL, "Error (Eval00043):  ~= operator not implemented yet.\n")
					return
				case gen.Tok_S_L_EQ:
					//TkL.CurValue = 0
					//TkL.DataType = CtxType_Int
					eval.SetErrorInfo(&TkL, "Error (Eval00044):  <<= operator not defined for floating point data.\n")
					return
				case gen.Tok_S_R_EQ:
					//TkL.CurValue = 0
					//TkL.DataType = CtxType_Int
					eval.SetErrorInfo(&TkL, "Error (Eval00045):  >>= operator not defined for floating point data.\n")
					return
				}
			}
			TkL.DataType = TkR.DataType
			eval.Ctx.SetInContext(TkL.Match, TkL.DataType, TkL.CurValue)
		}
	}
	return
}

// -------------------------------------------------------------------------------------------------------------------------------------------------------
// Computational IF ?:
func (eval *EvalType) PresF() (TkL tok.Token) {
	// fmt.Printf("At: %s\n", dbgo.LF())
	TkL = eval.PresE()
	if TkL.Error {
		return
	}

	opPos := 0
	opTk := 0
	_, _ = opPos, opTk

	Adv := func() {
		if eval.Pos < len(eval.Mm) {
			eval.Pos++
			opPos = eval.Pos
		}
		opTk = gen.Tok_null
		if eval.Pos < len(eval.Mm) {
			opTk = eval.Mm[eval.Pos].TokNo
		}
	}
	Set := func() {
		if eval.Pos < len(eval.Mm) {
			opPos = eval.Pos
		}
		opTk = gen.Tok_null
		if eval.Pos < len(eval.Mm) {
			opTk = eval.Mm[eval.Pos].TokNo
		}
	}

	// fmt.Printf("Before Loop At: %s\n", dbgo.LF())
	for eval.Pos < len(eval.Mm) {
		// fmt.Printf("Loop Top At: %s\n", dbgo.LF())
		Set()

		if opTk != gen.Tok_QUEST {
			// fmt.Printf("At: %s\n", dbgo.LF())
			return
		} else {
			Adv() // Move over '?'
			// fmt.Printf("At: %s\n", dbgo.LF())
			if TkL.DataType == CtxType_Bool && TkL.CurValue.(bool) { // Eval True Part
				// fmt.Printf("At: %s\n", dbgo.LF())
				TkL = eval.PresE()
				if TkL.Error {
					return
				}
				Set()
				if opTk == gen.Tok_COLON {
					Adv()
					_ = eval.PresE() // xyzzy - error - if error then what? -- Must take no action! No Side Effects
					Set()
				} else {
					// fmt.Printf("At: %s\n", dbgo.LF())
					//TkL.CurValue = 0
					//TkL.DataType = CtxType_Int
					eval.SetErrorInfo(&TkL, "Error (Eval00046):  Missing ':' in conditional if.\n")
					return
				}
			} else {
				// fmt.Printf("At: %s\n", dbgo.LF())
				_ = eval.PresE() // xyzzy - error - if error then what? -- Must take no action! No Side Effects
				Set()
				if opTk == gen.Tok_COLON {
					Adv()
					// fmt.Printf("At: %s\n", dbgo.LF())
					TkL = eval.PresE()
					if TkL.Error {
						return
					}
					Set()
				} else {
					// fmt.Printf("At: %s\n", dbgo.LF())
					//TkL.CurValue = 0
					//TkL.DataType = CtxType_Int
					eval.SetErrorInfo(&TkL, "Error (Eval00047):  Missing ':' in conditional if.\n")
					return
				}
			}
		}
	}
	return
}

// -------------------------------------------------------------------------------------------------------------------------------------------------------
// This is called form PresG (pipe) during error conditions.
func (eval *EvalType) ListTokens() string {
	// fmt.Printf("At: %s\n", dbgo.LF())
	s := "Tokens are ("
	com := ""
	for j := eval.Pos; j < len(eval.Mm); j++ {
		// s += fmt.Sprintf("%s [%d] = %d/%s", com, j, eval.Mm[j].TokNo, eval.Mm[j].TokNo)
		s += fmt.Sprintf("%s [%d] = %d", com, j, eval.Mm[j].TokNo)
	}
	s += " )"
	return s
}

// -------------------------------------------------------------------------------------------------------------------------------------------------------
//
// expr | function-name ( values, ... ) | function-name ( values, ... )
//
// x | Fx | Fy ( w, v ) | Fy ( m, n )
//
// a, b := Fx ( x )
// c, err := Fy ( a, b, w, v )
// d := Fy ( c, m, n )
//
// -------------------------------------------------------------------------------------------------------------------------------------------------------
func (eval *EvalType) PresG() (TkL tok.Token) {
	// fmt.Printf("At: %s\n", dbgo.LF())
	TkL = eval.PresF() // Eval 1st expression in possible pipe
	if TkL.Error {
		return
	}

	opPos := 0
	opTk := 0
	_, _ = opPos, opTk

	Adv := func() {
		if eval.Pos < len(eval.Mm) {
			eval.Pos++
			opPos = eval.Pos
		}
		opTk = gen.Tok_null
		if eval.Pos < len(eval.Mm) {
			opTk = eval.Mm[eval.Pos].TokNo
		}
	}
	Set := func() {
		if eval.Pos < len(eval.Mm) {
			opPos = eval.Pos
		}
		opTk = gen.Tok_null
		if eval.Pos < len(eval.Mm) {
			opTk = eval.Mm[eval.Pos].TokNo
		}
	}

	if eval.Pos < len(eval.Mm) {
		// fmt.Printf("At: %s\n", dbgo.LF())
		if eval.Mm[eval.Pos].TokNo == gen.Tok_PIPE {
			Tk0 := TkL
			// fmt.Printf("At: %s\n", dbgo.LF())
			Set()
			for eval.Pos < len(eval.Mm) && eval.Mm[eval.Pos].TokNo == gen.Tok_PIPE {
				Adv()
				// fmt.Printf("At: %s\n", dbgo.LF())
				if eval.Pos < len(eval.Mm) {
					match := eval.Mm[eval.Pos].Match
					// fmt.Printf("At: %s\n", dbgo.LF())
					if eval.Mm[eval.Pos].TokNo == gen.Tok_ID {
						// fmt.Printf("At: %s\n", dbgo.LF())
						if /*val0*/ _, t, f := eval.Ctx.GetFromContext(eval.Mm[eval.Pos].Match); f {
							// fmt.Printf("At: %s\n", dbgo.LF())
							// fmt.Printf("   Found, Type=%d/%s, %s\n", t, eval.Ctx.NameOfType(t), dbgo.LF())
							switch t {
							case CtxType_Func:
								// fmt.Printf("At: %s\n", dbgo.LF())
								Adv()
								if eval.Pos < len(eval.Mm) && eval.Mm[eval.Pos].TokNo == gen.Tok_OP {
									Plist := eval.ParsePlist()
									// fmt.Printf("Found a Func with a plist\n")
									P2 := make([]tok.Token, 0, len(Plist)+1) // Make the piped value the 1st parameter -- xyzzy is this how go templates work?
									P2 = append(P2, Tk0)
									P2 = append(P2, Plist...)
									Tk0 = eval.CallFunction(match, P2)
									Set()
								} else {
									// fmt.Printf("At: %s\n", dbgo.LF())
									Tk0 = eval.CallFunction(match, []tok.Token{Tk0})
									Set()
								}
							default:
								// fmt.Printf("At: %s\n", dbgo.LF())
								//TkL.CurValue = 0
								//TkL.DataType = CtxType_Int
								eval.SetErrorInfo(&TkL, "Error (Eval00048): Function refrenced %s is not a function -- can not call it\n", match)
								return
							}
						} else {
							// fmt.Printf("At: %s\n", dbgo.LF())
							//TkL.CurValue = 0
							//TkL.DataType = CtxType_Int
							eval.SetErrorInfo(&TkL, "Error (Eval00049): Function refrenced %s missing () to make call\n", match)
							return
						}
					} else {
						// fmt.Printf("At: %s\n", dbgo.LF())
						//TkL.CurValue = 0
						//TkL.DataType = CtxType_Int
						eval.SetErrorInfo(&TkL, "Error (Eval00050): Token found after pipe is not an ID, found %s\n", match)
						return
					}
				}
				// fmt.Printf("At: %s\n", dbgo.LF())
				TkL = Tk0
				if TkL.Error {
					return
				}
				// fmt.Printf("At: %s\n", dbgo.LF())
			}
		}
	}
	return
}

// -------------------------------------------------------------------------------------------------------------------------------------------------------
func (eval *EvalType) PresTop() (TkL tok.Token) {
	// fmt.Printf("At: %s\n", dbgo.LF())
	if false {
		fmt.Printf("PresTop INPUT %s\n", com.SVarI(eval))
	} else {
		fmt.Printf("PresTop INPUT %+v\n", eval)
	}
	TkL = eval.PresG()
	fmt.Printf("PresTop AFTER eval.Pos=%d len(eval.Mm)=%d\n", eval.Pos, len(eval.Mm))
	if eval.Pos < len(eval.Mm) {
		// fmt.Printf("At: %s\n", dbgo.LF())
		//TkL.CurValue = 0
		//TkL.DataType = CtxType_Int
		eval.SetErrorInfo(&TkL, "Error (Eval00051): Extra tokens found at end of expressions, %s\n", eval.ListTokens())
		return
	}
	return
}

/* vim: set noai ts=4 sw=4: */
