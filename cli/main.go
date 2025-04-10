package main

//
// C L I / T E S T 2 - Part of Lexie Lexical Generation System
//
// (C) Philip Schlump, 2014-2015.
// Version: 0.2.1
// BuildNo: 34
//
// Special Thanks to 2C-Why, LLC for supporting this project.
//

/*

Answer Question:
	1. What is left for this to be a static site generator? (5)
	2. What is left for this to be a static blog generator? (a bunch of builtins)
	3. What is left for this to be a static X dev system? (15,16,17 + 18:push to server + themes)

TODO :
	5. Try/Fix readJson op -> data															2h
		1. Read
		2. Iterate over some stuff
		3. Add in "if" functions so can check if "features" section is defined
		4. Use in www_lexie_com

	+15. watch-file-and-CLI/make (get "go" cli for this) -- Add watch DIR and -R flags

	16. watch-file and push to browser (run that on ./www)
		1. Get iPad working with (16)

	19. Backups - bk.5 - re-run

	18. Integrate with AWS or Linode - create static sites & update
		?? digital Ocean ??
		?? Linode + Docker Container ??

Quickie to get started with:
	3. "../../../go-lib/uuid" -> ../uuid
	4. Remove all mention of "pongo2" from RingGo
	5. Make all "Ringo" -> "RingGo"
	6. Shrink images to reasonable size (background jpg's)

Tools:

Projects:

	6. RingGo features																		1h

	9. Start documenting

Blog:
	1. Features
	2.

	11. SortByDate operations on list of files
		1. List of files
		2. Transform .md -> .html in list
		3. Sort in date order
		4. Generate URL from date/seq/title
	12. Implement a blog with this
		1. Most Recent X items from list
		2. Chop Html to X chars
		3. Tags
		4. Titles
	14. Checkout copy of this with themes
		1. A "mod-theme" - copy existing theme and rename
		2. A multi-user design for this


Notes:
	This tool uses code from:
		https://github.com/microcosm-cc/bluemonday
		https://github.com/russross/blackfriday

*/

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jessevdk/go-flags"
	"github.com/pschlump/dbgo"
	"github.com/pschlump/filelib"
	"github.com/pschlump/lexie/bluemonday"
	"github.com/pschlump/lexie/com"
	"github.com/pschlump/lexie/dfa"
	"github.com/pschlump/lexie/pbread"
	"github.com/pschlump/lexie/test01"
	"github.com/russross/blackfriday"
)

var opts struct {
	BaseAssets   string `short:"b" long:"baseAssets"    description:"Base Assets (static)"             default:"./assets"`                     //
	ForcedCpFlag bool   `short:"f" long:"forcedCopy"    description:"Copy files even if not changed"   default:"false"`                        //
	Input        string `short:"i" long:"input"         description:"Input File"                       default:""`                             //
	LexPat       string `short:"l" long:"lex"           description:"Lex Input File"                   default:"../in/django3.lex"`            //
	Output       string `short:"o" long:"output"        description:"Output File"                      default:""`                             //
	TemplatePath string `short:"p" long:"templatepath"  description:"Template Path"                    default:"./tmpl"`                       //
	Recursive    bool   `short:"R" long:"recursive"     description:"Recursive Processing"             default:"false"`                        //
	SiteName     string `short:"s" long:"siteName"      description:"Site Name (www_name_com)"         default:""`                             //
	SiteAssets   string `short:"S" long:"siteAssets"    description:"PerSite Assets (static)"          default:"./site_assets/%{site_name%}/"` //
	Theme        string `short:"T" long:"Theme"         description:"Name of a theme (A-Theme)"        default:""`                             //
	User         string `short:"U" long:"User"          description:"User ID (123)"                    default:""`                             //
	Debug        string `short:"X" long:"debug"         description:"Debug Flags"                      default:""`                             //
}

type ConfigOptionsType struct {
	MdExtensions    []string
	ConvertMdToHtml bool
	LeaveTmpFiles   bool
	TmpDir          string
}

var Options ConfigOptionsType

const db_debug3 = true

