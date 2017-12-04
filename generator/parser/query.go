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
	fields    []int
	cols      []*Token
	values    []*Token
	tableName *Token
}

func (q *InsertQuery) String() string {
	return ""
}
func (q *InsertQuery) Tokens() []*Token {
	return nil
}
func (q *InsertQuery) Type() QueryType {
	return INSERT_QUERY
}
func (q *InsertQuery) Table() string {
	return ""
}
func (q *InsertQuery) Fields() []string {
	return nil
}
func (q *InsertQuery) Args() []*Token {
	return nil
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

type DeleteQuery struct{}

func (q *DeleteQuery) String() string {
	return ""
}
func (q *DeleteQuery) Tokens() []*Token {
	return nil
}
func (q *DeleteQuery) Type() QueryType {
	return DELETE_QUERY
}
func (q *DeleteQuery) Table() string {
	return ""
}
func (q *DeleteQuery) Fields() []string {
	return nil
}
func (q *DeleteQuery) Args() []*Token {
	return nil
}

type UpdateQuery struct {
	tokens    []*Token
	fields    []int
	cols      []*Token
	values    []*Token
	tableName *Token
	pk        []*Token
}

func (q *UpdateQuery) String() string {
	return ""
}
func (q *UpdateQuery) Tokens() []*Token {
	return nil
}
func (q *UpdateQuery) Type() QueryType {
	return UPDATE_QUERY
}
func (q *UpdateQuery) Table() string {
	return ""
}
func (q *UpdateQuery) Fields() []string {
	return nil
}
func (q *UpdateQuery) Args() []*Token {
	return nil
}

type QueryType int

const (
	SELECT_QUERY QueryType = iota
	UPDATE_QUERY
	DELETE_QUERY
	INSERT_QUERY
)
