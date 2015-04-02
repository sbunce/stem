// Copyright 2015 Seth Bunce. All rights reserved. Use of this source code is
// governed by a BSD-style license that can be found in the LICENSE file.
package stem

import (
	"bytes"
	"testing"
)

// ttest contains tests for go data.
var ttest = []struct {
	name string                 // name of test printed with errors.
	tmpl string                 // tmpl is the raw template.
	data map[string]interface{} // data combined with template.
	want string                 // want is the output we want.
}{
	{
		name: "ifdef",
		tmpl: "0{{+a}}1{{/a}}2",
		data: map[string]interface{}{"a": ""},
		want:  "012",
	},
	{
		name: "ifndef",
		tmpl: "0{{-a}}1{{/a}}2",
		data: map[string]interface{}{"a": ""},
		want:  "02",
	},
	{
		name: "print",
		tmpl: "{{*a}}",
		data: map[string]interface{}{"a": "0"},
		want:  "0",
	},
	{
		name: "object",
		tmpl: "{{$a}}{{*b}}{{/a}}",
		data: map[string]interface{}{
			"a": map[string]interface{}{"b": "0"},
		},
		want:  "0",
	},
	{
		name: "object nest",
		tmpl: "{{$a}}{{$a}}{{*b}}{{/a}}{{/a}}",
		data: map[string]interface{}{
			"a": map[string]interface{}{"b": "0"},
		},
		want:  "0",
	},
	{
		name: "array object",
		tmpl: "{{#a}}{{*b}}{{/a}}",
		data: map[string]interface{}{
			"a": []interface{}{
				map[string]interface{}{"b": "0"},
				map[string]interface{}{"b": "1"},
			},
		},
		want: "01",
	},
	{
		name: "array",
		tmpl: "{{#a}}{{*}}{{/a}}",
		data: map[string]interface{}{
			"a": []interface{}{0,1},
		},
		want: "01",
	},
}

// ttestJSON contains tests for JSON data.
var ttestJSON = []struct {
	name string // name of test printed with errors.
	tmpl string // tmpl is the template.
	data string // data is JSON combined with the template.
	want string // want this output.
}{
	{
		name: "ifdef JSON",
		tmpl: "0{{+a}}1{{/a}}2",
		data: `{"a": ""}`,
		want: "012",
	},
	{
		name: "ifndef JSON",
		tmpl: "0{{-a}}1{{/a}}2",
		data: `{"a": ""}`,
		want: "02",
	},
	{
		name: "print JSON",
		tmpl: "{{*a}}",
		data: `{"a": "0"}`,
		want: "0",
	},
	{
		name: "object JSON",
		tmpl: "{{$a}}{{*b}}{{/a}}",
		data: `{"a": {"b": "0"}}`,
		want:  "0",
	},
	{
		name: "object nest JSON",
		tmpl: "{{$a}}{{$a}}{{*b}}{{/a}}{{/a}}",
		data: `{"a": {"b": "0"}}`,
		want: "0",
	},
	{
		name: "array object JSON",
		tmpl: "{{#a}}{{*b}}{{/a}}",
		data: `{"a": [{"b": "0"}, {"b": "1"}]}`,
		want: "01",
	},
	{
		name: "array JSON",
		tmpl: "{{#a}}{{*}}{{/a}}",
		data: `{"a": [0, 1]}`,
		want: "01",
	},
}

func TestTemplate(t *testing.T) {
	for _, test := range ttest {
		tmpl, err := Parse(test.tmpl)
		if err != nil {
			t.Fatalf("couldn't parse template: %v", err)
		}
		got := bytes.NewBuffer(make([]byte, 0))
		if err := tmpl.Execute(got, test.data); err != nil {
			t.Fatalf("couldn't execute template: %v", err)
		}
		if got.String() != test.want {
			t.Fatalf("test %q, got %q, want %q", test.name, got.String(), test.want)
		}
	}
}

func TestTemplateJSON(t *testing.T) {
	for _, test := range ttestJSON {
		tmpl, err := Parse(test.tmpl)
		if err != nil {
			t.Fatalf("couldn't parse template: %v", err)
		}
		got := bytes.NewBuffer(make([]byte, 0))
		if err := tmpl.ExecuteJSON(got, test.data); err != nil {
			t.Fatalf("couldn't execute template: %v", err)
		}
		if got.String() != test.want {
			t.Fatalf("test %q, got %q, want %q", test.name, got.String(), test.want)
		}
	}
}

func TestTemplateInclude(t *testing.T) {
	set := NewSet()

	// Create template foo.
	foo, err := Parse("0{{>bar}}2");
	if err != nil {
		t.Fatalf("couldn't parse template: %v", err)
	}
	foo.SetName("foo")
	set.Add(foo)

	// Create template bar which will be included in foo.
	bar, err := Parse("1")
	if err != nil {
		t.Fatalf("couldn't parse template: %v", err)
	}
	bar.SetName("bar")
	set.Add(bar)

	// Execute the template.
	got := bytes.NewBuffer(nil)
	if err := set.Execute(got, "foo", map[string]interface{}{}); err != nil {
		t.Fatalf("couldn't execute template: %v", err)
	}
	if got, want := got.String(), "012"; got != want {
		t.Fatalf("got %v, want %v", got, want)
	}
}
