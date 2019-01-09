// This file is generated by protoc-gen-persist
// Source File: tests/sql/basic/basic.proto
// DO NOT EDIT !
package basic

import (
	sql "database/sql"
	driver "database/sql/driver"
	fmt "fmt"
	io "io"

	proto "github.com/golang/protobuf/proto"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
	test "github.com/tcncloud/protoc-gen-persist/tests/test"
	context "golang.org/x/net/context"
	codes "google.golang.org/grpc/codes"
	gstatus "google.golang.org/grpc/status"
)

type PersistTx interface {
	Commit() error
	Rollback() error
	Runnable
}

func NopPersistTx(r Runnable) (PersistTx, error) {
	return &ignoreTx{r}, nil
}

type ignoreTx struct {
	r Runnable
}

func (this *ignoreTx) Commit() error   { return nil }
func (this *ignoreTx) Rollback() error { return nil }
func (this *ignoreTx) QueryContext(ctx context.Context, x string, ys ...interface{}) (*sql.Rows, error) {
	return this.r.QueryContext(ctx, x, ys...)
}
func (this *ignoreTx) ExecContext(ctx context.Context, x string, ys ...interface{}) (sql.Result, error) {
	return this.r.ExecContext(ctx, x, ys...)
}

type Runnable interface {
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
}

func DefaultClientStreamingPersistTx(ctx context.Context, db *sql.DB) (PersistTx, error) {
	return db.BeginTx(ctx, nil)
}
func DefaultServerStreamingPersistTx(ctx context.Context, db *sql.DB) (PersistTx, error) {
	return NopPersistTx(db)
}
func DefaultBidiStreamingPersistTx(ctx context.Context, db *sql.DB) (PersistTx, error) {
	return NopPersistTx(db)
}
func DefaultUnaryPersistTx(ctx context.Context, db *sql.DB) (PersistTx, error) {
	return NopPersistTx(db)
}

type alwaysScanner struct {
	i *interface{}
}

func (s *alwaysScanner) Scan(src interface{}) error {
	s.i = &src
	return nil
}

type scanable interface {
	Scan(...interface{}) error
	Columns() ([]string, error)
}

// Amazing_Queries holds all the queries found the proto service option as methods
type Amazing_Queries struct {
	opts Amazing_Opts
}

// AmazingPersistQueries returns all the known 'SQL' queires for the 'Amazing' service.
func AmazingPersistQueries(opts ...Amazing_Opts) *Amazing_Queries {
	var myOpts Amazing_Opts
	if len(opts) > 0 {
		myOpts = opts[0]
	} else {
		myOpts = AmazingOpts(nil, nil)
	}
	return &Amazing_Queries{
		opts: myOpts,
	}
}

// SelectByIdQuery returns a new struct wrapping the current Amazing_Opts
// that will perform 'Amazing' services 'select_by_id' on the database
// when executed
func (this *Amazing_Queries) SelectById(ctx context.Context, db Runnable) *Amazing_SelectByIdQuery {
	return &Amazing_SelectByIdQuery{
		opts: this.opts,
		ctx:  ctx,
		db:   db,
	}
}

type Amazing_SelectByIdQuery struct {
	opts Amazing_Opts
	db   Runnable
	ctx  context.Context
}

func (this *Amazing_SelectByIdQuery) QueryInTypeUser()  {}
func (this *Amazing_SelectByIdQuery) QueryOutTypeUser() {}

// Executes the query with parameters retrieved from x
func (this *Amazing_SelectByIdQuery) Execute(x Amazing_SelectByIdIn) *Amazing_SelectByIdIter {
	var setupErr error
	params := []interface{}{
		func() (out interface{}) {
			out = x.GetId()
			return
		}(),
		func() (out interface{}) {
			mapper := this.opts.MAPPINGS.TimestampTimestamp()
			out = mapper.ToSql(x.GetStartTime())
			return
		}(),
	}
	result := &Amazing_SelectByIdIter{
		tm:  this.opts.MAPPINGS,
		ctx: this.ctx,
	}
	if setupErr != nil {
		result.err = setupErr
		return result
	}
	result.rows, result.err = this.db.QueryContext(this.ctx, "SELECT * from example_table Where id=$1 AND start_time>$2", params...)
	return result
}

// SelectByNameQuery returns a new struct wrapping the current Amazing_Opts
// that will perform 'Amazing' services 'select_by_name' on the database
// when executed
func (this *Amazing_Queries) SelectByName(ctx context.Context, db Runnable) *Amazing_SelectByNameQuery {
	return &Amazing_SelectByNameQuery{
		opts: this.opts,
		ctx:  ctx,
		db:   db,
	}
}

