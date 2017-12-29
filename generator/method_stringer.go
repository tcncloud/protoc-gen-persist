package generator

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
)

type Printer struct {
	str string
}

func (p *Printer) P(formatString string, args ...interface{}) {
	p.str += fmt.Sprintf(formatString, args...)
}
func (p *Printer) PA(formatStrings []string, args ...interface{}) {
	s := strings.Join(formatStrings, "")
	p.P(s, args...)
}

func (p *Printer) PTemplate(t string, dot interface{}) {
	var buff bytes.Buffer

	tem, err := template.New("printTemplate").Parse(t)
	if err != nil {
		p.P("\nPARSE ERROR:<%v>\nPARSING:<%s>\n", err, t)
		return
	}
	if err := tem.Execute(&buff, dot); err != nil {
		p.P("\nEXEC ERROR:<%v>\nEXECUTING:<%s>\n", err, t)
		return
	}
	p.P(buff.String())
}

func (p *Printer) String() string {
	return p.str
}

// each method has these steps in common:
//
// - declare header
// - receive input
// - handle before hook
// - map input to persist_lib
// - call child method
// - receive row
// - marshal row to output type
// - handle after hook
// - return output
//
// Each of these things look similar to eachother, but have to be stringed
// in a different way depending on a varity of conditions
// - method type: (unary, client stream, server stream, bidi stream)
// - service type: spanner, or sql
// - whether before hook exists or not
// - whether after hook exists or not
// - input type
// - if input type fields are: (mapped, primitive, or message types)
// - if input type field is repeated or not for all of the above
// - if ouptut type is mapped primitive or message
// - if output type is repeated or not
//
// best approach seems to be to stub out each of these
// steps, and have each step be in charge of stringing a section
// based on the conditions of the method
type MethodStringer struct {
	method *Method
}

func (m *MethodStringer) HandlerString() string {
	printer := &Printer{}
	m.PrintHeader(printer)
	if m.method.IsSpanner() && m.method.GetMethodOption() != nil {
		m.PrintSpannerInput(printer)
		m.PrintSpannerBeforeHook(printer)
		m.PrintSpannerParams(printer)
		m.PrintSpannerPersistLibCall(printer)
		m.PrintSpannerRowReceiving(printer)
		m.PrintSpannerRowToOutputHandling(printer)
		m.PrintAfterHook(printer)
		m.PrintReturnOutput(printer)
	} else {
		// we are not a persist method, just forward it to the other receiver
		m.PrintForwardedMethod(printer)
	}
	return printer.String()
}

func (m *MethodStringer) PrintInterfaceDefinition(p *Printer) {
	if m.method.IsUnary() {
		p.P(
			"%s(ctx context.Context, req *%s) (*%s, error)\n",
			m.method.GetName(),
			m.method.GetInputType(),
			m.method.GetOutputType(),
		)
	} else if m.method.IsServerStreaming() {
		p.P(
			"%s(req *%s, stream %s) error\n",
			m.method.GetName(),
			m.method.GetInputType(),
			NewStreamType(m.method),
		)
	} else {
		p.P(
			"%s(stream %s) error\n",
			m.method.GetName(),
			NewStreamType(m.method),
		)
	}
}

func (m *MethodStringer) PrintHeader(p *Printer) {
	if m.method.IsUnary() {
		p.P(
			"func (s *%sImpl) %s(ctx context.Context, req *%s) (*%s, error) {\n",
			m.method.GetServiceName(),
			m.method.GetName(),
			m.method.GetInputType(),
			m.method.GetOutputType(),
		)
	} else if m.method.IsClientStreaming() {
		p.P(
			"func (s *%sImpl) %s(stream %s) error{\n",
			m.method.GetServiceName(),
			m.method.GetName(),
			NewStreamType(m.method),
		)
	} else if m.method.IsServerStreaming() {
		p.P(
			"func (s *%sImpl) %s(req *%s, stream %s) error{\n",
			m.method.GetServiceName(),
			m.method.GetName(),
			m.method.GetInputType(),
			NewStreamType(m.method),
		)
	} else if m.method.IsBidiStreaming() {
	}
}

func (m *MethodStringer) PrintSpannerInput(p *Printer) {
	p.P("\tvar err error\n_ = err\n")
	// if we are server streaming or unary, we already have input
	if m.method.IsClientStreaming() {
		p.PA([]string{
			"feed, stop := s.PERSIST.%s(stream.Context())\n",
			"for {\nreq, err := stream.Recv()\n if err == io.EOF {\nbreak\n} else if",
			" err != nil {\n %s\n}\n",
		},
			m.method.GetName(),
			NewErrPrinter(m, "error recieving input: %v"),
		)
	}
}

