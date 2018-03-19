// This file is generated by protoc-gen-persist
// Source File: tests/spanner/basic/basic_example.proto
// DO NOT EDIT !
package basic

import (
	fmt "fmt"
	io "io"

	spanner "cloud.google.com/go/spanner"
	proto "github.com/golang/protobuf/proto"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
	mytime "github.com/tcncloud/protoc-gen-persist/tests/mytime"
	persist_lib "github.com/tcncloud/protoc-gen-persist/tests/spanner/basic/persist_lib"
	hooks "github.com/tcncloud/protoc-gen-persist/tests/spanner/hooks"
	test "github.com/tcncloud/protoc-gen-persist/tests/test"
	context "golang.org/x/net/context"
	codes "google.golang.org/grpc/codes"
	gstatus "google.golang.org/grpc/status"
)

type ExtraSrvImpl struct {
	PERSIST   *persist_lib.ExtraSrvMethodReceiver
	FORWARDED RestOfExtraSrvHandlers
}
type RestOfExtraSrvHandlers interface {
	ExtraMethod(ctx context.Context, req *test.ExampleTable) (*test.ExampleTable, error)
}
type ExtraSrvImplBuilder struct {
	err           error
	rest          RestOfExtraSrvHandlers
	queryHandlers *persist_lib.ExtraSrvQueryHandlers
	i             *ExtraSrvImpl
	db            spanner.Client
}

func NewExtraSrvBuilder() *ExtraSrvImplBuilder {
	return &ExtraSrvImplBuilder{i: &ExtraSrvImpl{}}
}
func (b *ExtraSrvImplBuilder) WithRestOfGrpcHandlers(r RestOfExtraSrvHandlers) *ExtraSrvImplBuilder {
	b.rest = r
	return b
}
func (b *ExtraSrvImplBuilder) WithPersistQueryHandlers(p *persist_lib.ExtraSrvQueryHandlers) *ExtraSrvImplBuilder {
	b.queryHandlers = p
	return b
}
func (b *ExtraSrvImplBuilder) WithDefaultQueryHandlers() *ExtraSrvImplBuilder {
	accessor := persist_lib.NewSpannerClientGetter(&b.db)
	queryHandlers := &persist_lib.ExtraSrvQueryHandlers{
		ExtraUnaryHandler: persist_lib.DefaultExtraUnaryHandler(accessor),
	}
	b.queryHandlers = queryHandlers
	return b
}

// set the custom handlers you want to use in the handlers
// this method will make sure to use a default handler if
// the handler is nil.
func (b *ExtraSrvImplBuilder) WithNilAsDefaultQueryHandlers(p *persist_lib.ExtraSrvQueryHandlers) *ExtraSrvImplBuilder {
	accessor := persist_lib.NewSpannerClientGetter(&b.db)
	if p.ExtraUnaryHandler == nil {
		p.ExtraUnaryHandler = persist_lib.DefaultExtraUnaryHandler(accessor)
	}
	b.queryHandlers = p
	return b
}
func (b *ExtraSrvImplBuilder) WithSpannerClient(c *spanner.Client) *ExtraSrvImplBuilder {
	b.db = *c
	return b
}
func (b *ExtraSrvImplBuilder) WithSpannerURI(ctx context.Context, uri string) *ExtraSrvImplBuilder {
	cli, err := spanner.NewClient(ctx, uri)
	b.err = err
	b.db = *cli
	return b
}
func (b *ExtraSrvImplBuilder) Build() (*ExtraSrvImpl, error) {
	if b.err != nil {
		return nil, b.err
	}
	b.i.PERSIST = &persist_lib.ExtraSrvMethodReceiver{Handlers: *b.queryHandlers}
	b.i.FORWARDED = b.rest
	return b.i, nil
}
func (b *ExtraSrvImplBuilder) MustBuild() *ExtraSrvImpl {
	s, err := b.Build()
	if err != nil {
		panic("error in builder: " + err.Error())
	}
	return s
}
func NumRowsToExtraSrvPersistType(req *test.NumRows) (*persist_lib.Test_NumRowsForExtraSrv, error) {
	var err error
	_ = err
	params := &persist_lib.Test_NumRowsForExtraSrv{}
	// set 'NumRows.count' in params
	params.Count = req.Count
	return params, nil
}
func ExampleTableFromExtraSrvDatabaseRow(row *spanner.Row) (*test.ExampleTable, error) {
	res := &test.ExampleTable{}
	var Id_ int64
	{
		local := &spanner.NullInt64{}
		if err := row.ColumnByName("id", local); err != nil {
			return nil, err
		}
		if local.Valid {
			Id_ = local.Int64
		}
		res.Id = Id_
	}
	var StartTime_ []byte
	if err := row.ColumnByName("start_time", &StartTime_); err != nil {
		return nil, err
	}
	{
		local := new(timestamp.Timestamp)
		if err := proto.Unmarshal(StartTime_, local); err != nil {
			return nil, err
		}
		res.StartTime = local
	}
	var Name_ string
	{
		local := &spanner.NullString{}
		if err := row.ColumnByName("name", local); err != nil {
			return nil, err
		}
		if local.Valid {
			Name_ = local.StringVal
		}
		res.Name = Name_
	}
	return res, nil
}
func IterExtraSrvExampleTableProto(iter *spanner.RowIterator, next func(i *test.ExampleTable) error) error {
	return iter.Do(func(r *spanner.Row) error {
		item, err := ExampleTableFromExtraSrvDatabaseRow(r)
		if err != nil {
			return fmt.Errorf("error converting test.ExampleTable row to protobuf message: %s", err)
		}
		return next(item)
	})
}
func (s *ExtraSrvImpl) ExtraUnary(ctx context.Context, req *test.NumRows) (*test.ExampleTable, error) {
	var err error
	var res = &test.ExampleTable{}
	_ = err
	_ = res
	params, err := NumRowsToExtraSrvPersistType(req)
	if err != nil {
		return nil, err
	}
	var iterErr error
	err = s.PERSIST.ExtraUnary(ctx, params, func(row *spanner.Row) {
		if row == nil { // there was no return data
			return
		}
		res, err = ExampleTableFromExtraSrvDatabaseRow(row)
		if err != nil {
			iterErr = err
			return
		}
	})
	if err != nil {
		return nil, gstatus.Errorf(codes.Unknown, "error calling persist service: %v", err)
	} else if iterErr != nil {
		return nil, iterErr
	}
	return res, nil
}
func (s *ExtraSrvImpl) ExtraMethod(ctx context.Context, req *test.ExampleTable) (*test.ExampleTable, error) {
	return s.FORWARDED.ExtraMethod(ctx, req)
}

