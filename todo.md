
1. 



./dfa/nm.go: 
	// xyzzy3330
	// If "greedy" in this state then need to look ahead, on a terminal state if there are future states
	// and if the next character will allow you to move to a future state then save (stash) the current
	// result and position and continue running the machine.  A future longer match will result in replacing
	// the current match.  If it is "greedy" then repeat the process.  If not "greedy" then return.
	// If you reach a fail-to-match ( current rune has no future state ), then return the saved state
	// and re-position to the point at which the saved state had saved it.



./re/ - [:alpha:] - how is this mapped and how is it correct?  ( see ./re/re_test.go "xyzzy821" )

./dfa/dfa_03_test.go: // Lr2Type{StrTokNo: "Tok_CL", Match: "%}"}, // 0 ***correct*** xyzzy847





3. Documentation
5. Website for this
