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
	IsLegal(tkn *Token) bool
	Parse(s *Scanner) (Query, error)
}

type StartMode struct {
	query *Query
}

func (m *StartMode) Parse(scanner *Scanner) (Query, error) {
	return nil, fmt.Errorf("unimplemented")
}
