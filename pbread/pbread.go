package pbread

//
// P B B U F F E R - Push back buffer.
//
// Copyright (C) Philip Schlump, 2013-2025.
// Version: 1.0.0
//
// Push Back Buffer.
//
// This buffer allows for reading input and "pushing back" inputh that you want to look at again.  It is
// primarily designed for the processing of "macros" or templates.
//

import (
	"fmt"
	"io"
	"io/ioutil"
	"path/filepath"
	"unicode/utf8"

	"github.com/pschlump/dbgo"
)

// Data from file or push back save in a buffer
type ABuffer struct {
	Buffer      []rune
	FileName    string
	AbsFileName string
	LineNo      int
	ColNo       int
	Pos         int  // # of chars since start of this file
	EofOnFile   bool // End of buffer is EOF on this file.
}

// The read type to track file position and collect push backs.
type PBReadType struct {
	FileName    string          //
	AbsFileName string          //
	FilesOpened map[string]bool // Set of files that have been opened
	PbBuffer    []*ABuffer      //
	PbAFew      []rune          //
	PbTop       int             //
}

const (
	MaxAFew = 512
)

// Create a new push back buffer and return it.
func NewPbRead() (rv *PBReadType) {
	rv = &PBReadType{
		FilesOpened: make(map[string]bool),
		PbBuffer:    make([]*ABuffer, 0, 10),
		PbAFew:      make([]rune, MaxAFew, MaxAFew),
		PbTop:       0,
	}
	return
}

// Output debugging info
func (pb *PBReadType) Dump01(fo io.Writer) {
	fmt.Fprintf(fo, "Dump At: %s\n", dbgo.LF())
	fmt.Fprintf(fo, "N PbBuffer=%d\n", len(pb.PbBuffer))
	for ii := 0; ii < len(pb.PbBuffer); ii++ {
		fmt.Fprintf(fo, "  Buffer [%d] Len: %d Pos: %d\n", ii, len(pb.PbBuffer[ii].Buffer), pb.PbBuffer[ii].Pos)
		fmt.Fprintf(fo, "  Contents ->")
		for jj := pb.PbBuffer[ii].Pos; jj < len(pb.PbBuffer[ii].Buffer); jj++ {
			fmt.Fprintf(fo, "%s", string(pb.PbBuffer[ii].Buffer[jj]))
		}
		fmt.Fprintf(fo, "<-\n")
	}
	if pb.PbTop > 0 {
		fmt.Fprintf(fo, "PbTop=%d\n", pb.PbTop)
		fmt.Fprintf(fo, "  PbAFew ->")
		for jj := pb.PbTop - 1; jj >= 0; jj-- {
			fmt.Fprintf(fo, "%s", string(pb.PbAFew[jj]))
		}
		fmt.Fprintf(fo, "<-\n")
	}
}

// Open a file - this puts the file at the end of the input.   This is used on the command line for a list of files
// in order.  Each opened and added to the end of the list.
func (pb *PBReadType) OpenFile(fn string) (err error) {
	dbgo.DbPrintf("pbbuf01", "At: %s\n", dbgo.LF())
	pb.FileName = fn
	pb.AbsFileName, _ = filepath.Abs(fn)
	pb.FilesOpened[pb.AbsFileName] = true

	// read file -> PbBuffer
	b := &ABuffer{
		FileName:    fn,
		AbsFileName: pb.AbsFileName,
		LineNo:      1,
		ColNo:       1,
	}
	pb.PbBuffer = append(pb.PbBuffer, b)

	bb, err := ioutil.ReadFile(fn)
	if err != nil {
		return
	}
	b.EofOnFile = true
	b.Pos = 0
	var rn rune
	var sz int
	b.Buffer = make([]rune, 0, len(bb))
	for ii := 0; ii < len(bb); ii += sz {
		rn, sz = utf8.DecodeRune(bb[ii:])
		b.Buffer = append(b.Buffer, rn)
	}

	return nil
}

