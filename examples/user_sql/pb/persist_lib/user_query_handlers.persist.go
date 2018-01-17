package persist_lib

import "golang.org/x/net/context"

type UServMethodReceiver struct {
	Handlers UServQueryHandlers
}
type UServQueryHandlers struct {
	CreateTableHandler     func(context.Context, *EmptyForUServ, func(Scanable)) error
	InsertUsersHandler     func(context.Context) (func(*UserForUServ), func() (Scanable, error))
	GetAllUsersHandler     func(context.Context, *EmptyForUServ, func(Scanable)) error
	SelectUserByIdHandler  func(context.Context, *UserForUServ, func(Scanable)) error
	UpdateUserNamesHandler func(context.Context) (func(*UserForUServ) (Scanable, error), func() error)
	UpdateNameToFooHandler func(context.Context, *UserForUServ, func(Scanable)) error
	GetFriendsHandler      func(context.Context, *FriendsQueryForUServ, func(Scanable)) error
	DropTableHandler       func(context.Context, *EmptyForUServ, func(Scanable)) error
}

// next must be called on each result row
func (p *UServMethodReceiver) CreateTable(ctx context.Context, params *EmptyForUServ, next func(Scanable)) error {
	return p.Handlers.CreateTableHandler(ctx, params, next)
}

// given a context, returns two functions.  (feed, stop)
// feed will be called once for every row recieved by the handler
// stop will be called when the client is done streaming. it expects
//a  row to be returned, or nil.
func (p *UServMethodReceiver) InsertUsers(ctx context.Context) (func(*UserForUServ), func() (Scanable, error)) {
	return p.Handlers.InsertUsersHandler(ctx)
}

// next must be called on each result row
func (p *UServMethodReceiver) GetAllUsers(ctx context.Context, params *EmptyForUServ, next func(Scanable)) error {
	return p.Handlers.GetAllUsersHandler(ctx, params, next)
}

// next must be called on each result row
func (p *UServMethodReceiver) SelectUserById(ctx context.Context, params *UserForUServ, next func(Scanable)) error {
	return p.Handlers.SelectUserByIdHandler(ctx, params, next)
}

// returns two functions (feed, stop)
// feed needs to be called for every row received. It will run the query
// and return the result + error// stop needs to be called to signal the transaction has finished
func (p *UServMethodReceiver) UpdateUserNames(ctx context.Context) (func(*UserForUServ) (Scanable, error), func() error) {
	return p.Handlers.UpdateUserNamesHandler(ctx)
}

// next must be called on each result row
func (p *UServMethodReceiver) UpdateNameToFoo(ctx context.Context, params *UserForUServ, next func(Scanable)) error {
	return p.Handlers.UpdateNameToFooHandler(ctx, params, next)
}

// next must be called on each result row
func (p *UServMethodReceiver) GetFriends(ctx context.Context, params *FriendsQueryForUServ, next func(Scanable)) error {
	return p.Handlers.GetFriendsHandler(ctx, params, next)
}

// next must be called on each result row
func (p *UServMethodReceiver) DropTable(ctx context.Context, params *EmptyForUServ, next func(Scanable)) error {
	return p.Handlers.DropTableHandler(ctx, params, next)
}
func DefaultCreateTableHandler(accessor SqlClientGetter) func(context.Context, *EmptyForUServ, func(Scanable)) error {
	return func(ctx context.Context, req *EmptyForUServ, next func(Scanable)) error {
		sqlDB, err := accessor()
		if err != nil {
			return err
		}
		if res := UServCreateTableQuery(sqlDB, req); res.Err() != nil {
			return err
		}
		return nil
	}
}
func DefaultInsertUsersHandler(accessor SqlClientGetter) func(context.Context) (func(*UserForUServ), func() (Scanable, error)) {
	return func(ctx context.Context) (func(*UserForUServ), func() (Scanable, error)) {
		var feedErr error
		sqlDb, err := accessor()
		if err != nil {
			feedErr = err
		}
		tx, err := sqlDb.Begin()
		if err != nil {
			feedErr = err
		}
		feed := func(req *UserForUServ) {
			if feedErr != nil {
				return
			}
			if res := UServInsertUsersQuery(tx, req); res.Err() != nil {
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
func DefaultGetAllUsersHandler(accessor SqlClientGetter) func(context.Context, *EmptyForUServ, func(Scanable)) error {
	return func(ctx context.Context, req *EmptyForUServ, next func(Scanable)) error {
		sqlDB, err := accessor()
		if err != nil {
			return err
		}
		tx, err := sqlDB.Begin()
		if err != nil {
			return err
		}
		res := UServGetAllUsersQuery(tx, req)
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
func DefaultSelectUserByIdHandler(accessor SqlClientGetter) func(context.Context, *UserForUServ, func(Scanable)) error {
	return func(ctx context.Context, req *UserForUServ, next func(Scanable)) error {
		sqlDB, err := accessor()
		if err != nil {
			return err
		}
		res := UServSelectUserByIdQuery(sqlDB, req)
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
func DefaultUpdateUserNamesHandler(accessor SqlClientGetter) func(context.Context) (func(*UserForUServ) (Scanable, error), func() error) {
	return func(ctx context.Context) (func(*UserForUServ) (Scanable, error), func() error) {
		var feedErr error
		sqlDb, err := accessor()
		if err != nil {
			feedErr = err
		}
		tx, err := sqlDb.Begin()
		if err != nil {
			feedErr = err
		}
		feed := func(req *UserForUServ) (Scanable, error) {
			if feedErr != nil {
				return nil, feedErr
			}
			res := UServUpdateUserNamesQuery(tx, req)
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
func DefaultUpdateNameToFooHandler(accessor SqlClientGetter) func(context.Context, *UserForUServ, func(Scanable)) error {
	return func(ctx context.Context, req *UserForUServ, next func(Scanable)) error {
		sqlDB, err := accessor()
		if err != nil {
			return err
		}
		if res := UServUpdateNameToFooQuery(sqlDB, req); res.Err() != nil {
			return err
		}
		return nil
	}
}
func DefaultGetFriendsHandler(accessor SqlClientGetter) func(context.Context, *FriendsQueryForUServ, func(Scanable)) error {
	return func(ctx context.Context, req *FriendsQueryForUServ, next func(Scanable)) error {
		sqlDB, err := accessor()
		if err != nil {
			return err
		}
		tx, err := sqlDB.Begin()
		if err != nil {
			return err
		}
		res := UServGetFriendsQuery(tx, req)
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
func DefaultDropTableHandler(accessor SqlClientGetter) func(context.Context, *EmptyForUServ, func(Scanable)) error {
	return func(ctx context.Context, req *EmptyForUServ, next func(Scanable)) error {
		sqlDB, err := accessor()
		if err != nil {
			return err
		}
		if res := UServDropTableQuery(sqlDB, req); res.Err() != nil {
			return err
		}
		return nil
	}
}
