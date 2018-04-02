// This file is generated by protoc-gen-persist
// Source File: tests/sql/little_of_everything/little_of_everything.proto
// DO NOT EDIT !
package little_of_everything

import (
	sql "database/sql"
	fmt "fmt"
	io "io"

	proto "github.com/golang/protobuf/proto"
	mytime "github.com/tcncloud/protoc-gen-persist/tests/mytime"
	persist_lib "github.com/tcncloud/protoc-gen-persist/tests/sql/little_of_everything/persist_lib"
	test "github.com/tcncloud/protoc-gen-persist/tests/test"
	context "golang.org/x/net/context"
	codes "google.golang.org/grpc/codes"
	gstatus "google.golang.org/grpc/status"
)

type Testservice1Impl struct {
	PERSIST   *persist_lib.Testservice1MethodReceiver
	FORWARDED RestOfTestservice1Handlers
}
type RestOfTestservice1Handlers interface {
}
type Testservice1ImplBuilder struct {
	err           error
	rest          RestOfTestservice1Handlers
	queryHandlers *persist_lib.Testservice1QueryHandlers
	i             *Testservice1Impl
	db            sql.DB
}

func NewTestservice1Builder() *Testservice1ImplBuilder {
	return &Testservice1ImplBuilder{i: &Testservice1Impl{}}
}
func (b *Testservice1ImplBuilder) WithRestOfGrpcHandlers(r RestOfTestservice1Handlers) *Testservice1ImplBuilder {
	b.rest = r
	return b
}
func (b *Testservice1ImplBuilder) WithPersistQueryHandlers(p *persist_lib.Testservice1QueryHandlers) *Testservice1ImplBuilder {
	b.queryHandlers = p
	return b
}
func (b *Testservice1ImplBuilder) WithDefaultQueryHandlers() *Testservice1ImplBuilder {
	accessor := persist_lib.NewSqlClientGetter(&b.db)
	queryHandlers := &persist_lib.Testservice1QueryHandlers{
		UnaryExample1Handler:          persist_lib.DefaultUnaryExample1Handler(accessor),
		UnaryExample2Handler:          persist_lib.DefaultUnaryExample2Handler(accessor),
		ServerStreamSelectHandler:     persist_lib.DefaultServerStreamSelectHandler(accessor),
		ClientStreamingExampleHandler: persist_lib.DefaultClientStreamingExampleHandler(accessor),
	}
	b.queryHandlers = queryHandlers
	return b
}

