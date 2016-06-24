package main

import (
	"fmt"
	"io/ioutil"
	"unicode/utf8"
)

func main() {
	b, _ := ioutil.ReadFile("utf8.ify.1")
	r := ""
	// for c, sz := range s {
	for i := 0; i < len(b); {
		c, sz := utf8.DecodeRune(b[i:])
		i += sz
		fmt.Printf("c=%v sz=%d\n", c, sz)
		if sz == 1 && c <= 0x7f && c >= '@' {
			r += string(c)
		} else {
			r += fmt.Sprintf("\\u%04x", int(c))
		}
	}
	fmt.Printf("%s\n", r)
}