func main() {

	// ------------------------------------------------------ cli processing --------------------------------------------------------------
	ifnList, err := flags.ParseArgs(&opts, os.Args)

	if err != nil {
		fmt.Printf("Invalid Command Line: %s\n", err)
		os.Exit(1)
	}

	test01.Dbf = os.Stdout
	if opts.Debug != "" {
		s := strings.Split(opts.Debug, ";")
		com.DbOnFlags[opts.Debug] = true
		for _, v := range s {
			if len(v) > 5 && v[0:4] == "out:" {
				test01.Dbf, _ = filelib.Fopen(v[4:], "w")
			} else {
				com.DbOnFlags[v] = true
			}
		}
	}

	fmt.Fprintf(test01.Dbf, "Test Matcher test from %s file, %s\n", opts.LexPat, dbgo.LF())

	// ------------------------------------------------------ Options --------------------------------------------------------------
	// should be read in from a .json file!
	Options.MdExtensions = []string{".md", ".makrdown"}
	Options.ConvertMdToHtml = true
	Options.LeaveTmpFiles = false
	Options.TmpDir = "./tmp"

	if !com.Exists(Options.TmpDir) {
		os.Mkdir(Options.TmpDir, 0700)
	}

	// ------------------------------------------------------ setup Lexie --------------------------------------------------------------
	pt := test01.NewParse2Type()
	pt.Lex = dfa.NewLexie()
	pt.Lex.SetChanelOnOff(true) // Set for getting back stuff via Chanel

	// ------------------------------------------------------ input machine  --------------------------------------------------------------
	if opts.LexPat != "" {
		if !com.Exists(opts.LexPat) {
			fmt.Fprintf(os.Stderr, "Fatal: Must have -l <fn> lexical analyzer machine.  Missing file.\n")
			os.Exit(1)
		}
		pt.Lex.NewReadFile(opts.LexPat) // pstk.Lex.NewReadFile("../in/django3.lex")
	} else {
		fmt.Fprintf(os.Stderr, "Fatal: Must have -l <fn> lexical analyzer machine.\n")
		os.Exit(1)
	}

	// -------------------------------------------------- start scanning process  ----------------------------------------------------------
	tp := strings.Split(opts.TemplatePath, ";")
	for _, tps := range tp {
		pt.OpenLibraries(tps)
	}

	if opts.Recursive {

		CopyInAssets(opts.SiteName, opts.BaseAssets, opts.SiteAssets, opts.Output, opts.User, opts.Theme, opts.ForcedCpFlag)

		//fmt.Printf("After CopyInAssets: Not Implemented Yet\n")
		//os.Exit(1)

		data2 := make(map[string]string)
		data2["site_name"] = opts.SiteName

		if opts.Input == "" && opts.SiteName != "" {
			opts.Input = com.Qt("./site/%{site_name%}/", data2)
		}
		if opts.Output == "" && opts.SiteName != "" {
			opts.Output = com.Qt("./www/%{site_name%}/", data2)
		}

		// ---------------------------------------------------------------------------------------------------------------------------------

		// 1. Do the rsync copy ops
		// 2. Process the set of fiels from -i -> -o

		// -- Process the static files -----------------------------------------------------------------------------------------------------
		dp := make([]string, 0, 10)
		if opts.Input != "" {
			dp = append(dp, opts.Input)
		} else {
			for _, fn := range ifnList[1:] {
				dp = append(dp, fn)
			}
		}

		var fns, dirs []string

		for _, dir := range dp {
			if db_debug3 {
				fmt.Printf("Getting for %s\n", dir)
			}
			t_fns, t_dirs, err := com.GetFilenamesRecrusive(dir)
			if err != nil {
				if db_debug3 {
					fmt.Printf("Error: %s on %s\n", err, dir)
				}
			} else {
				fns = append(fns, t_fns...)
				dirs = append(dirs, t_dirs...)
			}
		}

		if db_debug3 {
			fmt.Printf("fns: %+v\n", fns)
			fmt.Printf("dirs: %+v\n", dirs)
		}

		mds := com.ReplaceEach(dirs, opts.Input, opts.Output)
		for _, aDir := range mds {
			if !com.Exists(aDir) {
				err := os.Mkdir(aDir, 0764)
				if err != nil {
					if db_debug3 {
						fmt.Printf("Error: Unable to create directory %s, error: %s\n", aDir, err)
					}
				}
			}
		}

		mf := com.ReplaceEach(fns, opts.Input, opts.Input+"/%{user%}/%{theme%}/")
		mO := com.ReplaceEach(fns, opts.Input, opts.Output)
		if db_debug3 {
			fmt.Printf("modded_files: %+v\n", mf)
		}

		final := make([]string, 0, len(mf))
		data := make(map[string]string)
		has_err := false

		for _, mff := range mf {
			data["user"] = opts.User
			data["theme"] = opts.Theme
			mfmod := com.Qt(mff, data)
			if com.Exists(mfmod) {
				final = append(final, mfmod)
			} else {

				data["user"] = ""
				// data["theme"] = "A-Theme"
				mfmod := com.Qt(mff, data)
				if com.Exists(mfmod) {
					final = append(final, mfmod)
				} else {

					data["user"] = opts.User
					// data["user"] = ""
					data["theme"] = ""
					mfmod := com.Qt(mff, data)
					if com.Exists(mfmod) {
						final = append(final, mfmod)
					} else {
						fmt.Printf("Error: File Missing %s\n", mfmod)
						has_err = true
					}
				}
			}
		}
		if has_err {
			fmt.Printf("Error occured...\n")
			os.Exit(1)
		}

		if db_debug3 {
			fmt.Printf("Final Files:%s\n", final)
		}
		tmpFiles := make([]string, 0, len(final))
		for ii, yy := range final {
			yyt := yy
			fmt.Printf("Process %s to %s\n", yy, mO[ii])
			ext := filepath.Ext(yy)
			if Options.ConvertMdToHtml && com.InArray(ext, Options.MdExtensions) { // if ext == ".md" || ext == ".markdown" {
				fmt.Printf("\t Convetting from MD to HTML\n")
				in := yy
				//yyt = "./tmp/" + com.Basename(yy) + ".html"	// old code - not using a Tempfile
				//err := ConvertMdToHtmlFile(in, yyt)
				yyt, err = ConvertMdToHtmlFileTmpFile(in)
				if err != nil {
					fmt.Printf("Error: In processing from markdown %s to HTML %s: %s\n", in, yyt, err)
					os.Exit(1)
				}
				mO[ii] = com.RmExt(mO[ii]) + ".html"
				tmpFiles = append(tmpFiles, yyt)
			}
			ProcessFileList(pt, []string{yyt}, mO[ii])
		}
		if !Options.LeaveTmpFiles {
			for _, fn := range tmpFiles {
				os.Remove(fn)
			}
		}

	} else {

		inList := make([]string, 0, 10)
		if opts.Input != "" {
			inList = append(inList, opts.Input)
		} else {
			inList = ifnList[1:]
		}

		ProcessFileList(pt, inList, opts.Output)

	}

}

