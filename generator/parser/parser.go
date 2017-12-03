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
		return (&SelectMode{tkns: []*Token{tkn}}).Parse(scanner)
	case INSERT:
		return (&InsertMode{}).Parse(scanner)
	case UPDATE:
		return (&UpdateMode{}).Parse(scanner)
	case DELETE:
		return (&DeleteMode{}).Parse(scanner)
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

type InsertMode struct{}

func (m *InsertMode) Parse(scanner *Scanner) (Query, error) {
	return nil, fmt.Errorf("unimplemented")
}

type UpdateMode struct{}

func (m *UpdateMode) Parse(scanner *Scanner) (Query, error) {
	return nil, fmt.Errorf("unimplemented")
}

type DeleteMode struct{}

func (m *DeleteMode) Parse(scanner *Scanner) (Query, error) {
	return nil, fmt.Errorf("unimplemented")
}
