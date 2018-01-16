package persist_lib

import "golang.org/x/net/context"
import "cloud.google.com/go/spanner"

type BobsMethodReceiver struct {
	Handlers BobsQueryHandlers
}
type BobsQueryHandlers struct {
	DeleteBobsHandler         func(context.Context, *BobForBobs, func(*spanner.Row)) error
	PutBobsHandler            func(context.Context) (func(*BobForBobs), func() (*spanner.Row, error))
	GetBobsHandler            func(context.Context, *EmptyForBobs, func(*spanner.Row)) error
	GetPeopleFromNamesHandler func(context.Context, *NamesForBobs, func(*spanner.Row)) error
}

// next must be called on each result row
func (p *BobsMethodReceiver) DeleteBobs(ctx context.Context, params *BobForBobs, next func(*spanner.Row)) error {
	return p.Handlers.DeleteBobsHandler(ctx, params, next)
}

// given a context, returns two functions.  (feed, stop)
// feed will be called once for every row recieved by the handler
// stop will be called when the client is done streaming. it expects
//a  row to be returned, or nil.
func (p *BobsMethodReceiver) PutBobs(ctx context.Context) (func(*BobForBobs), func() (*spanner.Row, error)) {
	return p.Handlers.PutBobsHandler(ctx)
}

// next must be called on each result row
func (p *BobsMethodReceiver) GetBobs(ctx context.Context, params *EmptyForBobs, next func(*spanner.Row)) error {
	return p.Handlers.GetBobsHandler(ctx, params, next)
}

// next must be called on each result row
func (p *BobsMethodReceiver) GetPeopleFromNames(ctx context.Context, params *NamesForBobs, next func(*spanner.Row)) error {
	return p.Handlers.GetPeopleFromNamesHandler(ctx, params, next)
}
func DefaultDeleteBobsHandler(accessor SpannerClientGetter) func(context.Context, *BobForBobs, func(*spanner.Row)) error {
	return func(ctx context.Context, req *BobForBobs, next func(*spanner.Row)) error {
		cli, err := accessor()
		if err != nil {
			return err
		}
		if _, err := cli.Apply(ctx, []*spanner.Mutation{BobsDeleteBobsQuery(req)}); err != nil {
			return err
		}
		next(nil) // this is an apply, it has no result
		return nil
	}
}
func DefaultPutBobsHandler(accessor SpannerClientGetter) func(context.Context) (func(*BobForBobs), func() (*spanner.Row, error)) {
	return func(ctx context.Context) (func(*BobForBobs), func() (*spanner.Row, error)) {
		var muts []*spanner.Mutation
		feed := func(req *BobForBobs) {
			muts = append(muts, BobsPutBobsQuery(req))
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
		return feed, done
	}
}
func DefaultGetBobsHandler(accessor SpannerClientGetter) func(context.Context, *EmptyForBobs, func(*spanner.Row)) error {
	return func(ctx context.Context, req *EmptyForBobs, next func(*spanner.Row)) error {
		cli, err := accessor()
		if err != nil {
			return err
		}
		iter := cli.Single().Query(ctx, BobsGetBobsQuery(req))
		if err := iter.Do(func(r *spanner.Row) error {
			next(r)
			return nil
		}); err != nil {
			return err
		}
		return nil
	}
}
func DefaultGetPeopleFromNamesHandler(accessor SpannerClientGetter) func(context.Context, *NamesForBobs, func(*spanner.Row)) error {
	return func(ctx context.Context, req *NamesForBobs, next func(*spanner.Row)) error {
		cli, err := accessor()
		if err != nil {
			return err
		}
		iter := cli.Single().Query(ctx, BobsGetPeopleFromNamesQuery(req))
		if err := iter.Do(func(r *spanner.Row) error {
			next(r)
			return nil
		}); err != nil {
			return err
		}
		return nil
	}
}
