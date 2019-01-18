package main

import (
	"fmt"
	"io"
	"net"

	"cloud.google.com/go/spanner"
	"github.com/tcncloud/protoc-gen-persist/examples/user_spanner/pb"
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

	service := pb.ImplUServ(conn, &RestOfImpl{DB: conn}, pb.Opts_UServ{
		HOOKS:    &HooksImpl{},
		MAPPINGS: &MappingImpl{},
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

func (m *MappingImpl) TimestampTimestamp() pb.MappingImpl_UServ_TimestampTimestamp {
	return &pb.TimeString{}
}
func (m *MappingImpl) SliceStringParam() pb.MappingImpl_UServ_SliceStringParam {
	return &pb.SliceStringConverter{}
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
func (h *HooksImpl) GetFriendsBeforeHook(context.Context, *pb.Friends) (*pb.User, error) {
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
	queries := pb.QueriesUServ(pb.Opts_UServ{MAPPINGS: &MappingImpl{}})

	var users []*pb.User

	_, err := d.DB.ReadWriteTransaction(ctx, func(ctx context.Context, tx *spanner.ReadWriteTransaction) error {
		users = make([]*pb.User, 0)
		renameToFoo := queries.UpdateNameToFoo(ctx, tx)
		selectUser := queries.SelectUserById(ctx, tx)
		allUsers := queries.GetAllUsers(ctx, tx).Execute(req)

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
			users = append(users, res)

			return nil
		})
	})

	if err != nil {
		return err
	}

	for _, user := range users {
		err := stream.Send(user)
		if err != nil {
			return err
		}
	}

	return nil
}

func (d *RestOfImpl) UpdateUserNames(stream pb.UServ_UpdateUserNamesServer) error {
	ctx := stream.Context()
	queries := pb.QueriesUServ(pb.Opts_UServ{
		MAPPINGS: &MappingImpl{},
	})

	users := make([]*pb.User, 0)
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		users = append(users, req)
	}

	responses := make([]*pb.User, 0)
	_, err := d.DB.ReadWriteTransaction(ctx, func(ctx context.Context, tx *spanner.ReadWriteTransaction) error {
		updateUserName := queries.UpdateUserName(ctx, tx)
		selectUser := queries.SelectUserById(ctx, tx)

		for _, user := range users {
			err := updateUserName.Execute(user).Zero()
			if err != nil {
				return err
			}
			resp, err := selectUser.Execute(user).One().User()
			if err != nil {
				return err
			}
			responses = append(responses, resp)
		}

		return nil
	})
	if err != nil {
		return err
	}

	for _, user := range responses {
		if err := stream.Send(user); err != nil {
			return err
		}
	}

	return nil
}
