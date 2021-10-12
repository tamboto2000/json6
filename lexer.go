package json6

import (
	"io"
	"unicode"
)

// TokenType is token type
type TokenType byte

// Token types
const (
	TokenIdentifier TokenType = iota
	TokenPunctuator           // '{', '}', '[', ']', ':', ','
	TokenString
	TokenNumber
	TokenNull
	TokenBool
	TokenUndefined
	TokenComment
)

// sub-type for TokenNumber
const (
	tokenNumInteger = iota
	tokenNumDouble
)

// Token types in string
var tokenTypeMap = map[TokenType]string{
	TokenIdentifier: "identifier",
	TokenString:     "string",
	TokenNumber:     "number",
	TokenNull:       "null",
	TokenBool:       "boolean",
	TokenUndefined:  "undefined",
	TokenPunctuator: "punctuator",
	TokenComment:    "comment",
}

// Token contain characters that form the token, position in file, and its type
type Token struct {
	Pos             *Position
	t               TokenType
	tokenNumSubType uint
	chars           []rune
}

// String return token string
func (t Token) String() string {
	return string(t.chars)
}

// Type return token type in TokenType
func (t Token) Type() TokenType {
	return t.t
}

// Type return token type name (string)
func (t Token) TypeString() string {
	return tokenTypeMap[t.t]
}

// lastChar fetch last pushed character in Token.chars
func (t Token) lastChar() rune {
	if len(t.chars) > 0 {
		return t.chars[len(t.chars)-1]
	}

	return 0
}

// Position indicating token's position
type Position struct {
	ln  int
	col int
}

func newPosition(ln, col int) *Position {
	return &Position{ln: ln, col: col}
}

// Line of the position of a token
func (pos *Position) Line() int {
	return pos.ln
}

// Column of the position of a token
func (pos *Position) Column() int {
	return pos.col
}

func (pos *Position) addLn(add int) {
	pos.ln += add
}

func (pos *Position) setCol(col int) {
	pos.col = col
}

func (pos *Position) addCol(add int) {
	pos.col += add
}

// Lexer fetch JSON6 tokens
type Lexer struct {
	pos       *Position
	r         io.RuneReader
	tokens    []Token
	token     Token // current token
	ignoreErr bool  // set to true to ignore lexical error
}

func NewLexer(r io.RuneReader) *Lexer {
	pos := newPosition(1, 0)
	r = newReader(r, pos)
	return &Lexer{
		pos: pos,
		r:   r,
	}
}

func (lx *Lexer) push() {
	lx.token.Pos = newPosition(lx.pos.Line(), lx.pos.Column())
	lx.tokens = append(lx.tokens, lx.token)
	lx.token = Token{}
}

func (lx *Lexer) pushWithPos(ln, cl int) {
	lx.token.Pos = newPosition(ln, cl)
	lx.tokens = append(lx.tokens, lx.token)
	lx.token = Token{}
}

// IgnoreError determine if lexer will be ignoring lexical error or not,
// default behavior is to not allow lexical error.
// Call IgnoreError(true) to ignore lexical error
func (lx *Lexer) IgnoreError(ignore bool) {
	lx.ignoreErr = ignore
}

// Tokens return fetched tokens
func (lx *Lexer) Tokens() ([]Token, error) {
	for {
		char, _, err := lx.r.ReadRune()
		if err != nil {
			if err == io.EOF {
				break
			}

			return nil, err
		}

		switch char {
		// comment
		case '/':
			if err := lx.fetchComment(); err != nil {
				if lx.ignoreErr {
					lx.token = Token{}
					continue
				}

				return nil, err
			}

			continue

		// true boolean
		case 't':
			if err := lx.fetchTrueBool(); err != nil {
				if lx.ignoreErr {
					lx.token = Token{}
					continue
				}

				return nil, err
			}

			continue

		// false boolean
		case 'f':
			if err := lx.fetchFalseBool(); err != nil {
				if lx.ignoreErr {
					lx.token = Token{}
					continue
				}

				return nil, err
			}

			continue

		// null
		case 'n':
			if err := lx.fetchNull(); err != nil {
				if lx.ignoreErr {
					lx.token = Token{}
					continue
				}

				return nil, err
			}

			continue

		// undefined
		case 'u':
			if err := lx.fetchUndefined(); err != nil {
				if lx.ignoreErr {
					lx.token = Token{}
					continue
				}

				return nil, err
			}

			continue

		// punctuator
		case '{', '}', '[', ']', ':', ',':
			lx.fetchPunct(char)
			continue

		// string
		case '"', '\'', '`':
			if err := lx.fetchString(char); err != nil {
				if lx.ignoreErr {
					lx.token = Token{}
					continue
				}

				return nil, err
			}

			continue

		// number
		case '-', '+', '.', 'I', 'N':
			if err := lx.fetchNumber(char); err != nil {
				if lx.ignoreErr {
					lx.token = Token{}
					continue
				}

				return nil, err
			}

			continue

		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			if err := lx.fetchNumber(char); err != nil {
				if lx.ignoreErr {
					lx.token = Token{}
					continue
				}

				return nil, err
			}

			continue

		default:
			// Check if char is whitespace
			if isCharWhitespace(char) {
				continue
			}

			// if char is not whitespace, try to fetch identifier token
			if err := lx.fetchIdentifier(true, char); err != nil {
				if lx.ignoreErr {
					lx.token = Token{}
					continue
				}

				return nil, err
			}
		}
	}

	return lx.tokens, nil
}

