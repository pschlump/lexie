//
// C O M - Part of Lexie Lexical Generation System
//
// (C) Philip Schlump, 2014-2015.
// Version: 1.0.8
// BuildNo: 203
//
// Special Thanks to 2C-Why, LLC for supporting this project.
//

package com

import (
	"fmt"
	"html"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/pschlump/json" //	"encoding/json"

	// "../../../go-lib/flags"  //		"github.com/jessevdk/go-flags"
	// "../../../go-lib/sizlib" //
)

func USortIntSlice(inputSet []int) (rv []int) {
	// sort.Sort(sort.IntSlice(inputSet))
	uniq := make(map[int]bool)
	for _, v := range inputSet {
		uniq[v] = true
	}
	for i := range uniq {
		rv = append(rv, i)
	}
	sort.Sort(sort.IntSlice(rv))

	//	for _, v := range inputSet {
	//		have := false
	//		for _, w := range rv {
	//			if w == v {
	//				have = true
	//				break
	//			}
	//		}
	//		if !have {
	//			rv = append(rv, v)
	//		}
	//	}

	return
}

func SortMapStringString(str map[string]string) (rv []string) {
	for ii := range str {
		rv = append(rv, ii)
	}
	rv = KeyStringSort(rv)
	return
}

func KeyStringSort(str []string) (rv []string) {
	rv = str
	sort.Sort(sort.StringSlice(rv))
	return
}

func NameOf(inputSet []int) string {
	// inputSet = USortIntSlice(inputSet)
	com := ""
	s := ""
	for _, v := range inputSet {
		s += com + fmt.Sprintf("%d", v)
		com = "-"
	}
	return s
}

func CompareSlices(X, Y []int) []int {
	m := make(map[int]int)

	for _, y := range Y {
		m[y]++
	}

	var ret []int
	for _, x := range X {
		if m[x] > 0 {
			m[x]--
			continue
		}
		ret = append(ret, x)
	}

	return ret
}

func EqualStringSlices(X, Y []string) bool {
	if len(X) != len(Y) {
		return false
	}
	for ii, vv := range X {
		if Y[ii] != vv {
			return false
		}
	}
	return true
}

// DbPrintf("db_DumpDFAPool", " %12s %12s \u2714              \tEdges", "StateName", "StateSet")

func ChkOrX(v bool) string {
	if v {
		return "\u2714"
	}
	return "\u2716"
}
func ChkOrBlank(v bool) string {
	if v {
		return "\u2714"
	}
	return " "
}

func ConvertActionFlagToString(kk int) (rv string) {

	if kk == 0 {
		rv = "**No A Flag**"
		return
	}

	rv = fmt.Sprintf("(%02x) ", kk)

	com := ""
	if (kk & A_Repl) != 0 {
		rv = rv + com + "A_Repl"
		com = "|"
	}
	if (kk & A_EOF) != 0 {
		rv = rv + com + "A_EOF"
		com = "|"
	}
	if (kk & A_Push) != 0 {
		rv = rv + com + "A_Push"
		com = "|"
	}
	if (kk & A_Pop) != 0 {
		rv = rv + com + "A_Pop"
		com = "|"
	}
	if (kk & A_Observe) != 0 {
		rv = rv + com + "A_Observe"
		com = "|"
	}
	if (kk & A_Greedy) != 0 {
		rv = rv + com + "A_Greedy"
		com = "|"
	}
	if (kk & A_Reset) != 0 {
		rv = rv + com + "A_Reset"
		com = "|"
	}
	if (kk & A_NotGreedy) != 0 {
		rv = rv + com + "A_NotGreedy"
		com = "|"
	}
	if (kk & A_Error) != 0 {
		rv = rv + com + "A_Error"
		com = "|"
	}
	if (kk & A_Warning) != 0 {
		rv = rv + com + "A_Warning"
		com = "|"
	}
	if (kk & A_Alias) != 0 {
		rv = rv + com + "A_Alias"
		com = "|"
	}

	return
}

// -----------------------------------------------------------------------------------------------------------------------------------
// These may need to be "bit" flags - and call them "A_" for Actions
// -----------------------------------------------------------------------------------------------------------------------------------

var ReservedActionNames = []string{"A_Repl", "A_EOF", "A_Push", "A_Pop", "A_Observe", "A_Greedy", "A_Reset", "A_NotGreedy", "A_Error", "A_Warning", "A_Alias"}
var ReservedActionToString map[int]string

