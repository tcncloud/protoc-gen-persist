package generator

import (
	"fmt"
)

type PersistStringer struct{}

func (per *PersistStringer) MessageInputDeclaration(method *Method) string {
	printer := &Printer{}
	printer.P("type %s struct{\n", NewPLInputName(method))

	getPersistLibTypeName := GetSqlPersistLibTypeName
	if method.IsSpanner() {
		getPersistLibTypeName = GetSpannerPersistLibTypeName
	}

	inputTypeDescs := method.GetTypeDescArrayForStruct(method.GetInputTypeStruct())
	for _, qf := range inputTypeDescs {
		typeName := getPersistLibTypeName(qf)
		printer.P("%s %s\n", qf.Name, typeName)
	}
	printer.P("}\n")
	printer.P("// this could be used in a query, so generate the getters/setters\n")
	for _, qf := range inputTypeDescs {
		typeName := getPersistLibTypeName(qf)
		plInputName := NewPLInputName(method)
		printer.P(
			"func(p *%s) Get%s() %s{ return p.%s }\n",
			plInputName, qf.Name, typeName, qf.Name,
		)
		printer.P(
			"func(p *%s) Set%s(param %s) { p.%s = param }\n",
			plInputName, qf.Name, typeName, qf.Name,
		)
	}

	return printer.String()
}

