package re

//
// 0. Check that the Sigma is correct
//
// 1. String function that outputs RE parse trees
// 2. Compare these with ref-parse trees
// 3. Hand check that they are OK
//

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/pschlump/dbgo"
	"github.com/pschlump/lexie/com"

	. "gopkg.in/check.v1"
)

// https://labix.org/gocheck

type TokByte struct {
	Tok LR_TokType // Node Type
	Dat []byte
}

type Lexie01DataType struct {
	Test         string
	Re           string
	Rv           int
	NExpectedErr int
	SkipTest     bool
	ELen         int
	SM           []TokByte
	Sigma        string
}

// Test of lr.Next() the RE Scanner -

var Lexie00Data = []Lexie01DataType{
	{Test: "0000", SkipTest: true}, //
	{Test: "0001", Re: "(x|y)*abb", Rv: 0001, SkipTest: false, ELen: 3,
		SM: []TokByte{
			TokByte{LR_OP_PAR, []byte{0x28}},
			TokByte{LR_Text, []byte{0x78}},
			TokByte{LR_OR, []byte{0x7c}},
			TokByte{LR_Text, []byte{0x79}},
			TokByte{LR_CL_PAR, []byte{0x29}},
			TokByte{LR_STAR, []byte{0x2a}},
			TokByte{LR_Text, []byte{0x61}},
			TokByte{LR_Text, []byte{0x62}},
			TokByte{LR_Text, []byte{0x62}},
		}},
	{Test: "0002", Re: "(\u2022|\u2318)+abb", Rv: 0002, SkipTest: false, ELen: 4,
		SM: []TokByte{
			TokByte{LR_OP_PAR, []byte{0x28}},
			TokByte{LR_Text, []byte{0xe2, 0x80, 0xa2}},
			TokByte{LR_OR, []byte{0x7c}},
			TokByte{LR_Text, []byte{0xe2, 0x8c, 0x98}},
			TokByte{LR_CL_PAR, []byte{0x29}},
			TokByte{LR_PLUS, []byte{0x2b}},
			TokByte{LR_Text, []byte{0x61}},
			TokByte{LR_Text, []byte{0x62}},
			TokByte{LR_Text, []byte{0x62}},
		}},
	{Test: "0003", Re: "(\u0428|\u0496|\u044a)+abb", Rv: 0003, SkipTest: false, ELen: 4}, // Len(4)	 // The Teepee is a lower case letter
	{Test: "0004", Re: "[abc]", Rv: 0004, SkipTest: false},
	{Test: "0005", Re: "[0[:alpha:]9]", Rv: 0004, SkipTest: false},
	{Test: "0006", Re: "c{2,3}", Rv: 0004, SkipTest: false},
	{Test: "0007", Re: "abc{2,3}", Rv: 0004, SkipTest: false},
	{Test: "0008", Re: "abc{,3}", Rv: 0004, SkipTest: false},
	{Test: "0009", Re: "abc{2}", Rv: 0004, SkipTest: false},
	{Test: "0010", Re: "abc{2,}", Rv: 0004, SkipTest: false},
	{Test: "0011", Re: `abc\{\*\(\|\\\[\]def`, Rv: 0004, SkipTest: false},
}

// Test of Parseing REs into RE-ParseTrees

