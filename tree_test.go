// Copyright 2015 Seth Bunce. All rights reserved. Use of this source code is
// governed by a BSD-style license that can be found in the LICENSE file.
package stem

import (
	"reflect"
	"testing"
)

func TestTree(t *testing.T) {
	tests := []struct{
		name string // name of test printed with errors.
		src  string // src is the template.
		want []node // want this tree.
	}{
		{
			name: "array",
			src:  "{{#a}}{{*b}}{{/a}}",
			want: []node{
				&nodeArray{
					name:  "a",
					nodes: []node{
						&nodePrint{
							name: "b",
						},
					},
				},
			},
		},
		{
			name: "comment",
			src:  "{{!a}}",
			want: []node{},
		},
		{
			name: "ifdef",
			src:  "{{+a}}{{*b}}{{/a}}",
			want: []node{
				&nodeIfdef{
					name:  "a",
					nodes: []node{
						&nodePrint{
							name: "b",
						},
					},
				},
			},
		},
		{
			name: "ifndef",
			src:  "{{-a}}{{*b}}{{/a}}",
			want: []node{
				&nodeIfndef{
					name:  "a",
					nodes: []node{
						&nodePrint{
							name: "b",
						},
					},
				},
			},
		},
		{
			name: "include",
			src:  "{{>a}}",
			want: []node{
				&nodeInclude{
					name: "a",
				},
			},
		},
		{
			name: "object",
			src:  "{{$a}}{{*b}}{{/a}}",
			want: []node{
				&nodeObject{
					name:  "a",
					nodes: []node{
						&nodePrint{
							name: "b",
						},
					},
				},
			},
		},
		{
			name: "string",
			src:  "abc",
			want: []node{
				&nodeString{
					val: "abc",
				},
			},
		},
		{
			name: "nested array",
			src:  "{{#a}}{{#a}}{{*b}}{{/a}}{{/a}}",
			want: []node{
				&nodeArray{
					name:  "a",
					nodes: []node{
						&nodeArray{
							name: "a",
							nodes: []node{
								&nodePrint{
									name: "b",
								},	
							},
						},
					},
				},
			},
		},
	}

	for _, test := range tests {
		got, err := parse(test.name, test.src)
		if err != nil {
			t.Fatalf("couldn't parse template: %v", err)
		}
		if !reflect.DeepEqual(got, test.want) {
			t.Fatalf("test %q, got %v, want %v", test.name, got, test.want)
		}
	}
}
