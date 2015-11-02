package re

import "strings"

const (
	X_DOT          = "\uF8FA" // Any char in Sigma
	X_BOL          = "\uF8F3" // Beginning of line
	X_EOL          = "\uF8F4" // End of line
	X_NUMERIC      = "\uF8F5"
	X_LOWER        = "\uF8F6"
	X_UPPER        = "\uF8F7"
	X_ALPHA        = "\uF8F8"
	X_ALPHNUM      = "\uF8F9"
	X_EOF          = "\uF8FB"
	X_not_CH       = "\uF8FC" // On input lookup if the char is NOT in Signa then it is returned as this.
	X_else_CH      = "\uF8FC" // If char is not matched in this state then take this path
	X_N_CCL        = "\uF8FD"
	X_LAMBDA_MATCH = "\uF8FE"
)

const (
	R_min_reserved = '\uF8FA' //
	R_DOT          = '\uF8FA' // Any char in Sigma
	R_BOL          = '\uF8F3' // Beginning of line
	R_EOL          = '\uF8F4' // End of line
	R_NUMERIC      = '\uF8F5'
	R_LOWER        = '\uF8F6'
	R_UPPER        = '\uF8F7'
	R_ALPHA        = '\uF8F8'
	R_ALPHNUM      = '\uF8F9'
	R_EOF          = '\uF8FB'
	R_not_CH       = '\uF8FC' // On input lookup if the char is NOT in Signa then it is returned as this.
	R_else_CH      = '\uF8FC' // If char is not matched in this state then take this path
	R_N_CCL        = '\uF8FD' // If char is not matched in this state then take this path
	R_LAMBDA_MATCH = '\uF8FE'
)

const InfiniteIteration = 9999999999

type LR_TokType int

const (
	LR_null   LR_TokType = iota //  0
	LR_Text                     //  1
	LR_EOF                      //  2
	LR_DOT                      //  3 .
	LR_STAR                     //  4 *
	LR_PLUS                     //  5 +
	LR_QUEST                    //  6 ?
	LR_OP_PAR                   //  7 (
	LR_CL_PAR                   //  8 )
	LR_CCL                      //  9 [...]
	LR_N_CCL                    // 10 [^...]
	LR_E_CCL                    // 11 ]
	LR_CARROT                   // 12 ^
	LR_MINUS                    // 13 -
	LR_DOLLAR                   // 14 $
	LR_OR                       // 15 |
	LR_OP_BR                    // 16 {
	LR_CL_BR                    // 17 }
	LR_COMMA                    // 18 ,
)

var LR_TokTypeLookup map[LR_TokType]string

func init() {
	LR_TokTypeLookup = make(map[LR_TokType]string)
	LR_TokTypeLookup[LR_null] = "LR_null"
	LR_TokTypeLookup[LR_Text] = "LR_Text"
	LR_TokTypeLookup[LR_EOF] = "LR_EOF"
	LR_TokTypeLookup[LR_DOT] = "LR_DOT"
	LR_TokTypeLookup[LR_STAR] = "LR_STAR"
	LR_TokTypeLookup[LR_PLUS] = "LR_PLUS"
	LR_TokTypeLookup[LR_QUEST] = "LR_QUEST"
	LR_TokTypeLookup[LR_OP_PAR] = "LR_OP_PAR"
	LR_TokTypeLookup[LR_CL_PAR] = "LR_CL_PAR"
	LR_TokTypeLookup[LR_CCL] = "LR_CCL"
	LR_TokTypeLookup[LR_N_CCL] = "LR_N_CCL"
	LR_TokTypeLookup[LR_E_CCL] = "LR_E_CCL"
	LR_TokTypeLookup[LR_CARROT] = "LR_CARROT"
	LR_TokTypeLookup[LR_MINUS] = "LR_MINUS"
	LR_TokTypeLookup[LR_DOLLAR] = "LR_DOLLAR"
	LR_TokTypeLookup[LR_OR] = "LR_OR"
	LR_TokTypeLookup[LR_OP_BR] = "LR_OP_BR"
	LR_TokTypeLookup[LR_CL_BR] = "LR_CL_BR"
	LR_TokTypeLookup[LR_COMMA] = "LR_COMMA"
}

