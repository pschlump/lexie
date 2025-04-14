//
// I N - Part of Lexie Lexical Generation System
//
// Copyright (C) Philip Schlump, 2014-2025.
// Version: 1.0.8
//
// Basic scanning and input to convert an input file into an internal data structure.
// It is ironic that this tool uses a hand-coded scanner when being a scanner generator.
//
// Error: TODO xyzzy108: This is not relly correct, try a comment inside a quoted string and see why

package in

// --------------------------------------------------------------------------------------------------------------
//
// Ideas
//		1. Investigate empty token in machines - xyzzy100
// 		2. Output - Awesome Output from generation engine - that can be read back in
//
// Enhancements
// 		1. Common Prefix - Boyer More Pattern Matching - Think RestMatch
//
// --------------------------------------------------------------------------------------------------------------

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/pschlump/dbgo"
	"github.com/pschlump/lexie/com"
)

type ImMachineType struct {
	Name   string        //
	Mixins []string      //
	Rules  []*ImRuleType //
	Defs   *ImDefsType   //
}

type ImRuleType struct {
	PatternType  int    // Pattern, Str0,1,2, $eof etc.  // Pattern Stuff --------------------------------------------------------------------------------
	Pattern      string //
	LineNo       int    // Error Reporintg Stuff ------------------------------------------------------------------------
	FileName     string //
	Rv           int    // ActionInfo Stuff -----------------------------------------------------------------------------
	RvName       string //
	Call         int    // Final machine number that is being called.
	CallName     string //
	Repl         bool   //
	ReplString   string //
	Ignore       bool   //
	ReservedWord bool   //
	Warn         bool   //
	Err          bool   //
	WEString     string //
	Return       bool   //
	Reset        bool   //
	NotGreedy    bool   //
}

type ImSeenAtType struct {
	LineNo   []int    //
	FileName []string //
}

type ImDefinedValueType struct {
	Seq          int                     //
	WhoAmI       string                  //
	NameValueStr map[string]string       //
	NameValue    map[string]int          //
	Reverse      map[int]string          //
	SeenAt       map[string]ImSeenAtType //
}

var Tok_map = make(map[int]string)

type ImDefsType struct {
	DefsAre map[string]ImDefinedValueType //
}

type ImType struct {
	Def        ImDefsType      //
	Machine    []ImMachineType //
	InputLines []string
}

func ReadFileIntoLines(fn string) (rv []string) {
	s, err := ioutil.ReadFile(fn)
	if err != nil {
		fmt.Printf("Unable to read %s, error: %s\n", fn, err)
		return
	}
	rv = strings.Split(string(s), "\n")
	return
}

func ReadFileIntoString(fn string) string {
	s, err := ioutil.ReadFile(fn)
	if err != nil {
		fmt.Printf("Unable to read %s, error:%s\n", fn, err)
		return ""
	}
	return string(s)
}

type ClsOfString int

const (
	ClsMachine ClsOfString = 1
	ClsEnd     ClsOfString = 2
	ClsEof     ClsOfString = 3
	ClsDef     ClsOfString = 4
	ClsStr0    ClsOfString = 5
	ClsStr1    ClsOfString = 6
	ClsStr2    ClsOfString = 7
	ClsPattern ClsOfString = 8
)

func (cls ClsOfString) String() string {
	switch cls {
	case ClsMachine:
		return "$machine"
	case ClsEnd:
		return "$end"
	case ClsEof:
		return "$eof"
	case ClsDef:
		return "$def"
	case ClsStr0:
		return "str0"
	case ClsStr1:
		return "str1"
	case ClsStr2:
		return "str2"
	case ClsPattern:
		return "pattern"
	}
	return "--Unknown--ClsOfString--"
}

var reDollarMachine = regexp.MustCompile(`^\$machine[ 	]*\(`)
var reDollarEnd = regexp.MustCompile(`^\$end[ 	]*$`)
var reDollarEof = regexp.MustCompile(`^\$eof[ 	]*$`)
var reDollarDef = regexp.MustCompile(`^\$def[ 	]*\(`)

