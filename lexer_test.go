package json6

import (
	"bytes"
	"fmt"
	"testing"
)

// -------------------- Tests for isCharValidHexa() --------------------

// TestIsCharValidHexa test isCharValidHexa() behavior
func TestIsCharValidHexa(t *testing.T) {
	chars := []rune{
		'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
		'a', 'b', 'c', 'd', 'e', 'f',
		'A', 'B', 'C', 'D', 'E', 'F',
	}

	for _, char := range chars {
		if !isCharValidHexa(char) {
			t.Errorf("character '%s' should be valid hexadecimal", string([]rune{char}))
		}
	}
}

// TestIsCharValidHexaInvalid test isCharValidHexa() behavior
// when checking invalid hexadecimal char
func TestIsCharValidHexaInvalid(t *testing.T) {
	chars := []rune{'q', 'w', 'u', 'r', 't', 'y'}

	for _, char := range chars {
		if isCharValidHexa(char) {
			t.Errorf("character '%s' should be invalid hexadecimal", string([]rune{char}))
		}
	}
}

// -------------------- Tests for isCharValidOctal() --------------------

// TestIsCharValidOctal test isCharValidOctal() behavior
func TestIsCharValidOctal(t *testing.T) {
	chars := []rune{'0', '1', '2', '3', '4', '5', '6', '7'}
	for _, char := range chars {
		if !isCharValidOctal(char) {
			t.Errorf("character '%s' should be valid octaldecimal", string([]rune{char}))
		}
	}
}

// TestIsCharValidOctalInvalid test isCharValidHexa() behavior
// when checking invalid octaldecimal char
func TestIsCharValidOctalInvalid(t *testing.T) {
	chars := []rune{'8', '9', 'a', 'b', 'c'}
	for _, char := range chars {
		if isCharValidOctal(char) {
			t.Errorf("character '%s' should be invalid octaldecimal", string([]rune{char}))
		}
	}
}

// -------------------- Tests for isCharPunct() --------------------

// TestIsCharPunct test isCharPunct() behavior
func TestIsCharPunct(t *testing.T) {
	chars := []rune{'{', '}', '[', ']', ':', ','}
	for _, char := range chars {
		if !isCharPunct(char) {
			t.Errorf("character '%s' should be valid punctuation", string([]rune{char}))
		}
	}
}

// TestIsCharPunctInvalid test isCharPunct() behavior
// when checking invalid punctuation
func TestIsCharPunctInvalid(t *testing.T) {
	chars := []rune{'!', '@', '#'}
	for _, char := range chars {
		if isCharPunct(char) {
			t.Errorf("character '%s' should be invalid punctuation", string([]rune{char}))
		}
	}
}

// TestFetchNull test Lexer.fetchNull() behavior
func TestFetchNull(t *testing.T) {
	reader := bytes.NewReader([]byte("ull"))
	lex := NewLexer(reader)
	if err := lex.fetchNull(); err != nil {
		t.Error(err.Error())
		return
	}

	expected := "null"
	token := lex.tokens[0]
	if string(token.chars) != expected {
		t.Errorf("unexpected %s, expecting %s", string(token.chars), expected)
	}
}

// TestFetchNullNotNullToken test Lexer.fetchNull() behavior when input chars is not a 'null' token, but instead
// will be automatically acknowledged as identifier, example: "nullbuthisident"
func TestFetchNullNotNullToken(t *testing.T) {
	input := "ull$_abc"
	expected := "null$_abc"
	reader := bytes.NewReader([]byte(input))
	lex := NewLexer(reader)
	if err := lex.fetchNull(); err != nil {
		t.Error(err.Error())
		return
	}

	if len(lex.tokens) != 1 {
		t.Error("expecting 1 token to be fetched")
		return
	}

	token := lex.tokens[0]
	if token.t != TokenIdentifier {
		t.Errorf("unexpected token type %d (%s), expecting token type %d (%s)", token.t, tokenTypeMap[token.t], TokenIdentifier, tokenTypeMap[TokenIdentifier])
		return
	}

	if token.String() != expected {
		t.Errorf("unexpected %s, expecting %s", token.String(), expected)
	}
}

// TestFetchNullEndByPunct test Lexer.fetchNull() behavior when null token ended by punctuator
func TestFetchNullEndByPunct(t *testing.T) {
	expect1 := "null"
	expect2 := "{"
	input := "ull{"

	reader := bytes.NewReader([]byte(input))
	lex := NewLexer(reader)
	if err := lex.fetchNull(); err != nil {
		t.Error(err.Error())
		return
	}

	if len(lex.tokens) != 2 {
		t.Error("expecting 2 tokens to be fetched")
		return
	}

	token1 := lex.tokens[0]
	token2 := lex.tokens[1]

	if token1.String() != expect1 {
		t.Errorf("unexpected %s, expecting %s", token1.String(), expect1)
		return
	}

	if token2.String() != expect2 {
		t.Errorf("unexpected %s, expecting %s", token2.String(), expect2)
		return
	}
}