// Return the next rune.  If runes have been pushed back then use those first.
func (pb *PBReadType) NextRune() (rn rune, done bool) {
	dbgo.DbPrintf("pbbuf01", "At: %s\n", dbgo.LF())
	done = false

	if pb.PbTop > 0 {
		pb.PbTop--
		rn = pb.PbAFew[pb.PbTop]
	} else if len(pb.PbBuffer) <= 0 {
		done = true
		// } else if len(pb.PbBuffer) == 1 && pb.PbBuffer[0].Pos >= len(pb.PbBuffer[0].Buffer) && !pb.PbBuffer[0].EofOnFile {
		// Xyzzy - read in more form file - append
		// so far case never happens because EofOnFile is constant true at init time.
	} else if len(pb.PbBuffer) == 1 && pb.PbBuffer[0].Pos >= len(pb.PbBuffer[0].Buffer) && pb.PbBuffer[0].EofOnFile {
		done = true
	} else if len(pb.PbBuffer) > 1 && pb.PbBuffer[0].Pos >= len(pb.PbBuffer[0].Buffer) && pb.PbBuffer[0].EofOnFile {
		pb.PbBuffer = pb.PbBuffer[1:]
		return pb.NextRune()
	} else {
		//fmt.Printf("Just before core, Pos=%d\n", pb.PbBuffer[0].Pos)
		//fmt.Printf("Just before core, Len 1=%d\n", len(pb.PbBuffer[0].Buffer))
		if pb.PbBuffer[0].Pos >= len(pb.PbBuffer[0].Buffer) { // xyzzy --------------------- pjs - new code - not certain if correct ---------------------------------
			done = true
		} else {
			rn = pb.PbBuffer[0].Buffer[pb.PbBuffer[0].Pos]
			pb.PbBuffer[0].Pos++
			if rn == '\n' {
				pb.PbBuffer[0].LineNo++
				pb.PbBuffer[0].ColNo = 1
			} else {
				pb.PbBuffer[0].ColNo++
			}
		}
	}
	return
}

// Take a peek at what is out there - return the next rune without advancing forward.
func (pb *PBReadType) PeekRune() (rn rune, done bool) {
	dbgo.DbPrintf("pbbuf01", "At: %s\n", dbgo.LF())
	rn, done = pb.NextRune()
	pb.PbRune(rn)
	return
}

// PeekPeekRune will take a peek, 2 in the future, at what is out there - return the next rune without advancing forward.
func (pb *PBReadType) PeekPeekRune() (rn rune, done bool) {
	dbgo.DbPrintf("pbbuf01", "At: %s\n", dbgo.LF())
	rn0, done0 := pb.NextRune()
	rn, done = pb.NextRune()
	done = done0 || done
	pb.PbRune(rn)
	pb.PbRune(rn0)
	return
}

// Take any pushed back stuff and put it into a buffer.
func (pb *PBReadType) pushbackIntoBuffer() {
	bl := pb.PbTop
	if bl == 0 {
		return
	}
	b := &ABuffer{
		Buffer:    make([]rune, bl, bl), // /*old*/ Buffer:    make([]rune, MaxAFew, MaxAFew),
		EofOnFile: true,
	}
	if len(pb.PbBuffer) > 0 {
		b.FileName = pb.PbBuffer[0].FileName
		b.AbsFileName = pb.PbBuffer[0].AbsFileName
		b.LineNo = pb.PbBuffer[0].LineNo
		b.ColNo = pb.PbBuffer[0].ColNo
		b.Pos = 0
	} else {
		b.FileName = ""
		b.AbsFileName = ""
		b.LineNo = 1
		b.ColNo = 1
		b.Pos = 0
	}
	for jj, ii := pb.PbTop-1, 0; jj >= 0; jj-- {
		b.Buffer[ii] = pb.PbAFew[jj]
		ii++
	}

	pb.PbTop = 0

	if len(pb.PbBuffer) > 0 {
		pb.PbBuffer = append([]*ABuffer{b}, pb.PbBuffer...) // prepend
	} else {
		pb.PbBuffer = append(pb.PbBuffer, b)
	}
}

// Push back a single rune onto input.  You can call this more than one time.
func (pb *PBReadType) PbRune(rn rune) {
	dbgo.DbPrintf("pbbuf01", "At: %s\n", dbgo.LF())

	if pb.PbTop >= MaxAFew { // Buffer is full
		pb.pushbackIntoBuffer()
	}

	pb.PbAFew[pb.PbTop] = rn
	pb.PbTop++
}

