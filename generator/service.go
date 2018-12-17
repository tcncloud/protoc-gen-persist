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
	"regexp"
	"strings"

	"github.com/golang/protobuf/proto"
	desc "github.com/golang/protobuf/protoc-gen-go/descriptor"
	_gen "github.com/golang/protobuf/protoc-gen-go/generator"
	"github.com/sirupsen/logrus"
	"github.com/tcncloud/protoc-gen-persist/persist"
)

type Service struct {
	Desc       *desc.ServiceDescriptorProto
	Package    string // protobuf package
	File       *FileStruct
	AllStructs *StructList
}

func (s *Service) GetName() string {
	return s.Desc.GetName()
}

func (s *Service) GetQueriesOption() *persist.QueryOpts {
	if s.Desc.Options != nil && proto.HasExtension(s.Desc.Options, persist.E_Ql) {
		ext, err := proto.GetExtension(s.Desc.Options, persist.E_Ql)
		if err == nil {
			return ext.(*persist.QueryOpts)
		}
	}
	return nil
}
func (s *Service) GetTypeMapping() *persist.TypeMapping {
	if s.Desc.Options != nil && proto.HasExtension(s.Desc.Options, persist.E_Mapping) {
		ext, err := proto.GetExtension(s.Desc.Options, persist.E_Mapping)
		if err == nil {
			return ext.(*persist.TypeMapping)
		}
	}
	return nil
}

func (s *Service) GetServiceType() *persist.PersistenceOptions {
	if s.Desc.Options != nil && proto.HasExtension(s.Desc.Options, persist.E_ServiceType) {
		ext, err := proto.GetExtension(s.Desc.Options, persist.E_ServiceType)
		if err == nil {
			return ext.(*persist.PersistenceOptions)
		}
	}
	return nil
}

func (s *Service) IsSQL() bool {
	if p := s.GetServiceType(); p != nil {
		if *p == persist.PersistenceOptions_SQL {
			return true
		}
	}
	return false
}

func (s *Service) IsSpanner() bool {
	if p := s.GetServiceType(); p != nil {
		if *p == persist.PersistenceOptions_SPANNER {
			return true
		}
	}
	return false
}

func (s *Service) GetUndoctoredQueryByName(queryName string) (*persist.QLImpl, error) {
	qopts := s.GetQueriesOption()
	queries := qopts.GetQueries()

	for _, query := range queries {
		if query.GetName() == queryName {
			return query, nil
		}
	}
	return nil, fmt.Errorf("query not found with name: %s on service: %s", queryName, s.GetName())
}

type Services []*Service

// we are a persist service if we have persist options. meaning we are either spanner
// or sql
func (s Services) HasPersistService() bool {
	for _, serv := range s {
		if serv.IsSQL() || serv.IsSpanner() {
			return true
		}
	}
	return false
}

func (s *Services) AddService(pkg string, desc *desc.ServiceDescriptorProto, allStructs *StructList, file *FileStruct) *Service {
	ret := &Service{
		Package:    pkg,
		Desc:       desc,
		AllStructs: allStructs,
		File:       file,
	}
	logrus.Debugf("created a service: %s", ret)
	*s = append(*s, ret)
	return ret
}

type QueryProtoOpts struct {
	query     *persist.QLImpl
	inMsg     *Struct
	outMsg    *Struct
	inFields  []*desc.FieldDescriptorProto
	outFields []*desc.FieldDescriptorProto
}

func NewQueryProtoOpts(qopt *persist.QLImpl, all *StructList) (*QueryProtoOpts, error) {
	in := all.GetStructByProtoName(qopt.GetIn())
	out := all.GetStructByProtoName(qopt.GetOut())
	if !in.IsMessage || !out.IsMessage {
		return nil, fmt.Errorf("in/out option must be proto messages for query: '%s'", qopt.GetName())
	}
	inFields, _ := in.GetFieldDescriptorsIfMessage()
	outFields, _ := out.GetFieldDescriptorsIfMessage()

	return &QueryProtoOpts{
		query:     qopt,
		inMsg:     in,
		outMsg:    out,
		inFields:  inFields,
		outFields: outFields,
	}, nil
}

type TypeMappingProtoOpts struct {
	tm *persist.TypeMapping_TypeDescriptor
}

func NewTypeMappingProtoOpts(opt *persist.TypeMapping_TypeDescriptor, all *StructList) (*TypeMappingProtoOpts, error) {
	return &TypeMappingProtoOpts{tm: opt}, nil
}

