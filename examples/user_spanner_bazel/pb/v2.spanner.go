// This file is generated by protoc-gen-persist
// Source File: pb/user.proto
// DO NOT EDIT !

// TXs, Queries, Hooks, TypeMappings, Handlers, Rows, Iters
package pb

import (
	"fmt"
  "time"
	io "io"

  spanner "cloud.google.com/go/spanner"
  "google.golang.org/api/iterator"

	proto "github.com/golang/protobuf/proto"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
	context "golang.org/x/net/context"
	codes "google.golang.org/grpc/codes"
	gstatus "google.golang.org/grpc/status"
)

// WriteHandlers
// RestOf<S>Handlers
type RestOfUServHandlers interface {
	UpdateAllNames(req *Empty, stream UServ_UpdateAllNamesServer) error
}

// WriteTypeMappigns
type UServTypeMappings interface {
	TimestampTimestamp() TimestampTimestampMappingImpl
	SliceStringParam() SliceStringParamMappingImpl
}

type TimestampTimestampMappingImpl interface {
	ToProto(**timestamp.Timestamp) error
	Empty() TimestampTimestampMappingImpl
	ToSpanner(*timestamp.Timestamp) TimestampTimestampMappingImpl
	SpannerScan(src *spanner.GenericColumnValue) error
	SpannerValue() (interface{}, error)
}

type SliceStringParamMappingImpl interface {
	ToProto(**SliceStringParam) error
	Empty() SliceStringParamMappingImpl
  ToSpanner(*SliceStringParam) SliceStringParamMappingImpl
	SpannerScan(src *spanner.GenericColumnValue) error
	SpannerValue() (interface{}, error)
}

type UServHooks interface {
	InsertUsersBeforeHook(*User) (*Empty, error)
	InsertUsersAfterHook(*User, *Empty) error
	GetAllUsersBeforeHook(*Empty) ([]*User, error)
	GetAllUsersAfterHook(*Empty, *User) error
}
type alwaysScanner struct {
	i *interface{}
}

func (s *alwaysScanner) Scan(src interface{}) error {
	s.i = &src
	return nil
}

type Result interface {
  LastInsertId() (int64, error)
  RowsAffected() (int64, error)
}
type SpannerResult struct {
  iter *spanner.RowIterator
}

func (sr *SpannerResult) LastInsertId() (int64, error) {
  // sr.iter.QueryStats or sr.iter.QueryPlan
  return -1, nil
}
func (sr *SpannerResult) RowsAffected() (int64, error) {
  // Execution statistics for the query. Available after RowIterator.Next returns iterator.Done
  return sr.iter.RowCount, nil
}

type Runable interface {
	// QueryContext(context.Context, string, ...interface{}) (*spanner.RowIterator, error)
	// ExecContext(context.Context, string, ...interface{}) (SpannerResult, error)
  ReadWriteTransaction(context.Context, func(context.Context, *spanner.ReadWriteTransaction) error)(time.Time, error)
}

// func DefaultClientStreamingPersistTx(ctx context.Context, db *spanner.Client) (PersistTx, error) {
//   // TODO dont need if spanner is taking care of it
// 	return db.BeginTx(ctx, nil)
// }
// func DefaultServerStreamingPersistTx(ctx context.Context, db *spanner.Client) (PersistTx, error) {
// 	return NopPersistTx(db)
// }
// func DefaultBidiStreamingPersistTx(ctx context.Context, db *spanner.Client) (PersistTx, error) {
// 	return NopPersistTx(db)
// }
// func DefaultUnaryPersistTx(ctx context.Context, db *spanner.Client) (PersistTx, error) {
// 	return NopPersistTx(db)
// }

// type ignoreTx struct {
// 	r Runable
// }

// func (this *ignoreTx) Commit() error   { return nil }
// func (this *ignoreTx) Rollback() error { return nil }
// func (this *ignoreTx) QueryContext(ctx context.Context, x string, ys ...interface{}) (*spanner.RowIterator, error) {
// 	return this.r.QueryContext(ctx, x, ys...)
// }
// func (this *ignoreTx) ExecContext(ctx context.Context, x string, ys ...interface{}) (SpannerResult, error) {
// 	return this.r.ExecContext(ctx, x, ys...)
// }

