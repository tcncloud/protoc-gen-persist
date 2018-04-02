package persist_lib

import "golang.org/x/net/context"
import "cloud.google.com/go/spanner"

type ExtraSrvMethodReceiver struct {
	Handlers ExtraSrvQueryHandlers
}
type ExtraSrvQueryHandlers struct {
	ExtraUnaryHandler func(context.Context, *Test_NumRowsForExtraSrv, func(*spanner.Row)) error
}

// next must be called on each result row
func (p *ExtraSrvMethodReceiver) ExtraUnary(ctx context.Context, params *Test_NumRowsForExtraSrv, next func(*spanner.Row)) error {
	return p.Handlers.ExtraUnaryHandler(ctx, params, next)
}
func DefaultExtraUnaryHandler(accessor SpannerClientGetter) func(context.Context, *Test_NumRowsForExtraSrv, func(*spanner.Row)) error {
	return func(ctx context.Context, req *Test_NumRowsForExtraSrv, next func(*spanner.Row)) error {
		cli, err := accessor()
		if err != nil {
			return err
		}
		iter := cli.Single().Query(ctx, ExtraSrvExtraUnaryQuery(req))
		if err := iter.Do(func(r *spanner.Row) error {
			next(r)
			return nil
		}); err != nil {
			return err
		}
		return nil
	}
}

type MySpannerMethodReceiver struct {
	Handlers MySpannerQueryHandlers
}
type MySpannerQueryHandlers struct {
	UniaryInsertHandler                func(context.Context, *Test_ExampleTableForMySpanner, func(*spanner.Row)) error
	UniarySelectHandler                func(context.Context, *Test_ExampleTableForMySpanner, func(*spanner.Row)) error
	TestNestHandler                    func(context.Context, *SomethingForMySpanner, func(*spanner.Row)) error
	TestEverythingHandler              func(context.Context, *HasTimestampForMySpanner, func(*spanner.Row)) error
	UniarySelectWithDirectivesHandler  func(context.Context, *Test_ExampleTableForMySpanner, func(*spanner.Row)) error
	UniaryUpdateHandler                func(context.Context, *Test_ExampleTableForMySpanner, func(*spanner.Row)) error
	UniaryDeleteRangeHandler           func(context.Context, *Test_ExampleTableRangeForMySpanner, func(*spanner.Row)) error
	UniaryDeleteSingleHandler          func(context.Context, *Test_ExampleTableForMySpanner, func(*spanner.Row)) error
	NoArgsHandler                      func(context.Context, *Test_ExampleTableForMySpanner, func(*spanner.Row)) error
	ServerStreamHandler                func(context.Context, *Test_NameForMySpanner, func(*spanner.Row)) error
	ClientStreamInsertHandler          func(context.Context) (func(*Test_ExampleTableForMySpanner) error, func() (*spanner.Row, error), error)
	ClientStreamDeleteHandler          func(context.Context) (func(*Test_ExampleTableForMySpanner) error, func() (*spanner.Row, error), error)
	ClientStreamUpdateHandler          func(context.Context) (func(*Test_ExampleTableForMySpanner) error, func() (*spanner.Row, error), error)
	UniaryInsertWithHooksHandler       func(context.Context, *Test_ExampleTableForMySpanner, func(*spanner.Row)) error
	UniarySelectWithHooksHandler       func(context.Context, *Test_ExampleTableForMySpanner, func(*spanner.Row)) error
	UniaryUpdateWithHooksHandler       func(context.Context, *Test_ExampleTableForMySpanner, func(*spanner.Row)) error
	UniaryDeleteWithHooksHandler       func(context.Context, *Test_ExampleTableRangeForMySpanner, func(*spanner.Row)) error
	ServerStreamWithHooksHandler       func(context.Context, *Test_NameForMySpanner, func(*spanner.Row)) error
	ClientStreamUpdateWithHooksHandler func(context.Context) (func(*Test_ExampleTableForMySpanner) error, func() (*spanner.Row, error), error)
}

// next must be called on each result row
func (p *MySpannerMethodReceiver) UniaryInsert(ctx context.Context, params *Test_ExampleTableForMySpanner, next func(*spanner.Row)) error {
	return p.Handlers.UniaryInsertHandler(ctx, params, next)
}

// next must be called on each result row
func (p *MySpannerMethodReceiver) UniarySelect(ctx context.Context, params *Test_ExampleTableForMySpanner, next func(*spanner.Row)) error {
	return p.Handlers.UniarySelectHandler(ctx, params, next)
}