func (lx *Lexer) fetchComment() error {
	lx.token.chars = append(lx.token.chars, '/')
	lx.token.t = TokenComment

	char, _, err := lx.r.ReadRune()
	if err != nil {
		if err == io.EOF {
			return errUnexpectedEOF(lx.pos, "'/' or '*'")
		}

		return err
	}

	if char == '/' {
		lx.token.chars = append(lx.token.chars, char)
		for {
			char, _, err := lx.r.ReadRune()
			if err != nil {
				if err == io.EOF {
					lx.push()
					return nil
				}

				return err
			}

			switch char {
			case '\r', '\n', '\u2028', '\u2029':
				lx.pushWithPos(lx.pos.ln, lx.pos.col-1)
				return nil

			default:
				lx.token.chars = append(lx.token.chars, char)
			}
		}
	} else if char == '*' {
		lx.token.chars = append(lx.token.chars, char)
		for {
			char, _, err := lx.r.ReadRune()
			if err != nil {
				if err == io.EOF {
					return errUnexpectedEOF(lx.pos, "'*'")
				}

				return err
			}

			lx.token.chars = append(lx.token.chars, char)
			if char == '*' {
				char, _, err := lx.r.ReadRune()
				if err != nil {
					if err == io.EOF {
						return errUnexpectedEOF(lx.pos, "'/'")
					}

					return err
				}

				lx.token.chars = append(lx.token.chars, char)
				if char == '/' {
					lx.push()
					return nil
				}
			}
		}
	}

	lx.token.chars = append(lx.token.chars, char)
	return errInvalidChar(char, lx.pos, lx.token.chars, "'/'")
}

var falseBoolChars = []rune{'a', 'l', 's', 'e'}

// fetchFalseBool fetch 'false' boolean
func (lx *Lexer) fetchFalseBool() error {
	lx.token.t = TokenBool
	lx.token.chars = append(lx.token.chars, 'f')
	for _, c := range falseBoolChars {
		char, _, err := lx.r.ReadRune()
		if err != nil {
			if err == io.EOF {
				if len(lx.token.chars) > 0 {
					lx.token.t = TokenIdentifier
					lx.push()
				}

				return nil
			}

			return err
		}

		if char != c {
			return lx.fetchIdentifier(false, char)
		}

		lx.token.chars = append(lx.token.chars, char)
	}

	char, _, err := lx.r.ReadRune()
	if err != nil {
		if err != io.EOF {
			return err
		}

		lx.push()
		return nil
	}

	if !isCharWhitespace(char) {
		if isCharPunct(char) {
			defer lx.fetchPunct(char)
		} else {
			return lx.fetchIdentifier(false, char)
		}
	}

	lx.push()

	return nil
}

var trueBoolChars = []rune{'r', 'u', 'e'}

// fetchTrueBool fetch 'true' boolean
func (lx *Lexer) fetchTrueBool() error {
	lx.token.t = TokenBool
	lx.token.chars = append(lx.token.chars, 't')
	for _, c := range trueBoolChars {
		char, _, err := lx.r.ReadRune()
		if err != nil {
			if err == io.EOF {
				if len(lx.token.chars) > 0 {
					lx.token.t = TokenIdentifier
					lx.push()
				}

				return nil
			}

			return err
		}

		if char != c {
			return lx.fetchIdentifier(false, char)
		}

		lx.token.chars = append(lx.token.chars, char)
	}

	char, _, err := lx.r.ReadRune()
	if err != nil {
		if err != io.EOF {
			return err
		}

		lx.push()
		return nil
	}

	if !isCharWhitespace(char) {
		if isCharPunct(char) {
			defer lx.fetchPunct(char)
			lx.pushWithPos(lx.pos.ln, lx.pos.col-1)
			return nil
		} else {
			return lx.fetchIdentifier(false, char)
		}
	}

	lx.push()

	return nil
}

var nullChars = []rune{'u', 'l', 'l'}