// Push back a slice of runes
func (pb *PBReadType) PbRuneArray(rns []rune) {
	dbgo.DbPrintf("pbbuf01", "At: %s\n", dbgo.LF())
	for ii := len(rns) - 1; ii >= 0; ii-- {
		pb.PbRune(rns[ii])
	}
}

// Push back a string - will be converted form string to array of runes
func (pb *PBReadType) PbString(s string) {
	dbgo.DbPrintf("pbbuf01", "At: %s\n", dbgo.LF())
	rns := make([]rune, 0, len(s))
	var rn rune
	var sz int
	for ii := 0; ii < len(s); ii += sz {
		rn, sz = utf8.DecodeRune([]byte(s[ii:]))
		rns = append(rns, rn)
	}
	pb.PbRuneArray(rns)
}

// Push back a string.  Will be converted from an array of byte to an array of runes.
func (pb *PBReadType) PbByteArray(s []byte) {
	dbgo.DbPrintf("pbbuf01", "At: %s\n", dbgo.LF())
	rns := make([]rune, 0, len(s))
	var rn rune
	var sz int
	for ii := 0; ii < len(s); ii += sz {
		rn, sz = utf8.DecodeRune(s[ii:])
		rns = append(rns, rn)
	}
	pb.PbRuneArray(rns)
}

// Place the contents of a file in buffers at the head so NextRune will pull from this next.
func (pb *PBReadType) PbFile(fn string) (err error) {
	dbgo.DbPrintf("pbbuf01", "At: %s\n", dbgo.LF())
	err = nil

	pb.pushbackIntoBuffer()

	pb.FileName = fn
	pb.AbsFileName, _ = filepath.Abs(fn)
	pb.FilesOpened[pb.AbsFileName] = true

	// read file -> PbBuffer
	b := &ABuffer{
		FileName:    fn,
		AbsFileName: pb.AbsFileName,
		LineNo:      1,
		ColNo:       1,
	}
	// pb.PbBuffer = append(pb.PbBuffer, b)
	// data = append([]string{"Prepend Item"}, data...)
	pb.PbBuffer = append([]*ABuffer{b}, pb.PbBuffer...) // prepend

	bb, err := ioutil.ReadFile(fn)
	if err != nil {
		return
	}
	b.EofOnFile = true
	b.Pos = 0
	var rn rune
	var sz int
	b.Buffer = make([]rune, 0, len(bb))
	for ii := 0; ii < len(bb); ii += sz {
		rn, sz = utf8.DecodeRune(bb[ii:])
		b.Buffer = append(b.Buffer, rn)
	}

	return
}

// Have we already seen the specified file.  Useful for require(fn)
func (pb *PBReadType) FileSeen(fn string) bool {
	dbgo.DbPrintf("pbbuf01", "At: %s\n", dbgo.LF())
	a, _ := filepath.Abs(fn)
	if t, ok := pb.FilesOpened[a]; ok && t {
		return true
	}
	return false
}

// Get the current line/col no and file name
func (pb *PBReadType) GetPos() (LineNo int, ColNo int, FileName string) {
	dbgo.DbPrintf("pbbuf02", "At: %s\n", dbgo.LF())
	if len(pb.PbBuffer) > 0 {
		dbgo.DbPrintf("pbbuf02", "From Buffer At: %s\n", dbgo.LF())
		LineNo = pb.PbBuffer[0].LineNo
		ColNo = pb.PbBuffer[0].ColNo
		FileName = pb.PbBuffer[0].FileName
	} else {
		dbgo.DbPrintf("pbbuf02", "Not set At: %s\n", dbgo.LF())
		LineNo = 1
		ColNo = 1
		FileName = ""
	}
	return
}

// Set the line/col/file-name for the current buffer - Useful for constructing something like C/Pre processor's #line
func (pb *PBReadType) SetPos(LineNo int, ColNo int, FileName string) {
	dbgo.DbPrintf("pbbuf01", "At: %s\n", dbgo.LF())
	pb.pushbackIntoBuffer()
	if len(pb.PbBuffer) > 0 {
		pb.PbBuffer[0].LineNo = LineNo
		pb.PbBuffer[0].ColNo = ColNo
		pb.PbBuffer[0].FileName = FileName
	}
	return
}