type Amazing_SelectByNameQuery struct {
	opts Amazing_Opts
	db   Runnable
	ctx  context.Context
}

func (this *Amazing_SelectByNameQuery) QueryInTypeUser()  {}
func (this *Amazing_SelectByNameQuery) QueryOutTypeUser() {}

// Executes the query with parameters retrieved from x
func (this *Amazing_SelectByNameQuery) Execute(x Amazing_SelectByNameIn) *Amazing_SelectByNameIter {
	var setupErr error
	params := []interface{}{
		func() (out interface{}) {
			out = x.GetName()
			return
		}(),
	}
	result := &Amazing_SelectByNameIter{
		tm:  this.opts.MAPPINGS,
		ctx: this.ctx,
	}
	if setupErr != nil {
		result.err = setupErr
		return result
	}
	result.rows, result.err = this.db.QueryContext(this.ctx, "SELECT * FROM example_table WHERE name=$1", params...)
	return result
}

// InsertQuery returns a new struct wrapping the current Amazing_Opts
// that will perform 'Amazing' services 'insert' on the database
// when executed
func (this *Amazing_Queries) Insert(ctx context.Context, db Runnable) *Amazing_InsertQuery {
	return &Amazing_InsertQuery{
		opts: this.opts,
		ctx:  ctx,
		db:   db,
	}
}

type Amazing_InsertQuery struct {
	opts Amazing_Opts
	db   Runnable
	ctx  context.Context
}

func (this *Amazing_InsertQuery) QueryInTypeUser()  {}
func (this *Amazing_InsertQuery) QueryOutTypeUser() {}

// Executes the query with parameters retrieved from x
func (this *Amazing_InsertQuery) Execute(x Amazing_InsertIn) *Amazing_InsertIter {
	var setupErr error
	params := []interface{}{
		func() (out interface{}) {
			out = x.GetId()
			return
		}(),
		func() (out interface{}) {
			mapper := this.opts.MAPPINGS.TimestampTimestamp()
			out = mapper.ToSql(x.GetStartTime())
			return
		}(),
		func() (out interface{}) {
			out = x.GetName()
			return
		}(),
	}
	result := &Amazing_InsertIter{
		tm:  this.opts.MAPPINGS,
		ctx: this.ctx,
	}
	if setupErr != nil {
		result.err = setupErr
		return result
	}
	result.result, result.err = this.db.ExecContext(this.ctx, "INSERT INTO example_table (id, start_time, name) VALUES ($1, $2, $3)", params...)
	return result
}

type Amazing_SelectByIdIter struct {
	result sql.Result
	rows   *sql.Rows
	err    error
	tm     Amazing_TypeMappings
	ctx    context.Context
}

func (this *Amazing_SelectByIdIter) IterOutTypeTestExampleTable() {}
func (this *Amazing_SelectByIdIter) IterInTypeTestPartialTable()  {}

// Each performs 'fun' on each row in the result set.
// Each respects the context passed to it.
// It will stop iteration, and returns this.ctx.Err() if encountered.
func (this *Amazing_SelectByIdIter) Each(fun func(*Amazing_SelectByIdRow) error) error {
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
}

// One returns the sole row, or ensures an error if there was not one result when this row is converted
func (this *Amazing_SelectByIdIter) One() *Amazing_SelectByIdRow {
	first, hasFirst := this.Next()
	if first != nil && first.err != nil {
		return &Amazing_SelectByIdRow{err: first.err}
	}
	_, hasSecond := this.Next()
	if !hasFirst || hasSecond {
		amount := "none"
		if hasSecond {
			amount = "multiple"
		}
		return &Amazing_SelectByIdRow{err: fmt.Errorf("expected exactly 1 result from query 'SelectById' found %s", amount)}
	}
	return first
}

// Zero returns an error if there were any rows in the result
func (this *Amazing_SelectByIdIter) Zero() error {
	row, ok := this.Next()
	if row != nil && row.err != nil {
		return row.err
	}
	if ok {
		return fmt.Errorf("expected exactly 0 results from query 'SelectById'")
	}
	return nil
}

