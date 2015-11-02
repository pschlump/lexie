
$machine(S_Test_01)
a012					: Rv(Tok_ID)
[0-9]+					: Rv(Tok_NUM)
`"`						: Rv(Tok_Str0)
`\`						: Call(S_Quote) Rv(Tok_Call)
.*\"					: Rv(Tok_StrBody)
$eof					: Error(Err_EOF_In_String) 	// Error - State Machine Exit
$end


