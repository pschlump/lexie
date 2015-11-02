package gen

// Defs - OutputDef
// ==========================================================================
// DefType: Tokens
// ==========================================================================
const (
	Tok_null        = 0
	Tok_L_EQ        = 1
	Tok_GE          = 2
	Tok_LE          = 3
	Tok_L_AND       = 4
	Tok_L_OR        = 5
	Tok_OP_VAR      = 6
	Tok_CL_VAR      = 7
	Tok_OP_BL       = 8
	Tok_CL_BL       = 9
	Tok_NE          = 10
	Tok_NE_LG       = 11
	Tok_OP          = 12
	Tok_CL          = 13
	Tok_PLUS        = 14
	Tok_MINUS       = 15
	Tok_STAR        = 16
	Tok_LT          = 17
	Tok_GT          = 18
	Tok_SLASH       = 19
	Tok_CARRET      = 20
	Tok_COMMA       = 21
	Tok_DOT         = 22
	Tok_EXCLAM      = 23
	Tok_OR          = 24
	Tok_COLON       = 25
	Tok_EQ          = 26
	Tok_PCT         = 27
	Tok_in          = 28
	Tok_and         = 29
	Tok_or          = 30
	Tok_not         = 31
	Tok_true        = 32
	Tok_false       = 33
	Tok_as          = 34
	Tok_export      = 35
	Tok_SS          = 36
	Tok_PIPE        = 37
	Tok_OP_SQ       = 38
	Tok_CL_SQ       = 39
	Tok_TILDE       = 40
	Tok_B_AND       = 41
	Tok_B_OR        = 42
	Tok_S_L         = 43
	Tok_S_R         = 44
	Tok_PLUS_EQ     = 45
	Tok_MINUS_EQ    = 46
	Tok_STAR_EQ     = 47
	Tok_DIV_EQ      = 48
	Tok_MOD_EQ      = 49
	Tok_CAROT_EQ    = 50
	Tok_B_OR_EQ     = 51
	Tok_B_AND_EQ    = 52
	Tok_OP_BRACE    = 53
	Tok_CL_BRACE    = 54
	Tok_TILDE_EQ    = 55
	Tok_TILDE_TILDE = 56
	Tok_EQ3         = 57
	Tok_APROX_EQ    = 58
	Tok_QUEST       = 59
	Tok_RE_MATCH    = 60
	Tok_PLUS_PLUS   = 61
	Tok_MINUS_MINUS = 62
	Tok_DCL_VAR     = 63
	Tok_XOR         = 64
	Tok_Call        = 65
	Tok_EOF         = 66
	Tok_Float       = 67
	Tok_HTML        = 68
	Tok_ID          = 69
	Tok_Ignore      = 70
	Tok_NUM         = 71
	Tok_S_L_EQ      = 72
	Tok_S_R_EQ      = 73
	Tok_Str0        = 74
)

// ==========================================================================
// DefType: Machines
// ==========================================================================
const (
	S_Init   = 0
	S_Common = 1
	S_TAG    = 2
	S_VAR    = 3
	S_Quote  = 4
	S_Str0   = 5
	S_Str1   = 6
	S_Str2   = 7
)

// ==========================================================================
// DefType: Errors
// ==========================================================================
const (
	Err_EOF_In_String       = 1
	Err_EOF_Tag             = 2
	Warn_End_Var_Unexpected = 3
)

// ==========================================================================
// DefType: ReservedWords
// ==========================================================================
const (
	RW_and    = 4
	RW_or     = 5
	RW_in     = 28
	RW_not    = 31
	RW_true   = 32
	RW_false  = 33
	RW_as     = 34
	RW_export = 35
	RW_band   = 41
	RW_bor    = 42
)