// Next returns the next scanned row out of the database, or (nil, false) if there are no more rows
func (this *Amazing_SelectByIdIter) Next() (*Amazing_SelectByIdRow, bool) {
	if this.rows == nil || this.err == io.EOF {
		return nil, false
	} else if this.err != nil {
		err := this.err
		this.err = io.EOF
		return &Amazing_SelectByIdRow{err: err}, true
	}
	cols, err := this.rows.Columns()
	if err != nil {
		return &Amazing_SelectByIdRow{err: err}, true
	}
	if !this.rows.Next() {
		if this.err = this.rows.Err(); this.err == nil {
			this.err = io.EOF
			return nil, false
		}
	}
	toScan := make([]interface{}, len(cols))
	scanned := make([]alwaysScanner, len(cols))
	for i := range scanned {
		toScan[i] = &scanned[i]
	}
	if this.err = this.rows.Scan(toScan...); this.err != nil {
		return &Amazing_SelectByIdRow{err: this.err}, true
	}
	res := &test.ExampleTable{}
	for i, col := range cols {
		_ = i
		switch col {
		case "id":
			r, ok := (*scanned[i].i).(int64)
			if !ok {
				return &Amazing_SelectByIdRow{err: fmt.Errorf("cant convert db column id to protobuf go type ")}, true
			}
			res.Id = r
		case "start_time":
			var converted = this.tm.TimestampTimestamp().Empty()
			if err := converted.Scan(*scanned[i].i); err != nil {
				return &Amazing_SelectByIdRow{err: fmt.Errorf("could not convert mapped db column start_time to type on test.ExampleTable.StartTime: %v", err)}, true
			}
			if err := converted.ToProto(&res.StartTime); err != nil {
				return &Amazing_SelectByIdRow{err: fmt.Errorf("could not convert mapped db column start_timeto type on test.ExampleTable.StartTime: %v", err)}, true
			}
		case "name":
			r, ok := (*scanned[i].i).(string)
			if !ok {
				return &Amazing_SelectByIdRow{err: fmt.Errorf("cant convert db column name to protobuf go type ")}, true
			}
			res.Name = r

		default:
			return &Amazing_SelectByIdRow{err: fmt.Errorf("unsupported column in output: %s", col)}, true
		}
	}
	return &Amazing_SelectByIdRow{item: res}, true
}

// Slice returns all rows found in the iterator as a Slice.
func (this *Amazing_SelectByIdIter) Slice() []*Amazing_SelectByIdRow {
	var results []*Amazing_SelectByIdRow
	for {
		if i, ok := this.Next(); ok {
			results = append(results, i)
		} else {
			break
		}
	}
	return results
}

// returns the known columns for this result
func (r *Amazing_SelectByIdIter) Columns() ([]string, error) {
	if r.err != nil {
		return nil, r.err
	}
	if r.rows != nil {
		return r.rows.Columns()
	}
	return nil, nil
}

type Amazing_SelectByNameIter struct {
	result sql.Result
	rows   *sql.Rows
	err    error
	tm     Amazing_TypeMappings
	ctx    context.Context
}

func (this *Amazing_SelectByNameIter) IterOutTypeTestExampleTable() {}
func (this *Amazing_SelectByNameIter) IterInTypeTestName()          {}

// Each performs 'fun' on each row in the result set.
// Each respects the context passed to it.
// It will stop iteration, and returns this.ctx.Err() if encountered.
func (this *Amazing_SelectByNameIter) Each(fun func(*Amazing_SelectByNameRow) error) error {
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
}

// One returns the sole row, or ensures an error if there was not one result when this row is converted
func (this *Amazing_SelectByNameIter) One() *Amazing_SelectByNameRow {
	first, hasFirst := this.Next()
	if first != nil && first.err != nil {
		return &Amazing_SelectByNameRow{err: first.err}
	}
	_, hasSecond := this.Next()
	if !hasFirst || hasSecond {
		amount := "none"
		if hasSecond {
			amount = "multiple"
		}
		return &Amazing_SelectByNameRow{err: fmt.Errorf("expected exactly 1 result from query 'SelectByName' found %s", amount)}
	}
	return first
}

// Zero returns an error if there were any rows in the result
func (this *Amazing_SelectByNameIter) Zero() error {
	row, ok := this.Next()
	if row != nil && row.err != nil {
		return row.err
	}
	if ok {
		return fmt.Errorf("expected exactly 0 results from query 'SelectByName'")
	}
	return nil
}

