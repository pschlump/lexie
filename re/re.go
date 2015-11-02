//
// R E - Part of Lexie Lexical Generation System
//
// (C) Philip Schlump, 2014-2015.
// Version: 1.0.8
// BuildNo: 203
//
// Special Thanks to 2C-Why, LLC for supporting this project.
//

package re

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"unicode/utf8"

	"../com"
)

type ReTreeNodeType struct {
	Item     string           // Set of Runes
	Mm       int              // { m, n }
	Nn       int              //  { m, n }
	LR_Tok   LR_TokType       // Node Type
	Children []ReTreeNodeType // Children of this node
}

type LexReType struct {
	Buf   string          // Holds the RE being parsed
	Pos   int             // Where we are
	Tree  *ReTreeNodeType // Pointer to the top of the tree
	Error []error         // Set of errors
	Sigma string          //
}

type LexReMatcherType struct {
	Sym  string
	Rv   LR_TokType
	Repl string
}

// -- Functions ------------------------------------------------------------------------------------------------------------

func init() {
	if false {
		fmt.Printf("", com.SVarI(nil), com.LF())
	}
}

func NameOfLR_TokType(x LR_TokType) string {
	if t, ok := LR_TokTypeLookup[x]; ok {
		return t
	}
	return fmt.Sprintf("**unk-LR_Token = %d **", x)
}

func NewLexReType() *LexReType {
	return &LexReType{
		Buf:  "",
		Pos:  0,
		Tree: &ReTreeNodeType{Item: "", LR_Tok: LR_null, Children: make([]ReTreeNodeType, 0, 10)},
	}
}

func (lr *LexReType) SetBuf(s string) {
	lr.Buf = s
	lr.Pos = 0
	lr.Tree.Children = make([]ReTreeNodeType, 0, 10)
}

func (lr *LexReType) Next() (ss string, cl LR_TokType) {
	var rr rune
	var sz int
	if lr.Pos < len(lr.Buf) {
		for _, vv := range LexReMatcher {
			if strings.HasPrefix(lr.Buf[lr.Pos:], vv.Sym) {
				ss = vv.Sym
				cl = vv.Rv
				lr.Pos += len(vv.Sym)
				if vv.Repl != "" {
					ss = vv.Repl
				}
				goto done
			}
		}
		// ss = lr.Buf[lr.Pos : lr.Pos+1]
		rr, sz = utf8.DecodeRune([]byte(lr.Buf[lr.Pos:]))
		ss += string(rr)
		cl = LR_Text
		// lr.Pos += 1
		lr.Pos += sz
	done:
	} else {
		ss = ""
		cl = LR_EOF
	}
	return
}

func (lr *LexReType) Warn(s string) {
	if com.DbOn("OutputErrors") {
		fmt.Printf("Warning: %s\n", s)
	}
	lr.Error = append(lr.Error, errors.New(fmt.Sprintf("Warning: %s, %s", s, com.LF(2))))
}

func (lr *LexReType) Err(s string) {
	if com.DbOn("OutputErrors") {
		fmt.Printf("Error: %s\n", s)
	}
	lr.Error = append(lr.Error, errors.New(fmt.Sprintf("Error: %s, %s", s, com.LF(2))))
}

func NewReTreeNodeType() *ReTreeNodeType {
	return &ReTreeNodeType{
		Item:   "",
		LR_Tok: LR_null,
	}
}

