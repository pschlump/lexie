
$def(Tokens, Tok_GET=1, Tok_GETER=2, Tok_SET=3, Tok_EOF=4, Tok_SETER=5)

$def(Machines, S_Init)

$def(Errors, Err_Invalid_Char) 	

$def(Options, GoPackageName=test03package)

$machine(S_Init)
`get`					: Rv(Tok_GET) 
`geter`					: Rv(Tok_GETER) 
`set`					: Rv(Tok_SET) 
`seter`					: Rv(Tok_SETER) 
.						: Error(Err_Invalid_Char)
$end

