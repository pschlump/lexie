//
// Sample for M5 parse of macro
//

[a-zA-Z_][a-zA-Z0-9_]*		: Rv(Tok_ID)
`((`						: Rv(Tok_PListStart)
`))`						: Rv(Tok_PListEnd)
.*							: Rv(Tok_HTML)

