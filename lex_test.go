// Copyright 2015 Seth Bunce. All rights reserved. Use of this source code is
// governed by a BSD-style license that can be found in the LICENSE file.

package stem

import (
	"reflect"
	"testing"
)

func TestChangeDelim(t *testing.T) {
	src := "{{*a}}{{=[[ ]]}}[[*b]][[=<< >>]]<<*c>>"
	tests := []*token{
		&token{
			tt:   ttPrint,
			val:  "a",
			line: 1,
		},
		&token{
			tt:   ttPrint,
			val:  "b",
			line: 1,
		},
		&token{
			tt:   ttPrint,
			val:  "c",
			line: 1,
		},
	}
	lex := newLexer("", src)
	for _, want := range tests {
		got, err := lex.Next()
		if err != nil {
			t.Fatalf("couldn't get next token: %v", err)
		}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("got token %v, want %v", got, want)
		}
	}
}	

func TestLex(t *testing.T) {
	src := "{{#a}}\n{{/c}}\n{{+d}}\n{{-e}}\n{{$f}}\n{{*g}}\n{{>h}}\n{{>i\n}}"
	tests := []*token{
		&token{
			tt:   ttArray,
			val:  "a",
			line: 1,
		},
		&token{
			tt:   ttString,
			val:  "\n",
			line: 1,
		},
		&token{
			tt:   ttEnd,
			val:  "c",
			line: 2,
		},
		&token{
			tt:   ttString,
			val:  "\n",
			line: 2,
		},
		&token{
			tt:   ttIfdef,
			val:  "d",
			line: 3,
		},
		&token{
			tt:   ttString,
			val:  "\n",
			line: 3,
		},
		&token{
			tt:   ttIfndef,
			val:  "e",
			line: 4,
		},
		&token{
			tt:   ttString,
			val:  "\n",
			line: 4,
		},
		&token{
			tt:   ttObject,
			val:  "f",
			line: 5,
		},
		&token{
			tt:   ttString,
			val:  "\n",
			line: 5,
		},
		&token{
			tt:   ttPrint,
			val:  "g",
			line: 6,
		},
		&token{
			tt:   ttString,
			val:  "\n",
			line: 6,
		},
		&token{
			tt:   ttInclude,
			val:  "h",
			line: 7,
		},
		&token{
			tt:   ttString,
			val:  "\n",
			line: 7,
		},
		&token{
			tt:   ttInclude,
			val:  "i\n",
			line: 8,
		},
	}
	lex := newLexer("", src)
	for _, want := range tests {
		got, err := lex.Next()
		if err != nil {
			t.Fatalf("couldn't get next token: %v", err)
		}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("got token %v, want %v", got, want)
		}
	}
}
