package main


import(
	"fmt"
	"io"
	"time"
	"context"
	"google.golang.org/grpc"
	pb "github.com/tcncloud/protoc-gen-persist/examples/spanner/basic"
	google_protobuf "github.com/golang/protobuf/ptypes/timestamp"
	ptypes "github.com/golang/protobuf/ptypes"
)


func setupClient() pb.MySpannerClient {
	conn, err := grpc.Dial("127.0.0.1:50051",  grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	return pb.NewMySpannerClient(conn)
}
var ctx = context.Background()

func main() {
	client := setupClient()
	var err error
	err = uniaryInsert(client)
	if err != nil { panic(err) }
	err = clientStreamInsert(client)
	if err != nil { panic(err) }
	err = uniarySelect(client)
	if err != nil { panic(err) }
	err = noArgs(client)
	if err != nil { panic(err) }
	err = serverStream(client)
	if err != nil { panic(err) }
	err = uniaryUpdate(client)
	if err != nil { panic(err) }
	err = serverStream(client)
	if err != nil { panic(err) }
	err = clientStreamUpdate(client)
	if err != nil { panic(err) }
	err = uniaryDelete(client)
	if err != nil { panic(err) }
	err = clientStreamDelete(client)
	if err != nil { panic(err) }
	err = serverStream(client)
}
// query: "Insert into example_table (id, start_time, name)  Values (?, ?, \"bananas\")"
func uniaryInsert(client pb.MySpannerClient) error {
	now := time.Now().Truncate(time.Millisecond)
	fmt.Println("performing uniary insert")
	_, err := client.UniaryInsert(ctx, &pb.ExampleTable{Id: int64(5), StartTime: ToProtobufTime(&now)})
	if err != nil {
		return err
	}
	fmt.Println("inserted with uniary insert")
	return nil
}

// query: "SELECT * from example_table Where id=? AND name=?"
func uniarySelect(client pb.MySpannerClient) error {
	fmt.Println("performing uniary select")
	res, err := client.UniarySelect(ctx, &pb.ExampleTable{Id: int64(5), Name: "bananas"})
	if err != nil {
		return err
	}
	fmt.Printf("recieved res: %+v\n\n", res)
	return nil
}
// query: "Update example_table set start_time=?, name=\"oranges\" where id=?",
func uniaryUpdate(client pb.MySpannerClient) error {
	fmt.Println("performing uniaryUpdate")
	now := time.Now().Truncate(time.Millisecond)
	_, err := client.UniaryUpdate(ctx, &pb.ExampleTable{
		StartTime: ToProtobufTime(&now),
		Id: int64(1),
	})
	if err != nil {
		return err
	}
	fmt.Println("updated with uniary update")
	return nil
}

// query: "DELETE FROM example_table WHERE id>? AND id<?",
func uniaryDelete(client pb.MySpannerClient) error {
	fmt.Println("performing uniaryDelete")
	_, err := client.UniaryDelete(ctx, &pb.ExampleTableRange{StartId: int64(0), EndId: int64(5)})
	if err != nil {
		return err
	}
	fmt.Println("deleted using uniary delete")
	return nil
}

// query: "select * from example_table limit 1",
func noArgs(client pb.MySpannerClient) error {
	fmt.Println("performing NoArgs query")
	res, err := client.NoArgs(ctx, &pb.ExampleTable{})
	if err != nil {
		return err
	}
	fmt.Printf("recieved this from noArgs: %+v\n", res)
	return nil
}

// query: "SELECT * FROM example_table"
func serverStream(client pb.MySpannerClient) error {
	fmt.Printf("Getting all docs with server stream\n\n")
	stream, err := client.ServerStream(ctx, &pb.Name{ })
	if err != nil {
		return  err
	}
	for {
		doc, err := stream.Recv()
		if err == io.EOF {
			break;
		} else if err != nil {
			return err
		}

		fmt.Printf("%+v\n", doc)
	}
	fmt.Printf("recieved all serverStreamed docs\n\n")
	return nil
}

// query: "INSERT INTO example_table (id, start_time, name) VALUES (?, ?, ?)"
func clientStreamInsert(client pb.MySpannerClient) error {
	fmt.Println("inserting docs with client stream")
	now := time.Now().Truncate(time.Millisecond)
	docs := []*pb.ExampleTable{
		&pb.ExampleTable{
			Id: int64(1),
			StartTime: ToProtobufTime(&now),
			Name: "george",
		},
		&pb.ExampleTable{
			Id: int64(2),
			StartTime: ToProtobufTime(&now),
			Name: "michelle",
		},
		&pb.ExampleTable{
			Id: int64(3),
			StartTime: ToProtobufTime(&now),
			Name: "frank",
		},
		&pb.ExampleTable{
			Id: int64(4),
			StartTime: ToProtobufTime(&now),
			Name: "amy",
		},
	}
	stream, err := client.ClientStreamInsert(ctx)
	for _, doc := range(docs) {
		fmt.Printf("clientStreaming doc: %+v\n", doc)
		err := stream.Send(doc)
		if err != nil {
			return err
		}
	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		return err
	}
	fmt.Printf("inserted docs with clientStream: %+v\n\n", res)
	return nil
}

// query: "delete from example_table where id=?",
func clientStreamDelete(client pb.MySpannerClient) error {
	fmt.Println("deleting docs with client stream")
	docs := []*pb.ExampleTable{
		&pb.ExampleTable{
			Id: int64(5),
		},
	}
	stream, err := client.ClientStreamDelete(ctx)
	for _, doc := range(docs) {
		fmt.Printf("clientStreaming doc: %+v\n", doc)
		err := stream.Send(doc)
		if err != nil {
			return err
		}
	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		return err
	}
	fmt.Printf("deleted docs with clientStream: %+v\n\n", res)
	return nil
}

// query: "update example_table set start_time=?, name=? where id=?",
func clientStreamUpdate(client pb.MySpannerClient) error {
	fmt.Println("updating docs with client stream")
	now := time.Now().Truncate(time.Millisecond)
	docs := []*pb.ExampleTable{
		&pb.ExampleTable{
			Id: int64(1),
			StartTime: ToProtobufTime(&now),
			Name: "notgeorge",
		},
		&pb.ExampleTable{
			Id: int64(2),
			StartTime: ToProtobufTime(&now),
			Name: "notmichelle",
		},
		&pb.ExampleTable{
			Id: int64(3),
			StartTime: ToProtobufTime(&now),
			Name: "notfrank",
		},
		&pb.ExampleTable{
			Id: int64(4),
			StartTime: ToProtobufTime(&now),
			Name: "notamy",
		},
	}
	stream, err := client.ClientStreamUpdate(ctx)
	for _, doc := range(docs) {
		fmt.Printf("clientStreaming doc: %+v\n", doc)
		err := stream.Send(doc)
		if err != nil {
			return err
		}
	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		return err
	}
	fmt.Printf("updated docs with clientStream: %+v\n\n", res)
	return nil
}



// Just here to make things convienient
func ToTime(entry *google_protobuf.Timestamp) *time.Time {
	if entry == nil {
		return nil
	}
	lTime, err := ptypes.Timestamp(entry)
	if err != nil {
		return nil
	}
	return &lTime
}

func ToProtobufTime(lTime *time.Time) *google_protobuf.Timestamp {
	if lTime == nil {
		return nil
	}
	res, err := ptypes.TimestampProto(*lTime)
	if err != nil {
		fmt.Printf("something wrong %+v", err)
		return nil
	}
	return res
}