// Clasify a line into it's basic type.
func clasifyLine(ln string) (cls ClsOfString) {
	// if strings.HasPrefix(ln, "$machine") {
	if reDollarMachine.MatchString(ln) {
		return ClsMachine
	}
	// if strings.HasPrefix(ln, "$end") {
	if reDollarEnd.MatchString(ln) {
		return ClsEnd
	}
	// if strings.HasPrefix(ln, "$eof") {
	if reDollarEof.MatchString(ln) {
		return ClsEof
	}
	// if strings.HasPrefix(ln, "$def") {
	if reDollarDef.MatchString(ln) {
		return ClsDef
	}
	if strings.HasPrefix(ln, "\"") {
		return ClsStr0
	}
	if strings.HasPrefix(ln, "'") {
		return ClsStr1
	}
	if strings.HasPrefix(ln, "`") {
		return ClsStr2
	}
	return ClsPattern
}

func EscapeLiteralString(in string) (rv string) {
	rv = ""
	for _, c := range in {
		switch c {
		case '{', '}', '*', '+', '^', '$', '.', '|', '[', '(': // , '(':
			rv += `\`
		}
		rv += string(c)
	}
	return
}

// pat = EscapeNormalString(pat)
func EscapeNormalString(in string) (rv string) {
	rv = ""
	var c rune
	var sz int

	for i := 0; i < len(in); i += sz {
		c, sz = utf8.DecodeRune([]byte(in[i:]))
		if c == '\\' {
			i += sz
			c, sz = utf8.DecodeRune([]byte(in[i:]))
			switch c {
			case 'n':
				rv += "\n"
			case 't':
				rv += "\t"
			case 'f':
				rv += "\f"
			case 'r':
				rv += "\r"
			//case 'v':
			//	rv += "\v"
			default:
				rv += string(c)
			}
		} else {
			rv += string(c)
		}
	}
	return
}

func PickOffPatternAtBeginning(cls ClsOfString, ln string) (pat string, rest string) {
	var ii int
	// cls := clasifyLine(ln)
	//fmt.Printf("cls = %s, %s\n", cls, dbgo.LF())
	switch cls {
	case ClsStr0: // "str0":
		pat = ""
		for ii = 1; ii < len(ln); ii++ {
			if ln[ii] == '\\' && ii+1 < len(ln) {
				if ln[ii+1] == '"' {
					pat += "\""
				} else {
					pat += "\\"
					pat += string(ln[ii+1])
				}
				ii++
			} else if ln[ii] == '"' {
				break
			} else {
				pat += ln[ii : ii+1]
			}
		}
		pat = EscapeNormalString(pat)
		// pat = ln[1:ii]
		if ii+1 < len(ln) {
			rest = ln[ii+1:]
		}
	case ClsStr1: // "str1":
		pat = ""
		for ii = 1; ii < len(ln); ii++ {
			if ln[ii] == '\\' && ii+1 < len(ln) {
				if ln[ii+1] == '\'' {
					pat += "'"
				} else {
					pat += "\\"
					pat += string(ln[ii+1])
				}
				ii++
			} else if ln[ii] == '\'' {
				break
			} else {
				pat += ln[ii : ii+1]
			}
		}
		pat = EscapeNormalString(pat)
		// pat = ln[1:ii]
		if ii+1 < len(ln) {
			rest = ln[ii+1:]
		}
		//fmt.Printf("ii = %d ln[]= ->%s<-, %s\n", ii, ln[1:ii], dbgo.LF())
		//fmt.Printf("pat ->%s<-, %s\n", pat, dbgo.LF())
		//fmt.Printf("rest ->%s<-, %s\n", rest, dbgo.LF())
	case ClsStr2: // "str2":
		pat = ""
		for ii = 1; ii < len(ln); ii++ {
			if ii+1 < len(ln) && ln[ii] == '`' && ln[ii+1] == '`' {
				pat += "`"
				ii++
			} else if ln[ii] == '`' {
				break
			} else {
				pat += ln[ii : ii+1]
			}
		}
		if ii+1 < len(ln) {
			rest = ln[ii+1:]
		}
		pat = EscapeLiteralString(pat)
		//fmt.Printf("ii = %d, %s\n", ii, dbgo.LF())
		//fmt.Printf("pat ->%s<-, %s\n", pat, dbgo.LF())
		//fmt.Printf("rest ->%s<-, %s\n", rest, dbgo.LF())
	case ClsPattern: // "pattern":
		for ii = 0; ii < len(ln); ii++ {
			if ln[ii] == ' ' || ln[ii] == '\t' {
				break
			}
		}
		pat = ln[0:ii]
		// fmt.Printf("pat ->%s<-, %s\n", pat, dbgo.LF())
		if ii+1 < len(ln) {
			rest = ln[ii:]
		}
		// fmt.Printf("rest ->%s<-, %s\n", rest, dbgo.LF())
	}
	return
}