// TestFetchUnicodeEscape test Lexer.fetchUnicodeEscape() behavior
func TestFetchUnicodeEscape(t *testing.T) {
	expected := "212A"
	reader := bytes.NewReader([]byte(expected))
	lex := NewLexer(reader)
	if err := lex.fetchUnicodeEscape(); err != nil {
		t.Error(err.Error())
		return
	}

	token := lex.token
	if token.String() != expected {
		t.Errorf("unexpected %s, expecting %s", token.String(), expected)
	}
}

// TestFetchUnicodeEscapeWithBrackets test Lexer.fetchUnicodeEscape() behavior when fetching
// unicode escape with '{}', example: \u{ffff}
func TestFetchUnicodeEscapeWithBrackets(t *testing.T) {
	expected := "{212A}"
	reader := bytes.NewReader([]byte(expected))
	lex := NewLexer(reader)
	if err := lex.fetchUnicodeEscape(); err != nil {
		t.Error(err.Error())
		return
	}

	token := lex.token
	if token.String() != expected {
		t.Errorf("unexpected %s, expecting %s", token.String(), expected)
	}
}

// TestFetchIdentifierIsBeginTrue test Lexer.fetchIdentifier() with isBegin = true
func TestFetchIdentifierIsBeginTrue(t *testing.T) {
	expected := "$_ident"
	input := "_ident "
	reader := bytes.NewReader([]byte(input))
	lex := NewLexer(reader)
	if err := lex.fetchIdentifier(true, '$'); err != nil {
		t.Error(err.Error())
		return
	}

	if len(lex.tokens) != 1 {
		t.Error("expecting 1 token to be fetched")
		return
	}

	token := lex.tokens[0]
	if token.String() != expected {
		t.Errorf("unexpected '%s', expecting '%s'", token.String(), expected)
	}
}

// TestFetchIdentifierIsBeginFalse test Lexer.fetchIdentifier() behavior with isBegin = false
func TestFetchIdentifierIsBeginFalse(t *testing.T) {
	expected := "$_ident"
	input := "_ident"
	reader := bytes.NewReader([]byte(input))
	lex := NewLexer(reader)
	if err := lex.fetchIdentifier(false, '$'); err != nil {
		t.Error(err.Error())
		return
	}

	if len(lex.tokens) != 1 {
		t.Error("expecting 1 token to be fetched")
		return
	}

	token := lex.tokens[0]
	if token.String() != expected {
		t.Errorf("unexpected %s, expecting %s", token.String(), expected)
	}
}

// TestFetchTrueBool test Lexer.fetchTrueBool behavior
func TestFetchTrueBool(t *testing.T) {
	expected := "true"
	input := "rue"
	reader := bytes.NewReader([]byte(input))
	lex := NewLexer(reader)
	if err := lex.fetchTrueBool(); err != nil {
		t.Error(err.Error())
		return
	}

	if len(lex.tokens) != 1 {
		t.Error("expecting 1 token to be fetched")
		return
	}

	token := lex.tokens[0]
	if token.t != TokenBool {
		t.Errorf("unexpected token type %d (%s), expecting token type %d (%s)", token.t, tokenTypeMap[token.t], TokenBool, tokenTypeMap[TokenBool])
		return
	}

	if token.String() != expected {
		t.Errorf("unexpected %s, expecting %s", token.String(), expected)
	}
}

// TestFetchTrueBoolEndByPunct test Lexer.fetchTrueBool behavior when boolean token end by punctuator
func TestFetchTrueBoolEndByPunct(t *testing.T) {
	expect1 := "true"
	expect2 := "{"
	input := "rue{"

	reader := bytes.NewReader([]byte(input))
	lex := NewLexer(reader)
	if err := lex.fetchTrueBool(); err != nil {
		t.Error(err.Error())
		return
	}

	if len(lex.tokens) != 2 {
		t.Error("expecting 2 tokens to be fetched")
		return
	}

	token1 := lex.tokens[0]
	token2 := lex.tokens[1]

	if token1.String() != expect1 {
		t.Errorf("unexpected %s, expecting %s", token1.String(), expect1)
		return
	}

	if token2.String() != expect2 {
		t.Errorf("unexpected %s, expecting %s", token2.String(), expect2)
		return
	}
}