type MethodProtoOpts struct {
	method    *desc.MethodDescriptorProto
	option    *persist.MOpts
	query     *persist.QLImpl
	inStruct  *Struct
	outStruct *Struct
	inMsg     *desc.DescriptorProto
	outMsg    *desc.DescriptorProto
	inFields  []*desc.FieldDescriptorProto
	outFields []*desc.FieldDescriptorProto
}

func NewMethodProtoOpts(opt *desc.MethodDescriptorProto, all *StructList) (*MethodProtoOpts, error) {
	in := all.GetStructByProtoName(opt.GetInputType())
	out := all.GetStructByProtoName(opt.GetOutputType())

	if !in.IsMessage || !out.IsMessage {
		return nil, fmt.Errorf("in/out option must be proto messages for query: '%s'", opt.GetName())
	}
	inFields, _ := in.GetFieldDescriptorsIfMessage()
	outFields, _ := out.GetFieldDescriptorsIfMessage()

	var option *persist.MOpts

	if opt.Options != nil && proto.HasExtension(opt.Options, persist.E_Opts) {
		ext, err := proto.GetExtension(opt.Options, persist.E_Opts)
		if err == nil {
			option = ext.(*persist.MOpts)
		}
	}

	return &MethodProtoOpts{
		method:    opt,
		option:    option,
		inMsg:     in.MsgDesc,
		outMsg:    out.MsgDesc,
		inStruct:  in,
		outStruct: out,
		inFields:  inFields,
		outFields: outFields,
	}, nil
}
func WriteQueries(p *Printer, s *Service) error {
	m := Matcher(s)
	sName := s.GetName()
	camelQ := func(q *QueryProtoOpts) string {
		return _gen.CamelCase(q.query.GetName())
	}
	qname := func(q *QueryProtoOpts) string {
		return q.query.GetName()
	}
	queryAndFields := func(q *QueryProtoOpts) (string, []string) {
		orig := strings.Join(q.query.GetQuery(), " ")
		pmStrat := q.query.GetPmStrategy()
		nextParamMarker := func() func(string) string {
			var count int
			return func(req string) string {
				if pmStrat == "$" {
					count++
					return fmt.Sprintf("$%d", count)
				} else if pmStrat == "?" {
					return "?"
				}
				return req
			}
		}()

		newQuery := ""
		r := regexp.MustCompile("@[a-zA-Z0-9_]*")
		potentialFieldNames := r.FindAllString(orig, -1)
		fieldsMap := make(map[string]bool)
		for _, v := range q.inFields {
			fieldsMap[v.GetName()] = true
		}
		params := make([]string, 0)
		for _, pf := range potentialFieldNames {
			start := strings.Index(orig, pf)
			stop := start + len(pf)
			// index into the map (removing the "@")
			exists := fieldsMap[pf[1:]]
			// eat up to the field name
			newQuery += orig[:start]
			if !exists { // it was just part of the query, not a field on input
				newQuery += pf
			} else { // it was a field, mark it
				newQuery += nextParamMarker(pf)
				params = append(params, pf[1:])
			}
			// remove the already written stuff
			orig = orig[stop:]
		}
		newQuery += orig
		return newQuery, params
	}
	qstring := func(q *QueryProtoOpts) string {
		res, _ := queryAndFields(q)
		return res
	}
	resultOrRows := func(q *QueryProtoOpts) string {
		if len(q.outFields) == 0 {
			return "result"
		}
		return "rows"
	}
	qmethod := func(q *QueryProtoOpts) string {
		if len(q.outFields) == 0 {
			return "Exec"
		}
		return "Query"
	}
	execParams := func(q *QueryProtoOpts) string {
		printer := &Printer{}
		paramStrings := make(map[string]string)
		// basic mappping
		m.EachQueryIn(func(f *desc.FieldDescriptorProto, q *QueryProtoOpts) {
			paramStrings[f.GetName()] = P(`func() (out interface{}) {
                out = x.Get`, _gen.CamelCase(f.GetName()), `()
                return
            }(),
            `)
		}, m.MatchQuery(q))
		// all the proto message types
		// will overwrite paramStrings if the type is a message
		m.EachQueryIn(func(f *desc.FieldDescriptorProto, q *QueryProtoOpts) {
			paramStrings[f.GetName()] = P(`func() (out interface{}) {
                raw, err := proto.Marshal(x.Get`, _gen.CamelCase(f.GetName()), `())
                if err != nil {
                    setupErr = err
                }
                out = raw
                return
            }(),
            `)
		}, m.MatchQuery(q), m.QueryFieldIsMessage)
		// will overwrite paramStrings if type is mapped
		m.EachQueryIn(func(f *desc.FieldDescriptorProto, q *QueryProtoOpts) {
			m.EachTM(func(tm *TypeMappingProtoOpts) {
				_, titled := getGoNamesForTypeMapping(tm.tm, s.File)
				paramStrings[f.GetName()] = P(`func() (out interface{}) {
                    mapper := this.opts.MAPPINGS.`, titled, `()
                    out = mapper.ToSql(x.Get`, _gen.CamelCase(f.GetName()), `())
                    return
                }(),
                `)
			}, m.MatchTypeMapping(f))
		}, m.MatchQuery(q), m.QueryFieldIsMapped)

		_, paramOrdering := queryAndFields(q)
		for _, paramName := range paramOrdering {
			printer.Q(paramStrings[paramName])
		}
		return printer.String()
	}

	p.Q("type ", sName, "_QueryOpts struct {\n")
	p.Q("MAPPINGS ", sName, "_TypeMappings\n")
	p.Q("db Runnable\n")
	p.Q("ctx context.Context\n")
	p.Q("}\n")
	p.Q("// Default", sName, "QueryOpts return the default options to be used with ", sName, "_Queries\n")
	p.Q("func Default", sName, "QueryOpts(db Runnable) ", sName, "_QueryOpts {\n")
	p.Q("return ", sName, "_QueryOpts{\n")
	p.Q("db: db,\n")
	p.Q("}\n")
	p.Q("}\n") // End DefaultQueryOpts

	p.Q("// ", sName, "_Queries holds all the queries found the proto service option as methods\n")
	p.Q("type ", sName, "_Queries struct {\n")
	p.Q("opts ", sName, "_QueryOpts\n")
	p.Q("}\n")

	p.Q(`// `, sName, `PersistQueries returns all the known 'SQL' queires for the '`, sName, `' service.
    func `, sName, `PersistQueries(db Runnable, opts ...`, sName, `_QueryOpts) *`, sName, `_Queries {
        var myOpts `, sName, `_QueryOpts
        if len(opts) > 0 {
            myOpts = opts[0]
        } else {
            myOpts = Default`, sName, `QueryOpts(db)
        }
        myOpts.db = db
        return &`, sName, `_Queries{
            opts: myOpts,
        }
    }
    `)
	m.EachQuery(func(q *QueryProtoOpts) {
		p.Q(`// `, camelQ(q), `Query returns a new struct wrapping the current `, sName, `_QueryOpts
        // that will perform '`, sName, `' services '`, qname(q), `' on the database
        // when executed
        func (this *`, sName, `_Queries) `, camelQ(q), `Query(ctx context.Context) *`, sName, `_`, camelQ(q), `Query {
            return &`, sName, `_`, camelQ(q), `Query{
                opts: `, sName, `_QueryOpts{
                    MAPPINGS: this.opts.MAPPINGS,
                    db:       this.opts.db,
                    ctx:      ctx,
                },
            }
        }
        type `, sName, `_`, camelQ(q), `Query struct {
            opts `, sName, `_QueryOpts
        }

        func (this *`, sName, `_`, camelQ(q), `Query) QueryInTypeUser()  {}
        func (this *`, sName, `_`, camelQ(q), `Query) QueryOutTypeUser() {}

        // Executes the query with parameters retrieved from x
        func (this *`, sName, `_`, camelQ(q), `Query) Execute(x `, sName, `_`, camelQ(q), `In) *`, sName, `_`, camelQ(q), `Iter {
            var setupErr error
            params := []interface{}{
            `, execParams(q), `
            }
            result := &`, sName, `_`, camelQ(q), `Iter{
                tm: this.opts.MAPPINGS,
                ctx: this.opts.ctx,
            }
            if setupErr != nil {
                result.err = setupErr
                return result
            }
            result.`, resultOrRows(q), `, result.err = this.opts.db.`, qmethod(q), `Context(this.opts.ctx, "`, qstring(q), `", params...)

            return result
        }
        `)
	})

	return nil
}
func WriteHooks(p *Printer, s *Service) error {
	sName := s.GetName()
	m := Matcher(s)
	inName := func(opt *MethodProtoOpts) string {
		return convertedMsgTypeByProtoName(opt.inStruct.GetProtoName(), s.File)
	}
	outName := func(opt *MethodProtoOpts) string {
		return convertedMsgTypeByProtoName(opt.outStruct.GetProtoName(), s.File)
	}
	name := func(mopt *MethodProtoOpts) string {
		return mopt.method.GetName()
	}

	p.Q("type ", sName, "_Hooks interface {\n")
	m.EachMethod(func(m *MethodProtoOpts) {
		p.Q(name(m), "BeforeHook(context.Context, *", inName(m), ") ([]*", outName(m), ", error)\n")
	}, func(m *MethodProtoOpts) bool {
		return false
	}, m.BeforeHook, m.ServerStreaming)
	m.EachMethod(func(m *MethodProtoOpts) {
		p.Q(name(m), "BeforeHook(context.Context, *", inName(m), ") (*", outName(m), ", error)\n")
	}, m.BeforeHook)
	m.EachMethod(func(m *MethodProtoOpts) {
		p.Q(name(m), "AfterHook(context.Context, *", inName(m), ",*", outName(m), ") error\n")
	}, m.AfterHook)
	p.Q("}\n")
	return m.Err()
}