var LexReMatcher = []LexReMatcherType{
	{Sym: `\[`, Rv: LR_Text, Repl: "["},
	{Sym: `\]`, Rv: LR_Text, Repl: "]"},
	{Sym: `\(`, Rv: LR_Text, Repl: "("},
	{Sym: `\)`, Rv: LR_Text, Repl: ")"},
	{Sym: `\^`, Rv: LR_Text, Repl: "^"},
	{Sym: `\?`, Rv: LR_Text, Repl: "?"},
	{Sym: `\*`, Rv: LR_Text, Repl: "*"},
	{Sym: `\+`, Rv: LR_Text, Repl: "+"},
	{Sym: `\.`, Rv: LR_Text, Repl: "."},
	{Sym: `\-`, Rv: LR_Text, Repl: "-"},
	{Sym: `\^`, Rv: LR_Text, Repl: "^"},
	{Sym: `\$`, Rv: LR_Text, Repl: "$"},
	{Sym: `\\`, Rv: LR_Text, Repl: "\\"},
	{Sym: `\|`, Rv: LR_Text, Repl: "|"},
	{Sym: `\{`, Rv: LR_Text, Repl: "{"},
	{Sym: `\}`, Rv: LR_Text, Repl: "}"},
	{Sym: `\,`, Rv: LR_Text, Repl: ","},
	{Sym: "[^", Rv: LR_N_CCL},
	{Sym: ".", Rv: LR_DOT},
	{Sym: "*", Rv: LR_STAR},
	{Sym: "+", Rv: LR_PLUS},
	{Sym: "?", Rv: LR_QUEST},
	{Sym: "[", Rv: LR_CCL},
	{Sym: "]", Rv: LR_E_CCL},
	{Sym: "(", Rv: LR_OP_PAR},
	{Sym: ")", Rv: LR_CL_PAR},
	{Sym: "^", Rv: LR_CARROT},
	{Sym: "-", Rv: LR_MINUS},
	{Sym: "$", Rv: LR_DOLLAR},
	{Sym: "|", Rv: LR_OR},
	{Sym: "{", Rv: LR_OP_BR},
	{Sym: "}", Rv: LR_CL_BR},
	{Sym: ",", Rv: LR_COMMA},
}

func EscapeStr(s string) string {
	if s == X_DOT {
		// return "\u2022"	// Middle Bullet
		return "\u25C9" // FishEye
	} else if s == X_NUMERIC {
		return "0-9"
	} else if s == X_LOWER {
		return "a-z"
	} else if s == X_UPPER {
		return "A-Z"
	} else if s == X_ALPHA {
		return "a-zA-Z"
	} else if s == X_ALPHNUM {
		return "a-zA-Z0-9"
	} else if s == X_EOF {
		return "EOF"
	} else if s == X_BOL {
		return "BOL"
	} else if s == X_EOL {
		return "EOL"
	}
	s = strings.Replace(s, `\`, `\\`, -1)
	s = strings.Replace(s, `"`, `\"`, -1)
	return s
}

func EscapeStrForGV(s string) string {
	if s == X_DOT {
		// return "\u2022"	// Middle Bullet
		return "\u25C9" // FishEye
	} else if s == X_NUMERIC {
		return "0-9"
	} else if s == X_LOWER {
		return "a-z"
	} else if s == X_UPPER {
		return "A-Z"
	} else if s == X_ALPHA {
		return "a-zA-Z"
	} else if s == X_ALPHNUM {
		return "a-zA-Z0-9"
	} else if s == X_EOF {
		return "EOF"
	} else if s == X_BOL {
		return "BOL"
	} else if s == X_EOL {
		return "EOL"
	}
	s = strings.Replace(s, `\`, `\\`, -1)
	s = strings.Replace(s, `"`, `\"`, -1)
	s = strings.Replace(s, "\t", `\t`, -1)
	s = strings.Replace(s, "\n", `\n`, -1)
	s = strings.Replace(s, "\f", `\f`, -1)
	s = strings.Replace(s, "\v", `\v`, -1)
	s = strings.Replace(s, "\r", `\r`, -1)
	s = strings.Replace(s, " ", `\B`, -1)
	return s
}
