package generator

type UnaryStringer struct {
	printer *Printer
	method  *Method
	backend BackendStringer
}

func NewUnaryStringer(method *Method, b BackendStringer) *UnaryStringer {
	return &UnaryStringer{
		printer: &Printer{},
		method:  method,
		backend: b,
	}
}

func (s *UnaryStringer) String() string {
	if s.method.GetMethodOption() == nil {
		return s.Unimplemented()
	}
	// print my scope
	// print child level scopes
	s.printer.P(
		"func (s *%sImpl) %s(ctx context.Context, req *%s) (*%s, error) {\n",
		s.method.Service.GetName(),
		s.method.GetName(),
		s.method.GetInputType(),
		s.method.GetOutputType(),
	)
	s.printer.P("var err error\nvar res = &%s{}\n _ = err\n_ = res\n", s.method.GetOutputType())
	s.BeforeHook()
	s.Params()
	s.PersistCall()
	s.AfterHook()
	s.printer.P("return res, nil\n}\n")

	return s.printer.String()
}
func (s *UnaryStringer) Unimplemented() string {
	s.printer.P(
		"func (s *%sImpl) %s(ctx context.Context, req *%s) (*%s, error) {\n",
		s.method.Service.GetName(),
		s.method.GetName(),
		s.method.GetInputType(),
		s.method.GetOutputType(),
	)
	s.printer.P("return s.FORWARDED.%s(ctx, req)\n}\n", s.method.GetName())

	return s.printer.String()
}

func (s *UnaryStringer) BeforeHook() {
	opts := s.method.GetMethodOption()
	if opts == nil {
		return
	}
	before := opts.GetBefore()
	if before == nil {
		return
	}
	hookName := s.method.GetHookName(before)

	s.printer.PA([]string{
		"beforeRes, err := %s(req)\n",
		"if err != nil {\n return nil, gstatus.Errorf(codes.Unknown, \"%s\", err)\n}",
		" else if beforeRes != nil {\n return beforeRes, nil\n}\n",
	},
		hookName,
		"error in before hook: %v",
	)
}

func (s *UnaryStringer) Params() {
	s.printer.P("params, err := %s(req)\n", ToParamsFuncName(s.method))
	s.printer.P("if err != nil {\n return nil, err\n}\n")
}

func (s *UnaryStringer) PersistCall() {
	s.printer.P("var iterErr error\n")
	s.printer.P("err = s.PERSIST.%s(ctx, params, ", s.method.GetName())
	s.HandleRow()
	s.printer.PA([]string{
		")\n", // closes the persist function call
		"if err != nil {\n return nil, gstatus.Errorf(codes.Unknown, \"%s\", err) \n}",
		"else if iterErr != nil {\n return nil, iterErr \n}\n",
	},
		"error calling persist service: %v",
	)
}

func (s *UnaryStringer) HandleRow() {
	s.printer.P("func(row %s) {\n", s.backend.RowType())
	s.printer.P("if row == nil { // there was no return data\n return\n}\n")
	s.ResultFromRow()
	// dont include the newline on purpose
	s.printer.P("}")
}

func (s *UnaryStringer) ResultFromRow() {
	s.printer.P("res, err = %s(row)\n", FromScanableFuncName(s.method))
	s.printer.P("if err != nil {\n iterErr = err\n return\n}\n")
}

func (s *UnaryStringer) AfterHook() {
	after := s.method.GetMethodOption().GetAfter()
	if after == nil {
		return
	}
	s.printer.PA([]string{
		"if err := %s(req, res); err != nil {\n",
		"return nil, gstatus.Errorf(codes.Unknown, \"%s\", err)\n}\n",
	},
		s.method.GetHookName(after),
		"error in after hook: %v",
	)
}

type BidiStreamStringer struct {
	printer *Printer
	method  *Method
	backend BackendStringer
}

func NewBidiStreamStringer(method *Method, b BackendStringer) *BidiStreamStringer {
	return &BidiStreamStringer{
		printer: &Printer{},
		method:  method,
		backend: b,
	}
}

