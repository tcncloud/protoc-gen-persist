package main


import(
	"fmt"
	"io"
	"time"
	"context"
	"google.golang.org/grpc"
	pb "github.com/tcncloud/protoc-gen-persist/examples/spanner/bob_example"
	ptypes "github.com/golang/protobuf/ptypes"
)


func setupClient() pb.BobsClient{
	conn, err := grpc.Dial("127.0.0.1:50051",  grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	return pb.NewBobsClient(conn)
}

func main() {
	cli := setupClient()
	nowProto, _ := ptypes.TimestampProto(time.Now())
	oneHourProto, _ := ptypes.TimestampProto(time.Now().Add(time.Hour))
	twoHoursProto, _ := ptypes.TimestampProto(time.Now().Add(time.Hour * 2))

	stream, err := cli.PutBobs(context.Background())
	if err != nil {
		panic(err)
	}
	for i := 0; i < 3; i++ {
		err := stream.Send(&pb.Bob{Id: int64(i), Name: "Bob", StartTime: nowProto})
		if err != nil {
			panic(err)
		}
	}
	stream.Send(&pb.Bob{Id: int64(3), Name: "Alice", StartTime: nowProto})
	stream.Send(&pb.Bob{Id: int64(4), Name: "Bob", StartTime: twoHoursProto})
	stream.CloseAndRecv()

	PrintBobs(cli)

	cli.DeleteBobs(context.Background(), &pb.Bob{StartTime: oneHourProto})
	PrintBobs(cli)

}

func PrintBobs(client pb.BobsClient) error {
	fmt.Printf("Getting all docs with server stream\n\n")
	stream, err := client.GetBobs(context.Background(), &pb.Empty{ })
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
	return nil
}