// Next returns the next scanned row out of the database, or (nil, false) if there are no more rows
func (this *Amazing_SelectByNameIter) Next() (*Amazing_SelectByNameRow, bool) {
	if this.rows == nil || this.err == io.EOF {
		return nil, false
	} else if this.err != nil {
		err := this.err
		this.err = io.EOF
		return &Amazing_SelectByNameRow{err: err}, true
	}
	cols, err := this.rows.Columns()
	if err != nil {
		return &Amazing_SelectByNameRow{err: err}, true
	}
	if !this.rows.Next() {
		if this.err = this.rows.Err(); this.err == nil {
			this.err = io.EOF
			return nil, false
		}
	}
	toScan := make([]interface{}, len(cols))
	scanned := make([]alwaysScanner, len(cols))
	for i := range scanned {
		toScan[i] = &scanned[i]
	}
	if this.err = this.rows.Scan(toScan...); this.err != nil {
		return &Amazing_SelectByNameRow{err: this.err}, true
	}
	res := &test.ExampleTable{}
	for i, col := range cols {
		_ = i
		switch col {
		case "id":
			r, ok := (*scanned[i].i).(int64)
			if !ok {
				return &Amazing_SelectByNameRow{err: fmt.Errorf("cant convert db column id to protobuf go type ")}, true
			}
			res.Id = r
		case "start_time":
			var converted = this.tm.TimestampTimestamp().Empty()
			if err := converted.Scan(*scanned[i].i); err != nil {
				return &Amazing_SelectByNameRow{err: fmt.Errorf("could not convert mapped db column start_time to type on test.ExampleTable.StartTime: %v", err)}, true
			}
			if err := converted.ToProto(&res.StartTime); err != nil {
				return &Amazing_SelectByNameRow{err: fmt.Errorf("could not convert mapped db column start_timeto type on test.ExampleTable.StartTime: %v", err)}, true
			}
		case "name":
			r, ok := (*scanned[i].i).(string)
			if !ok {
				return &Amazing_SelectByNameRow{err: fmt.Errorf("cant convert db column name to protobuf go type ")}, true
			}
			res.Name = r

		default:
			return &Amazing_SelectByNameRow{err: fmt.Errorf("unsupported column in output: %s", col)}, true
		}
	}
	return &Amazing_SelectByNameRow{item: res}, true
}

// Slice returns all rows found in the iterator as a Slice.
func (this *Amazing_SelectByNameIter) Slice() []*Amazing_SelectByNameRow {
	var results []*Amazing_SelectByNameRow
	for {
		if i, ok := this.Next(); ok {
			results = append(results, i)
		} else {
			break
		}
	}
	return results
}

// returns the known columns for this result
func (r *Amazing_SelectByNameIter) Columns() ([]string, error) {
	if r.err != nil {
		return nil, r.err
	}
	if r.rows != nil {
		return r.rows.Columns()
	}
	return nil, nil
}

type Amazing_InsertIter struct {
	result sql.Result
	rows   *sql.Rows
	err    error
	tm     Amazing_TypeMappings
	ctx    context.Context
}

func (this *Amazing_InsertIter) IterOutTypeEmpty()           {}
func (this *Amazing_InsertIter) IterInTypeTestExampleTable() {}

// Each performs 'fun' on each row in the result set.
// Each respects the context passed to it.
// It will stop iteration, and returns this.ctx.Err() if encountered.
func (this *Amazing_InsertIter) Each(fun func(*Amazing_InsertRow) error) error {
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
}

// One returns the sole row, or ensures an error if there was not one result when this row is converted
func (this *Amazing_InsertIter) One() *Amazing_InsertRow {
	first, hasFirst := this.Next()
	if first != nil && first.err != nil {
		return &Amazing_InsertRow{err: first.err}
	}
	_, hasSecond := this.Next()
	if !hasFirst || hasSecond {
		amount := "none"
		if hasSecond {
			amount = "multiple"
		}
		return &Amazing_InsertRow{err: fmt.Errorf("expected exactly 1 result from query 'Insert' found %s", amount)}
	}
	return first
}

// Zero returns an error if there were any rows in the result
func (this *Amazing_InsertIter) Zero() error {
	row, ok := this.Next()
	if row != nil && row.err != nil {
		return row.err
	}
	if ok {
		return fmt.Errorf("expected exactly 0 results from query 'Insert'")
	}
	return nil
}

// Next returns the next scanned row out of the database, or (nil, false) if there are no more rows
func (this *Amazing_InsertIter) Next() (*Amazing_InsertRow, bool) {
	if this.rows == nil || this.err == io.EOF {
		return nil, false
	} else if this.err != nil {
		err := this.err
		this.err = io.EOF
		return &Amazing_InsertRow{err: err}, true
	}
	cols, err := this.rows.Columns()
	if err != nil {
		return &Amazing_InsertRow{err: err}, true
	}
	if !this.rows.Next() {
		if this.err = this.rows.Err(); this.err == nil {
			this.err = io.EOF
			return nil, false
		}
	}
	toScan := make([]interface{}, len(cols))
	scanned := make([]alwaysScanner, len(cols))
	for i := range scanned {
		toScan[i] = &scanned[i]
	}
	if this.err = this.rows.Scan(toScan...); this.err != nil {
		return &Amazing_InsertRow{err: this.err}, true
	}
	res := &Empty{}
	for i, col := range cols {
		_ = i
		switch col {

		default:
			return &Amazing_InsertRow{err: fmt.Errorf("unsupported column in output: %s", col)}, true
		}
	}
	return &Amazing_InsertRow{item: res}, true
}