var pa_re = regexp.MustCompile("[ \t]*:([ \t]*)|([a-zA-Z]+([^ \t]*))*")
var pnv_re = regexp.MustCompile("([a-zA-Z_][a-zA-Z0-9_]*)(=(.*))?")
var fx_re = regexp.MustCompile("([a-zA-Z_][a-zA-Z0-9_]*)\\([ \t]*([^) \t]*[ \t]*)\\)")
var pl_re = regexp.MustCompile("((([a-zA-Z_][a-zA-Z0-9_]*)((=[^, ]*)?)))*")
var com_re = regexp.MustCompile("[ \t]*//.*$")
var empty_re = regexp.MustCompile("^[ \t]*$")
var def_left_re = regexp.MustCompile("^[ \t]*\\$def[ \t]*\\(")
var def_right_re = regexp.MustCompile("[ \t]*\\)[ \t]*$")
var mach_left_re = regexp.MustCompile("^[ \t]*\\$machine[ \t]*\\(")
var mach_right_re = regexp.MustCompile("[ \t]*\\)[ \t]*$")
var numeric_re = regexp.MustCompile("^[0-9]+$")

// IsEmptyLine will return true if th eline is empty.  Blanks and tabs only will return true.
func IsEmptyLine(ln string) bool {
	a := empty_re.FindAllStringSubmatch(ln, -1)
	if len(a) > 0 {
		return true
	}
	return false
}

// IsNumeric will return true if the line all numeric 0-9.
func IsNumeric(s string) bool {
	a := numeric_re.FindAllStringSubmatch(s, -1)
	if len(a) > 0 {
		return true
	}
	return false
}

func ParseAction(ln string) [][]string {
	Action := pa_re.FindAllStringSubmatch(ln, -1)
	return Action
}

func ParsePattern(cls ClsOfString, ln string) (pat string, flag string, opt []string) {
	flag = ""
	pat, rest := PickOffPatternAtBeginning(cls, ln)
	// fmt.Printf("pat >%s< rest >%s<, %s\n", pat, rest, dbgo.LF())
	re := ParseAction(rest)
	// fmt.Printf("ln ->%s<- re %s\n", ln, dbgo.SVarI(re))

	for i := 1; i < len(re); i++ {
		if re[i][0] != "" {
			opt = append(opt, re[i][0])
		}
	}
	return
}

// ParseNameValue will convert a string like Tok_Name=1 into a pair, name="Tok_Name", value="1".
func ParseNameValue(nv string) (name string, value string) {
	name, value = "", ""
	t1 := pnv_re.FindAllStringSubmatch(nv, -1)
	dbgo.DbPrintf("in", "t1=%s\n", dbgo.SVarI(t1))
	if t1 != nil && len(t1[0]) > 0 {
		name = t1[0][1]
		if len(t1[0]) > 3 {
			value = t1[0][3]
		}
	} else {
		name = nv
	}
	return
}

// RemoveComment will remove any comments from a line of text, returning the modified line.
// Trailing white space will be removed.
// Error: TODO xyzzy108: This is not relly correct, try a comment inside a quoted string and see why
func RemoveComment(ln string) (oln string) {
	// com_re = regexp.MustCompile("[ \t]*//.*$")
	oln = com_re.ReplaceAllLiteralString(ln, "")
	return
}

// Rv(Name) Ignore(Xxx)
func ParseActionItem(act string) (aa string, pp string) {
	aa, pp = "", ""
	t1 := fx_re.FindAllStringSubmatch(act, -1)
	if t1 != nil {
		dbgo.DbPrintf("in", "t1=%s\n", dbgo.SVarI(t1))
		aa = t1[0][1]
		if len(t1[0]) > 1 {
			pp = t1[0][2]
		}
	} else {
		aa = act
	}
	return
}

