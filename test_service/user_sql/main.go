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

func (d *ShutdownImpl) UpdateAllNames(r *pb.Empty, stream pb.UServ_UpdateAllNamesServer) error {
	db, err := sql.Open(
		"postgres",
		"user=postgres password=postgres dbname=postgres sslmode=disable",
	)
	if err != nil {
		return err
	}
	// the params part is turned into this function.
	params, err := pb.EmptyToUServPersistType(r)
	res := pl.UServGetAllUsersQuery(db, params)
	err = res.Do(func(s pl.Scanable) error {
		user, err := pb.UserFromUServRow(s)
		if err != nil {
			return err
		}
		params, err := pb.UserToUServPersistType(user)
		if err != nil {
			return err
		}
		res := pl.UServUpdateNameToFooQuery(db, params)
		if res.Err() != nil {
			return res.Err()
		}
		return nil
	})
	if err != nil {
		return err
	}
	res = pl.UServGetAllUsersQuery(db, params)
	if res.Err() != nil {
		return res.Err()
	}
	err = res.Do(func(s pl.Scanable) error {
		user, err := pb.UserFromUServRow(s)
		if err != nil {
			return err
		}
		if err := stream.Send(user); err != nil {
			return err
		}
		return nil
	})
	return err
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
	res := pl.UServDropTableQuery(db, req)
	if err != nil {
		return err
	}

	fmt.Printf("result? %+v\n", res)
	return nil
}