func WriteTypeMappings(p *Printer, s *Service) error {
	sName := s.GetName()
	// TODO google's WKT protobufs probably don't need the package prefix
	p.Q("type ", sName, "_TypeMappings interface{\n")
	tms := s.GetTypeMapping().GetTypes()
	for _, tm := range tms {
		// TODO implement these interfaces
		_, titled := getGoNamesForTypeMapping(tm, s.File)
		// p.Q(titled, "() ", sName, titled, "MappingImpl\n")
		p.Q(titled, "() ", sName, titled, "MappingImpl\n")
	}
	p.Q("}\n")

	for _, tm := range tms {
		name, titled := getGoNamesForTypeMapping(tm, s.File)
		p.Q("type ", sName, titled, "MappingImpl interface {\n")
		p.Q(fmt.Sprintf(`
            ToProto(**%[1]s) error
            Empty() %[3]s%[2]sMappingImpl
            ToSql(*%[1]s) sql.Scanner
            sql.Scanner
            driver.Valuer
        `, name, titled, sName))
		p.Q("}\n")
	}

	return nil
}

func WriteIters(p *Printer, s *Service) (outErr error) {
	m := Matcher(s)
	sName := s.GetName()
	camelQ := func(q *QueryProtoOpts) string {
		return _gen.CamelCase(q.query.GetName())
	}
	fName := func(f *desc.FieldDescriptorProto) string {
		return f.GetName()
	}
	camelF := func(f *desc.FieldDescriptorProto) string {
		return _gen.CamelCase(f.GetName())
	}
	inName := func(opt *QueryProtoOpts) string {
		return convertedMsgTypeByProtoName(opt.inMsg.GetProtoName(), s.File)
	}
	outName := func(opt *QueryProtoOpts) string {
		return convertedMsgTypeByProtoName(opt.outMsg.GetProtoName(), s.File)
	}
	mustDefaultMapping := func(f *desc.FieldDescriptorProto) string {
		typ, err := defaultMapping(f, s.File)
		if err != nil {
			outErr = err
		}
		return typ
	}
	mustDefaultMappingNoStar := func(f *desc.FieldDescriptorProto) string {
		return strings.Map(func(r rune) rune {
			if r == '*' {
				return -1
			}
			return r
		}, mustDefaultMapping(f))
	}
	colswitch := func(opt *QueryProtoOpts) string {
		cases := make(map[string]string)
		// message case
		m.EachQueryOut(func(f *desc.FieldDescriptorProto, q *QueryProtoOpts) {
			cases[fName(f)] = P(`case "`, fName(f), `":
                r, ok := (*scanned[i].i).([]byte)
                if !ok {
                    return &`, sName, `_`, camelQ(q), `Row{err: fmt.Errorf("cant convert db column `, fName(f), ` to protobuf go type *`, mustDefaultMappingNoStar(f), `")}, true
                }
                var converted = new(`, mustDefaultMappingNoStar(f), `)
                if err := proto.Unmarshal(r, converted); err != nil {
                    return &`, sName, `_`, camelQ(q), `Row{err: err}, true
                }
                res.`, camelF(f), ` = converted
            `)
		}, m.MatchQuery(opt))

		// fits case
		m.EachQueryOut(func(f *desc.FieldDescriptorProto, q *QueryProtoOpts) {
			typ, err := defaultMapping(f, s.File)
			// SET OUT ERR
			if err != nil {
				outErr = err
			}
			cases[f.GetName()] = P(`case "`, fName(f), `": r, ok := (*scanned[i].i).(`, typ, `)
            if !ok {
                return &`, sName, `_`, camelQ(q), `Row{err: fmt.Errorf("cant convert db column `, fName(f), ` to protobuf go type `, f.GetTypeName(), `")}, true
            }
            res.`, camelF(f), `= r
            `)
		}, m.MatchQuery(opt), m.QueryFieldFitsDB)

		// mapping case
		m.EachQueryOut(func(f *desc.FieldDescriptorProto, q *QueryProtoOpts) {
			m.EachTM(func(opt *TypeMappingProtoOpts) {
				_, titled := getGoNamesForTypeMapping(opt.tm, s.File)
				cases[fName(f)] = P(`case "`, fName(f), `":
                    var converted = this.tm.`, titled, `().Empty()
                    if err := converted.Scan(*scanned[i].i); err != nil {
                        return &`, sName, `_`, camelQ(q), `Row{err: fmt.Errorf("could not convert mapped db column `, fName(f), ` to type on `, outName(q), `.`, camelF(f), `: %v", err)}, true
                    }
                    if err := converted.ToProto(&res.`, camelF(f), `); err != nil {
                        return &`, sName, `_`, camelQ(q), `Row{err: fmt.Errorf("could not convert mapped db column `, fName(f), `to type on `, outName(q), `.`, camelF(f), `: %v", err)}, true
                    }
                `)
			}, m.MatchTypeMapping(f))
		}, m.MatchQuery(opt), m.QueryFieldIsMapped)

		printer := &Printer{}

		// loop this way to prevent random order write because map ordering iteration is random
		m.EachQueryOut(func(f *desc.FieldDescriptorProto, _ *QueryProtoOpts) {
			printer.Q(cases[fName(f)])
		}, m.MatchQuery(opt))

		return printer.String()
	}
	m.EachQuery(func(q *QueryProtoOpts) {
		p.Q(`type `, sName, `_`, camelQ(q), `Iter struct {
            result sql.Result
            rows   *sql.Rows
            err    error
            tm     `, sName, `_TypeMappings
            ctx    context.Context
        }

        func (this *`, sName, `_`, camelQ(q), `Iter) IterOutType`, outName(q), `() {}
        func (this *`, sName, `_`, camelQ(q), `Iter) IterInType`, inName(q), `()  {}

        // Each performs 'fun' on each row in the result set.
        // Each respects the context passed to it.
        // It will stop iteration, and returns this.ctx.Err() if encountered.
        func (this *`, sName, `_`, camelQ(q), `Iter) Each(fun func(*`, sName, `_`, camelQ(q), `Row) error) error {
            for {
                select {
                case <-this.ctx.Done():
                    return this.ctx.Err()
                default:
                    if row, ok := this.Next(); !ok {
                        return nil
                    } else if err := fun(row); err != nil {
                        return err
                    }
                }
            }
        }

        // One returns the sole row, or ensures an error if there was not one result when this row is converted
        func (this *`, sName, `_`, camelQ(q), `Iter) One() *`, sName, `_`, camelQ(q), `Row {
            first, hasFirst := this.Next()
            _, hasSecond := this.Next()
            if !hasFirst || hasSecond {
                amount := "none"
                if hasSecond {
                    amount = "multiple"
                }
                return &`, sName, `_`, camelQ(q), `Row{err: fmt.Errorf("expected exactly 1 result from query '`, camelQ(q), `' found %s", amount)}
            }
            return first
        }

        // Zero returns an error if there were any rows in the result
        func (this *`, sName, `_`, camelQ(q), `Iter) Zero() error {
            if _, ok := this.Next(); ok {
                return fmt.Errorf("expected exactly 0 results from query '`, camelQ(q), `'")
            }
            return nil
        }

        // Next returns the next scanned row out of the database, or (nil, false) if there are no more rows
        func (this *`, sName, `_`, camelQ(q), `Iter) Next() (*`, sName, `_`, camelQ(q), `Row, bool) {
            if this.rows == nil || this.err == io.EOF {
                return nil, false
            } else if this.err != nil {
                err := this.err
                this.err = io.EOF
                return &`, sName, `_`, camelQ(q), `Row{err: err}, true
            }
            cols, err := this.rows.Columns()
            if err != nil {
                return &`, sName, `_`, camelQ(q), `Row{err: err}, true
            }
            if !this.rows.Next() {
                if this.err = this.rows.Err(); this.err == nil {
                    this.err = io.EOF
                    return nil, false
                }
            }
            toScan := make([]interface{}, len(cols))
            scanned := make([]alwaysScanner, len(cols))
            for i := range scanned {
                toScan[i] = &scanned[i]
            }
            if this.err = this.rows.Scan(toScan...); this.err != nil {
                return &`, sName, `_`, camelQ(q), `Row{err: this.err}, true
            }
            res := &`, outName(q), `{}
            for i, col := range cols {
                _ = i
                switch col {
                `, colswitch(q), `
                default:
                    return &`, sName, `_`, camelQ(q), `Row{err: fmt.Errorf("unsupported column in output: %s", col)}, true
                }
            }
            return &`, sName, `_`, camelQ(q), `Row{item: res}, true
        }

        // Slice returns all rows found in the iterator as a Slice.
        func (this *`, sName, `_`, camelQ(q), `Iter) Slice() []*`, sName, `_`, camelQ(q), `Row {
            var results []*`, sName, `_`, camelQ(q), `Row
            for {
                if i, ok := this.Next(); ok {
                    results = append(results, i)
                } else {
                    break
                }
            }
            return results
        }

        // returns the known columns for this result
        func (r *`, sName, `_`, camelQ(q), `Iter) Columns() ([]string, error) {
            if r.err != nil {
                return nil, r.err
            }
            if r.rows != nil {
                return r.rows.Columns()
            }
            return nil, nil
        }
        `)
	})

	if outErr == nil {
		return m.Err()
	}
	return
}