// fetchNull fetch null token
func (lx *Lexer) fetchNull() error {
	lx.token.t = TokenNull
	lx.token.chars = append(lx.token.chars, 'n')
	for _, c := range nullChars {
		char, _, err := lx.r.ReadRune()
		if err != nil {
			if err == io.EOF {
				if len(lx.token.chars) > 0 {
					lx.token.t = TokenIdentifier
					lx.push()
				}

				return nil
			}

			return err
		}

		// if not null value, it's probably identifier
		if char != c {
			return lx.fetchIdentifier(false, char)
		}

		lx.token.chars = append(lx.token.chars, char)
	}

	// the next char will determine if this token is really is boolean or identifier
	char, _, err := lx.r.ReadRune()
	if err != nil {
		if err != io.EOF {
			return err
		}

		lx.push()
		return nil
	}

	if !isCharWhitespace(char) {
		if isCharPunct(char) {
			defer lx.fetchPunct(char)
			lx.pushWithPos(lx.pos.ln, lx.pos.col-1)
			return nil
		} else {
			return lx.fetchIdentifier(false, char)
		}
	}

	lx.push()

	return nil
}

// fetchIdentifier fetch identifier token
func (lx *Lexer) fetchIdentifier(isBegin bool, firstChar rune) error {
	lx.token.t = TokenIdentifier

	// if isBegin is true, check if firstChar is valid identifier start and
	// append to token chars if firstChar is valid
	switch firstChar {
	case '$', '_':
		lx.token.chars = append(lx.token.chars, firstChar)

	// punctuator
	case '{', '}', '[', ']', ':', ',':
		defer lx.fetchPunct(firstChar)
		lx.pushWithPos(lx.pos.ln, lx.pos.col-1)
		return nil

	case '\\':
		// if char is begin of escape sequence, check if escape sequence is unicode escape sequence
		char, _, err := lx.r.ReadRune()
		if err != nil {
			if err == io.EOF {
				return errUnexpectedEOF(lx.pos, "'u'")
			}

			return err
		}

		lx.token.chars = append(lx.token.chars, char)

		if char != 'u' {
			lx.token.chars = append(lx.token.chars, char)
			return errInvalidChar(char, lx.pos, lx.token.chars, "'u'")
		}

		if err := lx.fetchUnicodeEscape(); err != nil {
			return err
		}

	default:
		if !unicode.In(firstChar, unicode.Lu, unicode.Ll, unicode.Lt, unicode.Lm, unicode.Lo, unicode.Nl) {
			if isBegin {
				lx.token.chars = append(lx.token.chars, firstChar)
				return errInvalidChar(firstChar, lx.pos, lx.token.chars, "'$', '_', unicode escape sequence, or any charater in categories Uppercase letter (Lu), Lowercase letter (Ll), Titlecase letter (Lt), Modifier letter (Lm), Other letter (Lo), Letter number (Nl)")
			}

			if !unicode.In(firstChar, unicode.Mn, unicode.Mc, unicode.Nd, unicode.Pc) {
				if !isCharWhitespace(firstChar) {
					lx.token.chars = append(lx.token.chars, firstChar)
					return errInvalidChar(firstChar, lx.pos, lx.token.chars, "'$', '_', unicode escape sequence, or any charater in categories Uppercase letter (Lu), Lowercase letter (Ll), Titlecase letter (Lt), Modifier letter (Lm), Other letter (Lo), Letter number (Nl), Non-spacing mark (Mn), Combining spacing mark (Mc), Decimal number (Nd), Connector punctuation (Pc)")
				}
			}
		}

		lx.token.chars = append(lx.token.chars, firstChar)
	}

LOOP:
	for {
		char, _, err := lx.r.ReadRune()
		if err != nil {
			if err == io.EOF {
				lx.push()
				return nil
			}

			return err
		}

		switch char {
		case '$', '_':
			lx.token.chars = append(lx.token.chars, char)
			continue

		// punctuator
		case '{', '}', '[', ']', ':', ',':
			defer lx.fetchPunct(char)
			lx.pushWithPos(lx.pos.ln, lx.pos.col-1)
			return nil

		case '\\':
			lx.token.chars = append(lx.token.chars, char)
			// if char is begin of escape sequence, check if escape sequence is unicode escape sequence
			char, _, err := lx.r.ReadRune()
			if err != nil {
				if err == io.EOF {
					return errUnexpectedEOF(lx.pos, "'u'")
				}

				return err
			}

			if char != 'u' {
				lx.token.chars = append(lx.token.chars, char)
				return errInvalidChar(char, lx.pos, lx.token.chars, "'u'")
			}

			lx.token.chars = append(lx.token.chars, char)

			if err := lx.fetchUnicodeEscape(); err != nil {
				return err
			}

		default:
			if !unicode.In(char, unicode.Lu, unicode.Ll, unicode.Lt, unicode.Lm, unicode.Lo, unicode.Nl, unicode.Mn, unicode.Mc, unicode.Nd, unicode.Pc) {
				if isCharWhitespace(char) {
					break LOOP
				}

				lx.token.chars = append(lx.token.chars, char)
				return errInvalidChar(firstChar, lx.pos, lx.token.chars, "'$', '_', unicode escape sequence, or any charater in categories Uppercase letter (Lu), Lowercase letter (Ll), Titlecase letter (Lt), Modifier letter (Lm), Other letter (Lo), Letter number (Nl), Non-spacing mark (Mn), Combining spacing mark (Mc), Decimal number (Nd), Connector punctuation (Pc)")
			}

			lx.token.chars = append(lx.token.chars, char)
		}
	}

	lx.pushWithPos(lx.pos.ln, lx.pos.col-1)

	return nil
}

