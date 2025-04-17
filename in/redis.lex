
$def(Tokens, Tok_GET=1, Tok_SET=2, Tok_DEL=3, Tok_TTL=4, Tok_EXPIRE=5, Tok_QUIT=6, Tok_EOF=7, Tok_AN_OPT)

$def(Machines, S_Init, S_OPTS)

$def(Errors, Err_Invalid_Command, Err_EOF_In_String, Err_Invalid_Char) 	

$def(Options, GoPackage=redisScanner)

$machine(S_Init)
`get`					: Rv(Tok_GET) Call(S_OPTS)
`set`					: Rv(Tok_SET) Call(S_OPTS)
`del`					: Rv(Tok_DEL) Call(S_OPTS)
`ttl`					: Rv(Tok_TTL) Call(S_OPTS)
`expire`				: Rv(Tok_EXPIRE) Call(S_OPTS)
`quit`					: Rv(Tok_QUIT) 
"[^ ]*"					: Error(Err_Invalid_Command)
.						: Error(Err_Invalid_Char)
$eof					: Rv(Tok_EOF)
$end

$machine(S_OPTS)
"[a-zA-Z][a-zA-Z0-9]*"	: Rv(Tok_AN_OPT)
"[ \t\f]"				: Ignore()						// Implies a terminal Rv(Tok) + no return and go to state 0
.						: Return()
$eof					: Error(Err_EOF_In_String) 		// Error - State Machine Exit
$end

