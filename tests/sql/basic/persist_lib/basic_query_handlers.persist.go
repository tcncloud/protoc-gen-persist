package persist_lib

import (
	"fmt"

	"golang.org/x/net/context"
)

type AmazingMethodReceiver struct {
	Handlers AmazingQueryHandlers
}
type AmazingQueryHandlers struct {
	UniarySelectHandler           func(context.Context, *Test_PartialTableForAmazing, func(Scanable)) error
	UniarySelectWithHooksHandler  func(context.Context, *Test_PartialTableForAmazing, func(Scanable)) error
	ServerStreamHandler           func(context.Context, *Test_NameForAmazing, func(Scanable)) error
	ServerStreamWithHooksHandler  func(context.Context, *Test_NameForAmazing, func(Scanable)) error
	BidirectionalHandler          func(context.Context) (func(*Test_ExampleTableForAmazing) (Scanable, error), func() error)
	BidirectionalWithHooksHandler func(context.Context) (func(*Test_ExampleTableForAmazing) (Scanable, error), func() error)
	ClientStreamHandler           func(context.Context) (func(*Test_ExampleTableForAmazing) error, func() (Scanable, error), error)
	ClientStreamWithHookHandler   func(context.Context) (func(*Test_ExampleTableForAmazing) error, func() (Scanable, error), error)
}

// next must be called on each result row
func (p *AmazingMethodReceiver) UniarySelect(ctx context.Context, params *Test_PartialTableForAmazing, next func(Scanable)) error {
	return p.Handlers.UniarySelectHandler(ctx, params, next)
}

// next must be called on each result row
func (p *AmazingMethodReceiver) UniarySelectWithHooks(ctx context.Context, params *Test_PartialTableForAmazing, next func(Scanable)) error {
	return p.Handlers.UniarySelectWithHooksHandler(ctx, params, next)
}

// next must be called on each result row
func (p *AmazingMethodReceiver) ServerStream(ctx context.Context, params *Test_NameForAmazing, next func(Scanable)) error {
	return p.Handlers.ServerStreamHandler(ctx, params, next)
}

// next must be called on each result row
func (p *AmazingMethodReceiver) ServerStreamWithHooks(ctx context.Context, params *Test_NameForAmazing, next func(Scanable)) error {
	return p.Handlers.ServerStreamWithHooksHandler(ctx, params, next)
}

// returns two functions (feed, stop)
// feed needs to be called for every row received. It will run the query
// and return the result + error// stop needs to be called to signal the transaction has finished
func (p *AmazingMethodReceiver) Bidirectional(ctx context.Context) (func(*Test_ExampleTableForAmazing) (Scanable, error), func() error) {
	return p.Handlers.BidirectionalHandler(ctx)
}

// returns two functions (feed, stop)
// feed needs to be called for every row received. It will run the query
// and return the result + error// stop needs to be called to signal the transaction has finished
func (p *AmazingMethodReceiver) BidirectionalWithHooks(ctx context.Context) (func(*Test_ExampleTableForAmazing) (Scanable, error), func() error) {
	return p.Handlers.BidirectionalWithHooksHandler(ctx)
}

// given a context, returns two functions.  (feed, stop)
// feed will be called once for every row recieved by the handler
// stop will be called when the client is done streaming. it expects
//a  row to be returned, or nil.
func (p *AmazingMethodReceiver) ClientStream(ctx context.Context) (func(*Test_ExampleTableForAmazing) error, func() (Scanable, error), error) {
	return p.Handlers.ClientStreamHandler(ctx)
}

