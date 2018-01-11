package persist_lib

import "golang.org/x/net/context"
import "cloud.google.com/go/spanner"

type UServMethodReceiver struct {
	Handlers UServQueryHandlers
}
type UServQueryHandlers struct {
	InsertUsersHandler     func(context.Context) (func(*UserForUServ), func() (*spanner.Row, error))
	GetAllUsersHandler     func(context.Context, *EmptyForUServ, func(*spanner.Row)) error
	SelectUserByIdHandler  func(context.Context, *UserForUServ, func(*spanner.Row)) error
	UpdateUserNamesHandler func(context.Context) (func(*UserForUServ), func() (*spanner.Row, error))
	GetFriendsHandler      func(context.Context, *FriendsForUServ, func(*spanner.Row)) error
}

// given a context, returns two functions.  (feed, stop)
// feed will be called once for every row recieved by the handler
// stop will be called when the client is done streaming. it expects
//a  row to be returned, or nil.
func (p *UServMethodReceiver) InsertUsers(ctx context.Context) (func(*UserForUServ), func() (*spanner.Row, error)) {
	return p.Handlers.InsertUsersHandler(ctx)
}

// next must be called on each result row
func (p *UServMethodReceiver) GetAllUsers(ctx context.Context, params *EmptyForUServ, next func(*spanner.Row)) error {
	return p.Handlers.GetAllUsersHandler(ctx, params, next)
}

// next must be called on each result row
func (p *UServMethodReceiver) SelectUserById(ctx context.Context, params *UserForUServ, next func(*spanner.Row)) error {
	return p.Handlers.SelectUserByIdHandler(ctx, params, next)
}

// given a context, returns two functions.  (feed, stop)
// feed will be called once for every row recieved by the handler
// stop will be called when the client is done streaming. it expects
//a  row to be returned, or nil.
func (p *UServMethodReceiver) UpdateUserNames(ctx context.Context) (func(*UserForUServ), func() (*spanner.Row, error)) {
	return p.Handlers.UpdateUserNamesHandler(ctx)
}

// next must be called on each result row
func (p *UServMethodReceiver) GetFriends(ctx context.Context, params *FriendsForUServ, next func(*spanner.Row)) error {
	return p.Handlers.GetFriendsHandler(ctx, params, next)
}
func DefaultInsertUsersHandler(accessor SpannerClientGetter) func(context.Context) (func(*UserForUServ), func() (*spanner.Row, error)) {
	return func(ctx context.Context) (func(*UserForUServ), func() (*spanner.Row, error)) {
		var muts []*spanner.Mutation
		feed := func(req *UserForUServ) {
			muts = append(muts, UserFromInsertUsersQuery(req))
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
func DefaultGetAllUsersHandler(accessor SpannerClientGetter) func(context.Context, *EmptyForUServ, func(*spanner.Row)) error {
	return func(ctx context.Context, req *EmptyForUServ, next func(*spanner.Row)) error {
		cli, err := accessor()
		if err != nil {
			return err
		}
		iter := cli.Single().Query(ctx, EmptyFromGetAllUsersQuery(req))
		if err := iter.Do(func(r *spanner.Row) error {
			next(r)
			return nil
		}); err != nil {
			return err
		}
		return nil
	}
}
func DefaultSelectUserByIdHandler(accessor SpannerClientGetter) func(context.Context, *UserForUServ, func(*spanner.Row)) error {
	return func(ctx context.Context, req *UserForUServ, next func(*spanner.Row)) error {
		cli, err := accessor()
		if err != nil {
			return err
		}
		iter := cli.Single().Query(ctx, UserFromSelectUserByIdQuery(req))
		if err := iter.Do(func(r *spanner.Row) error {
			next(r)
			return nil
		}); err != nil {
			return err
		}
		return nil
	}
}
func DefaultUpdateUserNamesHandler(accessor SpannerClientGetter) func(context.Context) (func(*UserForUServ), func() (*spanner.Row, error)) {
	return func(ctx context.Context) (func(*UserForUServ), func() (*spanner.Row, error)) {
		var muts []*spanner.Mutation
		feed := func(req *UserForUServ) {
			muts = append(muts, UserFromUpdateUserNamesQuery(req))
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
func DefaultGetFriendsHandler(accessor SpannerClientGetter) func(context.Context, *FriendsForUServ, func(*spanner.Row)) error {
	return func(ctx context.Context, req *FriendsForUServ, next func(*spanner.Row)) error {
		cli, err := accessor()
		if err != nil {
			return err
		}
		iter := cli.Single().Query(ctx, FriendsFromGetFriendsQuery(req))
		if err := iter.Do(func(r *spanner.Row) error {
			next(r)
			return nil
		}); err != nil {
			return err
		}
		return nil
	}
}