func WriteRows(p *Printer, s *Service) (outErr error) {
	m := Matcher(s)
	sName := s.GetName()
	camelQ := func(q *QueryProtoOpts) string {
		return _gen.CamelCase(q.query.GetName())
	}
	camelF := func(f *desc.FieldDescriptorProto) string {
		return _gen.CamelCase(f.GetName())
	}
	methOutName := func(opt *MethodProtoOpts) string {
		return convertedMsgTypeByProtoName(opt.method.GetOutputType(), s.File)
	}
	outName := func(opt *QueryProtoOpts) string {
		return convertedMsgTypeByProtoName(opt.outMsg.GetProtoName(), s.File)
	}
	// remove the err checks from this one method
	mustDefaultMapping := func(f *desc.FieldDescriptorProto) string {
		typ, err := defaultMapping(f, s.File)
		if err != nil {
			outErr = err
		}
		return typ
	}

	inInterfaceFields := func(opt *QueryProtoOpts) string {
		printer := &Printer{}
		m.EachQueryIn(func(f *desc.FieldDescriptorProto, q *QueryProtoOpts) {
			printer.Q(`Get`, camelF(f), `() `, mustDefaultMapping(f), "\n")
		}, m.MatchQuery(opt))
		return printer.String()
	}

	outInterfaceFields := func(opt *QueryProtoOpts) string {
		printer := &Printer{}
		m.EachQueryOut(func(f *desc.FieldDescriptorProto, q *QueryProtoOpts) {
			printer.Q(`Get`, camelF(f), `() `, mustDefaultMapping(f), "\n")
		}, m.MatchQuery(opt))
		return printer.String()
	}

	eachSharedField := func(mopt *MethodProtoOpts, do func(*desc.FieldDescriptorProto, *QueryProtoOpts)) {
		m.EachMethodOut(func(mf *desc.FieldDescriptorProto, thisMopt *MethodProtoOpts) {
			m.EachQueryOut(func(qf *desc.FieldDescriptorProto, thisQopt *QueryProtoOpts) {
				if mf.GetName() == qf.GetName() &&
					mf.GetTypeName() == qf.GetTypeName() &&
					thisMopt.option.GetQuery() == thisQopt.query.GetName() &&
					mopt.method.GetName() == thisMopt.method.GetName() {
					do(qf, thisQopt)
				}
			})
		})
	}
	unwrapMarshelOut := func(qopt *QueryProtoOpts) string {
		// set the field only if it exists in both the method and the query messages
		setSharedOnPointer := func(mopt *MethodProtoOpts) string {
			printer := &Printer{}
			eachSharedField(mopt, func(qf *desc.FieldDescriptorProto, q *QueryProtoOpts) {
				printer.Q(`o.`, camelF(qf), ` = res.`, camelF(qf), "\n")
			})
			return printer.String()
		}
		p := &Printer{}
		m.EachMethod(func(mopt *MethodProtoOpts) {
			p.Q(`if o, ok := (pointerToMsg).(*`, methOutName(mopt), `); ok {
                if o == nil {
                    return fmt.Errorf("must initialize *`, methOutName(mopt), ` before giving to Unwrap()")
                }
                res, _ := this.`, methOutName(mopt), `()
                _ = res
                `, setSharedOnPointer(mopt), `
                return nil
            }`)
		}, m.MatchQueryName(qopt))
		return p.String()
	}
	setQueryOutFields := func(qopt *QueryProtoOpts) string {
		printer := &Printer{}
		m.EachQueryOut(func(f *desc.FieldDescriptorProto, q *QueryProtoOpts) {
			printer.Q(camelF(f), `: this.item.Get`, camelF(f), "(),\n")
		}, m.MatchQuery(qopt))
		return printer.String()
	}
	outMethods := func(q *QueryProtoOpts) string {
		setSharedFields := func(mopt *MethodProtoOpts) string {
			printer := &Printer{}
			eachSharedField(mopt, func(qf *desc.FieldDescriptorProto, q *QueryProtoOpts) {
				printer.Q(camelF(qf), `: this.item.Get`, camelF(qf), "(),\n")
			})
			return printer.String()
		}
		printer := &Printer{}
		did := make(map[string]bool)
		m.EachMethod(func(mopt *MethodProtoOpts) {
			printer.Q(`func (this *`, sName, `_`, camelQ(q), `Row) `, methOutName(mopt), `() (*`, methOutName(mopt), `, error) {
                if this.err != nil {
                    return nil, this.err
                }
                return &`, methOutName(mopt), `{
                `, setSharedFields(mopt), `
                }, nil
            }
            `)
			did[methOutName(mopt)] = true
		}, m.MatchQueryName(q), func(m *MethodProtoOpts) bool {
			// only write the method that hasnt been written yet
			return !did[methOutName(m)]
		})
		return printer.String()
	}
	m.EachQuery(func(q *QueryProtoOpts) {
		p.Q(`type `, sName, `_`, camelQ(q), `In interface {
            `, inInterfaceFields(q), `
        }

        type `, sName, `_`, camelQ(q), `Out interface {
            `, outInterfaceFields(q), `
        }

        type `, sName, `_`, camelQ(q), `Row struct {
            item `, sName, `_`, camelQ(q), `Out
            err  error
        }

        func new`, sName, `_`, camelQ(q), `Row(item `, sName, `_`, camelQ(q), `Out, err error) *`, sName, `_`, camelQ(q), `Row {
            return &`, sName, `_`, camelQ(q), `Row{item, err}
        }

        // Unwrap takes an address to a proto.Message as its only parameter
        // Unwrap can only set into output protos of that match method return types + the out option on the query itself
        func (this *`, sName, `_`, camelQ(q), `Row) Unwrap(pointerToMsg proto.Message) error {
            if this.err != nil {
                return this.err
            }
            `, unwrapMarshelOut(q), `
            return nil
        }
        `, outMethods(q), `

        func (this *`, sName, `_`, camelQ(q), `Row) Proto() (*`, outName(q), `, error) {
            if this.err != nil {
                return nil, this.err
            }
            return &`, outName(q), `{
                `, setQueryOutFields(q), `
            }, nil
        }
        `)
	})
	if outErr == nil {
		return m.Err()
	}
	return
}

