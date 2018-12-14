package main

import (
	"fmt"
	"net"

	"cloud.google.com/go/spanner"
	"github.com/coltonmorris/protoc-gen-persist/examples/user_spanner_bazel/pb"
	pl "github.com/coltonmorris/protoc-gen-persist/examples/user_spanner_bazel/pb/persist_lib"
  // pl "github.com/coltonmorris/protoc-gen-persist/examples/user_spanner_bazel/pb/persist_lib"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func main() {
	params := ReadSpannerParams()
  ctx := context.Background()
  conn, err := spanner.NewClient(ctx, params.URI())
  if err != nil {
    fmt.Printf("error connecting to db: %v\n", err)
    return
  }
  defer conn.Close()

  service := pb.UServPersistImpl(conn, pb.UServ_ImplOpts{
    HOOKS: &HooksImpl{},
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

type RestOfImpl struct {
	Params SpannerParams
}
type MappingImpl struct{}

func (m *MappingImpl) TimestampTimestamp() pb.TimestampTimestampMappingImpl {
	return &pb.TimeString{}
}
func (m *MappingImpl) SliceStringParam() pb.SliceStringParamMappingImpl {
  return &pb.SliceStringConverter{}
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
func (h *HooksImpl) GetFriendsBeforeHook(*pb.Friends) ([]*pb.User, error) {
	return nil, nil
}
func (h *HooksImpl) GetFriendsAfterHook(*pb.Friends, *pb.User) error {
	return nil
}

func (d *RestOfImpl) CreateTable(req *pb.Empty) (*pb.Empty, error) {
  out := new(pb.Empty)
  return out, nil
}

// using the persist lib queries to implement your own handlers.
func (d *RestOfImpl) UpdateAllNames(req *pb.Empty, stream pb.UServ_UpdateAllNamesServer) error {
	client, err := spanner.NewClient(stream.Context(), d.Params.URI())
	if err != nil {
		return err
	}
	// convert our request type to a persist's type to use it in the query.
	params, err := pb.EmptyToUServPersistType(&MappingImpl{}, req)
	if err != nil {
		return err
	}
	// create the query using the persist type we got above (params)
	iter := client.Single().Query(stream.Context(), pl.UServGetAllUsersQuery(params))
	muts := make([]*spanner.Mutation, 0)
	err = pb.IterUServUserProto(&MappingImpl{}, iter, func(user *pb.User) error {
		params, err := pb.UserToUServPersistType(&MappingImpl{}, user)
		if err != nil {
			return err
		}
		// get our mutation for this iteration's user and add it to the group
		muts = append(muts, pl.UServUpdateNameToFooQuery(params))
		return nil
	})
	if err != nil {
		return err
	}
	// apply our mutations
	if _, err := client.Apply(stream.Context(), muts); err != nil {
		return err
	}
	// get all our updated users, and stream them back to the client.
	params, _ = pb.EmptyToUServPersistType(&MappingImpl{}, req)
	iter = client.Single().Query(stream.Context(), pl.UServGetAllUsersQuery(params))

	return pb.IterUServUserProto(&MappingImpl{}, iter, stream.Send)
}