const (
	A_Repl      = 1 << iota // Replace input that matches with this - acts as a hard token
	A_EOF       = 1 << iota // Reached EOF
	A_Push      = 1 << iota // Call(x)
	A_Pop       = 1 << iota // Return
	A_Observe   = 1 << iota // Observe and report occurance of an item, continue processing
	A_Greedy    = 1 << iota // not used
	A_Reset     = 1 << iota // Reset stack to top level - restart machine (error recovery)
	A_NotGreedy = 1 << iota // Report token if if could be greedy and accumulate
	A_Error     = 1 << iota // Have an Error to report - often combined with A_Reset
	A_Warning   = 1 << iota // A warning to report
	A_Alias     = 1 << iota // An alias - replaces input and processes as if this was the original input (Different than A_Repl)
)

var ReservedActionValues = map[string]int{
	"A_Repl":      A_Repl,
	"A_EOF":       A_EOF,
	"A_Push":      A_Push,
	"A_Pop":       A_Pop,
	"A_Observe":   A_Observe,
	"A_Greedy":    A_Greedy,
	"A_Reset":     A_Reset,
	"A_NotGreedy": A_NotGreedy,
	"A_Error":     A_Error,
	"A_Warning":   A_Warning,
	"A_Alias":     A_Alias,
}

func init() {
	ReservedActionToString = make(map[int]string)
	ReservedActionToString[A_Repl] = "A_Repl"
	ReservedActionToString[A_EOF] = "A_EOF"
	ReservedActionToString[A_Push] = "A_Push"
	ReservedActionToString[A_Pop] = "A_Pop"
	ReservedActionToString[A_Observe] = "A_Observe"
	ReservedActionToString[A_Greedy] = "A_Greedy"
	ReservedActionToString[A_Reset] = "A_Reset"
	ReservedActionToString[A_NotGreedy] = "A_NotGreedy"
	ReservedActionToString[A_Error] = "A_Error"
	ReservedActionToString[A_Warning] = "A_Warning"
	ReservedActionToString[A_Alias] = "A_Alias"
}

// -----------------------------------------------------------------------------------------------------------------------------------
// Generate an array with all the files in the path
// -----------------------------------------------------------------------------------------------------------------------------------
func AllFilesInPath(path string) (filenames []string) {
	pa := strings.Split(path, ";")
	// fmt.Printf("pa=%+v\n", pa)
	for _, dir := range pa {
		t, dirs := GetFilenames(dir)
		// fmt.Printf("t=%+v dirs=%+v \n", t, dirs)
		for _, fn := range t {
			filenames = append(filenames, dir+"/"+fn)
		}
		// fmt.Printf("After append, filenames = %+v\n", filenames)
		for _, aDir := range dirs {
			t2 := AllFilesInPath(dir + "/" + aDir)
			for _, fn := range t2 {
				filenames = append(filenames, fn)
			}
		}
	}
	return
}

// -----------------------------------------------------------------------------------------------------------------------------------
// Exists reports whether the named file or directory exists.
// -----------------------------------------------------------------------------------------------------------------------------------
func DirExists(name string) bool {
	if fstat, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	} else {
		if fstat.IsDir() {
			return true
		}
	}
	return false
}

// -------------------------------------------------------------------------------------------------
// Exists reports whether the named file or directory exists.
// -------------------------------------------------------------------------------------------------
func Exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

// -------------------------------------------------------------------------------------------------
// Get a list of filenames and directorys.
// -------------------------------------------------------------------------------------------------
func GetFilenames(dir string) (filenames, dirs []string) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, nil
	}
	for _, fstat := range files {
		if !strings.HasPrefix(string(fstat.Name()), ".") {
			if fstat.IsDir() {
				dirs = append(dirs, fstat.Name())
			} else {
				filenames = append(filenames, fstat.Name())
			}
		}
	}
	return
}

// -------------------------------------------------------------------------------------------------
// -------------------------------------------------------------------------------------------------
func EscapeStr(v string, on bool) string {
	if on {
		return html.EscapeString(v)
	} else {
		return v
	}
}

