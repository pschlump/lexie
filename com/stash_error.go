//
// D B - Part of Lexie Lexical Generation System
//
// (C) Philip Schlump, 2014-2015.
// Version: 1.0.8
// BuildNo: 203
//
// Special Thanks to 2C-Why, LLC for supporting this project.
//

package com

import (
	"fmt"
	"os"

	"github.com/pschlump/dbgo"
)

type ErrorBufferType struct {
	Err []string
}

var ErrorBuffer ErrorBufferType

func StashError(s string) {
	if dbgo.IsDbOn("OutputErrors") {
		fmt.Fprintf(os.Stderr, "%s\n", s)
	}
	ErrorBuffer.Err = append(ErrorBuffer.Err, s)
}

func GetErrorStash() (s string) {
	for _, e := range ErrorBuffer.Err {
		s = s + fmt.Sprintf("%s\n", e)
	}
	return s
}

/* vim: set noai ts=4 sw=4: */
