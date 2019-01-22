package generator

import (
	"strings"
	"text/template"

	_gen "github.com/golang/protobuf/protoc-gen-go/generator"
)

type printerProxy struct {
	printer *Printer
}

type handlerParams struct {
	Service        string
	Method         string
	Query          string
	Request        string
	Response       string
	RespMethodCall string
	ZeroResponse   bool
	Before         bool
	After          bool
}

func OneOrZero(hp handlerParams) string {
	if hp.ZeroResponse {
		return strings.Join([]string{`
err := result.Zero()
res := &`, hp.Response, `{}
        `}, "")
	}
	return "res, err := result.One()." + hp.RespMethodCall + "()"
}

func (h *printerProxy) Write(data []byte) (int, error) {
	h.printer.Q(string(data))
	return len(data), nil
}

func NewPrinterProxy(printer *Printer) *printerProxy {
	return &printerProxy{
		printer: printer,
	}
}

func WritePersistServerStruct(printer *Printer, service, db string) error {
	printerProxy := NewPrinterProxy(printer)
	structFormat := `
type Opts_{{.Service}} struct {
    MAPPINGS TypeMappings_{{.Service}}
    HOOKS    Hooks_{{.Service}}
}

func Opts{{.Service}}(hooks Hooks_{{.Service}}, mappings TypeMappings_{{.Service}}) Opts_{{.Service}} {
    opts := Opts_{{.Service}}{
        HOOKS: &DefaultHooks_{{.Service}}{},
        MAPPINGS: &DefaultTypeMappings_{{.Service}}{},
    }
    if hooks != nil {
        opts.HOOKS = hooks
    }
    if mappings != nil {
        opts.MAPPINGS = mappings
    }
    return opts
}


type Impl_{{.Service}} struct {
    opts    *Opts_{{.Service}}
    QUERIES *Queries_{{.Service}}
    HANDLERS RestOfHandlers_{{.Service}}
    DB      *{{.DB}}
}

func Impl{{.Service}}(db *{{.DB}}, handlers RestOfHandlers_{{.Service}}, opts ...Opts_{{.Service}}) *Impl_{{.Service}} {
    var myOpts Opts_{{.Service}}
    if len(opts) > 0 {
        myOpts = opts[0]
    } else {
        myOpts = Opts{{.Service}}(&DefaultHooks_{{.Service}}{}, &DefaultTypeMappings_{{.Service}}{})
    }
    return &Impl_{{.Service}}{
        opts:    &myOpts,
        QUERIES: Queries{{.Service}}(myOpts),
        DB:      db,
        HANDLERS: handlers,
    }
}
    `
	t := template.Must(template.New("PersistServerStruct").Parse(structFormat))
	return t.Execute(printerProxy, map[string]string{
		"Service": service,
		"DB":      db,
	})
}