// fetchHexaEscape fetch hexadecimal escape sequence, example:
// \xff
func (lx *Lexer) fetchHexaEscape() error {
	for i := 0; i < 2; i++ {
		char, _, err := lx.r.ReadRune()
		if err != nil {
			if err == io.EOF {
				return errInvalidChar(char, lx.pos, lx.token.chars, "hexadecimal digit")
			}

			return err
		}

		lx.token.chars = append(lx.token.chars, char)

		if !isCharValidHexa(char) {
			return errInvalidChar(char, lx.pos, lx.token.chars, "hexadecimal digit")
		}
	}

	return nil
}

// fetchUnicodeEscape fetch unicode escape sequence
func (lx *Lexer) fetchUnicodeEscape() error {
	// In ECMAScript 6, there's 2 (two)	types of unicode escape sequence: the good ol' 4 digit hexa digit unicode
	// and higher or fewer digit with {} (example: \u{12344f}, \u{f}) to contain them, so we must check
	// the first char immediately after char 'u'
	char, _, err := lx.r.ReadRune()
	if err != nil {
		return err
	}

	lx.token.chars = append(lx.token.chars, char)
	if char == '{' {
		for {
			char, _, err := lx.r.ReadRune()
			if err != nil {
				if err == io.EOF {
					return errUnexpectedEOF(lx.pos, "hexadecimal digit or '}'")
				}

				return err
			}

			lx.token.chars = append(lx.token.chars, char)
			if !isCharValidHexa(char) {
				if char == '}' {
					return nil
				}

				return errInvalidChar(char, lx.pos, lx.token.chars, "hexadecimal digit or '}'")
			}
		}
	}

	if !isCharValidHexa(char) {
		return errInvalidChar(char, lx.pos, lx.token.chars, "'{' or hexadecimal digit")
	}

	for i := 0; i < 3; i++ {
		char, _, err := lx.r.ReadRune()
		if err != nil {
			if err == io.EOF {
				return errUnexpectedEOF(lx.pos, "hexadecimal digit")
			}

			return err
		}

		if !isCharValidHexa(char) {
			return errInvalidChar(char, lx.pos, lx.token.chars, "hexadecimal digit")
		}

		lx.token.chars = append(lx.token.chars, char)
	}

	return nil
}

// fetchPunct is not exactly for fetching, more like creating the token
func (lx *Lexer) fetchPunct(char rune) {
	lx.token.t = TokenPunctuator
	lx.token.chars = append(lx.token.chars, char)
	lx.push()
}

var undefinedChars = []rune{'n', 'd', 'e', 'f', 'i', 'n', 'e', 'd'}

// fetchUndefined fetch undefined token
func (lx *Lexer) fetchUndefined() error {
	lx.token.t = TokenUndefined
	lx.token.chars = append(lx.token.chars, 'u')

	for _, c := range undefinedChars {
		char, _, err := lx.r.ReadRune()
		if err != nil {
			if err == io.EOF {
				if len(lx.token.chars) > 0 {
					lx.token.t = TokenIdentifier
					lx.push()
				}

				return nil
			}

			return err
		}

		if char != c {
			return lx.fetchIdentifier(false, char)
		}

		lx.token.chars = append(lx.token.chars, char)
	}

	char, _, err := lx.r.ReadRune()
	if err != nil {
		if err != io.EOF {
			return err
		}

		lx.push()
		return nil
	}

	if !isCharWhitespace(char) {
		if isCharPunct(char) {
			defer lx.fetchPunct(char)
			lx.pushWithPos(lx.pos.ln, lx.pos.col-1)
			return nil
		} else {
			return lx.fetchIdentifier(false, char)
		}
	}

	lx.push()

	return nil
}

func (lx *Lexer) fetchString(firstChar rune) error {
	lx.token.t = TokenString
	lx.token.chars = append(lx.token.chars, firstChar)

LOOP:
	for {
		char, _, err := lx.r.ReadRune()
		if err != nil {
			if err == io.EOF {
				return errUnexpectedEOF(lx.pos, "any Unicode code point")
			}

			return err
		}

		lx.token.chars = append(lx.token.chars, char)
		switch char {
		// possible unicode escape or hexa escape
		case '\\':
			char, _, err := lx.r.ReadRune()
			if err != nil {
				if err == io.EOF {
					return errUnexpectedEOF(lx.pos, "any Unicode point")
				}

				return err
			}

			lx.token.chars = append(lx.token.chars, char)

			switch char {
			// unicode escape
			case 'u':
				if err := lx.fetchUnicodeEscape(); err != nil {
					return err
				}

				continue

			// hexa escape
			case 'x':
				if err := lx.fetchHexaEscape(); err != nil {
					return err
				}

				continue

			default:
				continue
			}

		case '"':
			if firstChar == char {
				break LOOP
			}

		case '\'':
			if firstChar == char {
				break LOOP
			}

		case '`':
			if firstChar == char {
				break LOOP
			}
		}
	}

	lx.push()
	return nil
}

