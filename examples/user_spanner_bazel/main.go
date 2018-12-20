package main

import (
	"fmt"
	"net"

	"cloud.google.com/go/spanner"
	"github.com/tcncloud/protoc-gen-persist/examples/user_spanner_bazel/pb"
  // pl "github.com/tcncloud/protoc-gen-persist/examples/user_spanner_bazel/pb/persist_lib"
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
  // defer conn.Close()

  service := pb.UServPersistImpl(conn, pb.UServ_ImplOpts{
    HOOKS: &HooksImpl{},
    MAPPINGS: &MappingImpl{},
    HANDLERS: &RestOfImpl{DB: conn},
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
  DB     *spanner.Client
}
type MappingImpl struct{}

func (m *MappingImpl) TimestampTimestamp() pb.TimestampTimestampMappingImpl {
	return &pb.TimeString{}
}
func (m *MappingImpl) SliceStringParam() pb.SliceStringParamMappingImpl {
  return &pb.SliceStringConverter{}
}

type HooksImpl struct{}

func (h *HooksImpl) InsertUsersBeforeHook(ctx context.Context, req *pb.User) (*pb.Empty, error) {
  pb.IncId(req)
  return nil, nil
}
func (h *HooksImpl) InsertUsersAfterHook(context.Context, *pb.User, *pb.Empty) error {
  return nil
}
func (h *HooksImpl) GetAllUsersBeforeHook(context.Context, *pb.Empty) ([]*pb.User, error) {
  return nil, nil
}
func (h *HooksImpl) GetAllUsersAfterHook(context.Context, *pb.Empty, *pb.User) error {
  return nil
}
func (h *HooksImpl) GetFriendsBeforeHook(context.Context, *pb.Friends) ([]*pb.User, error) {
	return nil, nil
}
func (h *HooksImpl) GetFriendsAfterHook(context.Context, *pb.Friends, *pb.User) error {
	return nil
}

func (d *RestOfImpl) CreateTable(req *pb.Empty) (*pb.Empty, error) {
  out := new(pb.Empty)
  return out, nil
}

// using the persist lib queries to implement your own handlers.
func (d *RestOfImpl) UpdateAllNames(req *pb.Empty, stream pb.UServ_UpdateAllNamesServer) error {
  ctx := stream.Context()
  queries := pb.UServPersistQueries(d.DB, pb.UServ_QueryOpts{
    MAPPINGS: &MappingImpl{},
  })
  // renameToFoo := queries.UpdateNameToFooQuery(ctx)
  allUsers := queries.GetAllUsersQuery(ctx).Execute(req)
  // selectUser := queries.SelectUserByIdQuery(ctx)

  return allUsers.Each(func(row *pb.UServ_GetAllUsersRow) error {
    user, err := row.User()
    if err != nil {
      return err
    }
    fmt.Println("this rows name: ", user.Name)
    return nil

    // err = renameToFoo.Execute(user).Zero()
    // if err != nil {
    //   return err
    // }

    // res, err := selectUser.Execute(user).One().User()
    // if err != nil {
    //   return err
    // }
    // return stream.Send(res)
  })
}
