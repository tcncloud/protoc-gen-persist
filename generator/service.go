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
	if in == nil || out == nil {
		existing := make([]string, 0)
		for _, v := range *all {
			existing = append(existing, v.GetProtoName())
		}

		return nil, fmt.Errorf("in/out message did not exist: in: %s\n out: %s\nexisting: %s", qopt.GetIn(), qopt.GetOut(), strings.Join(existing, "\n"))
	}
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
	qin := func(q *QueryProtoOpts) string {
		return q.inMsg.GetGoName()
	}
	qout := func(q *QueryProtoOpts) string {
		return q.outMsg.GetGoName()
	}
	runnable := func() string {
		if s.IsSQL() {
			return `persist.Runnable`
		} else if s.IsSpanner() {
			return `persist.SpannerRunnable`
		}
		return ``
	}

	createNextParamMarker := func(pmStrat string) func(string) string {
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
	}

	queryAndFields := func(q *QueryProtoOpts) (string, []string) {
		orig := strings.Join(q.query.GetQuery(), " ")
		pmStrat := q.query.GetPmStrategy()
		nextParamMarker := createNextParamMarker(pmStrat)
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
	qFieldDoc := func(q *QueryProtoOpts) string {
		_, res := queryAndFields(q)
		return P(res)
	}

	p.Q("// Queries_", sName, " holds all the queries found the proto service option as methods\n")
	p.Q("type Queries_", sName, " struct {\n")
	p.Q("opts Opts_", sName, "\n")
	p.Q("}\n")

	p.Q(`// Queries`, sName, ` returns all the known 'SQL' queires for the '`, sName, `' service.
// If no opts are provided default implementations are used.
func Queries`, sName, `(opts ... Opts_`, sName, `) * Queries_`, sName, ` {
    var myOpts Opts_`, sName, `
    if len(opts) > 0 {
        myOpts = opts[0]
    } else {
        myOpts = Opts`, sName, `(&DefaultHooks_`, sName, `{}, &DefaultTypeMappings_`, sName, `{})
    }
    return &Queries_`, sName, `{
        opts: myOpts,
    }
}
    `)

	if s.IsSQL() {
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

		// GetFriendsQuery returns a struct that will perform the 'get_friends' query.
		// When Execute is called, it will expect the following fields:
		m.EachQuery(func(q *QueryProtoOpts) {
			p.Q(`
// `, camelQ(q), ` returns a struct that will perform the '`, qname(q), `' query.
// When Execute is called, it will use the following fields:
// `, qFieldDoc(q), `
func (this *Queries_`, sName, `) `, camelQ(q), `(ctx context.Context, db `, runnable(), `) *Query_`, sName, `_`, camelQ(q), ` {
    return &Query_`, sName, `_`, camelQ(q), `{
        opts: this.opts,
        ctx: ctx,
        db: db,
    }
}

// Query_`, sName, `_`, camelQ(q), ` (future doc string needed) 
type Query_`, sName, `_`, camelQ(q), ` struct {
    opts Opts_`, sName, `
    db `, runnable(), `
    ctx context.Context
}

func (this *Query_`, sName, `_`, camelQ(q), `) QueryInType_`, qin(q), `()  {}
func (this *Query_`, sName, `_`, camelQ(q), `) QueryOutType_`, qout(q), `() {}

// Executes the query '`, qname(q), `' with parameters retrieved from x.
// Fields used: `, qFieldDoc(q), `
func (this *Query_`, sName, `_`, camelQ(q), `) Execute(x In_`, sName, `_`, camelQ(q), `) *Iter_`, sName, `_`, camelQ(q), ` {
    var setupErr error
    params := []interface{}{
    `, execParams(q), `
    }
    result := &Iter_`, sName, `_`, camelQ(q), `{
        tm: this.opts.MAPPINGS,
        ctx: this.ctx,
    }
    if setupErr != nil {
        result.err = setupErr
        return result
    }
    result.`, resultOrRows(q), `, result.err = this.db.`, qmethod(q), `Context(this.ctx, "`, qstring(q), `", params...)

    return result
}
        `)
		})
	}

	if s.IsSpanner() {

		populateParams := func(q *QueryProtoOpts) string {
			orig := strings.Join(q.query.GetQuery(), " ")
			result := make([]string, 0)

			r := regexp.MustCompile("@[a-zA-Z0-9_]*")
			potentialFieldNames := r.FindAllString(orig, -1)
			mappedField := make(map[string]string)
			fieldsMap := make(map[string]desc.FieldDescriptorProto_Type)
			repeatedMap := make(map[string]desc.FieldDescriptorProto_Label)

			for _, v := range q.inFields {
				m.EachTM(func(opts *TypeMappingProtoOpts) {
					if v.GetLabel() != opts.tm.GetProtoLabel() {
						return
					} else if v.GetTypeName() != opts.tm.GetProtoTypeName() {
						return
					} else if v.GetType() != opts.tm.GetProtoType() {
						return
					}
					_, titled := getGoNamesForTypeMapping(opts.tm, s.File)
					mappedField[v.GetName()] = titled
				})

				fieldsMap[v.GetName()] = v.GetType()
				repeatedMap[v.GetName()] = v.GetLabel()
			}

			for _, pf := range potentialFieldNames {
				start := strings.Index(orig, pf)
				stop := start + len(pf)
				fieldType, ok := fieldsMap[pf[1:]]
				if ok { // it was a field, mark it
					key := pf[1:]
					mappingType, isMapped := mappedField[key]
					label := repeatedMap[key]
					if isMapped {
						result = append(result,
							key+`, err := this.opts.MAPPINGS.`+mappingType+`().ToSpanner(x.Get`+_gen.CamelCase(key)+`()).SpannerValue()
                            if err != nil {
                                return nil, err
                            }
                            result["`+key+`"] = `+key+`
                            `,
						)
					} else if label == desc.FieldDescriptorProto_LABEL_REPEATED && fieldType == desc.FieldDescriptorProto_TYPE_MESSAGE {
						result = append(result, `
							`+key+` := make([][]byte, 0)
							for _, tmp := range x.Get`+_gen.CamelCase(key)+`() {
								bytes, err := proto.Marshal(tmp)
								if err != nil {
									return nil, err
								}
								`+key+` = append(`+key+`, bytes)
							}
							result["`+key+`"] = `+key+`
						`)
					} else if fieldType == desc.FieldDescriptorProto_TYPE_MESSAGE {
						result = append(result, `
                        `+key+`, err := proto.Marshal(x.Get`+_gen.CamelCase(key)+`())
                        if err != nil {
                            return nil, err
                        }
                         result["`+key+`"] = `+key)
					} else {
						result = append(result, `result["`+key+`"] = x.Get`+_gen.CamelCase(key)+`()`)
					}
				}
				orig = orig[stop:]
			}

			return strings.Join(result, "\n")
		}

		m.EachQuery(func(q *QueryProtoOpts) {
			p.Q(`
// `, camelQ(q), ` returns a struct that will perform the '`, qname(q), `' query.
// When Execute is called, it will use the following fields:
// `, qFieldDoc(q), `
func (this *Queries_`, sName, `) `, camelQ(q), `(ctx context.Context, db `, runnable(), `) *Query_`, sName, `_`, camelQ(q), ` {
    return &Query_`, sName, `_`, camelQ(q), `{
        opts: this.opts,
        ctx: ctx,
        db: db,
    }
}

// Query_`, sName, `_`, camelQ(q), ` (future doc string needed) 
type Query_`, sName, `_`, camelQ(q), ` struct {
    opts Opts_`, sName, `
    db `, runnable(), `
    ctx context.Context
}

func (this *Query_`, sName, `_`, camelQ(q), `) QueryInType_`, qin(q), `()  {}
func (this *Query_`, sName, `_`, camelQ(q), `) QueryOutType_`, qout(q), `() {}

// Executes the query '`, qname(q), `' with parameters retrieved from x.
// Fields used: `, qFieldDoc(q), `
func (this *Query_`, sName, `_`, camelQ(q), `) Execute(x In_`, sName, `_`, camelQ(q), `) *Iter_`, sName, `_`, camelQ(q), ` {
    ctx := this.ctx
    result := &Iter_`, sName, `_`, camelQ(q), `{
        result: &SpannerResult{},
        tm: this.opts.MAPPINGS,
        ctx: ctx,
    }
    params, err  := func() (map[string]interface{}, error) {
        result := make(map[string]interface{})
        `, populateParams(q), `
        return result, nil
    }()
    if err != nil {
        result.err = err
        return result
    }

    iter := this.db.QueryWithStats(ctx, spanner.Statement{
        SQL: "`, qstring(q), `",
        Params: params,
    })
    result.rows = iter

    result.err = err
    return result
}
        `)
		})
	}

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

	p.Q("type Hooks_", sName, " interface {\n")
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
	p.Q("type DefaultHooks_", sName, " struct{}\n")
	m.EachMethod(func(m *MethodProtoOpts) {
		p.Q(`func(*DefaultHooks_`, sName, `) `, name(m), `BeforeHook(context.Context, *`, inName(m), `) ([]*`, outName(m), `, error) {
            return nil, nil
        }
        `)
	}, func(m *MethodProtoOpts) bool {
		return false
	}, m.BeforeHook, m.ServerStreaming)
	m.EachMethod(func(m *MethodProtoOpts) {
		p.Q(`func(*DefaultHooks_`, sName, `) `, name(m), `BeforeHook(context.Context, *`, inName(m), `) (*`, outName(m), `, error) {
            return nil, nil
        }
        `)
	}, m.BeforeHook)
	m.EachMethod(func(m *MethodProtoOpts) {
		p.Q(`func(*DefaultHooks_`, sName, `) `, name(m), `AfterHook(context.Context, *`, inName(m), `,*`, outName(m), `)error {
            return nil
        }
        `)
	}, m.AfterHook)

	return m.Err()
}