// Slice returns all rows found in the iterator as a Slice.
func (this *Amazing_InsertIter) Slice() []*Amazing_InsertRow {
	var results []*Amazing_InsertRow
	for {
		if i, ok := this.Next(); ok {
			results = append(results, i)
		} else {
			break
		}
	}
	return results
}

// returns the known columns for this result
func (r *Amazing_InsertIter) Columns() ([]string, error) {
	if r.err != nil {
		return nil, r.err
	}
	if r.rows != nil {
		return r.rows.Columns()
	}
	return nil, nil
}

type Amazing_SelectByIdIn interface {
	GetId() int64
	GetStartTime() *timestamp.Timestamp
}
type Amazing_SelectByIdOut interface {
	GetId() int64
	GetStartTime() *timestamp.Timestamp
	GetName() string
}
type Amazing_SelectByIdRow struct {
	item Amazing_SelectByIdOut
	err  error
}

func newAmazing_SelectByIdRow(item Amazing_SelectByIdOut, err error) *Amazing_SelectByIdRow {
	return &Amazing_SelectByIdRow{item, err}
}

// Unwrap takes an address to a proto.Message as its only parameter
// Unwrap can only set into output protos of that match method return types + the out option on the query itself
func (this *Amazing_SelectByIdRow) Unwrap(pointerToMsg proto.Message) error {
	if this.err != nil {
		return this.err
	}
	if o, ok := (pointerToMsg).(*test.ExampleTable); ok {
		if o == nil {
			return fmt.Errorf("must initialize *test.ExampleTable before giving to Unwrap()")
		}
		res, _ := this.TestExampleTable()
		_ = res
		o.Id = res.Id
		o.StartTime = res.StartTime
		o.Name = res.Name
		return nil
	}
	if o, ok := (pointerToMsg).(*test.ExampleTable); ok {
		if o == nil {
			return fmt.Errorf("must initialize *test.ExampleTable before giving to Unwrap()")
		}
		res, _ := this.TestExampleTable()
		_ = res
		o.Id = res.Id
		o.StartTime = res.StartTime
		o.Name = res.Name
		return nil
	}

	return nil
}
func (this *Amazing_SelectByIdRow) TestExampleTable() (*test.ExampleTable, error) {
	if this.err != nil {
		return nil, this.err
	}
	return &test.ExampleTable{
		Id:        this.item.GetId(),
		StartTime: this.item.GetStartTime(),
		Name:      this.item.GetName(),
	}, nil
}

func (this *Amazing_SelectByIdRow) Proto() (*test.ExampleTable, error) {
	if this.err != nil {
		return nil, this.err
	}
	return &test.ExampleTable{
		Id:        this.item.GetId(),
		StartTime: this.item.GetStartTime(),
		Name:      this.item.GetName(),
	}, nil
}

type Amazing_SelectByNameIn interface {
	GetName() string
}
type Amazing_SelectByNameOut interface {
	GetId() int64
	GetStartTime() *timestamp.Timestamp
	GetName() string
}
type Amazing_SelectByNameRow struct {
	item Amazing_SelectByNameOut
	err  error
}

func newAmazing_SelectByNameRow(item Amazing_SelectByNameOut, err error) *Amazing_SelectByNameRow {
	return &Amazing_SelectByNameRow{item, err}
}

// Unwrap takes an address to a proto.Message as its only parameter
// Unwrap can only set into output protos of that match method return types + the out option on the query itself
func (this *Amazing_SelectByNameRow) Unwrap(pointerToMsg proto.Message) error {
	if this.err != nil {
		return this.err
	}
	if o, ok := (pointerToMsg).(*test.ExampleTable); ok {
		if o == nil {
			return fmt.Errorf("must initialize *test.ExampleTable before giving to Unwrap()")
		}
		res, _ := this.TestExampleTable()
		_ = res
		o.Id = res.Id
		o.StartTime = res.StartTime
		o.Name = res.Name
		return nil
	}
	if o, ok := (pointerToMsg).(*test.ExampleTable); ok {
		if o == nil {
			return fmt.Errorf("must initialize *test.ExampleTable before giving to Unwrap()")
		}
		res, _ := this.TestExampleTable()
		_ = res
		o.Id = res.Id
		o.StartTime = res.StartTime
		o.Name = res.Name
		return nil
	}

	return nil
}
func (this *Amazing_SelectByNameRow) TestExampleTable() (*test.ExampleTable, error) {
	if this.err != nil {
		return nil, this.err
	}
	return &test.ExampleTable{
		Id:        this.item.GetId(),
		StartTime: this.item.GetStartTime(),
		Name:      this.item.GetName(),
	}, nil
}

