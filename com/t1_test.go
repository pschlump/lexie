package com

import (
	"fmt"
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

	// Tests DirExists
	if DirExists("./t1/a.tpl") {
		t.Errorf("DirExists() failed.")
	}
	if !DirExists("./t1") {
		t.Errorf("DirExists() failed.")
	}
	if DirExists("./t5") {
		t.Errorf("DirExists() failed.")
	}

	StashError("s")
	StashError("t")

	ns := GetErrorStash()
	if ns == "" {
		t.Errorf("Falied to get stashed errors\n")
	}

	ss := ConvertActionFlagToString(0xFFFF)
	if ss != "(ffff) A_Repl|A_EOF|A_Push|A_Pop|A_Observe|A_Greedy|A_Reset|A_NotGreedy|A_Error|A_Warning|A_Alias" {
		t.Errorf("Dislaying flags\n")
	}
	ss = ConvertActionFlagToString(0)
	if ss != "**No A Flag**" {
		t.Errorf("Dislaying flags (2)\n")
	}

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

}
