# JSON6 - JSON for Human

[![build](https://github.com/tamboto2000/json6/actions/workflows/build.yml/badge.svg?branch=main)](https://github.com/tamboto2000/json6/actions/workflows/build.yml) [![codecov](https://codecov.io/gh/tamboto2000/json6/branch/main/graph/badge.svg?token=sLilLpAaYq)](https://codecov.io/gh/tamboto2000/json6)

_[Documentation base cloned from JSON6 project]()_

JSON6 is a proposed extension to JSON, [visit this repo](https://github.com/d3x0r/JSON6) for more information. It aims to make it easier for humans to write and maintain by hand. It does this by adding some minimal syntax features directly from ECMAScript 6. Currently, this parser is _**on alpha phase**_, so there will be a lot of improvement. Current functionality that can be used is ```Unmarshal()``` JSON6 string to Go value

## Install
```sh
go get github.com/tamboto2000/json6
```

## Usage example
```go
package main

type CurrentJob struct {
	Title   string `json6:"title"`
	Company string `json:"company"`
	Year    int    `json:"year"`
}

type Profile struct {
	Name       string `json6:"name"`
	Age        int    `json5:"age"`
	Sex        string
	Address    string     `json:"address"`
	CurrentJob CurrentJob `json:"currentJob"`
}

func main() {
    src :=
		`
	{
		name: "Franklin Collin Tamboto",
		age: 21,
		Sex: 'L',
		address: 'Bandung, Jawa Barat',
		currentJob: {
			title: 'Golang Developer',
			company: 'PT. Dwidasa Samsara Indonesia',
			year: 1
		}
	}
	`

	var profile Profile
	if err := Unmarshal([]byte(src), &profile); err != nil {
		t.Error(err)
		return
	}

	// do something
}
```

## Why

JSON isn’t the friendliest to *write*. Keys need to be quoted, objects and
arrays can’t have trailing commas, and comments aren’t allowed — even though
none of these are the case with regular JavaScript today.

That was fine when JSON’s goal was to be a great data format, but JSON’s usage
has expanded beyond *machines*. JSON is now used for writing [configs][ex1],
[manifests][ex2], even [tests][ex3] — all by *humans*.

[ex1]: http://plovr.com/docs.html
[ex2]: https://www.npmjs.org/doc/files/package.json.html
[ex3]: http://code.google.com/p/fuzztester/wiki/JSONFileFormat

There are other formats that are human-friendlier, like YAML _(personally, the author of this repo is trully hate YML and YAML...)_, but changing
from JSON to a completely different format is undesirable in many cases.
JSON6’s aim is to remain close to JSON and JavaScript.

## Features

The following is the exact list of additions to JSON’s syntax introduced by
JSON6. **All of these are optional**, and **MOST of these come from ES5/6**.

## Caveats

Does not include stringify, instead falling back to original (internal) JSON.stringify.
This will cause problems maintaining undefined, Infinity and NaN type values.

### Summary of Changes from JSON5

JSON6 includes all features of JSON5 plus the following.

  - Keyword `undefined`
  - Objects/Strings back-tick quoted strings (no template support, just uses same quote); Object key names can be unquoted.
  - Strings - generous multiline string definition; all javascript character escapes work. \(\0, \x##, \u####, \u\{\} \)
  - Numbers - underscore digit separation in numbers, octal `0o` and binary `0b` formats; all javascript number notations.
  - Arrays - empty members
  - Streaming reader interface
  - (Twice the speed of JSON5; subjective)

### Objects

- Object keys can be unquoted if they do not have ':', ']', '[', '{', '}', ',', any quote or whitespace; keywords will be interpreted as strings.

- Object keys can be single-quoted, (**JSON6**) or back-tick quoted; any valid string

- Object keys can be double-quoted (original JSON).

- Objects can have a single trailing comma. Excessive commas in objects will cause an exception. '{ a:123,,b:456 }' is invalid.

[mdn_variables]: https://developer.mozilla.org/en/Core_JavaScript_1.5_Guide/Core_Language_Features#Variables

### Arrays

- Arrays can have trailing commas. If more than 1 is found, additional empty elements will be added.

- (**JSON6**) Arrays can have comma ( ['test',,,'one'] ), which will result with empty values in the empty places.

### Strings

- Strings can be double-quoted (as per original JSON).

- Strings can be single-quoted.

- Strings can be back-tick (\`) ([grave accent](https://en.wikipedia.org/wiki/Grave_accent)) -quoted.

- Strings can be split across multiple lines; just prefix each newline with a
  backslash. [ES5 [§7.8.4](http://es5.github.com/#x7.8.4)]

- (**JSON6**) all strings will continue keeping every character between the start and end, this allows multi-line strings
  and keep the newlines in the string; if you do not want the newlines they can be escaped as previously mentioned.

- (**JSON5+?**) Strings can have characters emitted using 1 byte hex, interpreted as a utf8 codepoint `\xNN`, 2 and only 2 hex digits must follow `\x`; they may be 4 byte unicode characters `\uUUUU`, 4 and only 4 hex digits must follow `\u`; higher codepoints can be specified with `\u{HHHHH}`, (where H is a hex digit) This is permissive and may accept a single hex digit between `{` and `}`.  All other standard escape sequeneces are also recognized.  Any character that is not recognized as a valid escape character is emitted without the leading escape slash ( for example, `"\012"` will parse as `"012"`

- (**JSON6**) The interpretation of newline is dynamic treating `\r`, `\n`, and `\r\n` as valid combinations of line ending whitespace.  The `\` will behave approrpriately on those combinations.  Mixed line endings like `\n\r?` or `\n\r\n?` are two line endings; 1 for newline, 1 for the \r(follwed by any character), and 1 for the newline, and 1 for the \r\n pair in the second case.

### Numbers

- (**JSON6**) Numbers can have underscores separating digits '_' these are treated as zero-width-non-breaking-space. ([Proposal](https://github.com/tc39/proposal-numeric-separator) with the exception that \_ can preceed or follow . and may be trailing.)

- Numbers can be hexadecimal (base 16).  ( 0x prefix )

- (**JSON6**) Numbers can be binary (base 2).  (0b prefix)

- (**JSON6**) Numbers can be octal (base 8).  (0o prefix)

- (**JSON6**) Decimal Numbers can have leading zeros.  (0 prefix followed by more numbers, without a decimal)

- Numbers can begin or end with a (leading or trailing) decimal point.

- Numbers can include `Infinity`, `-Infinity`,  `NaN`, and `-NaN`. (-NaN results as NaN)

- Numbers can begin with an explicit plus sign.

- Numbers can begin with multiple minus signs. For example '----123' === 123.

### Keyword Values

- (**JSON6**) supports 'undefined' in addition to 'true', 'false', 'null'.


### Comments

- Both inline (single-line using '//' (todo:or '#'?) ) and block (multi-line using \/\* \*\/ ) comments are allowed.
  - `//` comments end at a `\r` or `\n` character; They MAY also end at the end of a document, although a warning is issued at this time.
  - `/*` comments should be closed before the end of a document or stream flush.
  - `/` followed by anything else other than `/` or `*` is an error.


## Example

The following is a contrived example, but it illustrates most of the features:

```js
{
	foo: 'bar',
	while: true,
	nothing : undefined, // why not?

	this: 'is a \
multi-line string',

	thisAlso: 'is a
multi-line string; but keeps newline',

	// this is an inline comment
	here: 'is another', // inline comment

	/* this is a block comment
	   that continues on another line */

	hex: 0xDEAD_beef,
	binary: 0b0110_1001,
	decimal: 123_456_789,
	octal: 0o123,
	half: .5,
	delta: +10,
	negative : ---123,
	to: Infinity,   // and beyond!

	finally: 'a trailing comma',
	oh: [
		"we shouldn't forget",
		'arrays can have',
		'trailing commas too',
	],
}
```

This implementation’s own [package.JSON6](https://github.com/d3x0r/JSON6/blob/master/package.JSON6) is more realistic:

```js
// This file is written in JSON6 syntax, naturally, but npm needs a regular
// JSON file, so compile via `npm run build`. Be sure to keep both in sync!

{
	name: 'JSON6',
	version: '0.1.105',
	description: 'JSON for the ES6 era.',
	keywords: ['json', 'es6'],
	author: 'd3x0r <d3x0r@github.com>',
	contributors: [
		// TODO: Should we remove this section in favor of GitHub's list?
		// https://github.com/d3x0r/JSON6/contributors
	],
	main: 'lib/JSON6.js',
	bin: 'lib/cli.js',
	files: ["lib/"],
	dependencies: {},
	devDependencies: {
		gulp: "^3.9.1",
		'gulp-jshint': "^2.0.0",
		jshint: "^2.9.1",
		'jshint-stylish': "^2.1.0",
		mocha: "^2.4.5"
	},
	scripts: {
		build: 'node ./lib/cli.js -c package.JSON6',
		test: 'mocha --ui exports --reporter spec',
			// TODO: Would it be better to define these in a mocha.opts file?
	},
	homepage: 'http://github.com/d3x0r/JSON6/',
	license: 'MIT',
	repository: {
		type: 'git',
		url: 'https://github.com/d3x0r/JSON6',
	},
}
```
## Community

Join the [Google Group](http://groups.google.com/group/JSON6) if you’re
interested in JSON6 news, updates, and general discussion.
Don’t worry, it’s very low-traffic.

The [GitHub wiki](https://github.com/d3x0r/JSON6/wiki) (will be) a good place to track
JSON6 support and usage. Contribute freely there!

[GitHub Issues](https://github.com/d3x0r/JSON6/issues) is the place to
formally propose feature requests and report bugs. Questions and general
feedback are better directed at the Google Group.

## License
MIT