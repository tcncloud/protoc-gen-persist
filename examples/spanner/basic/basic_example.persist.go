// This file is generated by protoc-gen-persist
// Source File: examples/spanner/basic/basic_example.proto
// DO NOT EDIT !
package basic

import (
	fmt "fmt"
	io "io"
	strings "strings"

	"cloud.google.com/go/spanner"
	mytime "github.com/tcncloud/protoc-gen-persist/examples/mytime"
	"github.com/tcncloud/protoc-gen-persist/examples/spanner/hooks"
	test "github.com/tcncloud/protoc-gen-persist/examples/test"
	context "golang.org/x/net/context"
	iterator "google.golang.org/api/iterator"
	"google.golang.org/api/option"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
)

type MySpannerImpl struct {
	SpannerDB *spanner.Client
}

func NewMySpannerImpl(d string, conf *spanner.ClientConfig, opts ...option.ClientOption) (*MySpannerImpl, error) {
	var client *spanner.Client
	var err error
	if conf != nil {
		client, err = spanner.NewClientWithConfig(context.Background(), d, *conf, opts...)
	} else {
		client, err = spanner.NewClient(context.Background(), d, opts...)
	}
	if err != nil {
		return nil, err
	}
	return &MySpannerImpl{SpannerDB: client}, nil
}

// spanner unary select UniaryInsert
func (s *MySpannerImpl) UniaryInsert(ctx context.Context, req *test.ExampleTable) (*test.ExampleTable, error) {
	var err error

	params := make([]interface{}, 0)
	var conv interface{}

	conv = req.Id

	if err != nil {
		return nil, grpc.Errorf(codes.Unknown, err.Error())
	}
	params = append(params, conv)

	conv, err = mytime.MyTime{}.ToSpanner(req.StartTime).SpannerValue()

	if err != nil {
		return nil, grpc.Errorf(codes.Unknown, err.Error())
	}
	params = append(params, conv)
	params = append(params, "bananas")
	muts := make([]*spanner.Mutation, 1)
	muts[0] = spanner.Insert("example_table", []string{"id", "start_time", "name"}, params)
	_, err = s.SpannerDB.Apply(ctx, muts)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			return nil, grpc.Errorf(codes.AlreadyExists, err.Error())
		} else {
			return nil, grpc.Errorf(codes.Unknown, err.Error())
		}
	}
	res := test.ExampleTable{}

	return &res, nil
}

// spanner unary select UniarySelect
func (s *MySpannerImpl) UniarySelect(ctx context.Context, req *test.ExampleTable) (*test.ExampleTable, error) {
	var err error
	var (
		Id        int64
		Name      string
		StartTime mytime.MyTime
	)

	params := make(map[string]interface{})
	var conv interface{}

	conv = req.Id

	if err != nil {
		return nil, grpc.Errorf(codes.Unknown, err.Error())
	}
	params["string0"] = conv

	conv = req.Name

	if err != nil {
		return nil, grpc.Errorf(codes.Unknown, err.Error())
	}
	params["string1"] = conv
	//stmt := spanner.Statement{SQL: "{ {.Spanner.Query} }", Params: params}
	stmt := spanner.Statement{SQL: "SELECT * from example_table Where id= @string0 AND name= @string1", Params: params}
	tx := s.SpannerDB.Single()
	defer tx.Close()
	iter := tx.Query(ctx, stmt)
	defer iter.Stop()
	row, err := iter.Next()
	if err == iterator.Done {
		return nil, grpc.Errorf(codes.NotFound, "no rows found")
	} else if err != nil {
		return nil, grpc.Errorf(codes.Unknown, err.Error())
	}
	// scan our values out of the row

	err = row.ColumnByName("id", &Id)
	if err != nil {
		return nil, grpc.Errorf(codes.Unknown, err.Error())
	}

	gcv := new(spanner.GenericColumnValue)
	err = row.ColumnByName("start_time", gcv)
	if err != nil {
		return nil, grpc.Errorf(codes.Unknown, err.Error())
	}
	err = StartTime.SpannerScan(gcv)
	if err != nil {
		return nil, grpc.Errorf(codes.Unknown, err.Error())
	}

	err = row.ColumnByName("name", &Name)
	if err != nil {
		return nil, grpc.Errorf(codes.Unknown, err.Error())
	}

	_, err = iter.Next()
	if err != iterator.Done {
		fmt.Println("Unary select that returns more than one row..")
	}
	res := test.ExampleTable{

		Id:        Id,
		Name:      Name,
		StartTime: StartTime.ToProto(),
	}

	return &res, nil
}

