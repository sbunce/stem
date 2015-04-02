# Go json templating package.

##Overview
JSON + Template = Output

##Inspiration
http://google-ctemplate.googlecode.com

http://www.json.org

http://code.google.com/p/json-template

http://json-template.googlecode.com/svn/trunk/doc/On-Design-Minimalism.html

##Features
Conditions.
Includes.
Lexical scope.
Repeated sections.
Change delimiter.

##Tags
	{{!a}}          Comment.
	{{#a}}...{{/a}} Enter array.
	{{$a}}...{{/a}} Enter object.
	{{+a}}...{{/a}} Render section if defined.
	{{-a}}...{{/a}} Render section if not defined.
	{{>a}}          Include template.
	{{*a}}          Print. To access element of array use "{{*}}".
	{{=<ld> <rd>}}  Change delimiters.

##Examples
	Print a symbol.
	JSON:
		{"a": "foo", "b": "bar"}
	Template:
		{{*a}}{{*b}}
	Output:
		foobar

	Conditional output.
	JSON:
		{"a": ""}
	Template:
		{{!"+" outputs if symbol defined, "-" is the opposite.}}
		{{+a}}foo{{/a}}
		{{-a}}bar{{/a}}{{!"bar" is not output because a is defined.}}
	Output:
		foo

	Enter array of documents.
	JSON:
	{"a": [{"b": "foo"}, {"b": "bar"}]}
	Template:
		{{!Evaluated once for every element.}}
		{{#a}}{{*b}}{{/a}}
	Output:
		foobar

	Enter array of non-documents.
	JSON:
		{"a": ["foo", "bar"]}
	Template:
		{{!Notice how we use a print tag without a name.}}
		{{#a}}{{*}}{{/a}}
	Output:
		foobar

	Enter object.
	JSON:
		{"a": {"b": "foo"}}
	Template:
		{{$a}}{{*b}}{{/a}}
	Output:
		foo

	Scoped symbol lookup example 1.
	JSON:
		{"a": {"b": "foo"}, "c": "bar"}
	Template:
		{{!Notice how "c" is found in the outer scope.
		Name lookup is from inner to outer scope, just like C.}}
	{{$a}}{{*b}}{{*c}}{{/a}}
	Output:
		foobar

	Scoped symbol lookup example 2.
	JSON:
		{"a": {"b": "foo", "c": "bar"}, "c": "baz"}
	Template:
		{{!Notice how the inner "c" shadows the outer "c".
		This only happens when in the "a" scope.}}
		{{$a}}{{*b}}{{*c}}{{/a}}{{*c}}
	Output:
		foobarbaz

	Change delimiters.
	This can be used when your document contains the default delimiters.
	JSON:
		{"a": "foo", "b": "bar", "c": "baz"}
	Template:
		{{*a}}{{=[[ ]]}}[[*b]][[=<< >>]]<<*c>>
	Output:
		foobarbaz