// next must be called on each result row
func (p *MySpannerMethodReceiver) TestNest(ctx context.Context, params *SomethingForMySpanner, next func(*spanner.Row)) error {
	return p.Handlers.TestNestHandler(ctx, params, next)
}

// next must be called on each result row
func (p *MySpannerMethodReceiver) TestEverything(ctx context.Context, params *HasTimestampForMySpanner, next func(*spanner.Row)) error {
	return p.Handlers.TestEverythingHandler(ctx, params, next)
}

// next must be called on each result row
func (p *MySpannerMethodReceiver) UniarySelectWithDirectives(ctx context.Context, params *Test_ExampleTableForMySpanner, next func(*spanner.Row)) error {
	return p.Handlers.UniarySelectWithDirectivesHandler(ctx, params, next)
}

// next must be called on each result row
func (p *MySpannerMethodReceiver) UniaryUpdate(ctx context.Context, params *Test_ExampleTableForMySpanner, next func(*spanner.Row)) error {
	return p.Handlers.UniaryUpdateHandler(ctx, params, next)
}

// next must be called on each result row
func (p *MySpannerMethodReceiver) UniaryDeleteRange(ctx context.Context, params *Test_ExampleTableRangeForMySpanner, next func(*spanner.Row)) error {
	return p.Handlers.UniaryDeleteRangeHandler(ctx, params, next)
}

// next must be called on each result row
func (p *MySpannerMethodReceiver) UniaryDeleteSingle(ctx context.Context, params *Test_ExampleTableForMySpanner, next func(*spanner.Row)) error {
	return p.Handlers.UniaryDeleteSingleHandler(ctx, params, next)
}

// next must be called on each result row
func (p *MySpannerMethodReceiver) NoArgs(ctx context.Context, params *Test_ExampleTableForMySpanner, next func(*spanner.Row)) error {
	return p.Handlers.NoArgsHandler(ctx, params, next)
}

// next must be called on each result row
func (p *MySpannerMethodReceiver) ServerStream(ctx context.Context, params *Test_NameForMySpanner, next func(*spanner.Row)) error {
	return p.Handlers.ServerStreamHandler(ctx, params, next)
}

// given a context, returns two functions.  (feed, stop)
// feed will be called once for every row recieved by the handler
// stop will be called when the client is done streaming. it expects
//a  row to be returned, or nil.
func (p *MySpannerMethodReceiver) ClientStreamInsert(ctx context.Context) (func(*Test_ExampleTableForMySpanner) error, func() (*spanner.Row, error), error) {
	return p.Handlers.ClientStreamInsertHandler(ctx)
}

// given a context, returns two functions.  (feed, stop)
// feed will be called once for every row recieved by the handler
// stop will be called when the client is done streaming. it expects
//a  row to be returned, or nil.
func (p *MySpannerMethodReceiver) ClientStreamDelete(ctx context.Context) (func(*Test_ExampleTableForMySpanner) error, func() (*spanner.Row, error), error) {
	return p.Handlers.ClientStreamDeleteHandler(ctx)
}

// given a context, returns two functions.  (feed, stop)
// feed will be called once for every row recieved by the handler
// stop will be called when the client is done streaming. it expects
//a  row to be returned, or nil.
func (p *MySpannerMethodReceiver) ClientStreamUpdate(ctx context.Context) (func(*Test_ExampleTableForMySpanner) error, func() (*spanner.Row, error), error) {
	return p.Handlers.ClientStreamUpdateHandler(ctx)
}

// next must be called on each result row
func (p *MySpannerMethodReceiver) UniaryInsertWithHooks(ctx context.Context, params *Test_ExampleTableForMySpanner, next func(*spanner.Row)) error {
	return p.Handlers.UniaryInsertWithHooksHandler(ctx, params, next)
}

// next must be called on each result row
func (p *MySpannerMethodReceiver) UniarySelectWithHooks(ctx context.Context, params *Test_ExampleTableForMySpanner, next func(*spanner.Row)) error {
	return p.Handlers.UniarySelectWithHooksHandler(ctx, params, next)
}

// next must be called on each result row
func (p *MySpannerMethodReceiver) UniaryUpdateWithHooks(ctx context.Context, params *Test_ExampleTableForMySpanner, next func(*spanner.Row)) error {
	return p.Handlers.UniaryUpdateWithHooksHandler(ctx, params, next)
}

