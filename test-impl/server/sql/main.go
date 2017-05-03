package main

import (
	"fmt"
	_ "github.com/lib/pq"
	pb "github.com/tcncloud/protoc-gen-persist/examples/sql/basic"
	"google.golang.org/grpc"
	"net"
)

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		panic(err)
	}
	grpcServer := grpc.NewServer()
	connString := "dbname=postgres user=postgres host=postgres port=5432 sslmode=disable"
	s, err := pb.NewAmazingImpl("postgres", connString)
	if err != nil {
		panic(err)
	}
	_, err = s.SqlDB.Exec(`CREATE TABLE example_table(
		id          bigserial,
		start_time  varchar(255) NOT NULL,
		name        varchar(255) NOT NULL,
		primary key(id)
	)`)
	if err != nil {
		fmt.Printf("Server err:  %+v", err)
	}
	pb.RegisterAmazingServer(grpcServer, s)
	fmt.Printf("server listening on 50051\n")
	grpcServer.Serve(lis)
}
