package json6

import (
	"fmt"
	"testing"
)

func TestDecodeObject(t *testing.T) {
	input :=
		`
		ident: 'single quote string',
		'sigle quote ident': "double quote string",
		"double quote ident":` + "`back tick string`," +
			"`back tick ident`: " + `{
			// sigle line comment
			innerObject: "inner val",
			identwith\u1234unicode: "string with\u1234unicode",
			identWith\x12hexaEsc: "string with\x12 hexa escape",
			moreString: "Lorem ipsum dolor\
 sit amet",
			nullVal: null,
			undefinedVal: undefined //comment
			,
			boolFalseVal: false//comment
			,
			boolTrueVal: true//comment
			,
			minusInt1: -123//comment
			,
			minusInt2: ---123//comment
			,
			plusInt1: 123//comment
			,
			plusInt2: +123//comment
			,
			plusInt3: --123,
			hexaDecimal1: 0x123,
			hexaDecimal2: 0X123//comment
			,
			binary1: 0b1010,
			binary2: 0B1010//comment
			,
			octalDecimal1: 0o123,
			octalDecimal2: 0O123//comment
			,
			double1: 0.123//comment
			,
			double2: .123//comment
			,
			double3: 123.,
			exponent1: 1e123//comment
			,
			exponent2: 1e-123,
			exponent3: 1e+123,
			exponent4: 1E123,
			exponent5: 1E-123,
			exponent6: 1E+123,
			exponent7: 1.e123,
			exponent8: 1.e-123,
			exponent9: 1.e+123,
			exponent10: 1.E123,
			exponent11: 1.E-123,
			exponent12: 1.E+123,
			exponent13: .1e123,
			exponent14: .1e-123,
			exponent15: .1e+123,
			exponent16: .1E123,
			exponent17: .1E-123,
			exponent18: .1E+123,
			NaNNum1: NaN,
			NaNNum2: -NaN//comment
			,
			NaNNum3: +NaN,
			InfinityNum1: Infinity//comment
			,
			InfinityNum2: +Infinity,
			InfinityNum3: -Infinity//comment
		},

		/*
			multiline comment
		*/
	}
	`

	var val struct{}
	dec, err := newDecoderFromBytes([]byte(input), &val)
	if err != nil {
		t.Error(err.Error())
		return
	}

	_, err = decodeObject(dec.lx.tokenReader)
	if err != nil {
		t.Error(err.Error())
	}

	// fmt.Printf("%#v\n", val.objVal)
}

func TestDecodeArray(t *testing.T) {
	input :=
		`
	
		{
			"amount": 20000,
			"timestamp": "2021-10-12",
			"description": "MO1922 - Teknisi EDC PT. Mahapay di Cibubur",
			"status": "Approved",
			"category": "incentive"
		},
		{
			"amount": 5000,
			"timestamp": "2021-10-12",
			"description": "MO1922 - Teknisi EDC PT. Mahapay di Jatinegara",
			"status": "Completed",
			"category": "fee"
		},
		{
			"amount": -30000,
			"timestamp": "2021-10-12",
			"description": "Bank Mandiri - Andini Septia",
			"status": "Transfer Out",
			"category": "price"
		},
		{
			"amount": 0,
			"timestamp": "2021-10-12",
			"description": "MO1920 - Nobar EPL Week 6",
			"status": "Canceled",
			"category": "others"
		},,,,
		{
			"amount": 0,
			"timestamp": "2021-10-11",
			"description": "MO1919 - Komando IDM Asem Baris",
			"status": "Revoked",
			"category": "incentive"
		},
		[],
	]
	`

	var val struct{}
	dec, err := newDecoderFromBytes([]byte(input), &val)
	if err != nil {
		t.Error(err.Error())
		return
	}

	_, err = decodeArray(dec.lx.tokenReader)
	if err != nil {
		t.Error(err.Error())
	}

	// fmt.Printf("%#v\n", val.arrVal)
}

func TestDecodeValue(t *testing.T) {
	input :=
		`
		{
			ident// comment
			: 'single quote string',
			'sigle quote ident': "double quote string",
			"double quote ident":` + "`back tick string`," +
			"`back tick ident`: " + `{
				// sigle line comment
				innerObject: "inner val",
				identwith\u1234unicode: "string with\u1234unicode",
				identWith\x12hexaEsc: "string with\x12 hexa escape",
				moreString: "Lorem ipsum dolor\
 sit amet",
				nullVal: null,
				undefinedVal: undefined //comment
				,
				boolFalseVal: false//comment
				,
				boolTrueVal: true//comment
				,
				minusInt1: -123//comment
				,
				minusInt2: ---123//comment
				,
				plusInt1: 123//comment
				,
				plusInt2: +123//comment
				,
				plusInt3: --123,
				hexaDecimal1: 0x123,
				hexaDecimal2: 0X123//comment
				,
				binary1: 0b1010,
				binary2: 0B1010//comment
				,
				octalDecimal1: 0o123,
				octalDecimal2: 0O123//comment
				,
				double1: 0.123//comment
				,
				double2: .123//comment
				,
				double3: 123.,
				exponent1: 1e123//comment
				,
				exponent2: 1e-123,
				exponent3: 1e+123,
				exponent4: 1E123,
				exponent5: 1E-123,
				exponent6: 1E+123,
				exponent7: 1.e123,
				exponent8: 1.e-123,
				exponent9: 1.e+123,
				exponent10: 1.E123,
				exponent11: 1.E-123,
				exponent12: 1.E+123,
				exponent13: .1e123,
				exponent14: .1e-123,
				exponent15: .1e+123,
				exponent16: .1E123,
				exponent17: .1E-123,
				exponent18: .1E+123,
				NaNNum1: NaN,
				NaNNum2: -NaN//comment
				,
				NaNNum3: +NaN,
				InfinityNum1: Infinity//comment
				,
				InfinityNum2: +Infinity,
				InfinityNum3: -Infinity//comment
			},

			/*
				multiline comment
			*/
		}
		`

	var val struct{}
	dec, err := newDecoderFromBytes([]byte(input), &val)
	if err != nil {
		t.Error(err.Error())
		return
	}

	if err := dec.decodeValue(); err != nil {
		t.Error(err.Error())
	}
}

