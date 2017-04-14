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
	"github.com/Sirupsen/logrus"
	"github.com/xwb1989/sqlparser"
	"strconv"
)

type PassedInArgPos int

type Args struct {
	Cur PassedInArgPos
}

func (a *Args) ParseValExpr(expr sqlparser.ValExpr) (interface{}, error) {
	switch value := expr.(type) {
	case sqlparser.StrVal: // a quoted string
		logrus.Debugf("StrVal %+v\n", value)
		return string(value[:]), nil
	case sqlparser.NumVal:
		logrus.Debugf("NumVal %+v\n", value)
		rv, err := strconv.ParseInt(string(value[:]), 10, 64)
		if err != nil {
			rv, err := strconv.ParseFloat(string(value[:]), 64)
			if err != nil {
				return nil, fmt.Errorf("could not parse number value as int or float")
			}
			return rv, nil
		} else {
			return rv, nil
		}
	case sqlparser.ValArg: // ? arg to be supplied by the user
		pa := a.Cur
		a.Cur += 1
		return pa, nil
	case *sqlparser.NullVal:
		return nil, nil
	case *sqlparser.ColName:
		logrus.Debugf("ColName %+v\n", value)
	case sqlparser.ValTuple:
		logrus.Debugf("ValTuple %+v\n", value)
	case *sqlparser.Subquery:
		logrus.Debugf("Subquery %+v\n", value)
	case sqlparser.ListArg:
		logrus.Debugf("ListArg %+v\n", value)
	case *sqlparser.BinaryExpr:
		logrus.Debugf("BinaryExpr %+v\n", value)
	case *sqlparser.UnaryExpr:
		logrus.Debugf("UnaryExpr %+v\n", value)
	case *sqlparser.FuncExpr:
		logrus.Debugf("FuncExpr %+v\n", value)
	case *sqlparser.CaseExpr:
		logrus.Debugf("CaseExpr %+v\n", value)
	}
	return nil, fmt.Errorf("unsupported value expression")
}

type partialArgSlice struct {
	args         []interface{}
	unfilled     map[int]PassedInArgPos
	expectedArgs int
}

type partialArgMap struct {
	args         map[string]interface{}
	unfilled     map[string]PassedInArgPos
	expectedArgs int
}

func newPartialArgSlice() *partialArgSlice {
	return &partialArgSlice{
		args:     make([]interface{}, 0),
		unfilled: make(map[int]PassedInArgPos),
	}
}

func newPartialArgMap() *partialArgMap {
	return &partialArgMap{
		args:     make(map[string]interface{}),
		unfilled: make(map[string]PassedInArgPos),
	}
}

func (p *partialArgSlice) AddArgs(args ...interface{}) {
	for _, arg := range args {
		p.args = append(p.args, arg)
		if ap, ok := arg.(PassedInArgPos); ok {
			p.unfilled[len(p.args)-1] = ap
			p.expectedArgs += 1
		}
	}
}

func (p *partialArgMap) AddArg(key string, val interface{}) {
	p.args[key] = val
	if ap, ok := val.(PassedInArgPos); ok {
		p.expectedArgs += 1
		p.unfilled[key] = ap
	}
}
