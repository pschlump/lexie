package pbread

//
// P B B U F F E R - Push back buffer.
//
// (C) Philip Schlump, 2014-2015.
// Version: 1.0.8
// BuildNo: 203
//
// Special Thanks to 2C-Why, LLC for supporting this project.
//

/*
// xyzzyErr1 - what happens if we PbFile when NPbAFew is > 0 - don't we need to pre-buffer the push back?
*/

import (
	"fmt"
	"io"
	"io/ioutil"
	"path/filepath"
	"unicode/utf8"

	"../com"

	"../../../go-lib/tr"

	// "../../../go-lib/tr"
	// "../../../go-lib/sizlib"
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
	FileName     string
	AbsFileName  string
	FilesOpened  map[string]bool // Set of files that have been opened
	PbBuffer     []*ABuffer
	PbAFew       []rune
	NPbAFew      int
	NPbAFewStart int
	NPbAFewEnd   int
}

//
const (
	MaxAFew = 512
)

// Create a new push back buffer and return it.
func NewPbRead() (rv *PBReadType) {
	rv = &PBReadType{
		FilesOpened:  make(map[string]bool),
		PbBuffer:     make([]*ABuffer, 0, 10),
		PbAFew:       make([]rune, MaxAFew, MaxAFew),
		NPbAFew:      0,
		NPbAFewStart: 0,
		NPbAFewEnd:   0,
	}
	return
}

// Output debugging info
func (pb *PBReadType) Dump01(fo io.Writer) {
	fmt.Fprintf(fo, "Dump At: %s\n", tr.LF())
	fmt.Fprintf(fo, "N PbBuffer=%d\n", len(pb.PbBuffer))
	for ii := 0; ii < len(pb.PbBuffer); ii++ {
		fmt.Fprintf(fo, "  Buffer [%d] Len: %d Pos: %d\n", ii, len(pb.PbBuffer[ii].Buffer), pb.PbBuffer[ii].Pos)
		fmt.Fprintf(fo, "  Contents ->")
		for jj := pb.PbBuffer[ii].Pos; jj < len(pb.PbBuffer[ii].Buffer); jj++ {
			fmt.Fprintf(fo, "%s", string(pb.PbBuffer[ii].Buffer[jj]))
		}
		fmt.Fprintf(fo, "<-\n")
	}
	if pb.NPbAFew > 0 {
		fmt.Fprintf(fo, "NPbAFew=%d\n", pb.NPbAFew)
		fmt.Fprintf(fo, "NPbAFewStart=%d\n", pb.NPbAFewStart)
		fmt.Fprintf(fo, "NPbAFewEnd=%d\n", pb.NPbAFewEnd)
		fmt.Fprintf(fo, "  PbAFew ->")
		for jj := pb.NPbAFewStart; jj < len(pb.PbAFew) && jj < pb.NPbAFewEnd; jj++ {
			fmt.Fprintf(fo, "%s", string(pb.PbAFew[jj]))
		}
		fmt.Fprintf(fo, "<-\n")
	}
}

// Open a file - this puts the file at the end of the input.   This is used on the command line for a list of files
// in order.  Each opened and added to the end of the list.
func (pb *PBReadType) OpenFile(fn string) (err error) {
	com.DbPrintf("pbbuf01", "At: %s\n", tr.LF())
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
	com.DbPrintf("pbbuf01", "At: %s\n", tr.LF())
	done = false

	if pb.NPbAFew > 0 {
		pb.NPbAFew--
		rn = pb.PbAFew[pb.NPbAFewStart]
		pb.NPbAFewStart++
		if pb.NPbAFewStart >= pb.NPbAFewEnd || pb.NPbAFew <= 0 {
			pb.NPbAFewStart = 0
			pb.NPbAFewEnd = 0
			pb.NPbAFew = 0
		}
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

		rn = pb.PbBuffer[0].Buffer[pb.PbBuffer[0].Pos] // xyzzyStack
		pb.PbBuffer[0].Pos++

		if rn == '\n' {
			pb.PbBuffer[0].LineNo++
			pb.PbBuffer[0].ColNo = 1
		} else {
			pb.PbBuffer[0].ColNo++
		}
	}
	return
}

// Take a peek at what is out there - return the next rune without advancing forward.
func (pb *PBReadType) PeekRune() (rn rune, done bool) {
	com.DbPrintf("pbbuf01", "At: %s\n", tr.LF())
	rn, done = pb.NextRune()
	if pb.NPbAFewStart > 0 && pb.NPbAFewStart <= pb.NPbAFewEnd { // xyzzyStack
		pb.NPbAFewStart--
		pb.NPbAFew++
	} else {
		// I am still concerned that I have missed a case at this point!   Logic seems wrong or incomplete.
		pb.PbRune(rn)
	}
	return
}

