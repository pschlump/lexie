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
	
	
{"Input":"Lex-Machine-2", "Rv":1, "Start": 0, "States":[
 { "Sn":0,  "Edge":[ { "On":"	", "Fr":0, "To":1 }, { "On":"
", "Fr":0, "To":1 }, { "On":"", "Fr":0, "To":1 }, { "On":"", "Fr":0, "To":1 }, { "On":" ", "Fr":0, "To":1 }, { "On":"!", "Fr":0, "To":2 }, { "On":""", "Fr":0, "To":3 }, { "On":"%", "Fr":0, "To":4 }, { "On":"&", "Fr":0, "To":5 }, { "On":"'", "Fr":0, "To":6 }, { "On":"(", "Fr":0, "To":7 }, { "On":")", "Fr":0, "To":8 }, { "On":"*", "Fr":0, "To":9 }, { "On":"+", "Fr":0, "To":10 }, { "On":",", "Fr":0, "To":11 }, { "On":"-", "Fr":0, "To":12 }, { "On":".", "Fr":0, "To":13 }, { "On":"/", "Fr":0, "To":14 }, { "On":":", "Fr":0, "To":15 }, { "On":"<", "Fr":0, "To":16 }, { "On":"=", "Fr":0, "To":17 }, { "On":">", "Fr":0, "To":18 }, { "On":"?", "Fr":0, "To":19 }, { "On":"[", "Fr":0, "To":20 }, { "On":"]", "Fr":0, "To":21 }, { "On":"^", "Fr":0, "To":22 }, { "On":"_", "Fr":0, "To":23 }, { "On":"`", "Fr":0, "To":24 }, { "On":"e", "Fr":0, "To":25 }, { "On":"f", "Fr":0, "To":25 }, { "On":"o", "Fr":0, "To":25 }, { "On":"{", "Fr":0, "To":26 }, { "On":"|", "Fr":0, "To":27 }, { "On":"}", "Fr":0, "To":28 }, { "On":"~", "Fr":0, "To":29 }, { "On":"", "Fr":0, "To":30 }, { "On":"", "Fr":0, "To":31 }, { "On":"", "Fr":0, "To":23 }, { "On":"", "Fr":0, "To":23 }, { "On":"", "Fr":0, "To":25 }]}
 { "Sn":1,  "Term":71,  "Edge":[ ]}
 { "Sn":2,  "Term":23,  "Edge":[ { "On":"=", "Fr":2, "To":32 }]}
 { "Sn":3,  "Term":71,  "Edge":[ ]}
 { "Sn":4,  "Term":27,  "Edge":[ { "On":"=", "Fr":4, "To":33 }, { "On":"}", "Fr":4, "To":34 }]}
 { "Sn":5,  "Term":41,  "Edge":[ { "On":"&", "Fr":5, "To":35 }, { "On":"=", "Fr":5, "To":36 }]}
 { "Sn":6,  "Term":71,  "Edge":[ ]}
 { "Sn":7,  "Term":12,  "Edge":[ ]}
 { "Sn":8,  "Term":13,  "Edge":[ ]}
 { "Sn":9,  "Term":16,  "Edge":[ { "On":"=", "Fr":9, "To":37 }]}
 { "Sn":10,  "Term":14,  "Edge":[ { "On":"+", "Fr":10, "To":38 }, { "On":"=", "Fr":10, "To":39 }]}
 { "Sn":11,  "Term":21,  "Edge":[ ]}
 { "Sn":12,  "Term":15,  "Edge":[ { "On":"-", "Fr":12, "To":40 }, { "On":"=", "Fr":12, "To":41 }]}
 { "Sn":13,  "Term":22,  "Edge":[ ]}
 { "Sn":14,  "Term":19,  "Edge":[ { "On":"=", "Fr":14, "To":42 }]}
 { "Sn":15,  "Term":25,  "Edge":[ { "On":"=", "Fr":15, "To":43 }]}
 { "Sn":16,  "Term":17,  "Edge":[ { "On":"<", "Fr":16, "To":44 }, { "On":"=", "Fr":16, "To":45 }, { "On":">", "Fr":16, "To":46 }]}
 { "Sn":17,  "Term":26,  "Edge":[ { "On":"=", "Fr":17, "To":47 }, { "On":"~", "Fr":17, "To":48 }]}
 { "Sn":18,  "Term":18,  "Edge":[ { "On":"=", "Fr":18, "To":49 }, { "On":">", "Fr":18, "To":50 }]}
 { "Sn":19,  "Term":59,  "Edge":[ { "On":"=", "Fr":19, "To":51 }]}
 { "Sn":20,  "Term":38,  "Edge":[ ]}
 { "Sn":21,  "Term":39,  "Edge":[ ]}
 { "Sn":22,  "Term":20,  "Edge":[ { "On":"=", "Fr":22, "To":52 }]}
 { "Sn":23,  "Term":70,  "Edge":[ { "On":"_", "Fr":23, "To":53 }, { "On":"", "Fr":23, "To":53 }, { "On":"", "Fr":23, "To":53 }, { "On":"", "Fr":23, "To":53 }]}
 { "Sn":24,  "Term":71,  "Edge":[ ]}
 { "Sn":25,  "Edge":[ ]}
 { "Sn":26,  "Term":53,  "Edge":[ { "On":"%", "Fr":26, "To":54 }, { "On":"{", "Fr":26, "To":55 }]}
 { "Sn":27,  "Term":24,  "Edge":[ { "On":"=", "Fr":27, "To":56 }, { "On":"|", "Fr":27, "To":57 }]}
 { "Sn":28,  "Term":54,  "Edge":[ { "On":"}", "Fr":28, "To":58 }]}
 { "Sn":29,  "Term":40,  "Edge":[ { "On":"=", "Fr":29, "To":59 }, { "On":"~", "Fr":29, "To":60 }]}
 { "Sn":30,  "Edge":[ { "On":"e", "Fr":30, "To":61 }]}
 { "Sn":31,  "Edge":[ ]}
 { "Sn":32,  "Term":10,  "Edge":[ ]}
 { "Sn":33,  "Term":49,  "Edge":[ ]}
 { "Sn":34,  "Term":9,  "Edge":[ ]}
 { "Sn":35,  "Term":4,  "Edge":[ ]}
 { "Sn":36,  "Term":52,  "Edge":[ ]}
 { "Sn":37,  "Term":47,  "Edge":[ ]}
 { "Sn":38,  "Term":61,  "Edge":[ ]}
 { "Sn":39,  "Term":45,  "Edge":[ ]}
 { "Sn":40,  "Term":62,  "Edge":[ ]}
 { "Sn":41,  "Term":46,  "Edge":[ ]}
 { "Sn":42,  "Term":48,  "Edge":[ ]}
 { "Sn":43,  "Term":63,  "Edge":[ ]}
 { "Sn":44,  "Term":43,  "Edge":[ { "On":"=", "Fr":44, "To":62 }]}
 { "Sn":45,  "Term":3,  "Edge":[ ]}
 { "Sn":46,  "Term":10,  "Edge":[ ]}
 { "Sn":47,  "Term":1,  "Edge":[ { "On":"=", "Fr":47, "To":63 }]}
 { "Sn":48,  "Edge":[ { "On":"=", "Fr":48, "To":64 }]}
 { "Sn":49,  "Term":2,  "Edge":[ ]}
 { "Sn":50,  "Term":44,  "Edge":[ { "On":"=", "Fr":50, "To":65 }]}
 { "Sn":51,  "Term":60,  "Edge":[ ]}
 { "Sn":52,  "Term":50,  "Edge":[ ]}
 { "Sn":53,  "Term":70,  "Edge":[ { "On":"_", "Fr":53, "To":53 }, { "On":"", "Fr":53, "To":53 }, { "On":"", "Fr":53, "To":53 }, { "On":"", "Fr":53, "To":53 }]}
 { "Sn":54,  "Term":8,  "Edge":[ ]}
 { "Sn":55,  "Term":6,  "Edge":[ ]}
 { "Sn":56,  "Term":51,  "Edge":[ ]}
 { "Sn":57,  "Term":5,  "Edge":[ ]}
 { "Sn":58,  "Term":65,  "Edge":[ ]}
 { "Sn":59,  "Term":55,  "Edge":[ ]}
 { "Sn":60,  "Term":56,  "Edge":[ ]}
 { "Sn":61,  "Edge":[ { "On":"o", "Fr":61, "To":66 }]}
 { "Sn":62,  "Term":73,  "Edge":[ ]}
 { "Sn":63,  "Term":57,  "Edge":[ ]}
 { "Sn":64,  "Term":58,  "Edge":[ ]}
 { "Sn":65,  "Term":74,  "Edge":[ ]}
 { "Sn":66,  "Edge":[ { "On":"f", "Fr":66, "To":67 }]}
 { "Sn":67,  "Edge":[ ]}
]}
