
$def(Tokens, Tok_CL=1, Tok_PCT=2, Tok_BB=3, Tok_EOF=4)

$def(Machines, S_Init)

$def(Errors, Err_Invalid_Char) 	

$def(Options, GoPackageName=test03package)

$machine(S_Init)
`%}`					: Rv(Tok_CL) 
`%`						: Rv(Tok_PCT) 
`bb`					: Rv(Tok_BB) 
.						: Error(Err_Invalid_Char)
$end

