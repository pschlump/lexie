package lexie

/*
type SetOfStrType struct {
	ARe    string
	ItemId int
}

var SetOfStr_000 = []SetOfStrType{
	{ARe: "a", ItemId: 1},
	{ARe: "ab", ItemId: 2},
	{ARe: "bb", ItemId: 3},
	{ARe: "abb", ItemId: 4},
}

*/
/*
func AddToRv(ndfa *[]NDFA_Type, ss string, isTerm bool) {
	last = len(ndfa) - 1
	ndfa[last].Next[ss] = 0
}

func MakeNDFA(input []SetOfReType) (rv []NDFA_Type) {
	rv = make([]NDFA_Type, 0, 100)
	for ii, vv := range input {
		for jj, ww := range vv.ARe {
			ss := ww[jj : jj+1] // Should be a rune
			ss_len = 1          // Rune Length
			AddToRv(&rv, ss, (jj+1) == len(vv.ARe))
		}
	}
}
*/
