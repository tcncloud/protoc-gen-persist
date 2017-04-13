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
	"database/sql/driver"
	"fmt"
	"github.com/Sirupsen/logrus"
	"cloud.google.com/go/spanner"
	"github.com/xwb1989/sqlparser"
)

// because spanner has multiple primary keys support,  EVERY field
// found in the query is assumed to be a primary key.  It will build the spanner.Key with the fields in the
// query in the order they are discovered. For Example if i have two queries:
//    q1 = "DELETE FROM test_table WHERE id = 1 AND simple_string="test_string"
//    q2 = "DELETE FROM test_table WHERE simple_string="test_string" AND id = 1
//    q1 would produce key: { 1, "test_string" }
//    q2 would produce key: { "test_string", 2 }
// Try to construct your queries with ANDs  instead of ORs. Because different fields are
// interpreted as primary keys,  it gets too difficult to parse what is meant by queries like:
//    q1 = "DELETE FROM test_table WHERE (id < 1 OR id >= 10) AND simple_string = "test_string"
// it might be possible in the future to parse the meaning of statements like this,  but for now it was
// easier to just drop support of statements for OR expressions
//   Other Rules:
// - NOT expressions are not supported, It is not possible to tell a spanner key what "not"  means.
// - currently only one key range, per primary key is permitted.  Just use two queries. ex.
//    not permitted: DELETE FROM test_table WHERE id > 1 AND id < 10 AND id > 20 AND id < 100
// - Does not support cross table queries
func extractSpannerKeyFromDelete(del *sqlparser.Delete) (*MergableKeyRange, error) {
	where := del.Where
	if where == nil {
		return nil, fmt.Errorf("Must include a where clause that contain primary keys in delete statement")
	}
	myArgs := &Args{}
	logrus.Debugf("where type: %+v\n", where.Type)
	aKeySet := &AwareKeySet{
		Args:     myArgs,
		Keys:     make(map[string]*Key),
		KeyOrder: make([]string, 0),
	}
	err := aKeySet.walkBoolExpr(where.Expr)
	if err != nil {
		return nil, err
	}
	return aKeySet.packKeySet()
}

type Key struct {
	Name       string
	LowerValue interface{}
	UpperValue interface{}
	LowerOpen  bool
	UpperOpen  bool
	HaveLower  bool
	HaveUpper  bool
}

type AwareKeySet struct {
	Keys     map[string]*Key
	KeyOrder []string
	Args     *Args
}

type MergableKeyRange struct {
	Start     *partialArgSlice
	End       *partialArgSlice
	LowerOpen bool
	UpperOpen bool
	HaveLower bool
	HaveUpper bool
}

// all lower bounds are turned into a key together.
// all upper bounds are turned into a key together.
// it is expected that all fields in a query belong together
func (a *AwareKeySet) packKeySet() (*MergableKeyRange, error) {
	var prev *MergableKeyRange
	//makes sure all we dont have holes in our key ranges,  that is undefined behaviour
	for i := len(a.KeyOrder) - 1; i > 0; i-- { // dont check before the first elem
		me := a.Keys[a.KeyOrder[i]]
		keyBeforeMe := a.Keys[a.KeyOrder[i-1]]
		if me.HaveLower {
			if !keyBeforeMe.HaveLower {
				return nil, fmt.Errorf("cannot have a lower bound on a key range without defining all higher priority lower bounds")
			}
		}
		if me.HaveUpper {
			if !keyBeforeMe.HaveUpper {
				return nil, fmt.Errorf("cannot have a upper bound on a key range without defining all higher priority upper bounds")
			}
		}
	}
	for _, k := range a.KeyOrder {
		key := a.Keys[k]
		if prev == nil {
			prev = &MergableKeyRange{Start: newPartialArgSlice(), End: newPartialArgSlice()}
			prev.fromKey(key)
		} else {
			logrus.Debugf("key that will populate prev: %#v\n\n", key)
			err := prev.mergeKey(key)
			logrus.Debugf("merged prev with m %#v\n\n", prev)
			if err != nil {
				return nil, err
			}
		}
	}
	return prev, nil
}

func (k1 *MergableKeyRange) fromKey(key *Key) {
	if key == nil {
		return
	}
	k1.LowerOpen = key.LowerOpen
	k1.UpperOpen = key.UpperOpen
	k1.HaveLower = key.HaveLower
	k1.HaveUpper = key.HaveUpper
	k1.Start.AddArgs(key.LowerValue)
	k1.End.AddArgs(key.UpperValue)
}

