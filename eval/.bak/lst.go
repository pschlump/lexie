package eval

import (
	"errors"
	"fmt"
	"reflect"
	"sync"

	"../com"

	// "../com"
	// "../../../go-lib/sizlib"
)

const (
	CtxType_NULL     = 0
	CtxType_Int      = 1
	CtxType_Str      = 2
	CtxType_Bool     = 3
	CtxType_Float    = 4
	CtxType_ArrayOf  = 5
	CtxType_MapOf    = 6
	CtxType_SMapOf   = 7
	CtxType_KMapOf   = 8
	CtxType_Error    = 9
	CtxType_ID       = 10
	CtxType_Func     = 11
	CtxType_TypeCast = 12
	CtxType_Token    = 13 // Id is equvilant token (and, or, xor etc)
)

// fmt.Printf("   Found, Type=%d\n", t, eval.Ctx.NameOfType(t))
func (ctx *ContextType) NameOfType(t int) (rv string) {
	switch t {
	case CtxType_NULL:
		rv = "CtxType_NULL"
	case CtxType_Int:
		rv = "CtxType_Int"
	case CtxType_Str:
		rv = "CtxType_Str"
	case CtxType_Bool:
		rv = "CtxType_Bool"
	case CtxType_Float:
		rv = "CtxType_Float"
	case CtxType_ArrayOf:
		rv = "CtxType_ArrayOf"
	case CtxType_MapOf:
		rv = "CtxType_MapOf"
	case CtxType_SMapOf:
		rv = "CtxType_SMapOf"
	case CtxType_KMapOf:
		rv = "CtxType_KMapOf"
	case CtxType_Error:
		rv = "CtxType_Error"
	case CtxType_ID:
		rv = "CtxType_ID"
	case CtxType_Func:
		rv = "CtxType_Func"
	case CtxType_TypeCast:
		rv = "CtxType_TypeCast"
	case CtxType_Token:
		rv = "CtxType_Token"
	default:
		rv = fmt.Sprintf("CtxType_%d", t)
	}
	return
}

// -------------------------------------------------------------------------------------------------------------------------------------------------------------------------
// Set a value in the context -- Context Implementation --
// xyzzySetGet
//	1. Locking
//	2. Tests
//	3. Read Json -> Context
// 	4. Context Search
//		1. With context
//		2. Arrays
//		3. Structs/Maps of data
//	5. Iteration over a set of data
//

type ContextValueType struct {
	TypeOf int
	Val    interface{}
	Func   reflect.Value
	Prev   *ContextValueType
}

type ContextType struct {
	Store map[string]*ContextValueType
	mutex sync.RWMutex
}

var (
	ErrParamsNotAdapted = errors.New("The number of params is not adapted.")
)

func NewContextType() (rv *ContextType) {
	rv = &ContextType{
		Store: make(map[string]*ContextValueType),
	}
	rv.Store["true"] = &ContextValueType{TypeOf: CtxType_Bool, Val: true}
	rv.Store["TRUE"] = &ContextValueType{TypeOf: CtxType_Bool, Val: true}
	rv.Store["false"] = &ContextValueType{TypeOf: CtxType_Bool, Val: false}
	rv.Store["FALSE"] = &ContextValueType{TypeOf: CtxType_Bool, Val: false}

	// rv.Store["nullFunc"] = &ContextValueType{TypeOf: CtxType_Func, Val: func() {}}
	// rv.Store["nullFuncSSN"] = &ContextValueType{TypeOf: CtxType_Func, Val: func() {}}
	rv.SetInContext("nullFunc", CtxType_Func, func() {})
	rv.SetInContext("nullFuncSSN", CtxType_Func, func(s, t string, n int) {})
	rv.SetInContext("nullFuncSSB", CtxType_Func, func(s, t string, b bool) {})

	// Xyzzy - add __version__, __file__, __line__, __col_no__
	return
}