func (lr *LexReType) parseCCL(depth int, ww LR_TokType) (tree ReTreeNodeType) {
	pos := 0
	com.DbPrintf("db2", "parseCCL Top: depth=%d,  %d=%s\n", depth, ww, NameOfLR_TokType(ww))
	s := ""
	dx := 0
	marked := false
	c, w := lr.Next()
	for w != LR_EOF {
		com.DbPrintf("re2", "Top of parseCCL dx=%d ->%s<-\n", dx, c)
		switch w {
		case LR_MINUS: // -		-- Text if not in CCL and not 1st char in CCL
			fallthrough
		case LR_Text: //			-- Add a node to list, move right
			fallthrough
		case LR_CARROT: // ^		-- BOL
			fallthrough
		case LR_DOT: // .		-- Match any char
			fallthrough
		case LR_STAR: // *		-- Error if 1st char, else take prev item from list, star and replace it.
			fallthrough
		case LR_PLUS: // +		-- Error if 1st char
			fallthrough
		case LR_QUEST: // ?		-- Error if 1st char
			fallthrough
		case LR_OP_PAR: // (		-- Start of Sub_Re
			fallthrough
		case LR_CL_PAR: // )
			fallthrough
		case LR_OR: // |
			fallthrough
		case LR_OP_BR: // {
			fallthrough
		case LR_CL_BR: // }
			fallthrough
		case LR_COMMA: // ,
			fallthrough
		case LR_DOLLAR: // $
			s += c

		case LR_N_CCL: // [^...]	-- N_CCL Node
			marked = true
			com.DbPrintf("re2", "    incr to %d, %s\n", dx, com.LF())
			s += c // Add To CCL

		case LR_CCL: // [...]	-- CCL Node (Above)
			marked = true
			com.DbPrintf("re2", "    incr to %d, %s\n", dx, com.LF())
			s += c // Add To CCL

		case LR_E_CCL:
			dx--
			com.DbPrintf("re2", "    decr to %d, %s\n", dx, com.LF())
			if dx < 0 {
				tree.Item = expandCCL(s) // Do Something, Return
				tree.LR_Tok = ww
				return
			} else {
				s += c //
			}

		default:
			lr.Error = append(lr.Error, errors.New(fmt.Sprintf("Unreacable Code, invalid token in parseCCL = %d, %s", w, com.LF())))
			tree.Item = expandCCL(s) // Do Something, Return
			tree.LR_Tok = ww
			return
		}
		pos++
		c, w = lr.Next()

		if c == ":" && marked {
			marked = false
			dx++
		}
	}
	if w == LR_EOF {
		lr.Err("EOF found in Character Class [...] or [^...].")
	}
	tree.Item = expandCCL(s)
	tree.LR_Tok = ww
	return
}

// tree.Item, tree.Mm, tree.Nn = parseIteratorString(s) // Do Something, Return
/*
Should build a test for this
	m,n= ->2,3<- 2 3
	m,n= ->3,2<- 3 2
	m,n= ->3,<- 3 9999999999
	m,n= ->,2<- 0 2
	m,n= ->2,3<- 2 3
	m,n= ->3,2<- 3 2
	m,n= ->3,<- 3 9999999999
	m,n= ->,2<- 0 2
	m,n= ->,2<- 0 2
	m,n= ->,<- 0 9999999999
*/
func (lr *LexReType) parseIteratorString(s string) (mm int, nn int) {
	var err error
	mm, nn = 1, 1
	com := strings.Index(s, ",")
	end := strings.Index(s, "}")
	if end != -1 {
		s = s[0:end]
	}
	if com == -1 {
		mm, err = strconv.Atoi(s)
		if err == nil {
			nn = mm
		}
	} else if com == 0 {
		mm = 0
		nn = InfiniteIteration
		p2 := s[com+1:]
		if len(p2) > 0 {
			nn, err = strconv.Atoi(p2)
			if err != nil {
				nn = 1
			}
		}
	} else if com+1 >= len(s) {
		nn = InfiniteIteration
		p1 := s[0:com]
		mm, err = strconv.Atoi(p1)
		if err != nil {
			mm = 0
		}
	} else {
		nn = InfiniteIteration
		p1 := s[0:com]
		p2 := s[com+1:]
		mm, err = strconv.Atoi(p1)
		if err != nil {
			mm = 0
		}
		if len(p2) > 0 {
			nn, err = strconv.Atoi(p2)
			if err != nil {
				nn = 1
			}
		}
	}
	// fmt.Printf("m,n= ->%s<- %d %d\n", s, mm, nn)
	return
}

