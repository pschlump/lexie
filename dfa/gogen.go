package dfa

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/pschlump/dbgo"
	"github.com/pschlump/filelib"
	"github.com/pschlump/lexie/in"
)

// SImilar to  ../in/in.go:964 func DumpTokenMap() {
func GenerateTokenMap(fn string) {
	pkgName := "xyzzy"

	dn := filepath.Dir(fn)
	os.MkdirAll(dn, 0755)

	fp, err := filelib.Fopen(fn, "w")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to open %s for output: %s\n", fn, err)
		return
	}

	// var Tok_map = make(map[int]string)
	fmt.Fprintf(fp, "package %s\n\n", pkgName)
	fmt.Fprintf(fp, "\nimport \"fmt\"\n")
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