func (m *MethodStringer) PrintSpannerBeforeHook(p *Printer) {
	before := m.method.GetMethodOption().GetBefore()
	if before == nil {
		return
	}
	pkg := m.method.GetGoPackage(before.GetPackage())
	if pkg == "" {
		p.P("beforeRes, err := %s(req)\n", before.GetName())
	} else {
		p.P("beforeRes, err := %s.%s(req)\n", pkg, before.GetName())
	}
	p.PA([]string{
		"if err != nil {\n %s\n}",
		"else if beforeRes != nil {\n",
	},
		NewErrPrinter(m, "error in before hook: %v"),
	)
	if m.method.IsUnary() {
		p.P("return beforeRes, nil\n")
	} else if m.method.IsClientStreaming() {
		p.P("continue\n")
	} else if m.method.IsServerStreaming() {
		p.PA([]string{
			"for _, res := range beforeRes {\n",
			"if err := stream.Send(res); err != nil {\n %s\n}\n}\n",
		},
			NewErrPrinter(m, "error sending back before hook result: %v"),
		)
	}
	p.P("}\n")
}
func (m *MethodStringer) PrintSpannerParams(p *Printer) {
	p.P("params := &persist_lib.%s{}\n", NewPLInputName(m.method))
	typeDescs := m.method.GetTypeDescForQueryFields()

	// if value is mapped, always use the mapped value
	// if value is primitive or repeated primitive, use it
	// else convert to []byte, or [][]byte for spanner
	for _, td := range typeDescs {
		p.P(
			"// set '%s.%s' in params\n",
			m.method.GetInputTypeName(),
			td.ProtoName,
		)
		if td.IsMapped {
			p.PA([]string{
				"if params.%s, err = (%s{}).ToSpanner(req.%s).SpannerValue(); err != nil {\n",
				"%s\n}\n",
			},
				td.Name,
				td.GoName,
				td.Name,
				NewErrPrinter(m, "could not convert type to persist_lib type: %v, err"),
			)
		} else if td.IsMessage {
			if td.IsRepeated {
				p.PA([]string{
					"{\nvar bytesOfBytes [][]byte\n",
					"for _, msg := range req.%s{\n",
					"raw, err := proto.Marshal(msg)\nif err != nil {\n",
					"%s\n}\n",
					"bytesOfBytes = append(bytesOfBytes, raw)\n}\n",
					"params.%s = bytesOfBytes\n}\n",
				},
					td.Name,
					NewErrPrinter(m, "could not convert type to [][]byte, err: %s"),
					td.Name,
				)
			} else {
				p.PA([]string{
					"{\nraw, err := proto.Marshal(req.%s)\nif err != nil {\n",
					"%s\n}\n",
					"params.%s = raw\n}\n",
				},
					td.Name,
					NewErrPrinter(m, "could not convert type to []byte, err: %s"),
					td.Name,
				)
			}
		} else {
			p.P("params.%s = req.%s\n", td.Name, td.Name)
		}
	}
}
func (m *MethodStringer) PrintSpannerPersistLibCall(p *Printer) {
	if m.method.IsClientStreaming() {
		p.P("feed(params)\n}\n")
	} else if m.method.IsUnary() {
		p.PA([]string{
			"var res = %s{}\nvar iterErr error\n_ = iterErr\n_ = res\n",
			"err = s.PERSIST.%s(ctx, params, func(row *spanner.Row) {\n",
		},
			m.method.GetOutputType(),
			m.method.GetName(),
		)
	} else if m.method.IsServerStreaming() {
		p.PA([]string{
			"var iterErr error\n",
			"err = s.PERSIST.%s(stream.Context(), params, func(row *spanner.Row) {\n",
		},
			m.method.GetName(),
		)
	}
}
func (m *MethodStringer) PrintSpannerRowReceiving(p *Printer) {
	if m.method.IsClientStreaming() {
		p.PA([]string{
			"row, err := stop()\nif err != nil {\n %s\n}\n",
			"res := %s{}\n",
			"if row != nil {\n",
		},
			NewErrPrinter(m, "error receiving result row: %v"),
			m.method.GetOutputType(),
		)
		return
	}
	p.P("if row == nil { // there was no return data\n return\n}\n")
}
func (m *MethodStringer) PrintSpannerRowToOutputHandling(p *Printer) {
	if m.method.IsServerStreaming() {
		p.P("res := %s{}\n", m.method.GetOutputType())
	}
	for _, td := range m.method.GetTypeDescArrayForStruct(m.method.GetOutputTypeStruct()) {
		if td.IsMapped {
			p.PA([]string{
				"var %s *spanner.GenericColumnValue\n",
				"if err := row.ColumnByName(\"%s\", %s); err != nil {\n%s\n}\n{\n",
				"local := &%s{}\n",
				"if err := local.SpannerScan(%s); err != nil {\n %s\n}\n",
				"res.%s = local.ToProto()\n}\n",
			},
				td.Name,
				td.ProtoName,
				td.Name,
				NewNestedErrPrinter(m),
				td.GoName,
				td.Name,
				NewNestedErrPrinter(m),
				td.Name,
			)
		} else if td.IsMessage {
			// this is super tacky.  But I can be sure I need this import at this point
			m.method.
				Service.
				File.ImportList.GetOrAddImport("proto", "github.com/golang/protobuf/proto")
			if td.IsRepeated {
				p.PA([]string{
					"var %s [][]byte\n",
					"if err := row.ColumnByName(\"%s\", &%s); err != nil {\n %s\n}\n{\n",
					"local := make(%s, len(%s))\n",
					"for i := range local {\nlocal[i] = new(%s)\n",
					"if err := proto.Unmarshal(%s[i], local[i]); err != nil {\n %s\n}\n}\n",
					"res.%s = local\n}\n",
				},
					td.Name,
					td.ProtoName,
					td.Name,
					NewNestedErrPrinter(m),
					td.GoName,
					td.Name,
					td.GoTypeName,
					td.Name,
					NewNestedErrPrinter(m),
					td.Name,
				)
			} else {
				p.PA([]string{
					"var %s[]byte\n",
					"if err := row.ColumnByName(\"%s\", &%s); err != nil {\n %s\n}\n{\n",
					"local := new(%s)\n",
					"if err := proto.Unmarshal(%s, local); err != nil {\n %s\n}\n",
					"res.%s = local\n}\n",
				},
					td.Name,
					td.ProtoName,
					td.Name,
					NewNestedErrPrinter(m),
					td.GoTypeName,
					td.Name,
					NewNestedErrPrinter(m),
					td.Name,
				)
			}
		} else if td.IsRepeated {
			p.PA([]string{
				"var %s %s\n{\n",
				"local := make(%s, 0)\n",
				"if err := row.ColumnByName(\"%s\", &local); err != nil {\n %s\n}\n",
				"for _, l := range local {\nif l.Valid {\n",
				"%s = append(%s, l.%s)\n",
				"res.%s = %s\n}\n}\n}\n",
			},
				td.Name,
				td.GoName,
				td.SpannerType,
				td.ProtoName,
				NewNestedErrPrinter(m),
				td.Name,
				td.Name,
				td.SpannerTypeFieldName,
				td.Name,
				td.Name,
			)
		} else {
			p.PA([]string{
				"var %s %s\n{\nlocal := &%s{}\n",
				"if err := row.ColumnByName(\"%s\", local); err != nil {\n %s\n}\n",
				"if local.Valid {\n %s = local.%s\n}\n",
				"res.%s = %s\n}\n",
			},
				td.Name,
				td.GoName,
				td.SpannerType,
				td.ProtoName,
				NewNestedErrPrinter(m),
				td.Name,
				td.SpannerTypeFieldName,
				td.Name,
				td.Name,
			)
		}
	}
	// close nested function
	if m.method.IsUnary() {
		p.P("})\nif err != nil {\n%s\n}\n", NewErrPrinter(m, "error in closure: %v"))
	} else if m.method.IsClientStreaming() {
		p.P("}\n") // close the if row != nil check
	}

}
func (m *MethodStringer) PrintAfterHook(p *Printer) {
	after := m.method.GetMethodOption().GetAfter()
	if after == nil {
		return
	}
	if m.method.IsClientStreaming() {
		p.PA([]string{
			"// NOTE: I dont want to store your requests in memory\n",
			"// So the after hook is called with an empty request\n",
			"req := &%s{}\n",
		},
			m.method.GetInputType(),
		)
	}
	p.P("if err := ")
	if m.method.GetGoPackage(after.GetPackage()) != "" {
		p.P("%s.", m.method.GetGoPackage(after.GetPackage()))
	}
	if m.method.IsServerStreaming() {
		p.P(
			"%s(req, &res); err != nil {\n iterErr = gstatus.Errorf(codes.Unknown, %s)\nreturn\n}\n",
			after.GetName(),
			"\"error in after hook: %v\", err",
		)
	} else {
		p.P(
			"%s(req, &res); err != nil {\n%s\n}\n",
			after.GetName(),
			NewErrPrinter(m, "error in after hook: %v"),
		)
	}
}

