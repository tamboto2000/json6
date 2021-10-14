package json6

import (
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

	dec, err := newDecoderFromBytes([]byte(input), nil)
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

	dec, err := newDecoderFromBytes([]byte(input), nil)
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