// given a context, returns two functions.  (feed, stop)
// feed will be called once for every row recieved by the handler
// stop will be called when the client is done streaming. it expects
//a  row to be returned, or nil.
func (p *AmazingMethodReceiver) ClientStreamWithHook(ctx context.Context) (func(*Test_ExampleTableForAmazing) error, func() (Scanable, error), error) {
	return p.Handlers.ClientStreamWithHookHandler(ctx)
}
func DefaultUniarySelectHandler(accessor SqlClientGetter) func(context.Context, *Test_PartialTableForAmazing, func(Scanable)) error {
	return func(ctx context.Context, req *Test_PartialTableForAmazing, next func(Scanable)) error {
		sqlDB, err := accessor()
		if err != nil {
			return err
		}
		res := AmazingUniarySelectQuery(sqlDB, req)
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
func DefaultUniarySelectWithHooksHandler(accessor SqlClientGetter) func(context.Context, *Test_PartialTableForAmazing, func(Scanable)) error {
	return func(ctx context.Context, req *Test_PartialTableForAmazing, next func(Scanable)) error {
		sqlDB, err := accessor()
		if err != nil {
			return err
		}
		res := AmazingUniarySelectWithHooksQuery(sqlDB, req)
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
func DefaultServerStreamHandler(accessor SqlClientGetter) func(context.Context, *Test_NameForAmazing, func(Scanable)) error {
	return func(ctx context.Context, req *Test_NameForAmazing, next func(Scanable)) error {
		sqlDB, err := accessor()
		if err != nil {
			return err
		}
		tx, err := sqlDB.Begin()
		if err != nil {
			return err
		}
		res := AmazingServerStreamQuery(tx, req)
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
func DefaultServerStreamWithHooksHandler(accessor SqlClientGetter) func(context.Context, *Test_NameForAmazing, func(Scanable)) error {
	return func(ctx context.Context, req *Test_NameForAmazing, next func(Scanable)) error {
		sqlDB, err := accessor()
		if err != nil {
			return err
		}
		tx, err := sqlDB.Begin()
		if err != nil {
			return err
		}
		res := AmazingServerStreamWithHooksQuery(tx, req)
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
func DefaultBidirectionalHandler(accessor SqlClientGetter) func(context.Context) (func(*Test_ExampleTableForAmazing) (Scanable, error), func() error) {
	return func(ctx context.Context) (func(*Test_ExampleTableForAmazing) (Scanable, error), func() error) {
		var feedErr error
		sqlDb, err := accessor()
		if err != nil {
			feedErr = err
		}
		tx, err := sqlDb.Begin()
		if err != nil {
			feedErr = err
		}
		feed := func(req *Test_ExampleTableForAmazing) (Scanable, error) {
			if feedErr != nil {
				return nil, feedErr
			}
			res := AmazingBidirectionalQuery(tx, req)
			return res, nil
		}
		done := func() error {
			if feedErr != nil {
				tx.Rollback()
			} else {
				feedErr = tx.Commit()
			}
			return feedErr
		}
		return feed, done
	}
}
func DefaultBidirectionalWithHooksHandler(accessor SqlClientGetter) func(context.Context) (func(*Test_ExampleTableForAmazing) (Scanable, error), func() error) {
	return func(ctx context.Context) (func(*Test_ExampleTableForAmazing) (Scanable, error), func() error) {
		var feedErr error
		sqlDb, err := accessor()
		if err != nil {
			feedErr = err
		}
		tx, err := sqlDb.Begin()
		if err != nil {
			feedErr = err
		}
		feed := func(req *Test_ExampleTableForAmazing) (Scanable, error) {
			if feedErr != nil {
				return nil, feedErr
			}
			res := AmazingBidirectionalWithHooksQuery(tx, req)
			return res, nil
		}
		done := func() error {
			if feedErr != nil {
				tx.Rollback()
			} else {
				feedErr = tx.Commit()
			}
			return feedErr
		}
		return feed, done
	}
}
func DefaultClientStreamHandler(accessor SqlClientGetter) func(context.Context) (func(*Test_ExampleTableForAmazing) error, func() (Scanable, error), error) {
	return func(ctx context.Context) (func(*Test_ExampleTableForAmazing) error, func() (Scanable, error), error) {
		sqlDb, err := accessor()
		if err != nil {
			return nil, nil, err
		}
		tx, err := sqlDb.Begin()
		if err != nil {
			return nil, nil, err
		}
		feed := func(req *Test_ExampleTableForAmazing) error {
			if res := AmazingClientStreamQuery(tx, req); res.Err() != nil {
				if err := tx.Rollback(); err != nil {
					return fmt.Errorf("%v, %v", err, res.Err())
				}
				return res.Err()
			}
			return nil
		}
		done := func() (Scanable, error) {
			if err := tx.Commit(); err != nil {
				return nil, err
			}
			return nil, nil
		}
		return feed, done, nil
	}
}
func DefaultClientStreamWithHookHandler(accessor SqlClientGetter) func(context.Context) (func(*Test_ExampleTableForAmazing) error, func() (Scanable, error), error) {
	return func(ctx context.Context) (func(*Test_ExampleTableForAmazing) error, func() (Scanable, error), error) {
		sqlDb, err := accessor()
		if err != nil {
			return nil, nil, err
		}
		tx, err := sqlDb.Begin()
		if err != nil {
			return nil, nil, err
		}
		feed := func(req *Test_ExampleTableForAmazing) error {
			if res := AmazingClientStreamWithHookQuery(tx, req); res.Err() != nil {
				if err := tx.Rollback(); err != nil {
					return fmt.Errorf("%v, %v", err, res.Err())
				}
				return res.Err()
			}
			return nil
		}
		done := func() (Scanable, error) {
			if err := tx.Commit(); err != nil {
				return nil, err
			}
			return nil, nil
		}
		return feed, done, nil
	}
}