// TestFetchFalseBool test Lexer.fetchFalseBool behavior
func TestFetchFalseBool(t *testing.T) {
	expected := "false"
	input := "alse"
	reader := bytes.NewReader([]byte(input))
	lex := NewLexer(reader)
	if err := lex.fetchFalseBool(); err != nil {
		t.Error(err.Error())
		return
	}

	if len(lex.tokens) != 1 {
		t.Error("expecting 1 token to be fetched")
		return
	}

	token := lex.tokens[0]
	if token.t != TokenBool {
		t.Errorf("unexpected token type %d (%s), expecting token type %d (%s)", token.t, tokenTypeMap[token.t], TokenBool, tokenTypeMap[TokenBool])
		return
	}

	if token.String() != expected {
		t.Errorf("unexpected %s, expecting %s", token.String(), expected)
	}
}

// TestFetchFalseBoolEndByPunct test Lexer.fetchFalseBool behavior when boolean token end by punctuator
func TestFetchFalseBoolEndByPunct(t *testing.T) {
	expect1 := "false"
	expect2 := "{"
	input := "alse{"

	reader := bytes.NewReader([]byte(input))
	lex := NewLexer(reader)
	if err := lex.fetchFalseBool(); err != nil {
		t.Error(err.Error())
		return
	}

	if len(lex.tokens) != 2 {
		t.Error("expecting 2 tokens to be fetched")
		return
	}

	token1 := lex.tokens[0]
	token2 := lex.tokens[1]

	if token1.String() != expect1 {
		t.Errorf("unexpected %s, expecting %s", token1.String(), expect1)
		return
	}

	if token2.String() != expect2 {
		t.Errorf("unexpected %s, expecting %s", token2.String(), expect2)
		return
	}
}

// TestFetchUndefined test Lexer.fetchUndefined() behavior
func TestFetchUndefined(t *testing.T) {
	expected := "undefined"
	input := "ndefined"
	reader := bytes.NewReader([]byte(input))
	lex := NewLexer(reader)
	if err := lex.fetchUndefined(); err != nil {
		t.Error(err.Error())
		return
	}

	if len(lex.tokens) != 1 {
		t.Error("expecting 1 token to be fetched")
		return
	}

	token := lex.tokens[0]
	if token.Type() != TokenUndefined {
		t.Errorf("unexpected token type %d (%s), expecting %d (%s)", token.Type(), token.TypeString(), TokenUndefined, tokenTypeMap[TokenUndefined])
		return
	}

	if token.String() != expected {
		t.Errorf("unexpected %s, expecting %s", token.String(), expected)
	}
}

// TestFetchHexaEscape test Lexer.fetchHexaEscape() behavior
func TestFetchHexaEscape(t *testing.T) {
	expected := "Ff"
	reader := bytes.NewReader([]byte(expected))
	lex := NewLexer(reader)
	if err := lex.fetchHexaEscape(); err != nil {
		t.Error(err.Error())
		return
	}

	token := lex.token
	if token.String() != expected {
		t.Errorf("unexpected %s, expecting %s", token.String(), expected)
	}
}

// TestFetchHexaNumber test Lexer.fetchHexaNumber behavior
func TestFetchHexaNumber(t *testing.T) {
	expected := "123abc"
	reader := bytes.NewReader([]byte(expected))
	lex := NewLexer(reader)
	if err := lex.fetchHexaNumber(); err != nil {
		t.Error(err.Error())
		return
	}

	tokens := lex.tokens

	if len(tokens) != 1 {
		t.Error("expecting 1 token to be fetched")
		return
	}

	token := tokens[0]
	if token.String() != expected {
		t.Errorf("unexpected '%s', expecting '%s'", token.String(), expected)
	}
}

// TestFetchHexaNumberWithSeparator test Lexer.fetchHexaNumber behavior with number separator
func TestFetchHexaNumberWithSeparator(t *testing.T) {
	expected := "1_2_3_a_b_c"
	reader := bytes.NewReader([]byte(expected))
	lex := NewLexer(reader)
	if err := lex.fetchHexaNumber(); err != nil {
		t.Error(err.Error())
		return
	}

	tokens := lex.tokens

	if len(tokens) != 1 {
		t.Error("expecting 1 token to be fetched")
		return
	}

	token := tokens[0]
	if token.String() != expected {
		t.Errorf("unexpected '%s', expecting '%s'", token.String(), expected)
	}
}