// spanner unary select UniaryUpdate
func (s *MySpannerImpl) UniaryUpdate(ctx context.Context, req *test.ExampleTable) (*test.PartialTable, error) {
	var err error

	params := make(map[string]interface{})
	var conv interface{}

	conv, err = mytime.MyTime{}.ToSpanner(req.StartTime).SpannerValue()

	if err != nil {
		return nil, grpc.Errorf(codes.Unknown, err.Error())
	}
	params["start_time"] = conv
	conv = "oranges"
	params["name"] = conv

	conv = req.Id

	if err != nil {
		return nil, grpc.Errorf(codes.Unknown, err.Error())
	}
	params["id"] = conv
	muts := make([]*spanner.Mutation, 1)
	muts[0] = spanner.UpdateMap("example_table", params)
	_, err = s.SpannerDB.Apply(ctx, muts)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			return nil, grpc.Errorf(codes.AlreadyExists, err.Error())
		} else {
			return nil, grpc.Errorf(codes.Unknown, err.Error())
		}
	}
	res := test.PartialTable{}

	return &res, nil
}

// spanner unary select UniaryDelete
func (s *MySpannerImpl) UniaryDelete(ctx context.Context, req *test.ExampleTableRange) (*test.ExampleTable, error) {
	var err error

	start := make([]interface{}, 0)
	end := make([]interface{}, 0)
	var conv interface{}
	conv = req.StartId
	if err != nil {
		return nil, grpc.Errorf(codes.Unknown, err.Error())
	}
	start = append(start, conv)
	conv = req.EndId
	if err != nil {
		return nil, grpc.Errorf(codes.Unknown, err.Error())
	}
	end = append(end, conv)
	key := spanner.KeyRange{
		Start: start,
		End:   end,
		Kind:  spanner.ClosedOpen,
	}
	muts := make([]*spanner.Mutation, 1)
	muts[0] = spanner.DeleteKeyRange("example_table", key)
	_, err = s.SpannerDB.Apply(ctx, muts)
	if err != nil {
		if strings.Contains(err.Error(), "does not exist") {
			return nil, grpc.Errorf(codes.NotFound, err.Error())
		}
	}
	res := test.ExampleTable{}

	return &res, nil
}

// spanner unary select NoArgs
func (s *MySpannerImpl) NoArgs(ctx context.Context, req *test.ExampleTable) (*test.ExampleTable, error) {
	var err error
	var (
		Id        int64
		Name      string
		StartTime mytime.MyTime
	)

	params := make(map[string]interface{})
	//stmt := spanner.Statement{SQL: "{ {.Spanner.Query} }", Params: params}
	stmt := spanner.Statement{SQL: "select * from example_table limit 1", Params: params}
	tx := s.SpannerDB.Single()
	defer tx.Close()
	iter := tx.Query(ctx, stmt)
	defer iter.Stop()
	row, err := iter.Next()
	if err == iterator.Done {
		return nil, grpc.Errorf(codes.NotFound, "no rows found")
	} else if err != nil {
		return nil, grpc.Errorf(codes.Unknown, err.Error())
	}
	// scan our values out of the row

	err = row.ColumnByName("id", &Id)
	if err != nil {
		return nil, grpc.Errorf(codes.Unknown, err.Error())
	}

	gcv := new(spanner.GenericColumnValue)
	err = row.ColumnByName("start_time", gcv)
	if err != nil {
		return nil, grpc.Errorf(codes.Unknown, err.Error())
	}
	err = StartTime.SpannerScan(gcv)
	if err != nil {
		return nil, grpc.Errorf(codes.Unknown, err.Error())
	}

	err = row.ColumnByName("name", &Name)
	if err != nil {
		return nil, grpc.Errorf(codes.Unknown, err.Error())
	}

	_, err = iter.Next()
	if err != iterator.Done {
		fmt.Println("Unary select that returns more than one row..")
	}
	res := test.ExampleTable{

		Id:        Id,
		Name:      Name,
		StartTime: StartTime.ToProto(),
	}

	return &res, nil
}

