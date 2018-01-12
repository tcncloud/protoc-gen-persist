package generator

import (
	"bytes"
	"fmt"
	"github.com/tcncloud/protoc-gen-persist/persist"
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

func GetHookName(hook *persist.QLImpl_CallbackFunction) string {
	var name string
	pkg := GetGoPackage(hook.GetPackage())
	if pkg != "" {
		name = pkg + "."
	}
	name += hook.GetName()
	return name
}