// merges custom defined handlers with our own
func (per *PersistStringer) PersistImplBuilder(service *Service) string {
	var dbType string
	var backend string
	if service.IsSpanner() {
		dbType = "spanner.Client"
		backend = "Spanner"
	} else {
		dbType = "sql.DB"
		backend = "Sql"
	}
	printer := &Printer{}
	printer.P(
		"type %sImpl struct{\nPERSIST *persist_lib.%s\nFORWARDED RestOf%sHandlers\n}\n",
		service.GetName(),
		NewPersistHelperName(service),
		service.GetName(),
	)
	printer.P("type RestOf%sHandlers interface{\n", service.GetName())
	for _, m := range *service.Methods {
		if m.GetMethodOption() == nil {
			if m.IsUnary() {
				printer.P(
					"%s(ctx context.Context, req *%s) (*%s, error)\n",
					m.GetName(),
					m.GetInputType(),
					m.GetOutputType(),
				)
			} else if m.IsServerStreaming() {
				printer.P(
					"%s(req *%s, stream %s) error\n",
					m.GetName(),
					m.GetInputType(),
					NewStreamType(m),
				)
			} else {
				printer.P(
					"%s(stream %s) error\n",
					m.GetName(),
					NewStreamType(m),
				)
			}
		}
	}
	printer.P("}\n")
	printer.PA([]string{
		"type %sImplBuilder struct {\n",
		"err error\n ",
		"rest RestOf%sHandlers\n",
		"queryHandlers *persist_lib.%sQueryHandlers\n",
		"i *%sImpl\n",
		"db *%s\n}\n",
		"func New%sBuilder() *%sImplBuilder {\nreturn &%sImplBuilder{i: &%sImpl{}}\n}\n",
	},
		service.GetName(),
		service.GetName(),
		service.GetName(),
		service.GetName(), dbType,
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
		"accessor := persist_lib.New%sClientGetter(b.db)\n",
		"queryHandlers := &persist_lib.%sQueryHandlers{\n",
	},
		service.GetName(), service.GetName(),
		backend,
		service.GetName(),
	)
	for _, m := range *service.Methods {
		if m.GetMethodOption() == nil {
			continue
		}
		printer.P(
			"%s: persist_lib.Default%s(accessor),\n",
			NewPersistHandlerName(m),
			NewPersistHandlerName(m),
		)
	}
	printer.P("}\n b.queryHandlers = queryHandlers\n return b\n}\n")

	printer.PA([]string{
		"func (b *%sImplBuilder) With%sClient(c *%s) *%sImplBuilder {\n",
		"b.db = c\n return b\n}\n",
	},
		service.GetName(), backend, dbType, service.GetName(),
	)

	if service.IsSpanner() {
		printer.PA([]string{
			"func (b *%sImplBuilder) WithSpannerURI(ctx context.Context, uri string) *%sImplBuilder {\n",
			"cli, err := spanner.NewClient(ctx, uri)\n b.err = err\n b.db = cli\n return b\n}\n",
		},
			service.GetName(), service.GetName(),
		)
	} else {
		printer.PA([]string{
			"func (b *%sImplBuilder) WithNewSqlDb(driverName, dataSourceName string) *%sImplBuilder {\n",
			"db, err := sql.Open(driverName, dataSourceName)\n",
			"b.err = err\n b.db = db\n return b\n}\n",
		},
			service.GetName(), service.GetName(),
		)
	}

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
		var rowType string
		if method.IsSpanner() {
			rowType = "*spanner.Row"
		} else {
			rowType = "Scanable"
		}
		if method.IsClientStreaming() {
			printer.P(
				"%s func(context.Context)(func(*%s), func() (%s, error))\n",
				NewPersistHandlerName(method),
				NewPLInputName(method),
				rowType,
			)
		} else if method.IsBidiStreaming() {
			printer.P(
				"%s func(context.Context) (func(*%s) (%s, error), func() error)\n",
				NewPersistHandlerName(method),
				NewPLInputName(method),
				rowType,
			)
		} else {
			printer.P(
				"%s func(context.Context, *%s, func(%s)) error\n",
				NewPersistHandlerName(method),
				NewPLInputName(method),
				rowType,
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
		var rowType string
		if method.IsSpanner() {
			rowType = "*spanner.Row"
		} else {
			rowType = "Scanable"
		}

		if method.IsClientStreaming() {
			printer.PA([]string{
				"// given a context, returns two functions.  (feed, stop)\n",
				"// feed will be called once for every row recieved by the handler\n",
				"// stop will be called when the client is done streaming. it expects\n",
				"//a  row to be returned, or nil.\n",
				"func (p *%s) %s(ctx context.Context)(func(*%s), func() (%s, error)) {\n",
				"return p.Handlers.%s(ctx)\n}\n",
			},
				NewPersistHelperName(service),
				method.GetName(),
				NewPLInputName(method),
				rowType,
				NewPersistHandlerName(method),
			)
		} else if method.IsBidiStreaming() {
			printer.PA([]string{
				"// returns two functions (feed, stop)\n",
				"// feed needs to be called for every row received. It will run the query\n",
				"// and return the result + error",
				"// stop needs to be called to signal the transaction has finished\n",
				"func (p *%s) %s(ctx context.Context)(func(*%s) (%s, error), func() error) {\n",
				"return p.Handlers.%s(ctx)\n}\n",
			},
				NewPersistHelperName(service),
				method.GetName(),
				NewPLInputName(method),
				rowType,
				NewPersistHandlerName(method),
			)
		} else {
			printer.PA([]string{
				"// next must be called on each result row\n",
				"func(p *%s) %s(ctx context.Context, params *%s, next func(%s)) error {\n",
				"return p.Handlers.%s(ctx, params, next)\n}\n",
			},
				NewPersistHelperName(service),
				method.GetName(),
				NewPLInputName(method),
				rowType,
				NewPersistHandlerName(method),
			)
		}
	}
	return printer.String()
}

func (per *PersistStringer) QueryInterfaceDefinition(method *Method) string {
	if method.GetMethodOption() == nil {
		return ""
	}
	printer := &Printer{}
	printer.P(
		"type %sParams interface{\n",
		NewPLQueryMethodName(method),
	)
	getPersistLibTypeName := GetSpannerPersistLibTypeName
	if method.IsSQL() {
		getPersistLibTypeName = GetSqlPersistLibTypeName
	}

	for _, t := range method.GetTypeDescForQueryFields() {
		interfaceType := getPersistLibTypeName(t)
		printer.P("Get%s() %s\n", t.Name, interfaceType)
	}
	printer.P("}\n")
	return printer.String()
}