// TestFetchBinaryNumber test Lexer.fetchBinaryNumber() behavior
func TestFetchBinaryNumber(t *testing.T) {
	expected := "10101010"
	reader := bytes.NewReader([]byte(expected))
	lex := NewLexer(reader)
	if err := lex.fetchBinaryNumber(); err != nil {
		t.Error(err.Error())
		return
	}

	tokens := lex.tokens

	if len(tokens) != 1 {
		t.Error("expecting 1 token to be fetched")
		return
	}

	token := tokens[0]
	if token.String() != expected {
		t.Errorf("unexpected '%s', expecting '%s'", token.String(), expected)
	}
}

// TestFetchBinaryNumberWithSeparator test Lexer.fetchBinaryNumber() behavior with separator
func TestFetchBinaryNumberWithSeparator(t *testing.T) {
	expected := "1_01_0_10_1_0"
	reader := bytes.NewReader([]byte(expected))
	lex := NewLexer(reader)
	if err := lex.fetchBinaryNumber(); err != nil {
		t.Error(err.Error())
		return
	}

	tokens := lex.tokens

	if len(tokens) != 1 {
		t.Error("expecting 1 token to be fetched")
		return
	}

	token := tokens[0]
	if token.String() != expected {
		t.Errorf("unexpected '%s', expecting '%s'", token.String(), expected)
	}
}

// -------------------- Tests for Lexer.fetchOctalNumber() --------------------

// TestFetchOctalNumber test Lexer.fetchOctalNumber() behavior
func TestFetchOctalNumber(t *testing.T) {
	expected := "01234567"
	reader := bytes.NewReader([]byte(expected))
	lex := NewLexer(reader)
	if err := lex.fetchOctalNumber(); err != nil {
		t.Error(err.Error())
		return
	}

	tokens := lex.tokens

	if len(tokens) != 1 {
		t.Error("expecting 1 token to be fetched")
		return
	}

	token := tokens[0]
	if token.String() != expected {
		t.Errorf("unexpected '%s', expecting '%s'", token.String(), expected)
	}
}

// TestFetchOctalNumberInvalidDigit test Lexer.fetchOctalNumber() behavior
// with invalid digit
func TestFetchOctalNumberInvalidDigit(t *testing.T) {
	expected := "012345678"
	reader := bytes.NewReader([]byte(expected))
	lex := NewLexer(reader)
	if err := lex.fetchOctalNumber(); err != nil {
		t.Log("expected error", err.Error())
	} else {
		t.Errorf("expecting error, got result '%s' instead", lex.tokens[0].String())
	}
}

// TestFetchOctalNumberWithSeparator test Lexer.fetchOctalNumber() behavior with separator
func TestFetchOctalNumberWithSeparator(t *testing.T) {
	expected := "0_1_2_3_4_5_6_7"
	reader := bytes.NewReader([]byte(expected))
	lex := NewLexer(reader)
	if err := lex.fetchOctalNumber(); err != nil {
		t.Error(err.Error())
		return
	}

	tokens := lex.tokens

	if len(tokens) != 1 {
		t.Error("expecting 1 token to be fetched")
		return
	}

	token := tokens[0]
	if token.String() != expected {
		t.Errorf("unexpected '%s', expecting '%s'", token.String(), expected)
	}
}

// TestFetchOctalNumberWithIvalidSeparator1 test Lexer.fetchOctalNumber() behavior with invalid separator
// placement (at the begining of token)
func TestFetchOctalNumberWithInvalidSeparator1(t *testing.T) {
	expected := "_0_1_2_3_4_5_6_7"
	reader := bytes.NewReader([]byte(expected))
	lex := NewLexer(reader)
	if err := lex.fetchOctalNumber(); err != nil {
		t.Log("expected error", err.Error())
		return
	} else {
		token := lex.tokens[0]
		t.Errorf("expecting error, got '%s' instead", token.String())
	}
}

// TestFetchOctalNumberWithIvalidSeparator2 test Lexer.fetchOctalNumber() behavior with invalid separator
// placement (at the end of token)
func TestFetchOctalNumberWithInvalidSeparator2(t *testing.T) {
	expected := "0_1_2_3_4_5_6_7_"
	reader := bytes.NewReader([]byte(expected))
	lex := NewLexer(reader)
	if err := lex.fetchOctalNumber(); err != nil {
		t.Log("expected error", err.Error())
		return
	} else {
		token := lex.tokens[0]
		t.Errorf("expecting error, got '%s' instead", token.String())
	}
}

// TestFetchOctalNumberWithInvalidDigitAfterSeparator test lexer.fetchOctalNumber() behavior
// with invalid digit after separator
func TestFetchOctalNumberWithInvalidDigitAfterSeparator(t *testing.T) {
	expected := "0_1_2_3_4_5_6_8"
	reader := bytes.NewReader([]byte(expected))
	lex := NewLexer(reader)
	if err := lex.fetchOctalNumber(); err != nil {
		t.Log("expected error", err.Error())
		return
	} else {
		token := lex.tokens[0]
		t.Errorf("expecting error, got '%s' instead", token.String())
	}
}