// mm, nn := lr.parseIterator ( depth+1 )
func (lr *LexReType) parseIterator(depth int) (tree ReTreeNodeType) {
	pos := 0
	com.DbPrintf("db2", "parseIterator Top: depth=%d\n", depth)
	s := ""
	c, w := lr.Next()
	for w != LR_EOF {
		switch w {
		case LR_MINUS: // -		-- Text if not in CCL and not 1st char in CCL
			fallthrough
		case LR_CARROT: // ^		-- BOL
			fallthrough
		case LR_DOT: // .		-- Match any char
			fallthrough
		case LR_STAR: // *		-- Error if 1st char, else take prev item from list, star and replace it.
			fallthrough
		case LR_PLUS: // +		-- Error if 1st char
			fallthrough
		case LR_QUEST: // ?		-- Error if 1st char
			fallthrough
		case LR_OP_PAR: // (		-- Start of Sub_Re
			fallthrough
		case LR_CL_PAR: // )
			fallthrough
		case LR_CCL: // [...]	-- CCL Node (Above)
			fallthrough
		case LR_OR: // |
			fallthrough
		case LR_OP_BR: // {
			fallthrough
		case LR_DOLLAR: // $
			fallthrough
		case LR_E_CCL:
			fallthrough
		case LR_N_CCL: // [^...]	-- N_CCL Node
			lr.Error = append(lr.Error, errors.New(fmt.Sprintf("in parseIterator, Invalid {m,n} - invalid chars found, %s", com.LF())))
			tree.Mm, tree.Nn = 1, 1
			tree.LR_Tok = LR_OP_BR
			return

		case LR_Text: //			-- Add a node to list, move right
			if c[0] >= '0' && c[0] <= '9' || c[0] == ',' {
				s += c // Add To Iterator
			} else {
				lr.Error = append(lr.Error, errors.New(fmt.Sprintf("Unreacable Code, invalid token in parseIterator, %s", com.LF())))
				tree.Mm, tree.Nn = lr.parseIteratorString(s) // Do Something, Return
				tree.Item = "{"
				tree.LR_Tok = LR_OP_BR
				return
			}
		case LR_COMMA: // ,
			s += c // Add To Iterator
		case LR_CL_BR: // }
			tree.Mm, tree.Nn = lr.parseIteratorString(s) // Do Something, Return
			tree.Item = "{"
			tree.LR_Tok = LR_OP_BR
			return
		default:
			lr.Error = append(lr.Error, errors.New(fmt.Sprintf("Unreacable Code, invalid token in parseIterator, %s", com.LF())))
			tree.Mm, tree.Nn = lr.parseIteratorString(s) // Do Something, Return
			tree.Item = "{"
			tree.LR_Tok = LR_OP_BR
			return
		}
		pos++
		c, w = lr.Next()
	}
	if w == LR_EOF {
		lr.Err("EOF found in Iterator {m,n}.")
	}
	tree.Mm, tree.Nn = lr.parseIteratorString(s) // Do Something, Return
	tree.Item = "{"
	tree.LR_Tok = LR_OP_BR
	return
}

func N4Blanks(n int) (rv string) {
	rv = ""
	for i := 0; i < n; i++ {
		rv += "    "
	}
	return
}

//
// What can I see at the top of a RE
//
//	LR_Text                     //
//	LR_EOF                      //
//	LR_DOT                      // .		-- Match any char
//	LR_STAR                     // *		-- Error if 1st char
//	LR_PLUS                     // +		-- Error if 1st char
//	LR_QUEST                    // ?		-- Error if 1st char
//	LR_B_CCL                    // [		-- Start of CCL Node
//	LR_E_CCL                    // ]
//	LR_OP_PAR                   // (		-- Start of Sub_Re
//	LR_CL_PAR                   // )
//	LR_CCL                      // [...]	-- CCL Node (Above)
//	LR_N_CCL                    // [^...]	-- N_CCL Node
//	LR_CARROT                   // ^		-- BOL
//	LR_MINUS                    // -		-- Text if not in CCL and not 1st char in CCL
//
//	Item     string
//	LR_Tok   LR_TokType
//	Children []*ReTreeNodeType
//	Next     *ReTreeNodeType
//
func (lr *LexReType) DumpParseNodesChild(ch []ReTreeNodeType, d int) {
	com.DbPrintf("DumpParseNodes", "\n%sDumpParseNodesChild: At %s\n", N4Blanks(d), com.LF())
	for ii, vv := range ch {
		com.DbPrintf("DumpParseNodes", "%sat %s [step %d] ", N4Blanks(d), com.LF(), ii)
		com.DbPrintf("DumpParseNodes", "Item: [%s] %d=%s, N-Children=%d\n", vv.Item, vv.LR_Tok, NameOfLR_TokType(vv.LR_Tok), len(vv.Children))
		if len(vv.Children) > 0 {
			lr.DumpParseNodesChild(vv.Children, d+1)
		}
	}
	com.DbPrintf("DumpParseNodes", "%sDumpParseNodesChild: Done %s\n\n", N4Blanks(d), com.LF())
}

func (lr *LexReType) DumpParseNodes() {
	com.DbPrintf("DumpParseNodes", "\nDumpParseNodes: At %s\n", com.LF())
	for ii, vv := range lr.Tree.Children {
		com.DbPrintf("DumpParseNodes", "at %s [step %d] ", com.LF(), ii)
		com.DbPrintf("DumpParseNodes", "Item: [%s] %d=%s, N-Children=%d\n", vv.Item, vv.LR_Tok, NameOfLR_TokType(vv.LR_Tok), len(vv.Children))
		if len(vv.Children) > 0 {
			lr.DumpParseNodesChild(vv.Children, 1)
		}
	}
	com.DbPrintf("DumpParseNodes", "DumpParseNodes: Done %s\n\n", com.LF())
	com.DbPrintf("DumpParseNodesX", "DumpParseNodes: %s\n\n", com.SVarI(lr.Tree))
}