// spanner server streaming ServerStream
func (s *MySpannerImpl) ServerStream(req *test.Name, stream MySpanner_ServerStreamServer) error {
	var (
		Id        int64
		Name      string
		StartTime mytime.MyTime
	)

	params := make(map[string]interface{})
	stmt := spanner.Statement{SQL: "SELECT * FROM example_table", Params: params}
	tx := s.SpannerDB.Single()
	defer tx.Close()
	iter := tx.Query(context.Background(), stmt)
	defer iter.Stop()
	for {
		row, err := iter.Next()
		if err == iterator.Done {
			break
		} else if err != nil {
			return grpc.Errorf(codes.Unknown, err.Error())
		}
		// scan our values out of the row

		err = row.ColumnByName("id", &Id)
		if err != nil {
			return grpc.Errorf(codes.Unknown, err.Error())
		}

		gcv := new(spanner.GenericColumnValue)
		err = row.ColumnByName("start_time", gcv)
		if err != nil {
			return grpc.Errorf(codes.Unknown, err.Error())
		}
		err = StartTime.SpannerScan(gcv)
		if err != nil {
			return grpc.Errorf(codes.Unknown, err.Error())
		}

		err = row.ColumnByName("name", &Name)
		if err != nil {
			return grpc.Errorf(codes.Unknown, err.Error())
		}

		res := test.ExampleTable{

			Id:        Id,
			Name:      Name,
			StartTime: StartTime.ToProto(),
		}

		stream.Send(&res)
	}
	return nil
}

// spanner client streaming ClientStreamInsert
func (s *MySpannerImpl) ClientStreamInsert(stream MySpanner_ClientStreamInsertServer) error {
	var err error
	res := test.NumRows{}

	muts := make([]*spanner.Mutation, 0)
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			return grpc.Errorf(codes.Unknown, err.Error())
		}

		//spanner client streaming insert
		params := make([]interface{}, 0)
		var conv interface{}

		conv = req.Id

		if err != nil {
			return grpc.Errorf(codes.Unknown, err.Error())
		}
		params = append(params, conv)

		conv, err = mytime.MyTime{}.ToSpanner(req.StartTime).SpannerValue()

		if err != nil {
			return grpc.Errorf(codes.Unknown, err.Error())
		}
		params = append(params, conv)

		conv = req.Name

		if err != nil {
			return grpc.Errorf(codes.Unknown, err.Error())
		}
		params = append(params, conv)
		muts = append(muts, spanner.Insert("example_table", []string{"id", "start_time", "name"}, params))

		////////////////////////////// NOTE //////////////////////////////////////
		// In the future, we might do apply if muts gets really big,  but for now,
		// we only do one apply on the database with all the records stored in muts
		//////////////////////////////////////////////////////////////////////////
	}
	_, err = s.SpannerDB.Apply(context.Background(), muts)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			return grpc.Errorf(codes.AlreadyExists, err.Error())
		} else {
			return grpc.Errorf(codes.Unknown, err.Error())
		}
	}

	stream.SendAndClose(&res)
	return nil
}

// spanner client streaming ClientStreamDelete
func (s *MySpannerImpl) ClientStreamDelete(stream MySpanner_ClientStreamDeleteServer) error {
	var err error
	res := test.NumRows{}

	muts := make([]*spanner.Mutation, 0)
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			return grpc.Errorf(codes.Unknown, err.Error())
		}

		//spanner client streaming delete
		start := make([]interface{}, 0)
		end := make([]interface{}, 0)
		var conv interface{}
		conv = req.Id
		if err != nil {
			return grpc.Errorf(codes.Unknown, err.Error())
		}
		start = append(start, conv)
		conv = req.Id
		if err != nil {
			return grpc.Errorf(codes.Unknown, err.Error())
		}
		end = append(end, conv)
		key := spanner.KeyRange{
			Start: start,
			End:   end,
			Kind:  spanner.ClosedClosed,
		}
		muts = append(muts, spanner.DeleteKeyRange("example_table", key))
		////////////////////////////// NOTE //////////////////////////////////////
		// In the future, we might do apply if muts gets really big,  but for now,
		// we only do one apply on the database with all the records stored in muts
		//////////////////////////////////////////////////////////////////////////
	}
	_, err = s.SpannerDB.Apply(context.Background(), muts)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			return grpc.Errorf(codes.AlreadyExists, err.Error())
		} else {
			return grpc.Errorf(codes.Unknown, err.Error())
		}
	}

	stream.SendAndClose(&res)
	return nil
}

