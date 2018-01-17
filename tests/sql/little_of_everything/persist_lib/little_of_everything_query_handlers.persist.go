package persist_lib

import "golang.org/x/net/context"

type Testservice1MethodReceiver struct {
	Handlers Testservice1QueryHandlers
}
type Testservice1QueryHandlers struct {
	UnaryExample1Handler          func(context.Context, *ExampleTable1ForTestservice1, func(Scanable)) error
	UnaryExample2Handler          func(context.Context, *Test_TestForTestservice1, func(Scanable)) error
	ServerStreamSelectHandler     func(context.Context, *ExampleTable1ForTestservice1, func(Scanable)) error
	ClientStreamingExampleHandler func(context.Context) (func(*ExampleTable1ForTestservice1), func() (Scanable, error))
}

// next must be called on each result row
func (p *Testservice1MethodReceiver) UnaryExample1(ctx context.Context, params *ExampleTable1ForTestservice1, next func(Scanable)) error {
	return p.Handlers.UnaryExample1Handler(ctx, params, next)
}

// next must be called on each result row
func (p *Testservice1MethodReceiver) UnaryExample2(ctx context.Context, params *Test_TestForTestservice1, next func(Scanable)) error {
	return p.Handlers.UnaryExample2Handler(ctx, params, next)
}

// next must be called on each result row
func (p *Testservice1MethodReceiver) ServerStreamSelect(ctx context.Context, params *ExampleTable1ForTestservice1, next func(Scanable)) error {
	return p.Handlers.ServerStreamSelectHandler(ctx, params, next)
}

// given a context, returns two functions.  (feed, stop)
// feed will be called once for every row recieved by the handler
// stop will be called when the client is done streaming. it expects
//a  row to be returned, or nil.
func (p *Testservice1MethodReceiver) ClientStreamingExample(ctx context.Context) (func(*ExampleTable1ForTestservice1), func() (Scanable, error)) {
	return p.Handlers.ClientStreamingExampleHandler(ctx)
}
func DefaultUnaryExample1Handler(accessor SqlClientGetter) func(context.Context, *ExampleTable1ForTestservice1, func(Scanable)) error {
	return func(ctx context.Context, req *ExampleTable1ForTestservice1, next func(Scanable)) error {
		sqlDB, err := accessor()
		if err != nil {
			return err
		}
		res := Testservice1UnaryExample1Query(sqlDB, req)
		err = res.Do(func(row Scanable) error {
			next(row)
			return nil
		})
		if err != nil {
			return err
		}
		return nil
	}
}
func DefaultUnaryExample2Handler(accessor SqlClientGetter) func(context.Context, *Test_TestForTestservice1, func(Scanable)) error {
	return func(ctx context.Context, req *Test_TestForTestservice1, next func(Scanable)) error {
		sqlDB, err := accessor()
		if err != nil {
			return err
		}
		res := Testservice1UnaryExample2Query(sqlDB, req)
		err = res.Do(func(row Scanable) error {
			next(row)
			return nil
		})
		if err != nil {
			return err
		}
		return nil
	}
}
func DefaultServerStreamSelectHandler(accessor SqlClientGetter) func(context.Context, *ExampleTable1ForTestservice1, func(Scanable)) error {
	return func(ctx context.Context, req *ExampleTable1ForTestservice1, next func(Scanable)) error {
		sqlDB, err := accessor()
		if err != nil {
			return err
		}
		tx, err := sqlDB.Begin()
		if err != nil {
			return err
		}
		res := Testservice1ServerStreamSelectQuery(tx, req)
		err = res.Do(func(row Scanable) error {
			next(row)
			return nil
		})
		if err != nil {
			return err
		}
		if err := tx.Commit(); err != nil {
			return err
		}
		return res.Err()
	}
}
func DefaultClientStreamingExampleHandler(accessor SqlClientGetter) func(context.Context) (func(*ExampleTable1ForTestservice1), func() (Scanable, error)) {
	return func(ctx context.Context) (func(*ExampleTable1ForTestservice1), func() (Scanable, error)) {
		var feedErr error
		sqlDb, err := accessor()
		if err != nil {
			feedErr = err
		}
		tx, err := sqlDb.Begin()
		if err != nil {
			feedErr = err
		}
		feed := func(req *ExampleTable1ForTestservice1) {
			if feedErr != nil {
				return
			}
			if res := Testservice1ClientStreamingExampleQuery(tx, req); res.Err() != nil {
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