func WriteHandlers(p *Printer, s *Service) (outErr error) {
	m := Matcher(s)
	serviceName := s.GetName()
	err := WritePersistServerStruct(p, s.GetName())
	if err != nil {
		return err
	}

	p.Q(`
    type RestOf`, serviceName, `Handlers interface {
    `)

	m.EachMethod(func(mpo *MethodProtoOpts) {
		if m.ServerStreaming(mpo) {
			method := mpo.method.GetName()
			inMsg := mpo.inMsg.GetName()
			p.Q(method, `(*`, inMsg, `, `, serviceName, `_`, method, `Server) error`)
		}
	}, func(mpo *MethodProtoOpts) bool {
		return !proto.HasExtension(mpo.method.Options, persist.E_Opts)
	})

	p.Q("}\n")

	m.EachMethod(func(mpo *MethodProtoOpts) {
		method := mpo.method.GetName()
		inMsg := mpo.inMsg.GetName()
		outMsg := mpo.outMsg.GetName()

		if m.ServerStreaming(mpo) {
			p.Q(`
func (this *`, serviceName, `_Impl) `, method, `(req *`, inMsg, `, stream `, serviceName, `_`, method, `Server) error {
    return this.opts.HANDLERS.`, method, `(req, stream)
}
        `)
		}

		if m.ClientStreaming(mpo) {
			p.Q(`
func (this *`, serviceName, `_Impl) `, method, `(stream `, serviceName, `_`, inMsg, `Server) error {
    return this.opts.HANDLERS.`, inMsg, `(stream)
}
        `)
		}

		if m.Unary(mpo) {
			p.Q(`
func (this *`, serviceName, `_Impl) `, method, `(ctx context.Context, req *`, inMsg, `) (*`, outMsg, `, error) {
    return this.opts.HANDLERS.`, method, `(ctx, req)
}
        `)
		}

		if m.BidiStreaming(mpo) {
			p.Q(`
func (this *`, serviceName, `_Impl) `, method, `(stream `, serviceName, `_`, method, `Server) error {
    return this.opts.HANDLERS.`, method, `(stream)
}
        `)
		}
	}, func(mpo *MethodProtoOpts) bool {
		return !proto.HasExtension(mpo.method.Options, persist.E_Opts)
	})

	m.EachMethod(func(mpo *MethodProtoOpts) {
		var queryOptions *QueryProtoOpts
		m.EachQuery(func(qpo *QueryProtoOpts) {
			queryOptions = qpo
		}, func(qpo *QueryProtoOpts) bool {
			if qpo.query.GetName() == mpo.option.GetQuery() {
				return true
			}
			return false
		})

		zeroResponse := len(queryOptions.outFields) == 0
		params := &handlerParams{
			Service:      serviceName,
			Method:       mpo.method.GetName(),
			Request:      mpo.inMsg.GetName(),
			Response:     mpo.outMsg.GetName(),
			ZeroResponse: zeroResponse,
			Query:        mpo.option.GetQuery(),
			Before:       mpo.option.GetBefore(),
			After:        mpo.option.GetAfter(),
		}

		if m.Unary(mpo) {
			err = WriteUnary(p, params)
			if err != nil {
				outErr = err
			}
		}
		if m.ClientStreaming(mpo) {
			err = WriteClientStreaming(p, params)
			if err != nil {
				outErr = err
			}
		}

		if m.ServerStreaming(mpo) {
			err = WriteSeverStream(p, params)
			if err != nil {
				outErr = err
			}
		}

		if m.BidiStreaming(mpo) {
			err = WriteBidirectionalStream(p, params)
			if err != nil {
				outErr = err
			}
		}

	}, func(mpo *MethodProtoOpts) bool {
		// Only methods that have persist options
		return proto.HasExtension(mpo.method.Options, persist.E_Opts)
	})

	return nil
}