func (ctx *ContextType) SetInContext(id string, ty int, val interface{}) (err error) {
	if ty == CtxType_Func {
		defer func() {
			if e := recover(); e != nil {
				err = errors.New(id + " is not callable.")
			}
		}()
		v := reflect.ValueOf(val)
		v.Type().NumIn() // 								Test if this is a function that can be called.
		w := &ContextValueType{
			TypeOf: ty,
			Val:    val,
			Func:   v,
		}
		ctx.mutex.Lock()
		ctx.Store[id] = w
		ctx.mutex.Unlock()
		return
	}
	w := &ContextValueType{
		TypeOf: ty,
		Val:    val,
	}
	ctx.mutex.Lock()
	ctx.Store[id] = w
	ctx.mutex.Unlock()
	return
}

// Call a function that has been placed in the table
func (ctx *ContextType) Call(name string, params ...interface{}) (result []reflect.Value, err error) {
	ctx.mutex.RLock()
	fx, ok := ctx.Store[name]
	ctx.mutex.RUnlock()
	if !ok { // 								Need to have a mutext - lock
		err = errors.New(name + " function does not exist.")
		return
	}
	np := len(params)
	fmt.Printf("np=%d, name=%s\n", np, name)
	if np != fx.Func.Type().NumIn() { // 		Posssibility of default params? // Should save # of params from SetInContext call?
		err = ErrParamsNotAdapted
		return
	}
	in := make([]reflect.Value, np) // 			Type check params for correctness?
	for k, param := range params {
		if params[k] == nil {
			in[k] = reflect.ValueOf((*string)(nil))
		} else {
			in[k] = reflect.ValueOf(param)
		}
	}
	fmt.Printf("Just Before %s %d %+v, %s\n", name, np, params[0], com.LF())
	result = fx.Func.Call(in)
	return
}

func (ctx *ContextType) GetFromContext(id string) (v interface{}, t int, f bool) {
	t = CtxType_Error
	f = false
	ctx.mutex.RLock()
	if x, ok := ctx.Store[id]; ok {
		v = x.Val
		t = x.TypeOf
		f = true
	}
	ctx.mutex.RUnlock()
	return
}

func (ctx *ContextType) DumpContext() {
	ctx.mutex.RLock()
	if ctx == nil {
		fmt.Printf("Context is NIL - never created\n")
	} else {
		// fmt.Printf("Context = %s\n", sizlib.SVarI(ctx.Store))
		for ii, vv := range ctx.Store {
			if vv.TypeOf != CtxType_Func {
				fmt.Printf("\t[%s] = Type %d/%s %+v\n", ii, vv.TypeOf, ctx.NameOfType(vv.TypeOf), vv.Val)
			} else {
				fmt.Printf("\t[%s] = Type %d/%s <<<function>>>\n", ii, vv.TypeOf, ctx.NameOfType(vv.TypeOf))
			}
		}
	}
	fmt.Printf("\n\n")
	ctx.mutex.RUnlock()
}

func (ctx *ContextType) PushInContext(id string, ty int, val interface{}) {
	var p *ContextValueType
	var ok bool
	ctx.mutex.Lock()
	if p, ok = ctx.Store[id]; !ok {
		p = nil
	}
	v := &ContextValueType{
		TypeOf: ty,
		Val:    val,
		Prev:   p,
	}
	ctx.Store[id] = v
	ctx.mutex.Unlock()
}

func (ctx *ContextType) PopFromContext(id string) {
	ctx.mutex.Lock()
	if t, ok := ctx.Store[id]; ok {
		ctx.Store[id] = t.Prev
	}
	ctx.mutex.Unlock()
}

func (ctx *ContextType) IsDefinedContext(id string) (f bool) {
	ctx.mutex.RLock()
	_, f = ctx.Store[id]
	ctx.mutex.RUnlock()
	return
}

// ------------------------------------------------------------------------------------------------------------
//
// ExtendItem
//
// Given: From{}, New{}, List of Params - produce New{} as extend of From
// From - has a set of Params in it.
//

// ------------------------------------------------------------------------------------------------------------
//
// DefineTemplate
//
// Input: Name{} + Set of Params with Defaults
//