type UServ_QueryOpts struct {
	MAPPINGS UServTypeMappings
	db       Runable
	ctx      context.Context
}

func DefaultUServQueryOpts(db Runable) UServ_QueryOpts {
	return UServ_QueryOpts{
		db: db,
	}
}

type UServ_Queries struct {
	opts UServ_QueryOpts
}
//type PersistTx interface {
//	Commit() error
//	Rollback() error
//	Runable
//}

//func (tx *PersistTx) Commit() error {
//  //TODO 
//}
//func (tx *PersistTx) Rollback() error {
//  //TODO 
//}

// func NopPersistTx(r Runable) (PersistTx, error) {
// 	return &ignoreTx{r}, nil
// }

type UServ_InsertUsersOut interface {
	GetId() int64
	GetName() string
	GetFriends() *Friends
	GetCreatedOn() *timestamp.Timestamp
}

type UServ_InsertUsersRow struct {
	item UServ_InsertUsersOut
	err  error
}

func newUServ_InsertUsersRow(item UServ_InsertUsersOut, err error) *UServ_InsertUsersRow {
	return &UServ_InsertUsersRow{item, err}
}

// Unwrap takes an address to a proto.Message as its only parameter
// Unwrap can only set into output protos of that match method return types + the out option on the query itself
func (this *UServ_InsertUsersRow) Unwrap(pointerToMsg proto.Message) error {
	if this.err != nil {
		return this.err
	}
	// for each known method result
	if o, ok := (pointerToMsg).(*User); ok {
		if o == nil {
			return fmt.Errorf("must initialize *User before giving to Unwrap()")
		}
		res, _ := this.User()
		// set shared fields
		o.Id = res.Id
		o.Name = res.Name
		o.Friends = res.Friends
		o.CreatedOn = res.CreatedOn
		return nil
	}
	if o, ok := (pointerToMsg).(*Friends); ok {
		if o == nil {
			return fmt.Errorf("must initialize *Friends before giving to Unwrap()")
		}
	}

	return nil
}

// one for each Output type of the methods that use this query + the output proto itself

func (this *UServ_InsertUsersRow) User() (*User, error) {
	return &User{
		Id:        this.item.GetId(),
		Name:      this.item.GetName(),
		Friends:   this.item.GetFriends(),
		CreatedOn: this.item.GetCreatedOn(),
	}, nil
}

// just for example
func (this *UServ_InsertUsersRow) Friends() (*Friends, error) {
	return nil, nil
}

// UServPersistQueries returns all the known 'SQL' queires for the 'UServ' service.
func UServPersistQueries(db Runable, opts ...UServ_QueryOpts) *UServ_Queries {
	var myOpts UServ_QueryOpts
	if len(opts) > 0 {
		myOpts = opts[0]
	} else {
		myOpts = DefaultUServQueryOpts(db)
	}
	return &UServ_Queries{
		opts: myOpts,
	}
}

// camel case the services query name
// method for every query

// InsertUsersQuery returns a new struct wrapping the current UServ_QueryOpts
// that will perform 'UServ' services 'insert_users_query' on the database
// when executed
func (this *UServ_Queries) InsertUsersQuery(ctx context.Context) *UServ_InsertUsersQuery {
	return &UServ_InsertUsersQuery{
		opts: UServ_QueryOpts{
			MAPPINGS: this.opts.MAPPINGS,
			db:       this.opts.db,
			ctx:      ctx,
		},
	}
}

// I dont know this is a insert query, I only know this is a query
type UServ_InsertUsersQuery struct {
	opts UServ_QueryOpts
	ctx  context.Context
}

func (this *UServ_InsertUsersQuery) QueryInTypeUser()  {}
func (this *UServ_InsertUsersQuery) QueryOutTypeUser() {}