func (lr *LexReType) CalcLengthChild(tree *ReTreeNodeType, d int) (x int, hard bool) {
	t := 0

	hard = false
	if d == 1 {
		com.DbPrintf("CalcLength", "CalcLengthChild at top: %s\n\n", com.SVarI(tree))
	}

	switch tree.LR_Tok {
	case LR_null: //
		for jj := range tree.Children {
			t, hard = lr.CalcLengthChild(&tree.Children[jj], d+1)
			x += t
		}
	case LR_Text: //
		// com.DbPrintf("CalcLength", "Len of item(%s) = %d, %s\n", tree.Item, len(tree.Item), com.LF())
		x += len(tree.Item)
		hard = true
	case LR_EOF: //
		hard = true
	case LR_DOT: // .
		x += 1
	case LR_STAR: // *
		x = 0
		// com.DbPrintf("CalcLength", "After * x = %d, hard=%v\n", x, hard)
	case LR_PLUS: // +
		// patch to fix the problem with [0-9]+ not working -- In reality the length is only if it is a "FIXED" length, 0 else
		//		if len(tree.Children) > 0 {
		//			t, hard = lr.CalcLengthChild(&tree.Children[0], d+1)
		//			x += t
		//		}
		//		hard = true
		x = 0
	case LR_QUEST: // ?
		x = 0
	case LR_OP_BR: // { 			// {m,n} - need to calculate length of ( m times, length of children
		x = 0
	case LR_OP_PAR: // (
		if len(tree.Children) > 0 {
			t, hard = lr.CalcLengthChild(&tree.Children[0], d+1)
			x += t
		}
		// com.DbPrintf("CalcLength", "After ( x = %d, hard=%v\n", x, hard)
	case LR_CL_PAR: // )
		x = 0
	case LR_CCL: // [...]
		x += 1
		hard = true
	case LR_N_CCL: // [^...]
		x += 1
	case LR_E_CCL: // ]
		x += 1
	case LR_CARROT: // ^
		x += 0
		hard = true
	case LR_MINUS: // -
		x += 1
		hard = true
	case LR_DOLLAR: // $
		hard = true

	case LR_OR: // |
		y := -1
		z := 0
		hard = false
		if len(tree.Children) > 0 {
			hard = true
			h := false
			for jj := range tree.Children {
				z, h = lr.CalcLengthChild(&tree.Children[jj], d+1)
				if y == -1 {
					y = z
				} else if y < z {
					y = z
				}
				if !h {
					hard = false
				}
			}
		}
		x += y
		// com.DbPrintf("CalcLength", "After | x = %d, hard = %v\n", x, hard)
	}

	return
}

func (lr *LexReType) CalcLength() (int, bool) {
	x, h := lr.CalcLengthChild(lr.Tree, 1)
	com.DbPrintf("CalcLength", "CalcLength Final Value for Tree = %d, hard=%v\n", x, h)
	return x, h
}