// TestFetchOctalNumberEndByWhiteSpace test Lexer.fetchOctalNumber() behavior
// with whitespace as end of token
func TestFetchOctalNumberEndByWhiteSpace(t *testing.T) {
	expected := "01234567"
	input := "01234567 "
	reader := bytes.NewReader([]byte(input))
	lex := NewLexer(reader)
	if err := lex.fetchOctalNumber(); err != nil {
		t.Errorf("unexpected error '%s', expecting result '%s'", err.Error(), expected)
		return
	}

	token := lex.tokens[0]
	if token.String() != expected {
		t.Errorf("unexpected result '%s', expecting '%s'", token.String(), expected)
	}
}

// TestFetchOctalNumberEndByPunct test Lexer.fetchOctalNumber() behavior
// with punctuator as end of token
func TestFetchOctalNumberEndByPunct(t *testing.T) {
	expected := "01234567"
	input := "01234567}"
	reader := bytes.NewReader([]byte(input))
	lex := NewLexer(reader)
	if err := lex.fetchOctalNumber(); err != nil {
		t.Errorf("unexpected error '%s', expecting result '%s'", err.Error(), expected)
		return
	}

	token := lex.tokens[0]
	if token.String() != expected {
		t.Errorf("unexpected result '%s', expecting '%s'", token.String(), expected)
	}
}

// -------------------- Tests for Lexer.fetchInfinityNumber() --------------------

// TestFetchInfinityNumber test Lexer.fetchInfinityNumber() behavior
func TestFetchInfinityNumber(t *testing.T) {
	input := "nfinity"
	expected := "Infinity"
	reader := bytes.NewReader([]byte(input))
	lex := NewLexer(reader)
	if err := lex.fetchInfinityNumber(); err != nil {
		t.Error(err.Error())
		return
	}

	tokens := lex.tokens

	if len(tokens) != 1 {
		t.Error("expecting 1 token to be fetched")
		return
	}

	token := tokens[0]
	if token.String() != expected {
		t.Errorf("unexpected '%s', expecting '%s'", token.String(), expected)
	}
}

// TestFetchInfinityNumberIdentResult1 test Lexer.fetchInfinityNumber() behavior
// when the characters is not forming an Infinity token because wrong set of characters
func TestFetchInfinityNumberIdentResult1(t *testing.T) {
	input := "nfinitybutisident"
	expected := "Infinitybutisident"
	reader := bytes.NewReader([]byte(input))
	lex := NewLexer(reader)
	if err := lex.fetchInfinityNumber(); err != nil {
		t.Error(err.Error())
		return
	}

	tokens := lex.tokens

	if len(tokens) != 1 {
		t.Error("expecting 1 token to be fetched")
		return
	}

	token := tokens[0]
	if token.String() != expected {
		t.Errorf("unexpected '%s', expecting '%s'", token.String(), expected)
	}
}

// TestFetchInfinityNumberIdentResult2 test Lexer.fetchInfinityNumber() behavior
// when the characters is not forming an Infinity token because wrong set of characters
func TestFetchInfinityNumberIdentResult2(t *testing.T) {
	input := "nfiniti"
	expected := "Infiniti"
	reader := bytes.NewReader([]byte(input))
	lex := NewLexer(reader)
	if err := lex.fetchInfinityNumber(); err != nil {
		t.Error(err.Error())
		return
	}

	tokens := lex.tokens

	if len(tokens) != 1 {
		t.Error("expecting 1 token to be fetched")
		return
	}

	token := tokens[0]
	if token.String() != expected {
		t.Errorf("unexpected '%s', expecting '%s'", token.String(), expected)
	}
}

// TestFetchInfinityNumberEndByWhitespace test Lexer.fetchInfinityNumber() behavior
// end by whitespace
func TestFetchInfinityNumberEndByWhitespace(t *testing.T) {
	input := "nfinity "
	expected := "Infinity"
	reader := bytes.NewReader([]byte(input))
	lex := NewLexer(reader)
	if err := lex.fetchInfinityNumber(); err != nil {
		t.Error(err.Error())
		return
	}

	tokens := lex.tokens

	if len(tokens) != 1 {
		t.Error("expecting 1 token to be fetched")
		return
	}

	token := tokens[0]
	if token.String() != expected {
		t.Errorf("unexpected '%s', expecting '%s'", token.String(), expected)
	}
}

