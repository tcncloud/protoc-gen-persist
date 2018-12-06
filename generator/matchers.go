package generator

import desc "github.com/golang/protobuf/protoc-gen-go/descriptor"

type Match struct {
	s   *Service
	err error
}

func Matcher(s *Service) Match {
	return Match{s: s}
}
func (Match) BeforeHook(mopt *MethodProtoOpts) bool {
	return mopt.option != nil && mopt.option.GetBefore()
}
func (Match) AfterHook(mopt *MethodProtoOpts) bool {
	return mopt.option != nil && mopt.option.GetAfter()
}
func (Match) ClientStreaming(mopt *MethodProtoOpts) bool {
	return mopt.method.GetClientStreaming() && !mopt.method.GetServerStreaming()
}
func (Match) ServerStreaming(mopt *MethodProtoOpts) bool {
	return !mopt.method.GetClientStreaming() && mopt.method.GetServerStreaming()
}
func (Match) BidiStreaming(mopt *MethodProtoOpts) bool {
	return mopt.method.GetClientStreaming() && mopt.method.GetServerStreaming()
}
func (Match) Unary(mopt *MethodProtoOpts) bool {
	return !mopt.method.GetClientStreaming() && !mopt.method.GetServerStreaming()
}
func (Match) QueryFieldIsMapped(field *desc.FieldDescriptorProto, q *QueryProtoOpts) bool {
	return true
}
func (Match) QueryFieldIsMessage(field *desc.FieldDescriptorProto, q *QueryProtoOpts) bool {
	return true
}
func (Match) QueryFieldIsRepeated(field *desc.FieldDescriptorProto, q *QueryProtoOpts) bool {
	return true
}
func (Match) QueryFieldFitsDB(field *desc.FieldDescriptorProto, q *QueryProtoOpts) bool {
	return true
}

func (Match) MatchQuery(opt *QueryProtoOpts) func(*desc.FieldDescriptorProto, *QueryProtoOpts) bool {
	if opt.query == nil {
		return func(*desc.FieldDescriptorProto, *QueryProtoOpts) bool { return false }
	}
	return func(f *desc.FieldDescriptorProto, opt2 *QueryProtoOpts) bool {
		return opt.query.GetName() == opt2.query.GetName()
	}
}
func (Match) MatchMethod(mopt *MethodProtoOpts) func(*QueryProtoOpts) bool {
	if mopt.option == nil {
		return func(*QueryProtoOpts) bool { return false }
	}
	return func(qopt *QueryProtoOpts) bool {
		return mopt.option.GetQuery() == qopt.query.GetName()
	}
}
func (Match) MatchTypeMapping(f *desc.FieldDescriptorProto) func(*TypeMappingProtoOpts) bool {
	return func(opt *TypeMappingProtoOpts) bool {
		tm := opt.tm
		ptn := tm.GetProtoTypeName()
		return tm.GetProtoType() == f.GetType() &&
			tm.GetProtoLabel() == f.GetLabel() &&
			ptn == f.GetTypeName() || ("."+ptn) == f.GetTypeName()
	}

}
func (Match) FilterFieldNames(names []string) func(*desc.FieldDescriptorProto, *QueryProtoOpts) bool {
	return func(f *desc.FieldDescriptorProto, _ *QueryProtoOpts) bool {
		for _, name := range names {
			if f.GetName() == name {
				return false
			}
		}
		return true
	}
}
func (Match) MatchingFieldNames(names []string) func(*desc.FieldDescriptorProto, *QueryProtoOpts) bool {
	return func(f *desc.FieldDescriptorProto, _ *QueryProtoOpts) bool {
		for _, name := range names {
			if f.GetName() == name {
				return true
			}
		}
		return false
	}
}

func (m *Match) EachQuery(do func(*QueryProtoOpts), matches ...func(*QueryProtoOpts) bool) {
	if m.err != nil {
		return
	}

	qopts := m.s.GetQueriesOption()
	for _, qopt := range qopts.GetQueries() {
		q, err := NewQueryProtoOpts(qopt, m.s.AllStructs)
		if err != nil {
			m.err = err
			return
		}
		matchAll := true
		for _, match := range matches {
			matchAll = matchAll && match(q)
		}
		if matchAll {
			do(q)
		}
	}
}
func (m *Match) EachQueryIn(do func(*desc.FieldDescriptorProto, *QueryProtoOpts), matches ...func(*desc.FieldDescriptorProto, *QueryProtoOpts) bool) {
	if m.err != nil {
		return
	}
	qopts := m.s.GetQueriesOption()
	for _, qopt := range qopts.GetQueries() {
		q, err := NewQueryProtoOpts(qopt, m.s.AllStructs)
		if err != nil {
			m.err = err
			return
		}
		for _, f := range q.inFields {
			matchAll := true
			for _, match := range matches {
				matchAll = matchAll && match(f, q)
			}
			if matchAll {
				do(f, q)
			}
		}
	}
}
func (m *Match) EachQueryOut(do func(*desc.FieldDescriptorProto, *QueryProtoOpts), matches ...func(*desc.FieldDescriptorProto, *QueryProtoOpts) bool) {
	if m.err != nil {
		return
	}
	qopts := m.s.GetQueriesOption()
	for _, qopt := range qopts.GetQueries() {
		q, err := NewQueryProtoOpts(qopt, m.s.AllStructs)
		if err != nil {
			m.err = err
			return
		}
		for _, f := range q.outFields {
			matchAll := true
			for _, match := range matches {
				matchAll = matchAll && match(f, q)
			}
			if matchAll {
				do(f, q)
			}
		}
	}
}

func (m *Match) EachMethod(do func(*MethodProtoOpts), matches ...func(*MethodProtoOpts) bool) {
	if m.err != nil {
		return
	}
	for _, me := range m.s.Desc.GetMethod() {
		mopt, err := NewMethodProtoOpts(me, m.s.AllStructs)
		if err != nil {
			m.err = err
			return
		}
		matchAll := true
		for _, match := range matches {
			matchAll = matchAll && match(mopt)
		}
		if matchAll {
			do(mopt)
		}
	}
}

func (m *Match) EachTM(do func(*TypeMappingProtoOpts), matches ...func(*TypeMappingProtoOpts) bool) {
	if m.err != nil {
		return
	}
	for _, me := range m.s.GetTypeMapping().GetTypes() {
		mopt, err := NewTypeMappingProtoOpts(me, m.s.AllStructs)
		if err != nil {
			m.err = err
			return
		}
		matchAll := true
		for _, match := range matches {
			matchAll = matchAll && match(mopt)
		}
		if matchAll {
			do(mopt)
		}
	}
}
func (m *Match) Err() error {
	return m.err
}