var Lexie01Data = []Lexie01DataType{
	{Test: "1000", Re: "(x|y)*abb", Rv: 1000, SkipTest: false, ELen: 3},                // Len(3)
	{Test: "1001", Re: "x*", Rv: 1001, SkipTest: false, ELen: 0},                       // Len(0)
	{Test: "1002", Re: "(xx)*", Rv: 1002, SkipTest: false, ELen: 0},                    // Len(0)
	{Test: "1003", Re: "(xx)+", Rv: 1003, SkipTest: false, ELen: 2},                    // Len(2)
	{Test: "1004", Re: "(xx)?", Rv: 1004, SkipTest: false, ELen: 0},                    // Len(0)
	{Test: "1005", Re: "(a|b)", Rv: 1005, SkipTest: false, ELen: 1},                    // Len(Min(len(1),Len(1)) = Len(1)
	{Test: "1006", Re: "(aa|bb)", Rv: 1006, SkipTest: false, ELen: 2},                  // Len(2)
	{Test: "1007", Re: "(a|b)*abb", Rv: 1007, SkipTest: false, ELen: 3},                // Len(3) Examle from Dragon Compiler Book and .pdf files
	{Test: "1008", Re: "(aa|bb|ccc)*abb", Rv: 1008, SkipTest: false, ELen: 3},          // Len(3)
	{Test: "1009", Re: "^abb$", Rv: 1009, SkipTest: false, ELen: 3},                    // Len(3)+Hard
	{Test: "1010", Re: "^abb", Rv: 1010, SkipTest: false, ELen: 3},                     // Len(3)+Hard
	{Test: "1011", Re: `a(bcd)*(ghi)+(jkl)*X`, Rv: 1011, SkipTest: false, ELen: 5},     // Len(1+3+1)
	{Test: "1012", Re: `a[.]d`, Rv: 1012, SkipTest: false, ELen: 3},                    // Len(3)
	{Test: "1013", Re: `a[^]d`, Rv: 1013, SkipTest: false, ELen: 0},                    // Len(?) TODO: -- Sigma should have an X_N_CCL char in it - missing
	{Test: "1014", Re: `a(def)*(klm(mno)+)?b`, Rv: 1014, SkipTest: false, ELen: 2},     // Len(2)
	{Test: "1015", Re: `a[a-zA-Z_][a-zA-Z_0-9]*d`, Rv: 1015, SkipTest: false, ELen: 3}, // Len(3)
	{Test: "1016", Re: `a.d`, Rv: 1016, SkipTest: false, ELen: 3},                      // Len(3)
	{Test: "1017", Re: "(aa|bb|ccc)abb", Rv: 1017, SkipTest: false, ELen: 5},           // Len(2+3=5)
	{Test: "1018", Re: "(||)", Rv: 1018, SkipTest: false, ELen: 0},                     // Len(0)
	{Test: "1019", Re: "||", Rv: 1019, SkipTest: false, ELen: 0},                       // Len(0)
	{Test: "1020", Re: "(||||||||||||||)", Rv: 1020, SkipTest: false, ELen: 0},         // Len(0)
	{Test: "1021", Re: "(||||||||a||||||)", Rv: 1021, SkipTest: false, ELen: 0},        // Len(0)
	{Test: "1022", Re: "(||||||||a|aa|||||)", Rv: 1022, SkipTest: false, ELen: 0},      // Len(0)
	{Test: "1023", Re: "(a|aa|aaa)", Rv: 1023, SkipTest: false, ELen: 1},               // Len(1)
	{Test: "1024", Re: "(ab|aab|aaab)", Rv: 1024, SkipTest: false, ELen: 2},            // Len(2)
	{Test: "1025", Re: "(ab|aab|aaab)c", Rv: 1025, SkipTest: false, ELen: 3},           // Len(3)
	{Test: "1026", Re: "(a*|aab|aaab)", Rv: 1026, SkipTest: false, ELen: 0},            // Len(0)
	{Test: "1027", Re: "(.*|-=-|-=#)", Rv: 1027, SkipTest: false, ELen: 0},             // Len(0)

	//	wabi-sabi (侘寂?)
	{Test: "1028", Re: "(\u0428|\u0496|\u044a)+abb", Rv: 1028, SkipTest: false, ELen: 4, Sigma: "abШъҖ"}, // Len(4)	 // The Teepee is a lower case letter
	{Test: "1029", Re: "(a\u03bbb|a\u0428b|aaab)", Rv: 1029, SkipTest: false, ELen: 2, Sigma: "abλШ"},    // Len(2)
}

type Test6DataType struct {
	Re           string
	TopTok       []LR_TokType
	TopVal       []string
	NExpectedErr int
}

// xyzzy - Add ability to check value of 1st child (recursive), 2nd, 3rd ... - and validate what it looks like