// spanner client streaming ClientStreamUpdate
func (s *MySpannerImpl) ClientStreamUpdate(stream MySpanner_ClientStreamUpdateServer) error {
	var err error
	res := test.NumRows{}

	muts := make([]*spanner.Mutation, 0)
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			return grpc.Errorf(codes.Unknown, err.Error())
		}

		//spanner client streaming update
		params := make(map[string]interface{})
		var conv interface{}

		conv, err = mytime.MyTime{}.ToSpanner(req.StartTime).SpannerValue()

		if err != nil {
			return grpc.Errorf(codes.Unknown, err.Error())
		}
		params["start_time"] = conv

		conv = req.Name

		if err != nil {
			return grpc.Errorf(codes.Unknown, err.Error())
		}
		params["name"] = conv

		conv = req.Id

		if err != nil {
			return grpc.Errorf(codes.Unknown, err.Error())
		}
		params["id"] = conv
		muts = append(muts, spanner.UpdateMap("example_table", params))

		////////////////////////////// NOTE //////////////////////////////////////
		// In the future, we might do apply if muts gets really big,  but for now,
		// we only do one apply on the database with all the records stored in muts
		//////////////////////////////////////////////////////////////////////////
	}
	_, err = s.SpannerDB.Apply(context.Background(), muts)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			return grpc.Errorf(codes.AlreadyExists, err.Error())
		} else {
			return grpc.Errorf(codes.Unknown, err.Error())
		}
	}

	stream.SendAndClose(&res)
	return nil
}

// spanner unary select UniaryInsertWithHooks
func (s *MySpannerImpl) UniaryInsertWithHooks(ctx context.Context, req *test.ExampleTable) (*test.ExampleTable, error) {
	var err error

	beforeRes, err := hooks.UniaryInsertBeforeHook(req)
	if err != nil {

		return nil, grpc.Errorf(codes.Unknown, err.Error())

	}
	if beforeRes != nil {

		return beforeRes, nil

	}

	params := make([]interface{}, 0)
	var conv interface{}

	conv = req.Id

	if err != nil {
		return nil, grpc.Errorf(codes.Unknown, err.Error())
	}
	params = append(params, conv)

	conv, err = mytime.MyTime{}.ToSpanner(req.StartTime).SpannerValue()

	if err != nil {
		return nil, grpc.Errorf(codes.Unknown, err.Error())
	}
	params = append(params, conv)
	params = append(params, "bananas")
	muts := make([]*spanner.Mutation, 1)
	muts[0] = spanner.Insert("example_table", []string{"id", "start_time", "name"}, params)
	_, err = s.SpannerDB.Apply(ctx, muts)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			return nil, grpc.Errorf(codes.AlreadyExists, err.Error())
		} else {
			return nil, grpc.Errorf(codes.Unknown, err.Error())
		}
	}
	res := test.ExampleTable{}

	err = hooks.UniaryInsertAfterHook(req, &res)
	if err != nil {

		return nil, grpc.Errorf(codes.Unknown, err.Error())

	}

	return &res, nil
}

