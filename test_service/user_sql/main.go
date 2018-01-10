package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/tcncloud/protoc-gen-persist/test_service/user_sql/pb"
	pl "github.com/tcncloud/protoc-gen-persist/test_service/user_sql/pb/persist_lib"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"net"
)

func main() {
	Serve(func(s *grpc.Server) {
		lis, err := net.Listen("tcp", "0.0.0.0:50051")
		if err != nil {
			panic(err)
		}
		if err := s.Serve(lis); err != nil {
			fmt.Printf("error serving: %v\n", err)
		}
	})
}

func Serve(servFunc func(s *grpc.Server)) {
	service := pb.NewUServBuilder().
		WithNilAsDefaultQueryHandlers(&pl.UServQueryHandlers{DropTableHandler: myDrop}).
		WithNewSqlDb("postgres", "user=postgres password=postgres dbname=postgres sslmode=disable").
		WithRestOfGrpcHandlers(&ShutdownImpl{}).
		MustBuild()
	server := grpc.NewServer()

	pb.RegisterUServServer(server, service)

	servFunc(server)
}

type ShutdownImpl struct{}

// an example of using your own grpc handlers mixed with persists auto generated ones
func (d *ShutdownImpl) Shutdown(ctx context.Context, req *pb.Empty) (*pb.Empty, error) {
	return nil, fmt.Errorf("Unimplemented")
}

// an example of using a custom query handler instead of the default query handler.
// you can do whatever you want here, and have the pb/persist_lib package available
func myDrop(ctx context.Context, req *pl.EmptyForUServ, _ func(pl.Scanable)) error {
	db, err := sql.Open(
		"postgres",
		"user=postgres password=postgres dbname=postgres sslmode=disable",
	)
	if err != nil {
		return err
	}
	res, err := pl.EmptyFromDropTableQuery(db, req)
	if err != nil {
		return err
	}
	fmt.Printf("result? %+v\n", res)
	return nil
}