var Test6Data = []Test6DataType{
	/* 000 */ {Re: `ab*c`, TopTok: []LR_TokType{LR_Text, LR_STAR, LR_Text}},
	/* 001 */ {Re: `ab+c`, TopTok: []LR_TokType{LR_Text, LR_PLUS, LR_Text}},
	/* 002 */ {Re: `x[aeiou]y`, TopTok: []LR_TokType{LR_Text, LR_CCL, LR_Text}, TopVal: []string{"x", "aeiou", "y"}},
	/* 003 */ {Re: `x[^aeiou]y`, TopTok: []LR_TokType{LR_Text, LR_N_CCL, LR_Text}, TopVal: []string{"x", "aeiou", "y"}},
	/* 004 */ {Re: `x[a-c]y`, TopTok: []LR_TokType{LR_Text, LR_CCL, LR_Text}, TopVal: []string{"x", "abc", "y"}},
	/* 005 */ {Re: `x[c-a]y`, TopTok: []LR_TokType{LR_Text, LR_CCL, LR_Text}, TopVal: []string{"x", "", "y"}},
	/* 006 */ {Re: `x[a-dA-D]y`, TopTok: []LR_TokType{LR_Text, LR_CCL, LR_Text}, TopVal: []string{"x", "abcdABCD", "y"}},
	/* 007 */ {Re: `x[a-dA-D_]y`, TopTok: []LR_TokType{LR_Text, LR_CCL, LR_Text}, TopVal: []string{"x", "abcdABCD_", "y"}},
	/* 008 */ {Re: `x[_a-dA-D_]y`, TopTok: []LR_TokType{LR_Text, LR_CCL, LR_Text}, TopVal: []string{"x", "_abcdABCD_", "y"}},
	/* 009 */ {Re: `x[_a-dA-D]y`, TopTok: []LR_TokType{LR_Text, LR_CCL, LR_Text}, TopVal: []string{"x", "_abcdABCD", "y"}},

	/* 010 */ {Re: `x[_a-dA-D0-9]y`, TopTok: []LR_TokType{LR_Text, LR_CCL, LR_Text}, TopVal: []string{"x", "_abcdABCD", "y"}},
	/* 011 */ {Re: `x[0-8]y`, TopTok: []LR_TokType{LR_Text, LR_CCL, LR_Text}, TopVal: []string{"x", "012345678", "y"}},
	/* 012 */ {Re: `x[-0-8]y`, TopTok: []LR_TokType{LR_Text, LR_CCL, LR_Text}, TopVal: []string{"x", "-012345678", "y"}},
	/* 013 */ {Re: `x[-x-z]y`, TopTok: []LR_TokType{LR_Text, LR_CCL, LR_Text}, TopVal: []string{"x", "-xyz", "y"}},
	/* 014 */ {Re: `x[^-x-z]y`, TopTok: []LR_TokType{LR_Text, LR_N_CCL, LR_Text}, TopVal: []string{"x", "-xyz", "y"}},
	/* 015 */ {Re: `x[x-z]*y`, TopTok: []LR_TokType{LR_Text, LR_STAR, LR_Text}, TopVal: []string{"x", "*", "y"}},
	/* 016 */ {Re: `x[-x-z]*y`, TopTok: []LR_TokType{LR_Text, LR_STAR, LR_Text}, TopVal: []string{"x", "*", "y"}},
	/* 017 */ {Re: `x[z-x]*y`, TopTok: []LR_TokType{LR_Text, LR_STAR, LR_Text}, TopVal: []string{"x", "*", "y"}},
	/* 018 */ {Re: `x[a-z]*y`, TopTok: []LR_TokType{LR_Text, LR_STAR, LR_Text}, TopVal: []string{"x", "*", "y"}},
	/* 019 */ {Re: `x[a-zA-Z]*y`, TopTok: []LR_TokType{LR_Text, LR_STAR, LR_Text}, TopVal: []string{"x", "*", "y"}},

	/* 020 */ {Re: `x[0-9]*y`, TopTok: []LR_TokType{LR_Text, LR_STAR, LR_Text}, TopVal: []string{"x", "*", "y"}},
	/* 021 */ {Re: `x[a-zA-Z0-9]*y`, TopTok: []LR_TokType{LR_Text, LR_STAR, LR_Text}, TopVal: []string{"x", "*", "y"}},
	/* 022 */ {Re: `x[0-9]+y`, TopTok: []LR_TokType{LR_Text, LR_PLUS, LR_Text}, TopVal: []string{"x", "+", "y"}},
	/* 023 */ {Re: `x[0-9]?y`, TopTok: []LR_TokType{LR_Text, LR_QUEST, LR_Text}, TopVal: []string{"x", "?", "y"}},
	/* 024 */ {Re: `x[^0-9]*y`, TopTok: []LR_TokType{LR_Text, LR_STAR, LR_Text}, TopVal: []string{"x", "*", "y"}},
	/* 025 */ {Re: `x[0-8]\?y`, TopTok: []LR_TokType{LR_Text, LR_CCL, LR_Text, LR_Text}, TopVal: []string{"x", "012345678", "?", "y"}},
	/* 026 */ {Re: `x[0-8]\+y`, TopTok: []LR_TokType{LR_Text, LR_CCL, LR_Text, LR_Text}, TopVal: []string{"x", "012345678", "+", "y"}},
	/* 027 */ {Re: `x[0-8]\*y`, TopTok: []LR_TokType{LR_Text, LR_CCL, LR_Text, LR_Text}, TopVal: []string{"x", "012345678", "*", "y"}},
	/* 028 */ {Re: `x[1-9*]\*y`, TopTok: []LR_TokType{LR_Text, LR_CCL, LR_Text, LR_Text}, TopVal: []string{"x", "123456789*", "*", "y"}},
	/* 029 */ {Re: `x[1-9*+?]\*y`, TopTok: []LR_TokType{LR_Text, LR_CCL, LR_Text, LR_Text}, TopVal: []string{"x", "123456789*+?", "*", "y"}},

	/* 030 */ {Re: `)c`, TopTok: []LR_TokType{LR_Text, LR_Text}, TopVal: []string{")", "c"}, NExpectedErr: 1},
	/* 031 */ {Re: `?c`, TopTok: []LR_TokType{LR_Text, LR_Text}, TopVal: []string{"?", "c"}, NExpectedErr: 1},
	/* 032 */ {Re: `c[`, TopTok: []LR_TokType{LR_Text, LR_CCL}, TopVal: []string{"c", ""}, NExpectedErr: 1},
	/* 033 */ {Re: `c]d`, TopTok: []LR_TokType{LR_Text, LR_Text, LR_Text}, TopVal: []string{"c", "]", "d"}, NExpectedErr: 1},
	/* 034 */ {Re: `*c`, TopTok: []LR_TokType{LR_Text, LR_Text}, TopVal: []string{"*", "c"}, NExpectedErr: 1},
	/* 035 */ {Re: `+c`, TopTok: []LR_TokType{LR_Text, LR_Text}, TopVal: []string{"+", "c"}, NExpectedErr: 1},
	/* 036 */ {Re: `.X`, TopTok: []LR_TokType{LR_DOT, LR_Text}, TopVal: []string{".", "X"}, NExpectedErr: 0},
	/* 037 */ {Re: `.*X`, TopTok: []LR_TokType{LR_STAR, LR_Text}, TopVal: []string{"*", "X"}, NExpectedErr: 0},
	/* 038 */ {Re: `a.*`, TopTok: []LR_TokType{LR_Text, LR_STAR}, TopVal: []string{"a", "*"}, NExpectedErr: 0},
	/* 039 */ {Re: `a.*X`, TopTok: []LR_TokType{LR_Text, LR_STAR, LR_Text}, TopVal: []string{"a", "*", "X"}, NExpectedErr: 0},

	/* 040 */ {Re: `a(bcd)X`, TopTok: []LR_TokType{LR_Text, LR_OP_PAR, LR_Text}, TopVal: []string{"a", "(", "X"}, NExpectedErr: 0},
	/* 041 */ {Re: `a(bcd)*X`, TopTok: []LR_TokType{LR_Text, LR_STAR, LR_Text}, TopVal: []string{"a", "*", "X"}, NExpectedErr: 0},
	/* 042 */ {Re: `a(bcd)+X`, TopTok: []LR_TokType{LR_Text, LR_PLUS, LR_Text}, TopVal: []string{"a", "+", "X"}, NExpectedErr: 0},
	/* 043 */ {Re: `a(bcd)?X`, TopTok: []LR_TokType{LR_Text, LR_QUEST, LR_Text}, TopVal: []string{"a", "?", "X"}, NExpectedErr: 0},
	/* 044 */ {Re: `a(b[0-9]d)?X`, TopTok: []LR_TokType{LR_Text, LR_QUEST, LR_Text}, TopVal: []string{"a", "?", "X"}, NExpectedErr: 0},
	/* 045 */ {Re: `a(bcd)*(ghi)+X`, TopTok: []LR_TokType{LR_Text, LR_STAR, LR_PLUS, LR_Text}, TopVal: []string{"a", "*", "+", "X"}, NExpectedErr: 0},
	/* 046 */ {Re: `a(bcd)*(ghi)+(jkl)*X`, TopTok: []LR_TokType{LR_Text, LR_STAR, LR_PLUS, LR_STAR, LR_Text}, TopVal: []string{"a", "*", "+", "*", "X"}, NExpectedErr: 0},
	/* 047 */ {Re: `a(bcd`, TopTok: []LR_TokType{LR_Text, LR_OP_PAR}, TopVal: []string{"a", "("}, NExpectedErr: 0},
	/* 048 */ {Re: `a(`, TopTok: []LR_TokType{LR_Text, LR_OP_PAR}, TopVal: []string{"a", "("}, NExpectedErr: 0},
	/* 049 */ {Re: `a()`, TopTok: []LR_TokType{LR_Text, LR_OP_PAR}, TopVal: []string{"a", "("}, NExpectedErr: 0},

	/* 050 */ {Re: `a[bbbCCC]d`, TopTok: []LR_TokType{LR_Text, LR_CCL, LR_Text}, TopVal: []string{"a", "bbbCCC", "d"}, NExpectedErr: 0},
	/* 051 */ {Re: `a[^]d`, TopTok: []LR_TokType{LR_Text, LR_N_CCL, LR_Text}, TopVal: []string{"a", "", "d"}, NExpectedErr: 0},
	/* 052 */ {Re: `a[.]d`, TopTok: []LR_TokType{LR_Text, LR_CCL, LR_Text}, TopVal: []string{"a", ".", "d"}, NExpectedErr: 0},
	/* 053 */ {Re: `a[^.]d`, TopTok: []LR_TokType{LR_Text, LR_N_CCL, LR_Text}, TopVal: []string{"a", ".", "d"}, NExpectedErr: 0},
	/* 054 */ {Re: `a\.d`, TopTok: []LR_TokType{LR_Text, LR_Text, LR_Text}, TopVal: []string{"a", ".", "d"}, NExpectedErr: 0},
	/* 055 */ {Re: `a\.\.d`, TopTok: []LR_TokType{LR_Text, LR_Text, LR_Text, LR_Text}, TopVal: []string{"a", ".", ".", "d"}, NExpectedErr: 0},
	/* 056 */ {Re: `a\.\[d`, TopTok: []LR_TokType{LR_Text, LR_Text, LR_Text, LR_Text}, TopVal: []string{"a", ".", "[", "d"}, NExpectedErr: 0},
	/* 057 */ {Re: `a\.\]d`, TopTok: []LR_TokType{LR_Text, LR_Text, LR_Text, LR_Text}, TopVal: []string{"a", ".", "]", "d"}, NExpectedErr: 0},
	/* 058 */ {Re: `a\.\(d`, TopTok: []LR_TokType{LR_Text, LR_Text, LR_Text, LR_Text}, TopVal: []string{"a", ".", "(", "d"}, NExpectedErr: 0},
	/* 059 */ {Re: `a\.\(\)d`, TopTok: []LR_TokType{LR_Text, LR_Text, LR_Text, LR_Text, LR_Text}, TopVal: []string{"a", ".", "(", ")", "d"}, NExpectedErr: 0},

	/* 060 */ {Re: ``, TopTok: []LR_TokType{}, TopVal: []string{}, NExpectedErr: 0},
	/* 061 */ {Re: `^a$`, TopTok: []LR_TokType{LR_CARROT, LR_Text, LR_DOLLAR}, TopVal: []string{"^", "a", "$"}, NExpectedErr: 0},
	/* 062 */ {Re: `^a`, TopTok: []LR_TokType{LR_CARROT, LR_Text}, TopVal: []string{"^", "a"}, NExpectedErr: 0},
	/* 063 */ {Re: `a$`, TopTok: []LR_TokType{LR_Text, LR_DOLLAR}, TopVal: []string{"a", "$"}, NExpectedErr: 0},
	/* 064 */ {Re: `$`, TopTok: []LR_TokType{LR_DOLLAR}, TopVal: []string{"$"}, NExpectedErr: 0},
	/* 065 */ {Re: `^`, TopTok: []LR_TokType{LR_CARROT}, TopVal: []string{"^"}, NExpectedErr: 0},

	// Verify that special chars can be escaped as text
	/* 066 */ {Re: `c\^\$\.\(\)\[\]\[a-z\]\\\|x`,
		TopTok:       []LR_TokType{LR_Text, LR_Text, LR_Text, LR_Text, LR_Text, LR_Text, LR_Text, LR_Text, LR_Text, LR_Text, LR_Text, LR_Text, LR_Text, LR_Text, LR_Text, LR_Text},
		TopVal:       []string{"c", "^", "$", ".", "(", ")", "[", "]", "[", "a", "-", "z", "]", "\\", "|", "x"},
		NExpectedErr: 0},

	// Verify special chars are not special inside a CCL
	/* 067 */ {Re: `a[-x|$^[.]d`, TopTok: []LR_TokType{LR_Text, LR_CCL, LR_Text}, TopVal: []string{"a", "-x|$^[.", "d"}, NExpectedErr: 0},

	// Verify escape of END CCL inside CCL
	/* 068 */ {Re: `a[-x|$^[.\]]d`, TopTok: []LR_TokType{LR_Text, LR_CCL, LR_Text}, TopVal: []string{"a", "-x|$^[.]", "d"}, NExpectedErr: 0},

	// Nested Parens ------------------------------------------------------------------------------------
	/* 069 */ {Re: `a(b)`, TopTok: []LR_TokType{LR_Text, LR_OP_PAR}, TopVal: []string{"a", "("}, NExpectedErr: 0},
	/* 070 */ {Re: `a(b)?`, TopTok: []LR_TokType{LR_Text, LR_QUEST}, TopVal: []string{"a", "?"}, NExpectedErr: 0},
	/* 071 */ {Re: `a(bbb)?`, TopTok: []LR_TokType{LR_Text, LR_QUEST}, TopVal: []string{"a", "?"}, NExpectedErr: 0},
	/* 072 */ {Re: `a(bbb)?e`, TopTok: []LR_TokType{LR_Text, LR_QUEST, LR_Text}, TopVal: []string{"a", "?", "e"}, NExpectedErr: 0},
	/* 073 */ {Re: `a(b(ddd)?e)*X`, TopTok: []LR_TokType{LR_Text, LR_STAR, LR_Text}, TopVal: []string{"a", "*", "X"}, NExpectedErr: 0},
	/* 074 */ {Re: `a(def)*(klm(mno)+)?b`, TopTok: []LR_TokType{LR_Text, LR_STAR, LR_QUEST, LR_Text}, TopVal: []string{"a", "*", "?", "b"}, NExpectedErr: 0},
	/* 075 */ {Re: `a(())`, TopTok: []LR_TokType{LR_Text, LR_OP_PAR}, TopVal: []string{"a", "("}, NExpectedErr: 0}, // Empty Paren ok - Matches 0 chars
	/* 076 */ {Re: `a((`, TopTok: []LR_TokType{LR_Text, LR_OP_PAR}, TopVal: []string{"a", "("}, NExpectedErr: 0}, // Parens auto-close at EOF
	/* 077 */ {Re: `a(()`, TopTok: []LR_TokType{LR_Text, LR_OP_PAR}, TopVal: []string{"a", "("}, NExpectedErr: 0}, // Parens auto-close at EOF

	// Involving OR == | ------------------------------------------------------------------------------------
	/* 078 */ {Re: `abc`, TopTok: []LR_TokType{LR_Text, LR_Text, LR_Text}, TopVal: []string{"a", "b", "c"}, NExpectedErr: 0},
	/* 079 */ {Re: `aaa|bbb`, TopTok: []LR_TokType{LR_OR}, TopVal: []string{"|"}, NExpectedErr: 0},
	/* 080 */ {Re: `(aaa|bbb)`, TopTok: []LR_TokType{LR_OP_PAR}, TopVal: []string{"("}, NExpectedErr: 0},
	/* 081 */ {Re: `(aaa|bbb)*`, TopTok: []LR_TokType{LR_STAR}, TopVal: []string{"*"}, NExpectedErr: 0},
	/* 082 */ {Re: `a|b|c`, TopTok: []LR_TokType{LR_OR}, TopVal: []string{"|"}, NExpectedErr: 0},
	/* 083 */ {Re: `aaa|bbb|ccc`, TopTok: []LR_TokType{LR_OR}, TopVal: []string{"|"}, NExpectedErr: 0},
	/* 084 */ {Re: `aa|bb|cc|dd|ee|ff|gg|hh|ii|jj|kk|ll|mm|nn|oo|pp|qq|rr|ss|tt`, TopTok: []LR_TokType{LR_OR}, TopVal: []string{"|"}, NExpectedErr: 0},
	/* 085 */ {Re: `(a|b)`, TopTok: []LR_TokType{LR_OP_PAR}, TopVal: []string{"("}, NExpectedErr: 0},
	/* 086 */ {Re: `(a|b)*a`, TopTok: []LR_TokType{LR_STAR, LR_Text}, TopVal: []string{"*", "a"}, NExpectedErr: 0},
	// `(a|b)*aab`			// Test case from Compilers Book, by Aho, Sethi, Ullman
	/* 087 */ {Re: `(a|b)*aab`, TopTok: []LR_TokType{LR_STAR, LR_Text, LR_Text, LR_Text}, TopVal: []string{"*", "a", "a", "b"}, NExpectedErr: 0},
	/* 088 */ {Re: `a(def|ghi)*(klm|(mno|pqr)?|stu)+b`, TopTok: []LR_TokType{LR_Text, LR_STAR, LR_PLUS, LR_Text}, TopVal: []string{"a", "*", "+", "b"}, NExpectedErr: 0},

	/* 089 */ {Re: `(aaa|bbb\|ccc)`, TopTok: []LR_TokType{LR_OP_PAR}, TopVal: []string{"("}, NExpectedErr: 0},
	/* 090 */ {Re: `c|`, TopTok: []LR_TokType{LR_OR}, TopVal: []string{"|"}, NExpectedErr: 0},

	// Involving {a,b} -------------------------------------------------------------------------------------
	/* 091 */ {Re: `c{2,3}`, TopTok: []LR_TokType{LR_OP_BR}, TopVal: []string{"{"}, NExpectedErr: 0},
	/* 092 */ {Re: `c{3,2}`, TopTok: []LR_TokType{LR_OP_BR}, TopVal: []string{"{"}, NExpectedErr: 1},
	/* 093 */ {Re: `c{3,}`, TopTok: []LR_TokType{LR_OP_BR}, TopVal: []string{"{"}, NExpectedErr: 0},
	/* 094 */ {Re: `c{,2}`, TopTok: []LR_TokType{LR_OP_BR}, TopVal: []string{"{"}, NExpectedErr: 0},
	/* 095 */ {Re: `abc{2,3}`, TopTok: []LR_TokType{LR_Text, LR_Text, LR_OP_BR}, TopVal: []string{"a", "b", "{"}, NExpectedErr: 0},
	/* 096 */ {Re: `abc{3,2}`, TopTok: []LR_TokType{LR_Text, LR_Text, LR_OP_BR}, TopVal: []string{"a", "b", "{"}, NExpectedErr: 1},
	/* 097 */ {Re: `abc{3,}`, TopTok: []LR_TokType{LR_Text, LR_Text, LR_OP_BR}, TopVal: []string{"a", "b", "{"}, NExpectedErr: 0},
	/* 098 */ {Re: `abc{,2}`, TopTok: []LR_TokType{LR_Text, LR_Text, LR_OP_BR}, TopVal: []string{"a", "b", "{"}, NExpectedErr: 0},
	/* 099 */ {Re: `abc{,2}|def`, TopTok: []LR_TokType{LR_OR}, TopVal: []string{"|"}, NExpectedErr: 0},
	/* 100 */ {Re: `abc{,}`, TopTok: []LR_TokType{LR_Text, LR_Text, LR_STAR}, TopVal: []string{"a", "b", "*"}, NExpectedErr: 0},

	// Involving Extended CCLs/NCCLs -------------------------------------------------------------------------------------
	/* 101 */ {Re: `a[3[:alpha:]4]d`, TopTok: []LR_TokType{LR_Text, LR_CCL, LR_Text}, TopVal: []string{"a", "34", "d"}, NExpectedErr: 0},

	/* 102 */ {Re: `\{\{`, TopTok: []LR_TokType{LR_Text, LR_Text}, TopVal: []string{"{", "{"}, NExpectedErr: 0},
	/* 103 */ {Re: `\{\|\*\(\)\+\?\[\]`, TopTok: []LR_TokType{LR_Text, LR_Text, LR_Text, LR_Text, LR_Text, LR_Text, LR_Text, LR_Text, LR_Text},
		TopVal: []string{"{", "|", "*", "(", ")", "+", "?", "[", "]"}, NExpectedErr: 0},
	// -----------------------------------------------------------------------------------------------------
	// Later
	// Extended CCLs / N_CCLs
	// -----------------------------------------------------------------------------------------------------

	// Extended CCL/NCCL
	// `c[\w]d`
	// `c[\s]d`
	// `c[\W]d`
	// `c[\S]d`

}