func (lr *LexReType) parseExpression(depth int, d_depth int, xTree *ReTreeNodeType) []ReTreeNodeType {
	//var first *ReTreeNodeType
	//var last *ReTreeNodeType
	pre := strings.Repeat("    ", depth)
	if depth == 0 {
		xTree = lr.Tree
		com.DbPrintf("parseExpression", "%sat %s !!!top!!!, depth=%d \n", pre, com.LF(), depth)
	}
	isFirst := true
	inOr := false
	com.DbPrintf("parseExpression", "%sat %s\n", pre, com.LF())
	c, w := lr.Next()
	for w != LR_EOF {
		com.DbPrintf("parseExpression", "%sat %s !!!top!!!, depth=%d c=->%s<- w=%d %s -- Loop Top -- xTree=%s\n\n",
			pre, com.LF(), depth, c, w, NameOfLR_TokType(w), com.SVarI(xTree))
		switch w {
		case LR_CL_BR: // }
			fallthrough
		case LR_COMMA: // ,
			fallthrough
		case LR_E_CCL:
			fallthrough
		case LR_MINUS: // -		-- Text if not in CCL and not 1st char in CCL
			fallthrough
		case LR_Text: //			-- Add a node to list, move right
			//if true {
			xTree.Children = append(xTree.Children, ReTreeNodeType{Item: c, LR_Tok: LR_Text})
			//} else {
			// // Bad Idea - mucks up '*' and other processing -  To Simplify Tree needs to be done post-generation with Simp-Rules
			//	ll := len(lr.Tree.Children) - 1
			//	if ll >= 0 && lr.Tree.Children[ll].LR_Tok == LR_Text {
			//		lr.Tree.Children[ll].Item += c
			//	} else {
			//		xTree.Children = append(xTree.Children, ReTreeNodeType{Item: c, LR_Tok: LR_Text})
			//	}
			//}

		case LR_CARROT: // ^		-- BOL		-- If at beginning, or after ( or | then BOL - else just text??
			fallthrough
		case LR_DOLLAR: // $		-- BOL		-- If at end, or just before ) or | the EOL - else just text??
			fallthrough
		case LR_DOT: // .		-- Match any char
			xTree.Children = append(xTree.Children, ReTreeNodeType{Item: c, LR_Tok: w})

		case LR_OP_BR: // {
			if isFirst {
				lr.Warn(fmt.Sprintf("Invalid '%s' at beginning of R.E. assumed to be a text character missing esacape.", c))
				xTree.Children = append(xTree.Children, ReTreeNodeType{Item: c, LR_Tok: LR_Text})
			} else {
				ll := len(xTree.Children) - 1
				tmp := xTree.Children[ll]
				newTree := lr.parseIterator(depth + 1)
				if newTree.Mm == 0 && newTree.Nn == InfiniteIteration {
					ll := len(xTree.Children) - 1
					tmp := xTree.Children[ll]
					com.DbPrintf("parseExpression", "%sAT %s, w=%d %s, ll=%d, xTree=%s tmp=%s\n", pre, com.LF(), w, NameOfLR_TokType(w), ll, com.SVarI(xTree), com.SVarI(tmp))
					xTree.Children[ll] = ReTreeNodeType{Item: "*", LR_Tok: LR_STAR, Children: []ReTreeNodeType{tmp}}
				} else {
					if newTree.Mm > newTree.Nn {
						lr.Error = append(lr.Error, errors.New(fmt.Sprintf("Invalid Range, Start is bigger than end, {%d,%d}, %s", newTree.Mm, newTree.Nn, com.LF())))
					}
					com.DbPrintf("parseExpression", "%sAT %s, w=%d %s, ll=%d, xTree=%s tmp=%s\n", pre, com.LF(), w, NameOfLR_TokType(w), ll, com.SVarI(xTree), com.SVarI(tmp))
					// xTree.Children[ll] = ReTreeNodeType{Item: c, LR_Tok: LR_OP_BR, Children: []ReTreeNodeType{tmp}, Mm: newTree.Mm, Nn: newTree.Nn}
					newTree.Children = []ReTreeNodeType{tmp}
					xTree.Children[ll] = newTree
					// CCL: xTree.Children = append(xTree.Children, lr.parseCCL(depth+1, w)) // xyzzy needs work ---------------------------------------------------
				}
				com.DbPrintf("parseExpression", "%sat %s\n", pre, com.LF())
			}

		case LR_STAR: // *		-- Error if 1st char, else take prev item from list, star and replace it.
			fallthrough
		case LR_PLUS: // +		-- Error if 1st char
			fallthrough
		case LR_QUEST: // ?		-- Error if 1st char
			if isFirst {
				lr.Warn(fmt.Sprintf("Invalid '%s' at beginning of R.E. assumed to be a text character missing esacape.", c))
				xTree.Children = append(xTree.Children, ReTreeNodeType{Item: c, LR_Tok: LR_Text})
			} else {
				ll := len(xTree.Children) - 1
				tmp := xTree.Children[ll]
				com.DbPrintf("parseExpression", "%sAT %s, w=%d %s, ll=%d, xTree=%s tmp=%s\n", pre, com.LF(), w, NameOfLR_TokType(w), ll, com.SVarI(xTree), com.SVarI(tmp))
				xTree.Children[ll] = ReTreeNodeType{Item: c, LR_Tok: w, Children: []ReTreeNodeType{tmp}}
				com.DbPrintf("parseExpression", "%sat %s\n", pre, com.LF())
			}

		case LR_OR: // |		n-ary or operator
			com.DbPrintf("parseExpression", "%sat %s\n", pre, com.LF())
			inOr = true

			// Left Machine is collected to be sub-machine == Beginnig-to-current
			// left := xTree.Children // change to be left section back to but not including "|" node - or all if no | node.
			kk := -1
			for jj := len(xTree.Children) - 1; jj >= 0; jj-- {
				if xTree.Children[jj].LR_Tok == LR_OR {
					kk = jj
					break
				}
			}
			if kk == -1 { // No OR tok found
				left := xTree.Children // change to be left section back to but not including "|" node - or all if no | node.
				ll := len(left)
				leftNode := ReTreeNodeType{Item: "", LR_Tok: LR_null, Children: make([]ReTreeNodeType, ll, ll)}
				for jj := range left {
					leftNode.Children[jj] = left[jj]
				}

				newTop := ReTreeNodeType{Item: "|", LR_Tok: LR_OR, Children: make([]ReTreeNodeType, 0, 10)}
				newTop.Children = append(newTop.Children, leftNode) // only if no "or" node, else ref to "or" node
				xTree.Children = xTree.Children[:0]
				xTree.Children = append(xTree.Children, newTop)
				com.DbPrintf("parseExpression", "%sAT %s, w=%d %s, left=%s\n", pre, com.LF(), w, NameOfLR_TokType(w), com.SVarI(left))
			} else {
				if kk >= 0 {
					if kk < len(xTree.Children) {
						tmp := xTree.Children[kk+1:]
						xTree.Children = xTree.Children[0 : kk+1]
						newNode := ReTreeNodeType{Item: "", LR_Tok: LR_null, Children: make([]ReTreeNodeType, len(tmp), len(tmp))}
						for i := 0; i < len(tmp); i++ {
							newNode.Children[i] = tmp[i]
						}
						xTree.Children[kk].Children = append(xTree.Children[kk].Children, newNode)
					}
				}
			}

			// Or node is created like (LR_STAR)
			// Recursive call to parse rest of items at this level
			//newNode := ReTreeNodeType{Item: "", LR_Tok: LR_null, Children: make([]ReTreeNodeType, 1, 10)} // No recursive call
			//lr.parseExpression(depth+1, depth, &newNode.Children[0])
			//newTop.Children = append(newTop.Children, newNode.Children[0])

			// Take results of recursion and put in as RIGHT machine under LR_OR (optimize for N-Tree OR at this point)
			//xTree.Children = append(xTree.Children, newTop)
			//if depth > d_depth {
			//com.DbPrintf("parseExpression", "%sat %s, depth=%d d_detph=%d\n", pre, com.LF(), depth, d_depth)
			//	return xTree.Children
			//}

		case LR_OP_PAR: // (		-- Start of Sub_Re
			com.DbPrintf("parseExpression", "%sat %s\n", pre, com.LF())

			newNode := ReTreeNodeType{Item: c, LR_Tok: LR_OP_PAR, Children: make([]ReTreeNodeType, 1, 10)}

			lr.parseExpression(depth+1, depth+1, &newNode.Children[0])

			newNode.Children[0].Item = c
			newNode.Children[0].LR_Tok = LR_OP_PAR

			xTree.Children = append(xTree.Children, newNode.Children[0])

			com.DbPrintf("parseExpression", "%sat %s\n", pre, com.LF())

		case LR_CL_PAR: // )
			// If in "or" node set - then collect last section to "or" ------------------------ <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
			com.DbPrintf("parseExpression", "%sat %s\n", pre, com.LF())
			if depth == 0 {
				com.DbPrintf("parseExpression", "%sat %s\n", pre, com.LF())
				lr.Warn(fmt.Sprintf("Invalid '%s' at not properly nested.   Assuming that this was to match a character.", c))
				xTree.Children = append(xTree.Children, ReTreeNodeType{Item: c, LR_Tok: LR_Text})
				com.DbPrintf("parseExpression", "%sat %s\n", pre, com.LF())
			} else {
				com.DbPrintf("parseExpression", "%sat %s\n", pre, com.LF())
				if inOr {
					com.DbPrintf("parseExpression", "%sAT Top of new code %s, BOTTOM xTree=%s\n", pre, com.LF(), com.SVarI(xTree))
					kk := -1
					for jj := len(xTree.Children) - 1; jj >= 0; jj-- {
						if xTree.Children[jj].LR_Tok == LR_OR {
							kk = jj
							break
						}
					}
					if kk >= 0 {
						if kk < len(xTree.Children) {
							tmp := xTree.Children[kk+1:]
							xTree.Children = xTree.Children[0 : kk+1]
							newNode := ReTreeNodeType{Item: "", LR_Tok: LR_null, Children: make([]ReTreeNodeType, len(tmp), len(tmp))}
							for i := 0; i < len(tmp); i++ {
								newNode.Children[i] = tmp[i]
							}
							xTree.Children[kk].Children = append(xTree.Children[kk].Children, newNode)
						}
					}
					com.DbPrintf("parseExpression", "%sAT Bo5 of new code %s, BOTTOM xTree=%s\n", pre, com.LF(), com.SVarI(xTree))
				}
				return xTree.Children
			}
			com.DbPrintf("parseExpression", "%sat %s\n", pre, com.LF())
			inOr = false

		case LR_CCL: // [...]	-- CCL Node (Above)
			fallthrough
		case LR_N_CCL: // [^...]	-- N_CCL Node
			xTree.Children = append(xTree.Children, lr.parseCCL(depth+1, w)) // xyzzy needs work ---------------------------------------------------

		default:
			lr.Error = append(lr.Error, errors.New(fmt.Sprintf("Invalid LR Token Type, '%d', '%s', %s", w, NameOfLR_TokType(w), com.LF())))
			return xTree.Children
		}
		isFirst = false
		com.DbPrintf("parseExpression", "%sAT %s, BOTTOM xTree=%s\n", pre, com.LF(), com.SVarI(xTree))
		c, w = lr.Next()
	}
	// If in "or" node set - then collect last section to "or" ------------------------ <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
	if inOr {
		com.DbPrintf("parseExpression", "%sAT Top of new code %s, BOTTOM xTree=%s\n", pre, com.LF(), com.SVarI(xTree))
		kk := -1
		for jj := len(xTree.Children) - 1; jj >= 0; jj-- {
			if xTree.Children[jj].LR_Tok == LR_OR {
				kk = jj
				break
			}
		}
		if kk >= 0 {
			if kk < len(xTree.Children) {
				tmp := xTree.Children[kk+1:]
				xTree.Children = xTree.Children[0 : kk+1]
				newNode := ReTreeNodeType{Item: "", LR_Tok: LR_null, Children: make([]ReTreeNodeType, len(tmp), len(tmp))}
				for i := 0; i < len(tmp); i++ {
					newNode.Children[i] = tmp[i]
				}
				xTree.Children[kk].Children = append(xTree.Children[kk].Children, newNode)
			}
		}
		com.DbPrintf("parseExpression", "%sAT Bo5 of new code %s, BOTTOM xTree=%s\n", pre, com.LF(), com.SVarI(xTree))
	}
	return xTree.Children
}

