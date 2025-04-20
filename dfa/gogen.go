package dfa

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/pschlump/dbgo"
	"github.com/pschlump/filelib"
	"github.com/pschlump/lexie/in"
	"github.com/pschlump/pluto/g_lib"
)

// Similar to  ../in/in.go:964 func DumpTokenMap() {
func (lex *Lexie) GenerateTokenMap(fn string) {
	pkgName := lex.getGoPackage()

	dn := filepath.Dir(fn)
	if !g_lib.InArray(dn, []string{".", "..", "/"}) {
		os.MkdirAll(dn, 0755)
	}

	fp, err := filelib.Fopen(fn, "w")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to open %s for output: %s\n", fn, err)
		return
	}

	// var Tok_map = make(map[int]string)
	fmt.Fprintf(fp, "package %s\n\n", pkgName)
	fmt.Fprintf(fp, "import \"fmt\"\n")
	fmt.Fprintf(fp, "\ntype TokenType int\n")
	fmt.Fprintf(fp, "\nconst (\n")
	for kk, vv := range in.Tok_map {
		fmt.Fprintf(fp, "\t%s TokenType = %d\n", vv, kk)
	}
	fmt.Fprintf(fp, ")\n\n")

	fmt.Fprintf(fp, `
func (tt TokenType)String() string {
	switch tt {
`)
	for kk, vv := range in.Tok_map {
		fmt.Fprintf(fp, `
	case %s: /* %d */
		return %q
`, vv, kk, vv)
	}
	fmt.Fprintf(fp, `
	default:
		return fmt.Sprintf ( "--unknown TokenType %%d--", int(tt) )
	}
}
`)

	fp.Close()

	out, err := exec.Command("goimorts", "-w", fn).Output()
	if err != nil {
		dbgo.Fprintf(os.Stderr, "%(red)Error%(reset) from goimports, %s, %(LF)\n", err)
		dbgo.Fprintf(os.Stderr, "Output: %s\n", out)
	}
}

/*
From ../in/in.go
// Lookup will search a specific defined type, Machine, Tokens etc, for the named item 't'.
func (Im *ImType) Lookup(DefType string, t string) (int, error) {
	if validateDefType(DefType) {
		dd := Im.Def.DefsAre[DefType]
		// fmt.Printf("In %s Looking for %s\n", s, t)
		if v, ok := dd.NameValue[t]; ok {
			// fmt.Printf("In %s Found for %s=%d\n", s, t, v)
			return v, nil
		}
	}
	return -1, fmt.Errorf("Missing - did not find the specified key >%s< in >%s<\n", t, DefType)
}
*/

func ValidateDefType(DefType string) bool {
	if !g_lib.InArray(DefType, []string{"Tokens", "Machines", "Errors", "ReservedWords", "Options"}) {
		fmt.Printf("Error Invalid $def type -->%s<--, should be one of \"Tokens\", \"Machines\", \"Errors\", \"ReservedWords\" \n", DefType)
		return false
	}
	return true
}

func (lex *Lexie) getGoPackage() string {
	if ValidateDefType("Options") {
		if dd, ok0 := lex.Im.Def.DefsAre["Options"]; ok0 {
			if v, ok := dd.NameValueStr["GoPackageName"]; ok {
				return v
			}
		}
	}
	return "lexieGenerated"
}
