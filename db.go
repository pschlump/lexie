package lexie

import (
	"errors"
	"fmt"
	"os"
	"sync"

	"github.com/mgutz/ansi"
)

// ------------------------------------------------------------------------------------------------------------------------------------------
// Debug Print - controllable with flags.
// ------------------------------------------------------------------------------------------------------------------------------------------
var dbOn map[string]bool
var dbOnLock sync.Mutex

func init() {
	dbOn = make(map[string]bool)
	dbOn["debug"] = true
}

func DbPrintf(db string, format string, args ...interface{}) {
	dbOnLock.Lock()
	defer dbOnLock.Unlock()
	if x, o := dbOn[db]; o && x {
		fmt.Printf(format, args...)
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
	dbOnLock.Lock()
	defer dbOnLock.Unlock()
	if x, o := dbOn[db]; o {
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

/* vim: set noai ts=4 sw=4: */