func WriteClientStreaming(printer *Printer, params *handlerParams, isSql bool) error {
	printerProxy := NewPrinterProxy(printer)
	sqlClientStreamingFormat := `
func (this *Impl_{{.Service}}) {{.Method}}(stream {{.Service}}_{{.Method}}Server) error {
    tx, err := DefaultClientStreamingPersistTx(stream.Context(), this.DB)
    if err != nil {
        return gstatus.Errorf(codes.Unknown, "error creating persist tx: %v", err)
    }
    if err := this.{{.Method}}Tx(stream, tx); err != nil {
        return gstatus.Errorf(codes.Unknown, "error executing '{{.Query}}' query: %v", err)
    }
    return nil
}

func (this *Impl_{{.Service}}) {{.Method}}Tx(stream {{.Service}}_{{.Method}}Server, tx PersistTx) error {
    query := this.QUERIES.{{camelCase .Query}}(stream.Context(), tx)
    var first *{{.Request}}
    for {
        req, err := stream.Recv()
        if err == io.EOF {
            break
        } else if err != nil {
            return gstatus.Errorf(codes.Unknown, "error receiving request: %v", err)
        }
        if first == nil {
            first = req
        }
        {{if .Before}}
        {
            beforeRes, err := this.opts.HOOKS.{{.Method}}BeforeHook(stream.Context(), req)
            if err != nil {
                return gstatus.Errorf(codes.Unknown, "error in before hook: %v", err)
            } else if beforeRes != nil {
                continue
            }
        }
        {{end}}
        result := query.Execute(req)
        if err := result.Zero(); err != nil {
            return err
        }
    }
    if err := tx.Commit(); err != nil {
        if rollbackErr := tx.Rollback(); rollbackErr != nil {
            return fmt.Errorf("error executing '{{.Query}}' query :::AND COULD NOT ROLLBACK::: rollback err: %v, query err: %v", rollbackErr, err)
        }
    }
    res := &{{.Response}}{}

    {{if .After}}
    {
        if err := this.opts.HOOKS.{{.Method}}AfterHook(stream.Context(), first, res); err != nil {
            return gstatus.Errorf(codes.Unknown, "error in after hook: %v", err)
        }
    }
    {{end}}
    if err := stream.SendAndClose(res); err != nil {
        return gstatus.Errorf(codes.Unknown, "error sending back response: %v", err)
    }

    return nil
}
        `

	spannerClientStreamingFormat := `
func (this *Impl_{{.Service}}) {{.Method}}(stream {{.Service}}_{{.Method}}Server) error {
    if err := this.{{.Method}}Tx(stream); err != nil {
        return gstatus.Errorf(codes.Unknown, "error executing '{{.Query}}' query: %v", err)
    }
    return nil
}

func (this *Impl_{{.Service}}) {{.Method}}Tx(stream {{.Service}}_{{.Method}}Server) error {
    items := make([]*{{.Request}}, 0)
    var first *{{.Request}}
    for {
        req, err := stream.Recv()
        if err == io.EOF {
            break
        } else if err != nil {
            return gstatus.Errorf(codes.Unknown, "error receiving request: %v", err)
        }
        if first == nil {
            first = req
        }
        {{if .Before}}
        {
            beforeRes, err := this.opts.HOOKS.{{.Method}}BeforeHook(stream.Context(), req)
            if err != nil {
                return gstatus.Errorf(codes.Unknown, "error in before hook: %v", err)
            } else if beforeRes != nil {
                continue
            }
        }
        {{end}}

        items = append(items, req)
    }

    _, err := this.DB.ReadWriteTransaction(stream.Context(), func(ctx context.Context, tx *spanner.ReadWriteTransaction) error {
        for _, item := range items {
            query := this.QUERIES.{{camelCase .Query}}(ctx, tx)
            result := query.Execute(item)
            if err := result.Zero(); err != nil {
                return err
            }
        }
        return nil
    })
	if err != nil {
		return gstatus.Errorf(codes.Unknown, "error in read write transaction: %v", err)
	}

    res := &{{.Response}}{}

    {{if .After}}
    {
        if err := this.opts.HOOKS.{{.Method}}AfterHook(stream.Context(), first, res); err != nil {
            return gstatus.Errorf(codes.Unknown, "error in after hook: %v", err)
        }
    }
    {{end}}
    if err := stream.SendAndClose(res); err != nil {
        return gstatus.Errorf(codes.Unknown, "error sending back response: %v", err)
    }

    return nil
}
        `

	var clientStreamingFormat string
	if isSql {
		clientStreamingFormat = sqlClientStreamingFormat
	} else {
		clientStreamingFormat = spannerClientStreamingFormat
	}

	funcMap := template.FuncMap{
		"camelCase": _gen.CamelCase,
	}
	t := template.Must(template.New("ClientStreaming").Funcs(funcMap).Parse(clientStreamingFormat))
	return t.Execute(printerProxy, params)
}

