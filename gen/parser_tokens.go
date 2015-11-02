package gen

const (
	Tok_Tree_List         = 400
	Tok_Tree_Bool         = 401
	Tok_Tree_Float        = 402
	Tok_Tree_Int          = 403
	Tok_Tree_HTML         = 404
	Tok_Tree_CallFx       = 405
	Tok_Tree_Item         = 406 // An item like {% csrf_token %}
	Tok_Tree_Begin        = 407 // An begin token like {% block <name> %} <name>==Value, ID="block"
	Tok_Tree_End          = 408 // An begin token like {% endblock <name> %} <name>==Value, ID="block"
	Tok_Tree_Macro        = 409 // {% define abc def %}
	Tok_Tree_If           = 410
	Tok_Tree_ElsIf        = 411
	Tok_Tree_Else         = 412
	Tok_Tree_Endif        = 413
	Tok_Tree_For          = 414
	Tok_Tree_Empty        = 415
	Tok_Tree_EndFor       = 416
	Tok_Tree_Ifequal      = 417
	Tok_Tree_Ifnotequal   = 418
	Tok_Tree_Ifchanged    = 419
	Tok_Tree_Ifnotchanged = 420
	Tok_Tree_Comment      = 421
	Tok_Template          = 422
	Tok_Expr              = 423
	Tok_ID_or_Str         = 424
	Tok_Match_Str         = 425
)