func debugPrintStrinSlice(s []string, pre string) {
	if db_debug4 {
		for _, vv := range s {
			fmt.Printf("%s%s\n", pre, vv)
		}
	}
}

func MkDirArray(dirs []string) (err error) {
	x_err := err
	for _, dir := range dirs {
		err = os.Mkdir(dir, 0764)
		if err != nil {
			x_err = x_err
		}
	}
	err = x_err
	return
}

func InputModified(in string, out string, forcedCpFlag bool) bool {
	if forcedCpFlag {
		return true
	}
	t_in := ModTimeOfFile(in).UnixNano()
	t_out := ModTimeOfFile(out).UnixNano()
	if db_debug4 {
		fmt.Printf("TimeCmp: %s %s %d %d\n", in, out, t_in, t_out)
	}
	if t_in > t_out {
		return true
	}
	return false
}

func FileSizeDiffers(in string, out string) bool {
	s_in := SizeOfFile(in)
	s_out := SizeOfFile(out)
	if db_debug4 {
		fmt.Printf("SizeCmp: %s %s %d %d\n", in, out, s_in, s_out)
	}
	if s_in != s_out {
		return true
	}
	return false
}

func SizeOfFile(fn string) (sz int64) {

	sz = -1
	fi, err := os.Stat(fn)
	if err != nil {
		// Could not obtain stat, handle error
		return
	}

	sz = fi.Size()
	return
}

