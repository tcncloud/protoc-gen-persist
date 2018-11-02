package main

import (
	"fmt"
	"net"

	"cloud.google.com/go/spanner"
	"github.com/tcncloud/protoc-gen-persist/examples/user_spanner_bazel/pb"
	pl "github.com/tcncloud/protoc-gen-persist/examples/user_spanner_bazel/pb/persist_lib"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func main() {
	params := ReadSpannerParams()
	service := pb.NewUServBuilder().
		WithDefaultQueryHandlers().
		WithSpannerURI(context.Background(), params.URI()).
		WithRestOfGrpcHandlers(&RestOfImpl{Params: params}).
		WithHooks(&HooksImpl{}).
		WithTypeMapping(&Mappings{}).
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

type RestOfImpl struct {
	Params SpannerParams
}
type Mappings struct{}

func (m *Mappings) TimestampTimestamp() pb.UServTimestampTimestampMappingImpl {
	return &pb.TimeString{}
}

type HooksImpl struct{}

func (h *HooksImpl) GetFriendsBeforeHook(*pb.Friends) ([]*pb.User, error) {
	return nil, nil
}
func (h *HooksImpl) GetFriendsAfterHook(*pb.Friends, *pb.User) error {
	return nil
}

// using the persist lib queries to implement your own handlers.
func (d *RestOfImpl) UpdateAllNames(req *pb.Empty, stream pb.UServ_UpdateAllNamesServer) error {
	client, err := spanner.NewClient(stream.Context(), d.Params.URI())
	if err != nil {
		return err
	}
	// convert our request type to a persist's type to use it in the query.
	params, err := pb.EmptyToUServPersistType(&Mappings{}, req)
	if err != nil {
		return err
	}
	// create the query using the persist type we got above (params)
	iter := client.Single().Query(stream.Context(), pl.UServGetAllUsersQuery(params))
	muts := make([]*spanner.Mutation, 0)
	err = pb.IterUServUserProto(&Mappings{}, iter, func(user *pb.User) error {
		params, err := pb.UserToUServPersistType(&Mappings{}, user)
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
	params, _ = pb.EmptyToUServPersistType(&Mappings{}, req)
	iter = client.Single().Query(stream.Context(), pl.UServGetAllUsersQuery(params))

	return pb.IterUServUserProto(&Mappings{}, iter, stream.Send)
}
