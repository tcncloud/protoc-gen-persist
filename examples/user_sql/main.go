package main

import (
	"database/sql"
	"fmt"
	"net"

	_ "github.com/lib/pq"
	"github.com/tcncloud/protoc-gen-persist/examples/user_sql/pb"
	pl "github.com/tcncloud/protoc-gen-persist/examples/user_sql/pb/persist_lib"
	"google.golang.org/grpc"
)

func main() {
	restOfHandlers := &RestOfImpl{
		Mappings: &MappingImpl{},
		Hooks:    &HooksImpl{},
	}
	service := pb.NewUServBuilder().
		WithDefaultQueryHandlers().
		WithNewSqlDb("postgres", "user=postgres password=postgres dbname=postgres sslmode=disable").
		WithRestOfGrpcHandlers(restOfHandlers).
		WithHooks(restOfHandlers.Hooks).
		WithTypeMapping(restOfHandlers.Mappings).
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

type HooksImpl struct{}

func (h *HooksImpl) InsertUsersBeforeHook(req *pb.User) (*pb.Empty, error) {
	pb.IncId(req)
	return nil, nil
}
func (h *HooksImpl) InsertUsersAfterHook(*pb.User, *pb.Empty) error {
	return nil
}
func (h *HooksImpl) GetAllUsersBeforeHook(*pb.Empty) ([]*pb.User, error) {
	return nil, nil
}
func (h *HooksImpl) GetAllUsersAfterHook(*pb.Empty, *pb.User) error {
	return nil
}

type MyTimestampImpl struct{}

type MappingImpl struct{}

func (m *MappingImpl) TimestampTimestamp() pb.TimestampTimestampMappingImpl {
	return &pb.TimeString{}
}
func (m *MappingImpl) SliceStringParam() pb.SliceStringParamMappingImpl {
	return &pb.SliceStringConverter{}
}

type RestOfImpl struct {
	Mappings *MappingImpl
	Hooks    *HooksImpl
}

func (d *RestOfImpl) UpdateAllNames(r *pb.Empty, stream pb.UServ_UpdateAllNamesServer) error {
	db, err := sql.Open(
		"postgres",
		"user=postgres password=postgres dbname=postgres sslmode=disable",
	)
	if err != nil {
		return err
	}
	params, err := pb.EmptyToUServPersistType(d.Mappings, r)
	if err != nil {
		return err
	}
	res := pl.UServGetAllUsersQuery(db, params)
	err = pb.IterUServUserProto(d.Mappings, res, func(user *pb.User) error {
		params, err := pb.UserToUServPersistType(d.Mappings, user)
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
	return pb.IterUServUserProto(d.Mappings, res, stream.Send)
}