// BUG:: ONLY WORKS WITH sql.DB, needs to work with spanner
func WriteImports(p *Printer, f *FileStruct) error {
	p.PA([]string{
		"// This file is generated by protoc-gen-persist\n",
		"// Source File: %s\n",
		"// DO NOT EDIT !\n",
		"package %s\n",
	}, f.GetOrigName(), f.GetImplPackage())
	f.SanatizeImports()
	p.P("import(\n")
	for _, i := range *f.ImportList {
		if f.NotSameAsMyPackage(i.GoImportPath) {
			p.P("%s \"%s\"\n", i.GoPackageName, i.GoImportPath)
		}
	}
	p.P(")\n")
	p.Q(`type alwaysScanner struct {
        i *interface{}
    }

    func (s *alwaysScanner) Scan(src interface{}) error {
        s.i = &src
        return nil
    }

    type scanable interface {
        Scan(...interface{}) error
        Columns() ([]string, error)
    }
    type Runnable interface {
        QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
        ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
    }

    func DefaultClientStreamingPersistTx(ctx context.Context, db *sql.DB) (PersistTx, error) {
        return db.BeginTx(ctx, nil)
    }
    func DefaultServerStreamingPersistTx(ctx context.Context, db *sql.DB) (PersistTx, error) {
        return NopPersistTx(db)
    }
    func DefaultBidiStreamingPersistTx(ctx context.Context, db *sql.DB) (PersistTx, error) {
        return NopPersistTx(db)
    }
    func DefaultUnaryPersistTx(ctx context.Context, db *sql.DB) (PersistTx, error) {
        return NopPersistTx(db)
    }

    type ignoreTx struct {
        r Runnable
    }

    func (this *ignoreTx) Commit() error   { return nil }
    func (this *ignoreTx) Rollback() error { return nil }
    func (this *ignoreTx) QueryContext(ctx context.Context, x string, ys ...interface{}) (*sql.Rows, error) {
        return this.r.QueryContext(ctx, x, ys...)
    }
    func (this *ignoreTx) ExecContext(ctx context.Context, x string, ys ...interface{}) (sql.Result, error) {
        return this.r.ExecContext(ctx, x, ys...)
    }
    type PersistTx interface {
        Commit() error
        Rollback() error
        Runnable
    }

    func NopPersistTx(r Runnable) (PersistTx, error) {
        return &ignoreTx{r}, nil
    }
    `)
	return nil
}