func (lr *LexReType) ParseRe(ss string) {
	com.DbPrintf("db2", "at %s\n", com.LF())
	lr.SetBuf(ss)
	com.DbPrintf("db2", "at %s\n", com.LF())
	lr.parseExpression(0, 0, nil)
	com.DbPrintf("db2", "at %s\n", com.LF())
}

func expandCCL(s string) (ccl string) {
	ccl = ""
	com.DbPrintf("db2", "at %s\n", com.LF())

	pos := 0
	if len(s) > 0 && s[0:1] == "-" { // Check for leading '-' include in CCL
		ccl += "-"
		pos = 1
	}

	for ii := pos; ii < len(s); ii++ {
		com.DbPrintf("re2", "ii=%d remaining ->%s<-, %s\n", ii, s[ii:], com.LF())
		if strings.HasPrefix(s[ii:], "[:alphnum:]") {
			com.DbPrintf("re2", "   Matched: %s\n", com.LF())
			// ccl += X_ALPHA
			ccl += X_LOWER
			ccl += X_UPPER
			ccl += X_NUMERIC
			ii += len("[:alphnum:]") - 1
		} else if strings.HasPrefix(s[ii:], "[:alpha:]") {
			com.DbPrintf("re2", "   Matched: %s\n", com.LF())
			// ccl += X_ALPHA
			ccl += X_LOWER
			ccl += X_UPPER
			ii += len("[:alpha:]") - 1
		} else if strings.HasPrefix(s[ii:], "[:lower:]") {
			com.DbPrintf("re2", "   Matched: %s\n", com.LF())
			ccl += X_LOWER
			ii += len("[:lower:]") - 1
		} else if strings.HasPrefix(s[ii:], "[:upper:]") {
			com.DbPrintf("re2", "   Matched: %s\n", com.LF())
			ccl += X_UPPER
			ii += len("[:upper:]") - 1
		} else if strings.HasPrefix(s[ii:], "[:numeric:]") {
			com.DbPrintf("re2", "   Matched: %s\n", com.LF())
			ccl += X_NUMERIC
			ii += len("[:numeric:]") - 1
		} else if ii+9 <= len(s) && s[ii:ii+9] == "a-zA-Z0-9" {
			com.DbPrintf("re2", "   Matched: %s\n", com.LF())
			// ccl += X_ALPHA
			ccl += X_LOWER
			ccl += X_UPPER
			ccl += X_NUMERIC
			ii += 8
		} else if ii+6 <= len(s) && s[ii:ii+6] == "a-zA-Z" {
			com.DbPrintf("re2", "   Matched: %s\n", com.LF())
			// ccl += X_ALPHA
			ccl += X_LOWER
			ccl += X_UPPER
			ii += 5
		} else if ii+3 <= len(s) && s[ii:ii+3] == "0-9" {
			com.DbPrintf("re2", "   Matched: %s\n", com.LF())
			// fmt.Printf("matched 0-9 pattern\n")
			ccl += X_NUMERIC
			ii += 2
		} else if ii+3 <= len(s) && s[ii:ii+3] == "a-z" {
			com.DbPrintf("re2", "   Matched: %s\n", com.LF())
			ccl += X_LOWER
			ii += 2
		} else if ii+3 <= len(s) && s[ii:ii+3] == "A-Z" {
			com.DbPrintf("re2", "   Matched: %s\n", com.LF())
			ccl += X_UPPER
			ii += 2
		} else if ii+2 < len(s) && s[ii+1:ii+2] == "-" {
			com.DbPrintf("re2", "   Matched: %s\n", com.LF())
			// xyzzyRune  TODO - this code is horribley non-utf8 compatable at this moment in time.
			e := s[ii+2]
			// fmt.Printf("matched a-b pattern, b=%s e=%s\n", string(s[ii]), string(e))
			if s[ii] >= e {
				fmt.Printf("Error: Poorly formatted character-class, beginning is larger than or equal to end of CCL\n") // Xyzzy - need line number etc
			}
			for b := s[ii]; b <= e; b++ {
				ccl += string(b)
			}
			ii += 2
		} else {
			ccl += s[ii : ii+1]
		}
		// fmt.Printf("bottom ccl now: ->%s<-\n", ccl)
	}

	return
}

