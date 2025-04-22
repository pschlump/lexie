
/*
// ------------------------------------------------------------------------------------------------------------------------------------------------
// Plan
// ------------------------------------------------------------------------------------------------------------------------------------------------

+1. Study Django templates
	+2. Understand what they do
	+2. Decide on replace or use code
		http://www.b-list.org/weblog/2007/sep/22/standalone-django-scripts/

+0. What I want is to be able to get a lit of tokens/class/values out as an output from the CLI

+0. Use CLI and turn off all extraneous output -- Validate matches on a bunch of machines

0. Output machine as a .mlex format - machine lex format -m <fn>
0. Output machine as a .go format --fmt go -o <fn>
0. add "$cfg" to input - for config params like size of smap
0. add a SetConfig() to code to set the config params

+1. Use in Ringo
+2. The Pongo2 v.s. Ringo test
+3. Coverage Testing
4. Benchmarks
4. Memory usage - eliminate dynamic allocations
+4. Internal comments (doc.go) etc.
+5. As a "go-routine"



// ------------------------------------------------------------------------------------------------------------------------------------------------
//
// TODO:
//		1. Failure to match causes a fall-out-the-bottom and quit behavior - change to a reset to 0 and continue to EOF
//				1. Warnings about bad tokens are not reported a {{ token should produce an error about invalid range.
// 					1. xyzzyStackEmpty - if stack not empty then - pop and continue running on different data??a
//		1. In nfa->dfa conversion may end up with non-distinct terminal values - use one that is first seen? -- Lex implies this.
//		1. Rune Fix ( ../re/re.go:800 ) + test cases
//					 xyzzyRune  TODO - this code is horribley non-utf8 compatable at this moment in time.
//
//	Err-Output:
//			1. Tests in ../in/in.go - are not automated - fix to check results of input to be correct
//				5. No test case for A_Reset
//				7. No Test case for this error: --- CCL --- 0-0 is an invalid CCL, 1-0 also etc. -- Report
//
//	Feature
//			1. Options on size of smap
//
// ------------------------------------------------------------------------------------------------------------------------------------------------
*/
/*
// ------------------------------------------------------------------------------------------------------------------------------------------------
// Plan
// ------------------------------------------------------------------------------------------------------------------------------------------------

+1. Study Django templates
	+2. Understand what they do
	+2. Decide on replace or use code
		http://www.b-list.org/weblog/2007/sep/22/standalone-django-scripts/

+0. What I want is to be able to get a lit of tokens/class/values out as an output from the CLI

+0. Use CLI and turn off all extraneous output -- Validate matches on a bunch of machines

0. Output machine as a .mlex format - machine lex format -m <fn>
0. Output machine as a .go format --fmt go -o <fn>
0. add "$cfg" to input - for config params like size of smap
0. add a SetConfig() to code to set the config params

+1. Use in Ringo
+2. The Pongo2 v.s. Ringo test
+3. Coverage Testing
4. Benchmarks
4. Memory usage - eliminate dynamic allocations
+4. Internal comments (doc.go) etc.
+5. As a "go-routine"



// ------------------------------------------------------------------------------------------------------------------------------------------------
//
// TODO:
//		1. Failure to match causes a fall-out-the-bottom and quit behavior - change to a reset to 0 and continue to EOF
//				1. Warnings about bad tokens are not reported a {{ token should produce an error about invalid range.
// 					1. xyzzyStackEmpty - if stack not empty then - pop and continue running on different data??a
//		1. In nfa->dfa conversion may end up with non-distinct terminal values - use one that is first seen? -- Lex implies this.
//		1. Rune Fix ( ../re/re.go:800 ) + test cases
//					 xyzzyRune  TODO - this code is horribley non-utf8 compatable at this moment in time.
//
//	Err-Output:
//			1. Tests in ../in/in.go - are not automated - fix to check results of input to be correct
//				5. No test case for A_Reset
//				7. No Test case for this error: --- CCL --- 0-0 is an invalid CCL, 1-0 also etc. -- Report
//
//	Feature
//			1. Options on size of smap
//
// ------------------------------------------------------------------------------------------------------------------------------------------------
*/
