# Lexie - A Lexical Analizer Genrator and Lexical Analyzer

Thanks for taking a look at *Lexie*.

Lexie is designed to generate fast lexical analyzer based on
transforming a set of regular expressions into a nondeterministic
finite state machine (NFA) and then taking that NFA and transforming
it into a deterministic finite state (DFA) machine.   Multiple DFAs
can be generated and can be switched between and pushed and popped
to allow context sensitive scanning of input.

Lexie has the ability to change the specification and then regenerate
the NFA and DFA at run time.

An Example.   You want to specify the scanner for a template language
that starts template with `{{` and ends them with `}}`.  Inside the
template you want to recognize strings.  You can embed `{{` and `}}`
inside your strings.  The is the context dependent part of lexie.
Also you need to be able to chagne the `{{` and `}}` to some other
tokens.   For example you may want yor tool to work with AngularJS
and it already uses `{{` and `}} as delimeters.   So you wat at
runtime to be able to say, `{{changequote "{{" "-=[" "}}" "]=-"}}`
and change from `{{` to `-=[` and chagne the closing template 
marker from `}}` to `]=-`.

Lexie can scan languages that involve nested items.  For example
you can specify a C-like comment and make it nest and contain
other C-like comments.  This make `/* commented out /* comment */ nests */`
a legitimate input.  This is easy to do and an example of this
is in the ./examles directory.

Lexie also has a concept of reserved words to that a pattern match can
pick out word tokens and then lookup specific values of that word as 
a reserved word.  Example:  The pattern `[a-zA-Z_][a-zA-Z_0-9]*` 
matches all identifieers in a language.   After the match a check
can be made to see if the identifer is one of the reserved words
`or`, `and` and return a different token for these.   

## Development status

Lexie is used inside [Ringo](https://github.com/pschlump/ringo),
the template/macro processor that implements a superset of Django
Templates in go.  Ringo is based on
[pongo2](https://github.com/flosch/pongo2) and I am thankful for
the wonderful work on pongo2 that lead to building this tool.
No code in Lexie is taken from pongo2.

Lexie is also used for a fast user-agent identification library,
ua-quick (not open soruce yet, but I am working on it).  This has
improved the speed of parsing and identifying user agens by a factor
of 10,000.

The set of regular expressions that lexie understands is limited
but growing.  It is adequate to specify simple languages.   This is *not*
a Perl-regular expression matcher - nor is it Posix.  That said...
It is fast and usable.

## Features

1) Works with UTF8 / Unicode.
2) Runs multiple sets of pattern matchers in a context dependent fasion.
3) Changeable on the fly at runtime.
4) Clear error reporting.
5) Embedable - can be used inside another proram.
6) Stand alone - can be run as a fast pattern maching tool in a stand alone configuraiton.
7) Fast.  Fast to perform pattern matches.  Fast to modify an existing matcher.
8) Clear error messages if a modification of a pattern matcher breaks the matcher.
9) Context Senstive Matching with multiple machines and push/pop of machines and states.
10) Ability to push-back onto input and re-scan input if necessary.
11) Extensible matchers that allow for non-text pattern matching (Think greenhouse control and real time systems control).
12) Runs as a goroutine - this improves performance.
13) Has a cute name and a cute mascot.
14) State machines can be cached in Redis so that they do not need to be regenrated every time.
15) Can directly generate state machines in Go code for fast static state macines.

## Lexie Definition

"The name Lexie is an American baby name. In American the meaning
of the name Lexie is: defender of mankind. ...  1. People with this
name tend to initiate events, to be leaders rather than followers,
with powerful personalities."
From: www.sheknows.com/baby-names/name/lexie

## Comparison with other similar tools

This is not ment to be a comprehensive list.  These are the tools that I use.

### lex

Lex is now a quite old tool.  I still refer to its documentation it when I am workin in C or C++.
It can be used with other languages.  The newer replacement flex is a much better
choice.  Lex apears to me to have a number of serious defects (atleast 20+ years
ago it had defects - then I switched to flex).

### flex

The open source version of Lex with lots of fixes.  Lots of languages are
supported.  It generates static tables and supports multiple input states.
Unicode support is bascially missing.  Dynamic re-configuration is not 
really a choice.   It works best in C and C++ with Bison for a parser
generator.

## Notes

1. Be able to output the DFA into a file for re-reading quicey - so as to not
need to re-build it every time.  
2. Output should be in "JSON" or ".go" code.