func (lr *LexReType) GenerateSigma() (rv string) {
	rv = lr.GenerateSigmaFromTree(lr.Tree, 1)
	lr.Sigma = rv

	uniq := make(map[string]bool)
	for _, rn := range rv {
		uniq[string(rn)] = true
	}

	//fmt.Printf("RE:GenerateSigma: uniq=%v\n", uniq)

	// To store the keys in slice in sorted order
	var keys []string
	for k := range uniq {
		keys = append(keys, k)
	}
	//fmt.Printf("RE:GenerateSigma: keys (unsorted)=%v\n", keys)
	sort.Strings(keys)
	//fmt.Printf("RE:GenerateSigma: keys (sorted)=%v\n", keys)

	rv = ""
	for _, k := range keys {
		rv += k
	}

	lr.Sigma = rv
	return
}

func (lr *LexReType) GenerateSigmaFromTree(tree *ReTreeNodeType, d int) (rv string) {

	switch tree.LR_Tok {
	case LR_OP_PAR: // (
		fallthrough
	case LR_QUEST: // ?
		fallthrough
	case LR_PLUS: // +
		fallthrough
	case LR_STAR: // *
		if len(tree.Children) > 0 {
			rv += lr.GenerateSigmaFromTree(&tree.Children[0], d+1)
		}
	case LR_OR: // |
		fallthrough
	case LR_null: //
		for jj := range tree.Children {
			rv += lr.GenerateSigmaFromTree(&tree.Children[jj], d+1)
		}
	case LR_CARROT: // ^
		fallthrough
	case LR_Text: //
		fallthrough
	case LR_EOF: //
		fallthrough
	case LR_DOT: // .
		fallthrough
	case LR_MINUS: // -
		fallthrough
	case LR_CCL: // [...]
		fallthrough
	case LR_N_CCL: // [^...]			// probably incorrect
		rv += tree.Item
	case LR_CL_PAR: // )
	case LR_E_CCL: // ]
	case LR_DOLLAR: // $
	}

	return
}

/*

// Set of possible input tokens
// Walk the NFA and collect all unique tokens that are not lambda and have a transition
func (nn *NFA_PoolType) GenerateSigma() (s string) {
	uniq := make(map[string]bool)
	s = ""
	for _, vv := range nn.Pool {
		if vv.IsUsed {
			for _, ww := range vv.Next2 {
				if !ww.IsLambda {
					uniq[ww.On] = true
				}
			}

		}
	}

	fmt.Printf("GenerateSigma: uniq=%v\n", uniq)

	// To store the keys in slice in sorted order
	var keys []string
	for k := range uniq {
		keys = append(keys, k)
	}
	fmt.Printf("GenerateSigma: keys (unsorted)=%v\n", keys)
	sort.Strings(keys)
	fmt.Printf("GenerateSigma: keys (sorted)=%v\n", keys)

	for _, k := range keys {
		s += k
	}

	return
}

*/

/* vim: set noai ts=4 sw=4: */