// TestFetchInfinityNumberEndByPunct test Lexer.fetchInfinityNumber() behavior
// end by punctuator
func TestFetchInfinityNumberEndByPunct(t *testing.T) {
	input := "nfinity}"
	expected := "Infinity"
	reader := bytes.NewReader([]byte(input))
	lex := NewLexer(reader)
	if err := lex.fetchInfinityNumber(); err != nil {
		t.Error(err.Error())
		return
	}

	tokens := lex.tokens

	token := tokens[0]
	if token.String() != expected {
		t.Errorf("unexpected '%s', expecting '%s'", token.String(), expected)
	}
}

// -------------------- Tests for Lexer.fetchExponentNumber() --------------------

// TestFetchExponentNumber test Lexer.fetchExponentNumber() behavior
func TestFetchExponentNumberWithSign(t *testing.T) {
	expects := []string{"123", "-123", "+123"}
	for _, expected := range expects {
		reader := bytes.NewReader([]byte(expected))
		lex := NewLexer(reader)
		if err := lex.fetchExponentNumber(); err != nil {
			t.Error(err.Error())
			return
		}

		token := lex.tokens[0]

		if token.String() != expected {
			t.Errorf("unexpected '%s', expecting '%s'", token.String(), expected)
			return
		}
	}
}

// TestFetchExponentNumberWithInvalidChar test Lexer.fetchExponentNumber() behavior
// when invalid character present
func TestFetchExponentNumberWithInvalidChar(t *testing.T) {
	inputs := []string{"123a", "a123"}
	for _, input := range inputs {
		reader := bytes.NewReader([]byte(input))
		lex := NewLexer(reader)
		if err := lex.fetchExponentNumber(); err == nil {
			token := lex.tokens[0]
			t.Errorf("unexpected result '%s', expecting invalid char error", token.String())
			return
		}
	}
}

// TestFetchExponentNumberWithSeparator test Lexer.fetchExponentNumber() behavior
// when separator present
func TestFetchExponentNumberWithSeparator(t *testing.T) {
	expected := "1_2_3"
	reader := bytes.NewReader([]byte(expected))
	lex := NewLexer(reader)
	if err := lex.fetchExponentNumber(); err != nil {
		t.Errorf(err.Error())
		return
	}

	token := lex.tokens[0]

	if token.String() != expected {
		t.Errorf("unexpected '%s', expecting '%s'", token.String(), expected)
	}
}

// TestFetchExponentNumberWithInvalidSeparator test Lexer.fetchExponentNumber() behavior
// when separator present but with invalid placement
func TestFetchExponentNumberWithInvalidSeparator(t *testing.T) {
	inputs := []string{"_123", "123_", "1__23"}
	for _, input := range inputs {
		reader := bytes.NewReader([]byte(input))
		lex := NewLexer(reader)
		if err := lex.fetchExponentNumber(); err == nil {
			token := lex.tokens[0]
			t.Errorf("unexpected result '%s', expecting invalid char error", token.String())
			return
		}
	}
}

// TestFetchExponentNumberEndByNonEOF test Lexer.fetchExponentNumber() behavior
// when token end by whitespace and punctuation, not an EOF
func TestFetchExponentNumberEndByNonEOF(t *testing.T) {
	inputs := []string{"123 ", "123{", "123}", "123[", "123]", "123:", "123,"}
	expected := "123"
	for _, input := range inputs {
		reader := bytes.NewReader([]byte(input))
		lex := NewLexer(reader)
		if err := lex.fetchExponentNumber(); err != nil {
			t.Errorf(err.Error())
			return
		}

		token := lex.tokens[0]

		if token.String() != expected {
			t.Errorf("unexpected '%s', expecting '%s'", token.String(), expected)
			return
		}
	}
}

// -------------------- Tests for Lexer.fetchNanNumber() --------------------

// TestFetchNanNumber test Lexer.fetchNanNumber() behavior
func TestFetchNanNumber(t *testing.T) {
	expected := "NaN"
	input := []byte("aN")
	lex := NewLexer(bytes.NewReader(input))
	err := lex.fetchNanNumber()
	if err != nil {
		t.Error(err.Error())
		return
	}

	token := lex.tokens[0]
	if token.String() != expected {
		t.Errorf("unexpected '%s', expecting '%s'", token.String(), expected)
	}
}