func (s *BidiStreamStringer) String() string {
	bidiSpanner := s.method.IsBidiStreaming() && s.method.IsSpanner()
	if s.method.GetMethodOption() == nil || bidiSpanner {
		return s.Unimplemented()
	}

	s.printer.P(
		"func (s *%sImpl) %s(stream %s) error {\n",
		s.method.Service.GetName(),
		s.method.GetName(),
		NewStreamType(s.method),
	)
	s.printer.P("var err error\n _ = err\n")

	s.PersistCall()

	s.printer.PA([]string{
		"for {\n req, err := stream.Recv()\n",
		"if err == io.EOF {\nbreak\n} else if err != nil {\n",
		" return gstatus.Errorf(codes.Unknown, \"%s\", err)\n}\n",
	},
		"error receiving request: %v",
	)
	s.BeforeHook()
	s.Params()
	s.HandleRow()
	s.printer.P("}\n")
	s.printer.P("return stop()\n}\n")

	return s.printer.String()
}
func (s *BidiStreamStringer) Unimplemented() string {
	srvName := s.method.Service.GetName()
	name := s.method.GetName()

	s.printer.PA([]string{
		"func (s *%sImpl) %s(stream %s) error {\n",
		"return s.FORWARDED.%s(stream)\n}\n",
	},
		srvName, name, NewStreamType(s.method),
		name,
	)
	return s.printer.String()
}

func (s *BidiStreamStringer) PersistCall() {
	s.printer.P(
		"feed, stop := s.PERSIST.%s(stream.Context())\n",
		s.method.GetName(),
	)
}
func (s *BidiStreamStringer) BeforeHook() {
	opts := s.method.GetMethodOption()
	if opts == nil {
		return
	}
	before := opts.GetBefore()
	if before == nil {
		return
	}
	hookName := s.method.GetHookName(before)
	s.printer.PA([]string{
		"beforeRes, err := %s(req)\n",
		"if err != nil {\n return gstatus.Errorf(codes.Unknown, \"%s\", err)\n} ",
		" else if beforeRes != nil {\ncontinue\n}\n",
	},
		hookName,
		"error in before hook: %v",
	)
}
func (s *BidiStreamStringer) Params() {
	s.printer.P("params, err := %s(req)\n", ToParamsFuncName(s.method))
	// s.printer.P("params := &persist_lib.%s{}\n", NewPLInputName(s.method))
	// s.printer.P("err = %s", s.backend.MapRequestToParams())
	s.printer.P("if err != nil {\n return err\n}\n")
}
func (s *BidiStreamStringer) HandleRow() {
	s.printer.PA([]string{
		"row, err := feed(params)\n",
		"if err != nil {\n return gstatus.Errorf(codes.Unknown, \"%s\", err)\n}\n",
		"if row != nil {\n",
	},
		"error receiving result row: %v",
	)
	s.printer.P("res, err := %s(row)\n", FromScanableFuncName(s.method))
	s.printer.P("if err != nil {\n return err \n}\n")
	s.AfterHook()
	s.printer.P("if err := stream.Send(res); err != nil {\n return err\n}\n")
	s.printer.P("}\n")
}
func (s *BidiStreamStringer) AfterHook() {
	opts := s.method.GetMethodOption()
	if opts == nil {
		return
	}
	after := opts.GetAfter()
	if after == nil {
		return
	}
	s.printer.PA([]string{
		"if err := %s(req, res); err != nil {\n",
		"return gstatus.Errorf(codes.Unknown, \"%s\", err)\n}\n",
	},
		s.method.GetHookName(after),
		"error in after hook: %v",
	)
}

type ClientStreamStringer struct {
	printer *Printer
	method  *Method
	backend BackendStringer
}

func NewClientStreamStringer(method *Method, b BackendStringer) *ClientStreamStringer {
	return &ClientStreamStringer{
		printer: &Printer{},
		method:  method,
		backend: b,
	}
}