func (lx *Lexer) fetchHexaNumber() error {
	isFirstChar := true
	for {
		char, _, err := lx.r.ReadRune()
		if err != nil {
			if err == io.EOF {
				if isFirstChar {
					return errUnexpectedEOF(lx.pos, "hexadecimal digit")
				}

				lx.push()
				return nil
			}

			return err
		}

		if !isCharValidHexa(char) {
			if isFirstChar {
				lx.token.chars = append(lx.token.chars, char)
				return errInvalidChar(char, lx.pos, lx.token.chars, "hexadecimal digit")
			}

			if isCharPunct(char) {
				defer lx.fetchPunct(char)
				lx.pushWithPos(lx.pos.ln, lx.pos.col-1)

				return nil
			} else if isCharWhitespace(char) {
				lx.pushWithPos(lx.pos.ln, lx.pos.col-1)
				return nil
			} else if char == '_' {
				lx.token.chars = append(lx.token.chars, char)
				// check if next character is valid hexadecimal digit
				char, _, err := lx.r.ReadRune()
				if err != nil {
					if err == io.EOF {
						return errUnexpectedEOF(lx.pos, "hexadecimal digit")
					}

					return err
				}

				lx.token.chars = append(lx.token.chars, char)
				if !isCharValidHexa(char) {
					return errInvalidChar(char, lx.pos, lx.token.chars, "hexadecimal digit")
				}

				continue
			}

			lx.token.chars = append(lx.token.chars, char)
			return errInvalidChar(char, lx.pos, lx.token.chars, "hexadecimal digit")
		}

		isFirstChar = false
		lx.token.chars = append(lx.token.chars, char)
	}
}

func (lx *Lexer) fetchBinaryNumber() error {
	isFirstChar := true
	for {
		char, _, err := lx.r.ReadRune()
		if err != nil {
			if err == io.EOF {
				if isFirstChar {
					return errUnexpectedEOF(lx.pos, "binary digit")
				}

				lx.push()
				return nil
			}

			return err
		}

		if char != '0' && char != '1' {
			if isFirstChar {
				lx.token.chars = append(lx.token.chars, char)
				return errInvalidChar(char, lx.pos, lx.token.chars, "binary digit")
			}

			if char == '_' {
				lx.token.chars = append(lx.token.chars, char)
				char, _, err := lx.r.ReadRune()
				if err != nil {
					if err == io.EOF {
						return errUnexpectedEOF(lx.pos, "binary digit")
					}

					return err
				}

				lx.token.chars = append(lx.token.chars, char)
				if char != '0' && char != '1' {
					return errInvalidChar(char, lx.pos, lx.token.chars, "binary digit")
				}

				continue
			} else if isCharPunct(char) {
				defer lx.fetchPunct(char)
				lx.pushWithPos(lx.pos.ln, lx.pos.col-1)

				return nil
			} else if isCharWhitespace(char) {
				lx.pushWithPos(lx.pos.ln, lx.pos.col-1)

				return nil
			}

			lx.token.chars = append(lx.token.chars, char)
			return errInvalidChar(char, lx.pos, lx.token.chars, "binary digit or '_'")
		}

		isFirstChar = false
		lx.token.chars = append(lx.token.chars, char)
	}
}

var infinityChars = []rune{'n', 'f', 'i', 'n', 'i', 't', 'y'}

// fetchInfinity fetch number token with Infinity as value
func (lx *Lexer) fetchInfinityNumber() error {
	lx.token.chars = append(lx.token.chars, 'I')
	for _, c := range infinityChars {
		char, _, err := lx.r.ReadRune()
		if err != nil {
			if err == io.EOF {
				return errUnexpectedEOF(lx.pos, string([]rune{c}))
			}

			return err
		}

		if char != c {
			return lx.fetchIdentifier(false, char)
		}

		lx.token.chars = append(lx.token.chars, char)
	}

	char, _, err := lx.r.ReadRune()
	if err != nil {
		if err == io.EOF {
			lx.push()
			return nil
		}

		return err
	}

	if isCharPunct(char) {
		defer lx.fetchPunct(char)
		lx.pushWithPos(lx.pos.ln, lx.pos.col-1)
		return nil
	} else if isCharWhitespace(char) {
		lx.pushWithPos(lx.pos.ln, lx.pos.col-1)
		return nil
	}

	// if its turn out to be not a number, the nearest possibility is that this token
	// might be an identifier, but Infinity can start with '-' and '+' sign, so we must check
	// if this token start with a sign
	if char := lx.token.chars[0]; char == '+' || char == '-' {
		pos := newPosition(lx.pos.ln, lx.pos.col-len(lx.token.chars)-1)
		return errInvalidChar(char, pos, lx.token.chars, "'$', '_', unicode escape sequence, or any charater in categories Uppercase letter (Lu), Lowercase letter (Ll), Titlecase letter (Lt), Modifier letter (Lm), Other letter (Lo), Letter number (Nl)")
	}

	return lx.fetchIdentifier(false, char)
}

