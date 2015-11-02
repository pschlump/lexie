package com

import (
	"fmt"
	"testing"
)

// "../../../go-lib/sizlib"

const db_flag2 = false

func Test_Com02(t *testing.T) {

	// DbFprintf("debug_test", os.Stdout, "ok: %s\n", LF(1))

	ref_filenames := []string{"./t2-001/a/b/b1", "./t2-001/a/b/b2", "./t2-001/a/b/b3", "./t2-001/f1.txt"}
	ref_dirs := []string{"./t2-001/a", "./t2-001/a/b", "./t2-001/aa"}

	filenames, dirs, err := GetFilenamesRecrusive("./t2-001")

	if db_flag2 {
		fmt.Printf("filenames = %+v\n", filenames)
		fmt.Printf("dirs = %+v\n", dirs)
		fmt.Printf("err = %s\n", err)
	}

	if err != nil {
		t.Errorf("Expected success for list of files/dirs got %s error\n", err)
	}
	if !EqualStringSlices(ref_filenames, filenames) {
		t.Errorf("Expected matching file names, expected %s got %s\n", ref_filenames, filenames)
	}
	if !EqualStringSlices(ref_dirs, dirs) {
		t.Errorf("Expected matching directory names, expected %s got %s\n", ref_dirs, dirs)
	}

	fmt.Printf("--------------------------------\n")
	filenames, dirs, err = GetFilenamesRecrusive("./site/www_lexie_com")
	fmt.Printf("filenames = %+v\n", filenames)
	fmt.Printf("dirs = %+v\n", dirs)
	fmt.Printf("err = %s\n", err)

	oa := ReplaceEach(dirs, "./site/www_lexie_com", "./www///www_lexie_com")
	fmt.Printf("oa = %+v\n", oa)

}
