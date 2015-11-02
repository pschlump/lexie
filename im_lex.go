package lexie

/*
//
// Lexie Input for the Django Teplate Superset, Ringo
//
// (C) Philip Schlump, 2010-2015.
// Version: 1.0.8
// BuildNo: 141
//

$def('type',"Name","Value",...)
$machine(Name[,Mixin...])
"token"				: 	ActionInfo
'token'				: 	ActionInfo
token				:	ActionInfo
`token`				:	ActionInfo
$end
$eof


$def('Tokens', Tok_null=0, Tok_ID=1 )

$def('Machines', S_Init, S_TAG, S_Common, S_Esc, S_Str0, S_Str1, S_Str2, S_VAR, S_Quote )

$def('Errors',  Warn_End_Var_Unexpected, Err_EOF_Tag, Err_EOF_Variable, Err_EOF_In_String, Err_EOF_Uneval_Quote )

$def('ReservedWords', and=Tok_L_AND, or=Tok_L_OR, true, false, not, export )

$machine(S_Init)
`{{`					: Rv(Tok_OP_Var) Call(S_VAR)
`{%`					: Rv(Tok_OP_Tag) Call(S_TAG)
`{\{`					: Repl(`{{`)					// Implies a terminal state of Rv(Tok)
`{\%`					: Repl(`{%`)					// Implies a terminal state of Rv(Tok)
.*						: Rv(Tok_HTML)
$eof					: Rv(Tok_EOF)
$end


$machine(S_Common)
[a-zA-Z_][a-zA-Z0-0_]*	: Rv(Tok_ID) ReservedWord()		// Both Return ID nd check if it is a reserved word
[0-9]+					: Rv(Tok_NUM)
<=						: Rv(Tok_LE)
==						: Rv(Tok_EQEQ)
>=						: Rv(Tok_GE)
&&						: Rv(Tok_L_AND)
`||`					: Rv(Tok_L_OR)
!=						: Rv(Tok_NE)
<>						: Rv(Tok_NE)
`"`						: Call(S_Str0)					// Implies a terminal state of Rv(Tok)
`'`						: Call(S_Str1)					// Implies a terminal state of Rv(Tok)
"`"						: Call(S_Str2)					// Implies a terminal state of Rv(Tok)
`^`						: Rv(Tok_CARROT)
`(`						: Rv(Tok_OP_PAR)
`)`						: Rv(Tok_CL_PAR)
+						: Rv(Tok_PLUS)
-						: Rv(Tok_MINUS)
`*`						: Rv(Tok_STAR)
<						: Rv(Tok_LT)
>						: Rv(Tok_GT)
/						: Rv(Tok_SLASH)
`%`						: Rv(Tok_PCT)
`|`						: Rv(Tok_OR)
`=`						: Rv(Tok_EQ)
"[ \t\n\f\r]"			: Ignore()						// Implies a terminal Rv(Tok) + no return and go to state 0
.						: Warn(Warn_Unrecog_Char)	Return()
$end


$machine(S_TAG,S_Common)
`%}`					: Return()
`}}`					: Warn(Warn_End_Var_Unexpected)
$eof					: Error(Err_EOF_Tag) 		// Error - State Machine Exit
$end

$machine(S_VAR,S_Common)
`}}`					: Return()
`%}`					: Warn(Warn_End_Tag_Unexpected)
$eof					: Error(Err_EOF_Var) 		// Error - State Machine Exit
$end

$machine(S_Quote)
.						: Return()
$eof					: Error(Err_EOF_In_String) 	// Error - State Machine Exit
$end

$machien(S_Str0)
`"`						: Rv(Tok_Str0)
`\`						: Call(S_Quote)
.						: Accept()
$eof					: Error(Err_EOF_In_String) 	// Error - State Machine Exit
$end

$machien(S_Str1)
`'`						: Rv(Tok_Str1)
`\`						: Call(S_Quote)
.						: Accept()
$eof					: Error(Err_EOF_In_String) 	// Error - State Machine Exit
$end

$machien(S_Str2)
"``"					: Repl("`")
"`"						: Rv(Tok_Str0)
.						: Accept()
$eof					: Error(Err_EOF_In_String) 	// Error - State Machine Exit
$end

*/