type MySpannerImpl struct {
	PERSIST   *persist_lib.MySpannerMethodReceiver
	FORWARDED RestOfMySpannerHandlers
}
type RestOfMySpannerHandlers interface {
}
type MySpannerImplBuilder struct {
	err           error
	rest          RestOfMySpannerHandlers
	queryHandlers *persist_lib.MySpannerQueryHandlers
	i             *MySpannerImpl
	db            spanner.Client
}

func NewMySpannerBuilder() *MySpannerImplBuilder {
	return &MySpannerImplBuilder{i: &MySpannerImpl{}}
}
func (b *MySpannerImplBuilder) WithRestOfGrpcHandlers(r RestOfMySpannerHandlers) *MySpannerImplBuilder {
	b.rest = r
	return b
}
func (b *MySpannerImplBuilder) WithPersistQueryHandlers(p *persist_lib.MySpannerQueryHandlers) *MySpannerImplBuilder {
	b.queryHandlers = p
	return b
}
func (b *MySpannerImplBuilder) WithDefaultQueryHandlers() *MySpannerImplBuilder {
	accessor := persist_lib.NewSpannerClientGetter(&b.db)
	queryHandlers := &persist_lib.MySpannerQueryHandlers{
		UniaryInsertHandler:                persist_lib.DefaultUniaryInsertHandler(accessor),
		UniarySelectHandler:                persist_lib.DefaultUniarySelectHandler(accessor),
		TestNestHandler:                    persist_lib.DefaultTestNestHandler(accessor),
		TestEverythingHandler:              persist_lib.DefaultTestEverythingHandler(accessor),
		UniarySelectWithDirectivesHandler:  persist_lib.DefaultUniarySelectWithDirectivesHandler(accessor),
		UniaryUpdateHandler:                persist_lib.DefaultUniaryUpdateHandler(accessor),
		UniaryDeleteRangeHandler:           persist_lib.DefaultUniaryDeleteRangeHandler(accessor),
		UniaryDeleteSingleHandler:          persist_lib.DefaultUniaryDeleteSingleHandler(accessor),
		NoArgsHandler:                      persist_lib.DefaultNoArgsHandler(accessor),
		ServerStreamHandler:                persist_lib.DefaultServerStreamHandler(accessor),
		ClientStreamInsertHandler:          persist_lib.DefaultClientStreamInsertHandler(accessor),
		ClientStreamDeleteHandler:          persist_lib.DefaultClientStreamDeleteHandler(accessor),
		ClientStreamUpdateHandler:          persist_lib.DefaultClientStreamUpdateHandler(accessor),
		UniaryInsertWithHooksHandler:       persist_lib.DefaultUniaryInsertWithHooksHandler(accessor),
		UniarySelectWithHooksHandler:       persist_lib.DefaultUniarySelectWithHooksHandler(accessor),
		UniaryUpdateWithHooksHandler:       persist_lib.DefaultUniaryUpdateWithHooksHandler(accessor),
		UniaryDeleteWithHooksHandler:       persist_lib.DefaultUniaryDeleteWithHooksHandler(accessor),
		ServerStreamWithHooksHandler:       persist_lib.DefaultServerStreamWithHooksHandler(accessor),
		ClientStreamUpdateWithHooksHandler: persist_lib.DefaultClientStreamUpdateWithHooksHandler(accessor),
	}
	b.queryHandlers = queryHandlers
	return b
}