func (this *Amazing_SelectByNameRow) Proto() (*test.ExampleTable, error) {
	if this.err != nil {
		return nil, this.err
	}
	return &test.ExampleTable{
		Id:        this.item.GetId(),
		StartTime: this.item.GetStartTime(),
		Name:      this.item.GetName(),
	}, nil
}

type Amazing_InsertIn interface {
	GetId() int64
	GetStartTime() *timestamp.Timestamp
	GetName() string
}
type Amazing_InsertOut interface {
}
type Amazing_InsertRow struct {
	item Amazing_InsertOut
	err  error
}

func newAmazing_InsertRow(item Amazing_InsertOut, err error) *Amazing_InsertRow {
	return &Amazing_InsertRow{item, err}
}

// Unwrap takes an address to a proto.Message as its only parameter
// Unwrap can only set into output protos of that match method return types + the out option on the query itself
func (this *Amazing_InsertRow) Unwrap(pointerToMsg proto.Message) error {
	if this.err != nil {
		return this.err
	}
	if o, ok := (pointerToMsg).(*test.NumRows); ok {
		if o == nil {
			return fmt.Errorf("must initialize *test.NumRows before giving to Unwrap()")
		}
		res, _ := this.TestNumRows()
		_ = res

		return nil
	}
	if o, ok := (pointerToMsg).(*test.Ids); ok {
		if o == nil {
			return fmt.Errorf("must initialize *test.Ids before giving to Unwrap()")
		}
		res, _ := this.TestIds()
		_ = res

		return nil
	}

	return nil
}
func (this *Amazing_InsertRow) TestNumRows() (*test.NumRows, error) {
	if this.err != nil {
		return nil, this.err
	}
	return &test.NumRows{}, nil
}
func (this *Amazing_InsertRow) TestIds() (*test.Ids, error) {
	if this.err != nil {
		return nil, this.err
	}
	return &test.Ids{}, nil
}

func (this *Amazing_InsertRow) Proto() (*Empty, error) {
	if this.err != nil {
		return nil, this.err
	}
	return &Empty{}, nil
}

type Amazing_Hooks interface {
	UniarySelectWithHooksBeforeHook(context.Context, *test.PartialTable) (*test.ExampleTable, error)
	ServerStreamWithHooksBeforeHook(context.Context, *test.Name) (*test.ExampleTable, error)
	ClientStreamWithHookBeforeHook(context.Context, *test.ExampleTable) (*test.Ids, error)
	UniarySelectWithHooksAfterHook(context.Context, *test.PartialTable, *test.ExampleTable) error
	ServerStreamWithHooksAfterHook(context.Context, *test.Name, *test.ExampleTable) error
	ClientStreamWithHookAfterHook(context.Context, *test.ExampleTable, *test.Ids) error
}
type Amazing_DefaultHooks struct{}

func (*Amazing_DefaultHooks) UniarySelectWithHooksBeforeHook(context.Context, *test.PartialTable) (*test.ExampleTable, error) {
	return nil, nil
}
func (*Amazing_DefaultHooks) ServerStreamWithHooksBeforeHook(context.Context, *test.Name) (*test.ExampleTable, error) {
	return nil, nil
}
func (*Amazing_DefaultHooks) ClientStreamWithHookBeforeHook(context.Context, *test.ExampleTable) (*test.Ids, error) {
	return nil, nil
}
func (*Amazing_DefaultHooks) UniarySelectWithHooksAfterHook(context.Context, *test.PartialTable, *test.ExampleTable) error {
	return nil
}
func (*Amazing_DefaultHooks) ServerStreamWithHooksAfterHook(context.Context, *test.Name, *test.ExampleTable) error {
	return nil
}
func (*Amazing_DefaultHooks) ClientStreamWithHookAfterHook(context.Context, *test.ExampleTable, *test.Ids) error {
	return nil
}

type Amazing_TypeMappings interface {
	TimestampTimestamp() AmazingTimestampTimestampMappingImpl
}
type Amazing_DefaultTypeMappings struct{}

func (this *Amazing_DefaultTypeMappings) TimestampTimestamp() AmazingTimestampTimestampMappingImpl {
	return &Amazing_DefaultTimestampTimestampMappingImpl{}
}

