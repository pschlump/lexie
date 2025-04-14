//
// D B - Part of Lexie Lexical Generation System
//
// Copyright (C) Philip Schlump, 2014-2025.
// Version: 1.0.8
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
