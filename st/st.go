package st

//
// S T - Symbol table
//
// (C) Philip Schlump, 2013-2015.
// Version: 1.0.0
// BuildNo: 28
//
// Special Thanks to 2C-Why, LLC for supporting this project.
//

import (
	"errors"
	"fmt"
	"io"
	"sync"

	"github.com/pschlump/dbgo"
	"github.com/pschlump/lexie/gen"
	"github.com/pschlump/lexie/mt"
)

// Add to this Fx - pointer to function for eval of builtins?

type SymbolType struct {
	Name      string      //
	Body      string      //
	SymType   int         //                    new - Macro / Item-Name / Begnning-Block-Name / Ending Block Name? / If/Eif/Else block name?
	FxId      int         //                    new
	AnyData   interface{} //
	NParam    int         //
	ParamName []string    //
	DefVal    []string    //  					new - Default value for params in param order
	Prev      *SymbolType //
}

type SymbolTable struct {
	Symbols map[string]*SymbolType
	mutex   sync.RWMutex
}

func NewSymbolTable() (st *SymbolTable) {
	st = &SymbolTable{
		Symbols: make(map[string]*SymbolType),
	}
	return
}

var NotFoundError = errors.New("Not Found")

func (st *SymbolTable) LookupSymbol(name string) (as *SymbolType, err error) {
	ok := false
	err = nil
	st.mutex.RLock()
	if as, ok = st.Symbols[name]; ok {
		st.mutex.RUnlock()
		return
	}
	st.mutex.RUnlock()
	as = nil
	err = NotFoundError
	return
}

func (st *SymbolTable) DefineSymbol(name, body string, plist []string) (ss *SymbolType) {
	as, err := st.LookupSymbol(name)
	x := &SymbolType{
		Name:      name,               //
		Body:      body,               //
		ParamName: plist,              //
		SymType:   gen.Tok_Tree_Macro, // Tok_Tree_Macro
		FxId:      0,                  // not a builtin
		NParam:    0,                  // no params defined
	}
	if err == nil {
		x.Prev = as
	}
	st.mutex.Lock()
	st.Symbols[name] = x
	st.mutex.Unlock()
	ss = x
	return
}

func (st *SymbolTable) DefineReservedWord(name string, fxid int) (ss *SymbolType) {
	as, err := st.LookupSymbol(name)
	x := &SymbolType{
		Name:    name,               //
		SymType: gen.Tok_Tree_Macro, // Tok_Tree_Macro
		FxId:    fxid,               //
	}
	if err == nil {
		x.Prev = as
	}
	st.mutex.Lock()
	st.Symbols[name] = x
	st.mutex.Unlock()
	ss = x
	return
}

func (st *SymbolTable) UnDefineSymbol(name string) {
	as, err := st.LookupSymbol(name)
	if err == nil {
		st.mutex.Lock()
		if as.Prev != nil {
			st.Symbols[name] = as.Prev
		} else {
			delete(st.Symbols, name)
		}
		st.mutex.Unlock()
	}
}

func (st *SymbolTable) Dump01(fo io.Writer) {
	fmt.Fprintf(fo, "Dump of symbol table\n")
	for ii, vv := range st.Symbols {
		fmt.Fprintf(fo, "[%s] %s \n", ii, vv.Body)
	}
}

func (st *SymbolTable) DumpSymbolTable(fo io.Writer) {
	for ii, vv := range st.Symbols {
		fmt.Fprintf(fo, "\t[%s] Body=%s SymType=%d FxId=%d\n", ii, vv.Body, vv.SymType, vv.FxId)
		if vv.SymType == gen.Tok_Template {
			fmt.Fprintf(fo, "\t\tTemplate\n")
			mtv := vv.AnyData.(*mt.MtType)
			fmt.Fprintf(fo, "AnyData = %s\n", dbgo.SVarI(mtv))
		}
	}
}