type Amazing_DefaultTimestampTimestampMappingImpl struct{}

func (this *Amazing_DefaultTimestampTimestampMappingImpl) ToProto(**timestamp.Timestamp) error {
	return nil
}
func (this *Amazing_DefaultTimestampTimestampMappingImpl) Empty() AmazingTimestampTimestampMappingImpl {
	return this
}
func (this *Amazing_DefaultTimestampTimestampMappingImpl) ToSql(*timestamp.Timestamp) sql.Scanner {
	return this
}
func (this *Amazing_DefaultTimestampTimestampMappingImpl) Scan(interface{}) error {
	return nil
}
func (this *Amazing_DefaultTimestampTimestampMappingImpl) Value() (driver.Value, error) {
	return "DEFAULT_TYPE_MAPPING_VALUE", nil
}

type AmazingTimestampTimestampMappingImpl interface {
	ToProto(**timestamp.Timestamp) error
	Empty() AmazingTimestampTimestampMappingImpl
	ToSql(*timestamp.Timestamp) sql.Scanner
	sql.Scanner
	driver.Valuer
}
type Amazing_Opts struct {
	MAPPINGS Amazing_TypeMappings
	HOOKS    Amazing_Hooks
}

func AmazingOpts(hooks Amazing_Hooks, mappings Amazing_TypeMappings) Amazing_Opts {
	opts := Amazing_Opts{
		HOOKS:    &Amazing_DefaultHooks{},
		MAPPINGS: &Amazing_DefaultTypeMappings{},
	}
	if hooks != nil {
		opts.HOOKS = hooks
	}
	if mappings != nil {
		opts.MAPPINGS = mappings
	}
	return opts
}

type Amazing_Impl struct {
	opts     *Amazing_Opts
	QUERIES  *Amazing_Queries
	HANDLERS RestOfAmazingHandlers
	DB       *sql.DB
}

func AmazingPersistImpl(db *sql.DB, handlers RestOfAmazingHandlers, opts ...Amazing_Opts) *Amazing_Impl {
	var myOpts Amazing_Opts
	if len(opts) > 0 {
		myOpts = opts[0]
	} else {
		myOpts = AmazingOpts(nil, nil)
	}
	return &Amazing_Impl{
		opts:     &myOpts,
		QUERIES:  AmazingPersistQueries(myOpts),
		DB:       db,
		HANDLERS: handlers,
	}
}

type RestOfAmazingHandlers interface {
	UnImplementedPersistMethod(context.Context, *test.ExampleTable) (*test.ExampleTable, error)
	NoGenerationForBadReturnTypes(context.Context, *test.ExampleTable) (*BadReturn, error)
}

func (this *Amazing_Impl) UnImplementedPersistMethod(ctx context.Context, req *test.ExampleTable) (*test.ExampleTable, error) {
	return this.HANDLERS.UnImplementedPersistMethod(ctx, req)
}

func (this *Amazing_Impl) NoGenerationForBadReturnTypes(ctx context.Context, req *test.ExampleTable) (*BadReturn, error) {
	return this.HANDLERS.NoGenerationForBadReturnTypes(ctx, req)
}

