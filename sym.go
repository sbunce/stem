// Copyright 2015 Seth Bunce. All rights reserved. Use of this source code is
// governed by a BSD-style license that can be found in the LICENSE file.

package stem

import (
	"fmt"
	"reflect"
)

// Symbol table.
type symtab struct {
	scope     []reflect.Value // scopes has inner most scope as the last elem.
	arrayElem reflect.Value   // arrayElem is zero value except when in array.
}

// indirect all interfaces/pointers.
func indirect(v reflect.Value) reflect.Value {
loop:
	for {
		switch v.Kind() {
		case reflect.Interface:
			v = v.Elem()
		case reflect.Ptr:
			v = v.Elem()
		default:
			break loop
		}
	}
	return v
}

func newsymtab(data map[string]interface{}) *symtab {
	return &symtab{
		scope: []reflect.Value{reflect.ValueOf(data)},
	}
}

// Array returns a slice or the zero value.
func (s *symtab) Array(key string) reflect.Value {
	v := reflect.ValueOf(key)
	for x := len(s.scope) - 1; x >= 0; x-- {
		if e := s.scope[x].MapIndex(v); e.IsValid() {
			e = indirect(e)
			if e.Kind() == reflect.Slice && e.IsValid() && !e.IsNil() {
				return e
			}
			break
		}
	}
	return reflect.Value{}
}

// EnterArrayElem sets the current scope as an array element.
func (s *symtab) EnterArrayElem(elem reflect.Value) *symtab {
	return &symtab{
		scope:     s.scope,
		arrayElem: elem,
	}
}

// EnterObject returns the symbol table with obj as the inner most scope.
func (s *symtab) EnterObject(obj reflect.Value) *symtab {
	return &symtab{
		scope: append(s.scope, obj),
	}
}

// Ifdef returns true if the key is defined.
func (s *symtab) Ifdef(key string) bool {
	if key == "" && s.arrayElem.IsValid() {
		return true
	}
	v := reflect.ValueOf(key)
	for x := len(s.scope) - 1; x >= 0; x-- {
		if e := s.scope[x].MapIndex(v); e.IsValid() {
			return true
		}
	}
	return false
}

// Ifndef returns true if the key is not defined.
func (s *symtab) Ifndef(symbol string) bool {
	return !s.Ifdef(symbol)
}

// Print returns the string representation of the value.
func (s *symtab) Print(symbol string) string {
	if symbol == "" && s.arrayElem.IsValid() {
		return fmt.Sprint(indirect(s.arrayElem).Interface())
	}
	v := reflect.ValueOf(symbol)
	for x := len(s.scope) - 1; x >= 0; x-- {
		if e := s.scope[x].MapIndex(v); e.IsValid() {
			return fmt.Sprint(indirect(e).Interface())
		}
	}
	return ""
}

// Object returns a map or the zero value.
func (s *symtab) Object(symbol string) reflect.Value {
	v := reflect.ValueOf(symbol)
	for x := len(s.scope) - 1; x >= 0; x-- {
		if e := s.scope[x].MapIndex(v); e.IsValid() {
			e = indirect(e)
			if e.Kind() == reflect.Map && e.IsValid() && !e.IsNil() {
				return e
			}
			break
		}
	}
	return reflect.Value{}
}
