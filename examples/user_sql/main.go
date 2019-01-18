package main

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"net"

	_ "github.com/lib/pq"
	"github.com/tcncloud/protoc-gen-persist/examples/user_sql/pb"
	"google.golang.org/grpc"
)

func main() {
	conn, err := sql.Open("postgres", "user=postgres password=postgres dbname=postgres sslmode=disable")
	if err != nil {
		panic(err)
	}

	hooks := &HooksImpl{}
	mapping := &MappingImpl{}
	opts := pb.OptsUServ(hooks, mapping)
	handlers := &RestOfImpl{
		DB:      conn,
		QUERIES: pb.QueriesUServ(opts),
	}
	service := pb.ImplUServ(conn, handlers, opts)
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

type HooksImpl struct{}

func (h *HooksImpl) InsertUsersBeforeHook(ctx context.Context, req *pb.User) (*pb.Empty2, error) {
	pb.IncId(req)
	return nil, nil
}
func (h *HooksImpl) InsertUsersAfterHook(context.Context, *pb.User, *pb.Empty2) error {
	return nil
}
func (h *HooksImpl) GetAllUsersBeforeHook(context.Context, *pb.Empty) (*pb.User, error) {
	return nil, nil
}
func (h *HooksImpl) GetAllUsersAfterHook(context.Context, *pb.Empty, *pb.User) error {
	return nil
}

type MyTimestampImpl struct{}

type MappingImpl struct{}

func (m *MappingImpl) TimestampTimestamp() pb.MappingImpl_UServ_TimestampTimestamp {
	return &pb.TimeString{}
}
func (m *MappingImpl) SliceStringParam() pb.MappingImpl_UServ_SliceStringParam {
	return &pb.SliceStringConverter{}
}

// Type Aliasing to remove redundency
type Queries = pb.Queries_UServ
type RestOfImpl struct {
	DB      *sql.DB
	QUERIES *Queries
}

func (d *RestOfImpl) UpdateUserNames(stream pb.UServ_UpdateUserNamesServer) error {
	query := d.QUERIES.UpdateUserName(stream.Context(), d.DB)
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		res := new(pb.User)
		if err := query.Execute(req).One().Unwrap(res); err != nil {
			return err
		}
		if err := stream.Send(res); err != nil {
			return err
		}
	}
	return nil

}

func (d *RestOfImpl) UpdateAllNames(req *pb.Empty, stream pb.UServ_UpdateAllNamesServer) error {
	ctx := stream.Context()
	// tests that we can use both queries made from two different calls
	testOpts := pb.OptsUServ(nil, &MappingImpl{})
	renameToFoo := pb.QueriesUServ(testOpts).UpdateNameToFoo(ctx, d.DB)
	allUsers := d.QUERIES.GetAllUsers(ctx, d.DB).Execute(req)
	selectUser := d.QUERIES.SelectUserById(ctx, d.DB)

	return allUsers.Each(func(row *pb.Row_UServ_GetAllUsers) error {
		user, err := row.User()
		if err != nil {
			return err
		}

		err = renameToFoo.Execute(user).Zero()
		if err != nil {
			return err
		}

		res, err := selectUser.Execute(user).One().User()
		if err != nil {
			return err
		}
		return stream.Send(res)
	})
}