func (this *Amazing_Impl) UniarySelect(ctx context.Context, req *test.PartialTable) (*test.ExampleTable, error) {
	query := this.QUERIES.SelectById(ctx, this.DB)

	result := query.Execute(req)
	res, err := result.One().TestExampleTable()
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (this *Amazing_Impl) UniarySelectWithHooks(ctx context.Context, req *test.PartialTable) (*test.ExampleTable, error) {
	query := this.QUERIES.SelectById(ctx, this.DB)

	beforeRes, err := this.opts.HOOKS.UniarySelectWithHooksBeforeHook(ctx, req)
	if err != nil {
		return nil, gstatus.Errorf(codes.Unknown, "error in before hook: %v", err)
	} else if beforeRes != nil {
		return beforeRes, nil
	}

	result := query.Execute(req)
	res, err := result.One().TestExampleTable()
	if err != nil {
		return nil, err
	}

	if err := this.opts.HOOKS.UniarySelectWithHooksAfterHook(ctx, req, res); err != nil {
		return nil, gstatus.Errorf(codes.Unknown, "error in after hook: %v", err)
	}

	return res, nil
}

func (this *Amazing_Impl) ServerStream(req *test.Name, stream Amazing_ServerStreamServer) error {
	tx, err := DefaultServerStreamingPersistTx(stream.Context(), this.DB)
	if err != nil {
		return gstatus.Errorf(codes.Unknown, "error creating persist tx: %v", err)
	}
	if err := this.ServerStreamTx(req, stream, tx); err != nil {
		return gstatus.Errorf(codes.Unknown, "error executing 'select_by_name' query: %v", err)
	}
	return nil
}
func (this *Amazing_Impl) ServerStreamTx(req *test.Name, stream Amazing_ServerStreamServer, tx PersistTx) error {
	ctx := stream.Context()
	query := this.QUERIES.SelectByName(ctx, tx)
	iter := query.Execute(req)
	return iter.Each(func(row *Amazing_SelectByNameRow) error {
		res, err := row.TestExampleTable()
		if err != nil {
			return err
		}
		return stream.Send(res)
	})
}

func (this *Amazing_Impl) ServerStreamWithHooks(req *test.Name, stream Amazing_ServerStreamWithHooksServer) error {
	tx, err := DefaultServerStreamingPersistTx(stream.Context(), this.DB)
	if err != nil {
		return gstatus.Errorf(codes.Unknown, "error creating persist tx: %v", err)
	}
	if err := this.ServerStreamWithHooksTx(req, stream, tx); err != nil {
		return gstatus.Errorf(codes.Unknown, "error executing 'select_by_name' query: %v", err)
	}
	return nil
}
func (this *Amazing_Impl) ServerStreamWithHooksTx(req *test.Name, stream Amazing_ServerStreamWithHooksServer, tx PersistTx) error {
	ctx := stream.Context()
	query := this.QUERIES.SelectByName(ctx, tx)
	iter := query.Execute(req)
	return iter.Each(func(row *Amazing_SelectByNameRow) error {
		res, err := row.TestExampleTable()
		if err != nil {
			return err
		}
		return stream.Send(res)
	})
}

func (this *Amazing_Impl) ClientStream(stream Amazing_ClientStreamServer) error {
	tx, err := DefaultClientStreamingPersistTx(stream.Context(), this.DB)
	if err != nil {
		return gstatus.Errorf(codes.Unknown, "error creating persist tx: %v", err)
	}
	if err := this.ClientStreamTx(stream, tx); err != nil {
		return gstatus.Errorf(codes.Unknown, "error executing 'insert' query: %v", err)
	}
	return nil
}
func (this *Amazing_Impl) ClientStreamTx(stream Amazing_ClientStreamServer, tx PersistTx) error {
	query := this.QUERIES.Insert(stream.Context(), tx)
	var first *test.ExampleTable
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

		result := query.Execute(req)
		if err := result.Zero(); err != nil {
			return err
		}
	}
	if err := tx.Commit(); err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return fmt.Errorf("error executing 'insert' query :::AND COULD NOT ROLLBACK::: rollback err: %v, query err: %v", rollbackErr, err)
		}
	}
	res := &test.NumRows{}

	if err := stream.SendAndClose(res); err != nil {
		return gstatus.Errorf(codes.Unknown, "error sending back response: %v", err)
	}
	return nil
}

func (this *Amazing_Impl) ClientStreamWithHook(stream Amazing_ClientStreamWithHookServer) error {
	tx, err := DefaultClientStreamingPersistTx(stream.Context(), this.DB)
	if err != nil {
		return gstatus.Errorf(codes.Unknown, "error creating persist tx: %v", err)
	}
	if err := this.ClientStreamWithHookTx(stream, tx); err != nil {
		return gstatus.Errorf(codes.Unknown, "error executing 'insert' query: %v", err)
	}
	return nil
}
func (this *Amazing_Impl) ClientStreamWithHookTx(stream Amazing_ClientStreamWithHookServer, tx PersistTx) error {
	query := this.QUERIES.Insert(stream.Context(), tx)
	var first *test.ExampleTable
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

		beforeRes, err := this.opts.HOOKS.ClientStreamWithHookBeforeHook(stream.Context(), req)
		if err != nil {
			return gstatus.Errorf(codes.Unknown, "error in before hook: %v", err)
		} else if beforeRes != nil {
			continue
		}

		result := query.Execute(req)
		if err := result.Zero(); err != nil {
			return err
		}
	}
	if err := tx.Commit(); err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return fmt.Errorf("error executing 'insert' query :::AND COULD NOT ROLLBACK::: rollback err: %v, query err: %v", rollbackErr, err)
		}
	}
	res := &test.Ids{}

	if err := this.opts.HOOKS.ClientStreamWithHookAfterHook(stream.Context(), first, res); err != nil {
		return gstatus.Errorf(codes.Unknown, "error in after hook: %v", err)
	}

	if err := stream.SendAndClose(res); err != nil {
		return gstatus.Errorf(codes.Unknown, "error sending back response: %v", err)
	}
	return nil
}