var nanChars = []rune{'a', 'N'}

// fetchNanNumber fetch number token with NaN as value
func (lx *Lexer) fetchNanNumber() error {
	lx.token.chars = append(lx.token.chars, 'N')
	for _, c := range nanChars {
		char, _, err := lx.r.ReadRune()
		if err != nil {
			if err == io.EOF {
				return errUnexpectedEOF(lx.pos, string([]rune{c}))
			}

			return err
		}

		if char != c {
			return lx.fetchIdentifier(false, char)
		}

		lx.token.chars = append(lx.token.chars, char)
	}

	char, _, err := lx.r.ReadRune()
	if err != nil {
		if err == io.EOF {
			lx.push()
			return nil
		}

		return err
	}

	if isCharPunct(char) {
		defer lx.fetchPunct(char)
		lx.pushWithPos(lx.pos.ln, lx.pos.col-1)
		return nil
	} else if isCharWhitespace(char) {
		lx.pushWithPos(lx.pos.ln, lx.pos.col-1)
		return nil
	}

	// if its turn out to be not a number, the nearest possibility is that this token
	// might be an identifier, but NaN can start with '-' and '+' sign, so we must check
	// if this token start with a sign
	if char := lx.token.chars[0]; char == '+' || char == '-' {
		pos := newPosition(lx.pos.ln, lx.pos.col-len(lx.token.chars)-1)
		return errInvalidChar(char, pos, lx.token.chars, "'$', '_', unicode escape sequence, or any charater in categories Uppercase letter (Lu), Lowercase letter (Ll), Titlecase letter (Lt), Modifier letter (Lm), Other letter (Lo), Letter number (Nl)")
	}

	return lx.fetchIdentifier(false, char)
}

// fetchOctalNumber fetch number token with octal numeric as value
func (lx *Lexer) fetchOctalNumber() error {
	isFirstChar := true
	for {
		char, _, err := lx.r.ReadRune()
		if err != nil {
			if err == io.EOF {
				if isFirstChar {
					return errUnexpectedEOF(lx.pos, "octal digit")
				}

				lx.push()
				return nil
			}

			return err
		}

		if !isCharValidOctal(char) {
			if isFirstChar {
				lx.token.chars = append(lx.token.chars, char)
				return errInvalidChar(char, lx.pos, lx.token.chars, "octal digit")
			}

			if isCharPunct(char) {
				defer lx.fetchPunct(char)
				lx.pushWithPos(lx.pos.ln, lx.pos.col-1)
				return nil
			} else if isCharWhitespace(char) {
				lx.pushWithPos(lx.pos.ln, lx.pos.col-1)
				return nil
			} else if char == '_' {
				lx.token.chars = append(lx.token.chars, char)
				char, _, err := lx.r.ReadRune()
				if err != nil {
					if err == io.EOF {
						return errUnexpectedEOF(lx.pos, "octal digit")
					}

					return err
				}

				lx.token.chars = append(lx.token.chars, char)
				if !isCharValidOctal(char) {
					return errInvalidChar(char, lx.pos, lx.token.chars, "octal digit")
				}

				continue
			}

			lx.token.chars = append(lx.token.chars, char)
			return errInvalidChar(char, lx.pos, lx.token.chars, "octal digit")
		}

		lx.token.chars = append(lx.token.chars, char)
		isFirstChar = false
	}
}

