package persist_lib

import "golang.org/x/net/context"

type ExampleService1MethodReceiver struct {
	Handlers ExampleService1QueryHandlers
}
type ExampleService1QueryHandlers struct {
	UnaryExample1Handler          func(context.Context, *ExampleTable1ForExampleService1, func(Scanable)) error
	UnaryExample2Handler          func(context.Context, *Test_TestForExampleService1, func(Scanable)) error
	ServerStreamSelectHandler     func(context.Context, *ExampleTable1ForExampleService1, func(Scanable)) error
	ClientStreamingExampleHandler func(context.Context) (func(*ExampleTable1ForExampleService1), func() (Scanable, error))
}

// next must be called on each result row
func (p *ExampleService1MethodReceiver) UnaryExample1(ctx context.Context, params *ExampleTable1ForExampleService1, next func(Scanable)) error {
	return p.Handlers.UnaryExample1Handler(ctx, params, next)
}

// next must be called on each result row
func (p *ExampleService1MethodReceiver) UnaryExample2(ctx context.Context, params *Test_TestForExampleService1, next func(Scanable)) error {
	return p.Handlers.UnaryExample2Handler(ctx, params, next)
}

// next must be called on each result row
func (p *ExampleService1MethodReceiver) ServerStreamSelect(ctx context.Context, params *ExampleTable1ForExampleService1, next func(Scanable)) error {
	return p.Handlers.ServerStreamSelectHandler(ctx, params, next)
}

// given a context, returns two functions.  (feed, stop)
// feed will be called once for every row recieved by the handler
// stop will be called when the client is done streaming. it expects
//a  row to be returned, or nil.
func (p *ExampleService1MethodReceiver) ClientStreamingExample(ctx context.Context) (func(*ExampleTable1ForExampleService1), func() (Scanable, error)) {
	return p.Handlers.ClientStreamingExampleHandler(ctx)
}
func DefaultUnaryExample1Handler(accessor SqlClientGetter) func(context.Context, *ExampleTable1ForExampleService1, func(Scanable)) error {
	return func(ctx context.Context, req *ExampleTable1ForExampleService1, next func(Scanable)) error {
		sqlDB, err := accessor()
		if err != nil {
			return err
		}
		row := ExampleTable1FromUnaryExample1Query(sqlDB, req)
		next(row)
		return nil
	}
}
func DefaultUnaryExample2Handler(accessor SqlClientGetter) func(context.Context, *Test_TestForExampleService1, func(Scanable)) error {
	return func(ctx context.Context, req *Test_TestForExampleService1, next func(Scanable)) error {
		sqlDB, err := accessor()
		if err != nil {
			return err
		}
		row := TestFromUnaryExample2Query(sqlDB, req)
		next(row)
		return nil
	}
}
func DefaultServerStreamSelectHandler(accessor SqlClientGetter) func(context.Context, *ExampleTable1ForExampleService1, func(Scanable)) error {
	return func(ctx context.Context, req *ExampleTable1ForExampleService1, next func(Scanable)) error {
		sqlDB, err := accessor()
		if err != nil {
			return err
		}
		tx, err := sqlDB.Begin()
		if err != nil {
			return err
		}
		rows, err := ExampleTable1FromServerStreamSelectQuery(tx, req)
		if err != nil {
			return err
		}
		defer rows.Close()
		for rows.Next() {
			next(rows)
		}
		if err := tx.Commit(); err != nil {
			return err
		}
		return rows.Err()
	}
}
func DefaultClientStreamingExampleHandler(accessor SqlClientGetter) func(context.Context) (func(*ExampleTable1ForExampleService1), func() (Scanable, error)) {
	return func(ctx context.Context) (func(*ExampleTable1ForExampleService1), func() (Scanable, error)) {
		var feedErr error
		sqlDb, err := accessor()
		if err != nil {
			feedErr = err
		}
		tx, err := sqlDb.Begin()
		if err != nil {
			feedErr = err
		}
		feed := func(req *ExampleTable1ForExampleService1) {
			if feedErr != nil {
				return
			}
			if _, err := ExampleTable1FromClientStreamingExampleQuery(tx, req); err != nil {
				feedErr = err
			}
		}
		done := func() (Scanable, error) {
			if err := tx.Commit(); err != nil {
				return nil, err
			}
			return nil, feedErr
		}
		return feed, done
	}
}

type NotEnabledServiceMethodReceiver struct {
	Handlers NotEnabledServiceQueryHandlers
}
type NotEnabledServiceQueryHandlers struct {
}
