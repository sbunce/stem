// Copyright 2015 Seth Bunce. All rights reserved. Use of this source code is
// governed by a BSD-style license that can be found in the LICENSE file.

package stem

import (
	"bytes"
	"testing"
)

func TestFilter(t *testing.T) {
	tests := []struct{
		name string // name of the test printed with errors.
		flag Filter // flag for filter to enable.
		tmpl string // tmpl is the template text.
		want string // want this output.
	}{
		{
			name: "blank lines",
			flag: NoBlankLines,
			tmpl: "foo\n\nbar",
			want: "foo\nbar",
		},
		{
			name: "left space",
			flag: TrimLeftSpace,
			tmpl: "foo\n bar",
			want: "foo\nbar",
		},
		{
			name: "left space",
			flag: TrimRightSpace,
			tmpl: "foo \nbar",
			want: "foo\nbar",
		},
	}
	for _, test := range tests {
		tmpl, err := Parse(test.tmpl)
		if err != nil {
			t.Fatalf("couldn't parse template: %v", err)
		}
		tmpl.Filter(test.flag)
		got := bytes.NewBuffer(make([]byte, 0))
		if err := tmpl.Execute(got, map[string]interface{}{}); err != nil {
			t.Fatalf("couldn't execute template: %v", err)
		}
		if got.String() != test.want {
			t.Fatalf("test %q, got %v, want %v", test.name, got.String(), test.want)
		}
	}
}