func ParsePlist(pl string) (aa []string) {
	t1 := pl_re.FindAllStringSubmatch(pl, -1)
	if t1 != nil {
		dbgo.DbPrintf("in", "t1=%s\n", dbgo.SVarI(t1))
		for _, vv := range t1 {
			if len(vv) > 3 && vv[2] != "" {
				aa = append(aa, vv[2])
			}
		}
	}
	return
}

func NewIm() (rv *ImType) {
	rv = &ImType{}
	rv.Def.DefsAre = make(map[string]ImDefinedValueType)
	return
}

// Found $def ->$def(Tokens, Tok_null=0, Tok_ID=1 )<-
func ParseDef(ln string) (aa []string) {
	ln = def_left_re.ReplaceAllLiteralString(ln, "")
	ln = def_right_re.ReplaceAllLiteralString(ln, "")
	aa = ParsePlist(ln)
	// fmt.Printf("ParseDef: Plist %v\n", aa)
	return
}

func ParseMachine(ln string) (aa []string) {
	ln = mach_left_re.ReplaceAllLiteralString(ln, "")
	ln = mach_right_re.ReplaceAllLiteralString(ln, "")
	aa = ParsePlist(ln)
	// fmt.Printf("ParseMachine: Plist %v\n", aa)
	return
}

func validateDefType(DefType string) bool {
	if !com.InArray(DefType, []string{"Tokens", "Machines", "Errors", "ReservedWords"}) {
		fmt.Printf("Error Invalid $def type -->%s<--, should be one of \"Tokens\", \"Machines\", \"Errors\", \"ReservedWords\" \n", DefType)
		return false
	}
	return true
}

// -----------------------------------------------------------------------------------------------------------------------------
// -----------------------------------------------------------------------------------------------------------------------------
func (Im *ImType) SaveDef(DefType string, Defs []string, line_no int, file_name string) {
	if validateDefType(DefType) {
		for _, nm := range Defs {
			dd, ok := Im.Def.DefsAre[DefType]
			if !ok {
				dd = ImDefinedValueType{
					Seq:          1,
					WhoAmI:       DefType,
					NameValue:    make(map[string]int),
					NameValueStr: make(map[string]string),
					Reverse:      make(map[int]string),
					SeenAt:       make(map[string]ImSeenAtType),
				}
			}
			// seq := dd.Seq
			n, v := ParseNameValue(nm)
			if n == "" && v == "" { // xyzzy100
				return
			}
			dbgo.DbPrintf("in", "Input: ->%s<- n >%s< v >%s<\n", nm, n, v)
			if v != "" {
				dd.NameValue[n] = -2
				if vv, ok1 := dd.NameValueStr[n]; !ok1 {
					dd.NameValueStr[n] = v
				} else {
					if vv != v {
						fmt.Printf("Error: Attempt to redfine %s from %s to %s - Probably an error\n", n, vv, v)
					}
				}
			} else {
				dd.NameValue[n] = -1
				if _, ok1 := dd.NameValueStr[n]; !ok1 {
					dd.NameValueStr[n] = ""
				}
			}
			sa := dd.SeenAt[n]
			sa.LineNo = append(sa.LineNo, line_no)
			sa.FileName = append(sa.FileName, file_name)
			dd.SeenAt[n] = sa
			Im.Def.DefsAre[DefType] = dd
		}
	}
}

