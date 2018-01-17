package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/tcncloud/protoc-gen-persist/test_service/user_sql/pb"
	pl "github.com/tcncloud/protoc-gen-persist/test_service/user_sql/pb/persist_lib"
	"google.golang.org/grpc"
	"net"
)

func main() {
	service := pb.NewUServBuilder().
		WithDefaultQueryHandlers().
		WithNewSqlDb("postgres", "user=postgres password=postgres dbname=postgres sslmode=disable").
		WithRestOfGrpcHandlers(&RestOfImpl{}).
		MustBuild()
	server := grpc.NewServer()

	pb.RegisterUServServer(server, service)

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		panic(err)
	}
	if err := server.Serve(lis); err != nil {
		fmt.Printf("error serving: %v\n", err)
	}
}

type RestOfImpl struct{}

func (d *RestOfImpl) UpdateAllNames(r *pb.Empty, stream pb.UServ_UpdateAllNamesServer) error {
	db, err := sql.Open(
		"postgres",
		"user=postgres password=postgres dbname=postgres sslmode=disable",
	)
	if err != nil {
		return err
	}
	params, err := pb.EmptyToUServPersistType(r)
	if err != nil {
		return err
	}
	res := pl.UServGetAllUsersQuery(db, params)
	err = pb.IterUServUserProto(res, func(user *pb.User) error {
		params, err := pb.UserToUServPersistType(user)
		if err != nil {
			return err
		}
		// unlike spanner, the query is actually run here.
		res := pl.UServUpdateNameToFooQuery(db, params)
		return res.Err()
	})
	if err != nil {
		return err
	}
	res = pl.UServGetAllUsersQuery(db, params)
	if res.Err() != nil {
		return res.Err()
	}
	return pb.IterUServUserProto(res, func(user *pb.User) error {
		return stream.Send(user)
	})
}