func (m *MethodStringer) PrintReturnOutput(p *Printer) {
	if m.method.IsServerStreaming() {
		p.PA([]string{
			"if err := stream.Send(&res); err != nil {\n %s\n}\n})\n",
			"if err != nil {\n %s} else if iterErr != nil {\n",
			"return iterErr\n}\n",
			"return nil\n}\n",
		},
			NewNestedErrPrinter(m),
			NewErrPrinter(m, "error during iteration: %v"),
		)
	} else if m.method.IsClientStreaming() {
		p.P(
			"if err := stream.SendAndClose(&res); err != nil {\n %s\n}\nreturn nil\n}\n",
			NewErrPrinter(m, "error sending response: %v"),
		)
	} else if m.method.IsUnary() {
		p.P("return &res, nil\n}\n")
	}
}

func (m *MethodStringer) PrintForwardedMethod(p *Printer) {
	p.P("return s.FORWARDED.%s", m.method.GetName())
	if m.method.IsUnary() {
		p.P("(ctx, req)")
	} else if m.method.IsServerStreaming() {
		p.P("(req, stream)")
	} else {
		p.P("(stream)")
	}
	p.P("\n}\n")
}

type ErrPrinter struct {
	m   *Method
	msg string
}

func NewErrPrinter(m *MethodStringer, msg string) ErrPrinter {
	return ErrPrinter{m: m.method, msg: msg}
}