// the main execute function
func (this *UServ_InsertUsersQuery) Execute(x UServ_InsertUsersOut) *UServ_InsertUsersIter {
	var setupErr error
	params := []interface{}{
		func() (out interface{}) {
			out = x.GetId()
			return
		}(),
		func() (out interface{}) {
			out = x.GetName()
			return
		}(),
		func() (out interface{}) {
			raw, err := proto.Marshal(x.GetFriends())
			if err != nil {
				setupErr = err
			}
			out = raw
			return
		}(),
		func() (out interface{}) {
			mapper := this.opts.MAPPINGS.TimestampTimestamp()
			out = mapper.ToSpanner(x.GetCreatedOn())
			return
		}(),
	}
	result := &UServ_InsertUsersIter{
		tm:  this.opts.MAPPINGS,
		ctx: this.ctx,
	}
	if setupErr != nil {
		result.err = setupErr
		return result
	}

  _, result.err = this.opts.db.ReadWriteTransaction(this.ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
    stmt := spanner.Statement{
      SQL: fmt.Sprintf("insert into users (id, name, friends, created_on) values (%v, %v, %v, %v);", params...)}
    iter := txn.QueryWithStats(ctx, stmt)
    result.rows = iter
    result.result = SpannerResult{
      iter: iter,
    }
    // TODO TODO TODO defer iter.Stop()
    return nil
  })

	return result
}

//<SERVICE><QUERY (camel)><MESSAGE TYPE>Iter
type UServ_InsertUsersIter struct {
	result SpannerResult
	rows   *spanner.RowIterator
	err    error
	tm     UServTypeMappings
	ctx    context.Context
}

func (this *UServ_InsertUsersIter) IterOutTypeUser() {}
func (this *UServ_InsertUsersIter) IterInTypeUser()  {}

// Each performs 'fun' on each row in the result set.
// Each respects the context passed to it.
// It will stop iteration, and returns ctx.Err() if encountered.
func (this *UServ_InsertUsersIter) Each(fun func(*UServ_InsertUsersRow) error) error {
	for {
		select {
		case <-this.ctx.Done():
			return this.ctx.Err()
		default:
			if row, ok := this.Next(); !ok {
				return nil
			} else if err := fun(row); err != nil {
				return err
			}
		}
	}
	return nil
}

// One returns the sole row, or ensures an error if there was not one result when this row is converted
func (this *UServ_InsertUsersIter) One() *UServ_InsertUsersRow {
	first, hasFirst := this.Next()
	_, hasSecond := this.Next()
	if !hasFirst || hasSecond {
		return newUServ_InsertUsersRow(first.item, fmt.Errorf("expected exactly 1 result from query 'InsertUsers'"))
	}
	return first
}

// Zero returns an error if there were any rows in the result
func (this *UServ_InsertUsersIter) Zero() error {
	if _, ok := this.Next(); ok {
		return fmt.Errorf("expected exactly 0 results from query 'InsertUsers'")
	}
	return nil
}

// Next returns the next scanned row out of the database, or (nil, false) if there are no more rows
func (this *UServ_InsertUsersIter) Next() (*UServ_InsertUsersRow, bool) {
	if this.rows == nil || this.err == io.EOF {
		return nil, false
	} else if this.err != nil {
		err := this.err
		this.err = io.EOF
		return &UServ_InsertUsersRow{err: err}, true
	}
  row, err := this.rows.Next()
  if err == iterator.Done {
    if this.err == nil {
      this.err = io.EOF
      return nil, false
    }
  }
  if err != nil {
		return &UServ_InsertUsersRow{err: err}, true
  }

  var id int64
  if err := row.ColumnByName("id", &id); err != nil {
      return &UServ_InsertUsersRow{err: fmt.Errorf("cant convert db column id to protobuf go type int64")}, true
  }
  var name string
  if err := row.ColumnByName("name", &name); err != nil {
      return &UServ_InsertUsersRow{err: fmt.Errorf("cant convert db column name to protobuf go type string")}, true
  }
  var friends *Friends
  if err := row.ColumnByName("friends", &friends); err != nil {
      return &UServ_InsertUsersRow{err: fmt.Errorf("cant convert db column friends to protobuf go type *Friends")}, true
  }
  var created_on *timestamp.Timestamp
  if err := row.ColumnByName("created_on", &created_on); err != nil {
      return &UServ_InsertUsersRow{err: fmt.Errorf("could not convert mapped db column created_on to type on User.CreatedOn: %v", err)}, true
  }

  return &UServ_InsertUsersRow{item: &User{Id: id, Name: name, Friends: friends, CreatedOn: created_on}}, true
}

// Slice returns all rows found in the iterator as a Slice.
func (this *UServ_InsertUsersIter) Slice() []*UServ_InsertUsersRow {
	var results []*UServ_InsertUsersRow
	for {
		if i, ok := this.Next(); ok {
			results = append(results, i)
		} else {
			break
		}
	}
	return results
}

