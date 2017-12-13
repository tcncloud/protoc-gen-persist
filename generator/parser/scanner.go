package parser

import (
	"bufio"
	"fmt"
	"io"
	"unicode"
)

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
		return &Token{
			tk:  ILLEGAL,
			raw: fmt.Sprintf("error around pos: %d, err: %v", s.pos, s.err.Error()),
		}
	}
	ch, stop := s.Read()
	if stop {
		return &Token{tk: EOF, raw: io.EOF.Error()}
	}
	if unicode.IsSpace(ch) {
		s.Unread()
		s.ScanWhitespace() // ignore the whitespace token
		return s.Scan()
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
	return s.ScanIdentifier()
}

func (s *Scanner) ScanWhitespace() *Token {
	var whitespace string
	for {
		ch, stop := s.Read()
		if stop || !unicode.IsSpace(ch) {
			if !stop {
				s.Unread()
			}
			return &Token{tk: WS, raw: whitespace}
		}
		whitespace += string(ch)
	}
}

func (s *Scanner) ScanIdentifier() *Token {
	ch, stop := s.Read()
	if stop {
		return &Token{
			tk:  ILLEGAL,
			raw: fmt.Sprintf("asked to scan identifier, but stop signal found at: %d", s.pos),
		}
	}
	s.Unread()
	if ch == '\'' || ch == '"' { // string literals
		return s.ScanString()
	} else if ch == '@' {
		return s.ScanDirective()
	} else if unicode.IsLetter(ch) { // SELECT, UPDATE, INTO, table_name, etc.
		return s.ScanSpecial()
	} else if unicode.IsNumber(ch) || ch == '-' { // floats and ints
		return s.ScanNumber()
	}

	return &Token{
		tk:  ILLEGAL,
		raw: fmt.Sprintf("invalid character \"%s\" at position: %d", string(ch), s.pos),
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
func (s *Scanner) ScanDirective() *Token {
	var raw string
	ch, stop := s.Read()
	if ch != '@' || stop {
		return &Token{
			tk:  ILLEGAL,
			raw: fmt.Sprintf("expected a field name or directive begining '@' near pos: %d", s.pos),
		}
	}
	raw += string(ch)

	ch, stop = s.Read()
	if stop {
		return &Token{
			tk: ILLEGAL,
			raw: fmt.Sprintf(
				"expected a letter, or '{', near pos: %d, instead found: %s",
				s.pos,
				string(ch),
			),
		}
	}
	startPos := s.pos
	raw += string(ch)
	if raw[1] == '{' {
		for {
			ch, stop := s.Read()
			if stop {
				return &Token{
					tk: ILLEGAL,
					raw: fmt.Sprintf(
						"asked to scan a directive '@{...}' starting at %d but stop signal found",
						startPos,
					),
				}
			}
			raw += string(ch)
			if ch == '}' {
				break
			}
		}
		return &Token{
			tk:  IDENT_DIRECTIVE,
			raw: raw,
		}
	}
	acceptedChar := func(r rune) bool {
		if unicode.IsNumber(r) || unicode.IsLetter(r) {
			return true
		}
		switch r {
		case '_', '-':
			return true
		}
		return false
	}
	for {
		ch, stop := s.Read()
		if stop {
			break
		} else if unicode.IsSpace(ch) || !acceptedChar(ch) {
			s.Unread()
			break
		}
		raw += string(ch)
	}
	return &Token{
		tk:  IDENT_FIELD,
		raw: raw,
	}
}
func (s *Scanner) ScanSpecial() *Token {
	var raw string

	acceptedSpecialChar := func(r rune) bool {
		if unicode.IsNumber(r) || unicode.IsLetter(r) {
			return true
		}
		switch r {
		case '_', '-':
			return true
		}
		return false
	}

	for {
		ch, stop := s.Read()
		if stop {
			break
		} else if unicode.IsSpace(ch) || !acceptedSpecialChar(ch) {
			s.Unread()
			break
		}
		raw += string(ch)
	}
	switch raw {
	case "SELECT", "Select", "select":
		return &Token{tk: SELECT, raw: raw}
	case "INSERT", "Insert", "insert":
		return &Token{tk: INSERT, raw: raw}
	case "DELETE", "Delete", "delete":
		return &Token{tk: DELETE, raw: raw}
	case "UPDATE", "Update", "update":
		return &Token{tk: UPDATE, raw: raw}
	case "INTO", "Into", "into":
		return &Token{tk: INTO, raw: raw}
	case "FROM", "From", "from":
		return &Token{tk: FROM, raw: raw}
	case "VALUES", "Values", "values":
		return &Token{tk: VALUES, raw: raw}
	case "true", "false":
		return &Token{tk: IDENT_BOOL, raw: raw}
	case "SET", "Set", "set":
		return &Token{tk: SET, raw: raw}
	case "KIND", "Kind", "kind":
		return &Token{tk: KIND, raw: raw}
	case "AND", "And", "and":
		return &Token{tk: AND, raw: raw}
	case "OR", "Or", "or":
		return &Token{tk: OR, raw: raw}
	case "START", "Start", "start":
		return &Token{tk: START, raw: raw}
	case "END", "End", "end":
		return &Token{tk: END, raw: raw}
	case "CO", "co", "CLOSED_OPEN", "closed_open", "closedOpen", "ClosedOpen":
		return &Token{tk: CLOSED_OPEN_KIND, raw: raw}
	case "CC", "cc", "CLOSED_CLOSED", "closed_closed", "closedClosed", "ClosedClosed":
		return &Token{tk: CLOSED_CLOSED_KIND, raw: raw}
	case "OC", "oc", "OPEN_CLOSED", "opend_closed", "openClosed", "OpenClosed":
		return &Token{tk: OPEN_CLOSED_KIND, raw: raw}
	case "OO", "oo", "OPEN_OPEN", "open_open", "openOpen", "OpenOpen":
		return &Token{tk: OPEN_OPEN_KIND, raw: raw}
	case "PK", "pk", "PRIMARY_KEY", "primary_key", "primaryKey", "PrimaryKey":
		return &Token{tk: PRIMARY_KEY, raw: raw}
	}
	if len(raw) > 0 {
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
	if stop || (ch != '\'' && ch != '"') {
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
		} else if ch == '\'' || ch == '"' {
			return &Token{tk: IDENT_STRING, raw: str}
		}
		str += string(ch)
	}
}
