package com

import (
	"fmt"
	"os"
	"testing"

	"github.com/pschlump/dbgo"
)

const db_flag = false

func Test_Com01(t *testing.T) {

	return

	v := []string{
		"./t1/a.tpl",
		"./t1/t2/t3/c.tpl",
		"./t1/t2/t3/c2.tpl",
		"./t2/aaa.tpl",
	}

	fx := AllFilesInPath("./t1;./t2;./t3")

	if db_flag {
		fmt.Printf("%s\n", dbgo.SVarI(fx))
	}

	if len(fx) != len(v) {
		t.Errorf("Length of responce did not match expeced length, should be 3, got %d\n", len(fx))
	} else {
		for ii, ww := range v {
			if ww != fx[ii] {
				t.Errorf("Expected %s got %s in list of files\n", ww, fx[ii])
			}
		}
	}

	// Tests Exists
	if Exists("./t1/a.tpl") {
		fmt.Printf("ok: %s %s\n", LINE(1), FILE(1))
	} else {
		t.Errorf("Exists() %s %s\n", LINE(), FILE())
	}
	if Exists("./t1/ab.tpl") {
		t.Errorf("Exists() %s %s\n", LINE(), FILE())
	} else {
		fmt.Printf("ok: %d %s\n", LINEn(1), FILE())
	}

	DbOnFlags["debug_test"] = true

	// Tests DirExists
	if DirExists("./t1/a.tpl") {
		t.Errorf("DirExists() %s %s\n", LINE(), FILE())
	} else {
		DbPrintf("debug_test", "ok: %s\n", LF(1))
	}
	if DirExists("./t1") {
		if DbOn("debug_test") {
			fmt.Printf("ok: %s\n", LF())
		}
	} else {
		t.Errorf("DirExists() %s %s\n", LINE(), FILE())
	}
	if DirExists("./t5") {
		t.Errorf("DirExists() %s %s\n", LINE(), FILE())
	} else {
		fmt.Printf("ok: %s\n", LF())
	}

	//
	//
	//
	//
	//
	//
	//
	//
	//
	//
	//
	//
	//
	//
	//
	//
	//
	//
	//
	//
	//
	//
	//
	//
	//
	//
	//
	//  Should be line 100 in this file - change comments above to adjust
	ln := LINE()
	fn := FILE()
	lnn := LINEn(1)
	ln1 := LINE(1)
	fn1 := FILE(1)
	lf := LF()
	lf1 := LF(1)

	if ln != "101" {
		t.Errorf("LINE() failed\n")
	}
	if fn != "/Users/corwin/Projects/pongo2/lexie/com/t1_test.go" {
		t.Errorf("FILE() failed, got >%s<\n", fn)
	}
	if lnn != 103 {
		t.Errorf("LINEn() failed\n")
	}
	if ln1 != "104" {
		t.Errorf("LINEn(1) failed\n")
	}
	if fn1 != "/Users/corwin/Projects/pongo2/lexie/com/t1_test.go" {
		t.Errorf("FILE(1) failed, got >%s<\n", fn)
	}
	if lf != "File: /Users/corwin/Projects/pongo2/lexie/com/t1_test.go LineNo:106" {
		t.Errorf("LF() failed, got >%s<\n", lf)
	}
	if lf1 != "File: /Users/corwin/Projects/pongo2/lexie/com/t1_test.go LineNo:107" {
		t.Errorf("LF(1) failed, got >%s<\n", lf1)
	}

	StashError("s")
	DbOnFlags["OutputErrors"] = true
	StashError("t")

	fp1, err := Fopen("tF/a.a", "w")
	if err != nil {
		t.Errorf("Fopen 1\n")
	}
	fp1.Close()

	fp1, err = Fopen("tF/a.a", "r")
	if err != nil {
		t.Errorf("Fopen 2\n")
	}
	fp1.Close()

	fp1, err = Fopen("tF/a.a", "a")
	if err != nil {
		t.Errorf("Fopen 3\n")
	}
	fp1.Close()

	fp1, err = Fopen("tF/a.a", "x")
	if err == nil {
		t.Errorf("Fopen 4\n")
	}

	ss := ConvertActionFlagToString(0xFFFF)
	if ss != "(ffff) A_Repl|A_EOF|A_Push|A_Pop|A_Observe|A_Greedy|A_Reset|A_NotGreedy|A_Error|A_Warning|A_Alias" {
		t.Errorf("Dislaying flags\n")
	}
	ss = ConvertActionFlagToString(0)
	if ss != "**No A Flag**" {
		t.Errorf("Dislaying flags (2)\n")
	}
	// fmt.Printf("ss= >%s<-\n", ss)

	x1 := ChkOrX(true)
	if x1 != "\u2714" {
		t.Errorf("err")
	}
	x1 = ChkOrX(false)
	if x1 != "\u2716" {
		t.Errorf("err")
	}

	x1 = ChkOrBlank(false)
	if x1 == "\u2714" {
		t.Errorf("err")
	}
	x1 = ChkOrBlank(true)
	if x1 == " " {
		t.Errorf("err")
	}

	vv := CompareSlices([]int{1, 2, 4, 4}, []int{2, 1, 2, 3})
	// fmt.Printf("vv=%+v\n", vv)
	if len(vv) != 2 {
		t.Errorf("err")
	} else {
		if vv[0] != 4 || vv[1] != 4 {
			t.Errorf("err")
		}
	}

	qq := NameOf([]int{1, 4, 7})
	if qq != "1-4-7" {
		t.Errorf("err")
	}

	q2 := USortIntSlice([]int{1, 5, 2, 1})
	// fmt.Printf("q2=%+v\n", q2)
	if len(q2) != 3 {
		t.Errorf("err")
	} else {
		if q2[0] != 1 || q2[1] != 2 || q2[2] != 5 {
			t.Errorf("err")
		}
	}

	ms := make(map[string]string)
	ms["a"] = "1"
	ms["c"] = "2"
	ms["b"] = "3"
	s2 := SortMapStringString(ms)
	// fmt.Printf("%+v\n", s2)
	if len(s2) != 3 {
		t.Errorf("err")
	} else {
		if s2[0] != "a" || s2[1] != "b" || s2[2] != "c" {
			t.Errorf("err")
		}
	}

	aa := EscapeStr("<>", true)
	if aa != "&lt;&gt;" {
		t.Errorf("err")
	}
	aa = EscapeStr("<>", false)
	if aa != "<>" {
		t.Errorf("err")
	}

	DbFprintf("debug_test", os.Stdout, "ok: %s\n", LF(1))

}