// -------------------------------------------------------------------------------------------------
// xyzzy - need "fast" version of "CompareFiles" with some limits on what it will use "fast"
// compare for - .jpg,.gif,.png fiels - a fiel size before uses fast etc.   Compare Size?
// Compare name?  What is the "fast" compare for rsync? -- Calculate Hashes for each and
// keep them around?
// -------------------------------------------------------------------------------------------------
func CompareFiles(cmpFile string, refFile string) bool {
	cmp, err := ioutil.ReadFile(cmpFile)
	if err != nil {
		fmt.Printf("Unable to read %s\n", cmpFile)
		return false
	}

	if Exists(refFile) {
		ref, err := ioutil.ReadFile(refFile)
		if err != nil {
			fmt.Printf("Unable to read %s\n", refFile)
			return false
		}
		if len(ref) != len(cmp) { // xyzzy - Could be faster - just check lenths on disk - if diff then return false
			return false
		}
		if string(ref) != string(cmp) {
			return false
		}
	} else {
		return false
	}
	return true
}

// -------------------------------------------------------------------------------------------------
// Get a list of filenames and directories.
// -------------------------------------------------------------------------------------------------
func GetFilenamesRecrusive(dir string) (filenames, dirs []string, err error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, nil, err
	}
	//for ii, fstat := range files {
	//	fmt.Printf("Top files %d:[%s]\n", ii, fstat.Name())
	//}
	for _, fstat := range files {
		if !strings.HasPrefix(string(fstat.Name()), ".") {
			if fstat.IsDir() {
				name := fstat.Name()
				dirs = append(dirs, dir+"/"+name)
				// fmt.Printf("Recursive dir [%s]\n", dir+"/"+name)
				tf, td, err := GetFilenamesRecrusive(dir + "/" + name)
				if err != nil {
					return nil, nil, err
				}
				filenames = append(filenames, tf...)
				dirs = append(dirs, td...)
			} else {
				name := fstat.Name()
				name = dir + "/" + name
				// fmt.Printf("dir %s ->%s<-\n", dir, name)
				filenames = append(filenames, name)
			}
		}
	}
	return
}

// -------------------------------------------------------------------------------------------------
// -------------------------------------------------------------------------------------------------
func FilterArray(re string, inArr []string) (outArr []string) {
	var validID = regexp.MustCompile(re)

	outArr = make([]string, 0, len(inArr))
	for k := range inArr {
		if validID.MatchString(inArr[k]) {
			outArr = append(outArr, inArr[k])
		}
	}
	// fmt.Printf ( "output = %v\n", outArr )
	return
}

// -------------------------------------------------------------------------------------------------
// -------------------------------------------------------------------------------------------------
func GetTemplateFiles(dir string) (fns []string, err error) {

	fns, _, err = GetFilenamesRecrusive(dir)
	if err != nil {
		return
	}
	fns = FilterArray(".*\\.tpl$", fns)

	return
}

// -------------------------------------------------------------------------------------------------
// 1.  For each dir - Create destination directies -o <name>/+/...		ReplaceEach ( []string, pat, repl )
// -------------------------------------------------------------------------------------------------
func ReplaceEach(data []string, pat, repl string) (outArr []string) {
	for ii, vv := range data {
		_ = ii
		t := strings.Replace(vv, pat, repl, 1)
		outArr = append(outArr, t)
	}
	return
}

// From: http://stackoverflow.com/questions/21060945/simple-way-to-copy-a-file-in-golang
// Modified.

// CopyFile copies a file from src to dst. If src and dst files exist, and are
// the same, then return success. Otherise, if useHardLink is true, an attempt
// to create a hard link between the two files.
// If that fail, copy the file contents from src to dst.
func CopyFile(src, dst string, useHardLink bool) (err error) {
	sfi, err := os.Stat(src)
	if err != nil {
		return
	}
	if !sfi.Mode().IsRegular() {
		// cannot copy non-regular files (e.g., directories, symlinks, devices, etc.)
		return fmt.Errorf("Error: non-regular source file %s (%q)", sfi.Name(), sfi.Mode().String())
	}
	dfi, err := os.Stat(dst)
	if err != nil {
		if !os.IsNotExist(err) {
			return
		}
	} else {
		if !(dfi.Mode().IsRegular()) {
			return fmt.Errorf("Error: non-regular destination file %s (%q)", dfi.Name(), dfi.Mode().String())
		}
		if os.SameFile(sfi, dfi) {
			return
		}
	}
	if useHardLink {
		if err = os.Link(src, dst); err == nil {
			return
		}
	}
	err = copyFileContents(src, dst)
	return
}

// copyFileContents copies the contents of the file named src to the file named
// by dst. The file will be created if it does not already exist. If the
// destination file exists, all it's contents will be replaced by the contents
// of the source file.
func copyFileContents(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return
	}
	err = out.Sync()
	return
}

