package parser

import (
	"fmt"
)

type Query interface {
	String() string // golang syntax of the query
	Tokens() []*Token
	Type() QueryType
	Table() string
	Fields() []string
	Args() []*Token
	// map of field, to go syntax string
	SetParams(map[string]string)
	AddParam(key, val string)
}

type InsertQuery struct {
	tokens    []*Token
	cols      []*Token
	values    []*Token
	tableName *Token
	params    map[string]string
}

func (q *InsertQuery) String() string {
	valMap := make(map[string]string)
	for i, tkn := range q.cols {
		val := q.values[i]
		if val.tk == IDENT_FIELD {
			valMap[tkn.raw] = q.params[val.raw]
		} else {
			valMap[tkn.raw] = SyntaxStringFromIdent(val)
		}
	}
	valMapString := "map[string]interface{}{"
	for k, v := range valMap {
		valMapString += fmt.Sprintf("\n\t\"%s\": %s,", k, v)
	}
	valMapString += "\n}"
	return fmt.Sprintf(`spanner.InsertMap("%s", %s)`, q.tableName.raw, valMapString)
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
func (q *InsertQuery) SetParams(p map[string]string) {
	q.params = p
}
func (q *InsertQuery) AddParam(key, val string) {
	if q.params == nil {
		q.params = make(map[string]string)
	}
	q.params[key] = val
}

type SelectQuery struct {
	query  string
	params map[string]string
	fields []*Token
}

func NewSelectQuery(q string, fields []*Token) *SelectQuery {
	return &SelectQuery{query: q, fields: fields}
}
func (q *SelectQuery) String() string {
	params := "map[string]interface{}{"
	for _, v := range q.fields {
		params += fmt.Sprintf("\n\t\"%s\": %s,", v.raw[1:], q.params[v.raw])
	}
	params += "\n}"
	return fmt.Sprintf(`spanner.Statement{
	SQL: "%s",
	Params: %s,
}`, q.query, params)
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
func (q *SelectQuery) Fields() (out []string) {
	for _, t := range q.fields {
		out = append(out, t.raw)
	}
	return
}
func (q *SelectQuery) Args() []*Token {
	return q.fields
}
func (q *SelectQuery) SetParams(p map[string]string) {
	q.params = p
}
func (q *SelectQuery) AddParam(key, val string) {
	if q.params == nil {
		q.params = make(map[string]string)
	}
	q.params[key] = val
}

type DeleteQuery struct {
	tokens       []*Token
	start        []*Token
	end          []*Token
	kind         *Token
	values       []*Token
	table        *Token
	usesKeyRange bool
	params       map[string]string
}

func (q *DeleteQuery) String() string {
	var key string
	if q.usesKeyRange {
		key = q.StringKR()
	} else {
		key = q.StringSingle()
	}
	return fmt.Sprintf("spanner.Delete(\"%s\", %s)", q.Table(), key)
}
func (q *DeleteQuery) addToSyntaxStr(str *string, arr []*Token) {
	for _, v := range arr {
		var val string
		if v.tk == IDENT_FIELD {
			val = q.params[v.raw]
		} else {
			val = SyntaxStringFromIdent(v)
		}
		*str += fmt.Sprintf("\n\t%s,", val)
	}
}

func (q *DeleteQuery) StringKR() string {
	keyRange := "spanner.KeyRange{\n\tStart: spanner.Key{"
	q.addToSyntaxStr(&keyRange, q.start)
	keyRange += "\n},\nEnd: spanner.Key{"
	q.addToSyntaxStr(&keyRange, q.end)
	keyRange += fmt.Sprintf("\n},\nKind: %s,\n}", q.KindSyntaxString())
	return keyRange
}
func (q *DeleteQuery) StringSingle() string {
	key := "spanner.Key{"
	q.addToSyntaxStr(&key, q.values)
	key += "\n}"
	return key
}
func (q *DeleteQuery) KindSyntaxString() string {
	switch q.kind.tk {
	case CLOSED_CLOSED_KIND:
		return "spanner.ClosedClosed"
	case CLOSED_OPEN_KIND:
		return "spanner.ClosedOpen"
	case OPEN_CLOSED_KIND:
		return "spanner.OpenClosed"
	case OPEN_OPEN_KIND:
		return "spanner.OpenOpen"
	default: // default is open closed
		return "spanner.OpenClosed"
	}
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
func (q *DeleteQuery) Fields() (out []string) {
	if q.usesKeyRange {
		m := make(map[string]struct{})
		for _, tkn := range q.start {
			if tkn.tk == IDENT_FIELD {
				m[tkn.raw] = struct{}{}
			}
		}
		for _, tkn := range q.end {
			if tkn.tk == IDENT_FIELD {
				m[tkn.raw] = struct{}{}
			}
		}
		for k := range m {
			out = append(out, k)
		}
		return
	}
	for _, tkn := range q.values {
		if tkn.tk == IDENT_FIELD {
			out = append(out, tkn.raw)
		}
	}
	return
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
func (q *DeleteQuery) SetParams(p map[string]string) {
	q.params = p
}
func (q *DeleteQuery) AddParam(key, val string) {
	if q.params == nil {
		q.params = make(map[string]string)
	}
	q.params[key] = val
}

type UpdateQuery struct {
	tokens    []*Token
	cols      []*Token
	values    []*Token
	tableName *Token
	pk        map[Token]*Token
	params    map[string]string
}

func (q *UpdateQuery) String() string {
	update := fmt.Sprintf("spanner.UpdateMap(\"%s\", map[string]interface{}{", q.Table())
	for i, name := range q.cols {
		v := q.values[i]
		var val string
		if v.tk == IDENT_FIELD {
			val = q.params[v.raw]
		} else {
			val = SyntaxStringFromIdent(v)
		}
		update += fmt.Sprintf("\n\t\"%s\": %s,", name.raw, val)
	}
	for k, v := range q.pk {
		var val string
		if v.tk == IDENT_FIELD {
			val = q.params[v.raw]
		} else {
			val = SyntaxStringFromIdent(v)
		}
		update += fmt.Sprintf("\n\t\"%s\": %s,", k.raw, val)
	}
	update += "\n})"
	return update
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
func (q *UpdateQuery) Fields() (out []string) {
	names := make(map[string]struct{})
	for _, tkn := range q.cols {
		names[tkn.raw] = struct{}{}
	}
	for _, tkn := range q.pk {
		names[tkn.raw] = struct{}{}
	}
	for name := range names {
		out = append(out, name)
	}
	return
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
func (q *UpdateQuery) SetParams(p map[string]string) {
	q.params = p
}
func (q *UpdateQuery) AddParam(key, val string) {
	if q.params == nil {
		q.params = make(map[string]string)
	}
	q.params[key] = val
}
func SyntaxStringFromIdent(tkn *Token) string {
	switch tkn.tk {
	case IDENT_STRING:
		return fmt.Sprintf(`"%s"`, tkn.raw)
	case IDENT_INT, IDENT_FLOAT, IDENT_BOOL:
		return fmt.Sprintf(`%s`, tkn.raw)
	default:
		return fmt.Sprintf(`%s`, tkn.raw)
	}
}

type QueryType int

const (
	SELECT_QUERY QueryType = iota
	UPDATE_QUERY
	DELETE_QUERY
	INSERT_QUERY
)