// Take any pushed back stuff and put it into a buffer.
func (pb *PBReadType) pushbackIntoBuffer() {
	bl := pb.NPbAFew
	b := &ABuffer{
		Buffer:    make([]rune, bl, bl), // Buffer:    make([]rune, MaxAFew, MaxAFew),
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
	jj := pb.NPbAFewStart
	// for ii := 0; ii < MaxAFew && jj < pb.NPbAFew; ii++ {
	for ii := 0; ii < bl && jj < pb.NPbAFew; ii++ { // xyzzyStack - must reverse if stack
		b.Buffer[ii] = pb.PbAFew[jj]
		jj++
	}

	pb.NPbAFewStart = 0
	pb.NPbAFewEnd = 0
	pb.NPbAFew = 0

	if len(pb.PbBuffer) > 0 {
		pb.PbBuffer = append([]*ABuffer{b}, pb.PbBuffer...) // prepend
	} else {
		pb.PbBuffer = append(pb.PbBuffer, b)
	}
}

// Push back a single rune onto input.  You can call this more than one time.
func (pb *PBReadType) PbRune(rn rune) {
	com.DbPrintf("pbbuf01", "At: %s\n", tr.LF())

	if pb.NPbAFew >= MaxAFew { // Buffer is full
		pb.pushbackIntoBuffer()
	}

	pb.NPbAFew++
	pb.PbAFew[pb.NPbAFewEnd] = rn // xyzzyStack - should insert on top of stack
	pb.NPbAFewEnd++
}

// Push back a slice of runes
func (pb *PBReadType) PbRuneArray(rns []rune) {
	com.DbPrintf("pbbuf01", "At: %s\n", tr.LF())
	//for ii := len(rns) - 1; ii >= 1; ii-- {
	//	pb.PbRune(rns[ii])
	//}
	for _, vv := range rns {
		pb.PbRune(vv)
	}
}

// Push back a string - will be converted form string to array of runes
func (pb *PBReadType) PbString(s string) {
	com.DbPrintf("pbbuf01", "At: %s\n", tr.LF())
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
	com.DbPrintf("pbbuf01", "At: %s\n", tr.LF())
	rns := make([]rune, 0, len(s))
	var rn rune
	var sz int
	for ii := 0; ii < len(s); ii += sz {
		rn, sz = utf8.DecodeRune(s[ii:])
		rns = append(rns, rn)
	}
	pb.PbRuneArray(rns)
}

// xyzzyErr1 - what happens if we PbFile when NPbAFew is > 0 - don't we need to pre-buffer the push back?
// Place the contents of a file in buffers at the head so NextRune will pull from this next.
func (pb *PBReadType) PbFile(fn string) (err error) {
	com.DbPrintf("pbbuf01", "At: %s\n", tr.LF())
	err = nil

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
	com.DbPrintf("pbbuf01", "At: %s\n", tr.LF())
	a, _ := filepath.Abs(fn)
	if t, ok := pb.FilesOpened[a]; ok && t {
		return true
	}
	return false
}

// Get the current line/col no and file name
func (pb *PBReadType) GetPos() (LineNo int, ColNo int, FileName string) {
	com.DbPrintf("pbbuf01", "At: %s\n", tr.LF())
	if len(pb.PbBuffer) > 0 {
		com.DbPrintf("pbbuf01", "From Buffer At: %s\n", tr.LF())
		LineNo = pb.PbBuffer[0].LineNo
		ColNo = pb.PbBuffer[0].ColNo
		FileName = pb.PbBuffer[0].FileName
	} else {
		com.DbPrintf("pbbuf01", "Not set At: %s\n", tr.LF())
		LineNo = 1
		ColNo = 1
		FileName = ""
	}
	return
}

// Set the line/col/file-name for the current buffer - Useful for constructing something like C/Pre processor's #line
func (pb *PBReadType) SetPos(LineNo int, ColNo int, FileName string) {
	com.DbPrintf("pbbuf01", "At: %s\n", tr.LF())
	if len(pb.PbBuffer) > 0 {
		if pb.NPbAFew > 0 {
			pb.pushbackIntoBuffer()
		}
		pb.PbBuffer[0].LineNo = LineNo
		pb.PbBuffer[0].ColNo = ColNo
		pb.PbBuffer[0].FileName = FileName
	}
	return
}
