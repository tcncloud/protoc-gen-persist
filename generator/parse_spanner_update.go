//
// Copyright 2017, TCN Inc.
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are
// met:
//
//     * Redistributions of source code must retain the above copyright
// notice, this list of conditions and the following disclaimer.
//     * Redistributions in binary form must reproduce the above
// copyright notice, this list of conditions and the following disclaimer
// in the documentation and/or other materials provided with the
// distribution.
//     * Neither the name of TCN Inc. nor the names of its
// contributors may be used to endorse or promote products derived from
// this software without specific prior written permission.
//
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
	"github.com/xwb1989/sqlparser"
)

type updateMap struct {
	updatedVals *partialArgMap
	myArgs      *Args
}

// Spanner updates a particular row by being able to find the row by the primary key.
// Because of this,  updates on primary key fields are not supported.
// A valid sql update query for spanner currently means  listing the the values to to update
// Using SET clauses,  and the a WHERE  clause that specifies the primary key. Multi-column primary
// keys should be joined with AND expressions  in the Where clause.
// Example:
// 	db.Exec("UPDATE test_table1 SET simple_string="hello_world" WHERE id=1")
// Would reference a spanner table with two fields, (id, simple_string)
// with "id" being the primary key
// 	db.Exec("UPDATE test_table2 SET simple_string="hello_world" WHERE id=1 AND other_id=2")
// would reference a spanner table with 3 fields: (id, other_id, simple_string)
// with "id" and "other_id"  being the primary keys
func extractUpdateClause(update *sqlparser.Update) (*partialArgMap, error) {
	myArgs := &Args{}
	updatedVals := newPartialArgMap()
	updateExprs := ([]*sqlparser.UpdateExpr)(update.Exprs)
	for _, updateExpr := range updateExprs {
		if updateExpr.Name == nil {
			return nil, fmt.Errorf("No column name associated with expression %+v", updateExpr.Expr)
		}
		if len(updateExpr.Name.Qualifier) > 0 {
			return nil, fmt.Errorf("qualifiers on column names not allowed for update clause")
		}
		if len(updateExpr.Name.Name) <= 0 {
			return nil, fmt.Errorf("No column name associated with expression %+v", updateExpr.Expr)
		}
		name := string(updateExpr.Name.Name[:])
		arg, err := myArgs.ParseValExpr(updateExpr.Expr)
		if err != nil {
			return nil, err
		}
		updatedVals.AddArg(name, arg)
	}
	upMap := updateMap{updatedVals: updatedVals, myArgs: myArgs}
	err := upMap.walkBoolExpr(update.Where.Expr)
	if err != nil {
		return nil, err
	}
	return upMap.updatedVals, nil
}

func (u *updateMap) walkBoolExpr(boolExpr sqlparser.BoolExpr) error {
	switch expr := boolExpr.(type) {
	case *sqlparser.AndExpr:
		err := u.walkBoolExpr(expr.Left)
		if err != nil {
			return err
		}
		err = u.walkBoolExpr(expr.Right)
		if err != nil {
			return err
		}
	case *sqlparser.ComparisonExpr:
		name, err := u.validColNameFromValExpr(expr.Left)
		if err != nil {
			return err
		}
		if expr.Operator != "=" {
			return fmt.Errorf("only =  operator is supported in update query's Where clause")
		}
		val, err := u.myArgs.ParseValExpr(expr.Right)
		if err != nil {
			return err
		}
		//passed all the tests,  put the value in the map
		u.updatedVals.AddArg(name, val)
	case *sqlparser.NullCheck:
		name, err := u.validColNameFromValExpr(expr.Expr)
		if err != nil {
			return err
		}
		if expr.Operator != "is null" {
			return fmt.Errorf(`only "is null" checks are supported in update query's Where clause`)
		}
		u.updatedVals.AddArg(name, nil)
	default:
		return fmt.Errorf("Unsupported Boolexpr, only support AndExpr, NullCheck, or ComparisonExpr with =")
	}
	return nil
}

func (u *updateMap) validColNameFromValExpr(expr sqlparser.ValExpr) (string, error) {
	col, ok := expr.(*sqlparser.ColName)
	if !ok {
		return "", fmt.Errorf("problem with converting ValExpr to ColName %+v", expr)
	}

	if len(col.Qualifier) > 0 {
		return "", fmt.Errorf("qualifiers not supported in update queries")
	}
	name := string(col.Name[:])
	if _, present := u.updatedVals.args[name]; present {
		return "", fmt.Errorf("update query's where clause cannot have a column that overrides a row being upated")
	}
	return name, nil
}