func (Im *ImType) ParseFile(data []string) {
	var st = 0
	var MNo = 0
	// Im.SaveDef("Tokens", []string{"Tok_null=0", "Tok_ID=1", "Tok_Ignore=2"}, com.LINEn(), com.FILE())
	Im.SaveDef("Tokens", []string{"Tok_null=0"}, dbgo.LINEn(), dbgo.FILE())
	for line_no_m1, line := range data {
		line_no := line_no_m1 + 1
		line = RemoveComment(line)
		if !IsEmptyLine(line) {
			cls := clasifyLine(line)
			switch cls {
			case ClsMachine: // "$machine":
				st = 1
				// fmt.Printf("Found $machine ->%s<-\n", line)
				m := ParseMachine(line) // parse machine
				MNo = Im.SaveMachine(m) //  save machine
			case ClsEnd: // "$end":
				st = 0
				// fmt.Printf("Found $end ->%s<-\n", line)
			case ClsDef: // "$def":
				// fmt.Printf("Found $def ->%s<-\n", line)
				if st != 0 {
					fmt.Printf("Error: $def found inside of a machine specificaiton, Line: %d\n", line_no)
				}
				d := ParseDef(line)
				Im.SaveDef(d[0], d[1:], line_no, "unk-file")
			case ClsStr0: // "str0":
				fallthrough
			case ClsStr1: // "str1":
				fallthrough
			case ClsStr2: // "str2":
				fallthrough
			case ClsPattern: // "pattern":
				pat, _, opt := ParsePattern(cls, line)
				// fmt.Printf("pat >%s< opt >%s<\n", pat, opt)
				Im.SavePattern(MNo, pat, false, opt, line_no, "unk-file")
			case ClsEof: // "$eof":
				// fmt.Printf("Found $eof ->%s<-\n", line)
				if st != 1 {
					fmt.Printf("Error: $eof found outside of a machine specificaiton, Line: %d\n", line_no)
				}
				_, _, opt := ParsePattern(ClsPattern, line[1:]) // parse $eof pattern -
				Im.SavePattern(MNo, "", true, opt, line_no, "unk-file")
			default:
				panic("Unreacable Code")
			}
		}
	}
	Im.SaveDef("Tokens", []string{"Tok_ID", "Tok_Ignore"}, dbgo.LINEn(), dbgo.FILE())
	Im.FinializeFile()
	return
}

func (Im *ImType) SavePattern(MNo int, pat string, isEof bool, opt []string, line_no int, file_name string) {
	pp := 1 // Pattern
	if isEof {
		pp = 2 // EOF
	}
	// fmt.Printf("opt: %v\n", opt)
	x := &ImRuleType{
		Pattern:     pat,
		PatternType: pp,
		LineNo:      line_no, // Error Reporintg Stuff ------------------------------------------------------------------------
		FileName:    file_name,
	}
	for ii, vv := range opt {
		_ = ii
		nm, param := ParseActionItem(vv)
		// fmt.Printf("opt[%d] nm %s param >%s<\n", ii, nm, param)
		switch nm {
		case "Rv":
			x.RvName = param
			Im.SaveDef("Tokens", []string{param}, line_no, file_name)
		case "Call":
			x.CallName = param
			Im.SaveDef("Machines", []string{param}, line_no, file_name)
		case "Repl":
			x.Repl = true
			x.ReplString = param
		case "Ignore":
			x.Ignore = true
			x.RvName = "Tok_Ignore"
		case "NotGreedy":
			x.NotGreedy = true
		case "Error":
			x.Err = true
			x.WEString = param
			Im.SaveDef("Errors", []string{param}, line_no, file_name)
		case "ReservedWord":
			x.ReservedWord = true
			Im.SaveDef("ReservedWords", []string{param}, line_no, file_name)
		case "Return":
			x.Return = true
		case "Warn":
			x.Warn = true
			x.WEString = param
			Im.SaveDef("Errors", []string{param}, line_no, file_name)
		case "Reset":
			// TODO xyzzy - not implemented yet
		default:
			fmt.Printf("Error: %s is not a defined operation, line %d file %s\n", nm, line_no, file_name)
		}
	}
	Im.Machine[MNo].Rules = append(Im.Machine[MNo].Rules, x)
}

func (Im *ImType) SaveMachine(opt []string) int {
	ap := len(Im.Machine)
	Mt := ImMachineType{
		Name:   opt[0],
		Mixins: opt[1:],
		Defs:   &Im.Def,
	}
	Im.Machine = append(Im.Machine, Mt)
	Im.Def.DefsAre["Machines"].NameValueStr[opt[0]] = fmt.Sprintf("%d", ap)
	return ap
}