func (e ErrPrinter) String() string {
	if e.m.IsUnary() {
		return fmt.Sprintf("return nil, gstatus.Errorf(codes.Unknown, \"%s\", err)", e.msg)
	} else {
		return fmt.Sprintf("return gstatus.Errorf(codes.Unknown, \"%s\", err)", e.msg)
	}
}

type NestedErrPrinter struct {
	m *Method
}

func NewNestedErrPrinter(m *MethodStringer) NestedErrPrinter {
	return NestedErrPrinter{m: m.method}
}

func (e NestedErrPrinter) String() string {
	p := &Printer{}
	if e.m.IsUnary() || e.m.IsServerStreaming() {
		p.P("iterErr = gstatus.Errorf(codes.Unknown, \"%v\", err)\n", "couldnt scan out message err: %v")
	} else if e.m.IsClientStreaming() {
		p.P("%s", ErrPrinter{m: e.m, msg: "couldnt scan out message, err: %v"})
	}
	return p.String()
}

// the stream struct name used by streaming grpc methods
type StreamType struct {
	m *Method
}

func NewStreamType(m *Method) StreamType {
	return StreamType{m: m}
}
func (s StreamType) String() string {
	p := &Printer{}
	p.P(
		"%s%s_%sServer",
		s.m.GetFilePackage(),
		s.m.GetServiceName(),
		s.m.GetName(),
	)
	return p.String()
}

type PLInputName struct {
	m *Method
}

// name of the structs that are inputs to a method
// these need to be specific to the package and service.
// Each service can have their own type mapping of the same message.
// if the mappings of a message don't match then they will override
// eachother unless they are specific to the service.
func NewPLInputName(m *Method) PLInputName {
	return PLInputName{m: m}
}

func (p PLInputName) String() string {
	pr := &Printer{}
	spl := strings.Split(p.m.GetInputType(), ".")
	name := strings.Title(spl[0])
	spl = spl[1:]
	for _, s := range spl {
		name += "_" + strings.Title(s)
	}
	pr.P("%sFor%s", name, p.m.Service.GetName())
	return pr.String()
}

// the method name that returns our stringed query
type PLQueryMethodName struct {
	m *Method
}

func NewPLQueryMethodName(m *Method) PLQueryMethodName {
	return PLQueryMethodName{m: m}
}

func (pl PLQueryMethodName) String() string {
	p := &Printer{}
	p.P("%sFrom%sQuery", pl.m.GetInputTypeName(), pl.m.Desc.GetName())
	return p.String()
}

// name of the method reciever struct in the persist lib
// implements all persist methods on a service
type PersistHelperName struct {
	s *Service
}

func NewPersistHelperName(s *Service) PersistHelperName {
	return PersistHelperName{s: s}
}

func (p PersistHelperName) String() string {
	return p.s.GetName() + "MethodReceiver"
}

// name of the handler on the persist helper struct
type PersistHandlerName struct {
	m *Method
}

func NewPersistHandlerName(m *Method) PersistHandlerName {
	return PersistHandlerName{m: m}
}

func (p PersistHandlerName) String() string {
	return p.m.GetName() + "Handler"
}