// statinfo.ModTime()
func ModTimeOfFile(fn string) (mt time.Time) {

	fi, err := os.Stat(fn)
	if err != nil {
		// Could not obtain stat, handle error
		return
	}

	mt = fi.ModTime()
	return
}

func CopyFilesInHash(fnList map[string]string) {
	// func CopyFile(src, dst string, useHardLink bool) (err error) {
	for ii, vv := range fnList {
		_ = com.CopyFile(vv, ii, false)
	}
}

func ConvertMdToHtml(input []byte) (html []byte) {

	unsafe := blackfriday.MarkdownCommon(input)
	html = bluemonday.UGCPolicy().SanitizeBytes(unsafe)
	return

}

// Not used - old test code
func ConvertMdToHtmlFile(infn, outfn string) (err error) {

	input, err := ioutil.ReadFile(infn)
	if err != nil {
		return
	}

	output := ConvertMdToHtml(input)
	err = ioutil.WriteFile(outfn, output, 0600)

	return
}

func ConvertMdToHtmlFileTmpFile(infn string) (outfn string, err error) {

	// xyzzy - Should Use func ioutil.TempFile(dir, prefix string) (f *os.File, err error)

	input, err := ioutil.ReadFile(infn)
	if err != nil {
		return
	}

	output := ConvertMdToHtml(input)
	// f, err := ioutil.TempFile("./tmp", com.Basename(infn)+"___")
	f, err := ioutil.TempFile(Options.TmpDir, com.Basename(infn)+"___")
	if err != nil {
		return
	}
	n, err := f.Write(output)
	if err == nil && n < len(output) {
		err = io.ErrShortWrite
	}
	if err != nil {
		return
	}
	outfn = f.Name()
	if err1 := f.Close(); err == nil {
		err = err1
	}
	return
}

//

func ProcessFileList(pt *test01.Parse2Type, inList []string, outFn string) (err error) {

	var fp *os.File

	if outFn != "" {
		fp, err = filelib.Fopen(outFn, "w")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Fatal: Unable to create output file\n")
			err = fmt.Errorf("Fatal: Unable to create output file")
			return
		}
		defer fp.Close()
	} else {
		fp = os.Stdout
	}

	for _, fn := range inList {
		if !com.Exists(fn) {
			fmt.Fprintf(os.Stderr, "Fatal: Missing input file %s\n", fn)
			err = fmt.Errorf("Fatal: Missing input file %s", fn)
			return
		}
	}

	go func() {
		r := pbread.NewPbRead()
		for _, fn := range inList {
			r.OpenFile(fn)
		}
		pt.Lex.MatcherLexieTable(r, "S_Init")
	}()

	// ------------------------------------------------------ process tokens --------------------------------------------------------------
	// Generate a parse tree and print out.
	xpt := pt.GenParseTree(0)
	pt.TheTree = xpt
	xpt.DumpMtType(test01.Dbf, 0, 0)
	pt.ExecuteFunctions(0)
	if false {
		fmt.Fprintf(test01.Dbf, "----------------------------------- debug output ----------------------------------------------------\n")
		fmt.Fprintf(test01.Dbf, "%s\n", dbgo.SVarI(xpt))
	}
	fmt.Fprintf(test01.Dbf, "----------------------------------- errors ----------------------------------------------------\n")
	pp := pt.CollectErrorNodes(0)
	for ii, vv := range pp {
		fmt.Fprintf(test01.Dbf, "Error [%3d]: msg=%s\n", ii, vv.ErrorMsg)
	}
	fmt.Fprintf(test01.Dbf, "----------------------------------- final template results  ----------------------------------------------------\n")
	pt.OutputTree(test01.Dbf, 0)

	pt.OutputTree(fp, 0)

	return
}

const db_debug4 = false

//	SiteName     string `short:"s" long:"siteName"      description:"Site Name (www_name_com)"   default:""`                             // new