func (k *MergableKeyRange) ToKeyRange(args []driver.Value) (*spanner.KeyRange, error) {
	low := k.LowerOpen
	up := k.UpperOpen

	var kind spanner.KeyRangeKind

	if low && up {
		kind = spanner.OpenOpen
	} else if low && !up {
		kind = spanner.OpenClosed
	} else if !low && up {
		kind = spanner.ClosedOpen
	} else {
		kind = spanner.ClosedClosed
	}
	start, err := k.Start.GetFilledArgs(args)
	if err != nil {
		return nil, err
	}
	end, err := k.End.GetFilledArgs(args)
	if err != nil {
		return nil, err
	}
	return &spanner.KeyRange{
		Start: start,
		End:   end,
		Kind:  kind,
	}, nil
}

func (k1 *MergableKeyRange) mergeKey(k2 *Key) error {
	logrus.Debug("\nmerging into k1: %#v\n  k2: %#v\n\n", k1, k2)
	if k2.HaveLower {
		if k1.LowerOpen != k2.LowerOpen {
			return fmt.Errorf("Kinds in ranges must all match")
		}
		k1.Start.AddArgs(k2.LowerValue)
	}
	if k2.HaveUpper {
		if k1.UpperOpen != k2.UpperOpen {
			return fmt.Errorf("Kinds in ranges must all match")
		}
		k1.End.AddArgs(k2.UpperValue)
	}
	return nil
}

func (a *AwareKeySet) addKeyFromValExpr(valExpr sqlparser.ValExpr) (*Key, error) {
	col, ok := valExpr.(*sqlparser.ColName)
	if !ok {
		return nil, fmt.Errorf("not a valid column name")
	}
	if len(col.Qualifier) != 0 {
		return nil, fmt.Errorf("qualifiers not allowed")
	}
	keyName := string(col.Name[:])
	if a.Keys[keyName] == nil {
		a.KeyOrder = append(a.KeyOrder, keyName)
		a.Keys[keyName] = &Key{Name: keyName}
	}
	return a.Keys[keyName], nil
}

func (a *AwareKeySet) walkBoolExpr(boolExpr sqlparser.BoolExpr) error {
	switch expr := boolExpr.(type) {
	case *sqlparser.AndExpr:
		logrus.Debugf("AndExpr %#v\n", expr)
		err := a.walkBoolExpr(expr.Left)
		if err != nil {
			return err
		}
		err = a.walkBoolExpr(expr.Right)
		if err != nil {
			return err
		}
		return nil
	case *sqlparser.OrExpr:
		logrus.Debugf("OrExpr %#v\n", expr)
		return fmt.Errorf("Or Expressions are not currently supported")
	case *sqlparser.ParenBoolExpr:
		logrus.Debugf("ParenBoolExpr %#v\n", expr)
	case *sqlparser.ComparisonExpr:
		logrus.Debugf("ComparisonExpr %#v\n", expr)
		myKey, err := a.addKeyFromValExpr(expr.Left)
		if err != nil {
			return err
		}
		val, err := a.Args.ParseValExpr(expr.Right)
		if err != nil {
			return err
		}
		logrus.Debugf("OPERTATOR %#v\n", expr.Operator)
		switch expr.Operator {
		case "=":
			myKey.LowerValue = val
			myKey.UpperValue = val
			myKey.LowerOpen = false
			myKey.UpperOpen = false
			myKey.HaveUpper = true
			myKey.HaveLower = true
			return nil
		case ">":
			myKey.LowerValue = val
			myKey.LowerOpen = true
			myKey.HaveLower = true
			return nil
		case "<":
			myKey.UpperValue = val
			myKey.UpperOpen = true
			myKey.HaveUpper = true
			return nil
		case ">=":
			myKey.LowerValue = val
			myKey.LowerOpen = false
			myKey.HaveLower = true
			return nil
		case "<=":
			myKey.UpperValue = val
			myKey.UpperOpen = false
			myKey.HaveUpper = true
			return nil
		case "!=":
			return fmt.Errorf("!= comparisons are not supported")
		case "not in", "in":
			return fmt.Errorf("in, and not in  comparisons are not supported")
		default:
			return fmt.Errorf("%#v  is not a supported operator", expr.Operator)
		}
	case *sqlparser.RangeCond:
		logrus.Debugf("RangeCond %#v\n", expr)
		myKey, err := a.addKeyFromValExpr(expr.Left)
		if err != nil {
			return err
		}
		from, err := a.Args.ParseValExpr(expr.From)
		if err != nil {
			return err
		}
		to, err := a.Args.ParseValExpr(expr.To)
		if err != nil {
			return err
		}
		switch expr.Operator {
		case "between":
			myKey.LowerValue = from
			myKey.LowerOpen = true
			myKey.UpperValue = to
			myKey.UpperOpen = true
		case "not between":
			return fmt.Errorf("not between operator is not supported")
		}
	case *sqlparser.ExistsExpr:
		logrus.Debugf("ExistsExpr %#v\n", expr)
		return fmt.Errorf("Exists Expressions are not supported")
	}

	return fmt.Errorf("not a boolexpr %#v\n", boolExpr)
}
