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
	Fields    map[string]int
	TableName string
}

func (i *InsertQuery) String() string {
	return ""
}
func (i *InsertQuery) Token() []*Token {
	return nil
}
func (i *InsertQuery) Type() QueryType {
	return INSERT_QUERY
}
func (i *InsertQuery) Table() string {
	return ""
}
func (i *InsertQuery) Fields() []string {
	return nil
}
func (i *InsertQuery) Args() []*Token {
	return nil
}

type SelectQuery struct{}

func (i *SelectQuery) String() string {
	return ""
}
func (i *SelectQuery) Token() []*Token {
	return nil
}
func (i *SelectQuery) Type() QueryType {
	return SELECT_QUERY
}
func (i *SelectQuery) Table() string {
	return ""
}
func (i *SelectQuery) Fields() []string {
	return nil
}
func (i *SelectQuery) Args() []*Token {
	return nil
}

type DeleteQuery struct{}

func (i *DeleteQuery) String() string {
	return ""
}
func (i *DeleteQuery) Token() []*Token {
	return nil
}
func (i *DeleteQuery) Type() QueryType {
	return DELETE_QUERY
}
func (i *DeleteQuery) Table() string {
	return ""
}
func (i *DeleteQuery) Fields() []string {
	return nil
}
func (i *DeleteQuery) Args() []*Token {
	return nil
}

type UpdateQuery struct{}

func (i *UpdateQuery) String() string {
	return ""
}
func (i *UpdateQuery) Token() []*Token {
	return nil
}
func (i *UpdateQuery) Type() QueryType {
	return UPDATE_QUERY
}
func (i *UpdateQuery) Table() string {
	return ""
}
func (i *UpdateQuery) Fields() []string {
	return nil
}
func (i *UpdateQuery) Args() []*Token {
	return nil
}

type QueryType int

const (
	SELECT_QUERY QueryType = iota
	UPDATE_QUERY
	DELETE_QUERY
	INSERT_QUERY
)