//	BaseAssets   string `short:"b" long:"baseAssets"    description:"Base Assets (static)"       default:"./assets"`                     // new

//	SiteAssets   string `short:"S" long:"siteAssets"    description:"PerSite Assets (static)"    default:"./site_assets/%{site_name%}/"` // new
//			./site_assets/%{user%}/{%theme%}/%{site_name%}
//			./site_assets/{%theme%}/%{site_name%}
//			./site_assets/%{site_name%}

func CopyInAssets(optsSiteName string, BaseAssets string, SiteAssets string, optsOutput string, User string, Theme string, forcedCpFlag bool) {

	// ---------------------------------------------------------------------------------------------------------------------------------
	data := make(map[string]string)
	data["site_name"] = optsSiteName
	data["user"] = User
	data["theme"] = Theme

	//if opts.Input == "" && optsSiteName != "" {
	//	opts.Input = com.Qt("./site/%{site_name%}/", data2)
	//}

	if optsOutput == "" && optsSiteName != "" {
		optsOutput = com.Qt("./www/%{site_name%}/", data)
	}

	// Generate list of top directories to search
	top := make([]string, 0, 10)
	topDirs := com.Qt(SiteAssets+"/%{user%}/%{theme%}/", data) // ./site_assets/%{site_name%} ->  ./site_assets/%{site_name%}/%{user%}/%{theme%}/
	if com.Exists(topDirs) {
		top = append(top, topDirs)
	}

	data["user"] = ""
	// data["theme"] = "A-Theme"
	topDirs = com.Qt(SiteAssets+"/%{user%}/%{theme%}/", data) // ./site_assets/%{site_name%} ->  ./site_assets/%{site_name%}/%{user%}/%{theme%}/
	if com.Exists(topDirs) {
		top = append(top, topDirs)
	}

	data["user"] = User
	// data["user"] = ""
	data["theme"] = ""
	topDirs = com.Qt(SiteAssets+"/%{user%}/%{theme%}/", data) // ./site_assets/%{site_name%} ->  ./site_assets/%{site_name%}/%{user%}/%{theme%}/
	if com.Exists(topDirs) {
		top = append(top, topDirs)
	}

	top = append(top, BaseAssets)
	// top has array in order of top level dirs to copy from.

	var infn, outfn, dirs []string

	cpList2 := make(map[string]string)

	for ii := range top {
		jj := (len(top) - 1) - ii
		dir := top[jj]

		t_fns, t_dirs, err := com.GetFilenamesRecrusive(dir)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
		} else {

			// if base-fiel-name is not in fns, then add it
			infn = append(infn, t_fns...)
			t_fns = com.ReplaceEach(t_fns, opts.Input, optsOutput)
			outfn = append(outfn, t_fns...)
			outfn = com.ReplaceEach(outfn, dir, optsOutput)

			for kk, in := range infn {
				// fmt.Printf("Loop %2d: in=>%s<- dir=%s\n", kk, in, dir)
				if com.Exists(in) {
					// xyzzy - compare time stamps and size
					if InputModified(in, outfn[kk], forcedCpFlag) || FileSizeDiffers(in, outfn[kk]) {
						cpList2[outfn[kk]] = in
					}
				}
			}

			// if base-dir-name is not in dirs then add base dir name
			t_dirs = com.ReplaceEach(t_dirs, dir, optsOutput)
			for _, x := range t_dirs {
				if !com.InArray(x, dirs) {
					dirs = append(dirs, x)
				}
			}

			if db_debug4 {
				fmt.Printf("\tinfn=\n")
				debugPrintStrinSlice(infn, "\t\t")
				fmt.Printf("\toutfn=\n")
				debugPrintStrinSlice(outfn, "\t\t")
				fmt.Printf("\tdirs 0=\n")
				debugPrintStrinSlice(dirs, "\t\t")
			}
		}
	}
	MkDirArray(dirs)
	if db_debug4 {
		fmt.Printf("cp: %s\n", dbgo.SVarI(cpList2))
	}
	CopyFilesInHash(cpList2)

}
