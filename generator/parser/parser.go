package parser

import (
	"bufio"
	"fmt"
	"io"
)

type Parser struct {
	s *Scanner
}

func NewParser(r io.Reader) *Parser {
	buf := bufio.NewReader(r)
	s := NewScanner(buf)

	return &Parser{s: s}
}

func (p *Parser) Parse() (Query, error) {
	q, err := (&StartMode{}).Parse(p.s)
	if err != nil {
		return nil, fmt.Errorf("Failure parsing query.\n err: %v", err)
	}
	return q, nil
}

type Mode interface {
	Parse(s *Scanner) (Query, error)
}

type StartMode struct{}

func (m *StartMode) Parse(scanner *Scanner) (Query, error) {
	tkn := scanner.Scan()
	switch tkn.tk {
	case SELECT:
		return NewSelectMode([]*Token{tkn}).Parse(scanner)
	case INSERT:
		return NewInsertMode([]*Token{tkn}).Parse(scanner)
	case UPDATE:
		return NewUpdateMode([]*Token{tkn}).Parse(scanner)
	case DELETE:
		return NewDeleteMode([]*Token{tkn}).Parse(scanner)
	case EOF:
		return nil, fmt.Errorf("no query found")
	case ILLEGAL:
		return nil, fmt.Errorf("scanned illegal token: %s", tkn.raw)
	default:
		return nil, fmt.Errorf("token kind %d, raw: '%s' not acceptable here", tkn.tk, tkn.raw)
	}
}

// we don't parse anything for select mode
type SelectMode struct {
	tkns []*Token
}

func NewSelectMode(tkns []*Token) *SelectMode {
	return &SelectMode{tkns: tkns}
}
func (m *SelectMode) Parse(scanner *Scanner) (Query, error) {
	// we have scanned out the select word, so add it back
	query := m.tkns[0].raw
	for {
		ch, stop := scanner.Read()
		if stop {
			break
		}
		query += string(ch)
	}
	return NewSelectQuery(query), nil
}

type InsertMode struct {
	tkns []*Token
}

func NewInsertMode(tkns []*Token) *InsertMode {
	return &InsertMode{tkns: tkns}
}

//INSERT INTO tablename <optional> (colname, colname) VALUES(...)
func (m *InsertMode) Parse(scanner *Scanner) (Query, error) {
	var table *Token
	var cols []*Token
	var values []*Token

	eater := NewEater(scanner)

	eater.Eat(INTO)
	if eater.Eat(IDENT_TABLE_OR_COL) {
		table = eater.Top()
	}
	// we have an optional array here
	if scanner.Peek(1)[0].tk == OPEN_PARAN {
		// declaring column ordering here
		cols, _ = eater.EatArrayOf(IDENT_TABLE_OR_COL)
	}
	eater.Eat(VALUES)
	values, _ = eater.EatArrayOf(
		IDENT_STRING,
		IDENT_FLOAT,
		IDENT_INT,
		IDENT_FIELD,
		IDENT_BOOL,
	)
	if !eater.Eat(EOF) {
		return nil, eater.TakeErr()
	}

	m.tkns = append(m.tkns, eater.TakeTokens()...)

	// validate columns
	if len(cols) != len(values) {
		return nil, fmt.Errorf(
			"columns len: %d, does not match values len: %d",
			len(cols),
			len(values),
		)
	}
	// compute the field positions
	var fields []int
	for i, v := range values {
		if v.tk == IDENT_FIELD {
			fields = append(fields, i)
		}
	}
	return &InsertQuery{
		tokens:    m.tkns,
		fields:    fields,
		cols:      cols,
		values:    values,
		tableName: table,
	}, nil
}

type UpdateMode struct {
	tkns []*Token
}

func NewUpdateMode(tkns []*Token) *UpdateMode {
	return &UpdateMode{tkns: tkns}
}

// Normal Update sql query syntax does not really make a ton of sense here
// for now, just support
func (m *UpdateMode) Parse(scanner *Scanner) (Query, error) {
	return nil, fmt.Errorf("unimplemented")
}

type DeleteMode struct {
	tkns []*Token
}

func NewDeleteMode(tkns []*Token) *DeleteMode {
	return &DeleteMode{tkns: tkns}
}
func (m *DeleteMode) Parse(scanner *Scanner) (Query, error) {
	return nil, fmt.Errorf("unimplemented")
}

type eater struct {
	scanner    *Scanner
	lastTokens []*Token
	lastErr    error
}

func NewEater(scanner *Scanner) *eater {
	return &eater{scanner: scanner}
}

func (e *eater) Eat(expected ...TokenKind) bool {
	if e.lastErr != nil {
		return false
	}
	tkn := e.scanner.Scan()
	for _, kind := range expected {
		if tkn.tk == kind {
			e.lastTokens = append(e.lastTokens, tkn)
			return true
		}
	}
	e.lastErr = fmt.Errorf(
		"unexpected token: %+v, expected one in kinds: %+v",
		*tkn,
		expected,
	)
	return false
}

func (e *eater) EatArrayOf(expected ...TokenKind) ([]*Token, bool) {
	// for returning the value tokens in the array
	var values []*Token

	e.Eat(OPEN_PARAN)
	e.Eat(expected...)
	for {
		if e.scanner.Peek(1)[0].tk == CLOSE_PARAN {
			break
		}
		if !e.Eat(COMMA) {
			return nil, false
		}
		if !e.Eat(expected...) {
			return nil, false
		}
		// populate values with the top of the array
		values = append(values, e.Top())
	}
	e.Eat(CLOSE_PARAN)

	return values, e.lastErr == nil
}

func (e *eater) Top() *Token {
	if len(e.lastTokens) == 0 {
		return &Token{tk: ILLEGAL, raw: "cannot get top of empty token slice"}
	}
	return e.lastTokens[len(e.lastTokens)-1]
}

func (e *eater) TakeTokens() []*Token {
	tkns := e.lastTokens
	e.lastTokens = nil
	return tkns
}

func (e *eater) TakeErr() error {
	err := e.lastErr
	e.lastErr = nil
	return err
}
