
//
// Lexie Input for the Django Teplate Superset, Ringo
//
// (C) Philip Schlump, 2010-2015.
// Version: 1.0.8
// BuildNo: 141
//
// Notes -----------------------------------------------------------------------------------------------------------------
//
// 1. report all unused tokens
// 2. report any unused and defiend "Machines"
// 3. report all unused errors
// 4. report line number where token first seen ( Cross ref on Token Line Numbers )
// 5. Report any machine that fails to accep the entire Sigma
// 6. Report any machine that is not callable S_Init is top level and always callable.
// 7. Report any non-S_Init that has no Return
// 8. report any errors/warnings that are not declared - Allow but list
//
// 
// $def('type',"Name","Value",...)
// $machine(Name[,Mixin...])
// "token"				: 	ActionInfo
// 'token'				: 	ActionInfo
// token				:	ActionInfo
// `token`				:	ActionInfo
// $end
// $eof
//
// Test: [+-]?[0-9]+(\.[0-9]*([eE][+-][0-9]*)?)?
// 

$def(Tokens, Tok_L_EQ=1, Tok_GE=2, Tok_LE=3, Tok_L_AND=4, Tok_L_OR=5 , Tok_OP_VAR=6, Tok_CL_VAR=7, Tok_OP_BL=8, Tok_CL_BL=9 , Tok_NE=10, Tok_NE_LG=11, Tok_OP=12, Tok_CL=13, Tok_PLUS=14, Tok_MINUS=15, Tok_STAR=16, Tok_LT=17, Tok_GT=18, Tok_SLASH=19, Tok_CARRET=20, Tok_COMMA=21, Tok_DOT=22, Tok_EXCLAM=23, Tok_OR=24, Tok_COLON=25, Tok_EQ=26, Tok_PCT=27, Tok_in=28, Tok_and=29, Tok_or=30, Tok_not=31, Tok_true=32, Tok_false=33, Tok_as=34, Tok_export=35, Tok_SS)

$def(Machines, S_Init, S_TAG, S_Common, S_Esc, S_Str0, S_Str1, S_Str2, S_VAR, S_Quote )

$def(Errors,  Warn_End_Var_Unexpected, Err_EOF_Tag, Err_EOF_In_String )

$def(ReservedWords, and=Tok_L_AND, or=Tok_L_OR, true=Tok_true, false=Tok_false, not=Tok_not, export=Tok_export, in=Tok_in, not=Tok_not, as=Tok_as )

$machine(S_Init)
`{{`					: Rv(Tok_OP_VAR) Call(S_VAR)
`{%`					: Rv(Tok_OP_BL) Call(S_TAG)
`{\{`					: Repl(`{{`)					// Implies a terminal state of Rv(Tok)
`{\%`					: Repl(`{%`)					// Implies a terminal state of Rv(Tok)
.*						: Rv(Tok_HTML)
$eof					: Rv(Tok_EOF)
$end


$machine(S_Common)
<=										: Rv(Tok_LE)
==										: Rv(Tok_L_EQ)
>=										: Rv(Tok_GE)
&&										: Rv(Tok_L_AND)
`||`									: Rv(Tok_L_OR)
!=										: Rv(Tok_NE)
<>										: Rv(Tok_NE)
`{{`									: Rv(Tok_OP_VAR) Call(S_VAR)
`{%`									: Rv(Tok_OP_BL) Call(S_TAG)
`%}`									: Return()	Rv(Tok_CL_BL)
`}}`									: Return()	Rv(Tok_CL_VAR)
[a-zA-Z_][a-zA-Z0-9_]*					: Rv(Tok_ID) ReservedWord()		// Both Return ID nd check if it is a reserved word
[0-9]+									: Rv(Tok_NUM)
[0-9]+\.[0-9]+([eE][0-9]+(\.[0-9]*)?)?	: Rv(Tok_Float)
`"`										: Call(S_Str0) Ignore()					// Implies a terminal state of Rv(Tok)
`'`										: Call(S_Str1) Ignore()					// Implies a terminal state of Rv(Tok)
"`"										: Call(S_Str2) Ignore()					// Implies a terminal state of Rv(Tok)
`^`										: Rv(Tok_CARRET)
`(`										: Rv(Tok_OP)
`)`										: Rv(Tok_CL)
`+`										: Rv(Tok_PLUS)
-										: Rv(Tok_MINUS)
`*`										: Rv(Tok_STAR)
<										: Rv(Tok_LT)
>										: Rv(Tok_GT)
/										: Rv(Tok_SLASH)
`%`										: Rv(Tok_PCT)
`|`										: Rv(Tok_OR)
,										: Rv(Tok_COMMA)
`.`										: Rv(Tok_DOT)
!										: Rv(Tok_EXCLAM)
:										: Rv(Tok_COLON)
=										: Rv(Tok_EQ)
"[ \t\n\f\r]"							: Ignore()						// Implies a terminal Rv(Tok) + no return and go to state 0
.										: Warn(Warn_Unrecog_Char)	Return()
$eof									: Error(Err_EOF_Tag) 		// Error - State Machine Exit
$end


$machine(S_TAG,S_Common)
$end

$machine(S_VAR,S_Common)
$end

$machine(S_Quote)
.						: Return()
$eof					: Error(Err_EOF_In_String) 	// Error - State Machine Exit
$end

$machine(S_Str0)
`\`						: Call(S_Quote) Rv(Tok_Call)
.*"						: Rv(Tok_Str0) NotGreedy() Return()
$eof					: Error(Err_EOF_In_String) 	// Error - State Machine Exit
$end

$machine(S_Str1)
`\`						: Call(S_Quote) Rv(Tok_Call)
.*'						: Rv(Tok_Str1) NotGreedy() Return()
$eof					: Error(Err_EOF_In_String) 	// Error - State Machine Exit
$end

$machine(S_Str2)
"``"					: Repl("`")
.*`						: Rv(Tok_Str0)
$eof					: Error(Err_EOF_In_String) 	// Error - State Machine Exit
$end


