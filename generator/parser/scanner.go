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

	//Operators
	EQUAL_OP
	GREATER_OP
	LESS_OP
	GREATER_EQUAL_OP
	LESS_EQUAL_OP

	IDENT_STRING
	IDENT_FLOAT
	IDENT_INT
	IDENT_TABLE_OR_COL
	IDENT_FIELD
	IDENT_BOOL

	// Keywords
	INSERT
	UPDATE
	DELETE
	SELECT
	FROM
	INTO
	VALUES
	SET
	AND
	OR
	START
	END
	KIND
	CLOSED_OPEN_KIND
	CLOSED_CLOSED_KIND
	OPEN_OPEN_KIND
	OPEN_CLOSED_KIND
)

type Token struct {
	tk  TokenKind
	raw string
}

type Scanner struct {
	r      *bufio.Reader
	err    error
	pos    int
	peeked []*Token
}

func NewScanner(r *bufio.Reader) *Scanner {
	return &Scanner{r: r}
}
func (s *Scanner) Read() (rune, bool) {
	ch, _, err := s.r.ReadRune()
	if err == io.EOF {
		return '0', true
	}
	s.pos++
	if err != nil {
		s.err = err
		return '0', true
	}
	return ch, false
}
func (s *Scanner) Unread() {
	if err := s.r.UnreadRune(); err != nil {
		s.err = err
	}
	s.pos--
}

// peek the number ahead of tokens
// since we cant unscan tokens, set these peeked tokens
// as next to be scanned
func (s *Scanner) Peek(num int) []*Token {
	var peeked []*Token
	for i := 0; i < num; i++ {
		peeked = append(peeked, s.Scan())
	}
	s.peeked = append(s.peeked, peeked...)
	return peeked
}

// if there is peeked tokens, return those first.
// otherwise if there is an error, return that as an illegal token,
// otherwise perform scan
func (s *Scanner) Scan() *Token {
	if len(s.peeked) != 0 {
		peeked := s.peeked[0]

		s.peeked = s.peeked[1:]
		return peeked
	}
	if s.err != nil {
		return &Token{tk: ILLEGAL, raw: s.err.Error()}
	}
	ch, stop := s.Read()
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
	case '=':
		return &Token{tk: EQUAL_OP, raw: string(ch)}
	case '<':
		next, _ := s.Read()
		if next == '=' {
			return &Token{tk: LESS_EQUAL_OP, raw: string(ch) + string(next)}
		}
		s.Unread()
		return &Token{tk: LESS_OP, raw: string(ch)}
	case '>':
		next, _ := s.Read()
		if next == '=' {
			return &Token{tk: GREATER_EQUAL_OP, raw: string(ch) + string(next)}
		}
		s.Unread()
		return &Token{tk: GREATER_OP, raw: string(ch)}
	}
	s.Unread()
	if unicode.IsSpace(ch) {
		s.ScanWhitespace() // ignore the whitespace token
	}
	return s.ScanIdentifier()
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
	if len(rawNum) == 0 || rawNum == "." { // we dont have a number, and we were asked to parse one
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
	case "SET", "set":
		return &Token{tk: SET, raw: raw}
	case "KIND", "kind":
		return &Token{tk: KIND, raw: raw}
	case "AND", "and":
		return &Token{tk: AND, raw: raw}
	case "OR", "or":
		return &Token{tk: OR, raw: raw}
	case "START", "start":
		return &Token{tk: START, raw: raw}
	case "END", "end":
		return &Token{tk: END, raw: raw}
	case "CO", "co", "CLOSED_OPEN", "closed_open", "closedOpen", "ClosedOpen":
		return &Token{tk: CLOSED_OPEN_KIND, raw: raw}
	case "CC", "cc", "CLOSED_CLOSED", "closed_closed", "closedClosed", "ClosedClosed":
		return &Token{tk: CLOSED_CLOSED_KIND, raw: raw}
	case "OC", "oc", "OPEN_CLOSED", "opend_closed", "openClosed", "OpenClosed":
		return &Token{tk: OPEN_CLOSED_KIND, raw: raw}
	case "OO", "oo", "OPEN_OPEN", "open_open", "openOpen", "OpenOpen":
		return &Token{tk: OPEN_OPEN_KIND, raw: raw}
	}
	if len(raw) > 1 && raw[0] == '@' {
		return &Token{tk: IDENT_FIELD, raw: raw}
	} else if len(raw) > 0 {
		return &Token{tk: IDENT_TABLE_OR_COL, raw: raw}
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