// TestFetchNanNumberEndByNonEOF test Lexer.fetchNanNumber() behavior
// when token end by not an EOF
func TestFetchNanNumberEndByNonEOF(t *testing.T) {
	inputs := [][]byte{[]byte("aN "), []byte("aN{"), []byte("aN}"), []byte("aN["), []byte("aN]"), []byte("aN:"), []byte("aN,")}
	expected := "NaN"

	for _, input := range inputs {
		lex := NewLexer(bytes.NewReader(input))
		err := lex.fetchNanNumber()
		if err != nil {
			t.Error(err.Error())
			return
		}

		token := lex.tokens[0]
		if token.String() != expected {
			t.Errorf("unexpected '%s', expecting '%s'", token.String(), expected)
			return
		}
	}
}

// TestFetchNanNumberButIdent test Lexer.fetchNanNumber()
// when token is not forming a NaN but instead and identifier
func TestFetchNanNumberButIdent(t *testing.T) {
	inputs := [][]byte{[]byte("aNButIsIdentifier"), []byte("ahThisIsNotNaN")}
	expects := []string{"NaNButIsIdentifier", "NahThisIsNotNaN"}
	for i, input := range inputs {
		lex := NewLexer(bytes.NewReader(input))
		if err := lex.fetchNanNumber(); err != nil {
			t.Error(err.Error())
			return
		}

		token := lex.tokens[0]
		if token.Type() != TokenIdentifier {
			t.Errorf("unexpected token type '%s', expecting type '%s'", token.TypeString(), tokenTypeMap[TokenIdentifier])
			return
		}

		expected := expects[i]
		if token.String() != expected {
			t.Errorf("unexpected '%s', expecting '%s'", token.String(), expected)
			return
		}
	}
}

// -------------------- Tests for Lexer.fetchDoubleNumber() --------------------

// TestFetchDoubleNumber test Lexer.fetchDoubleNumber() behavior
func TestFetchDoubleNumber(t *testing.T) {
	input := []byte("123")
	expected := ".123"
	lex := NewLexer(bytes.NewReader(input))
	if err := lex.fetchDoubleNumber(); err != nil {
		t.Error(err.Error())
		return
	}

	token := lex.tokens[0]
	if token.String() != expected {
		t.Errorf("unexpected '%s', expecting '%s'", token.String(), expected)
	}
}

// TestFetchDoubleNumberInvalid test Lexer.fetchDoubleNumber() behavior
// when there's invalid character
func TestFetchDoubleNumberInvalid(t *testing.T) {
	input := []byte("12a")
	lex := NewLexer(bytes.NewReader(input))
	if err := lex.fetchDoubleNumber(); err == nil {
		token := lex.tokens[0]
		t.Errorf("unexpected result '%s', expecting error invalid character", token.String())
		return
	}
}

// TestFetchDoubleNumberWithSeparator test Lexer.fetchDoubleNumber() behavior
// when using separator
func TestFetchDoubleNumberWithSeparator(t *testing.T) {
	input := []byte("1_2_3")
	expected := ".1_2_3"
	lex := NewLexer(bytes.NewReader(input))
	if err := lex.fetchDoubleNumber(); err != nil {
		t.Error(err.Error())
		return
	}

	token := lex.tokens[0]
	if token.String() != expected {
		t.Errorf("unexpected result '%s', expecting '%s'", token.String(), expected)
	}
}

// TestFetchDoubleNumberWithInvalidSeparator test Lexer.fetchDoubleNumber() behavior
// when using invalid separator placement
func TestFetchDoubleNumberWithInvalidSeparator(t *testing.T) {
	inputs := [][]byte{[]byte("_123"), []byte("123_")}
	for _, input := range inputs {
		lex := NewLexer(bytes.NewReader(input))
		if err := lex.fetchDoubleNumber(); err == nil {
			token := lex.tokens[0]
			t.Errorf("unexpected result '%s', expecting error invalid character", token.String())
			return
		}
	}
}

func TestFetchDoubleNumberWithExponent(t *testing.T) {
	inputs := []string{"1e123", "1E123", "1e-123", "1e+123", "1E-123", "1E+123"}
	expects := []string{".1e123", ".1E123", ".1e-123", ".1e+123", ".1E-123", ".1E+123"}
	for i, input := range inputs {
		lex := NewLexer(bytes.NewReader([]byte(input)))
		if err := lex.fetchDoubleNumber(); err != nil {
			t.Error(err.Error())
			return
		}

		token := lex.tokens[0]
		expected := expects[i]
		if token.String() != expected {
			t.Errorf("unexpected '%s', expecting '%s'", token.String(), expected)
			return
		}
	}
}

// -------------------- Tests for Lexer.fetchNumber() --------------------

