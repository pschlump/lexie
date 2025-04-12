module github.com/pschlump/lexie

go 1.24.0

require (
	github.com/jessevdk/go-flags v1.5.0
	github.com/microcosm-cc/bluemonday v1.0.27
	github.com/pschlump/ansi v1.0.9
	github.com/pschlump/dbgo v1.0.9
	github.com/pschlump/filelib v1.0.12
	github.com/pschlump/uuid v1.0.3
	github.com/russross/blackfriday v1.6.0
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c
)

require (
	github.com/kr/pretty v0.2.1 // indirect
	github.com/kr/text v0.1.0 // indirect
	github.com/pschlump/go-colorable v0.0.24 // indirect
	github.com/pschlump/go-isatty v0.0.24 // indirect
	github.com/pschlump/json v1.12.1 // indirect
	golang.org/x/net v0.39.0 // indirect
	golang.org/x/sys v0.32.0 // indirect
)

replace github.com/pschlump/dbgo => ../dbgo

replace github.com/microcosm-cc/bluemonday => ../../microcosm-cc/bluemonday

replace github.com/russross/blackfriday => ../../russross/blackfriday