// set the custom handlers you want to use in the handlers
// this method will make sure to use a default handler if
// the handler is nil.
func (b *Testservice1ImplBuilder) WithNilAsDefaultQueryHandlers(p *persist_lib.Testservice1QueryHandlers) *Testservice1ImplBuilder {
	accessor := persist_lib.NewSqlClientGetter(&b.db)
	if p.UnaryExample1Handler == nil {
		p.UnaryExample1Handler = persist_lib.DefaultUnaryExample1Handler(accessor)
	}
	if p.UnaryExample2Handler == nil {
		p.UnaryExample2Handler = persist_lib.DefaultUnaryExample2Handler(accessor)
	}
	if p.ServerStreamSelectHandler == nil {
		p.ServerStreamSelectHandler = persist_lib.DefaultServerStreamSelectHandler(accessor)
	}
	if p.ClientStreamingExampleHandler == nil {
		p.ClientStreamingExampleHandler = persist_lib.DefaultClientStreamingExampleHandler(accessor)
	}
	b.queryHandlers = p
	return b
}
func (b *Testservice1ImplBuilder) WithSqlClient(c *sql.DB) *Testservice1ImplBuilder {
	b.db = *c
	return b
}
func (b *Testservice1ImplBuilder) WithNewSqlDb(driverName, dataSourceName string) *Testservice1ImplBuilder {
	db, err := sql.Open(driverName, dataSourceName)
	b.err = err
	if b.err == nil {
		b.db = *db
	}
	return b
}
func (b *Testservice1ImplBuilder) Build() (*Testservice1Impl, error) {
	if b.err != nil {
		return nil, b.err
	}
	b.i.PERSIST = &persist_lib.Testservice1MethodReceiver{Handlers: *b.queryHandlers}
	b.i.FORWARDED = b.rest
	return b.i, nil
}
func (b *Testservice1ImplBuilder) MustBuild() *Testservice1Impl {
	s, err := b.Build()
	if err != nil {
		panic("error in builder: " + err.Error())
	}
	return s
}
func ExampleTable1ToTestservice1PersistType(req *ExampleTable1) (*persist_lib.ExampleTable1ForTestservice1, error) {
	params := &persist_lib.ExampleTable1ForTestservice1{}
	params.TableId = req.TableId
	params.Key = req.Key
	params.Value = req.Value
	if req.InnerMessage == nil {
		req.InnerMessage = new(ExampleTable1_InnerMessage)
	}
	{
		raw, err := proto.Marshal(req.InnerMessage)
		if err != nil {
			return nil, err
		}
		params.InnerMessage = raw
	}
	params.InnerEnum = int32(req.InnerEnum)
	params.StringArray = req.StringArray
	params.BytesField = req.BytesField
	params.StartTime = (mytime.MyTime{}).ToSql(req.StartTime)
	if req.TestField == nil {
		req.TestField = new(test.Test)
	}
	{
		raw, err := proto.Marshal(req.TestField)
		if err != nil {
			return nil, err
		}
		params.TestField = raw
	}
	params.Myyenum = int32(req.Myyenum)
	params.Testsenum = int32(req.Testsenum)
	params.Mappedenum = (MyMappedEnum{}).ToSql(req.Mappedenum)
	return params, nil
}
func ExampleTable1FromTestservice1DatabaseRow(row persist_lib.Scanable) (*ExampleTable1, error) {
	res := &ExampleTable1{}
	var TableId_ int32
	var Key_ string
	var Value_ string
	var InnerMessage_ []byte
	var InnerEnum_ int32
	var StringArray_ []string
	var BytesField_ []byte
	var StartTime_ mytime.MyTime
	var TestField_ []byte
	var Myyenum_ int32
	var Testsenum_ int32
	var Mappedenum_ MyMappedEnum
	if err := row.Scan(
		&TableId_,
		&Key_,
		&Value_,
		&InnerMessage_,
		&InnerEnum_,
		&StringArray_,
		&BytesField_,
		&StartTime_,
		&TestField_,
		&Myyenum_,
		&Testsenum_,
		&Mappedenum_,
	); err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	res.TableId = TableId_
	res.Key = Key_
	res.Value = Value_
	{
		var converted = new(ExampleTable1_InnerMessage)
		if err := proto.Unmarshal(InnerMessage_, converted); err != nil {
			return nil, err
		}
		res.InnerMessage = converted
	}
	res.InnerEnum = ExampleTable1_InnerEnum(InnerEnum_)
	res.StringArray = StringArray_
	res.BytesField = BytesField_
	res.StartTime = StartTime_.ToProto()
	{
		var converted = new(test.Test)
		if err := proto.Unmarshal(TestField_, converted); err != nil {
			return nil, err
		}
		res.TestField = converted
	}
	res.Myyenum = MyEnum(Myyenum_)
	res.Testsenum = test.TestEnum(Testsenum_)
	res.Mappedenum = Mappedenum_.ToProto()
	return res, nil
}
func IterTestservice1ExampleTable1Proto(iter *persist_lib.Result, next func(i *ExampleTable1) error) error {
	return iter.Do(func(r persist_lib.Scanable) error {
		item, err := ExampleTable1FromTestservice1DatabaseRow(r)
		if err != nil {
			return fmt.Errorf("error converting ExampleTable1 row to protobuf message: %s", err)
		}
		return next(item)
	})
}
func TestToTestservice1PersistType(req *test.Test) (*persist_lib.Test_TestForTestservice1, error) {
	params := &persist_lib.Test_TestForTestservice1{}
	params.Id = req.Id
	params.Name = req.Name
	return params, nil
}
func CountRowsFromTestservice1DatabaseRow(row persist_lib.Scanable) (*CountRows, error) {
	res := &CountRows{}
	var Count_ int64
	if err := row.Scan(
		&Count_,
	); err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	res.Count = Count_
	return res, nil
}
func IterTestservice1CountRowsProto(iter *persist_lib.Result, next func(i *CountRows) error) error {
	return iter.Do(func(r persist_lib.Scanable) error {
		item, err := CountRowsFromTestservice1DatabaseRow(r)
		if err != nil {
			return fmt.Errorf("error converting CountRows row to protobuf message: %s", err)
		}
		return next(item)
	})
}
func (s *Testservice1Impl) UnaryExample1(ctx context.Context, req *ExampleTable1) (*ExampleTable1, error) {
	var err error
	var res = &ExampleTable1{}
	_ = err
	_ = res
	params, err := ExampleTable1ToTestservice1PersistType(req)
	if err != nil {
		return nil, err
	}
	var iterErr error
	err = s.PERSIST.UnaryExample1(ctx, params, func(row persist_lib.Scanable) {
		if row == nil { // there was no return data
			return
		}
		res, err = ExampleTable1FromTestservice1DatabaseRow(row)
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
func (s *Testservice1Impl) UnaryExample2(ctx context.Context, req *test.Test) (*ExampleTable1, error) {
	var err error
	var res = &ExampleTable1{}
	_ = err
	_ = res
	params, err := TestToTestservice1PersistType(req)
	if err != nil {
		return nil, err
	}
	var iterErr error
	err = s.PERSIST.UnaryExample2(ctx, params, func(row persist_lib.Scanable) {
		if row == nil { // there was no return data
			return
		}
		res, err = ExampleTable1FromTestservice1DatabaseRow(row)
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
func (s *Testservice1Impl) ServerStreamSelect(req *ExampleTable1, stream Testservice1_ServerStreamSelectServer) error {
	var err error
	_ = err
	params, err := ExampleTable1ToTestservice1PersistType(req)
	if err != nil {
		return err
	}
	var iterErr error
	err = s.PERSIST.ServerStreamSelect(stream.Context(), params, func(row persist_lib.Scanable) {
		if row == nil { // there was no return data
			return
		}
		res, err := ExampleTable1FromTestservice1DatabaseRow(row)
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
func (s *Testservice1Impl) ClientStreamingExample(stream Testservice1_ClientStreamingExampleServer) error {
	var err error
	_ = err
	res := &CountRows{}
	feed, stop, err := s.PERSIST.ClientStreamingExample(stream.Context())
	if err != nil {
		return err
	}
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			return gstatus.Errorf(codes.Unknown, "error receiving request: %v", err)
		}
		params, err := ExampleTable1ToTestservice1PersistType(req)
		if err != nil {
			return err
		}
		if err := feed(params); err != nil {
			return err
		}
	}
	row, err := stop()
	if err != nil {
		return gstatus.Errorf(codes.Unknown, "error receiving result row: %v", err)
	}
	if row != nil {
		res, err = CountRowsFromTestservice1DatabaseRow(row)
		if err != nil {
			return err
		}
	}
	if err := stream.SendAndClose(res); err != nil {
		return gstatus.Errorf(codes.Unknown, "error sending back response: %v", err)
	}
	return nil
}
