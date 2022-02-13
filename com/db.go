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
	"errors"
	"fmt"
	"io"
	"os"
	"runtime"
	"sync"

	"github.com/pschlump/ansi"
)

// ------------------------------------------------------------------------------------------------------------------------------------------
// Debug Print - controllable with flags.
// ------------------------------------------------------------------------------------------------------------------------------------------
var DbOnFlags map[string]bool
var DbOnFlagsLock sync.Mutex

func init() {
	DbOnFlags = make(map[string]bool)
	DbOnFlags["debug"] = true
}

func DbPrintf(db string, format string, args ...interface{}) {
	DbOnFlagsLock.Lock()
	defer DbOnFlagsLock.Unlock()
	if x, o := DbOnFlags[db]; o && x {
		fmt.Printf(format, args...)
	}
}

func DbFprintf(db string, w io.Writer, format string, args ...interface{}) {
	DbOnFlagsLock.Lock()
	defer DbOnFlagsLock.Unlock()
	if x, o := DbOnFlags[db]; o && x {
		fmt.Fprintf(w, format, args...)
	}
}

var (
	Red    = ansi.ColorCode("red")
	Yellow = ansi.ColorCode("yellow")
	Green  = ansi.ColorCode("green")
	Reset  = ansi.ColorCode("reset")
)

func DbOn(db string) (ok bool) {
	ok = false
	DbOnFlagsLock.Lock()
	defer DbOnFlagsLock.Unlock()
	if x, o := DbOnFlags[db]; o {
		ok = x
	}
	return
}

func Fopen(fn string, mode string) (file *os.File, err error) {
	file = nil
	if mode == "r" {
		file, err = os.Open(fn) // For read access.
	} else if mode == "w" {
		file, err = os.OpenFile(fn, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	} else if mode == "a" {
		file, err = os.OpenFile(fn, os.O_RDWR|os.O_APPEND, 0660)
	} else {
		err = errors.New("Invalid Mode")
	}
	return
}

type ErrorBufferType struct {
	Err []string
}

var ErrorBuffer ErrorBufferType

func StashError(s string) {
	if DbOn("OutputErrors") {
		fmt.Printf("%s\n", s)
	}
	ErrorBuffer.Err = append(ErrorBuffer.Err, s)
}

// Return the File name and Line no as a string.
func LF(d ...int) (rv string) {
	depth := 1
	rv = "File: Unk LineNo:Unk"
	if len(d) > 0 {
		depth = d[0]
	}
	_, file, line, ok := runtime.Caller(depth)
	if ok {
		rv = fmt.Sprintf("File: %s LineNo:%d", file, line)
	}
	return
}

// ----------------------------------------------------------------------------------------------------------
// Return the current line number as a string.
func LINE(d ...int) (rv string) {
	rv = "LineNo:Unk"
	depth := 1
	if len(d) > 0 {
		depth = d[0]
	}
	_, _, line, ok := runtime.Caller(depth)
	if ok {
		rv = fmt.Sprintf("%d", line)
	}
	return
}

// Return the current file name.
func FILE(d ...int) (rv string) {
	rv = "File:Unk"
	depth := 1
	if len(d) > 0 {
		depth = d[0]
	}
	_, file, _, ok := runtime.Caller(depth)
	if ok {
		rv = file
	}
	return
}

func LINEn(d ...int) (rv int) {
	depth := 1
	if len(d) > 0 {
		depth = d[0]
	}
	_, _, line, ok := runtime.Caller(depth)
	if ok {
		rv = line
	}
	return
}

/* vim: set noai ts=4 sw=4: */
