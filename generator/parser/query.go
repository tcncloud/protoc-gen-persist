package parser

type Query interface {
	String() string // golang syntax of the query
	Tokens() []*Token
	Type() QueryType
	Table() string
	Fields() []string
	Args() []*Token
}

type InsertQuery struct {
	tokens    []*Token
	cols      []*Token
	values    []*Token
	tableName *Token
}

func (q *InsertQuery) String() string {
	return ""
}
func (q *InsertQuery) Tokens() []*Token {
	return q.tokens
}
func (q *InsertQuery) Type() QueryType {
	return INSERT_QUERY
}
func (q *InsertQuery) Table() string {
	return q.tableName.raw
}
func (q *InsertQuery) Fields() []string {
	cs := make([]string, len(q.cols))
	for i, tkn := range q.cols {
		cs[i] = tkn.raw
	}
	return cs
}
func (q *InsertQuery) Args() []*Token {
	var args []*Token
	for _, tkn := range q.values {
		if tkn.tk == IDENT_FIELD {
			args = append(args, tkn)
		}
	}
	return args
}

type SelectQuery struct {
	query string
}

func NewSelectQuery(q string) *SelectQuery {
	return &SelectQuery{query: q}
}
func (q *SelectQuery) String() string {
	return q.query
}
func (q *SelectQuery) Tokens() []*Token {
	return nil
}
func (q *SelectQuery) Type() QueryType {
	return SELECT_QUERY
}
func (q *SelectQuery) Table() string {
	return ""
}
func (q *SelectQuery) Fields() []string {
	return nil
}
func (q *SelectQuery) Args() []*Token {
	return nil
}

type DeleteQuery struct {
	tokens       []*Token
	start        []*Token
	end          []*Token
	kind         *Token
	cols         []*Token
	values       []*Token
	pk           []*Token
	table        *Token
	usesKeyRange bool
}

func (q *DeleteQuery) String() string {
	return ""
}
func (q *DeleteQuery) Tokens() []*Token {
	return q.tokens
}
func (q *DeleteQuery) Type() QueryType {
	return DELETE_QUERY
}
func (q *DeleteQuery) Table() string {
	return q.table.raw
}
func (q *DeleteQuery) Fields() []string {
	if q.usesKeyRange {
		return nil
	}
	fields := make([]string, len(q.cols))
	for i, tkn := range q.cols {
		fields[i] = tkn.raw
	}
	return fields
}
func (q *DeleteQuery) Args() []*Token {
	var args []*Token
	putInArgsFrom := func(arr []*Token) {
		for _, tkn := range arr {
			if tkn.tk == IDENT_FIELD {
				args = append(args, tkn)
			}
		}
	}
	if q.usesKeyRange {
		putInArgsFrom(q.start)
		putInArgsFrom(q.end)
	} else {
		putInArgsFrom(q.values)
	}
	return args
}

type UpdateQuery struct {
	tokens    []*Token
	cols      []*Token
	values    []*Token
	tableName *Token
	pk        []*Token
}

func (q *UpdateQuery) String() string {
	return ""
}
func (q *UpdateQuery) Tokens() []*Token {
	return q.tokens
}
func (q *UpdateQuery) Type() QueryType {
	return UPDATE_QUERY
}
func (q *UpdateQuery) Table() string {
	return q.tableName.raw
}
func (q *UpdateQuery) Fields() []string {
	cs := make([]string, len(q.cols))
	for i, tkn := range q.cols {
		cs[i] = tkn.raw
	}
	return cs
}
func (q *UpdateQuery) Args() []*Token {
	var fields []*Token
	for _, tkn := range q.values {
		if tkn.tk == IDENT_FIELD {
			fields = append(fields, tkn)
		}
	}
	return fields
}

type QueryType int

const (
	SELECT_QUERY QueryType = iota
	UPDATE_QUERY
	DELETE_QUERY
	INSERT_QUERY
)
