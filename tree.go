// Copyright 2015 Seth Bunce. All rights reserved. Use of this source code is
// governed by a BSD-style license that can be found in the LICENSE file.

package stem

import (
	"fmt"
)

// depthLimit is the max recurse depth used to stop pathological cases.
const depthLimit = 32

// node in the parse tree.
type node interface {
	// Prints node type.
	String() string
}

// nodeArray is a repeated section.
type nodeArray struct {
	name  string
	nodes []node
}

// nodeIfdef renders if the name is defined.
type nodeIfdef struct {
	name  string
	nodes []node
}

// nodeIfndef renders if the name is not defined.
type nodeIfndef struct {
	name  string
	nodes []node
}

// nodeInclude includes another template by name.
type nodeInclude struct {
	name string
}

// nodeObject enters a JSON object.
type nodeObject struct {
	name  string
	nodes []node
}

// nodePrint prints a symbol.
type nodePrint struct {
	name string
}

// nodeString is a string literal.
type nodeString struct {
	val string
}

func (n *nodeArray) String() string {
	return "array"
}

func (n *nodeIfdef) String() string {
	return "ifdef"
}

func (n *nodeIfndef) String() string {
	return "ifndef"
}

func (n *nodeInclude) String() string {
	return "include"
}

func (n *nodeObject) String() string {
	return "object"
}

func (n *nodePrint) String() string {
	return "print"
}

func (n *nodeString) String() string {
	return "string"
}

// parse creates a parse tree.
func parse(name, src string) ([]node, error) {
	return parseRecurse(make([]node, 0), newLexer(name, src), nil, 0)
}

// parseRecurse recursively builds a parse tree.
func parseRecurse(tree []node, l *lexer, end *token, depth int) ([]node, error) {
	depth++
	if depth > depthLimit {
		return nil, fmt.Errorf("depth limit %v", depthLimit)
	}
	for {
		t, err := l.Next()
		if err != nil {
			return nil, err
		}
		if t == nil {
			if end != nil {
				return nil, fmt.Errorf("unclosed scope %v", end)
			}
			return tree, nil
		}
		switch t.tt {
		case ttArray:
			nodes, err := parseRecurse(make([]node, 0), l, t, depth)
			if err != nil {
				return nil, err
			}
			tree = append(tree, &nodeArray{name: t.val, nodes: nodes})
		case ttEnd:
			if end == nil {
				return nil, fmt.Errorf("unopened scope %v", t)
			}
			if t.val != end.val {
				return nil, fmt.Errorf("unmatched tag %v", t)
			}
			return tree, nil
		case ttIfdef:
			nodes, err := parseRecurse(make([]node, 0), l, t, depth)
			if err != nil {
				return nil, err
			}
			tree = append(tree, &nodeIfdef{name: t.val, nodes: nodes})
		case ttIfndef:
			nodes, err := parseRecurse(make([]node, 0), l, t, depth)
			if err != nil {
				return nil, err
			}
			tree = append(tree, &nodeIfndef{name: t.val, nodes: nodes})
		case ttInclude:
			tree = append(tree, &nodeInclude{name: t.val})
		case ttObject:
			nodes, err := parseRecurse(make([]node, 0), l, t, depth)
			if err != nil {
				return nil, err
			}
			tree = append(tree, &nodeObject{name: t.val, nodes: nodes})
		case ttPrint:
			tree = append(tree, &nodePrint{name: t.val})
		case ttString:
			tree = append(tree, &nodeString{val: t.val})
		default:
			panic(fmt.Sprintf("unknown type %q, programmer error", t.tt))
		}
	}
}
