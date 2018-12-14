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
	service := pb.UServPersistImpl(conn, pb.UServ_ImplOpts{
		HOOKS:    &HooksImpl{},
		MAPPINGS: &MappingImpl{},
		HANDLERS: &RestOfImpl{},
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
	Mappings *MappingImpl
	Hooks    *HooksImpl
}

func (d *RestOfImpl) UpdateAllNames(r *pb.Empty, stream pb.UServ_UpdateAllNamesServer) error {
	// db, err := sql.Open(
	// 	"postgres",
	// 	"user=postgres password=postgres dbname=postgres sslmode=disable",
	// )
	// if err != nil {
	// 	return err
	// }
	// params, err := pb.EmptyToUServPersistType(d.Mappings, r)
	// if err != nil {
	// 	return err
	// }
	// res := pl.UServGetAllUsersQuery(db, params)

	// err = pb.IterUServUserProto(d.Mappings, res, func(user *pb.User) error {
	// 	params, err := pb.UserToUServPersistType(d.Mappings, user)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	// unlike spanner, the query is actually run here.
	// 	res := pl.UServUpdateNameToFooQuery(db, params)
	// 	return res.Err()
	// })
	// if err != nil {
	// 	return err
	// }
	// res = pl.UServGetAllUsersQuery(db, params)
	// if res.Err() != nil {
	// 	return res.Err()
	// }
	// return pb.IterUServUserProto(d.Mappings, res, stream.Send)

	return nil
}
