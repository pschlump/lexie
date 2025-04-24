Machine:
	
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
	
	
{"Input":"Lex-Machine-0", "Rv":1, "Start": 0, "States":[
 { "Sn":0,  "Edge":[ { "On":"e", "Fr":0, "To":1 }, { "On":"g", "Fr":0, "To":2 }, { "On":"r", "Fr":0, "To":1 }, { "On":"s", "Fr":0, "To":3 }, { "On":"t", "Fr":0, "To":1 }, { "On":"", "Fr":0, "To":1 }]}
 { "Sn":1,  "Edge":[ ]}
 { "Sn":2,  "Edge":[ { "On":"e", "Fr":2, "To":4 }]}
 { "Sn":3,  "Edge":[ { "On":"e", "Fr":3, "To":5 }]}
 { "Sn":4,  "Edge":[ { "On":"t", "Fr":4, "To":6 }]}
 { "Sn":5,  "Edge":[ { "On":"t", "Fr":5, "To":7 }]}
 { "Sn":6,  "Term":1,  "Edge":[ { "On":"e", "Fr":6, "To":8 }]}
 { "Sn":7,  "Term":3,  "Edge":[ { "On":"e", "Fr":7, "To":9 }]}
 { "Sn":8,  "Edge":[ { "On":"r", "Fr":8, "To":10 }]}
 { "Sn":9,  "Edge":[ { "On":"r", "Fr":9, "To":11 }]}
 { "Sn":10,  "Term":2,  "Edge":[ ]}
 { "Sn":11,  "Term":5,  "Edge":[ ]}
]}
