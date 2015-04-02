// Copyright 2015 Seth Bunce. All rights reserved. Use of this source code is
// governed by a BSD-style license that can be found in the LICENSE file.
package stem

import (
	"regexp"
)

// Filter input before rendering.
// Run before execute for performance.
type Filter int
const (
	// Replace consecutive \n with single \n.
	// Remove blanks lines that result from having tags on separate lines.
	NoBlankLines Filter = 1 << iota

	// Remove space at start of line.
	TrimLeftSpace Filter = 1 << iota

	// Remove space at end of line.
	TrimRightSpace Filter = 1 << iota
)

// Regexps for filters.
var (
	blankLines = regexp.MustCompile("\n+")
	leftSpace  = regexp.MustCompile("\n[ \t]+")
	rightSpace = regexp.MustCompile("[ \t]+\n")
)

// filter recurses through a parse tree and filters all string nodes.
func filter(tree []node, filters Filter) {
	for _, n := range tree {
		switch nt := n.(type) {
		case *nodeArray:
			filter(nt.nodes, filters)
		case *nodeIfdef:
			filter(nt.nodes, filters)
		case *nodeIfndef:
			filter(nt.nodes, filters)
		case *nodeObject:
			filter(nt.nodes, filters)
		case *nodeString:
			if filters & TrimLeftSpace != 0 {
				nt.val = leftSpace.ReplaceAllString(nt.val, "\n")
			}
			if filters & TrimRightSpace != 0 {
				nt.val = rightSpace.ReplaceAllString(nt.val, "\n")
			}
			if filters & NoBlankLines != 0 {
				nt.val = blankLines.ReplaceAllString(nt.val, "\n")
			}
		}
	}
}