// fetchDoubleNumber fetch double number (number with decimal point, example: .123, 0.123, 1.234)
func (lx *Lexer) fetchDoubleNumber() error {
	lx.token.tokenNumSubType = tokenNumDouble
	lx.token.chars = append(lx.token.chars, '.')
	isFirstChar := true
	for {
		char, _, err := lx.r.ReadRune()
		if err != nil {
			if err == io.EOF {
				if isFirstChar {
					return errUnexpectedEOF(lx.pos, "decimal digit")
				}

				lx.push()
				return nil
			}

			return err
		}

		// possible exponent number
		if !unicode.IsDigit(char) {
			if isFirstChar {
				if char != 'e' && char != 'E' {
					lx.token.chars = append(lx.token.chars, char)
					return errInvalidChar(char, lx.pos, lx.token.chars, "decimal digit or exponent indicator")
				}

				lx.token.chars = append(lx.token.chars, char)
				return lx.fetchExponentNumber()
			}

			if char == 'e' || char == 'E' {
				lx.token.chars = append(lx.token.chars, char)
				return lx.fetchExponentNumber()
			} else if char == '_' {
				lx.token.chars = append(lx.token.chars, char)
				char, _, err := lx.r.ReadRune()
				if err != nil {
					if err == io.EOF {
						return errUnexpectedEOF(lx.pos, "decimal digit")
					}

					return err
				}

				lx.token.chars = append(lx.token.chars, char)
				if !unicode.IsDigit(char) {
					return errInvalidChar(char, lx.pos, lx.token.chars, "decimal digit")
				}

				continue
			} else if isCharPunct(char) {
				defer lx.fetchPunct(char)
				lx.pushWithPos(lx.pos.ln, lx.pos.col-1)
				return nil
			} else if isCharWhitespace(char) {
				lx.pushWithPos(lx.pos.ln, lx.pos.col-1)
				return nil
			}

			lx.token.chars = append(lx.token.chars, char)
			return errInvalidChar(char, lx.pos, lx.token.chars, "decimal digit, separator, or exponent indicator")
		}

		lx.token.chars = append(lx.token.chars, char)
		isFirstChar = false
	}
}

// fetchExponentNumber fetch number with exponent sign
func (lx *Lexer) fetchExponentNumber() error {
	lx.token.tokenNumSubType = tokenNumDouble
	// check first char
	char, _, err := lx.r.ReadRune()
	if err != nil {
		if err == io.EOF {
			return errUnexpectedEOF(lx.pos, "decimal digit")
		}

		return err
	}

	if !unicode.IsDigit(char) {
		if char == '-' || char == '+' {
			lx.token.chars = append(lx.token.chars, char)
			char, _, err := lx.r.ReadRune()
			if err != nil {
				if err == io.EOF {
					return errUnexpectedEOF(lx.pos, "decimal digit")
				}

				return err
			}

			if !unicode.IsDigit(char) {
				lx.token.chars = append(lx.token.chars, char)
				return errInvalidChar(char, lx.pos, lx.token.chars, "decimal digit")
			}

			lx.token.chars = append(lx.token.chars, char)
		} else {
			lx.token.chars = append(lx.token.chars, char)
			return errInvalidChar(char, lx.pos, lx.token.chars, "'+', '-', or decimal digit")
		}
	} else {
		lx.token.chars = append(lx.token.chars, char)
	}

	for {
		char, _, err := lx.r.ReadRune()
		if err != nil {
			if err == io.EOF {
				lx.push()
				return nil
			}

			return err
		}

		if !unicode.IsDigit(char) {
			if char == '_' {
				lx.token.chars = append(lx.token.chars, char)
				char, _, err := lx.r.ReadRune()
				if err != nil {
					if err == io.EOF {
						return errUnexpectedEOF(lx.pos, "decimal digit")
					}

					return err
				}

				lx.token.chars = append(lx.token.chars, char)
				if !unicode.IsDigit(char) {
					return errInvalidChar(char, lx.pos, lx.token.chars, "decimal digit")
				}

				continue
			} else if isCharWhitespace(char) {
				lx.pushWithPos(lx.pos.ln, lx.pos.col-1)
				return nil
			} else if isCharPunct(char) {
				defer lx.fetchPunct(char)
				lx.pushWithPos(lx.pos.ln, lx.pos.col-1)
				return nil
			}

			lx.token.chars = append(lx.token.chars, char)
			return errInvalidChar(char, lx.pos, lx.token.chars, "decimal digit, separator, whitespace, or punctuator")
		}

		lx.token.chars = append(lx.token.chars, char)
	}
}