// spanner unary select UniarySelectWithHooks
func (s *MySpannerImpl) UniarySelectWithHooks(ctx context.Context, req *test.ExampleTable) (*test.ExampleTable, error) {
	var err error
	var (
		Id        int64
		Name      string
		StartTime mytime.MyTime
	)

	beforeRes, err := hooks.UniaryInsertBeforeHook(req)
	if err != nil {

		return nil, grpc.Errorf(codes.Unknown, err.Error())

	}
	if beforeRes != nil {

		return beforeRes, nil

	}

	params := make(map[string]interface{})
	var conv interface{}

	conv = req.Id

	if err != nil {
		return nil, grpc.Errorf(codes.Unknown, err.Error())
	}
	params["string0"] = conv
	//stmt := spanner.Statement{SQL: "{ {.Spanner.Query} }", Params: params}
	stmt := spanner.Statement{SQL: "SELECT * from example_table Where id= @string0", Params: params}
	tx := s.SpannerDB.Single()
	defer tx.Close()
	iter := tx.Query(ctx, stmt)
	defer iter.Stop()
	row, err := iter.Next()
	if err == iterator.Done {
		return nil, grpc.Errorf(codes.NotFound, "no rows found")
	} else if err != nil {
		return nil, grpc.Errorf(codes.Unknown, err.Error())
	}
	// scan our values out of the row

	err = row.ColumnByName("id", &Id)
	if err != nil {
		return nil, grpc.Errorf(codes.Unknown, err.Error())
	}

	gcv := new(spanner.GenericColumnValue)
	err = row.ColumnByName("start_time", gcv)
	if err != nil {
		return nil, grpc.Errorf(codes.Unknown, err.Error())
	}
	err = StartTime.SpannerScan(gcv)
	if err != nil {
		return nil, grpc.Errorf(codes.Unknown, err.Error())
	}

	err = row.ColumnByName("name", &Name)
	if err != nil {
		return nil, grpc.Errorf(codes.Unknown, err.Error())
	}

	_, err = iter.Next()
	if err != iterator.Done {
		fmt.Println("Unary select that returns more than one row..")
	}
	res := test.ExampleTable{

		Id:        Id,
		Name:      Name,
		StartTime: StartTime.ToProto(),
	}

	err = hooks.UniaryInsertAfterHook(req, &res)
	if err != nil {

		return nil, grpc.Errorf(codes.Unknown, err.Error())

	}

	return &res, nil
}

// spanner unary select UniaryUpdateWithHooks
func (s *MySpannerImpl) UniaryUpdateWithHooks(ctx context.Context, req *test.ExampleTable) (*test.PartialTable, error) {
	var err error

	beforeRes, err := hooks.UniaryUpdateBeforeHook(req)
	if err != nil {

		return nil, grpc.Errorf(codes.Unknown, err.Error())

	}
	if beforeRes != nil {

		return beforeRes, nil

	}

	params := make(map[string]interface{})
	var conv interface{}

	conv, err = mytime.MyTime{}.ToSpanner(req.StartTime).SpannerValue()

	if err != nil {
		return nil, grpc.Errorf(codes.Unknown, err.Error())
	}
	params["start_time"] = conv
	conv = "oranges"
	params["name"] = conv

	conv = req.Id

	if err != nil {
		return nil, grpc.Errorf(codes.Unknown, err.Error())
	}
	params["id"] = conv
	muts := make([]*spanner.Mutation, 1)
	muts[0] = spanner.UpdateMap("example_table", params)
	_, err = s.SpannerDB.Apply(ctx, muts)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			return nil, grpc.Errorf(codes.AlreadyExists, err.Error())
		} else {
			return nil, grpc.Errorf(codes.Unknown, err.Error())
		}
	}
	res := test.PartialTable{}

	err = hooks.UniaryUpdateAfterHook(req, &res)
	if err != nil {

		return nil, grpc.Errorf(codes.Unknown, err.Error())

	}

	return &res, nil
}

// spanner unary select UniaryDeleteWithHooks
func (s *MySpannerImpl) UniaryDeleteWithHooks(ctx context.Context, req *test.ExampleTableRange) (*test.ExampleTable, error) {
	var err error

	beforeRes, err := hooks.UniaryDeleteBeforeHook(req)
	if err != nil {

		return nil, grpc.Errorf(codes.Unknown, err.Error())

	}
	if beforeRes != nil {

		return beforeRes, nil

	}

	start := make([]interface{}, 0)
	end := make([]interface{}, 0)
	var conv interface{}
	conv = req.StartId
	if err != nil {
		return nil, grpc.Errorf(codes.Unknown, err.Error())
	}
	start = append(start, conv)
	conv = req.EndId
	if err != nil {
		return nil, grpc.Errorf(codes.Unknown, err.Error())
	}
	end = append(end, conv)
	key := spanner.KeyRange{
		Start: start,
		End:   end,
		Kind:  spanner.ClosedOpen,
	}
	muts := make([]*spanner.Mutation, 1)
	muts[0] = spanner.DeleteKeyRange("example_table", key)
	_, err = s.SpannerDB.Apply(ctx, muts)
	if err != nil {
		if strings.Contains(err.Error(), "does not exist") {
			return nil, grpc.Errorf(codes.NotFound, err.Error())
		}
	}
	res := test.ExampleTable{}

	err = hooks.UniaryDeleteAfterHook(req, &res)
	if err != nil {

		return nil, grpc.Errorf(codes.Unknown, err.Error())

	}

	return &res, nil
}