// TestFetchNumber test Lexer.fetchNumber() behavior
func TestFetchNumber(t *testing.T) {
	inputBegins := []rune{
		'1',
		'-',
		'+',
		'-',
		'.',
		'-',
		'+',
		'-',
		'0',
		'-',
		'+',
		'-',
		'1',
		'1',
		'0',
		'0',
		'0',
		'0',
		'0',
		'0',
		'0',
		'-',
		'+',
		'I',
		'-',
		'+',
		'N',
	}

	inputs := []string{
		"23",
		"123",
		"123",
		"-123",
		"123",
		".123",
		".123",
		"-.123",
		".123",
		"0.123",
		"0.123",
		"-0.123",
		".e123",
		".E123",
		"x123",
		"X123",
		"b1010",
		"B1010",
		"o123",
		"O123",
		"123",
		"Infinity",
		"Infinity",
		"nfinity",
		"NaN",
		"NaN",
		"aN",
	}

	expects := []string{
		"123",
		"-123",
		"+123",
		"--123",
		".123",
		"-.123",
		"+.123",
		"--.123",
		"0.123",
		"-0.123",
		"+0.123",
		"--0.123",
		"1.e123",
		"1.E123",
		"0x123",
		"0X123",
		"0b1010",
		"0B1010",
		"0o123",
		"0O123",
		"0123",
		"-Infinity",
		"+Infinity",
		"Infinity",
		"-NaN",
		"+NaN",
		"NaN",
	}

	for i, input := range inputs {
		lex := NewLexer(bytes.NewReader([]byte(input)))
		if err := lex.fetchNumber(inputBegins[i]); err != nil {
			t.Error(err.Error())
			return
		}

		expected := expects[i]
		token := lex.tokens[0]
		if token.String() != expected {
			t.Errorf("unexpected '%s', expecting '%s'", token.String(), expected)
			return
		}
	}
}

// -------------------- Tests for Lexer.fetchComment() --------------------

// TestFetchComment test Lexer.fetchComment() behavior
func TestFetchComment(t *testing.T) {
	inputs := []string{
		"/ inline comment",
		`*
			multiline comment
		*/`,
	}

	expects := []string{
		"// inline comment",
		`/*
			multiline comment
		*/`,
	}

	for i, input := range inputs {
		lex := NewLexer(bytes.NewReader([]byte(input)))
		if err := lex.fetchComment(); err != nil {
			t.Error(err.Error())
			return
		}

		token := lex.tokens[0]

		if token.Type() != TokenComment {
			t.Errorf("unexpected token type '%s', expecting type '%s'", token.TypeString(), tokenTypeMap[TokenComment])
			return
		}

		expected := expects[i]
		if token.String() != expected {
			t.Errorf("unexpected '%s', expecting '%s'", token.String(), expected)
			return
		}
	}
}

// TestFetchCommentInvalid test Lexer.fetchComment() behavior
// when trying to fetch invalid comment token
func TestFetchCommentInvalid(t *testing.T) {
	inputs := []string{
		"invalid comment",
		`*
			invalid multiline comment
		*
		`,
	}

	for _, input := range inputs {
		lex := NewLexer(bytes.NewReader([]byte(input)))
		if err := lex.fetchComment(); err == nil {
			token := lex.tokens[0]
			t.Errorf("unexpected result '%s', expecting error", token.String())
			return
		}
	}
}

func TestTokens(t *testing.T) {
	input := `
	{		
		nullVal: null, // comment 
		undef: undefined,
		trueBool: true,
		falseBool: false,
		array: [
			null,
			undefined,
			true,
			false,
			"double-quote string",
			'single-quote string',
		],
		"doubleStrIden": "\u1234 \u{12345678} \xFf",
		'singleStrIden': 'Howdy',
		/*
			multiline comment
		*/
		arrayOfNumbers: [
			123,
			-123,
			+123,
			--123,
			.123,
			-.123,
			+.123,
			--.123,
			0.123,
			-0.123,
			+0.123,
			--0.123,			
			1.e123,
			1.E123,
			0x123,
			0X123,
			0b1010,
			0B1010,
			0o123,
			0O123,
			0123,
			-Infinity,
			+Infinity,
			Infinity,
			-NaN,
			+NaN,
			NaN,
		],
	}
`
	reader := bytes.NewReader([]byte(input))
	lex := NewLexer(reader)
	err := lex.FetchTokens()
	if err != nil {
		t.Error(err.Error())
		return
	}

	for {
		token, err := lex.ReadToken()
		if err != nil {
			if err == ErrNoMoreToken {
				break
			}

			t.Error(err.Error())
			return
		}

		fmt.Println("type:", token.TypeString())
		fmt.Println("string:", token.String())
		fmt.Printf("position: %d:%d \n", token.StartPos.Line(), token.StartPos.Column())
		fmt.Print("\n")
	}
}
