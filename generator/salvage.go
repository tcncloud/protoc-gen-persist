package generator

import (
	"fmt"
	"strings"

	_gen "github.com/golang/protobuf/protoc-gen-go/generator"
)

func WriteBuilderHookInterfaceAndFunc(p *Printer, s *Service) {
	p.Q("type ", s.GetName(), "Hooks interface{\n")
	for _, m := range *s.Methods {
		opt := m.GetMethodOption()
		if opt == nil {
			continue
		}
		if opt.GetBefore() {
			sliceStarOrStar := "*"
			if m.IsServerStreaming() {
				sliceStarOrStar = "[]*"
			}

			p.Q("\t", m.GetBeforeHookName(), "(*", m.GetInputType(), ") (", sliceStarOrStar, m.GetOutputType(), ", error)\n")
		}
		if opt.GetAfter() {
			p.Q("\t", m.GetAfterHookName(), "(*", m.GetInputType(), ", *", m.GetOutputType(), ") error\n")
		}
	}
	p.Q("}\n")
}
func WriteSqlBuilderHooksAcceptingFunc(p *Printer, serv *Service) {
	s := serv.GetName()
	p.Q(
		"func(b *", s, "ImplBuilder) WithHooks(hs ", s, "Hooks) *", s, "ImplBuilder {\n",
		"b.hooks = hs\n",
		"return b\n",
		"}\n",
	)
}

func WriteBuilderTypeMappingsAcceptingFunc(p *Printer, serv *Service) {
	s := serv.GetName()
	p.Q("func(b *", s, "ImplBuilder) WithTypeMapping(ts ", s, "TypeMapping) *", s, "ImplBuilder {\n")
	p.Q("\tb.mappings = ts\n")
	p.Q("\treturn b\n")
	p.Q("}\n")
}

func WriteBuilderTypeMappingsInterface(p *Printer, s *Service) {
	sName := s.GetName()
	// TODO google's WKT protobufs probably don't need the package prefix
	p.Q("type ", sName, "TypeMapping interface{\n")
	tms := s.GetTypeMapping().GetTypes()
	for _, tm := range tms {
		// TODO implement these interfaces
		_, titled := getGoNamesForTypeMapping(tm, s.File)
		// p.Q(titled, "() ", sName, titled, "MappingImpl\n")
		p.Q(titled, "() ", titled, "MappingImpl\n")
	}
	p.Q("}\n")

}
func WriteScanValuerInterface(p *Printer, s *Service) {
	if s.IsSQL() {
		p.Q("type ScanValuer interface {\n")
		p.Q("\tsql.Scanner\n")
		p.Q("\tdriver.Valuer\n")
		p.Q("}\n")
	} else if s.IsSpanner() {
		p.Q("type ScanValuer interface {\n")
		p.Q("\tSpannerScan(src *spanner.GenericColumnValue) error\n")
		p.Q("\tSpannerValue() (interface{}, error)\n")
		p.Q("}\n")
	}
}

func WriteSqlQueries(p *Printer, s *Service) error {
	qopts := s.GetQueriesOption()
	sName := s.GetName()
	// make a struct of all the queries on this service
	p.Q("type Persist", sName, "Queries struct {\n")
	for _, q := range qopts.GetQueries() {
		// get the structs represented by the in and out tags
		in := s.AllStructs.GetStructByProtoName(q.GetIn())
		out := s.AllStructs.GetStructByProtoName(q.GetOut())
		if !in.IsMessage || !out.IsMessage {
			return fmt.Errorf("in/out option must be proto messages for query: '%s', on service: '%s'", q.GetName(), s.GetName())
		}

		queryNameCamelCase := _gen.CamelCase(q.GetName())
		removeDot := func(r rune) rune {
			if r == '.' {
				return -1
			}
			return r
		}
		queryInputInterface := strings.Map(removeDot, s.File.GetGoTypeName(q.GetIn()))
		p.Q("\t", queryNameCamelCase, "func (", queryInputInterface)
	}
	p.Q("}\n")
	return nil
}
