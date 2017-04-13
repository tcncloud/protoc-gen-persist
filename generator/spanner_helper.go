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
	"strings"
	"github.com/xwb1989/sqlparser"
	"github.com/Sirupsen/logrus"
)

type QueryArg struct {
	Name         string      // name in the map,
	Value        interface{} // generic value of the argument. If is a field, this will be empty
	IsFieldValue bool        // Whether this refers to a field passed in
	Field        TypeDesc    // if IsFieldValue is true, this will describe the Field
}

type SpannerHelper struct {
	RawQuery string
	Query string
	ParsedQuery sqlparser.Statement
	TableName string
	OptionArguments []string
	IsSelect bool
	IsUpdate bool
	IsInsert bool
	IsDelete bool
	QueryArgs []QueryArg
	InsertCols []string // the column names for insert queries
	Parent *Method
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
	if opts == nil {
		return nil, fmt.Errorf("no options found on proto method")
	}
	args := opts.GetArguments()
	query := opts.GetQuery()
	logrus.Debugf("query: %#v", query)
	pquery, err := sqlparser.Parse(query)
	if err != nil {
		return nil, fmt.Errorf("parsing error in spanner_helper: %s", err)
	}
	// get the fields descriptions to construct query args
	input := p.GetInputTypeStruct()
	fieldsMap := p.GetTypeDescForFieldsInStructSnakeCase(input)


	sh := &SpannerHelper{
		RawQuery: query,
		ParsedQuery: pquery,
		OptionArguments: args,
		Parent: p,
		ProtoFieldDescs: fieldsMap,
	}
	err = sh.Parse()
	if err != nil {
		return nil, err
	}
	return sh, nil
}

func (sh *SpannerHelper) Parse() error {
	// parse our query
	switch pq := sh.ParsedQuery.(type) {
	case *sqlparser.Select:
		sh.IsSelect = true
		spl := strings.Split(sh.RawQuery, "?")
		var updatedQuery string

		if len(sh.OptionArguments) != len(spl) - 1 {
			errStr := "err parsing spanner query: not correct number of option arguments"
			errStr += " for method: %s of service: %s  want: %d have: %d"
			return fmt.Errorf(errStr, sh.Parent.GetName(), sh.Parent.Service.GetName(), len(spl) - 1, len(sh.OptionArguments))
		}
		for i := 0; i < len(spl)-1; i++ {
			name := fmt.Sprintf("@%d", i)
			field := sh.ProtoFieldDescs[sh.OptionArguments[i]]
			qa := QueryArg{
				Name: name,
				IsFieldValue: true,
				Field: field,
			}
			sh.QueryArgs = append(sh.QueryArgs, qa)
			updatedQuery += (spl[i] + name)
		}
		updatedQuery += spl[len(spl)-1]
		sh.Query = updatedQuery
	case *sqlparser.Insert:
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
		for _, arg := range pas.args {
			var qa QueryArg
			if ap, ok := arg.(PassedInArgPos); ok {
				index := int(ap)
				argName := sh.OptionArguments[index]
				qa = QueryArg{
					IsFieldValue: true,
					Field: sh.ProtoFieldDescs[argName],
				}
			} else {
				qa = QueryArg{
					Value: arg,
					IsFieldValue: false,
				}
			}
			sh.QueryArgs = append(sh.QueryArgs, qa)
		}
		sh.InsertCols = cols
		sh.TableName = table
	case *sqlparser.Delete:
		sh.IsUpdate = true
		table, err := extractIUDTableName(pq)
		if err != nil {
			return err
		}
		sh.TableName = table
	case *sqlparser.Update:
		sh.IsDelete = true
		table, err := extractIUDTableName(pq)
		if err != nil {
			return err
		}
		sh.TableName = table
	}
	return nil
}

func (sh *SpannerHelper) InsertColsAsString() string {
	return fmt.Sprintf("%#v", sh.InsertCols)
}


