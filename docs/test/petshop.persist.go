// This file is generated by protoc-gen-persist
// Source File: petshop.proto
// DO NOT EDIT !
package test

import (
	sql "database/sql"
	driver "database/sql/driver"
	fmt "fmt"
	io "io"

	proto "github.com/golang/protobuf/proto"
	persist "github.com/tcncloud/protoc-gen-persist/v4/persist"
	context "golang.org/x/net/context"
	codes "google.golang.org/grpc/codes"
	gstatus "google.golang.org/grpc/status"
)

func NopPersistTx(r persist.Runnable) (persist.PersistTx, error) {
	return &ignoreTx{r}, nil
}

type ignoreTx struct {
	r persist.Runnable
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

func DefaultClientStreamingPersistTx(ctx context.Context, db *sql.DB) (persist.PersistTx, error) {
	return db.BeginTx(ctx, nil)
}
func DefaultServerStreamingPersistTx(ctx context.Context, db *sql.DB) (persist.PersistTx, error) {
	return NopPersistTx(db)
}
func DefaultBidiStreamingPersistTx(ctx context.Context, db *sql.DB) (persist.PersistTx, error) {
	return NopPersistTx(db)
}
func DefaultUnaryPersistTx(ctx context.Context, db *sql.DB) (persist.PersistTx, error) {
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

// Queries_PetShop holds all the queries found the proto service option as methods
type Queries_PetShop struct {
	opts Opts_PetShop
}

// QueriesPetShop returns all the known 'SQL' queires for the 'PetShop' service.
// If no opts are provided default implementations are used.
func QueriesPetShop(opts ...Opts_PetShop) *Queries_PetShop {
	var myOpts Opts_PetShop
	if len(opts) > 0 {
		myOpts = opts[0]
	} else {
		myOpts = OptsPetShop(&DefaultHooks_PetShop{}, &DefaultTypeMappings_PetShop{})
	}
	return &Queries_PetShop{
		opts: myOpts,
	}
}

// GetCatByName returns a struct that will perform the 'GetCatByName' query.
// When Execute is called, it will use the following fields:
// [cat_name]
func (this *Queries_PetShop) GetCatByName(ctx context.Context, db persist.Runnable) *Query_PetShop_GetCatByName {
	return &Query_PetShop_GetCatByName{
		opts: this.opts,
		ctx:  ctx,
		db:   db,
	}
}

// Query_PetShop_GetCatByName (future doc string needed)
type Query_PetShop_GetCatByName struct {
	opts Opts_PetShop
	db   persist.Runnable
	ctx  context.Context
}

func (this *Query_PetShop_GetCatByName) QueryInType_CatName() {}
func (this *Query_PetShop_GetCatByName) QueryOutType_Cat()    {}

// Executes the query 'GetCatByName' with parameters retrieved from x.
// Fields used: [cat_name]
func (this *Query_PetShop_GetCatByName) Execute(x In_PetShop_GetCatByName) *Iter_PetShop_GetCatByName {
	var setupErr error
	params := []interface{}{
		func() (out interface{}) {
			out = x.GetCatName()
			return
		}(),
	}
	result := &Iter_PetShop_GetCatByName{
		tm:  this.opts.MAPPINGS,
		ctx: this.ctx,
	}
	if setupErr != nil {
		result.err = setupErr
		return result
	}
	result.rows, result.err = this.db.QueryContext(this.ctx, "SELECT name, age, cost FROM cats WHERE name = $1", params...)
	return result
}

// InsertFish returns a struct that will perform the 'InsertFish' query.
// When Execute is called, it will use the following fields:
// [species cost]
func (this *Queries_PetShop) InsertFish(ctx context.Context, db persist.Runnable) *Query_PetShop_InsertFish {
	return &Query_PetShop_InsertFish{
		opts: this.opts,
		ctx:  ctx,
		db:   db,
	}
}

// Query_PetShop_InsertFish (future doc string needed)
type Query_PetShop_InsertFish struct {
	opts Opts_PetShop
	db   persist.Runnable
	ctx  context.Context
}

func (this *Query_PetShop_InsertFish) QueryInType_Fish()   {}
func (this *Query_PetShop_InsertFish) QueryOutType_Empty() {}

// Executes the query 'InsertFish' with parameters retrieved from x.
// Fields used: [species cost]
func (this *Query_PetShop_InsertFish) Execute(x In_PetShop_InsertFish) *Iter_PetShop_InsertFish {
	var setupErr error
	params := []interface{}{
		func() (out interface{}) {
			out = x.GetSpecies()
			return
		}(),
		func() (out interface{}) {
			out = x.GetCost()
			return
		}(),
	}
	result := &Iter_PetShop_InsertFish{
		tm:  this.opts.MAPPINGS,
		ctx: this.ctx,
	}
	if setupErr != nil {
		result.err = setupErr
		return result
	}
	result.result, result.err = this.db.ExecContext(this.ctx, "INSERT INTO fish( species, cost ) VALUES( $1, $2 )", params...)
	return result
}

// GetAllOwners returns a struct that will perform the 'GetAllOwners' query.
// When Execute is called, it will use the following fields:
// []
func (this *Queries_PetShop) GetAllOwners(ctx context.Context, db persist.Runnable) *Query_PetShop_GetAllOwners {
	return &Query_PetShop_GetAllOwners{
		opts: this.opts,
		ctx:  ctx,
		db:   db,
	}
}

// Query_PetShop_GetAllOwners (future doc string needed)
type Query_PetShop_GetAllOwners struct {
	opts Opts_PetShop
	db   persist.Runnable
	ctx  context.Context
}

func (this *Query_PetShop_GetAllOwners) QueryInType_Empty()  {}
func (this *Query_PetShop_GetAllOwners) QueryOutType_Owner() {}

// Executes the query 'GetAllOwners' with parameters retrieved from x.
// Fields used: []
func (this *Query_PetShop_GetAllOwners) Execute(x In_PetShop_GetAllOwners) *Iter_PetShop_GetAllOwners {
	var setupErr error
	params := []interface{}{}
	result := &Iter_PetShop_GetAllOwners{
		tm:  this.opts.MAPPINGS,
		ctx: this.ctx,
	}
	if setupErr != nil {
		result.err = setupErr
		return result
	}
	result.rows, result.err = this.db.QueryContext(this.ctx, "SELECT id, aquarium, dog_ids FROM dog_and_fish_owners", params...)
	return result
}

type Iter_PetShop_GetCatByName struct {
	result sql.Result
	rows   *sql.Rows
	err    error
	tm     TypeMappings_PetShop
	ctx    context.Context
}

func (this *Iter_PetShop_GetCatByName) IterOutTypeCat()    {}
func (this *Iter_PetShop_GetCatByName) IterInTypeCatName() {}

// Each performs 'fun' on each row in the result set.
// Each respects the context passed to it.
// It will stop iteration, and returns this.ctx.Err() if encountered.
func (this *Iter_PetShop_GetCatByName) Each(fun func(*Row_PetShop_GetCatByName) error) error {
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
func (this *Iter_PetShop_GetCatByName) One() *Row_PetShop_GetCatByName {
	first, hasFirst := this.Next()
	if first != nil && first.err != nil && first.err != io.EOF {
		return &Row_PetShop_GetCatByName{err: first.err}
	}
	_, hasSecond := this.Next()
	if !hasFirst || hasSecond {
		amount := "none"
		if hasSecond {
			amount = "multiple"
		}
		return &Row_PetShop_GetCatByName{err: persist.NotFound{Msg: fmt.Sprintf("expected exactly 1 result from query 'GetCatByName' found %s", amount)}}
	}
	return first
}

// Zero returns an error if there were any rows in the result
func (this *Iter_PetShop_GetCatByName) Zero() error {
	row, ok := this.Next()
	if row != nil && row.err != nil && row.err != io.EOF {
		return row.err
	}
	if ok {
		return fmt.Errorf("expected exactly 0 results from query 'GetCatByName'")
	}
	return nil
}

// Next returns the next scanned row out of the database, or (nil, false) if there are no more rows
func (this *Iter_PetShop_GetCatByName) Next() (*Row_PetShop_GetCatByName, bool) {
	if this.err == io.EOF {
		return nil, false
	}
	if this.err != nil {
		err := this.err
		this.err = io.EOF
		return &Row_PetShop_GetCatByName{err: err}, true
	}
	if this.rows == nil {
		this.err = io.EOF
		return nil, false
	}
	cols, err := this.rows.Columns()
	if err != nil {
		return &Row_PetShop_GetCatByName{err: err}, true
	}
	if !this.rows.Next() {
		if this.err = this.rows.Err(); this.err == nil {
			this.err = io.EOF
			return nil, false
		} else if this.err != nil {
			return &Row_PetShop_GetCatByName{err: err}, true
		}
	}
	toScan := make([]interface{}, len(cols))
	scanned := make([]alwaysScanner, len(cols))
	for i := range scanned {
		toScan[i] = &scanned[i]
	}
	if this.err = this.rows.Scan(toScan...); this.err != nil {
		return &Row_PetShop_GetCatByName{err: this.err}, true
	}
	res := &Cat{}
	for i, col := range cols {
		_ = i
		switch col {
		case "id":
			r, ok := (*scanned[i].i).(string)
			if !ok {
				return &Row_PetShop_GetCatByName{err: fmt.Errorf("cant convert db column id to protobuf go type ")}, true
			}
			res.Id = r
		case "age":
			r, ok := (*scanned[i].i).(float64)
			if !ok {
				return &Row_PetShop_GetCatByName{err: fmt.Errorf("cant convert db column age to protobuf go type ")}, true
			}
			res.Age = r
		case "cost":
			r, ok := (*scanned[i].i).(float64)
			if !ok {
				return &Row_PetShop_GetCatByName{err: fmt.Errorf("cant convert db column cost to protobuf go type ")}, true
			}
			res.Cost = r
		case "name":
			r, ok := (*scanned[i].i).(string)
			if !ok {
				return &Row_PetShop_GetCatByName{err: fmt.Errorf("cant convert db column name to protobuf go type ")}, true
			}
			res.Name = r

		default:
			return &Row_PetShop_GetCatByName{err: fmt.Errorf("unsupported column in output: %s", col)}, true
		}
	}
	return &Row_PetShop_GetCatByName{item: res}, true
}

// Slice returns all rows found in the iterator as a Slice.
func (this *Iter_PetShop_GetCatByName) Slice() []*Row_PetShop_GetCatByName {
	var results []*Row_PetShop_GetCatByName
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
func (r *Iter_PetShop_GetCatByName) Columns() ([]string, error) {
	if r.err != nil {
		return nil, r.err
	}
	if r.rows != nil {
		return r.rows.Columns()
	}
	return nil, nil
}

type Iter_PetShop_InsertFish struct {
	result sql.Result
	rows   *sql.Rows
	err    error
	tm     TypeMappings_PetShop
	ctx    context.Context
}

func (this *Iter_PetShop_InsertFish) IterOutTypeEmpty() {}
func (this *Iter_PetShop_InsertFish) IterInTypeFish()   {}

// Each performs 'fun' on each row in the result set.
// Each respects the context passed to it.
// It will stop iteration, and returns this.ctx.Err() if encountered.
func (this *Iter_PetShop_InsertFish) Each(fun func(*Row_PetShop_InsertFish) error) error {
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
func (this *Iter_PetShop_InsertFish) One() *Row_PetShop_InsertFish {
	first, hasFirst := this.Next()
	if first != nil && first.err != nil && first.err != io.EOF {
		return &Row_PetShop_InsertFish{err: first.err}
	}
	_, hasSecond := this.Next()
	if !hasFirst || hasSecond {
		amount := "none"
		if hasSecond {
			amount = "multiple"
		}
		return &Row_PetShop_InsertFish{err: persist.NotFound{Msg: fmt.Sprintf("expected exactly 1 result from query 'InsertFish' found %s", amount)}}
	}
	return first
}

// Zero returns an error if there were any rows in the result
func (this *Iter_PetShop_InsertFish) Zero() error {
	row, ok := this.Next()
	if row != nil && row.err != nil && row.err != io.EOF {
		return row.err
	}
	if ok {
		return fmt.Errorf("expected exactly 0 results from query 'InsertFish'")
	}
	return nil
}

// Next returns the next scanned row out of the database, or (nil, false) if there are no more rows
func (this *Iter_PetShop_InsertFish) Next() (*Row_PetShop_InsertFish, bool) {
	if this.err == io.EOF {
		return nil, false
	}
	if this.err != nil {
		err := this.err
		this.err = io.EOF
		return &Row_PetShop_InsertFish{err: err}, true
	}
	if this.rows == nil {
		this.err = io.EOF
		return nil, false
	}
	cols, err := this.rows.Columns()
	if err != nil {
		return &Row_PetShop_InsertFish{err: err}, true
	}
	if !this.rows.Next() {
		if this.err = this.rows.Err(); this.err == nil {
			this.err = io.EOF
			return nil, false
		} else if this.err != nil {
			return &Row_PetShop_InsertFish{err: err}, true
		}
	}
	toScan := make([]interface{}, len(cols))
	scanned := make([]alwaysScanner, len(cols))
	for i := range scanned {
		toScan[i] = &scanned[i]
	}
	if this.err = this.rows.Scan(toScan...); this.err != nil {
		return &Row_PetShop_InsertFish{err: this.err}, true
	}
	res := &Empty{}
	for i, col := range cols {
		_ = i
		switch col {

		default:
			return &Row_PetShop_InsertFish{err: fmt.Errorf("unsupported column in output: %s", col)}, true
		}
	}
	return &Row_PetShop_InsertFish{item: res}, true
}

// Slice returns all rows found in the iterator as a Slice.
func (this *Iter_PetShop_InsertFish) Slice() []*Row_PetShop_InsertFish {
	var results []*Row_PetShop_InsertFish
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
func (r *Iter_PetShop_InsertFish) Columns() ([]string, error) {
	if r.err != nil {
		return nil, r.err
	}
	if r.rows != nil {
		return r.rows.Columns()
	}
	return nil, nil
}

type Iter_PetShop_GetAllOwners struct {
	result sql.Result
	rows   *sql.Rows
	err    error
	tm     TypeMappings_PetShop
	ctx    context.Context
}

func (this *Iter_PetShop_GetAllOwners) IterOutTypeOwner() {}
func (this *Iter_PetShop_GetAllOwners) IterInTypeEmpty()  {}

// Each performs 'fun' on each row in the result set.
// Each respects the context passed to it.
// It will stop iteration, and returns this.ctx.Err() if encountered.
func (this *Iter_PetShop_GetAllOwners) Each(fun func(*Row_PetShop_GetAllOwners) error) error {
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
func (this *Iter_PetShop_GetAllOwners) One() *Row_PetShop_GetAllOwners {
	first, hasFirst := this.Next()
	if first != nil && first.err != nil && first.err != io.EOF {
		return &Row_PetShop_GetAllOwners{err: first.err}
	}
	_, hasSecond := this.Next()
	if !hasFirst || hasSecond {
		amount := "none"
		if hasSecond {
			amount = "multiple"
		}
		return &Row_PetShop_GetAllOwners{err: persist.NotFound{Msg: fmt.Sprintf("expected exactly 1 result from query 'GetAllOwners' found %s", amount)}}
	}
	return first
}

// Zero returns an error if there were any rows in the result
func (this *Iter_PetShop_GetAllOwners) Zero() error {
	row, ok := this.Next()
	if row != nil && row.err != nil && row.err != io.EOF {
		return row.err
	}
	if ok {
		return fmt.Errorf("expected exactly 0 results from query 'GetAllOwners'")
	}
	return nil
}

// Next returns the next scanned row out of the database, or (nil, false) if there are no more rows
func (this *Iter_PetShop_GetAllOwners) Next() (*Row_PetShop_GetAllOwners, bool) {
	if this.err == io.EOF {
		return nil, false
	}
	if this.err != nil {
		err := this.err
		this.err = io.EOF
		return &Row_PetShop_GetAllOwners{err: err}, true
	}
	if this.rows == nil {
		this.err = io.EOF
		return nil, false
	}
	cols, err := this.rows.Columns()
	if err != nil {
		return &Row_PetShop_GetAllOwners{err: err}, true
	}
	if !this.rows.Next() {
		if this.err = this.rows.Err(); this.err == nil {
			this.err = io.EOF
			return nil, false
		} else if this.err != nil {
			return &Row_PetShop_GetAllOwners{err: err}, true
		}
	}
	toScan := make([]interface{}, len(cols))
	scanned := make([]alwaysScanner, len(cols))
	for i := range scanned {
		toScan[i] = &scanned[i]
	}
	if this.err = this.rows.Scan(toScan...); this.err != nil {
		return &Row_PetShop_GetAllOwners{err: this.err}, true
	}
	res := &Owner{}
	for i, col := range cols {
		_ = i
		switch col {
		case "id":
			r, ok := (*scanned[i].i).(string)
			if !ok {
				return &Row_PetShop_GetAllOwners{err: fmt.Errorf("cant convert db column id to protobuf go type ")}, true
			}
			res.Id = r
		case "cats":
			r, ok := (*scanned[i].i).([]byte)
			if !ok {
				return &Row_PetShop_GetAllOwners{err: fmt.Errorf("cant convert db column cats to protobuf go type *[]Cat")}, true
			}
			var converted = new([]Cat)
			if err := proto.Unmarshal(r, converted); err != nil {
				return &Row_PetShop_GetAllOwners{err: err}, true
			}
			res.Cats = converted
		case "aquarium":
			r, ok := (*scanned[i].i).([]byte)
			if !ok {
				return &Row_PetShop_GetAllOwners{err: fmt.Errorf("cant convert db column aquarium to protobuf go type *FishTank")}, true
			}
			var converted = new(FishTank)
			if err := proto.Unmarshal(r, converted); err != nil {
				return &Row_PetShop_GetAllOwners{err: err}, true
			}
			res.Aquarium = converted
		case "dog_ids":
			var converted = this.tm.DogIds()
			if err := converted.Scan(*scanned[i].i); err != nil {
				return &Row_PetShop_GetAllOwners{err: fmt.Errorf("could not convert mapped db column dog_ids to type on Owner.DogIds: %v", err)}, true
			}
			if err := converted.ToProto(&res.DogIds); err != nil {
				return &Row_PetShop_GetAllOwners{err: fmt.Errorf("could not convert mapped db column dog_idsto type on Owner.DogIds: %v", err)}, true
			}
		case "money":
			r, ok := (*scanned[i].i).(float64)
			if !ok {
				return &Row_PetShop_GetAllOwners{err: fmt.Errorf("cant convert db column money to protobuf go type ")}, true
			}
			res.Money = r

		default:
			return &Row_PetShop_GetAllOwners{err: fmt.Errorf("unsupported column in output: %s", col)}, true
		}
	}
	return &Row_PetShop_GetAllOwners{item: res}, true
}

// Slice returns all rows found in the iterator as a Slice.
func (this *Iter_PetShop_GetAllOwners) Slice() []*Row_PetShop_GetAllOwners {
	var results []*Row_PetShop_GetAllOwners
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
func (r *Iter_PetShop_GetAllOwners) Columns() ([]string, error) {
	if r.err != nil {
		return nil, r.err
	}
	if r.rows != nil {
		return r.rows.Columns()
	}
	return nil, nil
}

type In_PetShop_GetCatByName interface {
	GetCatName() string
}
type Out_PetShop_GetCatByName interface {
	GetId() string
	GetAge() float64
	GetCost() float64
	GetName() string
}
type Row_PetShop_GetCatByName struct {
	item Out_PetShop_GetCatByName
	err  error
}

func newRowPetShopGetCatByName(item Out_PetShop_GetCatByName, err error) *Row_PetShop_GetCatByName {
	return &Row_PetShop_GetCatByName{item, err}
}

// Unwrap takes an address to a proto.Message as its only parameter
// Unwrap can only set into output protos of that match method return types + the out option on the query itself
func (this *Row_PetShop_GetCatByName) Unwrap(pointerToMsg proto.Message) error {
	if this.err != nil {
		return this.err
	}
	if o, ok := (pointerToMsg).(*Cat); ok {
		if o == nil {
			return fmt.Errorf("must initialize *Cat before giving to Unwrap()")
		}
		res, _ := this.Cat()
		_ = res
		o.Id = res.Id
		o.Age = res.Age
		o.Cost = res.Cost
		o.Name = res.Name
		return nil
	}

	if o, ok := (pointerToMsg).(*Cat); ok {
		if o == nil {
			return fmt.Errorf("must initialize *Cat before giving to Unwrap()")
		}
		res, _ := this.Cat()
		_ = res
		o.Id = res.Id
		o.Age = res.Age
		o.Cost = res.Cost
		o.Name = res.Name
		return nil
	}

	return nil
}
func (this *Row_PetShop_GetCatByName) Cat() (*Cat, error) {
	if this.err != nil {
		return nil, this.err
	}
	return &Cat{
		Id:   this.item.GetId(),
		Age:  this.item.GetAge(),
		Cost: this.item.GetCost(),
		Name: this.item.GetName(),
	}, nil
}

func (this *Row_PetShop_GetCatByName) Proto() (*Cat, error) {
	if this.err != nil {
		return nil, this.err
	}
	return &Cat{
		Id:   this.item.GetId(),
		Age:  this.item.GetAge(),
		Cost: this.item.GetCost(),
		Name: this.item.GetName(),
	}, nil
}

type In_PetShop_InsertFish interface {
	GetId() string
	GetCost() float64
	GetSpecies() string
}
type Out_PetShop_InsertFish interface {
}
type Row_PetShop_InsertFish struct {
	item Out_PetShop_InsertFish
	err  error
}

func newRowPetShopInsertFish(item Out_PetShop_InsertFish, err error) *Row_PetShop_InsertFish {
	return &Row_PetShop_InsertFish{item, err}
}

// Unwrap takes an address to a proto.Message as its only parameter
// Unwrap can only set into output protos of that match method return types + the out option on the query itself
func (this *Row_PetShop_InsertFish) Unwrap(pointerToMsg proto.Message) error {
	if this.err != nil {
		return this.err
	}
	if o, ok := (pointerToMsg).(*Empty); ok {
		if o == nil {
			return fmt.Errorf("must initialize *Empty before giving to Unwrap()")
		}
		res, _ := this.Empty()
		_ = res

		return nil
	}

	if o, ok := (pointerToMsg).(*Empty); ok {
		if o == nil {
			return fmt.Errorf("must initialize *Empty before giving to Unwrap()")
		}
		res, _ := this.Empty()
		_ = res

		return nil
	}

	return nil
}
func (this *Row_PetShop_InsertFish) Empty() (*Empty, error) {
	if this.err != nil {
		return nil, this.err
	}
	return &Empty{}, nil
}

func (this *Row_PetShop_InsertFish) Proto() (*Empty, error) {
	if this.err != nil {
		return nil, this.err
	}
	return &Empty{}, nil
}

type In_PetShop_GetAllOwners interface {
}
type Out_PetShop_GetAllOwners interface {
	GetId() string
	GetCats() []*Cat
	GetAquarium() *FishTank
	GetDogIds() *DogIds
	GetMoney() float64
}
type Row_PetShop_GetAllOwners struct {
	item Out_PetShop_GetAllOwners
	err  error
}

func newRowPetShopGetAllOwners(item Out_PetShop_GetAllOwners, err error) *Row_PetShop_GetAllOwners {
	return &Row_PetShop_GetAllOwners{item, err}
}

// Unwrap takes an address to a proto.Message as its only parameter
// Unwrap can only set into output protos of that match method return types + the out option on the query itself
func (this *Row_PetShop_GetAllOwners) Unwrap(pointerToMsg proto.Message) error {
	if this.err != nil {
		return this.err
	}
	if o, ok := (pointerToMsg).(*Owner); ok {
		if o == nil {
			return fmt.Errorf("must initialize *Owner before giving to Unwrap()")
		}
		res, _ := this.Owner()
		_ = res
		o.Id = res.Id
		o.Cats = res.Cats
		o.Aquarium = res.Aquarium
		o.DogIds = res.DogIds
		o.Money = res.Money
		return nil
	}

	return nil
}
func (this *Row_PetShop_GetAllOwners) Owner() (*Owner, error) {
	if this.err != nil {
		return nil, this.err
	}
	return &Owner{
		Id:       this.item.GetId(),
		Cats:     this.item.GetCats(),
		Aquarium: this.item.GetAquarium(),
		DogIds:   this.item.GetDogIds(),
		Money:    this.item.GetMoney(),
	}, nil
}

func (this *Row_PetShop_GetAllOwners) Proto() (*Owner, error) {
	if this.err != nil {
		return nil, this.err
	}
	return &Owner{
		Id:       this.item.GetId(),
		Cats:     this.item.GetCats(),
		Aquarium: this.item.GetAquarium(),
		DogIds:   this.item.GetDogIds(),
		Money:    this.item.GetMoney(),
	}, nil
}

type Hooks_PetShop interface {
}
type DefaultHooks_PetShop struct{}
type TypeMappings_PetShop interface {
	DogIds() MappingImpl_PetShop_DogIds
}
type DefaultTypeMappings_PetShop struct{}

func (this *DefaultTypeMappings_PetShop) DogIds() MappingImpl_PetShop_DogIds {
	return &DefaultMappingImpl_PetShop_DogIds{}
}

type DefaultMappingImpl_PetShop_DogIds struct{}

func (this *DefaultMappingImpl_PetShop_DogIds) ToProto(**DogIds) error {
	return nil
}
func (this *DefaultMappingImpl_PetShop_DogIds) ToSql(*DogIds) sql.Scanner {
	return this
}
func (this *DefaultMappingImpl_PetShop_DogIds) Scan(interface{}) error {
	return nil
}
func (this *DefaultMappingImpl_PetShop_DogIds) Value() (driver.Value, error) {
	return "DEFAULT_TYPE_MAPPING_VALUE", nil
}

type MappingImpl_PetShop_DogIds interface {
	ToProto(**DogIds) error
	ToSql(*DogIds) sql.Scanner
	sql.Scanner
	driver.Valuer
}

type Opts_PetShop struct {
	MAPPINGS TypeMappings_PetShop
	HOOKS    Hooks_PetShop
}

func OptsPetShop(hooks Hooks_PetShop, mappings TypeMappings_PetShop) Opts_PetShop {
	opts := Opts_PetShop{
		HOOKS:    &DefaultHooks_PetShop{},
		MAPPINGS: &DefaultTypeMappings_PetShop{},
	}
	if hooks != nil {
		opts.HOOKS = hooks
	}
	if mappings != nil {
		opts.MAPPINGS = mappings
	}
	return opts
}

type Impl_PetShop struct {
	opts     *Opts_PetShop
	QUERIES  *Queries_PetShop
	HANDLERS RestOfHandlers_PetShop
	DB       *sql.DB
}

func ImplPetShop(db *sql.DB, handlers RestOfHandlers_PetShop, opts ...Opts_PetShop) *Impl_PetShop {
	var myOpts Opts_PetShop
	if len(opts) > 0 {
		myOpts = opts[0]
	} else {
		myOpts = OptsPetShop(&DefaultHooks_PetShop{}, &DefaultTypeMappings_PetShop{})
	}
	return &Impl_PetShop{
		opts:     &myOpts,
		QUERIES:  QueriesPetShop(myOpts),
		DB:       db,
		HANDLERS: handlers,
	}
}

type RestOfHandlers_PetShop interface {
	PetDog(context.Context, *Dog) (*Empty, error)
}

func (this *Impl_PetShop) PetDog(ctx context.Context, req *Dog) (*Empty, error) {
	return this.HANDLERS.PetDog(ctx, req)
}

func (this *Impl_PetShop) GetCatByName(ctx context.Context, req *CatName) (*Cat, error) {
	query := this.QUERIES.GetCatByName(ctx, this.DB)

	result := query.Execute(req)
	res, err := result.One().Cat()
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (this *Impl_PetShop) ShipFish(stream PetShop_ShipFishServer) error {
	tx, err := DefaultClientStreamingPersistTx(stream.Context(), this.DB)
	if err != nil {
		return gstatus.Errorf(codes.Unknown, "error creating persist tx: %v", err)
	}
	if err := this.ShipFishTx(stream, tx); err != nil {
		return gstatus.Errorf(codes.Unknown, "error executing 'InsertFish' query: %v", err)
	}
	return nil
}
func (this *Impl_PetShop) ShipFishTx(stream PetShop_ShipFishServer, tx persist.PersistTx) error {
	query := this.QUERIES.InsertFish(stream.Context(), tx)
	var first *Fish
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
			return fmt.Errorf("error executing 'InsertFish' query :::AND COULD NOT ROLLBACK::: rollback err: %v, query err: %v", rollbackErr, err)
		}
	}
	res := &Empty{}

	if err := stream.SendAndClose(res); err != nil {
		return gstatus.Errorf(codes.Unknown, "error sending back response: %v", err)
	}
	return nil
}
