package main

import (
	"context"
	"database/sql"
	"fmt"
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
	service := pb.UServPersistImpl(conn, pb.UServ_ImplOpts{
		HOOKS:    hooks,
		MAPPINGS: mapping,
		HANDLERS: &RestOfImpl{
			DB: conn,
		},
	})
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

func (h *HooksImpl) InsertUsersBeforeHook(ctx context.Context, req *pb.User) (*pb.Empty, error) {
	pb.IncId(req)
	return nil, nil
}
func (h *HooksImpl) InsertUsersAfterHook(context.Context, *pb.User, *pb.Empty) error {
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

func (m *MappingImpl) TimestampTimestamp() pb.UServTimestampTimestampMappingImpl {
	return &pb.TimeString{}
}
func (m *MappingImpl) SliceStringParam() pb.UServSliceStringParamMappingImpl {
	return &pb.SliceStringConverter{}
}

type RestOfImpl struct {
	DB *sql.DB
}

func (d *RestOfImpl) UpdateAllNames(req *pb.Empty, stream pb.UServ_UpdateAllNamesServer) error {
	ctx := stream.Context()
	queries := pb.UServPersistQueries(d.DB, pb.UServ_QueryOpts{
		MAPPINGS: &MappingImpl{},
	})
	renameToFoo := queries.UpdateNameToFooQuery(ctx)
	allUsers := queries.GetAllUsersQuery(ctx).Execute(req)
	selectUser := queries.SelectUserByIdQuery(ctx)

	return allUsers.Each(func(row *pb.UServ_GetAllUsersRow) error {
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