/*
	type ImSeenAtType struct {
		LineNo   []int    //
		FileName []string //
	}

	type ImDefinedValueType struct {
		Seq       int                     //
		WhoAmI    string                  //
		NameValue map[string]int          //
		Reverse   map[int]string          //
		SeenAt    map[string]ImSeenAtType //
	}

	type ImDefsType struct {
		DefsAre map[string]ImDefinedValueType //
	}

	type ImDefinedValueType struct {
		Seq          int                     //
		WhoAmI       string                  //
		NameValueStr map[string]string       //
		NameValue    map[string]int          //
		Reverse      map[int]string          //
		SeenAt       map[string]ImSeenAtType //
	}
*/
func (Im *ImType) FindValueFor(t string) int {
	//for s, dd := range Im.Def.DefsAre {
	//	_ = s
	for _, DefType := range []string{"Machines", "Errors", "ReservedWords", "Tokens"} {
		dd := Im.Def.DefsAre[DefType]
		// fmt.Printf("In %s Looking for %s\n", s, t)
		if v, ok := dd.NameValue[t]; ok {
			// fmt.Printf("In %s Found for %s=%d\n", s, t, v)
			return v
		}
	}
	return -1
}

func (Im *ImType) Lookup(DefType string, t string) int {
	if validateDefType(DefType) {
		dd := Im.Def.DefsAre[DefType]
		// fmt.Printf("In %s Looking for %s\n", s, t)
		if v, ok := dd.NameValue[t]; ok {
			// fmt.Printf("In %s Found for %s=%d\n", s, t, v)
			return v
		}
	}
	return -1
}

// -----------------------------------------------------------------------------------------------------------------------------
// -----------------------------------------------------------------------------------------------------------------------------
//
//	type ImRuleType struct {
//		PatternType  int    // Pattern, Str0,1,2, $eof etc.  // Pattern Stuff --------------------------------------------------------------------------------
//		Pattern      string //
func (Im *ImType) LocatePattern(ff *ImRuleType, in []*ImRuleType) (rv int) {
	rv = -1
	// fmt.Printf("LocatePattern for %s %d\n", ff.Pattern, ff.PatternType)
	for kk, tt := range in {
		// fmt.Printf("    Compare to %s %d\n", tt.Pattern, tt.PatternType)
		if tt.PatternType == ff.PatternType && tt.Pattern == ff.Pattern {
			dbgo.DbPrintf("in", "    Found\n")
			return kk
		}
	}
	// fmt.Printf("    NOT NOT NOT Found\n")
	// Xyzzy - should add an error at this point?
	//		2. Errors about missing machines in mixin not reported at all - see xyzzyMixin01
	return
}

