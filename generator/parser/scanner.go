package parser

import (
	"bufio"
	"fmt"
	"io"
	"unicode"
)

type TokenKind int

const (
	// special
	ILLEGAL TokenKind = iota
	EOF
	WS

	COMMA
	OPEN_PARAN
	CLOSE_PARAN

	IDENT_STRING
	IDENT_FLOAT
	IDENT_INT
	IDENT_TABLE
	IDENT_FIELD
	IDENT_BOOL

	// Keywords
	INSERT
	UPDATE
	DELETE
	SELECT // (SELECT ...) is the only allowed format.  we dont parse anything else.
	FROM
	INTO
	VALUES
)

type Token struct {
	tk  TokenKind
	raw string
}

type Scanner struct {
	r    bufio.Reader
	err  error
	mode Mode
	pos  int
}

func (s *Scanner) Read() (rune, bool) {
	ch, _, err := s.r.ReadRune()
	if err == io.EOF {
		return '0', true
	} else if err != nil {
		s.err = err
		return '0', true
	}
	return ch, false
}
func (s *Scanner) Unread() {
	if err := s.r.UnreadRune(); err != nil {
		s.err = err
	}
}

func (s *Scanner) Scan() *Token {
	if s.err != nil {
		return &Token{tk: ILLEGAL, raw: s.err.Error()}
	}
	return s.mode.Scan(s)

}

func (s *Scanner) ScanWhitespace() *Token {
	var whitespace string
	for {
		ch, stop := s.Read()
		if stop || !unicode.IsSpace(ch) {
			s.Unread()
			return &Token{tk: WS, raw: whitespace}
		}
		whitespace += string(ch)
	}
}

func (s *Scanner) ScanIdentifier() *Token {
	for {
		ch, stop := s.Read()
		if stop {
			return &Token{
				tk:  ILLEGAL,
				raw: fmt.Sprintf("asked to scan identifier, but stop signal found at: %d", s.pos),
			}
		}
		s.Unread()
		if ch == '\'' { // string literals
			return s.ScanString()
		} else if unicode.IsLetter(ch) { // SELECT, UPDATE, INTO, table_name, etc.
			return s.ScanSpecial()
		} else if unicode.IsNumber(ch) { // floats and ints
			return s.ScanNumber()
		}

		return &Token{
			tk:  ILLEGAL,
			raw: fmt.Sprintf("invalid character \"%s\" at position: %d", string(ch), s.pos),
		}
	}
}

func (s *Scanner) ScanNumber() *Token {
	var rawNum string
	var isFloat bool
	for {
		ch, stop := s.Read()
		if stop {
			break
		}
		if ch == '.' && isFloat {
			return &Token{
				tk:  ILLEGAL,
				raw: fmt.Sprintf("additional \".\" found in float identifier at pos: %d", s.pos),
			}
		} else if ch == '.' {
			rawNum += string(ch)
			isFloat = true
		} else if !unicode.IsNumber(ch) {
			s.Unread()
			break
		} else {
			// we finally have a number character
			rawNum += string(ch)
		}
	}
	if len(rawNum) == 0 || rawNum == "." { // we do not have a number, and we were asked to parse one
		return &Token{
			tk: ILLEGAL,
			raw: fmt.Sprintf(
				"expected a number at pos: %d, instead found illegal char, or stop signal",
				s.pos,
			),
		}
	}
	if isFloat {
		return &Token{tk: IDENT_FLOAT, raw: rawNum}
	}
	return &Token{tk: IDENT_INT, raw: rawNum}
}
func (s *Scanner) ScanSpecial() *Token {
	var raw string

	acceptedSpecialChar := func(r rune) bool {
		if unicode.IsNumber(r) || unicode.IsLetter(r) {
			return true
		}
		switch r {
		case '@', '_', '-':
			return true
		}
		return false
	}

	for {
		ch, stop := s.Read()
		if stop {
			break
		} else if unicode.IsSpace(ch) {
			s.Unread()
			break
		} else if !acceptedSpecialChar(ch) {
			break
		}
		raw += string(ch)
	}
	switch raw {
	case "SELECT", "select":
		return &Token{tk: SELECT, raw: raw}
	case "INSERT", "insert":
		return &Token{tk: INSERT, raw: raw}
	case "DELETE", "delete":
		return &Token{tk: DELETE, raw: raw}
	case "UPDATE", "update":
		return &Token{tk: UPDATE, raw: raw}
	case "INTO", "into":
		return &Token{tk: INTO, raw: raw}
	case "FROM", "from":
		return &Token{tk: FROM, raw: raw}
	case "VALUES", "values":
		return &Token{tk: VALUES, raw: raw}
	case "true", "false":
		return &Token{tk: IDENT_BOOL, raw: raw}
	}
	if len(raw) > 1 && raw[0] == '@' {
		return &Token{tk: IDENT_FIELD, raw: raw}
	} else if len(raw) > 0 {
		return &Token{tk: IDENT_TABLE, raw: raw}
	}
	return &Token{
		tk:  ILLEGAL,
		raw: fmt.Sprintf("unknown token: '%s', at pos: %d", raw, s.pos),
	}
}

func (s *Scanner) ScanString() *Token {
	var str string

	ch, stop := s.Read()
	if stop || ch != '\'' {
		return &Token{
			tk:  ILLEGAL,
			raw: fmt.Sprintf("expected \"'\"  at position: %d", s.pos),
		}
	}
	for {
		ch, stop := s.Read()
		if stop {
			return &Token{
				tk:  ILLEGAL,
				raw: fmt.Sprintf("asked to scan string but stop signal found at: %d", s.pos),
			}
		} else if ch == '\'' {
			return &Token{tk: IDENT_STRING, raw: str}
		}
		str += string(ch)
	}
}

// a mode is used to direct whether we should be skipping whitespace
type Mode interface {
	Scan(*Scanner) *Token
}

type NormalMode struct{}

// scan all whitespace out and return one whitespace token
func (m *NormalMode) Scan(scanner *Scanner) *Token {
	ch, stop := scanner.Read()
	if stop {
		return &Token{tk: EOF, raw: io.EOF.Error()}
	}
	// add more cases later
	switch ch {
	case ',':
		return &Token{tk: COMMA, raw: string(ch)}
	case '(':
		return &Token{tk: OPEN_PARAN, raw: string(ch)}
	case ')':
		return &Token{tk: CLOSE_PARAN, raw: string(ch)}
	}

	scanner.Unread()
	if unicode.IsSpace(ch) {
		return scanner.ScanWhitespace()
	} else {
		return scanner.ScanIdentifier()
	}
}
