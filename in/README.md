
# in - Input

This is the input for Lexie, the lexical analyzer generator.    Ironically
this is a hand coded input lexer.

At the top level the syntax is a series of definitions followed by a set of machins:

```
lexie ::= 
    definitions
    machines
    ;

definitions ::= definitions definition
    | 
    ;


definition ::=
    | def_tokens
    | def_machines
    | def_errors
    | def_reserved_words
    | def_go_package_name
    | def_options
    ;

// $def(Tokens, Tok_L_EQ=1, Tok_GE=2, Tok_LE=3, Tok_L_AND=4, Tok_L_OR=5, Tok_OP_VAR=6, Tok_CL_VAR=7, Tok_OP_BL=8, Tok_CL_BL=9, Tok_NE=10, Tok_NE_LG=11, Tok_OP=12, Tok_CL=13, Tok_PLUS=14, Tok_MINUS=15, Tok_STAR=16, Tok_LT=17, Tok_GT=18, Tok_SLASH=19, Tok_CARRET=20, Tok_COMMA=21, Tok_DOT=22, Tok_EXCLAM=23, Tok_OR=24, Tok_COLON=25, Tok_EQ=26, Tok_PCT=27, Tok_in=28, Tok_and=29, Tok_or=30, Tok_not=31, Tok_true=32, Tok_false=33, Tok_as=34, Tok_export=35, Tok_SS=36, Tok_PIPE=37, Tok_OP_SQ=38, Tok_CL_SQ=39, Tok_TILDE=40, Tok_B_AND=41, Tok_B_OR=42, Tok_S_L=43, Tok_S_R=44, Tok_PLUS_EQ=45, Tok_MINUS_EQ=46, Tok_STAR_EQ=47, Tok_DIV_EQ=48, Tok_MOD_EQ=49, Tok_CAROT_EQ=50, Tok_B_OR_EQ=51, Tok_B_AND_EQ=52, Tok_OP_BRACE=53, Tok_CL_BRACE=54, Tok_TILDE_EQ=55, Tok_TILDE_TILDE=56, Tok_EQ3=57, Tok_APROX_EQ=58, Tok_QUEST=59, Tok_RE_MATCH=60, Tok_PLUS_PLUS=61, Tok_MINUS_MINUS=62, Tok_DCL_VAR=63, Tok_XOR=64)
def_tokens ::= '$def(Tokens,' token_list)
    ;

// $def(Machines, S_Init, S_TAG, S_Common, S_Str0, S_Str1, S_Str2, S_VAR, S_Quote )
def_machines ::= '$def(Machines,' A_list)
    ;

// $def(Errors,  Warn_End_Var_Unexpected, Err_EOF_Tag, Err_EOF_In_String )
def_machines ::= '$def(Errors,' A_list)
    ;

// $def(ReservedWords, and=Tok_L_AND, or=Tok_L_OR, true=Tok_true, false=Tok_false, not=Tok_not, export=Tok_export, in=Tok_in, not=Tok_not, as=Tok_as, bor=Tok_B_OR, band=Tok_B_AND, xor=Tok_XOR )
def_machines ::= '$def(Errors,' rw_list)
    ;

// Comma Seperated list of tokens, each token is NAME=VALUE
token_list ::= token ',' token_list
    | token
    ;


token ::= TOKEN_NAME '=' TOKEN_VALUE
    ;

// List of all of the defined machines.
machines ::= definitions machine
    | machine 
    ;

// a single macine defined in '$def(Machine, MACHINE_NAME)'
machine ::= '$machine(MACINE_NAME'
    ;

``