// set the custom handlers you want to use in the handlers
// this method will make sure to use a default handler if
// the handler is nil.
func (b *MySpannerImplBuilder) WithNilAsDefaultQueryHandlers(p *persist_lib.MySpannerQueryHandlers) *MySpannerImplBuilder {
	accessor := persist_lib.NewSpannerClientGetter(&b.db)
	if p.UniaryInsertHandler == nil {
		p.UniaryInsertHandler = persist_lib.DefaultUniaryInsertHandler(accessor)
	}
	if p.UniarySelectHandler == nil {
		p.UniarySelectHandler = persist_lib.DefaultUniarySelectHandler(accessor)
	}
	if p.TestNestHandler == nil {
		p.TestNestHandler = persist_lib.DefaultTestNestHandler(accessor)
	}
	if p.TestEverythingHandler == nil {
		p.TestEverythingHandler = persist_lib.DefaultTestEverythingHandler(accessor)
	}
	if p.UniarySelectWithDirectivesHandler == nil {
		p.UniarySelectWithDirectivesHandler = persist_lib.DefaultUniarySelectWithDirectivesHandler(accessor)
	}
	if p.UniaryUpdateHandler == nil {
		p.UniaryUpdateHandler = persist_lib.DefaultUniaryUpdateHandler(accessor)
	}
	if p.UniaryDeleteRangeHandler == nil {
		p.UniaryDeleteRangeHandler = persist_lib.DefaultUniaryDeleteRangeHandler(accessor)
	}
	if p.UniaryDeleteSingleHandler == nil {
		p.UniaryDeleteSingleHandler = persist_lib.DefaultUniaryDeleteSingleHandler(accessor)
	}
	if p.NoArgsHandler == nil {
		p.NoArgsHandler = persist_lib.DefaultNoArgsHandler(accessor)
	}
	if p.ServerStreamHandler == nil {
		p.ServerStreamHandler = persist_lib.DefaultServerStreamHandler(accessor)
	}
	if p.ClientStreamInsertHandler == nil {
		p.ClientStreamInsertHandler = persist_lib.DefaultClientStreamInsertHandler(accessor)
	}
	if p.ClientStreamDeleteHandler == nil {
		p.ClientStreamDeleteHandler = persist_lib.DefaultClientStreamDeleteHandler(accessor)
	}
	if p.ClientStreamUpdateHandler == nil {
		p.ClientStreamUpdateHandler = persist_lib.DefaultClientStreamUpdateHandler(accessor)
	}
	if p.UniaryInsertWithHooksHandler == nil {
		p.UniaryInsertWithHooksHandler = persist_lib.DefaultUniaryInsertWithHooksHandler(accessor)
	}
	if p.UniarySelectWithHooksHandler == nil {
		p.UniarySelectWithHooksHandler = persist_lib.DefaultUniarySelectWithHooksHandler(accessor)
	}
	if p.UniaryUpdateWithHooksHandler == nil {
		p.UniaryUpdateWithHooksHandler = persist_lib.DefaultUniaryUpdateWithHooksHandler(accessor)
	}
	if p.UniaryDeleteWithHooksHandler == nil {
		p.UniaryDeleteWithHooksHandler = persist_lib.DefaultUniaryDeleteWithHooksHandler(accessor)
	}
	if p.ServerStreamWithHooksHandler == nil {
		p.ServerStreamWithHooksHandler = persist_lib.DefaultServerStreamWithHooksHandler(accessor)
	}
	if p.ClientStreamUpdateWithHooksHandler == nil {
		p.ClientStreamUpdateWithHooksHandler = persist_lib.DefaultClientStreamUpdateWithHooksHandler(accessor)
	}
	b.queryHandlers = p
	return b
}
func (b *MySpannerImplBuilder) WithSpannerClient(c *spanner.Client) *MySpannerImplBuilder {
	b.db = *c
	return b
}
func (b *MySpannerImplBuilder) WithSpannerURI(ctx context.Context, uri string) *MySpannerImplBuilder {
	cli, err := spanner.NewClient(ctx, uri)
	b.err = err
	b.db = *cli
	return b
}
func (b *MySpannerImplBuilder) Build() (*MySpannerImpl, error) {
	if b.err != nil {
		return nil, b.err
	}
	b.i.PERSIST = &persist_lib.MySpannerMethodReceiver{Handlers: *b.queryHandlers}
	b.i.FORWARDED = b.rest
	return b.i, nil
}
func (b *MySpannerImplBuilder) MustBuild() *MySpannerImpl {
	s, err := b.Build()
	if err != nil {
		panic("error in builder: " + err.Error())
	}
	return s
}
func ExampleTableToMySpannerPersistType(req *test.ExampleTable) (*persist_lib.Test_ExampleTableForMySpanner, error) {
	var err error
	_ = err
	params := &persist_lib.Test_ExampleTableForMySpanner{}
	// set 'ExampleTable.id' in params
	params.Id = req.Id
	// set 'ExampleTable.start_time' in params
	if params.StartTime, err = (mytime.MyTime{}).ToSpanner(req.StartTime).SpannerValue(); err != nil {
		return nil, err
	}
	// set 'ExampleTable.name' in params
	params.Name = req.Name
	return params, nil
}
func ExampleTableFromMySpannerDatabaseRow(row *spanner.Row) (*test.ExampleTable, error) {
	res := &test.ExampleTable{}
	var Id_ int64
	{
		local := &spanner.NullInt64{}
		if err := row.ColumnByName("id", local); err != nil {
			return nil, err
		}
		if local.Valid {
			Id_ = local.Int64
		}
		res.Id = Id_
	}
	var StartTime_ = new(spanner.GenericColumnValue)
	if err := row.ColumnByName("start_time", StartTime_); err != nil {
		return nil, err
	}
	{
		local := &mytime.MyTime{}
		if err := local.SpannerScan(StartTime_); err != nil {
			return nil, err
		}
		res.StartTime = local.ToProto()
	}
	var Name_ string
	{
		local := &spanner.NullString{}
		if err := row.ColumnByName("name", local); err != nil {
			return nil, err
		}
		if local.Valid {
			Name_ = local.StringVal
		}
		res.Name = Name_
	}
	return res, nil
}
func IterMySpannerExampleTableProto(iter *spanner.RowIterator, next func(i *test.ExampleTable) error) error {
	return iter.Do(func(r *spanner.Row) error {
		item, err := ExampleTableFromMySpannerDatabaseRow(r)
		if err != nil {
			return fmt.Errorf("error converting test.ExampleTable row to protobuf message: %s", err)
		}
		return next(item)
	})
}
func SomethingToMySpannerPersistType(req *Something) (*persist_lib.SomethingForMySpanner, error) {
	var err error
	_ = err
	params := &persist_lib.SomethingForMySpanner{}
	// set 'Something.thing' in params
	if req.Thing == nil {
		req.Thing = new(Something_SomethingElse)
	}
	{
		raw, err := proto.Marshal(req.Thing)
		if err != nil {
			return nil, err
		}
		params.Thing = raw
	}
	// set 'Something.myenum' in params
	params.Myenum = int32(req.Myenum)
	// set 'Something.mappedenum' in params
	params.Mappedenum = int32(req.Mappedenum)
	return params, nil
}
func SomethingFromMySpannerDatabaseRow(row *spanner.Row) (*Something, error) {
	res := &Something{}
	var Thing_ []byte
	if err := row.ColumnByName("thing", &Thing_); err != nil {
		return nil, err
	}
	{
		local := new(Something_SomethingElse)
		if err := proto.Unmarshal(Thing_, local); err != nil {
			return nil, err
		}
		res.Thing = local
	}
	var Myenum_ int64
	if err := row.ColumnByName("myenum", &Myenum_); err != nil {
		return nil, err
	}
	res.Myenum = MyEnum(Myenum_)
	var Mappedenum_ int64
	if err := row.ColumnByName("mappedenum", &Mappedenum_); err != nil {
		return nil, err
	}
	res.Mappedenum = MappedEnum(Mappedenum_)
	return res, nil
}
func IterMySpannerSomethingProto(iter *spanner.RowIterator, next func(i *Something) error) error {
	return iter.Do(func(r *spanner.Row) error {
		item, err := SomethingFromMySpannerDatabaseRow(r)
		if err != nil {
			return fmt.Errorf("error converting Something row to protobuf message: %s", err)
		}
		return next(item)
	})
}
func HasTimestampToMySpannerPersistType(req *HasTimestamp) (*persist_lib.HasTimestampForMySpanner, error) {
	var err error
	_ = err
	params := &persist_lib.HasTimestampForMySpanner{}
	// set 'HasTimestamp.time' in params
	if params.Time, err = (mytime.MyTime{}).ToSpanner(req.Time).SpannerValue(); err != nil {
		return nil, err
	}
	// set 'HasTimestamp.some' in params
	if req.Some == nil {
		req.Some = new(Something)
	}
	{
		raw, err := proto.Marshal(req.Some)
		if err != nil {
			return nil, err
		}
		params.Some = raw
	}
	// set 'HasTimestamp.str' in params
	params.Str = req.Str
	// set 'HasTimestamp.table' in params
	if req.Table == nil {
		req.Table = new(test.ExampleTable)
	}
	{
		raw, err := proto.Marshal(req.Table)
		if err != nil {
			return nil, err
		}
		params.Table = raw
	}
	// set 'HasTimestamp.strs' in params
	params.Strs = req.Strs
	// set 'HasTimestamp.tables' in params
	{
		var bytesOfBytes [][]byte
		for _, msg := range req.Tables {
			raw, err := proto.Marshal(msg)
			if err != nil {
				return nil, err
			}
			bytesOfBytes = append(bytesOfBytes, raw)
		}
		params.Tables = bytesOfBytes
	}
	// set 'HasTimestamp.somes' in params
	{
		var bytesOfBytes [][]byte
		for _, msg := range req.Somes {
			raw, err := proto.Marshal(msg)
			if err != nil {
				return nil, err
			}
			bytesOfBytes = append(bytesOfBytes, raw)
		}
		params.Somes = bytesOfBytes
	}
	// set 'HasTimestamp.times' in params
	{
		var bytesOfBytes [][]byte
		for _, msg := range req.Times {
			raw, err := proto.Marshal(msg)
			if err != nil {
				return nil, err
			}
			bytesOfBytes = append(bytesOfBytes, raw)
		}
		params.Times = bytesOfBytes
	}
	return params, nil
}
func HasTimestampFromMySpannerDatabaseRow(row *spanner.Row) (*HasTimestamp, error) {
	res := &HasTimestamp{}
	var Time_ = new(spanner.GenericColumnValue)
	if err := row.ColumnByName("time", Time_); err != nil {
		return nil, err
	}
	{
		local := &mytime.MyTime{}
		if err := local.SpannerScan(Time_); err != nil {
			return nil, err
		}
		res.Time = local.ToProto()
	}
	var Some_ []byte
	if err := row.ColumnByName("some", &Some_); err != nil {
		return nil, err
	}
	{
		local := new(Something)
		if err := proto.Unmarshal(Some_, local); err != nil {
			return nil, err
		}
		res.Some = local
	}
	var Str_ string
	{
		local := &spanner.NullString{}
		if err := row.ColumnByName("str", local); err != nil {
			return nil, err
		}
		if local.Valid {
			Str_ = local.StringVal
		}
		res.Str = Str_
	}
	var Table_ []byte
	if err := row.ColumnByName("table", &Table_); err != nil {
		return nil, err
	}
	{
		local := new(test.ExampleTable)
		if err := proto.Unmarshal(Table_, local); err != nil {
			return nil, err
		}
		res.Table = local
	}
	var Strs_ []string
	{
		local := make([]spanner.NullString, 0)
		if err := row.ColumnByName("strs", &local); err != nil {
			return nil, err
		}
		for _, l := range local {
			if l.Valid {
				Strs_ = append(Strs_, l.StringVal)
				res.Strs = Strs_
			}
		}
	}
	var Tables_ [][]byte
	if err := row.ColumnByName("tables", &Tables_); err != nil {
		return nil, err
	}
	{
		local := make([]*test.ExampleTable, len(Tables_))
		for i := range local {
			local[i] = new(test.ExampleTable)
			if err := proto.Unmarshal(Tables_[i], local[i]); err != nil {
				return nil, err
			}
		}
		res.Tables = local
	}
	var Somes_ [][]byte
	if err := row.ColumnByName("somes", &Somes_); err != nil {
		return nil, err
	}
	{
		local := make([]*Something, len(Somes_))
		for i := range local {
			local[i] = new(Something)
			if err := proto.Unmarshal(Somes_[i], local[i]); err != nil {
				return nil, err
			}
		}
		res.Somes = local
	}
	var Times_ [][]byte
	if err := row.ColumnByName("times", &Times_); err != nil {
		return nil, err
	}
	{
		local := make([]*timestamp.Timestamp, len(Times_))
		for i := range local {
			local[i] = new(timestamp.Timestamp)
			if err := proto.Unmarshal(Times_[i], local[i]); err != nil {
				return nil, err
			}
		}
		res.Times = local
	}
	return res, nil
}
func IterMySpannerHasTimestampProto(iter *spanner.RowIterator, next func(i *HasTimestamp) error) error {
	return iter.Do(func(r *spanner.Row) error {
		item, err := HasTimestampFromMySpannerDatabaseRow(r)
		if err != nil {
			return fmt.Errorf("error converting HasTimestamp row to protobuf message: %s", err)
		}
		return next(item)
	})
}
func PartialTableFromMySpannerDatabaseRow(row *spanner.Row) (*test.PartialTable, error) {
	res := &test.PartialTable{}
	var Id_ int64
	{
		local := &spanner.NullInt64{}
		if err := row.ColumnByName("id", local); err != nil {
			return nil, err
		}
		if local.Valid {
			Id_ = local.Int64
		}
		res.Id = Id_
	}
	var StartTime_ = new(spanner.GenericColumnValue)
	if err := row.ColumnByName("start_time", StartTime_); err != nil {
		return nil, err
	}
	{
		local := &mytime.MyTime{}
		if err := local.SpannerScan(StartTime_); err != nil {
			return nil, err
		}
		res.StartTime = local.ToProto()
	}
	return res, nil
}
func IterMySpannerPartialTableProto(iter *spanner.RowIterator, next func(i *test.PartialTable) error) error {
	return iter.Do(func(r *spanner.Row) error {
		item, err := PartialTableFromMySpannerDatabaseRow(r)
		if err != nil {
			return fmt.Errorf("error converting test.PartialTable row to protobuf message: %s", err)
		}
		return next(item)
	})
}
func ExampleTableRangeToMySpannerPersistType(req *test.ExampleTableRange) (*persist_lib.Test_ExampleTableRangeForMySpanner, error) {
	var err error
	_ = err
	params := &persist_lib.Test_ExampleTableRangeForMySpanner{}
	// set 'ExampleTableRange.start_id' in params
	params.StartId = req.StartId
	// set 'ExampleTableRange.end_id' in params
	params.EndId = req.EndId
	return params, nil
}
func NameToMySpannerPersistType(req *test.Name) (*persist_lib.Test_NameForMySpanner, error) {
	var err error
	_ = err
	params := &persist_lib.Test_NameForMySpanner{}
	// set 'Name.name' in params
	params.Name = req.Name
	return params, nil
}
func NumRowsFromMySpannerDatabaseRow(row *spanner.Row) (*test.NumRows, error) {
	res := &test.NumRows{}
	var Count_ int64
	{
		local := &spanner.NullInt64{}
		if err := row.ColumnByName("count", local); err != nil {
			return nil, err
		}
		if local.Valid {
			Count_ = local.Int64
		}
		res.Count = Count_
	}
	return res, nil
}
func IterMySpannerNumRowsProto(iter *spanner.RowIterator, next func(i *test.NumRows) error) error {
	return iter.Do(func(r *spanner.Row) error {
		item, err := NumRowsFromMySpannerDatabaseRow(r)
		if err != nil {
			return fmt.Errorf("error converting test.NumRows row to protobuf message: %s", err)
		}
		return next(item)
	})
}
func (s *MySpannerImpl) UniaryInsert(ctx context.Context, req *test.ExampleTable) (*test.ExampleTable, error) {
	var err error
	var res = &test.ExampleTable{}
	_ = err
	_ = res
	params, err := ExampleTableToMySpannerPersistType(req)
	if err != nil {
		return nil, err
	}
	var iterErr error
	err = s.PERSIST.UniaryInsert(ctx, params, func(row *spanner.Row) {
		if row == nil { // there was no return data
			return
		}
		res, err = ExampleTableFromMySpannerDatabaseRow(row)
		if err != nil {
			iterErr = err
			return
		}
	})
	if err != nil {
		return nil, gstatus.Errorf(codes.Unknown, "error calling persist service: %v", err)
	} else if iterErr != nil {
		return nil, iterErr
	}
	return res, nil
}
func (s *MySpannerImpl) UniarySelect(ctx context.Context, req *test.ExampleTable) (*test.ExampleTable, error) {
	var err error
	var res = &test.ExampleTable{}
	_ = err
	_ = res
	params, err := ExampleTableToMySpannerPersistType(req)
	if err != nil {
		return nil, err
	}
	var iterErr error
	err = s.PERSIST.UniarySelect(ctx, params, func(row *spanner.Row) {
		if row == nil { // there was no return data
			return
		}
		res, err = ExampleTableFromMySpannerDatabaseRow(row)
		if err != nil {
			iterErr = err
			return
		}
	})
	if err != nil {
		return nil, gstatus.Errorf(codes.Unknown, "error calling persist service: %v", err)
	} else if iterErr != nil {
		return nil, iterErr
	}
	return res, nil
}
func (s *MySpannerImpl) TestNest(ctx context.Context, req *Something) (*Something, error) {
	var err error
	var res = &Something{}
	_ = err
	_ = res
	params, err := SomethingToMySpannerPersistType(req)
	if err != nil {
		return nil, err
	}
	var iterErr error
	err = s.PERSIST.TestNest(ctx, params, func(row *spanner.Row) {
		if row == nil { // there was no return data
			return
		}
		res, err = SomethingFromMySpannerDatabaseRow(row)
		if err != nil {
			iterErr = err
			return
		}
	})
	if err != nil {
		return nil, gstatus.Errorf(codes.Unknown, "error calling persist service: %v", err)
	} else if iterErr != nil {
		return nil, iterErr
	}
	return res, nil
}
func (s *MySpannerImpl) TestEverything(ctx context.Context, req *HasTimestamp) (*HasTimestamp, error) {
	var err error
	var res = &HasTimestamp{}
	_ = err
	_ = res
	params, err := HasTimestampToMySpannerPersistType(req)
	if err != nil {
		return nil, err
	}
	var iterErr error
	err = s.PERSIST.TestEverything(ctx, params, func(row *spanner.Row) {
		if row == nil { // there was no return data
			return
		}
		res, err = HasTimestampFromMySpannerDatabaseRow(row)
		if err != nil {
			iterErr = err
			return
		}
	})
	if err != nil {
		return nil, gstatus.Errorf(codes.Unknown, "error calling persist service: %v", err)
	} else if iterErr != nil {
		return nil, iterErr
	}
	return res, nil
}
func (s *MySpannerImpl) UniarySelectWithDirectives(ctx context.Context, req *test.ExampleTable) (*test.ExampleTable, error) {
	var err error
	var res = &test.ExampleTable{}
	_ = err
	_ = res
	params, err := ExampleTableToMySpannerPersistType(req)
	if err != nil {
		return nil, err
	}
	var iterErr error
	err = s.PERSIST.UniarySelectWithDirectives(ctx, params, func(row *spanner.Row) {
		if row == nil { // there was no return data
			return
		}
		res, err = ExampleTableFromMySpannerDatabaseRow(row)
		if err != nil {
			iterErr = err
			return
		}
	})
	if err != nil {
		return nil, gstatus.Errorf(codes.Unknown, "error calling persist service: %v", err)
	} else if iterErr != nil {
		return nil, iterErr
	}
	return res, nil
}
func (s *MySpannerImpl) UniaryUpdate(ctx context.Context, req *test.ExampleTable) (*test.PartialTable, error) {
	var err error
	var res = &test.PartialTable{}
	_ = err
	_ = res
	params, err := ExampleTableToMySpannerPersistType(req)
	if err != nil {
		return nil, err
	}
	var iterErr error
	err = s.PERSIST.UniaryUpdate(ctx, params, func(row *spanner.Row) {
		if row == nil { // there was no return data
			return
		}
		res, err = PartialTableFromMySpannerDatabaseRow(row)
		if err != nil {
			iterErr = err
			return
		}
	})
	if err != nil {
		return nil, gstatus.Errorf(codes.Unknown, "error calling persist service: %v", err)
	} else if iterErr != nil {
		return nil, iterErr
	}
	return res, nil
}
func (s *MySpannerImpl) UniaryDeleteRange(ctx context.Context, req *test.ExampleTableRange) (*test.ExampleTable, error) {
	var err error
	var res = &test.ExampleTable{}
	_ = err
	_ = res
	params, err := ExampleTableRangeToMySpannerPersistType(req)
	if err != nil {
		return nil, err
	}
	var iterErr error
	err = s.PERSIST.UniaryDeleteRange(ctx, params, func(row *spanner.Row) {
		if row == nil { // there was no return data
			return
		}
		res, err = ExampleTableFromMySpannerDatabaseRow(row)
		if err != nil {
			iterErr = err
			return
		}
	})
	if err != nil {
		return nil, gstatus.Errorf(codes.Unknown, "error calling persist service: %v", err)
	} else if iterErr != nil {
		return nil, iterErr
	}
	return res, nil
}
func (s *MySpannerImpl) UniaryDeleteSingle(ctx context.Context, req *test.ExampleTable) (*test.ExampleTable, error) {
	var err error
	var res = &test.ExampleTable{}
	_ = err
	_ = res
	params, err := ExampleTableToMySpannerPersistType(req)
	if err != nil {
		return nil, err
	}
	var iterErr error
	err = s.PERSIST.UniaryDeleteSingle(ctx, params, func(row *spanner.Row) {
		if row == nil { // there was no return data
			return
		}
		res, err = ExampleTableFromMySpannerDatabaseRow(row)
		if err != nil {
			iterErr = err
			return
		}
	})
	if err != nil {
		return nil, gstatus.Errorf(codes.Unknown, "error calling persist service: %v", err)
	} else if iterErr != nil {
		return nil, iterErr
	}
	return res, nil
}
func (s *MySpannerImpl) NoArgs(ctx context.Context, req *test.ExampleTable) (*test.ExampleTable, error) {
	var err error
	var res = &test.ExampleTable{}
	_ = err
	_ = res
	params, err := ExampleTableToMySpannerPersistType(req)
	if err != nil {
		return nil, err
	}
	var iterErr error
	err = s.PERSIST.NoArgs(ctx, params, func(row *spanner.Row) {
		if row == nil { // there was no return data
			return
		}
		res, err = ExampleTableFromMySpannerDatabaseRow(row)
		if err != nil {
			iterErr = err
			return
		}
	})
	if err != nil {
		return nil, gstatus.Errorf(codes.Unknown, "error calling persist service: %v", err)
	} else if iterErr != nil {
		return nil, iterErr
	}
	return res, nil
}
func (s *MySpannerImpl) ServerStream(req *test.Name, stream MySpanner_ServerStreamServer) error {
	var err error
	_ = err
	params, err := NameToMySpannerPersistType(req)
	if err != nil {
		return err
	}
	var iterErr error
	err = s.PERSIST.ServerStream(stream.Context(), params, func(row *spanner.Row) {
		if row == nil { // there was no return data
			return
		}
		res, err := ExampleTableFromMySpannerDatabaseRow(row)
		if err != nil {
			iterErr = err
			return
		}
		if err := stream.Send(res); err != nil {
			iterErr = gstatus.Errorf(codes.Unknown, "error during iteration: %v", err)
		}
	})
	if err != nil {
		return gstatus.Errorf(codes.Unknown, "error during iteration: %v", err)
	} else if iterErr != nil {
		return iterErr
	}
	return nil
}
func (s *MySpannerImpl) ClientStreamInsert(stream MySpanner_ClientStreamInsertServer) error {
	var err error
	_ = err
	res := &test.NumRows{}
	feed, stop := s.PERSIST.ClientStreamInsert(stream.Context())
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			return gstatus.Errorf(codes.Unknown, "error receiving request: %v", err)
		}
		params, err := ExampleTableToMySpannerPersistType(req)
		if err != nil {
			return err
		}
		feed(params)
	}
	row, err := stop()
	if err != nil {
		return gstatus.Errorf(codes.Unknown, "error receiving result row: %v", err)
	}
	if row != nil {
		res, err = NumRowsFromMySpannerDatabaseRow(row)
		if err != nil {
			return err
		}
	}
	if err := stream.SendAndClose(res); err != nil {
		return gstatus.Errorf(codes.Unknown, "error sending back response: %v", err)
	}
	return nil
}
func (s *MySpannerImpl) ClientStreamDelete(stream MySpanner_ClientStreamDeleteServer) error {
	var err error
	_ = err
	res := &test.NumRows{}
	feed, stop := s.PERSIST.ClientStreamDelete(stream.Context())
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			return gstatus.Errorf(codes.Unknown, "error receiving request: %v", err)
		}
		params, err := ExampleTableToMySpannerPersistType(req)
		if err != nil {
			return err
		}
		feed(params)
	}
	row, err := stop()
	if err != nil {
		return gstatus.Errorf(codes.Unknown, "error receiving result row: %v", err)
	}
	if row != nil {
		res, err = NumRowsFromMySpannerDatabaseRow(row)
		if err != nil {
			return err
		}
	}
	if err := stream.SendAndClose(res); err != nil {
		return gstatus.Errorf(codes.Unknown, "error sending back response: %v", err)
	}
	return nil
}
func (s *MySpannerImpl) ClientStreamUpdate(stream MySpanner_ClientStreamUpdateServer) error {
	var err error
	_ = err
	res := &test.NumRows{}
	feed, stop := s.PERSIST.ClientStreamUpdate(stream.Context())
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			return gstatus.Errorf(codes.Unknown, "error receiving request: %v", err)
		}
		params, err := ExampleTableToMySpannerPersistType(req)
		if err != nil {
			return err
		}
		feed(params)
	}
	row, err := stop()
	if err != nil {
		return gstatus.Errorf(codes.Unknown, "error receiving result row: %v", err)
	}
	if row != nil {
		res, err = NumRowsFromMySpannerDatabaseRow(row)
		if err != nil {
			return err
		}
	}
	if err := stream.SendAndClose(res); err != nil {
		return gstatus.Errorf(codes.Unknown, "error sending back response: %v", err)
	}
	return nil
}
func (s *MySpannerImpl) UniaryInsertWithHooks(ctx context.Context, req *test.ExampleTable) (*test.ExampleTable, error) {
	var err error
	var res = &test.ExampleTable{}
	_ = err
	_ = res
	beforeRes, err := hooks.UniaryInsertBeforeHook(req)
	if err != nil {
		return nil, gstatus.Errorf(codes.Unknown, "error in before hook: %v", err)
	} else if beforeRes != nil {
		return beforeRes, nil
	}
	params, err := ExampleTableToMySpannerPersistType(req)
	if err != nil {
		return nil, err
	}
	var iterErr error
	err = s.PERSIST.UniaryInsertWithHooks(ctx, params, func(row *spanner.Row) {
		if row == nil { // there was no return data
			return
		}
		res, err = ExampleTableFromMySpannerDatabaseRow(row)
		if err != nil {
			iterErr = err
			return
		}
	})
	if err != nil {
		return nil, gstatus.Errorf(codes.Unknown, "error calling persist service: %v", err)
	} else if iterErr != nil {
		return nil, iterErr
	}
	if err := hooks.UniaryInsertAfterHook(req, res); err != nil {
		return nil, gstatus.Errorf(codes.Unknown, "error in after hook: %v", err)
	}
	return res, nil
}
func (s *MySpannerImpl) UniarySelectWithHooks(ctx context.Context, req *test.ExampleTable) (*test.ExampleTable, error) {
	var err error
	var res = &test.ExampleTable{}
	_ = err
	_ = res
	beforeRes, err := hooks.UniaryInsertBeforeHook(req)
	if err != nil {
		return nil, gstatus.Errorf(codes.Unknown, "error in before hook: %v", err)
	} else if beforeRes != nil {
		return beforeRes, nil
	}
	params, err := ExampleTableToMySpannerPersistType(req)
	if err != nil {
		return nil, err
	}
	var iterErr error
	err = s.PERSIST.UniarySelectWithHooks(ctx, params, func(row *spanner.Row) {
		if row == nil { // there was no return data
			return
		}
		res, err = ExampleTableFromMySpannerDatabaseRow(row)
		if err != nil {
			iterErr = err
			return
		}
	})
	if err != nil {
		return nil, gstatus.Errorf(codes.Unknown, "error calling persist service: %v", err)
	} else if iterErr != nil {
		return nil, iterErr
	}
	if err := hooks.UniaryInsertAfterHook(req, res); err != nil {
		return nil, gstatus.Errorf(codes.Unknown, "error in after hook: %v", err)
	}
	return res, nil
}
func (s *MySpannerImpl) UniaryUpdateWithHooks(ctx context.Context, req *test.ExampleTable) (*test.PartialTable, error) {
	var err error
	var res = &test.PartialTable{}
	_ = err
	_ = res
	beforeRes, err := hooks.UniaryUpdateBeforeHook(req)
	if err != nil {
		return nil, gstatus.Errorf(codes.Unknown, "error in before hook: %v", err)
	} else if beforeRes != nil {
		return beforeRes, nil
	}
	params, err := ExampleTableToMySpannerPersistType(req)
	if err != nil {
		return nil, err
	}
	var iterErr error
	err = s.PERSIST.UniaryUpdateWithHooks(ctx, params, func(row *spanner.Row) {
		if row == nil { // there was no return data
			return
		}
		res, err = PartialTableFromMySpannerDatabaseRow(row)
		if err != nil {
			iterErr = err
			return
		}
	})
	if err != nil {
		return nil, gstatus.Errorf(codes.Unknown, "error calling persist service: %v", err)
	} else if iterErr != nil {
		return nil, iterErr
	}
	if err := hooks.UniaryUpdateAfterHook(req, res); err != nil {
		return nil, gstatus.Errorf(codes.Unknown, "error in after hook: %v", err)
	}
	return res, nil
}
func (s *MySpannerImpl) UniaryDeleteWithHooks(ctx context.Context, req *test.ExampleTableRange) (*test.ExampleTable, error) {
	var err error
	var res = &test.ExampleTable{}
	_ = err
	_ = res
	beforeRes, err := hooks.UniaryDeleteBeforeHook(req)
	if err != nil {
		return nil, gstatus.Errorf(codes.Unknown, "error in before hook: %v", err)
	} else if beforeRes != nil {
		return beforeRes, nil
	}
	params, err := ExampleTableRangeToMySpannerPersistType(req)
	if err != nil {
		return nil, err
	}
	var iterErr error
	err = s.PERSIST.UniaryDeleteWithHooks(ctx, params, func(row *spanner.Row) {
		if row == nil { // there was no return data
			return
		}
		res, err = ExampleTableFromMySpannerDatabaseRow(row)
		if err != nil {
			iterErr = err
			return
		}
	})
	if err != nil {
		return nil, gstatus.Errorf(codes.Unknown, "error calling persist service: %v", err)
	} else if iterErr != nil {
		return nil, iterErr
	}
	if err := hooks.UniaryDeleteAfterHook(req, res); err != nil {
		return nil, gstatus.Errorf(codes.Unknown, "error in after hook: %v", err)
	}
	return res, nil
}
func (s *MySpannerImpl) ServerStreamWithHooks(req *test.Name, stream MySpanner_ServerStreamWithHooksServer) error {
	var err error
	_ = err
	beforeRes, err := hooks.ServerStreamBeforeHook(req)
	if err != nil {
		return gstatus.Errorf(codes.Unknown, "error in before hook: %v", err)
	} else if beforeRes != nil {
		for _, res := range beforeRes {
			if err := stream.Send(res); err != nil {
				return gstatus.Errorf(codes.Unknown, "error sending back before hook result: %v", err)
			}
		}
	}
	params, err := NameToMySpannerPersistType(req)
	if err != nil {
		return err
	}
	var iterErr error
	err = s.PERSIST.ServerStreamWithHooks(stream.Context(), params, func(row *spanner.Row) {
		if row == nil { // there was no return data
			return
		}
		res, err := ExampleTableFromMySpannerDatabaseRow(row)
		if err != nil {
			iterErr = err
			return
		}
		if err := hooks.ServerStreamAfterHook(req, res); err != nil {
			iterErr = gstatus.Errorf(codes.Unknown, "error in after hook: %v", err)
			return
		}
		if err := stream.Send(res); err != nil {
			iterErr = gstatus.Errorf(codes.Unknown, "error during iteration: %v", err)
		}
	})
	if err != nil {
		return gstatus.Errorf(codes.Unknown, "error during iteration: %v", err)
	} else if iterErr != nil {
		return iterErr
	}
	return nil
}
func (s *MySpannerImpl) ClientStreamUpdateWithHooks(stream MySpanner_ClientStreamUpdateWithHooksServer) error {
	var err error
	_ = err
	res := &test.NumRows{}
	feed, stop := s.PERSIST.ClientStreamUpdateWithHooks(stream.Context())
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			return gstatus.Errorf(codes.Unknown, "error receiving request: %v", err)
		}
		beforeRes, err := hooks.ClientStreamUpdateBeforeHook(req)
		if err != nil {
			return gstatus.Errorf(codes.Unknown, "error in before hook: %v", err)
		} else if beforeRes != nil {
			continue
		}
		params, err := ExampleTableToMySpannerPersistType(req)
		if err != nil {
			return err
		}
		feed(params)
	}
	row, err := stop()
	if err != nil {
		return gstatus.Errorf(codes.Unknown, "error receiving result row: %v", err)
	}
	if row != nil {
		res, err = NumRowsFromMySpannerDatabaseRow(row)
		if err != nil {
			return err
		}
	}
	// NOTE: I dont want to store your requests in memory
	// so the after hook for client streaming calls
	// is called with an empty request struct
	fakeReq := &test.ExampleTable{}
	if err := hooks.ClientStreamUpdateAfterHook(fakeReq, res); err != nil {
		return gstatus.Errorf(codes.Unknown, "error in after hook: %v", err)
	}
	if err := stream.SendAndClose(res); err != nil {
		return gstatus.Errorf(codes.Unknown, "error sending back response: %v", err)
	}
	return nil
}