// fetchNumber fetch number token
func (lx *Lexer) fetchNumber(beginChar rune) error {
	lx.token.t = TokenNumber
	lx.token.tokenNumSubType = tokenNumInteger
BEGIN_CHAR_CHECK:
	switch beginChar {
	case '-':
		lx.token.chars = append(lx.token.chars, beginChar)
		for {
			char, _, err := lx.r.ReadRune()
			if err != nil {
				if err == io.EOF {
					return errUnexpectedEOF(lx.pos, "decimal digit or '-'")
				}

				return err
			}

			if char == '-' {
				lx.token.chars = append(lx.token.chars, char)
				continue
			}

			beginChar = char
			goto BEGIN_CHAR_CHECK
		}

	case '+':
		lx.token.chars = append(lx.token.chars, beginChar)
		char, _, err := lx.r.ReadRune()
		if err != nil {
			if err == io.EOF {
				return errUnexpectedEOF(lx.pos, "decimal digit")
			}

			return err
		}

		beginChar = char
		goto BEGIN_CHAR_CHECK

	// double
	case '.':
		return lx.fetchDoubleNumber()

	// Infinity
	case 'I':
		return lx.fetchInfinityNumber()

	// NaN
	case 'N':
		return lx.fetchNanNumber()

	// possible hexadecimal, binary, octaldecimal, or double
	case '0':
		lx.token.chars = append(lx.token.chars, beginChar)
		char, _, err := lx.r.ReadRune()
		if err != nil {
			if err == io.EOF {
				lx.push()
				return nil
			}

			return nil
		}

		switch char {
		// hexadecimal
		case 'x', 'X':
			lx.token.chars = append(lx.token.chars, char)
			return lx.fetchHexaNumber()

		// binary
		case 'b', 'B':
			lx.token.chars = append(lx.token.chars, char)
			return lx.fetchBinaryNumber()

		// octaldecimal
		case 'o', 'O':
			lx.token.chars = append(lx.token.chars, char)
			return lx.fetchOctalNumber()

		// exponent
		case 'e', 'E':
			lx.token.chars = append(lx.token.chars, char)
			return lx.fetchExponentNumber()

		// double
		case '.':
			return lx.fetchDoubleNumber()

		// separator
		case '_':
			lx.token.chars = append(lx.token.chars, char)
			char, _, err := lx.r.ReadRune()
			if err != nil {
				if err == io.EOF {
					return errUnexpectedEOF(lx.pos, "decimal digit")
				}

				return err
			}

			lx.token.chars = append(lx.token.chars, char)
			if !unicode.IsDigit(char) {
				return errInvalidChar(char, lx.pos, lx.token.chars, "decimal digit")
			}

		default:
			if unicode.IsDigit(char) {
				lx.token.chars = append(lx.token.chars, char)
				break
			} else if isCharPunct(char) {
				defer lx.fetchPunct(char)
				lx.pushWithPos(lx.pos.ln, lx.pos.col-1)
				return nil
			} else if isCharWhitespace(char) {
				lx.pushWithPos(lx.pos.ln, lx.pos.col-1)
				return nil
			}

			lx.token.chars = append(lx.token.chars, char)
			return errInvalidChar(char, lx.pos, lx.token.chars, "decimal digit, hexadecimal indicator, octaldecimal indicator, binary indicator, decimal point, exponent indicator, punctuator, whitespace, or separator")
		}

	// decimal digit
	case '1', '2', '3', '4', '5', '6', '7', '8', '9':
		lx.token.chars = append(lx.token.chars, beginChar)

	default:
		lx.token.chars = append(lx.token.chars, beginChar)
		return errInvalidChar(beginChar, lx.pos, lx.token.chars, "decimal digit, decimal point, 'I' (Infinity), or 'N' (NaN)")
	}

	for {
		char, _, err := lx.r.ReadRune()
		if err != nil {
			if err == io.EOF {
				lx.push()
				return nil
			}

			return err
		}

		if !unicode.IsDigit(char) {
			switch char {
			case '.':
				return lx.fetchDoubleNumber()

			case 'e', 'E':
				lx.token.chars = append(lx.token.chars, char)
				return lx.fetchExponentNumber()

			case '_':
				lx.token.chars = append(lx.token.chars, char)
				char, _, err := lx.r.ReadRune()
				if err != nil {
					if err == io.EOF {
						return errUnexpectedEOF(lx.pos, "decimal digit")
					}

					return err
				}

				lx.token.chars = append(lx.token.chars, char)
				if !unicode.IsDigit(char) {
					return errInvalidChar(char, lx.pos, lx.token.chars, "decimal digit")
				}
			}

			if isCharPunct(char) {
				defer lx.fetchPunct(char)
				lx.pushWithPos(lx.pos.ln, lx.pos.col-1)
				return nil
			} else if isCharWhitespace(char) {
				lx.pushWithPos(lx.pos.ln, lx.pos.col-1)
				return nil
			}

			lx.token.chars = append(lx.token.chars, char)
			return errInvalidChar(char, lx.pos, lx.token.chars, "decimal digit, decimal point, exponent indicator, punctuator, whitespace, or separator")
		}

		lx.token.chars = append(lx.token.chars, char)
	}
}

// isCharValidHexa check if char is valid hexadecimal digit
func isCharValidHexa(char rune) bool {
	switch char {
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return true

	case 'a', 'b', 'c', 'd', 'e', 'f':
		return true

	case 'A', 'B', 'C', 'D', 'E', 'F':
		return true
	}

	return false
}

// isCharValidOctal check if char is valid octal digit
func isCharValidOctal(char rune) bool {
	switch char {
	case '0', '1', '2', '3', '4', '5', '6', '7':
		return true
	}

	return false
}

func isCharWhitespace(char rune) bool {
	switch char {
	case '\t', '\n', '\v', '\f', '\r', ' ', '\u00A0', '\u2028', '\u2029', '\uFEFF':
		return true

	default:
		if unicode.Is(unicode.Zs, char) {
			return true
		}
	}

	return false
}

func isCharPunct(char rune) bool {
	switch char {
	case '{', '}', '[', ']', ':', ',':
		return true
	}

	return false
}