// spanner server streaming ServerStreamWithHooks
func (s *MySpannerImpl) ServerStreamWithHooks(req *test.Name, stream MySpanner_ServerStreamWithHooksServer) error {
	var (
		Id        int64
		Name      string
		StartTime mytime.MyTime
	)

	beforeRes, err := hooks.ServerStreamBeforeHook(req)
	if err != nil {

		return grpc.Errorf(codes.Unknown, err.Error())

	}
	if beforeRes != nil {

		for _, res := range beforeRes {
			err = stream.Send(res)
			if err != nil {
				return err
			}
		}
		return nil

	}

	params := make(map[string]interface{})
	stmt := spanner.Statement{SQL: "SELECT * FROM example_table", Params: params}
	tx := s.SpannerDB.Single()
	defer tx.Close()
	iter := tx.Query(context.Background(), stmt)
	defer iter.Stop()
	for {
		row, err := iter.Next()
		if err == iterator.Done {
			break
		} else if err != nil {
			return grpc.Errorf(codes.Unknown, err.Error())
		}
		// scan our values out of the row

		err = row.ColumnByName("id", &Id)
		if err != nil {
			return grpc.Errorf(codes.Unknown, err.Error())
		}

		gcv := new(spanner.GenericColumnValue)
		err = row.ColumnByName("start_time", gcv)
		if err != nil {
			return grpc.Errorf(codes.Unknown, err.Error())
		}
		err = StartTime.SpannerScan(gcv)
		if err != nil {
			return grpc.Errorf(codes.Unknown, err.Error())
		}

		err = row.ColumnByName("name", &Name)
		if err != nil {
			return grpc.Errorf(codes.Unknown, err.Error())
		}

		res := test.ExampleTable{

			Id:        Id,
			Name:      Name,
			StartTime: StartTime.ToProto(),
		}

		err = hooks.ServerStreamAfterHook(req, &res)
		if err != nil {

			return grpc.Errorf(codes.Unknown, err.Error())

		}

		stream.Send(&res)
	}
	return nil
}

// spanner client streaming ClientStreamUpdateWithHooks
func (s *MySpannerImpl) ClientStreamUpdateWithHooks(stream MySpanner_ClientStreamUpdateWithHooksServer) error {
	var err error
	res := test.NumRows{}

	reqs := make([]*test.ExampleTable, 0)

	muts := make([]*spanner.Mutation, 0)
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			return grpc.Errorf(codes.Unknown, err.Error())
		}

		beforeRes, err := hooks.ClientStreamUpdateBeforeHook(req)
		if err != nil {

			return grpc.Errorf(codes.Unknown, err.Error())

		}
		if beforeRes != nil {

			continue

		}

		reqs = append(reqs, req)

		//spanner client streaming update
		params := make(map[string]interface{})
		var conv interface{}

		conv, err = mytime.MyTime{}.ToSpanner(req.StartTime).SpannerValue()

		if err != nil {
			return grpc.Errorf(codes.Unknown, err.Error())
		}
		params["name"] = conv

		conv = req.Name

		if err != nil {
			return grpc.Errorf(codes.Unknown, err.Error())
		}
		params["id"] = conv
		muts = append(muts, spanner.UpdateMap("example_table", params))

		////////////////////////////// NOTE //////////////////////////////////////
		// In the future, we might do apply if muts gets really big,  but for now,
		// we only do one apply on the database with all the records stored in muts
		//////////////////////////////////////////////////////////////////////////
	}
	_, err = s.SpannerDB.Apply(context.Background(), muts)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			return grpc.Errorf(codes.AlreadyExists, err.Error())
		} else {
			return grpc.Errorf(codes.Unknown, err.Error())
		}
	}

	for _, req := range reqs {

		err = hooks.ClientStreamUpdateAfterHook(req, &res)
		if err != nil {

			return grpc.Errorf(codes.Unknown, err.Error())

		}

	}

	stream.SendAndClose(&res)
	return nil
}