// next must be called on each result row
func (p *MySpannerMethodReceiver) UniaryDeleteWithHooks(ctx context.Context, params *Test_ExampleTableRangeForMySpanner, next func(*spanner.Row)) error {
	return p.Handlers.UniaryDeleteWithHooksHandler(ctx, params, next)
}

// next must be called on each result row
func (p *MySpannerMethodReceiver) ServerStreamWithHooks(ctx context.Context, params *Test_NameForMySpanner, next func(*spanner.Row)) error {
	return p.Handlers.ServerStreamWithHooksHandler(ctx, params, next)
}

// given a context, returns two functions.  (feed, stop)
// feed will be called once for every row recieved by the handler
// stop will be called when the client is done streaming. it expects
//a  row to be returned, or nil.
func (p *MySpannerMethodReceiver) ClientStreamUpdateWithHooks(ctx context.Context) (func(*Test_ExampleTableForMySpanner) error, func() (*spanner.Row, error), error) {
	return p.Handlers.ClientStreamUpdateWithHooksHandler(ctx)
}
func DefaultUniaryInsertHandler(accessor SpannerClientGetter) func(context.Context, *Test_ExampleTableForMySpanner, func(*spanner.Row)) error {
	return func(ctx context.Context, req *Test_ExampleTableForMySpanner, next func(*spanner.Row)) error {
		cli, err := accessor()
		if err != nil {
			return err
		}
		if _, err := cli.Apply(ctx, []*spanner.Mutation{MySpannerUniaryInsertQuery(req)}); err != nil {
			return err
		}
		next(nil) // this is an apply, it has no result
		return nil
	}
}
func DefaultUniarySelectHandler(accessor SpannerClientGetter) func(context.Context, *Test_ExampleTableForMySpanner, func(*spanner.Row)) error {
	return func(ctx context.Context, req *Test_ExampleTableForMySpanner, next func(*spanner.Row)) error {
		cli, err := accessor()
		if err != nil {
			return err
		}
		iter := cli.Single().Query(ctx, MySpannerUniarySelectQuery(req))
		if err := iter.Do(func(r *spanner.Row) error {
			next(r)
			return nil
		}); err != nil {
			return err
		}
		return nil
	}
}
func DefaultTestNestHandler(accessor SpannerClientGetter) func(context.Context, *SomethingForMySpanner, func(*spanner.Row)) error {
	return func(ctx context.Context, req *SomethingForMySpanner, next func(*spanner.Row)) error {
		cli, err := accessor()
		if err != nil {
			return err
		}
		iter := cli.Single().Query(ctx, MySpannerTestNestQuery(req))
		if err := iter.Do(func(r *spanner.Row) error {
			next(r)
			return nil
		}); err != nil {
			return err
		}
		return nil
	}
}
func DefaultTestEverythingHandler(accessor SpannerClientGetter) func(context.Context, *HasTimestampForMySpanner, func(*spanner.Row)) error {
	return func(ctx context.Context, req *HasTimestampForMySpanner, next func(*spanner.Row)) error {
		cli, err := accessor()
		if err != nil {
			return err
		}
		iter := cli.Single().Query(ctx, MySpannerTestEverythingQuery(req))
		if err := iter.Do(func(r *spanner.Row) error {
			next(r)
			return nil
		}); err != nil {
			return err
		}
		return nil
	}
}
func DefaultUniarySelectWithDirectivesHandler(accessor SpannerClientGetter) func(context.Context, *Test_ExampleTableForMySpanner, func(*spanner.Row)) error {
	return func(ctx context.Context, req *Test_ExampleTableForMySpanner, next func(*spanner.Row)) error {
		cli, err := accessor()
		if err != nil {
			return err
		}
		iter := cli.Single().Query(ctx, MySpannerUniarySelectWithDirectivesQuery(req))
		if err := iter.Do(func(r *spanner.Row) error {
			next(r)
			return nil
		}); err != nil {
			return err
		}
		return nil
	}
}
func DefaultUniaryUpdateHandler(accessor SpannerClientGetter) func(context.Context, *Test_ExampleTableForMySpanner, func(*spanner.Row)) error {
	return func(ctx context.Context, req *Test_ExampleTableForMySpanner, next func(*spanner.Row)) error {
		cli, err := accessor()
		if err != nil {
			return err
		}
		if _, err := cli.Apply(ctx, []*spanner.Mutation{MySpannerUniaryUpdateQuery(req)}); err != nil {
			return err
		}
		next(nil) // this is an apply, it has no result
		return nil
	}
}
func DefaultUniaryDeleteRangeHandler(accessor SpannerClientGetter) func(context.Context, *Test_ExampleTableRangeForMySpanner, func(*spanner.Row)) error {
	return func(ctx context.Context, req *Test_ExampleTableRangeForMySpanner, next func(*spanner.Row)) error {
		cli, err := accessor()
		if err != nil {
			return err
		}
		if _, err := cli.Apply(ctx, []*spanner.Mutation{MySpannerUniaryDeleteRangeQuery(req)}); err != nil {
			return err
		}
		next(nil) // this is an apply, it has no result
		return nil
	}
}
func DefaultUniaryDeleteSingleHandler(accessor SpannerClientGetter) func(context.Context, *Test_ExampleTableForMySpanner, func(*spanner.Row)) error {
	return func(ctx context.Context, req *Test_ExampleTableForMySpanner, next func(*spanner.Row)) error {
		cli, err := accessor()
		if err != nil {
			return err
		}
		if _, err := cli.Apply(ctx, []*spanner.Mutation{MySpannerUniaryDeleteSingleQuery(req)}); err != nil {
			return err
		}
		next(nil) // this is an apply, it has no result
		return nil
	}
}
func DefaultNoArgsHandler(accessor SpannerClientGetter) func(context.Context, *Test_ExampleTableForMySpanner, func(*spanner.Row)) error {
	return func(ctx context.Context, req *Test_ExampleTableForMySpanner, next func(*spanner.Row)) error {
		cli, err := accessor()
		if err != nil {
			return err
		}
		iter := cli.Single().Query(ctx, MySpannerNoArgsQuery(req))
		if err := iter.Do(func(r *spanner.Row) error {
			next(r)
			return nil
		}); err != nil {
			return err
		}
		return nil
	}
}
func DefaultServerStreamHandler(accessor SpannerClientGetter) func(context.Context, *Test_NameForMySpanner, func(*spanner.Row)) error {
	return func(ctx context.Context, req *Test_NameForMySpanner, next func(*spanner.Row)) error {
		cli, err := accessor()
		if err != nil {
			return err
		}
		iter := cli.Single().Query(ctx, MySpannerServerStreamQuery(req))
		if err := iter.Do(func(r *spanner.Row) error {
			next(r)
			return nil
		}); err != nil {
			return err
		}
		return nil
	}
}
func DefaultClientStreamInsertHandler(accessor SpannerClientGetter) func(context.Context) (func(*Test_ExampleTableForMySpanner) error, func() (*spanner.Row, error), error) {
	return func(ctx context.Context) (func(*Test_ExampleTableForMySpanner) error, func() (*spanner.Row, error), error) {
		var muts []*spanner.Mutation
		feed := func(req *Test_ExampleTableForMySpanner) error {
			muts = append(muts, MySpannerClientStreamInsertQuery(req))
			return nil
		}
		done := func() (*spanner.Row, error) {
			cli, err := accessor()
			if err != nil {
				return nil, err
			}
			if _, err := cli.Apply(ctx, muts); err != nil {
				return nil, err
			}
			return nil, nil // we dont have a row, because we are an apply
		}
		return feed, done, nil
	}
}
func DefaultClientStreamDeleteHandler(accessor SpannerClientGetter) func(context.Context) (func(*Test_ExampleTableForMySpanner) error, func() (*spanner.Row, error), error) {
	return func(ctx context.Context) (func(*Test_ExampleTableForMySpanner) error, func() (*spanner.Row, error), error) {
		var muts []*spanner.Mutation
		feed := func(req *Test_ExampleTableForMySpanner) error {
			muts = append(muts, MySpannerClientStreamDeleteQuery(req))
			return nil
		}
		done := func() (*spanner.Row, error) {
			cli, err := accessor()
			if err != nil {
				return nil, err
			}
			if _, err := cli.Apply(ctx, muts); err != nil {
				return nil, err
			}
			return nil, nil // we dont have a row, because we are an apply
		}
		return feed, done, nil
	}
}
func DefaultClientStreamUpdateHandler(accessor SpannerClientGetter) func(context.Context) (func(*Test_ExampleTableForMySpanner) error, func() (*spanner.Row, error), error) {
	return func(ctx context.Context) (func(*Test_ExampleTableForMySpanner) error, func() (*spanner.Row, error), error) {
		var muts []*spanner.Mutation
		feed := func(req *Test_ExampleTableForMySpanner) error {
			muts = append(muts, MySpannerClientStreamUpdateQuery(req))
			return nil
		}
		done := func() (*spanner.Row, error) {
			cli, err := accessor()
			if err != nil {
				return nil, err
			}
			if _, err := cli.Apply(ctx, muts); err != nil {
				return nil, err
			}
			return nil, nil // we dont have a row, because we are an apply
		}
		return feed, done, nil
	}
}
func DefaultUniaryInsertWithHooksHandler(accessor SpannerClientGetter) func(context.Context, *Test_ExampleTableForMySpanner, func(*spanner.Row)) error {
	return func(ctx context.Context, req *Test_ExampleTableForMySpanner, next func(*spanner.Row)) error {
		cli, err := accessor()
		if err != nil {
			return err
		}
		if _, err := cli.Apply(ctx, []*spanner.Mutation{MySpannerUniaryInsertWithHooksQuery(req)}); err != nil {
			return err
		}
		next(nil) // this is an apply, it has no result
		return nil
	}
}
func DefaultUniarySelectWithHooksHandler(accessor SpannerClientGetter) func(context.Context, *Test_ExampleTableForMySpanner, func(*spanner.Row)) error {
	return func(ctx context.Context, req *Test_ExampleTableForMySpanner, next func(*spanner.Row)) error {
		cli, err := accessor()
		if err != nil {
			return err
		}
		iter := cli.Single().Query(ctx, MySpannerUniarySelectWithHooksQuery(req))
		if err := iter.Do(func(r *spanner.Row) error {
			next(r)
			return nil
		}); err != nil {
			return err
		}
		return nil
	}
}
func DefaultUniaryUpdateWithHooksHandler(accessor SpannerClientGetter) func(context.Context, *Test_ExampleTableForMySpanner, func(*spanner.Row)) error {
	return func(ctx context.Context, req *Test_ExampleTableForMySpanner, next func(*spanner.Row)) error {
		cli, err := accessor()
		if err != nil {
			return err
		}
		if _, err := cli.Apply(ctx, []*spanner.Mutation{MySpannerUniaryUpdateWithHooksQuery(req)}); err != nil {
			return err
		}
		next(nil) // this is an apply, it has no result
		return nil
	}
}
func DefaultUniaryDeleteWithHooksHandler(accessor SpannerClientGetter) func(context.Context, *Test_ExampleTableRangeForMySpanner, func(*spanner.Row)) error {
	return func(ctx context.Context, req *Test_ExampleTableRangeForMySpanner, next func(*spanner.Row)) error {
		cli, err := accessor()
		if err != nil {
			return err
		}
		if _, err := cli.Apply(ctx, []*spanner.Mutation{MySpannerUniaryDeleteWithHooksQuery(req)}); err != nil {
			return err
		}
		next(nil) // this is an apply, it has no result
		return nil
	}
}
func DefaultServerStreamWithHooksHandler(accessor SpannerClientGetter) func(context.Context, *Test_NameForMySpanner, func(*spanner.Row)) error {
	return func(ctx context.Context, req *Test_NameForMySpanner, next func(*spanner.Row)) error {
		cli, err := accessor()
		if err != nil {
			return err
		}
		iter := cli.Single().Query(ctx, MySpannerServerStreamWithHooksQuery(req))
		if err := iter.Do(func(r *spanner.Row) error {
			next(r)
			return nil
		}); err != nil {
			return err
		}
		return nil
	}
}
func DefaultClientStreamUpdateWithHooksHandler(accessor SpannerClientGetter) func(context.Context) (func(*Test_ExampleTableForMySpanner) error, func() (*spanner.Row, error), error) {
	return func(ctx context.Context) (func(*Test_ExampleTableForMySpanner) error, func() (*spanner.Row, error), error) {
		var muts []*spanner.Mutation
		feed := func(req *Test_ExampleTableForMySpanner) error {
			muts = append(muts, MySpannerClientStreamUpdateWithHooksQuery(req))
			return nil
		}
		done := func() (*spanner.Row, error) {
			cli, err := accessor()
			if err != nil {
				return nil, err
			}
			if _, err := cli.Apply(ctx, muts); err != nil {
				return nil, err
			}
			return nil, nil // we dont have a row, because we are an apply
		}
		return feed, done, nil
	}
}
