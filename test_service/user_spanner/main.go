package main

import (
	"cloud.google.com/go/spanner"
	"fmt"
	"github.com/tcncloud/protoc-gen-persist/test_service/user_spanner/pb"
	pl "github.com/tcncloud/protoc-gen-persist/test_service/user_spanner/pb/persist_lib"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"net"
)

func main() {
	params := ReadSpannerParams()
	service := pb.NewUServBuilder().
		WithDefaultQueryHandlers().
		WithSpannerURI(context.Background(), params.URI()).
		WithRestOfGrpcHandlers(&RestOfImpl{Params: params}).
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

// using the persist lib queries to implement your own handlers.
func (d *RestOfImpl) UpdateAllNames(req *pb.Empty, stream pb.UServ_UpdateAllNamesServer) error {
	client, err := spanner.NewClient(stream.Context(), d.Params.URI())
	if err != nil {
		return err
	}
	// convert our request type to a persist's type to use it in the query.
	params, err := pb.EmptyToUServPersistType(req)
	if err != nil {
		return err
	}
	// create the query using the persist type we got above (params)
	iter := client.Single().Query(stream.Context(), pl.UServGetAllUsersQuery(params))
	muts := make([]*spanner.Mutation, 0)
	err = pb.IterUServUserProto(iter, func(user *pb.User) error {
		params, err := pb.UserToUServPersistType(user)
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
	params, _ = pb.EmptyToUServPersistType(req)
	iter = client.Single().Query(stream.Context(), pl.UServGetAllUsersQuery(params))

	return pb.IterUServUserProto(iter, func(user *pb.User) error {
		return stream.Send(user)
	})
}
