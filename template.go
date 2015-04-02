// Copyright 2015 Seth Bunce. All rights reserved. Use of this source code is
// governed by a BSD-style license that can be found in the LICENSE file.
package stem

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"path"
	"reflect"
)

// Compiled template ready to be combined with data.
// Multiple goroutines can use tmpl concurrently.
type Template struct {
	name string
	tree []node
}

// executeRecurse recursively parses. Every time we encounter a node which
// contains other nodes we recurse and push a new symbol table on to the stack
// for recursive lookup.
func executeRecurse(wr io.Writer, set *Set, sym *symtab, tree []node) error {
	for _, n := range tree {
		switch nt := n.(type) {
		case *nodeArray:
			array := sym.Array(nt.name)
			if array.IsValid() {
				for i := 0; i < array.Len(); i++ {
					elem := indirect(array.Index(i))
					if elem.Kind() == reflect.Map && !elem.IsNil() {
						if err := executeRecurse(wr, set, sym.EnterObject(elem), nt.nodes); err != nil {
							return err
						}
					} else {
						if err := executeRecurse(wr, set, sym.EnterArrayElem(elem), nt.nodes); err != nil {
							return err
						}
					}
				}
			}
		case *nodeIfdef:
			if sym.Ifdef(nt.name) {
				if err := executeRecurse(wr, set, sym, nt.nodes); err != nil {
					return err
				}
			}
		case *nodeIfndef:
			if sym.Ifndef(nt.name) {
				if err := executeRecurse(wr, set, sym, nt.nodes); err != nil {
					return err
				}
			}
		case *nodeInclude:
			if set != nil {
				if t := set.template(nt.name); t != nil {
					if err := executeRecurse(wr, set, sym, t.tree); err != nil {
						return err
					}
				}
			}
		case *nodeObject:
			obj := sym.Object(nt.name)
			if obj.IsValid() {
				if err := executeRecurse(wr, set, sym.EnterObject(obj), nt.nodes); err != nil {
					return err
				}
			}
		case *nodePrint:
			if _, err := wr.Write([]byte(sym.Print(nt.name))); err != nil {
				return err
			}
		case *nodeString:
			if _, err := wr.Write([]byte(nt.val)); err != nil {
				return err
			}
		default:
			panic("unknown node type, programmer error")
		}
	}
	return nil
}

// Parse template.
func Parse(text string) (*Template, error) {
	tree, err := parse("", text)
	if err != nil {
		return nil, err
	}
	return &Template{
		tree: tree,
	}, nil
}

// ParseFile parses a template file. The template name will be file name.
func ParseFile(filename string) (*Template, error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("couldn't open file %q, error: %v", filename, err)
	}
	t, err := Parse(string(b))
	if err != nil {
		return nil, err
	}
	t.SetName(TemplateName(filename))
	return t, nil
}

// MustParse parses a template and panics if there's an error.
func MustParse(text string) *Template {
	t, err := Parse(text)
	if err != nil {
		panic(err)
	}
	return t
}

// Execute combines the template with data and writes the result to wr.
func (tmpl *Template) Execute(wr io.Writer, data map[string]interface{}) error {
	return executeRecurse(wr, nil, newsymtab(data), tmpl.tree)
}

// ExecuteJSON combines the template with JSON data and writes the result to wr.
func (tmpl *Template) ExecuteJSON(wr io.Writer, JSON string) error {
	data := make(map[string]interface{})
	if err := json.Unmarshal([]byte(JSON), &data); err != nil {
		return fmt.Errorf("couldn't unmarshal json: %v", err)
	}
	return executeRecurse(wr, nil, newsymtab(data), tmpl.tree)
}

// Filter all strings in the template.
func (tmpl *Template) Filter(filters Filter) {
	filter(tmpl.tree, filters)
}

// SetName that template can be included by in a Set.
func (tmpl *Template) SetName(name string) {
	tmpl.name = name
}

// TemplateName for template with specified filename.
func TemplateName(filename string) string {
	return path.Base(filename)
}