// -----------------------------------------------------------------------------------------------------------------------------
// -----------------------------------------------------------------------------------------------------------------------------
func (Im *ImType) FinializeFile() {
	ADef := ""
	AKey := ""

	// for DefType, dd := range Im.Def.DefsAre {
	for _, DefType := range []string{"Machines", "Errors", "Tokens", "ReservedWords"} {
		dd := Im.Def.DefsAre[DefType]
		// fmt.Printf("DefType (FinializeFiile): %s\n", DefType)
		// Pass 1 - Take Numbers and put in
		ss := com.SortMapStringString(dd.NameValueStr)
		// for AKey, ADef := range dd.NameValueStr {
		for _, AKey = range ss {
			ADef = dd.NameValueStr[AKey]
			if len(ADef) > 0 && IsNumeric(ADef) {
				// fmt.Printf("Found numeric for %s=%s\n", AKey, ADef)
				v, err := strconv.Atoi(ADef)
				if err != nil {
					fmt.Printf("Error: Invalid numeric value for a token, >%s<, error=%s\n", ADef, err)
				} else {
					dd.NameValue[AKey] = v
					dd.Reverse[v] = AKey
				}
			}
		}
		// Pass 2 - Assign All Others
		seq := dd.Seq
		// for AKey, ADef := range dd.NameValueStr {
		for _, AKey = range ss {
			ADef = dd.NameValueStr[AKey]
			if len(ADef) > 0 && IsNumeric(ADef) {
			} else if len(ADef) == 0 {
				// fmt.Printf("Found sequence assign for for %s=%s\n", AKey, ADef)
				for _, ok := dd.Reverse[seq]; ok; {
					seq++
					_, ok = dd.Reverse[seq]
				}
				dd.NameValue[AKey] = seq
				dd.Reverse[seq] = AKey
				dd.NameValueStr[AKey] = fmt.Sprintf("%d", seq)
				// fmt.Printf("     Assigned %d\n", seq)
				seq++
			}
		}
		dd.Seq = seq

		// fmt.Printf("dd.NameValue = %v\n", dd.NameValue)

		// Pass 3 - Assign Tokens
		// for AKey, ADef := range dd.NameValueStr {
		for _, AKey = range ss {
			ADef = dd.NameValueStr[AKey]
			if len(ADef) > 0 && IsNumeric(ADef) {
			} else if len(ADef) > 0 {
				// fmt.Printf("Found Name Assign for for %s=%s\n", AKey, ADef)
				if v, ok := dd.NameValue[ADef]; ok {
					dd.NameValue[AKey] = v
					dd.Reverse[v] = AKey
				}
			}
		}
		// Pass 4 - Look for any unsigned
		// for AKey, ADef := range dd.NameValueStr {
		for _, AKey = range ss {
			ADef = dd.NameValueStr[AKey]
			if len(ADef) > 0 && IsNumeric(ADef) {
			} else if len(ADef) > 0 {
				if v, ok := dd.NameValue[ADef]; ok {
					dd.NameValue[AKey] = v
					dd.Reverse[v] = AKey
				} else {
					v := Im.FindValueFor(ADef)
					if v < 0 {
						fmt.Printf("Warning: Token is not defined, Automatically defining!, ADef/AKey %s=%s=%d\n", ADef, AKey, seq) // !! !! Requries a 4th pass - after all defined !! !!
						dd.NameValue[AKey] = seq
						dd.Reverse[seq] = AKey
						seq++
					} else {
						dd.NameValue[AKey] = v
						dd.Reverse[v] = AKey
					}
				}
			}
		}
		Im.Def.DefsAre[DefType] = dd
	}

	for _, vv := range Im.Machine {
		for _, ww := range vv.Rules {
			if len(ww.RvName) > 0 {
				// fmt.Printf("%-20s", fmt.Sprintf(" Rv:%d=%s ", ww.Rv, ww.RvName))
				ww.Rv = Im.Lookup("Tokens", ww.RvName)
			}
			if len(ww.CallName) > 0 {
				// fmt.Printf("%-20s", fmt.Sprintf(" Call:%d=%s ", ww.Call, ww.CallName))
				ww.Call = Im.Lookup("Machines", ww.CallName)
			}
		}
	}

	// xyzzy-Machine--Mixins---
	for kk, vv := range Im.Machine {
		var tRules []*ImRuleType
		tRules = make([]*ImRuleType, 0, 100)
		for _, rr := range vv.Rules {
			p := Im.LocatePattern(rr, tRules) // A merge operation - if not found then append, else replace
			if p >= 0 {
				tRules[p] = rr
			} else {
				tRules = append(tRules, rr)
			}
		}
		for _, ww := range vv.Mixins {
			ii := Im.Lookup("Machines", ww)
			if ii >= 0 && ii < len(Im.Machine) {
				for _, uu := range Im.Machine[ii].Rules {
					p := Im.LocatePattern(uu, tRules) // A merge operation - if not found then append, else replace
					if p < 0 {
						tRules = append(tRules, uu)
					}
				}
			} else {
				fmt.Printf("Error: Mixin - did not find %s as a machine name\n", ww)
			}
		}
		Im.Machine[kk].Rules = tRules
	}

	dd := Im.Def.DefsAre["Tokens"]
	Tok_map = dd.Reverse
}

func (Im *ImType) OutputDef() {
	fmt.Printf("Defs - OutputDef\n")
	for _, DefType := range []string{"Machines", "Errors", "ReservedWords", "Tokens"} {
		dd := Im.Def.DefsAre[DefType]
		fmt.Printf("==========================================================================\n")
		fmt.Printf("DefType: %s\n", DefType)
		fmt.Printf("==========================================================================\n")

		ss := com.SortMapStringString(dd.NameValueStr)
		// for AKey, ADef := range dd.NameValueStr {
		for _, AKey := range ss {
			ADef := dd.NameValueStr[AKey]
			fmt.Printf("    %s=%v\n", AKey, ADef)
		}
	}
}

