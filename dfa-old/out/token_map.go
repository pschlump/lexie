package test03package

import "fmt"

type TokenType int

const (
	Tok_CL     TokenType = 1
	Tok_EOF    TokenType = 4
	Tok_PCT    TokenType = 2
	Tok_null   TokenType = 0
	Tok_ID     TokenType = 5
	Tok_Ignore TokenType = 6
	Tok_BB     TokenType = 3
)

func (tt TokenType) String() string {
	switch tt {

	case Tok_EOF: /* 4 */
		return "Tok_EOF"

	case Tok_PCT: /* 2 */
		return "Tok_PCT"

	case Tok_null: /* 0 */
		return "Tok_null"

	case Tok_ID: /* 5 */
		return "Tok_ID"

	case Tok_Ignore: /* 6 */
		return "Tok_Ignore"

	case Tok_BB: /* 3 */
		return "Tok_BB"

	case Tok_CL: /* 1 */
		return "Tok_CL"

	default:
		return fmt.Sprintf("--unknown TokenType %d--", int(tt))
	}
}

/* Printed At:File: /Users/philip/go/src/github.com/pschlump/lexie/dfa-old/gogen.go LineNo:57 defs are = lex.Im.Def.DefsAre = {
	"Errors": {
		"Seq": 2,
		"WhoAmI": "Errors",
		"NameValueStr": {
			"Err_Invalid_Char": "1"
		},
		"NameValue": {
			"Err_Invalid_Char": 1
		},
		"Reverse": {
			"1": "Err_Invalid_Char"
		},
		"SeenAt": {
			"Err_Invalid_Char": {
				"LineNo": [
					6,
					14
				],
				"FileName": [
					"unk-file",
					"unk-file"
				]
			}
		}
	},
	"Machines": {
		"Seq": 1,
		"WhoAmI": "Machines",
		"NameValueStr": {
			"S_Init": "0"
		},
		"NameValue": {
			"S_Init": 0
		},
		"Reverse": {
			"0": "S_Init"
		},
		"SeenAt": {
			"S_Init": {
				"LineNo": [
					4
				],
				"FileName": [
					"unk-file"
				]
			}
		}
	},
	"Options": {
		"Seq": 1,
		"WhoAmI": "Options",
		"NameValueStr": {
			"GoPackageName": "test03package"
		},
		"NameValue": {
			"GoPackageName": -2
		},
		"Reverse": {},
		"SeenAt": {
			"GoPackageName": {
				"LineNo": [
					8
				],
				"FileName": [
					"unk-file"
				]
			}
		}
	},
	"ReservedWords": {
		"Seq": 0,
		"WhoAmI": "",
		"NameValueStr": null,
		"NameValue": null,
		"Reverse": null,
		"SeenAt": null
	},
	"Tokens": {
		"Seq": 7,
		"WhoAmI": "Tokens",
		"NameValueStr": {
			"Tok_BB": "3",
			"Tok_CL": "1",
			"Tok_EOF": "4",
			"Tok_ID": "5",
			"Tok_Ignore": "6",
			"Tok_PCT": "2",
			"Tok_null": "0"
		},
		"NameValue": {
			"Tok_BB": 3,
			"Tok_CL": 1,
			"Tok_EOF": 4,
			"Tok_ID": 5,
			"Tok_Ignore": 6,
			"Tok_PCT": 2,
			"Tok_null": 0
		},
		"Reverse": {
			"0": "Tok_null",
			"1": "Tok_CL",
			"2": "Tok_PCT",
			"3": "Tok_BB",
			"4": "Tok_EOF",
			"5": "Tok_ID",
			"6": "Tok_Ignore"
		},
		"SeenAt": {
			"Tok_BB": {
				"LineNo": [
					2,
					13
				],
				"FileName": [
					"unk-file",
					"unk-file"
				]
			},
			"Tok_CL": {
				"LineNo": [
					2,
					11
				],
				"FileName": [
					"unk-file",
					"unk-file"
				]
			},
			"Tok_EOF": {
				"LineNo": [
					2
				],
				"FileName": [
					"unk-file"
				]
			},
			"Tok_ID": {
				"LineNo": [
					543
				],
				"FileName": [
					"/Users/philip/go/src/github.com/pschlump/lexie/in/in.go"
				]
			},
			"Tok_Ignore": {
				"LineNo": [
					543
				],
				"FileName": [
					"/Users/philip/go/src/github.com/pschlump/lexie/in/in.go"
				]
			},
			"Tok_PCT": {
				"LineNo": [
					2,
					12
				],
				"FileName": [
					"unk-file",
					"unk-file"
				]
			},
			"Tok_null": {
				"LineNo": [
					499
				],
				"FileName": [
					"/Users/philip/go/src/github.com/pschlump/lexie/in/in.go"
				]
			}
		}
	}
} */
