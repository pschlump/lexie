package main

import (
	"fmt"

	smap "../"
)

func main() {
	var mm smap.SMapType
	mm.MinV = 22
	mm.MaxV = 100
	fmt.Printf("%s\n", mm.String())
	xv := smap.NewSMapType("def ghi\u2022", rune(0xf812))
	fmt.Printf("Results = %s\n", xv.String())
}
