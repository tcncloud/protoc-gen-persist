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
	tkn := scanner.Peek(1)[0]
	switch tkn.tk {
	case SELECT:
		return NewSelectMode().Parse(scanner)
	case INSERT:
		return NewInsertMode().Parse(scanner)
	case UPDATE:
		return NewUpdateMode().Parse(scanner)
	case DELETE:
		return NewDeleteMode().Parse(scanner)
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

func NewSelectMode() *SelectMode {
	return &SelectMode{}
}
func (m *SelectMode) Parse(scanner *Scanner) (Query, error) {
	// we have scanned out the select word, so add it back
	query := ""
	eater := NewEater(scanner)
	if eater.Eat(SELECT) {
		query += eater.TakeTokens()[0].raw
	} else {
		return nil, fmt.Errorf("told to parse a select query, but SELECT not found")
	}
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

func NewInsertMode() *InsertMode {
	return &InsertMode{}
}

//INSERT INTO tablename <optional> (colname, colname) VALUES(...)
func (m *InsertMode) Parse(scanner *Scanner) (Query, error) {
	var table *Token
	var cols []*Token
	var values []*Token
	var fields []int

	eater := NewEater(scanner)
	eater.Eat(INSERT)
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

	// validate columns
	if len(cols) != len(values) {
		return nil, fmt.Errorf(
			"columns len: %d, does not match values len: %d",
			len(cols),
			len(values),
		)
	}
	// compute the field positions
	for i, v := range values {
		if v.tk == IDENT_FIELD {
			fields = append(fields, i)
		}
	}
	return &InsertQuery{
		tokens:    eater.TakeTokens(),
		fields:    fields,
		cols:      cols,
		values:    values,
		tableName: table,
	}, nil
}

type UpdateMode struct{}

func NewUpdateMode() *UpdateMode {
	return &UpdateMode{}
}

// Normal Update sql query syntax does not really make a ton of sense here
// for now, just support:
// Update table <set loop> | <(col), VALUES (values)> PK(cols)
func (m *UpdateMode) Parse(scanner *Scanner) (Query, error) {
	var cols []*Token
	var values []*Token
	var fields []int
	var table *Token
	var primaryKey []*Token

	eater := NewEater(scanner)
	eater.Eat(UPDATE)
	eater.Eat(IDENT_TABLE_OR_COL)

	//either array of col names, or set loop
	switch scanner.Peek(1)[0].tk {
	case OPEN_PARAN:
		cs, _ := eater.EatArrayOf(IDENT_TABLE_OR_COL)
		cols = append(cols, cs...)
		eater.Eat(VALUES)
		vs, _ := eater.EatArrayOf(
			IDENT_STRING,
			IDENT_FLOAT,
			IDENT_INT,
			IDENT_FIELD,
			IDENT_BOOL,
		)
		values = append(values, vs...)
	case SET:
		for {
			eater.Eat(SET)
			if eater.Eat(IDENT_TABLE_OR_COL) {
				cols = append(cols, eater.Top())
			}
			eater.Eat(EQUAL_OP)
			if eater.Eat(
				IDENT_STRING,
				IDENT_FLOAT,
				IDENT_INT,
				IDENT_FIELD,
				IDENT_BOOL,
			) {
				values = append(values, eater.Top())
			}
			if scanner.Peek(1)[0].tk != COMMA {
				break
			}
			if !eater.Eat(COMMA) {
				break
			}
		}
	}
	eater.Eat(PRIMARY_KEY)
	pk, _ := eater.EatArrayOf(IDENT_TABLE_OR_COL)
	primaryKey = append(primaryKey, pk...)

	if !eater.Eat(EOF) {
		return nil, eater.TakeErr()
	}

	for i, tkn := range values {
		if tkn.tk == IDENT_FIELD {
			fields = append(fields, i)
		}
	}
	return &UpdateQuery{
		tokens:    eater.TakeTokens(),
		fields:    fields,
		cols:      cols,
		values:    values,
		tableName: table,
		pk:        primaryKey,
	}, nil
}

type DeleteMode struct {
}

func NewDeleteMode() *DeleteMode {
	return &DeleteMode{}
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

// eats each token in group, in order till the FIRST token in the group does not match
// a comma resets.  It only fails if a noncomplete group is eaten, or nothing is eaten
func (e *eater) EatCommaSeperatedGroupOf(group ...TokenKind) ([][]*Token, bool) {
	if e.lastErr != nil {
		return nil, false
	}
	tokenGroups := make([][]*Token, 0)
	scanOneGroup := func() (g []*Token) {
		for _, kind := range group {
			if e.Eat(kind) {
				g = append(g, e.Top())
			}
		}
		return
	}
	if e.scanner.Peek(1)[0].tk == group[0] {
		e.lastErr = fmt.Errorf("asked to eat a group of %+v, but none was found", group)
		return nil, false
	}
	for {
		group := scanOneGroup()
		if e.lastErr != nil {
			return nil, false
		}
		tokenGroups = append(tokenGroups, group)
		if e.scanner.Peek(1)[0].tk != COMMA {
			break
		}
		e.Eat(COMMA)
	}

	return tokenGroups, true
}

// eats the pattern: <(expected, expected, ...)> and returns the expected
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
