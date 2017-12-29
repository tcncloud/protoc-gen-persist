package generator

type PersistStringer struct{}

func (per *PersistStringer) MessageInputDeclaration(method *Method) string {
	printer := &Printer{}
	printer.P("type %s struct{\n", NewPLInputName(method))
	for _, qf := range method.GetTypeDescArrayForStruct(method.GetInputTypeStruct()) {
		printer.P("%s ", qf.Name)
		if qf.IsMapped {
			printer.P("interface{}\n")
		} else if qf.IsMessage && qf.IsRepeated {
			printer.P("[][]byte\n")
		} else if qf.IsMessage {
			printer.P("[]byte\n")
		} else {
			printer.P("%s\n", qf.GoName) // go name has package + repeated + type
		}
	}
	printer.P("}\n")

	return printer.String()
}

// merges custom defined handlers with our own
func (per *PersistStringer) PersistImplBuilder(service *Service) string {
	printer := &Printer{}
	printer.P(
		"type %sImpl struct{\nPERSIST *persist_lib.%s\nFORWARDED RestOf%sHandlers\n}\n",
		service.GetName(),
		NewPersistHelperName(service),
		service.GetName(),
	)
	printer.P("type RestOf%sHandlers interface{\n", service.GetName())
	for _, m := range *service.Methods {
		if !m.IsSpanner() || m.GetMethodOption() == nil {
			str := MethodStringer{method: m}
			str.PrintInterfaceDefinition(printer)
		}
	}
	printer.P("}\n")
	printer.PA([]string{
		"type %sImplBuilder struct {\n",
		"err error\n ",
		"rest RestOf%sHandlers\n",
		"queryHandlers *persist_lib.%sQueryHandlers\n",
		"i *%sImpl\n",
		"db *spanner.Client\n}\n",
		"func New%sBuilder() *%sImplBuilder {\nreturn &%sImplBuilder{i: &%sImpl{}}\n}\n",
	},
		service.GetName(),
		service.GetName(),
		service.GetName(),
		service.GetName(),
		service.GetName(), service.GetName(), service.GetName(), service.GetName(),
	)
	printer.PA([]string{
		"func (b *%sImplBuilder) WithRestOfGrpcHandlers(r RestOf%sHandlers) *%sImplBuilder {\n",
		"b.rest = r\n return b\n}\n",
	},
		service.GetName(),
		service.GetName(),
		service.GetName(),
	)
	printer.PA([]string{
		"func (b *%sImplBuilder) WithPersistQueryHandlers(p *persist_lib.%sQueryHandlers)",
		"*%sImplBuilder {\n",
		"b.queryHandlers = p\n return b\n}\n",
	},
		service.GetName(),
		service.GetName(),
		service.GetName(),
	)

	// setup default query functions
	printer.PA([]string{
		"func (b *%sImplBuilder) WithDefaultQueryHandlers() *%sImplBuilder {\n",
		"accessor := persist_lib.NewSpannerClientGetter(b.db)\n",
		"queryHandlers := &persist_lib.%sQueryHandlers{\n",
	},
		service.GetName(), service.GetName(),
		service.GetName(),
	)
	for _, m := range *service.Methods {
		if m.GetMethodOption() == nil {
			continue
		}
		printer.P("%s: persist_lib.Default%s(accessor),\n", NewPersistHandlerName(m), NewPersistHandlerName(m))
	}
	printer.P("}\n b.queryHandlers = queryHandlers\n return b\n}\n")

	printer.PA([]string{
		"func (b *%sImplBuilder) WithSpannerClient(c *spanner.Client) *%sImplBuilder {\n",
		"b.db = c\n return b\n}\n",
	},
		service.GetName(), service.GetName(),
	)
	printer.PA([]string{
		"func (b *%sImplBuilder) WithSpannerURI(ctx context.Context, uri string) *%sImplBuilder {\n",
		"cli, err := spanner.NewClient(ctx, uri)\n b.err = err\n b.db = cli\n return b\n}\n",
	},
		service.GetName(), service.GetName(),
	)
	printer.PA([]string{
		"func (b *%sImplBuilder) Build() (*%sImpl, error) {\n",
		"if b.err != nil {\n return nil, b.err\n}\n",
		"b.i.PERSIST = &persist_lib.%s{Handlers: *b.queryHandlers}\n",
		"b.i.FORWARDED = b.rest\n",
		"return b.i, nil\n}\n",
	},
		service.GetName(), service.GetName(),
		NewPersistHelperName(service),
	)

	return printer.String()
}

func (per *PersistStringer) HandlersStructDeclaration(service *Service) string {
	printer := &Printer{}
	// contains our query handlers struct, and is reciever of our methods
	printer.P(
		"type %s struct{\nHandlers %sQueryHandlers}\n",
		NewPersistHelperName(service),
		service.GetName(),
	)
	// actually runs the queries
	printer.P("type %sQueryHandlers struct {\n", service.GetName())
	for _, method := range *service.Methods {
		if method.GetMethodOption() == nil {
			continue
		}
		if method.IsClientStreaming() {
			printer.P(
				"%s func(context.Context)(func(*%s), func() (*spanner.Row, error))\n",
				NewPersistHandlerName(method),
				NewPLInputName(method),
			)
		} else {
			printer.P(
				"%s func(context.Context, *%s, func(*spanner.Row)) error\n",
				NewPersistHandlerName(method),
				NewPLInputName(method),
			)
		}
	}
	printer.P("}\n")
	return printer.String()
}