// -------------------------------------------------------------------------------------------------------
//func main() {
//	 fmt.Printf("Copying %s to %s\n", os.Args[1], os.Args[2])
//	 err := CopyFile(os.Args[1], os.Args[2])
//	 if err != nil {
//		  fmt.Printf("CopyFile failed %q\n", err)
//	 } else {
//		  fmt.Printf("CopyFile succeeded\n")
//	 }
//}

var trueValues map[string]bool

func init() {
	trueValues = make(map[string]bool)
	trueValues["t"] = true
	trueValues["T"] = true
	trueValues["yes"] = true
	trueValues["Yes"] = true
	trueValues["YES"] = true
	trueValues["1"] = true
	trueValues["true"] = true
	trueValues["True"] = true
	trueValues["TRUE"] = true
	trueValues["on"] = true
	trueValues["On"] = true
	trueValues["ON"] = true
}

func ParseBool(s string) (b bool) {
	_, b = trueValues[s]
	return
}

// -------------------------------------------------------------------------------------------------
func SVar(v interface{}) string {
	s, err := json.Marshal(v)
	// s, err := json.MarshalIndent ( v, "", "\t" )
	if err != nil {
		return fmt.Sprintf("Error:%s", err)
	} else {
		return string(s)
	}
}

// -------------------------------------------------------------------------------------------------
func SVarI(v interface{}) string {
	// s, err := json.Marshal ( v )
	s, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		return fmt.Sprintf("Error:%s", err)
	} else {
		return string(s)
	}
}

// -------------------------------------------------------------------------------------------------
// Exists returns true if the named file or directory exists.
//func Exists(name string) bool {
//	if _, err := os.Stat(name); err != nil {
//		if os.IsNotExist(err) {
//			return false
//		}
//	}
//	return true
//}

// -------------------------------------------------------------------------------------------------
// xyzzy - str.
// Return the basename from a file path.  This is the last component with the directory path
// stripped off.  File extension removed.
// -------------------------------------------------------------------------------------------------
func Basename(fn string) (bn string) {
	i, j := strings.LastIndex(fn, "/"), strings.LastIndex(fn, path.Ext(fn)) // xyzzy windoz
	// fmt.Printf ( "i=%d j=%d\n", i, j )
	if i < 0 && j < 0 {
		bn = fn
	} else if i < 0 {
		bn = fn[0:j]
	} else {
		bn = fn[i+1 : j]
	}
	return
}

// -------------------------------------------------------------------------------------------------
// xyzzy - str.
// With file extension
// -------------------------------------------------------------------------------------------------
func BasenameExt(fn string) (bn string) {
	i, j := strings.LastIndex(fn, "/"), len(fn) // xyzzy windoz
	// fmt.Printf ( "i=%d j=%d\n", i, j )
	if i < 0 && j < 0 {
		bn = fn
	} else if i < 0 {
		bn = fn[0:j]
	} else {
		bn = fn[i+1 : j]
	}
	return
}

func InArray(lookFor string, inArr []string) bool {
	for _, v := range inArr {
		if lookFor == v {
			return true
		}
	}
	return false
}

func RmExt(filename string) string {
	var extension = filepath.Ext(filename)
	var name = filename[0 : len(filename)-len(extension)]
	return name
}

var qtRegEx *regexp.Regexp

func init() {
	qtRegEx = regexp.MustCompile("%{([A-Za-z0-9_]*)%}")
}

// QT: Quick template
// %{name%} gets replace with substitution from map if it is in map, else ""
func Qt(format string, data map[string]string) string {
	// re := regexp.MustCompile("%{([A-Za-z0-9_]*)%}")
	return string(qtRegEx.ReplaceAllFunc([]byte(format), func(s []byte) []byte {
		return []byte(data[string(s[2:len(s)-2])])
	}))
}

func QtR(format string, data map[string]interface{}) string {
	return string(qtRegEx.ReplaceAllFunc([]byte(format), func(s []byte) []byte {
		t := string(s[2 : len(s)-2])
		u, ok := data[t]
		if !ok {
			return []byte("")
		}
		switch u.(type) {
		case string:
			return []byte(u.(string))
		default:
			sb := fmt.Sprintf("%v", u)
			return []byte(sb)
		}
	}))
}

/* vim: set noai ts=4 sw=4: */
