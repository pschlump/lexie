Machine:
	
	//
	// Lexie Input for the Django Teplate Superset, Ringo
	//
	// (C) Philip Schlump, 2010-2025.
	// Version: 1.0.8
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
	
	$def(Tokens, Tok_L_EQ=1, Tok_GE=2, Tok_LE=3, Tok_L_AND=4, Tok_L_OR=5, Tok_OP_VAR=6, Tok_CL_VAR=7, Tok_OP_BL=8, Tok_CL_BL=9, Tok_NE=10, Tok_NE_LG=11, Tok_OP=12, Tok_CL=13, Tok_PLUS=14, Tok_MINUS=15, Tok_STAR=16, Tok_LT=17, Tok_GT=18, Tok_SLASH=19, Tok_CARRET=20, Tok_COMMA=21, Tok_DOT=22, Tok_EXCLAM=23, Tok_OR=24, Tok_COLON=25, Tok_EQ=26, Tok_PCT=27, Tok_in=28, Tok_and=29, Tok_or=30, Tok_not=31, Tok_true=32, Tok_false=33, Tok_as=34, Tok_export=35, Tok_SS=36, Tok_PIPE=37, Tok_OP_SQ=38, Tok_CL_SQ=39, Tok_TILDE=40, Tok_B_AND=41, Tok_B_OR=42, Tok_S_L=43, Tok_S_R=44, Tok_PLUS_EQ=45, Tok_MINUS_EQ=46, Tok_STAR_EQ=47, Tok_DIV_EQ=48, Tok_MOD_EQ=49, Tok_CAROT_EQ=50, Tok_B_OR_EQ=51, Tok_B_AND_EQ=52, Tok_OP_BRACE=53, Tok_CL_BRACE=54, Tok_TILDE_EQ=55, Tok_TILDE_TILDE=56, Tok_EQ3=57, Tok_APROX_EQ=58, Tok_QUEST=59, Tok_RE_MATCH=60, Tok_PLUS_PLUS=61, Tok_MINUS_MINUS=62, Tok_DCL_VAR=63, Tok_XOR=64)
	
	$def(Machines, S_Init, S_TAG, S_Common, S_Str0, S_Str1, S_Str2, S_VAR, S_Quote )
	
	$def(Errors,  Warn_End_Var_Unexpected, Err_EOF_Tag, Err_EOF_In_String )
	
	$def(ReservedWords, and=Tok_L_AND, or=Tok_L_OR, true=Tok_true, false=Tok_false, not=Tok_not, export=Tok_export, in=Tok_in, not=Tok_not, as=Tok_as, bor=Tok_B_OR, band=Tok_B_AND, xor=Tok_XOR )
	
	$machine(S_Init)
	`{{`					: Rv(Tok_OP_VAR) Call(S_VAR)
	`{%`					: Rv(Tok_OP_BL) Call(S_TAG)
	`{\\{`					: Repl(`{{`)					// Implies a terminal state of Rv(Tok)
	`{\\%`					: Repl(`{%`)					// Implies a terminal state of Rv(Tok)
	.*						: Rv(Tok_HTML)
	$eof					: Rv(Tok_EOF)
	$end
	
	
	$machine(S_Common)
	<<=										: Rv(Tok_S_L_EQ)
	>>=										: Rv(Tok_S_R_EQ)
	===										: Rv(Tok_EQ3)
	=~=										: Rv(Tok_APROX_EQ)
	<=										: Rv(Tok_LE)
	<<										: Rv(Tok_S_L)
	>>										: Rv(Tok_S_R)
	==										: Rv(Tok_L_EQ)
	>=										: Rv(Tok_GE)
	&&										: Rv(Tok_L_AND)
	`||`									: Rv(Tok_L_OR)
	!=										: Rv(Tok_NE)
	<>										: Rv(Tok_NE)
	`+=`									: Rv(Tok_PLUS_EQ)
	`-=`									: Rv(Tok_MINUS_EQ)
	`*=`									: Rv(Tok_STAR_EQ)
	`/=`									: Rv(Tok_DIV_EQ)
	`%=`									: Rv(Tok_MOD_EQ)
	`^=`									: Rv(Tok_CAROT_EQ)
	`|=`									: Rv(Tok_B_OR_EQ)
	`&=`									: Rv(Tok_B_AND_EQ)
	`~=`									: Rv(Tok_TILDE_EQ)
	`~~`									: Rv(Tok_TILDE_TILDE)
	`?=`									: Rv(Tok_RE_MATCH)
	`{{`									: Rv(Tok_OP_VAR) Call(S_VAR)
	`{%`									: Rv(Tok_OP_BL) Call(S_TAG)
	`%}`									: Rv(Tok_CL_BL) Return()
	`}}`									: Rv(Tok_CL_CL) Return()
	`++`									: Rv(Tok_PLUS_PLUS) 
	--										: Rv(Tok_MINUS_MINUS) 
	:=										: Rv(Tok_DCL_VAR)
	[a-zA-Z_][a-zA-Z0-9_]*					: Rv(Tok_ID) ReservedWord()		// Both Return ID nd check if it is a reserved word
	[0-9]									: Call(S_Num) 
	`"`										: Call(S_Str0) Ignore()					// Implies a terminal state of Rv(Tok)
	`'`										: Call(S_Str1) Ignore()					// Implies a terminal state of Rv(Tok)
	"`"										: Call(S_Str2) Ignore()					// Implies a terminal state of Rv(Tok)
	`^`										: Rv(Tok_CARRET)
	`(`										: Rv(Tok_OP)
	`)`										: Rv(Tok_CL)
	`+`										: Rv(Tok_PLUS)
	`[`										: Rv(Tok_OP_SQ)
	`.`										: Rv(Tok_DOT)
	`]`										: Rv(Tok_CL_SQ)
	`~`										: Rv(Tok_TILDE)
	`&`										: Rv(Tok_B_AND)
	`|`										: Rv(Tok_PIPE)
	`?`										: Rv(Tok_QUEST)
	-										: Rv(Tok_MINUS)
	`*`										: Rv(Tok_STAR)
	<										: Rv(Tok_LT)
	>										: Rv(Tok_GT)
	/										: Rv(Tok_SLASH)
	`%`										: Rv(Tok_PCT)
	`|`										: Rv(Tok_OR)
	`{`										: Rv(Tok_OP_BRACE)
	`}`										: Rv(Tok_CL_BRACE)
	,										: Rv(Tok_COMMA)
	`.`										: Rv(Tok_DOT)
	!										: Rv(Tok_EXCLAM)
	:										: Rv(Tok_COLON)
	=										: Rv(Tok_EQ)
	"[ \t\n\f\r]"							: Ignore()						// Implies a terminal Rv(Tok) + no return and go to state 0
	.										: Warn(Warn_Unrecog_Char)	Reset()
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
	
	$machine(S_Num)
	[0-9]*\.[0-9]+([eE][-+]?[0-9]+(\.[0-9]*)?)?		: Rv(Tok_Float) Return()
	x[0-9a-fA-F]*									: Rv(Tok_NUM) Return()
	[0-9]*											: Rv(Tok_NUM) Return()
	$eof											: Error(Err_EOF_In_String) 	// Error - State Machine Exit
	$end
	
	
{"Input":"Lex-Machine-6", "Rv":1, "Start": 0, "States":[
 { "Sn":0,  "Edge":[ { "On":"'", "Fr":0, "To":1 }, { "On":"\", "Fr":0, "To":2 }, { "On":"e", "Fr":0, "To":0 }, { "On":"f", "Fr":0, "To":0 }, { "On":"o", "Fr":0, "To":0 }, { "On":"", "Fr":0, "To":3 }, { "On":"", "Fr":0, "To":0 }]}
 { "Sn":1,  "Term":76,  "Edge":[ { "On":"'", "Fr":1, "To":1 }, { "On":"\", "Fr":1, "To":2 }, { "On":"e", "Fr":1, "To":0 }, { "On":"f", "Fr":1, "To":0 }, { "On":"o", "Fr":1, "To":0 }, { "On":"", "Fr":1, "To":3 }, { "On":"", "Fr":1, "To":0 }]}
 { "Sn":2,  "Term":66,  "Edge":[ { "On":"'", "Fr":2, "To":1 }, { "On":"\", "Fr":2, "To":2 }, { "On":"e", "Fr":2, "To":0 }, { "On":"f", "Fr":2, "To":0 }, { "On":"o", "Fr":2, "To":0 }, { "On":"", "Fr":2, "To":3 }, { "On":"", "Fr":2, "To":0 }]}
 { "Sn":3,  "Edge":[ { "On":"'", "Fr":3, "To":1 }, { "On":"\", "Fr":3, "To":2 }, { "On":"e", "Fr":3, "To":4 }, { "On":"f", "Fr":3, "To":0 }, { "On":"o", "Fr":3, "To":0 }, { "On":"", "Fr":3, "To":3 }, { "On":"", "Fr":3, "To":0 }]}
 { "Sn":4,  "Edge":[ { "On":"'", "Fr":4, "To":1 }, { "On":"\", "Fr":4, "To":2 }, { "On":"e", "Fr":4, "To":0 }, { "On":"f", "Fr":4, "To":0 }, { "On":"o", "Fr":4, "To":5 }, { "On":"", "Fr":4, "To":3 }, { "On":"", "Fr":4, "To":0 }]}
 { "Sn":5,  "Edge":[ { "On":"'", "Fr":5, "To":1 }, { "On":"\", "Fr":5, "To":2 }, { "On":"e", "Fr":5, "To":0 }, { "On":"f", "Fr":5, "To":6 }, { "On":"o", "Fr":5, "To":0 }, { "On":"", "Fr":5, "To":3 }, { "On":"", "Fr":5, "To":0 }]}
 { "Sn":6,  "Edge":[ { "On":"'", "Fr":6, "To":1 }, { "On":"\", "Fr":6, "To":2 }, { "On":"e", "Fr":6, "To":0 }, { "On":"f", "Fr":6, "To":0 }, { "On":"o", "Fr":6, "To":0 }, { "On":"", "Fr":6, "To":3 }, { "On":"", "Fr":6, "To":0 }]}
]}
