// Copyright 2017, TCN Inc.
// All rights reserved.

// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are
// met:

//     * Redistributions of source code must retain the above copyright
// notice, this list of conditions and the following disclaimer.
//     * Redistributions in binary form must reproduce the above
// copyright notice, this list of conditions and the following disclaimer
// in the documentation and/or other materials provided with the
// distribution.
//     * Neither the name of TCN Inc. nor the names of its
// contributors may be used to endorse or promote products derived from
// this software without specific prior written permission.

// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
// A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
// OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
// LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package generator

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/tcncloud/protoc-gen-persist/generator/delete_parser"
	"github.com/xwb1989/sqlparser"
	"strings"
)

type QueryArg struct {
	Name         string      // name in the map,
	Value        interface{} // generic value of the argument. If is a field, this will be empty
	IsFieldValue bool        // Whether this refers to a field passed in
	IsValue      bool        // If the Value field is set, this will be true
	Field        TypeDesc    // if IsFieldValue is true, this will describe the Field
}

type KeyRangeDesc struct {
	Start []QueryArg
	End   []QueryArg
	Kind  string // a string of of a spanner.KeyRangeKind (ClosedOpen, ClosedClosed ex.)
	Table string
}

type SpannerHelper struct {
	RawQuery        string
	Query           string
	ParsedQuery     sqlparser.Statement
	DeleteParser    *delete_parser.Parser
	TableName       string
	OptionArguments []string
	IsSelect        bool
	IsUpdate        bool
	IsInsert        bool
	IsDelete        bool
	QueryArgs       []QueryArg
	KeyRangeDesc    *KeyRangeDesc // used for delete queries, will be set if IsDelete is true
	InsertCols      []string      // the column names for insert queries
	Parent          *Method
	ProtoFieldDescs map[string]TypeDesc
}

func (sh *SpannerHelper) String() string {
	if sh != nil {
		return fmt.Sprintf("SpannerHelper\n\tQuery: %s\n\tIsSelect: %t\n\tIsUpdate: %t\n\tIsInsert: %t\n\tIsDelete: %t\n\n",
			sh.Query, sh.IsSelect, sh.IsUpdate, sh.IsInsert, sh.IsDelete)
	}
	return "<nil>"
}

func NewSpannerHelper(p *Method) (*SpannerHelper, error) {
	// get the query, and parse it
	opts := p.GetMethodOption()
	args := opts.GetArguments()
	query := p.GetQuery()
	logrus.Debugf("query: %#v", query)
	// get the fields descriptions to construct query args
	input := p.GetInputTypeStruct()
	fieldsMap := p.GetTypeDescForFieldsInStructSnakeCase(input)

	sh := &SpannerHelper{
		RawQuery:        query,
		OptionArguments: args,
		Parent:          p,
		ProtoFieldDescs: fieldsMap,
	}
	err := sh.Parse()
	if err != nil {
		return nil, err
	}
	return sh, nil
}

func (sh *SpannerHelper) Parse() error {
	// parse our RAW query
	// if this is a delete query, (starts with DELETE)  then create a delete parser"
	// otherwise use the sqlparser
	if strings.HasPrefix(sh.RawQuery, "DELETE") {
		sh.DeleteParser = delete_parser.NewParser(sh.RawQuery)
		pdq, err := sh.DeleteParser.Expr()
		if err != nil {
			return err
		}
		return sh.HandleDelete(pdq)
	} else if strings.HasPrefix(strings.ToUpper(sh.RawQuery), "SELECT") {
		return sh.ParseSelect()
	} else {
		parsed, err := sqlparser.Parse(sh.RawQuery)
		if err != nil {
			logrus.Debugf("got error trying to parse spanner query: %s, err: %s\n", sh.RawQuery, err)
			return err
		}
		sh.ParsedQuery = parsed
		switch pq := sh.ParsedQuery.(type) {
		case *sqlparser.Insert:
			return sh.ParseInsert(pq)
		case *sqlparser.Delete:
			return fmt.Errorf("delete querys must start with DELETE, instead saw: %s", sh.RawQuery)
		case *sqlparser.Update:
			return sh.ParseUpdate(pq)
		default:
			return fmt.Errorf("not a query we can parse")
		}
	}
}

func (sh *SpannerHelper) InsertColsAsString() string {
	return fmt.Sprintf("%#v", sh.InsertCols)
}