type UServ_ImplOpts struct {
	MAPPINGS UServTypeMappings
	HOOKS    UServHooks
  HANDLERS RestOfUServHandlers
}

func DefaultUServImplOpts() UServ_ImplOpts {
	return UServ_ImplOpts{}
}

type UServ_Impl struct {
	opts    *UServ_ImplOpts
	QUERIES *UServ_Queries
	DB      *spanner.Client
}

func UServPersistImpl(db *spanner.Client, opts ...UServ_ImplOpts) *UServ_Impl {
	var myOpts UServ_ImplOpts
	if len(opts) > 0 {
		myOpts = opts[0]
	} else {
		myOpts = DefaultUServImplOpts()
	}
	return &UServ_Impl{
		opts:    &myOpts,
		QUERIES: UServPersistQueries(db, UServ_QueryOpts{MAPPINGS: myOpts.MAPPINGS}),
		DB:      db,
	}
}

// THIS is the grpc handler
func (this *UServ_Impl) InsertUsers(stream UServ_InsertUsersServer) error {
	// tx, err := DefaultClientStreamingPersistTx(stream.Context(), this.DB)
	// if err != nil {
	// 	return gstatus.Errorf(codes.Unknown, "error creating persist tx: %v", err)
	// }
	if err := this.InsertUsersTx(stream); err != nil {
		return gstatus.Errorf(codes.Unknown, "error executing 'insert_users' query: %v", err)
	}
	return nil
}

func (this *UServ_Impl) CreateTable(Empty) Empty {
  // tx, err := DefaultClientStreamingPersistTx(stream.Context(), this.DB)
  // if err != nil {
  // 	return gstatus.Errorf(codes.Unknown, "error creating persist tx: %v", err)
  // }
  if err := this.InsertUsersTx(stream); err != nil {
    return gstatus.Errorf(codes.Unknown, "error executing 'insert_users' query: %v", err)
  }
  return nil
}

func (this *UServ_Impl) InsertUsersTx(stream UServ_InsertUsersServer) error {
	query := this.QUERIES.InsertUsersQuery(stream.Context())
	var first *User
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
		// TODO UPDATE HOOKS FOR CTX
		beforeRes, err := this.opts.HOOKS.InsertUsersBeforeHook( /*stream.Context(),*/ req)
		if err != nil {
			return gstatus.Errorf(codes.Unknown, "error in before hook: %v", err)
		} else if beforeRes != nil {
			continue
		}
		result := query.Execute(req)
		/*for {
		      res := new(User)
		      if row, ok := result.Next(); !ok {
		          break
		      } else if err := row.Unwrap(res); err != nil {
		          return err
		      } else {
		          stream.Send(res)
		      }
		  }
		  err := result.Each(stream.Context(), func(row *UServ_InsertUsersRow) error {
		      res, err := row.User()
		      if err != nil {
		          return err
		      }
		      return stream.Send(res)
		  })
		  users := result.Slice()
		  for _, row := range users {
		      user, err := row.User()
		      if err != nil {
		          return err
		      }
		      stream.Send(user)
		  }
		  res, err := result.One().Friends()
		  err := result.Zero()*/
		// TODO allow results to be returned here?
		if err := result.Zero(); err != nil {
			return gstatus.Errorf(codes.InvalidArgument, "client streaming queries must return zero results")
		}
	}
  // TODO might need to handle commits and rollbacks uniquely.
	// if err := tx.Commit(); err != nil {
		// return fmt.Errorf("executed 'insert_users' query without error, but received error on commit: %v", err)
		// if rollbackErr := tx.Rollback(); rollbackErr != nil {
			// return fmt.Errorf("error executing 'insert_users' query :::AND COULD NOT ROLLBACK::: rollback err: %v, query err: %v", rollbackErr, err)
		// }
	// }
	res := &Empty{}
	if err := this.opts.HOOKS.InsertUsersAfterHook( /*stream.Context(),*/ first, res); err != nil {
		return gstatus.Errorf(codes.Unknown, "error in after hook: %v", err)
	}
	if err := stream.SendAndClose(res); err != nil {
		return gstatus.Errorf(codes.Unknown, "error sending back response: %v", err)
	}

	return nil
}