func (per *PersistStringer) HelperFunctionImpl(service *Service) string {
	printer := &Printer{}
	for _, method := range *service.Methods {
		if method.GetMethodOption() == nil {
			continue // we do not have any persist options
		}
		if method.IsClientStreaming() {
			printer.PA([]string{
				"// given a context, returns two functions.  (feed, stop)\n",
				"// feed will be called once for every row recieved by the handler\n",
				"// stop will be called when the client is done streaming. it expects\n",
				"//a  *spanner.Row to be returned, or nil.\n",
				"func (p *%s) %s(ctx context.Context)(func(*%s), func() (*spanner.Row, error)) {\n",
				"return p.Handlers.%s(ctx)\n}\n",
			},
				NewPersistHelperName(service),
				method.GetName(),
				NewPLInputName(method),
				NewPersistHandlerName(method),
			)
		} else {
			printer.PA([]string{
				"// next must be called on each result row\n",
				"func(p *%s) %s(ctx context.Context, params *%s, next func(*spanner.Row)) error {\n",
				"return p.Handlers.%s(ctx, params, next)\n}\n",
			},
				NewPersistHelperName(service),
				method.GetName(),
				NewPLInputName(method),
				NewPersistHandlerName(method),
			)
		}
	}
	return printer.String()
}

func (per *PersistStringer) QueryInterfaceDefinition(method *Method) string {
	return ""
}

func (per *PersistStringer) QueryFunction(method *Method) string {
	// we do not have a persist query
	if method.GetMethodOption() == nil {
		return ""
	}
	printer := &Printer{}
	if method.IsSelect() {
		printer.P(
			"func %s(req *%s) spanner.Statement {\nreturn %s\n}\n",
			NewPLQueryMethodName(method),
			NewPLInputName(method),
			method.Query,
		)
	} else {
		printer.P(
			"func %s(req *%s) *spanner.Mutation {\nreturn %s\n}\n",
			NewPLQueryMethodName(method),
			NewPLInputName(method),
			method.Query,
		)
	}
	return printer.String()
}

func (per *PersistStringer) DefaultFunctionsImpl(service *Service) string {
	printer := &Printer{}
	for _, method := range *service.Methods {
		if method.GetMethodOption() == nil {
			continue
		}
		if method.IsClientStreaming() {
			printer.PA([]string{
				"func Default%sHandler(accessor SpannerClientGetter) func(context.Context) ",
				"(func(*%s), func()(*spanner.Row, error)) {\n",
				"return func(ctx context.Context) (func(*%s), func()(*spanner.Row, error)) {\n",
				"var muts []*spanner.Mutation\n",
				"feed := func(req *%s) {\nmuts = append(muts, %s(req))\n}\n",
				"done := func() (*spanner.Row, error) {\n",
				"cli, err := accessor()\nif err != nil {\n return nil, err\n}\n",
				"if _, err := cli.Apply(ctx, muts); err != nil {\nreturn nil, err\n}\n",
				"return nil, nil // we dont have a row, because we are an apply\n",
				"}\n return feed, done\n}\n}\n",
			},
				method.GetName(),
				NewPLInputName(method),
				NewPLInputName(method),
				NewPLInputName(method),
				NewPLQueryMethodName(method),
			)
		} else {
			printer.PA([]string{
				"func Default%sHandler(accessor SpannerClientGetter) ",
				"func(context.Context, *%s, func(*spanner.Row)) error {\n",
				"return func(ctx context.Context, req *%s, next func(*spanner.Row)) error {\n",
			},
				method.GetName(),
				NewPLInputName(method),
				NewPLInputName(method),
			)
			printer.P("cli, err := accessor()\n if err != nil {\n return err\n}\n")

			if method.IsSelect() {
				printer.PA([]string{
					"iter := cli.Single().Query(ctx, %s(req))\n",
					"if err := iter.Do(func(r *spanner.Row) error {\n",
					"next(r)\nreturn nil\n}); err != nil {\nreturn err\n}\n",
				},
					NewPLQueryMethodName(method),
				)
			} else {
				printer.PA([]string{
					"if _, err := cli.Apply(ctx, []*spanner.Mutation{%s(req)}); err != nil {\n",
					"return err\n}\n next(nil) // this is an apply, it has no result\n",
				},
					NewPLQueryMethodName(method),
				)
			}
			printer.P("return nil\n}\n}\n")
		}
	}
	return printer.String()
}

func (per *PersistStringer) DeclareSpannerGetter() string {
	printer := &Printer{}
	printer.P("type SpannerClientGetter func() (*spanner.Client, error)\n")
	printer.PA([]string{
		"func NewSpannerClientGetter(cli *spanner.Client) SpannerClientGetter {\n",
		"return func() (*spanner.Client, error) {\n return cli, nil \n}\n}\n",
	})

	return printer.String()
}