func (sh *SpannerHelper) PopulateArgSlice(slice []interface{}) []QueryArg {
	qas := make([]QueryArg, len(slice))
	for i, arg := range slice {
		var qa QueryArg
		if ap, ok := arg.(PassedInArgPos); ok {
			index := int(ap)
			argName := sh.OptionArguments[index]
			qa = QueryArg{
				IsFieldValue: true,
				IsValue:      false,
				Field:        sh.ProtoFieldDescs[argName],
			}
		} else {
			qa = QueryArg{
				Value:        fmt.Sprintf("%#v", arg),
				IsFieldValue: false,
				IsValue:      true,
			}
			if arg != nil {
				qa.Value = fmt.Sprintf("%#v", arg)
			} else {
				qa.Value = "nil"
			}
		}
		qas[i] = qa
	}
	return qas
}

func (sh *SpannerHelper) ParseInsert(pq *sqlparser.Insert) error {
	sh.IsInsert = true
	cols, err := extractInsertColumns(pq)
	if err != nil {
		return err
	}
	table, err := extractIUDTableName(pq)
	if err != nil {
		return err
	}
	pas, err := prepareInsertValues(pq)
	if err != nil {
		return err
	}
	qas := sh.PopulateArgSlice(pas.args)
	if err != nil {
		return err
	}
	sh.QueryArgs = qas
	sh.InsertCols = cols
	sh.TableName = table
	return nil
}
func (sh *SpannerHelper) ParseSelect() error {
	sh.IsSelect = true
	spl := strings.Split(sh.RawQuery, "?")
	var updatedQuery string

	if len(sh.OptionArguments) != len(spl)-1 {
		errStr := "err parsing spanner query: not correct number of option arguments"
		errStr += " for method: %s of service: %s  want: %d have: %d"
		return fmt.Errorf(errStr, sh.Parent.GetName(), sh.Parent.Service.GetName(), len(spl)-1, len(sh.OptionArguments))
	}
	for i := 0; i < len(spl)-1; i++ {
		name := fmt.Sprintf(` @string%d`, i)
		index := fmt.Sprintf("string%d", i)
		field := sh.ProtoFieldDescs[sh.OptionArguments[i]]
		qa := QueryArg{
			Name:         index,
			IsFieldValue: true,
			IsValue:      false,
			Field:        field,
		}
		sh.QueryArgs = append(sh.QueryArgs, qa)
		updatedQuery += (spl[i] + name)
	}
	updatedQuery += spl[len(spl)-1]
	sh.Query = updatedQuery
	return nil
}

func (sh *SpannerHelper) HandleDelete(pdq *delete_parser.ParsedKeyRange) error {
	sh.IsDelete = true
	sh.TableName = pdq.Table

	start := make([]QueryArg, len(pdq.Start))
	end := make([]QueryArg, len(pdq.End))
	index := 0
	if err := sh.PopulateDeleteSlice(pdq.Start, start, &index); err != nil {
		return err
	}
	if err := sh.PopulateDeleteSlice(pdq.End, end, &index); err != nil {
		return err
	}
	sh.KeyRangeDesc = &KeyRangeDesc{
		Start: start,
		End:   end,
		Kind:  pdq.Kind,
	}
	return nil
}

func (sh *SpannerHelper) PopulateDeleteSlice(source []*delete_parser.Token, dest []QueryArg, index *int) error {
	for i, tok := range source {
		var qa QueryArg
		if tok.Type == "?" {
			if *index >= len(sh.OptionArguments) {
				return fmt.Errorf("too many ? for arguments")
			}
			argName := sh.OptionArguments[*index]
			*index += 1
			qa = QueryArg{
				IsValue:      false,
				IsFieldValue: true,
				Field:        sh.ProtoFieldDescs[argName],
			}
		} else {
			qa = QueryArg{
				IsValue:      true,
				IsFieldValue: false,
				Value:        tok.Value,
			}
		}
		dest[i] = qa
	}
	return nil
}

func (sh *SpannerHelper) ParseUpdate(pq *sqlparser.Update) error {
	sh.IsUpdate = true
	table, err := extractIUDTableName(pq)
	if err != nil {
		return err
	}
	pam, err := extractUpdateClause(pq)
	if err != nil {
		return err
	}
	for key, arg := range pam.args {
		var qa QueryArg
		if ap, ok := arg.(PassedInArgPos); ok {
			index := int(ap)
			argName := sh.OptionArguments[index]
			qa = QueryArg{
				Name:         key,
				IsFieldValue: true,
				IsValue:      false,
				Field:        sh.ProtoFieldDescs[argName],
			}
		} else {
			qa = QueryArg{
				Name:         key,
				Value:        fmt.Sprintf("%#v", arg),
				IsFieldValue: false,
				IsValue:      true,
			}
		}
		sh.QueryArgs = append(sh.QueryArgs, qa)
	}
	sh.TableName = table
	return nil
}