// min, max := com.RangeOfIntKeys(dd.Reverse)
func RangeOfIntKeys(x map[int]string) (min int, max int) {
	init := true
	for ii := range x {
		if init {
			init = false
			min, max = ii, ii
		} else {
			if ii > max {
				max = ii
			}
			if ii < min {
				min = ii
			}
		}
	}
	return
}

func (Im *ImType) OutputDefAsGoCode(fo io.Writer) {
	fmt.Fprintf(fo, "\n// Defs - OutputDef\n")
	// for _, DefType := range []string{"Machines", "Errors", "ReservedWords", "Tokens"} {
	for _, DefType := range []string{"Tokens", "Machines", "Errors", "ReservedWords"} {
		dd := Im.Def.DefsAre[DefType]
		fmt.Fprintf(fo, "// ==========================================================================\n")
		fmt.Fprintf(fo, "// DefType: %s\n", DefType)
		fmt.Fprintf(fo, "// ==========================================================================\n")
		fmt.Fprintf(fo, "const (\n")

		min, max := RangeOfIntKeys(dd.Reverse)
		// for AKey, ADef := range dd.NameValueStr {
		for ii := min; ii < max; ii++ {
			if AKey, ok := dd.Reverse[ii]; ok {
				ADef := ii
				if DefType == "ReservedWords" {
					fmt.Fprintf(fo, "    RW_%s = %v\n", AKey, ADef)
				} else {
					fmt.Fprintf(fo, "    %s = %v\n", AKey, ADef)
				}
			}
		}

		fmt.Fprintf(fo, ")\n\n")
	}
}

// Output the Im structure
func (Im *ImType) OutputImType() {

	dpt := []string{"???", "Pat", "EOF", "???"}
	if dbgo.IsDbOn("in-echo-machine") {

		Im.OutputDef()
		Im.OutputDefAsGoCode(os.Stdout)
		for ii, vv := range Im.Machine {
			fmt.Printf("Machine[%d] Name[%s]-----------------------------------------------------------------\n", ii, vv.Name)
			fmt.Printf("    Mixins: %v\n", vv.Mixins)
			// Rules  []*ImRuleType //
			for jj, ww := range vv.Rules {
				s := fmt.Sprintf("%q", ww.Pattern)
				s = s[1:]
				s = s[0 : len(s)-1]
				fmt.Printf("      %3d: %3s   %-30s ", jj, dpt[ww.PatternType], s)
				if len(ww.RvName) > 0 {
					fmt.Printf("%-20s", fmt.Sprintf(" Rv:%d=%s ", ww.Rv, ww.RvName))
				} else {
					fmt.Printf("%-20s", "")
				}
				if len(ww.CallName) > 0 {
					fmt.Printf("%-20s", fmt.Sprintf(" Call:%d=%s ", ww.Call, ww.CallName))
				} else {
					fmt.Printf("%-20s", "")
				}
				if ww.Return {
					fmt.Printf(" Return ")
				}
				if ww.Repl {
					fmt.Printf(" Repl:%s ", ww.ReplString)
				}
				if ww.Ignore {
					fmt.Printf(" [Ignore] ")
				}
				if ww.ReservedWord {
					fmt.Printf(" [ReservedWord] ")
				}
				if ww.Err {
					fmt.Printf(" [Err=%s] ", ww.WEString)
				}
				if ww.Warn {
					fmt.Printf(" [Warn=%s] ", ww.WEString)
				}
				fmt.Printf("\n")
			}
		}
	}

}

func (Im *ImType) LookupMachine(name string) int {
	x := Im.Lookup("Machines", name)
	return x
}

func ImReadFile(fn string) (Im *ImType) {
	Im = NewIm()
	fd := ReadFileIntoLines(fn)
	Im.InputLines = fd
	if len(fd) > 0 {
		Im.ParseFile(fd)
	}
	// fmt.Printf("%+v\n", Im)
	Im.OutputImType()
	return
}

func Lookup_Tok_Name(Tok int) (rv string) {
	ok := false
	if rv, ok = Tok_map[Tok]; ok {
		return
	}
	rv = fmt.Sprintf("Unknown(Tok=%d)", Tok)
	return
}

func Add_Lookup_Token(Tok int, Name string) {
	Tok_map[Tok] = Name
}

/* vim: set noai ts=4 sw=4: */
