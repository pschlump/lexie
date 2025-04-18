Machine:
	
	$def(Tokens, Tok_CL=1, Tok_PCT=2, Tok_BB=3, Tok_EOF=4)
	
	$def(Machines, S_Init)
	
	$def(Errors, Err_Invalid_Char) 	
	
	$machine(S_Init)
	`%}`					: Rv(Tok_CL) 
	`%`						: Rv(Tok_PCT) 
	`bb`					: Rv(Tok_BB) 
	.						: Error(Err_Invalid_Char)
	$end
	
	
{"Input":"Lex-Machine-0", "Rv":1, "Start": 0, "States":[
 { "Sn":0,  "Edge":[ { "On":"%", "Fr":0, "To":1 }, { "On":"b", "Fr":0, "To":2 }, { "On":"}", "Fr":0, "To":3 }, { "On":"", "Fr":0, "To":3 }]}
 { "Sn":1,  "Term":2,  "Edge":[ { "On":"}", "Fr":1, "To":4 }]}
 { "Sn":2,  "Edge":[ { "On":"b", "Fr":2, "To":5 }]}
 { "Sn":3,  "Edge":[ ]}
 { "Sn":4,  "Term":1,  "Edge":[ ]}
 { "Sn":5,  "Term":3,  "Edge":[ ]}
]}
