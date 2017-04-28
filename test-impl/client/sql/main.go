package main

import(
	"fmt"
	"io"
	"time"
	"context"
	"google.golang.org/grpc"
	pbclient "github.com/tcncloud/protoc-gen-persist/examples/sql/basic"
	pb "github.com/tcncloud/protoc-gen-persist/examples/test"
	google_protobuf "github.com/golang/protobuf/ptypes/timestamp"
	ptypes "github.com/golang/protobuf/ptypes"
)
func setupClient() pbclient.AmazingClient {
	conn, err := grpc.Dial("s:50051",  grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	return pbclient.NewAmazingClient(conn)
}

func clientStreamInsert(client pbclient.AmazingClient, name string) error {
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
			Name: name,
		},
		&pb.ExampleTable{
			Id: int64(3),
			StartTime: ToProtobufTime(&now),
			Name: name,
		},
		&pb.ExampleTable{
			Id: int64(4),
			StartTime: ToProtobufTime(&now),
			Name: name,
		},
	}
	stream, err := client.ClientStreamWithHook(context.Background())
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

func serverStreamFromName(client pbclient.AmazingClient, name string) (*[]*pb.ExampleTable, error) {
	res := make([]*pb.ExampleTable, 0)
	fmt.Printf("Getting all docs that match name %s with server stream\n", name)
	stream, err := client.ServerStreamWithHooks(context.Background(), &pb.Name{ Name: name })
	if err != nil {
		return nil, err
	}
	for {
		doc, err := stream.Recv()
		if err == io.EOF {
			break;
		} else if err != nil {
			return nil, err
		}

		fmt.Printf("recieved doc: %+v\n", doc)
		res = append(res, doc)
	}
	fmt.Printf("recieved all serverStreamed docs\n\n")
	return &res, nil
}

func bidirectionalStream(client pbclient.AmazingClient, recs []*pb.ExampleTable) error {
	stream, err := client.BidirectionalWithHooks(context.Background())
	if err != nil {
		return err
	}

	tomorrow := time.Now().Add(time.Hour * 24)

	for _, rec := range(recs) {
		if rec != nil {
			fmt.Printf("Before Bidirectional Update: %+v\n", rec)
			rec.Name = "jenkins"
			rec.StartTime = ToProtobufTime(&tomorrow)
			err := stream.Send(rec)

			if err != nil {
				return err
			}

			changed, err := stream.Recv()
			if err == io.EOF {
				break;
			} else if err != nil {
				return err
			}

			fmt.Printf("After Bidirectional Update: %+v\n", changed)
		} else {
			fmt.Printf("nil record\n")
		}
	}
	return nil
}

func uniaryCacheTest(client pbclient.AmazingClient) error {
	yesterday := time.Now().Add(time.Hour * -24).Truncate(time.Millisecond)
	req := &pb.PartialTable{
		Id: int64(1),
		StartTime: ToProtobufTime(&yesterday),
	}
	res1, err := client.UniarySelectWithHooks(context.Background(), req)
	if err != nil {
		return err
	}
	res2, err := client.UniarySelectWithHooks(context.Background(), req)
	if err != nil {
		return err
	}
	fmt.Printf("UniaryCacheTest: req: %+v\nresp1:%+v\nresp2:%+v\n\n", *req, *res1, *res2)
	return nil
}

func main() {
	client := setupClient()

	err := clientStreamInsert(client, "bill")
	if err != nil {
		panic(err)
	}
	err = uniaryCacheTest(client)
	if err != nil {
		panic(err)
	}
	docs, err := serverStreamFromName(client, "bill")
	if err != nil {
		panic(err)
	}
	err = bidirectionalStream(client, *docs)
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
