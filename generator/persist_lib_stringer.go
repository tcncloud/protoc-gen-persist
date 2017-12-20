package generator

type PersistStringer struct {
	file *FileStruct
}

func (per *PersistStringer) MessageInputDeclarations() string {
	printer := &Printer{}
	processed := make(map[string]bool)
	for _, service := range *per.file.ServiceList {
		for _, method := range *service.Methods {
			if !processed[method.GetInputType()] {
				processed[method.GetInputType()] = true
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
			}
		}
	}

	return printer.String()
}

func (per *PersistStringer) HandlersDeclarations() string {
	printer := &Printer{}
	for _, service := range *per.file.ServiceList {
		if !service.IsSpanner() {
			continue
		}
		printer.P(
			"type %s struct{\nHandlers %sHandlers}\n",
			NewPersistHelperName(service),
			service.GetName(),
		)
		printer.P("type %sHandlers struct {\n", service.GetName())
		for _, method := range *service.Methods {
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
	}
	return printer.String()
}

func (per *PersistStringer) HelperHandlers() string {
	printer := &Printer{}
	for _, service := range *per.file.ServiceList {
		if !service.IsSpanner() {
			continue
		}
		for _, method := range *service.Methods {
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
	}
	return printer.String()
}

func (per *PersistStringer) QueryFunctions() string {
	printer := &Printer{}
	for _, service := range *per.file.ServiceList {
		if !service.IsSpanner() {
			continue
		}
		for _, method := range *service.Methods {
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
		}
	}
	return printer.String()
}

func (per *PersistStringer) DefaultFunctions() string {
	printer := &Printer{}
	for _, service := range *per.file.ServiceList {
		if !service.IsSpanner() {
			continue
		}
		for _, method := range *service.Methods {
			if method.IsClientStreaming() {
				printer.PA([]string{
					"func Default%sHandler(cli *spanner.Client) func(context.Context) ",
					"(func(*%s), func()(*spanner.Row, error)) {\n",
					"return func(ctx context.Context) (func(*%s), func()(*spanner.Row, error)) {\n",
					"var muts []*spanner.Mutation\n",
					"feed := func(req *%s) {\nmuts = append(muts, %s(req))\n}\n",
					"done := func() (*spanner.Row, error) {\n",
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
					"func Default%sHandler(cli *spanner.Client) ",
					"func(context.Context, *%s, func(*spanner.Row)) error {\n",
					"return func(ctx context.Context, req *%s, next func(*spanner.Row)) error {\n",
				},
					method.GetName(),
					NewPLInputName(method),
					NewPLInputName(method),
				)
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
	}
	return printer.String()
}
