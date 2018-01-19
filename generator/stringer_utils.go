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
	p.P("%s", buff.String())
}

func (p *Printer) String() string {
	return p.str
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
		s.m.Service.GetName(),
		s.m.GetName(),
	)
	return p.String()
}

type PLOutputName struct {
	m *Method
}

func NewPLOutputName(m *Method) PLOutputName {
	return PLOutputName{m: m}
}

func (p PLOutputName) String() string {
	return PersistLibTypeString(p.m, p.m.GetOutputType())
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
	return PersistLibTypeString(p.m, p.m.GetInputType())
}
func PersistLibTypeString(m *Method, typ string) string {
	pr := &Printer{}
	spl := strings.Split(typ, ".")
	name := strings.Title(spl[0])
	spl = spl[1:]
	for _, s := range spl {
		name += "_" + strings.Title(s)
	}
	pr.P("%sFor%s", name, m.Service.GetName())
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

	p.P("%s%sQuery", pl.m.Service.GetName(), pl.m.Desc.GetName())
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

func ToParamsFuncName(m *Method) string {
	printer := &Printer{}
	printer.P("%sTo%sPersistType", m.GetInputTypeMinusPackage(), m.Service.GetName())
	return printer.String()
}

func FromScanableFuncName(m *Method) string {
	printer := &Printer{}
	printer.P("%sFrom%sDatabaseRow", m.GetOutputTypeMinusPackage(), m.Service.GetName())
	return printer.String()
}

func IterProtoName(m *Method) string {
	return fmt.Sprintf("Iter%s%sProto", m.Service.GetName(), m.GetOutputTypeMinusPackage())
}
func IterPersistName(m *Method) string {
	return fmt.Sprintf("Iter%s%sPersist", m.Service.GetName(), m.GetOutputTypeMinusPackage())
}