func (per *PersistStringer) SqlQueryFunction(method *Method) string {
	opts := method.GetMethodOption()
	if opts == nil {
		return ""
	}
	query := func() (out string) {
		for _, q := range opts.GetQuery() {
			out += q + " "
		}
		return
	}()

	args := opts.GetArguments()
	tds := method.GetTypeDescForFieldsInStructSnakeCase(method.GetInputTypeStruct())

	var argParams string
	for _, a := range args {
		argParams += fmt.Sprintf("req.Get%s(),\n", tds[a].Name)
	}
	// if we are an empty result, then perform an exec, not a query
	lenOfResult := len(method.GetTypeDescArrayForStruct(method.GetOutputTypeStruct()))

	printer := &Printer{}
	queryMethodName := NewPLQueryMethodName(method)
	printer.P("func %s(tx Runable, req %sParams)", queryMethodName, queryMethodName)
	if lenOfResult == 0 || method.IsClientStreaming() { // use an exec
		printer.PA([]string{
			"(sql.Result, error) {\n",
			"return tx.Exec(\n\"%s\",\n%s)",
		},
			query, argParams,
		)
	} else if method.IsServerStreaming() {
		printer.PA([]string{
			"(*sql.Rows, error) {\n",
			"return tx.Query(\n\"%s\",\n%s)",
		},
			query, argParams,
		)
	} else {
		printer.PA([]string{
			"*sql.Row {\n",
			"return tx.QueryRow(\n\"%s\",\n%s)",
		},
			query, argParams,
		)
	}
	printer.P("\n}\n")

	return printer.String()
}
func (per *PersistStringer) SpannerQueryFunction(method *Method) string {
	// we do not have a persist query
	if method.GetMethodOption() == nil {
		return ""
	}
	printer := &Printer{}
	if method.IsSelect() {
		printer.P(
			"func %s(req %sParams) spanner.Statement {\nreturn %s\n}\n",
			NewPLQueryMethodName(method),
			NewPLQueryMethodName(method),
			method.Query,
		)
	} else {
		printer.P(
			"func %s(req %sParams) *spanner.Mutation {\nreturn %s\n}\n",
			NewPLQueryMethodName(method),
			NewPLQueryMethodName(method),
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
		if method.IsSQL() {
			printer.P("%s", per.DefaultSqlFunctionsImpl(method))
		} else if method.IsSpanner() {
			printer.P("%s", per.DefaultSpannerFunctionsImpl(method))
		}
	}
	return printer.String()
}
func (per *PersistStringer) DefaultSqlFunctionsImpl(method *Method) string {
	printer := &Printer{}
	lenOfOutFields := len(method.GetTypeDescArrayForStruct(method.GetOutputTypeStruct()))
	if method.IsClientStreaming() { // use exec
		printer.PA([]string{
			"func Default%sHandler(accessor SqlClientGetter) func(context.Context) ",
			"(func(*%s), func() (Scanable, error)) {\n",
			"return func(ctx context.Context) (func(*%s), func() (Scanable, error)) {\n",
			"var feedErr error\n",
			"sqlDb, err := accessor()\n",
			"if err != nil {\n feedErr = err\n}\n",
			"tx, err := sqlDb.Begin()\n",
			"if err != nil {\n feedErr = err\n}\n",
			"feed := func(req *%s) {\n",
			"if feedErr != nil {\n return \n}\n",
			"if _, err := %s(tx, req); err != nil {\n feedErr = err\n}\n}\n",
			"done := func() (Scanable, error) {\n if err := tx.Commit();err != nil {\n",
			"return nil, err\n}\n return nil, feedErr\n}\n",
			"return feed, done\n}\n}\n",
		},
			method.GetName(),
			NewPLInputName(method),
			NewPLInputName(method),
			NewPLInputName(method),
			NewPLQueryMethodName(method),
		)
	} else if lenOfOutFields == 0 { // use exec
		printer.PA([]string{
			"func Default%sHandler(accessor SqlClientGetter) ",
			"func (context.Context, *%s, func(Scanable)) error {\n",
			"return func(ctx context.Context, req *%s, next func(Scanable)) error {\n",
			"sqlDB, err := accessor()\n if err != nil {\n return err \n}\n",
			"if _, err := %s(sqlDB, req); err != nil {\n return err \n}\n",
			"return nil\n}\n}\n",
		},
			method.GetName(),
			NewPLInputName(method),
			NewPLInputName(method),
			NewPLInputName(method),
			NewPLQueryMethodName(method),
		)
	} else if method.IsServerStreaming() { // use query
		printer.PA([]string{
			"func Default%sHandler(accessor SqlClientGetter) ",
			"func(context.Context, *%s, func(Scanable)) error {\n",
			"return func(ctx context.Context, req *%s, next func(Scanable)) error {\n",
			"sqlDB, err := accessor()\n if err != nil {\n return err\n}\n",
			"tx, err := sqlDB.Begin()\n",
			"if err != nil {\n return err\n}\n",
			"rows, err := %s(tx, req)\n",
			"if err != nil {\n return err \n}\n",
			"defer rows.Close()\n",
			"for rows.Next() {\n",
			"next(rows)\n",
			"}\n if err := tx.Commit(); err != nil { return err \n}\n",
			"return rows.Err()\n}\n}\n",
		},
			method.GetName(),
			NewPLInputName(method),
			NewPLInputName(method),
			NewPLQueryMethodName(method),
		)
	} else if method.IsUnary() { // use queryRow
		printer.PA([]string{
			"func Default%sHandler(accessor SqlClientGetter)  ",
			"func(context.Context, *%s, func(Scanable)) error {\n",
			"return func(ctx context.Context, req *%s, next func(Scanable)) error {\n",
			"sqlDB, err := accessor()\n if err != nil {\n return err\n}\n",
			"row := %s(sqlDB, req)\n",
			"next(row)\nreturn nil}\n}\n",
		},
			method.GetName(),
			NewPLInputName(method),
			NewPLInputName(method),
			NewPLQueryMethodName(method),
		)
	} else if method.IsBidiStreaming() {
		printer.PA([]string{
			"func Default%sHandler(accessor SqlClientGetter) ",
			"func(context.Context) (func(*%s) (Scanable, error), func() error) {\n",
			"return func(ctx context.Context) (func(*%s) (Scanable, error), func() error) {\n",
			"var feedErr error\n",
			"sqlDb, err := accessor()\n",
			"if err != nil {\n feedErr = err\n}\n",
			"tx, err := sqlDb.Begin()\n",
			"if err != nil {\n feedErr = err\n}\n",
			"feed := func(req *%s) (Scanable, error) {\n",
			"if feedErr != nil{\n return nil, feedErr\n}\n row := %s(tx, req)\n",
			"return row, nil\n}\n",
			"done := func() error {\n if feedErr != nil {\n tx.Rollback()\n} else {\n feedErr = tx.Commit()\n}\n",
			"return feedErr\n}\n",
			"return feed,done\n}\n}\n",
		},
			method.GetName(),
			NewPLInputName(method),
			NewPLInputName(method),
			NewPLInputName(method),
			NewPLQueryMethodName(method),
		)
	}
	return printer.String()
}
func (per *PersistStringer) DefaultSpannerFunctionsImpl(method *Method) string {
	printer := &Printer{}
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

// package level definitions for sql implemented libraries
// SqlClientGetter
// Scanable interface
// Runable interface
func (per *PersistStringer) DeclareSqlPackageDefs() string {
	printer := &Printer{}
	printer.P("type SqlClientGetter func() (*sql.DB, error)\n")
	printer.PA([]string{
		"func NewSqlClientGetter(cli *sql.DB) SqlClientGetter {\n",
		"return func() (*sql.DB, error) {\n return cli, nil \n}\n}\n",
	})
	printer.P("type Scanable interface{\nScan(dest ...interface{}) error\n}\n")
	printer.PA([]string{
		"type Runable interface{\n",
		"Query(string, ...interface{}) (*sql.Rows, error)\n",
		"QueryRow(string, ...interface{}) *sql.Row\n",
		"Exec(string, ...interface{}) (sql.Result, error)\n}\n",
	})
	return printer.String()
}
func GetSqlPersistLibTypeName(t TypeDesc) string {
	if t.IsMapped {
		return "interface{}"
	} else if t.IsMessage {
		return "[]byte"
	} else {
		return t.GoName
	}
}
func GetSpannerPersistLibTypeName(t TypeDesc) string {
	if t.IsMapped {
		return "interface{}"
	} else if t.IsMessage && t.IsRepeated {
		return "[][]byte"
	} else if t.IsMessage {
		return "[]byte"
	} else {
		return t.GoName
	}
}
