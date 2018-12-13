package generator

import (
	"text/template"

	_gen "github.com/golang/protobuf/protoc-gen-go/generator"
)

type printerProxy struct {
	printer *Printer
}

type handlerParams struct {
	Service  string
	Method   string
	Query    string
	Request  string
	Response string
	Before   bool
	After    bool
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

func WritePersistServerStruct(printer *Printer, service string) error {
	printerProxy := NewPrinterProxy(printer)
	structFormat := `
type {{.Service}}_ImplOpts struct {
    MAPPINGS {{.Service}}_TypeMappings
    HOOKS    {{.Service}}_Hooks
}

func Default{{.Service}}ImplOpts() {{.Service}}_ImplOpts {
    return {{.Service}}_ImplOpts{}
}

type {{.Service}}_Impl struct {
    opts    *{{.Service}}_ImplOpts
    QUERIES *{{.Service}}_Queries
    DB      *sql.DB
}

func {{.Service}}PersistImpl(db *sql.DB, opts ...{{.Service}}_ImplOpts) *{{.Service}}_Impl {
    var myOpts {{.Service}}_ImplOpts
    if len(opts) > 0 {
        myOpts = opts[0]
    } else {
        myOpts = Default{{.Service}}ImplOpts()
    }
    return &{{.Service}}_Impl{
        opts:    &myOpts,
        QUERIES: {{.Service}}PersistQueries(db, {{.Service}}_QueryOpts{MAPPINGS: myOpts.MAPPINGS}),
        DB:      db,
    }
}
    `
	t := template.Must(template.New("PersistServerStruct").Parse(structFormat))
	return t.Execute(printerProxy, map[string]string{
		"Service": service,
	})
}

func WriteClientStreaming(printer *Printer, params *handlerParams) error {
	printerProxy := NewPrinterProxy(printer)
	clientStreamingFormat := `
func (this *{{.Service}}_Impl) {{.Method}}(stream {{.Service}}_{{.Method}}Server) error {
    tx, err := DefaultClientStreamingPersistTx(stream.Context(), this.DB)
    if err != nil {
        return gstatus.Errorf(codes.Unknown, "error creating persist tx: %v", err)
    }
    if err := this.{{.Method}}Tx(stream, tx); err != nil {
        return gstatus.Errorf(codes.Unknown, "error executing '{{.Query}}' query: %v", err)
    }
    return nil
}

func (this *{{.Service}}_Impl) {{.Method}}Tx(stream {{.Service}}_{{.Method}}Server, tx PersistTx) error {
    query := this.QUERIES.{{camelCase .Query}}Query(stream.Context())
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
        beforeRes, err := this.opts.HOOKS.{{.Method}}BeforeHook(stream.Context(), req)
        if err != nil {
            return gstatus.Errorf(codes.Unknown, "error in before hook: %v", err)
        } else if beforeRes != nil {
            continue
        }
        {{end}}
        result := query.Execute(req)
        if err := result.Zero(); err != nil {
            return gstatus.Errorf(codes.InvalidArgument, "client streaming queries must return zero results")
        }
    }
    if err := tx.Commit(); err != nil {
        return fmt.Errorf("executed '{{.Query}}' query without error, but received error on commit: %v", err)
        if rollbackErr := tx.Rollback(); rollbackErr != nil {
            return fmt.Errorf("error executing '{{.Query}}' query :::AND COULD NOT ROLLBACK::: rollback err: %v, query err: %v", rollbackErr, err)
        }
    }
    res := &Empty{}

    {{if .After}}
    if err := this.opts.HOOKS.{{.Method}}AfterHook(stream.Context(), first, res); err != nil {
        return gstatus.Errorf(codes.Unknown, "error in after hook: %v", err)
    }
    {{end}}
    if err := stream.SendAndClose(res); err != nil {
        return gstatus.Errorf(codes.Unknown, "error sending back response: %v", err)
    }

    return nil
}
        `

	funcMap := template.FuncMap{
		"camelCase": _gen.CamelCase,
	}
	t := template.Must(template.New("ClientStreaming").Funcs(funcMap).Parse(clientStreamingFormat))
	return t.Execute(printerProxy, params)
}

func WriteUnary(printer *Printer, params *handlerParams) error {
	printerProxy := NewPrinterProxy(printer)
	unaryFormat := `
func (this *{{.Service}}_Impl) {{.Method}}(ctx context.Context, req *{{.Request}}) (*{{.Response}}, error) {
    query := this.QUERIES.{{camelCase .Query}}Query(ctx)
    {{if .Before}}
    beforeRes, err := this.opts.HOOKS.{{.Method}}BeforeHook(ctx, req)
    if err != nil {
        return nil, gstatus.Errorf(codes.Unknown, "error in before hook: %v", err)
    } else if beforeRes != nil {
        return gstatus.Error(codes.Unknown, "before hook returned nil")
    }
    req = beforeRes
    {{end}}

    result := query.Execute(req)
    res, err := result.One().{{.Response}}()
    if err != nil {
        return nil, err
    }

    {{if .After}}
    if err := this.opts.HOOKS.{{.Method}}AfterHook(ctx, req, res); err != nil {
        return nil, gstatus.Errorf(codes.Unknown, "error in after hook: %v", err)
    }
    {{end}}

    return res, nil
}
    `
	funcMap := template.FuncMap{
		"camelCase": _gen.CamelCase,
	}
	t := template.Must(template.New("UnaryRequest").Funcs(funcMap).Parse(unaryFormat))
	return t.Execute(printerProxy, params)
}

func WriteSeverStream(printer *Printer, params *handlerParams) error {
	printerProxy := NewPrinterProxy(printer)
	serverFormat := `
func (this *{{.Service}}_Impl) {{.Method}}(req *{{.Request}}, stream {{.Service}}_{{.Method}}Server) error {
    tx, err := DefaultServerStreamingPersistTx(stream.Context(), this.DB)
    if err != nil {
        return gstatus.Errorf(codes.Unknown, "error creating persist tx: %v", err)
    }
    if err := this.{{.Method}}Tx(req, stream, tx); err != nil {
        return gstatus.Errorf(codes.Unknown, "error executing '{{.Query}}' query: %v", err)
    }
    return nil
}

func (this *{{.Service}}_Impl) {{.Method}}Tx(req *{{.Request}}, stream {{.Service}}_{{.Method}}Server, tx PersistTx) error {
    ctx := stream.Context()
    query := this.QUERIES.{{camelCase .Query}}Query(ctx)

    iter := query.Execute(req)
    return iter.Each(func(row *{{.Service}}_{{camelCase .Query}}Row) error {
        res, err := row.{{.Response}}()
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
	t := template.Must(template.New("ServerStream").Funcs(funcMap).Parse(serverFormat))
	return t.Execute(printerProxy, params)
}

func WriteBidirectionalStream(printer *Printer, params *handlerParams) error {
	printerProxy := NewPrinterProxy(printer)
	biFormat := `
func (this *{{.Service}}_Impl) {{.Method}}(stream {{.Service}}_{{.Method}}Server) error {
    tx, err := DefaultBidiStreamingPersistTx(stream.Context(), this.DB)
    if err != nil {
        return gstatus.Errorf(codes.Unknown, "error creating persist tx: %v", err)
    }
    if err := this.{{.Method}}Tx(stream, tx); err != nil {
        return gstatus.Errorf(codes.Unknown, "error executing '{{.Query}}' query: %v", err)
    }
    return nil
}

func (this *{{.Service}}_Impl) {{.Method}}Tx(stream {{.Service}}_{{.Method}}Server, tx PersistTx) error {
    ctx := stream.Context()
    for {
        req, err := stream.Recv()
        if err == io.EOF {
            err = tx.Commit()
            if err != nil {
                return tx.Rollback()
            }
            return nil
        } else if err != nil {
            return gstatus.Errorf(codes.Unknown, "error receiving request: %v", err)
        }
        iter := this.QUERIES.{{camelCase .Query}}Query(ctx).Execute(req)
        err = iter.Each(func(row *{{.Service}}_{{camelCase .Query}}Row) error {
            res, err := row.{{.Response}}()
            if err != nil {
                return err
            }
            return stream.Send(res)
        })
        if err != nil {
            return err
        }
    }
    return nil
}
    `
	funcMap := template.FuncMap{
		"camelCase": _gen.CamelCase,
	}
	t := template.Must(template.New("BidirectionalStream").Funcs(funcMap).Parse(biFormat))
	return t.Execute(printerProxy, params)
}