func (s *ClientStreamStringer) String() string {
	if s.method.GetMethodOption() == nil {
		return s.Unimplemented()
	}

	s.printer.P(
		"func (s *%sImpl) %s(stream %s) error {\n",
		s.method.Service.GetName(),
		s.method.GetName(),
		NewStreamType(s.method),
	)
	s.printer.P("var err error\n _ = err\n")
	s.printer.P("res := &%s{}\n", s.method.GetOutputType())

	s.PersistCall()

	s.printer.PA([]string{
		"for {\n req, err := stream.Recv()\n",
		"if err == io.EOF {\nbreak\n} else if err != nil {\n",
		" return gstatus.Errorf(codes.Unknown, \"%s\", err)\n}\n",
	},
		"error receiving request: %v",
	)
	s.BeforeHook()
	s.Params()

	s.printer.P("if err := feed(params); err != nil {\nreturn err\n}\n")
	s.printer.P("}\n")

	s.HandleRow()
	s.AfterHook()
	s.printer.PA([]string{
		"if err := stream.SendAndClose(res); err != nil {\n",
		"return gstatus.Errorf(codes.Unknown, \"%s\", err)\n}\n",
		"return nil\n}\n",
	},
		"error sending back response: %v",
	)
	return s.printer.String()
}
func (s *ClientStreamStringer) Unimplemented() string {
	srvName := s.method.Service.GetName()
	name := s.method.GetName()

	s.printer.PA([]string{
		"func (s *%sImpl) %s(stream %s) error {\n",
		"return s.FORWARDED.%s(stream)\n}\n",
	},
		srvName, name, NewStreamType(s.method),
		name,
	)
	return s.printer.String()
}

func (s *ClientStreamStringer) PersistCall() {
	s.printer.P(
		"feed, stop, err := s.PERSIST.%s(stream.Context())\nif err != nil {\nreturn err\n}\n",
		s.method.GetName(),
	)
}
func (s *ClientStreamStringer) BeforeHook() {
	opts := s.method.GetMethodOption()
	if opts == nil {
		return
	}
	before := opts.GetBefore()
	if before == nil {
		return
	}
	hookName := s.method.GetHookName(before)
	s.printer.PA([]string{
		"beforeRes, err := %s(req)\n",
		"if err != nil {\n return gstatus.Errorf(codes.Unknown, \"%s\", err)\n} ",
		" else if beforeRes != nil {\ncontinue\n}\n",
	},
		hookName,
		"error in before hook: %v",
	)
}
func (s *ClientStreamStringer) Params() {
	s.printer.P("params, err := %s(req)\n", ToParamsFuncName(s.method))
	s.printer.P("if err != nil {\n return err\n}\n")
}
func (s *ClientStreamStringer) HandleRow() {
	s.printer.PA([]string{
		"row, err := stop()\n",
		"if err != nil {\n return gstatus.Errorf(codes.Unknown, \"%s\", err)\n}\n",
		"if row != nil {\n",
		"res, err = %s(row)\n if err != nil {\nreturn err\n}\n",
		"}\n",
	},
		"error receiving result row: %v",
		FromScanableFuncName(s.method),
	)
}
func (s *ClientStreamStringer) AfterHook() {
	after := s.method.GetMethodOption().GetAfter()
	if after == nil {
		return
	}
	s.printer.PA([]string{
		"// NOTE: I dont want to store your requests in memory\n",
		"// so the after hook for client streaming calls\n",
		"// is called with an empty request struct\n",
		"fakeReq := &%s{}\n",
		"if err := %s(fakeReq, res); err != nil {\n",
		"return gstatus.Errorf(codes.Unknown, \"%s\", err)\n}\n",
	},
		s.method.GetInputType(),
		s.method.GetHookName(after),
		"error in after hook: %v",
	)
}

type ServerStreamStringer struct {
	printer *Printer
	method  *Method
	backend BackendStringer
}

func NewServerStreamStringer(method *Method, b BackendStringer) *ServerStreamStringer {
	return &ServerStreamStringer{
		printer: &Printer{},
		method:  method,
		backend: b,
	}
}

