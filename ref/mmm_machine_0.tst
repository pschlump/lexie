Machine:
	
	$def(Tokens, Tok_CL=1, Tok_PCT=2, Tok_BB=3, Tok_EOF=4)
	
	$def(Machines, S_Init)
	
	$def(Errors, Err_Invalid_Char) 	
	
	$machine(S_Init)
	`%}`					: Rv(Tok_CL) 
	`%`						: Rv(Tok_PCT) 
	`bb`					: Rv(Tok_BB) 
	.						: Error(Err_Invalid_Char)
	$eof					: Rv(Tok_EOF)
	$end
	
	
{"Input":"Lex-Machine-0", "Rv":1, "Start": 0, "States":[
 { "Sn":0,  "Edge":[ { "On":"%", "Fr":0, "To":1 }, { "On":"b", "Fr":0, "To":2 }, { "On":"e", "Fr":0, "To":3 }, { "On":"f", "Fr":0, "To":3 }, { "On":"o", "Fr":0, "To":3 }, { "On":"}", "Fr":0, "To":3 }, { "On":"", "Fr":0, "To":4 }, { "On":"", "Fr":0, "To":3 }]}
 { "Sn":1,  "Term":2,  "Edge":[ { "On":"}", "Fr":1, "To":5 }]}
 { "Sn":2,  "Edge":[ { "On":"b", "Fr":2, "To":6 }]}
 { "Sn":3,  "Edge":[ ]}
 { "Sn":4,  "Edge":[ { "On":"e", "Fr":4, "To":7 }]}
 { "Sn":5,  "Term":1,  "Edge":[ ]}
 { "Sn":6,  "Term":3,  "Edge":[ ]}
 { "Sn":7,  "Edge":[ { "On":"o", "Fr":7, "To":8 }]}
 { "Sn":8,  "Edge":[ { "On":"f", "Fr":8, "To":9 }]}
 { "Sn":9,  "Term":4,  "Edge":[ ]}
]}
