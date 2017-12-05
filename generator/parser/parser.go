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

	eater := NewEater(scanner)
	eater.Eat(INSERT)
	eater.Eat(INTO)
	if eater.Eat(IDENT_TABLE_OR_COL) {
		table = eater.Top()
	}
	// we have an optional array here
	peeked := scanner.Peek(1)[0]
	fmt.Printf("peeked ahead and got: %+v\n", peeked)
	if scanner.Peek(1)[0].tk == OPEN_PARAN {
		fmt.Printf("we have an open paran\n")
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
	return &InsertQuery{
		tokens:    eater.TakeTokens(),
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
	return &UpdateQuery{
		tokens:    eater.TakeTokens(),
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

// supports two formats:
// delete key range: DELETE FROM table_name START(...) END(...) KIND(...)
// delete record pk: DELETE FROM table_name (...) VALUES (...) PRIMARY_KEY(...)
func (m *DeleteMode) Parse(scanner *Scanner) (Query, error) {
	var table *Token
	var kind *Token
	var start []*Token
	var end []*Token
	var cols []*Token
	var values []*Token
	var primaryKey []*Token
	var usesKeyRange bool

	eater := NewEater(scanner)
	eater.Eat(DELETE)
	eater.Eat(FROM)
	if eater.Eat(IDENT_TABLE_OR_COL) {
		table = eater.Top()
	}
	switch scanner.Peek(1)[0].tk {
	case START: // is a key range query
		usesKeyRange = true
		eater.Eat(START)
		s, _ := eater.EatArrayOf(
			IDENT_STRING,
			IDENT_FLOAT,
			IDENT_INT,
			IDENT_FIELD,
			IDENT_BOOL,
		)
		start = append(start, s...)

		eater.Eat(END)
		e, _ := eater.EatArrayOf(
			IDENT_STRING,
			IDENT_FLOAT,
			IDENT_INT,
			IDENT_FIELD,
			IDENT_BOOL,
		)
		end = append(end, e...)

		eater.Eat(KIND)
		eater.Eat(OPEN_PARAN)
		if eater.Eat(
			CLOSED_OPEN_KIND,
			CLOSED_CLOSED_KIND,
			OPEN_OPEN_KIND,
			OPEN_CLOSED_KIND,
		) {
			kind = eater.Top()
		}
		eater.Eat(CLOSE_PARAN)
	case OPEN_PARAN: // is a delete single record query
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

		eater.Eat(PRIMARY_KEY)
		pk, _ := eater.EatArrayOf(IDENT_TABLE_OR_COL)
		primaryKey = append(primaryKey, pk...)
	}
	if !eater.Eat(EOF) {
		return nil, eater.TakeErr()
	}

	return &DeleteQuery{
		tokens:       eater.TakeTokens(),
		start:        start,
		end:          end,
		kind:         kind,
		cols:         cols,
		pk:           primaryKey,
		table:        table,
		usesKeyRange: usesKeyRange,
	}, nil
}
