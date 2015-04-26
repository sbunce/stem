// Copyright 2015 Seth Bunce. All rights reserved. Use of this source code is
// governed by a BSD-style license that can be found in the LICENSE file.

package stem

import (
	"fmt"
	"strings"
)

// Token types.
type ttype int
const (
	ttArray ttype = iota
	ttChangeDelim // Never returned by the lexer.
	ttComment     // Never returned by the lexer.
	ttEnd
	ttIfdef
	ttIfndef
	ttInclude
	ttObject
	ttPrint
	ttString
)

// Default delimiters.
const (
	ldel = "{{"
	rdel = "}}"
)

// token returned by the lexer.
type token struct {
	tt   ttype  // tt is the token type.
	val  string // val is the value of the token (name or string literal).
	line int    // line number token starts on.
}

// lexer that operates on the raw template.
type lexer struct {
	name string // name of the template.
	ldel string // ldel is the current left delimiter.
	rdel string // rdel is the current right delimiter.
	line int    // line is the current line.
	src  string // src is the remaining input.
}

// Lookup map for token types.
var tagType = map[byte]ttype{
	'#': ttArray,
	'=': ttChangeDelim,
	'!': ttComment,
	'/': ttEnd,
	'+': ttIfdef,
	'-': ttIfndef,
	'>': ttInclude,
	'$': ttObject,
	'*': ttPrint,
}

// Error returns an error that includes information about where the error
// occurred.
func (l *lexer) Error(a ...interface{}) error {
	return fmt.Errorf("%v:%v %v", l.name, l.line, fmt.Sprint(a...))
}

// Next returns the next token or io.EOF.
func (l *lexer) Next() (*token, error) {
	for {
		t, err := l.next()
		if err != nil {
			return nil, err
		}
		if t == nil {
			return nil, nil
		}
		switch t.tt {
		case ttChangeDelim:
			tok := strings.Split(t.val, " ")
			if len(tok) != 2 {
				l.src = ""
				return nil, l.Error("malformed tag")
			}
			if tok[0] == "" || tok[1] == "" {
				l.src = ""
				return nil, l.Error("malformed tag")
			}
			l.ldel = tok[0]
			l.rdel = tok[1]
			continue
		case ttComment:
			continue
		}
		return t, nil
	}
}

func (l *lexer) next() (*token, error) {
	if l.src == "" {
		return nil, nil
	}
	if strings.HasPrefix(l.src, l.ldel) {
		return l.lexTag()
	}
	return l.lexString()
}

// lexString is called when the src starts with a string.
func (l *lexer) lexString() (*token, error) {
	i := strings.Index(l.src, l.ldel)
	if i == -1 {
		// Remainder of source is string.
		t := &token{tt: ttString, val: l.src, line: l.line}
		l.src = ""
		return t, nil
	}
	t := &token{tt: ttString, val: l.src[:i], line: l.line}
	l.src = l.src[i:]
	l.line += strings.Count(t.val, "\n")
	return t, nil
}

// lexTag is called when the src starts with a tag.
func (l *lexer) lexTag() (*token, error) {
	l.src = l.src[len(l.ldel):]
	if len(l.src) < 1 {
		return nil, l.Error("incomplete tag")
	}
	tt, ok := tagType[l.src[0]]
	if !ok {
		return nil, l.Error("unrecognized tag")
	}
	l.src = l.src[1:]
	i := strings.Index(l.src, l.rdel)
	if i == -1 {
		return nil, l.Error("incomplete tag")
	}
	t := &token{tt: tt, val: l.src[:i], line: l.line}
	l.src = l.src[i+len(l.rdel):]
	// A tag name may have a newline in it.
	l.line += strings.Count(t.val, "\n")
	return t, nil
}

// newLexer returns a new lexer.
func newLexer(name, src string) *lexer {
	return &lexer{
		name: name,
		ldel: ldel,
		rdel: rdel,
		src:  src,
		line: 1,
	}
}

// String pretty prints the token.
func (l *token) String() string {
	if len(l.val) > 10 {
		return fmt.Sprintf("type:%v val:%.10q... line:%v", l.tt, l.val, l.line)
	}
	return fmt.Sprintf("type:%v val:%q line:%v", l.tt, l.val, l.line)
}

// String returns the token type.
func (l ttype) String() string {
	switch l {
	case ttArray:
		return "array"
	case ttComment:
		return "comment"
	case ttEnd:
		return "end"
	case ttIfdef:
		return "ifdef"
	case ttIfndef:
		return "ifndef"
	case ttInclude:
		return "include"
	case ttObject:
		return "object"
	case ttPrint:
		return "print"
	case ttString:
		return "string"
	}
	return "unknown"
}
