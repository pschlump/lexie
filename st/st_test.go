package st

//
// S T - Symbol table - test
//
// Copyright (C) Philip Schlump, 2013-2025.
//

import (
	"fmt"
	"os"
	"testing"

	"github.com/pschlump/filelib"
)

const (
	CmdInsert   = 1
	CmdLookup   = 2
	CmdDelete   = 3
	CmdDump     = 4
	CmdDefRW    = 5
	CmdLookupRW = 6
)

type ActType struct {
	OpCode      int    //
	Item        string //
	Data        string //
	IfFoundFlag rune   //
	No          int    //
}

type Pb01TestType struct {
	Test         string
	SkipTest     bool
	Actions      []ActType
	Results      string
	ResultsFound bool
}

var Pb01Test = []Pb01TestType{

	// Simple test with 2 fiels
	{Test: "0001", SkipTest: false, Actions: []ActType{
		ActType{OpCode: CmdInsert, Item: "aa", Data: "<<this is aa 1>>"},
		ActType{OpCode: CmdInsert, Item: "ab", Data: "<<this is ab>>"},
		ActType{OpCode: CmdInsert, Item: "aa", Data: "<<this is aa 2>>"},
		ActType{OpCode: CmdDump, Data: "test/1-4.out"},
		ActType{OpCode: CmdLookup, Item: "aa", Data: "<<this is aa 2>>"},
		ActType{OpCode: CmdDelete, Item: "aa"},
		ActType{OpCode: CmdLookup, Item: "aa", Data: "<<this is aa 1>>"},
		ActType{OpCode: CmdDelete, Item: "aa"},
		ActType{OpCode: CmdLookup, Item: "aa", Data: ""},
		ActType{OpCode: CmdDefRW, Item: "rw", No: 12},
		ActType{OpCode: CmdDefRW, Item: "rw", No: 14},
		ActType{OpCode: CmdLookupRW, Item: "rw", No: 14},
	}, Results: ``},
}

var db_test01 = false

func Test_St01(t *testing.T) {

	os.MkdirAll("./test", 0755)

	SymbolTable := NewSymbolTable()

	for ii, vv := range Pb01Test {
		_ = ii

		if !vv.SkipTest {

			// Implement a quick - fetch execute macine to test - the SymbolTable
			for pc, ww := range vv.Actions {

				switch ww.OpCode {
				case CmdInsert:
					SymbolTable.DefineSymbol(ww.Item, ww.Data, []string{})
				case CmdDefRW:
					SymbolTable.DefineReservedWord(ww.Item, ww.No)
				case CmdLookup:
					as, err := SymbolTable.LookupSymbol(ww.Item)
					if err == nil {
						if db_test01 {
							fmt.Printf("%s: found, value %s\n", ww.Item, as.Body)
						}
						if ww.Data == "" {
							t.Errorf("%04s: %d error, expected to have symbol in table, missing, %s\n", vv.Test, pc, ww.Item)
						}
						if as.Body != ww.Data {
							t.Errorf("%04s: %d error, expected value %s got %s invalid for, %s\n", vv.Test, pc, ww.Data, as.Body, ww.Item)
						}
					} else {
						if ww.Data != "" {
							t.Errorf("%04s: %d error, expected to NOT have symbol, found it, %s\n", vv.Test, pc, ww.Item)
						}
						if db_test01 {
							fmt.Printf("%s: not found.\n", ww.Item)
						}
					}
				case CmdLookupRW:
					as, err := SymbolTable.LookupSymbol(ww.Item)
					if err == nil {
						if db_test01 {
							fmt.Printf("%s: found, value %s\n", ww.Item, as.Body)
						}
						if as.FxId != ww.No {
							t.Errorf("%04s: %d error, expected value %d got %d invalid for, %s\n", vv.Test, pc, ww.No, as.FxId, ww.Item)
						}
					} else {
						if db_test01 {
							fmt.Printf("%s: not found.\n", ww.Item)
						}
					}
				case CmdDelete:
					SymbolTable.UnDefineSymbol(ww.Item)
				case CmdDump:
					if ww.Data == "" {
						SymbolTable.Dump01(os.Stdout)
					} else {
						fp, err := filelib.Fopen(ww.Data, "w")
						if err == nil {
							SymbolTable.Dump01(fp)
							fp.Close()
						} else {
							SymbolTable.Dump01(os.Stdout)
							t.Errorf("%04s: Unable to open file for output ->%s<-, error: %s\n", vv.Test, ww.Data, err)
						}
					}
				}
			}
		}
	}

}