func (s *ServerStreamStringer) String() string {
	if s.method.GetMethodOption() == nil {
		return s.Unimplemented()
	}
	// print my scope
	// print child level scopes
	s.printer.P(
		"func (s *%sImpl) %s(req *%s, stream %s) error{\n",
		s.method.Service.GetName(),
		s.method.GetName(),
		s.method.GetInputType(),
		NewStreamType(s.method),
	)
	s.printer.P("var err error\n _ = err\n")
	s.BeforeHook()
	s.Params()
	s.PersistCall()
	s.printer.P("return nil\n}\n")

	return s.printer.String()
}

func (s *ServerStreamStringer) Unimplemented() string {
	srvName := s.method.Service.GetName()
	name := s.method.GetName()
	in := s.method.GetInputType()

	s.printer.PA([]string{
		"func (s *%sImpl) %s(req *%s, stream %s) error {\n",
		"return s.FORWARDED.%s(req, stream)\n}\n",
	},
		srvName, name, in, NewStreamType(s.method),
		name,
	)
	return s.printer.String()
}

func (s *ServerStreamStringer) BeforeHook() {
	opts := s.method.GetMethodOption()
	if opts == nil {
		return
	}
	before := opts.GetBefore()
	if before == nil {
		return
	}
	hookName := s.method.GetHookName(before)
	s.printer.PA([]string{
		"beforeRes, err := %s(req)\n",
		"if err != nil {\n return gstatus.Errorf(codes.Unknown, \"%s\", err)\n} ",
		"else if beforeRes != nil {\n",
		"for _, res := range beforeRes {\n",
		"if err := stream.Send(res); err != nil {\n",
		"return gstatus.Errorf(codes.Unknown, \"%s\", err)\n}\n}\n}\n",
	},
		hookName,
		"error in before hook: %v",
		"error sending back before hook result: %v",
	)
}

func (s *ServerStreamStringer) Params() {
	s.printer.P("params, err := %s(req)\n", ToParamsFuncName(s.method))
	s.printer.P("if err != nil {\n return err\n}\n")
}

func (s *ServerStreamStringer) PersistCall() {
	s.printer.P("var iterErr error\n")
	s.printer.P("err = s.PERSIST.%s(stream.Context(), params, ", s.method.GetName())
	s.HandleRow()
	s.printer.PA([]string{
		")\n",
		"if err != nil {\n return gstatus.Errorf(codes.Unknown, \"%s\", err) \n}",
		" else if iterErr != nil {\n return iterErr \n}\n",
	},
		"error during iteration: %v",
	)
}

func (s *ServerStreamStringer) HandleRow() {
	s.printer.P("func(row %s) {\n", s.backend.RowType())
	s.printer.P("if row == nil { // there was no return data\n return\n}\n")
	s.ResultFromRow()
	s.AfterHook()
	s.printer.PA([]string{
		"if err := stream.Send(res); err != nil {\n",
		"iterErr = gstatus.Errorf(codes.Unknown, \"%s\", err)\n}\n",
	},
		"error during iteration: %v",
	)
	// dont include the newline on purpose
	s.printer.P("}")
}

func (s *ServerStreamStringer) ResultFromRow() {
	if len(s.method.GetTypeDescForFieldsInStruct(s.method.GetOutputTypeStruct())) > 0 {
		s.printer.P(
			"res, err := %s(row)\n if err != nil {\n iterErr = err\n return\n}\n",
			FromScanableFuncName(s.method),
		)
	}
}

func (s *ServerStreamStringer) AfterHook() {
	opts := s.method.GetMethodOption()
	if opts == nil {
		return
	}
	after := opts.GetAfter()
	if after == nil {
		return
	}
	s.printer.PA([]string{
		"if err := %s(req, res); err != nil {\n",
		"iterErr = gstatus.Errorf(codes.Unknown, \"%s\", err)\n return\n}\n",
	},
		s.method.GetHookName(after),
		"error in after hook: %v",
	)
}