// -----------------------------------------------------------------------------------------------------------------------------------------
// From: https://labix.org/gocheck
// -----------------------------------------------------------------------------------------------------------------------------------------

func TestLexie(t *testing.T) { TestingT(t) }

type ReTesteSuite struct{}

var _ = Suite(&ReTesteSuite{})

func (s *ReTesteSuite) TestLexie(c *C) {

	// Test of Parseing REs into RE-TParseTrees
	fmt.Fprintf(os.Stderr, "Test Parsing of REs, %s\n", dbgo.LF())

	dbgo.SetADbFlag("db_DumpPool", true)
	dbgo.SetADbFlag("parseExpression", true)
	dbgo.SetADbFlag("CalcLength", true)
	dbgo.SetADbFlag("DumpParseNodes", true)
	dbgo.SetADbFlag("DumpParseNodesX", true)

	fmt.Printf("**** In Test RE\n")

	n_err := 0
	n_skip := 0

	for ii, vv := range Lexie00Data {
		if !vv.SkipTest {
			fmt.Printf("\n\n--- %d Test: %s -----------------------------------------------------------------------------\n\n", ii, vv.Test)
			lr := NewLexReType()
			lr.SetBuf(vv.Re)
			tn := 0
			cc, ww := lr.Next()
			for ww != LR_EOF {
				fmt.Printf("    %s % x = %d %s\n", string(cc), string(cc), ww, LR_TokTypeLookup[ww])

				if len(vv.SM) > 0 {
					// Check correct token
					if vv.SM[tn].Tok != ww {
						fmt.Printf("     Failed to return the correct token, expecting %d/%s, got %d/%s\n", vv.SM[tn].Tok, LR_TokTypeLookup[vv.SM[tn].Tok], ww,
							LR_TokTypeLookup[ww])
						c.Check(int(vv.SM[tn].Tok), Equals, int(ww))
						n_err++
					}
					// Check correct string/rune returned
					if !CmpByteArr(vv.SM[tn].Dat, []byte(cc)) {
						fmt.Printf("     The returned runes did not match\n")
						c.Check(string(vv.SM[tn].Dat), Equals, cc)
						n_err++
					}

				}

				tn++
				cc, ww = lr.Next()
			}
			fmt.Printf("\n--- %d End : %s -----------------------------------------------------------------------------\n\n", ii, vv.Test)
		}
	}

	for ii, vv := range Lexie01Data {
		if !vv.SkipTest {
			fmt.Printf("\n\n--- %d Test: %s -----------------------------------------------------------------------------\n\n", ii, vv.Test)

			lr := NewLexReType()
			lr.ParseRe(vv.Re)

			lr.Sigma = lr.GenerateSigma()
			fmt.Printf("Sigma: ->%s<-\n", lr.Sigma)

			if vv.Sigma != "" {
				if vv.Sigma != lr.Sigma {
					fmt.Printf("     The calculated and reference Sigma did not match\n")
					c.Check(vv.Sigma, Equals, lr.Sigma)
					n_err++
				}
			}

			fmt.Printf("\n--- %d End : %s -----------------------------------------------------------------------------\n\n", ii, vv.Test)
		}
	}

	// -------------------------------------------------------------------------------------------------------------------------------------------------
	lr := NewLexReType()
	dbgo.SetADbFlag("DumpParseNodes", true)
	dbgo.SetADbFlag("DumpParseNodesX", true)
	dbgo.SetADbFlag("db_DumpPool", true)
	dbgo.SetADbFlag("parseExpression", true)
	for i, v := range Test6Data {
		dbgo.DbPrintf("debug", "\nTest[%03d]: `%s` %s\n\n", i, v.Re, strings.Repeat("-", 120-len(v.Re)))
		lr.ParseRe(v.Re)
		lr.DumpParseNodes()
		if len(v.TopTok) > 0 {
			if len(v.TopTok) != len(lr.Tree.Children) {
				dbgo.DbPrintf("debug", "%(red)Error%(reset): wrong number of tokens prsed, Expected: %d Got %d\n",  len(v.TopTok), len(lr.Tree.Children))
				n_err++
			} else {
				for i := 0; i < len(v.TopTok); i++ {
					if v.TopTok[i] != lr.Tree.Children[i].LR_Tok {
						dbgo.DbPrintf("debug", "%(red)Error%(reset): invalid token returnd at postion %d\n",  i)
						c.Check(v.TopTok[i], Equals, lr.Tree.Children[i].LR_Tok)
						n_err++
					}
				}
			}
		}
		if len(v.TopVal) > 0 {
			if len(v.TopVal) != len(lr.Tree.Children) {
				dbgo.DbPrintf("debug", "%(red)Error%(reset): wrong number of tokens prsed, Expected: %d Got %d - Based on number of values, TopVal\n",  len(v.TopVal), len(lr.Tree.Children))
				n_err++
			} else {
				for i := 0; i < len(v.TopVal); i++ {
					if v.TopVal[i] != lr.Tree.Children[i].Item {
						dbgo.DbPrintf("debug", "%(red)Error%(reset): invalid value at postion %d, %s\n",  i, dbgo.LF())
						n_err++
					}
				}
			}
		}
		if len(lr.Error) > v.NExpectedErr {
			dbgo.DbPrintf("debug", "%(red)Error%(reset): Errors reported in R.E. parsing %d\n",  len(lr.Error))
			n_err++
		} else if len(lr.Error) > 0 {
			dbgo.DbPrintf("debug", "%(green)Note%(reset): Errors reported in R.E. parsing %d\n",  len(lr.Error))
		}
		lr.Error = lr.Error[:0]
		dbgo.DbPrintf("debug", "\nDone[%03d]: `%s` %s\n\n", i, v.Re, strings.Repeat("-", 120-len(v.Re)))
	}
	if n_err > 0 {
		fmt.Fprintf(os.Stderr, "%(red)Failed, # of errors = %d\n",  n_err )
		dbgo.DbPrintf("debug", "\n\n%(red)Failed, # of errors = %d\n",  n_err)
	} else {
		fmt.Fprintf(os.Stderr, "%(green)PASS\n")
		dbgo.DbPrintf("debug", "\n\n%(green)PASS\n")
	}

	if n_skip > 0 {
		fmt.Fprintf(os.Stderr, "%(yellow)Skipped, # of files without automated checks = %d\n", n_skip)
		dbgo.DbPrintf("debug", "\n\n%(yellow)Skipped, # of files without automated checks = %d\n",  n_skip)
	}
	if n_err > 0 {
		c.Check(n_err, Equals, 0)
		fmt.Fprintf(os.Stderr, "%(red)Failed, # of errors = %d\n",  n_err)
		dbgo.DbPrintf("debug", "\n\n%(red)Failed, # of errors = %d\n",  n_err)
	} else {
		fmt.Fprintf(os.Stderr, "%(red)PASS\n", 
		dbgo.DbPrintf("debug", "\n\n%(red)PASS\n"
	}
}

func CmpByteArr(a []byte, b []byte) (rv bool) {
	rv = false
	if len(a) != len(b) {
		return
	}
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return
		}
	}
	rv = true
	return
}

/* vim: set noai ts=4 sw=4: */