func WriteTypeMappings(p *Printer, s *Service) error {
	sName := s.GetName()
	// TODO google's WKT protobufs probably don't need the package prefix
	p.Q("type TypeMappings_", sName, " interface{\n")
	tms := s.GetTypeMapping().GetTypes()
	for _, tm := range tms {
		_, titled := getGoNamesForTypeMapping(tm, s.File)
		p.Q(titled, "() MappingImpl_", sName, "_", titled, "\n")
	}
	p.Q("}\n")
	m := Matcher(s)
	m.EachTM(func(tm *TypeMappingProtoOpts) {
		p.Q(`
		`)
	})
	p.Q(`type DefaultTypeMappings_`, sName, ` struct{}
    `)

	if s.IsSQL() {
		for _, tm := range tms {
			name, titled := getGoNamesForTypeMapping(tm, s.File)
			p.Q(`func (this *DefaultTypeMappings_`, sName, `) `, titled, `() MappingImpl_`, sName, `_`, titled, ` {
            return &DefaultMappingImpl_`, sName, `_`, titled, `{}
        }
        type DefaultMappingImpl_`, sName, `_`, titled, ` struct{}

        func (this *DefaultMappingImpl_`, sName, `_`, titled, `) ToProto(**`, name, `) error {
            return nil
        }
        func (this *DefaultMappingImpl_`, sName, `_`, titled, `) ToSql(*`, name, `) sql.Scanner {
            return this
        }
        func (this *DefaultMappingImpl_`, sName, `_`, titled, `) Scan(interface{}) error {
            return nil
        }
        func (this *DefaultMappingImpl_`, sName, `_`, titled, `) Value() (driver.Value, error) {
            return "DEFAULT_TYPE_MAPPING_VALUE", nil
		}
		type MappingImpl_`, sName, `_`, titled, ` interface {
			ToProto(**`, name, `) error
			ToSql(*`, name, `) sql.Scanner
			sql.Scanner
			driver.Valuer
		}
        `)
		}
	}

	if s.IsSpanner() {
		for _, tm := range tms {
			name, titled := getGoNamesForTypeMapping(tm, s.File)

			p.Q(`func (this *DefaultTypeMappings_`, sName, `) `, titled, `() MappingImpl_`, sName, `_`, titled, ` {
            return &DefaultMappingImpl_`, sName, `_`, titled, `{}
        }
        type DefaultMappingImpl_`, sName, `_`, titled, ` struct{}

        func (this *DefaultMappingImpl_`, sName, `_`, titled, `) ToProto(**`, name, `) error {
            return nil
        }
        func (this *DefaultMappingImpl_`, sName, `_`, titled, `) ToSpanner(*`, name, `) persist.SpannerScanValuer {
            return this
        }
        func (this *DefaultMappingImpl_`, sName, `_`, titled, `) SpannerScan(*spanner.GenericColumnValue) error {
            return nil
        }
        func (this *DefaultMappingImpl_`, sName, `_`, titled, `) SpannerValue() (interface{}, error) {
            return "DEFAULT_TYPE_MAPPING_VALUE", nil
		}

		type MappingImpl_`, sName, `_`, titled, ` interface{
			ToProto(**`, name, `) error
            ToSpanner(*`, name, `) persist.SpannerScanValuer 
            SpannerScan(*spanner.GenericColumnValue) error
			SpannerValue() (interface{}, error)
		}
        `)
		}
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
	inNamePkg := func(opt *QueryProtoOpts) string {
		return _gen.CamelCase(strings.Map(func(r rune) rune {
			if r == '.' {
				return -1
			}
			return r
		}, convertedMsgTypeByProtoName(opt.inMsg.GetProtoName(), s.File)))
	}
	outNamePkg := func(opt *QueryProtoOpts) string {
		return _gen.CamelCase(strings.Map(func(r rune) rune {
			if r == '.' {
				return -1
			}
			return r
		}, convertedMsgTypeByProtoName(opt.outMsg.GetProtoName(), s.File)))
	}
	_ = outNamePkg
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
                    return &Row_`, sName, `_`, camelQ(q), `{err: fmt.Errorf("cant convert db column `, fName(f), ` to protobuf go type *`, mustDefaultMappingNoStar(f), `")}, true
                }
                var converted = new(`, mustDefaultMappingNoStar(f), `)
                if err := proto.Unmarshal(r, converted); err != nil {
                    return &Row_`, sName, `_`, camelQ(q), `{err: err}, true
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
			cases[f.GetName()] = P(`case "`, fName(f), `":
            r, ok := (*scanned[i].i).(`, typ, `)
            if !ok {
                return &Row_`, sName, `_`, camelQ(q), `{err: fmt.Errorf("cant convert db column `, fName(f), ` to protobuf go type `, f.GetTypeName(), `")}, true
            }
            res.`, camelF(f), `= r
            `)
		}, m.MatchQuery(opt), m.QueryFieldFitsDB)
		// enum case
		m.EachQueryOut(func(f *desc.FieldDescriptorProto, q *QueryProtoOpts) {
			ename := convertedMsgTypeByProtoName(f.GetTypeName(), s.File)
			cases[fName(f)] = P(`case "`, fName(f), `":
                r, ok := (*scanned[i].i).(int64)
                if !ok {
                    return &Row_`, sName, `_`, camelQ(q), `{err: fmt.Errorf("cant convert db column `, fName(f), ` to protobuf go type *`, mustDefaultMappingNoStar(f), `")}, true
                }
                var converted = (`, ename, `)(int32(r))
                res.`, camelF(f), ` = converted
            `)
		}, m.MatchQuery(opt), func(f *desc.FieldDescriptorProto, q *QueryProtoOpts) bool {
			str := s.AllStructs.GetStructByProtoName(f.GetTypeName())
			return str != nil && str.EnumDesc != nil
		})
		// mapping case
		m.EachQueryOut(func(f *desc.FieldDescriptorProto, q *QueryProtoOpts) {
			m.EachTM(func(opt *TypeMappingProtoOpts) {
				_, titled := getGoNamesForTypeMapping(opt.tm, s.File)
				cases[fName(f)] = P(`case "`, fName(f), `":
                    var converted = this.tm.`, titled, `()
                    if err := converted.Scan(*scanned[i].i); err != nil {
                        return &Row_`, sName, `_`, camelQ(q), `{err: fmt.Errorf("could not convert mapped db column `, fName(f), ` to type on `, outName(q), `.`, camelF(f), `: %v", err)}, true
                    }
                    if err := converted.ToProto(&res.`, camelF(f), `); err != nil {
                        return &Row_`, sName, `_`, camelQ(q), `{err: fmt.Errorf("could not convert mapped db column `, fName(f), `to type on `, outName(q), `.`, camelF(f), `: %v", err)}, true
                    }
                `)
			}, m.MatchTypeMapping(f))
		}, m.MatchQuery(opt), m.QueryFieldIsMapped)
		// mapped enum
		m.EachQueryOut(func(f *desc.FieldDescriptorProto, q *QueryProtoOpts) {
			m.EachTM(func(opt *TypeMappingProtoOpts) {
				_, titled := getGoNamesForTypeMapping(opt.tm, s.File)
				cases[fName(f)] = P(`case "`, fName(f), `":
                    var converted = this.tm.`, titled, `()
                    if err := converted.Scan(*scanned[i].i); err != nil {
                        return &Row_`, sName, `_`, camelQ(q), `{err: fmt.Errorf("could not convert mapped db column `, fName(f), ` to type on `, outName(q), `.`, camelF(f), `: %v", err)}, true
                    }
                    pToRes := &res.`, camelF(f), `

                    if err := converted.ToProto(&pToRes); err != nil {
                        return &Row_`, sName, `_`, camelQ(q), `{err: fmt.Errorf("could not convert mapped db column `, fName(f), `to type on `, outName(q), `.`, camelF(f), `: %v", err)}, true
                    }
                `)
			}, m.MatchTypeMapping(f))
		}, m.MatchQuery(opt), m.QueryFieldIsMapped, func(f *desc.FieldDescriptorProto, q *QueryProtoOpts) bool {
			str := s.AllStructs.GetStructByProtoName(f.GetTypeName())
			return str != nil && str.EnumDesc != nil
		})

		printer := &Printer{}

		// loop this way to prevent random order write because map ordering iteration is random
		m.EachQueryOut(func(f *desc.FieldDescriptorProto, _ *QueryProtoOpts) {
			printer.Q(cases[fName(f)])
		}, m.MatchQuery(opt))

		return printer.String()
	}
	// SQL Iterators
	if s.IsSQL() {
		m.EachQuery(func(q *QueryProtoOpts) {
			p.Q(`
        type Iter_`, sName, `_`, camelQ(q), ` struct {
            result sql.Result
            rows   *sql.Rows
            err    error
            tm     TypeMappings_`, sName, `
            ctx    context.Context
        }

        func (this *Iter_`, sName, `_`, camelQ(q), `) IterOutType`, outNamePkg(q), `() {}
        func (this *Iter_`, sName, `_`, camelQ(q), `) IterInType`, inNamePkg(q), `()  {}

        // Each performs 'fun' on each row in the result set.
        // Each respects the context passed to it.
        // It will stop iteration, and returns this.ctx.Err() if encountered.
        func (this *Iter_`, sName, `_`, camelQ(q), `) Each(fun func(*Row_`, sName, `_`, camelQ(q), `) error) error {
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
        func (this *Iter_`, sName, `_`, camelQ(q), `) One() *Row_`, sName, `_`, camelQ(q), ` {
            first, hasFirst := this.Next()
            if first != nil && first.err != nil {
                return &Row_`, sName, `_`, camelQ(q), `{err: first.err}
            }

            _, hasSecond := this.Next()
            if !hasFirst || hasSecond {
                amount := "none"
                if hasSecond {
                    amount = "multiple"
                }
                return &Row_`, sName, `_`, camelQ(q), `{err: fmt.Errorf("expected exactly 1 result from query '`, camelQ(q), `' found %s", amount)}
            }
            return first
        }

        // Zero returns an error if there were any rows in the result
        func (this *Iter_`, sName, `_`, camelQ(q), `) Zero() error {
            row, ok := this.Next()
            if row != nil && row.err != nil {
                return row.err
            }
            if ok {
                return fmt.Errorf("expected exactly 0 results from query '`, camelQ(q), `'")
            }
            return nil
        }

        // Next returns the next scanned row out of the database, or (nil, false) if there are no more rows
        func (this *Iter_`, sName, `_`, camelQ(q), `) Next() (*Row_`, sName, `_`, camelQ(q), `, bool) {
            if this.rows == nil || this.err == io.EOF {
                return nil, false
            } else if this.err != nil {
                err := this.err
                this.err = io.EOF
                return &Row_`, sName, `_`, camelQ(q), `{err: err}, true
            }
            cols, err := this.rows.Columns()
            if err != nil {
                return &Row_`, sName, `_`, camelQ(q), `{err: err}, true
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
                return &Row_`, sName, `_`, camelQ(q), `{err: this.err}, true
            }
            res := &`, outName(q), `{}
            for i, col := range cols {
                _ = i
                switch col {
                `, colswitch(q), `
                default:
                    return &Row_`, sName, `_`, camelQ(q), `{err: fmt.Errorf("unsupported column in output: %s", col)}, true
                }
            }
            return &Row_`, sName, `_`, camelQ(q), `{item: res}, true
        }

        // Slice returns all rows found in the iterator as a Slice.
        func (this *Iter_`, sName, `_`, camelQ(q), `) Slice() []*Row_`, sName, `_`, camelQ(q), ` {
            var results []*Row_`, sName, `_`, camelQ(q), `
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
        func (r *Iter_`, sName, `_`, camelQ(q), `) Columns() ([]string, error) {
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

	}

	// Spanner Iterators
	if s.IsSpanner() {

		rowToProto := func(q *QueryProtoOpts) string {
			acc := make([]string, 0)
			names := make([]string, 0)
			for _, field := range q.outFields {
				name := field.GetName()
				names = append(names, name)
				goType := mustDefaultMapping(field)

				if m.QueryFieldIsMapped(field, q) {
					m.EachTM(func(opt *TypeMappingProtoOpts) {
						_, titled := getGoNamesForTypeMapping(opt.tm, s.File)
						acc = append(acc, `
                        var `+name+` `+goType+`
                        var `+name+`_col spanner.GenericColumnValue
                        if err := row.ColumnByName("`+name+`", &`+name+`_col); err != nil {
                            return &Row_`+sName+`_`+camelQ(q)+`{err: fmt.Errorf("failed to convert db column `+name+` to spanner.GenericColumnValue")}, true
                        }

                        convert_`+name+` := this.tm.`+titled+`()
                        if err := convert_`+name+`.SpannerScan(&`+name+`_col); err != nil {
                            return &Row_`+sName+`_`+camelQ(q)+`{err: fmt.Errorf("SpannerScan failed for `+name+`")}, true
                        }

                        if err := convert_`+name+`.ToProto(&`+name+`); err != nil {
                            return &Row_`+sName+`_`+camelQ(q)+`{err: fmt.Errorf("ToProto for `+name+` when reading from spanner")}, true
                        }
                        `)

					}, m.MatchTypeMapping(field))
				} else if field.GetType() == desc.FieldDescriptorProto_TYPE_MESSAGE && field.GetLabel() != desc.FieldDescriptorProto_LABEL_REPEATED {
					msg := mustDefaultMappingNoStar(field)
					acc = append(acc, `
                    `+name+` :=  &`+msg+`{}
                    `+name+`Bytes := make([]byte, 0)
                    if err := row.ColumnByName("`+name+`", &`+name+`Bytes); err != nil {
                        return &Row_`+sName+`_`+camelQ(q)+`{err: fmt.Errorf("failed to convert db column `+name+` to []byte")}, true
                    }

                    if err := proto.Unmarshal(`+name+`Bytes, `+name+`); err != nil {
                        return &Row_`+sName+`_`+camelQ(q)+`{err: fmt.Errorf("failed to unmarshal column `+name+` to proto message")}, true
                    }
                    `)
				} else if field.GetType() == desc.FieldDescriptorProto_TYPE_MESSAGE && field.GetLabel() == desc.FieldDescriptorProto_LABEL_REPEATED {
					goType := s.File.GetGoTypeName(field.GetTypeName())
					msg := mustDefaultMapping(field)
					acc = append(acc, `
                    `+name+` := make(`+msg+`, 0)
                    `+name+`Bytes := make([][]byte, 0)
                    if err := row.ColumnByName("`+name+`", &`+name+`Bytes); err != nil {
                        return &Row_`+sName+`_`+camelQ(q)+`{err: fmt.Errorf("failed to convert db column `+name+` to [][]byte")}, true
                    }
                    for _, x := range `+name+`Bytes {
                        tmp := &`+goType+`{}
                        if err := proto.Unmarshal(x, tmp); err != nil {
                            return &Row_`+sName+`_`+camelQ(q)+`{err: fmt.Errorf("failed to unmarshal column table to proto message")}, true
                        }
                        `+name+` = append(`+name+`, tmp)
                    }
                    `)
				} else {
					acc = append(acc, "var "+name+" "+goType)
					acc = append(acc, `
            if err := row.ColumnByName("`+name+`", &`+name+`); err != nil {
                return &Row_`+sName+`_`+camelQ(q)+`{err: fmt.Errorf("cant convert db column `+name+` to protobuf go type `+goType+`")}, true
            }
                `)

				}
			}

			acc = append(acc, `res := &`+outName(q)+`{`)
			for _, name := range names {
				acc = append(acc, _gen.CamelCase(name)+": "+name+",")
			}
			acc = append(acc, "}")

			return strings.Join(acc, "\n")
		}

		m.EachQuery(func(q *QueryProtoOpts) {
			p.Q(`
        type Iter_`, sName, `_`, camelQ(q), ` struct {
            result *SpannerResult
            rows   *spanner.RowIterator
            err    error
            tm     TypeMappings_`, sName, `
            ctx    context.Context
        }

        func (this *Iter_`, sName, `_`, camelQ(q), `) IterOutType`, outNamePkg(q), `() {}
        func (this *Iter_`, sName, `_`, camelQ(q), `) IterInType`, inNamePkg(q), `()  {}

        // Each performs 'fun' on each row in the result set.
        // Each respects the context passed to it.
        // It will stop iteration, and returns this.ctx.Err() if encountered.
        func (this *Iter_`, sName, `_`, camelQ(q), `) Each(fun func(*Row_`, sName, `_`, camelQ(q), `) error) error {
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
        func (this *Iter_`, sName, `_`, camelQ(q), `) One() *Row_`, sName, `_`, camelQ(q), ` {
            first, hasFirst := this.Next()
            if first != nil && first.err != nil {
                return &Row_`, sName, `_`, camelQ(q), `{err: first.err}
            }

            _, hasSecond := this.Next()
            if !hasFirst || hasSecond {
                amount := "none"
                if hasSecond {
                    amount = "multiple"
                }
                return &Row_`, sName, `_`, camelQ(q), `{err: fmt.Errorf("expected exactly 1 result from query '`, camelQ(q), `' found %s", amount)}
            }
            return first
        }

        // Zero returns an error if there were any rows in the result
        func (this *Iter_`, sName, `_`, camelQ(q), `) Zero() error {
            row, ok := this.Next()
            if row != nil && row.err != nil {
                return row.err
            }
            if ok {
                return fmt.Errorf("expected exactly 0 results from query '`, camelQ(q), `'")
            }
            return nil
        }

        // Next returns the next scanned row out of the database, or (nil, false) if there are no more rows
        func (this *Iter_`, sName, `_`, camelQ(q), `) Next() (*Row_`, sName, `_`, camelQ(q), `, bool) {
            row, err := this.rows.Next()
            _ = row
            if err == iterator.Done {
                return nil, false
            }
            if err != nil {
                return &Row_`, sName, `_`, camelQ(q), `{err: err}, true
            }

            `, rowToProto(q), `

            return &Row_`, sName, `_`, camelQ(q), `{item: res}, true
        }

        // Slice returns all rows found in the iterator as a Slice.
        func (this *Iter_`, sName, `_`, camelQ(q), `) Slice() []*Row_`, sName, `_`, camelQ(q), ` {
            var results []*Row_`, sName, `_`, camelQ(q), `
            for {
                if i, ok := this.Next(); ok {
                    results = append(results, i)
                } else {
                    break
                }
            }
            return results
        }

        `)
		})
	}

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
	methOutNamePkg := func(opt *MethodProtoOpts) string {
		return _gen.CamelCase(strings.Map(func(r rune) rune {
			if r == '.' {
				return -1
			}
			return r
		}, convertedMsgTypeByProtoName(opt.method.GetOutputType(), s.File)))
	}
	_ = methOutNamePkg
	outNamePkg := func(opt *QueryProtoOpts) string {
		return _gen.CamelCase(strings.Map(func(r rune) rune {
			if r == '.' {
				return -1
			}
			return r
		}, convertedMsgTypeByProtoName(opt.outMsg.GetProtoName(), s.File)))
	}
	_ = outNamePkg
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

	unwrapQueryOut := func(qopt *QueryProtoOpts) string {
		p := &Printer{}

		setFields := func() string {
			printer := &Printer{}
			for _, field := range qopt.outFields {
				printer.Q(`o.`, camelF(field), ` = res.`, camelF(field), "\n")
			}
			return printer.String()
		}

		p.Q(`if o, ok := (pointerToMsg).(*`, outName(qopt), `); ok {
            if o == nil {
                return fmt.Errorf("must initialize *`, outName(qopt), ` before giving to Unwrap()")
            }
            res, _ := this.`, outNamePkg(qopt), `()
            _ = res
            `, setFields(), `
            return nil
        }
        `)

		return p.String()
	}

	unwrapMarshalOut := func(qopt *QueryProtoOpts) string {
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
                res, _ := this.`, methOutNamePkg(mopt), `()
                _ = res
                `, setSharedOnPointer(mopt), `
                return nil
            }
            `)
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
		setQueryOutFields := func(q *QueryProtoOpts) string {
			printer := &Printer{}
			for _, field := range q.outFields {
				printer.Q(camelF(field), `: this.item.Get`, camelF(field), "(),\n")
			}
			return printer.String()
		}
		setSharedFields := func(mopt *MethodProtoOpts) string {
			printer := &Printer{}
			eachSharedField(mopt, func(qf *desc.FieldDescriptorProto, q *QueryProtoOpts) {
				printer.Q(camelF(qf), `: this.item.Get`, camelF(qf), "(),\n")
			})
			return printer.String()
		}
		printer := &Printer{}
		did := make(map[string]bool)

		printer.Q(`func (this *Row_`, sName, `_`, camelQ(q), `) `, outNamePkg(q), `() (*`, outName(q), `, error) {
			if this.err != nil {
				return nil, this.err
			}
			return &`, outName(q), `{
				`, setQueryOutFields(q), `
			}, nil
		}
		`)
		did[outName(q)] = true

		m.EachMethod(func(mopt *MethodProtoOpts) {
			printer.Q(`func (this *Row_`, sName, `_`, camelQ(q), `) `, methOutNamePkg(mopt), `() (*`, methOutName(mopt), `, error) {
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
		p.Q(`type In_`, sName, `_`, camelQ(q), ` interface {
            `, inInterfaceFields(q), `
        }

        type Out_`, sName, `_`, camelQ(q), ` interface {
            `, outInterfaceFields(q), `
        }

        type Row_`, sName, `_`, camelQ(q), ` struct {
            item Out_`, sName, `_`, camelQ(q), `
            err  error
        }

        func newRow`, sName, camelQ(q), `(item Out_`, sName, `_`, camelQ(q), `, err error) *Row_`, sName, `_`, camelQ(q), ` {
            return &Row_`, sName, `_`, camelQ(q), `{item, err}
        }

        // Unwrap takes an address to a proto.Message as its only parameter
        // Unwrap can only set into output protos of that match method return types + the out option on the query itself
        func (this *Row_`, sName, `_`, camelQ(q), `) Unwrap(pointerToMsg proto.Message) error {
            if this.err != nil {
                return this.err
            }
            `, unwrapQueryOut(q), `
            `, unwrapMarshalOut(q), `
            return nil
        }
        `, outMethods(q), `

        func (this *Row_`, sName, `_`, camelQ(q), `) Proto() (*`, outName(q), `, error) {
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
	methOutNamePkg := func(opt *MethodProtoOpts) string {
		return _gen.CamelCase(strings.Map(func(r rune) rune {
			if r == '.' {
				return -1
			}
			return r
		}, convertedMsgTypeByProtoName(opt.method.GetOutputType(), s.File)))
	}
	db := "sql.DB"
	if s.IsSpanner() {
		db = "spanner.Client"
	}
	err := WritePersistServerStruct(p, s.GetName(), db)
	if err != nil {
		return err
	}

	p.Q(`
    type RestOfHandlers_`, serviceName, ` interface {
    `)

	m.EachMethod(func(mpo *MethodProtoOpts) {
		method := mpo.method.GetName()
		inMsg := s.File.GetGoTypeName(mpo.inStruct.GetProtoName())
		outMsg := s.File.GetGoTypeName(mpo.outStruct.GetProtoName())
		if m.ServerStreaming(mpo) {
			p.Q(method, `(*`, inMsg, `, `, serviceName, `_`, method, `Server) error`, "\n")
		}
		if m.ClientStreaming(mpo) || m.BidiStreaming(mpo) {
			p.Q(method, `(`, serviceName, `_`, method, `Server) error`, "\n")
		}
		if m.Unary(mpo) {
			p.Q(method, `(context.Context, *`, inMsg, `) (*`, outMsg, `, error)`, "\n")
		}
	}, func(mpo *MethodProtoOpts) bool {
		return !proto.HasExtension(mpo.method.Options, persist.E_Opts) || m.BidiStreaming(mpo)
	})

	p.Q("}\n")

	m.EachMethod(func(mpo *MethodProtoOpts) {
		method := mpo.method.GetName()
		inMsg := s.File.GetGoTypeName(mpo.inStruct.GetProtoName())
		outMsg := s.File.GetGoTypeName(mpo.outStruct.GetProtoName())

		if m.ServerStreaming(mpo) {
			p.Q(`
func (this *Impl_`, serviceName, `) `, method, `(req *`, inMsg, `, stream `, serviceName, `_`, method, `Server) error {
    return this.HANDLERS.`, method, `(req, stream)
}
        `)
		}

		if m.ClientStreaming(mpo) {
			p.Q(`
func (this *Impl_`, serviceName, `) `, method, `(stream `, serviceName, `_`, method, `Server) error {
    return this.HANDLERS.`, method, `(stream)
}
        `)
		}

		if m.Unary(mpo) {
			p.Q(`
func (this *Impl_`, serviceName, `) `, method, `(ctx context.Context, req *`, inMsg, `) (*`, outMsg, `, error) {
    return this.HANDLERS.`, method, `(ctx, req)
}
        `)
		}

		if m.BidiStreaming(mpo) {
			p.Q(`
func (this *Impl_`, serviceName, `) `, method, `(stream `, serviceName, `_`, method, `Server) error {
    return this.HANDLERS.`, method, `(stream)
}
        `)
		}
	}, func(mpo *MethodProtoOpts) bool {
		return !proto.HasExtension(mpo.method.Options, persist.E_Opts) || m.BidiStreaming(mpo)
	})

	m.EachMethod(func(mpo *MethodProtoOpts) {
		var queryOptions *QueryProtoOpts
		inMsg := s.File.GetGoTypeName(mpo.inStruct.GetProtoName())
		outMsg := s.File.GetGoTypeName(mpo.outStruct.GetProtoName())
		m.EachQuery(func(qpo *QueryProtoOpts) {
			queryOptions = qpo
		}, func(qpo *QueryProtoOpts) bool {
			if qpo.query.GetName() == mpo.option.GetQuery() {
				return true
			}
			return false
		})

		zeroResponse := queryOptions != nil && len(queryOptions.outFields) == 0
		params := &handlerParams{
			Service:        serviceName,
			Method:         mpo.method.GetName(),
			Request:        inMsg,  //mpo.inMsg.GetName(),
			Response:       outMsg, //mpo.outMsg.GetName(),
			RespMethodCall: methOutNamePkg(mpo),
			ZeroResponse:   zeroResponse,
			Query:          mpo.option.GetQuery(),
			Before:         mpo.option.GetBefore(),
			After:          mpo.option.GetAfter(),
		}

		if m.Unary(mpo) {
			err = WriteUnary(p, params, s.IsSQL())
			if err != nil {
				outErr = err
			}
		}
		if m.ClientStreaming(mpo) {
			err = WriteClientStreaming(p, params, s.IsSQL())
			if err != nil {
				outErr = err
			}
		}

		if m.ServerStreaming(mpo) {
			err = WriteServerStream(p, params, s.IsSQL())
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
	hasSQL := false
	hasSpanner := false
	for _, service := range *f.ServiceList {
		if service.IsSQL() {
			hasSQL = true
		}

		if service.IsSpanner() {
			hasSpanner = true
		}
	}

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
	p.P("%s \"%s\"\n", "proto", "github.com/golang/protobuf/proto")
	p.Q("persist ", "\"github.com/tcncloud/protoc-gen-persist/persist\"\n")

	if hasSpanner {
		p.P("%s \"%s\"\n", "iterator", "google.golang.org/api/iterator")
	}

	p.P(")\n")

	if hasSQL {
		p.Q(`

func NopPersistTx(r persist.Runnable) (persist.PersistTx, error) {
    return &ignoreTx{r}, nil
}

type ignoreTx struct {
    r persist.Runnable
}

func (this *ignoreTx) Commit() error   { return nil }
func (this *ignoreTx) Rollback() error { return nil }
func (this *ignoreTx) QueryContext(ctx context.Context, x string, ys ...interface{}) (*sql.Rows, error) {
    return this.r.QueryContext(ctx, x, ys...)
}
func (this *ignoreTx) ExecContext(ctx context.Context, x string, ys ...interface{}) (sql.Result, error) {
    return this.r.ExecContext(ctx, x, ys...)
}

type Runnable interface {
    QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
    ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
}

func DefaultClientStreamingPersistTx(ctx context.Context, db *sql.DB) (persist.PersistTx, error) {
    return db.BeginTx(ctx, nil)
}
func DefaultServerStreamingPersistTx(ctx context.Context, db *sql.DB) (persist.PersistTx, error) {
    return NopPersistTx(db)
}
func DefaultBidiStreamingPersistTx(ctx context.Context, db *sql.DB) (persist.PersistTx, error) {
    return NopPersistTx(db)
}
func DefaultUnaryPersistTx(ctx context.Context, db *sql.DB) (persist.PersistTx, error) {
    return NopPersistTx(db)
}

type alwaysScanner struct {
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

        `)
	} else if hasSpanner {
		p.Q(`
type Result interface {
    LastInsertId() (int64, error)
    RowsAffected() (int64, error)
}
type SpannerResult struct {
    // TODO shouldn't be an iter
    iter *spanner.RowIterator
}

func (sr *SpannerResult) LastInsertId() (int64, error) {
    // sr.iter.QueryStats or sr.iter.QueryPlan
    return -1, nil
}
func (sr *SpannerResult) RowsAffected() (int64, error) {
    // Execution statistics for the query. Available after RowIterator.Next returns iterator.Done
    return sr.iter.RowCount, nil
}


        `)
	}
	return nil
}