func WriteUnary(printer *Printer, params *handlerParams, isSql bool) error {
	printerProxy := NewPrinterProxy(printer)
	sqlUnaryFormat := `
func (this *Impl_{{.Service}}) {{.Method}}(ctx context.Context, req *{{.Request}}) (*{{.Response}}, error) {
    query := this.QUERIES.{{camelCase .Query}}(ctx, this.DB)
    {{if .Before}}
    {
        beforeRes, err := this.opts.HOOKS.{{.Method}}BeforeHook(ctx, req)
        if err != nil {
            return nil, gstatus.Errorf(codes.Unknown, "error in before hook: %v", err)
        } else if beforeRes != nil {
            return beforeRes, nil
        }
    }
    {{end}}

    result := query.Execute(req)
    {{oneOrZero .}}
    if err != nil {
        return nil, err
    }

    {{if .After}}
    {
        if err := this.opts.HOOKS.{{.Method}}AfterHook(ctx, req, res); err != nil {
            return nil, gstatus.Errorf(codes.Unknown, "error in after hook: %v", err)
        }
    }
    {{end}}

    return res, nil
}
  `
	spannerUnaryFormat := `
func (this *Impl_{{.Service}}) {{.Method}}(ctx context.Context, req *{{.Request}}) (*{{.Response}}, error) {
    query := this.QUERIES.{{camelCase .Query}}(ctx, this.DB.Single())
    {{if .Before}}
    beforeRes, err := this.opts.HOOKS.{{.Method}}BeforeHook(ctx, req)
    if err != nil {
        return nil, gstatus.Errorf(codes.Unknown, "error in before hook: %v", err)
    } else if beforeRes != nil {
        return beforeRes, nil
    }
    {{end}}

    result := query.Execute(req)
    {{oneOrZero .}}
    if err != nil {
        return nil, err
    }

    {{if .After}}
    {
        if err := this.opts.HOOKS.{{.Method}}AfterHook(ctx, req, res); err != nil {
            return nil, gstatus.Errorf(codes.Unknown, "error in after hook: %v", err)
        }
    }
    {{end}}

    return res, nil
}
    `
	funcMap := template.FuncMap{
		"camelCase": _gen.CamelCase,
		"oneOrZero": OneOrZero,
	}

	var unaryFormat string
	if isSql {
		unaryFormat = sqlUnaryFormat
	} else {
		unaryFormat = spannerUnaryFormat
	}

	t := template.Must(template.New("UnaryRequest").Funcs(funcMap).Parse(unaryFormat))
	return t.Execute(printerProxy, params)
}

func WriteServerStream(printer *Printer, params *handlerParams, isSql bool) error {
	printerProxy := NewPrinterProxy(printer)
	sqlServerFormat := `
func (this *Impl_{{.Service}}) {{.Method}}(req *{{.Request}}, stream {{.Service}}_{{.Method}}Server) error {
    tx, err := DefaultServerStreamingPersistTx(stream.Context(), this.DB)
    if err != nil {
        return gstatus.Errorf(codes.Unknown, "error creating persist tx: %v", err)
    }
    if err := this.{{.Method}}Tx(req, stream, tx); err != nil {
        return gstatus.Errorf(codes.Unknown, "error executing '{{.Query}}' query: %v", err)
    }
    return nil
}

func (this *Impl_{{.Service}}) {{.Method}}Tx(req *{{.Request}}, stream {{.Service}}_{{.Method}}Server, tx PersistTx) error {
    ctx := stream.Context()
    query := this.QUERIES.{{camelCase .Query}}(ctx, tx)

    iter := query.Execute(req)
    return iter.Each(func(row *Row_{{.Service}}_{{camelCase .Query}}) error {
        res, err := row.{{.RespMethodCall}}()
        if err != nil {
            return err
        }
        return stream.Send(res)
    })
}
  `

	spannerServerFormat := `
func (this *Impl_{{.Service}}) {{.Method}}(req *{{.Request}}, stream {{.Service}}_{{.Method}}Server) error {
    if err := this.{{.Method}}Tx(req, stream, this.DB.Single()); err != nil {
        return gstatus.Errorf(codes.Unknown, "error executing '{{.Query}}' query: %v", err)
    }
    return nil
}

func (this *Impl_{{.Service}}) {{.Method}}Tx(req *{{.Request}}, stream {{.Service}}_{{.Method}}Server, tx PersistTx) error {
    ctx := stream.Context()
    query := this.QUERIES.{{camelCase .Query}}(ctx, tx)

    iter := query.Execute(req)
    return iter.Each(func(row *Row_{{.Service}}_{{camelCase .Query}}) error {
        res, err := row.{{.RespMethodCall}}()
        if err != nil {
            return err
        }
        return stream.Send(res)
    })
}
    `
	funcMap := template.FuncMap{
		"camelCase": _gen.CamelCase,
	}

	var serverFormat string
	if isSql {
		serverFormat = sqlServerFormat
	} else {
		serverFormat = spannerServerFormat
	}

	t := template.Must(template.New("ServerStream").Funcs(funcMap).Parse(serverFormat))
	return t.Execute(printerProxy, params)
}