func TestUnmarshal(t *testing.T) {
	src := []byte("-3000")
	expected := -3000
	var val int
	if err := Unmarshal(src, &val); err != nil {
		t.Error(err.Error())
		return
	}

	if val != expected {
		t.Errorf("unexpected value '%d', expecting '%d'", val, expected)
	}
}

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

func TestUnmarshalObject(t *testing.T) {
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

	var val Profile
	if err := Unmarshal([]byte(src), &val); err != nil {
		t.Error(err)
		return
	}

	fmt.Printf("%#v\n", val)
}

func TestUnmarshalObjectMismatchType(t *testing.T) {
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
			year: '1'
		}
	}
	`

	var val Profile
	if err := Unmarshal([]byte(src), &val); err == nil {
		t.Error("expecting type mismatch error")
		return
	} else {
		fmt.Println(err.Error())
	}
}

func TestUnmarshalObjectToMap(t *testing.T) {
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

	val := make(map[string]interface{})
	if err := Unmarshal([]byte(src), &val); err != nil {
		t.Error(err)
		return
	}

	fmt.Printf("%#v\n", val)
}

func TestUnmarshalObjectToInvalidMap(t *testing.T) {
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

	val := make(map[int]interface{})
	if err := Unmarshal([]byte(src), &val); err == nil {
		t.Error("expecting error invalid map")
		return
	} else {
		fmt.Println(err.Error())
	}
}

func TestUnmarshalObjectToInvalidType(t *testing.T) {
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

	val := 0
	if err := Unmarshal([]byte(src), &val); err == nil {
		t.Error("expecting error invalid type")
		return
	} else {
		fmt.Println(err.Error())
	}
}

func TestUnmarshalObjectToInterface(t *testing.T) {
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

	var val interface{}
	if err := Unmarshal([]byte(src), &val); err != nil {
		t.Error(err)
		return
	}

	fmt.Printf("%#v\n", val)
}

func TestUnmarshalArrayToSlice(t *testing.T) {
	src :=
		`
	[
		1,
		-2,
		3,
		0x4,
		0e5,
	]
	`

	var val []int
	if err := Unmarshal([]byte(src), &val); err != nil {
		t.Error(err)
		return
	}

	fmt.Printf("%#v\n", val)
}

func TestUnmarshalArrayWithNullElemToSlice(t *testing.T) {
	src :=
		`
	[
		1,
		null,
		-2,
		3,
		0x4,
		0e5,
	]
	`

	var val []int
	if err := Unmarshal([]byte(src), &val); err != nil {
		t.Error(err)
		return
	}

	fmt.Printf("%#v\n", val)
}

func TestUnmarshalArrayWithUndefinedElemToSlice(t *testing.T) {
	src :=
		`
	[
		1,
		undefined,
		-2,
		3,
		0x4,
		0e5,
	]
	`

	var val []int
	if err := Unmarshal([]byte(src), &val); err != nil {
		t.Error(err)
		return
	}

	fmt.Printf("%#v\n", val)
}

func TestUnmarshalArrayToArray(t *testing.T) {
	src :=
		`
	[
		1,
		-2,
		null,
		0x4,
		0e5,
	]
	`

	var val [3]int
	if err := Unmarshal([]byte(src), &val); err != nil {
		t.Error(err)
		return
	}

	fmt.Printf("%#v\n", val)
}

func TestUnmarshalArrayToInterface(t *testing.T) {
	src :=
		`
	[
		1,
		-2,
		null,
		0x4,
		0e5,
	]
	`

	var val interface{}
	if err := Unmarshal([]byte(src), &val); err != nil {
		t.Error(err)
		return
	}

	fmt.Printf("%#v\n", val)
}

func TestUnmarshalStrFromd3x0r_JSON6ToMap1(t *testing.T) {
	src :=
		`
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
	`

	val := make(map[string]interface{})
	if err := Unmarshal([]byte(src), &val); err != nil {
		t.Error(err)
		return
	}

	fmt.Printf("%#v\n", val)
}

func TestUnmarshalStrFromd3x0r_JSON6ToMap2(t *testing.T) {
	src :=
		`
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
	`

	val := make(map[string]interface{})
	if err := Unmarshal([]byte(src), &val); err != nil {
		t.Error(err)
		return
	}

	fmt.Printf("%#v\n", val)
}
