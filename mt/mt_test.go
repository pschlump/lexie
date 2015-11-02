package mt

import (
	"fmt"
	"testing"

	"github.com/pschlump/lexie/com"
	"github.com/pschlump/lexie/gen"
)

// "../../../go-lib/sizlib"

// test tree duplicate -------------------------------------------------------------------------------------------------------------------------------
func Test_Mt01(t *testing.T) {

	var bob *MtType
	bob = NewMtType(1, "bob")

	jane := DuplicateTree(bob)

	fmt.Printf("part 1\n")
	fmt.Printf("bob :%s\n", bob.XValue)
	fmt.Printf("jane:%s\n", jane.XValue)

	fmt.Printf("part 2\n")
	bob.XValue = "part 2"
	fmt.Printf("jane:%s\n", jane.XValue)

	if bob.XValue == jane.XValue {
		t.Errorf("Failed to duplicate the tree, case 1\n")
	}

	bob.List = append(bob.List, NewMtType(1, "bob.1"))
	bob.List = append(bob.List, NewMtType(1, "bob.2"))
	bob.List = append(bob.List, NewMtType(1, "bob.3"))
	bob.List[2].List = append(bob.List[2].List, NewMtType(1, "bob.3.1"))
	bob.List[2].List = append(bob.List[2].List, NewMtType(1, "bob.3.2"))
	bob.List[2].List = append(bob.List[2].List, NewMtType(1, "bob.3.3"))
	bob.List[2].List = append(bob.List[2].List, NewMtType(1, "bob.3.4"))
	bob.List[2].List = append(bob.List[2].List, NewMtType(1, "bob.3.5"))
	jane = DuplicateTree(bob)

	if len(jane.List) != len(bob.List) {
		t.Errorf("Failed to duplicate the tree, case 2\n")
	}
	if len(jane.List[2].List) != len(bob.List[2].List) {
		t.Errorf("Failed to duplicate the tree, case 2\n")
	}

	bob.List = bob.List[:0]
	if len(jane.List) != 3 {
		t.Errorf("Failed to duplicate the tree, case 3\n")
	}
	if len(jane.List[2].List) != 5 {
		t.Errorf("Failed to duplicate the tree, case 3\n")
	}

}

// test Replace sub-tree  -------------------------------------------------------------------------------------------------------------------------------
// func ReplaceBlocksWithNew(search_in_tree, new_block *MtType) {
func Test_Mt02(t *testing.T) {
	var bob *MtType
	bob = NewMtType(1, "bob")
	bob.List = append(bob.List, NewMtType(1, "bob.1"))
	bob.List = append(bob.List, NewMtType(1, "bob.2"))
	bob.List = append(bob.List, NewMtType(1, "bob.3"))
	bob.List[2].List = append(bob.List[2].List, NewMtType(gen.Fx_block, "bob.3.1"))
	bob.List[2].List[0].SVal = make([]string, 1, 1)
	bob.List[2].List[0].FxId = gen.Fx_block
	bob.List[2].List[0].SVal[0] = "mike"
	bob.List[2].List[0].List = append(bob.List[2].List[0].List, NewMtType(1, "bob.3.1.1"))
	bob.List[2].List[0].List[0].HTML_Output = "Original Chunk 1"
	bob.List[2].List[0].List = append(bob.List[2].List[0].List, NewMtType(1, "bob.3.1.2"))
	bob.List[2].List[0].List[1].HTML_Output = "Original Chunk 2"
	bob.List[2].List[0].List = append(bob.List[2].List[0].List, NewMtType(1, "bob.3.1.3"))
	bob.List[2].List = append(bob.List[2].List, NewMtType(1, "bob.3.2"))
	bob.List[2].List = append(bob.List[2].List, NewMtType(1, "bob.3.3"))
	bob.List[2].List = append(bob.List[2].List, NewMtType(1, "bob.3.4"))
	bob.List[2].List = append(bob.List[2].List, NewMtType(1, "bob.3.5"))
	bob.List = append(bob.List, NewMtType(1, "bob.4"))
	bob.List = append(bob.List, NewMtType(1, "bob.5"))
	fmt.Printf("bob before change =%s\n\n", com.SVarI(bob))

	var repl *MtType
	repl = NewMtType(gen.Fx_block, "bob")
	repl.SVal = make([]string, 1, 1)
	repl.SVal[0] = "mike"
	repl.List = append(repl.List, NewMtType(1, "repl.1"))
	repl.List = append(repl.List, NewMtType(1, "repl.1.1"))
	repl.List[0].HTML_Output = "Replacement Text 1.1"
	repl.List = append(repl.List, NewMtType(1, "repl.1.2"))
	fmt.Printf("repl before change =%s\n\n", com.SVarI(repl))

	ReplaceBlocksWithNew(&bob, repl)
	fmt.Printf("bob after change =%s\n\n", com.SVarI(bob))

}
